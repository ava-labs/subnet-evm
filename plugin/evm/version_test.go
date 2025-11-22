// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type rpcChainCompatibility struct {
	RPCChainVMProtocolVersion map[string]uint `json:"rpcChainVMProtocolVersion"`
}

const compatibilityFile = "../../compatibility.json"

func TestCompatibility(t *testing.T) {
	compat, err := os.ReadFile(compatibilityFile)
	require.NoError(t, err, "reading compatibility file")

	var parsedCompat rpcChainCompatibility
	err = json.Unmarshal(compat, &parsedCompat)
	require.NoError(t, err, "json decoding compatibility file")
	require.Contains(t, parsedCompat.RPCChainVMProtocolVersion, Version, "subnet-evm version %s missing from rpcChainVMProtocolVersion object", Version)
}
