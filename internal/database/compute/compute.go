package compute

import (
	"errors"
	"github.com/rs/zerolog"
	"strings"
)

var (
	errInvalidQuery     = errors.New("empty query")
	errInvalidCommand   = errors.New("invalid command")
	errInvalidArguments = errors.New("invalid arguments")
	errLoggerIsInvalid  = errors.New("logger is invalid")
)

type Compute struct {
	logger *zerolog.Logger
}

func NewCompute(logger *zerolog.Logger) (*Compute, error) {
	if logger == nil {
		return nil, errLoggerIsInvalid
	}
	return &Compute{
		logger: logger,
	}, nil
}

func (d *Compute) Parse(queryStr string) (Query, error) {
	tokens := strings.Fields(queryStr)
	if len(tokens) == 0 {
		d.logger.Debug().Str("query", queryStr).Msg("empty tokens")
		return Query{}, errInvalidQuery
	}

	command := tokens[0]
	commandID := commandNameToCommandID(command)
	if commandID == UnknownCommandID {
		d.logger.Debug().Str("query", queryStr).Msg("invalid command")
		return Query{}, errInvalidCommand
	}

	query := NewQuery(commandID, tokens[1:])
	argumentsNumber := commandArgumentsNumber(commandID)
	if len(query.Arguments()) != argumentsNumber {
		d.logger.Debug().Str("query", queryStr).Msg("invalid arguments for query")
		return Query{}, errInvalidArguments
	}

	return query, nil

}
