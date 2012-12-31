package stomp

// A Header represents a STOMP header mapping 
// keys to sets of values. 
//
// Normally a STOMP header
// only has one value, but the STOMP standard does
// allow multiple values for diagnostic purposes.
//
// This type is very similar to textproto.MIMEHeader. The main
// difference is that STOMP header keys are case-sensitive.
type Header map[string][]string

// Add adds the key, value pair to the header.
// It appends to any existing values associated with the key.
func (h Header) Add(key, value string) {
	h[key] = append(h[key], value)
}

// Set sets the header entries associated with 
func (h Header) Set(key, value string) {
	h[key] = []string{value}
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns "".
// Get is a convenience method. For more complex queries, access
// the map directly.
func (h Header) Get(key string) string {
	if h == nil {
		return ""
	}
	values := h[key]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

// Del deletes the values associated with key.
func (h Header) Del(key string) {
	delete(h, key)
}

// Clone returns a deep copy of a Header
func (h Header) Clone() Header {
	hc := Header{}
	for k, v := range h {
		if len(v) > 0 {
			vc := make([]string, len(v))
			copy(vc, v)
			hc[k] = vc
		}
	}
	return hc
}
