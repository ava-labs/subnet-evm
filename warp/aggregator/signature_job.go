// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package aggregator

import (
	"context"
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/plugin/evm/message"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ClientBackend defines the minimum network interface to perform signature aggregation
type ClientBackend interface {
	SendAppRequest(nodeID ids.NodeID, message []byte) ([]byte, error)
}

type signatureJob struct {
	client ClientBackend
	msg    *avalancheWarp.UnsignedMessage

	nodeID    ids.NodeID
	publicKey *bls.PublicKey
	weight    uint64
}

func (s *signatureJob) String() string {
	return fmt.Sprintf("(NodeID: %s, UnsignedMsgID: %s)", s.nodeID, s.msg.ID())
}

func newSignatureJob(client ClientBackend, validator *avalancheWarp.Validator, msg *avalancheWarp.UnsignedMessage) *signatureJob {
	return &signatureJob{
		client:    client,
		msg:       msg,
		nodeID:    validator.NodeIDs[0], // XXX: should we attempt to fetch from all nodeIDs and use the first valid response?
		publicKey: validator.PublicKey,
		weight:    validator.Weight,
	}
}

func (s *signatureJob) Execute(ctx context.Context) (*bls.Signature, error) {
	signatureReq := message.SignatureRequest{
		MessageID: s.msg.ID(),
	}
	signatureReqBytes, err := message.RequestToBytes(message.Codec, signatureReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal signature request: %w", err)
	}

	for ctx.Err() == nil {
		signatureRes, err := s.client.SendAppRequest(s.nodeID, signatureReqBytes)
		if err != nil {
			return nil, err
		}

		var response message.SignatureResponse
		if _, err := message.Codec.Unmarshal(signatureRes, &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal signature res: %w", err)
		}

		blsSignature, err := bls.SignatureFromBytes(response.Signature[:])
		if err != nil {
			return nil, fmt.Errorf("failed to parse signature from res: %w", err)
		}
		if !bls.Verify(s.publicKey, blsSignature, s.msg.Bytes()) {
			return nil, fmt.Errorf("node %s returned invalid signature %s for msg %s", s.nodeID, hexutil.Bytes(response.Signature[:]), s.msg.ID())
		}
		return blsSignature, nil
	}

	return nil, fmt.Errorf("ctx expired fetching signature for message %s from %s: %w", s.msg.ID(), s.nodeID, ctx.Err())
}
