// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package messages

// Payload provides a common interface for all payloads implemented by this
// package.
type Payload interface {
	// Bytes returns the binary representation of this payload.
	Bytes() []byte

	// initialize the payload with the provided binary representation.
	initialize(b []byte)
}

func Parse(bytes []byte) (Payload, error) {
	var payload Payload
	if _, err := Codec.Unmarshal(bytes, &payload); err != nil {
		return nil, err
	}
	payload.initialize(bytes)
	return payload, nil
}
