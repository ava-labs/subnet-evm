// Copyright (C) 2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package validators

import (
	"fmt"
	"sync"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/subnet-evm/plugin/evm/validators/interfaces"
	stateinterfaces "github.com/ava-labs/subnet-evm/plugin/evm/validators/state/interfaces"
)

type lockedReader struct {
	manager interfaces.Manager
	lock    sync.Locker
}

func NewLockedValidatorReader(
	manager interfaces.Manager,
	lock sync.Locker,
) interfaces.ValidatorReader {
	return &lockedReader{
		lock:    lock,
		manager: manager,
	}
}

// GetValidatorAndUptime returns the calculated uptime of the validator specified by validationID
// and the last updated time.
// GetValidatorAndUptime holds the chain context lock while performing the operation and can be called concurrently.
func (l *lockedReader) GetValidatorAndUptime(validationID ids.ID) (stateinterfaces.Validator, time.Duration, time.Time, error) {
	// lock the state
	l.lock.Lock()
	defer l.lock.Unlock()

	// Get validator first
	vdr, err := l.manager.GetValidator(validationID)
	if err != nil {
		return stateinterfaces.Validator{}, 0, time.Time{}, fmt.Errorf("failed to get validator: %w", err)
	}

	uptime, lastUpdated, err := l.manager.CalculateUptime(vdr.NodeID)
	if err != nil {
		return stateinterfaces.Validator{}, 0, time.Time{}, fmt.Errorf("failed to get uptime: %w", err)
	}

	return vdr, uptime, lastUpdated, nil
}
