// Code generated
// This file is a generated precompile contract config with stubbed abstract functions.
// The file is generated by a template. Please inspect every code and comment in this file before use.

package rewardmanager

import (
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
)

const (
	FeeRecipientsAllowedEventGasCost = contract.LogGas + contract.LogTopicGas*1                               // 1 indexed topic
	RewardAddressChangedEventGasCost = contract.LogGas + contract.LogTopicGas*2 + contract.ReadGasCostPerSlot // 1 indexed topic + reading oldRewardAddress from state
	RewardsDisabledEventGasCost      = contract.LogGas + contract.LogTopicGas*1                               // 1 indexed topic
)

// PackFeeRecipientsAllowedEvent packs the event into the appropriate arguments for FeeRecipientsAllowed.
// It returns topic hashes and the encoded non-indexed data.
func PackFeeRecipientsAllowedEvent(sender common.Address) ([]common.Hash, []byte, error) {
	return RewardManagerABI.PackEvent("FeeRecipientsAllowed", sender)
}

// PackRewardAddressChangedEvent packs the event into the appropriate arguments for RewardAddressChanged.
// It returns topic hashes and the encoded non-indexed data.
func PackRewardAddressChangedEvent(sender common.Address, oldRewardAddress common.Address, newRewardAddress common.Address) ([]common.Hash, []byte, error) {
	return RewardManagerABI.PackEvent("RewardAddressChanged", sender, oldRewardAddress, newRewardAddress)
}

// PackRewardsDisabledEvent packs the event into the appropriate arguments for RewardsDisabled.
// It returns topic hashes and the encoded non-indexed data.
func PackRewardsDisabledEvent(sender common.Address) ([]common.Hash, []byte, error) {
	return RewardManagerABI.PackEvent("RewardsDisabled", sender)
}
