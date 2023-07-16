package orderbook

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/eth"
	"github.com/ava-labs/subnet-evm/metrics"
	"github.com/ava-labs/subnet-evm/plugin/evm/orderbook/abis"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

var OrderBookContractAddress = common.HexToAddress("0x0300000000000000000000000000000000000000")
var MarginAccountContractAddress = common.HexToAddress("0x0300000000000000000000000000000000000001")
var ClearingHouseContractAddress = common.HexToAddress("0x0300000000000000000000000000000000000002")
var IOCOrderBookContractAddress = common.HexToAddress("0x635c5F96989a4226953FE6361f12B96c5d50289b")

// var IOCOrderBookContractAddress = common.HexToAddress("0x0300000000000000000000000000000000000006")

type LimitOrderTxProcessor interface {
	PurgeLocalTx()
	CheckIfOrderBookContractCall(tx *types.Transaction) bool
	ExecuteMatchedOrdersTx(incomingOrder Order, matchedOrder Order, fillAmount *big.Int) error
	ExecuteFundingPaymentTx() error
	ExecuteLiquidation(trader common.Address, matchedOrder Order, fillAmount *big.Int) error
	UpdateMetrics(block *types.Block)
	ExecuteLimitOrderCancel(orderIds []LimitOrder) error
}

type ValidatorTxFeeConfig struct {
	baseFeeEstimate *big.Int
	blockNumber     uint64
}

type limitOrderTxProcessor struct {
	txPool                       *core.TxPool
	memoryDb                     LimitOrderDatabase
	orderBookABI                 abi.ABI
	clearingHouseABI             abi.ABI
	marginAccountABI             abi.ABI
	orderBookContractAddress     common.Address
	clearingHouseContractAddress common.Address
	marginAccountContractAddress common.Address
	backend                      *eth.EthAPIBackend
	validatorAddress             common.Address
	validatorPrivateKey          string
	validatorTxFeeConfig         ValidatorTxFeeConfig
}

func NewLimitOrderTxProcessor(txPool *core.TxPool, memoryDb LimitOrderDatabase, backend *eth.EthAPIBackend, validatorPrivateKey string) LimitOrderTxProcessor {
	orderBookABI, err := abi.FromSolidityJson(string(abis.OrderBookAbi))
	if err != nil {
		panic(err)
	}

	clearingHouseABI, err := abi.FromSolidityJson(string(abis.ClearingHouseAbi))
	if err != nil {
		panic(err)
	}

	marginAccountABI, err := abi.FromSolidityJson(string(abis.MarginAccountAbi))
	if err != nil {
		panic(err)
	}

	if validatorPrivateKey == "" {
		panic("private key is not supplied")
	}
	validatorAddress, err := getAddressFromPrivateKey(validatorPrivateKey)
	if err != nil {
		panic("Unable to get address from validator private key")
	}

	lotp := &limitOrderTxProcessor{
		txPool:                       txPool,
		orderBookABI:                 orderBookABI,
		clearingHouseABI:             clearingHouseABI,
		marginAccountABI:             marginAccountABI,
		memoryDb:                     memoryDb,
		orderBookContractAddress:     OrderBookContractAddress,
		clearingHouseContractAddress: ClearingHouseContractAddress,
		marginAccountContractAddress: MarginAccountContractAddress,
		backend:                      backend,
		validatorAddress:             validatorAddress,
		validatorPrivateKey:          validatorPrivateKey,
		validatorTxFeeConfig:         ValidatorTxFeeConfig{baseFeeEstimate: big.NewInt(0), blockNumber: 0},
	}
	lotp.updateValidatorTxFeeConfig()
	return lotp
}

func (lotp *limitOrderTxProcessor) ExecuteLiquidation(trader common.Address, matchedOrder Order, fillAmount *big.Int) error {
	orderBytes, err := matchedOrder.RawOrder.EncodeToABI()
	if err != nil {
		log.Error("EncodeLimitOrder failed in ExecuteLiquidation", "order", matchedOrder, "err", err)
		return err
	}
	txHash, err := lotp.executeLocalTx(lotp.orderBookContractAddress, lotp.orderBookABI, "liquidateAndExecuteOrder", trader, orderBytes, fillAmount)
	log.Info("ExecuteLiquidation", "trader", trader, "matchedOrder", matchedOrder, "fillAmount", prettifyScaledBigInt(fillAmount, 18), "txHash", txHash.String(), "err", err)
	// log.Info("ExecuteLiquidation", "trader", trader, "matchedOrder", matchedOrder, "fillAmount", prettifyScaledBigInt(fillAmount, 18), "orderBytes", hex.EncodeToString(orderBytes), "txHash", txHash.String(), "err", err)
	return err
}

func (lotp *limitOrderTxProcessor) ExecuteFundingPaymentTx() error {
	txHash, err := lotp.executeLocalTx(lotp.orderBookContractAddress, lotp.orderBookABI, "settleFunding")
	log.Info("ExecuteFundingPaymentTx", "txHash", txHash.String(), "err", err)
	return err
}

func (lotp *limitOrderTxProcessor) ExecuteMatchedOrdersTx(longOrder Order, shortOrder Order, fillAmount *big.Int) error {
	var err error
	orders := make([][]byte, 2)
	orders[0], err = longOrder.RawOrder.EncodeToABI()
	if err != nil {
		log.Error("EncodeLimitOrder failed for longOrder", "order", longOrder, "err", err)
		return err
	}
	orders[1], err = shortOrder.RawOrder.EncodeToABI()
	if err != nil {
		log.Error("EncodeLimitOrder failed for shortOrder", "order", shortOrder, "err", err)
		return err
	}

	txHash, err := lotp.executeLocalTx(lotp.orderBookContractAddress, lotp.orderBookABI, "executeMatchedOrders", orders, fillAmount)
	log.Info("ExecuteMatchedOrdersTx", "LongOrder", longOrder, "ShortOrder", shortOrder, "fillAmount", prettifyScaledBigInt(fillAmount, 18), "txHash", txHash.String(), "err", err)
	return err
}

func (lotp *limitOrderTxProcessor) ExecuteLimitOrderCancel(orders []LimitOrder) error {
	txHash, err := lotp.executeLocalTx(lotp.orderBookContractAddress, lotp.orderBookABI, "cancelOrders", orders)
	log.Info("ExecuteLimitOrderCancel", "orders", orders, "txHash", txHash.String(), "err", err)
	return err
}

func (lotp *limitOrderTxProcessor) executeLocalTx(contract common.Address, contractABI abi.ABI, method string, args ...interface{}) (common.Hash, error) {
	var txHash common.Hash
	lotp.updateValidatorTxFeeConfig()
	nonce := lotp.txPool.GetOrderBookTxNonce(common.HexToAddress(lotp.validatorAddress.Hex())) // admin address

	data, err := contractABI.Pack(method, args...)
	if err != nil {
		log.Error("abi.Pack failed", "method", method, "args", args, "err", err)
		return txHash, err
	}
	key, err := crypto.HexToECDSA(lotp.validatorPrivateKey) // admin private key
	if err != nil {
		log.Error("HexToECDSA failed", "err", err)
		return txHash, err
	}
	tx := types.NewTransaction(nonce, contract, big.NewInt(0), 1500000, lotp.validatorTxFeeConfig.baseFeeEstimate, data)
	signer := types.NewLondonSigner(lotp.backend.ChainConfig().ChainID)
	signedTx, err := types.SignTx(tx, signer, key)
	if err != nil {
		log.Error("types.SignTx failed", "err", err)
		return txHash, err
	}
	txHash = signedTx.Hash()
	err = lotp.txPool.AddOrderBookTx(signedTx)
	if err != nil {
		log.Error("lop.txPool.AddOrderBookTx failed", "err", err, "tx", signedTx.Hash().String(), "nonce", nonce)
		return txHash, err
	}
	// log.Info("executeLocalTx - AddOrderBookTx success", "tx", signedTx.Hash().String(), "nonce", nonce)

	return txHash, nil
}

func (lotp *limitOrderTxProcessor) getBaseFeeEstimate() *big.Int {
	baseFeeEstimate, err := lotp.backend.EstimateBaseFee(context.TODO())
	if err != nil {
		baseFeeEstimate = big.NewInt(0).Abs(lotp.backend.CurrentBlock().BaseFee())
		log.Error("Error in calculating updated bassFee, using last header's baseFee", "baseFeeEstimate", baseFeeEstimate)
	}
	return baseFeeEstimate
}

func (lotp *limitOrderTxProcessor) updateValidatorTxFeeConfig() {
	currentBlockNumber := lotp.backend.CurrentBlock().NumberU64()
	if lotp.validatorTxFeeConfig.blockNumber < currentBlockNumber {
		baseFeeEstimate := lotp.getBaseFeeEstimate()
		// log.Info("inside lotp updating txFeeConfig", "blockNumber", currentBlockNumber, "baseFeeEstimate", baseFeeEstimate)
		lotp.validatorTxFeeConfig.baseFeeEstimate = baseFeeEstimate
		lotp.validatorTxFeeConfig.blockNumber = currentBlockNumber
	}
}

func (lotp *limitOrderTxProcessor) PurgeLocalTx() {
	pending := lotp.txPool.Pending(true)
	for _, txs := range pending {
		for _, tx := range txs {
			method, err := getOrderBookContractCallMethod(tx, lotp.orderBookABI, lotp.orderBookContractAddress)
			if err == nil {
				if method.Name == "executeMatchedOrders" || method.Name == "settleFunding" || method.Name == "liquidateAndExecuteOrder" {
					lotp.txPool.RemoveTx(tx.Hash())
				}
			}
		}
	}
	lotp.txPool.PurgeOrderBookTxs()
}

func (lotp *limitOrderTxProcessor) CheckIfOrderBookContractCall(tx *types.Transaction) bool {
	return checkIfOrderBookContractCall(tx, lotp.orderBookABI, lotp.orderBookContractAddress)
}

func getPositionTypeBasedOnBaseAssetQuantity(baseAssetQuantity *big.Int) PositionType {
	if baseAssetQuantity.Sign() == 1 {
		return LONG
	}
	return SHORT
}

func checkIfOrderBookContractCall(tx *types.Transaction, orderBookABI abi.ABI, orderBookContractAddress common.Address) bool {
	input := tx.Data()
	if tx.To() != nil && tx.To().Hash() == orderBookContractAddress.Hash() && len(input) > 3 {
		return true
	}
	return false
}

func getOrderBookContractCallMethod(tx *types.Transaction, orderBookABI abi.ABI, orderBookContractAddress common.Address) (*abi.Method, error) {
	if checkIfOrderBookContractCall(tx, orderBookABI, orderBookContractAddress) {
		input := tx.Data()
		method := input[:4]
		m, err := orderBookABI.MethodById(method)
		return m, err
	} else {
		err := errors.New("tx is not an orderbook contract call")
		return nil, err
	}
}

func getAddressFromPrivateKey(key string) (common.Address, error) {
	privateKey, err := crypto.HexToECDSA(key) // admin private key
	if err != nil {
		return common.Address{}, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, errors.New("unable to get address from private key")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address, nil
}

func (lotp *limitOrderTxProcessor) UpdateMetrics(block *types.Block) {
	// defer func(start time.Time) { log.Info("limitOrderTxProcessor.UpdateMetrics", "time", time.Since(start)) }(time.Now())

	transactionsPerBlockHistogram.Update(int64(len(block.Transactions())))
	gasUsedPerBlockHistogram.Update(int64(block.GasUsed()))
	blockGasCostPerBlockHistogram.Update(block.BlockGasCost().Int64())

	ctx := context.Background()
	txs := block.Transactions()

	receipts, err := lotp.backend.GetReceipts(ctx, block.Hash())
	if err != nil {
		log.Error("UpdateMetrics - lotp.backend.GetReceipts failed", "err", err)
		return
	}

	bigblock := new(big.Int).SetUint64(block.NumberU64())
	timestamp := new(big.Int).SetUint64(block.Header().Time)
	signer := types.MakeSigner(lotp.backend.ChainConfig(), bigblock, timestamp)

	for i := 0; i < len(txs); i++ {
		tx := txs[i]
		receipt := receipts[i]
		from, _ := types.Sender(signer, tx)
		contractAddress := tx.To()
		input := tx.Data()
		if contractAddress == nil || len(input) < 4 {
			continue
		}
		method_ := input[:4]
		method, _ := lotp.orderBookABI.MethodById(method_)

		if method == nil {
			continue
		}

		if from == lotp.validatorAddress {
			if receipt.Status == 0 {
				orderBookTransactionsFailureTotalCounter.Inc(1)
			} else if receipt.Status == 1 {
				orderBookTransactionsSuccessTotalCounter.Inc(1)
			}

			if contractAddress != nil && lotp.orderBookContractAddress == *contractAddress {
				note := "success"
				if receipt.Status == 0 {
					note = "failure"
				}
				counterName := fmt.Sprintf("orderbooktxs/%s/%s", method.Name, note)
				metrics.GetOrRegisterCounter(counterName, nil).Inc(1)
			}

		}

		// measure the gas usage irrespective of whether the tx is from this validator or not
		if contractAddress != nil {
			var contractName string
			switch *contractAddress {
			case lotp.orderBookContractAddress:
				contractName = "OrderBook"
			case lotp.clearingHouseContractAddress:
				contractName = "ClearingHouse"
			case lotp.marginAccountContractAddress:
				contractName = "MarginAccount"
			default:
				continue
			}

			gasUsageMetric := fmt.Sprintf("orderbooktxs/%s/%s/gas", contractName, method.Name)
			sampler := metrics.ResettingSample(metrics.NewExpDecaySample(1028, 0.015))
			metrics.GetOrRegisterHistogram(gasUsageMetric, nil, sampler).Update(int64(receipt.GasUsed))
		}
	}
}

func EncodeLimitOrder(order LimitOrder) ([]byte, error) {
	limitOrderType, _ := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "ammIndex", Type: "uint256"},
		{Name: "trader", Type: "address"},
		{Name: "baseAssetQuantity", Type: "int256"},
		{Name: "price", Type: "uint256"},
		{Name: "salt", Type: "uint256"},
		{Name: "reduceOnly", Type: "bool"},
	})

	encodedLimitOrder, err := abi.Arguments{{Type: limitOrderType}}.Pack(order)
	if err != nil {
		return nil, fmt.Errorf("limit order packing failed: %w", err)
	}

	orderType, _ := abi.NewType("uint8", "uint8", nil)
	orderBytesType, _ := abi.NewType("bytes", "bytes", nil)
	encodedOrder, err := abi.Arguments{{Type: orderType}, {Type: orderBytesType}}.Pack(uint8(0), encodedLimitOrder)
	if err != nil {
		return nil, fmt.Errorf("order encoding failed: %w", err)
	}

	return encodedOrder, nil
}
