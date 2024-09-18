// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/validators"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
)

var (
	DefaultStartTime      = uint64(time.Date(2024, time.July, 30, 0, 0, 0, 0, time.UTC).Unix())
	DefaultSetWeightNonce = uint64(0)
	DefaultIsActive       = true
)

type MockValidatorOutput struct {
	NodeID         ids.NodeID
	VID            ids.ID
	IsActive       bool
	StartTime      uint64
	SetWeightNonce uint64
	Weight         uint64
	BLSPublicKey   *bls.PublicKey
}

type MockedValidatorState interface {
	validators.State
	// GetCurrentValidatorSet returns the current validator set for the provided subnet
	// Returned map contains the ValidationID as the key and the ValidatorOutput as the value
	GetCurrentValidatorSet(ctx context.Context, subnetID ids.ID) (map[ids.ID]*MockValidatorOutput, error)
}

type MockValidatorState struct {
	validators.State
}

func NewMockValidatorState(pState validators.State) MockedValidatorState {
	return &MockValidatorState{
		State: pState,
	}
}

func (t *MockValidatorState) GetCurrentValidatorSet(ctx context.Context, subnetID ids.ID) (map[ids.ID]*MockValidatorOutput, error) {
	currentPHeight, err := t.GetCurrentHeight(ctx)
	if err != nil {
		return nil, err
	}
	validatorSet, err := t.GetValidatorSet(ctx, currentPHeight, subnetID)
	if err != nil {
		return nil, err
	}
	output := make(map[ids.ID]*MockValidatorOutput, len(validatorSet))
	for key, value := range validatorSet {
		// Converts the 20 bytes nodeID to a 32-bytes validationID
		// TODO: This is a temporary solution until we can use the correct ID type
		// fill bytes with 0s to make it 32 bytes
		keyBytes := make([]byte, 32)
		copy(keyBytes[:], key.Bytes())
		validationID, err := ids.ToID(keyBytes)
		if err != nil {
			return nil, err
		}
		output[validationID] = &MockValidatorOutput{
			VID:            validationID,
			NodeID:         value.NodeID,
			IsActive:       DefaultIsActive,
			StartTime:      DefaultStartTime,
			SetWeightNonce: DefaultSetWeightNonce,
			Weight:         value.Weight,
			BLSPublicKey:   value.PublicKey,
		}
	}
	return output, nil
}
