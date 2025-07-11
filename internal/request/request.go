package request

import (
	"io"
	"fmt"
	"strings"
	"unicode"
	"github.com/trolioSFG/http-go/internal/headers"

)

type parserState int
const (
	Initialized parserState = iota
	Done
	ParsingHeaders
)

type Request struct {
	RequestLine RequestLine
	pState parserState
	Headers headers.Headers
}

type RequestLine struct {
	HttpVersion string
	RequestTarget string
	Method string
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	buf := make([]byte, 8, 8)
	bytesRead := 0

	rq := &Request{
		pState: Initialized,
		Headers: headers.NewHeaders(),
	}

	for rq.pState != Done {
		if bytesRead >= len(buf) {
			aux := make([]byte, 2 * len(buf))
			copy(aux, buf)
			buf = aux
		}

		numRead, err := reader.Read(buf[bytesRead:])
		/***
		fmt.Printf("nread: %d bytesRead: %d\nBuf:\n%v\n",
			numRead, bytesRead, string(buf))
		*****/

		if err != nil {
			if err == io.EOF {
				// if rq.pState != Done {
				/** TODO:
				if bytesRead > 0 {
					fmt.Printf("Buffer: %v\n", buf)
					fmt.Printf("%#v\n", rq)
					return nil, fmt.Errorf("Incomplete request: in state %d read %d bytes at EOF", rq.pState, bytesRead)
				}
				**/
				rq.pState = Done
				break
			}

			return nil, err
		}
		bytesRead += numRead
		numParsed, err := rq.parse(buf[:bytesRead])
		if err != nil {
			fmt.Printf("Error parsing: %v\n", err)
			return nil, err
		}

		// fmt.Printf("bytesRead: %d\n", bytesRead)
		if rq.pState == Done {
			fmt.Printf("Finished bytesRead: %d bytesParsed: %d\n", bytesRead, numParsed)
		}

		// HERE!
		copy(buf, buf[numParsed:])
		bytesRead -= numParsed
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

	// fmt.Printf("Version: %v\n", version)

	if version[1] != "1.1" {
		return 0, nil, fmt.Errorf("Unsupported version")
	}

	parts[2] = version[1]

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
			fmt.Printf("Done!")
			r.pState = Done
		}

		return bytesHdr, nil
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
