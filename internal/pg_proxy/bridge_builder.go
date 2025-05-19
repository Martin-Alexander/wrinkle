package pg_proxy

import (
	"io"
	"net"
)

type TlsNegotiationError struct {
	message string
}

func (e *TlsNegotiationError) Error() string {
	return e.message
}

type Handshaker interface {
	Handshake(net.Conn) (net.Conn, error)
}

type BridgeBuilder struct {
	clientHandshaker Handshaker
	serverHandshaker Handshaker
}

func NewBridgeBuilder(
	tlsClientHandshaker Handshaker,
	tlsServerHandshaker Handshaker,
) *BridgeBuilder {
	return &BridgeBuilder{
		clientHandshaker: tlsClientHandshaker,
		serverHandshaker: tlsServerHandshaker,
	}
}

func (b *BridgeBuilder) Build(
	frontendConn net.Conn,
	backendConn net.Conn,
) (*Bridge, error) {
	firstEightBytes, err := readNBytes(frontendConn, 8)
	if err != nil {
		return nil, err
	}

	if _, err := backendConn.Write(firstEightBytes); err != nil {
		return nil, err
	}

	responseByte, err := readNBytes(backendConn, 1)
	if err != nil {
		return nil, err
	}

	if _, err := frontendConn.Write(responseByte); err != nil {
		return nil, err
	}

	if responseByte[0] != 'S' {
		return nil, &TlsNegotiationError{
			message: "Backend rejection",
		}
	}

	frontendTlsConn, err := b.serverHandshaker.Handshake(frontendConn)
	if err != nil {
		return nil, err
	}

	backendTlsConn, err := b.clientHandshaker.Handshake(backendConn)
	if err != nil {
		return nil, err
	}

	bridge := NewBridge(frontendTlsConn, backendTlsConn)

	return bridge, nil
}

func readNBytes(conn net.Conn, n int) ([]byte, error) {
	buffer := make([]byte, n)
	if _, err := io.ReadFull(conn, buffer); err != nil {
		return nil, err
	}

	return buffer, nil
}
