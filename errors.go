package stomp

import (
	"github.com/jjeffery/stomp/message"
)

// StompError implements the Error interface, and provides
// additional information about a STOMP error.
type Error struct {
	Message string
	Frame   *message.Frame
}

func (e Error) Error() string {
	return e.Message
}

func NewError(f *message.Frame) Error {
	e := Error{Frame: f}

	if f.Command == message.ERROR {
		if message, ok := f.Contains(message.Message); ok {
			e.Message = message
		} else {
			e.Message = "ERROR frame, missing message header"
		}
	} else {
		e.Message = "Unexpected frame: " + f.Command
	}
	return e
}
