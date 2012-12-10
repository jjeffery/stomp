package stomp

import (
	//"github.com/jjeffery/stomp/message"
	"time"
)

const (
	// default address for listening for connections
	DefaultListenAddr = ":61613"
)

type Authenticator interface {
	Authenticate(login, passcode string) bool
}

// Contains configurable parameters that modify the behaviour
// of the STOMP server.
type ServerConfig struct {
	// Authenticates login/passcode pairs. If nil, no authentication is performed.
	Authenticator Authenticator

	// Preferred value for heart-beat read timeout. Zero indicates no read heart-beat.
	HeartBeatReadTimeout time.Duration

	// Preferred value for heart-beat write timeout. Zero indicates no write heart-beat.
	HeartBeatWriteTimeout time.Duration

	// Maximum size of stomp headers
	MaxHeaderBytes int

	// Maximum size of stomp body
	MaxBodyBytes int
}

func NewServerConfig() {
	c := new(ServerConfig)
	c.HeartBeatReadTimeout = time.Duration(time.Minute)
	c.HeartBeatWriteTimeout = time.Duration(time.Minute)
	c.MaxHeaderBytes = 4096
	c.MaxBodyBytes = 1024 * 1024
}

type Server struct {
	Addr   string        // TCP address to listen on, DefaultListenAddr if empty
	Config *ServerConfig // Configuration
}
