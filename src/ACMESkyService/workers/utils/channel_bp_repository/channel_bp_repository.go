package chanBPRepo

import (
	"sync"
)

var channelStore sync.Map

func SetContext(id string) {
	var key string = id
	channelStore.Store(key, make(chan map[string]interface{}))
}

func GetContext(id string) chan map[string]interface{} {

	var key string = id

	c, _ := channelStore.Load(key)
	if c != nil {
		return c.(chan map[string]interface{})
	}
	return nil
}

func UnsetContext(id string) {
	var key string = id
	c, hasValue := channelStore.LoadAndDelete(key)
	if hasValue {
		close(c.(chan map[string]interface{}))
	}
}
