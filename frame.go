package stomp

import (
//"github.com/jjeffery/stomp/frame"
)

// A Frame represents a STOMP frame. A frame consists of a command
// followed by a collection of header elements, and then an optional
// body.
//
// Users of this package will not normally need to make use of the Frame
// type directly. It is a lower level type useful for implementing 
// STOMP protocol handlers.
type Frame struct {
	Command string
	*Header
	Body []byte
}

// NewFrame creates a new STOMP frame with the specified command and headers. 
// The headers should contain an even number of entries. Each even index is 
// the header name, and the odd indexes are the assocated header values.
func NewFrame(command string, headers ...string) *Frame {
	f := &Frame{Command: command, Header: &Header{}}
	for index := 0; index < len(headers); index += 2 {
		f.Add(headers[index], headers[index+1])
	}
	return f
}

// Clone creates a deep copy of the frame and its header. The cloned
// frame shares the body with the original frame.
func (f *Frame) Clone() *Frame {
	return &Frame{Command: f.Command, Header: f.Header.Clone(), Body: f.Body}
}

// ValidateFor checks to see if the STOMP frame is valid for the 
// given version of the STOMP protocol. It checks for mandatory
// header entries.
func (f *Frame) ValidateFor(version Version) error {
	/*
		if version == Version("") {
			// unknown version, only valid for CONNECT or STOMP frame
			switch f.Command {
			case frame.CONNECT:
				return f.verifyConnect(version, false)
			case frame.STOMP:
				return f.verif
			}
		}
		if version != V10 &&
			version != V11 &&
			version != V12 {
			return invalidVersion
		}
		switch f.Command {
		case frame.CONNECT:
			return f.verifyConnect(version, false)
		case frame.STOMP:
			return f.verifyConnect(version, true)

		}*/
	return nil
}

func (f *Frame) verifyConnect(version Version, isStomp bool) error {
	switch version {
	case V10:
		if isStomp {

		}
	case V11:
	case V12:
	}
	return nil
}

func (f *Frame) verifyMandatory(keys ...string) error {
	for _, key := range keys {
		if _, ok := f.Header.index(key); !ok {
			return &Error{
				Message: "missing header: " + key,
				Frame:   f,
			}
		}
	}
	return nil
}
