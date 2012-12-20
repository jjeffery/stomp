package queue

import (
	"github.com/jjeffery/stomp/message"
	. "launchpad.net/gocheck"
)

type MemoryQueueSuite struct{}

var _ = Suite(&MemoryQueueSuite{})

func (s *MemoryQueueSuite) Test1(c *C) {
	mq := NewMemoryQueueStorage()
	mq.Start()

	f1 := message.NewFrame(message.MESSAGE,
		message.Destination, "/queue/test",
		message.MessageId, "msg-001",
		message.Subscription, "1")

	err := mq.Enqueue("/queue/test", f1)
	c.Assert(err, IsNil)

	f2 := message.NewFrame(message.MESSAGE,
		message.Destination, "/queue/test",
		message.MessageId, "msg-002",
		message.Subscription, "1")

	err = mq.Enqueue("/queue/test", f2)
	c.Assert(err, IsNil)

	f3 := message.NewFrame(message.MESSAGE,
		message.Destination, "/queue/test2",
		message.MessageId, "msg-003",
		message.Subscription, "2")

	err = mq.Enqueue("/queue/test2", f3)
	c.Assert(err, IsNil)

	// attempt to dequeue from a different queue
	f, err := mq.Dequeue("/queue/other-queue")
	c.Check(err, IsNil)
	c.Assert(f, IsNil)

	f, err = mq.Dequeue("/queue/test2")
	c.Check(err, IsNil)
	c.Assert(f, Equals, f3)

	f, err = mq.Dequeue("/queue/test")
	c.Check(err, IsNil)
	c.Assert(f, Equals, f1)

	f, err = mq.Dequeue("/queue/test")
	c.Check(err, IsNil)
	c.Assert(f, Equals, f2)

	f, err = mq.Dequeue("/queue/test")
	c.Check(err, IsNil)
	c.Assert(f, IsNil)

	f, err = mq.Dequeue("/queue/test2")
	c.Check(err, IsNil)
	c.Assert(f, IsNil)
}
