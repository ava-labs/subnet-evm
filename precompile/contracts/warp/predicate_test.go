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
	nodeIDs            []ids.NodeID
	blsSecretKeys      []*bls.SecretKey
	blsPublicKeys      []*bls.PublicKey
	unsignedMsg        *avalancheWarp.UnsignedMessage
	addressedPayload   *warpPayload.AddressedPayload
	blsSignatures      []*bls.Signature
	predicateTests     = make(map[string]testutils.PredicateTest)
	sourceChainID      = ids.GenerateTestID()
	sourceSubnetID     = ids.GenerateTestID()
	destinationChainID = ids.GenerateTestID()
)

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
	for i := 0; i < 10_000; i++ {
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

	for _, totalNodes := range []int{10, 100, 1_000, 10_000} {
		testName := fmt.Sprintf("valid warp predicate %d nodes, N validators and N signers", totalNodes)
		predicateTests[testName] = createNValidatorsAndSignersTest(totalNodes)
	}

	for _, totalNodes := range []int{10, 100, 1_000, 10_000} {
		testName := fmt.Sprintf("valid warp predicate %d nodes, 10 heavily weighted keys", totalNodes)
		predicateTests[testName] = createMissingPublicKeyTest(10, totalNodes)
	}

	for _, totalNodes := range []int{10, 100, 1_000, 10_000} {
		testName := fmt.Sprintf("valid warp predicate %d nodes, 10 duplicated keys", totalNodes)
		predicateTests[testName] = createDuplicateKeyTest(10, totalNodes)
	}
}

func createPredicate(numKeys int) []byte {
	aggregateSignature, err := bls.AggregateSignatures(blsSignatures[0:numKeys])
	if err != nil {
		panic(err)
	}
	bitSet := set.NewBits()
	for i := 0; i < numKeys; i++ {
		bitSet.Add(i)
	}
	warpSignature := &avalancheWarp.BitSetSignature{
		Signers: bitSet.Bytes(),
	}
	copy(warpSignature.Signature[:], bls.SignatureToBytes(aggregateSignature))
	warpMsg, err := avalancheWarp.NewMessage(unsignedMsg, warpSignature)
	if err != nil {
		panic(err)
	}
	predicateBytes := utils.PackPredicate(warpMsg.Bytes())
	return predicateBytes
}

func createNValidatorsAndSignersTest(numKeys int) testutils.PredicateTest {
	predicateBytes := createPredicate(numKeys)

	validatorOutput := make(map[ids.NodeID]*validators.GetValidatorOutput)
	for i := 0; i < numKeys; i++ {
		validatorOutput[nodeIDs[i]] = &validators.GetValidatorOutput{
			NodeID:    nodeIDs[i],
			PublicKey: blsPublicKeys[i],
			Weight:    20,
		}
	}
	snowCtx := snow.DefaultContextTest()
	state := &validators.TestState{
		GetSubnetIDF: func(ctx context.Context, chainID ids.ID) (ids.ID, error) {
			return sourceSubnetID, nil
		},
		GetValidatorSetF: func(ctx context.Context, height uint64, subnetID ids.ID) (map[ids.NodeID]*validators.GetValidatorOutput, error) {
			return validatorOutput, nil
		},
	}
	snowCtx.ValidatorState = state

	return testutils.PredicateTest{
		Config: NewConfig(big.NewInt(0), 0),
		ProposerVMBlockContext: &block.Context{
			PChainHeight: 1,
		},
		SnowContext:  snowCtx,
		StorageSlots: predicateBytes,
		Gas:          GasCostPerSignatureVerification + uint64(len(predicateBytes))*GasCostPerWarpMessageBytes + uint64(numKeys)*GasCostPerWarpSigner,
		GasErr:       nil,
		PredicateErr: nil,
	}
}

func createMissingPublicKeyTest(numKeys int, numValidators int) testutils.PredicateTest {
	predicateBytes := createPredicate(numKeys)

	validatorOutput := make(map[ids.NodeID]*validators.GetValidatorOutput)
	for i := 0; i < numKeys; i++ {
		validatorOutput[nodeIDs[i]] = &validators.GetValidatorOutput{
			NodeID:    nodeIDs[i],
			PublicKey: blsPublicKeys[i],
			Weight:    10_000_000,
		}
	}
	// Add remaining nodes with no BLS Public Key and negligible weight
	for i := 10; i < numValidators; i++ {
		validatorOutput[nodeIDs[i]] = &validators.GetValidatorOutput{
			NodeID:    nodeIDs[i],
			PublicKey: nil,
			Weight:    20,
		}
	}
	snowCtx := snow.DefaultContextTest()
	state := &validators.TestState{
		GetSubnetIDF: func(ctx context.Context, chainID ids.ID) (ids.ID, error) {
			return sourceSubnetID, nil
		},
		GetValidatorSetF: func(ctx context.Context, height uint64, subnetID ids.ID) (map[ids.NodeID]*validators.GetValidatorOutput, error) {
			return validatorOutput, nil
		},
	}
	snowCtx.ValidatorState = state

	return testutils.PredicateTest{
		Config: NewConfig(big.NewInt(0), 0),
		ProposerVMBlockContext: &block.Context{
			PChainHeight: 1,
		},
		SnowContext:  snowCtx,
		StorageSlots: predicateBytes,
		Gas:          GasCostPerSignatureVerification + uint64(len(predicateBytes))*GasCostPerWarpMessageBytes + uint64(numKeys)*GasCostPerWarpSigner,
		GasErr:       nil,
		PredicateErr: nil,
	}
}

func createDuplicateKeyTest(numKeys int, numValidators int) testutils.PredicateTest {
	predicateBytes := createPredicate(numKeys)

	validatorOutput := make(map[ids.NodeID]*validators.GetValidatorOutput)
	for i := 0; i < numKeys; i++ {
		validatorOutput[nodeIDs[i]] = &validators.GetValidatorOutput{
			NodeID:    nodeIDs[i],
			PublicKey: blsPublicKeys[i],
			Weight:    10_000_000,
		}
	}

	// Add remaining nodes with a duplicate BLS Public Key and negligible weight
	for i := 10; i < numValidators; i++ {
		validatorOutput[nodeIDs[i]] = &validators.GetValidatorOutput{
			NodeID:    nodeIDs[i],
			PublicKey: blsPublicKeys[i%10],
			Weight:    20,
		}
	}
	snowCtx := snow.DefaultContextTest()
	state := &validators.TestState{
		GetSubnetIDF: func(ctx context.Context, chainID ids.ID) (ids.ID, error) {
			return sourceSubnetID, nil
		},
		GetValidatorSetF: func(ctx context.Context, height uint64, subnetID ids.ID) (map[ids.NodeID]*validators.GetValidatorOutput, error) {
			return validatorOutput, nil
		},
	}
	snowCtx.ValidatorState = state

	return testutils.PredicateTest{
		Config: NewConfig(big.NewInt(0), 0),
		ProposerVMBlockContext: &block.Context{
			PChainHeight: 1,
		},
		SnowContext:  snowCtx,
		StorageSlots: predicateBytes,
		Gas:          GasCostPerSignatureVerification + uint64(len(predicateBytes))*GasCostPerWarpMessageBytes + uint64(numKeys)*GasCostPerWarpSigner,
		GasErr:       nil,
		PredicateErr: nil,
	}
}

func TestWarpPredicate(t *testing.T) {
	testutils.RunPredicateTests(t, predicateTests)
}

func BenchmarkWarpPredicate(b *testing.B) {
	testutils.RunPredicateBenchmarks(b, predicateTests)
}
