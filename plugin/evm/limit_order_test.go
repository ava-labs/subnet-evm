package evm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetOrderBookContractFileLocation(t *testing.T) {
	newFileLocation := "new/location"
	SetOrderBookContractFileLocation(newFileLocation)
	assert.Equal(t, newFileLocation, orderBookContractFileLocation)
}

func TestNewLimitOrderProcesser(t *testing.T) {
	txFeeCap := float64(11)
	enabledEthAPIs := []string{"debug"}
	configJSON := fmt.Sprintf("{\"rpc-tx-fee-cap\": %g,\"eth-apis\": %s}", txFeeCap, fmt.Sprintf("[%q]", enabledEthAPIs[0]))
	_, vm, _, _ := GenesisVM(t, false, "", configJSON, "")
	lop := NewLimitOrderProcesser(
		vm.ctx,
		vm.chainConfig,
		vm.txPool,
		vm.shutdownChan,
		&vm.shutdownWg,
		vm.eth.APIBackend,
		vm.eth.BlockChain(),
	)
	assert.NotNil(t, lop)
}
