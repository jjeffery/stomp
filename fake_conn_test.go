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
	writer     io.WriteCloser
	reader     io.ReadCloser
	localAddr  net.Addr
	remoteAddr net.Addr
}

var (
	errClosing = errors.New("use of closed network connection")
)

func newFakeConn(c *C) (client *fakeConn, server *fakeConn) {
	clientReader, serverWriter := io.Pipe()
	serverReader, clientWriter := io.Pipe()
	const bufferSize = 10240

	clientConn := &fakeConn{
		C:      c,
		reader: clientReader,
		writer: clientWriter,
	}

	serverConn := &fakeConn{
		C:      c,
		reader: serverReader,
		writer: serverWriter,
	}

	return clientConn, serverConn
}

func (fc *fakeConn) Read(p []byte) (n int, err error) {
	return fc.reader.Read(p)
}

func (fc *fakeConn) Write(p []byte) (n int, err error) {
	return fc.writer.Write(p)
}

func (fc *fakeConn) Close() error {
	err1 := fc.reader.Close()
	err2 := fc.writer.Close()

	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
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
	panic("not implemented")
}

func (fc *fakeConn) SetReadDeadline(t time.Time) error {
	panic("not implemented")
}

func (fc *fakeConn) SetWriteDeadline(t time.Time) error {
	panic("not implemented")
}

type FakeConnSuite struct{}

var _ = Suite(&FakeConnSuite{})

func (s *FakeConnSuite) TestFakeConn(c *C) {
	//c.Skip("temporary")
	fc1, fc2 := newFakeConn(c)

	one := []byte{1, 2, 3, 4}
	two := []byte{5, 6, 7, 8, 9, 10, 11, 12, 13}
	stop := make(chan struct{})

	go func() {
		defer func() {
			fc2.Close()
			close(stop)
		}()

		rx1 := make([]byte, 6)
		n, err := fc2.Read(rx1)
		c.Assert(n, Equals, 4)
		c.Assert(err, IsNil)
		c.Assert(rx1[0:n], DeepEquals, one)

		rx2 := make([]byte, 5)
		n, err = fc2.Read(rx2)
		c.Assert(n, Equals, 5)
		c.Assert(err, IsNil)
		c.Assert(rx2, DeepEquals, []byte{5, 6, 7, 8, 9})

		rx3 := make([]byte, 10)
		n, err = fc2.Read(rx3)
		c.Assert(n, Equals, 4)
		c.Assert(err, IsNil)
		c.Assert(rx3[0:n], DeepEquals, []byte{10, 11, 12, 13})
	}()

	c.Assert(fc1.C, Equals, c)
	c.Assert(fc2.C, Equals, c)

	n, err := fc1.Write(one)
	c.Assert(n, Equals, 4)
	c.Assert(err, IsNil)

	n, err = fc1.Write(two)
	c.Assert(n, Equals, len(two))
	c.Assert(err, IsNil)

	<-stop
}
