// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package handlers

import (
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/messages"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
	"github.com/stretchr/testify/require"
)

func TestValidatorUptimeHandler(t *testing.T) {
	require := require.New(t)

	v := &ValidatorUptimeHandler{}

	validationID := ids.GenerateTestID()
	totalUptime := uint64(1_000_000) // arbitrary value
	vdrUptime, err := messages.NewValidatorUptime(validationID, totalUptime)
	require.NoError(err)

	addressedCall, err := payload.NewAddressedCall(nil, vdrUptime.Bytes())
	require.NoError(err)

	networkID := uint32(0)
	sourceChain := ids.Empty
	message, err := avalancheWarp.NewUnsignedMessage(networkID, sourceChain, addressedCall.Bytes())
	require.NoError(err)

	require.NoError(v.ValidateMessage(message))
}
