package queue

import (
	"github.com/jjeffery/stomp/client"
	"github.com/jjeffery/stomp/message"
)

// Subscribes a client connection to a queue
func Subscribe(destination string, conn *client.Connection) {
	panic("not implemented")
}

// Unsubscribe a client connection from a queue
func Unsubscribe(destination string, conn *client.Connection) {
	panic("not implemented")
}

// Indicates a client has disconnected. Any unacknowledged
// frames will be re-queued.
func Disconnected(conn *client.Connection) {
	panic("not implemented")
}

func Ack(conn *client.Connection, frame *message.Frame) {
	panic("not implemented")
}

func Nack(conn *client.Connection, frame *message.Frame) {

}

func Stop() {
	// TODO
}

type Manager struct {
	qstore Storage // handles queue storage
}

// Create a queue manager with the specified queue storage mechanism
func NewQueueManager(qstore QueueStorage) *QueueManager {
	qm := new(QueueManager)
	qm.qstore = qstore
	return qm
}
