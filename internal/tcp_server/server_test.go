package tcp_server

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ConnectionHandlerMock struct {
	mock.Mock
}

func (m *ConnectionHandlerMock) HandleConnection(conn net.Conn) error {
	args := m.Called(conn)
	return args.Error(0)
}

func TestNewServer(t *testing.T) {
	connectionHandlerMock := new(ConnectionHandlerMock)

	server := NewServer("tcp4", "54321", connectionHandlerMock)

	assert.NotNil(t, server)
}

func TestStart(t *testing.T) {
	connectionHandlerMock := new(ConnectionHandlerMock)

	ready := make(chan bool)
	done := make(chan bool)
	errCh := make(chan error)

	connectionHandlerMock.
		On("HandleConnection", mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			done <- true
		})

	server := NewServer("tcp4", "54321", connectionHandlerMock)

	go func() {
		onReady := func() {
			ready <- true
		}
		errCh <- server.Listen(onReady)
	}()

	select {
	case <-ready:
	case <-errCh:
		t.Fatalf("Server returned an error: %v", <-errCh)
	case <-time.After(time.Second):
		t.Fatal("Timed out waiting for server to be ready")
	}

	conn, err := net.Dial("tcp4", "127.0.0.1:54321")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Timed out waiting for handler to be called")
	}

	connectionHandlerMock.AssertExpectations(t)
}
