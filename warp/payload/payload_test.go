// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package payload

import (
	"encoding/base64"
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils"
<<<<<<< HEAD
=======
	"github.com/ethereum/go-ethereum/common"
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	"github.com/stretchr/testify/require"
)

func TestAddressedPayload(t *testing.T) {
	require := require.New(t)

	addressedPayload, err := NewAddressedPayload(
<<<<<<< HEAD
		ids.GenerateTestID(),
		ids.GenerateTestID(),
=======
		common.Address(ids.GenerateTestShortID()),
		common.Hash(ids.GenerateTestID()),
		common.Address(ids.GenerateTestShortID()),
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
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
<<<<<<< HEAD
	base64Payload := "AAAAAAAAAQIDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEBQYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMHCAk="
	payload := &AddressedPayload{
		SourceAddress:      ids.ID{1, 2, 3},
		DestinationAddress: ids.ID{4, 5, 6},
		Payload:            []byte{7, 8, 9},
	}

	require.NoError(t, payload.Initialize())
=======
	base64Payload := "AAAAAAAAAQIDAAAAAAAAAAAAAAAAAAAAAAAEBQYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAcICQAAAAAAAAAAAAAAAAAAAAAAAAAAAwoLDA=="
	payload := &AddressedPayload{
		SourceAddress:      common.Address{1, 2, 3},
		DestinationChainID: common.Hash{4, 5, 6},
		DestinationAddress: common.Address{7, 8, 9},
		Payload:            []byte{10, 11, 12},
	}

	require.NoError(t, payload.initialize())
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d

	require.Equal(t, base64Payload, base64.StdEncoding.EncodeToString(payload.Bytes()))

	parsedPayload, err := ParseAddressedPayload(payload.Bytes())
	require.NoError(t, err)
	require.Equal(t, payload, parsedPayload)
}
