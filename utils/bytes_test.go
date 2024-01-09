// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"bytes"
	"math/big"
	"math/rand"
	"testing"

	"github.com/ava-labs/avalanchego/utils"
	"github.com/ava-labs/avalanchego/utils/wrappers"
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

	copiedBytes := HashSliceToBytes(hashSlice)

	if len(b)%32 == 0 {
		require.Equal(t, b, copiedBytes)
	} else {
		require.Equal(t, b, copiedBytes[:len(b)])
		// Require that any additional padding is all zeroes
		padding := copiedBytes[len(b):]
		require.Equal(t, bytes.Repeat([]byte{0x00}, len(padding)), padding)
	}
}

func testBytesToBigInt(t testing.TB, b []byte) {
	p := wrappers.Packer{
		Bytes: append([]byte{BigIntPositive}[:], b...),
	}
	number, err := UnpackBigInt(&p)
	if len(b) < 32 {
		require.Error(t, err, "")
	} else {
		require.Nil(t, err)
		buf := make([]byte, 32)
		require.Equal(t, number.FillBytes(buf), b[0:32])
	}
}

func testBigInt(t testing.TB, rawNumber int64) {
	p := wrappers.Packer{
		Bytes:   []byte{},
		MaxSize: 100,
	}
	number := big.NewInt(rawNumber)
	err := PackBigInt(&p, number)
	require.Nil(t, err)
	if rawNumber == 0 {
		require.Equal(t, 1, len(p.Bytes))
	} else {
		require.Equal(t, 33, len(p.Bytes))
	}
	newNumber, err := UnpackBigInt(&wrappers.Packer{Bytes: p.Bytes})
	require.Nil(t, err)
	require.Equal(t, number, newNumber)
}

func testBytesToBigIntNil(t testing.TB, b []byte) {
	p := wrappers.Packer{
		Bytes: append([]byte{BigIntNil}[:], b...),
	}
	number, err := UnpackBigInt(&p)
	require.Nil(t, err)
	require.Nil(t, number)
}

func FuzzHashSliceToBytes(f *testing.F) {
	for i := 0; i < 100; i++ {
		f.Add(utils.RandomBytes(i))
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		testBytesToHashSlice(t, b)
	})
}

func FuzzBigInt(f *testing.F) {
	f.Add(int64(0))
	for i := 0; i < 200; i++ {
		f.Add(int64(rand.Uint64()))
		f.Add(-1 * int64(rand.Uint64()))
	}

	f.Fuzz(func(t *testing.T, b int64) {
		testBigInt(t, b)
	})

}

func FuzzSliceOfBytesToBigInt(f *testing.F) {
	for i := 0; i < 200; i++ {
		f.Add(utils.RandomBytes(i))
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		testBytesToBigInt(t, b)
		testBytesToBigIntNil(t, b)
	})
}
