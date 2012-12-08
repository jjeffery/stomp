package message

import (
	"errors"
	"strconv"
)

// Represents a single STOMP frame.
type Frame struct {
	// The frame command. Should be one of the standard STOMP commands. Note that
	// STOMP commands are case sensitive.
	Command string

	// Collection of frame headers.
	Headers

	// The frame body. Only SEND, MESSAGE and ERROR frames may have a body.
	// All other frames must not have a body.
	Body []byte
}

// Returns the value of the "content-length" header, and whether it was
// found or not. Used for deserializing a frame. If the content length
// is specified in the header, then the body can contain null characters.
// Otherwise the body is read until a null character is encountered.
// If an error is returned, then the content-length header is malformed.
func (f *Frame) ContentLength() (contentLength int, ok bool, err error) {
	text, ok := f.Headers.Contains(ContentLength)
	if !ok {
		return
	}

	value, err := strconv.ParseInt(text, 10, 32)
	if err != nil {
		ok = false
		return
	}

	contentLength = int(value)
	ok = true
	return
}

// Check frame for required headers
func (f *Frame) Validate() error {
	switch f.Command {
	case CONNECT, STOMP:
		return f.validateConnect()
	case CONNECTED:
		return f.validateConnected()
	case SEND:
		return f.validateSend()
	case SUBSCRIBE:
		return f.validateSubscribe()
	case UNSUBSCRIBE:
		return f.validateUnsubscribe()
	case ACK:
		return f.validateAck()
	case NACK:
		return f.validateNack()
	case BEGIN:
		return f.validateBegin()
	case COMMIT:
		return f.validateCommit()
	case ABORT:
		return f.validateAbort()
	case DISCONNECT:
		return f.validateDisconnect()
	case MESSAGE:
		return f.validateMessage()
	case RECEIPT:
		return f.validateReceipt()
	case ERROR:
		return f.validateError()
	}
	return invalidCommand
}

func (f *Frame) verifyRequiredHeaders(names ...string) error {
	for _, name := range names {
		if _, ok := f.Headers.Contains(name); !ok {
			return errors.New("missing header: " + name)
		}
	}
	return nil
}

func (f *Frame) validateConnect() error {
	// TODO: check for valid version
	// TODO: if version is >= 1.1 need to have accept-version and host
	return nil
}

func (f *Frame) validateConnected() error {
	return nil
}

func (f *Frame) validateSend() error {
	return f.verifyRequiredHeaders(Destination)
}

func (f *Frame) validateSubscribe() error {
	return f.verifyRequiredHeaders(Destination, Id)
}

func (f *Frame) validateUnsubscribe() error {
	return f.verifyRequiredHeaders(Id)
}

func (f *Frame) validateAck() error {
	return f.verifyRequiredHeaders(Id)
}

func (f *Frame) validateNack() error {
	return f.verifyRequiredHeaders(Id)
}

func (f *Frame) validateBegin() error {
	return f.verifyRequiredHeaders(Transaction)
}

func (f *Frame) validateAbort() error {
	return f.verifyRequiredHeaders(Transaction)
}

func (f *Frame) validateCommit() error {
	return f.verifyRequiredHeaders(Transaction)
}

func (f *Frame) validateDisconnect() error {
	return nil
}

func (f *Frame) validateMessage() error {
	return f.verifyRequiredHeaders(Destination, MessageId, Subscription)
}

func (f *Frame) validateReceipt() error {
	return f.verifyRequiredHeaders(ReceiptId)
}

func (f *Frame) validateError() error {
	return nil
}
