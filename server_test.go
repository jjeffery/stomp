package stomp

import (
	"fmt"
	. "launchpad.net/gocheck"
	_ "log"
	"net"
	"runtime"
)

type ServerSuite struct{}

var _ = Suite(&ServerSuite{})

func (s *ServerSuite) SetUpSuite(c *C) {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func (s *ServerSuite) TearDownSuite(c *C) {
	runtime.GOMAXPROCS(1)
}

func (s *ServerSuite) TestConnectAndDisconnect(c *C) {
	addr := ":59091"
	l, err := net.Listen("tcp", addr)
	c.Assert(err, IsNil)
	defer func() { l.Close() }()
	go Serve(l)

	conn, err := net.Dial("tcp", "127.0.0.1"+addr)
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

	addr := ":59092"

	l, err := net.Listen("tcp", addr)
	c.Assert(err, IsNil)
	defer func() { l.Close() }()
	go Serve(l)

	count := 100
	go runReceiver(c, ch, count, "/queue/test-1", addr)
	go runSender(c, ch, count, "/queue/test-1", addr)
	go runSender(c, ch, count, "/queue/test-2", addr)
	go runReceiver(c, ch, count, "/queue/test-2", addr)
	go runReceiver(c, ch, count, "/queue/test-3", addr)
	go runSender(c, ch, count, "/queue/test-3", addr)
	go runSender(c, ch, count, "/queue/test-4", addr)
	go runReceiver(c, ch, count, "/queue/test-4", addr)

	for i := 0; i < 8; i++ {
		<-ch
	}
}

func runSender(c *C, ch chan bool, count int, destination, addr string) {
	conn, err := net.Dial("tcp", "127.0.0.1"+addr)
	c.Assert(err, IsNil)

	client := &Client{}
	err = client.Connect(conn, nil)
	c.Assert(err, IsNil)

	for i := 0; i < count; i++ {
		client.Send(SendMessage{
			Destination: destination,
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("%s test message %d", destination, i)),
		})
		//log.Println("sent", i)
	}

	ch <- true
}

func runReceiver(c *C, ch chan bool, count int, destination, addr string) {
	conn, err := net.Dial("tcp", "127.0.0.1"+addr)
	c.Assert(err, IsNil)

	client := &Client{}
	err = client.Connect(conn, nil)
	c.Assert(err, IsNil)

	sub, err := client.Subscribe(destination, AckAuto)
	c.Assert(err, IsNil)
	c.Assert(sub, NotNil)

	for i := 0; i < count; i++ {
		msg := <-sub.C
		expectedText := fmt.Sprintf("%s test message %d", destination, i)
		c.Assert(msg.Body, DeepEquals, []byte(expectedText))
		//log.Println("received", i)
	}
	ch <- true
}
