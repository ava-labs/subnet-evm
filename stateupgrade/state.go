package stateupgrade

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type StateUpgradeStateConfigStruct struct {
	// The address of the contract that will be upgraded.
	Address common.Address `json:"address,omitempty"`

	// The memory slot.
	Slot *big.Int `json:"slot,omitempty"`

	// The new cap to be set.
	Value *big.Int `json:"value,omitempty"`
}

// StateUpgradeStateConfig implements the StateUpgradeConfig interface.
type StateUpgradeStateConfig struct {
	UpgradeableConfig
	InitialStateUpgradeStateConfig *StateUpgradeStateConfigStruct `json:"config,omitempty"`
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

	// Set the storage slot.
	log.Info("Running State Upgrader State", "config", c)
	if c.InitialStateUpgradeStateConfig != nil {
		log.Info("Setting the storage", "address", c.InitialStateUpgradeStateConfig.Address, "slot", c.InitialStateUpgradeStateConfig.Slot, "value", c.InitialStateUpgradeStateConfig.Value)
		state.SetState(c.InitialStateUpgradeStateConfig.Address, common.BigToHash(c.InitialStateUpgradeStateConfig.Slot), common.BigToHash(c.InitialStateUpgradeStateConfig.Value))
	} else {
		log.Error("State Upgrader Config is not set")
	}
}

// Verify tries to verify StateUpgradeStateConfig and returns an error accordingly.
func (c *StateUpgradeStateConfig) Verify() error {
	return nil
}
