package chatserver

import (
	"fmt"
	"net"
	"os"
)

type Server struct {
	host   string
	port   uint16
	hub    *Hub
	logger *Logger
}

func (s *Server) listen() {
	go s.hub.start()
	go s.logger.start()
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't Listen on  %s:%d, error %s.", s.host, s.port, err)
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		fmt.Printf("Server started,accepting connection\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't accept connection, error %s.", err)
			os.Exit(1)
		}
		go s.handleRequest(conn)
	}
}

func (s *Server) handleRequest(conn net.Conn) {
	session := NewChatSession(conn, s.hub)
	session.Start()
}
func (s *Server) Start() {
	s.listen()
}

func NewChatServer(config *Config) *Server {
	logger, err := NewLogger(config.LogPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening the log file, %s.\n", err)
		os.Exit(1)
	}
	return &Server{
		host:   config.Host,
		port:   config.ChatPort,
		hub:    NewHub(logger),
		logger: logger,
	}
}
