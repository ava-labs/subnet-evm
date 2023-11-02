// (c) 2020-2021, Ava Labs, Inc.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package evm

import (
	"context"
	_ "embed"
	"math/big"
	"testing"
	"time"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/accounts/abi/bind/backends"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed TrieStressTest.bin
	stressBinStr string
	//go:embed TrieStressTest.abi
	stressABIStr string
)
var testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")

func BenchmarkTrie(t *testing.B) {
	require := require.New(t)
	vm, address := createVM(t, "TestTrieStreeMessage")

	head, _ := vm.HeaderByNumber(context.Background(), nil) // Should be child's, good enough
	stressABI := contract.ParseABI(stressABIStr)
	gasPrice := new(big.Int).Add(head.BaseFee, big.NewInt(1))
	txPayload, err := stressABI.Pack(
		"writeValues",
		big.NewInt(5000),
	)
	require.NoError(err)
	tx := types.NewTransaction(1, address, big.NewInt(0), 3000000, gasPrice, txPayload)
	signer := types.NewLondonSigner(big.NewInt(1337))
	tx, err = types.SignTx(tx, signer, testKey)
	require.NoError(err)

	var (
		mined = make(chan struct{})
		ctx   = context.Background()
	)
	go func() {
		address, err = bind.WaitDeployed(ctx, vm, tx)
		close(mined)
	}()

	err = vm.SendTransaction(ctx, tx)
	require.NoError(err)
	vm.Commit(true)

	select {
	case <-mined:
		require.NoError(err)
	case <-time.After(2 * time.Second):
		t.Errorf("test timeout waiting for function call")
	}

	vm.GetInternalDB() // Maybe we would needs this?
}

func createVM(t *testing.B, name string) (*backends.SimulatedBackend, common.Address) {
	require := require.New(t)
	backend := backends.NewSimulatedBackend(
		core.GenesisAlloc{
			crypto.PubkeyToAddress(testKey.PublicKey): {Balance: new(big.Int).Mul(big.NewInt(10000000000000000), big.NewInt(1000))},
		},
		10000000,
	)

	head, _ := backend.HeaderByNumber(context.Background(), nil) // Should be child's, good enough
	gasPrice := new(big.Int).Add(head.BaseFee, big.NewInt(1))

	tx := types.NewContractCreation(0, big.NewInt(0), 3000000, gasPrice, common.FromHex(stressBinStr))
	signer := types.NewLondonSigner(big.NewInt(1337))
	tx, _ = types.SignTx(tx, signer, testKey)

	// Wait for it to get mined in the background.
	var (
		err     error
		address common.Address
		mined   = make(chan struct{})
		ctx     = context.Background()
	)
	go func() {
		address, err = bind.WaitDeployed(ctx, backend, tx)
		close(mined)
	}()

	// Send and mine the transaction.
	if err := backend.SendTransaction(ctx, tx); err != nil {
		t.Errorf("Failed to send transaction: %s", err)
	}
	backend.Commit(true)

	select {
	case <-mined:
		require.NoError(err)
	case <-time.After(2 * time.Second):
		t.Errorf("test %q: timeout", name)
	}
	return backend, address
}
