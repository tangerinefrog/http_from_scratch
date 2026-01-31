package response

import (
	"fmt"
	"net"

	"github.com/tangerinefrog/http_from_scratch/internal/headers"
)

type StatusCode int32

const (
	StatusOK                  StatusCode = 200
	StatusCreated             StatusCode = 201
	StatusBadRequest          StatusCode = 400
	StatusUnauthorized        StatusCode = 401
	StatusNotFound            StatusCode = 404
	StatusMethodNotAllowed    StatusCode = 405
	StatusInternalServerError StatusCode = 500
)

func GetDefaultHeaders() headers.Headers {
	h := headers.NewHeaders()
	h.Set("Connection", "close")

	return h
}

type Writer struct {
	headers headers.Headers
	conn    net.Conn
	state   ResponseState
}

type ResponseState int

const (
	stateHeaders = iota
	stateDone
)

func NewWriter(conn net.Conn) *Writer {
	return &Writer{
		headers: GetDefaultHeaders(),
		conn:    conn,
		state:   stateHeaders,
	}
}

func (w *Writer) Write(body []byte) (int, error) {
	// if no headers were written - assume 200 OK
	if w.state <= stateHeaders {
		err := w.WriteHeaders(StatusOK)
		if err != nil {
			return 0, err
		}
	}

	n, err := w.conn.Write(body)
	return n, err
}

func (w *Writer) WriteHeaders(code StatusCode) error {
	if w.state > stateHeaders {
		return fmt.Errorf("headers have been already written to the client")
	}

	sl := getStatusLine(code)
	_, err := w.conn.Write(sl)
	if err != nil {
		return err
	}

	cType := w.headers.Get("Content-Type")
	if cType == "" {
		w.headers.Set("Content-Type", "text/plain")
	}

	for n, v := range w.headers {
		name := headers.CapitalizeName(n)
		header := fmt.Sprintf("%s: %s\r\n", name, v)

		_, err := w.conn.Write([]byte(header))
		if err != nil {
			return err
		}
	}

	w.conn.Write([]byte("\r\n"))

	w.state = stateDone

	return nil
}

func (w *Writer) Error(message string, code StatusCode) error {
	if w.state > stateHeaders {
		return fmt.Errorf("headers have been already written to the client")
	}
	w.DeleteHeader("Content-Type")
	w.AddHeader("Content-Type", "text/plain")
	err := w.WriteHeaders(code)
	if err != nil {
		return err
	}

	_, err = w.conn.Write([]byte(message))
	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) AddHeader(name, value string) {
	w.headers.Set(name, value)
}

func (w *Writer) DeleteHeader(name string) {
	w.headers.Remove(name)
}

func getStatusLine(code StatusCode) []byte {
	reason := statusReason(code)
	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", code, reason))
}

func statusReason(code StatusCode) string {
	switch code {
	case StatusOK:
		return "OK"
	case StatusCreated:
		return "Created"
	case StatusBadRequest:
		return "Bad Request"
	case StatusUnauthorized:
		return "Unauthorized"
	case StatusNotFound:
		return "Not Found"
	case StatusMethodNotAllowed:
		return "Method Not Allowed"
	case StatusInternalServerError:
		return "Internal Server Error"
	default:
		return ""
	}
}
