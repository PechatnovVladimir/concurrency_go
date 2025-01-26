package main

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog"
	"kvdatabase/internal/database"
	"kvdatabase/internal/database/compute"
	"kvdatabase/internal/database/storage"
	"kvdatabase/internal/database/storage/engine"
	"os"
	"time"
)

func CreateLogger() *zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	return &logger
}

func main() {

	logger := CreateLogger()

	e := engine.NewEngine()
	c, err := compute.NewCompute(logger)
	if err != nil {
		fmt.Println(err)
		return
	}
	s, err := storage.NewStorage(e, logger)
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := database.NewDB(c, s, logger)
	if err != nil {
		fmt.Println(err)
		return
	}

	var command string
	for {
		fmt.Print(">")
		command, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		res := db.HandleQuery(command)
		fmt.Println(res)
	}
}
