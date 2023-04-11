// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/engine/snowman/block"
	"github.com/ava-labs/avalanchego/snow/validators"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/set"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/ava-labs/subnet-evm/utils"
	warpPayload "github.com/ava-labs/subnet-evm/warp/payload"
)

var (
	nodeIDs             []ids.NodeID
	blsSecretKeys       []*bls.SecretKey
	blsPublicKeys       []*bls.PublicKey
	unsignedMsg         *avalancheWarp.UnsignedMessage
	addressedPayload    *warpPayload.AddressedPayload
	blsSignatures       []*bls.Signature
	predicateTests      = make(map[string]testutils.PredicateTest)
	sourceChainID       = ids.GenerateTestID()
	sourceSubnetID      = ids.GenerateTestID()
	destinationChainID  = ids.GenerateTestID()
	getExpectedSubnetID = func(ctx context.Context, chainID ids.ID) (ids.ID, error) {
		if chainID == sourceChainID {
			return sourceSubnetID, nil
		} else {
			return ids.ID{}, fmt.Errorf("unexpected blockchainID: %s", chainID)
		}
	}
)

func produceGetValidatorSetF(
	expectedHeight uint64,
	expectedSubnetID ids.ID,
	res map[ids.NodeID]*validators.GetValidatorOutput,
	err error,
) func(ctx context.Context, height uint64, subnetID ids.ID) (map[ids.NodeID]*validators.GetValidatorOutput, error) {
	return func(ctx context.Context, height uint64, subnetID ids.ID) (map[ids.NodeID]*validators.GetValidatorOutput, error) {
		if height != expectedHeight {
			return nil, fmt.Errorf("height (expected: %d, actual: %d)", expectedHeight, height)
		}
		if subnetID != expectedSubnetID {
			return nil, fmt.Errorf("subnetID (expected: %s, actual: %s)", expectedSubnetID, subnetID)
		}

		return res, err
	}
}

func init() {
	var err error
	addressedPayload, err = warpPayload.NewAddressedPayload(ids.GenerateTestID(), ids.GenerateTestID(), []byte{1, 2, 3})
	if err != nil {
		panic(err)
	}
	unsignedMsg, err = avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, addressedPayload.Bytes())
	if err != nil {
		panic(err)
	}
	for i := 0; i < 1000; i++ {
		nodeIDs = append(nodeIDs, ids.GenerateTestNodeID())

		blsSecretKey, err := bls.NewSecretKey()
		if err != nil {
			panic(err)
		}
		blsPublicKey := bls.PublicFromSecretKey(blsSecretKey)
		blsSignature := bls.Sign(blsSecretKey, unsignedMsg.Bytes())
		blsSecretKeys = append(blsSecretKeys, blsSecretKey)
		blsPublicKeys = append(blsPublicKeys, blsPublicKey)
		blsSignatures = append(blsSignatures, blsSignature)
	}

	fiveRegisteredNodes := make(map[ids.NodeID]*validators.GetValidatorOutput)
	for i := 0; i < 5; i++ {
		fiveRegisteredNodes[nodeIDs[i]] = &validators.GetValidatorOutput{
			NodeID:    nodeIDs[i],
			PublicKey: blsPublicKeys[i],
			Weight:    20,
		}
	}

	fiveNodeAggregateSignature, err := bls.AggregateSignatures(blsSignatures[0:5])
	if err != nil {
		panic(err)
	}
	fiveNodeBitSet := set.NewBits(0, 1, 2, 3, 4)
	fiveNodeWarpSignature := &avalancheWarp.BitSetSignature{
		Signers: fiveNodeBitSet.Bytes(),
	}
	copy(fiveNodeWarpSignature.Signature[:], bls.SignatureToBytes(fiveNodeAggregateSignature))
	fiveNodeWarpMessage, err := avalancheWarp.NewMessage(unsignedMsg, fiveNodeWarpSignature)
	if err != nil {
		panic(err)
	}

	fiveNodeSnowContext := snow.DefaultContextTest()
	fiveNodeSnowContext.ValidatorState = &validators.TestState{
		GetSubnetIDF:     getExpectedSubnetID,
		GetValidatorSetF: produceGetValidatorSetF(1, sourceSubnetID, fiveRegisteredNodes, nil),
	}

	fiveNodeWarpMessagePredicateBytes := utils.PackPredicate(fiveNodeWarpMessage.Bytes())

	predicateTests["valid warp predicate"] = testutils.PredicateTest{
		Config: NewConfig(big.NewInt(1), 0),
		ProposerVMBlockContext: &block.Context{
			PChainHeight: 1,
		},
		SnowContext:  fiveNodeSnowContext,
		StorageSlots: fiveNodeWarpMessagePredicateBytes,
		Gas:          5*GasCostPerWarpSigner + uint64(len(fiveNodeWarpMessagePredicateBytes))*GasCostPerWarpMessageBytes,
		GasErr:       nil,
		PredicateErr: nil,
	}

	// TODO: implement the following test cases:
	// copy all test cases from avalanchego/vms/platformvm/warp/signature_test.go (https://github.com/ava-labs/avalanchego/blob/master/vms/platformvm/warp/signature_test.go#L165)

	// Add the following cases with the following numbers of total validators and BLS keys
	// 10 validators 10 keys (10 heavily weighted with no duplicate validator keys)
	// 100 validators 10 keys
	// 1000 validators 10 keys
	// 10 validators 10 keys (duplicate keys as necessary for every validator to be assigned a key)
	// 100 validators 10 keys
	// 1000 validators 10 keys
	// 10 validators 10 keys
	// 100 validators 100 keys
	// 1000 validators 1000 keys
	// 10000 validators 10000 keys (optional)
}

func TestWarpPredicate(t *testing.T) {
	testutils.RunPredicateTests(t, predicateTests)
}

func BenchmarkWarpPredicate(b *testing.B) {
	testutils.RunPredicateBenchmarks(b, predicateTests)
}
