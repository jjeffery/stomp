package stomp

import (
	"bufio"
	"io"
)

// Writes STOMP frames to an underlying io.Writer.
type Writer struct {
	writer *bufio.Writer
}

func NewWriter(writer io.Writer) *Writer {
	sw := new(Writer)
	sw.writer = bufio.NewWriterSize(writer, bufferSize)
	return sw
}

func (w *Writer) Write(frame *Frame) error {
	_, err := w.writer.Write([]byte(frame.Command))
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

	// write the final nul (0) byte	
	_, err = w.writer.Write(nullSlice)
	if err != nil {
		return err
	}

	err = w.writer.Flush()
	return err
}
