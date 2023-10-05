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
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/ava-labs/subnet-evm/utils/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestSharedMemoryRun(t *testing.T) {
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

	tests := map[string]testutils.PrecompileTest{
		"exportAVAX": {
			SnowCtx: snowCtx,
			Caller:  caller,
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				state.AddBalance(ContractAddress, big.NewInt(params.Ether)) // Simulate having sent 1ETH to the precompile
			},
			InputFn: func(t testing.TB) []byte {
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
			SuppliedGas: ExportAVAXGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, s contract.StateDB) {
				// TODO: fix this to not require a type cast
				logs := s.(*state.StateDB).Logs()
				require.Len(t, logs, 1)
				exportAVAXLog := logs[0]

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

				// TODO: clean this up and move it into helper function
				txHash := common.Hash{1, 2, 3}
				logIndex := 0 // TODO: should we change this type to uint32 in the function signature as well?
				chainID, reqs, err := acceptedLogsToSharedMemoryOps(snowCtx, txHash, logIndex, exportAVAXLog.Topics, exportAVAXLog.Data)
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
			},
		},
	}
	testutils.RunPrecompileTests(t, Module, state.NewTestStateDB, tests)
}
