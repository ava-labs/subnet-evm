// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package predicate

import (
	"bytes"
	"testing"

	"github.com/ava-labs/avalanchego/utils"
	"github.com/stretchr/testify/require"
)

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

func FuzzHashSliceToBytes(f *testing.F) {
	for i := 0; i < 100; i++ {
		f.Add(utils.RandomBytes(i))
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		testBytesToHashSlice(t, b)
	})
}

func testPackPredicate(t testing.TB, b []byte) {
	packedPredicate := PackPredicate(b)
	unpackedPredicated, err := UnpackPredicate(packedPredicate)
	require.NoError(t, err)
	require.Equal(t, b, unpackedPredicated)
}

func FuzzPackPredicate(f *testing.F) {
	for i := 0; i < 100; i++ {
		f.Add(utils.RandomBytes(i))
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		testPackPredicate(t, b)
	})
}

func TestUnpackInvalidPredicate(t *testing.T) {
	require := require.New(t)
<<<<<<< HEAD
=======
	// Predicate encoding requires a 0xff delimiter byte followed by padding of all zeroes, so any other
	// excess padding should invalidate the predicate.
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	paddingCases := make([][]byte, 0, 200)
	for i := 1; i < 100; i++ {
		paddingCases = append(paddingCases, bytes.Repeat([]byte{0xee}, i))
		paddingCases = append(paddingCases, make([]byte, i))
	}

	for _, l := range []int{0, 1, 31, 32, 33, 63, 64, 65} {
		validPredicate := PackPredicate(utils.RandomBytes(l))

		for _, padding := range paddingCases {
			invalidPredicate := append(validPredicate, padding...)
			_, err := UnpackPredicate(invalidPredicate)
			require.Error(err, "Predicate length %d, Padding length %d (0x%x)", len(validPredicate), len(padding), invalidPredicate)
		}
	}
}
