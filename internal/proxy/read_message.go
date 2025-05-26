package proxy

import (
	"encoding/binary"
	"io"
	"wrinkle/internal/pg"
)

func ReadMessage(source io.Reader) (*pg.Message, error) {
	typeBuff := make([]byte, 1)
	lengthBuff := make([]byte, 4)

	if _, err := io.ReadFull(source, typeBuff); err != nil {
		return nil, err
	}

	if _, err := io.ReadFull(source, lengthBuff); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lengthBuff)

	dataBuff := make([]byte, length-4)
	if _, err := io.ReadFull(source, dataBuff); err != nil {
		return nil, err
	}

	message := pg.Message{
		Type:   typeBuff[0],
		Length: int32(length),
		Data:   dataBuff,
	}

	return &message, nil
}
