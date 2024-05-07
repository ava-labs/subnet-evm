// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"crypto/rand"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ava-labs/avalanchego/genesis"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/tests/antithesis"
	"github.com/ava-labs/avalanchego/tests/fixture/tmpnet"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/set"
)

const NumKeys = 5

func main() {
	c, err := antithesis.NewConfig(os.Args)
	if err != nil {
		log.Fatalf("invalid config: %s", err)
	}

	ctx := context.Background()
	antithesis.AwaitHealthyNodes(ctx, c.URIs)

	if len(c.ChainIDs) != 1 {
		log.Fatalf("expected 1 chainID, saw %d", len(c.ChainIDs))
	}
	chainID, err := ids.FromString(c.ChainIDs[0])
	if err != nil {
		log.Fatalf("failed to parse chainID: %s", err)
	}

	genesisWorkload := &workload{
		id:      0,
		chainID: chainID,
		key:     genesis.VMRQKey,
		addrs:   set.Of(tmpnet.HardhatKey.Address()),
		uris:    c.URIs,
	}

	workloads := make([]*workload, NumKeys)
	workloads[0] = genesisWorkload

	for i := 1; i < NumKeys; i++ {
		key, err := secp256k1.NewPrivateKey()
		if err != nil {
			log.Fatalf("failed to generate key: %s", err)
		}

		// TODO(marun) Transfer funds to the new key

		workloads[i] = &workload{
			id:      i,
			chainID: chainID,
			key:     key,
			addrs:   set.Of(key.Address()),
			uris:    c.URIs,
		}
	}

	for _, w := range workloads[1:] {
		go w.run(ctx)
	}
	genesisWorkload.run(ctx)
}

type workload struct {
	id      int
	chainID ids.ID
	key     *secp256k1.PrivateKey
	addrs   set.Set[ids.ShortID]
	uris    []string
}

func (w *workload) run(ctx context.Context) {
	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}

	// TODO(marun) Check initial balance

	for {
		// TODO(marun) Exercise evm operations

		val, err := rand.Int(rand.Reader, big.NewInt(int64(time.Second)))
		if err != nil {
			log.Fatalf("failed to read randomness: %s", err)
		}

		timer.Reset(time.Duration(val.Int64()))
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		}
	}
}
