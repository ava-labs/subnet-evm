// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements solidity tests.
package solidity

import (
	"context"

	"github.com/ava-labs/avalanchego/api/health"
	"github.com/ava-labs/subnet-evm/tests/utils"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("[Warp]", ginkgo.Ordered, func() {
	ginkgo.It("Setup Subnets", ginkgo.Label("Warp"), func() {
		client := health.NewClient(utils.DefaultLocalNodeURI)
		healthy, err := client.Readiness(context.Background())
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(healthy.Healthy).Should(gomega.BeTrue())
	})

	// Steps:
	// Set up two subnets with the same genesis
	// send a message from blockchainA to blockchainB
	// verify message on blockchainB from blockchainA
})
