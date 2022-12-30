package limitorders

import (
	"context"
	"errors"
	"math/big"
	"math/rand"
	"time"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/eth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

// using multiple private keys to make executeMatchedOrders contract call.
// This will be replaced by validator's private key and address
var userAddress1 = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
var privateKey1 = "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027"
var userAddress2 = "0x4Cf2eD3665F6bFA95cE6A11CFDb7A2EF5FC1C7E4"
var privateKey2 = "31b571bf6894a248831ff937bb49f7754509fe93bbd2517c9c73c4144c0e97dc"

type LimitOrderTxProcessor struct {
	txPool                   *core.TxPool
	orderBookABI             abi.ABI
	memoryDb                 *InMemoryDatabase
	orderBookContractAddress common.Address
}

func NewLimitOrderTxProcessor(txPool *core.TxPool, orderBookABI abi.ABI, memoryDb *InMemoryDatabase, orderBookContractAddress common.Address) *LimitOrderTxProcessor {
	return &LimitOrderTxProcessor{
		txPool:                   txPool,
		orderBookABI:             orderBookABI,
		memoryDb:                 memoryDb,
		orderBookContractAddress: orderBookContractAddress,
	}
}

func (lotp *LimitOrderTxProcessor) HandleOrderBookTx(tx *types.Transaction, blockNumber uint64, backend eth.EthAPIBackend) {
	m, err := getOrderBookContractCallMethod(tx, lotp.orderBookABI, lotp.orderBookContractAddress)
	if err == nil {
		input := tx.Data()
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

			baseAssetQuantity := int(order.BaseAssetQuantity.Int64())
			if baseAssetQuantity == 0 {
				log.Error("order not saved because baseAssetQuantity is zero")
				return
			}
			positionType := getPositionTypeBasedOnBaseAssetQuantity(baseAssetQuantity)
			price, _ := new(big.Float).SetInt(order.Price).Float64()
			limitOrder := &LimitOrder{
				PositionType:      positionType,
				UserAddress:       order.Trader.Hash().String(),
				BaseAssetQuantity: baseAssetQuantity,
				Price:             price,
				Salt:              order.Salt.String(),
				Status:            "unfulfilled",
				Signature:         signature,
				BlockNumber:       blockNumber,
				RawOrder:          in["order"],
				RawSignature:      in["signature"],
			}
			lotp.memoryDb.Add(limitOrder)
		}
		if m.Name == "executeMatchedOrders" && checkTxStatusSucess(backend, tx.Hash()) {
			signature1 := in["signature1"].([]byte)
			lotp.memoryDb.Delete(signature1)
			signature2 := in["signature2"].([]byte)
			lotp.memoryDb.Delete(signature2)
		}
	}
}

func (lotp *LimitOrderTxProcessor) ExecuteMatchedOrdersTx(incomingOrder LimitOrder, matchedOrder LimitOrder) error {
	//randomly selecting private key to get different validator profile on different nodes
	rand.Seed(time.Now().UnixNano())
	var privateKey, userAddress string
	if rand.Intn(10000)%2 == 0 {
		privateKey = privateKey1
		userAddress = userAddress1
	} else {
		privateKey = privateKey2
		userAddress = userAddress2
	}

	nonce := lotp.txPool.Nonce(common.HexToAddress(userAddress)) // admin address

	data, err := lotp.orderBookABI.Pack("executeMatchedOrders", incomingOrder.RawOrder, incomingOrder.Signature, matchedOrder.RawOrder, matchedOrder.Signature)
	if err != nil {
		log.Error("abi.Pack failed", "err", err)
		return err
	}
	key, err := crypto.HexToECDSA(privateKey) // admin private key
	if err != nil {
		log.Error("HexToECDSA failed", "err", err)
		return err
	}
	executeMatchedOrdersTx := types.NewTransaction(nonce, lotp.orderBookContractAddress, big.NewInt(0), 5000000, big.NewInt(80000000000), data)
	signer := types.NewLondonSigner(big.NewInt(321123))
	signedTx, err := types.SignTx(executeMatchedOrdersTx, signer, key)
	if err != nil {
		log.Error("types.SignTx failed", "err", err)
	}
	err = lotp.txPool.AddLocal(signedTx)
	if err != nil {
		log.Error("lop.txPool.AddLocal failed", "err", err)
		return err
	}
	return nil
}

func (lotp *LimitOrderTxProcessor) PurgeLocalTx() {
	pending := lotp.txPool.Pending(true)
	localAccounts := []common.Address{common.HexToAddress(userAddress1), common.HexToAddress(userAddress2)}

	for _, account := range localAccounts {
		if txs := pending[account]; len(txs) > 0 {
			for _, tx := range txs {
				m, err := getOrderBookContractCallMethod(tx, lotp.orderBookABI, lotp.orderBookContractAddress)
				if err == nil && m.Name == "executeMatchedOrders" {
					lotp.txPool.RemoveTx(tx.Hash())
				}
			}
		}
	}
}
func (lotp *LimitOrderTxProcessor) CheckIfOrderBookContractCall(tx *types.Transaction) bool {
	return checkIfOrderBookContractCall(tx, lotp.orderBookABI, lotp.orderBookContractAddress)
}

func getPositionTypeBasedOnBaseAssetQuantity(baseAssetQuantity int) string {
	if baseAssetQuantity > 0 {
		return "long"
	}
	return "short"
}

func checkTxStatusSucess(backend eth.EthAPIBackend, hash common.Hash) bool {
	ctx := context.Background()
	defer ctx.Done()

	_, blockHash, _, index, err := backend.GetTransaction(ctx, hash)
	if err != nil {
		log.Error("err in lop.backend.GetTransaction", "err", err)
		return false
	}
	receipts, err := backend.GetReceipts(ctx, blockHash)
	if err != nil {
		log.Error("err in lop.backend.GetReceipts", "err", err)
		return false
	}
	if len(receipts) <= int(index) {
		return false
	}
	receipt := receipts[index]
	return receipt.Status == uint64(1)
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
