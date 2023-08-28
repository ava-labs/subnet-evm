// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"context"
	"crypto/ecdsa"
	"errors"
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
	"github.com/ava-labs/subnet-evm/tests/utils/runner"
	predicateutils "github.com/ava-labs/subnet-evm/utils/predicate"
	warpclient "github.com/ava-labs/subnet-evm/warp"
	"github.com/ava-labs/subnet-evm/x/warp"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/sync/errgroup"
)

func MkSendWarpTxGenerator(chainID *big.Int, dstChainID ids.ID, gasFeeCap, gasTipCap *big.Int) txs.CreateTx {
	txGenerator := func(key *ecdsa.PrivateKey, nonce uint64) (*types.Transaction, error) {
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
		return types.SignNewTx(key, signer, &types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			To:        &warp.Module.Address,
			Gas:       200_000, // sufficient gas
			GasFeeCap: gasFeeCap,
			GasTipCap: gasTipCap,
			Value:     common.Big0,
			Data:      packedInput,
		})
	}
	return txGenerator
}

// getTestWarpPayload returns dstChain+addr+nonce (as an arbitrary choice).
// We use this in tests to verify the warp message was sent correctly.
func getTestWarpPayload(dstChainID ids.ID, addr common.Address, nonce uint64) []byte {
	length := len(ids.Empty) + common.AddressLength + wrappers.LongLen
	p := wrappers.Packer{Bytes: make([]byte, length)}
	p.PackFixedBytes(dstChainID[:])
	p.PackFixedBytes(addr.Bytes())
	p.PackLong(nonce)
	return p.Bytes
}

type warpRelayClient struct {
	client     ethclient.Client
	warpClient warpclient.WarpClient
	aggregator chan<- warpSignature
	nodeID     ids.NodeID
}

func (wr *warpRelayClient) doLoop(ctx context.Context) error {
	log.Info("starting warp relay client", "nodeID", wr.nodeID)

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
			unsignedMsg, err := avalancheWarp.ParseUnsignedMessage(txLog.Data)
			if err != nil {
				return err
			}
			unsignedWarpMessageID := unsignedMsg.ID()

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
				signer:    wr.nodeID,
				message:   unsignedMsg,
			}
		}
	}
}

type warpSignature struct {
	message   *avalancheWarp.UnsignedMessage
	signature *bls.Signature
	signer    ids.NodeID
}

type warpMessage struct {
	weight     uint64
	signers    set.Bits
	signatures []*bls.Signature
	sent       bool
}

type warpRelay struct {
	// TODO: should be an LRU to avoid getting larger forever
	messages       map[ids.ID]*warpMessage     // map of messages to signed weight
	validatorInfo  validatorInfo               // validator info needed to aggregate signatures
	threshold      uint64                      // threshold for quorum
	signatures     <-chan warpSignature        // channel of signatures
	signedMessages chan *avalancheWarp.Message // channel of signed messages
	eg             *errgroup.Group             // tracks warp relay clients
	cancelFn       context.CancelFunc          // invoking this stops the warp relay clients
	done           <-chan struct{}             // channel to signal shutdown
}

func NewWarpRelay(
	ctx context.Context,
	subnetA *runner.Subnet,
	thresholdNumerator int,
	done <-chan struct{},
) (*warpRelay, error) {
	// We need the validator set of subnet A to determine the index of
	// each validator in the bit set.
	validatorInfo, totalWeight, err := getValidatorInfo(ctx, subnetA.ValidatorURIs[0], subnetA.SubnetID)
	if err != nil {
		return nil, err
	}

	var eg errgroup.Group
	egCtx, cancel := context.WithCancel(ctx)
	signatures := make(chan warpSignature) // channel for incoming signatures
	// We will need to aggregate signatures for messages that are sent on
	// subnet A. So we will subscribe to the subnet A's accepted logs.
	endpoints := toWebsocketURIs(subnetA)
	for i, endpoint := range endpoints {
		// Skip the node if its BLS key is not in the validator info map
		// this means the node shares a BLS key with another node which
		// is in the validator set instead.
		if _, ok := validatorInfo[subnetA.NodeIDs[i]]; !ok {
			continue
		}

		client, err := ethclient.Dial(endpoint)
		if err != nil {
			cancel() // shutdown any warp relay clients that have already been started
			return nil, fmt.Errorf("failed to dial client at %s: %w", endpoint, err)
		}
		log.Info("Connected to client", "client", endpoint, "idx", i)

		warpClient, err := warpclient.NewWarpClient(
			subnetA.ValidatorURIs[i], subnetA.BlockchainID.String())
		if err != nil {
			cancel() // shutdown any warp relay clients that have already been started
			return nil, err
		}
		relayClient := &warpRelayClient{
			client:     client,
			warpClient: warpClient,
			aggregator: signatures,
			nodeID:     subnetA.NodeIDs[i],
		}
		eg.Go(func() error {
			return relayClient.doLoop(egCtx)
		})
	}

	return &warpRelay{
		messages:       make(map[ids.ID]*warpMessage),
		validatorInfo:  validatorInfo,
		threshold:      totalWeight * uint64(thresholdNumerator) / 100,
		signatures:     signatures,
		signedMessages: make(chan *avalancheWarp.Message),
		eg:             &eg,
		cancelFn:       cancel,
		done:           done,
	}, nil
}

func (wr *warpRelay) Run(ctx context.Context) error {
	defer func() {
		wr.cancelFn() // shutdown the warp relay clients
		if err := wr.eg.Wait(); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Error("warp relay client failed", "err", err)
			}
		}
		close(wr.signedMessages)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-wr.done:
			return nil

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
			vdr, ok := wr.validatorInfo[signature.signer]
			if !ok {
				return fmt.Errorf("received signature from unknown validator %s", signature.signer)
			}
			message.signers.Add(vdr.index)
			message.signatures = append(message.signatures, signature.signature)
			message.weight += vdr.weight
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
			log.Info(
				"Signatures aggregated",
				"messageID", messageID,
				"signers", message.signers.Len(),
			)
			wr.signedMessages <- msg
			message.sent = true
		}
	}
}

type warpRelayTxSequence struct {
	messages <-chan *avalancheWarp.Message
	chainID  *big.Int
	key      *ecdsa.PrivateKey
	nonce    uint64

	txs chan *types.Transaction
}

func NewWarpRelayTxSequence(
	ctx context.Context,
	messages <-chan *avalancheWarp.Message,
	chainID *big.Int,
	key *ecdsa.PrivateKey,
	startingNonce uint64,
) txs.TxSequence[*types.Transaction] {
	wr := &warpRelayTxSequence{
		messages: messages,
		chainID:  chainID,
		key:      key,
		nonce:    startingNonce,
		txs:      make(chan *types.Transaction, 1),
	}
	go func() {
		err := wr.doLoop(ctx)
		if err != nil {
			log.Error("warp relay tx sequence failed", "err", err)
		}
	}()
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
			wr.nonce++
			wr.txs <- signedTx
		}
	}
}

func (wr *warpRelayTxSequence) Chan() <-chan *types.Transaction {
	return wr.txs
}
