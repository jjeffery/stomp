package frame

import (
	"bytes"
	"strings"
)

var (
	replacerForEncodeValue = strings.NewReplacer(
		"\\", "\\\\",
		"\r", "\\r",
		"\n", "\\n",
		":", "\\c",
	)
	replacerForUnencodeValue = strings.NewReplacer(
		"\\r", "\r",
		"\\n", "\n",
		"\\c", ":",
		"\\\\", "\\",
	)
)

// Unencodes a header value using STOMP value encoding
// TODO: return error if invalid sequences found (eg "\t")
func unencodeValue(b []byte) (string, error) {
	s := replacerForUnencodeValue.Replace(string(b))
	return s, nil
}
