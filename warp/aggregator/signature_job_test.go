// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package aggregator

import (
	"context"
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/plugin/evm/message"
	"github.com/stretchr/testify/require"
)

type mockClient struct {
	handleReq func(nodeID ids.NodeID, msg []byte) ([]byte, error)
}

func (m *mockClient) SendAppRequest(nodeID ids.NodeID, message []byte) ([]byte, error) {
	return m.handleReq(nodeID, message)
}

type signatureJobTest struct {
	ctx               context.Context
	job               *signatureJob
	expectedSignature *bls.Signature
	expectedErr       error
}

func executeSignatureJobTest(t testing.TB, test signatureJobTest) {
	blsSignature, err := test.job.Execute(test.ctx)
	if test.expectedErr != nil {
		require.ErrorIs(t, err, test.expectedErr)
		return
	}
	require.NoError(t, err)
	require.Equal(t, bls.SignatureToBytes(blsSignature), bls.SignatureToBytes(test.expectedSignature))
}

func TestSignatureRequestSuccess(t *testing.T) {
	nodeID := ids.GenerateTestNodeID()
	blsSecretKey, err := bls.NewSecretKey()
	require.NoError(t, err)
	blsPublicKey := bls.PublicFromSecretKey(blsSecretKey)
	unsignedMsg, err := avalancheWarp.NewUnsignedMessage(ids.GenerateTestID(), ids.GenerateTestID(), []byte{1, 2, 3})
	require.NoError(t, err)
	blsSignature := bls.Sign(blsSecretKey, unsignedMsg.Bytes())

	job := newSignatureJob(
		&mockClient{
			handleReq: func(_ ids.NodeID, _ []byte) ([]byte, error) {
				var response message.SignatureResponse
				signatureBytes := bls.SignatureToBytes(blsSignature)
				copy(response.Signature[:], signatureBytes)
				res, err := message.Codec.Marshal(message.Version, response)
				if err != nil {
					panic(err)
				}
				return res, nil
			},
		},
		&avalancheWarp.Validator{
			NodeIDs:   []ids.NodeID{nodeID},
			PublicKey: blsPublicKey,
			Weight:    10,
		},
		unsignedMsg,
	)

	executeSignatureJobTest(t, signatureJobTest{
		ctx:               context.Background(),
		job:               job,
		expectedSignature: blsSignature,
	})
}
