package main

import (
	"lamport_demo/server"
	"log"
)

func main() {
	server, err := server.CreateNewServer()
	if err != nil {
		log.Fatalf("Failed to create sever %v", server)
	}
	server.ProcessMessages()
}
