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
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	rLine, err := parseRequestLine(data)
	if err != nil {
		return nil, err
	}

	rl := RequestLine {
		HttpVersion: rLine[2],
		RequestTarget: rLine[1],
		Method: rLine[0],
	}

	req := Request {
		RequestLine: rl,
	}

	return &req, nil
}

func parseRequestLine(data []byte) (int, []string, error) {
	lines := strings.Split(string(data), "\r\n")
	// fmt.Printf("Request line: %v\n", lines[0])
	if len(lines) == 1 && lines[0] == string(data) {
		return 0, nil, nil
	}

	parts := strings.Split(lines[0], " ")

	if len(parts) != 3 {
		return nil, fmt.Errorf("Bad request line")
	}

	// fmt.Printf("%v\n", parts)

	for _, l := range parts[0] {
		if !unicode.IsLetter(l) {
			return nil, fmt.Errorf("No alphabetic")
		}
	}

	if parts[0] != strings.ToUpper(parts[0]) {
		return nil, fmt.Errorf("Not capital method")
	}

	version := strings.Split(parts[2], "/")
	if len(version) != 2 {
		return nil, fmt.Errorf("Bad version")
	}

	// fmt.Printf("Version: %v\n", version)

	if version[1] != "1.1" {
		return nil, fmt.Errorf("Unsupported version")
	}

	parts[2] = version[1]

	return parts, nil
}

func (r *Request) parse(data []byte) (int, error) {

