package stomp

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

	// Frame headers. Note that this is an array and not a map. The reason is
	// that STOMP 1.2 allows multiple headers with the same name. When there are
	// multiple headers with the same name, the first one has the value and any 
	// subsequent headers are for historical information only.
	Headers

	// The frame body. Only SEND, MESSAGE and ERROR frames may have a body.
	// All other frames must not have a body.
	Body []byte
}

func (f *Frame) ContentLength() (contentLength int, ok bool) {
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

/*
func (f *Frame) WriteTo(w io.Writer) (n int64, err error) {
	count, err := w.Write([]byte(f.Command))
	n += int64(count)
	if err != nil {
		return
	}

	count, err = w.Write(newlineSlice)
	n += int64(count)
	if err != nil {
		return
	}

	for _, h := range f.Headers {
		var count64 int64
		count64, err = h.WriteTo(w)
		n += count64
		if err != nil {
			return
		}
	}

	count, err = w.Write(newlineSlice)
	n += int64(count)
	if err != nil {
		return
	}

	if len(f.Body) > 0 {
		count, err = w.Write(f.Body)
		n += int64(count)
		if err != nil {
			return
		}
	}

	// write the final nul (0) byte	
	count, err = w.Write(nullSlice)
	n += int64(count)
	return
}

func (f *Frame) findHeader(name string) (index int, ok bool) {
	for i, v := range f.Headers {
		if v.Name == name {
			index = i
			ok = true
			return
		}
	}

	return
}
*/
