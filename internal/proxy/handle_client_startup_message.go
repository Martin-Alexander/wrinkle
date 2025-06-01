package proxy

import (
	"encoding/binary"
	"io"
	"net"
	"slices"

	"github.com/pkg/errors"
)

func HandleClientStartupMessage(feConn net.Conn, beConn net.Conn) error {
	lengthBuff := make([]byte, 4)
	if _, err := io.ReadFull(feConn, lengthBuff); err != nil {
		return errors.WithStack(err)
	}

	length := binary.BigEndian.Uint32(lengthBuff)

	messageBuff := make([]byte, length-4)
	if _, err := io.ReadFull(feConn, messageBuff); err != nil {
		return errors.WithStack(err)
	}

	beConn.Write(slices.Concat(lengthBuff, messageBuff))

	return nil
}
