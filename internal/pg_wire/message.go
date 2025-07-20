package pg_wire

import (
	"encoding/binary"
	"errors"
)

func FromBytes(data []byte) (*Message, error) {
	if len(data) < 5 {
		return nil, errors.New("data too short to be a valid message")
	}

	messageType := data[0]
	messageLength := binary.BigEndian.Uint32(data[1:5])

	if uint32(len(data)) < messageLength {
		return nil, errors.New("data length does not match message length")
	}

	return &Message{
		Type:   messageType,
		Length: messageLength,
		Data:   data[5 : messageLength+1],
	}, nil
}

type Message struct {
	Type   byte
	Length uint32
	Data   []byte
}

func (m *Message) Bytes() []byte {
	buffer := make([]byte, 5+len(m.Data))

	buffer[0] = m.Type
	binary.BigEndian.PutUint32(buffer[1:], m.Length)
	copy(buffer[5:], m.Data)

	return buffer
}
