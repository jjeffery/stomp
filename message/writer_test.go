package message

import (
	"bytes"
	. "launchpad.net/gocheck"
	"strings"
)

type WriterSuite struct{}

var _ = Suite(&WriterSuite{})

func (s *WriterSuite) TestWrites(c *C) {
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
	var writer = NewWriter(&b)
	err = writer.Write(frame)
	c.Assert(err, IsNil)
	newFrameText := b.String()
	c.Check(newFrameText, Equals, frameText)
	c.Check(b.String(), Equals, frameText)
}
