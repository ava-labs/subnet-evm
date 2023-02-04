package stateupgrade

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type StateUpgradeCodeConfigStruct struct {
	// The source contract address.
	Source common.Address `json:"source,omitempty"`

	// The address of the contract that will be upgraded.
	Target common.Address `json:"target,omitempty"`
}

// StateUpgradeCodeConfig implements the StateUpgradeConfig interface.
type StateUpgradeCodeConfig struct {
	UpgradeableConfig
	InitialStateUpgradeCodeConfig *StateUpgradeCodeConfigStruct `json:"initialStateUpgradeCodeConfig,omitempty"`
}

func init() {
	// Nothing to do here
}

// NewStateUpgradeCodeConfig returns a config for a network upgrade at [blockTimestamp] that enables
// StateUpgradeCode .
func NewStateUpgradeCodeConfig(blockTimestamp *big.Int) *StateUpgradeCodeConfig {
	return &StateUpgradeCodeConfig{
		UpgradeableConfig: UpgradeableConfig{BlockTimestamp: blockTimestamp},
	}
}

// NewDisableStateUpgradeCodeConfig returns config for a network upgrade at [blockTimestamp]
// that disables StateUpgradeCode.
func NewDisableStateUpgradeCodeConfig(blockTimestamp *big.Int) *StateUpgradeCodeConfig {
	return &StateUpgradeCodeConfig{
		UpgradeableConfig: UpgradeableConfig{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

// Equal returns true if [s] is a [*StateUpgradeCodeConfig] and it has been configured identical to [c].
func (c *StateUpgradeCodeConfig) Equal(s StateUpgradeConfig) bool {
	// typecast before comparison
	other, ok := (s).(*StateUpgradeCodeConfig)
	if !ok {
		return false
	}
	// modify this boolean accordingly with your custom StateUpgradeCodeConfig, to check if [other] and the current [c] are equal
	// if StateUpgradeCodeConfig contains only UpgradeableConfig  you can skip modifying it.
	equals := c.UpgradeableConfig.Equal(&other.UpgradeableConfig)
	return equals
}

// String returns a string representation of the StateUpgradeCodeConfig.
func (c *StateUpgradeCodeConfig) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

// Configure configures [state] with the initial configuration.
func (c *StateUpgradeCodeConfig) RunUpgrade(_ ChainConfig, state StateDB, _ BlockContext) {
	// This will be called in the first block where it is enabled.
	// 1) If BlockTimestamp is nil, this will not be called
	// 2) If BlockTimestamp is 0, this will be called while setting up the genesis block
	// 3) If BlockTimestamp is 1000, this will be called while processing the first block
	// whose timestamp is >= 1000

	// Set the code.
	log.Info("Running StateUpgradeCode Config", "config", c)
	if c.InitialStateUpgradeCodeConfig != nil {
		// Load the account of the source contract.
		log.Info("Setting the code", "source", c.InitialStateUpgradeCodeConfig.Source, "target", c.InitialStateUpgradeCodeConfig.Target, "code", state.GetCode(c.InitialStateUpgradeCodeConfig.Source))
		state.SetCode(c.InitialStateUpgradeCodeConfig.Target, state.GetCode(c.InitialStateUpgradeCodeConfig.Source))
	} else {
		log.Error("Code Upgrader Config is not set")
	}
}

// Verify tries to verify StateUpgradeCodeConfig and returns an error accordingly.
func (c *StateUpgradeCodeConfig) Verify() error {
	return nil
}
