// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"testing"

	"github.com/ava-labs/avalanchego/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncrOne(t *testing.T) {
	type test struct {
		input    []byte
		expected []byte
	}
	for name, test := range map[string]test{
		"increment no overflow no carry": {
			input:    []byte{0, 0},
			expected: []byte{0, 1},
		},
		"increment overflow": {
			input:    []byte{255, 255},
			expected: []byte{0, 0},
		},
		"increment carry": {
			input:    []byte{0, 255},
			expected: []byte{1, 0},
		},
	} {
		t.Run(name, func(t *testing.T) {
			output := common.CopyBytes(test.input)
			IncrOne(output)
			assert.Equal(t, output, test.expected)
		})
	}
}

func testBytesToHashSlice(t testing.TB, b []byte) {
	hashSlice := BytesToHashSlice(b)

	copiedBytes, ok := HashSliceToBytes(hashSlice)
	require.True(t, ok)
	require.Equal(t, b, copiedBytes)
}

func FuzzHashSliceToBytes(f *testing.F) {
	for i := 0; i < 100; i++ {
		f.Add(utils.RandomBytes(i))
	}

	f.Fuzz(func(t *testing.T, a []byte) {
		testBytesToHashSlice(t, a)
	})
}
