package main

import (
	"lamport_demo/client"
	"lamport_demo/constants"
	"lamport_demo/network"

	"log"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatalf("Did not pass in an ip address!")
	}
	ip := os.Args[1]

	client, err := client.CreateClient(ip, constants.SERVER_PORT, network.GetLocalIP())
	if err != nil {
		log.Fatalf("Could not start client : %v", err)
	}
	client.StartClient(constants.CLIENT_PORT)
	client.SendFromStdIn()
}
