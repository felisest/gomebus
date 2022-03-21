package gomebus_tests

import (
	"fmt"
	"testing"

	"felis.est/gomebus"
)

func TestEncode(t *testing.T) {

	template_packet := gomebus.Message{Message_size: 34, Version: 0, Message_type: 2, Address: "test.client", Payload: []byte{'t', 'e', 's', 't', ' ', 'p', 'a', 'y', 'l', 'o', 'a', 'd'}}

	template_buff := []byte{0x0, 0x0, 0x0, 0x22, 										//size
		0x0,                                                                            //version
		0x2,                                                                            //packet type = register
		't', 'e', 's', 't', '.', 'c', 'l', 'i', 'e', 'n', 't', 0x0, 0x0, 0x0, 0x0, 0x0, //address - 16 byte
		't', 'e', 's', 't', ' ', 'p', 'a', 'y', 'l', 'o', 'a', 'd', 					//payload
	}

	codec := gomebus.Codec{}

	out_buff, err := codec.Encode(&template_packet, "test.client")

	if err != nil {
		t.Errorf("Encode error: %s", err.Error())
	}

	err = nil

	for i, b := range template_buff {
		if b != out_buff[i] {
			err = fmt.Errorf("buffers not equial, pos %d, %c != %c", i, b, out_buff[i])
			break
		}
	}

	if err != nil {
		t.Errorf("Out buffer want: %s, got: \n%s, \nerr: \n%s", template_buff, out_buff, err)
	}
}

func TestDecode(t *testing.T) {

	template_buff := []byte{0x0, 0x0, 0x0, 0x22, 										//size
		0x0,                                                                            //version
		0x2,                                                                            //packet type = register
		't', 'e', 's', 't', '.', 'c', 'l', 'i', 'e', 'n', 't', 0x0, 0x0, 0x0, 0x0, 0x0, //address - 16 byte
		't', 'e', 's', 't', ' ', 'p', 'a', 'y', 'l', 'o', 'a', 'd', 					//payload
	}

	codec := gomebus.Codec{}

	out_packet, err := codec.Decode(template_buff)

	if err != nil {
		t.Errorf("Decode error: %s", err.Error())
	}

	fmt.Println(out_packet.Address)

	if out_packet.Message_size != 34 {
		t.Errorf("Decode size error want: %d, got: %d\n", 34, out_packet.Message_size)
	}
	if out_packet.Version != 0 {
		t.Errorf("Decode version error want: %d, got: %d\n", 0, out_packet.Version)
	}
	if out_packet.Message_type != 2 {
		t.Errorf("Decode type error want: %d, got: %d\n", 2, out_packet.Message_type)
	}
	if out_packet.Address != "test.client" {
		t.Errorf("Decode address error want: 'test.client', got: '%s' \n", out_packet.Address)
	}
	if len(out_packet.Payload) != 12 {
		t.Errorf("Payload len error want: %d, got: %d \n", 12, len(out_packet.Payload))
	}
}
