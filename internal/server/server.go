package server

import (
	"fmt"
	"log"
	"net"

	"github.com/tangerinefrog/http_from_scratch/internal/request"
	"github.com/tangerinefrog/http_from_scratch/internal/response"
)

type Server struct {
	listener net.Listener
}

func Serve(addr string) (*Server, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: listener,
	}

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) listen() {

	for {
		conn, err := s.listener.Accept()
		if err != nil {
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
		log.Printf("Error while parsing request: %v", err)
		return
	}
	print(r)

	writer.Write([]byte("pong"))
}

func print(r *request.Request) {
	fmt.Printf("Got new HTTP request:\n\n   - Version: %s\n   - Method: %s\n   - Target: %s", r.RequestLine.HttpVersion, r.RequestLine.Method, r.RequestLine.RequestTarget)
	if len(r.Headers) > 0 {
		fmt.Printf("\n\nHeaders:\n\n")
		for k, v := range r.Headers {
			fmt.Printf("   %s:%s\n", k, v)
		}
	}
}
