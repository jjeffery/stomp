package stomp

import (
	"strconv"
	"sync/atomic"
)

// BUG(jpj): All unique identifiers for receipt, subscriptions
// and transactions are allocated from the same namespace, and
// are shared between clients. This could give a hostile STOMP
// server information about the client.

var _lastId uint64

// allocateId returns a unique number for the current
// process. Starts at one and increases. Used for 
// allocating subscription ids, receipt ids, etc.
func allocateId() string {
	id := atomic.AddUint64(&_lastId, 1)
	return strconv.FormatUint(id, 10)
}
