package initialization

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"kvdatabase/internal/config"
	"kvdatabase/internal/database"
	compute2 "kvdatabase/internal/database/compute"
	storage2 "kvdatabase/internal/database/storage"
	engine2 "kvdatabase/internal/database/storage/engine"
	wal2 "kvdatabase/internal/database/storage/wal"
	"kvdatabase/internal/network"
)

type Init struct {
	logger *zerolog.Logger
	engine *engine2.Engine
	wal    *wal2.WAL
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

	wal, err := CreateWAL(cfg.WAL, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize wal: %w", err)
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
		wal:    wal,
		server: server,
	}, nil
}

func (i *Init) StartApp(ctx context.Context) error {
	compute, err := compute2.NewCompute(i.logger)
	if err != nil {
		return err
	}

	storage, err := storage2.NewStorage(i.engine, i.wal, i.logger)
	if err != nil {
		return err
	}

	db, err := database.NewDB(compute, storage, i.logger)
	if err != nil {
		return err
	}

	group, groupCtx := errgroup.WithContext(ctx)
	if i.wal != nil {
		group.Go(func() error {
			i.wal.Start(groupCtx)
			return nil
		})
	}

	group.Go(func() error {

		i.server.HandleQueries(ctx, func(ctx context.Context, data []byte) []byte {
			res := db.HandleQuery(ctx, string(data))
			return []byte(res)
		})
		return nil
	})

	err = group.Wait()

	return err
}
