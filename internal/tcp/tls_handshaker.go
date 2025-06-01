package tcp

import (
	"crypto/tls"
	"net"
)

func NewTlsHandshaker(config *tls.Config, server bool) *TlsHandshaker {
	return &TlsHandshaker{
		server: server,
		config: config,
	}
}

type TlsHandshaker struct {
	server bool
	config *tls.Config
}

func (h *TlsHandshaker) Handshake(conn net.Conn) (net.Conn, error) {
	var tlsConn *tls.Conn

	if h.server {
		tlsConn = tls.Server(conn, h.config)
	} else {
		tlsConn = tls.Client(conn, h.config)
	}

	if err := tlsConn.Handshake(); err != nil {
		return nil, err
	}
	return tlsConn, nil
}
