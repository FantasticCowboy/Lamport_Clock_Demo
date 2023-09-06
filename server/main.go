package main

import (
	"lamport_demo/constants"
	"lamport_demo/network"
	"log"
)

func pingClient(ip string) {
	conn, err := network.CreateConnection(ip, constants.CLIENT_PORT)
	if err != nil {
		return
	}
	conn.SendMessage("Heartbeat")
}

func main() {
	ip := network.GetLocalIP()

	newMessages, err := network.StartListening(ip, constants.SERVER_PORT)
	if err != nil {
		log.Panicf("Error starting server: %s", err.Error())
	}
	log.Printf("Starting server at %s:%s", ip, constants.SERVER_PORT)

	for {
		message := <-newMessages
		log.Printf("Recieved message %v", message)
		go pingClient(message.SenderIpAddress)
	}
}
