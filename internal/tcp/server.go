package tcp

import (
	"net"

	"github.com/pkg/errors"
)

type ConnectionEvent struct {
	Conn net.Conn
	Err  error
}

type Server struct {
	network     string
	port        string
	listener    net.Listener
	readyCh     chan net.Addr
	connEventCh chan ConnectionEvent
	done        chan struct{}
}

func NewServer(network string, port string) (*Server, <-chan net.Addr) {
	readyCh := make(chan net.Addr, 1)
	connEventCh := make(chan ConnectionEvent, 16)
	done := make(chan struct{})

	server := &Server{
		network:     network,
		port:        port,
		readyCh:     readyCh,
		connEventCh: connEventCh,
		done:        done,
	}

	return server, readyCh
}

func (s *Server) Listen() error {
	address := net.JoinHostPort("", s.port)
	listener, err := net.Listen(s.network, address)
	if err != nil {
		return errors.WithStack(err)
	}
	defer listener.Close()

	s.listener = listener
	s.readyCh <- listener.Addr()

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			select {
			case <-s.done:
				return nil
			default:
				s.connEventCh <- ConnectionEvent{
					Conn: conn,
					Err:  errors.WithStack(err),
				}

				continue
			}
		}

		s.connEventCh <- ConnectionEvent{
			Conn: conn,
			Err:  err,
		}
	}
}

func (s *Server) Accept() <-chan ConnectionEvent {
	return s.connEventCh
}

func (s *Server) Close() {
	close(s.done)
	s.listener.Close()
}
