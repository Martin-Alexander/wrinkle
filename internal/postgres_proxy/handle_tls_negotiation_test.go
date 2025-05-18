package postgres_proxy

import (
	"errors"
	"net"
	"os"
	"testing"
	"time"
)

func TestHandleTlsNegotiation(t *testing.T) {
	testCases := []struct {
		name            string
		frontendPacket  []byte
		backendResponse []byte
		expectedResult  bool
		expectedError   error
	}{
		{
			name:            "Successful handshake",
			frontendPacket:  []byte{0x00, 0x00, 0x00, 0x08, 0x04, 0xd2, 0x16, 0x2f},
			backendResponse: []byte("S"),
			expectedResult:  true,
			expectedError:   nil,
		},
		{
			name:            "Backend rejects with N response",
			frontendPacket:  []byte{0x00, 0x00, 0x00, 0x08, 0x04, 0xd2, 0x16, 0x2f},
			backendResponse: []byte("N"),
			expectedResult:  false,
			expectedError:   nil,
		},
		{
			name:            "Malformed frontend packet",
			frontendPacket:  []byte{0x00, 0x00, 0x00},
			backendResponse: nil,
			expectedResult:  false,
			expectedError:   os.ErrDeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			frontendClient, frontendServer := net.Pipe()
			backendClient, backendServer := net.Pipe()
			setIoTimeout(t, frontendClient)
			setIoTimeout(t, frontendServer)
			setIoTimeout(t, backendClient)
			setIoTimeout(t, backendServer)
			defer frontendClient.Close()
			defer frontendServer.Close()
			defer backendClient.Close()
			defer backendServer.Close()

			resultChan := make(chan bool, 1)
			errorChan := make(chan error, 1)

			go func() {
				result, err := HandleTlsNegotiation(frontendServer, backendClient)
				resultChan <- result
				errorChan <- err
			}()

			frontendClient.Write(tc.frontendPacket)

			buffer := make([]byte, len(tc.frontendPacket))
			backendServer.Read(buffer)
			if string(buffer) != string(tc.frontendPacket) {
				t.Errorf("Expected %s but got %s", string(tc.frontendPacket), string(buffer))
			}

			backendServer.Write(tc.backendResponse)

			buffer = make([]byte, len(tc.backendResponse))
			frontendClient.Read(buffer)
			if string(buffer) != string(tc.backendResponse) {
				t.Errorf("Expected %s but got %s", string(tc.backendResponse), string(buffer))
			}

			result := <-resultChan
			if result != tc.expectedResult {
				t.Errorf("Expected %v but got %v", tc.expectedResult, result)
			}
			err := <-errorChan
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("Expected %v but got %v", tc.expectedError, err)
			}
		})
	}
}

func setIoTimeout(t *testing.T, conn net.Conn) {
	t.Helper()
	if err := conn.SetDeadline(time.Now().Add(1 * time.Millisecond)); err != nil {
		t.Fatal(err)
	}
}
