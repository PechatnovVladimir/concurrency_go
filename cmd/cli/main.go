package main

import (
	"bufio"
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

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	address := flag.String("address", "localhost:3223", "Address of the spider")
	idleTimeout := flag.Duration("idle_timeout", 5*time.Minute, "Idle timeout for connection")
	maxMessageSize := flag.Int("max_message_size", 4096, "Max message size for connection")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	client, err := network.NewTCPClient(*address, *idleTimeout, *maxMessageSize)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect with server")
	}

	for {
		fmt.Print("[kv-db] > ")
		request, err := reader.ReadString('\n')
		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal().Err(err).Msg("connection was closed")
		} else if err != nil {
			logger.Fatal().Err(err).Msg("failed to read query")
		}

		response, err := client.Send([]byte(request))
		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal().Err(err).Msg("connection was closed")
		} else if err != nil {
			logger.Error().Err(err).Msg("failed to send query")
		}

		fmt.Println(string(response))
	}
}
