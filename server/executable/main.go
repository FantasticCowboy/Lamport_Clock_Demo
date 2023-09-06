package main

import (
	"lamport_demo/network"
	"lamport_demo/server"
	"log"
)

func main() {
	ip := network.GetLocalIP()
	server, err := server.CreateNewServer(ip)
	if err != nil {
		log.Fatalf("Failed to create sever %v", server)
	}
	server.ProcessMessages()
}
