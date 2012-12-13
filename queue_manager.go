package stomp

import (
	"container/list"
	"github.com/jjeffery/stomp/message"
	"time"
)

type queueManager struct {
	qstore  QueueStorage // handles queue storage
	subs    subscriptionManager
	pending *list.List // list of pendingMsg
}

// Information about a message pending
type pendingMsg struct {
	subscription *subscription
	frame        *message.Frame
	sent         time.Time
}

// Create a queue manager with the specified queue storage mechanism
func newQueueManager(qstore QueueStorage) *queueManager {
	qm := new(queueManager)
	qm.qstore = qstore
	qm.pending = list.New()
	return qm
}

func (qm *queueManager) handleConnect(c *conn) error {
	return notImplementedYet
}

func (qm *queueManager) handleDisconnect(c *conn) error {
	return notImplementedYet
}

func (qm *queueManager) handleSubscribe(conn *conn, frame *message.Frame) error {
	return notImplementedYet
}

func (qm *queueManager) handleUnsubscribe(conn *conn, frame *message.Frame) error {
	return notImplementedYet
}

func (qm *queueManager) handleAck(conn *conn, frame *message.Frame) error {
	return notImplementedYet
}

func (qm *queueManager) handleNack(conn *conn, frame *message.Frame) error {
	return notImplementedYet
}

func (qm *queueManager) handleSend(conn *conn, frame *message.Frame) error {
	return notImplementedYet
}
