package pg_middleman

import (
	"crypto/tls"
	"net"

	"github.com/pkg/errors"
)

func PerformTlsHandshakes(
	feConn net.Conn,
	beConn net.Conn,
	feTlsConfig *tls.Config,
	beTlsConfig *tls.Config,
) (net.Conn, net.Conn, error) {
	tlsConn := tls.Server(feConn, feTlsConfig)
	if err := tlsConn.Handshake(); err != nil {
		return nil, nil, errors.WithStack(err)
	}

	beTlsConn := tls.Client(beConn, beTlsConfig)
	if err := beTlsConn.Handshake(); err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return tlsConn, beTlsConn, nil
}
