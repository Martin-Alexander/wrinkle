package proxy

import (
	"bytes"
	"testing"
	"wrinkle/internal/pg"

	"github.com/stretchr/testify/assert"
)

func TestReadMessage(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected pg.Message
	}{
		{
			name: "Valid message",
			input: []byte{
				byte(pg.ClientSimpleQuery),
				0x00, 0x00, 0x00, 0x09,
				0x04, 0xd2, 0x16, 0x2f, 0x00,
			},
			expected: pg.Message{
				Type:   byte(pg.ClientSimpleQuery),
				Length: 9,
				Data:   []byte{0x04, 0xd2, 0x16, 0x2f, 0x00},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			source := bytes.NewReader(tc.input)

			result, err := ReadMessage(source)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assert.Equal(t, tc.expected.Type, result.Type)
			assert.Equal(t, tc.expected.Length, result.Length)

			if !bytes.Equal(result.Data, tc.expected.Data) {
				t.Fatalf("expected %v, got %v", tc.expected.Data, result.Data)
			}
		})
	}
}
