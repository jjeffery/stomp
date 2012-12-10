package stomp

import (
//"github.com/jjeffery/stomp/message"
)

type QueueManager struct {
	qstore QueueStorage // handles queue storage
}

// Create a queue manager with the specified queue storage mechanism
func NewQueueManager(qstore QueueStorage) *QueueManager {
	qm := new(QueueManager)
	qm.qstore = qstore
	return qm
}
