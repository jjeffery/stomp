package client

import (
	"fmt"
	"github.com/jjeffery/stomp/message"
	"io"
	"log"
	"net"
	"time"
)

const (
	// Maximum number of pending frames allowed to a client.
	// before a disconnect occurs. If the client cannot keep
	// up with the server, we do not want the server to backlog
	// pending frames indefinitely.
	MaxPendingWrites = 16
)

// Connection with client
type Connection struct {
	conn           net.Conn
	requestChannel chan Request
	writeChannel   chan *message.Frame
	stateFunc      func(c *Connection, f *message.Frame) error
	readTimeout    time.Duration
	writeTimeout   time.Duration
	version        message.StompVersion
}

func newConnection(conn net.Conn, channel chan Request) *Connection {
	c := new(Connection)
	c.conn = conn
	c.requestChannel = channel
	c.writeChannel = make(chan *message.Frame, 32)
	channel <- Request{Type: Create, Connection: c}
	go c.readLoop()
	go c.writeLoop()
	return c
}

// Write a frame to the connection.
func (c *Connection) Send(f *message.Frame) {
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
func (c *Connection) SendError(err error) {
	f := new(message.Frame)
	f.Command = message.ERROR
	f.Headers.Append(message.Message, err.Error())
	c.Send(f) // will close after successful send
}

func (c *Connection) readLoop() {
	reader := message.NewReader(c.conn)
	for {
		if c.readTimeout == time.Duration(0) {
			// infinite timeout
			c.conn.SetReadDeadline(time.Time{})
		} else {
			c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
		}
		f, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				log.Println("connection closed:", c.conn.RemoteAddr())
			} else {
				log.Println("read failed:", err, c.conn.RemoteAddr())
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

func (c *Connection) writeLoop() {
	writer := message.NewWriter(c.conn)
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

func (c *Connection) Close() {
	c.conn.Close()
	c.requestChannel <- Request{Type: Disconnect, Connection: c}
}

func (c *Connection) handleConnect(f *message.Frame) error {
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

func connecting(c *Connection, f *message.Frame) error {
	switch f.Command {
	case message.CONNECT, message.STOMP:
		return c.handleConnect(f)
	}
	return notConnected
}
