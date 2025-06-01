package main

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"os"
	"wrinkle/internal/proxy"
	"wrinkle/internal/tcp"

	"github.com/pkg/errors"
)

func main() {
	cert, err := tls.LoadX509KeyPair("/.ssl/proxy.crt", "/.ssl/proxy.key")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %+v\n", errors.WithStack(err))
	}

	connectionConfig := proxy.ConnectionConfig{
		BackendHostname: "pg",
		BackendPort:     "5432",
		FrontendTlsConfig: &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		},
		BackendTlsConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	server, readyCh := tcp.NewServer("tcp4", "5432")

	go func() {
		if err := server.Listen(); err != nil {
			fmt.Fprintf(os.Stderr, "Error %+v\n", err)

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

		go func() {
			if err := proxy.HandleConnection(connEvent.Conn, connectionConfig); err != nil {
				fmt.Fprintf(os.Stderr, "Error %+v\n", err)
			}
		}()
	}
}
