package main

import (
	"crypto/tls"
	"log"
	"log/slog"
	"wrinkle/internal/pg_proxy"
	"wrinkle/internal/tcp_server"
)

func main() {
	cert, err := tls.LoadX509KeyPair("/app/.ssl/server.crt", "/app/.ssl/server.key")
	if err != nil {
		log.Fatalf("Failed to load cert: %v", err)
	}

	tlsServerHandshaker := tcp_server.NewServerTlsHandshaker(&tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	})

	tlsClientHandshaker := tcp_server.NewClientTlsHandshaker(&tls.Config{
		InsecureSkipVerify: true,
	})

	bridgeBuilder := pg_proxy.NewBridgeBuilder(tlsClientHandshaker, tlsServerHandshaker)

	backendConnecter := pg_proxy.NewBackendConnecter("tcp4", "postgres", "5432")

	connectionHandler := pg_proxy.NewConnectionHandler(backendConnecter, bridgeBuilder)

	server := tcp_server.NewServer("tcp4", "54321", connectionHandler)

	onReady := func() {
		slog.Info("Server ready and listening on port 54321...")
	}

	if err = server.Listen(onReady); err != nil {
		slog.Error("Error starting server", "error", err)
		return
	}
}
