// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"context"

	"github.com/ava-labs/avalanchego/api/metrics"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/validators"
	"github.com/ava-labs/avalanchego/snow/validators/validatorstest"
	"github.com/ava-labs/avalanchego/upgrade/upgradetest"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/logging"
)

func TestSnowContext() *snow.Context {
	sk, err := bls.NewSecretKey()
	if err != nil {
		panic(err)
	}
	pk := bls.PublicFromSecretKey(sk)
	return &snow.Context{
		NetworkID:       constants.UnitTestID,
		SubnetID:        ids.Empty,
		ChainID:         ids.Empty,
		NodeID:          ids.EmptyNodeID,
		NetworkUpgrades: upgradetest.GetConfig(upgradetest.Latest),
		PublicKey:       pk,
		Log:             logging.NoLog{},
		BCLookup:        ids.NewAliaser(),
		Metrics:         metrics.NewPrefixGatherer(),
		ChainDataDir:    "",
		ValidatorState:  NewTestValidatorState(),
	}
}

func NewTestValidatorState() *validatorstest.State {
	return &validatorstest.State{
		GetCurrentHeightF: func(context.Context) (uint64, error) {
			return 0, nil
		},
		GetValidatorSetF: func(context.Context, uint64, ids.ID) (map[ids.NodeID]*validators.GetValidatorOutput, error) {
			return make(map[ids.NodeID]*validators.GetValidatorOutput), nil
		},
	}
}
