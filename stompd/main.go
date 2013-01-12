/*
A simple, stand-alone STOMP server.
*/
package main

import (
	"github.com/jjeffery/stomp/server"
	"log"
	"net"
)

/*
func main() {
	// create a channel for listening for termination signals
	stopChannel := newStopChannel()

	for {
		select {
		case sig := <-stopChannel:
			log.Println("received signal:", sig)
			break
		}
	}

}
*/

func main() {
	addr := ":61613"
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %s", err.Error())
	}
	defer func() { l.Close() }()

	log.Println("listening on", l.Addr().Network(), l.Addr().String())
	server.Serve(l)

}
