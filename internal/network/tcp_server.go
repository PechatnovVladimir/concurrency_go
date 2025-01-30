package network

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"kvdatabase/internal/concurrency"
	"net"
	"sync"
	"time"
)

type TCPHandler = func([]byte) []byte

type TCPServer struct {
	listener  net.Listener
	semaphore *concurrency.Semaphore

	idleTimeout    time.Duration
	bufferSize     int
	maxConnections int

	logger *zerolog.Logger
}

// TODO: переделать параметры сервера через опции
func NewTCPServer(address string, idleTimeout time.Duration, bufferSize int, maxConnections int, logger *zerolog.Logger) (*TCPServer, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	server := &TCPServer{
		listener:       listener,
		logger:         logger,
		maxConnections: maxConnections,
		idleTimeout:    idleTimeout,
		bufferSize:     bufferSize,
	}

	if server.maxConnections != 0 {
		server.semaphore = concurrency.NewSemaphore(server.maxConnections)
	}
	if server.bufferSize == 0 {
		server.bufferSize = 4 << 10
	}

	return server, nil
}

func (s *TCPServer) handleConnection(connection net.Conn, handler TCPHandler) {
	defer func() {
		err := connection.Close()

		if err != nil {
			s.logger.Warn().Err(err).Msg("failed to close connection")
		}
	}()

	request := make([]byte, s.bufferSize)

	for {
		if s.idleTimeout != 0 {
			if err := connection.SetReadDeadline(time.Now().Add(s.idleTimeout)); err != nil {
				s.logger.Warn().Err(err).Msg("failed to set read deadline")
				break
			}
		}

		count, err := connection.Read(request)
		if err != nil && err != io.EOF {
			s.logger.Warn().Err(err).Str("address", connection.RemoteAddr().String()).Msg("failed to read data")
			break
		} else if count == s.bufferSize {
			s.logger.Warn().Err(err).Int("buffer_size", s.bufferSize).Msg("small buffer size")
			break
		}

		if s.idleTimeout != 0 {
			if err := connection.SetWriteDeadline(time.Now().Add(s.idleTimeout)); err != nil {
				s.logger.Warn().Err(err).Msg("failed to set read deadline")
				break
			}
		}

		response := handler(request[:count])
		if _, err := connection.Write(response); err != nil {
			s.logger.Warn().Err(err).Str("address", connection.RemoteAddr().String()).Msg("failed to write data")
			break
		}
	}
}

func (s *TCPServer) Start(handler TCPHandler) {
	var wg sync.WaitGroup

	s.logger.Info().Str("adrress", s.listener.Addr().String()).Msg("start server")

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}

			s.logger.Error().Err(err).Msg("failed to accept")
			continue
		}

		s.semaphore.Acquire()
		wg.Add(1)

		go func(conn net.Conn) {
			defer func() {
				conn.Close()
				s.semaphore.Release()
				wg.Done()
			}()
			s.handleConnection(conn, handler)
		}(conn)
	}
	wg.Wait()
}
