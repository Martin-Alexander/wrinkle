package pg_wire

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

type MessageReader struct{}

func (mr *MessageReader) ReadMessage(reader io.Reader) ([]byte, error) {
	messageTypeBuff, err := read(reader, 1)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	messageLengthBuff, err := read(reader, 4)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	messageLength := binary.BigEndian.Uint32(messageLengthBuff)

	dataSectionLength := int(messageLength - 4)

	if dataSectionLength < 0 {
		return nil, errors.Errorf("invalid message length: %d", messageLength)
	}

	dataBuff, err := read(reader, dataSectionLength)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	message := Message{
		Type:   messageTypeBuff[0],
		Length: messageLength,
		Data:   dataBuff,
	}

	return message.Bytes(), nil
}

func read(reader io.Reader, length int) ([]byte, error) {
	data := make([]byte, length)
	_, err := io.ReadFull(reader, data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}
