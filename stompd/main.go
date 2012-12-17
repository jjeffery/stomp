package main

import (
	"github.com/jjeffery/stomp"
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
		case request := <-client.Requests:
			handleRequest(request)
		}
	}

	client.StopListening()
}

func handleRequest(request client.Request) {
	switch request.Type {
	case client.Disconnect:
		queue.Disconnected(client.Connection)
	case client.Create:
		// do nothing at the moment
	case client.Subscribe:
		handleSubscribeRequest(request)
	}
}

func handleSubscribeRequest(request client.Request) {
	// frame has already been checked, so we know it has a destination and an id
	destination, _ := request.Frame.Headers.Contains(message.Destination)
	id, _ := request.Frame.Headers.Contains(message.Id)
	
	// type of acknowledgement required
	ack, hasAck := 
}
