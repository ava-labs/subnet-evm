// Copyright 2018 The go-ethereum Authors
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

package apitypes

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/holiman/uint256"
)

type ValidationInfo struct {
	Typ     string `json:"type"`
	Message string `json:"message"`
}
type ValidationMessages struct {
	Messages []ValidationInfo
}

const (
	WARN = "WARNING"
	CRIT = "CRITICAL"
	INFO = "Info"
)

func (vs *ValidationMessages) Crit(msg string) {
	vs.Messages = append(vs.Messages, ValidationInfo{CRIT, msg})
}
func (vs *ValidationMessages) Warn(msg string) {
	vs.Messages = append(vs.Messages, ValidationInfo{WARN, msg})
}
func (vs *ValidationMessages) Info(msg string) {
	vs.Messages = append(vs.Messages, ValidationInfo{INFO, msg})
}

// GetWarnings returns an error with all messages of type WARN of above, or nil if no warnings were present
func (vs *ValidationMessages) GetWarnings() error {
	var messages []string
	for _, msg := range vs.Messages {
		if msg.Typ == WARN || msg.Typ == CRIT {
			messages = append(messages, msg.Message)
		}
	}
	if len(messages) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(messages, ","))
	}
	return nil
}

// SendTxArgs represents the arguments to submit a transaction
// This struct is identical to ethapi.TransactionArgs, except for the usage of
// common.MixedcaseAddress in From and To
type SendTxArgs struct {
	From                 common.MixedcaseAddress  `json:"from"`
	To                   *common.MixedcaseAddress `json:"to"`
	Gas                  hexutil.Uint64           `json:"gas"`
	GasPrice             *hexutil.Big             `json:"gasPrice"`
	MaxFeePerGas         *hexutil.Big             `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big             `json:"maxPriorityFeePerGas"`
	Value                hexutil.Big              `json:"value"`
	Nonce                hexutil.Uint64           `json:"nonce"`

	// We accept "data" and "input" for backwards-compatibility reasons.
	// "input" is the newer name and should be preferred by clients.
	// Issue detail: https://github.com/ethereum/go-ethereum/issues/15628
	Data  *hexutil.Bytes `json:"data,omitempty"`
	Input *hexutil.Bytes `json:"input,omitempty"`

	// For non-legacy transactions
	AccessList *types.AccessList `json:"accessList,omitempty"`
	ChainID    *hexutil.Big      `json:"chainId,omitempty"`

	// For BlobTxType
	BlobFeeCap *hexutil.Big  `json:"maxFeePerBlobGas,omitempty"`
	BlobHashes []common.Hash `json:"blobVersionedHashes,omitempty"`

	// For BlobTxType transactions with blob sidecar
	Blobs       []kzg4844.Blob       `json:"blobs,omitempty"`
	Commitments []kzg4844.Commitment `json:"commitments,omitempty"`
	Proofs      []kzg4844.Proof      `json:"proofs,omitempty"`
}

func (args SendTxArgs) String() string {
	s, err := json.Marshal(args)
	if err == nil {
		return string(s)
	}
	return err.Error()
}

// data retrieves the transaction calldata. Input field is preferred.
func (args *SendTxArgs) data() []byte {
	if args.Input != nil {
		return *args.Input
	}
	if args.Data != nil {
		return *args.Data
	}
	return nil
}

// ToTransaction converts the arguments to a transaction.
func (args *SendTxArgs) ToTransaction() (*types.Transaction, error) {
	// Add the To-field, if specified
	var to *common.Address
	if args.To != nil {
		dstAddr := args.To.Address()
		to = &dstAddr
	}
	if err := args.validateTxSidecar(); err != nil {
		return nil, err
	}
	var data types.TxData
	switch {
	case args.BlobHashes != nil:
		al := types.AccessList{}
		if args.AccessList != nil {
			al = *args.AccessList
		}
		data = &types.BlobTx{
			To:         *to,
			ChainID:    uint256.MustFromBig((*big.Int)(args.ChainID)),
			Nonce:      uint64(args.Nonce),
			Gas:        uint64(args.Gas),
			GasFeeCap:  uint256.MustFromBig((*big.Int)(args.MaxFeePerGas)),
			GasTipCap:  uint256.MustFromBig((*big.Int)(args.MaxPriorityFeePerGas)),
			Value:      uint256.MustFromBig((*big.Int)(&args.Value)),
			Data:       args.data(),
			AccessList: al,
			BlobHashes: args.BlobHashes,
			BlobFeeCap: uint256.MustFromBig((*big.Int)(args.BlobFeeCap)),
		}
		if args.Blobs != nil {
			data.(*types.BlobTx).Sidecar = &types.BlobTxSidecar{
				Blobs:       args.Blobs,
				Commitments: args.Commitments,
				Proofs:      args.Proofs,
			}
		}

	case args.MaxFeePerGas != nil:
		al := types.AccessList{}
		if args.AccessList != nil {
			al = *args.AccessList
		}
		data = &types.DynamicFeeTx{
			To:         to,
			ChainID:    (*big.Int)(args.ChainID),
			Nonce:      uint64(args.Nonce),
			Gas:        uint64(args.Gas),
			GasFeeCap:  (*big.Int)(args.MaxFeePerGas),
			GasTipCap:  (*big.Int)(args.MaxPriorityFeePerGas),
			Value:      (*big.Int)(&args.Value),
			Data:       args.data(),
			AccessList: al,
		}
	case args.AccessList != nil:
		data = &types.AccessListTx{
			To:         to,
			ChainID:    (*big.Int)(args.ChainID),
			Nonce:      uint64(args.Nonce),
			Gas:        uint64(args.Gas),
			GasPrice:   (*big.Int)(args.GasPrice),
			Value:      (*big.Int)(&args.Value),
			Data:       args.data(),
			AccessList: *args.AccessList,
		}
	default:
		data = &types.LegacyTx{
			To:       to,
			Nonce:    uint64(args.Nonce),
			Gas:      uint64(args.Gas),
			GasPrice: (*big.Int)(args.GasPrice),
			Value:    (*big.Int)(&args.Value),
			Data:     args.data(),
		}
	}

	return types.NewTx(data), nil
}

// validateTxSidecar validates blob data, if present
func (args *SendTxArgs) validateTxSidecar() error {
	// No blobs, we're done.
	if args.Blobs == nil {
		return nil
	}

	n := len(args.Blobs)
	// Assume user provides either only blobs (w/o hashes), or
	// blobs together with commitments and proofs.
	if args.Commitments == nil && args.Proofs != nil {
		return errors.New(`blob proofs provided while commitments were not`)
	} else if args.Commitments != nil && args.Proofs == nil {
		return errors.New(`blob commitments provided while proofs were not`)
	}

	// len(blobs) == len(commitments) == len(proofs) == len(hashes)
	if args.Commitments != nil && len(args.Commitments) != n {
		return fmt.Errorf("number of blobs and commitments mismatch (have=%d, want=%d)", len(args.Commitments), n)
	}
	if args.Proofs != nil && len(args.Proofs) != n {
		return fmt.Errorf("number of blobs and proofs mismatch (have=%d, want=%d)", len(args.Proofs), n)
	}
	if args.BlobHashes != nil && len(args.BlobHashes) != n {
		return fmt.Errorf("number of blobs and hashes mismatch (have=%d, want=%d)", len(args.BlobHashes), n)
	}

	if args.Commitments == nil {
		// Generate commitment and proof.
		commitments := make([]kzg4844.Commitment, n)
		proofs := make([]kzg4844.Proof, n)
		for i, b := range args.Blobs {
			c, err := kzg4844.BlobToCommitment(&b)
			if err != nil {
				return fmt.Errorf("blobs[%d]: error computing commitment: %v", i, err)
			}
			commitments[i] = c
			p, err := kzg4844.ComputeBlobProof(&b, c)
			if err != nil {
				return fmt.Errorf("blobs[%d]: error computing proof: %v", i, err)
			}
			proofs[i] = p
		}
		args.Commitments = commitments
		args.Proofs = proofs
	} else {
		for i, b := range args.Blobs {
			if err := kzg4844.VerifyBlobProof(&b, args.Commitments[i], args.Proofs[i]); err != nil {
				return fmt.Errorf("failed to verify blob proof: %v", err)
			}
		}
	}

	hashes := make([]common.Hash, n)
	hasher := sha256.New()
	for i, c := range args.Commitments {
		hashes[i] = kzg4844.CalcBlobHashV1(hasher, &c)
	}
	if args.BlobHashes != nil {
		for i, h := range hashes {
			if h != args.BlobHashes[i] {
				return fmt.Errorf("blob hash verification failed (have=%s, want=%s)", args.BlobHashes[i], h)
			}
		}
	} else {
		args.BlobHashes = hashes
	}
	return nil
}
