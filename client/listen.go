package client

import (
	"log"
	"net"
)

// Channel for receiving client requests
var Requests = make(chan Request, 128)

var listener net.Listener
var stopped = false

// Start listening for client connections.
func StartListening() error {
	stopped = false
	l, err := net.Listen("tcp", ":61613")
	if err != nil {
		return err
	}
	listener = l
	go listen()
	log.Println("listening on ", l.Addr())
	return nil
}

// Stop listening for client connections.
func StopListening() {
	stopped = true
	listener.Close()
}

func listen() {
	for {
		conn, err := listener.Accept()
		if stopped {
			// request to stop
			return
		}

		if err != nil {
			// TODO: is there better error handling here
			// than exiting
			log.Fatal(err)
		}

		log.Println("accepted connection from", conn.RemoteAddr())

		newConnection(conn, Requests)
	}
}
