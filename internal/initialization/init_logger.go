package initialization

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"kvdatabase/internal/config"
	"os"
	"time"
)

func CreateLogger(cfg *config.LoggingConfig) (*zerolog.Logger, error) {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	file, err := os.OpenFile(cfg.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	//defer file.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("Не удалось открыть файл для записи логов")
		return nil, err
	}

	multiLogger := zerolog.MultiLevelWriter(consoleWriter, file)

	logger := zerolog.New(multiLogger).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	return &logger, nil
}
