package evm

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type rpcChainCompatibility struct {
	RPCChainVMProtocolVersion map[string]int `json:"rpcChainVMProtocolVersion"`
}

const compatibilityFile = "../../compatibility.json"

func TestCompatibility(t *testing.T) {
	subevmVersion := "v0.4.0"
	expectedRPCVersion := 17

	compat, err := os.ReadFile(compatibilityFile)
	assert.NoError(t, err)

	var parsedCompat rpcChainCompatibility
	err = json.Unmarshal(compat, &parsedCompat)
	assert.NoError(t, err)
	assert.Equal(t, expectedRPCVersion, parsedCompat.RPCChainVMProtocolVersion[subevmVersion])
}

func TestCompatibilityCurrentVersion(t *testing.T) {
	compat, err := os.ReadFile(compatibilityFile)
	assert.NoError(t, err)

	var parsedCompat rpcChainCompatibility
	err = json.Unmarshal(compat, &parsedCompat)
	assert.NoError(t, err)

	_, valueInJSON := parsedCompat.RPCChainVMProtocolVersion[Version]
	assert.True(t, valueInJSON)
}
