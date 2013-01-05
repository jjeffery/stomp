package stomp

import (
	"fmt"
	"github.com/jjeffery/stomp/frame"
	"github.com/jjeffery/stomp/testutil"
	. "launchpad.net/gocheck"
)

func (s *StompSuite) Test_unsuccessful_connect(c *C) {
	fc1, fc2 := testutil.NewFakeConn(c)
	stop := make(chan struct{})

	go func() {
		defer func() {
			fc2.Close()
			close(stop)
		}()

		reader := NewReader(fc2)
		writer := NewWriter(fc2)
		f1, err := reader.Read()
		c.Assert(err, IsNil)
		c.Assert(f1.Command, Equals, "CONNECT")
		f2 := NewFrame("ERROR", "message", "auth-failed")
		writer.Write(f2)
	}()

	conn, err := Connect(fc1, Options{})
	c.Assert(conn, IsNil)
	c.Assert(err, ErrorMatches, "auth-failed")
}

func (s *StompSuite) Test_successful_connect_and_disconnect(c *C) {
	testcases := []struct {
		Options           Options
		NegotiatedVersion string
		ExpectedVersion   Version
		ExpectedSession   string
		ExpectedHost      string
	}{
		{
			Options:         Options{},
			ExpectedVersion: "1.0",
			ExpectedSession: "",
			ExpectedHost:    "the-server",
		},
		{
			Options:           Options{},
			NegotiatedVersion: "1.1",
			ExpectedVersion:   "1.1",
			ExpectedSession:   "the-session",
			ExpectedHost:      "the-server",
		},
		{
			Options:           Options{Host: "xxx"},
			NegotiatedVersion: "1.2",
			ExpectedVersion:   "1.2",
			ExpectedSession:   "the-session",
			ExpectedHost:      "xxx",
		},
	}

	for _, tc := range testcases {
		resetId()
		fc1, fc2 := testutil.NewFakeConn(c)
		stop := make(chan struct{})

		go func() {
			defer func() {
				fc2.Close()
				close(stop)
			}()
			reader := NewReader(fc2)
			writer := NewWriter(fc2)

			f1, err := reader.Read()
			c.Assert(err, IsNil)
			c.Assert(f1.Command, Equals, "CONNECT")
			host, _ := f1.Contains("host")
			c.Check(host, Equals, tc.ExpectedHost)
			connectedFrame := NewFrame("CONNECTED")
			if tc.NegotiatedVersion != "" {
				connectedFrame.Add("version", tc.NegotiatedVersion)
			}
			if tc.ExpectedSession != "" {
				connectedFrame.Add("session", tc.ExpectedSession)
			}
			writer.Write(connectedFrame)

			f2, err := reader.Read()
			c.Assert(err, IsNil)
			c.Assert(f2.Command, Equals, "DISCONNECT")
			receipt, _ := f2.Contains("receipt")
			c.Check(receipt, Equals, "1")

			writer.Write(NewFrame("RECEIPT", frame.ReceiptId, "1"))
		}()

		client, err := Connect(fc1, tc.Options)
		c.Assert(err, IsNil)
		c.Assert(client, NotNil)
		c.Assert(client.Version(), Equals, tc.ExpectedVersion)
		c.Assert(client.Session(), Equals, tc.ExpectedSession)

		err = client.Disconnect()
		c.Assert(err, IsNil)

		<-stop
	}
}

func (s *StompSuite) Test_subscribe(c *C) {
	Helper_subscribe(c, AckAuto, V10)
	Helper_subscribe(c, AckAuto, V11)
	Helper_subscribe(c, AckAuto, V12)
	Helper_subscribe(c, AckClient, V10)
	Helper_subscribe(c, AckClient, V11)
	Helper_subscribe(c, AckClient, V12)
	Helper_subscribe(c, AckClientIndividual, V10)
	Helper_subscribe(c, AckClientIndividual, V11)
	Helper_subscribe(c, AckClientIndividual, V12)
}

func Helper_subscribe(c *C, ackMode AckMode, version Version) {
	fc1, fc2 := testutil.NewFakeConn(c)
	stop := make(chan struct{})

	go func() {
		defer func() {
			fc2.Close()
			close(stop)
		}()

		reader := NewReader(fc2)
		writer := NewWriter(fc2)
		f1, err := reader.Read()
		c.Assert(err, IsNil)
		c.Assert(f1.Command, Equals, "CONNECT")
		f2 := NewFrame("CONNECTED", "version", version.String())
		writer.Write(f2)
		f3, err := reader.Read()
		c.Assert(err, IsNil)
		c.Assert(f3.Command, Equals, "SUBSCRIBE")
		id, ok := f3.Contains("id")
		c.Assert(ok, Equals, true)
		destination := f3.Get("destination")
		c.Assert(destination, Equals, "/queue/test-1")
		ack := f3.Get("ack")
		c.Assert(ack, Equals, ackMode.String())

		for i := 1; i <= 5; i++ {
			messageId := fmt.Sprintf("message-%d", i)
			bodyText := fmt.Sprintf("Message body %d", i)
			f4 := NewFrame("MESSAGE",
				frame.Subscription, id,
				frame.MessageId, messageId,
				frame.Destination, destination)
			if version == V12 {
				f4.Add(frame.Ack, messageId)
			}
			f4.Body = []byte(bodyText)
			writer.Write(f4)

			if ackMode.ShouldAck() {
				f5, _ := reader.Read()
				c.Assert(f5.Command, Equals, "ACK")
				if version == V12 {
					c.Assert(f5.Get("id"), Equals, messageId)
				} else {
					c.Assert(f5.Get("subscription"), Equals, id)
					c.Assert(f5.Get("message-id"), Equals, messageId)
				}
			}
		}

		f6, _ := reader.Read()
		c.Assert(f6.Command, Equals, "UNSUBSCRIBE")
		c.Assert(f6.Get(frame.Receipt), Not(Equals), "")
		c.Assert(f6.Get(frame.Id), Equals, id)
		writer.Write(NewFrame(frame.RECEIPT,
			frame.ReceiptId, f6.Get(frame.Receipt)))

		f7, _ := reader.Read()
		c.Assert(f7.Command, Equals, "DISCONNECT")
		writer.Write(NewFrame(frame.RECEIPT,
			frame.ReceiptId, f7.Get(frame.Receipt)))
	}()

	conn, err := Connect(fc1, Options{})
	c.Assert(conn, NotNil)
	c.Assert(err, IsNil)
	sub, err := conn.Subscribe("/queue/test-1", ackMode)
	c.Assert(sub, NotNil)
	c.Assert(err, IsNil)

	for i := 1; i <= 5; i++ {
		msg := <-sub.C
		messageId := fmt.Sprintf("message-%d", i)
		bodyText := fmt.Sprintf("Message body %d", i)
		c.Assert(msg.Subscription, Equals, sub)
		c.Assert(msg.Body, DeepEquals, []byte(bodyText))
		c.Assert(msg.Destination, Equals, "/queue/test-1")
		c.Assert(msg.Header.Get(frame.MessageId), Equals, messageId)

		c.Assert(msg.ShouldAck(), Equals, ackMode.ShouldAck())
		if msg.ShouldAck() {
			msg.Conn.Ack(msg)
		}
	}

	err = sub.Unsubscribe()
	c.Assert(err, IsNil)

	conn.Disconnect()
}
