package proxy

import (
	"net"
	"testing"
	"time"
	"wrinkle/internal/pg"

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

	frontendMessage := pg.Message{
		Type:   byte(pg.ClientSimpleQuery),
		Length: 9,
		Data:   []byte{0x04, 0xd2, 0x16, 0x2f},
	}
	frontendClient.Write(frontendMessage.Binary())
	buffer := make([]byte, len(frontendMessage.Binary()))
	backendServer.Read(buffer)

	assert.Equal(t, string(frontendMessage.Binary()), string(buffer))

	backendMessage := pg.Message{
		Type:   byte(pg.ServerCommandComplete),
		Length: 9,
		Data:   []byte{0x04, 0xd2, 0x16, 0x2f},
	}
	backendServer.Write(backendMessage.Binary())
	backendBuffer := make([]byte, len(backendMessage.Binary()))
	n, _ := frontendClient.Read(backendBuffer)

	assert.Equal(t, string(backendMessage.Binary()), string(backendBuffer[:n]))
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
