package request

import (
	"io"
	"fmt"
	"strings"
	"unicode"
)

type parserState int
const (
	Initialized parserState = iota
	Done
)

type Request struct {
	RequestLine RequestLine
	pState parserState
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

	return len(data), parts, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.pState == Done {
		return 0, fmt.Errorf("Error: trying to read data in done state")
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
	r.pState = Done

	return n, nil
}
