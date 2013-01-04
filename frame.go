package stomp

// A Frame represents a STOMP frame. A frame consists of a command
// followed by a collection of header elements, and then an optional
// body.
type Frame struct {
	Command string
	*Header
	Body []byte
}

// Creates a new frame with the specified command and headers. The headers
// should contain an even number of entries. Each even index is the header 
// name, and the odd indexes are the assocated header values.
func NewFrame(command string, headers ...string) *Frame {
	f := &Frame{Command: command, Header: &Header{}}
	for index := 0; index < len(headers); index += 2 {
		f.Add(headers[index], headers[index+1])
	}
	return f
}

func (f *Frame) Validate() error {
	// TODO(jpj): implement
	return nil
}
