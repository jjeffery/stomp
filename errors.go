package stomp

const (
//	notConnected         = errorMessage("expected CONNECT or STOMP frame")
//	unexpectedCommand    = errorMessage("unexpected frame command")
//	unknownCommand       = errorMessage("unknown command")
//	missingHeader        = errorMessage("missing header")
//	receiptInConnect     = errorMessage("receipt header prohibited in CONNECT or STOMP frame")
//	authenticationFailed = errorMessage("authentication failed")
//	txAlreadyInProgress  = errorMessage("transaction already in progress")
//	txUnknown            = errorMessage("unknown transaction")
//	subscriptionInUse    = errorMessage("another subscription has the same id")
//	subscriptionNotFound = errorMessage("subscription not found")
//	notImplementedYet    = errorMessage("not implemented yet")
)

type errorMessage string

func (e errorMessage) Error() string {
	return string(e)
}
