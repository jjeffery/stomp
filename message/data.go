package message

const (
	// client frames
	ConnectCommand = "CONNECT"
	Send           = "SEND"
	Subscribe      = "SUBSCRIBE"
	Unsubscribe    = "UNSUBSCRIBE"
	Ack            = "ACK"
	Nack           = "NACK"
	Begin          = "BEGIN"
	Commit         = "COMMIT"
	Abort          = "ABORT"
	Disconnect     = "DISCONNECT"

	// server frames
	Message = "MESSAGE"
	Receipt = "RECEIPT"
	Error   = "ERROR"

	// header names
	ContentLength = "content-length"
	ContentType   = "content-type"
	ReceiptHeader = "receipt"
	AcceptVersion = "accept-version"
	Host          = "host"
	Version       = "version"
	Login         = "login"
	Passcode      = "passcode"
	HeartBeat     = "heart-beat"
	Session       = "session"
	Server        = "server"
	Destination   = "destination"
	Id            = "id"
	AckHeader     = "ack"
	Transaction   = "transaction"
	ReceiptId     = "receipt-id"
	Subscription  = "subscription"
	MessageId     = "message-id"
	MessageHeader = "message"
)
