package database

import (
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
	Set(key string, value string) error
	Get(key string) (string, error)
	Del(key string) error
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

func (d *DB) HandleQuery(queryStr string) string {
	d.logger.Debug().Str("query", queryStr).Msg("handling query")
	query, err := d.computeLayer.Parse(queryStr)
	if err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	switch query.CommandID() {
	case compute.SetCommandID:
		return d.handleSetQuery(query)
	case compute.GetCommandID:
		return d.handleGetQuery(query)
	case compute.DelCommandID:
		return d.handleDelQuery(query)
	default:
		d.logger.Error().Int("command_id", query.CommandID()).Msg("compute layer is incorrect")
	}

	return "[error] internal error"
}

func (d *DB) handleSetQuery(query compute.Query) string {
	arguments := query.Arguments()
	if err := d.storageLayer.Set(arguments[0], arguments[1]); err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return "[SET OK]"
}

func (d *DB) handleGetQuery(query compute.Query) string {
	arguments := query.Arguments()
	value, err := d.storageLayer.Get(arguments[0])
	if errors.Is(err, storage.ErrNotFound) {
		return "[not found]"
	} else if err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return fmt.Sprintf("[GET OK] %s", value)
}

func (d *DB) handleDelQuery(query compute.Query) string {
	arguments := query.Arguments()
	if err := d.storageLayer.Del(arguments[0]); err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return "[DEL OK]"
}
