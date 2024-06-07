package types

import (
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

var MultiCoinEnabled = false

func (obj *StateAccount) EncodeRLP(_w io.Writer) error {
	w := rlp.NewEncoderBuffer(_w)
	_tmp0 := w.List()
	w.WriteUint64(obj.Nonce)
	if obj.Balance == nil {
		w.Write(rlp.EmptyString)
	} else {
		if obj.Balance.Sign() == -1 {
			return rlp.ErrNegativeBigInt
		}
		w.WriteBigInt(obj.Balance)
	}
	w.WriteBytes(obj.Root[:])
	w.WriteBytes(obj.CodeHash)
	if MultiCoinEnabled {
		w.WriteBool(obj.IsMultiCoin)
	}
	w.ListEnd(_tmp0)
	return w.Flush()
}

func (obj *StateAccount) DecodeRLP(dec *rlp.Stream) error {
	var _tmp0 StateAccount
	{
		if _, err := dec.List(); err != nil {
			return err
		}
		// Nonce:
		_tmp1, err := dec.Uint64()
		if err != nil {
			return err
		}
		_tmp0.Nonce = _tmp1
		// Balance:
		_tmp2, err := dec.BigInt()
		if err != nil {
			return err
		}
		_tmp0.Balance = _tmp2
		// Root:
		var _tmp3 common.Hash
		if err := dec.ReadBytes(_tmp3[:]); err != nil {
			return err
		}
		_tmp0.Root = _tmp3
		// CodeHash:
		_tmp4, err := dec.Bytes()
		if err != nil {
			return err
		}
		_tmp0.CodeHash = _tmp4
		// IsMultiCoin:
		if MultiCoinEnabled {
			_tmp5, err := dec.Bool()
			if err != nil {
				return err
			}
			_tmp0.IsMultiCoin = _tmp5
		}
		if err := dec.ListEnd(); err != nil {
			return err
		}
	}
	*obj = _tmp0
	return nil
}
