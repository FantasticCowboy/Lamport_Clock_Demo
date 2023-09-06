package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

type Connection struct {
	address string
}

type Message struct {
	Msg             string
	LamportClock    int
	WallClock       time.Time
	SenderIpAddress string
}

func CreateConnection(ip string, port string) (Connection, error) {
	connection := Connection{}
	connection.address = fmt.Sprintf("%s:%s", ip, port)
	return connection, nil
}

func (connection *Connection) SendMessage(message *Message) error {
	log.Printf("Starting Send Message: %v to %v", message, connection.address)
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
	log.Printf("Ending Send Message: %v", message)
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
	msg.SenderIpAddress = strings.Split(conn.RemoteAddr().String(), ":")[0]

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
			log.Printf("Received new connection!")
			go handleNewConnection(conn, newMessages)
		}
	}()
	log.Printf("Starting to listen at %s:%s", ip, port)
	return newMessages, nil
}
