package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"kvdatabase/internal/network"
	"os"
	"syscall"
	"time"
)

func main() {
	//это не в рамках ДЗ, для себя...
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	address := flag.String("address", "localhost:3323", "Address of the spider")
	idleTimeout := flag.Duration("idle_timeout", 5*time.Minute, "Idle timeout for connection")
	maxMessageSize := flag.Int("max_message_size", 4096, "Max message size for connection")
	flag.Parse()

	client, err := network.NewTCPClient(*address, *idleTimeout, *maxMessageSize)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect with server")
	}

	//заливаем тысячу ключей, чтобы руками не вбивать с клиента, посмотреть как создаются файлы wal
	//можно не обращать на это внимание, это для себя
	for i := 1; i < 10000; i++ {
		request := fmt.Sprintf("SET %d %d", i, i)

		_, err := client.Send([]byte(request))
		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal().Err(err).Msg("connection was closed")
		} else if err != nil {
			logger.Error().Err(err).Msg("failed to send query")
		}
	}

	//удаляем каждый 50-й ключ
	//можно не обращать на это внимание, это для себя
	for i := 50; i < 10000; i += 50 {
		request := fmt.Sprintf("DEL %d", i)

		_, err := client.Send([]byte(request))
		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal().Err(err).Msg("connection was closed")
		} else if err != nil {
			logger.Error().Err(err).Msg("failed to send query")
		}
	}
}
