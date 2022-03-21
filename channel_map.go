package gomebus

import (
	"sync"
)

var cm_instance *ChannelMap
var cm_once sync.Once

type ChannelMap struct {
	channels map[string]*Channel

	sync.RWMutex
}

func GetChannelMap() *ChannelMap {

	cm_once.Do(func() {
		cm_instance = &ChannelMap{}
		cm_instance.channels = make(map[string]*Channel)
	})

	return cm_instance
}

func (cm *ChannelMap) Add(key string, value *Channel) error {

	cm.Lock()
	cm.channels[key] = value
	cm.Unlock()

	return nil
}

func (cm *ChannelMap) Get(key string) *Channel {

	cm.RLock()
	defer cm.RUnlock()

	return cm.channels[key]
}

func (cm *ChannelMap) Count() int {

	cm.RLock()
	defer cm.RUnlock()

	return len(cm.channels)
}
