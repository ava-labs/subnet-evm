// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"github.com/ava-labs/avalanchego/codec"
	"github.com/ava-labs/avalanchego/codec/linearcodec"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/utils/wrappers"
)

const (
	codecVersion   uint16 = 0
	maxMessageSize        = 512 * units.KiB
	maxSliceLen           = maxMessageSize
)

// Codec does serialization and deserialization
var Codec codec.Manager

// TODO: switch from using Codec to use ABI serialization
func init() {
	Codec = codec.NewManager(maxMessageSize)
	lc := linearcodec.NewCustomMaxLength(maxSliceLen)

	errs := wrappers.Errs{}
	errs.Add(
		lc.RegisterType(&WarpMessage{}),
		Codec.RegisterCodec(codecVersion, lc),
	)
	if errs.Errored() {
		panic(errs.Err)
	}
}
