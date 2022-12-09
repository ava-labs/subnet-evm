package evm

import (
	"context"
	"io/ioutil"
	"math/big"

	"sync"
	"time"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/eth"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/plugin/evm/limitorders"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

var orderBookContractFileLocation = "contract-examples/artifacts/contracts/OrderBook.sol/OrderBook.json"

type LimitOrderProcesser interface {
	ListenAndProcessTransactions()
	AddMatchingOrdersToTxPool()
}

type limitOrderProcesser struct {
	ctx          *snow.Context
	chainConfig  *params.ChainConfig
	txPool       *core.TxPool
	shutdownChan <-chan struct{}
	shutdownWg   *sync.WaitGroup
	backend      *eth.EthAPIBackend
	memoryDb     limitorders.InMemoryDatabase
	orderBookABI abi.ABI
}

func SetOrderBookContractFileLocation(location string) {
	orderBookContractFileLocation = location
}

func NewLimitOrderProcesser(ctx *snow.Context, chainConfig *params.ChainConfig, txPool *core.TxPool, shutdownChan <-chan struct{}, shutdownWg *sync.WaitGroup, backend *eth.EthAPIBackend) LimitOrderProcesser {
	jsonBytes, _ := ioutil.ReadFile(orderBookContractFileLocation)
	orderBookAbi, err := abi.FromSolidityJson(string(jsonBytes))
	if err != nil {
		panic(err)
	}

	return &limitOrderProcesser{
		ctx:          ctx,
		chainConfig:  chainConfig,
		txPool:       txPool,
		shutdownChan: shutdownChan,
		shutdownWg:   shutdownWg,
		backend:      backend,
		memoryDb:     limitorders.NewInMemoryDatabase(),
		orderBookABI: orderBookAbi,
	}
}

func (lop *limitOrderProcesser) ListenAndProcessTransactions() {
	lop.listenAndStoreLimitOrderTransactions()
}

func (lop *limitOrderProcesser) AddMatchingOrdersToTxPool() {
	orders := lop.memoryDb.GetAllOrders()
	for _, order := range orders {
		nonce := lop.txPool.Nonce(common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")) // admin address

		data, err := lop.orderBookABI.Pack("executeTestOrder", order.RawOrder, order.Signature)
		if err != nil {
			log.Error("abi.Pack failed", "err", err)
		}
		key, err := crypto.HexToECDSA("56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027") // admin private key
		if err != nil {
			log.Error("HexToECDSA failed", "err", err)
		}
		executeOrderTx := types.NewTransaction(nonce, common.HexToAddress("0x52C84043CD9c865236f11d9Fc9F56aa003c1f922"), big.NewInt(0), 8000000, big.NewInt(250000000), data)
		signer := types.NewLondonSigner(big.NewInt(99999))
		signedTx, err := types.SignTx(executeOrderTx, signer, key)
		if err != nil {
			log.Error("types.SignTx failed", "err", err)
		}
		err = lop.txPool.AddLocal(signedTx)
		if err != nil {
			log.Error("lop.txPool.AddLocal failed", "err", err)
		}
		log.Info("#### AddMatchingOrdersToTxPool", "executeOrder", signedTx.Hash().String(), "from signature", string(order.Signature))
		lop.memoryDb.Delete(order.Signature)
	}
}

func (lop *limitOrderProcesser) listenAndStoreLimitOrderTransactions() {

	type Order struct {
		Trader            common.Address `json:"trader"`
		BaseAssetQuantity *big.Int       `json:"baseAssetQuantity"`
		Price             *big.Int       `json:"price"`
		Salt              *big.Int       `json:"salt"`
	}

	newHeadChan := make(chan core.NewTxPoolHeadEvent)
	lop.txPool.SubscribeNewHeadEvent(newHeadChan)

	lop.shutdownWg.Add(1)
	go lop.ctx.Log.RecoverAndPanic(func() {
		defer lop.shutdownWg.Done()

		for {
			select {
			case newHeadEvent := <-newHeadChan:
				tsHashes := []string{}
				for _, tx := range newHeadEvent.Block.Transactions() {
					tsHashes = append(tsHashes, tx.Hash().String())
					parseTx(lop.orderBookABI, lop.memoryDb, *lop.backend, tx) // parse update in memory db
				}
				log.Info("$$$$$ New head event", "number", newHeadEvent.Block.Header().Number, "tx hashes", tsHashes,
					"miner", newHeadEvent.Block.Coinbase().String(),
					"root", newHeadEvent.Block.Header().Root.String(), "gas used", newHeadEvent.Block.Header().GasUsed,
					"nonce", newHeadEvent.Block.Header().Nonce)
			case <-lop.shutdownChan:
				return
			}
		}
	})
}

func parseTx(orderBookABI abi.ABI, memoryDb limitorders.InMemoryDatabase, backend eth.EthAPIBackend, tx *types.Transaction) {
	input := tx.Data()
	if len(input) < 4 {
		log.Info("transaction data has less than 3 fields")
		return
	}
	method := input[:4]
	m, err := orderBookABI.MethodById(method)
	if err == nil {
		in := make(map[string]interface{})
		_ = m.Inputs.UnpackIntoMap(in, input[4:])
		if m.Name == "placeOrder" {
			log.Info("##### in ParseTx", "placeOrder tx hash", tx.Hash().String())
			order, _ := in["order"].(struct {
				Trader            common.Address `json:"trader"`
				BaseAssetQuantity *big.Int       `json:"baseAssetQuantity"`
				Price             *big.Int       `json:"price"`
				Salt              *big.Int       `json:"salt"`
			})
			signature := in["signature"].([]byte)

			var positionType string
			if order.BaseAssetQuantity.Int64() > 0 {
				positionType = "long"
			} else {
				positionType = "short"
			}
			price, _ := new(big.Float).SetInt(order.Price).Float64()
			limitOrder := limitorders.LimitOrder{
				PositionType:      positionType,
				UserAddress:       order.Trader.Hash().String(),
				BaseAssetQuantity: int(order.BaseAssetQuantity.Uint64()),
				Price:             price,
				Salt:              order.Salt.String(),
				Signature:         signature,
				RawOrder:          in["order"],
				RawSignature:      in["signature"],
			}
			memoryDb.Add(limitOrder)
		}
		if m.Name == "executeTestOrder" {
			log.Info("##### in ParseTx", "executeTestOrder tx hash", tx.Hash().String())
			go pollForReceipt(backend, tx.Hash())
			signature := in["signature"].([]byte)
			memoryDb.Delete(signature)
		}
	}
}

func pollForReceipt(backend eth.EthAPIBackend, txHash common.Hash) {
	for i := 0; i < 10; i++ {
		receipt := getTxReceipt(backend, txHash)
		if receipt != nil {
			log.Info("receipt found", "tx", txHash.String(), "receipt", receipt)
			return
		}
		time.Sleep(time.Second * 5)
	}
	log.Info("receipt not found", "tx", txHash.String())
}

func getTxReceipt(backend eth.EthAPIBackend, hash common.Hash) *types.Receipt {
	ctx := context.Background()
	_, blockHash, _, index, err := backend.GetTransaction(ctx, hash)
	if err != nil {
		log.Error("err in lop.backend.GetTransaction", "err", err)
	}
	receipts, err := backend.GetReceipts(ctx, blockHash)
	if err != nil {
		log.Error("err in lop.backend.GetReceipts", "err", err)
	}
	if len(receipts) <= int(index) {
		// log.Info("len(receipts) <= int(index)", "len(receipts)", len(receipts), "index", index)
		return nil
	}
	receipt := receipts[index]
	return receipt
}
