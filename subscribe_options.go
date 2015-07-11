package stomp

import (
	"github.com/jjeffery/stomp/frame"
)

// SubscribeOpt contains options for for the Conn.Subscribe function.
var SubscribeOpt struct {
	// Header provides the opportunity to include custom header entries
	// in the SUBSCRIBE frame that the client sends to the server.
	Header func(header *Header) func(*Frame) error
}

func init() {
	SubscribeOpt.Header = func(header *Header) func(*Frame) error {
		return func(f *Frame) error {
			if f.Command != frame.SUBSCRIBE {
				return ErrInvalidCommand
			}
			f.Header.AddHeader(header)
			return nil
		}
	}
}
