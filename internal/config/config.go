package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"time"
)

type Config struct {
	Engine  *EngineConfig  `yaml:"engine"`
	WAL     *WALConfig     `yaml:"wal"`
	Network *NetworkConfig `yaml:"network"`
	Logging *LoggingConfig `yaml:"logging"`
}

type EngineConfig struct {
	Type string `yaml:"type"`
}

type NetworkConfig struct {
	Address        string        `yaml:"address"`
	MaxConnections int           `yaml:"max_connections"`
	MaxMessageSize int           `yaml:"max_message_size"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

type WALConfig struct {
	FlushingBatchSize    int           `yaml:"flushing_batch_size"`
	FlushingBatchTimeout time.Duration `yaml:"flushing_batch_timeout"`
	MaxSegmentSize       int64         `yaml:"max_segment_size"`
	DataDirectory        string        `yaml:"data_directory"`
}

func Load(reader io.Reader) (*Config, error) {
	if reader == nil {
		return nil, errors.New("incorrect reader")
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.New("falied to read buffer")
	}

	var config Config

	err = yaml.Unmarshal(data, &config)

	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}
