package wal

import (
	"context"
	"errors"
	"kvdatabase/internal/common"
	"kvdatabase/internal/concurrency"
	"kvdatabase/internal/database/compute"
	"sync"
	"time"
)

type transactionLogWriter interface {
	Write([]WriteRequest)
}

type transactionLogReader interface {
	ReadWALFiles() ([]TransactionLog, error)
}

type WAL struct {
	transactionLogWriter transactionLogWriter
	transactionLogReader transactionLogReader

	flushTimeout time.Duration
	maxBatchSize int

	batches chan []WriteRequest
	mx      sync.Mutex
	batch   []WriteRequest
}

func NewWAL(writer transactionLogWriter, reader transactionLogReader, flushTimeout time.Duration, maxBatchSize int) (*WAL, error) {
	if writer == nil {
		return nil, errors.New("writer is invalid")
	}
	if reader == nil {
		return nil, errors.New("reader is invalid")
	}
	return &WAL{
		transactionLogWriter: writer,
		transactionLogReader: reader,
		flushTimeout:         flushTimeout,
		maxBatchSize:         maxBatchSize,
		batches:              make(chan []WriteRequest, 1),
	}, nil
}

func (w *WAL) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(w.flushTimeout)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				w.flushBatch()
				return
			default:
			}

			select {
			case <-ctx.Done():
				w.flushBatch()
				return
			case batch := <-w.batches:
				w.transactionLogWriter.Write(batch)
				ticker.Reset(w.flushTimeout)
			case <-ticker.C:
				w.flushBatch()
			}
		}
	}()
}

func (w *WAL) Recover() ([]TransactionLog, error) {
	return w.transactionLogReader.ReadWALFiles()
}

func (w *WAL) Set(ctx context.Context, key string, value string) concurrency.Future {
	return w.push(ctx, compute.SetCommandID, []string{key, value})
}

func (w *WAL) Del(ctx context.Context, key string) concurrency.Future {
	return w.push(ctx, compute.DelCommandID, []string{key})
}

func (w *WAL) push(ctx context.Context, commandID int, args []string) concurrency.Future {
	id := common.GetIDFromContext(ctx)
	op := NewWriteRequest(id, commandID, args)

	concurrency.WithLock(&w.mx, func() {
		w.batch = append(w.batch, op)
		if len(w.batch) == w.maxBatchSize {
			w.batches <- w.batch
			w.batch = nil
		}
	})

	return op.FutureResponse()
}

func (w *WAL) flushBatch() {
	var batch []WriteRequest
	concurrency.WithLock(&w.mx, func() {
		batch = w.batch
		w.batch = nil
	})

	if len(batch) != 0 {
		w.transactionLogWriter.Write(batch)
	}
}
