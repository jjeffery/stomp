package server

import (
	"github.com/jjeffery/stomp"
)

// input channel for receiving requests
var inputChannel chan Request

var processes map[*Connection]*Process

// state functions
type requestHandler func(Request)

// not a good name -- it is the connection state?
type Process struct {
	handleRequest requestHandler
}

func Run() {

	for {
		select {
		case request := <-inputChannel:
			handleRequest(request)
		}
	}
}

func handleRequest(r Request) {
	process := processes[r.Connection]
	if process == nil {
		process = new(Process)
		process.handleRequest = waitingForConnect
		processes[r.Connection] = process
	}
	process.handleRequest(r)
	
	// if an error was received, remove the process
	if r.Error != nil {
		delete(processes, r.Connection)
	}
}

func waitingForConnect(r Request) {
	if r.Error != nil {
		// no cleanup required, as nothing happened yet
		return
	}
	
	if frame == nil || frame.Command != stomp.Connect {
		
	}
}
