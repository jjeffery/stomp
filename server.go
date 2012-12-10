package stomp

import (
	//"github.com/jjeffery/stomp/message"
	"log"
	"net"
	"time"
)

// Default address for listening for connections.
const DefaultListenAddr = ":61613"

// Default maximum bytes for the header part of STOMP frames.
// Override by setting Server.MaxHeaderBytes.
const DefaultMaxHeaderBytes = 1 << 12 // 4KB

// Default maximum bytes for the body part of STOMP frames.
// Override by setting Server.MaxBodyBytes
const DefaultMaxBodyBytes = 1 << 20 // 1MB

// Default read timeout for heart-beat.
// Override by setting Server.HeartBeatReadTimeout.
const DefaultHeartBeatReadTimeout = time.Duration(30) * time.Duration(time.Second)

// Default write timeout for heart-beat.
// Override by setting Server.HeartBeatWriteTimeout.
const DefaultHeartBeatWriteTimeout = time.Duration(30) * time.Duration(time.Second)

// Interface for authenticating STOMP clients.
type Authenticator interface {
	Authenticate(login, passcode string) bool
}

// A Server defines parameters for running a STOMP server.
type Server struct {
	Addr                  string        // TCP address to listen on, DefaultListenAddr if empty
	Authenticator         Authenticator // Authenticates login/passcodes. If nil no authentication is performed
	QueueStorage          QueueStorage  // Implementation of queue storage. If nil, in-memory queues are used.
	HeartBeatReadTimeout  time.Duration // Preferred value for heart-beat read timeout. Zero indicates no read heart-beat.
	HeartBeatWriteTimeout time.Duration // Preferred value for heart-beat write timeout. Zero indicates no write heart-beat.
	MaxHeaderBytes        int           // Maximum size of STOMP headers in bytes
	MaxBodyBytes          int           // Maximum size of STOMP body in bytes
}

func ListenAndServe(addr string) error {
	s := &Server{Addr: addr}
	return s.ListenAndServe()
}

func Serve(l net.Listener) error {
	s := &Server{}
	return s.Serve(l)
}

// Listens on the TCP network address s.Addr and then calls
// Serve to handle requests on the incoming connections. If
// s.Addr is blank, then DefaultListenAddr is used.
func (s *Server) ListenAndServe() error {
	addr := s.Addr
	if addr == "" {
		addr = DefaultListenAddr
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return s.Serve(l)
}

// Accepts incoming connections on the Listener l, creating a new
// service thread for each connection. The service threads read
// requests and then process each request.
func (s *Server) Serve(l net.Listener) error {
	defer l.Close()
	ch := make(chan request, 32)
	go s.processRequests(ch, s.QueueStorage)
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
			return err
		}
		timeout = 0
		_ = newConn(rw, ch)
	}
	panic("not reached")
}

func (s *Server) processRequests(ch chan request, queueStorage QueueStorage) {
	if queueStorage == nil {
		// TODO allocate in-memory storage
	}

	for {
		r := <-ch
		switch r.op {
		case stopOp:
			return
		case disconnectOp:
			// TODO
		case frameOp:
			// TODO	
		}
	}
}
