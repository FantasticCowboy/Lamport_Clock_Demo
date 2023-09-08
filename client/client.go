package client

import (
	"bufio"
	"fmt"
	"lamport_demo/io/network"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Client struct {
	LamportClock int
	lock         *sync.Mutex
	connection   network.Connection
	ip           string
	id           int
}

func generateRandomNumber() int {
	s1 := rand.NewSource(time.Now().UnixNano())
	return rand.New(s1).Int()
}

func (client *Client) getClockValueAndIncrement() int {
	client.lock.Lock()
	defer client.lock.Unlock()
	client.LamportClock++
	return client.LamportClock
}

func (client *Client) updateClockValue(val int) {
	client.lock.Lock()
	defer client.lock.Unlock()
	if client.LamportClock < val {
		client.LamportClock = val
	}
}

func CreateClient(serverIp string, serverPort string, clientIP string) (Client, error) {
	client := Client{}
	conn, err := network.CreateConnection(serverIp, serverPort)
	if err != nil {
		return client, err
	}
	client.LamportClock = 0
	client.lock = new(sync.Mutex)
	client.connection = conn
	client.ip = clientIP
	client.id = generateRandomNumber()
	return client, nil
}

func (client *Client) StartClient(clientPort string) {

	go func() {
		incomingMessages, err := network.StartListening(client.ip, clientPort)
		if err != nil {
			log.Panicf("Error starting to listen %v", err)
		}
		for {
			msg := <-incomingMessages
			log.Printf("Received message: %+v", msg)
		}
	}()
}

func (client *Client) SendMessage(text string) {
	client.connection.SendMessage(
		&network.Message{
			Text:      text,
			WallClock: time.Now(),
			SenderId:  client.id,
		},
	)
}

func (client *Client) SendFromStdIn() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter message to send: \n")
		text, _ := reader.ReadString('\n')
		client.SendMessage(text)
	}
}
