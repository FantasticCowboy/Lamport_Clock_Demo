package client

import (
	"bufio"
	"fmt"
	"lamport_demo/network"
	"log"
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
			client.updateClockValue(msg.LamportClock)
		}
	}()
}

func (client *Client) SendMessage(text string) {
	client.connection.SendMessage(
		&network.Message{
			Msg:             text,
			LamportClock:    client.getClockValueAndIncrement(),
			WallClock:       time.Now(),
			SenderIpAddress: client.ip,
		},
	)
}

func (client *Client) SendFromStdIn() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter message to send: ")
		text, _ := reader.ReadString('\n')
		client.SendMessage(text)
	}
}
