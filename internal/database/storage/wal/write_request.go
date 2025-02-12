package wal

import "kvdatabase/internal/concurrency"

type WriteRequest struct {
	transactionLog TransactionLog
	promise        concurrency.Promise
}

func NewWriteRequest(lsn int64, commandID int, args []string) WriteRequest {
	return WriteRequest{
		transactionLog: TransactionLog{
			LSN:       lsn,
			CommandID: commandID,
			Args:      args,
		},
		promise: concurrency.NewPromise(),
	}
}

func (l *WriteRequest) TransactionLog() TransactionLog {
	return l.transactionLog
}

func (l *WriteRequest) SetResponse(err error) {
	l.promise.Set(err)
}

func (l *WriteRequest) FutureResponse() concurrency.Future {
	return l.promise.GetFuture()
}
