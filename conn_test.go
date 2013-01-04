package stomp

import (
	"github.com/jjeffery/stomp/frame"
	"github.com/jjeffery/stomp/testutil"
	. "launchpad.net/gocheck"
)

func (s *StompSuite) Test_successful_connect_and_disconnect(c *C) {
	testcases := []struct {
		Options           Options
		NegotiatedVersion string
		ExpectedVersion   string
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
