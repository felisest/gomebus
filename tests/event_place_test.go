package gomebus_tests

import (
	"testing"
	"time"

	"felis.est/gomebus"
)

func TestAddGetEvent(t *testing.T) {

	em := gomebus.GetEventMap()

	ch := gomebus.NewChannel(nil, "test.channel")
	ch2 := gomebus.NewChannel(nil, "test.channel2")

	em.Add("test.eventplace", "test.channel", ch)
	em.Add("test.eventplace", "test.channel2", ch2)

	cnt := em.Count()

	if cnt != 1 {
		t.Errorf("Map count want 1, got %d ", cnt)
	}

	listeners := em.Get("test.eventplace")

	if len(listeners) != 2 {
		t.Errorf("Connection count want 2, got %d ", len(listeners))
	}

	em = nil
}

func TestSendEvent(t *testing.T) {

	payload := []byte("NEW EVENT")

	ln, err := gomebus.NewDispatcher("tcp", ":8080")
	sender_client := gomebus.RemoteClient{}
	receiver_client := gomebus.RemoteClient{}

	if ln == nil {
		t.Error("Empty listener")
	}
	if err != nil {
		t.Error("Listener error: %w", err)
	}

	go ln.Accept()

	sender_client.Register("tcp", "localhost:8080", "sender.client")
	receiver_client.Register("tcp", "localhost:8080", "receiver.client")

	sender_client.Subscribe("event.place1")
	receiver_client.Subscribe("event.place1")

	go sender_client.SendEvent("event.place1", payload)

	go receiver_client.Receive(func(msg *gomebus.Message) {
		if msg.Message_size != 31 {
			t.Errorf("Transmission error, size: '%d', got '%d'", 31, msg.Message_size)
		}
		if msg.Version != 0 {
			t.Errorf("Transmission error, version: '%d', got '%d'", 0, msg.Version)
		}
		if msg.Message_type != gomebus.EVENT {
			t.Errorf("Transmission error, size: '%d', got '%d'", gomebus.EVENT, msg.Message_type)
		}
		if msg.Address != "event.place1" {
			t.Errorf("Transmission error, src: '%s', got '%s'", "event.place1", msg.Address)
		}
		if string(msg.Payload) != string(payload) {
			t.Errorf("Transmission error, send: '%s', got '%s'", "NEW EVENT", string(msg.Payload))
		}
	})

	go sender_client.Receive(func(msg *gomebus.Message) {
		if msg.Message_size != 31 {
			t.Errorf("Transmission error, size: '%d', got '%d'", 31, msg.Message_size)
		}
		if msg.Version != 0 {
			t.Errorf("Transmission error, version: '%d', got '%d'", 0, msg.Version)
		}
		if msg.Message_type != gomebus.EVENT {
			t.Errorf("Transmission error, size: '%d', got '%d'", gomebus.EVENT, msg.Message_type)
		}
		if msg.Address != "event.place1" {
			t.Errorf("Transmission error, src: '%s', got '%s'", "event.place1", msg.Address)
		}
		if string(msg.Payload) != string(payload) {
			t.Errorf("Transmission error, send: '%s', got '%s'", "NEW EVENT", string(msg.Payload))
		}
	})

	time.Sleep(1000 * time.Millisecond)

	sender_client.Close()
	receiver_client.Close()
	ln.Close()
}
