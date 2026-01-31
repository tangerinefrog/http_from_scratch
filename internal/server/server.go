package server

import (
	"log"
	"net"

	"github.com/tangerinefrog/http_from_scratch/internal/request"
	"github.com/tangerinefrog/http_from_scratch/internal/response"
)

type Server struct {
	listener net.Listener
	state    serverState
	handlers map[string]Handler
}

type serverState string

const (
	serverStateRunning  serverState = "running"
	serverStateStopping serverState = "stopping"
)

type Handler func(w *response.Writer, r *request.Request)

func Serve(addr string) (*Server, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: listener,
		handlers: make(map[string]Handler),
	}

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	s.state = serverStateStopping
	err := s.listener.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Get(path string, h Handler) {
	s.handlers["GET"+path] = h
}

func (s *Server) Post(path string, h Handler) {
	s.handlers["POST"+path] = h
}

func (s *Server) getHandler(path, method string) (Handler, bool) {
	h, ok := s.handlers[method+path]
	return h, ok
}

func (s *Server) listen() {
	s.state = serverStateRunning

	for s.state == serverStateRunning {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.state == serverStateStopping {
				return
			}

			log.Printf("error opening a TCP connection: %v\n", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	writer := response.NewWriter(conn)

	r, err := request.RequestFromReader(conn)
	if err != nil {
		writer.WriteHeaders(response.StatusBadRequest)
		return
	}

	h, ok := s.getHandler(r.RequestLine.RequestTarget, r.RequestLine.Method)
	if !ok {
		writer.WriteHeaders(response.StatusNotFound)
		return
	}

	h(writer, r)
}
