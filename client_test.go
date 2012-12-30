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

func (s *ClientSuite) TestConnectAndDisconnect(c *C) {
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
		c.Check(host, Equals, "the-server")
		writer.Write(message.NewFrame("CONNECTED"))

		f2, err := reader.Read()
		c.Assert(err, IsNil)
		c.Assert(f2.Command, Equals, "DISCONNECT")
		receipt, _ := f2.Contains("receipt")
		c.Check(receipt, Equals, "1")

		writer.Write(message.NewFrame("RECEIPT", message.ReceiptId, "1"))
	}()

	client := NewClient()
	err := client.Connect(fc1, nil)
	c.Assert(err, IsNil)

	err = client.Disconnect()
	c.Assert(err, IsNil)

	<-stop
}
