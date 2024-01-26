// (c) 2022 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompileconfig

import (
	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/docker/docker/pkg/units"
)

// Upgrade contains the timestamp for the upgrade along with
// a boolean [Disable]. If [Disable] is set, the upgrade deactivates
// the precompile and clears its storage.
type Upgrade struct {
	BlockTimestamp *uint64 `json:"blockTimestamp"`
	Disable        bool    `json:"disable,omitempty"`
}

// Timestamp returns the timestamp this network upgrade goes into effect.
func (u *Upgrade) Timestamp() *uint64 {
	return u.BlockTimestamp
}

// IsDisabled returns true if the network upgrade deactivates the precompile.
func (u *Upgrade) IsDisabled() bool {
	return u.Disable
}

// Equal returns true iff [other] has the same blockTimestamp and has the
// same on value for the Disable flag.
func (u *Upgrade) Equal(other *Upgrade) bool {
	if other == nil {
		return false
	}
	return u.Disable == other.Disable && utils.Uint64PtrEqual(u.BlockTimestamp, other.BlockTimestamp)
}

func (u *Upgrade) MarshalBinary() ([]byte, error) {
	p := wrappers.Packer{
		Bytes:   []byte{},
		MaxSize: 1 * units.MiB,
	}
	if u.BlockTimestamp == nil {
		p.PackBool(true)
	} else {
		p.PackBool(false)
		if p.Err != nil {
			return nil, p.Err
		}
		p.PackLong(*u.BlockTimestamp)
	}
	if p.Err != nil {
		return nil, p.Err
	}
	p.PackBool(u.Disable)
	return p.Bytes, p.Err
}

func (u *Upgrade) UnmarshalBinary(data []byte) error {
	p := wrappers.Packer{
		Bytes: data,
	}
	isNil := p.UnpackBool()
	if !isNil {
		timestamp := p.UnpackLong()
		u.BlockTimestamp = &timestamp
	}
	u.Disable = p.UnpackBool()
	return nil
}
