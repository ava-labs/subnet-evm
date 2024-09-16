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

type ValidatorOutput struct {
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
	GetCurrentValidatorSet(ctx context.Context, subnetID ids.ID) (map[ids.ID]*ValidatorOutput, error)
}

type recordedValidator struct {
	StartTime      uint64
	SetWeightNonce uint64
	IsActive       bool
}

type MockValidatorState struct {
	validators.State
	recordedValidators map[ids.NodeID]recordedValidator
}

func NewMockValidatorState(pState validators.State) MockedValidatorState {
	return &MockValidatorState{
		State:              pState,
		recordedValidators: make(map[ids.NodeID]recordedValidator),
	}
}

func (t *MockValidatorState) RecordValidator(nodeID ids.NodeID, startTime, setWeightNonce uint64) {
	t.recordedValidators[nodeID] = recordedValidator{
		StartTime:      startTime,
		SetWeightNonce: setWeightNonce,
		IsActive:       true,
	}
}

func (t *MockValidatorState) GetCurrentValidatorSet(ctx context.Context, subnetID ids.ID) (map[ids.ID]*ValidatorOutput, error) {
	currentPHeight, err := t.GetCurrentHeight(ctx)
	if err != nil {
		return nil, err
	}
	validatorSet, err := t.GetValidatorSet(ctx, currentPHeight, subnetID)
	if err != nil {
		return nil, err
	}
	output := make(map[ids.ID]*ValidatorOutput, len(validatorSet))
	for key, value := range validatorSet {
		startTime, isActive, setWeightNonce := DefaultStartTime, DefaultIsActive, DefaultSetWeightNonce
		if recordedValidator, ok := t.recordedValidators[key]; ok {
			startTime = recordedValidator.StartTime
			isActive = recordedValidator.IsActive
			setWeightNonce = recordedValidator.SetWeightNonce
		}
		// Converts the key to a validationID
		// TODO: This is a temporary solution until we can use the correct ID type
		validationID, err := ids.ToID(key.Bytes())
		if err != nil {
			return nil, err
		}
		output[validationID] = &ValidatorOutput{
			NodeID:         value.NodeID,
			IsActive:       isActive,
			StartTime:      startTime,
			SetWeightNonce: setWeightNonce,
			Weight:         value.Weight,
			BLSPublicKey:   value.PublicKey,
		}
	}
	return output, nil
}
