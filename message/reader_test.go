package message

import (
	"io"
	. "launchpad.net/gocheck"
	"strings"
)

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
	c.Assert(frame.Headers.Count(), Equals, 1)
	k, v := frame.Headers.GetAt(0)
	c.Assert(k, Equals, "destination")
	c.Assert(v, Equals, "xxx")
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
	c.Assert(frame.Headers.Count(), Equals, 2)
	k, _ := frame.Headers.GetAt(0)
	c.Assert(k, Equals, "destination")
	c.Assert(frame.Body, DeepEquals, []byte{0x00, 0x01, 0x02, 0x03, 0x04})

	// ensure we are at the end of input
	frame, err = reader.Read()
	c.Assert(frame, IsNil)
	c.Assert(err, Equals, io.EOF)
}

func (s *ReaderSuite) TestInvalidCommand(c *C) {
	reader := NewReader(strings.NewReader("sEND\ndestination:xxx\ncontent-length:5\n\n\x00\x01\x02\x03\x04\x00"))

	frame, err := reader.Read()
	c.Check(frame, IsNil)
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "invalid command")
}

func (s *ReaderSuite) TestSendWithoutDestination(c *C) {
	reader := NewReader(strings.NewReader("SEND\ndeestination:xxx\ncontent-length:5\n\n\x00\x01\x02\x03\x04\x00"))

	frame, err := reader.Read()
	c.Check(frame, IsNil)
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "missing header: destination")
}

func (s *ReaderSuite) TestSubscribeWithoutDestination(c *C) {
	reader := NewReader(strings.NewReader("SUBSCRIBE\ndeestination:xxx\nid:7\n\n\x00"))

	frame, err := reader.Read()
	c.Check(frame, IsNil)
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "missing header: destination")
}

func (s *ReaderSuite) TestSubscribeWithoutId(c *C) {
	reader := NewReader(strings.NewReader("SUBSCRIBE\ndestination:xxx\nIId:7\n\n\x00"))

	frame, err := reader.Read()
	c.Check(frame, IsNil)
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "missing header: id")
}

func (s *ReaderSuite) TestUnsubscribeWithoutId(c *C) {
	reader := NewReader(strings.NewReader("UNSUBSCRIBE\nIId:7\n\n\x00"))

	frame, err := reader.Read()
	c.Check(frame, IsNil)
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "missing header: id")
}
