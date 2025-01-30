package initialization

import (
	"errors"
	"github.com/rs/zerolog"
	"kvdatabase/internal/config"
	"kvdatabase/internal/database/storage/engine"
)

func CreateEngine(cfg *config.EngineConfig, logger *zerolog.Logger) (*engine.Engine, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	if cfg == nil {
		return engine.NewEngine(logger)
	}

	logger.Debug().Str("type", cfg.Type).Msg("create engine")

	return engine.NewEngine(logger)

}
