package handshake

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type Object2 struct {
	Field map[string]string `json:"field"`
}

type Object struct {
	Field  map[string]string `json:"b-field"`
	FieldB Object2           `json:"a-field"`
	FieldA int32             `json:"Z-field"`
}

func TestDeterministicJsonEncoding(t *testing.T) {
	object := Object{
		FieldA: 1,
		FieldB: Object2{
			Field: map[string]string{
				"test": "a",
				"z":    "a",
				"a":    "a",
			},
		},
		Field: map[string]string{
			"test": "a",
			"z":    "a",
			"a":    "a",
		},
	}

	regularBytes, err := json.Marshal(object)
	require.NoError(t, err)

	bytes, err := DeterministicJsonEncoding(object)
	require.NoError(t, err)
	require.NotEqual(t, regularBytes, bytes)
	require.Equal(
		t,
		"{\"Z-field\":1,\"a-field\":{\"field\":{\"a\":\"a\",\"test\":\"a\",\"z\":\"a\"}},\"b-field\":{\"a\":\"a\",\"test\":\"a\",\"z\":\"a\"}}",
		string(bytes[:]),
	)
}
