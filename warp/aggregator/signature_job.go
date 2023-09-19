// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package aggregator

import (
	"context"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
)

// SignatureBackend defines the minimum network interface to perform signature aggregation
type SignatureBackend interface {
	// FetchWarpSignature attempts to fetch a BLS Signature from [nodeID] for [unsignedWarpMessage]
	FetchWarpSignature(ctx context.Context, nodeID ids.NodeID, unsignedWarpMessage *avalancheWarp.UnsignedMessage) (*bls.Signature, error)
}

// signatureJob fetches a single signature using the injected dependency SignatureBackend and returns a verified signature of the requested message.
type signatureJob struct {
	msg       *avalancheWarp.UnsignedMessage
	nodeID    ids.NodeID
	publicKey *bls.PublicKey
	weight    uint64
}
