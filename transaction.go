package stomp

import (
	"github.com/jjeffery/stomp/frame"
)

// A Transaction applies to the sending of messages to the STOMP server,
// and the acknowledgement of messages received from the STOMP server.
// All messages sent and and acknowledged in the context of a transaction
// are processed atomically by the STOMP server.
//
// Transactions are committed with the Commit method. When a transaction is
// committed, all sent messages, acknowledgements and negative acknowledgements, 
// are processed by the STOMP server. Alternatively transactions can be aborted,
// in which case all sent messages, acknowledgements and negative
// acknowledgements are discarded by the STOMP server.
type Transaction struct {
	id        string
	conn      *Conn
	completed bool
}

// Id returns the unique identifier for the transaction.
func (tx *Transaction) Id() string {
	return tx.id
}

// Conn returns the connection associated with this transaction.
func (tx *Transaction) Conn() *Conn {
	return tx.conn
}

// Abort will abort the transaction. Any calls to Send, SendWithReceipt,
// Ack and Nack on this transaction will be discarded.
func (tx *Transaction) Abort() error {
	if tx.completed {
		return completedTransaction
	}

	f := NewFrame(frame.ABORT, frame.Transaction, tx.id)
	tx.conn.sendFrame(f)
	tx.completed = true

	return nil
}

// Commit will commit the transaction. All messages and acknowledgements
// sent to the STOMP server on this transaction will be processed atomically.
func (tx *Transaction) Commit() error {
	if tx.completed {
		return completedTransaction
	}

	f := NewFrame(frame.COMMIT, frame.Transaction, tx.id)
	tx.conn.sendFrame(f)
	tx.completed = true

	return nil
}

// Send a message to the server as part of the transaction. The STOMP
// server will not process the message until the transaction has been
// committed.
func (tx *Transaction) Send(msg Message) error {
	if tx.completed {
		return completedTransaction
	}

	f, err := msg.createSendFrame()
	if err != nil {
		return err
	}

	f.Header.Set(frame.Transaction, tx.id)
	tx.conn.sendFrame(f)
	return nil
}

// SendWithReceipt sends a message to the server and waits for the 
// server to acknowledge receipt. Although the STOMP server has received
// the message, it will not process the message until the transaction 
// has been committed.
func (tx *Transaction) SendWithReceipt(msg *Message) error {
	if tx.completed {
		return completedTransaction
	}

	f, err := msg.createSendFrame()
	if err != nil {
		return err
	}

	f.Set(frame.Transaction, tx.id)
	return tx.conn.sendFrameWithReceipt(f)
}

// Ack sends an acknowledgement for the message to the server. The STOMP
// server will not process the acknowledgement until the transaction
// has been committed. If the subscription has an AckMode of AckAuto, calling
// this function has no effect.
func (tx *Transaction) Ack(msg *Message) error {
	if tx.completed {
		return completedTransaction
	}

	f, err := tx.conn.createAckNackFrame(msg, true)
	if err != nil {
		return err
	}

	if f != nil {
		f.Set(frame.Transaction, tx.id)
		tx.conn.sendFrame(f)
	}

	return nil
}

// Nack sends a negative acknowledgement for the message to the server,
// indicating that this client cannot or will not process the message and
// that it should be processed elsewhere. The STOMP server will not process
// the negative acknowledgement until the transaction has been committed.
// It is an error to call this method if the subscription has an AckMode
// of AckAuto, because the STOMP server will not be expecting any kind
// of acknowledgement (positive or negative) for this message.
func (tx *Transaction) Nack(msg *Message) error {
	if tx.completed {
		return completedTransaction
	}

	f, err := tx.conn.createAckNackFrame(msg, false)
	if err != nil {
		return err
	}

	if f != nil {
		f.Set(frame.Transaction, tx.id)
		tx.conn.sendFrame(f)
	}

	return nil
}
