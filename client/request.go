package client

import (
	"github.com/jjeffery/stomp/message"
)

// Represents a request received from the client,
// consisting of a frame and the connection it
// was received from
type Request struct {
	Type       RequestType    // type of request
	Connection *Connection    // connection originating request
	Frame      *message.Frame // frame associated with request, might be nil
}

// Indicates the type of request received from the client.
type RequestType int

const (
	Create = RequestType(iota)
	Send
	Subscribe
	Unsubscribe
	Disconnect
)
