package client

import (
	"github.com/jjeffery/stomp/message"
)

type Subscription struct {
	conn    *Conn
	dest    string
	id      string            // client's subscription id
	ack     string            // auto, client, client-individual
	subList *SubscriptionList // am I in a list
	frame   *message.Frame    // message allocated to subscription
}

func newSubscription(c *Conn, dest string, id string, ack string) *Subscription {
	return &Subscription{
		conn: c,
		dest: dest,
		id:   id,
		ack:  ack,
	}
}

func (s *Subscription) Destination() string {
	return s.dest
}

func (s *Subscription) Send(f *message.Frame) {
	if s.frame != nil {
		panic("subscription already has a frame pending")
	}
	s.frame = f
	f.Set(message.Id, s.id)

	// let the connection deal with the sub
	s.conn.subChannel <- s
}
