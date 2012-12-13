package stomp

import (
	"github.com/jjeffery/stomp/message"
	"strconv"
)

// Opcode used in requests.
type requestOp int

func (r requestOp) String() string {
	return strconv.Itoa(int(r))
}

// Valid value for request opcodes.
const (
	stopOp       requestOp = iota // server stop
	connectOp                     // client has connected
	disconnectOp                  // client has disconnected
	frameOp                       // process a frame
)

// requests received to be processed by main processing loop
type request struct {
	op    requestOp      // opcode for request
	conn  *conn          // connectOp, disconnectOp, frameOp
	frame *message.Frame // frameOp
}
