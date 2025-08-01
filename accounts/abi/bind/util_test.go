// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2016 The go-ethereum Authors
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

package bind_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/types"
	"github.com/ava-labs/libevm/crypto"
	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/ethclient/simulated"
	"github.com/ava-labs/subnet-evm/params"
)

var testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")

var waitDeployedTests = map[string]struct {
	code        string
	gas         uint64
	wantAddress common.Address
	wantErr     error
}{
	"successful deploy": {
		code:        `6060604052600a8060106000396000f360606040526008565b00`,
		gas:         3000000,
		wantAddress: common.HexToAddress("0x3a220f351252089d385b29beca14e27f204c296a"),
	},
	"empty code": {
		code:        ``,
		gas:         300000,
		wantErr:     bind.ErrNoCodeAfterDeploy,
		wantAddress: common.HexToAddress("0x3a220f351252089d385b29beca14e27f204c296a"),
	},
}

func TestWaitDeployed(t *testing.T) {
	t.Parallel()
	for name, test := range waitDeployedTests {
		backend := simulated.NewBackend(
			types.GenesisAlloc{
				crypto.PubkeyToAddress(testKey.PublicKey): {Balance: new(big.Int).Mul(big.NewInt(10000000000000000), big.NewInt(1000))},
			},
		)
		defer backend.Close()

		// Create the transaction
		head, _ := backend.Client().HeaderByNumber(context.Background(), nil) // Should be child's, good enough
		gasPrice := new(big.Int).Add(head.BaseFee, big.NewInt(params.GWei))

		tx := types.NewContractCreation(0, big.NewInt(0), test.gas, gasPrice, common.FromHex(test.code))
		tx, _ = types.SignTx(tx, types.LatestSignerForChainID(big.NewInt(1337)), testKey)

		// Wait for it to get mined in the background.
		var (
			err     error
			address common.Address
			mined   = make(chan struct{})
			ctx     = context.Background()
		)
		go func() {
			address, err = bind.WaitDeployed(ctx, backend.Client(), tx)
			close(mined)
		}()

		// Send and mine the transaction.
		if err := backend.Client().SendTransaction(ctx, tx); err != nil {
			t.Fatalf("Failed to send transaction: %s", err)
		}
		backend.Commit(true)

		select {
		case <-mined:
			if err != test.wantErr {
				t.Errorf("test %q: error mismatch: want %q, got %q", name, test.wantErr, err)
			}
			if address != test.wantAddress {
				t.Errorf("test %q: unexpected contract address %s", name, address.Hex())
			}
		case <-time.After(2 * time.Second):
			t.Errorf("test %q: timeout", name)
		}
	}
}

func TestWaitDeployedCornerCases(t *testing.T) {
	backend := simulated.NewBackend(
		types.GenesisAlloc{
			crypto.PubkeyToAddress(testKey.PublicKey): {Balance: big.NewInt(1000000000000000000)},
		},
	)
	defer backend.Close()

	head, _ := backend.Client().HeaderByNumber(context.Background(), nil) // Should be child's, good enough
	gasPrice := new(big.Int).Add(head.BaseFee, big.NewInt(1))

	// Create a transaction to an account.
	code := "6060604052600a8060106000396000f360606040526008565b00"
	tx := types.NewTransaction(0, common.HexToAddress("0x01"), big.NewInt(0), 3000000, gasPrice, common.FromHex(code))
	tx, _ = types.SignTx(tx, types.LatestSignerForChainID(big.NewInt(1337)), testKey)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := backend.Client().SendTransaction(ctx, tx); err != nil {
		t.Fatalf("Failed to send transaction: %s", err)
	}
	backend.Commit(true)
	notContractCreation := errors.New("tx is not contract creation")
	if _, err := bind.WaitDeployed(ctx, backend.Client(), tx); err.Error() != notContractCreation.Error() {
		t.Errorf("error mismatch: want %q, got %q, ", notContractCreation, err)
	}

	// Create a transaction that is not mined.
	tx = types.NewContractCreation(1, big.NewInt(0), 3000000, gasPrice, common.FromHex(code))
	tx, _ = types.SignTx(tx, types.LatestSignerForChainID(big.NewInt(1337)), testKey)

	go func() {
		contextCanceled := errors.New("context canceled")
		if _, err := bind.WaitDeployed(ctx, backend.Client(), tx); err.Error() != contextCanceled.Error() {
			t.Errorf("error mismatch: want %q, got %q, ", contextCanceled, err)
		}
	}()

	if err := backend.Client().SendTransaction(ctx, tx); err != nil {
		t.Fatalf("Failed to send transaction: %s", err)
	}
	cancel()
}
