package gomebus_tests

import (
	"testing"
	"time"

	"felis.est/gomebus"
)

func TestHB(t *testing.T) {

	ln, err := gomebus.NewDispatcher("tcp", ":8080")
	sender_client := gomebus.RemoteClient{}

	if ln == nil {
		t.Error("Empty listener")
	}
	if err != nil {
		t.Error("Listener error: %w", err)
	}

	go ln.Accept()

	sender_client.Register("tcp", "localhost:8080", "sender.client")

	sender_client.RunHB()

	var counter int = 0

	go sender_client.Receive(func(msg *gomebus.Message) {
		if msg.Message_size != 22 {
			t.Errorf("Transmission error, size: '%d', got '%d'", 41, msg.Message_size)
		}
		if msg.Version != 0 {
			t.Errorf("Transmission error, version: '%d', got '%d'", 0, msg.Version)
		}
		if msg.Message_type != gomebus.PONG {
			t.Errorf("Transmission error, size: '%d', got '%d'", gomebus.SEND, msg.Message_type)
		}
		counter++
	})

	time.Sleep(3000 * time.Millisecond)

	if counter != 3 {
		t.Errorf("HB count error, want: '%d', got '%d'", 3, counter)
	}

	sender_client.Close()
	ln.Close()
}

func TestTransmitMessage(t *testing.T) {

	payload := []byte("PAYLOAD from sender")

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

	go sender_client.Send(gomebus.SEND, "receiver.client", payload)

	go receiver_client.Receive(func(msg *gomebus.Message) {
		if msg.Message_size != 41 {
			t.Errorf("Transmission error, size: '%d', got '%d'", 41, msg.Message_size)
		}
		if msg.Version != 0 {
			t.Errorf("Transmission error, version: '%d', got '%d'", 0, msg.Version)
		}
		if msg.Message_type != gomebus.SEND {
			t.Errorf("Transmission error, size: '%d', got '%d'", gomebus.SEND, msg.Message_type)
		}
		if msg.Address != "sender.client" {
			t.Errorf("Transmission error, src: '%s', got '%s'", "sender.client", msg.Address)
		}
		if string(msg.Payload) != string(payload) {
			t.Errorf("Transmission error, send: '%s', got '%s'", "PAYLOAD from sender", string(msg.Payload))
		}
	})
	time.Sleep(1000 * time.Millisecond)

	sender_client.Close()
	receiver_client.Close()
	ln.Close()
}

func TestTransmitReplyMessage(t *testing.T) {

	payload := []byte("PAYLOAD from sender")
	reply_payload := []byte("REPLY from receiver")

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

	go sender_client.Send(gomebus.SEND, "receiver.client", payload)

	go receiver_client.Receive(func(msg *gomebus.Message) {
		if msg.Message_size != 41 {
			t.Errorf("Transmission error, size: '%d', got '%d'", 41, msg.Message_size)
		}
		if msg.Version != 0 {
			t.Errorf("Transmission error, version: '%d', got '%d'", 0, msg.Version)
		}
		if msg.Message_type != gomebus.SEND {
			t.Errorf("Transmission error, size: '%d', got '%d'", gomebus.SEND, msg.Message_type)
		}
		if msg.Address != "sender.client" {
			t.Errorf("Transmission error, src: '%s', got '%s'", "sender.client", msg.Address)
		}
		if string(msg.Payload) != string(payload) {
			t.Errorf("Transmission error, send: '%s', got '%s'", string(payload), string(msg.Payload))
		}
	})

	go receiver_client.Send(gomebus.REPLY, "sender.client", reply_payload)

	go sender_client.Receive(func(msg *gomebus.Message) {
		if msg.Message_size != 41 {
			t.Errorf("Transmission error, size: '%d', got '%d'", 41, msg.Message_size)
		}
		if msg.Version != 0 {
			t.Errorf("Transmission error, version: '%d', got '%d'", 0, msg.Version)
		}
		if msg.Message_type != gomebus.SEND {
			t.Errorf("Transmission error, size: '%d', got '%d'", gomebus.REPLY, msg.Message_type)
		}
		if msg.Address != "receiver.client" {
			t.Errorf("Transmission error, src: '%s', got '%s'", "receiver.client", msg.Address)
		}
		if string(msg.Payload) != string(reply_payload) {
			t.Errorf("Transmission error, send: '%s', got '%s'", string(reply_payload), string(msg.Payload))
		}
	})

	time.Sleep(1000 * time.Millisecond)

	sender_client.Close()
	receiver_client.Close()
	ln.Close()
}
