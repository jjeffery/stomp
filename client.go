package stomp

import (
	"github.com/jjeffery/stomp/message"
)

// Indicates the acknowledgement mode for a STOMP subscription.
type AckMode int

// AckMode constants
const (
	AckAuto AckMode = iota
	AckClient
	AckClientIndividual
)

// A Client is a STOMP client.
type Client struct {
	// TODO
}

type Subscription struct {
	C chan *message.Frame
	// TODO other members: client, id, etc
}

func (s *Subscription) Unsubscribe() error {
	panic("not implemented")
}

func NewClient(addr string) *Client {
	panic("not implemented")
}

func (c *Client) Connect() error {
	panic("not implemented")
}

func (c *Client) ConnectAuth(login, passcode string) error {
	panic("not implemented")
}

func (c *Client) Disconnect() error {
	panic("not implemented")
}

// Subscribe to a destination. Returns a channel for receiving message frames.
func (c *Client) Subscribe(destination string, ack AckMode) (*Subscription, error) {
	panic("not implemented")
}
