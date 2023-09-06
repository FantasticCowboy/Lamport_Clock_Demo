package main

import (
	"lamport_demo/server"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Did not pass in an ip address!")
	}
	ip := os.Args[1]
	server, err := server.CreateNewServer(ip)
	if err != nil {
		log.Fatalf("Failed to create sever %v", server)
	}
	server.ProcessMessages()
}
