package pg_proxy

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRelay(t *testing.T) {
	_, frontendServer, backendClient, _, closeConnections := createConnections(t)
	defer closeConnections()

	relay := NewRelay(frontendServer, backendClient)

	assert.NotNil(t, relay)
	assert.Equal(t, frontendServer, relay.frontendConn)
	assert.Equal(t, backendClient, relay.backendConn)
}

func TestStart(t *testing.T) {
	frontendClient, frontendServer := createTcpConnection(t)
	backendClient, backendServer := createTcpConnection(t)
	defer frontendClient.Close()
	defer frontendServer.Close()
	defer backendClient.Close()
	defer backendServer.Close()

	relay := NewRelay(frontendServer, backendClient)

	go relay.Start()

	frontendDataSent := []byte("frontend test data")
	frontendClient.Write(frontendDataSent)
	buffer := make([]byte, len(frontendDataSent))
	backendServer.Read(buffer)

	assert.Equal(t, string(frontendDataSent), string(buffer))

	backendDataSent := []byte("backend test data")
	backendServer.Write(backendDataSent)
	backendBuffer := make([]byte, len(backendDataSent))
	n, _ := frontendClient.Read(backendBuffer)

	assert.Equal(t, string(backendDataSent), string(backendBuffer[:n]))
}

func createTcpConnection(t *testing.T) (net.Conn, net.Conn) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create TCP listener: %v", err)
	}
	defer listener.Close()

	done := make(chan net.Conn, 1)

	go func() {
		serverConn, _ := listener.Accept()

		done <- serverConn
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatalf("Failed to create TCP connection: %v", err)
	}

	select {
	case serverConn := <-done:
		return conn, serverConn
	case <-time.After(time.Millisecond):
		t.Fatalf("Timeout waiting for server connection")
		return nil, nil
	}
}
