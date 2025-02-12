package wal

import (
	"bytes"
	"encoding/json"
)

type TransactionLog struct {
	LSN       int64
	CommandID int
	Args      []string
}

func (l *TransactionLog) Encode(buffer *bytes.Buffer) error {
	encoder := json.NewEncoder(buffer)
	return encoder.Encode(*l)
}

func (l *TransactionLog) Decode(buffer *bytes.Buffer) error {
	decoder := json.NewDecoder(buffer)
	return decoder.Decode(l)
}
