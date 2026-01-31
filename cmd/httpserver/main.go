package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tangerinefrog/http_from_scratch/internal/handlers"
	"github.com/tangerinefrog/http_from_scratch/internal/server"
)

const addr = ":8080"

func main() {
	server, err := server.Serve(addr)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Printf("Server started on '%s'...", addr)

	setupHandlers(server)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server stopped")
}

func setupHandlers(s *server.Server) {
	s.Get("/", handlers.TestGoodHandler)
	s.Get("/yourproblem", handlers.TestYourProblemHandler)
	s.Get("/myproblem", handlers.TestMyProblemHandler)
}
