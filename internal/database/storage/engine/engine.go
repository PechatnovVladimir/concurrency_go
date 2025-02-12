package engine

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
)

type Engine struct {
	hashTable *HashTable
	logger    *zerolog.Logger
}

func NewEngine(logger *zerolog.Logger) (*Engine, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	c := NewHashTable()
	return &Engine{
		hashTable: c,
		logger:    logger,
	}, nil
}

func (e *Engine) Set(ctx context.Context, key string, value string) {
	e.hashTable.Set(key, value)

	e.logger.Debug().Str("key", key).Str("value", value).Msg("engine set query")
}

func (e *Engine) Get(ctx context.Context, key string) (string, bool) {
	value, ok := e.hashTable.Get(key)

	e.logger.Debug().Str("key", key).Str("value", value).Msg("engine get query")
	return value, ok
}

func (e *Engine) Del(ctx context.Context, key string) {
	e.hashTable.Del(key)

	e.logger.Debug().Str("key", key).Msg("engine del query")
}
