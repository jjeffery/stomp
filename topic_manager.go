package stomp

import (
	"github.com/jjeffery/stomp/message"
)

// TODO: topic manager is similar to queue manager, but is always memory based.
type topicManager struct {
}

func newTopicManager() *topicManager {
	tm := new(topicManager)
	return tm
}

func (tm *topicManager) handleConnect(c *conn) error {
	return notImplementedYet
}

func (tm *topicManager) handleDisconnect(c *conn) error {
	return notImplementedYet
}

func (tm *topicManager) handleSubscribe(conn *conn, frame *message.Frame) error {
	return notImplementedYet
}

func (tm *topicManager) handleUnsubscribe(conn *conn, frame *message.Frame) error {
	return notImplementedYet
}

func (tm *topicManager) handleSend(conn *conn, frame *message.Frame) error {
	return notImplementedYet
}

