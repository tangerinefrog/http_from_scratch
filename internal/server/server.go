package server

import (
	"log"
	"net"
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

	conn.Write([]byte("pong"))
}
