package client

import (
	. "launchpad.net/gocheck"
)

type SubscriptionListSuite struct{}

var _ = Suite(&SubscriptionListSuite{})

func (s *SubscriptionListSuite) TestAddAndGet(c *C) {
	sub1 := newSubscription(nil, "/dest", "1", "client")
	sub2 := newSubscription(nil, "/dest", "2", "client")
	sub3 := newSubscription(nil, "/dest", "3", "client")

	sl := NewSubscriptionList()
	sl.Add(sub1)
	sl.Add(sub2)
	sl.Add(sub3)

	c.Check(sl.Get(), Equals, sub1)

	// add the subscription again, should go to the back
	sl.Add(sub1)

	c.Check(sl.Get(), Equals, sub2)
	c.Check(sl.Get(), Equals, sub3)
	c.Check(sl.Get(), Equals, sub1)

	c.Check(sl.Get(), IsNil)
}

func (s *SubscriptionListSuite) TestAddAndRemove(c *C) {
	sub1 := newSubscription(nil, "/dest", "1", "client")
	sub2 := newSubscription(nil, "/dest", "2", "client")
	sub3 := newSubscription(nil, "/dest", "3", "client")

	sl := NewSubscriptionList()
	sl.Add(sub1)
	sl.Add(sub2)
	sl.Add(sub3)

	c.Check(sl.subs.Len(), Equals, 3)

	// now remove the second subscription
	sl.Remove(sub2)

	c.Check(sl.Get(), Equals, sub1)
	c.Check(sl.Get(), Equals, sub3)
	c.Check(sl.Get(), IsNil)
}
