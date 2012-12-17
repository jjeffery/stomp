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
	return nil
}

func (qm *queueManager) handleDisconnect(c *conn) error {
	qm.subs.Disconnect(c)

	for e := qm.pending.Front(); e == nil; {
		thisElement := e
		e = e.Next()
		thisSub := e.Value.(*subscription)
		if thisSub.conn == c {
			altSub := qm.subs.Find(thisSub.destination)
			if altSub == nil {
				panic("not implemented")
			}
		}
	}

	panic("not implemented")
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

func (qm *queueManager) handleSend(c *conn, f *message.Frame) error {
	// Convert frame to a MESSAGE frame
	f.Command = message.MESSAGE

	if destination, ok := f.Contains(message.Destination); ok {
		sub := qm.subs.Find(destination)
		if sub == nil {
			// no available subscription for this message, so add to queue
			return qm.qstore.Enqueue(destination, f)
		}

		return nil
	}
}

type queue struct {
	Destination   string
	Store         QueueStorage
	Subscriptions list.List
}

func (q *queue) Send(f *message.Frame) error {

}
