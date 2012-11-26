package stomp

import (
	"bufio"
	"io"
)

const (
	bufferSize = 4096
)

type Reader struct {
	reader *bufio.Reader
}

func NewReader(reader io.Reader) *Reader {
	sr := new(StompReader)
	sr.reader = bufio.NewReaderSize(reader, bufferSize)
	return sr
}

// Read a STOMP frame from the input. If the input contains one
// or more heart-beat characters and no frame, then nil will
// be returned for the frame. Calling programs should always check
// for a nil frame.
func (r *Reader) Read() (f *Frame, err error) {
	panic("not implemented: Reader.Read")
}
