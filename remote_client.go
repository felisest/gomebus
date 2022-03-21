package gomebus

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type RemoteClient struct {
	conn     net.Conn
	is_close bool
	Codec

	client_name string
}

func (c *RemoteClient) RunHB() error {

	buff := []byte{0x0, 0x0, 0x0, byte(HEADER_SIZE), //size
		0x0,                                                                            //version
		0x0,                                                                            //packet type = register
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, //address - 16 byte
	}

	go func() {
		for {
			time.Sleep(990 * time.Millisecond)
			if _, err := c.conn.Write(buff); err != nil {
				break
			}
		}
	}()

	return nil
}

func (c *RemoteClient) Register(network string, address string, clent_name string) error {

	c.is_close = false

	var err error
	c.conn, err = net.Dial(network, address)
	if err != nil {
		return fmt.Errorf("dial error: %w", err)
	}

	c.client_name = clent_name

	size := HEADER_SIZE
	buff := make([]byte, size)

	binary.BigEndian.PutUint32(buff, uint32(size))
	buff[4] = 0        //version
	buff[5] = REGISTER //packet type

	if len(clent_name) > 16 {
		return fmt.Errorf("length of the address cannot be more than 16")
	}

	register_buff := []byte(clent_name)
	copy(buff[6:], register_buff)

	_, err = c.conn.Write(buff)
	if err != nil {
		return fmt.Errorf("write error: %w", err)
	}
	return nil
}

func (c *RemoteClient) Subscribe(evnt string) error {
	buff, err := c.Wrap(SUBSCRIBE, evnt, nil)
	if err != nil {
		return fmt.Errorf("codec wrap error: %w", err)
	}

	_, err = c.conn.Write(buff)
	if err != nil {
		return fmt.Errorf("write error: %w", err)
	}
	return nil
}

func (c *RemoteClient) SendEvent(dst string, msg_payload []byte) error {

	buff, err := c.Wrap(EVENT, dst, msg_payload)
	if err != nil {
		return fmt.Errorf("codec wrap error: %w", err)
	}

	_, err = c.conn.Write(buff)
	if err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	return nil
}

func (c *RemoteClient) Send(msg_type uint8, dst string, msg_payload []byte) ([]byte, error) {

	pkt_buff, err := c.Wrap(SEND, dst, msg_payload)
	if err != nil {
		return nil, fmt.Errorf("codec wrap error: %w", err)
	}

	_, err = c.conn.Write(pkt_buff)
	if err != nil {
		return nil, fmt.Errorf("write error: %w", err)
	}
	return nil, nil
}

func (c *RemoteClient) Receive(callback func(*Message)) error {

	reader := bufio.NewReader(c.conn)
	header_buff := make([]byte, HEADER_SIZE)

	for !c.is_close {

		if _, err := reader.Read(header_buff); err != nil {
			return fmt.Errorf("header read error: %w", err)
		}

		msg, err := c.DecodeHeader(header_buff)
		if err != nil {
			return fmt.Errorf("decode error: %w", err)
		}

		body_buff := make([]byte, msg.Message_size-uint32(HEADER_SIZE))

		if _, err := reader.Read(body_buff); err != nil {
			return fmt.Errorf("body read error: %w", err)
		}

		msg.Payload = body_buff

		callback(msg)
	}

	return nil
}

func (c *RemoteClient) Close() {
	c.is_close = true
}
