package pg_middleman

import (
	"crypto/tls"
	"net"

	"github.com/pkg/errors"
)

type ConnectionCreator struct {
	BackendHostname   string
	BackendPort       string
	FrontendTlsConfig *tls.Config
	BackendTlsConfig  *tls.Config
}

func NewConnectionCreator(
	BackendHostname string,
	BackendPort string,
	FrontendTlsConfig *tls.Config,
	BackendTlsConfig *tls.Config,
) *ConnectionCreator {
	return &ConnectionCreator{
		BackendHostname:   BackendHostname,
		BackendPort:       BackendPort,
		FrontendTlsConfig: FrontendTlsConfig,
		BackendTlsConfig:  BackendTlsConfig,
	}
}

func (c *ConnectionCreator) CreateConnection(feConn net.Conn) (net.Conn, net.Conn, error) {
	address := net.JoinHostPort(c.BackendHostname, c.BackendPort)
	beConn, err := net.Dial("tcp4", address)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	if err := HandleTlsPreNegotiation(feConn, beConn); err != nil {
		return nil, nil, errors.WithStack(err)
	}

	feConn, beConn, err = PerformTlsHandshakes(
		feConn,
		beConn,
		c.FrontendTlsConfig,
		c.BackendTlsConfig,
	)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	if err := HandleClientStartupMessage(feConn, beConn); err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return feConn, beConn, nil
}
