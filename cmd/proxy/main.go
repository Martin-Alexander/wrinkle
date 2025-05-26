package main

import (
	"crypto/tls"
	"log"
	"log/slog"
	"os"
	"wrinkle/internal/proxy"
	"wrinkle/internal/tcp"
)

func main() {
	cert, err := tls.LoadX509KeyPair("/.ssl/proxy.crt", "/.ssl/proxy.key")
	if err != nil {
		log.Fatalf("Failed to load cert: %v", err)
	}

	tlsServerHandshaker := tcp.NewServerTlsHandshaker(&tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	})

	tlsClientHandshaker := tcp.NewClientTlsHandshaker(&tls.Config{
		InsecureSkipVerify: true,
	})

	pgProxyTlsNegotiator := proxy.NewTlsNegotiator(tlsClientHandshaker, tlsServerHandshaker)

	pgProxyBackendConnecter := proxy.NewBackendConnecter("tcp4", "pg", "5432")

	pgProxyConnectionHandler := proxy.NewConnectionHandler(pgProxyBackendConnecter, pgProxyTlsNegotiator)

	server, readyCh := tcp.NewServer("tcp4", "5432")

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
