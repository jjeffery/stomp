package stomp

import (
	//"github.com/jjeffery/stomp/message"
	"container/list"
)

type queueManager struct {
	qstore QueueStorage // handles queue storage

}

type subscriber struct {
	id          int
	destination string
	conn        *conn
	ack         AckMode
}

type subscriberStore struct {
	destinations map[string]*list.List // maps queue name to a list of subscribers
	conns        map[*conn]*list.List  // maps connection to a list of subscribers
}

func (s *subscriberStore) Subscribe(conn *conn, destination string, ack AckMode) {
	panic("not implemented")
}

func (s *subscriberStore) Unsubscribe(conn *conn, id int) {
	panic("not implemented")
}

// Create a queue manager with the specified queue storage mechanism
func newQueueManager(qstore QueueStorage) *queueManager {
	qm := new(queueManager)
	qm.qstore = qstore
	return qm
}
