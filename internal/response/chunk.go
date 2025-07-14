package response

import (
	"fmt"
)

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {

	chunkLen := fmt.Sprintf("%X", len(p)) + "\r\n"
	count := 0
	n, err := w.WriteBody([]byte(chunkLen))
	if err != nil {
		return n, err
	}
	count += n
	n, err = w.WriteBody(p)
	if err != nil {
		return count, err
	}
	count += n
	n, err = w.WriteBody([]byte("\r\n"))
	if err != nil {
		return count, err
	}
	count += n
	return count, err
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	n, err := w.WriteBody([]byte("0\r\n\r\n"))
	return n, err
}


