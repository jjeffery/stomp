package stomp_test

import (
	"github.com/jjeffery/stomp"
)

func ExampleDial() error {
	client, err := stomp.Dial("tcp", "192.168.1.1:61613", stomp.ConnectOptions{})
	if err != nil {
		return err
	}

	err = client.Send(stomp.SendMessage{
		Destination: "/queue/test-1",
		ContentType: "text/plain",
		Body:        []byte("Test message #1"),
	})
	if err != nil {
		return err
	}

	return client.Disconnect()
}
