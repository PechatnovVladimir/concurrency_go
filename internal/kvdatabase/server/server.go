package server

import (
	"bufio"
	"fmt"
	"net"
)

type Parser interface {
	Execute(command string) (string, bool)
}

type Server struct {
	parser Parser
}

func NewServer(p Parser) *Server {
	return &Server{
		parser: p,
	}
}

func (s *Server) Start(address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	defer func(ln net.Listener) {
		err := ln.Close()
		if err != nil {
			fmt.Printf("Failed to stop server: %v\n", err)
		}
	}(ln)

	fmt.Printf("Server is listening on %s\n", address)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("Failed to close connection: %v\n", err)
		}
	}(conn)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		response, _ := s.parser.Execute(line)
		_, err := conn.Write([]byte(response + "\n"))
		if err != nil {
			return
		}
	}
}
