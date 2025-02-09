package wal

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type TransactionLogReader struct {
	Dir string
}

func NewTransactionLogReader(dir string) (*TransactionLogReader, error) {
	return &TransactionLogReader{
		Dir: dir,
	}, nil
}

func (r *TransactionLogReader) Read() ([]TransactionLog, error) {
	var transactionLogs []TransactionLog

	fileNames, err := r.ListWALFiles()
	if err != nil {
		return nil, err
	}

	for _, filename := range fileNames {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Bytes()
			buffer := bytes.NewBuffer(line)

			var log TransactionLog
			if err := log.Decode(buffer); err != nil {
				return nil, fmt.Errorf("failed to parse logs data: %w", err)
			}
			transactionLogs = append(transactionLogs, log)
		}
		file.Close()
	}

	return transactionLogs, nil
}

func (r *TransactionLogReader) ListWALFiles() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(r.Dir, "wal_*.log"))
	if err != nil {
		fmt.Println("Ошибка при поиске файлов журнала:", err)
		return nil, err
	}

	// Сортировка файлов по имени (по порядку LSN)
	sort.Strings(files)

	return files, nil
}
