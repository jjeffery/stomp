package stomp

import (
	"errors"
	"github.com/jjeffery/stomp/message"
	"strconv"
)

// A Message is a message that is sent to or received from the server.
type Message struct {
	// Destination the message is sent to. The STOMP server should
	// in turn send this message to a STOMP clients that has subscribed 
	// to the destination.
	Destination string

	// MIME content type.
	ContentType string // MIME content

	// Connection that the message was received on. 
	// Ignored for message sent to the server.
	Conn *Conn

	// Subscription associated with the message.
	// Ignored for messages sent to the server.
	Subscription *Subscription

	// Optional header entries. When received from the server,
	// these are the header entries received with the message.
	// When sending to the server, these are optional header entries
	// that accompany the message to its destination.
	Header

	// The message body, which is an arbitrary sequence of bytes.
	// The ContentType indicates the format of this body.
	Body []byte // Content of message
}

func (msg *Message) createSendFrame() (*message.Frame, error) {
	if msg.Destination == "" {
		return nil, errors.New("no destination specififed")
	}
	f := message.NewFrame(message.SEND, message.Destination, msg.Destination)
	if msg.ContentType != "" {
		f.Append(message.ContentType, msg.ContentType)
	}
	f.Append(message.ContentLength, strconv.Itoa(len(msg.Body)))
	f.Body = msg.Body

	for key, values := range msg.Header {
		for _, value := range values {
			f.Append(key, value)
		}
	}

	return f, nil
}
