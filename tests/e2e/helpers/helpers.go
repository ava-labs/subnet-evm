// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package helpers

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"math/rand"
	"time"

	anrapi "github.com/ava-labs/avalanche-network-runner/api"
	"github.com/ava-labs/avalanchego/api"
	corethTypes "github.com/ava-labs/coreth/core/types"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	testPassword              = "ajwdygakjydwadawdg12121112121324123123123123" // #nosec G101 not a real credential
	seed                      = 0
	usernameSize              = 16
	letters                   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	txAcceptedTimeout         = 20 * time.Second
	defaultGasPriceMultiplier = 1.5
)

var (
	r = rand.New(rand.NewSource(seed)) // #nosec G404
)

// MultiplyMaxGasPrice returns a multiplied gas price as returned from
// SuggestGasPrice to account for dynamic fees and EIP-1559
func MultiplyMaxGasPrice(gasPrice *big.Int) *big.Int {
	fMultiPlier := big.NewFloat(defaultGasPriceMultiplier)
	fGasPrice := big.NewFloat(float64(gasPrice.Int64()))
	fMaxGasPrice := fGasPrice.Mul(fGasPrice, fMultiPlier)
	iMaxGasPrice := new(big.Int)
	fMaxGasPrice.Int(iMaxGasPrice)
	return iMaxGasPrice
}

// NewUserPass returns a new struct with a random username/password
func NewUserPass() api.UserPass {
	username := make([]byte, usernameSize)
	for i := 0; i < usernameSize; i++ {
		idx := r.Intn(len(letters)) // #nosec G404
		username[i] = letters[idx]
	}
	return api.UserPass{Username: string(username), Password: testPassword}
}

// AwaitedSendTransaction starts and waits tx in cchain
func AwaitedSendTransaction(
	ctx context.Context,
	client anrapi.EthClient,
	senderKey *ecdsa.PrivateKey,
	senderNonce uint64,
	address common.Address,
	amount *big.Int,
	data []byte,
	ethCChainID *big.Int,
	gasLimit uint64,
	gasPrice *big.Int,
) error {
	tx, err := SendTransaction(ctx, client, senderKey, senderNonce, address, amount, data, ethCChainID, gasLimit, gasPrice)
	if err != nil {
		return err
	}
	return AwaitTransaction(ctx, client, tx)
}

// AwaitTransaction waits tx completion receipt in cchain, postprocess status to err
func AwaitTransaction(ctx context.Context, client anrapi.EthClient, tx *corethTypes.Transaction) error {
	var r *corethTypes.Receipt
	var err error
	for startTime := time.Now(); time.Since(startTime) < txAcceptedTimeout; {
		r, err = client.TransactionReceipt(ctx, tx.Hash())
		if r != nil {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
	if r != nil && r.Status != types.ReceiptStatusSuccessful && err == nil {
		err = errors.New("receipt status problem")
	}
	return err
}

// SendTransaction creates and sends to cchain a [senderKey]-signed evm transaction
// that transfers [amount] to [addresss], including extra bytes [data] on it
// Both [senderNonce] and [senderKey] should have values consistent with the underlying
// senderAddr
func SendTransaction(
	ctx context.Context,
	client anrapi.EthClient,
	senderKey *ecdsa.PrivateKey,
	senderNonce uint64,
	address common.Address,
	amount *big.Int,
	data []byte,
	ethCChainID *big.Int,
	gasLimit uint64,
	gasPrice *big.Int,
) (*corethTypes.Transaction, error) {
	tx := corethTypes.NewTx(&corethTypes.LegacyTx{
		Nonce:    senderNonce,
		To:       &address,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})
	signedTx, err := corethTypes.SignTx(tx, corethTypes.NewEIP155Signer(ethCChainID), senderKey)
	if err != nil {
		return signedTx, err
	}
	return signedTx, client.SendTransaction(ctx, signedTx)
}
