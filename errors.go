package stomp

import (
	"github.com/jjeffery/stomp/frame"
)

var (
	ErrInvalidCommand        = newErrorMessage("invalid command")
	ErrInvalidFrameFormat    = newErrorMessage("invalid frame format")
	ErrUnsupportedVersion    = newErrorMessage("unsupported version")
	ErrCompletedTransaction  = newErrorMessage("transaction is completed")
	ErrNackNotSupported      = newErrorMessage("NACK not supported in STOMP 1.0")
	ErrNotReceivedMessage    = newErrorMessage("cannot ack/nack a message, not from server")
	ErrCannotNackAutoSub     = newErrorMessage("cannot send NACK for a subscription with ack:auto")
	ErrCompletedSubscription = newErrorMessage("subscription is unsubscribed")
	ErrClosed                = newErrorMessage("connection closed unexpectedly")
	ErrNilOption             = newErrorMessage("nil option")
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

func missingHeader(name string) Error {
	return newErrorMessage("missing header: " + name)
}

func newErrorMessage(msg string) Error {
	return Error{Message: msg}
}

func newError(f *Frame) Error {
	e := Error{Frame: f}

	if f.Command == frame.ERROR {
		if message := f.Header.Get(frame.Message); message != "" {
			e.Message = message
		} else {
			e.Message = "ERROR frame, missing message header"
		}
	} else {
		e.Message = "Unexpected frame: " + f.Command
	}
	return e
}
