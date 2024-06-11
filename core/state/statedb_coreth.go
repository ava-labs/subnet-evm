package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Retrieve the balance from the given address or 0 if object not found
func (s *StateDB) GetBalanceMultiCoin(addr common.Address, coinID common.Hash) *big.Int {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.BalanceMultiCoin(coinID, s.db)
	}
	return new(big.Int).Set(common.Big0)
}

// GetCommittedStateAP1 retrieves a value from the given account's committed storage trie.
func (s *StateDB) GetCommittedStateAP1(addr common.Address, hash common.Hash) common.Hash {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		NormalizeStateKey(&hash)
		return stateObject.GetCommittedState(hash)
	}
	return common.Hash{}
}

// AddBalance adds amount to the account associated with addr.
func (s *StateDB) AddBalanceMultiCoin(addr common.Address, coinID common.Hash, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.AddBalanceMultiCoin(coinID, amount, s.db)
	}
}

// SubBalance subtracts amount from the account associated with addr.
func (s *StateDB) SubBalanceMultiCoin(addr common.Address, coinID common.Hash, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SubBalanceMultiCoin(coinID, amount, s.db)
	}
}

func (s *StateDB) SetBalanceMultiCoin(addr common.Address, coinID common.Hash, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetBalanceMultiCoin(coinID, amount, s.db)
	}
}
