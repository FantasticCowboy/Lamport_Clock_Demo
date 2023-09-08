package file

import (
	"fmt"
	"lamport_demo/io/network"
	"log"
	"os"
)

func WriteMessagesToFile(file string, messages []network.Message) {
	f, err := os.Create("./sorted_by_lamport_clock.txt")
	defer f.Close()
	if err != nil {
		log.Fatalf("Couldn't open file!")
	}
	f.WriteString(fmt.Sprintf("Sender Id, Wall Clock, Lamport Clock, Text\n"))
	for _, val := range messages {
		f.WriteString(fmt.Sprintf("%d,%v,%v,%s", val.SenderId, val.WallClock, val.LamportClock, val.Text))
	}
}
