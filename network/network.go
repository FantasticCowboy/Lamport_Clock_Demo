package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type ClientConnection struct {
	address string
}

type Message struct {
	Msg          string
	LamportClock int
	WallClock    time.Time
}

func CreateConnection(ip string, port string) (ClientConnection, error) {
	connection := ClientConnection{}
	connection.address = fmt.Sprintf("%s:%s", ip, port)
	return connection, nil
}

func (connection *ClientConnection) SendMessage(message string) error {
	log.Printf("Starting Send Message: %s", message)

	conn, err := net.Dial("tcp", connection.address)
	defer conn.Close()
	if err != nil {
		log.Printf("Could not setup connection correctly: %s", err.Error())
		return err
	}
	msg := Message{
		Msg:          message,
		LamportClock: 0,
		WallClock:    time.Now(),
	}
	err = gob.NewEncoder(conn).Encode(msg)
	if err != nil {
		log.Printf("Could not encode correctly: %s", err.Error())
		return err
	}
	log.Printf("Ending Send Message: %s", message)
	return nil
}

func readBytesFromConnection(conn net.Conn) []byte {
	defer conn.Close()
	bytesReceived := make([]byte, 0)
	for {
		packet := make([]byte, 1024)
		_, err := conn.Read(packet)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error reading bytes from connection %s", err.Error())
			return nil
		}
		bytesReceived = append(bytesReceived, packet...)
	}
	return bytesReceived
}

func handleNewConnection(conn net.Conn, output chan Message) {
	bytesReceived := readBytesFromConnection(conn)
	if bytesReceived == nil {
		return
	}

	msg := Message{}
	err := gob.NewDecoder(bytes.NewBuffer(bytesReceived)).Decode(&msg)

	if err != nil {
		log.Printf("Could not decode message: %s", err.Error())
		return
	}

	output <- msg
}

func StartListening(port string) (chan Message, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return nil, err
	}
	newMessages := make(chan Message)
	go func() {
		for {
			conn, err := ln.Accept()
			log.Printf("Received new connection!")
			if err != nil {
				continue
			}
			go handleNewConnection(conn, newMessages)
		}
	}()
	log.Printf("Starting to listen at %s", port)
	return newMessages, nil
}
