// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package payload

import (
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
)

// AddressedPayload defines the format for delivering a point to point message across VMs
// ie. (ChainA, AddressA) -> (ChainB, AddressB)
type AddressedPayload struct {
	SourceAddress      ids.ID `serialize:"true"`
	DestinationChainID ids.ID `serialize:"true"`
	DestinationAddress ids.ID `serialize:"true"`
	Payload            []byte `serialize:"true"`

	bytes []byte
}

// NewAddressedPayload creates a new *AddressedPayload and initializes it.
func NewAddressedPayload(sourceAddress ids.ID, destinationChainID ids.ID, destinationAddress ids.ID, payload []byte) (*AddressedPayload, error) {
	ap := &AddressedPayload{
		SourceAddress:      sourceAddress,
		DestinationChainID: destinationChainID,
		DestinationAddress: destinationAddress,
		Payload:            payload,
	}
	return ap, ap.initialize()
}

// ParseAddressedPayload converts a slice of bytes into an initialized
// AddressedPayload.
func ParseAddressedPayload(b []byte) (*AddressedPayload, error) {
	var unmarshalledPayloadIntf any
	if _, err := c.Unmarshal(b, &unmarshalledPayloadIntf); err != nil {
		return nil, err
	}
	payload, ok := unmarshalledPayloadIntf.(*AddressedPayload)
	if !ok {
		return nil, fmt.Errorf("failed to parse unexpected type %T as addressed payload", unmarshalledPayloadIntf)
	}
	payload.bytes = b
	return payload, nil
}

// initialize recalculates the result of Bytes().
func (a *AddressedPayload) initialize() error {
	aIntf := any(a)
	bytes, err := c.Marshal(codecVersion, &aIntf)
	if err != nil {
		return fmt.Errorf("couldn't marshal warp addressed payload: %w", err)
	}
	a.bytes = bytes
	return nil
}

// Bytes returns the binary representation of this payload. It assumes that the
// payload is initialized from either NewAddressedPayload or ParseAddressedPayload.
func (a *AddressedPayload) Bytes() []byte {
	return a.bytes
}
