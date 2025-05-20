package pg_proxy

import (
	"crypto/tls"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHandshaker struct {
	mock.Mock
}

func (m *MockHandshaker) Handshake(conn net.Conn) (net.Conn, error) {
	args := m.Called(conn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(net.Conn), args.Error(1)
}

type Result struct {
	frontendConn net.Conn
	backendConn  net.Conn
	err          error
}

func TestNewTlsNegotiator(t *testing.T) {
	tlsClientHandshaker := new(MockHandshaker)
	tlsServerHandshaker := new(MockHandshaker)

	tlsNegotiator := NewTlsNegotiator(tlsClientHandshaker, tlsServerHandshaker)

	assert.NotNil(t, tlsNegotiator)
	assert.Equal(t, tlsClientHandshaker, tlsNegotiator.clientHandshaker)
	assert.Equal(t, tlsServerHandshaker, tlsNegotiator.serverHandshaker)
}

func TestNegotiate(t *testing.T) {
	testCases := []struct {
		name             string
		frontendPacket   []byte
		backendResponse  []byte
		handshakeSuccess bool
		expectResult     bool
		expectError      bool
	}{
		{
			name:             "Successful negotiation",
			frontendPacket:   []byte{0x00, 0x00, 0x00, 0x08, 0x04, 0xd2, 0x16, 0x2f},
			backendResponse:  []byte("S"),
			handshakeSuccess: true,
			expectResult:     true,
			expectError:      false,
		},
		{
			name:             "Successful negotiation but unsuccessful handshake",
			frontendPacket:   []byte{0x00, 0x00, 0x00, 0x08, 0x04, 0xd2, 0x16, 0x2f},
			backendResponse:  []byte("S"),
			handshakeSuccess: false,
			expectResult:     false,
			expectError:      true,
		},
		{
			name:             "Backend rejects with N response",
			frontendPacket:   []byte{0x00, 0x00, 0x00, 0x08, 0x04, 0xd2, 0x16, 0x2f},
			backendResponse:  []byte("N"),
			handshakeSuccess: false,
			expectResult:     false,
			expectError:      true,
		},
		{
			name:             "Malformed frontend packet",
			frontendPacket:   []byte{0x00, 0x00, 0x00},
			backendResponse:  nil,
			handshakeSuccess: false,
			expectResult:     false,
			expectError:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tlsClientHandshakerMock := new(MockHandshaker)
			tlsServerHandshakerMock := new(MockHandshaker)

			frontendTlsConn := &tls.Conn{}
			backendTlsConn := &tls.Conn{}

			negotiator := NewTlsNegotiator(tlsClientHandshakerMock, tlsServerHandshakerMock)

			frontendClient, frontendServer,
				backendClient, backendServer, closeConnections := createConnections(t)
			defer closeConnections()

			if tc.handshakeSuccess {
				tlsClientHandshakerMock.On("Handshake", mock.Anything).Return(frontendTlsConn, nil)
				tlsServerHandshakerMock.On("Handshake", mock.Anything).Return(backendTlsConn, nil)
			} else {
				tlsClientHandshakerMock.On("Handshake", mock.Anything).Return(nil, errors.New("handshake error"))
				tlsServerHandshakerMock.On("Handshake", mock.Anything).Return(nil, errors.New("handshake error"))
			}

			resultCh := make(chan *Result, 1)

			go func() {
				frontendConn, backendConn, err := negotiator.Negotiate(frontendServer, backendClient)
				resultCh <- &Result{
					frontendConn: frontendConn,
					backendConn:  backendConn,
					err:          err,
				}
			}()

			frontendClient.Write(tc.frontendPacket)
			buffer := make([]byte, len(tc.frontendPacket))
			backendServer.Read(buffer)

			assert.Equal(t, string(tc.frontendPacket), string(buffer))

			backendServer.Write(tc.backendResponse)
			buffer = make([]byte, len(tc.backendResponse))
			frontendClient.Read(buffer)

			assert.Equal(t, string(tc.backendResponse), string(buffer))

			result := <-resultCh

			if tc.expectError {
				assert.Error(t, result.err)
			} else {
				assert.NoError(t, result.err)
			}

			if tc.expectResult {
				assert.NotNil(t, result)
				assert.Equal(t, frontendTlsConn, result.frontendConn)
				assert.Equal(t, backendTlsConn, result.backendConn)
			} else {
				assert.Nil(t, result.frontendConn)
				assert.Nil(t, result.backendConn)
			}
		})
	}
}

func createConnections(t *testing.T) (net.Conn, net.Conn, net.Conn, net.Conn, func()) {
	t.Helper()

	frontendClient, frontendServer := net.Pipe()
	backendClient, backendServer := net.Pipe()
	timeout := time.Now().Add(time.Millisecond)
	frontendClient.SetDeadline(timeout)
	frontendServer.SetDeadline(timeout)
	backendClient.SetDeadline(timeout)
	backendServer.SetDeadline(timeout)

	teardown := func() {
		frontendClient.Close()
		frontendServer.Close()
		backendClient.Close()
		backendServer.Close()
	}

	return frontendClient, frontendServer, backendClient, backendServer, teardown
}
