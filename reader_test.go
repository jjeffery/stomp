package stomp

import (
	. "launchpad.net/gocheck"
	"strings"
	"testing"
	"io"
)

func Test(t *testing.T) { TestingT(t) }

type ReaderSuite struct{}

var _ = Suite(&ReaderSuite{})

func (s *ReaderSuite) TestConnect(c *C) {
	reader := NewReader(strings.NewReader("CONNECT\nlogin:xxx\npasscode:yyy\n\n\x00"))

	frame, err := reader.Read()
	c.Assert(err, IsNil)
	c.Assert(frame, NotNil)
	c.Assert(len(frame.Body), Equals, 0)
	
	// ensure we are at the end of input
	frame, err = reader.Read()
	c.Assert(frame, IsNil)
	c.Assert(err, Equals, io.EOF)
}

func (s *ReaderSuite) TestSendWithoutContentLength(c *C) {
	reader := NewReader(strings.NewReader("SEND\ndestination:xxx\n\nPayload\x00"))

	frame, err := reader.Read()
	c.Assert(err, IsNil)
	c.Assert(frame, NotNil)
	c.Assert(frame.Command, Equals, "SEND")
	c.Assert(len(frame.Headers), Equals, 1)
	c.Assert(frame.Headers[0].Name, Equals, "destination")
	c.Assert(frame.Headers[0].Value(), Equals, "xxx")
	c.Assert(string(frame.Body), Equals, "Payload")
	
	// ensure we are at the end of input
	frame, err = reader.Read()
	c.Assert(frame, IsNil)
	c.Assert(err, Equals, io.EOF)
}

func (s *ReaderSuite) TestSendWithContentLength(c *C) {
	reader := NewReader(strings.NewReader("SEND\ndestination:xxx\ncontent-length:5\n\n\x00\x01\x02\x03\x04\x00"))

	frame, err := reader.Read()
	c.Assert(err, IsNil)
	c.Assert(frame, NotNil)
	c.Assert(frame.Command, Equals, "SEND")
	c.Assert(len(frame.Headers), Equals, 2)
	c.Assert(frame.Headers[0].Name, Equals, "destination")
	c.Assert(frame.Body, DeepEquals, []byte{0x00,0x01,0x02,0x03,0x04})
	
	// ensure we are at the end of input
	frame, err = reader.Read()
	c.Assert(frame, IsNil)
	c.Assert(err, Equals, io.EOF)
}
