// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package payload

import (
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
)

// AddressedPayload defines an optional format for the bytes payload of a Warp message.
type AddressedPayload struct {
	SourceAddress      ids.ID `serialize:"true"`
	DestinationAddress ids.ID `serialize:"true"`
	Payload            []byte `serialize:"true"`

	bytes []byte
}

// NewAddressedPayload creates a new *AddressedPayload and initializes it.
func NewAddressedPayload(sourceAddress ids.ID, destinationAddress ids.ID, payload []byte) (*AddressedPayload, error) {
	ap := &AddressedPayload{
		SourceAddress:      sourceAddress,
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

// Initialize recalculates the result of Bytes().
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
// payload is initialized from either NewAddressedPayload, ParseAddressedPayload, or an explicit call to
// Initialize.
func (a *AddressedPayload) Bytes() []byte {
	return a.bytes
}
