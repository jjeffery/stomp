package stomp

// Version is the STOMP protocol version.
type Version string

const (
	V10 Version = "1.0"
	V11 Version = "1.1"
	V12 Version = "1.2"
)

func (v Version) String() string {
	return string(v)
}
