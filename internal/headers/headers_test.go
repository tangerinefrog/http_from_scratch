package headers

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	s := "Host: localhost:8080\r\n\r\n"
	data := []byte(s)
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:8080", headers.Get("Host"))
	assert.Equal(t, len(s), n)
	assert.True(t, done)

	// Test: Valid 2 headers
	headers = NewHeaders()
	s = "Host: localhost:8080\r\nAccept: */*\r\n\r\n"
	data = []byte(s)
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 2, len(headers))
	assert.Equal(t, "localhost:8080", headers.Get("Host"))
	assert.Equal(t, "*/*", headers.Get("Accept"))
	assert.Equal(t, len(s), n)
	assert.True(t, done)

	// Test: No headers
	headers = NewHeaders()
	s = "\r\n{\"key\":\"value\"}"
	data = []byte(s)
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 0, len(headers))
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Consumes only headers
	headers = NewHeaders()
	s = "Host: localhost:8080\r\n\r\n{\"key\":\"value\"}"
	data = []byte(s)
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:8080", headers.Get("Host"))
	assert.Equal(t, 1, len(headers))
	assert.Equal(t, 24, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	s = "       Host : localhost:8080       \r\n\r\n"
	data = []byte(s)
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid field name
	headers = NewHeaders()
	s = "H©st: localhost:8080\r\n\r\n"
	data = []byte(s)
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid case insensitive field name
	headers = NewHeaders()
	s = "HoSt: localhost:8080\r\naCCept: */*\r\n\r\n"
	data = []byte(s)
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 2, len(headers))
	assert.Equal(t, "localhost:8080", headers.Get("hOsT"))
	assert.Equal(t, "*/*", headers.Get("ACCEPT"))
	assert.Equal(t, len(s), n)
	assert.True(t, done)

	// Test: Valid multiple values with the same field name
	headers = NewHeaders()
	values := []string{"lane-loves-go", "prime-loves-zig", "tj-loves-ocaml"}
	s = ""
	for _, v := range values {
		s += fmt.Sprintf("Set-Person: %s\r\n", v)
	}
	s += "\r\n"
	data = []byte(s)
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 1, len(headers))
	assert.Equal(t, strings.Join(values, ", "), headers.Get("set-person"))
	assert.Equal(t, len(s), n)
	assert.True(t, done)
}
