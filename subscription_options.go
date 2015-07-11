package stomp

var SubscriptionOpt struct {
	Header func(header *Header) func(*Frame) error
}

func init() {
	SubscriptionOpt.Header = func(header *Header) func(*Frame) error {
		return func(f *Frame) error {
			f.AddHeader(header)
			return nil
		}
	}
}
