package stomp

import (
	"errors"
	"io"
	. "launchpad.net/gocheck"
	"net"
	"time"
)

// fakeConn is a fake connection used for testing.
type fakeConn struct {
	C          *C
	Closed     bool
	ch         chan []byte
	buf        []byte
	remoteAddr net.Addr
	localAddr  net.Addr
}

var (
	errClosing = errors.New("use of closed network connection")
)

func newFakeConn(c *C) *fakeConn {
	return &fakeConn{
		C:  c,
		ch: make(chan []byte, 16),
	}
}

func (fc *fakeConn) Read(p []byte) (n int, err error) {
	if fc.Closed {
		err = errClosing
		return
	}

	if len(fc.buf) == 0 {
		var ok bool
		fc.buf, ok = <-fc.ch
		if !ok {
			err = io.EOF
			return
		}
	}

	if len(fc.buf) <= len(p) {
		copy(p, fc.buf)
		n = len(fc.buf)
		fc.buf = nil
	} else {
		copy(p, fc.buf)
		n = len(p)
		fc.buf = fc.buf[n:]
	}

	return
}

func (fc *fakeConn) Write(p []byte) (n int, err error) {
	if fc.Closed {
		err = errClosing
		return
	}

	pcopy := make([]byte, len(p))
	copy(pcopy, p)
	fc.ch <- pcopy
	n = len(p)

	return
}

func (fc *fakeConn) Close() error {
	if fc.Closed {
		return errClosing
	}
	fc.Closed = true
	close(fc.ch)
	return nil
}

func (fc *fakeConn) LocalAddr() net.Addr {
	return fc.localAddr
}

func (fc *fakeConn) RemoteAddr() net.Addr {
	return fc.remoteAddr
}

func (fc *fakeConn) SetLocalAddr(addr net.Addr) {
	fc.localAddr = addr
}

func (fc *fakeConn) SetRemoteAddr(addr net.Addr) {
	fc.remoteAddr = addr
}

func (fc *fakeConn) SetDeadline(t time.Time) error {
	fc.C.Assert(fc.Closed, Equals, false)
	panic("not implemented")
}

func (fc *fakeConn) SetReadDeadline(t time.Time) error {
	fc.C.Assert(fc.Closed, Equals, false)
	panic("not implemented")
}

func (fc *fakeConn) SetWriteDeadline(t time.Time) error {
	fc.C.Assert(fc.Closed, Equals, false)
	panic("not implemented")
}

type FakeConnSuite struct{}

var _ = Suite(&FakeConnSuite{})

func (s *FakeConnSuite) TestFakeConn(c *C) {
	fc := newFakeConn(c)

	one := []byte{1, 2, 3, 4}

	n, err := fc.Write(one)
	c.Assert(n, Equals, 4)
	c.Assert(err, IsNil)

	two := []byte{5, 6, 7, 8, 9, 10, 11, 12, 13}

	n, err = fc.Write(two)
	c.Assert(n, Equals, len(two))
	c.Assert(err, IsNil)

	rx1 := make([]byte, 256)
	n, err = fc.Read(rx1)
	c.Assert(n, Equals, 4)
	c.Assert(err, IsNil)
	c.Assert(rx1[0:n], DeepEquals, one)

	rx2 := make([]byte, 5)
	n, err = fc.Read(rx2)
	c.Assert(n, Equals, 5)
	c.Assert(err, IsNil)
	c.Assert(rx2, DeepEquals, []byte{5, 6, 7, 8, 9})

	rx3 := make([]byte, 10)
	n, err = fc.Read(rx3)
	c.Assert(n, Equals, 4)
	c.Assert(err, IsNil)
	c.Assert(rx3[0:n], DeepEquals, []byte{10, 22, 12, 13})
}
