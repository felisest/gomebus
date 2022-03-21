package gomebus

import (
	"sync"
)

var em_instance *EventMap
var em_once sync.Once

type event_key string
type channel_key string

type EventMap struct {
	event_place map[event_key]map[channel_key]*Channel

	sync.RWMutex
}

func GetEventMap() *EventMap {

	em_once.Do(func() {
		em_instance = &EventMap{}
		em_instance.event_place = make(map[event_key]map[channel_key]*Channel)
	})

	return em_instance
}

func (em *EventMap) Add(event string, channel string, value *Channel) error {

	em.Lock()
	chan_map := em.event_place[event_key(event)]
	if chan_map == nil {
		chan_map = make(map[channel_key]*Channel)
	}
	chan_map[channel_key(channel)] = value
	em.event_place[event_key(event)] = chan_map
	em.Unlock()

	return nil
}

func (em *EventMap) Send(event_place string, src string, payload []byte) error {

	em.RLock()
	subscribers := em.Get(event_place)
	em.RUnlock()

	message := Message{Message_size: uint32(len(payload)) + uint32(HEADER_SIZE), Version: 0, Message_type: EVENT, Address: src, Payload: payload}

	for _, v := range subscribers {

		s := &Slide{Msg: &message, SrcAddress: event_place, DstAddress: v.address}
		v.receiver <- s
	}

	return nil
}

func (em *EventMap) Get(event string) []*Channel {

	em.RLock()
	defer em.RUnlock()

	chan_list := em.event_place[event_key(event)]
	size := len(chan_list)
	if size <= 0 {
		return nil
	}

	chan_vec := make([]*Channel, size)

	i := 0
	for _, v := range chan_list {
		chan_vec[i] = v
		i++
	}

	return chan_vec
}

func (em *EventMap) Count() int {

	em.RLock()
	defer em.RUnlock()

	return len(em.event_place)
}
