package message

// STOMP frame commands. Used upper case naming
// convention to avoid clashing with STOMP header names.
const (
	// connect frames
	CONNECT   = "CONNECT"
	STOMP     = "STOMP"
	CONNECTED = "CONNECTED"

	// client frames
	SEND        = "SEND"
	SUBSCRIBE   = "SUBSCRIBE"
	UNSUBSCRIBE = "UNSUBSCRIBE"
	ACK         = "ACK"
	NACK        = "NACK"
	BEGIN       = "BEGIN"
	COMMIT      = "COMMIT"
	ABORT       = "ABORT"
	DISCONNECT  = "DISCONNECT"

	// server frames
	MESSAGE = "MESSAGE"
	RECEIPT = "RECEIPT"
	ERROR   = "ERROR"
)
