package engine

import "sync"

type Engine struct {
	mx   sync.RWMutex
	data map[string]string
}

func NewEngine() *Engine {
	return &Engine{
		data: make(map[string]string),
	}
}

func (e *Engine) Set(key string, value string) {
	e.mx.Lock()
	defer e.mx.Unlock()

	e.data[key] = value
}

func (e *Engine) Get(key string) (string, bool) {
	e.mx.RLock()
	defer e.mx.RUnlock()

	value, ok := e.data[key]

	return value, ok
}

func (e *Engine) Del(key string) {
	e.mx.Lock()
	defer e.mx.Unlock()

	delete(e.data, key)
}
