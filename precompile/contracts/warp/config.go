// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	warpMessages "github.com/ava-labs/subnet-evm/warp/messages"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/log"
)

const (
	QuorumDenominator      uint64 = 100
	DefaultQuorumNumerator uint64 = 67
	MinQuorumNumerator     uint64 = 33
)

var (
	_ precompileconfig.Config             = &Config{}
	_ precompileconfig.ProposerPredicater = &Config{}
	_ precompileconfig.Accepter           = &Config{}
)

var errOverflowSignersGasCost = errors.New("overflow calculating warp signers gas cost")

// Config implements the precompileconfig.Config interface and
// adds specific configuration for Warp.
type Config struct {
	precompileconfig.Upgrade
	QuorumNumerator uint64 `json:"quorumNumerator"`
}

// NewConfig returns a config for a network upgrade at [blockTimestamp] that enables
// Warp with the given quorum numerator.
func NewConfig(blockTimestamp *big.Int, quorumNumerator uint64) *Config {
	return &Config{
		Upgrade:         precompileconfig.Upgrade{BlockTimestamp: blockTimestamp},
		QuorumNumerator: quorumNumerator,
	}
}

// NewDefaultConfig returns a config for a network upgrade at [blockTimestamp] that enables
// Warp with the default quorum numerator (0 denotes using the default).
func NewDefaultConfig(blockTimestamp *big.Int) *Config {
	return NewConfig(blockTimestamp, 0)
}

// NewDisableConfig returns config for a network upgrade at [blockTimestamp]
// that disables Warp.
func NewDisableConfig(blockTimestamp *big.Int) *Config {
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
func (c *Config) Verify() error {
	if c.QuorumNumerator > QuorumDenominator {
		return fmt.Errorf("cannot specify quorum numerator (%d) > quorum denominator (%d)", c.QuorumNumerator, QuorumDenominator)
	}
	// If a non-default quorum numerator is specified and it is less than the minimum, return an error
	if c.QuorumNumerator != 0 && c.QuorumNumerator < MinQuorumNumerator {
		return fmt.Errorf("cannot specify quorum numerator (%d) < min quorum numerator (%d)", c.QuorumNumerator, MinQuorumNumerator)
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

func (c *Config) verifyWarpMessage(predicateContext *precompileconfig.ProposerPredicateContext, warpMsg *warp.Message) (uint64, error) {
	numSigners, err := warpMsg.Signature.NumSigners()
	if err != nil {
		return 0, fmt.Errorf("failed to get num signers from warp message: %w", err)
	}
	msgGas, overflow := math.SafeMul(uint64(numSigners), GasCostPerWarpSigner)
	if overflow {
		return 0, errOverflowSignersGasCost
	}

	// Use default quourum numerator unless config specifies a non-default option
	quorumNumerator := DefaultQuorumNumerator
	if c.QuorumNumerator != 0 {
		quorumNumerator = c.QuorumNumerator
	}

	// Verify the warp payload can be decoded to the expected type
	_, err = warpMessages.ParseAddressedPayload(warpMsg.UnsignedMessage.Payload)
	if err != nil {
		return 0, fmt.Errorf("failed to parse warp payload into addressed payload: %w", err)
	}

	log.Info("verifyingWarpMessage", "warpMsg", warpMsg, "quorumNum", quorumNumerator, "quorumDenom", QuorumDenominator)
	if err := warpMsg.Signature.Verify(
		context.Background(),
		&warpMsg.UnsignedMessage,
		predicateContext.SnowCtx.ValidatorState, // TODO: special case messages from the C-Chain
		predicateContext.ProposerVMBlockCtx.PChainHeight,
		quorumNumerator,
		QuorumDenominator,
	); err != nil {
		return 0, fmt.Errorf("warp signature verification failed: %w", err)
	}

	return msgGas, nil
}

// TODO: move to general package, cleanup, and test
var predicateEndByte = byte(0xff)

func PackPredicate(predicate []byte) []byte {
	predicate = append(predicate, predicateEndByte)
	return common.RightPadBytes(predicate, (len(predicate)+31/32)*32)
}

func UnpackPredicate(paddedPredicate []byte) ([]byte, error) {
	trimmedPredicateBytes := common.TrimRightZeroes(paddedPredicate)
	if len(trimmedPredicateBytes) == 0 {
		return nil, fmt.Errorf("warp predicate specified invalid all zero bytes: 0x%x", paddedPredicate)
	}

	if trimmedPredicateBytes[len(trimmedPredicateBytes)-1] != predicateEndByte {
		return nil, fmt.Errorf("invalid end delimiter")
	}

	return trimmedPredicateBytes[:len(trimmedPredicateBytes)-1], nil
}

func (c *Config) VerifyPredicate(predicateContext *precompileconfig.ProposerPredicateContext, predicateBytes []byte) error {
	if predicateContext.ProposerVMBlockCtx == nil {
		return fmt.Errorf("cannot specify a proposer predicate for %s in a block before ProposerVM activation", ConfigKey)
	}
	totalGas, overflow := math.SafeMul(GasCostPerWarpMessageBytes, uint64(len(predicateBytes)))
	if overflow {
		return fmt.Errorf("overflow calculating gas cost for warp message bytes of size %d", len(predicateBytes))
	}

	unpackedPredicateBytes, err := UnpackPredicate(predicateBytes)
	if err != nil {
		return err
	}
	warpMessage, err := warp.ParseMessage(unpackedPredicateBytes)
	if err != nil {
		return err
	}
	verificationGas, err := c.verifyWarpMessage(predicateContext, warpMessage)
	if err != nil {
		return err
	}
	totalGas, overflow = math.SafeAdd(totalGas, verificationGas)
	if overflow {
		return fmt.Errorf("overflow adding verification gas (PrevTotal: %d, VerificationGas: %d)", totalGas, verificationGas)
	}
	return nil
}
