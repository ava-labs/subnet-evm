package stateupgrade

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type StateUpgradeStateConfigStruct struct {
	// The source contract address.
	Source common.Address `json:"source,omitempty"`

	// The address of the contract that will be upgraded.
	Target common.Address `json:"target,omitempty"`
}

// StateUpgradeStateConfig implements the StateUpgradeConfig interface.
type StateUpgradeStateConfig struct {
	UpgradeableConfig
	InitialStateUpgradeStateConfig *StateUpgradeStateConfigStruct `json:"initialStateUpgradeStateConfig,omitempty"`
}

func init() {
	// Nothing to do here
}

// NewStateUpgradeStateConfig returns a config for a network upgrade at [blockTimestamp] that enables
// StateUpgradeState .
func NewStateUpgradeStateConfig(blockTimestamp *big.Int) *StateUpgradeStateConfig {
	return &StateUpgradeStateConfig{
		UpgradeableConfig: UpgradeableConfig{BlockTimestamp: blockTimestamp},
	}
}

// NewDisableStateUpgradeStateConfig returns config for a network upgrade at [blockTimestamp]
// that disables StateUpgradeState.
func NewDisableStateUpgradeStateConfig(blockTimestamp *big.Int) *StateUpgradeStateConfig {
	return &StateUpgradeStateConfig{
		UpgradeableConfig: UpgradeableConfig{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

// Equal returns true if [s] is a [*StateUpgradeStateConfig] and it has been configured identical to [c].
func (c *StateUpgradeStateConfig) Equal(s StateUpgradeConfig) bool {
	// typecast before comparison
	other, ok := (s).(*StateUpgradeStateConfig)
	if !ok {
		return false
	}
	// modify this boolean accordingly with your custom StateUpgradeStateConfig, to check if [other] and the current [c] are equal
	// if StateUpgradeStateConfig contains only UpgradeableConfig  you can skip modifying it.
	equals := c.UpgradeableConfig.Equal(&other.UpgradeableConfig)
	return equals
}

// String returns a string representation of the StateUpgradeStateConfig.
func (c *StateUpgradeStateConfig) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

// Configure configures [state] with the initial configuration.
func (c *StateUpgradeStateConfig) RunUpgrade(_ ChainConfig, state StateDB, _ BlockContext) {
	// This will be called in the first block where it is enabled.
	// 1) If BlockTimestamp is nil, this will not be called
	// 2) If BlockTimestamp is 0, this will be called while setting up the genesis block
	// 3) If BlockTimestamp is 1000, this will be called while processing the first block
	// whose timestamp is >= 1000

	// Set the code.
	log.Info("Running StateUpgradeState Config", "config", c)
	if c.InitialStateUpgradeStateConfig != nil {
		// Load the account of the source contract.
		log.Info("Setting the code", "source", c.InitialStateUpgradeStateConfig.Source, "target", c.InitialStateUpgradeStateConfig.Target, "code", state.GetCode(c.InitialStateUpgradeStateConfig.Source))
		state.SetCode(c.InitialStateUpgradeStateConfig.Target, state.GetCode(c.InitialStateUpgradeStateConfig.Source))
	} else {
		log.Error("Code Upgrader Config is not set")
	}
}

// Verify tries to verify StateUpgradeStateConfig and returns an error accordingly.
func (c *StateUpgradeStateConfig) Verify() error {
	return nil
}
