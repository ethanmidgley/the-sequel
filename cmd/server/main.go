package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethanmidgley/the-sequel/in-memory/pkg/server"
)

func main() {

	s, err := server.New("127.0.0.1:6379")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s.Start()

	// Wait for a SIGINT or SIGTERM signal to gracefully shut down the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down server...")
	s.Stop()
	fmt.Println("Server stopped.")

}
