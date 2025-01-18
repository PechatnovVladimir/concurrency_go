package parser

import (
	"github.com/rs/zerolog/log"
	"strings"
)

type Engineer interface {
	Set(key string, value string)
	Get(key string) (string, bool)
	Del(key string)
}

type CommandParser struct {
	engine Engineer
}

func NewCommandParser(engine Engineer) *CommandParser {
	return &CommandParser{
		engine: engine,
	}
}

func (cp *CommandParser) Execute(command string) (string, bool) {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		log.Error().Msg("Empty command")
		return "ERROR: Empty command", false
	}

	switch strings.ToUpper(fields[0]) {
	case "SET":
		if len(fields) != 3 {
			log.Error().Msg("SET command requires 2 arguments")
			return "ERROR: SET command requires 2 arguments", false
		}
		key, value := fields[1], fields[2]
		cp.engine.Set(key, value)
		log.Info().Msg("Success " + command)
		return "OK", true

	case "GET":
		if len(fields) != 2 {
			log.Error().Msg("GET command requires 1 argument")
			return "ERROR: GET command requires 1 argument", false
		}
		key := fields[1]
		if value, exists := cp.engine.Get(key); exists {
			log.Info().Msg("Success " + command)
			return value, true
		}
		log.Info().Msg("Not found " + command)
		return "(not found)", false

	case "DEL":
		if len(fields) != 2 {
			return "ERROR: DEL command requires 1 argument", false
		}
		key := fields[1]
		cp.engine.Del(key)
		log.Info().Msg("Success " + command)
		return "(deleted)", true

	default:
		log.Error().Msg("Unknown command:" + command)
		return "ERROR: Unknown command", false
	}
}
