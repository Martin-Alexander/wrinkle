package proxy

import (
	"net"

	"github.com/pkg/errors"
)

func ConnectToBackend(host string, port string) (net.Conn, error) {
	address := net.JoinHostPort(host, port)
	conn, err := net.Dial("tcp4", address)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return conn, nil
}
