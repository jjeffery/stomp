package stomp

import (
	"time"
)

// Convert a time.Duration to milliseconds in an integer.
// If the time duration is too large, it is truncated to
// the largest number of milliseconds that can be represented
// as a 32 bit signed integer (about 40 days).
func asMilliseconds(d time.Duration, max int) int {
	if max < 0 {
		max = 0
	}
	max64 := int64(max)
	msec64 := int64(d / time.Millisecond)
	if msec64 > max64 {
		msec64 = max64
	}
	msec := int(msec64)
	return msec
}