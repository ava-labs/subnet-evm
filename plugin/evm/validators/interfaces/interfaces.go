// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package interfaces

import (
	"context"

	avalancheuptime "github.com/ava-labs/avalanchego/snow/uptime"
	stateinterfaces "github.com/ava-labs/subnet-evm/plugin/evm/validators/state/interfaces"
)

type ValidatorReader interface {
	stateinterfaces.StateReader
	avalancheuptime.Calculator
}

type Manager interface {
	stateinterfaces.State
	avalancheuptime.Manager
	Sync(ctx context.Context) error
	DispatchSync(ctx context.Context)
}
