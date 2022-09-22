// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	ginkgo "github.com/onsi/ginkgo/v2"
)

// DescribeLocal annotates the tests that requires local network-runner.
// Can only run with local cluster.
func DescribeLocal(text string, body func()) bool {
	return ginkgo.Describe("[Local] "+text, body)
}

func DescribePrecompile(body func()) bool {
	return ginkgo.Describe("[Precompiles]", ginkgo.Ordered, body)
}
