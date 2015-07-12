package stomp

import (
	"github.com/go-stomp/stomp/frame"
)

// SubscribeOpt contains options for for the Conn.Subscribe function.
var SubscribeOpt struct {
	// Header provides the opportunity to include custom header entries
	// in the SUBSCRIBE frame that the client sends to the server.
	Header func(header *frame.Header) func(*frame.Frame) error
}

func init() {
	SubscribeOpt.Header = func(header *frame.Header) func(*frame.Frame) error {
		return func(f *frame.Frame) error {
			if f.Command != frame.SUBSCRIBE {
				return ErrInvalidCommand
			}
			f.Header.AddHeader(header)
			return nil
		}
	}
}
