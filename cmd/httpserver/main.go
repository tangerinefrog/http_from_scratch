package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server stopped")
}
