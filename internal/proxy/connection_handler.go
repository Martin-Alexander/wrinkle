package proxy

import (
	"encoding/binary"
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

	lengthBuff := make([]byte, 4)
	if _, err := frontendTlsConn.Read(lengthBuff); err != nil {
		return err
	}

	length := binary.BigEndian.Uint32(lengthBuff)

	messageBuff := make([]byte, length-4)
	if _, err := frontendTlsConn.Read(messageBuff); err != nil {
		return err
	}

	backendTlsConn.Write(lengthBuff)
	backendTlsConn.Write(messageBuff)

	relay := NewRelay(frontendTlsConn, backendTlsConn)

	return relay.Start()
}
