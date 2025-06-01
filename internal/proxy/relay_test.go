package proxy

import (
	"testing"
	"wrinkle/internal/pg_wire"
	"wrinkle/internal/test_helper"

	"github.com/stretchr/testify/assert"
)

func TestNewRelay(t *testing.T) {
	_, feServer, feCloseFn := test_helper.CreateTcpConnection(t)
	beClient, _, beCloseFn := test_helper.CreateTcpConnection(t)
	defer feCloseFn()
	defer beCloseFn()

	relay := NewRelay(feServer, beClient)

	assert.NotNil(t, relay)
	assert.Equal(t, feServer, relay.frontendConn)
	assert.Equal(t, beClient, relay.backendConn)
}

func TestStart(t *testing.T) {
	feClient, feServer, feCloseFn := test_helper.CreateTcpConnection(t)
	beClient, beServer, beCloseFn := test_helper.CreateTcpConnection(t)
	defer feCloseFn()
	defer beCloseFn()

	relay := NewRelay(feServer, beClient)

	go relay.Start()

	feMessage := pg_wire.Message{
		Type:   byte(pg_wire.ClientSimpleQuery),
		Length: 8,
		Data:   []byte{0x04, 0xd2, 0x16, 0x2f},
	}
	feClient.Write(feMessage.Bytes())
	buffer := make([]byte, len(feMessage.Bytes()))
	beServer.Read(buffer)

	assert.Equal(t, string(feMessage.Bytes()), string(buffer))

	beMessage := pg_wire.Message{
		Type:   byte(pg_wire.ServerCommandComplete),
		Length: 8,
		Data:   []byte{0x04, 0xd2, 0x16, 0x2f},
	}
	beServer.Write(beMessage.Bytes())
	beBuffer := make([]byte, len(beMessage.Bytes()))
	n, _ := feClient.Read(beBuffer)

	assert.Equal(t, string(beMessage.Bytes()), string(beBuffer[:n]))
}
