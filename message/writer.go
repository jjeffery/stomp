package message

import (
	"bufio"
	"io"
)

// slices used to write frames
var (
	colonSlice   = []byte{58}     // colon ':'
	crlfSlice    = []byte{13, 10} // CR-LF
	newlineSlice = []byte{10}     // newline (LF)
	nullSlice    = []byte{0}      // null character
)

// Writes STOMP frames to an underlying io.Writer.
type Writer struct {
	writer *bufio.Writer
}

// Creates a new Writer object, which writes to an underlying io.Writer.
func NewWriter(writer io.Writer) *Writer {
	sw := new(Writer)
	sw.writer = bufio.NewWriterSize(writer, bufferSize)
	return sw
}

// Write the contents of a frame to the underlying io.Writer.
func (w *Writer) Write(frame *Frame) error {
	var err error

	if frame == nil {
		// nil frame means send a heart-beat LF
		_, err = w.writer.Write(newlineSlice)
		if err != nil {
			return err
		}
	} else {
		_, err = w.writer.Write([]byte(frame.Command))
		if err != nil {
			return err
		}

		_, err = w.writer.Write(newlineSlice)
		if err != nil {
			return err
		}

		headerCount := frame.Headers.Count()
		for i := 0; i < headerCount; i++ {
			h, v := frame.Headers.GetAt(i)
			// TODO: encode if STOMP 1.1 or later
			_, err = w.writer.Write([]byte(h))
			if err != nil {
				return err
			}

			_, err = w.writer.Write(colonSlice)
			if err != nil {
				return err
			}

			_, err = w.writer.Write([]byte(v))
			if err != nil {
				return err
			}

			_, err = w.writer.Write(newlineSlice)
			if err != nil {
				return err
			}
		}

		_, err = w.writer.Write(newlineSlice)
		if err != nil {
			return err
		}

		if len(frame.Body) > 0 {
			_, err = w.writer.Write(frame.Body)
			if err != nil {
				return err
			}
		}

		// write the final null (0) byte	
		_, err = w.writer.Write(nullSlice)
		if err != nil {
			return err
		}
	}

	err = w.writer.Flush()
	return err
}
