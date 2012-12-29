package stomp

import (
	"fmt"
	_ "github.com/gmallard/stompngo"
	. "launchpad.net/gocheck"
	"net"
)

type ServerSuite struct{}

var _ = Suite(&ServerSuite{})

func (s *ServerSuite) TestConnectAndDisconnect(c *C) {
	l, err := net.Listen("tcp", ":9091")
	c.Assert(err, IsNil)
	defer func() { l.Close() }()
	go Serve(l)

	conn, err := net.Dial("tcp", "127.0.0.1:9091")
	c.Assert(err, IsNil)

	client := &Client{}
	err = client.Connect(conn, map[string]string{})
	c.Assert(err, IsNil)

	err = client.Disconnect()
	c.Assert(err, IsNil)

	conn.Close()

}

func (s *ServerSuite) TestSendMessages(c *C) {
	ch := make(chan bool, 2)

	l, err := net.Listen("tcp", ":9091")
	c.Assert(err, IsNil)
	defer func() { l.Close() }()
	go Serve(l)

	count := 10
	go runSender(c, ch, count)
	go runReceiver(c, ch, count)

	<-ch
	<-ch
}

func runSender(c *C, ch chan bool, count int) {
	conn, err := net.Dial("tcp", "127.0.0.1:9091")
	c.Assert(err, IsNil)

	client := &Client{}
	err = client.Connect(conn, nil)
	c.Assert(err, IsNil)

	for i := 0; i < count; i++ {
		client.Send(SendMessage{
			Destination: "/queue/test",
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("test message %d", i)),
		})
	}

	ch <- true
}

func runReceiver(c *C, ch chan bool, count int) {
	conn, err := net.Dial("tcp", "127.0.0.1:9091")
	c.Assert(err, IsNil)

	client := &Client{}
	err = client.Connect(conn, nil)
	c.Assert(err, IsNil)

	sub, err := client.Subscribe("/queue/test", AckAuto)
	c.Assert(err, IsNil)
	c.Assert(sub, NotNil)

	for i := 0; i < count; i++ {
		msg := <-sub.C
		expectedText := fmt.Sprintf("test message %d", i)
		c.Assert(msg.Body, DeepEquals, []byte(expectedText))
	}
	ch <- true
}
