package file

import (
	"fmt"
	"lamport_demo/io/network"
	"log"
	"os"
	"strings"
)

func WriteMessagesToFile(file string, messages []network.Message) {
	f, err := os.Create(file)
	defer f.Close()
	if err != nil {
		log.Fatalf("Couldn't open file!")
	}
	f.WriteString(fmt.Sprintf("Wall Clock, Logical Clock, Text, Sender Id\n"))
	for _, val := range messages {
		val.Text = strings.Replace(val.Text, "\n", "", -1)
		clockTime := fmt.Sprintf("%v:%v:%v", val.WallClock.Hour(), val.WallClock.Minute(), val.WallClock.Second())
		f.WriteString(fmt.Sprintf("%v,%v,%v,%v\n", clockTime, val.LogicalClock, val.Text, val.SenderId))
	}
}
