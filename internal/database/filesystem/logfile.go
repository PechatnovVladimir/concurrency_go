package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	LogFileName = "wal"
	LogFileExt  = "log"
)

type LogFile struct {
	file     *os.File
	dir      string
	fileSize int64
	maxSize  int64
	index    int
}

func findLastIndex(dir string) (int, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	var indices []int

	prefixLogFile := LogFileName + "_"
	suffixLogFile := "." + LogFileExt

	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), prefixLogFile) && strings.HasSuffix(file.Name(), suffixLogFile) {
			indexStr := strings.TrimPrefix(file.Name(), prefixLogFile)
			indexStr = strings.TrimSuffix(indexStr, suffixLogFile)
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				continue
			}
			indices = append(indices, index)
		}
	}

	if len(indices) == 0 {
		return 0, nil
	}

	sort.Ints(indices)
	return indices[len(indices)-1], nil
}

func NewLogFile(dir string, maxSize int64) (*LogFile, error) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	lf := &LogFile{
		dir:     dir,
		maxSize: maxSize,
	}

	lastIndex, err := findLastIndex(lf.dir)
	if err != nil {
		return nil, err
	}
	lf.index = lastIndex

	if err := lf.openLastFile(); err != nil {
		return nil, err
	}

	return lf, nil
}

func (lf *LogFile) openLastFile() error {
	if lf.index == 0 {
		return lf.openNextFile()
	}

	fileName := filepath.Join(lf.dir, fmt.Sprintf("%s_%05d.%s", LogFileName, lf.index, LogFileExt))
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	lf.file = file
	lf.fileSize = fileInfo.Size()

	return nil
}

func (lf *LogFile) openNextFile() error {
	if lf.file != nil {
		lf.file.Close()
	}

	lf.index++
	fileName := filepath.Join(lf.dir, fmt.Sprintf("%s_%05d.%s", LogFileName, lf.index, LogFileExt))
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	lf.file = file
	lf.fileSize = 0

	return nil
}

func (lf *LogFile) Close() error {

	if lf.file != nil {
		return lf.file.Close()
	}
	return nil
}

func (lf *LogFile) Write(data []byte) error {

	if lf.fileSize+int64(len(data)) > lf.maxSize {
		if err := lf.openNextFile(); err != nil {
			return err
		}
	}

	n, err := lf.file.Write(data)
	if err != nil {
		return err
	}

	if err = lf.file.Sync(); err != nil {
		return err
	}

	lf.fileSize += int64(n)
	return nil
}
