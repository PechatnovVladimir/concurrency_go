package wal

import (
	"bytes"
	"errors"
	"github.com/rs/zerolog"
)

type logfile interface {
	Write([]byte) error
}

type TransactionLogWriter struct {
	logfile logfile
	logger  *zerolog.Logger
}

func NewTransactionLogWriter(logfile logfile, logger *zerolog.Logger) (*TransactionLogWriter, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}
	if logfile == nil {
		return nil, errors.New("logfile is invalid")
	}

	return &TransactionLogWriter{
		logfile: logfile,
		logger:  logger,
	}, nil

}

func (w *TransactionLogWriter) Write(requests []WriteRequest) {
	var buffer bytes.Buffer
	for idx := range requests {
		log := requests[idx].TransactionLog()
		if err := log.Encode(&buffer); err != nil {
			w.logger.Warn().Err(err).Msg("failed to encode logs data")
			w.acknowledgeWrite(requests, err)
			return
		}
	}

	err := w.logfile.Write(buffer.Bytes())
	if err != nil {
		w.logger.Warn().Err(err).Msg("failed to write logs data")
	}

	w.acknowledgeWrite(requests, err)
}

func (w *TransactionLogWriter) acknowledgeWrite(requests []WriteRequest, err error) {
	for idx := range requests {
		requests[idx].SetResponse(err)
	}
}
