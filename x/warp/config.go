// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"context"
	"errors"
	"fmt"
<<<<<<< HEAD
	"math/big"

	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
=======

	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/params"
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	predicateutils "github.com/ava-labs/subnet-evm/utils/predicate"
	warpPayload "github.com/ava-labs/subnet-evm/warp/payload"
	warpValidators "github.com/ava-labs/subnet-evm/warp/validators"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/log"
)

<<<<<<< HEAD
const (
	QuorumDenominator      uint64 = 100
	DefaultQuorumNumerator uint64 = 67
	MinQuorumNumerator     uint64 = 33
)

=======
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
var (
	_ precompileconfig.Config             = &Config{}
	_ precompileconfig.ProposerPredicater = &Config{}
	_ precompileconfig.Accepter           = &Config{}
)

var (
	errOverflowSignersGasCost  = errors.New("overflow calculating warp signers gas cost")
<<<<<<< HEAD
	errNoProposerPredicate     = errors.New("cannot verify warp predicate without proposer context")
	errInvalidPredicateBytes   = errors.New("cannot unpack predicate bytes")
	errInvalidWarpMsg          = errors.New("cannot unpack warp message")
	errInvalidAddressedPayload = errors.New("cannot unpack addressed payload")
	errCannotNumSigners        = errors.New("cannot fetch num signers from warp message")
=======
	errNoProposerCtxPredicate  = errors.New("cannot verify warp predicate without proposer context")
	errInvalidPredicateBytes   = errors.New("cannot unpack predicate bytes")
	errInvalidWarpMsg          = errors.New("cannot unpack warp message")
	errInvalidAddressedPayload = errors.New("cannot unpack addressed payload")
	errCannotGetNumSigners     = errors.New("cannot fetch num signers from warp message")
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
)

// Config implements the precompileconfig.Config interface and
// adds specific configuration for Warp.
type Config struct {
	precompileconfig.Upgrade
	QuorumNumerator uint64 `json:"quorumNumerator"`
}

// NewConfig returns a config for a network upgrade at [blockTimestamp] that enables
// Warp with the given quorum numerator.
<<<<<<< HEAD
func NewConfig(blockTimestamp *big.Int, quorumNumerator uint64) *Config {
=======
func NewConfig(blockTimestamp *uint64, quorumNumerator uint64) *Config {
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	return &Config{
		Upgrade:         precompileconfig.Upgrade{BlockTimestamp: blockTimestamp},
		QuorumNumerator: quorumNumerator,
	}
}

// NewDefaultConfig returns a config for a network upgrade at [blockTimestamp] that enables
// Warp with the default quorum numerator (0 denotes using the default).
<<<<<<< HEAD
func NewDefaultConfig(blockTimestamp *big.Int) *Config {
=======
func NewDefaultConfig(blockTimestamp *uint64) *Config {
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	return NewConfig(blockTimestamp, 0)
}

// NewDisableConfig returns config for a network upgrade at [blockTimestamp]
// that disables Warp.
<<<<<<< HEAD
func NewDisableConfig(blockTimestamp *big.Int) *Config {
=======
func NewDisableConfig(blockTimestamp *uint64) *Config {
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	return &Config{
		Upgrade: precompileconfig.Upgrade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

// Key returns the key for the Warp precompileconfig.
// This should be the same key as used in the precompile module.
func (*Config) Key() string { return ConfigKey }

// Verify tries to verify Config and returns an error accordingly.
<<<<<<< HEAD
func (c *Config) Verify() error {
	if c.QuorumNumerator > QuorumDenominator {
		return fmt.Errorf("cannot specify quorum numerator (%d) > quorum denominator (%d)", c.QuorumNumerator, QuorumDenominator)
	}
	// If a non-default quorum numerator is specified and it is less than the minimum, return an error
	if c.QuorumNumerator != 0 && c.QuorumNumerator < MinQuorumNumerator {
		return fmt.Errorf("cannot specify quorum numerator (%d) < min quorum numerator (%d)", c.QuorumNumerator, MinQuorumNumerator)
=======
func (c *Config) Verify(precompileconfig.ChainConfig) error {
	if c.QuorumNumerator > params.WarpQuorumDenominator {
		return fmt.Errorf("cannot specify quorum numerator (%d) > quorum denominator (%d)", c.QuorumNumerator, params.WarpQuorumDenominator)
	}
	// If a non-default quorum numerator is specified and it is less than the minimum, return an error
	if c.QuorumNumerator != 0 && c.QuorumNumerator < params.WarpQuorumNumeratorMinimum {
		return fmt.Errorf("cannot specify quorum numerator (%d) < min quorum numerator (%d)", c.QuorumNumerator, params.WarpQuorumNumeratorMinimum)
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	}
	return nil
}

// Equal returns true if [s] is a [*Config] and it has been configured identical to [c].
func (c *Config) Equal(s precompileconfig.Config) bool {
	// typecast before comparison
	other, ok := (s).(*Config)
	if !ok {
		return false
	}
	equals := c.Upgrade.Equal(&other.Upgrade)
	return equals && c.QuorumNumerator == other.QuorumNumerator
}

func (c *Config) Accept(acceptCtx *precompileconfig.AcceptContext, txHash common.Hash, logIndex int, topics []common.Hash, logData []byte) error {
	unsignedMessage, err := warp.ParseUnsignedMessage(logData)
	if err != nil {
		return fmt.Errorf("failed to parse warp log data into unsigned message (TxHash: %s, LogIndex: %d): %w", txHash, logIndex, err)
	}
	log.Info("Accepted warp unsigned message", "txHash", txHash, "logIndex", logIndex, "logData", common.Bytes2Hex(logData))
	if err := acceptCtx.Warp.AddMessage(unsignedMessage); err != nil {
		return fmt.Errorf("failed to add warp message during accept (TxHash: %s, LogIndex: %d): %w", txHash, logIndex, err)
	}
	return nil
}

// verifyWarpMessage checks that [warpMsg] can be parsed as an addressed payload and verifies the Warp Message Signature
// within [predicateContext].
func (c *Config) verifyWarpMessage(predicateContext *precompileconfig.ProposerPredicateContext, warpMsg *warp.Message) error {
<<<<<<< HEAD
	// Use default quourum numerator unless config specifies a non-default option
	quorumNumerator := DefaultQuorumNumerator
=======
	// Use default quorum numerator unless config specifies a non-default option
	quorumNumerator := params.WarpDefaultQuorumNumerator
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	if c.QuorumNumerator != 0 {
		quorumNumerator = c.QuorumNumerator
	}

	// Verify the warp payload can be decoded to the expected type
	_, err := warpPayload.ParseAddressedPayload(warpMsg.UnsignedMessage.Payload)
	if err != nil {
		return fmt.Errorf("%w: %s", errInvalidAddressedPayload, err)
	}

<<<<<<< HEAD
	log.Info("verifyingWarpMessage", "warpMsg", warpMsg, "quorumNum", quorumNumerator, "quorumDenom", QuorumDenominator)
	if err := warpMsg.Signature.Verify(
		context.Background(),
		&warpMsg.UnsignedMessage,
		warpValidators.NewState(predicateContext.SnowCtx), // Wrap validators.State on the chain snow context to special case the Primary Network
		predicateContext.ProposerVMBlockCtx.PChainHeight,
		quorumNumerator,
		QuorumDenominator,
=======
	log.Debug("verifying warp message", "warpMsg", warpMsg, "quorumNum", quorumNumerator, "quorumDenom", params.WarpQuorumDenominator)
	if err := warpMsg.Signature.Verify(
		context.Background(),
		&warpMsg.UnsignedMessage,
		predicateContext.SnowCtx.NetworkID,
		warpValidators.NewState(predicateContext.SnowCtx), // Wrap validators.State on the chain snow context to special case the Primary Network
		predicateContext.ProposerVMBlockCtx.PChainHeight,
		quorumNumerator,
		params.WarpQuorumDenominator,
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	); err != nil {
		return fmt.Errorf("warp signature verification failed: %w", err)
	}

	return nil
}

// PredicateGas returns the amount of gas necessary to verify the predicate
<<<<<<< HEAD
=======
// PredicateGas charges for:
// 1. Base cost of the message
// 2. Size of the message
// 3. Number of signers
// 4. TODO: Lookup of the validator set
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
func (c *Config) PredicateGas(predicateBytes []byte) (uint64, error) {
	totalGas := GasCostPerSignatureVerification
	bytesGasCost, overflow := math.SafeMul(GasCostPerWarpMessageBytes, uint64(len(predicateBytes)))
	if overflow {
		return 0, fmt.Errorf("overflow calculating gas cost for warp message bytes of size %d", len(predicateBytes))
	}
	totalGas, overflow = math.SafeAdd(totalGas, bytesGasCost)
	if overflow {
		return 0, fmt.Errorf("overflow adding bytes gas cost of size %d", len(predicateBytes))
	}

	unpackedPredicateBytes, err := predicateutils.UnpackPredicate(predicateBytes)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", errInvalidPredicateBytes, err)
	}
	warpMessage, err := warp.ParseMessage(unpackedPredicateBytes)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", errInvalidWarpMsg, err)
	}

	numSigners, err := warpMessage.Signature.NumSigners()
	if err != nil {
<<<<<<< HEAD
		return 0, fmt.Errorf("%w: %s", errCannotNumSigners, err)
=======
		return 0, fmt.Errorf("%w: %s", errCannotGetNumSigners, err)
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	}
	signerGas, overflow := math.SafeMul(uint64(numSigners), GasCostPerWarpSigner)
	if overflow {
		return 0, errOverflowSignersGasCost
	}
	totalGas, overflow = math.SafeAdd(totalGas, signerGas)
	if overflow {
		return 0, fmt.Errorf("overflow adding signer gas (PrevTotal: %d, VerificationGas: %d)", totalGas, signerGas)
	}

	// TODO: charge for the Subnet validator set lookup
	// ctx := context.Background()
	// subnetID, err := predicateContext.SnowCtx.ValidatorState.GetSubnetID(ctx, warpMessage.SourceChainID)
	// if err != nil {
	// 	return 0, fmt.Errorf("failed to look up SubnetID for SourceChainID: %s", warpMessage.SourceChainID)
	// }
	// validatorSet, err := predicateContext.SnowCtx.ValidatorState.GetValidatorSet(ctx, predicateContext.ProposerVMBlockCtx.PChainHeight, subnetID)
	// if err != nil {
	// 	return 0, fmt.Errorf("failed to look up validator set verifying warp message: %w", err)
	// }
	// subnetLookupGasCost, overflow := math.SafeMul(uint64(len(validatorSet)), GasCostPerSourceSubnetValidator)
	// if overflow {
	// 	return 0, fmt.Errorf("overflow calculating gas cost for subnet (%s) validator set lookup of size %d", subnetID, len(validatorSet))
	// }
	// totalGas, overflow = math.SafeAdd(totalGas, subnetLookupGasCost)
	// if overflow {
	// 	return 0, fmt.Errorf("overflow adding subnet lookup gas (PrevTotal: %d, SubnetLookupGas: %d)", totalGas, subnetLookupGasCost)
	// }

	return totalGas, nil
}

// VerifyPredicate verifies the predicate represents a valid signed and properly formatted Avalanche Warp Message.
func (c *Config) VerifyPredicate(predicateContext *precompileconfig.ProposerPredicateContext, predicateBytes []byte) error {
	if predicateContext.ProposerVMBlockCtx == nil {
<<<<<<< HEAD
		return errNoProposerPredicate
=======
		return errNoProposerCtxPredicate
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	}
	// Note: PredicateGas should be called before VerifyPredicate, so we should never reach an error case here.
	unpackedPredicateBytes, err := predicateutils.UnpackPredicate(predicateBytes)
	if err != nil {
		return err
	}

	// Note: PredicateGas should be called before VerifyPredicate, so we should never reach an error case here.
	warpMessage, err := warp.ParseMessage(unpackedPredicateBytes)
	if err != nil {
		return fmt.Errorf("%w: %s", errInvalidWarpMsg, err)
	}
	return c.verifyWarpMessage(predicateContext, warpMessage)
}
