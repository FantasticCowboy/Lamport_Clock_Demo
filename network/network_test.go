package network

import "testing"

func TestBasic(t *testing.T) {
	newMessages, err := StartListening("", "8080")
	if err != nil {
		t.Errorf("Did not start listening correctly: %d", err)
	}
	connection, err := CreateConnection("", "8080")
	if err != nil {
		t.Errorf("Did not establish a connection to the server successfully : %d", err)
	}
	connection.SendMessage(Message{Msg: "hello!"})
	msg := <-newMessages

	if msg.Msg != "Hello!" {
		t.Error("message received does not equal Hello")
	}
}
