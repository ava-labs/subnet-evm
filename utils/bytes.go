// (c) 2021-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
)

// IncrOne increments bytes value by one
func IncrOne(bytes []byte) {
	index := len(bytes) - 1
	for index >= 0 {
		if bytes[index] < 255 {
			bytes[index]++
			break
		} else {
			bytes[index] = 0
			index--
		}
	}
}

const hashSlicePrefixByte byte = 0xff

// zeroStrippedSlice returns a sub-slice of input with all leading zero bytes stripped
func zeroStrippedSlice(input []byte) []byte {
	zeroStrippedBytes := input
	for i, b := range zeroStrippedBytes {
		if b != 0 {
			return zeroStrippedBytes[i:]
		}
	}
	return zeroStrippedBytes
}

// HashSliceToBytes serializes a []common.Hash into a byte slice
// Strips all zero padding from the first hash and [hashSlicePrefixByte] to
// confirm that it has been encoded correctly.
func HashSliceToBytes(hashes []common.Hash) ([]byte, bool) {
	if len(hashes) == 0 {
		return nil, false
	}

	zeroStrippedBytes := zeroStrippedSlice(hashes[0][:])

	prefixStrippedBytes, hasPrefix := bytes.CutPrefix(zeroStrippedBytes, []byte{hashSlicePrefixByte})
	if !hasPrefix {
		return nil, false
	}

	copiedBytes := make([]byte, len(prefixStrippedBytes)+common.HashLength*(len(hashes)-1))
	copy(copiedBytes, prefixStrippedBytes)
	offset := len(prefixStrippedBytes)
	for _, hash := range hashes[1:] {
		copy(copiedBytes[offset:], hash[:])
		offset += common.HashLength
	}
	return copiedBytes, true
}

// BytesToHashSlice packs input into a slice of hashes.
// Packs with zero padding and a prefix of hashSlicePrefixByte to
// indicate the start of the actual bytes.
func BytesToHashSlice(input []byte) []common.Hash {
	var output []common.Hash

	// Calculate the number of bytes to add for zero padding at the beginning
	totalLen := (len(input) + 1 + 31) / 32 * 32
	paddingLen := totalLen - (len(input) + 1)

	// Create a new slice with the padded bytes at the beginning
	paddedInput := make([]byte, totalLen)
	paddedInput[paddingLen] = hashSlicePrefixByte
	copy(paddedInput[paddingLen+1:], input)

	// Loop through the input bytes
	for i := 0; i < len(paddedInput); i += common.HashLength {
		hash := common.Hash{}
		copy(hash[:], paddedInput[i:i+common.HashLength])

		// Add the padded block to the output
		output = append(output, hash)
	}

	return output
}
