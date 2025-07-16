package pg_wire

import (
	"encoding/binary"
	"net"
	"time"

	"github.com/pkg/errors"
)

func ReadMessage(source net.Conn) (*Message, error) {
	source.SetReadDeadline(time.Now().Add(time.Millisecond))
	typeBuff := make([]byte, 1)
	_, err := source.Read(typeBuff)
	source.SetReadDeadline(time.Time{})

	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}

	lengthBuff, lengthBuffErr := readExactly(source, 4)
	if lengthBuffErr != nil {
		return nil, errors.WithStack(lengthBuffErr)
	}

	length := binary.BigEndian.Uint32(lengthBuff)

	dataLength := int(length - 4)

	if dataLength < 0 {
		return nil, errors.Errorf("invalid message length: %d", length)
	}

	dataBuff, dataBuffErr := readExactly(source, dataLength)
	if dataBuffErr != nil {
		return nil, errors.WithStack(dataBuffErr)
	}

	message := Message{
		Type:   typeBuff[0],
		Length: length,
		Data:   dataBuff,
	}

	return &message, nil
}

func WriteMessage(destination net.Conn, message *Message) error {
	if _, err := destination.Write(message.Bytes()); err != nil {
		return errors.WithStack(err)
	}

	return nil
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

func readExactly(source net.Conn, length int) ([]byte, error) {
	data := make([]byte, length)

	n, err := source.Read(data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if n != length {
		return nil, errors.Errorf("expected to read %d bytes, but got %d", length, n)
	}

	return data, nil
}
