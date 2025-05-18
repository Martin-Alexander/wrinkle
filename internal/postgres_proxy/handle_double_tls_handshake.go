package postgres_proxy

import (
	"crypto/tls"
	"net"
)

func HandleDoubleTlsHandshake(
	frontendConn net.Conn,
	backendConn net.Conn,
	frontendConfig *tls.Config,
	backendConfig *tls.Config,
) (*tls.Conn, *tls.Conn, error) {

	frontendTlsCon, err := frontendTlsHandshake(frontendConn, frontendConfig)
	if err != nil {
		return nil, nil, err
	}

	backendTlsCon, err := backendTlsHandshake(backendConn, backendConfig)
	if err != nil {
		return nil, nil, err
	}

	return frontendTlsCon, backendTlsCon, nil
}

func frontendTlsHandshake(conn net.Conn, config *tls.Config) (*tls.Conn, error) {
	tlsConn := tls.Server(conn, config)
	err := tlsConn.Handshake()
	if err != nil {
		return nil, err
	}
	return tlsConn, nil
}

func backendTlsHandshake(conn net.Conn, config *tls.Config) (*tls.Conn, error) {
	tlsConn := tls.Client(conn, config)
	err := tlsConn.Handshake()
	if err != nil {
		return nil, err
	}
	return tlsConn, nil
}
