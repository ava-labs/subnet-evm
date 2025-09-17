// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package acp224feemanager

import (
	"github.com/ava-labs/libevm/common"

	"github.com/ava-labs/subnet-evm/commontype"
)

// ACP224FeeManagerFeeConfigUpdated represents a FeeConfigUpdated non-indexed event data raised by the ACP224FeeManager contract.
type FeeConfigUpdatedEventData struct {
	OldFeeConfig commontype.ACP224FeeConfig
	NewFeeConfig commontype.ACP224FeeConfig
}

// PackFeeConfigUpdatedEvent packs the event into the appropriate arguments for FeeConfigUpdated.
// It returns topic hashes and the encoded non-indexed data.
func PackFeeConfigUpdatedEvent(sender common.Address, data FeeConfigUpdatedEventData) ([]common.Hash, []byte, error) {
	return ACP224FeeManagerABI.PackEvent("FeeConfigUpdated", sender, data.OldFeeConfig, data.NewFeeConfig)
}

// UnpackFeeConfigUpdatedEventData attempts to unpack non-indexed [dataBytes].
func UnpackFeeConfigUpdatedEventData(dataBytes []byte) (FeeConfigUpdatedEventData, error) {
	eventData := FeeConfigUpdatedEventData{}
	err := ACP224FeeManagerABI.UnpackIntoInterface(&eventData, "FeeConfigUpdated", dataBytes)
	return eventData, err
}
