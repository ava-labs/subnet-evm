package limitorders

import (
	"context"
	"errors"
	"math/big"

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

var OrderBookContractAddress = common.HexToAddress("0x0300000000000000000000000000000000000069")
var MarginAccountContractAddress = common.HexToAddress("0x0300000000000000000000000000000000000070")
var ClearingHouseContractAddress = common.HexToAddress("0x0300000000000000000000000000000000000071")

// func SetContractFilesLocation(orderBook string, marginAccount string, clearingHouse string) {
// 	orderBookContractFileLocation = orderBook
// 	marginAccountContractFileLocation = marginAccount
// 	clearingHouseContractFileLocation = clearingHouse
// }

type LimitOrderTxProcessor interface {
	ExecuteMatchedOrdersTx(incomingOrder LimitOrder, matchedOrder LimitOrder, fillAmount *big.Int) error
	PurgeLocalTx()
	CheckIfOrderBookContractCall(tx *types.Transaction) bool
	ExecuteFundingPaymentTx() error
	ExecuteLiquidation(trader common.Address, matchedOrder LimitOrder, fillAmount *big.Int) error
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
}

// Order type is copy of Order struct defined in Orderbook contract
type Order struct {
	Trader            common.Address `json:"trader"`
	AmmIndex          *big.Int       `json:"ammIndex"`
	BaseAssetQuantity *big.Int       `json:"baseAssetQuantity"`
	Price             *big.Int       `json:"price"`
	Salt              *big.Int       `json:"salt"`
}

func NewLimitOrderTxProcessor(txPool *core.TxPool, memoryDb LimitOrderDatabase, backend *eth.EthAPIBackend) LimitOrderTxProcessor {
	orderBookABI, err := abi.FromSolidityJson(string(orderBookAbi))
	if err != nil {
		panic(err)
	}

	clearingHouseABI, err := abi.FromSolidityJson(string(clearingHouseAbi))
	if err != nil {
		panic(err)
	}

	marginAccountABI, err := abi.FromSolidityJson(string(marginAccountAbi))
	if err != nil {
		panic(err)
	}

	return &limitOrderTxProcessor{
		txPool:                       txPool,
		orderBookABI:                 orderBookABI,
		clearingHouseABI:             clearingHouseABI,
		marginAccountABI:             marginAccountABI,
		memoryDb:                     memoryDb,
		orderBookContractAddress:     OrderBookContractAddress,
		clearingHouseContractAddress: ClearingHouseContractAddress,
		marginAccountContractAddress: MarginAccountContractAddress,
		backend:                      backend,
	}
}

func (lotp *limitOrderTxProcessor) ExecuteLiquidation(trader common.Address, matchedOrder LimitOrder, fillAmount *big.Int) error {
	return lotp.executeLocalTx(lotp.orderBookContractAddress, lotp.orderBookABI, "liquidateAndExecuteOrder", trader.String(), matchedOrder.RawOrder, matchedOrder.Signature, fillAmount)
}

func (lotp *limitOrderTxProcessor) ExecuteFundingPaymentTx() error {
	return lotp.executeLocalTx(lotp.orderBookContractAddress, lotp.orderBookABI, "settleFunding")
}

func (lotp *limitOrderTxProcessor) ExecuteMatchedOrdersTx(incomingOrder LimitOrder, matchedOrder LimitOrder, fillAmount *big.Int) error {
	orders := make([]Order, 2)
	orders[0], orders[1] = getOrderFromRawOrder(incomingOrder.RawOrder), getOrderFromRawOrder(matchedOrder.RawOrder)

	log.Info("matching", "long order price", orders[0].Price, "short order price", orders[1].Price)
	signatures := make([][]byte, 2)
	signatures[0] = incomingOrder.Signature
	signatures[1] = matchedOrder.Signature

	return lotp.executeLocalTx(lotp.orderBookContractAddress, lotp.orderBookABI, "executeMatchedOrders", orders, signatures, fillAmount)
}

func (lotp *limitOrderTxProcessor) executeLocalTx(contract common.Address, contractABI abi.ABI, method string, args ...interface{}) error {
	nonce := lotp.txPool.Nonce(common.HexToAddress(userAddress1)) // admin address

	data, err := contractABI.Pack(method, args...)
	if err != nil {
		log.Error("abi.Pack failed", "err", err)
		return err
	}
	key, err := crypto.HexToECDSA(privateKey1) // admin private key
	if err != nil {
		log.Error("HexToECDSA failed", "err", err)
		return err
	}
	tx := types.NewTransaction(nonce, contract, big.NewInt(0), 5000000, big.NewInt(60000000000), data)
	signer := types.NewLondonSigner(lotp.backend.ChainConfig().ChainID)
	signedTx, err := types.SignTx(tx, signer, key)
	if err != nil {
		log.Error("types.SignTx failed", "err", err)
	}
	err = lotp.txPool.AddLocal(signedTx)
	if err != nil {
		log.Error("lop.txPool.AddLocal failed", "err", err, "tx", signedTx.Hash().String(), "nonce", nonce)
		return err
	}
	log.Info("executeLocalTx - AddLocal success", "tx", signedTx.Hash().String(), "nonce", nonce)

	return nil
}

func (lotp *limitOrderTxProcessor) PurgeLocalTx() {
	pending := lotp.txPool.Pending(true)
	localAccounts := []common.Address{common.HexToAddress(userAddress1), common.HexToAddress(userAddress2)}

	for _, account := range localAccounts {
		if txs := pending[account]; len(txs) > 0 {
			for _, tx := range txs {
				_, err := getOrderBookContractCallMethod(tx, lotp.orderBookABI, lotp.orderBookContractAddress)
				if err == nil {
					lotp.txPool.RemoveTx(tx.Hash())
				}
			}
		}
	}
}

func (lotp *limitOrderTxProcessor) CheckIfOrderBookContractCall(tx *types.Transaction) bool {
	return checkIfOrderBookContractCall(tx, lotp.orderBookABI, lotp.orderBookContractAddress)
}

func getPositionTypeBasedOnBaseAssetQuantity(baseAssetQuantity *big.Int) string {
	if baseAssetQuantity.Sign() == 1 {
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
