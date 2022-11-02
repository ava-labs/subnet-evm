package utils

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type rpcChainCompatibility struct {
	RPCChainVMProtocolVersion map[string]int `json:"rpcChainVMProtocolVersion"`
}

func TestCompatibility(t *testing.T) {
	subevmVersion := "v0.4.0"
	expectedRPCVersion := 17

	compat, err := os.ReadFile("../compatibility.json")
	assert.NoError(t, err)

	var parsedCompat rpcChainCompatibility
	err = json.Unmarshal(compat, &parsedCompat)
	assert.NoError(t, err)
	assert.Equal(t, expectedRPCVersion, parsedCompat.RPCChainVMProtocolVersion[subevmVersion])
}
