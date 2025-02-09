package storage

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"kvdatabase/internal/common"
	"kvdatabase/internal/concurrency"
	"kvdatabase/internal/database/compute"
	"kvdatabase/internal/database/storage/wal"
)

var (
	ErrEngineIsInvalid = errors.New("engine is invalid")
	ErrWALIsInvalid    = errors.New("wal is invalid")
	ErrLoggerIsInvalid = errors.New("logger is invalid")
	ErrNotFound        = errors.New("not found")
)

type Engine interface {
	Set(context.Context, string, string)
	Get(context.Context, string) (string, bool)
	Del(context.Context, string)
}

type WAL interface {
	Recover() ([]wal.TransactionLog, error)
	Set(context.Context, string, string) concurrency.Future
	Del(context.Context, string) concurrency.Future
}

type Storage struct {
	engine Engine
	wal    WAL
	id     *ID
	logger *zerolog.Logger
}

func NewStorage(engine Engine, wal WAL, logger *zerolog.Logger) (*Storage, error) {
	if engine == nil {
		return nil, ErrEngineIsInvalid
	}

	if wal == nil {
		return nil, ErrWALIsInvalid
	}

	if logger == nil {
		return nil, ErrLoggerIsInvalid
	}

	storage := &Storage{
		engine: engine,
		wal:    wal,
		logger: logger,
	}

	var lastLSN int64

	if storage.wal != nil {
		logs, err := storage.wal.Recover()
		if err != nil {
			logger.Error().Err(err).Msg("failed to recover data from WAL")
		} else {
			lastLSN = storage.recoveryDataFromLog(logs)
		}
	}

	storage.id = NewID(lastLSN)
	return storage, nil
}

func (s *Storage) Set(ctx context.Context, key, value string) error {

	id := s.id.Generate()
	ctx = common.ContextWithID(ctx, id)

	if s.wal != nil {
		futureResponse := s.wal.Set(ctx, key, value)
		if err := futureResponse.Get(); err != nil {
			return err
		}
	}

	s.engine.Set(ctx, key, value)
	return nil
}

func (s *Storage) Get(ctx context.Context, key string) (string, error) {

	id := s.id.Generate()
	ctx = common.ContextWithID(ctx, id)

	value, found := s.engine.Get(ctx, key)
	if !found {
		return "", ErrNotFound
	}
	return value, nil
}

func (s *Storage) Del(ctx context.Context, key string) error {

	id := s.id.Generate()
	ctx = common.ContextWithID(ctx, id)

	if s.wal != nil {
		futureResponse := s.wal.Del(ctx, key)
		if err := futureResponse.Get(); err != nil {
			return err
		}
	}

	s.engine.Del(ctx, key)
	return nil
}

func (s *Storage) recoveryDataFromLog(logs []wal.TransactionLog) int64 {
	var lastLSN int64
	for _, log := range logs {
		lastLSN = max(lastLSN, log.LSN)
		ctx := common.ContextWithID(context.Background(), log.LSN)
		switch log.CommandID {
		case compute.SetCommandID:
			s.engine.Set(ctx, log.Args[0], log.Args[1])
		case compute.DelCommandID:
			s.engine.Del(ctx, log.Args[0])
		}
	}

	return lastLSN
}
