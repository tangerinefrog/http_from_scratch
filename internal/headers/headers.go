package headers

import (
	"bytes"
	"errors"
	"strings"
)

type Headers map[string]string

var sep = []byte(": ")
var clrf = []byte("\r\n")

var ErrMalformedFieldLine = errors.New("malformed field line")
var ErrMalformedFieldName = errors.New("malformed field name")

func NewHeaders() Headers {
	h := make(map[string]string)
	return h
}

func (h Headers) Get(name string) string {
	name = strings.ToLower(name)
	v := h[name]
	return v
}

func (h Headers) Set(name, value string) {
	name = strings.ToLower(name)
	val, ok := h[name]
	if ok {
		value = strings.Join([]string{val, value}, ", ")
	}

	h[name] = value
}

func (h Headers) Remove(name string) {
	delete(h, name)
}

func (h Headers) Length() int {
	return len(h)
}

func (h *Headers) Parse(data []byte) (n int, done bool, err error) {
	read := 0

	for {
		clrfIdx := bytes.Index(data[read:], clrf)
		if clrfIdx == -1 {
			// needs more data
			break
		}

		fieldLine := data[read : clrfIdx+read]

		sepIdx := bytes.Index(fieldLine, sep)
		read += clrfIdx + len(clrf)

		if sepIdx == -1 {
			// end of headers
			if len(fieldLine) == 0 {
				return read, true, nil
			}

			// no 'key: value' in request line
			return 0, false, ErrMalformedFieldLine
		}

		name, value, err := parseFieldLine(fieldLine)
		if err != nil {
			return 0, false, err
		}

		h.Set(name, value)
	}

	return read, false, nil
}

func parseFieldLine(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", ErrMalformedFieldLine
	}

	name := parts[0]
	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", ErrMalformedFieldName
	}

	name = bytes.ToLower(bytes.TrimSpace(name))
	if !isValidName(name) {
		return "", "", ErrMalformedFieldName
	}

	value := bytes.TrimSpace(parts[1])

	return string(name), string(value), nil
}

func isValidName(fieldName []byte) bool {
	for _, c := range fieldName {
		if !allowedChars[c] {
			return false
		}
	}

	return true
}

var allowedChars = map[byte]bool{
	'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true, 'g': true,
	'h': true, 'i': true, 'j': true, 'k': true, 'l': true, 'm': true, 'n': true,
	'o': true, 'p': true, 'q': true, 'r': true, 's': true, 't': true, 'u': true,
	'v': true, 'w': true, 'x': true, 'y': true, 'z': true,
	'0': true, '1': true, '2': true, '3': true, '4': true, '5': true, '6': true,
	'7': true, '8': true, '9': true,
	'!': true, '#': true, '$': true, '%': true, '&': true, '\'': true, '*': true,
	'+': true, '-': true, '.': true, '^': true, '_': true, '`': true, '|': true, '~': true,
}
