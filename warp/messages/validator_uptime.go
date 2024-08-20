// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package messages

import (
	"errors"
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ethereum/go-ethereum/log"
)

var errInvalidRequest = errors.New("invalid request")

// ValidatorUptime is signed when the ValidationID is known and the validator
// has been up for TotalUptime seconds.
type ValidatorUptime struct {
	ValidationID ids.ID `serialize:"true"`
	TotalUptime  uint64 `serialize:"true"`

	bytes []byte
}

// NewValidatorUptime creates a new *ValidatorUptime and initializes it.
func NewValidatorUptime(validationID ids.ID, totalUptime uint64) (*ValidatorUptime, error) {
	bhp := &ValidatorUptime{
		ValidationID: validationID,
		TotalUptime:  totalUptime,
	}
	return bhp, initialize(bhp)
}

// ParseValidatorUptime converts a slice of bytes into an initialized ValidatorUptime.
func ParseValidatorUptime(b []byte) (*ValidatorUptime, error) {
	payloadIntf, err := Parse(b)
	if err != nil {
		return nil, err
	}
	payload, ok := payloadIntf.(*ValidatorUptime)
	if !ok {
		return nil, fmt.Errorf("%w: %T", errWrongType, payloadIntf)
	}
	return payload, nil
}

// Bytes returns the binary representation of this payload. It assumes that the
// payload is initialized from either NewValidatorUptime or Parse.
func (b *ValidatorUptime) Bytes() []byte {
	return b.bytes
}

func (b *ValidatorUptime) initialize(bytes []byte) {
	b.bytes = bytes
}

// VerifyMesssage returns nil if the message is valid.
func (b *ValidatorUptime) VerifyMesssage(sourceAddress []byte) error {
	// DO NOT USE THIS CODE AS IS IN PRODUCTION.

	// TODO: Does nil/empty SourceAddress matter?
	if len(sourceAddress) != 0 {
		return errInvalidRequest
	}

	log.Info("Received validator uptime message", "validationID", b.ValidationID, "totalUptime", b.TotalUptime)
	log.Warn("Signing validator uptime message by default, not production behavior", "validationID", b.ValidationID, "totalUptime", b.TotalUptime)

	return nil
}
