package stomp

import (
	"github.com/jjeffery/stomp/frame"
)

var (
	invalidCommand     = newErrorMessage("invalid command")
	invalidFrameFormat = newErrorMessage("invalid frame format")
)

// StompError implements the Error interface, and provides
// additional information about a STOMP error.
type Error struct {
	Message string
	Frame   *Frame
}

func (e Error) Error() string {
	return e.Message
}

func newErrorMessage(msg string) Error {
	return Error{Message: msg}
}

func newError(f *Frame) Error {
	e := Error{Frame: f}

	if f.Command == frame.ERROR {
		if message := f.Get(frame.Message); message != "" {
			e.Message = message
		} else {
			e.Message = "ERROR frame, missing message header"
		}
	} else {
		e.Message = "Unexpected frame: " + f.Command
	}
	return e
}
