package pg_wire

import (
	"encoding/binary"
)

type Sender int

const (
	Frontend Sender = iota
	Backend
)

type Message struct {
	Sender Sender
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

func (m *Message) Name() string {
	if m.Sender == Frontend {
		return ClientMessageType(m.Type).String()
	}

	return ServerMessageType(m.Type).String()
}
