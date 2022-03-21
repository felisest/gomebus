package gomebus

const (
	PING uint8 = iota
	PONG
	ERR
	REGISTER
	SEND
	REPLY
	SUBSCRIBE
	EVENT
)

const HEADER_SIZE uint16 = 22

type Message struct {
	Message_size uint32
	Version      uint8
	Message_type uint8
	Address      string
	Payload      []byte
}

type Slide struct {
	Msg        *Message
	SrcAddress string
	DstAddress string
}
