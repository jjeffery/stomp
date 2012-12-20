package stomp

import (
	//"github.com/jjeffery/stomp/message"
	"net"
	"time"
)

const QueuePrefix = "/queue"

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
const DefaultHeartBeat = time.Minute

// Interface for authenticating STOMP clients.
type Authenticator interface {
	// Authenticate based on the given login and passcode, either of which might be nil.
	// Returns true if authentication is successful, false otherwise.
	Authenticate(login, passcode string) bool
}

// A Server defines parameters for running a STOMP server.
type Server struct {
	Addr           string        // TCP address to listen on, DefaultListenAddr if empty
	Authenticator  Authenticator // Authenticates login/passcodes. If nil no authentication is performed
	QueueStorage   QueueStorage  // Implementation of queue storage. If nil, in-memory queues are used.
	HeartBeat      time.Duration // Preferred value for heart-beat read/write timeout.
	MaxHeaderBytes int           // Maximum size of STOMP headers in bytes
	MaxBodyBytes   int           // Maximum size of STOMP body in bytes
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
	proc := newRequestProcessor(s)
	return proc.Serve(l)
}
