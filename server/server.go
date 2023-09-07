package server

import (
	"fmt"
	"lamport_demo/constants"
	"lamport_demo/network"
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

	programDone := make(chan os.Signal, 1)
	signal.Notify(programDone, syscall.SIGTERM, syscall.SIGINT)

	go func(srv *Server) {
		<-programDone
		srv.lock.Lock()
		defer srv.lock.Unlock()
		sort.Slice(srv.messagesReceived,
			func(i, j int) bool {
				return (srv.messagesReceived[i].LamportClock < srv.messagesReceived[j].LamportClock ||
					(srv.messagesReceived[i].LamportClock == srv.messagesReceived[j].LamportClock &&
						srv.messagesReceived[i].SenderId < srv.messagesReceived[j].SenderId))
			})
		f, err := os.Create("./sorted_by_lamport_clock.txt")
		if err != nil {
			log.Fatalf("Couldn't open file!")
		}

		f.WriteString(fmt.Sprintf("Sender Id, wall clock, lamport clock, message"))
		for _, val := range srv.messagesReceived {
			f.WriteString(fmt.Sprintf("%d,%v,%v,%s", val.SenderId, val.WallClock, val.LamportClock, val.Msg))
		}
		f.Close()

		sort.Slice(srv.messagesReceived,
			func(i, j int) bool {
				return (srv.messagesReceived[i].WallClock.Before(srv.messagesReceived[j].WallClock))
			})

		f, err = os.Create("./sorted_by_timestamp.txt")
		if err != nil {
			log.Fatalf("Couldn't open file!")
		}
		f.WriteString(fmt.Sprintf("Sender Id, wall clock, lamport clock, message"))
		for _, val := range srv.messagesReceived {
			f.WriteString(fmt.Sprintf("%d,%v,%v,%s", val.SenderId, val.WallClock, val.LamportClock, val.Msg))
		}
		f.Close()
		os.Exit(0)
	}(&srv)

	return &srv, err
}

func (srv *Server) ProcessMessages() {
	log.Printf("Server processing messages")
	for {
		message := <-srv.incomingMessages
		srv.updateMaxClockValueSeen(message.LamportClock)
		log.Printf("Received message: %v", message)
		go srv.PingClient(message.SenderIpAddress, srv.getMaxClockValueSeen())
		srv.lock.Lock()
		srv.messagesReceived = append(srv.messagesReceived, message)
		srv.lock.Unlock()
	}
}
