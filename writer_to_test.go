package stomp

import (
	"bytes"
	. "launchpad.net/gocheck"
	"strings"
	"testing"
)

func TestWriterTo(t *testing.T) { TestingT(t) }

type WriterToSuite struct{}

var _ = Suite(&WriterToSuite{})

func (s *WriterToSuite) Test1(c *C) {
	var frameTexts = []string{
		"CONNECT\nlogin:xxx\npasscode:yyy\n\n\x00",
		"SEND\ndestination:/queue/request\ntx:1\ncontent-length:5\n\n\x00\x01\x02\x03\x04\x00",
		"SEND\n\nABCD\x00",
	}

	for _, frameText := range frameTexts {
		writeToBufferAndCheck(c, frameText)
	}
}

func writeToBufferAndCheck(c *C, frameText string) {
	reader := NewReader(strings.NewReader(frameText))

	frame, err := reader.Read()
	c.Assert(err, IsNil)
	c.Assert(frame, NotNil)

	var b bytes.Buffer
	length, err := frame.WriteTo(&b)
	c.Assert(err, IsNil)
	newFrameText := b.String()
	c.Check(newFrameText, Equals, frameText)
	c.Check(length, Equals, int64(len(newFrameText)))
	c.Check(b.String(), Equals, frameText)
}
