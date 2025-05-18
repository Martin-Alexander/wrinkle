package postgres_proxy

import (
	"io"
	"net"
)

func HandleTlsNegotiation(frontendConn net.Conn, backendConn net.Conn) (bool, error) {
	firstEightBytes, err := readNBytes(frontendConn, 8)
	if err != nil {
		return false, err
	}

	if _, err := backendConn.Write(firstEightBytes); err != nil {
		return false, err
	}

	responseByte, err := readNBytes(backendConn, 1)
	if err != nil {
		return false, err
	}

	if _, err := frontendConn.Write(responseByte); err != nil {
		return false, err
	}

	if responseByte[0] != 'S' {
		return false, nil
	}

	return true, nil
}

func readNBytes(conn net.Conn, n int) ([]byte, error) {
	buffer := make([]byte, n)
	if _, err := io.ReadFull(conn, buffer); err != nil {
		return nil, err
	}

	return buffer, nil
}
