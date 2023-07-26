// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package payload

import (
	"encoding/base64"
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils"
	"github.com/stretchr/testify/require"
)

func TestAddressedPayload(t *testing.T) {
	require := require.New(t)

	addressedPayload, err := NewAddressedPayload(
		ids.GenerateTestID(),
		ids.GenerateTestID(),
		ids.GenerateTestID(),
		[]byte("payload"),
	)
	require.NoError(err)

	addressedPayloadBytes := addressedPayload.Bytes()
	addressedPayload2, err := ParseAddressedPayload(addressedPayloadBytes)
	require.NoError(err)
	require.Equal(addressedPayload, addressedPayload2)
}

func TestParseAddressedPayloadJunk(t *testing.T) {
	_, err := ParseAddressedPayload(utils.RandomBytes(1024))
	require.Error(t, err)
}

func TestParseAddressedPayload(t *testing.T) {
	base64Payload := "AAAAAAAAAQIDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEBQYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAcICQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwoLDA=="
	payload := &AddressedPayload{
		SourceAddress:      ids.ID{1, 2, 3},
		DestinationChainID: ids.ID{4, 5, 6},
		DestinationAddress: ids.ID{7, 8, 9},
		Payload:            []byte{10, 11, 12},
	}

	require.NoError(t, payload.initialize())

	require.Equal(t, base64Payload, base64.StdEncoding.EncodeToString(payload.Bytes()))

	parsedPayload, err := ParseAddressedPayload(payload.Bytes())
	require.NoError(t, err)
	require.Equal(t, payload, parsedPayload)
}
