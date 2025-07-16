package pg_wire

import (
	"testing"
	"wrinkle/internal/test_helper"

	"github.com/stretchr/testify/assert"
)

func TestReadMessage(t *testing.T) {
	client, server, closeFn := test_helper.CreateTcpConnection(t)
	defer closeFn()

	message := &Message{
		Type:   byte(ClientSimpleQuery),
		Length: 8,
		Data:   []byte{0x04, 0xd2, 0x16, 0x00},
	}

	if _, err := client.Write(message.Bytes()); err != nil {
		t.Fatal(err)
	}

	message, err := ReadMessage(server)

	assert.Nil(t, err)
	assert.NotNil(t, message)

	assert.Equal(t, byte(ClientSimpleQuery), message.Type)
	assert.Equal(t, uint32(8), message.Length)
	assert.Equal(t, []byte{0x04, 0xd2, 0x16, 0x00}, message.Data)
}

func TestReadMessageNoMessage(t *testing.T) {
	_, server, closeFn := test_helper.CreateTcpConnection(t)
	defer closeFn()

	message, err := ReadMessage(server)

	assert.Nil(t, message)
	assert.Nil(t, err)
}

func TestReadMessageInvalidLengthTooLong(t *testing.T) {
	client, server, closeFn := test_helper.CreateTcpConnection(t)
	defer closeFn()

	message := &Message{
		Type:   byte(ClientSimpleQuery),
		Length: 10,
		Data:   []byte{0x04, 0xd2, 0x16, 0x00},
	}

	client.Write(message.Bytes())

	message, err := ReadMessage(server)

	assert.Error(t, err)
	assert.Nil(t, message)
}

func TestReadMessageInvalidLengthTooShort(t *testing.T) {
	client, server, closeFn := test_helper.CreateTcpConnection(t)
	defer closeFn()

	message := &Message{
		Type:   byte(ClientSimpleQuery),
		Length: 2, // Cannot be less that four
		Data:   []byte{0x04, 0xd2, 0x16, 0x00},
	}

	client.Write(message.Bytes())

	message, err := ReadMessage(server)

	assert.Error(t, err)
	assert.Nil(t, message)
}

func TestWriteMessage(t *testing.T) {
	client, server, closeFn := test_helper.CreateTcpConnection(t)
	defer closeFn()

	message := Message{
		Type:   byte(ClientSimpleQuery),
		Length: 8,
		Data:   []byte{0x04, 0xd2, 0x16, 0x00},
	}

	err := WriteMessage(client, &message)
	if err != nil {
		t.Fatal(err)
	}

	data := make([]byte, 9)
	_, err = server.Read(data)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte{byte(ClientSimpleQuery), 0, 0, 0, 8, 0x04, 0xd2, 0x16, 0x00}

	assert.Equal(t, expected, data)
}
