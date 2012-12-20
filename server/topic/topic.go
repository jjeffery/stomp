/*
Package topic provides implementations of server-side topics.
*/
package topic

import (
	"github.com/jjeffery/stomp/message"
	"github.com/jjeffery/stomp/server/client"
)

// A Topic is used for broadcasting to subscribed clients.
// In contrast to a queue, when a message is sent to a topic,
// that message is transmitted to all subscribed clients.
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

// Subscribe adds a subscription to a topic. Any message sent to the
// topic will be transmitted to the subscription's client until 
// unsubscription occurs.
func (t *Topic) Subscribe(sub *client.Subscription) {
	t.subs.Add(sub)
}

// Unsubscribe causes a subscription to be removed from the topic.
func (t *Topic) Unsubscribe(sub *client.Subscription) {
	t.subs.Remove(sub)
}

// Enqueue send a message to the topic. All subscriptions receive a copy
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
