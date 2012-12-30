package stomp

import (
	"github.com/jjeffery/stomp/message"
	. "launchpad.net/gocheck"
)

type ClientSuite struct{}

var _ = Suite(&ClientSuite{})

func (s *ClientSuite) SetUpTest(c *C) {
	resetId()
}

func (s *ClientSuite) Test_successful_connect_and_disconnect(c *C) {
	testcases := []struct {
		Options           ConnectOptions
		NegotiatedVersion string
		ExpectedVersion   string
		ExpectedSession   string
		ExpectedHost      string
	}{
		{
			Options:         ConnectOptions{},
			ExpectedVersion: "1.0",
			ExpectedSession: "",
			ExpectedHost:    "the-server",
		},
		{
			Options:           ConnectOptions{},
			NegotiatedVersion: "1.1",
			ExpectedVersion:   "1.1",
			ExpectedSession:   "the-session",
			ExpectedHost:      "the-server",
		},
		{
			Options:           ConnectOptions{Host: "xxx"},
			NegotiatedVersion: "1.2",
			ExpectedVersion:   "1.2",
			ExpectedSession:   "the-session",
			ExpectedHost:      "xxx",
		},
	}

	for _, tc := range testcases {
		resetId()
		fc1, fc2 := newFakeConn(c)
		stop := make(chan struct{})

		go func() {
			defer func() {
				fc2.Close()
				close(stop)
			}()
			reader := message.NewReader(fc2)
			writer := message.NewWriter(fc2)

			f1, err := reader.Read()
			c.Assert(err, IsNil)
			c.Assert(f1.Command, Equals, "CONNECT")
			host, _ := f1.Contains("host")
			c.Check(host, Equals, tc.ExpectedHost)
			connectedFrame := message.NewFrame("CONNECTED")
			if tc.NegotiatedVersion != "" {
				connectedFrame.Append("version", tc.NegotiatedVersion)
			}
			if tc.ExpectedSession != "" {
				connectedFrame.Append("session", tc.ExpectedSession)
			}
			writer.Write(connectedFrame)

			f2, err := reader.Read()
			c.Assert(err, IsNil)
			c.Assert(f2.Command, Equals, "DISCONNECT")
			receipt, _ := f2.Contains("receipt")
			c.Check(receipt, Equals, "1")

			writer.Write(message.NewFrame("RECEIPT", message.ReceiptId, "1"))
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
