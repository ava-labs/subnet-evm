// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package aggregator

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/validators"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/stretchr/testify/require"
)

func newValidator(t testing.TB, weight uint64) (*bls.SecretKey, *avalancheWarp.Validator) {
	sk, err := bls.NewSecretKey()
	require.NoError(t, err)
	pk := bls.PublicFromSecretKey(sk)
	return sk, &avalancheWarp.Validator{
		PublicKey:      pk,
		PublicKeyBytes: bls.PublicKeyToBytes(pk),
		Weight:         weight,
		NodeIDs:        []ids.NodeID{ids.GenerateTestNodeID()},
	}
}

func TestAggregateSignatures(t *testing.T) {
	subnetID := ids.GenerateTestID()
	errTest := errors.New("test error")
	pChainHeight := uint64(1337)
	networkID = uint32(1338)
	unsignedMsg := &avalancheWarp.UnsignedMessage{
		NetworkID:     1337,
		SourceChainID: ids.ID{'y', 'e', 'e', 't'},
		Payload:       []byte("hello world"),
	}
	require.NoError(t, unsignedMsg.Initialize())

	nodeID1, nodeID2, nodeID3 := ids.GenerateTestNodeID(), ids.GenerateTestNodeID(), ids.GenerateTestNodeID()
	vdrWeight := uint64(10001)
	vdr1sk, vdr1 := newValidator(t, vdrWeight)
	vdr2sk, vdr2 := newValidator(t, vdrWeight+1)
	vdr3sk, vdr3 := newValidator(t, vdrWeight-1)
	sig1 := bls.Sign(vdr1sk, unsignedMsg.Bytes())
	sig2 := bls.Sign(vdr2sk, unsignedMsg.Bytes())
	sig3 := bls.Sign(vdr3sk, unsignedMsg.Bytes())
	nonVdrSk, err := bls.NewSecretKey()
	require.NoError(t, err)
	nonVdrSig := bls.Sign(nonVdrSk, unsignedMsg.Bytes())
	vdrSet := map[ids.NodeID]*validators.GetValidatorOutput{
		nodeID1: {
			NodeID:    nodeID1,
			PublicKey: vdr1.PublicKey,
			Weight:    vdr1.Weight,
		},
		nodeID2: {
			NodeID:    nodeID2,
			PublicKey: vdr2.PublicKey,
			Weight:    vdr2.Weight,
		},
		nodeID3: {
			NodeID:    nodeID3,
			PublicKey: vdr3.PublicKey,
			Weight:    vdr3.Weight,
		},
	}

	type test struct {
		name               string
		aggregatorFunc     func(*gomock.Controller) *Aggregator
		unsignedMsg        *avalancheWarp.UnsignedMessage
		quorumNum          uint64
		assertResponseFunc func(*require.Assertions, *AggregateSignatureResult)
		expectedErr        error
	}

	tests := []test{
		{
			name: "can't get height",
			aggregatorFunc: func(ctrl *gomock.Controller) *Aggregator {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetCurrentHeight(gomock.Any()).Return(uint64(0), errTest)
				return NewAggregator(subnetID, state, nil)
			},
			unsignedMsg:        nil,
			quorumNum:          0,
			assertResponseFunc: nil,
			expectedErr:        errTest,
		},
		{
			name: "can't get validator set",
			aggregatorFunc: func(ctrl *gomock.Controller) *Aggregator {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetCurrentHeight(gomock.Any()).Return(pChainHeight, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errTest)
				return NewAggregator(subnetID, state, nil)
			},
			unsignedMsg:        nil,
			assertResponseFunc: nil,
			expectedErr:        errTest,
		},
		{
			name: "no validators exist",
			aggregatorFunc: func(ctrl *gomock.Controller) *Aggregator {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetCurrentHeight(gomock.Any()).Return(pChainHeight, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				return NewAggregator(subnetID, state, nil)
			},
			unsignedMsg:        nil,
			quorumNum:          0,
			assertResponseFunc: nil,
			expectedErr:        errNoValidators,
		},
		{
			name: "0/3 validators reply with signature",
			aggregatorFunc: func(ctrl *gomock.Controller) *Aggregator {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetCurrentHeight(gomock.Any()).Return(pChainHeight, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					vdrSet, nil,
				)

				client := NewMockSignatureGetter(ctrl)
				client.EXPECT().GetSignature(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errTest).AnyTimes()
				return NewAggregator(subnetID, state, client)
			},
			unsignedMsg:        unsignedMsg,
			quorumNum:          1,
			assertResponseFunc: nil,
			expectedErr:        errInsufficientWeight,
		},
		{
			name: "1/3 validators reply with signature; insufficient weight",
			aggregatorFunc: func(ctrl *gomock.Controller) *Aggregator {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetCurrentHeight(gomock.Any()).Return(pChainHeight, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					vdrSet, nil,
				)

				client := NewMockSignatureGetter(ctrl)
				client.EXPECT().GetSignature(gomock.Any(), nodeID1, gomock.Any()).Return(sig1, nil)
				client.EXPECT().GetSignature(gomock.Any(), nodeID2, gomock.Any()).Return(nil, errTest)
				client.EXPECT().GetSignature(gomock.Any(), nodeID3, gomock.Any()).Return(nil, errTest)
				return NewAggregator(subnetID, state, client)
			},
			unsignedMsg:        unsignedMsg,
			quorumNum:          35, // Require >1/3 of weight
			assertResponseFunc: nil,
			expectedErr:        errInsufficientWeight,
		},
		{
			name: "2/3 validators reply with signature; insufficient weight",
			aggregatorFunc: func(ctrl *gomock.Controller) *Aggregator {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetCurrentHeight(gomock.Any()).Return(pChainHeight, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					vdrSet, nil,
				)

				client := NewMockSignatureGetter(ctrl)
				client.EXPECT().GetSignature(gomock.Any(), nodeID1, gomock.Any()).Return(sig1, nil)
				client.EXPECT().GetSignature(gomock.Any(), nodeID2, gomock.Any()).Return(sig2, nil)
				client.EXPECT().GetSignature(gomock.Any(), nodeID3, gomock.Any()).Return(nil, errTest)
				return NewAggregator(subnetID, state, client)
			},
			unsignedMsg:        unsignedMsg,
			quorumNum:          69, // Require >2/3 of weight
			assertResponseFunc: nil,
			expectedErr:        errInsufficientWeight,
		},
		{
			name: "2/3 validators reply with signature; sufficient weight",
			aggregatorFunc: func(ctrl *gomock.Controller) *Aggregator {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetCurrentHeight(gomock.Any()).Return(pChainHeight, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					vdrSet, nil,
				)

				client := NewMockSignatureGetter(ctrl)
				client.EXPECT().GetSignature(gomock.Any(), nodeID1, gomock.Any()).Return(sig1, nil)
				client.EXPECT().GetSignature(gomock.Any(), nodeID2, gomock.Any()).Return(sig2, nil)
				client.EXPECT().GetSignature(gomock.Any(), nodeID3, gomock.Any()).Return(nil, errTest)
				return NewAggregator(subnetID, state, client)
			},
			unsignedMsg: unsignedMsg,
			quorumNum:   65, // Require <2/3 of weight
			assertResponseFunc: func(require *require.Assertions, res *AggregateSignatureResult) {
				require.Equal(vdr1.Weight+vdr2.Weight, res.SignatureWeight)
				require.Equal(vdr1.Weight+vdr2.Weight+vdr3.Weight, res.TotalWeight)
				require.Equal(unsignedMsg, &res.Message.UnsignedMessage)

				expectedSig, err := bls.AggregateSignatures([]*bls.Signature{sig1, sig2})
				require.NoError(err)

				gotBLSSig, ok := res.Message.Signature.(*avalancheWarp.BitSetSignature)
				require.True(ok)

				require.Equal(bls.SignatureToBytes(expectedSig), gotBLSSig.Signature[:])

				numSigners, err := res.Message.Signature.NumSigners()
				require.NoError(err)
				require.Equal(2, numSigners)
			},
			expectedErr: nil,
		},
		{
			name: "3/3 validators reply with signature; sufficient weight",
			aggregatorFunc: func(ctrl *gomock.Controller) *Aggregator {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetCurrentHeight(gomock.Any()).Return(pChainHeight, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					vdrSet, nil,
				)

				client := NewMockSignatureGetter(ctrl)
				client.EXPECT().GetSignature(gomock.Any(), nodeID1, gomock.Any()).Return(sig1, nil)
				client.EXPECT().GetSignature(gomock.Any(), nodeID2, gomock.Any()).Return(sig2, nil)
				client.EXPECT().GetSignature(gomock.Any(), nodeID3, gomock.Any()).Return(sig3, nil)
				return NewAggregator(subnetID, state, client)
			},
			unsignedMsg: unsignedMsg,
			quorumNum:   100, // Require all weight
			assertResponseFunc: func(require *require.Assertions, res *AggregateSignatureResult) {
				require.Equal(vdr1.Weight+vdr2.Weight+vdr3.Weight, res.SignatureWeight)
				require.Equal(vdr1.Weight+vdr2.Weight+vdr3.Weight, res.TotalWeight)
				require.Equal(unsignedMsg, &res.Message.UnsignedMessage)

				expectedSig, err := bls.AggregateSignatures([]*bls.Signature{sig1, sig2, sig3})
				require.NoError(err)

				gotBLSSig, ok := res.Message.Signature.(*avalancheWarp.BitSetSignature)
				require.True(ok)

				require.Equal(bls.SignatureToBytes(expectedSig), gotBLSSig.Signature[:])

				numSigners, err := res.Message.Signature.NumSigners()
				require.NoError(err)
				require.Equal(3, numSigners)
			},
			expectedErr: nil,
		},
		{
			name: "3/3 validators reply with signature; 1 invalid signature; sufficient weight",
			aggregatorFunc: func(ctrl *gomock.Controller) *Aggregator {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetCurrentHeight(gomock.Any()).Return(pChainHeight, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					vdrSet, nil,
				)

				client := NewMockSignatureGetter(ctrl)
				client.EXPECT().GetSignature(gomock.Any(), nodeID1, gomock.Any()).Return(nonVdrSig, nil)
				client.EXPECT().GetSignature(gomock.Any(), nodeID2, gomock.Any()).Return(sig2, nil)
				client.EXPECT().GetSignature(gomock.Any(), nodeID3, gomock.Any()).Return(sig3, nil)
				return NewAggregator(subnetID, state, client)
			},
			unsignedMsg: unsignedMsg,
			quorumNum:   64,
			assertResponseFunc: func(require *require.Assertions, res *AggregateSignatureResult) {
				require.Equal(vdr2.Weight+vdr3.Weight, res.SignatureWeight)
				require.Equal(vdr1.Weight+vdr2.Weight+vdr3.Weight, res.TotalWeight)
				require.Equal(unsignedMsg, &res.Message.UnsignedMessage)

				expectedSig, err := bls.AggregateSignatures([]*bls.Signature{sig2, sig3})
				require.NoError(err)

				gotBLSSig, ok := res.Message.Signature.(*avalancheWarp.BitSetSignature)
				require.True(ok)

				require.Equal(bls.SignatureToBytes(expectedSig), gotBLSSig.Signature[:])

				numSigners, err := res.Message.Signature.NumSigners()
				require.NoError(err)
				require.Equal(2, numSigners)
			},
			expectedErr: nil,
		},
		{
			name: "3/3 validators reply with signature; 3 invalid signatures; insufficient weight",
			aggregatorFunc: func(ctrl *gomock.Controller) *Aggregator {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetCurrentHeight(gomock.Any()).Return(pChainHeight, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					vdrSet, nil,
				)

				client := NewMockSignatureGetter(ctrl)
				client.EXPECT().GetSignature(gomock.Any(), nodeID1, gomock.Any()).Return(nonVdrSig, nil)
				client.EXPECT().GetSignature(gomock.Any(), nodeID2, gomock.Any()).Return(nonVdrSig, nil)
				client.EXPECT().GetSignature(gomock.Any(), nodeID3, gomock.Any()).Return(nonVdrSig, nil)
				return NewAggregator(subnetID, state, client)
			},
			unsignedMsg:        unsignedMsg,
			quorumNum:          1,
			assertResponseFunc: nil,
			expectedErr:        errInsufficientWeight,
		},
		{
			name: "3/3 validators reply with signature; 2 invalid signatures; insufficient weight",
			aggregatorFunc: func(ctrl *gomock.Controller) *Aggregator {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetCurrentHeight(gomock.Any()).Return(pChainHeight, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					vdrSet, nil,
				)

				client := NewMockSignatureGetter(ctrl)
				client.EXPECT().GetSignature(gomock.Any(), nodeID1, gomock.Any()).Return(nonVdrSig, nil)
				client.EXPECT().GetSignature(gomock.Any(), nodeID2, gomock.Any()).Return(nonVdrSig, nil)
				client.EXPECT().GetSignature(gomock.Any(), nodeID3, gomock.Any()).Return(sig3, nil)
				return NewAggregator(subnetID, state, client)
			},
			unsignedMsg:        unsignedMsg,
			quorumNum:          40,
			assertResponseFunc: nil,
			expectedErr:        errInsufficientWeight,
		},
		{
			name: "2/3 validators reply with signature; 1 invalid signature; sufficient weight",
			aggregatorFunc: func(ctrl *gomock.Controller) *Aggregator {
				state := validators.NewMockState(ctrl)
				state.EXPECT().GetCurrentHeight(gomock.Any()).Return(pChainHeight, nil)
				state.EXPECT().GetValidatorSet(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					vdrSet, nil,
				)

				client := NewMockSignatureGetter(ctrl)
				client.EXPECT().GetSignature(gomock.Any(), nodeID1, gomock.Any()).Return(nonVdrSig, nil)
				client.EXPECT().GetSignature(gomock.Any(), nodeID2, gomock.Any()).Return(nil, errTest)
				client.EXPECT().GetSignature(gomock.Any(), nodeID3, gomock.Any()).Return(sig3, nil)
				return NewAggregator(subnetID, state, client)
			},
			unsignedMsg: unsignedMsg,
			quorumNum:   30,
			assertResponseFunc: func(require *require.Assertions, res *AggregateSignatureResult) {
				require.Equal(vdr3.Weight, res.SignatureWeight)
				require.Equal(vdr1.Weight+vdr2.Weight+vdr3.Weight, res.TotalWeight)
				require.Equal(unsignedMsg, &res.Message.UnsignedMessage)

				expectedSig, err := bls.AggregateSignatures([]*bls.Signature{sig3})
				require.NoError(err)

				gotBLSSig, ok := res.Message.Signature.(*avalancheWarp.BitSetSignature)
				require.True(ok)

				require.Equal(bls.SignatureToBytes(expectedSig), gotBLSSig.Signature[:])

				numSigners, err := res.Message.Signature.NumSigners()
				require.NoError(err)
				require.Equal(1, numSigners)
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			require := require.New(t)

			a := tt.aggregatorFunc(ctrl)

			res, err := a.AggregateSignatures(context.Background(), tt.unsignedMsg, tt.quorumNum)
			require.ErrorIs(err, tt.expectedErr)
			if err != nil {
				return
			}
			if tt.assertResponseFunc != nil {
				tt.assertResponseFunc(require, res)
			}
		})
	}
}
