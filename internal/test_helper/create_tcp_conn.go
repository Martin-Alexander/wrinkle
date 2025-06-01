package test_helper

import (
	"net"
	"testing"
	"time"
)

func CreateTcpConnection(t *testing.T) (net.Conn, net.Conn, func()) {
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

	clientConn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatalf("Failed to create TCP connection: %v", err)
	}

	select {
	case serverConn := <-done:
		if err := clientConn.SetDeadline(time.Now().Add(time.Second)); err != nil {
			t.Fatalf("Failed to set deadline on client connection: %v", err)
		}
		if err := serverConn.SetDeadline(time.Now().Add(time.Second)); err != nil {
			t.Fatalf("Failed to set deadline on server connection: %v", err)
		}

		close := func() {
			clientConn.Close()
			serverConn.Close()
			listener.Close()
		}

		return clientConn, serverConn, close
	case <-time.After(time.Millisecond):
		t.Fatalf("Timeout waiting for server connection")
		return nil, nil, nil
	}
}
