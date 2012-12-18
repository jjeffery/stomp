package client

import (
	. "launchpad.net/gocheck"
)

// Test suite for testing that channels work the way I expect.
type ChannelSuite struct {} 

var _ = Suite(&ChannelSuite{})

func (s *ChannelSuite) TestChannelsWorkAsExpected(c *C) {

	ch := make(chan int, 10)
	
	ch <- 1
	ch <- 2
	
	select {
		case i, ok := <-ch:
			c.Assert(i, Equals, 1)
			c.Assert(ok, Equals, true)
		default:
			c.Error("expected value on channel")
	}
	
	select {
		case i := <-ch:
			c.Assert(i, Equals, 2)
		default:
			c.Error("expected value on channel")
	}
	
	select {
		case _ = <-ch:
			c.Error("not expecting anything on the channel")
		default:
	}
	
	ch <- 3
	close(ch)
	
	select {
		case i := <-ch:
			c.Assert(i, Equals, 3)
		default:
			c.Error("expected value on channel")
	}
	
	select {
		case _, ok := <-ch:
			c.Assert(ok, Equals, false)
		default:
			c.Error("expected value on channel")
	}

	
	select {
		case _, ok := <-ch:
			c.Assert(ok, Equals, false)
		default:
			c.Error("expected value on channel")
	}
}