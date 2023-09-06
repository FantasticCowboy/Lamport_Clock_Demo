package main

import (
	"lamport_demo/client"
	"lamport_demo/constants"

	"log"
	"os"
)

func main() {

	if len(os.Args) < 3 {
		log.Fatalf("Did not pass in an ip address!")
	}
	serverIp := os.Args[1]
	clientIp := os.Args[2]

	client, err := client.CreateClient(serverIp, constants.SERVER_PORT, clientIp)
	if err != nil {
		log.Fatalf("Could not start client : %v", err)
	}
	client.StartClient(constants.CLIENT_PORT)
	client.SendFromStdIn()
}
