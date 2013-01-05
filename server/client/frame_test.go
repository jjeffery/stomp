package client

import (
	"github.com/jjeffery/stomp"
	"github.com/jjeffery/stomp/frame"
	. "launchpad.net/gocheck"
	_ "strconv"
)

type FrameSuite struct{}

var _ = Suite(&FrameSuite{})

func (s *FrameSuite) TestDetermineVersion_V10_Connect(c *C) {
	f := stomp.NewFrame(frame.CONNECT)
	version, err := determineVersion(f)
	c.Check(err, IsNil)
	c.Check(version, Equals, stomp.V10)
}

func (s *FrameSuite) TestDetermineVersion_V10_Stomp(c *C) {
	// the "STOMP" command was introduced in V1.1, so it must
	// have an accept-version header
	f := stomp.NewFrame(frame.STOMP)
	_, err := determineVersion(f)
	c.Check(err, Equals, missingHeader(frame.AcceptVersion))
}

func (s *FrameSuite) TestDetermineVersion_V11_Connect(c *C) {
	f := stomp.NewFrame(frame.CONNECT)
	f.Header.Add(frame.AcceptVersion, "1.1")
	version, err := determineVersion(f)
	c.Check(version, Equals, stomp.V11)
	c.Check(err, IsNil)
}

func (s *FrameSuite) TestDetermineVersion_MultipleVersions(c *C) {
	f := stomp.NewFrame(frame.CONNECT)
	f.Header.Add(frame.AcceptVersion, "1.2,1.1,1.0,2.0")
	version, err := determineVersion(f)
	c.Check(version, Equals, stomp.V12)
	c.Check(err, IsNil)
}

func (s *FrameSuite) TestDetermineVersion_IncompatibleVersions(c *C) {
	f := stomp.NewFrame(frame.CONNECT)
	f.Header.Add(frame.AcceptVersion, "0.2,0.1,1.3,2.0")
	version, err := determineVersion(f)
	c.Check(version, Equals, stomp.Version(""))
	c.Check(err, Equals, unknownVersion)
}
