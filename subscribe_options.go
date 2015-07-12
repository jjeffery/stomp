package stomp

import (
	"github.com/go-stomp/stomp/frame"
)

// SubscribeOpt contains options for for the Conn.Subscribe function.
var SubscribeOpt struct {
	// Header provides the opportunity to include custom header entries
	// in the SUBSCRIBE frame that the client sends to the server.
	Header func(key, value string) func(*frame.Frame) error
}

func init() {
	SubscribeOpt.Header = func(key, value string) func(*frame.Frame) error {
		return func(f *frame.Frame) error {
			if f.Command != frame.SUBSCRIBE {
				return ErrInvalidCommand
			}
			if f.Header == nil {
				f.Header = frame.NewHeader(key, value)
			} else {
				f.Header.Add(key, value)
			}
			return nil
		}
	}
}
