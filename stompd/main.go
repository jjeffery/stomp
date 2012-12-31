/*
A simple, stand-alone STOMP server.
*/
package main

import (
	"log"
)

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
