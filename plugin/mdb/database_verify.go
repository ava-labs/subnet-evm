// (c) 2020-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package mdb

import (
	"context"
	"fmt"
	"time"

	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/trie"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

func (w *WithMerkleDB) VerifyMerkleRoot() error {
	it := w.merkleDB.NewIterator()
	defer it.Release()

	var (
		st    *trie.StackTrie
		acc   types.StateAccount
		owner common.Hash
	)
	verifyLast := func() error {
		if st == nil {
			return nil
		}
		if got := st.Hash(); got != acc.Root {
			return fmt.Errorf("account %x: expected root %x, got %x", owner, acc.Root, got)
		}
		if acc.Root != types.EmptyRootHash {
			log.Info("verified account", "account", owner, "root", acc.Root)
		}
		return nil
	}

	stMain := trie.NewStackTrie(nil)

	stKeys, mainKeys := 0, 0
	lastLog := time.Now()

	for it.Next() {
		if time.Since(lastLog) > 5*time.Second {
			log.Info("verifying merkle root", "stKeys", stKeys, "mainKeys", mainKeys)
			lastLog = time.Now()
		}

		keyLen := len(it.Key())
		if keyLen == 65 {
			k := it.Key()[33:]
			if err := st.Update(k, it.Value()); err != nil {
				return err
			}
			stKeys++
			continue
		}

		if err := verifyLast(); err != nil {
			return err
		}

		// should get account first
		owner = common.BytesToHash(it.Key())
		err := rlp.DecodeBytes(it.Value(), &acc)
		if err != nil {
			return fmt.Errorf("could not decode account %x: %w", owner, err)
		}
		st = trie.NewStackTrie(nil)
		if err := stMain.Update(it.Key(), it.Value()); err != nil {
			return err
		}
		mainKeys++
	}
	if err := it.Error(); err != nil {
		return err
	}
	if err := verifyLast(); err != nil {
		return err
	}

	ctx := context.Background()
	root, err := w.GetAltMerkleRoot(ctx)
	if err != nil {
		return err
	}
	if got := stMain.Hash(); got != root {
		return fmt.Errorf("main trie: expected root %x, got %x", root, got)
	}
	log.Info("verified main trie", "root", root)
	return nil
}
