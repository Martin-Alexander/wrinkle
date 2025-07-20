package pg_middleman

import (
	"net"

	"github.com/pkg/errors"
)

type TlsNegotiationError struct {
	message string
}

func (e *TlsNegotiationError) Error() string {
	return e.message
}

func HandleTlsPreNegotiation(feConn net.Conn, beConn net.Conn) error {
	initMessageBuff := make([]byte, 8)
	if _, err := feConn.Read(initMessageBuff); err != nil {
		return errors.WithStack(err)
	}

	if _, err := beConn.Write(initMessageBuff); err != nil {
		return errors.WithStack(err)
	}

	responseBuff := make([]byte, 1)
	if _, err := beConn.Read(responseBuff); err != nil {
		return errors.WithStack(err)
	}

	if _, err := feConn.Write(responseBuff); err != nil {
		return errors.WithStack(err)
	}

	if responseBuff[0] != 'S' {
		err := &TlsNegotiationError{message: "Backend rejection"}

		return errors.WithStack(err)
	}

	return nil
}
