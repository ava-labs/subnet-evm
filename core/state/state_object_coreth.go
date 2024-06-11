package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// AddBalanceMultiCoin adds amount of coinID to s's balance.
// It is used to add multicoin funds to the destination account of a transfer.
func (s *stateObject) AddBalanceMultiCoin(coinID common.Hash, amount *big.Int, db Database) {
	if amount.Sign() == 0 {
		if s.empty() {
			s.touch()
		}

		return
	}
	s.SetBalanceMultiCoin(coinID, new(big.Int).Add(s.BalanceMultiCoin(coinID, db), amount), db)
}

// SubBalanceMultiCoin removes amount of coinID from s's balance.
// It is used to remove multicoin funds from the origin account of a transfer.
func (s *stateObject) SubBalanceMultiCoin(coinID common.Hash, amount *big.Int, db Database) {
	if amount.Sign() == 0 {
		return
	}
	s.SetBalanceMultiCoin(coinID, new(big.Int).Sub(s.BalanceMultiCoin(coinID, db), amount), db)
}

func (s *stateObject) SetBalanceMultiCoin(coinID common.Hash, amount *big.Int, db Database) {
	s.EnableMultiCoin()
	NormalizeCoinID(&coinID)
	s.SetState(coinID, common.BigToHash(amount))
}
func (s *stateObject) enableMultiCoin() {
	s.data.IsMultiCoin = true
}

// NormalizeCoinID ORs the 0th bit of the first byte in
// [coinID], which ensures this bit will be 1 and all other
// bits are left the same.
// This partitions multicoin storage from normal state storage.
func NormalizeCoinID(coinID *common.Hash) {
	coinID[0] |= 0x01
}

// NormalizeStateKey ANDs the 0th bit of the first byte in
// [key], which ensures this bit will be 0 and all other bits
// are left the same.
// This partitions normal state storage from multicoin storage.
func NormalizeStateKey(key *common.Hash) {
	key[0] &= 0xfe
}

func (s *stateObject) BalanceMultiCoin(coinID common.Hash, db Database) *big.Int {
	NormalizeCoinID(&coinID)
	return s.GetState(coinID).Big()
}

func (s *stateObject) EnableMultiCoin() bool {
	if s.data.IsMultiCoin {
		return false
	}
	s.db.journal.append(multiCoinEnable{
		account: &s.address,
	})
	s.enableMultiCoin()
	return true
}
