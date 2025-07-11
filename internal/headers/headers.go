package headers

import (
	"strings"
	"fmt"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(map[string]string)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	lines := strings.Split(string(data), "\r\n")
	if len(lines) == 1 && lines[0] == string(data) {
		return 0, false, nil
	}

	if lines[0] == "" {
		fmt.Printf("End of Header found!\n")
		return 2, true, nil
	}

	hdrName, hdrValue, found := strings.Cut(lines[0], ":")
	if !found {
		return 0, false, nil
	}

	if hdrName != strings.TrimRight(hdrName, " ") {
		return 0, false, fmt.Errorf("Blank in header name before :")
	}

	
	hdrName = strings.ToLower(strings.TrimLeft(hdrName, " "))

	for _, char := range hdrName {
		if !unicode.IsLetter(char) {
			if !unicode.IsDigit(char) {
				if index := strings.Index("!#$%&'*+-.^_`|~", string(char)); index == -1 {
					return 0, false, fmt.Errorf("Error header name: invalid character %v(%s)", char, string(char))
				}
			}
		}
	}


	hdrValue = strings.Trim(hdrValue, " ")

	value, ok := h[hdrName]
	if ok {
		h[hdrName] = value + ", " + hdrValue
	} else {
		h[hdrName] = hdrValue
	}

	// Count \r\n in bytes consumed
	return len(lines[0]) + 2, false, nil
}

