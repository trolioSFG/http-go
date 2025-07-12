package response

import (
	"fmt"
	"io"
	"strconv"
	"github.com/trolioSFG/http-go/internal/headers"
)


type StatusCode int
const (
	StatusOK StatusCode = 200
	StatusBadRequest StatusCode = 400
	StatusError StatusCode = 500
)


func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusOK:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		return err
	case StatusBadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		return err
	case StatusError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		return err
	default:
		_, err := w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d \r\n", statusCode)))
		return err
	}
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


