package main

import (
	"crypto/tls"
	"encoding/hex"
	"log"
	"net"
	"wrinkle/internal/postgres_proxy"
)

type Server struct {
	tlsServerConfig *tls.Config
	tlsClientConfig *tls.Config
}

func (s *Server) Start() {
	cert, err := tls.LoadX509KeyPair("/app/.ssl/server.crt", "/app/.ssl/server.key")
	if err != nil {
		log.Fatalf("Failed to load cert: %v", err)
	}

	s.tlsServerConfig = &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	s.tlsClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	listener, err := net.Listen("tcp4", ":54321")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	for {
		frontendConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		backendConn, err := connectToBackend()
		if err != nil {
			log.Printf("Failed to connect to backend: %v", err)
			return
		}

		log.Printf("Accepted connection from %s", frontendConn.RemoteAddr())

		go s.handleConnection(frontendConn, backendConn)
	}
}

func (s *Server) handleConnection(frontendConn net.Conn, backendConn net.Conn) {
	success, err := postgres_proxy.HandleTlsNegotiation(frontendConn, backendConn)
	if err != nil {
		log.Printf("Failed to handle TLS negotiation: %v", err)
		return
	}
	if !success {
		log.Printf("TLS negotiation failed")
		return
	}

	frontendTlsConn, backendTlsConn, err := postgres_proxy.HandleDoubleTlsHandshake(
		frontendConn,
		backendConn,
		s.tlsServerConfig,
		s.tlsClientConfig,
	)
	if err != nil {
		log.Printf("Failed to upgrade connection to TLS: %v", err)
		return
	}

	go relay(frontendTlsConn, backendTlsConn, "PostgreSQL client")
	go relay(backendTlsConn, frontendTlsConn, "PostgreSQL server")
}

func main() {
	server := &Server{}
	server.Start()
}

func readUpTo(conn net.Conn, maxBytes int) ([]byte, error) {
	buf := make([]byte, maxBytes)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func connectToBackend() (net.Conn, error) {
	conn, err := net.Dial("tcp4", "postgres:5432")
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func relay(sourceConn net.Conn, destinationConn net.Conn, logLabel string) {
	for {
		data, err := readUpTo(sourceConn, 4096)
		if err != nil {
			log.Printf("Failed to read from connection: %v", err)
			return
		}

		if _, err := destinationConn.Write(data); err != nil {
			log.Printf("Error writing to destination: %v", err)
			return
		}

		logPacket(logLabel, data)
	}
}

func logPacket(label string, data []byte) {
	log.Printf(" -- [%s]:\n%s\n", label, hex.Dump(data))
}
