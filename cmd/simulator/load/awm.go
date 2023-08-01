// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/avalanchego/utils/wrappers"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/cmd/simulator/txs"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/interfaces"
	"github.com/ava-labs/subnet-evm/params"
	predicateutils "github.com/ava-labs/subnet-evm/utils/predicate"
	warpclient "github.com/ava-labs/subnet-evm/warp"
	"github.com/ava-labs/subnet-evm/warp/payload"
	"github.com/ava-labs/subnet-evm/x/warp"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

func MkSendWarpTxGenerator(chainID *big.Int, dstChainID ids.ID, gasFeeCap, gasTipCap *big.Int) txs.CreateTx[*AwmTx] {
	txGenerator := func(key *ecdsa.PrivateKey, nonce uint64) (*AwmTx, error) {
		addr := ethcrypto.PubkeyToAddress(key.PublicKey)
		input := warp.SendWarpMessageInput{
			DestinationChainID: common.Hash(dstChainID),
			DestinationAddress: addr,
			Payload:            getTestWarpPayload(dstChainID, addr, nonce),
		}
		packedInput, err := warp.PackSendWarpMessage(input)
		if err != nil {
			return nil, err
		}
		signer := types.LatestSignerForChainID(chainID)
		tx, err := types.SignNewTx(key, signer, &types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			To:        &warp.Module.Address,
			Gas:       200_000,
			GasFeeCap: gasFeeCap,
			GasTipCap: gasTipCap,
			Value:     common.Big0,
			Data:      packedInput,
		})
		if err != nil {
			return nil, err
		}

		// Compute a unique ID to track this AWM message
		// TODO: alternatively, we can pass some additional parameters
		// here and compute the AWM message ID, or we can use the
		// tx id on this side and track it via txHash which is available
		// on accepted logs.
		awmTx := &AwmTx{
			Tx:    tx,
			AwmID: ethcrypto.Keccak256Hash(input.Payload),
		}
		return awmTx, nil
	}
	return txGenerator
}

// getTestWarpPayload returns dstChain+addr+nonce (as an arbitrary choice).
// We use this in tests to verify the warp message was sent correctly.
func getTestWarpPayload(dstChainID ids.ID, addr common.Address, nonce uint64) []byte {
	length := len(ids.Empty) + common.AddressLength + wrappers.LongLen
	p := wrappers.Packer{Bytes: make([]byte, 0, length)}
	p.PackFixedBytes(dstChainID[:])
	p.PackFixedBytes(addr.Bytes())
	p.PackLong(nonce)
	return p.Bytes
}

type warpRelayClient struct {
	client     ethclient.Client
	warpClient warpclient.WarpClient
	aggregator chan<- warpSignature
	index      int
}

func NewWarpRelayClient(
	ctx context.Context,
	client ethclient.Client,
	warpClient warpclient.WarpClient,
	aggregator chan<- warpSignature,
	index int,
) *warpRelayClient {
	wr := &warpRelayClient{
		client:     client,
		warpClient: warpClient,
		aggregator: aggregator,
		index:      index,
	}
	go wr.doLoop(ctx)
	return wr
}

func (wr *warpRelayClient) doLoop(ctx context.Context) error {
	log.Info("starting warp relay client", "index", wr.index)

	logsCh := make(chan types.Log, 1)
	sub, err := wr.client.SubscribeFilterLogs(
		ctx,
		interfaces.FilterQuery{
			Addresses: []common.Address{warp.ContractAddress},
		},
		logsCh,
	)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case txLog, ok := <-logsCh:
			if !ok {
				return nil
			}
			log.Info("Parsing logData as unsigned warp message", "logData", common.Bytes2Hex(txLog.Data))
			unsignedMsg, err := avalancheWarp.ParseUnsignedMessage(txLog.Data)
			if err != nil {
				return err
			}
			unsignedWarpMessageID := unsignedMsg.ID()
			log.Info("Parsed unsignedWarpMsg", "unsignedWarpMessageID", unsignedWarpMessageID)

			signature, err := wr.warpClient.GetSignature(ctx, unsignedWarpMessageID)
			if err != nil {
				return err
			}

			blsSignature, err := bls.SignatureFromBytes(signature)
			if err != nil {
				return fmt.Errorf("failed to parse signature: %w", err)
			}

			wr.aggregator <- warpSignature{
				signature: blsSignature,
				signer:    wr.index,
				message:   unsignedMsg,
			}
		}
	}
}

type warpSignature struct {
	message   *avalancheWarp.UnsignedMessage
	signature *bls.Signature
	signer    int
}

type warpMessage struct {
	weight     uint64
	signers    set.Bits
	signatures []*bls.Signature
	sent       bool
}

type warpRelay struct {
	// TODO: should be an LRU to avoid getting larger forever
	messages         map[ids.ID]*warpMessage     // map of messages to signed weight
	threshold        uint64                      // threshold for quorum
	signatures       <-chan warpSignature        // channel of signatures
	signedMessages   chan *avalancheWarp.Message // channel of signed messages
	expectedMessages int                         // close signedMessages when this many messages are received
}

func NewWarpRelay(
	ctx context.Context,
	threshold uint64,
	signatures <-chan warpSignature,
	expectedMessages int,
) *warpRelay {
	wr := &warpRelay{
		messages:         make(map[ids.ID]*warpMessage),
		threshold:        threshold,
		signatures:       signatures,
		signedMessages:   make(chan *avalancheWarp.Message),
		expectedMessages: expectedMessages,
	}
	go wr.doLoop(ctx)
	return wr
}

func (wr *warpRelay) doLoop(ctx context.Context) error {
	defer close(wr.signedMessages)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case signature, ok := <-wr.signatures:
			if !ok {
				return nil
			}
			messageID := signature.message.ID()

			// If this not a known message, initialize it
			if _, ok := wr.messages[messageID]; !ok {
				wr.messages[messageID] = &warpMessage{
					signers: set.NewBits(),
				}
			}
			message := wr.messages[messageID]

			// If the message is already sent, ignore this signature
			if message.sent {
				continue
			}
			message.signers.Add(signature.signer)
			message.signatures = append(message.signatures, signature.signature)
			message.weight += 1 // TODO: use actual weights
			if message.weight < wr.threshold {
				continue
			}

			// Send the message if we have enough signatures
			aggregateSignature, err := bls.AggregateSignatures(message.signatures)
			if err != nil {
				return fmt.Errorf("failed to aggregate BLS signatures: %w", err)
			}
			warpSignature := &avalancheWarp.BitSetSignature{
				Signers: message.signers.Bytes(),
			}
			copy(warpSignature.Signature[:], bls.SignatureToBytes(aggregateSignature))
			msg, err := avalancheWarp.NewMessage(signature.message, warpSignature)
			if err != nil {
				return fmt.Errorf("failed to construct warp message: %w", err)
			}

			// Send the message on the result channel and mark it as sent
			wr.signedMessages <- msg
			message.sent = true
			wr.expectedMessages--
			if wr.expectedMessages == 0 {
				return nil
			}
		}
	}
}

type warpRelayTxSequence struct {
	messages <-chan *avalancheWarp.Message
	chainID  *big.Int
	key      *ecdsa.PrivateKey
	nonce    uint64

	txs chan *AwmTx
}

func NewWarpRelayTxSequence(
	ctx context.Context,
	messages <-chan *avalancheWarp.Message,
	chainID *big.Int,
	key *ecdsa.PrivateKey,
	startingNonce uint64,
) txs.TxSequence[*AwmTx] {
	wr := &warpRelayTxSequence{
		messages: messages,
		chainID:  chainID,
		key:      key,
		nonce:    startingNonce,
		txs:      make(chan *AwmTx, 1),
	}
	go wr.doLoop(ctx)
	return wr
}

func (wr *warpRelayTxSequence) doLoop(ctx context.Context) error {
	defer close(wr.txs)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-wr.messages:
			if !ok {
				return nil
			}
			packedInput, err := warp.PackGetVerifiedWarpMessage()
			if err != nil {
				return err
			}
			tx := predicateutils.NewPredicateTx(
				wr.chainID,
				wr.nonce,
				&warp.Module.Address,
				5_000_000,
				big.NewInt(225*params.GWei),
				big.NewInt(params.GWei),
				common.Big0,
				packedInput,
				types.AccessList{},
				warp.ContractAddress,
				msg.Bytes(),
			)
			signer := types.LatestSignerForChainID(wr.chainID)
			signedTx, err := types.SignTx(tx, signer, wr.key)
			if err != nil {
				return err
			}
			// Recompute the unique ID used to track this AWM message aka
			// the hash of the payload before wrapping in AddressedPayload.
			payload, err := payload.ParseAddressedPayload(msg.Payload)
			if err != nil {
				return err
			}
			awmTx := &AwmTx{
				Tx:    signedTx,
				AwmID: ethcrypto.Keccak256Hash(payload.Payload),
			}
			wr.nonce++
			wr.txs <- awmTx
		}
	}
}

func (wr *warpRelayTxSequence) Chan() <-chan *AwmTx {
	return wr.txs
}
