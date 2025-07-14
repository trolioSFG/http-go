package request

import (
	"io"
	"fmt"
//	"log"
	"strings"
	"unicode"
	"github.com/trolioSFG/http-go/internal/headers"
	"strconv"
)

type parserState int
const (
	Initialized parserState = iota
	Done
	ParsingHeaders
	ParsingBody
)

type Request struct {
	RequestLine RequestLine
	pState parserState
	Headers headers.Headers
	Body []byte
	bodyLengthRead int
}

type RequestLine struct {
	HttpVersion string
	RequestTarget string
	Method string
}

func (rql RequestLine) String() string {
	return rql.Method + " " + rql.RequestTarget + " " + "HTTP/" + rql.HttpVersion
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	buf := make([]byte, 8, 8)
	bytesRead := 0
	bytesParsed := 0

	rq := &Request{
		pState: Initialized,
		Headers: headers.NewHeaders(),
		Body: []byte{},
	}

	for rq.pState != Done {
		if bytesRead >= len(buf) {
			aux := make([]byte, 2 * len(buf))
			copy(aux, buf)
			buf = aux
			fmt.Printf("Buffer extended to %d\n", len(buf))
		}


		numRead, err := reader.Read(buf[bytesRead:])
		fmt.Printf("Read %d bytes\n", numRead)

		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("Incomplete request")
			}
			return nil, err
		}

		bytesRead += numRead
		bytesParsed, err = rq.parse(buf[:bytesRead])
		if err != nil {
			fmt.Printf("Error parsing: %v\n", err)
			return nil, err
		}
		copy(buf, buf[bytesParsed:])
		bytesRead -= bytesParsed

		for rq.pState != Done && bytesParsed > 0 {
			bytesParsed, err = rq.parse(buf[:bytesRead])
			if err != nil {
				fmt.Printf("Error parsing: %v\n", err)
				return nil, err
			}
			// HERE!
			fmt.Printf("PRE bytesRead: %d bytesParsed: %d\nBuffer:%v\n", bytesRead, bytesParsed, buf[:bytesRead])
			copy(buf, buf[bytesParsed:])
			bytesRead -= bytesParsed
			fmt.Printf("bytesRead: %d bytesParsed: %d\nBuffer:%v\n", bytesRead, bytesParsed, buf[:bytesRead])
		}
	}

	return rq, nil
}

func parseRequestLine(data []byte) (int, []string, error) {
	lines := strings.Split(string(data), "\r\n")
	if len(lines) == 1 && lines[0] == string(data) {
		return 0, nil, nil
	}

	parts := strings.Split(lines[0], " ")

	if len(parts) != 3 {
		return 0, nil, fmt.Errorf("Bad request line")
	}

	for _, l := range parts[0] {
		if !unicode.IsLetter(l) {
			return 0, nil, fmt.Errorf("No alphabetic")
		}
	}

	if parts[0] != strings.ToUpper(parts[0]) {
		return 0, nil, fmt.Errorf("Not capital method")
	}

	version := strings.Split(parts[2], "/")
	if len(version) != 2 {
		return 0, nil, fmt.Errorf("Bad version")
	}

	if version[1] != "1.1" {
		return 0, nil, fmt.Errorf("Unsupported version")
	}

	parts[2] = version[1]

	fmt.Printf("Finished request line\n")
	return len(lines[0]) + 2, parts, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.pState == Done {
		return 0, fmt.Errorf("Error: trying to read data in done state")
	}

	if r.pState == ParsingHeaders {
		bytesHdr, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}

		if done {
			fmt.Printf("Finished parsing headers\n")
			r.pState = ParsingBody
			// r.pState = Done
		}

		return bytesHdr, nil
	}

	if r.pState == ParsingBody {
		clStr, body := r.Headers["content-length"]
		if !body {
			fmt.Printf("No content-length header\n")
			r.pState = Done
			return len(data), nil
		}
		contentLen, err := strconv.Atoi(clStr)
		if err != nil {
			return 0, fmt.Errorf("malformed Content-Length: %s", err)
		}
		// TIL: Variadic
		r.Body = append(r.Body, data...)
		r.bodyLengthRead += len(data)

		if r.bodyLengthRead > contentLen {
			return 0, fmt.Errorf("Content-Length too large")
		}

		if r.bodyLengthRead == contentLen {
			r.pState = Done
		}

		return len(data), nil
	}


	if r.pState != Initialized {
		return 0, fmt.Errorf("Error: unknown state")
	}

	n, parts, err := parseRequestLine(data)
	if err != nil {
		return 0, err
	}

	if n == 0 {
		return 0, nil
	}

	r.RequestLine.HttpVersion = parts[2]
	r.RequestLine.RequestTarget = parts[1]
	r.RequestLine.Method = parts[0]
	r.pState = ParsingHeaders

	return n, nil
}
