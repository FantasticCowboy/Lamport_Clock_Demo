package network

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)

type Connection struct {
	address string
}

type Message struct {
	Text         string
	WallClock    time.Time
	SenderId     int
	LogicalClock int
}

func CreateConnection(ip string, port string) (Connection, error) {
	connection := Connection{}
	connection.address = fmt.Sprintf("%s:%s", ip, port)
	return connection, nil
}

func (connection *Connection) SendMessage(message *Message) error {
	log.Printf("Sending Message: %+v", message)
	conn, err := net.Dial("tcp", connection.address)
	if err != nil {
		log.Printf("Could not setup connection correctly: %s", err.Error())
		return err
	}
	defer conn.Close()
	err = gob.NewEncoder(conn).Encode(message)
	if err != nil {
		log.Printf("Could not encode correctly: %s", err.Error())
		return err
	}
	log.Printf("Message sent!")

	return nil
}

func handleNewConnection(conn net.Conn, output chan Message) {

	msg := Message{}
	err := gob.NewDecoder(conn).Decode(&msg)

	if err != nil {
		log.Printf("Could not decode message: %s", err.Error())
		return
	}

	output <- msg
}

func StartListening(ip string, port string) (chan Message, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return nil, err
	}
	newMessages := make(chan Message)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Printf("Error: %v", err)
				continue
			}
			go handleNewConnection(conn, newMessages)
		}
	}()
	log.Printf("Starting to listen at %s:%s", ip, port)
	return newMessages, nil
}
