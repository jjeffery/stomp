package main

import (
	"github.com/jjeffery/stomp"
)

func main() {
	_ = stomp.NewFrame("MESSAGE")
	println("hello")
}
