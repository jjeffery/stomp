package stomp

import (
	"github.com/jjeffery/stomp/message"
	"log"
	"net"
	"strings"
	"time"
)

type requestProcessor struct {
	server *Server
	ch     chan request
	qm     *queueManager
	tm     *topicManager
	stop   bool  // has stop been requested
	msgId1 int64 // for generating message-id headers
	msgId2 int64 // for generating message-id headers
}

func newRequestProcessor(server *Server) *requestProcessor {
	proc := &requestProcessor{
		server: server,
		ch:     make(chan request, 32),
		tm:     newTopicManager(),
		msgId1: time.Now().Unix(),
	}

	if server.QueueStorage == nil {
		proc.qm = newQueueManager(NewMemoryQueueStorage())
	} else {
		proc.qm = newQueueManager(server.QueueStorage)
	}

	return proc
}

func (proc *requestProcessor) Serve(l net.Listener) error {
	go proc.Listen(l)

	for {
		r := <-proc.ch
		switch r.op {
		case stopOp:
			l.Close() // stop listening
			// TODO: would be good to gracefully shutdown connections
			return nil
		case connectOp:
			proc.handleConnect(r.conn)
		case disconnectOp:
			proc.handleDisconnect(r.conn)
		case frameOp:
			proc.handleFrame(r.conn, r.frame)
		default:
			panic("unknown request: " + r.op.String())
		}
	}
	panic("not reached")
}

func (proc *requestProcessor) handleConnect(c *conn) {
	proc.qm.handleConnect(c)
	proc.tm.handleConnect(c)
}

func (proc *requestProcessor) handleDisconnect(c *conn) {
	proc.qm.handleDisconnect(c)
	proc.tm.handleDisconnect(c)
}

func (proc *requestProcessor) handleFrame(c *conn, f *message.Frame) {
	switch f.Command {
	case message.UNSUBSCRIBE:
		// We cannot easily tell whether this is for 
		// a topic or queue, so just send to both.
		proc.qm.handleUnsubscribe(c, f)
		proc.tm.handleUnsubscribe(c, f)
	case message.SUBSCRIBE:
		if isQueueFrame(f) {
			proc.qm.handleSubscribe(c, f)
		} else {
			proc.tm.handleSubscribe(c, f)
		}
	case message.ACK:
		// only queues require ACK
		proc.qm.handleAck(c, f)
	case message.NACK:
		// only queues require NACK
		proc.qm.handleNack(c, f)
	case message.SEND:
		// convert to a MESSAGE frame
		f.Command = message.MESSAGE
		if isQueueFrame(f) {
			proc.qm.handleSend(c, f)
		} else {
			proc.tm.handleSend(c, f)
		}
	default:
		log.Println("unhandled command:", f.Command)
	}
}

func (proc *requestProcessor) generateMessageId() string {

}

func isQueueFrame(f *message.Frame) bool {
	if destination, ok := f.Contains(message.Destination); ok {
		// If the frame has a destination header, then it applies
		// to a queue if it starts with the queue prefix. Otherwise
		// the destination is considered a topic.
		return strings.HasPrefix(destination, QueuePrefix)
	}
	return false
}

func (proc *requestProcessor) Listen(l net.Listener) {
	timeout := time.Duration(0) // how long to sleep on accept failure
	for {
		rw, err := l.Accept()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				if timeout == 0 {
					timeout = 5 * time.Millisecond
				} else {
					timeout *= 2
				}
				if max := 5 * time.Second; timeout > max {
					timeout = max
				}
				log.Printf("stomp: Accept error: %v; retrying in %v", err, timeout)
				time.Sleep(timeout)
				continue
			}
			return
		}
		timeout = 0
		// TODO: need to pass Server to connection so it has access to
		// configuration parameters.
		_ = newConn(proc.server, rw, proc.ch)
	}
	panic("not reached")
}
