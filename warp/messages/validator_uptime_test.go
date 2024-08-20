// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package messages

import (
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/stretchr/testify/require"
)

func TestVerifyMessage(t *testing.T) {
	require := require.New(t)

	validationID := ids.GenerateTestID()
	totalUptime := uint64(1_000_000) // arbitrary value
	vdrUptime, err := NewValidatorUptime(validationID, totalUptime)
	require.NoError(err)

	require.NoError(vdrUptime.VerifyMesssage(nil))
}
