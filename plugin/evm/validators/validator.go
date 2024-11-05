// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package validators

import (
	"time"

	ids "github.com/ava-labs/avalanchego/ids"
)

type Validator struct {
	ValidationID   ids.ID     `json:"validationID"`
	NodeID         ids.NodeID `json:"nodeID"`
	Weight         uint64     `json:"weight"`
	StartTimestamp uint64     `json:"startTimestamp"`
	IsActive       bool       `json:"isActive"`
	IsSoV          bool       `json:"isSoV"`
}

func (v *Validator) StartTime() time.Time { return time.Unix(int64(v.StartTimestamp), 0) }
