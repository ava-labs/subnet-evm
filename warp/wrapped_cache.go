// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"github.com/ava-labs/avalanchego/cache"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/proto/pb/sdk"
	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/protobuf/proto"
)

type wrappedCache struct {
	cache.Cacher[ids.ID, []byte]
}

// NewWrappedCache takes a SDK cache that caches SignatureResponses and wraps it
// to return the Signature from the SignatureResponse.
func NewWrappedCache(sdkCache cache.Cacher[ids.ID, []byte]) cache.Cacher[ids.ID, []byte] {
	return &wrappedCache{
		Cacher: sdkCache,
	}
}

func (w *wrappedCache) Get(key ids.ID) ([]byte, bool) {
	responseBytes, ok := w.Cacher.Get(key)
	if !ok {
		return responseBytes, false
	}
	response := sdk.SignatureResponse{}
	err := proto.Unmarshal(responseBytes, &response)
	if err != nil {
		log.Error("failed to unmarshal cached SignatureResponse", "error", err)
		return nil, false
	}

	return response.Signature, true
}

func (w *wrappedCache) Put(key ids.ID, value []byte) {
	response := sdk.SignatureResponse{
		Signature: value,
	}
	responseBytes, err := proto.Marshal(&response)
	if err != nil {
		log.Error("failed to marshal SignatureResponse", "error", err)
		return
	}

	w.Cacher.Put(key, responseBytes)
}
