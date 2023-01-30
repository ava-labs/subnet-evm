// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/types"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/choices"
	"github.com/ava-labs/avalanchego/snow/engine/snowman/block"
)

// Block implements the snowman.Block interface
type Block struct {
	id       ids.ID
	ethBlock *types.Block
	vm       *VM
	status   choices.Status
}

// newBlock returns a new Block wrapping the ethBlock type and implementing the snowman.Block interface
func (vm *VM) newBlock(ethBlock *types.Block) *Block {
	return &Block{
		id:       ids.ID(ethBlock.Hash()),
		ethBlock: ethBlock,
		vm:       vm,
	}
}

// ID implements the snowman.Block interface
func (b *Block) ID() ids.ID { return b.id }

// Accept implements the snowman.Block interface
func (b *Block) Accept(context.Context) error {
	vm := b.vm

	// Although returning an error from Accept is considered fatal, it is good
	// practice to cleanup the batch we were modifying in the case of an error.
	defer vm.db.Abort()

	b.status = choices.Accepted
	log.Debug(fmt.Sprintf("Accepting block %s (%s) at height %d", b.ID().Hex(), b.ID(), b.Height()))
	if err := vm.blockChain.Accept(b.ethBlock); err != nil {
		return fmt.Errorf("chain could not accept %s: %w", b.ID(), err)
	}
	if err := vm.acceptedBlockDB.Put(lastAcceptedKey, b.id[:]); err != nil {
		return fmt.Errorf("failed to put %s as the last accepted block: %w", b.ID(), err)
	}

	return vm.db.Commit()
}

// Reject implements the snowman.Block interface
func (b *Block) Reject(context.Context) error {
	b.status = choices.Rejected
	log.Debug(fmt.Sprintf("Rejecting block %s (%s) at height %d", b.ID().Hex(), b.ID(), b.Height()))
	return b.vm.blockChain.Reject(b.ethBlock)
}

// SetStatus implements the InternalBlock interface allowing ChainState
// to set the status on an existing block
func (b *Block) SetStatus(status choices.Status) { b.status = status }

// Status implements the snowman.Block interface
func (b *Block) Status() choices.Status {
	return b.status
}

// Parent implements the snowman.Block interface
func (b *Block) Parent() ids.ID {
	return ids.ID(b.ethBlock.ParentHash())
}

// Height implements the snowman.Block interface
func (b *Block) Height() uint64 {
	return b.ethBlock.NumberU64()
}

// Timestamp implements the snowman.Block interface
func (b *Block) Timestamp() time.Time {
	return time.Unix(int64(b.ethBlock.Time()), 0)
}

// syntacticVerify verifies that a *Block is well-formed.
func (b *Block) syntacticVerify() error {
	if b == nil || b.ethBlock == nil {
		return errInvalidBlock
	}

	header := b.ethBlock.Header()
	rules := b.vm.chainConfig.AvalancheRules(header.Number, new(big.Int).SetUint64(header.Time))
	return b.vm.syntacticBlockValidator.SyntacticVerify(b, rules)
}

// Verify implements the snowman.Block interface
// Since VerifyWithContext is implemented, we never expect this to be called.
func (b *Block) Verify(context.Context) error {
	return b.verify(nil, true)
}

// ShouldVerifyWithContext implements the block.WithVerifyContext interface
// TODO: Cache the result such that if called multiple times for the same block,
// we do not need to recompute the value.
func (b *Block) ShouldVerifyWithContext(context.Context) (bool, error) {
	precompileConfigs := b.vm.currentRules().Precompiles
	for _, tx := range b.ethBlock.Transactions() {
		for _, accessTuple := range tx.AccessList() {
			if _, ok := precompileConfigs[accessTuple.Address]; ok {
				log.Debug("Should verify block with proposer VM context", "block", b.ID(), "height", b.Height())
				return true, nil
			}
		}
	}
	log.Debug("Block does not require proposer VM context for verification.", "block", b.ID(), "height", b.Height())

	return false, nil
}

// VerifyWithContext implements the block.WithVerifyContext interface
func (b *Block) VerifyWithContext(ctx context.Context, proposerVMBlockCtx *block.Context) error {
	if proposerVMBlockCtx != nil {
		log.Debug("Verifying block with context", "block", b.ID(), "height", b.Height())
	} else {
		log.Debug("Verifying block without context", "block", b.ID(), "height", b.Height())
	}
	return b.verify(proposerVMBlockCtx, true)
}

func (b *Block) verify(proposerVMBlockCtx *block.Context, writes bool) error {
	if err := b.syntacticVerify(); err != nil {
		return fmt.Errorf("syntactic block verification failed: %w", err)
	}

	// If the chain is not yet bootstrapped, we do not need to verify the transaction predicates
	// because the block is already included in the chain, so we know that predicates must have
	// been valid at the time the block was accepted.
	if b.vm.bootstrapped {
		rules := b.vm.chainConfig.AvalancheRules(b.ethBlock.Number(), b.ethBlock.Timestamp())
		if _, err := core.CheckPredicatesForSenderTxs(rules, b.vm.ctx, proposerVMBlockCtx, b.ethBlock.Transactions()); err != nil {
			return fmt.Errorf("predicate transaction verification failed: %w", err)
		}
	}

	return b.vm.blockChain.InsertBlockManual(b.ethBlock, writes)
}

// Bytes implements the snowman.Block interface
func (b *Block) Bytes() []byte {
	res, err := rlp.EncodeToBytes(b.ethBlock)
	if err != nil {
		panic(err)
	}
	return res
}

func (b *Block) String() string { return fmt.Sprintf("EVM block, ID = %s", b.ID()) }
