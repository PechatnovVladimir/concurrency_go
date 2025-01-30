package initialization

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"kvdatabase/internal/config"
	"kvdatabase/internal/database"
	compute2 "kvdatabase/internal/database/compute"
	storage2 "kvdatabase/internal/database/storage"
	engine2 "kvdatabase/internal/database/storage/engine"
	"kvdatabase/internal/network"
)

type Init struct {
	logger *zerolog.Logger
	engine *engine2.Engine
	server *network.TCPServer
}

func NewInit(cfg *config.Config) (*Init, error) {
	if cfg == nil {
		return nil, errors.New("filed to initialization: config is invalid")
	}

	logger, err := CreateLogger(cfg.Logging)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	engine, err := CreateEngine(cfg.Engine, logger)
	if err != nil {

		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	server, err := CreateServer(cfg.Network, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize network: %w", err)
	}

	return &Init{
		logger: logger,
		engine: engine,
		server: server,
	}, nil
}

func (i *Init) StartApp() error {
	compute, err := compute2.NewCompute(i.logger)
	if err != nil {
		return err
	}

	storage, err := storage2.NewStorage(i.engine, i.logger)
	if err != nil {
		return err
	}

	db, err := database.NewDB(compute, storage, i.logger)
	if err != nil {
		return err
	}

	i.server.Start(func(data []byte) []byte {
		res := db.HandleQuery(string(data))
		return []byte(res)
	})

	return nil
}
