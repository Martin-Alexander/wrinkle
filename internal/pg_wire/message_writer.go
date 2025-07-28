package pg_wire

import (
	"io"

	"github.com/pkg/errors"
)

type MessageWriter struct{}

func (mw *MessageWriter) WriteMessage(writer io.Writer, message *Message) error {
	if _, err := writer.Write(message.Bytes()); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
