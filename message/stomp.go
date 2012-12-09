package message

const (
	bufferSize = 4096
	newline    = byte(10)
	colon      = byte(58)
	nullByte   = byte(0)
)

// STOMP protocol version
type StompVersion string

func (v StompVersion) GreaterThan(other StompVersion) bool {
	return v > other
}

// supported STOMP protocol versions
const (
	V1_0 = StompVersion("1.0")
	V1_1 = StompVersion("1.1")
	V1_2 = StompVersion("1.2")
)

type errorMessage string

const (
	invalidFrameFormat       = errorMessage("invalid frame format")
	invalidCommand           = errorMessage("invalid command")
	unknownVersion           = errorMessage("incompatible version")
	notConnectFrame          = errorMessage("operation valid for STOMP and CONNECT frames only")
	invalidHeartBeat         = errorMessage("invalid format for heart-beat")
	invalidOperationForFrame = errorMessage("invalid operation for frame")
	exceededMaxFrameSize     = errorMessage("exceeded max frame size")
)

func missingHeader(name string) errorMessage {
	return errorMessage("missing header: " + name)
}

func (e errorMessage) Error() string {
	return string(e)
}
