package engine

import (
	"sync"
)

type HashTable struct {
	mx   sync.RWMutex
	data map[string]string
}

func NewHashTable() *HashTable {

	return &HashTable{
		data: make(map[string]string),
	}
}

func (c *HashTable) Set(key string, value string) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.data[key] = value
}

func (c *HashTable) Get(key string) (string, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	value, ok := c.data[key]

	return value, ok
}

func (c *HashTable) Del(key string) {
	c.mx.Lock()
	defer c.mx.Unlock()
	delete(c.data, key)

}
