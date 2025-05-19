package tcp_server

import (
	"log/slog"
	"net"
)

type ConnectionHandler interface {
	HandleConnection(conn net.Conn) error
}

type Server struct {
	network     string
	port        string
	connHandler ConnectionHandler
}

func NewServer(network string, port string, connectionHandler ConnectionHandler) *Server {
	return &Server{
		network:     network,
		port:        port,
		connHandler: connectionHandler,
	}
}

func (s *Server) Listen(onReady func()) error {
	address := net.JoinHostPort("", s.port)
	listener, err := net.Listen(s.network, address)
	if err != nil {
		return err
	}
	onReady()
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("Error accepting connection", "error", err)
			continue
		}

		slog.Info("Accepted connection", "remoteAddr", conn.RemoteAddr().String())

		go func() {
			if err := s.connHandler.HandleConnection(conn); err != nil {
				slog.Error("Error handling connection", "error", err)
			}
		}()
	}
}
