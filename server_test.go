package stomp

import (
	_ "github.com/gmallard/stompngo"
	. "launchpad.net/gocheck"
	"net"
)

type ServerSuite struct {}

var _ = Suite(&ServerSuite{})

func (s *ServerSuite) TestServerListens(c *C) {

	go ListenAndServe(":9091")
	
	conn, err := net.Dial("tcp", "127.0.0.1:9091")
	c.Assert(err, IsNil)
	
	conn.Close()
	
	
	
}


