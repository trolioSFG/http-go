package response

import (
//	"bytes"
	"fmt"
	"io"
	"net"
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
	StateBlank WriterState = iota
	StateStatus
	StateHeaders
	StateBody
)

type Writer struct {
	Buf net.Conn
	State WriterState
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.State != StateBlank {
		return fmt.Errorf("Trying to WriteStatusLine: state != Blank")
	}

	switch statusCode {
	case StatusOK:
		_, err := w.Buf.Write([]byte("HTTP/1.1 200 OK\r\n"))
		return err
	case StatusBadRequest:
		_, err := w.Buf.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		return err
	case StatusError:
		_, err := w.Buf.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		return err
	default:
		_, err := w.Buf.Write([]byte(fmt.Sprintf("HTTP/1.1 %d \r\n", statusCode)))
		return err
	}
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	for k, v := range headers {
		_, err := w.Buf.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	w.Buf.Write([]byte("\r\n"))
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


