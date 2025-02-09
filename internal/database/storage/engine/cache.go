package engine

import (
	"sync"
)

type Cache struct {
	mx   sync.RWMutex
	data map[string]string
}

func NewCache() *Cache {

	return &Cache{
		data: make(map[string]string),
	}
}

func (c *Cache) Set(key string, value string) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.data[key] = value
}

func (c *Cache) Get(key string) (string, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	value, ok := c.data[key]

	return value, ok
}

func (c *Cache) Del(key string) {
	c.mx.Lock()
	defer c.mx.Unlock()
	delete(c.data, key)

}
