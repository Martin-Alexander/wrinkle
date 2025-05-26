package proxy

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

type TlsNegotiator struct {
	clientHandshaker Handshaker
	serverHandshaker Handshaker
}

func NewTlsNegotiator(
	tlsClientHandshaker Handshaker,
	tlsServerHandshaker Handshaker,
) *TlsNegotiator {
	return &TlsNegotiator{
		clientHandshaker: tlsClientHandshaker,
		serverHandshaker: tlsServerHandshaker,
	}
}

func (t *TlsNegotiator) Negotiate(
	frontendConn net.Conn,
	backendConn net.Conn,
) (net.Conn, net.Conn, error) {
	firstEightBytes, err := readNBytes(frontendConn, 8)
	if err != nil {
		return nil, nil, err
	}

	if _, err := backendConn.Write(firstEightBytes); err != nil {
		return nil, nil, err
	}

	responseByte, err := readNBytes(backendConn, 1)
	if err != nil {
		return nil, nil, err
	}

	if _, err := frontendConn.Write(responseByte); err != nil {
		return nil, nil, err
	}

	if responseByte[0] != 'S' {
		return nil, nil, &TlsNegotiationError{
			message: "Backend rejection",
		}
	}

	frontendTlsConn, err := t.serverHandshaker.Handshake(frontendConn)
	if err != nil {
		return nil, nil, err
	}

	backendTlsConn, err := t.clientHandshaker.Handshake(backendConn)
	if err != nil {
		return nil, nil, err
	}

	return frontendTlsConn, backendTlsConn, nil
}

func readNBytes(conn net.Conn, n int) ([]byte, error) {
	buffer := make([]byte, n)
	if _, err := io.ReadFull(conn, buffer); err != nil {
		return nil, err
	}

	return buffer, nil
}
