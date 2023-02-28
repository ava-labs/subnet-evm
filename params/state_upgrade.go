// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
)

//go:generate go run github.com/fjl/gencodec -type StateUpgradeAccount -field-override stateUpgradeAccountMarshaling -out gen_state_upgrade_account.go

// StateUpgrade describes the modifications to be made to the state during
// a state upgrade.
type StateUpgrade struct {
	BlockTimestamp *big.Int `json:"blockTimestamp,omitempty"`

	// map from account address to the modification to be made to the account.
	StateUpgradeAccounts StateUpgradeAccounts `json:"accounts"`
}

// StateUpgradeAccount describes the modifications to be made to an account during
// a state upgrade.
type StateUpgradeAccount struct {
	Code          []byte                      `json:"code,omitempty"`
	Storage       map[common.Hash]common.Hash `json:"storage,omitempty"`
	BalanceChange *big.Int                    `json:"balanceChange,omitempty"`
}

func (s *StateUpgrade) Equal(other *StateUpgrade) bool {
	return reflect.DeepEqual(s, other)
}

type stateUpgradeAccountMarshaling struct {
	Code          hexutil.Bytes
	Storage       map[storageJSON]storageJSON
	BalanceChange *math.HexOrDecimal256
}

// storageJSON represents a 256 bit byte array, but allows less than 256 bits when
// unmarshaling from hex.
type storageJSON common.Hash

func (h *storageJSON) UnmarshalText(text []byte) error {
	text = bytes.TrimPrefix(text, []byte("0x"))
	if len(text) > 64 {
		return fmt.Errorf("too many hex characters in storage key/value %q", text)
	}
	offset := len(h) - len(text)/2 // pad on the left
	if _, err := hex.Decode(h[offset:], text); err != nil {
		fmt.Println(err)
		return fmt.Errorf("invalid hex storage key/value %q", text)
	}
	return nil
}

func (h storageJSON) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}

// StateUpgradeAlloc specifies the state upgrade type.
type StateUpgradeAccounts map[common.Address]StateUpgradeAccount

func (s *StateUpgradeAccounts) UnmarshalJSON(data []byte) error {
	m := make(map[common.UnprefixedAddress]StateUpgradeAccount)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	*s = make(StateUpgradeAccounts)
	for addr, a := range m {
		(*s)[common.Address(addr)] = a
	}
	return nil
}
