package gomebus_tests

import (
	"testing"
	"time"

	"felis.est/gomebus"
)

func TestAddChannel(t *testing.T) {

	ln, err := gomebus.NewDispatcher("tcp", ":8080")
	client := gomebus.RemoteClient{}

	if ln == nil {
		t.Error("Empty listener")
	} else {

		if err != nil {
			t.Error("Listener error: " + err.Error())
		}

		client.Register("tcp", "localhost:8080", "sender.client")

		go ln.Accept()

		time.Sleep(300 * time.Millisecond)

		cnt := ln.ActiveChannels()

		if cnt != 1 {
			t.Errorf("Map Count want 1, got %d ", cnt)
		}
	}
	ln.Close()
}

func TestAddFewChannels(t *testing.T) {

	ln, err := gomebus.NewDispatcher("tcp", ":8080")
	client1 := gomebus.RemoteClient{}
	client2 := gomebus.RemoteClient{}
	client3 := gomebus.RemoteClient{}

	if ln == nil {
		t.Error("Empty listener")
	} else {

		if err != nil {
			t.Error("Listener error: " + err.Error())
		}

		client1.Register("tcp", "localhost:8080", "client1")
		client2.Register("tcp", "localhost:8080", "client2")
		client3.Register("tcp", "localhost:8080", "client3")

		go ln.Accept()

		time.Sleep(100 * time.Millisecond)

		cnt := ln.ActiveChannels()

		if cnt != 3 {
			t.Errorf("Map Count want 3, got %d ", cnt)
		}
	}
	ln.Close()
}
