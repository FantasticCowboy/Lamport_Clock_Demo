package main

import (
	"bufio"
	"fmt"
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

	conn, err := network.CreateConnection(ip, constants.SERVER_PORT)

	if err != nil {
		log.Fatalf("Error connecting to server: %s", err.Error())
	}

	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter message to send: ")
			text, _ := reader.ReadString('\n')
			conn.SendMessage(text)
		}
	}()

	go func() {
		incomingMessages, err := network.StartListening(ip, constants.CLIENT_PORT)
		if err != nil {
			log.Panicf("Error starting to listen %v", err)
		}
		for {
			msg := <-incomingMessages
			log.Printf("%s", msg.Msg)
		}
	}()

	for {
	}
}
