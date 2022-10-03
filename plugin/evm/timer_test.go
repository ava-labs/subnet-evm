// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"sync"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()

	timer := NewTimer(wg.Done)
	go timer.Dispatch()

	timer.SetTimeoutIn(time.Millisecond)
}
