package main

import (
	"os"
)

func signals(signalChannel chan os.Signal) {
	// Windows has no other signals other than os.Interrupt
}
