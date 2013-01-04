package stomp

import (
	"github.com/jjeffery/stomp/frame"
	"strconv"
)

// A Header represents a STOMP header mapping 
// keys to sets of values. 
//
// Normally a STOMP header
// only has one value, but the STOMP standard does
// allow multiple values for diagnostic purposes.
//
// This type is very similar to textproto.MIMEHeader. The main
// difference is that STOMP header keys are case-sensitive.
type Header struct {
	slice []string
}

func NewHeader(headerEntries ...string) *Header {
	h := &Header{}
	h.slice = append(h.slice, headerEntries...)
	if len(h.slice)%2 != 0 {
		h.slice = append(h.slice, "")
	}
	return h
}

// Add adds the key, value pair to the header.
// It appends to any existing values associated with the key.
func (h *Header) Add(key, value string) {
	h.slice = append(h.slice, key, value)
}

// Set sets the header entries associated with 
func (h *Header) Set(key, value string) {
	if i, ok := h.index(key); ok {
		h.slice[i+1] = value
	} else {
		h.slice = append(h.slice, key, value)
	}
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns "".
func (h *Header) Get(key string) string {
	value, _ := h.Contains(key)
	return value
}

// GetAll returns all of the values associated with a given key.
// Normally there is only one header entry per key, but it is permitted
// to have multiple entries according to the STOMP standard.
func (h *Header) GetAll(key string) []string {
	var values []string
	for i := 0; i < len(h.slice); i += 2 {
		if h.slice[i] == key {
			values = append(values, h.slice[i+1])
		}
	}
	return values
}

// Returns the header name and value at the specified index in
// the collection. The index should be in the range 0 <= index < Len(),
// a panic will occur if it is outside this range.
func (h *Header) GetAt(index int) (key, value string) {
	index *= 2
	return h.slice[index], h.slice[index+1]
}

// Contains gets the first value associated with the given key, 
// and also returns a bool indicating whether the header entry 
// exists.
//
// If there are no values associated with the key, Get returns ""
// for the value, and ok is false.
func (h *Header) Contains(key string) (value string, ok bool) {
	var i int
	if i, ok = h.index(key); ok {
		value = h.slice[i+1]
	}
	return
}

// Del deletes the values associated with key.
func (h *Header) Del(key string) {
	for i, ok := h.index(key); ok; i, ok = h.index(key) {
		h.slice = append(h.slice[:i], h.slice[i+2:]...)
	}
}

func (h *Header) Len() int {
	return len(h.slice) / 2
}

// Clone returns a deep copy of a Header.
func (h *Header) Clone() *Header {
	hc := &Header{slice: make([]string, len(h.slice))}
	copy(hc.slice, h.slice)
	return hc
}

// ContentLength returns the value of the "content-length" header entry.
// If the "content-length" header is missing, then ok is false. If the 
// "content-length" entry is present but is not a valid non-negative integer
// then err is non-nil.
func (h *Header) ContentLength() (value int, ok bool, err error) {
	text := h.Get(frame.ContentLength)
	if text == "" {
		return
	}

	n, err := strconv.ParseUint(text, 10, 32)
	if err != nil {
		return
	}

	value = int(n)
	ok = true
	return
}

// Returns the index of a header key in Headers, and a bool to indicate
// whether it was found or not.
func (h *Header) index(key string) (int, bool) {
	for i := 0; i < len(h.slice); i += 2 {
		if h.slice[i] == key {
			return i, true
		}
	}
	return -1, false
}
