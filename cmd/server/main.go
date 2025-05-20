package main

import (
	"crypto/tls"
	"log"
	"log/slog"
	"os"
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

	pgProxyTlsNegotiator := pg_proxy.NewTlsNegotiator(tlsClientHandshaker, tlsServerHandshaker)

	pgProxyBackendConnecter := pg_proxy.NewBackendConnecter("tcp4", "postgres", "5432")

	pgProxyConnectionHandler := pg_proxy.NewConnectionHandler(pgProxyBackendConnecter, pgProxyTlsNegotiator)

	server, readyCh := tcp_server.New("tcp4", "54321")

	go func() {
		if err = server.Listen(); err != nil {
			slog.Error("Error starting server", "error", err)

			os.Exit(1)
		}
	}()

	addr := <-readyCh
	slog.Info("Server is ready", "address", addr)

	for connEvent := range server.Accept() {
		if connEvent.Err != nil {
			slog.Error("Error accepting connection", "error", connEvent.Err)
			continue
		}

		go pgProxyConnectionHandler.HandleConnection(connEvent.Conn)
	}
}
