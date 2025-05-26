package tcp

import (
	"crypto/tls"
	"net"
)

type TlsHandshaker struct {
	server bool
	config *tls.Config
}

func NewClientTlsHandshaker(config *tls.Config) *TlsHandshaker {
	return &TlsHandshaker{
		server: false,
		config: config,
	}
}

func NewServerTlsHandshaker(config *tls.Config) *TlsHandshaker {
	return &TlsHandshaker{
		server: true,
		config: config,
	}
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
