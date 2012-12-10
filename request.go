package stomp

import (
	"github.com/jjeffery/stomp/message"
)

// client request operation
type requestOp int

// client requests operation code
const (
	stopOp       requestOp = iota // server stop
	disconnectOp                  // client connectionClosed
	frameOp                       // process a frame
)

// requests received
type request struct {
	op    requestOp
	conn  *conn
	frame *message.Frame
}
