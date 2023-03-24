// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package sharedmemory

import (
	"bytes"
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/validators"
	"github.com/ava-labs/avalanchego/vms/components/avax"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/utils/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// TODO: structure this within the precompile test
func TestSharedMemoryRun(t *testing.T) {
	type test struct {
		caller       common.Address
		preCondition func(t *testing.T, state *state.StateDB)
		input        func() []byte
		suppliedGas  uint64
		readOnly     bool
		config       *Config

		expectedRes []byte
		expectedErr string

		assertState           func(t *testing.T, state *state.StateDB)
		validatePrecompileLog func(t *testing.T, log *types.Log)
	}

	caller := common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")
	receiver := common.HexToAddress("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B")

	snowCtx := snow.DefaultContextTest()
	snowCtx.AVAXAssetID = ids.GenerateTestID()
	snowCtx.SubnetID = ids.GenerateTestID()
	snowCtx.ChainID = ids.GenerateTestID()
	destinationChainID := ids.GenerateTestID()
	snowCtx.ValidatorState = &validators.TestState{
		GetSubnetIDF: func(_ context.Context, chainID ids.ID) (ids.ID, error) {
			subnetID, ok := map[ids.ID]ids.ID{
				snowCtx.ChainID:    snowCtx.SubnetID,
				destinationChainID: snowCtx.SubnetID,
			}[chainID]
			if !ok {
				return ids.Empty, errors.New("unknown chain")
			}
			return subnetID, nil
		},
	}

	for name, test := range map[string]test{
		"exportAVAX": {
			caller: caller,
			preCondition: func(t *testing.T, state *state.StateDB) {
				state.SetBalance(ContractAddress, big.NewInt(params.Ether)) // Simulate having sent 1ETH to the precompile
			},
			input: func() []byte {
				input, err := PackExportAVAX(
					ExportAVAXInput{
						destinationChainID,
						0,
						1,
						[]common.Address{receiver},
					},
				)
				require.NoError(t, err)

				return input
			},
			suppliedGas: ExportAVAXGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
			},
			validatePrecompileLog: func(t *testing.T, exportAVAXLog *types.Log) {
				// Validate topics
				require.Len(t, exportAVAXLog.Topics, 2)
				event, err := SharedMemoryABI.EventByID(exportAVAXLog.Topics[0])
				require.NoError(t, err)
				require.Equal(t, event.Name, "ExportAVAX")
				require.Equal(t, exportAVAXLog.Topics[1], common.Hash(destinationChainID))

				// Validate data
				ev := &exportAVAXEvent{}
				err = SharedMemoryABI.UnpackInputIntoInterface(ev, "ExportAVAX", exportAVAXLog.Data)
				require.NoError(t, err)
				require.Equal(t, uint64(params.GWei), ev.Amount)
				require.Equal(t, uint64(0), ev.Locktime)
				require.Equal(t, uint64(1), ev.Threshold)
				require.Len(t, ev.Addrs, 1, "expected 1 address in exportAVAX log")
				require.Equal(t, receiver, ev.Addrs[0])
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			db := rawdb.NewMemoryDatabase()
			state, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
			require.NoError(t, err)

			if test.preCondition != nil {
				test.preCondition(t, state)
			}

			blockContext := contract.NewMockBlockContext(big.NewInt(0), 0)
			accessibleState := contract.NewMockAccessibleState(state, blockContext, snowCtx)
			if test.config != nil {
				err := Module.Configure(nil, test.config, state, blockContext)
				require.NoError(t, err)
			}
			ret, remainingGas, err := SharedMemoryPrecompile.Run(accessibleState, test.caller, ContractAddress, test.input(), test.suppliedGas, test.readOnly)
			if len(test.expectedErr) != 0 {
				require.ErrorContains(t, err, test.expectedErr)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, uint64(0), remainingGas)
			require.Equal(t, test.expectedRes, ret)

			if test.assertState != nil {
				test.assertState(t, state)
			}
			if test.validatePrecompileLog == nil {
				require.Len(t, state.Logs(), 0)
			} else {
				logs := state.Logs()
				require.Len(t, logs, 1)
				log := logs[0]
				test.validatePrecompileLog(t, log)

				// TODO: clean this up and move it into helper function
				txHash := common.Hash{1, 2, 3}
				logIndex := 0 // TODO: should we change this type to uint32 in the function signature as well?
				chainID, reqs, err := acceptedLogsToSharedMemoryOps(snowCtx, txHash, logIndex, log.Topics, log.Data)
				require.NoError(t, err)
				require.Equal(t, chainID, destinationChainID)
				require.Len(t, reqs.PutRequests, 1)
				require.Len(t, reqs.RemoveRequests, 0)

				elem := reqs.PutRequests[0]
				require.Len(t, elem.Traits, 1)
				require.True(t, bytes.Equal(elem.Traits[0], receiver[:]))
				// Skip verification of the key
				utxo := &avax.UTXO{}
				v, err := codec.Codec.Unmarshal(elem.Value, utxo)
				require.NoError(t, err)
				require.Equal(t, uint16(0), v)

				// Verify the generated UTXO
				require.Equal(t, ids.ID(txHash), utxo.UTXOID.TxID)
				require.Equal(t, uint32(logIndex), utxo.UTXOID.OutputIndex)
				require.Equal(t, snowCtx.AVAXAssetID, utxo.AssetID())

				utxoOutput := utxo.Out.(*secp256k1fx.TransferOutput)
				require.Equal(t, uint64(params.GWei), utxoOutput.Amt)
				require.Equal(t, uint64(0), utxoOutput.OutputOwners.Locktime)
				require.Equal(t, uint32(1), utxoOutput.OutputOwners.Threshold)
				require.Len(t, utxoOutput.OutputOwners.Addrs, 1)
				require.Equal(t, ids.ShortID(receiver), utxoOutput.OutputOwners.Addrs[0])
			}
		})
	}
}
