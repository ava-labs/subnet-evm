package hubbleutils

import (
	"encoding/hex"
	"math/big"
	"strings"

	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestDecodeSignedOrder(t *testing.T) {
	t.Run("long order", func(t *testing.T) {
		SetChainIdAndVerifyingSignedOrdersContract(321123, "0x4c5859f0F772848b2D91F1D83E2Fe57935348029")
		orderHash := strings.TrimPrefix("0x73d5196ac9576efaccb6e54b193b894e2cc0afd68ce5af519c901fec7e588595", "0x")
		signature := strings.TrimPrefix("0x3027ae4ab98663490d0facab04c71665e41da867a44b7ddc29e14cb8de3a3cfa12985be54945ce040196b2fcdcc4dafc56f7955ee72628bc9e7a634a7f258ce61c", "0x")
		encodedOrder := strings.TrimPrefix("0x00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000064ac0426000000000000000000000000000000000000000000000000000000000000000000000000000000000000000070997970c51812dc3a010c7d01b50e0d17dc79c80000000000000000000000000000000000000000000000004563918244f40000000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000001893fef795900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000000413027ae4ab98663490d0facab04c71665e41da867a44b7ddc29e14cb8de3a3cfa12985be54945ce040196b2fcdcc4dafc56f7955ee72628bc9e7a634a7f258ce61c00000000000000000000000000000000000000000000000000000000000000", "0x")
		typeEncodedOrder := strings.TrimPrefix("0x0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001c000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000064ac0426000000000000000000000000000000000000000000000000000000000000000000000000000000000000000070997970c51812dc3a010c7d01b50e0d17dc79c80000000000000000000000000000000000000000000000004563918244f40000000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000001893fef795900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000000413027ae4ab98663490d0facab04c71665e41da867a44b7ddc29e14cb8de3a3cfa12985be54945ce040196b2fcdcc4dafc56f7955ee72628bc9e7a634a7f258ce61c00000000000000000000000000000000000000000000000000000000000000", "0x")

		sig, err := hex.DecodeString(signature)
		assert.Nil(t, err)
		order := &SignedOrder{
			LimitOrder: LimitOrder{
				BaseOrder: BaseOrder{
					AmmIndex:          big.NewInt(0),
					Trader:            common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8"),
					BaseAssetQuantity: big.NewInt(5000000000000000000),
					Price:             big.NewInt(1000000000),
					Salt:              big.NewInt(1688994806105),
					ReduceOnly:        false,
				},
				PostOnly: true,
			},
			OrderType: 2,
			ExpireAt:  big.NewInt(1688994854),
			Sig:       sig,
		}
		h, err := order.Hash()
		assert.Nil(t, err)
		assert.Equal(t, orderHash, strings.TrimPrefix(h.Hex(), "0x"))

		b, err := order.EncodeToABIWithoutType()
		assert.Nil(t, err)
		assert.Equal(t, encodedOrder, hex.EncodeToString(b))

		b, err = order.EncodeToABI()
		assert.Nil(t, err)
		assert.Equal(t, typeEncodedOrder, hex.EncodeToString(b))

		testDecodeTypeAndEncodedSignedOrder(t, typeEncodedOrder, encodedOrder, Signed, order)

		data, err := hex.DecodeString(orderHash)
		assert.Nil(t, err)
		signer, err := ECRecover(data, sig)
		assert.Nil(t, err)
		assert.Equal(t, order.Trader, signer)

		sig_, _ := hex.DecodeString(signature)
		assert.Equal(t, sig_, sig)       // sig is not changed
		assert.Equal(t, sig_, order.Sig) // sig is not changed
	})

	t.Run("short order", func(t *testing.T) {
		SetChainIdAndVerifyingSignedOrdersContract(321123, "0x809d550fca64d94Bd9F66E60752A544199cfAC3D")
		orderHash := strings.TrimPrefix("0xee4b26ae386d1c88f89eb2f8b4b4205271576742f5ff4e0488633612f7a9a5e7", "0x")
		signature := strings.TrimPrefix("0xb2704b73b99f2700ecc90a218f514c254d1f5d46af47117f5317f6cc0348ce962dcfb024c7264fdeb1f1513e4564c2a7cd9c1d0be33d7b934cd5a73b96440eaf1c", "0x")
		encodedOrder := strings.TrimPrefix("0x00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000064ac0426000000000000000000000000000000000000000000000000000000000000000000000000000000000000000070997970c51812dc3a010c7d01b50e0d17dc79c8ffffffffffffffffffffffffffffffffffffffffffffffffba9c6e7dbb0c0000000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000001893fef79590000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000041b2704b73b99f2700ecc90a218f514c254d1f5d46af47117f5317f6cc0348ce962dcfb024c7264fdeb1f1513e4564c2a7cd9c1d0be33d7b934cd5a73b96440eaf1c00000000000000000000000000000000000000000000000000000000000000", "0x")
		typeEncodedOrder := strings.TrimPrefix("0x0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001c000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000064ac0426000000000000000000000000000000000000000000000000000000000000000000000000000000000000000070997970c51812dc3a010c7d01b50e0d17dc79c8ffffffffffffffffffffffffffffffffffffffffffffffffba9c6e7dbb0c0000000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000001893fef79590000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000041b2704b73b99f2700ecc90a218f514c254d1f5d46af47117f5317f6cc0348ce962dcfb024c7264fdeb1f1513e4564c2a7cd9c1d0be33d7b934cd5a73b96440eaf1c00000000000000000000000000000000000000000000000000000000000000", "0x")

		sig, err := hex.DecodeString(signature)
		assert.Nil(t, err)
		order := &SignedOrder{
			LimitOrder: LimitOrder{
				BaseOrder: BaseOrder{
					AmmIndex:          big.NewInt(0),
					Trader:            common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8"),
					BaseAssetQuantity: big.NewInt(-5000000000000000000),
					Price:             big.NewInt(1000000000),
					Salt:              big.NewInt(1688994806105),
					ReduceOnly:        false,
				},
				PostOnly: true,
			},
			OrderType: 2,
			ExpireAt:  big.NewInt(1688994854),
			Sig:       sig,
		}
		h, err := order.Hash()
		assert.Nil(t, err)
		assert.Equal(t, orderHash, strings.TrimPrefix(h.Hex(), "0x"))

		b, err := order.EncodeToABIWithoutType()
		assert.Nil(t, err)
		assert.Equal(t, encodedOrder, hex.EncodeToString(b))

		b, err = order.EncodeToABI()
		assert.Nil(t, err)
		assert.Equal(t, typeEncodedOrder, hex.EncodeToString(b))

		testDecodeTypeAndEncodedSignedOrder(t, typeEncodedOrder, encodedOrder, Signed, order)

		data, err := hex.DecodeString(orderHash)
		assert.Nil(t, err)
		signer, err := ECRecover(data, sig)
		assert.Nil(t, err)
		assert.Equal(t, order.Trader, signer)

		sig_, _ := hex.DecodeString(signature)
		assert.Equal(t, sig_, sig)       // sig is not changed
		assert.Equal(t, sig_, order.Sig) // sig is not changed
	})
}

func testDecodeTypeAndEncodedSignedOrder(t *testing.T, typedEncodedOrder string, encodedOrder string, orderType OrderType, expectedOutput *SignedOrder) {
	testData, err := hex.DecodeString(typedEncodedOrder)
	assert.Nil(t, err)

	decodeStep, err := DecodeTypeAndEncodedOrder(testData)
	assert.Nil(t, err)

	assert.Equal(t, orderType, decodeStep.OrderType)
	assert.Equal(t, encodedOrder, hex.EncodeToString(decodeStep.EncodedOrder))
	assert.Nil(t, err)
	testDecodeSignedOrder(t, decodeStep.EncodedOrder, expectedOutput)
}

func testDecodeSignedOrder(t *testing.T, encodedOrder []byte, expectedOutput *SignedOrder) {
	result, err := DecodeSignedOrder(encodedOrder)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assertSignedOrderEquality(t, expectedOutput, result)
}

func assertSignedOrderEquality(t *testing.T, expected, actual *SignedOrder) {
	assert.Equal(t, expected.OrderType, actual.OrderType)
	assert.Equal(t, expected.ExpireAt.Int64(), actual.ExpireAt.Int64())
	assert.Equal(t, expected.Sig, actual.Sig)
	assertLimitOrderEquality(t, expected.BaseOrder, actual.BaseOrder)
}

func assertLimitOrderEquality(t *testing.T, expected, actual BaseOrder) {
	assert.Equal(t, expected.AmmIndex.Int64(), actual.AmmIndex.Int64())
	assert.Equal(t, expected.Trader, actual.Trader)
	assert.Equal(t, expected.BaseAssetQuantity, actual.BaseAssetQuantity)
	assert.Equal(t, expected.Price, actual.Price)
	assert.Equal(t, expected.Salt, actual.Salt)
	assert.Equal(t, expected.ReduceOnly, actual.ReduceOnly)
}
