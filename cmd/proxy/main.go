package main

import (
	"crypto/tls"
	"log/slog"
	"os"
	"wrinkle/internal/pg_middleman"
	"wrinkle/internal/pg_wire"
	"wrinkle/internal/tcp"
	"wrinkle/internal/wrinkle"
)

func main() {
	slog.Debug("Starting the proxy server")

	cert, err := tls.LoadX509KeyPair("/.ssl/proxy.crt", "/.ssl/proxy.key")
	if err != nil {
		slog.Error("Error loading TLS certificates", "error", err)
	}

	connectionCreator := pg_middleman.NewConnectionCreator(
		"pg",
		"5432",
		&tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		},
		&tls.Config{
			InsecureSkipVerify: true,
		},
	)

	controller := wrinkle.NewController()

	errorCh := make(chan error, 1)

	broker := wrinkle.NewBroker(
		&pg_wire.MessageReader{},
		&pg_wire.MessageWriter{},
		controller,
		errorCh,
	)

	controller.Start()

	server, readyCh := tcp.NewServer("tcp4", "5432")

	go func() {
		if err := server.Listen(); err != nil {
			slog.Error("TCP server error", "error", err)

			os.Exit(1)
		}
	}()

	go func() {
		for err := range errorCh {
			slog.Error("Broker error", "error", err)

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

		slog.Info("Accepted new connection", "address", connEvent.Conn.RemoteAddr())

		if err := wrinkle.HandleConnection(connEvent.Conn, connectionCreator, broker); err != nil {
			slog.Error("Connection handling error", "error", err)
		}
	}
}
