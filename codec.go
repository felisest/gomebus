package gomebus

import (
	"encoding/binary"
	"fmt"
)

type Codec struct {
}

func (c *Codec) Decode(raw_message []byte) (*Message, error) {

	msg := &Message{}

	msg.Message_size = binary.BigEndian.Uint32(raw_message[0:4])
	msg.Version = uint8(raw_message[4])
	msg.Message_type = uint8(raw_message[5])

	var str_end int = 0
	for i, v := range raw_message[6:22] {
		if v == 0 {
			str_end = i
			break
		}
	}
	msg.Address = string(raw_message[6:str_end + 6])

	msg.Payload = raw_message[22:]

	return msg, nil
}

func (c *Codec) DecodeHeader(raw_message []byte) (*Message, error) {

	msg := &Message{}

	msg.Message_size = binary.BigEndian.Uint32(raw_message[0:4])
	msg.Version = uint8(raw_message[4])
	msg.Message_type = uint8(raw_message[5])

	var str_end int = 0
	for i, v := range raw_message[6:22] {
		if v == 0 {
			str_end = i
			break
		}
	}
	msg.Address = string(raw_message[6:6+str_end])

	return msg, nil
}

func (c *Codec) Encode(pkt *Message, adrs string) ([]byte, error) {

	size := len(pkt.Payload) + 22

	out_buff := make([]byte, len(pkt.Payload)+22)

	binary.BigEndian.PutUint32(out_buff, uint32(size))
	out_buff[4] = pkt.Version
	out_buff[5] = pkt.Message_type

	buff := []byte(adrs)
	copy(out_buff[6:], buff)

	if pkt.Payload == nil || len(pkt.Payload) == 0 {
		return out_buff, nil	
	}
	copy(out_buff[22:], pkt.Payload)

	return out_buff, nil
}

func (c *Codec) Wrap(pkt_type uint8, pkt_addr string, payload []byte) ([]byte, error) {

	size := len(payload) + 22
	buff := make([]byte, size)

	binary.BigEndian.PutUint32(buff, uint32(size))
	buff[4] = 0        //version
	buff[5] = pkt_type //packet type

	if len(pkt_addr) > 16 {
		return nil, fmt.Errorf("length of the address cannot be more than 16")
	}

	addr_buff := []byte(pkt_addr)
	copy(buff[6:], addr_buff)

	if payload == nil || len(payload) == 0 {
		return buff, nil		
	} 
	copy(buff[22:], payload)

	return buff, nil
}
