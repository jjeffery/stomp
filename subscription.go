package stomp

import (
	"fmt"
	"github.com/jjeffery/stomp/frame"
	"github.com/jjeffery/stomp/message"
	"log"
)

// The Subscription type represents a client subscription to
// a destination. The subscription is created by calling Conn.Subscribe.
//
// Once a client has subscribed, it can receive messages from the C channel.
type Subscription struct {
	C           chan *Message
	id          string
	destination string
	client      *Conn
	ackMode     AckMode
}

// BUG(jpj): If the client does not read messages from the Subscription.C 
// channel quickly enough, the client will stop reading messages from the 
// server.

// Identification for this subscription. Unique among
// all subscriptions for the same Client.
func (s *Subscription) Id() string {
	return s.id
}

// Destination for which the subscription applies.
func (s *Subscription) Destination() string {
	return s.destination
}

// The Acknowledgement mode for the subscription.
func (s *Subscription) AckMode() AckMode {
	return s.ackMode
}

// Unsubscribes and closes the channel C.
func (s *Subscription) Unsubscribe() error {
	_ = message.NewFrame(frame.UNSUBSCRIBE, frame.Id, s.id)
	panic("not implemented")
}

// Read a message from the subscription
func (s *Subscription) Read() (*Message, error) {
	panic("not implemented")
}

func (s *Subscription) readLoop(ch chan *Frame) {
	for {
		f, ok := <-ch
		if !ok {
			return
		}

		if f.Command == frame.MESSAGE {
			destination, _ := f.Contains(frame.Destination)
			contentType, _ := f.Contains(frame.ContentType)
			msg := &Message{
				Destination:  destination,
				ContentType:  contentType,
				Conn:         s.client,
				Subscription: s,
				Header:       f.Header.Clone(),
				Body:         f.Body,
			}
			s.C <- msg
		} else if f.Command == frame.ERROR {
			message, _ := f.Contains(frame.Message)
			text := fmt.Sprintf("ERROR message:%s", message)
			log.Println(text)
			return
		}

	}
}
