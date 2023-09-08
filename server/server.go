package server

import (
	"lamport_demo/constants"
	"lamport_demo/io/file"
	"lamport_demo/io/network"
	"log"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	maxClockValueSeen int
	messagesReceived  []network.Message
	incomingMessages  chan network.Message
	lock              *sync.Mutex
}

func (srv *Server) DumpMessagesReceived() []network.Message {
	return srv.messagesReceived
}

func (srv *Server) getMaxClockValueSeen() int {
	return srv.maxClockValueSeen
}

func (srv *Server) updateMaxClockValueSeen(val int) {
	if srv.maxClockValueSeen < val {
		srv.maxClockValueSeen = val
	}
	srv.maxClockValueSeen++
}

func (srv *Server) PingClient(ip string, clock int) {
	conn, err := network.CreateConnection(ip, constants.CLIENT_PORT)
	if err != nil {
		return
	}
	msg := network.Message{
		Text:      "",
		WallClock: time.Now(),
	}
	conn.SendMessage(&msg)
}

func CreateNewServer(ip string) (*Server, error) {
	srv := Server{}
	newMessages, err := network.StartListening(ip, constants.SERVER_PORT)
	if err != nil {
		log.Printf("Error starting server: %s", err.Error())
		return nil, err
	}
	srv.maxClockValueSeen = 0
	srv.incomingMessages = newMessages
	srv.messagesReceived = make([]network.Message, 0)
	srv.lock = new(sync.Mutex)
	log.Printf("Server created at %s:%s", ip, constants.SERVER_PORT)
	srv.registerTeardown()
	return &srv, err
}

func (srv *Server) ProcessMessages() {
	log.Printf("Server processing messages")
	for {
		log.Printf("Waiting for a message")
		message := <-srv.incomingMessages
		srv.lock.Lock()
		message.LogicalClock = srv.maxClockValueSeen
		srv.updateMaxClockValueSeen(0)
		log.Printf("Received message: %+v", message)
		srv.messagesReceived = append(srv.messagesReceived, message)
		srv.lock.Unlock()
	}
}

func (srv *Server) registerTeardown() {
	programDone := make(chan os.Signal, 1)
	signal.Notify(programDone, syscall.SIGTERM, syscall.SIGINT)

	go func(srv *Server) {
		<-programDone
		srv.lock.Lock()
		defer srv.lock.Unlock()
		srv.sortMessagesByLamportClock(srv.messagesReceived)
		file.WriteMessagesToFile("./sorted_by_logical_clock.txt", srv.messagesReceived)
		srv.sortMessagesByTimestamp(srv.messagesReceived)
		file.WriteMessagesToFile("./sorted_by_timestamp.txt", srv.messagesReceived)
		os.Exit(0)
	}(srv)
}

func (srv *Server) sortMessagesByTimestamp(messages []network.Message) {
	sort.Slice(srv.messagesReceived,
		func(i, j int) bool {
			return (srv.messagesReceived[i].WallClock.Before(srv.messagesReceived[j].WallClock))
		})
}

func (srv *Server) sortMessagesByLamportClock(messages []network.Message) {
	sort.Slice(srv.messagesReceived,
		func(i, j int) bool {
			return (srv.messagesReceived[i].LogicalClock < srv.messagesReceived[j].LogicalClock ||
				(srv.messagesReceived[i].LogicalClock == srv.messagesReceived[j].LogicalClock &&
					srv.messagesReceived[i].SenderId < srv.messagesReceived[j].SenderId))
		})
}
