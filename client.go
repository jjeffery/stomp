package stomp

import (
	"github.com/jjeffery/stomp/message"
	"io"
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
	rw      io.ReadWriter // Underlying network connection 
}

func (c *Client) Connect(rw io.ReadWriter) error {
	panic("not implemented")
}

func readLoop(c *Client) {
	reader := message.NewReader(c.rw)
	for {
		f, err := reader.Read()
		if err != nil {
			close(c.readCh)
			return
		}
		c.readCh <- f
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
