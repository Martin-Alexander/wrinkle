package proxy

import (
	"net"
)

type BackendConnecter struct {
	network string
	host    string
	port    string
}

func NewBackendConnecter(network, host string, port string) *BackendConnecter {
	return &BackendConnecter{
		network: network,
		host:    host,
		port:    port,
	}
}

func (bc *BackendConnecter) Dial() (net.Conn, error) {
	address := net.JoinHostPort(bc.host, bc.port)
	conn, err := net.Dial(bc.network, address)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
