package gomebus

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Channel struct {
	conn     net.Conn
	chan_map *ChannelMap
	event_place *EventMap
	Codec
	address string

	receiver chan *Slide
}

const TIMEOUT time.Duration = 1000

func NewChannel(conn net.Conn, address string) *Channel {
	c := &Channel{
		conn:     conn,
		chan_map: GetChannelMap(),
		event_place: GetEventMap(),
		address:  address,

		receiver: make(chan *Slide),
	}
	return c
}

func (c *Channel) Open() {

	go c.reader()
	go c.writer()

	c.conn.SetDeadline(time.Now().Add(time.Millisecond * TIMEOUT))
}

func (c *Channel) Close() {
	fmt.Printf("channel %s is closing\n", c.address)
	close(c.receiver)
	c.conn.Close()
	delete(c.chan_map.channels, c.address)
}

func (c *Channel) fail(err error) error {

	error_buf := []byte(err.Error())

	msg := &Message{Message_size: uint32(HEADER_SIZE), Version: 0, Message_type: ERR, Payload: error_buf}
	c.receiver <- &Slide{Msg: msg}

	return nil
}

func (c *Channel) pong() error {

	msg := &Message{Message_size: uint32(HEADER_SIZE), Version: 0, Message_type: PONG, Address: c.address}
	c.receiver <- &Slide{Msg: msg}

	return nil
}

func (c *Channel) send(msg *Message) error {

	dst_receiver := c.chan_map.Get(msg.Address)
	if dst_receiver == nil {
		return fmt.Errorf("endpoint not found")
	}
	sld := &Slide{Msg: msg, DstAddress: msg.Address, SrcAddress: c.address}
	dst_receiver.receiver <- sld

	return nil
}

func (c *Channel) reply(msg *Message) error {

	dst_receiver := c.chan_map.Get(msg.Address)
	if dst_receiver == nil {
		return fmt.Errorf("endpoint not found")
	}
	sld := &Slide{Msg: msg, DstAddress: msg.Address, SrcAddress: c.address}
	dst_receiver.receiver <- sld

	return nil
}

func (c *Channel) suscribe(msg *Message) error {

	return c.event_place.Add(msg.Address, c.address, c)
}

func (c *Channel) send_event(msg *Message) error {

	return c.event_place.Send(msg.Address, c.address, msg.Payload)
}

func (c *Channel) checkMessage(msg *Message) error {

	return nil
}

func (c *Channel) packetHandler(msg *Message) error {

	if err := c.checkMessage(msg); err != nil {
		return err
	}

	switch msg.Message_type {
	case PING:
		if err := c.pong(); err != nil {
			return err
		}
	case SEND:
		if err := c.send(msg); err != nil {
			return err
		}
	case REPLY:
		if err := c.reply(msg); err != nil {
			return err
		}
	case SUBSCRIBE:
		if err := c.suscribe(msg); err != nil {
			return err
		}		
	case EVENT:
		if err := c.send_event(msg); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected packet type")
	}

	return nil
}

func (c *Channel) reader() {

	header_buff := make([]byte, HEADER_SIZE)

	reader := bufio.NewReader(c.conn)

	for {

		if _, err := reader.Read(header_buff); err != nil {
			fmt.Printf("%s read error: %s\n", c.address, err.Error())
			c.Close()
			break
		}
		c.conn.SetDeadline(time.Now().Add(time.Millisecond * TIMEOUT))

		msg, err := c.DecodeHeader(header_buff)
		if err != nil {
			c.fail(fmt.Errorf("wrong header %w", err))
		}

		if msg.Message_size > uint32(HEADER_SIZE) {

			payload_buff := make([]byte, msg.Message_size-uint32(HEADER_SIZE))

			if _, err := reader.Read(payload_buff); err != nil {
				c.fail(fmt.Errorf("payload read error: %w", err))
			}
			msg.Payload = payload_buff
		}

		if err = c.packetHandler(msg); err != nil {
			c.fail(err)
		}
	}
}

func (c *Channel) writer() {

	buff_writer := bufio.NewWriter(c.conn)

	for sld := range c.receiver {

		buff_out, err := c.Encode(sld.Msg, sld.SrcAddress)
		if err != nil {
			_ = fmt.Errorf("encode error: %w", err)
		}

		if _, err := buff_writer.Write(buff_out); err != nil {
			_ = fmt.Errorf("write error: %w", err)
		}
		c.conn.SetDeadline(time.Now().Add(time.Millisecond * TIMEOUT))

		buff_writer.Flush()
	}
}
