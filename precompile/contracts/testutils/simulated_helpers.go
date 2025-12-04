// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package testutils

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/types"
	"github.com/ava-labs/libevm/crypto"
	"github.com/stretchr/testify/require"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/params/extras"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"

	sim "github.com/ava-labs/subnet-evm/ethclient/simulated"
)

// NewAuth creates a new transactor with the given private key and chain ID.
func NewAuth(t *testing.T, key *ecdsa.PrivateKey, chainID *big.Int) *bind.TransactOpts {
	t.Helper()
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	require.NoError(t, err)
	return auth
}

// NewBackendWithPrecompile creates a simulated backend with the given precompile enabled
// at genesis and funds the specified addresses with 1 ETH each.
func NewBackendWithPrecompile(t *testing.T, precompileCfg precompileconfig.Config, fundedAddrs ...common.Address) *sim.Backend {
	t.Helper()
	chainCfg := params.Copy(params.TestChainConfig)
	params.GetExtra(&chainCfg).GenesisPrecompiles = extras.Precompiles{
		precompileCfg.Key(): precompileCfg,
	}

	genesisAlloc := make(types.GenesisAlloc)
	for _, addr := range fundedAddrs {
		genesisAlloc[addr] = types.Account{Balance: big.NewInt(1000000000000000000)}
	}

	return sim.NewBackend(genesisAlloc, sim.WithChainConfig(&chainCfg))
}

// WaitReceipt commits the simulated backend and waits for the transaction receipt.
func WaitReceipt(t *testing.T, b *sim.Backend, tx *types.Transaction) *types.Receipt {
	t.Helper()
	b.Commit(true)
	receipt, err := b.Client().TransactionReceipt(t.Context(), tx.Hash())
	require.NoError(t, err, "failed to get transaction receipt")
	return receipt
}

// WaitReceiptSuccessful commits the backend, waits for the receipt, and asserts success.
func WaitReceiptSuccessful(t *testing.T, b *sim.Backend, tx *types.Transaction) *types.Receipt {
	t.Helper()
	receipt := WaitReceipt(t, b, tx)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status, "transaction should succeed")
	return receipt
}

// SendSimpleTx sends a simple ETH transfer transaction
// See ethclient/simulated/backend_test.go newTx() for the source of this code
// TODO(jonathanoppenheimer): after libevmifiying the geth code, investigate whether we can use the same code for both
func SendSimpleTx(t *testing.T, b *sim.Backend, key *ecdsa.PrivateKey) *types.Transaction {
	t.Helper()
	client := b.Client()
	addr := crypto.PubkeyToAddress(key.PublicKey)

	chainID, err := client.ChainID(t.Context())
	require.NoError(t, err)

	nonce, err := client.NonceAt(t.Context(), addr, nil)
	require.NoError(t, err)

	head, err := client.HeaderByNumber(t.Context(), nil)
	require.NoError(t, err)

	gasPrice := new(big.Int).Add(head.BaseFee, big.NewInt(params.GWei))

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: big.NewInt(params.GWei),
		GasFeeCap: gasPrice,
		Gas:       21000,
		To:        &addr,
		Value:     big.NewInt(0),
	})

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), key)
	require.NoError(t, err)

	err = client.SendTransaction(t.Context(), signedTx)
	require.NoError(t, err)

	return signedTx
}
