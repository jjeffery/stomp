package stomp

import (
	"errors"
	"github.com/jjeffery/stomp/frame"
	"github.com/jjeffery/stomp/message"
	"io"
	"log"
	"net"
)

// Options for connecting to the STOMP server
type Options struct {
	// Login and passcode for authentication with the STOMP server.
	// If no authentication is required, leave blank.
	Login, Passcode string

	// Value for the "host" header entry when connecting to the
	// STOMP server. Leave blank for default value.
	Host string

	// Comma-separated list of acceptable STOMP versions. 
	// Leave blank for default protocol negotiation, which is
	// the recommended setting.
	AcceptVersion string

	// Value to pass in the "heart-beat" header entry when connecting
	// to the STOMP server. Format is two non-negative integer values
	// separated by a comma. Leave blank for default heart-beat negotiation,
	// which is the recommended setting.
	HeartBeat string

	// Other header entries for STOMP servers that accept non-standard
	// header entries in the CONNECT frame.
	NonStandard map[string]string
}

// A Conn is a connection to a STOMP server. Create a Conn using either
// the Dial or Connect function.
type Conn struct {
	conn    io.ReadWriteCloser
	readCh  chan *message.Frame
	writeCh chan writeRequest
	version string
	session string
}

type writeRequest struct {
	Frame *message.Frame      // frame to send
	C     chan *message.Frame // response channel
}

func Dial(network, addr string, opts Options) (*Conn, error) {
	c, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	if opts.Host == "" {
		host, _, err := net.SplitHostPort(c.RemoteAddr().String())
		if err != nil {
			c.Close()
			return nil, err
		}
		opts.Host = host
	}

	return Connect(c, opts)
}

func Connect(conn io.ReadWriteCloser, opts Options) (*Conn, error) {
	reader := message.NewReader(conn)
	writer := message.NewWriter(conn)

	// set default values
	if opts.AcceptVersion == "" {
		opts.AcceptVersion = "1.0,1.1,1.2"
	}
	if opts.HeartBeat == "" {
		opts.HeartBeat = "60000,60000"
	}
	if opts.Host == "" {
		// Attempt to get host from net.Conn object if possible
		if connection, ok := conn.(net.Conn); ok {
			host, _, err := net.SplitHostPort(connection.RemoteAddr().String())
			if err == nil {
				opts.Host = host
			}
		}

		// If host is still blank, use default
		if opts.Host == "" {
			opts.Host = "default"
		}
	}

	connectFrame := message.NewFrame(frame.CONNECT,
		frame.Host, opts.Host,
		frame.AcceptVersion, opts.AcceptVersion,
		frame.HeartBeat, opts.HeartBeat)
	if opts.Login != "" || opts.Passcode != "" {
		connectFrame.Append(frame.Login, opts.Login)
		connectFrame.Append(frame.Passcode, opts.Passcode)
	}
	for key, value := range opts.NonStandard {
		connectFrame.Append(key, value)
	}

	writer.Write(connectFrame)
	response, err := reader.Read()
	if err != nil {
		return nil, err
	}

	if response.Command != frame.CONNECTED {
		return nil, newError(response)
	}

	c := &Conn{
		conn:    conn,
		readCh:  make(chan *message.Frame, 8),
		writeCh: make(chan writeRequest, 8),
	}

	if version, ok := response.Contains(frame.Version); ok {
		c.version = version
	} else {
		c.version = "1.0"
	}

	c.session, _ = response.Contains(frame.Session)

	// TODO(jpj): make any non-standard headers in the CONNECTED
	// frame available.

	go readLoop(c, reader)
	go processLoop(c, writer)

	return c, nil
}

func (c *Conn) Version() string {
	return c.version
}

func (c *Conn) Session() string {
	return c.session
}

// readLoop is a goroutine that reads frames from the
// reader and places them onto a channel for processing
// by the processLoop goroutine
func readLoop(c *Conn, reader *message.Reader) {
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
func processLoop(c *Conn, writer *message.Writer) {
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
			case frame.RECEIPT:
				if id, ok := f.Contains(frame.ReceiptId); ok {
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

			case frame.ERROR:
				log.Println("received ERROR")
				for _, ch := range channels {
					ch <- f
					close(ch)
				}

				return

			case frame.MESSAGE:
				if id, ok := f.Contains(frame.Subscription); ok {
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
				if receipt, ok := req.Frame.Contains(frame.Receipt); ok {
					// remember the channel for this receipt
					channels[receipt] = req.C
				}
			}

			switch req.Frame.Command {
			case frame.SUBSCRIBE:
				id, _ := req.Frame.Contains(frame.Id)
				channels[id] = req.C
			case frame.UNSUBSCRIBE:
				id, _ := req.Frame.Contains(frame.Id)
				// is this trying to be too clever -- add a receipt
				// header so that when the server responds with a 
				// RECEIPT frame, the corresponding channel will be closed
				req.Frame.Set(frame.Receipt, id)
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
	frame := message.NewFrame(frame.ERROR, frame.Message, err.Error())
	for _, ch := range m {
		ch <- frame
	}
}

// Disconnect will disconnect from the STOMP server. This function
// follows the suggested protocol for graceful disconnection: it
// sends a DISCONNECT frame with a receipt header element. Once the
// RECEIPT frame has been received, the connection with the STOMP
// server is closed and any further attempt to write to the server
// will fail.
func (c *Conn) Disconnect() error {
	ch := make(chan *message.Frame)
	c.writeCh <- writeRequest{
		Frame: message.NewFrame(frame.DISCONNECT, frame.Receipt, allocateId()),
		C:     ch,
	}

	response := <-ch
	if response.Command != frame.RECEIPT {
		return newError(response)
	}

	return c.conn.Close()
}

// SendWithReceipt sends a message to the STOMP server,
// and does not return until the STOMP server acknowledges
// receipt of the message.
//
// Note that this does not guarantee that the message has been
// delivered for processing, only that the STOMP server has received
// the message. Upon return, the message may be on a queue waiting
// to be processed.
func (c *Conn) SendWithReceipt(msg Message) error {
	f, err := msg.createSendFrame()
	if err != nil {
		return err
	}

	receipt := allocateId()
	f.Set(frame.Receipt, receipt)

	request := writeRequest{Frame: f}

	request.C = make(chan *message.Frame)
	c.writeCh <- request
	response := <-request.C
	if response.Command == frame.RECEIPT {
		// TODO: check receipt-id
		return nil
	}
	return newError(response)
}

// Send sends a message to the STOMP server, and does not
// wait for acknowledgement of receipt by the STOMP server.
func (c *Conn) Send(msg Message) error {
	f, err := msg.createSendFrame()
	if err != nil {
		return err
	}

	request := writeRequest{Frame: f}

	// no receipt required, so send and assume success
	c.writeCh <- request
	return nil
}

// Subscribe to a destination. Returns a channel for receiving message frames.
func (c *Conn) Subscribe(destination string, ack AckMode) (*Subscription, error) {
	ch := make(chan *message.Frame)
	id := allocateId()
	request := writeRequest{
		Frame: message.NewFrame(frame.SUBSCRIBE,
			frame.Id, id,
			frame.Destination, destination,
			frame.Ack, ack.String()),
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

func (c *Conn) Ack(m *Message) error {
	panic("not implemented")
}

func (c *Conn) Nack(m *Message) error {
	panic("not implemented")
}

func (c *Conn) Begin() (*Transaction, error) {
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

func (tx *Transaction) Send(msg *Message) error {
	panic("not implemented")
}

func (tx *Transaction) Ack(m *Message) error {
	panic("not implemented")
}

func (tx *Transaction) Nack(m *Message) error {
	panic("not implemented")
}
