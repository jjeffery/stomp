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

// represents a connection with the client
type conn struct {
	rw             net.Conn
	requestChannel chan request
	writeChannel   chan *message.Frame
	stateFunc      func(c *conn, f *message.Frame) error
	readTimeout    time.Duration
	writeTimeout   time.Duration
	version        message.StompVersion
}

func newConn(rw net.Conn, channel chan request) *conn {
	c := new(conn)
	c.rw = rw
	c.requestChannel = channel
	c.writeChannel = make(chan *message.Frame, maxPendingWrites)
	go c.readLoop()
	go c.writeLoop()
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

func (c *conn) readLoop() {
	reader := message.NewReader(c.rw)
	for {
		if c.readTimeout == time.Duration(0) {
			// infinite timeout
			c.rw.SetReadDeadline(time.Time{})
		} else {
			c.rw.SetReadDeadline(time.Now().Add(c.readTimeout))
		}
		f, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				log.Println("connection closed:", c.rw.RemoteAddr())
			} else {
				log.Println("read failed:", err, c.rw.RemoteAddr())
			}
			c.Close()
			return
		}

		if f == nil {
			// if the frame is nil, then it is a heartbeat
			continue
		}

		err = c.stateFunc(c, f)

		if err != nil {
			c.SendError(err)
			c.Close()
		}
	}
}

func (c *conn) writeLoop() {
	writer := message.NewWriter(c.rw)
	for {
		var timerChannel <-chan time.Time
		var timer *time.Timer

		if c.writeTimeout > 0 {
			timer = time.NewTimer(c.writeTimeout)
			timerChannel = timer.C
		}

		select {
		case f := <-c.writeChannel:
			// stop the heart-beat timer
			if timer != nil {
				timer.Stop()
				timer = nil
			}

			// write the frame to the client
			err := writer.Write(f)
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

		case _ = <-timerChannel:
			// write a heart-beat
			err := writer.Write(nil)
			if err != nil {
				c.Close()
				return
			}
		}
	}
}

func (c *conn) Close() {
	c.rw.Close()
	c.requestChannel <- request{op: disconnectOp, conn: c}
}

func (c *conn) handleConnect(f *message.Frame) error {
	var err error

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

	// apply a minimum heartbeat time of 30 seconds
	if cx > 0 && cx < 30000 {
		cx = 30000
	}
	if cy > 0 && cy < 30000 {
		cy = 30000
	}

	c.readTimeout = time.Duration(cx) * time.Millisecond
	c.writeTimeout = time.Duration(cy) * time.Millisecond

	response := message.NewFrame(message.CONNECTED,
		message.Version, string(c.version),
		message.Server, "stompd/x.y.z") // TODO: get version

	if c.version > message.V1_0 {
		value := fmt.Sprintf("%d,%d", cy, cx)
		response.Append(message.HeartBeat, value)
	}

	c.Send(response)

	return nil
}

func connecting(c *conn, f *message.Frame) error {
	switch f.Command {
	case message.CONNECT, message.STOMP:
		return c.handleConnect(f)
	}
	return notConnected
}
