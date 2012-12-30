package message

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	// Maximum content length allowed 16MB
	MaxContentLength = 16 * 1024 * 1024
)

var (
	// regexp for heart-beat header value
	heartBeatRegexp = regexp.MustCompile("^[0-9]{1,9},[0-9]{1,9}$")
)

// Maximum allowed value for heart-beat times in millisecond. 
// Approximately 11.5 days. Anything longer than this will be 
// rejected as an invalid frame.
const MaxHeartBeat = 999999999

// Represents a single STOMP frame.
type Frame struct {
	// The frame command. Should be one of the standard STOMP commands. 
	// Note that STOMP commands are case sensitive.
	Command string

	// Collection of frame headers.
	*Headers

	// The frame body. Only SEND, MESSAGE and ERROR frames may have a body.
	// All other frames must not have a body.
	Body []byte
}

// Creates a new frame with the specified command and headers. The headers
// should contain an even number of entries. Each even index is the header 
// name, and the odd indexes are the assocated header values.
func NewFrame(command string, headers ...string) *Frame {
	f := &Frame{Command: command, Headers: NewHeaders()}
	for index := 0; index < len(headers); index += 2 {
		f.Append(headers[index], headers[index+1])
	}
	return f
}

// Makes a copy of the frame. The command and headers are copied, but the
// body points to the same underlying array as the source frame.
func (f *Frame) Clone() *Frame {
	clone := &Frame{
		Command: f.Command,
		Headers: f.Headers.Clone(),
		Body:    f.Body,
	}
	return clone
}

// Returns the value of the "content-length" header, and whether it was
// found or not. Used for deserializing a frame. If the content length
// is specified in the header, then the body can contain null characters.
// Otherwise the body is read until a null character is encountered.
// If an error is returned, then the content-length header is malformed.
func (f *Frame) ContentLength() (contentLength int, ok bool, err error) {
	text, ok := f.Contains(ContentLength)
	if !ok {
		return
	}

	value, err := strconv.ParseUint(text, 10, 32)
	if err != nil {
		ok = false
		return
	}

	if value > MaxContentLength {
		err = exceededMaxFrameSize
		ok = false
		return
	}

	contentLength = int(value)
	return
}

// Determine the most acceptable version based on the accept-version
// header of a CONNECT or STOMP frame.
//
// Returns V1_0 if a CONNECT frame and the accept-version header is missing.
//
// Returns an error if the frame is not a CONNECT or STOMP frame, or
// if the accept-header is malformed or does not contain any compatible
// version numbers. Also returns an error if the accept-header is missing
// for a STOMP frame.
//
// Otherwise, returns the highest compatible version specified in the
// accept-version header. Compatible versions are V1_0, V1_1 and V1_2.
func (f *Frame) AcceptVersion() (version StompVersion, err error) {
	// frame can be CONNECT or STOMP with slightly different
	// handling of accept-verion for each
	isConnect := f.Command == CONNECT

	if !isConnect && f.Command != STOMP {
		err = notConnectFrame
		return
	}

	// start with an error, and remove if successful
	err = unknownVersion

	if acceptVersion, ok := f.Headers.Contains(AcceptVersion); ok {
		// sort the versions so that the latest version comes last
		versions := strings.Split(acceptVersion, ",")
		sort.Strings(versions)
		for _, v := range versions {
			switch StompVersion(v) {
			case V1_0:
				version = V1_0
				err = nil
			case V1_1:
				version = V1_1
				err = nil
			case V1_2:
				version = V1_2
				err = nil
			}
		}
	} else {
		// CONNECT frames can be missing the accept-version header,
		// we assume V1.0 in this case. STOMP frames were introduced
		// in V1.1, so they must have an accept-version header.
		if isConnect {
			// no "accept-version" header, so we assume 1.0
			version = V1_0
			err = nil
		} else {
			err = missingHeader(AcceptVersion)
		}
	}
	return
}

// Determine the heart-beat values in a CONNECT or STOMP frame.
//
// Returns 0,0 if the heart-beat header is missing. Otherwise
// returns the cx and cy values in the frame.
//
// Returns an error if the heart-beat header is malformed, or if
// the frame is not a CONNECT or STOMP frame. In this implementation,
// a heart-beat header is considered malformed if either cx or cy
// is greater than MaxHeartBeat.
func (f *Frame) HeartBeat() (cx, cy int, err error) {
	if f.Command != CONNECT && f.Command != STOMP && f.Command != CONNECTED {
		err = invalidOperationForFrame
		return
	}
	if heartBeat, ok := f.Headers.Contains(HeartBeat); ok {
		if !heartBeatRegexp.MatchString(heartBeat) {
			err = invalidHeartBeat
			return
		}

		// no error checking here because we are confident
		// that everything will work because the regexp matches.
		slice := strings.Split(heartBeat, ",")
		value1, _ := strconv.ParseUint(slice[0], 10, 32)
		value2, _ := strconv.ParseUint(slice[1], 10, 32)
		cx = int(value1)
		cy = int(value2)
	} else {
		// heart-beat header not present
		// this else clause is not necessary, but
		// included for clarity.
		cx = 0
		cy = 0
	}
	return
}

// Check frame for required headers. Returns
// nil if the frame is valid, non-nil otherwise.
//
// The frame is checked for a valid command (client or server),
// and each different command is verified for mandatory headers
// for that command. Also checks for prohibited headers. For 
// example, returns an error if the transaction header is present
// in a frame that does not support transactions (eg SUBSCRIBE,
// UNSUBSCRIBE, CONNECT)
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
			return missingHeader(name)
		}
	}
	return nil
}

func (f *Frame) verifyProhibitedHeaders(names ...string) error {
	for _, name := range names {
		if _, ok := f.Headers.Contains(name); ok {
			return prohibitedHeader(name)
		}
	}
	return nil
}

func (f *Frame) validateConnect() error {
	version, err := f.AcceptVersion()
	if err != nil {
		return err
	}
	if version == V1_0 {
		// no mandatory headers in V1.0
		return nil
	}

	// The STOMP specification mandates that this header must
	// be present for STOMP 1.1 and later. It is checked for
	// here, but the data is never used.
	err = f.verifyRequiredHeaders(Host)
	if err != nil {
		return err
	}

	if heartBeat, ok := f.Contains(HeartBeat); ok {
		if !heartBeatRegexp.MatchString(heartBeat) {
			return invalidHeartBeat
		}
	}

	err = f.verifyProhibitedHeaders(Destination, Transaction, Ack, Id, Receipt)
	if err != nil {
		return err
	}

	return nil
}

func (f *Frame) validateConnected() error {
	return nil
}

func (f *Frame) validateSend() error {
	return f.verifyRequiredHeaders(Destination)
}

func (f *Frame) validateSubscribe() error {
	err := f.verifyRequiredHeaders(Destination, Id)
	if err != nil {
		return err
	}

	if ack, ok := f.Contains(Ack); ok {
		switch ack {
		case AckAuto, AckClient, AckClientIndividual:
			// acceptable values, do nothing
		default:
			return invalidHeaderValue
		}
	}

	return f.verifyProhibitedHeaders(Transaction)
}

func (f *Frame) validateUnsubscribe() error {
	err := f.verifyRequiredHeaders(Id)
	if err != nil {
		return err
	}
	return f.verifyProhibitedHeaders(Transaction)
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
	return f.verifyProhibitedHeaders(Transaction)
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
