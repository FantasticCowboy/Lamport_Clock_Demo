package test

import (
	"lamport_demo/client"
	"lamport_demo/constants"
	"lamport_demo/server"
	"log"
	"testing"
)

func TestBasic(t *testing.T) {
	ip := "localhost"
	srv, err := server.CreateNewServer(ip)
	if err != nil {
		log.Fatalf("Error : %v", err)
	}

	client, err := client.CreateClient(ip, constants.SERVER_PORT, ip)
	client.StartClient(constants.CLIENT_PORT)

	if err != nil {
		log.Fatalf("Error : %v", err)
	}
	go srv.ProcessMessages()

	client.SendMessage("Hello")
	client.SendMessage("Hello")

	for len(srv.DumpMessagesReceived()) != 2 {
	}

}
