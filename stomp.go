package stomp


const (
	bufferSize = 4096
	newline    = byte(10)
	colon      = byte(58)
	nullByte   = byte(0)
)

type errorMessage string

const (
	invalidFrameFormat = errorMessage("invalid frame format")
)

func (e errorMessage) Error() string {
	return string(e)
}

