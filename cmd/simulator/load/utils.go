// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"fmt"
	"strings"

	"github.com/ava-labs/subnet-evm/tests/utils/runner"
)

func toWebsocketURIs(subnet *runner.Subnet) []string {
	nodeURIs := subnet.ValidatorURIs
	wsEndpoints := make([]string, len(nodeURIs))
	for i, uri := range nodeURIs {
		wsEndpoints[i] = fmt.Sprintf(
			"ws://%s/ext/bc/%s/ws",
			strings.TrimPrefix(uri, "http://"),
			subnet.BlockchainID,
		)
	}
	return wsEndpoints
}
