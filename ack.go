package stomp

import (
	"github.com/jjeffery/stomp/frame"
)

// The AckMode type is an enumeration of the acknowledgement modes for a 
// STOMP subscription.
type AckMode int

// String returns the string representation of the AckMode value.
func (a AckMode) String() string {
	switch a {
	case AckAuto:
		return frame.AckAuto
	case AckClient:
		return frame.AckClient
	case AckClientIndividual:
		return frame.AckClientIndividual
	}
	panic("invalid AckMode value")
}

// These constants are the valid values for the AckMode type. When a STOMP
// client subscribes to a destination on the server, it specifies how it
// will acknowledge messages it receives from the server.
const (
	// No acknowledgement is required, the server assumes that the client 
	// received the message.
	AckAuto AckMode = iota

	// Client acknowledges messages. When a client acknowledges a message,
	// any previously received messages are also acknowledged.
	AckClient

	// Client acknowledges message. Each message is acknowledged individually.
	AckClientIndividual
)
