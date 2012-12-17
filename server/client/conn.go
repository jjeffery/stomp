package client

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
type Conn struct {
	config         Config
	rw             net.Conn                              // Network connection to client
	writer         *message.Writer                       // Writes STOMP frames directly to the network connection
	requestChannel chan Request                          // For sending requests
	writeChannel   chan *message.Frame                   // For sending frames to the client
	subChannel     chan *Subscription                    // For sending messages to client
	readChannel    chan *message.Frame                   // For reading frames from the client
	stateFunc      func(c *Conn, f *message.Frame) error // State processing function
	readTimeout    time.Duration                         // Heart beat read timeout
	writeTimeout   time.Duration                         // Heart beat write timeout
	version        message.StompVersion                  // Negotiated STOMP protocol version
	closed         bool                                  // Is the connection closed
	txStore        *txStore                              // Stores transactions in progress
}

func NewConn(config Config, rw net.Conn, channel chan Request) *Conn {
	c := &Conn{
		config:         config,
		rw:             rw,
		requestChannel: channel,
		subChannel:     make(chan *Subscription, maxPendingWrites),
		writeChannel:   make(chan *message.Frame, maxPendingWrites),
		readChannel:    make(chan *message.Frame, maxPendingReads),
		txStore:        &txStore{},
	}
	go c.readLoop()
	go c.processLoop()
	return c
}

// Write a frame to the connection.
func (c *Conn) Send(f *message.Frame) {
	// place the frame on the write channel, or
	// close the connection if the write channel is full,
	// as this means the client is not keeping up.
	select {
	case c.writeChannel <- f:
	default:
		// write channel is full
		c.close()
	}
}

// TODO: should send other information, such as receipt-id
func (c *Conn) SendError(err error) {
	f := new(message.Frame)
	f.Command = message.ERROR
	f.Headers.Append(message.Message, err.Error())
	c.Send(f) // will close after successful send
}

// Send an ERROR frame to the client and immediately close the connection.
// Include the receipt-id header if the frame contains a receipt header.
func (c *Conn) sendErrorImmediatelyAndClose(err error, f *message.Frame) {
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
	c.close()
}

// Sends a STOMP frame to the client immediately, does not push onto the
// write channel to be processed in turn.
func (c *Conn) sendImmediately(f *message.Frame) error {
	return c.writer.Write(f)
}

// Go routine for reading bytes from a client and assembling into
// STOMP frames. Also handles heart-beat read timeout. All read
// frames are pushed onto the read channel to be processed by the
// processLoop go-routine. This keeps all processing of frames for
// this connection on the one go-routine and avoids race conditions.
func (c *Conn) readLoop() {
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
					log.Println("read failed:", err, ":", c.rw.RemoteAddr())
				}
				c.close()
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
func (c *Conn) processLoop() {
	defer c.cleanupSubscriptions()

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
				c.close()
				return
			}

			// if the frame just sent to the client is an error
			// frame, we disconnect
			if f.Command == message.ERROR {
				// sent an ERROR frame, so disconnect
				c.close()
				return
			}

		case f := <-c.readChannel:
			// Just received a frame from the client.
			// Validate the frame, checking for mandatory
			// headers and prohibited headers.
			err := f.Validate()
			if err != nil {
				c.sendErrorImmediatelyAndClose(err, f)
				return
			}

			// Pass to the appropriate function for handling
			// according to the current state of the connection.
			err = c.stateFunc(c, f)
			if err != nil {
				c.sendErrorImmediatelyAndClose(err, f)
				return
			}

		case _ = <-timerChannel:
			// write a heart-beat
			err := c.writer.Write(nil)
			if err != nil {
				c.close()
				return
			}
		}
	}
}

// Called when the connection is closing, and takes care of
// unsubscribing all subscriptions with the upper layer, and
// re-queueing all unacknowledged messages to the upper layer.
func (c *Conn) cleanupSubscriptions() {
	// tell the upper layer we have disconnected

	// TODO: get all subscriptions and send appropriate UnsubscribeOp,
	// followed by RequeueOps.

	//c.requestChannel <- request{op: disconnectOp, conn: c}

	panic("not implemented yet: cleanupSubscriptions")
}

func (c *Conn) close() {
	c.closed = true  // records that the server closed the connection
	c.txStore.Init() // clean out pending transactions
	c.rw.Close()     // close the socket
}

func (c *Conn) handleConnect(f *message.Frame) error {
	var err error

	if _, ok := f.Contains(message.Receipt); ok {
		// CONNNECT and STOMP frames are not allowed to have
		// a receipt header.
		return receiptInConnect
	}

	// if either of these fields are absent, pass nil to the
	// authenticator function.
	login, _ := f.Contains(message.Login)
	passcode, _ := f.Contains(message.Passcode)
	if !c.config.Authenticate(login, passcode) {
		// sleep to slow down a rogue client a little bit
		time.Sleep(time.Second)
		return authenticationFailed
	}

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
	min := asMilliseconds(c.config.HeartBeat(), message.MaxHeartBeat)

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

	// tell the upper layer we are connected
	//	c.requestChannel <- request{op: connectOp, conn: c}

	return nil
}

func connecting(c *Conn, f *message.Frame) error {
	switch f.Command {
	case message.CONNECT, message.STOMP:
		return c.handleConnect(f)
	}
	return notConnected
}

// Sends a RECEIPT frame to the client if the frame f contains
// a receipt header. If the frame does contain a receipt header,
// it will be removed from the frame.
func (c *Conn) sendReceiptImmediately(f *message.Frame) error {
	if receipt, ok := f.Contains(message.Receipt); ok {
		// Remove the receipt header from the frame. This is handy
		// for transactions, because the frame has its receipt 
		// header removed prior to entering the transaction store.
		// When the frame is processed upon transaction commit, it
		// will not have a receipt header anymore.
		f.Remove(message.Receipt)
		return c.sendImmediately(message.NewFrame(message.RECEIPT, message.ReceiptId, receipt))
	}
	return nil
}

func (c *Conn) handleDisconnect(f *message.Frame) error {
	// As soon as we receive a DISCONNECT frame from a client, we do
	// not want to send any more frames to that client, with the exception
	// of a RECEIPT frame if the client has requested one.
	// Ignore the error condition if we cannot send a RECEIPT frame,
	// as the connection is about to close anyway.
	_ = c.sendReceiptImmediately(f)
	c.close()
	return nil
}

func (c *Conn) handleBegin(f *message.Frame) error {
	// the frame should already have been validated for the
	// transaction header, but we check again here.
	if transaction, ok := f.Contains(message.Transaction); ok {
		// Send a receipt and remove the header
		err := c.sendReceiptImmediately(f)
		if err != nil {
			return err
		}

		return c.txStore.Begin(transaction)
	}
	return missingHeader
}

func (c *Conn) handleCommit(f *message.Frame) error {
	// the frame should already have been validated for the
	// transaction header, but we check again here.
	if transaction, ok := f.Contains(message.Transaction); ok {
		// Send a receipt and remove the header
		err := c.sendReceiptImmediately(f)
		if err != nil {
			return err
		}
		return c.txStore.Commit(transaction, func(f *message.Frame) error {
			// Call the state function (again) for each frame in the
			// transaction. This time each frame is stripped of its transaction
			// header (and its receipt header as well, if it had one).
			return c.stateFunc(c, f)
		})
	}
	return missingHeader
}

func (c *Conn) handleAbort(f *message.Frame) error {
	// the frame should already have been validated for the
	// transaction header, but we check again here.
	if transaction, ok := f.Contains(message.Transaction); ok {
		// Send a receipt and remove the header
		err := c.sendReceiptImmediately(f)
		if err != nil {
			return err
		}
		return c.txStore.Abort(transaction)
	}
	return missingHeader
}

// Handle a SEND frame received from the client. Note that
// this method is called after a SEND message is received,
// but also after a transaction commit.
func (c *Conn) handleSend(f *message.Frame) error {
	// Send a receipt and remove the header
	err := c.sendReceiptImmediately(f)
	if err != nil {
		return err
	}

	if tx, ok := f.Contains(message.Transaction); ok {
		// the transaction header is removed from the frame
		err = c.txStore.Add(tx, f)
		if err != nil {
			return err
		}
	} else {
		// not in a transaction, send to be processed
		c.requestChannel <- Request{Op: EnqueueOp, Frame: f}
	}

	return nil
}

// Send the frame to the request channel. Remove receipt header
// and send a RECEIPT frame to the client if necessary.
func (c *Conn) sendFrameRequest(f *message.Frame) error {
	// Send a receipt and remove the header
	err := c.sendReceiptImmediately(f)
	if err != nil {
		return err
	}

	// Handled by next level
	c.requestChannel <- Request{Op: EnqueueOp, Frame: f}
	return nil
}

func (c *Conn) handleSubscribe(f *message.Frame) error {
	return c.sendFrameRequest(f)
}

func (c *Conn) handleUnsubscribe(f *message.Frame) error {
	return c.sendFrameRequest(f)
}

func (c *Conn) handleAck(f *message.Frame) error {
	return c.sendFrameRequest(f)
}

func (c *Conn) handleNack(f *message.Frame) error {
	return c.sendFrameRequest(f)
}

func connected(c *Conn, f *message.Frame) error {
	switch f.Command {
	case message.CONNECT, message.STOMP:
		return unexpectedCommand
	case message.DISCONNECT:
		return c.handleDisconnect(f)
	case message.BEGIN:
		return c.handleBegin(f)
	case message.ABORT:
		return c.handleAbort(f)
	case message.COMMIT:
		return c.handleCommit(f)
	case message.SEND:
		return c.handleSend(f)
	case message.SUBSCRIBE:
		return c.handleSubscribe(f)
	case message.UNSUBSCRIBE:
		return c.handleUnsubscribe(f)
	case message.ACK:
		return c.handleAck(f)
	case message.NACK:
		return c.handleNack(f)
	case message.MESSAGE, message.RECEIPT, message.ERROR:
		// should only be sent by the server, should not come from the client
		return unexpectedCommand
	default:
		return unknownCommand
	}
	panic("not reached")
}
