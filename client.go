package stomp

import (
	_ "github.com/jjeffery/stomp/message"
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
	// TODO
}

func NewClient(addr string) *Client {
	panic("not implemented")
}

func (c *Client) SetLogin(login, passcode string) {
	panic("not implemented")
}

func (c *Client) Connect() error {
	panic("not implemented")
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

// The Subscription type represents a client subscription to
// a destination. The subscription is created by calling Client.Subscribe.
//
// Once a client has subscribed, it can receive messages from the C channel.
type Subscription struct {
	C      chan *Message
	client *Client
}

// BUG(jpj): If the client does not read messages from the Subscription.C channel quickly
// enough, the client will stop reading messages from the server.

// Identification for this subscription. Unique among
// all subscriptions for the same Client.
func (s *Subscription) Id() {
	panic("not implemented")
}

// Destination for which the subscription applies.
func (s *Subscription) Destination() {
	panic("not implemented")
}

// The Ack mode for the subscription: auto, client or client-individual.
func (s *Subscription) Ack() AckMode {
	panic("not implemented")
}

// Unsubscribes and closes the channel C.
func (s *Subscription) Unsubscribe() error {
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

// A Message is a message that is received from the server.
type Message struct {
	Destination  string        // Destination the message was sent to.
	ContentType  string        // MIME content
	Client       *Client       // Associated client
	Subscription *Subscription // Associated subscription
	Headers      Headers       // Optional headers
	Body         []byte        // Content of message
}

// The Headers interface represents a collection of headers, each having a key 
// and a value. There may be more than one header in the collection 
// with the same key, in which case the first header's value is used.
type Headers interface {
	// Returns the value associated with the specified key, and whether it was
	// found or not.
	Contains(key string) (string, bool)

	// Remove all headers with the specified key.
	Remove(key string)

	// Append the header to the end of the collection.
	Append(key, value string)

	// Replace any existing header with the same key, or append
	// if no header has the same key.
	Set(key, value string)

	// Get the header at the specified index.
	GetAt(index int) (key, value string)

	// Number of headers in the collection.
	Len() int
}
