// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package abi

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Note: This file contains tests in addition to those found in go-ethereum.

const TEST_ABI = `[ { "type" : "function", "name" : "receive", "inputs" : [ { "name" : "memo", "type" : "bytes" }], "outputs":[{"internalType":"bool","name":"isAllowed","type":"bool"}] } ]`

func TestUnpackInputIntoInterface(t *testing.T) {
	abi, err := JSON(strings.NewReader(TEST_ABI))
	require.NoError(t, err)

	input := []byte("hello")
	data, err := abi.Pack("receive", input)
	require.NoError(t, err)

	var v []byte
	err = abi.UnpackInputIntoInterface(&v, "receive", data[4:]) // skips 4 byte selector
	require.NoError(t, err)

	require.True(t, bytes.Equal(v, input))
}

func TestPackOutput(t *testing.T) {
	abi, err := JSON(strings.NewReader(TEST_ABI))
	require.NoError(t, err)

	bytes, err := abi.PackOutput("receive", true)
	require.NoError(t, err)

	vals, err := abi.Methods["receive"].Outputs.Unpack(bytes)
	require.NoError(t, err)

	require.Len(t, vals, 1)
	require.True(t, vals[0].(bool))
}
