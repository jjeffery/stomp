package topic

import (
	"github.com/jjeffery/stomp/message"
	. "launchpad.net/gocheck"
)

type TopicSuite struct{}

var _ = Suite(&TopicSuite{})

func (s *TopicSuite) TestTopicWithoutSubscription(c *C) {
	topic := newTopic("destination")

	f := message.NewFrame(message.MESSAGE,
		message.Destination, "destination")

	topic.Enqueue(f)
}

func (s *TopicSuite) TestTopicWithOneSubscription(c *C) {
	sub := &fakeSubscription{}

	topic := newTopic("destination")
	topic.Subscribe(sub)

	f := message.NewFrame(message.MESSAGE,
		message.Destination, "destination")

	topic.Enqueue(f)

	c.Assert(len(sub.Frames), Equals, 1)
	c.Assert(sub.Frames[0], Equals, f)
}

func (s *TopicSuite) TestTopicWithTwoSubscriptions(c *C) {
	sub1 := &fakeSubscription{}
	sub2 := &fakeSubscription{}

	topic := newTopic("destination")
	topic.Subscribe(sub1)
	topic.Subscribe(sub2)

	f := message.NewFrame(message.MESSAGE,
		message.Destination, "destination",
		"xxx", "yyy")

	topic.Enqueue(f)

	c.Assert(len(sub1.Frames), Equals, 1)
	c.Assert(len(sub2.Frames), Equals, 1)
	c.Assert(sub1.Frames[0], Not(Equals), f)
	c.Assert(sub2.Frames[0], Equals, f)
}

type fakeSubscription struct {
	// frames received by the subscription
	Frames []*message.Frame
}

func (s *fakeSubscription) SendTopicFrame(f *message.Frame) {
	s.Frames = append(s.Frames, f)
}
