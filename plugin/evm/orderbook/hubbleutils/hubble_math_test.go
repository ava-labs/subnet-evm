package hubbleutils

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestECRecovers(t *testing.T) {
	// 1. Test case from
	orderHash := "0xee4b26ae386d1c88f89eb2f8b4b4205271576742f5ff4e0488633612f7a9a5e7"
	address, err := ECRecover(common.FromHex(orderHash), common.FromHex("0xb2704b73b99f2700ecc90a218f514c254d1f5d46af47117f5317f6cc0348ce962dcfb024c7264fdeb1f1513e4564c2a7cd9c1d0be33d7b934cd5a73b96440eaf1c"))
	assert.Nil(t, err)
	assert.Equal(t, "0x70997970C51812dc3A010C7d01b50e0d17dc79C8", address.String())
}
