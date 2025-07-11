package headers

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"fmt"
)

func TestHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 37, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	headers["user-agent"] = "curl-any-version"
	data = []byte("       Host: localhost:42069       \r\n Accept: *.txt\r\n\r\n")
	done = false
	for !done {
		n, done, err = headers.Parse(data)
		fmt.Printf("Bytes parsed: %d Done: %v Headers:\n%#v\n", n, done, headers)
		copy(data, data[n:])
	}

	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "*.txt", headers["accept"])
	assert.Equal(t, "curl-any-version", headers["user-agent"])
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Valid done empty header
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 0, len(headers))
	assert.True(t, done)

	// Test: Invalid char in header name
	headers = NewHeaders()
	data = []byte("       H@st: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: VALID chars in header name
	headers = NewHeaders()
	data = []byte("!#$%&'*+-.^_`|~H9: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	// assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Append value to existing header
	headers = NewHeaders()
	headers["set-person"] = "lane-loves-go"
	data = []byte("Set-Person: prime-loves-zig\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, 29, n)
	assert.Equal(t, headers["set-person"], "lane-loves-go, prime-loves-zig")

}

