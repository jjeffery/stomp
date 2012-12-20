package topic

import (
	"github.com/jjeffery/stomp/message"
	"github.com/jjeffery/stomp/server/client"
)

// Topic for broadcasting to subscriptions.
type Topic struct {
	destination   string
	subs *client.SubscriptionList
}

// Create a new topic -- called from the topic manager only.
func newTopic(destination string) *Topic {
	return &Topic{
		destination: destination,
		subs: client.NewSubscriptionList(),
	}
}

// Add a subscription to a topic.
func (t *Topic) Subscribe(sub *client.Subscription) {
	t.subs.Add(sub)
}

// Unsubscribe a subscription.
func (t *Topic) Unsubscribe(sub *client.Subscription) {
	t.subs.Remove(sub)
}

// Send a message to the topic. All subscriptions receive a copy
// of the message.
func (t *Topic) Enqueue(f *message.Frame) {
	// find a subscription ready to receive the frame
	t.subs.ForEach(func(sub *client.Subscription, last bool) {
		if last {
			// can send without copying for the final subscription
			sub.SendTopicFrame(f)
		} else {
			// send a copy of the frame
			sub.SendTopicFrame(f.Clone())
		}
	})
}
