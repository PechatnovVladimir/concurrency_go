package engine

import (
	"errors"
	"github.com/rs/zerolog"
	"sync"
)

type Engine struct {
	mx     sync.Mutex
	data   map[string]string
	logger *zerolog.Logger
}

func NewEngine(logger *zerolog.Logger) (*Engine, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &Engine{
		data:   make(map[string]string),
		logger: logger,
	}, nil
}

func (e *Engine) Set(key string, value string) {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.data[key] = value

	e.logger.Debug().Str("key", key).Str("value", value).Msg("engine set query")
}

func (e *Engine) Get(key string) (string, bool) {
	e.mx.Lock()
	defer e.mx.Unlock()
	value, ok := e.data[key]

	e.logger.Debug().Str("key", key).Str("value", value).Msg("engine get query")
	return value, ok
}

func (e *Engine) Del(key string) {
	e.mx.Lock()
	defer e.mx.Unlock()
	delete(e.data, key)

	e.logger.Debug().Str("key", key).Msg("engine del query")
}
