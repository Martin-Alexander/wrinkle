package wrinkle

import (
	"net"

	"github.com/pkg/errors"
)

type ConnectionCreator interface {
	CreateConnection(feConn net.Conn) (net.Conn, net.Conn, error)
}

func HandleConnection(
	feConn net.Conn,
	connectionCreator ConnectionCreator,
	router *Router,
) error {
	feConn, beConn, err := connectionCreator.CreateConnection(feConn)
	if err != nil {
		return errors.WithStack(err)
	}

	router.Start(feConn, beConn)

	return nil
}
