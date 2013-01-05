package topic

import (
	"github.com/jjeffery/stomp"
)

// Subscription is the interface that wraps a subscriber to a topic.
type Subscription interface {
	// Send a message frame to the topic subscriber.
	SendTopicFrame(f *stomp.Frame)
}
