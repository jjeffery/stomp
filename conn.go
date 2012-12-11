package stomp

import (
	"fmt"
	"github.com/jjeffery/stomp/message"
	"io"
	"log"
	"net"
	"time"
)

// Maximum number of pending frames allowed to a client.
// before a disconnect occurs. If the client cannot keep
// up with the server, we do not want the server to backlog
// pending frames indefinitely.
const maxPendingWrites = 16

// Maximum number of pending frames allowed before the read
// go routine starts blocking.
const maxPendingReads = 16

// Represents a connection with the STOMP client.
type conn struct {
	server         *Server                               // Server configuration
	rw             net.Conn                              // Network connection to client
	writer         *message.Writer                       // Writes STOMP frames directly to the network connection
	requestChannel chan request                          // For sending requests
	writeChannel   chan *message.Frame                   // For sending frames to the client
	readChannel    chan *message.Frame                   // For reading frames from the client
	stateFunc      func(c *conn, f *message.Frame) error // State processing function
	readTimeout    time.Duration                         // Heart beat read timeout
	writeTimeout   time.Duration                         // Heart beat write timeout
	version        message.StompVersion                  // Negotiated STOMP protocol version
	closed         bool                                  // Is the connection closed
	txStore        txStore                               // Stores transactions in progress
}

func newConn(server *Server, rw net.Conn, channel chan request) *conn {
	c := new(conn)
	c.server = server
	c.rw = rw
	c.requestChannel = channel
	c.writeChannel = make(chan *message.Frame, maxPendingWrites)
	c.readChannel = make(chan *message.Frame, maxPendingReads)
	go c.readLoop()
	go c.processLoop()
	return c
}

// Write a frame to the connection.
func (c *conn) Send(f *message.Frame) {
	// place the frame on the write channel, or
	// close the connection if the write channel is full,
	// as this means the client is not keeping up.
	select {
	case c.writeChannel <- f:
	default:
		// write channel is full
		c.Close()
	}
}

// TODO: should send other information, such as receipt-id
func (c *conn) SendError(err error) {
	f := new(message.Frame)
	f.Command = message.ERROR
	f.Headers.Append(message.Message, err.Error())
	c.Send(f) // will close after successful send
}

// Send an ERROR frame to the client and immediately close the connection.
// Include the receipt-id header if the frame contains a receipt header.
func (c *conn) sendErrorImmediatelyAndClose(err error, f *message.Frame) {
	errorFrame := message.NewFrame(message.ERROR,
		message.Message, err.Error())

	// Include a receipt-id header if the frame that prompted the error had
	// a receipt header (as suggested by the STOMP protocol spec).
	if f != nil {
		if receipt, ok := f.Contains(message.Receipt); ok {
			errorFrame.Append(message.ReceiptId, receipt)
		}
	}

	// send the frame to the client, ignore any error condition
	// because we are about to close the connection anyway
	_ = c.sendImmediately(errorFrame)

	// close connection with the client
	c.Close()
}

// Sends a STOMP frame to the client immediately, does not push onto the
// write channel to be processed in turn.
func (c *conn) sendImmediately(f *message.Frame) error {
	return c.writer.Write(f)
}

// Go routine for reading bytes from a client and assembling into
// STOMP frames. Also handles heart-beat read timeout. All read
// frames are pushed onto the read channel to be processed by the
// processLoop go-routine. This keeps all processing of frames for
// this connection on the one go-routine and avoids race conditions.
func (c *conn) readLoop() {
	reader := message.NewReader(c.rw)
	for !c.closed {
		if c.readTimeout == time.Duration(0) {
			// infinite timeout
			c.rw.SetReadDeadline(time.Time{})
		} else {
			c.rw.SetReadDeadline(time.Now().Add(c.readTimeout))
		}
		f, err := reader.Read()
		if err != nil {
			if c.closed {
				log.Println("connection closed by server:", c.rw.RemoteAddr())
			} else {
				if err == io.EOF {
					log.Println("connection closed by client:", c.rw.RemoteAddr())
				} else {
					log.Println("read failed:", err, c.rw.RemoteAddr())
				}
				c.Close()
			}
			return
		}

		if f == nil {
			// if the frame is nil, then it is a heartbeat
			continue
		}

		// Add the frame to the read channel. Note that this will block
		// if we are reading from the client quicker than the server
		// can process frames.
		c.readChannel <- f
	}
}

// Go routine that processes all read frames and all write frames.
// Having all processing in one go routine helps eliminate any race conditions.
func (c *conn) processLoop() {
	c.writer = message.NewWriter(c.rw)
	c.stateFunc = connecting
	for !c.closed {
		var timerChannel <-chan time.Time
		var timer *time.Timer

		if c.writeTimeout > 0 {
			timer = time.NewTimer(c.writeTimeout)
			timerChannel = timer.C
		}

		select {
		case f := <-c.writeChannel:
			// have a frame to the client
			// stop the heart-beat timer
			if timer != nil {
				timer.Stop()
				timer = nil
			}

			// write the frame to the client
			err := c.writer.Write(f)
			if err != nil {
				c.Close()
				return
			}

			// if the frame just sent to the client is an error
			// frame, we disconnect
			if f.Command == message.ERROR {
				// sent an ERROR frame, so disconnect
				c.Close()
				return
			}

		case f := <-c.readChannel:
			// have just received a frame from the client
			err := c.stateFunc(c, f)
			if err != nil {
				c.sendErrorImmediatelyAndClose(err, f)
				return
			}

		case _ = <-timerChannel:
			// write a heart-beat
			err := c.writer.Write(nil)
			if err != nil {
				c.Close()
				return
			}
		}
	}
}

func (c *conn) Close() {
	c.closed = true
	c.rw.Close()
	c.requestChannel <- request{op: disconnectOp, conn: c}
}

func (c *conn) handleConnect(f *message.Frame) error {
	var err error

	if _, ok := f.Contains(message.Receipt); ok {
		// CONNNECT and STOMP frames are not allowed to have
		// a receipt header.
		return receiptInConnect
	}

	// Authenticate if an authenticator has been provided
	if c.server.Authenticator != nil {
		// if either of these fields are absent, pass nil to the
		// authenticator function.
		login, _ := f.Contains(message.Login)
		passcode, _ := f.Contains(message.Passcode)
		if !c.server.Authenticator.Authenticate(login, passcode) {
			// sleep to slow down a rogue client a little bit
			time.Sleep(time.Second)
			return authenticationFailed
		}
	}

	// TODO: add functionality for authentication.
	// currently no authentication checks are made

	c.version, err = f.AcceptVersion()
	if err != nil {
		return err
	}

	cx, cy, err := f.HeartBeat()
	if err != nil {
		return err
	}

	// Minimum value as per server config. If the client
	// has requested shorter periods than this value, the
	// server will insist on the longer time period.
	min := asMilliseconds(c.server.HeartBeat, message.MaxHeartBeat)

	// apply a minimum heartbeat 
	if cx > 0 && cx < min {
		cx = min
	}
	if cy > 0 && cy < min {
		cy = min
	}

	c.readTimeout = time.Duration(cx) * time.Millisecond
	c.writeTimeout = time.Duration(cy) * time.Millisecond

	// Note that the heart-beat header is included even if the
	// client is V1.0 and did not send a header. This should not
	// break V1.0 clients.
	response := message.NewFrame(message.CONNECTED,
		message.Version, string(c.version),
		message.Server, "stompd/x.y.z", // TODO: get version
		message.HeartBeat, fmt.Sprintf("%d,%d", cy, cx))

	c.Send(response)
	c.stateFunc = connected

	return nil
}

func connecting(c *conn, f *message.Frame) error {
	switch f.Command {
	case message.CONNECT, message.STOMP:
		return c.handleConnect(f)
	}
	return notConnected
}

func (c *conn) sendReceiptImmediately(f *message.Frame) error {
	if receipt, ok := f.Contains(message.Receipt); ok {
		return c.sendImmediately(message.NewFrame(message.RECEIPT, message.ReceiptId, receipt))
	}
	return nil
}

func (c *conn) handleDisconnect(f *message.Frame) error {
	// As soon as we receive a DISCONNECT frame from a client, we do
	// not want to send any more frames to that client, with the exception
	// of a RECEIPT frame if the client has requested one.
	// Ignore the error condition if we cannot send a RECEIPT frame,
	// as the connection is about to close anyway.
	_ = c.sendReceiptImmediately(f)
	c.Close()
	return nil
}

// Handle a SEND frame received from the client.
func (c *conn) handleSend(f *message.Frame) error {
	// This frame will be converted into a MESSAGE frame
	// and sent for distribution.

	// send a receipt and remove the header (don't want it in the MESSAGE frame)
	err := c.sendReceiptImmediately(f)
	f.Remove(message.Receipt)

	f.Command = message.MESSAGE

	request := request{op: frameOp, conn: c, frame: f}

	if tx, ok := f.Contains(message.Transaction); ok {
		// remove the transaction header from the frame, don't want it in MESSAGE
		f.Remove(message.Transaction)
		err = c.txStore.Add(tx, request)
		if err != nil {
			return err
		}
	} else {
		// not in a transaction, send to be processed
		c.requestChannel <- request
	}

	return nil
}

func connected(c *conn, f *message.Frame) error {
	switch f.Command {
	case message.CONNECT, message.STOMP:
		return unexpectedCommand
	case message.DISCONNECT:
		return c.handleDisconnect(f)
	case message.SEND:
		return c.handleSend(f)
	case message.MESSAGE, message.RECEIPT, message.ERROR:
		// should only be sent by the server, should not come from the client
		return unexpectedCommand
	default:
		return unknownCommand
	}
	panic("not reached")
}
