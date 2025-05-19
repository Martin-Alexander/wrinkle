package pg_proxy

import (
	"net"
)

type ConnectionHandler struct {
	backendConnecter *BackendConnecter
	bridgeBuilder    *BridgeBuilder
}

func NewConnectionHandler(
	backendConnecter *BackendConnecter,
	bridgeBuilder *BridgeBuilder,
) *ConnectionHandler {
	return &ConnectionHandler{
		backendConnecter: backendConnecter,
		bridgeBuilder:    bridgeBuilder,
	}
}

func (ch *ConnectionHandler) HandleConnection(conn net.Conn) error {
	defer conn.Close()

	backendConn, err := ch.backendConnecter.Dial()
	if err != nil {
		return err
	}
	defer backendConn.Close()

	bridge, err := ch.bridgeBuilder.Build(conn, backendConn)
	if err != nil {
		return err
	}

	return bridge.StartRelaying()
}
