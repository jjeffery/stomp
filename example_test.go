package stomp_test

import (
	"github.com/jjeffery/stomp"
)

// Creates a new Header.
func ExampleNewHeader() {
	/*
		Creates a header that looks like the following:

			login:scott
			passcode:tiger
			host:stompserver
			accept-version:1.1,1.2
	*/
	h := stomp.NewHeader(
		"login", "scott",
		"passcode", "tiger",
		"host", "stompserver",
		"accept-version", "1.1,1.2")
	doSomethingWith(h)
}

// Creates a STOMP frame.
func ExampleNewFrame() {
	/*
		Creates a STOMP frame that looks like the following:

			CONNECT
			login:scott
			passcode:tiger
			host:stompserver
			accept-version:1.1,1.2
	*/
	f := stomp.NewFrame("CONNECT",
		"login", "scott",
		"passcode", "tiger",
		"host", "stompserver",
		"accept-version", "1.1,1.2")
	doSomethingWith(f)
}

func doSomethingWith(f interface{}) {

}

// Connect to a STOMP server using default options.
func ExampleDial_1() error {
	conn, err := stomp.Dial("tcp", "192.168.1.1:61613", stomp.Options{})
	if err != nil {
		return err
	}

	err = conn.Send(stomp.Message{
		Destination: "/queue/test-1",
		ContentType: "text/plain",
		Body:        []byte("Test message #1"),
	})
	if err != nil {
		return err
	}

	return conn.Disconnect()
}

// Connect to a STOMP server that requires authentication. In addition,
// we are only prepared to use STOMP protocol version 1.1 or 1.2, and
// the virtual host is named "dragon".
func ExampleDial_2() error {
	conn, err := stomp.Dial("tcp", "192.168.1.1:61613", stomp.Options{
		Login:         "scott",
		Passcode:      "leopard",
		AcceptVersion: "1.1,1.2",
		Host:          "dragon",
	})
	if err != nil {
		return err
	}

	err = conn.Send(stomp.Message{
		Destination: "/queue/test-1",
		ContentType: "text/plain",
		Body:        []byte("Test message #1"),
	})
	if err != nil {
		return err
	}

	return conn.Disconnect()
}
