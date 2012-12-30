package stomp

import (
	. "launchpad.net/gocheck"
	"runtime"
)

type IdSuite struct{}

var _ = Suite(&IdSuite{})

// only used during testing, does not need to be thread-safe
func resetId() {
	_lastId = 0
}

func (s *IdSuite) SetUpSuite(c *C) {
	resetId()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func (s *IdSuite) TearDownSuite(c *C) {
	runtime.GOMAXPROCS(1)
}

func (s *IdSuite) TestAllocateId(c *C) {
	c.Assert(allocateId(), Equals, "1")
	c.Assert(allocateId(), Equals, "2")

	ch := make(chan bool, 50)
	for i := 0; i < 50; i++ {
		go doAllocate(100, ch)
	}

	for i := 0; i < 50; i++ {
		<-ch
	}

	c.Assert(allocateId(), Equals, "5003")
}

func doAllocate(count int, ch chan bool) {
	for i := 0; i < count; i++ {
		_ = allocateId()
	}
	ch <- true
}
