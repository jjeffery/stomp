package client

const (
	notConnected = errorMessage("expected CONNECT or STOMP frame")
)

type errorMessage string

func (e errorMessage) Error() string {
	return string(e)
}
