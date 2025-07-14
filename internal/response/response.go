package response

import (
//	"bytes"
	"fmt"
	"io"
//	"net"
	"strconv"
	"github.com/trolioSFG/http-go/internal/headers"
)


type StatusCode int
const (
	StatusOK StatusCode = 200
	StatusBadRequest StatusCode = 400
	StatusError StatusCode = 500
)

type WriterState int
const (
	StateStatus WriterState = iota
	StateHeaders
	StateBody
)

type Writer struct {
	Buf io.Writer
	State WriterState
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.State != StateStatus {
		return fmt.Errorf("StatusLine: wrong state %d", w.State)
	}

	switch statusCode {
	case StatusOK:
		_, err := w.Buf.Write([]byte("HTTP/1.1 200 OK\r\n"))
		w.State = StateHeaders
		return err
	case StatusBadRequest:
		_, err := w.Buf.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		w.State = StateHeaders
		return err
	case StatusError:
		_, err := w.Buf.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		w.State = StateHeaders
		return err
	default:
		_, err := w.Buf.Write([]byte(fmt.Sprintf("HTTP/1.1 %d \r\n", statusCode)))
		w.State = StateHeaders
		return err
	}
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.State != StateHeaders {
		return fmt.Errorf("WriteHeaders: wrong state %d", w.State)
	}
	for k, v := range headers {
		_, err := w.Buf.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}

	// Better here, just in case there is NO body
	w.Buf.Write([]byte("\r\n"))
	w.State = StateBody
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.State != StateBody {
		return 0, fmt.Errorf("WriteBody: wrong state %d", w.State)
	}
	n, err := w.Buf.Write(p)
	return n, err
}


func GetDefaultHeaders(contentLen int) headers.Headers {
	hdrs := headers.NewHeaders()
	hdrs["Content-Length"] = strconv.Itoa(contentLen)
	hdrs["Connection"] = "close"
	hdrs["Content-Type"] = "text/plain"

	return hdrs
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}

	return nil
}


