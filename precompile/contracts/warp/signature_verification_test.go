// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"context"
	"math"
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/validators"
	"github.com/ava-labs/avalanchego/snow/validators/validatorstest"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/set"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type testValidatorState struct {
	subnetID        ids.ID
	subnetIDErr     error
	validatorSet    map[ids.NodeID]*validators.GetValidatorOutput
	validatorSetErr error
}

type signatureTest struct {
	name      string
	state     testValidatorState
	quorumNum uint64
	quorumDen uint64
	msgF      func(*require.Assertions) *avalancheWarp.Message
	err       error
}

// This test copies the test coverage from https://github.com/ava-labs/avalanchego/blob/v1.10.0/vms/platformvm/warp/signature_test.go#L137.
// These tests are only expected to fail if there is a breaking change in AvalancheGo that unexpectedly changes behavior.
func TestSignatureVerification(t *testing.T) {
	tests := []signatureTest{
		{
			name: "can't get subnetID",
			state: testValidatorState{
				subnetIDErr: errTest,
			},
			quorumDen: 2,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					networkID,
					sourceChainID,
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
			state: testValidatorState{
				subnetID:        sourceSubnetID,
				validatorSetErr: errTest,
			},
			quorumNum: 1,
			quorumDen: 2,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					networkID,
					sourceChainID,
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
			state: testValidatorState{
				subnetID: sourceSubnetID,
				validatorSet: map[ids.NodeID]*validators.GetValidatorOutput{
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
				},
			},
			quorumNum: 1,
			quorumDen: 2,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					networkID,
					sourceChainID,
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
			state: testValidatorState{
				subnetID:     sourceSubnetID,
				validatorSet: vdrs,
			},
			quorumNum: 1,
			quorumDen: 2,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					networkID,
					sourceChainID,
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
			state: testValidatorState{
				subnetID:     sourceSubnetID,
				validatorSet: vdrs,
			},
			quorumNum: 1,
			quorumDen: 2,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					networkID,
					sourceChainID,
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
			state: testValidatorState{
				subnetID:     sourceSubnetID,
				validatorSet: vdrs,
			},
			quorumNum: 1,
			quorumDen: 1,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					networkID,
					sourceChainID,
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
			state: testValidatorState{
				subnetID:     sourceSubnetID,
				validatorSet: vdrs,
			},
			quorumNum: 1,
			quorumDen: 2,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					networkID,
					sourceChainID,
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
			state: testValidatorState{
				subnetID: sourceSubnetID,
			},
			quorumNum: 1,
			quorumDen: 2,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					networkID,
					sourceChainID,
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
			state: testValidatorState{
				subnetID:     sourceSubnetID,
				validatorSet: vdrs,
			},
			quorumNum: 3,
			quorumDen: 5,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					networkID,
					sourceChainID,
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
			state: testValidatorState{
				subnetID:     sourceSubnetID,
				validatorSet: vdrs,
			},
			quorumNum: 3,
			quorumDen: 5,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					networkID,
					sourceChainID,
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
			state: testValidatorState{
				subnetID:     sourceSubnetID,
				validatorSet: vdrs,
			},
			quorumNum: 3,
			quorumDen: 5,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					networkID,
					sourceChainID,
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
			state: testValidatorState{
				subnetID:     sourceSubnetID,
				validatorSet: vdrs,
			},
			quorumNum: 1,
			quorumDen: 2,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					networkID,
					sourceChainID,
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
			state: testValidatorState{
				subnetID:     sourceSubnetID,
				validatorSet: vdrs,
			},
			quorumNum: 2,
			quorumDen: 3,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					networkID,
					sourceChainID,
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
			state: testValidatorState{
				subnetID: sourceSubnetID,
				validatorSet: map[ids.NodeID]*validators.GetValidatorOutput{
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
				},
			},
			quorumNum: 1,
			quorumDen: 3,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					networkID,
					sourceChainID,
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
			state: testValidatorState{
				subnetID: sourceSubnetID,
				validatorSet: map[ids.NodeID]*validators.GetValidatorOutput{
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
				},
			},
			quorumNum: 2,
			quorumDen: 3,
			msgF: func(require *require.Assertions) *avalancheWarp.Message {
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(
					networkID,
					sourceChainID,
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
			pChainState := &validatorstest.State{
				GetSubnetIDF: func(ctx context.Context, chainID ids.ID) (ids.ID, error) {
					return tt.state.subnetID, tt.state.subnetIDErr
				},
				GetValidatorSetF: func(ctx context.Context, height uint64, subnetID ids.ID) (map[ids.NodeID]*validators.GetValidatorOutput, error) {
					return tt.state.validatorSet, tt.state.validatorSetErr
				},
			}

			err := msg.Signature.Verify(
				context.Background(),
				&msg.UnsignedMessage,
				networkID,
				pChainState,
				pChainHeight,
				tt.quorumNum,
				tt.quorumDen,
			)
			require.ErrorIs(err, tt.err)
		})
	}
}
