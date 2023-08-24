// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package payload

import (
	"fmt"

<<<<<<< HEAD
	"github.com/ava-labs/avalanchego/ids"
)

// AddressedPayload defines an optional format for the bytes payload of a Warp message.
type AddressedPayload struct {
	SourceAddress      ids.ID `serialize:"true"`
	DestinationAddress ids.ID `serialize:"true"`
	Payload            []byte `serialize:"true"`
=======
	"github.com/ethereum/go-ethereum/common"
)

// AddressedPayload defines the format for delivering a point to point message across VMs
// ie. (ChainA, AddressA) -> (ChainB, AddressB)
type AddressedPayload struct {
	SourceAddress      common.Address `serialize:"true"`
	DestinationChainID common.Hash    `serialize:"true"`
	DestinationAddress common.Address `serialize:"true"`
	Payload            []byte         `serialize:"true"`
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d

	bytes []byte
}

// NewAddressedPayload creates a new *AddressedPayload and initializes it.
<<<<<<< HEAD
func NewAddressedPayload(sourceAddress ids.ID, destinationAddress ids.ID, payload []byte) (*AddressedPayload, error) {
	ap := &AddressedPayload{
		SourceAddress:      sourceAddress,
		DestinationAddress: destinationAddress,
		Payload:            payload,
	}
	return ap, ap.Initialize()
=======
func NewAddressedPayload(sourceAddress common.Address, destinationChainID common.Hash, destinationAddress common.Address, payload []byte) (*AddressedPayload, error) {
	ap := &AddressedPayload{
		SourceAddress:      sourceAddress,
		DestinationChainID: destinationChainID,
		DestinationAddress: destinationAddress,
		Payload:            payload,
	}
	return ap, ap.initialize()
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
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

<<<<<<< HEAD
// Initialize recalculates the result of Bytes().
func (a *AddressedPayload) Initialize() error {
=======
// initialize recalculates the result of Bytes().
func (a *AddressedPayload) initialize() error {
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	aIntf := any(a)
	bytes, err := c.Marshal(codecVersion, &aIntf)
	if err != nil {
		return fmt.Errorf("couldn't marshal warp addressed payload: %w", err)
	}
	a.bytes = bytes
	return nil
}

// Bytes returns the binary representation of this payload. It assumes that the
<<<<<<< HEAD
// payload is initialized from either NewAddressedPayload, ParseAddressedPayload, or an explicit call to
// Initialize.
=======
// payload is initialized from either NewAddressedPayload or ParseAddressedPayload.
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
func (a *AddressedPayload) Bytes() []byte {
	return a.bytes
}
