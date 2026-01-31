package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/tangerinefrog/http_from_scratch/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       RequestState
	contentLen  int
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type RequestState string

var StateInit RequestState = "init"
var StateHeaders RequestState = "headers"
var StateBody RequestState = "body"
var StateDone RequestState = "done"

var sep = []byte("\r\n")

var ErrMalformedRequest = errors.New("malformed request")
var ErrBodyTooLong = errors.New("request body is greater than Content-Length")

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	buf := make([]byte, 8)
	bufIdx := 0

	for !request.done() {
		if bufIdx >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[bufIdx:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if request.state != StateDone {
					return nil, fmt.Errorf("incomplete request, in state: %s, read n bytes on EOF: %d", request.state, n)
				}
				break
			}
			return nil, err
		}

		bufIdx += n
		parsedN, err := request.parse(buf[:bufIdx])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[parsedN:bufIdx])
		bufIdx -= parsedN
	}

	return request, nil
}

func (r *Request) done() bool {
	return r.state == StateDone
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer:
	for {
		switch r.state {
		case StateInit:
			n, rl, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err
			}
			// needs more data
			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n
			r.state = StateHeaders

		case StateHeaders:
			n, done, err := r.Headers.Parse(data[read:])
			if err != nil {
				return 0, err
			}
			// needs more data
			if n == 0 {
				break outer
			}

			read += n

			if done {
				getContentLen(r)
				if r.contentLen > 0 {
					r.Body = make([]byte, 0, r.contentLen)
					r.state = StateBody
				} else {
					r.state = StateDone
				}
			}

		case StateBody:
			remaining := min(r.contentLen-len(r.Body), len(data[read:]))
			r.Body = append(r.Body, data[read:remaining+read]...)
			read += remaining
			if len(r.Body) > r.contentLen {
				return 0, ErrBodyTooLong
			} else if len(r.Body) == r.contentLen {
				r.state = StateDone
			}
			break outer

		case StateDone:
			break outer
		}
	}

	return read, nil
}

func parseRequestLine(data []byte) (int, *RequestLine, error) {
	sepIdx := bytes.Index(data, sep)
	if sepIdx == -1 {
		return 0, nil, nil
	}

	startLine := data[:sepIdx]
	read := sepIdx + len(sep)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return 0, nil, ErrMalformedRequest
	}

	method := parts[0]
	if !validateMethod(method) {
		return 0, nil, errors.Join(ErrMalformedRequest, fmt.Errorf("invalid method"))
	}

	target := parts[1]

	version, err := parseHttpVersion(parts[2])
	if err != nil {
		return 0, nil, errors.Join(ErrMalformedRequest, err)
	}

	return read,
		&RequestLine{
			Method:        string(method),
			RequestTarget: string(target),
			HttpVersion:   string(version),
		}, nil
}

func validateMethod(m []byte) bool {
	for _, v := range m {
		if v < 'A' || v > 'Z' {
			return false
		}
	}

	return true
}

func parseHttpVersion(s []byte) ([]byte, error) {
	parts := bytes.Split(s, []byte("/"))
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid HTTP version")
	}

	if string(parts[0]) != "HTTP" || string(parts[1]) != "1.1" {
		return nil, fmt.Errorf("invalid HTTP version")
	}

	return parts[1], nil
}

func getContentLen(r *Request) {
	lenStr := r.Headers.Get("content-length")
	if lenStr == "" {
		return
	}

	length, err := strconv.Atoi(lenStr)
	if err != nil || length <= 0 {
		return
	}

	r.contentLen = length
}
