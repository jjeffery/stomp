package main

import (
	"github.com/jjeffery/stomp/client"
	"log"
)

func main() {
	// create a channel for listening for termination signals
	stopChannel := newStopChannel()

	err := client.StartListening()
	if err != nil {
		log.Fatal(err)
	}

main_loop:
	for {
		select {
		case sig := <-stopChannel:
			log.Println("received signal:", sig)
			break main_loop
		case request := <- client.Requests
			handleRequest(request)
		}
	}

	client.StopListening()
}

func handleRequest(request client.Request) {
	// TODO
}
