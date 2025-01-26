package storage

import (
	"errors"
	"github.com/rs/zerolog"
)

var (
	ErrEngineIsInvalid = errors.New("engine is invalid")
	ErrLoggerIsInvalid = errors.New("logger is invalid")
	ErrNotFound        = errors.New("not found")
)

type Engine interface {
	Set(string, string)
	Get(string) (string, bool)
	Del(string)
}

type Storage struct {
	engine Engine
	logger *zerolog.Logger
}

func NewStorage(engine Engine, logger *zerolog.Logger) (*Storage, error) {
	if engine == nil {
		return nil, ErrEngineIsInvalid
	}

	if logger == nil {
		return nil, ErrLoggerIsInvalid
	}

	return &Storage{
		engine: engine,
		logger: logger,
	}, nil
}

func (s *Storage) Set(key, value string) error {
	s.engine.Set(key, value)
	return nil
}

func (s *Storage) Get(key string) (string, error) {
	value, found := s.engine.Get(key)
	if !found {
		return "", ErrNotFound
	}
	return value, nil
}

func (s *Storage) Del(key string) error {
	s.engine.Del(key)
	return nil
}
