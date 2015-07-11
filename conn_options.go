package stomp

import (
	"fmt"
	"github.com/jjeffery/stomp/frame"
	"strings"
	"time"
)

// ConnOptions is an opaque structure used to collection options
// for connecting to the other server.
type connOptions struct {
	FrameCommand    string
	Host            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	HeartBeatError  time.Duration
	Login, Passcode string
	AcceptVersions  []string
	Header          *Header
}

func newConnOptions(conn *Conn, opts []func(*Conn) error) (*connOptions, error) {
	co := &connOptions{
		FrameCommand:   frame.CONNECT,
		ReadTimeout:    time.Minute,
		WriteTimeout:   time.Minute,
		HeartBeatError: DefaultHeartBeatError,
	}

	// This is a slight of hand, attach the options to the Conn long
	// enough to run the options functions and then detach again.
	// The reason we do this is to allow for future options to be able
	// to modify the Conn object itself, in case that becomes desirable.
	conn.options = co
	defer func() { conn.options = nil }()

	for _, opt := range opts {
		if opt == nil {
			return nil, ErrNilOption
		}
		err := opt(conn)
		if err != nil {
			return nil, err
		}
	}

	if len(co.AcceptVersions) == 0 {
		co.AcceptVersions = append(co.AcceptVersions, string(V10), string(V11), string(V12))
	}

	return co, nil
}

func (co *connOptions) NewFrame() (*Frame, error) {
	f := NewFrame(co.FrameCommand)
	if co.Host != "" {
		f.Header.Set(frame.Host, co.Host)
	}

	// heart-beat
	{
		send := co.WriteTimeout / time.Millisecond
		recv := co.ReadTimeout / time.Millisecond
		f.Header.Set(frame.HeartBeat, fmt.Sprintf("%d,%d", send, recv))
	}

	// login, passcode
	if co.Login != "" || co.Passcode != "" {
		f.Header.Set(frame.Login, co.Login)
		f.Header.Set(frame.Passcode, co.Passcode)
	}

	// accept-version
	f.Header.Set(frame.AcceptVersion, strings.Join(co.AcceptVersions, ","))

	// custom header entries -- note that these do not override
	// header values already set as they are added to the end of
	// the header array
	f.Header.AddHeader(co.Header)

	return f, nil
}

// Options for connecting to the STOMP server. Used with the
// stomp.Dial and stomp.Connect functions, both of which have examples.
var ConnOpt struct {
	Login          func(login, passcode string) func(*Conn) error
	Host           func(host string) func(*Conn) error
	AcceptVersion  func(version Version) func(*Conn) error
	HeartBeat      func(sendTimeout, recvTimeout time.Duration) func(*Conn) error
	HeartBeatError func(errorTimeout time.Duration) func(*Conn) error
	Header         func(header *Header) func(*Conn) error
}

func init() {
	ConnOpt.Login = func(login, passcode string) func(*Conn) error {
		return func(c *Conn) error {
			c.options.Login = login
			c.options.Passcode = passcode
			return nil
		}
	}

	ConnOpt.Host = func(host string) func(*Conn) error {
		return func(c *Conn) error {
			c.options.Host = host
			return nil
		}
	}

	ConnOpt.AcceptVersion = func(version Version) func(*Conn) error {
		return func(c *Conn) error {
			if err := version.CheckSupported(); err != nil {
				return err
			}
			c.options.AcceptVersions = append(c.options.AcceptVersions, string(version))
			return nil
		}
	}

	ConnOpt.HeartBeat = func(sendTimeout, recvTimeout time.Duration) func(*Conn) error {
		return func(c *Conn) error {
			c.options.WriteTimeout = sendTimeout
			c.options.ReadTimeout = recvTimeout
			return nil
		}
	}

	ConnOpt.HeartBeatError = func(errorTimeout time.Duration) func(*Conn) error {
		return func(c *Conn) error {
			c.options.HeartBeatError = errorTimeout
			return nil
		}
	}

	ConnOpt.Header = func(header *Header) func(*Conn) error {
		return func(c *Conn) error {
			c.options.Header = header
			return nil
		}
	}
}
