package message

import (
	"bufio"
	"bytes"
	"io"
)

// Reads STOMP frames from an underlying io.Reader.
type Reader struct {
	reader *bufio.Reader
}

// Creates a new Reader object which reads from the
// underlying io.Reader.
func NewReader(reader io.Reader) *Reader {
	sr := new(Reader)
	sr.reader = bufio.NewReaderSize(reader, bufferSize)
	return sr
}

// Read a STOMP frame from the input. If the input contains one
// or more heart-beat characters and no frame, then nil will
// be returned for the frame. Calling programs should always check
// for a nil frame.
func (r *Reader) Read() (*Frame, error) {
	commandSlice, err := r.readLine()
	if err != nil {
		return nil, err
	}

	if len(commandSlice) == 0 {
		// received a heart-beat newline char (or cr-lf)
		return nil, nil
	}

	frame := NewFrame(string(commandSlice))
	switch frame.Command {
	case CONNECT, STOMP, SEND, SUBSCRIBE, UNSUBSCRIBE, ACK, NACK, BEGIN, COMMIT, ABORT, DISCONNECT:
		// valid command
	default:
		return nil, invalidCommand
	}

	// read headers
	for {
		headerSlice, err := r.readLine()
		if err != nil {
			return nil, err
		}

		if len(headerSlice) == 0 {
			// empty line means end of headers
			break
		}

		index := bytes.IndexByte(headerSlice, colon)
		if index <= 0 {
			// colon is missing or header name is zero length
			return nil, invalidFrameFormat
		}

		name := string(headerSlice[0:index])
		value := string(headerSlice[index+1:])

		// TODO: need to decode if STOMP 1.1 or later

		frame.Headers.Append(name, value)
	}

	err = frame.Validate()
	if err != nil {
		return nil, err
	}

	// get content length from the headers
	if contentLength, ok, err := frame.ContentLength(); err != nil {
		// happens if the content is malformed
		return nil, err
	} else if ok {
		// content length specified in the header, so use that
		frame.Body = make([]byte, contentLength)
		for bytesRead := 0; bytesRead < contentLength; {
			n, err := r.reader.Read(frame.Body[bytesRead:contentLength])
			if err != nil {
				return nil, err
			}
			bytesRead += n
		}

		// TODO! need to read the null byte here!
	} else {
		frame.Body, err = r.reader.ReadBytes(nullByte)
		if err != nil {
			return nil, err
		}
		// remove trailing null
		frame.Body = frame.Body[0 : len(frame.Body)-1]
	}

	// pass back frame
	return frame, nil
}

// read one line from input and strip off terminating LF or terminating CR-LF
func (r *Reader) readLine() (line []byte, err error) {
	line, err = r.reader.ReadBytes(newline)
	if err != nil {
		return
	}

	switch {
	case bytes.HasSuffix(line, crlfSlice):
		line = line[0 : len(line)-len(crlfSlice)]
	case bytes.HasSuffix(line, newlineSlice):
		line = line[0 : len(line)-len(newlineSlice)]
	}

	return
}
