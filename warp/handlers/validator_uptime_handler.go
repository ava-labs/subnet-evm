// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package handlers

import (
	"errors"
	"fmt"

	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
	"github.com/ava-labs/subnet-evm/warp/messages"
	"github.com/ethereum/go-ethereum/log"
)

var errInvalidRequest = errors.New("invalid request")

type ValidatorUptimeHandler struct{}

func (v *ValidatorUptimeHandler) ValidateMessage(unsignedMessage *avalancheWarp.UnsignedMessage) error {
	parsed, err := payload.ParseAddressedCall(unsignedMessage.Payload)
	if err != nil {
		return fmt.Errorf("failed to parse payload: %w", err)
	}
	// TODO: Does nil/empty SourceAddress matter?
	if len(parsed.SourceAddress) != 0 {
		return errInvalidRequest
	}

	vdr, err := messages.ParseValidatorUptime(parsed.Payload)
	if err != nil {
		return fmt.Errorf("failed to parse validator uptime message: %w", err)
	}

	log.Info("Received validator uptime message", "validationID", vdr.ValidationID, "totalUptime", vdr.TotalUptime)
	log.Warn("Signing validator uptime message by default, not production behavior", "validationID", vdr.ValidationID, "totalUptime", vdr.TotalUptime)
	return nil
}
