package proxy

import (
	"crypto/tls"
	"net"

	"github.com/pkg/errors"
)

type ConnectionConfig struct {
	BackendHostname   string
	BackendPort       string
	FrontendTlsConfig *tls.Config
	BackendTlsConfig  *tls.Config
}

func HandleConnection(
	feConn net.Conn,
	config ConnectionConfig,
) error {
	defer feConn.Close()

	beConn, err := ConnectToBackend(config.BackendHostname, config.BackendPort)
	if err != nil {
		return errors.WithStack(err)
	}
	defer beConn.Close()

	if err := HandleTlsNegotiation(feConn, beConn); err != nil {
		return errors.WithStack(err)
	}

	feConn, beConn, err = PerformTlsHandshakes(
		feConn,
		beConn,
		config.FrontendTlsConfig,
		config.BackendTlsConfig,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := HandleClientStartupMessage(feConn, beConn); err != nil {
		return errors.WithStack(err)
	}

	relay := NewRelay(feConn, beConn)

	return relay.Start()
}
