package evm

import (
	"context"
	"io/ioutil"
	"math/big"
	"sync"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/eth"
	"github.com/ava-labs/subnet-evm/params"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

type LimitOrderProcesser interface {
	ListenAndProcessLimitOrderTransactions()
}

type limitOrderProcesser struct {
	ctx          *snow.Context
	chainConfig  *params.ChainConfig
	txPool       *core.TxPool
	shutdownChan <-chan struct{}
	shutdownWg   *sync.WaitGroup
	backend      *eth.EthAPIBackend
}

func NewLimitOrderProcesser(ctx *snow.Context, chainConfig *params.ChainConfig, txPool *core.TxPool, shutdownChan <-chan struct{}, shutdownWg *sync.WaitGroup, backend *eth.EthAPIBackend) LimitOrderProcesser {
	return &limitOrderProcesser{
		ctx:          ctx,
		chainConfig:  chainConfig,
		txPool:       txPool,
		shutdownChan: shutdownChan,
		shutdownWg:   shutdownWg,
		backend:      backend,
	}
}

func (lop *limitOrderProcesser) ListenAndProcessLimitOrderTransactions() {
	lop.listenAndStoreLimitOrderTransactions()
}

type Order struct {
	Trader            common.Address
	BaseAssetQuantity *big.Int
	Price             *big.Int
	Salt              *big.Int
}

func (lop *limitOrderProcesser) listenAndStoreLimitOrderTransactions() {
	txSubmitChan := make(chan core.NewTxsEvent)
	lop.txPool.SubscribeNewTxsEvent(txSubmitChan)
	lop.shutdownWg.Add(1)
	go lop.ctx.Log.RecoverAndPanic(func() {
		defer lop.shutdownWg.Done()

		jsonBytes, _ := ioutil.ReadFile("contract-examples/artifacts/contracts/OrderBook.sol/OrderBook.json")
		orderBookAbi, err := abi.FromSolidityJson(string(jsonBytes))
		jsonBytes, _ = ioutil.ReadFile("contract-examples/artifacts/contracts/ERC20NativeMinter.sol/ERC20NativeMinter.json")
		minterAbi, err := abi.FromSolidityJson(string(jsonBytes))
		if err != nil {
			panic(err)
		}
		// log.Info("####", "Abi.Methods", orderBookAbi.Methods)

		for {
			select {
			case txsEvent := <-txSubmitChan:
				log.Info("New transaction event detected")

				for i := 0; i < len(txsEvent.Txs); i++ {
					tx := txsEvent.Txs[i]
					if tx.To() != nil && tx.Data() != nil && len(tx.Data()) != 0 {
						log.Info("transaction", "to is", tx.To().String())
						input := tx.Data() // "input" field above
						log.Info("transaction", "data is", input)
						if len(input) < 4 {
							log.Info("transaction data has less than 3 fields")
							continue
						}
						method := input[:4]
						m, err := minterAbi.MethodById(method)
						if err == nil {
							log.Info("transaction was made by minter contract")
							log.Info("transaction", "method name", m.Name)
							in := make(map[string]interface{})
							_ = m.Inputs.UnpackIntoMap(in, input[4:])
							log.Info("transaction", "amount in is: %+v\n", in["amount"])
							log.Info("transaction", "to is: %+v\n", in["to"])


							nonce := lop.txPool.Nonce(common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"))
							log.Info("###", "nonce", nonce)
			
							data, err := orderBookAbi.Pack("placeOrder", &Order{
								Trader:            common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"),
								BaseAssetQuantity: big.NewInt(3),
								Price:             big.NewInt(3),
								Salt:              big.NewInt(12837918),
							}, []byte("hash1"))
							if err != nil {
								log.Error("abi.Pack failed", "err", err)
							}
							log.Info("####", "data", data)
							key, err := crypto.HexToECDSA("56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027")
							if err != nil {
								log.Error("HexToECDSA failed", "err", err)
							}
							tx := types.NewTransaction(nonce, common.HexToAddress("0x52C84043CD9c865236f11d9Fc9F56aa003c1f922"), big.NewInt(0), 8000000, big.NewInt(250000000), data)
							signer := types.NewLondonSigner(big.NewInt(99999))
							signedTx, err := types.SignTx(tx, signer, key)
							if err != nil {
								log.Error("types.SignTx failed", "err", err)
							}
			
							// UNCOMMENT TO SEND TX ON EVERY TX
							log.Trace("##", "signedTx", signedTx)
							err = lop.backend.SendTx(context.Background(), signedTx)
							if err != nil {
								log.Error("SendTx failed", "err", err, "ctx", context.Background())
							}
			
						} else {
							log.Info("### transaction - incorrect ABI -  received on another contract")
						}
					}
				}
			case <-lop.shutdownChan:
				return
			}
		}
	})
}
