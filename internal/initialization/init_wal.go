package initialization

import (
	"errors"
	"github.com/rs/zerolog"
	"kvdatabase/internal/config"
	"kvdatabase/internal/database/filesystem"
	"kvdatabase/internal/database/storage/wal"
	"time"
)

const (
	defaultFlushingBatchSize    = 50
	defaultFlushingBatchTimeout = 10 * time.Millisecond
	defaultMaxSegmentSize       = 10 << 20
	defaultWALDataDirectory     = "./wal"
)

func CreateWAL(cfg *config.WALConfig, logger *zerolog.Logger) (*wal.WAL, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	flushingBatchSize := defaultFlushingBatchSize
	flushingBatchTimeout := defaultFlushingBatchTimeout
	maxSegmentSize := int64(defaultMaxSegmentSize)
	dataDirectory := defaultWALDataDirectory

	if cfg != nil {

		if cfg.FlushingBatchSize != 0 {
			flushingBatchSize = cfg.FlushingBatchSize
		}

		if cfg.FlushingBatchTimeout != 0 {
			flushingBatchTimeout = cfg.FlushingBatchTimeout
		}

		if cfg.MaxSegmentSize != 0 {
			maxSegmentSize = cfg.MaxSegmentSize
		}

		if cfg.DataDirectory != "" {
			dataDirectory = cfg.DataDirectory
		}

	}

	lf, err := filesystem.NewLogFile(dataDirectory, maxSegmentSize)
	if err != nil {
		return nil, err
	}

	w, err := wal.NewTransactionLogWriter(lf, logger)

	if err != nil {
		return nil, err
	}

	r, err := wal.NewTransactionLogReader(dataDirectory)
	if err != nil {
		return nil, err
	}

	return wal.NewWAL(w, r, flushingBatchTimeout, flushingBatchSize)
}
