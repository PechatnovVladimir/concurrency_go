package main

import (
	"fmt"
	"github.com/PechatnovVladimir/concurrency_go/config"
	"github.com/PechatnovVladimir/concurrency_go/internal/kvdatabase/compute/parser"
	"github.com/PechatnovVladimir/concurrency_go/internal/kvdatabase/server"
	"github.com/PechatnovVladimir/concurrency_go/internal/kvdatabase/storage/engine"
	"github.com/PechatnovVladimir/concurrency_go/pkg/logger"
	"github.com/rs/zerolog/log"
)

func main() {

	c, err := config.New()
	if err != nil {
		c.App.Name = "my-app"
		c.App.Version = "v0.1.0"
		c.App.Host = "localhost:12345"
		c.Logger.AppVersion = c.App.Version
		c.Logger.AppName = c.App.Name
		c.Logger.Level = "info"
		c.Logger.PrettyConsole = true
	}

	logger.Init(c.Logger)

	e := engine.NewEngine()
	p := parser.NewCommandParser(e)
	s := server.NewServer(p)

	if err := s.Start(c.App.Host); err != nil {
		log.Error().Msg(err.Error())
		fmt.Printf("Error: %v\n", err)
	}
}
