package main

import (
	"os"
	"os/signal"
	"syscall"
)

func setupStopSignals(signalChannel chan os.Signal) {
	// TODO: provide unix-specific option for daemonizing the 
	signal.Notify(signalChannel, syscall.SIGHUP)
	signal.Notify(signalChannel, syscall.SIGTERM)
}
