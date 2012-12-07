package message

import (
	"strconv"
)

// slices used to write frames
var (
	colonSlice   = []byte{58}     // colon ':'
	crlfSlice    = []byte{13, 10} // CR-LF
	newlineSlice = []byte{10}     // newline (LF)
	nullSlice    = []byte{0}      // null character
)

// Represents a single STOMP frame.
type Frame struct {
	// The frame command. Should be one of the standard STOMP commands. Note that
	// STOMP commands are case sensitive.
	Command string

	// Collection of frame headers.
	Headers

	// The frame body. Only SEND, MESSAGE and ERROR frames may have a body.
	// All other frames must not have a body.
	Body []byte
}

// Returns the value of the "content-length" header, and whether it was
// found or not. Used for deserializing a frame. If the content length
// is specified in the header, then the body can contain null characters.
// Otherwise the body is read until a null character is encountered.
// If an error is returned, then the content-length header is malformed.
func (f *Frame) ContentLength() (contentLength int, ok bool, err error) {
	text, ok := f.Headers.Contains(ContentLength)
	if !ok {
		return
	}

	value, err := strconv.ParseInt(text, 10, 32)
	if err != nil {
		ok = false
		return
	}

	contentLength = int(value)
	ok = true
	return
}
