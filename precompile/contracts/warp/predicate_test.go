// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/engine/snowman/block"
	"github.com/ava-labs/avalanchego/snow/validators"
	"github.com/ava-labs/avalanchego/utils"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/set"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	precompileUtils "github.com/ava-labs/subnet-evm/utils"
	warpPayload "github.com/ava-labs/subnet-evm/warp/payload"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const pChainHeight uint64 = 1337

var (
	_ utils.Sortable[*testValidator] = (*testValidator)(nil)

	errTest            = errors.New("non-nil error")
	sourceChainID      = ids.GenerateTestID()
	sourceSubnetID     = ids.GenerateTestID()
	destinationChainID = ids.GenerateTestID()

	unsignedMsg           *avalancheWarp.UnsignedMessage
	addressedPayload      *warpPayload.AddressedPayload
	addressedPayloadBytes []byte
	blsSignatures         []*bls.Signature

	numTestVdrs = 10_000
	testVdrs    []*testValidator
	vdrs        map[ids.NodeID]*validators.GetValidatorOutput
	tests       []signatureTest

	predicateTests = make(map[string]testutils.PredicateTest)
)

type testValidator struct {
	nodeID ids.NodeID
	sk     *bls.SecretKey
	vdr    *avalancheWarp.Validator
}

func (v *testValidator) Less(o *testValidator) bool {
	return v.vdr.Less(o.vdr)
}

func newTestValidator() *testValidator {
	sk, err := bls.NewSecretKey()
	if err != nil {
		panic(err)
	}

	nodeID := ids.GenerateTestNodeID()
	pk := bls.PublicFromSecretKey(sk)
	return &testValidator{
		nodeID: nodeID,
		sk:     sk,
		vdr: &avalancheWarp.Validator{
			PublicKey:      pk,
			PublicKeyBytes: pk.Serialize(),
			Weight:         3,
			NodeIDs:        []ids.NodeID{nodeID},
		},
	}
}

type signatureTest struct {
	name      string
	stateF    func(*gomock.Controller) validators.State
	quorumNum uint64
	quorumDen uint64
	msgF      func(*require.Assertions) *avalancheWarp.Message
	err       error
}

func init() {
	testVdrs = make([]*testValidator, 0, numTestVdrs)
	for i := 0; i < numTestVdrs; i++ {
		testVdrs = append(testVdrs, newTestValidator())
	}
	utils.Sort(testVdrs)

	vdrs = map[ids.NodeID]*validators.GetValidatorOutput{
		testVdrs[0].nodeID: {
			NodeID:    testVdrs[0].nodeID,
			PublicKey: testVdrs[0].vdr.PublicKey,
			Weight:    testVdrs[0].vdr.Weight,
		},
		testVdrs[1].nodeID: {
			NodeID:    testVdrs[1].nodeID,
			PublicKey: testVdrs[1].vdr.PublicKey,
			Weight:    testVdrs[1].vdr.Weight,
		},
		testVdrs[2].nodeID: {
			NodeID:    testVdrs[2].nodeID,
			PublicKey: testVdrs[2].vdr.PublicKey,
			Weight:    testVdrs[2].vdr.Weight,
		},
	}

	var err error
	addressedPayload, err = warpPayload.NewAddressedPayload(ids.GenerateTestID(), ids.GenerateTestID(), []byte{1, 2, 3})
	if err != nil {
		panic(err)
	}
	addressedPayloadBytes = addressedPayload.Bytes()
	unsignedMsg, err = avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, addressedPayload.Bytes())
	if err != nil {
		panic(err)
	}

	for _, testVdr := range testVdrs {
		blsSignature := bls.Sign(testVdr.sk, unsignedMsg.Bytes())
		blsSignatures = append(blsSignatures, blsSignature)
	}

	for _, totalNodes := range []int{10, 100, 1_000, 10_000} {
		testName := fmt.Sprintf("%d nodes %d signers", totalNodes, totalNodes)
		predicateTests[testName] = createNValidatorsAndSignersTest(totalNodes)
	}

	for _, totalNodes := range []int{10, 100, 1_000, 10_000} {
		testName := fmt.Sprintf("%d nodes 10 heavily weighted keys", totalNodes)
		predicateTests[testName] = createMissingPublicKeyTest(10, totalNodes)
	}

	for _, totalNodes := range []int{10, 100, 1_000, 10_000} {
		testName := fmt.Sprintf("%d nodes 10 duplicated keys", totalNodes)
		predicateTests[testName] = createDuplicateKeyTest(10, totalNodes)
	}
}

func createWarpMessage(numKeys int) *avalancheWarp.Message {
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
	return warpMsg
}

func createPredicate(numKeys int) []byte {
	warpMsg := createWarpMessage(numKeys)
	predicateBytes := precompileUtils.PackPredicate(warpMsg.Bytes())
	return predicateBytes
}

type validatorRange struct {
	start  int
	end    int
	weight uint64
}

func createSnowCtx(validatorRanges []validatorRange) *snow.Context {
	validatorOutput := make(map[ids.NodeID]*validators.GetValidatorOutput)

	for _, validatorRange := range validatorRanges {
		for i := validatorRange.start; i < validatorRange.end; i++ {
			validatorOutput[testVdrs[i].nodeID] = &validators.GetValidatorOutput{
				NodeID:    testVdrs[i].nodeID,
				PublicKey: testVdrs[i].vdr.PublicKey,
				Weight:    validatorRange.weight,
			}
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
	return snowCtx
}

func createNValidatorsAndSignersTest(numKeys int) testutils.PredicateTest {
	predicateBytes := createPredicate(numKeys)

	snowCtx := createSnowCtx([]validatorRange{
		{
			start:  0,
			end:    numKeys,
			weight: 20,
		},
	})

	return testutils.PredicateTest{
		Config: NewConfig(big.NewInt(0), 0),
		ProposerPredicateContext: &precompileconfig.ProposerPredicateContext{
			PrecompilePredicateContext: precompileconfig.PrecompilePredicateContext{
				SnowCtx: snowCtx,
			},
			ProposerVMBlockCtx: &block.Context{
				PChainHeight: 1,
			},
		},
		StorageSlots: predicateBytes,
		Gas:          GasCostPerSignatureVerification + uint64(len(predicateBytes))*GasCostPerWarpMessageBytes + uint64(numKeys)*GasCostPerWarpSigner,
		GasErr:       nil,
		PredicateErr: nil,
	}
}

func createMissingPublicKeyTest(numKeys int, numValidators int) testutils.PredicateTest {
	predicateBytes := createPredicate(numKeys)
	snowCtx := createSnowCtx([]validatorRange{
		{
			start:  0,
			end:    numKeys,
			weight: 10_000_000,
		},
		{
			start:  10,
			end:    numValidators,
			weight: 20,
		},
	})

	return testutils.PredicateTest{
		Config: NewConfig(big.NewInt(0), 0),
		ProposerPredicateContext: &precompileconfig.ProposerPredicateContext{
			PrecompilePredicateContext: precompileconfig.PrecompilePredicateContext{
				SnowCtx: snowCtx,
			},
			ProposerVMBlockCtx: &block.Context{
				PChainHeight: 1,
			},
		},
		StorageSlots: predicateBytes,
		Gas:          GasCostPerSignatureVerification + uint64(len(predicateBytes))*GasCostPerWarpMessageBytes + uint64(numKeys)*GasCostPerWarpSigner,
		GasErr:       nil,
		PredicateErr: nil,
	}
}

func createDuplicateKeyTest(numKeys int, numValidators int) testutils.PredicateTest {
	predicateBytes := createPredicate(numKeys)

	snowCtx := createSnowCtx([]validatorRange{
		{
			0,
			numKeys,
			10_000_000,
		},
		{
			10,
			numValidators,
			20,
		},
	})

	return testutils.PredicateTest{
		Config: NewConfig(big.NewInt(0), 0),
		ProposerPredicateContext: &precompileconfig.ProposerPredicateContext{
			PrecompilePredicateContext: precompileconfig.PrecompilePredicateContext{
				SnowCtx: snowCtx,
			},
			ProposerVMBlockCtx: &block.Context{
				PChainHeight: 1,
			},
		},
		StorageSlots: predicateBytes,
		Gas:          GasCostPerSignatureVerification + uint64(len(predicateBytes))*GasCostPerWarpMessageBytes + uint64(numKeys)*GasCostPerWarpSigner,
		GasErr:       nil,
		PredicateErr: nil,
	}
}

// This test copies the test coverage from https://github.com/ava-labs/avalanchego/blob/v1.10.0/vms/platformvm/warp/signature_test.go#L137.
// These tests are only expected to fail if there is a breaking change in AvalancheGo that unexpectedly changes behavior.
func TestSignatureVerification(t *testing.T) {
	tests = []signatureTest{
		{
			name: "can't get subnetID",
			stateF: func(ctrl *gomock.Controller) validators.State {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetSubnetID(gomock.Any(), sourceChainID).Return(sourceSubnetID, errTest)
				return state
			},
			quorumNum: 1,
			quorumDen: 2,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					sourceChainID,
					ids.Empty,
					addressedPayloadBytes,
				)
				require.NoError(err)

				msg, err := avalancheWarp.NewMessage(
					unsignedMsg,
					&avalancheWarp.BitSetSignature{},
				)
				require.NoError(err)
				return msg
			},
			err: errTest,
		},
		{
			name: "can't get validator set",
			stateF: func(ctrl *gomock.Controller) validators.State {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetSubnetID(gomock.Any(), sourceChainID).Return(sourceSubnetID, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), pChainHeight, sourceSubnetID).Return(nil, errTest)
				return state
			},
			quorumNum: 1,
			quorumDen: 2,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					sourceChainID,
					ids.Empty,
					addressedPayloadBytes,
				)
				require.NoError(err)

				msg, err := avalancheWarp.NewMessage(
					unsignedMsg,
					&avalancheWarp.BitSetSignature{},
				)
				require.NoError(err)
				return msg
			},
			err: errTest,
		},
		{
			name: "weight overflow",
			stateF: func(ctrl *gomock.Controller) validators.State {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetSubnetID(gomock.Any(), sourceChainID).Return(sourceSubnetID, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), pChainHeight, sourceSubnetID).Return(map[ids.NodeID]*validators.GetValidatorOutput{
					testVdrs[0].nodeID: {
						NodeID:    testVdrs[0].nodeID,
						PublicKey: testVdrs[0].vdr.PublicKey,
						Weight:    math.MaxUint64,
					},
					testVdrs[1].nodeID: {
						NodeID:    testVdrs[1].nodeID,
						PublicKey: testVdrs[1].vdr.PublicKey,
						Weight:    math.MaxUint64,
					},
				}, nil)
				return state
			},
			quorumNum: 1,
			quorumDen: 2,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					sourceChainID,
					ids.Empty,
					addressedPayloadBytes,
				)
				require.NoError(err)

				msg, err := avalancheWarp.NewMessage(
					unsignedMsg,
					&avalancheWarp.BitSetSignature{
						Signers: make([]byte, 8),
					},
				)
				require.NoError(err)
				return msg
			},
			err: avalancheWarp.ErrWeightOverflow,
		},
		{
			name: "invalid bit set index",
			stateF: func(ctrl *gomock.Controller) validators.State {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetSubnetID(gomock.Any(), sourceChainID).Return(sourceSubnetID, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), pChainHeight, sourceSubnetID).Return(vdrs, nil)
				return state
			},
			quorumNum: 1,
			quorumDen: 2,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					sourceChainID,
					ids.Empty,
					addressedPayloadBytes,
				)
				require.NoError(err)

				msg, err := avalancheWarp.NewMessage(
					unsignedMsg,
					&avalancheWarp.BitSetSignature{
						Signers:   make([]byte, 1),
						Signature: [bls.SignatureLen]byte{},
					},
				)
				require.NoError(err)
				return msg
			},
			err: avalancheWarp.ErrInvalidBitSet,
		},
		{
			name: "unknown index",
			stateF: func(ctrl *gomock.Controller) validators.State {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetSubnetID(gomock.Any(), sourceChainID).Return(sourceSubnetID, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), pChainHeight, sourceSubnetID).Return(vdrs, nil)
				return state
			},
			quorumNum: 1,
			quorumDen: 2,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					sourceChainID,
					ids.Empty,
					addressedPayloadBytes,
				)
				require.NoError(err)

				signers := set.NewBits()
				signers.Add(3) // vdr oob

				msg, err := avalancheWarp.NewMessage(
					unsignedMsg,
					&avalancheWarp.BitSetSignature{
						Signers:   signers.Bytes(),
						Signature: [bls.SignatureLen]byte{},
					},
				)
				require.NoError(err)
				return msg
			},
			err: avalancheWarp.ErrUnknownValidator,
		},
		{
			name: "insufficient weight",
			stateF: func(ctrl *gomock.Controller) validators.State {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetSubnetID(gomock.Any(), sourceChainID).Return(sourceSubnetID, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), pChainHeight, sourceSubnetID).Return(vdrs, nil)
				return state
			},
			quorumNum: 1,
			quorumDen: 1,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					sourceChainID,
					ids.Empty,
					addressedPayloadBytes,
				)
				require.NoError(err)

				// [signers] has weight from [vdr[0], vdr[1]],
				// which is 6, which is less than 9
				signers := set.NewBits()
				signers.Add(0)
				signers.Add(1)

				unsignedBytes := unsignedMsg.Bytes()
				vdr0Sig := bls.Sign(testVdrs[0].sk, unsignedBytes)
				vdr1Sig := bls.Sign(testVdrs[1].sk, unsignedBytes)
				aggSig, err := bls.AggregateSignatures([]*bls.Signature{vdr0Sig, vdr1Sig})
				require.NoError(err)
				aggSigBytes := [bls.SignatureLen]byte{}
				copy(aggSigBytes[:], bls.SignatureToBytes(aggSig))

				msg, err := avalancheWarp.NewMessage(
					unsignedMsg,
					&avalancheWarp.BitSetSignature{
						Signers:   signers.Bytes(),
						Signature: aggSigBytes,
					},
				)
				require.NoError(err)
				return msg
			},
			err: avalancheWarp.ErrInsufficientWeight,
		},
		{
			name: "can't parse sig",
			stateF: func(ctrl *gomock.Controller) validators.State {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetSubnetID(gomock.Any(), sourceChainID).Return(sourceSubnetID, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), pChainHeight, sourceSubnetID).Return(vdrs, nil)
				return state
			},
			quorumNum: 1,
			quorumDen: 2,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					sourceChainID,
					ids.Empty,
					addressedPayloadBytes,
				)
				require.NoError(err)

				signers := set.NewBits()
				signers.Add(0)
				signers.Add(1)

				msg, err := avalancheWarp.NewMessage(
					unsignedMsg,
					&avalancheWarp.BitSetSignature{
						Signers:   signers.Bytes(),
						Signature: [bls.SignatureLen]byte{},
					},
				)
				require.NoError(err)
				return msg
			},
			err: avalancheWarp.ErrParseSignature,
		},
		{
			name: "no validators",
			stateF: func(ctrl *gomock.Controller) validators.State {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetSubnetID(gomock.Any(), sourceChainID).Return(sourceSubnetID, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), pChainHeight, sourceSubnetID).Return(nil, nil)
				return state
			},
			quorumNum: 1,
			quorumDen: 2,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					sourceChainID,
					ids.Empty,
					addressedPayloadBytes,
				)
				require.NoError(err)

				unsignedBytes := unsignedMsg.Bytes()
				vdr0Sig := bls.Sign(testVdrs[0].sk, unsignedBytes)
				aggSigBytes := [bls.SignatureLen]byte{}
				copy(aggSigBytes[:], bls.SignatureToBytes(vdr0Sig))

				msg, err := avalancheWarp.NewMessage(
					unsignedMsg,
					&avalancheWarp.BitSetSignature{
						Signers:   nil,
						Signature: aggSigBytes,
					},
				)
				require.NoError(err)
				return msg
			},
			err: bls.ErrNoPublicKeys,
		},
		{
			name: "invalid signature (substitute)",
			stateF: func(ctrl *gomock.Controller) validators.State {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetSubnetID(gomock.Any(), sourceChainID).Return(sourceSubnetID, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), pChainHeight, sourceSubnetID).Return(vdrs, nil)
				return state
			},
			quorumNum: 3,
			quorumDen: 5,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					sourceChainID,
					ids.Empty,
					addressedPayloadBytes,
				)
				require.NoError(err)

				signers := set.NewBits()
				signers.Add(0)
				signers.Add(1)

				unsignedBytes := unsignedMsg.Bytes()
				vdr0Sig := bls.Sign(testVdrs[0].sk, unsignedBytes)
				// Give sig from vdr[2] even though the bit vector says it
				// should be from vdr[1]
				vdr2Sig := bls.Sign(testVdrs[2].sk, unsignedBytes)
				aggSig, err := bls.AggregateSignatures([]*bls.Signature{vdr0Sig, vdr2Sig})
				require.NoError(err)
				aggSigBytes := [bls.SignatureLen]byte{}
				copy(aggSigBytes[:], bls.SignatureToBytes(aggSig))

				msg, err := avalancheWarp.NewMessage(
					unsignedMsg,
					&avalancheWarp.BitSetSignature{
						Signers:   signers.Bytes(),
						Signature: aggSigBytes,
					},
				)
				require.NoError(err)
				return msg
			},
			err: avalancheWarp.ErrInvalidSignature,
		},
		{
			name: "invalid signature (missing one)",
			stateF: func(ctrl *gomock.Controller) validators.State {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetSubnetID(gomock.Any(), sourceChainID).Return(sourceSubnetID, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), pChainHeight, sourceSubnetID).Return(vdrs, nil)
				return state
			},
			quorumNum: 3,
			quorumDen: 5,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					sourceChainID,
					ids.Empty,
					addressedPayloadBytes,
				)
				require.NoError(err)

				signers := set.NewBits()
				signers.Add(0)
				signers.Add(1)

				unsignedBytes := unsignedMsg.Bytes()
				vdr0Sig := bls.Sign(testVdrs[0].sk, unsignedBytes)
				// Don't give the sig from vdr[1]
				aggSigBytes := [bls.SignatureLen]byte{}
				copy(aggSigBytes[:], bls.SignatureToBytes(vdr0Sig))

				msg, err := avalancheWarp.NewMessage(
					unsignedMsg,
					&avalancheWarp.BitSetSignature{
						Signers:   signers.Bytes(),
						Signature: aggSigBytes,
					},
				)
				require.NoError(err)
				return msg
			},
			err: avalancheWarp.ErrInvalidSignature,
		},
		{
			name: "invalid signature (extra one)",
			stateF: func(ctrl *gomock.Controller) validators.State {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetSubnetID(gomock.Any(), sourceChainID).Return(sourceSubnetID, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), pChainHeight, sourceSubnetID).Return(vdrs, nil)
				return state
			},
			quorumNum: 3,
			quorumDen: 5,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					sourceChainID,
					ids.Empty,
					addressedPayloadBytes,
				)
				require.NoError(err)

				signers := set.NewBits()
				signers.Add(0)
				signers.Add(1)

				unsignedBytes := unsignedMsg.Bytes()
				vdr0Sig := bls.Sign(testVdrs[0].sk, unsignedBytes)
				vdr1Sig := bls.Sign(testVdrs[1].sk, unsignedBytes)
				// Give sig from vdr[2] even though the bit vector doesn't have
				// it
				vdr2Sig := bls.Sign(testVdrs[2].sk, unsignedBytes)
				aggSig, err := bls.AggregateSignatures([]*bls.Signature{vdr0Sig, vdr1Sig, vdr2Sig})
				require.NoError(err)
				aggSigBytes := [bls.SignatureLen]byte{}
				copy(aggSigBytes[:], bls.SignatureToBytes(aggSig))

				msg, err := avalancheWarp.NewMessage(
					unsignedMsg,
					&avalancheWarp.BitSetSignature{
						Signers:   signers.Bytes(),
						Signature: aggSigBytes,
					},
				)
				require.NoError(err)
				return msg
			},
			err: avalancheWarp.ErrInvalidSignature,
		},
		{
			name: "valid signature",
			stateF: func(ctrl *gomock.Controller) validators.State {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetSubnetID(gomock.Any(), sourceChainID).Return(sourceSubnetID, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), pChainHeight, sourceSubnetID).Return(vdrs, nil)
				return state
			},
			quorumNum: 1,
			quorumDen: 2,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					sourceChainID,
					ids.Empty,
					addressedPayloadBytes,
				)
				require.NoError(err)

				// [signers] has weight from [vdr[1], vdr[2]],
				// which is 6, which is greater than 4.5
				signers := set.NewBits()
				signers.Add(1)
				signers.Add(2)

				unsignedBytes := unsignedMsg.Bytes()
				vdr1Sig := bls.Sign(testVdrs[1].sk, unsignedBytes)
				vdr2Sig := bls.Sign(testVdrs[2].sk, unsignedBytes)
				aggSig, err := bls.AggregateSignatures([]*bls.Signature{vdr1Sig, vdr2Sig})
				require.NoError(err)
				aggSigBytes := [bls.SignatureLen]byte{}
				copy(aggSigBytes[:], bls.SignatureToBytes(aggSig))

				msg, err := avalancheWarp.NewMessage(
					unsignedMsg,
					&avalancheWarp.BitSetSignature{
						Signers:   signers.Bytes(),
						Signature: aggSigBytes,
					},
				)
				require.NoError(err)
				return msg
			},
			err: nil,
		},
		{
			name: "valid signature (boundary)",
			stateF: func(ctrl *gomock.Controller) validators.State {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetSubnetID(gomock.Any(), sourceChainID).Return(sourceSubnetID, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), pChainHeight, sourceSubnetID).Return(vdrs, nil)
				return state
			},
			quorumNum: 2,
			quorumDen: 3,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					sourceChainID,
					ids.Empty,
					addressedPayloadBytes,
				)
				require.NoError(err)

				// [signers] has weight from [vdr[1], vdr[2]],
				// which is 6, which meets the minimum 6
				signers := set.NewBits()
				signers.Add(1)
				signers.Add(2)

				unsignedBytes := unsignedMsg.Bytes()
				vdr1Sig := bls.Sign(testVdrs[1].sk, unsignedBytes)
				vdr2Sig := bls.Sign(testVdrs[2].sk, unsignedBytes)
				aggSig, err := bls.AggregateSignatures([]*bls.Signature{vdr1Sig, vdr2Sig})
				require.NoError(err)
				aggSigBytes := [bls.SignatureLen]byte{}
				copy(aggSigBytes[:], bls.SignatureToBytes(aggSig))

				msg, err := avalancheWarp.NewMessage(
					unsignedMsg,
					&avalancheWarp.BitSetSignature{
						Signers:   signers.Bytes(),
						Signature: aggSigBytes,
					},
				)
				require.NoError(err)
				return msg
			},
			err: nil,
		},
		{
			name: "valid signature (missing key)",
			stateF: func(ctrl *gomock.Controller) validators.State {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetSubnetID(gomock.Any(), sourceChainID).Return(sourceSubnetID, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), pChainHeight, sourceSubnetID).Return(map[ids.NodeID]*validators.GetValidatorOutput{
					testVdrs[0].nodeID: {
						NodeID:    testVdrs[0].nodeID,
						PublicKey: nil,
						Weight:    testVdrs[0].vdr.Weight,
					},
					testVdrs[1].nodeID: {
						NodeID:    testVdrs[1].nodeID,
						PublicKey: testVdrs[1].vdr.PublicKey,
						Weight:    testVdrs[1].vdr.Weight,
					},
					testVdrs[2].nodeID: {
						NodeID:    testVdrs[2].nodeID,
						PublicKey: testVdrs[2].vdr.PublicKey,
						Weight:    testVdrs[2].vdr.Weight,
					},
				}, nil)
				return state
			},
			quorumNum: 1,
			quorumDen: 3,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					sourceChainID,
					ids.Empty,
					addressedPayloadBytes,
				)
				require.NoError(err)

				// [signers] has weight from [vdr2, vdr3],
				// which is 6, which is greater than 3
				signers := set.NewBits()
				// Note: the bits are shifted because vdr[0]'s key was zeroed
				signers.Add(0) // vdr[1]
				signers.Add(1) // vdr[2]

				unsignedBytes := unsignedMsg.Bytes()
				vdr1Sig := bls.Sign(testVdrs[1].sk, unsignedBytes)
				vdr2Sig := bls.Sign(testVdrs[2].sk, unsignedBytes)
				aggSig, err := bls.AggregateSignatures([]*bls.Signature{vdr1Sig, vdr2Sig})
				require.NoError(err)
				aggSigBytes := [bls.SignatureLen]byte{}
				copy(aggSigBytes[:], bls.SignatureToBytes(aggSig))

				msg, err := avalancheWarp.NewMessage(
					unsignedMsg,
					&avalancheWarp.BitSetSignature{
						Signers:   signers.Bytes(),
						Signature: aggSigBytes,
					},
				)
				require.NoError(err)
				return msg
			},
			err: nil,
		},
		{
			name: "valid signature (duplicate key)",
			stateF: func(ctrl *gomock.Controller) validators.State {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetSubnetID(gomock.Any(), sourceChainID).Return(sourceSubnetID, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), pChainHeight, sourceSubnetID).Return(map[ids.NodeID]*validators.GetValidatorOutput{
					testVdrs[0].nodeID: {
						NodeID:    testVdrs[0].nodeID,
						PublicKey: nil,
						Weight:    testVdrs[0].vdr.Weight,
					},
					testVdrs[1].nodeID: {
						NodeID:    testVdrs[1].nodeID,
						PublicKey: testVdrs[2].vdr.PublicKey,
						Weight:    testVdrs[1].vdr.Weight,
					},
					testVdrs[2].nodeID: {
						NodeID:    testVdrs[2].nodeID,
						PublicKey: testVdrs[2].vdr.PublicKey,
						Weight:    testVdrs[2].vdr.Weight,
					},
				}, nil)
				return state
			},
			quorumNum: 2,
			quorumDen: 3,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					sourceChainID,
					ids.Empty,
					addressedPayloadBytes,
				)
				require.NoError(err)

				// [signers] has weight from [vdr2, vdr3],
				// which is 6, which meets the minimum 6
				signers := set.NewBits()
				// Note: the bits are shifted because vdr[0]'s key was zeroed
				// Note: vdr[1] and vdr[2] were combined because of a shared pk
				signers.Add(0) // vdr[1] + vdr[2]

				unsignedBytes := unsignedMsg.Bytes()
				// Because vdr[1] and vdr[2] share a key, only one of them sign.
				vdr2Sig := bls.Sign(testVdrs[2].sk, unsignedBytes)
				aggSigBytes := [bls.SignatureLen]byte{}
				copy(aggSigBytes[:], bls.SignatureToBytes(vdr2Sig))

				msg, err := avalancheWarp.NewMessage(
					unsignedMsg,
					&avalancheWarp.BitSetSignature{
						Signers:   signers.Bytes(),
						Signature: aggSigBytes,
					},
				)
				require.NoError(err)
				return msg
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			msg := tt.msgF(require)
			pChainState := tt.stateF(ctrl)

			err := msg.Signature.Verify(
				context.Background(),
				&msg.UnsignedMessage,
				pChainState,
				pChainHeight,
				tt.quorumNum,
				tt.quorumDen,
			)
			require.ErrorIs(err, tt.err)
		})
	}
}

func TestWarpPredicate(t *testing.T) {
	testutils.RunPredicateTests(t, predicateTests)
}

func BenchmarkWarpPredicate(b *testing.B) {
	testutils.RunPredicateBenchmarks(b, predicateTests)
}

func TestWarpNilProposerCtx(t *testing.T) {
	numKeys := 1
	snowCtx := createSnowCtx([]validatorRange{
		{
			start:  0,
			end:    numKeys,
			weight: 20,
		},
	})
	predicateBytes := createPredicate(numKeys)
	test := testutils.PredicateTest{
		Config: NewConfig(big.NewInt(0), 0),
		ProposerPredicateContext: &precompileconfig.ProposerPredicateContext{
			PrecompilePredicateContext: precompileconfig.PrecompilePredicateContext{
				SnowCtx: snowCtx,
			},
			ProposerVMBlockCtx: nil,
		},
		StorageSlots: predicateBytes,
		Gas:          GasCostPerSignatureVerification + uint64(len(predicateBytes))*GasCostPerWarpMessageBytes + uint64(numKeys)*GasCostPerWarpSigner,
		GasErr:       nil,
		PredicateErr: errNoProposerPredicate,
	}

	test.Run(t)
}

func TestInvalidPredicatePacking(t *testing.T) {
	numKeys := 1
	snowCtx := createSnowCtx([]validatorRange{
		{
			start:  0,
			end:    numKeys,
			weight: 20,
		},
	})
	predicateBytes := createPredicate(numKeys)
	predicateBytes = append(predicateBytes, byte(0x01)) // Invalidate the predicate byte packing

	test := testutils.PredicateTest{
		Config: NewConfig(big.NewInt(0), 0),
		ProposerPredicateContext: &precompileconfig.ProposerPredicateContext{
			PrecompilePredicateContext: precompileconfig.PrecompilePredicateContext{
				SnowCtx: snowCtx,
			},
			ProposerVMBlockCtx: &block.Context{
				PChainHeight: 1,
			},
		},
		StorageSlots: predicateBytes,
		Gas:          GasCostPerSignatureVerification + uint64(len(predicateBytes))*GasCostPerWarpMessageBytes + uint64(numKeys)*GasCostPerWarpSigner,
		GasErr:       errInvalidPredicateBytes,
		PredicateErr: nil, // Won't be reached
	}

	test.Run(t)
}

func TestInvalidWarpMessage(t *testing.T) {
	numKeys := 1
	snowCtx := createSnowCtx([]validatorRange{
		{
			start:  0,
			end:    numKeys,
			weight: 20,
		},
	})
	warpMsg := createWarpMessage(1)
	warpMsgBytes := warpMsg.Bytes()
	warpMsgBytes = append(warpMsgBytes, byte(0x01)) // Invalidate warp message packing
	predicateBytes := precompileUtils.PackPredicate(warpMsgBytes)

	test := testutils.PredicateTest{
		Config: NewConfig(big.NewInt(0), 0),
		ProposerPredicateContext: &precompileconfig.ProposerPredicateContext{
			PrecompilePredicateContext: precompileconfig.PrecompilePredicateContext{
				SnowCtx: snowCtx,
			},
			ProposerVMBlockCtx: &block.Context{
				PChainHeight: 1,
			},
		},
		StorageSlots: predicateBytes,
		Gas:          GasCostPerSignatureVerification + uint64(len(predicateBytes))*GasCostPerWarpMessageBytes + uint64(numKeys)*GasCostPerWarpSigner,
		GasErr:       errInvalidWarpMsg,
		PredicateErr: nil, // Won't be reached
	}

	test.Run(t)
}

func TestInvalidAddressedPayload(t *testing.T) {
	numKeys := 1
	snowCtx := createSnowCtx([]validatorRange{
		{
			start:  0,
			end:    numKeys,
			weight: 20,
		},
	})
	aggregateSignature, err := bls.AggregateSignatures(blsSignatures[0:numKeys])
	require.NoError(t, err)
	bitSet := set.NewBits()
	for i := 0; i < numKeys; i++ {
		bitSet.Add(i)
	}
	warpSignature := &avalancheWarp.BitSetSignature{
		Signers: bitSet.Bytes(),
	}
	copy(warpSignature.Signature[:], bls.SignatureToBytes(aggregateSignature))
	// Create an unsigned message with an invalid addressed payload
	unsignedMsg, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, []byte{1, 2, 3})
	require.NoError(t, err)
	warpMsg, err := avalancheWarp.NewMessage(unsignedMsg, warpSignature)
	require.NoError(t, err)
	warpMsgBytes := warpMsg.Bytes()
	predicateBytes := precompileUtils.PackPredicate(warpMsgBytes)

	test := testutils.PredicateTest{
		Config: NewConfig(big.NewInt(0), 0),
		ProposerPredicateContext: &precompileconfig.ProposerPredicateContext{
			PrecompilePredicateContext: precompileconfig.PrecompilePredicateContext{
				SnowCtx: snowCtx,
			},
			ProposerVMBlockCtx: &block.Context{
				PChainHeight: 1,
			},
		},
		StorageSlots: predicateBytes,
		Gas:          GasCostPerSignatureVerification + uint64(len(predicateBytes))*GasCostPerWarpMessageBytes + uint64(numKeys)*GasCostPerWarpSigner,
		GasErr:       nil,
		PredicateErr: errInvalidAddressedPayload,
	}

	test.Run(t)
}

func TestInvalidBitSet(t *testing.T) {
	unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
		sourceChainID,
		ids.Empty,
		[]byte{1, 2, 3},
	)
	require.NoError(t, err)

	msg, err := avalancheWarp.NewMessage(
		unsignedMsg,
		&avalancheWarp.BitSetSignature{
			Signers:   make([]byte, 1),
			Signature: [bls.SignatureLen]byte{},
		},
	)
	require.NoError(t, err)

	numKeys := 1
	snowCtx := createSnowCtx([]validatorRange{
		{
			start:  0,
			end:    numKeys,
			weight: 20,
		},
	})
	predicateBytes := precompileUtils.PackPredicate(msg.Bytes())
	test := testutils.PredicateTest{
		Config: NewConfig(big.NewInt(0), 0),
		ProposerPredicateContext: &precompileconfig.ProposerPredicateContext{
			PrecompilePredicateContext: precompileconfig.PrecompilePredicateContext{
				SnowCtx: snowCtx,
			},
			ProposerVMBlockCtx: &block.Context{
				PChainHeight: 1,
			},
		},
		StorageSlots: predicateBytes,
		Gas:          GasCostPerSignatureVerification + uint64(len(predicateBytes))*GasCostPerWarpMessageBytes + uint64(numKeys)*GasCostPerWarpSigner,
		GasErr:       errCannotNumSigners,
		PredicateErr: nil, // Won't be reached
	}

	test.Run(t)
}

func TestWarpSignatureWeightsDefaultQuorumNumerator(t *testing.T) {
	snowCtx := createSnowCtx([]validatorRange{
		{
			start:  0,
			end:    100,
			weight: 20,
		},
	})

	tests := make(map[string]testutils.PredicateTest)
	for _, numSigners := range []int{1, int(DefaultQuorumNumerator) - 1, int(DefaultQuorumNumerator), int(DefaultQuorumNumerator) + 1, 99, 100, 101} {
		var (
			predicateBytes       = createPredicate(numSigners)
			expectedPredicateErr error
		)
		// If the number of signers is less than the DefaultQuorumNumerator (67)
		if numSigners < int(DefaultQuorumNumerator) {
			expectedPredicateErr = avalancheWarp.ErrInsufficientWeight
		}
		if numSigners > int(QuorumDenominator) {
			expectedPredicateErr = avalancheWarp.ErrUnknownValidator
		}
		tests[fmt.Sprintf("default quorum %d signature(s)", numSigners)] = testutils.PredicateTest{
			Config: NewConfig(big.NewInt(0), 0),
			ProposerPredicateContext: &precompileconfig.ProposerPredicateContext{
				PrecompilePredicateContext: precompileconfig.PrecompilePredicateContext{
					SnowCtx: snowCtx,
				},
				ProposerVMBlockCtx: &block.Context{
					PChainHeight: 1,
				},
			},
			StorageSlots: predicateBytes,
			Gas:          GasCostPerSignatureVerification + uint64(len(predicateBytes))*GasCostPerWarpMessageBytes + uint64(numSigners)*GasCostPerWarpSigner,
			GasErr:       nil,
			PredicateErr: expectedPredicateErr,
		}
	}
	testutils.RunPredicateTests(t, tests)
}

func TestWarpSignatureWeightsNonDefaultQuorumNumerator(t *testing.T) {
	snowCtx := createSnowCtx([]validatorRange{
		{
			start:  0,
			end:    100,
			weight: 20,
		},
	})

	tests := make(map[string]testutils.PredicateTest)
	nonDefaultQuorumNumerator := 50
	// Ensure this test fails if the DefaultQuroumNumerator is changed to an unexpected value during development
	require.NotEqual(t, nonDefaultQuorumNumerator, int(DefaultQuorumNumerator))
	// Add cases with default quorum
	for _, numSigners := range []int{nonDefaultQuorumNumerator, nonDefaultQuorumNumerator + 1, 99, 100, 101} {
		var (
			predicateBytes       = createPredicate(numSigners)
			expectedPredicateErr error
		)
		// If the number of signers is less than the quorum numerator, expect ErrInsufficientWeight
		if numSigners < nonDefaultQuorumNumerator {
			expectedPredicateErr = avalancheWarp.ErrInsufficientWeight
		}
		if numSigners > int(QuorumDenominator) {
			expectedPredicateErr = avalancheWarp.ErrUnknownValidator
		}
		name := fmt.Sprintf("non-default quorum %d signature(s)", numSigners)
		tests[name] = testutils.PredicateTest{
			Config: NewConfig(big.NewInt(0), uint64(nonDefaultQuorumNumerator)),
			ProposerPredicateContext: &precompileconfig.ProposerPredicateContext{
				PrecompilePredicateContext: precompileconfig.PrecompilePredicateContext{
					SnowCtx: snowCtx,
				},
				ProposerVMBlockCtx: &block.Context{
					PChainHeight: 1,
				},
			},
			StorageSlots: predicateBytes,
			Gas:          GasCostPerSignatureVerification + uint64(len(predicateBytes))*GasCostPerWarpMessageBytes + uint64(numSigners)*GasCostPerWarpSigner,
			GasErr:       nil,
			PredicateErr: expectedPredicateErr,
		}
	}

	testutils.RunPredicateTests(t, tests)
}
