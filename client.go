package stomp

import (
	"errors"
	"github.com/jjeffery/stomp/message"
	"io"
	"log"
	"net"
	"strconv"
)

// The AckMode type is an enumeration of the acknowledgement modes for a 
// STOMP subscription. Valid values are AckAuto, AckClient and 
// AckClientIndividual, which are documented in the constants section.
type AckMode string

// These constants are the valid values for the AckMode type. When a STOMP
// client subscribes to a destination on the server, it specifies how it
// will acknowledge messages it receives from the server.
const (
	// No acknowledgement is required, the server assumes that the client 
	// received the message.
	AckAuto = AckMode("auto")

	// Client acknowledges messages. When a client acknowledges a message,
	// any previously received messages are also acknowledged.
	AckClient = AckMode("client")

	// Client acknowledges message. Each message is acknowledged individually.
	AckClientIndividual = AckMode("client-individual")
)

// A Client is a STOMP client.
type Client struct {
	Login    string // Login for authentication
	Passcode string // Passcode for authentication

	version string // STOMP protocol version
	readCh  chan *message.Frame
	writeCh chan writeRequest
}

type writeRequest struct {
	Frame *message.Frame      // frame to send
	C     chan *message.Frame // response channel
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(rw io.ReadWriter, headers map[string]string) error {
	c.readCh = make(chan *message.Frame, 8)
	c.writeCh = make(chan writeRequest, 8)
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
		return err
	}

	if response.Command != message.CONNECTED {
		return newError(response)
	}

	if version, ok := response.Contains(message.Version); ok {
		c.version = version
	} else {
		c.version = "1.0"
	}

	go readLoop(c, reader)
	go processLoop(c, writer)

	return nil
}

// readLoop is a goroutine that reads frames from the
// reader and places them onto a channel for processing
// by the processLoop goroutine
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

// processLoop is a goroutine that handles io with
// the server.
func processLoop(c *Client, writer *message.Writer) {
	channels := make(map[string]chan *message.Frame)

	for {
		select {
		case f, ok := <-c.readCh:
			if !ok {
				err := newErrorMessage("connection closed")
				sendError(channels, err)
				return
			}

			switch f.Command {
			case message.RECEIPT:
				if id, ok := f.Contains(message.ReceiptId); ok {
					if ch, ok := channels[id]; ok {
						ch <- f
						delete(channels, id)
						close(ch)
					}

				} else {
					err := &Error{Message: "missing receipt-id", Frame: f}
					sendError(channels, err)
					return
				}

			case message.ERROR:
				log.Println("received ERROR")
				for _, ch := range channels {
					ch <- f
					close(ch)
				}

				return

			case message.MESSAGE:
				if id, ok := f.Contains(message.Subscription); ok {
					if ch, ok := channels[id]; ok {
						ch <- f
					} else {
						log.Println("ignored MESSAGE for subscription", id)
					}
				}
			}

		case req, ok := <-c.writeCh:
			if !ok {
				sendError(channels, errors.New("write channel closed"))
			}
			if req.C != nil {
				if receipt, ok := req.Frame.Contains(message.Receipt); ok {
					// remember the channel for this receipt
					channels[receipt] = req.C
				}
			}

			switch req.Frame.Command {
			case message.SUBSCRIBE:
				id, _ := req.Frame.Contains(message.Id)
				channels[id] = req.C
			case message.UNSUBSCRIBE:
				id, _ := req.Frame.Contains(message.Id)
				// is this trying to be too clever -- add a receipt
				// header so that when the server responds with a 
				// RECEIPT frame, the corresponding channel will be closed
				req.Frame.Set(message.Receipt, id)
			}

			// frame to send
			err := writer.Write(req.Frame)
			if err != nil {
				sendError(channels, err)
				return
			}
		}
	}
}

// Send an error to all receipt channels.
func sendError(m map[string]chan *message.Frame, err error) {
	frame := message.NewFrame(message.ERROR, message.Message, err.Error())
	for _, ch := range m {
		ch <- frame
	}
}

func (c *Client) Disconnect() error {
	ch := make(chan *message.Frame)
	c.writeCh <- writeRequest{
		Frame: message.NewFrame(message.DISCONNECT, message.Receipt, allocateId()),
		C:     ch,
	}

	response := <-ch
	if response.Command != message.RECEIPT {
		return newError(response)
	}

	// TODO: should we do anything to close the connection?
	// not easy to do, seeing as we only have a ReadWriter.

	return nil
}

func (c *Client) Send(msg SendMessage) error {
	if msg.Destination == "" {
		return errors.New("no destination specififed")
	}
	f := message.NewFrame(message.SEND, message.Destination, msg.Destination)
	if msg.ContentType != "" {
		f.Append(message.ContentType, msg.ContentType)
	}
	f.Append(message.ContentLength, strconv.Itoa(len(msg.Body)))
	f.Body = msg.Body

	request := writeRequest{Frame: f}

	if msg.Receipt {
		request.C = make(chan *message.Frame)
		c.writeCh <- request
		response := <-request.C
		if response.Command == message.RECEIPT {
			return nil
		}
		return newError(response)
	}

	// no receipt required, so send and assume success
	c.writeCh <- request
	return nil
}

// Subscribe to a destination. Returns a channel for receiving message frames.
func (c *Client) Subscribe(destination string, ack AckMode) (*Subscription, error) {
	ch := make(chan *message.Frame)
	id := allocateId()
	request := writeRequest{
		Frame: message.NewFrame(message.SUBSCRIBE,
			message.Id, id,
			message.Destination, destination,
			message.Ack, string(ack)),
		C: ch,
	}

	sub := &Subscription{
		id:          id,
		destination: destination,
		client:      c,
		ackMode:     ack,
		C:           make(chan *Message, 16),
	}
	go sub.readLoop(ch)

	c.writeCh <- request
	return sub, nil
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
	Destination string            // Destination
	ContentType string            // MIME content type
	Receipt     bool              // Is a receipt required
	Headers     map[string]string // Optional headers
	Body        []byte            // Content of message
}
