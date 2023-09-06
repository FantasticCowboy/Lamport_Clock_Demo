package server

import (
	"lamport_demo/constants"
	"lamport_demo/network"
	"log"
	"sync"
	"time"
)

type Server struct {
	maxClockValueSeen int
	messagesReceived  []network.Message
	incomingMessages  chan network.Message
	lock              *sync.Mutex
}

func (srv *Server) getMaxClockValueSeen() int {
	return srv.maxClockValueSeen
}

func (srv *Server) updateMaxClockValueSeen(val int) {
	if srv.maxClockValueSeen < val {
		srv.maxClockValueSeen = val
	}
}

func (srv *Server) PingClient(ip string, clock int) {
	conn, err := network.CreateConnection(ip, constants.CLIENT_PORT)
	if err != nil {
		return
	}
	msg := network.Message{
		Msg:             "",
		LamportClock:    clock,
		WallClock:       time.Now(),
		SenderIpAddress: "",
	}
	conn.SendMessage(msg)
}

func CreateNewServer() (Server, error) {
	srv := Server{}
	ip := network.GetLocalIP()
	newMessages, err := network.StartListening(ip, constants.SERVER_PORT)
	if err != nil {
		log.Printf("Error starting server: %s", err.Error())
		return srv, err
	}
	srv.maxClockValueSeen = 0
	srv.incomingMessages = newMessages
	srv.messagesReceived = make([]network.Message, 0)
	srv.lock = new(sync.Mutex)
	log.Printf("Server created at %s:%s", ip, constants.SERVER_PORT)
	return srv, err
}

func (srv *Server) ProcessMessages() {
	log.Printf("Server processing messages")
	for {
		message := <-srv.incomingMessages
		srv.updateMaxClockValueSeen(message.LamportClock)
		log.Printf("Received message!")
		go srv.PingClient(message.SenderIpAddress, srv.getMaxClockValueSeen())
		srv.messagesReceived = append(srv.messagesReceived, message)
	}
}
