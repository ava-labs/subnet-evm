// (c) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/ava-labs/avalanchego/database/memdb"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/ava-labs/avalanchego/vms/platformvm/teleporter"
	"github.com/stretchr/testify/require"
)

var (
	sourceChainID      = ids.GenerateTestID()
	destinationChainID = ids.GenerateTestID()
	payload            = []byte("test")

	errTest = errors.New("non-nil error")
)

func TestInterfaceStructOneToOne(t *testing.T) {
	// checks struct provides at least the methods signatures in the interface
	var _ WarpBackend = (*warpBackend)(nil)
	// checks interface and struct have the same number of methods
	backendType := reflect.TypeOf(&warpBackend{})
	BackendType := reflect.TypeOf((*WarpBackend)(nil)).Elem()
	if backendType.NumMethod() != BackendType.NumMethod() {
		t.Fatalf("no 1 to 1 compliance between struct methods (%v) and interface methods (%v)", backendType.NumMethod(), BackendType.NumMethod())
	}
}

func TestWarpBackend_ValidMessage(t *testing.T) {
	db := memdb.New()

	snowCtx := snow.DefaultContextTest()
	snowCtx.TeleporterSigner = getTestSigner(t, sourceChainID)
	be := NewWarpBackend(snowCtx, db, 500)

	// Create a new unsigned message and add it to the warp backend.
	unsignedMsg, err := teleporter.NewUnsignedMessage(sourceChainID, destinationChainID, payload)
	require.NoError(t, err)
	err = be.AddMessage(context.Background(), unsignedMsg)
	require.NoError(t, err)

	// Verify that a signature is returned successfully, and compare to expected signature.
	messageID := hashing.ComputeHash256Array(unsignedMsg.Bytes())
	signature, err := be.GetSignature(context.Background(), messageID)
	require.NoError(t, err)

	expectedSig, err := snowCtx.TeleporterSigner.Sign(unsignedMsg)
	require.NoError(t, err)
	require.Equal(t, expectedSig, signature)
}

func TestWarpBackend_InvalidMessage(t *testing.T) {
	db := memdb.New()

	be := NewWarpBackend(snow.DefaultContextTest(), db, 500)
	unsignedMsg, err := teleporter.NewUnsignedMessage(sourceChainID, destinationChainID, payload)
	require.NoError(t, err)

	// Try getting a signature for a message that was not added.
	messageID := hashing.ComputeHash256Array(unsignedMsg.Bytes())
	_, err = be.GetSignature(context.Background(), messageID)
	require.Error(t, err)
}

func getTestSigner(t *testing.T, sourceID ids.ID) teleporter.Signer {
	sk, err := bls.NewSecretKey()
	require.NoError(t, err)

	return teleporter.NewSigner(sk, sourceID)
}
