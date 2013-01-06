package stomp_test

import (
	"fmt"
	"github.com/jjeffery/stomp"
	"net"
	"time"
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

			^@
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

func doAnotherThingWith(f interface{}, g interface{})

func ExampleSubscription() error {
	conn, err := stomp.Dial("tcp", "localhost:61613", stomp.Options{})
	if err != nil {
		return err
	}

	sub, err := conn.Subscribe("/queue/test-2", stomp.AckClient)
	if err != nil {
		return err
	}

	// receive 5 messages and then quit
	for i := 0; i < 5; i++ {
		msg := <-sub.C
		doSomethingWith(msg)

		// acknowledge the message
		err = conn.Ack(msg)
		if err != nil {
			return err
		}
	}

	err = sub.Unsubscribe()
	if err != nil {
		return err
	}

	return conn.Disconnect()
}

func ExampleTransaction() error {
	conn, err := stomp.Dial("tcp", "localhost:61613", stomp.Options{})
	if err != nil {
		return err
	}
	defer conn.Disconnect()

	sub, err := conn.Subscribe("/queue/test-2", stomp.AckClient)
	if err != nil {
		return err
	}

	// receive 5 messages and then quit
	for i := 0; i < 5; i++ {
		msg := <-sub.C

		tx := conn.Begin()

		doAnotherThingWith(msg, tx)

		tx.Send(stomp.Message{
			Destination: "/queue/another-one",
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("Message #%d", i)),
		})

		// acknowledge the message
		err = tx.Ack(msg)
		if err != nil {
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	err = sub.Unsubscribe()
	if err != nil {
		return err
	}

	return nil
}

// Example of connecting to a STOMP server using an existing network connection.
func ExampleConnect() error {
	netConn, err := net.DialTimeout("tcp", "stomp.server.com:61613", 10*time.Second)
	if err != nil {
		return err
	}

	stompConn, err := stomp.Connect(netConn, stomp.Options{})
	if err != nil {
		return err
	}

	defer stompConn.Disconnect()

	doSomethingWith(stompConn)
	return nil
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
