package stomp

import (
	"errors"
	"net"
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
	conn         net.Conn
	channel      chan Request
	writeChannel chan *Frame
}

// Represents a request received from the client,
// consisting of a frame and the connection it
// was received from
type Request struct {
	Frame      *Frame
	Connection *Connection
	Error      error
}

func NewConnection(conn net.Conn, channel chan Request) *Connection {
	c := new(Connection)
	c.conn = conn
	c.channel = channel
	c.writeChannel = make(chan *Frame, 32)
	go c.ReadLoop()
	go c.WriteLoop()
	return c
}

// Write a frame to the connection. TODO: caller blocks, need to introduce
// another channel and a go routine to read from the channel and write to
// the other party.
func (c *Connection) Send(f *Frame) {
	// place the frame on the write channel, or
	// close the connection if the write channel is full,
	// as this means the client is not keeping up.
	select {
	case c.writeChannel <- f:
	default:
		// write channel is full
		c.conn.Close()
		c.channel <- Request{
			Connection: c,
			Error:      errors.New("client blocked, connection closed"),
		}
		return
	}
}

// TODO: should send other information, such as receipt-id
func (c *Connection) SendError(err error) {
	f := new(Frame)
	f.Command = Error
	f.Headers.Append(Message, err.Error())
	c.Send(f) // will close after successful send
}

func (c *Connection) ReadLoop() {
	reader := NewReader(c.conn)
	for {
		f, err := reader.Read()
		if err != nil {
			c.conn.Close()
			c.channel <- Request{Connection: c, Error: err}
			return
		}

		if f == nil {
			// TODO: received a heart-beat from the client,
			// so restart the read timer
		} else {
			c.channel <- Request{Frame: f, Connection: c}
		}
	}
}

func (c *Connection) WriteLoop() {
	writer := NewWriter(c.conn)
	for {
		f := <-c.writeChannel
		err := writer.Write(f)
		if err != nil {
			c.conn.Close()
			c.channel <- Request{Connection: c, Error: err}
			return
		}
		if f.Command == Error {
			// sent an ERROR frame, so disconnect
			c.conn.Close()
			c.channel <- Request{Connection: c, Error: errors.New("closed after ERROR frame sent")}
			return
		}
	}
}
