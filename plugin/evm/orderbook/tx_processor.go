package orderbook

import (
	"context"
	"crypto/ecdsa"
	"time"

	// "encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/core/txpool"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/eth"
	"github.com/ava-labs/subnet-evm/metrics"
	"github.com/ava-labs/subnet-evm/plugin/evm/orderbook/abis"
	"github.com/ava-labs/subnet-evm/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

var OrderBookContractAddress = common.HexToAddress("0x03000000000000000000000000000000000000b0")
var MarginAccountContractAddress = common.HexToAddress("0x03000000000000000000000000000000000000b1")
var ClearingHouseContractAddress = common.HexToAddress("0x03000000000000000000000000000000000000b2")
var LimitOrderBookContractAddress = common.HexToAddress("0x03000000000000000000000000000000000000b3")
var IOCOrderBookContractAddress = common.HexToAddress("0x03000000000000000000000000000000000000b4")

type LimitOrderTxProcessor interface {
	GetOrderBookTxsCount() uint64
	SetOrderBookTxsBlockNumber(blockNumber uint64)
	PurgeOrderBookTxs()
	ExecuteMatchedOrdersTx(incomingOrder Order, matchedOrder Order, fillAmount *big.Int) error
	ExecuteFundingPaymentTx() error
	ExecuteSamplePITx() error
	ExecuteLiquidation(trader common.Address, matchedOrder Order, fillAmount *big.Int) error
	UpdateMetrics(block *types.Block)
	ExecuteLimitOrderCancel(orderIds []LimitOrder) error
}

type ValidatorTxFeeConfig struct {
	baseFeeEstimate *big.Int
	blockNumber     uint64
}

type limitOrderTxProcessor struct {
	txPool                       *txpool.TxPool
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

func NewLimitOrderTxProcessor(txPool *txpool.TxPool, memoryDb LimitOrderDatabase, backend *eth.EthAPIBackend, validatorPrivateKey string) LimitOrderTxProcessor {
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
	txHash, err := lotp.executeLocalTx(lotp.clearingHouseContractAddress, lotp.clearingHouseABI, "settleFunding")
	log.Info("ExecuteFundingPaymentTx", "txHash", txHash.String(), "err", err)
	return err
}

func (lotp *limitOrderTxProcessor) ExecuteSamplePITx() error {
	txHash, err := lotp.executeLocalTx(lotp.clearingHouseContractAddress, lotp.clearingHouseABI, "samplePI")
	log.Info("ExecuteSamplePITx", "txHash", txHash.String(), "err", err)
	if err == nil {
		lotp.memoryDb.SignalSamplePIAttempted(uint64(time.Now().Unix()))
	}
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
	txHash, err := lotp.executeLocalTx(lotp.orderBookContractAddress, lotp.orderBookABI, "cancelOrdersWithLowMargin", orders)
	log.Info("ExecuteLimitOrderCancel", "orders", orders, "txHash", txHash.String(), "err", err)
	return err
}

func (lotp *limitOrderTxProcessor) executeLocalTx(contract common.Address, contractABI abi.ABI, method string, args ...interface{}) (common.Hash, error) {
	var txHash common.Hash
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
	txFee := lotp.getTransactionFee()
	tx := types.NewTransaction(nonce, contract, big.NewInt(0), 1500000, txFee, data)
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

	return txHash, nil
}

func (lotp *limitOrderTxProcessor) getTransactionFee() *big.Int {
	latest := lotp.backend.CurrentHeader()
	latestBlockNumber := latest.Number.Uint64()

	// if the fee is already calculated for this block, then return it
	if lotp.validatorTxFeeConfig.blockNumber == latestBlockNumber {
		return lotp.validatorTxFeeConfig.baseFeeEstimate
	}

	baseFeeEstimate, err := lotp.backend.SuggestPrice(context.Background())
	if err != nil {
		log.Error("getBaseFeeEstimate - SuggestPrice failed", "err", err)
		return big.NewInt(65_000000000) // hardcoded to 65 gwei
	}
	// add 10%
	baseFeeEstimate.Add(baseFeeEstimate, big.NewInt(0).Div(baseFeeEstimate, big.NewInt(10)))

	feeConfig, _, err := lotp.backend.GetFeeConfigAt(latest)
	if err != nil {
		log.Error("getBaseFeeEstimate - GetFeeConfigAt failed", "err", err)
		// if feeConfig can't be obtained, then add another 10% to the baseFeeEstimate
		baseFeeEstimate.Add(baseFeeEstimate, big.NewInt(0).Div(baseFeeEstimate, big.NewInt(10)))
		return baseFeeEstimate
	}
	// assuming pessimistically that the block is being produced within a second of the latest block
	// we calculate the block gas cost as the latest block gas cost + the block gas cost step
	blockGasCost := big.NewInt(0).Add(latest.BlockGasCost, feeConfig.BlockGasCostStep)

	// assuming a minimum gas usage of 200k for a tx, we calculate the tip such that the entire block has an effective tip above the threshold
	// example calculation for blockGasCost = 10,000, baseFeeEstimate = 60 gwei, tx gas usage = 200,000
	// tip = (10000 * 60 * 1e9) / 200000 = 3 gwei
	tip := big.NewInt(0).Div(big.NewInt(0).Mul(blockGasCost, baseFeeEstimate), big.NewInt(200000))

	totalFee := baseFeeEstimate.Add(baseFeeEstimate, tip)

	lotp.validatorTxFeeConfig.baseFeeEstimate = totalFee
	lotp.validatorTxFeeConfig.blockNumber = latestBlockNumber
	return totalFee
}

func (lotp *limitOrderTxProcessor) PurgeOrderBookTxs() {
	lotp.txPool.PurgeOrderBookTxs()
}

func (lotp *limitOrderTxProcessor) GetOrderBookTxsCount() uint64 {
	return lotp.txPool.GetOrderBookTxsCount()
}

func (lotp *limitOrderTxProcessor) SetOrderBookTxsBlockNumber(blockNumber uint64) {
	lotp.txPool.SetOrderBookTxsBlockNumber(blockNumber)
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
	// blank key is allowed for non-validators
	if key == "" {
		return common.Address{}, nil
	}
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

func formatReceiptForLogging(receipt *types.Receipt) string {
	return fmt.Sprintf("Receipt{Status: %d, CumulativeGasUsed: %d, GasUsed: %d, EffectiveGasPrice: %d, BlockNumber: %d}",
		receipt.Status, receipt.CumulativeGasUsed, receipt.GasUsed, receipt.EffectiveGasPrice, receipt.BlockNumber)
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
	timestamp := block.Header().Time
	signer := types.MakeSigner(lotp.backend.ChainConfig(), bigblock, timestamp)

	currentBlock := lotp.backend.CurrentBlock() // head block
	headBlockLagHistogram.Update(int64(currentBlock.Number.Uint64() - block.NumberU64()))

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

			if contractAddress != nil && (lotp.orderBookContractAddress == *contractAddress || lotp.clearingHouseContractAddress == *contractAddress) {
				note := "success"
				if receipt.Status == 0 {
					log.Error("this validator's tx failed", "method", method.Name, "tx", tx.Hash().String(),
						"receipt", formatReceiptForLogging(receipt), "from", from.String())
					note = "failure"
				}
				counterName := fmt.Sprintf("orderbooktxs/%s/%s", method.Name, note)
				metrics.GetOrRegisterCounter(counterName, nil).Inc(1)
			}
		}

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

			// measure the gas usage irrespective of whether the tx is from this validator or not
			gasUsageMetric := fmt.Sprintf("orderbooktxs/%s/%s/gas", contractName, method.Name)
			sampler := metrics.ResettingSample(metrics.NewExpDecaySample(1028, 0.015))
			metrics.GetOrRegisterHistogram(gasUsageMetric, nil, sampler).Update(int64(receipt.GasUsed))

			// log the failure for validator txs irrespective of whether the tx is from this validator or not
			// this will help us identify tx failures that are not due to a hubble's validator
			validatorMethods := []string{"liquidateAndExecuteOrder", "executeMatchedOrders", "settleFunding", "samplePI", "cancelOrdersWithLowMargin"}
			if receipt.Status == 0 && utils.ContainsString(validatorMethods, method.Name) {
				log.Error("validator tx failed", "method", method.Name, "contractName", contractName, "tx", tx.Hash().String(), "from", from.String(), "receipt", formatReceiptForLogging(receipt))
			}
		}
	}
}
