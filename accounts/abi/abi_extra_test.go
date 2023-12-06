// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package abi

import (
	"bytes"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// Note: This file contains tests in addition to those found in go-ethereum.

const TEST_ABI = `[{"type":"function","name":"receive","inputs":[{"name":"sender","type":"address"},{"name":"amount","type":"uint256"},{"name":"memo","type":"bytes"}],"outputs":[{"internalType":"bool","name":"isAllowed","type":"bool"}]}]`

func TestUnpackInputIntoInterface(t *testing.T) {
	abi, err := JSON(strings.NewReader(TEST_ABI))
	require.NoError(t, err)

	type inputType struct {
		Sender common.Address
		Amount *big.Int
		Memo   []byte
	}
	input := inputType{
		Sender: common.HexToAddress("0x02"),
		Amount: big.NewInt(100),
		Memo:   []byte("hello"),
	}
	data, err := abi.Pack("receive", input.Sender, input.Amount, input.Memo)
	require.NoError(t, err)

	// Unpack into interface
	var v inputType
	err = abi.UnpackInputIntoInterface(&v, "receive", data[4:], true) // skips 4 byte selector
	require.NoError(t, err)

	// Verify unpacked values match input
	require.Equal(t, v.Amount, input.Amount)
	require.EqualValues(t, v.Amount, input.Amount)
	require.True(t, bytes.Equal(v.Memo, input.Memo))

	// add 32-bytes extra padding to data
	// this should work because it is divisible by 32
	data = append(data, make([]byte, 32)...)
	err = abi.UnpackInputIntoInterface(&v, "receive", data[4:], true)
	require.NoError(t, err)
	// Verify unpacked values match input
	require.Equal(t, v.Amount, input.Amount)
	require.EqualValues(t, v.Amount, input.Amount)
	require.True(t, bytes.Equal(v.Memo, input.Memo))

	// add 33-bytes extra padding to data
	// this should fail because it is not divisible by 32
	data = append(data, make([]byte, 33)...)
	err = abi.UnpackInputIntoInterface(&v, "receive", data[4:], true)
	require.ErrorContains(t, err, "abi: improperly formatted input:")

	// this should not fail because we don't use strict mode
	err = abi.UnpackInputIntoInterface(&v, "receive", data[4:], false)
	require.NoError(t, err)
	require.Equal(t, v.Amount, input.Amount)
	require.EqualValues(t, v.Amount, input.Amount)
	require.True(t, bytes.Equal(v.Memo, input.Memo))
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
