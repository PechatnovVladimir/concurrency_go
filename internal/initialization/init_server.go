package initialization

import (
	"errors"
	"github.com/rs/zerolog"
	"kvdatabase/internal/config"
	"kvdatabase/internal/network"
	"time"
)

func CreateServer(cfg *config.NetworkConfig, logger *zerolog.Logger) (*network.TCPServer, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	address := ":3333"
	maxConnection := 100
	maxMessageSize := 4 << 10
	idleTimeout := time.Duration(5 * time.Minute)

	if cfg != nil {
		if cfg.Address != "" {
			address = cfg.Address
		}

		if cfg.MaxConnections != 0 {
			maxConnection = cfg.MaxConnections
		}

		if cfg.MaxMessageSize != 0 {

			maxMessageSize = cfg.MaxMessageSize
		}

		if cfg.IdleTimeout != 0 {
			idleTimeout = cfg.IdleTimeout
		}
	}

	return network.NewTCPServer(address, idleTimeout, maxMessageSize, maxConnection, logger)
}
