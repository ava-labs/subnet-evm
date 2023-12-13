// Code generated
// This file is a generated precompile contract config with stubbed abstract functions.
// The file is generated by a template. Please inspect every code and comment in this file before use.

package nativeminter

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
)

const (
	NativeCoinMintedEventGasCost = contract.LogGas + contract.LogTopicGas*3
)

// PackNativeCoinMintedEvent packs the event into the appropriate arguments for NativeCoinMinted.
// It returns topic hashes and the encoded non-indexed data.
func PackNativeCoinMintedEvent(sender common.Address, recipient common.Address, amount *big.Int) ([]common.Hash, []byte, error) {
	return NativeMinterABI.PackEvent("NativeCoinMinted", sender, recipient, amount)
}
