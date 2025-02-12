package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"kvdatabase/internal/database/compute"
	"kvdatabase/internal/database/storage"
)

var (
	ErrComputeInvalid = errors.New("compute is invalid")
	ErrStorageInvalid = errors.New("storage is invalid")
	ErrLoggerInvalid  = errors.New("logger is invalid")
)

type storageLayer interface {
	Set(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

type computeLayer interface {
	Parse(queryStr string) (compute.Query, error)
}

type DB struct {
	computeLayer computeLayer
	storageLayer storageLayer
	logger       zerolog.Logger
}

func NewDB(computeLayer computeLayer, storageLayer storageLayer, logger *zerolog.Logger) (*DB, error) {
	if computeLayer == nil {
		return nil, ErrComputeInvalid
	}

	if storageLayer == nil {
		return nil, ErrStorageInvalid
	}

	if logger == nil {
		return nil, ErrLoggerInvalid
	}

	return &DB{
		computeLayer: computeLayer,
		storageLayer: storageLayer,
		logger:       *logger,
	}, nil
}

func (d *DB) HandleQuery(ctx context.Context, queryStr string) string {
	d.logger.Debug().Str("query", queryStr).Msg("handling query")
	query, err := d.computeLayer.Parse(queryStr)
	if err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	switch query.CommandID() {
	case compute.SetCommandID:
		return d.handleSetQuery(ctx, query)
	case compute.GetCommandID:
		return d.handleGetQuery(ctx, query)
	case compute.DelCommandID:
		return d.handleDelQuery(ctx, query)
	default:
		d.logger.Error().Int("command_id", query.CommandID()).Msg("compute layer is incorrect")
	}

	return "[error] internal error"
}

func (d *DB) handleSetQuery(ctx context.Context, query compute.Query) string {
	arguments := query.Arguments()
	if err := d.storageLayer.Set(ctx, arguments[0], arguments[1]); err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return "[SET OK]"
}

func (d *DB) handleGetQuery(ctx context.Context, query compute.Query) string {
	arguments := query.Arguments()
	value, err := d.storageLayer.Get(ctx, arguments[0])
	if errors.Is(err, storage.ErrNotFound) {
		return "[not found]"
	} else if err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return fmt.Sprintf("[GET OK] %s", value)
}

func (d *DB) handleDelQuery(ctx context.Context, query compute.Query) string {
	arguments := query.Arguments()
	if err := d.storageLayer.Del(ctx, arguments[0]); err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return "[DEL OK]"
}
