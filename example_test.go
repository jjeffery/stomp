package stomp_test

import (
	"github.com/jjeffery/stomp"
)

func ExampleDial() error {
	conn, err := stomp.Dial("tcp", "192.168.1.1:61613", stomp.ConnectOptions{})
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
