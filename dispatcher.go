package gomebus

import (
	"fmt"

	"net"
)

type Dispatcher struct {
	network string
	address string

	Codec

	listener   net.Listener
	channelMap *ChannelMap
	is_close   bool
}

func NewDispatcher(network, address string) (*Dispatcher, error) {

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, fmt.Errorf("listener error: %w", err)
	}
	return &Dispatcher{network: network, address: address, listener: ln, channelMap: GetChannelMap(), is_close: false}, nil
}

func (li *Dispatcher) newConnHandler(conn net.Conn) {

	buff := make([]byte, HEADER_SIZE)
	conn.Read(buff)

	msg, err := li.DecodeHeader(buff)
	if err != nil {
		conn.Close()
	}

	if msg.Message_size == uint32(HEADER_SIZE) && msg.Version == 0 && msg.Message_type == REGISTER {

		ch := NewChannel(conn, msg.Address)
		li.channelMap.Add(msg.Address, ch)
		ch.Open()

	} else {
		conn.Close()
	}
}

func (li *Dispatcher) Accept() error {

	for !li.is_close {
		conn, err := li.listener.Accept()
		if err != nil {
			return fmt.Errorf("Accept error: %w", err)
		}
		li.newConnHandler(conn)
	}
	return nil
}

func (li *Dispatcher) Close() error {
	li.is_close = true
	for k := range li.channelMap.channels {
		delete(li.channelMap.channels, k)
	}
	return li.listener.Close()
}

func (li *Dispatcher) Addr() net.Addr {
	return li.listener.Addr()
}

func (li *Dispatcher) ActiveChannels() int {
	return li.channelMap.Count()
}
