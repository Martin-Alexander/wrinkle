package tcp

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	server, _ := NewServer("tcp4", "")

	assert.NotNil(t, server)
}

func TestStart(t *testing.T) {
	server, readyCh := NewServer("tcp4", "")

	go server.Listen()
	defer server.Close()

	var addr net.Addr

	select {
	case addr = <-readyCh:
		assert.NotNil(t, addr)
	case <-time.After(time.Second):
		t.Fatal("Timed out waiting for server to be ready")
	}
}

func TestAccept(t *testing.T) {
	server, readyCh := NewServer("tcp4", "")

	go server.Listen()
	defer server.Close()

	var addr net.Addr

	select {
	case addr = <-readyCh:
	case <-time.After(time.Second):
		t.Fatal("Timed out waiting for server to be ready")
	}

	conn, err := net.Dial("tcp4", addr.String())
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	select {
	case connEvent := <-server.Accept():
		assert.NotNil(t, connEvent)
	case <-time.After(time.Second):
		t.Fatal("Timed out waiting for connection handler to be called")
	}
}
