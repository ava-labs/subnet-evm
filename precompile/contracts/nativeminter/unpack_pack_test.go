// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativeminter

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/constants"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func FuzzPackMintNativeCoinEqualTest(f *testing.F) {
	key, err := crypto.GenerateKey()
	require.NoError(f, err)
	addr := crypto.PubkeyToAddress(key.PublicKey)
	testAddrBytes := addr.Bytes()
	f.Add(testAddrBytes, common.Big0.Bytes())
	f.Add(testAddrBytes, common.Big1.Bytes())
	f.Add(testAddrBytes, abi.MaxUint256.Bytes())
	f.Add(testAddrBytes, new(big.Int).Sub(abi.MaxUint256, common.Big1).Bytes())
	f.Add(testAddrBytes, new(big.Int).Add(abi.MaxUint256, common.Big1).Bytes())
	f.Add(constants.BlackholeAddr.Bytes(), common.Big2.Bytes())
	f.Fuzz(func(t *testing.T, b []byte, bigIntBytes []byte) {
		bigIntVal := new(big.Int).SetBytes(bigIntBytes)
		doCheckOutputs := true
		// we can only check if outputs are correct if the value is less than MaxUint256
		// otherwise the value will be truncated when packed,
		// and thus unpacked output will not be equal to the value
		if bigIntVal.Cmp(abi.MaxUint256) > 0 {
			doCheckOutputs = false
		}
		testOldPackMintNativeCoinEqual(t, common.BytesToAddress(b), bigIntVal, doCheckOutputs)
	})
}

func TestUnpackMintNativeCoinInputEdgeCases(t *testing.T) {
	input, err := PackMintNativeCoin(constants.BlackholeAddr, common.Big2)
	require.NoError(t, err)
	// exclude 4 bytes for function selector
	input = input[4:]
	// add extra padded bytes
	input = append(input, make([]byte, 32)...)

	_, _, err = OldUnpackMintNativeCoinInput(input)
	require.ErrorIs(t, err, ErrInvalidLen)

	_, _, err = UnpackMintNativeCoinInput(input, false)
	require.ErrorIs(t, err, ErrInvalidLen)

	addr, value, err := UnpackMintNativeCoinInput(input, true)
	require.NoError(t, err)
	require.Equal(t, constants.BlackholeAddr, addr)
	require.Equal(t, common.Big2.Bytes(), value.Bytes())

	input = append(input, make([]byte, 1)...)
	// now it is not divisible by 32
	_, _, err = UnpackMintNativeCoinInput(input, true)
	require.Error(t, err)
}

func TestFunctionSignatures(t *testing.T) {
	// Test that the mintNativeCoin signature is correct
	abiMintNativeCoin := NativeMinterABI.Methods["mintNativeCoin"]
	require.Equal(t, mintSignature, abiMintNativeCoin.ID)
}

func testOldPackMintNativeCoinEqual(t *testing.T, addr common.Address, amount *big.Int, checkOutputs bool) {
	t.Helper()
	t.Run(fmt.Sprintf("TestUnpackAndPacks, addr: %s, amount: %s", addr.String(), amount.String()), func(t *testing.T) {
		input, err := OldPackMintNativeCoinInput(addr, amount)
		input2, err2 := PackMintNativeCoin(addr, amount)
		if err != nil {
			require.ErrorContains(t, err2, err.Error())
			return
		}
		require.NoError(t, err2)
		require.Equal(t, input, input2)

		input = input[4:]
		to, assetAmount, err := OldUnpackMintNativeCoinInput(input)
		unpackedAddr, unpackedAmount, err2 := UnpackMintNativeCoinInput(input, false)
		if err != nil {
			require.ErrorContains(t, err2, err.Error())
			return
		}
		require.NoError(t, err2)
		require.Equal(t, to, unpackedAddr)
		require.Equal(t, assetAmount.Bytes(), unpackedAmount.Bytes())
		if checkOutputs {
			require.Equal(t, addr, to)
			require.Equal(t, amount.Bytes(), assetAmount.Bytes())
		}
	})
}

func OldPackMintNativeCoinInput(address common.Address, amount *big.Int) ([]byte, error) {
	// function selector (4 bytes) + input(hash for address + hash for amount)
	res := make([]byte, contract.SelectorLen+mintInputLen)
	err := contract.PackOrderedHashesWithSelector(res, mintSignature, []common.Hash{
		address.Hash(),
		common.BigToHash(amount),
	})

	return res, err
}

func OldUnpackMintNativeCoinInput(input []byte) (common.Address, *big.Int, error) {
	mintInputAddressSlot := 0
	mintInputAmountSlot := 1
	if len(input) != mintInputLen {
		return common.Address{}, nil, fmt.Errorf("%w: %d", ErrInvalidLen, len(input))
	}
	to := common.BytesToAddress(contract.PackedHash(input, mintInputAddressSlot))
	assetAmount := new(big.Int).SetBytes(contract.PackedHash(input, mintInputAmountSlot))
	return to, assetAmount, nil
}
