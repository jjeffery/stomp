package stomp

import (
	"github.com/jjeffery/stomp/message"
	"io"
	"net"
)

// The AckMode type is an enumeration of the acknowledgement modes for a STOMP subscription. 
// Valid values are AckAuto, AckClient and AckClientIndividual, which are documented in the
// constants section.
type AckMode string

// These constants are the valid values for the AckMode type. When a STOMP
// client subscribes to a destination on the server, it specifies how it
// will acknowledge messages it receives from the server.
const (
	// No acknowledgement is required, the server assumes that the client received the message.
	AckAuto = AckMode("auto")

	// Client acknowledges messages. When a client acknowledges a message, any previously 
	// received messages are also acknowledged.
	AckClient = AckMode("client")

	// Client acknowledges message. Each message is acknowledged individually.
	AckClientIndividual = AckMode("client-individual")
)

// A Client is a STOMP client.
type Client struct {
	Login    string // Login for authentication
	Passcode string // Passcode for authentication

	readCh  chan *message.Frame
	writeCh chan *message.Frame
}

func (c *Client) Connect(rw io.ReadWriter, headers map[string]string) error {
	c.readCh = make(chan *message.Frame, 8)
	c.writeCh = make(chan *message.Frame, 8)
	reader := message.NewReader(rw)
	writer := message.NewWriter(rw)

	connectFrame := message.NewFrame(message.CONNECT)
	for key, value := range headers {
		connectFrame.Append(key, value)
	}

	// ensure mandatory header "accept-version" is set
	if _, ok := connectFrame.Contains(message.AcceptVersion); !ok {
		connectFrame.Append(message.AcceptVersion, "1.1,1.2")
	}

	// ensure mandatory header "host" is set
	if _, ok := connectFrame.Contains(message.Host); !ok {
		// no host, try to get it from the network connection
		if conn, ok := rw.(net.Conn); ok {
			host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
			if err != nil {
				return err
			}
			connectFrame.Append(message.Host, host)
		} else {
			// not a network connection, host is unknown
			connectFrame.Append(message.Host, "unknown")
		}
	}
	writer.Write(connectFrame)
	response, err := reader.Read()
	if err != nil {
		println("reader.Read failed")
		return err
	}

	if response.Command != message.CONNECTED {
		return NewError(response)
	}

	go readLoop(c, reader)

	return nil
}

func readLoop(c *Client, reader *message.Reader) {
	for {
		f, err := reader.Read()
		if err != nil {
			close(c.readCh)
			return
		}
		c.readCh <- f
	}
}

func processLoop(c *Client, writer *message.Writer) {
	for {
		select {
		case f, ok := <-c.readCh:
			// TODO process incoming frame

		case f, ok := <-c.writeCh:
			// frame to send
			err := writer.Write(f)
			if err != nil {

			}
		}
	}
}

func (c *Client) Disconnect() error {
	panic("not implemented")
}

// Subscribe to a destination. Returns a channel for receiving message frames.
func (c *Client) Subscribe(destination string, ack AckMode) (*Subscription, error) {
	panic("not implemented")
}

func (c *Client) Ack(m *Message) error {
	panic("not implemented")
}

func (c *Client) Nack(m *Message) error {
	panic("not implemented")
}

func (c *Client) Begin() (*Transaction, error) {
	panic("not implemented")
}

type Transaction struct {
}

func (tx *Transaction) Abort() error {
	panic("not implemented")
}

func (tx *Transaction) Commit() error {
	panic("not implemented")
}

func (tx *Transaction) Send(msg *SendMessage) error {
	panic("not implemented")
}

func (tx *Transaction) Ack(m *Message) error {
	panic("not implemented")
}

func (tx *Transaction) Nack(m *Message) error {
	panic("not implemented")
}

// A SendMessage is a message that is sent to the server.
type SendMessage struct {
	Destination string  // Destination
	ContentType string  // MIME content type
	Receipt     bool    // Is a receipt required
	Headers     Headers // Optional headers
	Body        []byte  // Content of message
}
