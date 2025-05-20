package pg_proxy

import (
	"net"
)

type ConnectionHandler struct {
	backendConnecter *BackendConnecter
	tlsNegotiator    *TlsNegotiator
}

func NewConnectionHandler(
	backendConnecter *BackendConnecter,
	tlsNegotiator *TlsNegotiator,
) *ConnectionHandler {
	return &ConnectionHandler{
		backendConnecter: backendConnecter,
		tlsNegotiator:    tlsNegotiator,
	}
}

func (ch *ConnectionHandler) HandleConnection(conn net.Conn) error {
	defer conn.Close()

	backendConn, err := ch.backendConnecter.Dial()
	if err != nil {
		return err
	}
	defer backendConn.Close()

	frontendTlsConn, backendTlsConn, err := ch.tlsNegotiator.Negotiate(conn, backendConn)
	if err != nil {
		return err
	}

	relay := NewRelay(frontendTlsConn, backendTlsConn)

	return relay.Start()
}
