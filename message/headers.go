package message

// Collection of STOMP headers. Each header consists of a key value pair.
type Headers struct {
	headers []string
}

// Returns a new Headers collection.
func NewHeaders() *Headers {
	return new(Headers)
}

// Perform a deep copy of the Headers collection. Any changes
// to the cloned object will not affect the original object.
func (h *Headers) Clone() *Headers {
	clone := new(Headers)
	clone.headers = make([]string, len(h.headers))
	copy(clone.headers, h.headers)
	return clone
}

// Returns the number of headers in the collection.
func (h *Headers) Count() int {
	return len(h.headers) / 2
}

// Returns the number of headers in the collection.
func (h *Headers) Len() int {
	return len(h.headers) / 2
}

// Returns the header name and value at the specified index in
// the collection. The index should be in the range 0 <= index < Count().
func (h *Headers) GetAt(index int) (key, value string) {
	index *= 2
	return h.headers[index], h.headers[index+1]
}

// Sets the value for a key/value pair in the Headers collection. 
// If the key already exists its value is replaced. If the key does 
// not already exist it is added.
func (h *Headers) Set(key, value string) {
	if i, ok := h.index(key); ok {
		h.headers[i+1] = value
	} else {
		h.headers = append(h.headers, key, value)
	}
}

// Appends the key/value pair to the Headers collection without
// checking if the key already exists in the collection. Use this
// method when de-serializing the headers from the frame data, as
// a frame may contain multiple values for the same key. When this
// happens, the value for the first key is used and the other values
// are ignored.
func (h *Headers) Append(key, value string) {
	h.headers = append(h.headers, key, value)
}

// Removes the key/value pair from the Headers collection. Takes
// no action if the key does not already exist in the colleciton.
// If the key appears more than once in the collection, all values
// are removed.
func (h *Headers) Remove(key string) {
	for i, ok := h.index(key); ok; i, ok = h.index(key) {
		h.headers = append(h.headers[:i], h.headers[i+2:]...)
	}
}

// Returns the associated value and true if the Headers collection contains 
// the specified key.
func (h *Headers) Contains(key string) (string, bool) {
	if i, ok := h.index(key); ok {
		return h.headers[i+1], true
	}
	return "", false
}

// Returns the index of a header key in Headers, and a bool to indicate
// whether it was found or not.
func (h *Headers) index(key string) (int, bool) {
	for i := 0; i < len(h.headers); i += 2 {
		if h.headers[i] == key {
			return i, true
		}
	}
	return -1, false
}
