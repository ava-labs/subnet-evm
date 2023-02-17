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
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/modules"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/choices"
	"github.com/ava-labs/avalanchego/snow/consensus/snowman"
	"github.com/ava-labs/avalanchego/snow/engine/snowman/block"
)

var (
	_ snowman.Block           = (*Block)(nil)
	_ block.WithVerifyContext = (*Block)(nil)
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
	// Call Accept for relevant precompile logs. This should apply DB operations to the VM's versionDB
	// to be committed atomically with marking this block as accepted.
	if err := b.handlePrecompileAccept(); err != nil {
		return err
	}
	if err := vm.acceptedBlockDB.Put(lastAcceptedKey, b.id[:]); err != nil {
		return fmt.Errorf("failed to put %s as the last accepted block: %w", b.ID(), err)
	}

	return vm.db.Commit()
}

// handlePrecompileAccept calls Accept on any logs generated with an active precompile address that implements
// contract.Accepter
// This function assumes that the Accept function will ONLY operate on state maintained in the VM's versiondb.
// This ensures that any DB operations are performed atomically with marking the block as accepted.
func (b *Block) handlePrecompileAccept() error {
	rules := b.vm.chainConfig.AvalancheRules(b.ethBlock.Number(), b.ethBlock.Timestamp())
	receipts := rawdb.ReadReceipts(b.vm.chaindb, b.ethBlock.Hash(), b.ethBlock.NumberU64(), b.vm.chainConfig)
	if receipts == nil {
		return fmt.Errorf("failed to read receipts for accepted block %s, height %d", b.ethBlock.Hash(), b.ethBlock.NumberU64())
	}

	for txIndex, receipt := range receipts {
		for _, log := range receipt.Logs {
			_, ok := rules.ActivePrecompiles[log.Address]
			if !ok {
				continue
			}

			module, ok := modules.GetPrecompileModuleByAddress(log.Address)
			if !ok {
				return fmt.Errorf("accepter accessed precompile config under address %s with no registered module", log.Address)
			}

			accepter, ok := module.Contract.(contract.Accepter)
			if !ok {
				continue
			}

			if err := accepter.Accept(log.TxHash, txIndex, log.Topics, log.Data); err != nil {
				return err
			}
		}
	}

	return nil
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
func (b *Block) Verify(context.Context) error {
	return b.verify(nil, true)
}

// ShouldVerifyWithContext implements the block.WithVerifyContext interface
func (b *Block) ShouldVerifyWithContext(context.Context) (bool, error) {
	precompileConfigs := b.vm.chainConfig.AvalancheRules(b.ethBlock.Number(), b.ethBlock.Timestamp()).ActivePrecompiles

	// Check if any of the transactions in the block list a precompile address in the access list.
	for _, tx := range b.ethBlock.Transactions() {
		for _, accessTuple := range tx.AccessList() {
			if _, ok := precompileConfigs[accessTuple.Address]; ok {
				log.Debug("Block verification requires proposerVM context", "block", b.ID(), "height", b.Height())
				return true, nil
			}
		}
	}

	log.Debug("Block verification does not require proposerVM context", "block", b.ID(), "height", b.Height())
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

// Verify the block is valid.
// Enforces that the predicates are valid within [proposerVMBlockCtx].
// Writes the block to disk iff writes=true
func (b *Block) verify(proposerVMBlockCtx *block.Context, writes bool) error {
	if err := b.syntacticVerify(); err != nil {
		return fmt.Errorf("syntactic block verification failed: %w", err)
	}

	// Only enforce predicates if the chain has already bootstrapped.
	// If the chain is still bootstrapping, we can assume that all blocks we are verifying have
	// been accepted by the network (so the predicate was validated by the network when the
	// block was originally verified).
	if b.vm.bootstrapped {
		if err := b.verifyPredicates(proposerVMBlockCtx); err != nil {
			return fmt.Errorf("failed to verify predicates: %w", err)
		}
	}

	// The engine may call VerifyWithContext multiple times on the same block with different contexts.
	// Since the engine will only call Accept/Reject once, we should only call InsertBlockManual once.
	// Additionally, if a block is already in processing, then it has already passed verification and
	// at this point we have checked the predicates are still valid in the different context so we
	// can return nil.
	if b.vm.State.IsProcessing(b.id) {
		return nil
	}

	return b.vm.blockChain.InsertBlockManual(b.ethBlock, writes)
}

// verifyPredicates verifies the predicates in the block are valid according to proposerVMBlockCtx.
func (b *Block) verifyPredicates(proposerVMBlockCtx *block.Context) error {
	rules := b.vm.chainConfig.AvalancheRules(b.ethBlock.Number(), b.ethBlock.Timestamp())
	predicateCtx := &contract.PredicateContext{
		SnowCtx:            b.vm.ctx,
		ProposerVMBlockCtx: proposerVMBlockCtx,
	}

	for _, tx := range b.ethBlock.Transactions() {
		if err := core.CheckPredicates(rules, predicateCtx, tx); err != nil {
			return err
		}
	}
	return nil
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
