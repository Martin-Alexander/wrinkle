package pg_wire

import (
	"io"

	"github.com/pkg/errors"
)

type MessageWriter struct{}

func (mw *MessageWriter) WriteMessage(writer io.Writer, message []byte) error {
	if _, err := writer.Write(message); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
