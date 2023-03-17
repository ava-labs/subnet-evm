// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package messages

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
	return ap, ap.Initialize()
}

// ParseAddressedPayload converts a slice of bytes into an initialized
// *AddressedPayload.
func ParseAddressedPayload(b []byte) (*AddressedPayload, error) {
	payload := new(AddressedPayload)
	if _, err := c.Unmarshal(b, payload); err != nil {
		return nil, err
	}
	payload.bytes = b
	return payload, nil
}

// Initialize recalculates the result of Bytes().
func (a *AddressedPayload) Initialize() error {
	bytes, err := c.Marshal(codecVersion, a)
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
