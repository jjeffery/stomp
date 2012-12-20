package message

import (
	. "launchpad.net/gocheck"
	"strconv"
)

type FrameSuite struct{}

var _ = Suite(&FrameSuite{})

func (s *FrameSuite) TestAcceptVersion_V10_Connect(c *C) {
	f := NewFrame(CONNECT)
	version, err := f.AcceptVersion()
	c.Check(err, IsNil)
	c.Check(version, Equals, V1_0)
}

func (s *FrameSuite) TestAcceptVersion_V10_Stomp(c *C) {
	// the "STOMP" command was introduced in V1.1, so it must
	// have an accept-version header
	f := NewFrame(STOMP)
	_, err := f.AcceptVersion()
	c.Check(err, Equals, missingHeader(AcceptVersion))
}

func (s *FrameSuite) TestAcceptVersion_V11_Connect(c *C) {
	f := NewFrame(CONNECT)
	f.Headers.Append(AcceptVersion, "1.1")
	version, err := f.AcceptVersion()
	c.Check(version, Equals, V1_1)
	c.Check(err, IsNil)
}

func (s *FrameSuite) TestAcceptVersion_MultipleVersions(c *C) {
	f := NewFrame(CONNECT)
	f.Headers.Append(AcceptVersion, "1.2,1.1,1.0,2.0")
	version, err := f.AcceptVersion()
	c.Check(version, Equals, V1_2)
	c.Check(err, IsNil)
}

func (s *FrameSuite) TestAcceptVersion_IncompatibleVersions(c *C) {
	f := NewFrame(CONNECT)
	f.Headers.Append(AcceptVersion, "0.2,0.1,1.3,2.0")
	version, err := f.AcceptVersion()
	c.Check(version, Equals, StompVersion(""))
	c.Check(err, Equals, unknownVersion)
}

func (s *FrameSuite) TestValidate_Connect(c *C) {
	f := NewFrame(CONNECT)

	// CONNECT without accept-version can be missing host header
	err := f.Validate()
	c.Check(err, IsNil)

	// Once the CONNECT states it is V1.1 compatible, it must have
	// a host header
	f.Headers.Append(AcceptVersion, "1.1")
	err = f.Validate()
	c.Check(err, Equals, missingHeader(Host))

	f.Headers.Append(Host, "")
	err = f.Validate()
	c.Check(err, IsNil)

	f.Headers.Append(HeartBeat, "0,0")
	err = f.Validate()
	c.Check(err, IsNil)

	f.Headers.Set(HeartBeat, "garbage")
	err = f.Validate()
	c.Check(err, Equals, invalidHeartBeat)

	f.Headers.Set(HeartBeat, "60000,120000")
	err = f.Validate()
	c.Check(err, IsNil)

	// only allow 9 digits per value
	f.Headers.Set(HeartBeat, "9999999999,999")
	err = f.Validate()
	c.Check(err, Equals, invalidHeartBeat)
}

func (s *FrameSuite) TestValidate_Stomp(c *C) {
	f := NewFrame(STOMP)

	// STOMP must have an accept-version header
	err := f.Validate()
	c.Check(err, Equals, missingHeader(AcceptVersion))

	// STOMP must have a host header
	f.Headers.Append(AcceptVersion, "1.1")
	err = f.Validate()
	c.Check(err, Equals, missingHeader(Host))

	f.Headers.Append(Host, "anything")
	err = f.Validate()
	c.Check(err, IsNil)
}

func (s *FrameSuite) TestHeartBeat(c *C) {
	f := NewFrame(CONNECT,
		AcceptVersion, "1.2",
		Host, "XX")

	// no heart-beat header means zero values
	x, y, err := f.HeartBeat()
	c.Check(x, Equals, 0)
	c.Check(y, Equals, 0)
	c.Check(err, IsNil)

	f.Headers.Append("heart-beat", "123,456")
	x, y, err = f.HeartBeat()
	c.Check(x, Equals, 123)
	c.Check(y, Equals, 456)
	c.Check(err, IsNil)

	f.Headers.Set(HeartBeat, "invalid")
	x, y, err = f.HeartBeat()
	c.Check(x, Equals, 0)
	c.Check(y, Equals, 0)
	c.Check(err, Equals, invalidHeartBeat)

	f.Headers.Remove(HeartBeat)
	_, _, err = f.HeartBeat()
	c.Check(err, IsNil)

	f.Command = SEND
	_, _, err = f.HeartBeat()
	c.Check(err, Equals, invalidOperationForFrame)
}

func (s *FrameSuite) TestContentLength(c *C) {
	f := NewFrame(SEND,
		Destination, "/queue/test",
		Transaction, "tx1")

	contentLength, ok, err := f.ContentLength()
	c.Check(contentLength, Equals, 0)
	c.Check(ok, Equals, false)
	c.Check(err, IsNil)

	f.Append(ContentLength, "12")
	contentLength, ok, err = f.ContentLength()
	c.Check(contentLength, Equals, 12)
	c.Check(ok, Equals, true)
	c.Check(err, IsNil)

	f.Set(ContentLength, "-12")
	contentLength, ok, err = f.ContentLength()
	c.Check(contentLength, Equals, 0)
	c.Check(ok, Equals, false)
	c.Check(err, NotNil)

	f.Set(ContentLength, strconv.Itoa(16*1024*1024+1))
	contentLength, ok, err = f.ContentLength()
	c.Check(contentLength, Equals, 0)
	c.Check(ok, Equals, false)
	c.Check(err, Equals, exceededMaxFrameSize)
}
