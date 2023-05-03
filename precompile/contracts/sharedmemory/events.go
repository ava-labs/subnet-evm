// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package sharedmemory

import (
	"fmt"

	"github.com/ava-labs/avalanchego/chains/atomic"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/vms/components/avax"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/subnet-evm/utils/codec"
	"github.com/ethereum/go-ethereum/common"
)

// exportAVAXEvent includes the non-indexed fields of the ExportAVAX event.
// This allows us to unpack the log data into this struct. Note indexed
// fields appear in the topics array and are not included in the log data.
type exportAVAXEvent struct {
	Amount    uint64
	Locktime  uint64
	Threshold uint64
	Addrs     []common.Address
}

func ExportAVAXEventToUTXO(assetID ids.ID, txHash common.Hash, logIndex int, topics []common.Hash, logData []byte) (*avax.UTXO, error) {
	// Parse the log data into the exportAVAXEvent struct
	ev := &exportAVAXEvent{}
	err := SharedMemoryABI.UnpackIntoInterface(ev, "ExportAVAX", logData)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack ExportAVAX event data: %w", err)
	}

	addrs := make([]ids.ShortID, 0, len(ev.Addrs))
	for _, addr := range ev.Addrs {
		addrs = append(addrs, ids.ShortID(addr))
	}
	utxo := &avax.UTXO{
		// Derive unique UTXOID from txHash and log index
		UTXOID: avax.UTXOID{
			TxID:        ids.ID(txHash),
			OutputIndex: uint32(logIndex),
		},
		Asset: avax.Asset{ID: assetID},
		Out: &secp256k1fx.TransferOutput{
			Amt: ev.Amount,
			OutputOwners: secp256k1fx.OutputOwners{
				Locktime:  ev.Locktime,
				Threshold: uint32(ev.Threshold), // TODO make the actual type uint32 to correspond to this
				Addrs:     addrs,
			},
		},
	}
	return utxo, nil
}

func handleExportAVAX(snowCtx *snow.Context, txHash common.Hash, logIndex int, topics []common.Hash, logData []byte) (ids.ID, *atomic.Requests, error) {
	// Parse the topics data.
	// TODO: Improve this by using the ABI to unpack the topics.
	destinationChainID := ids.ID(topics[1])
	utxo, err := ExportAVAXEventToUTXO(snowCtx.AVAXAssetID, txHash, logIndex, topics, logData)
	if err != nil {
		return ids.ID{}, nil, err
	}

	utxoBytes, err := codec.Codec.Marshal(codec.CodecVersion, utxo)
	if err != nil {
		return ids.ID{}, nil, err
	}
	utxoID := utxo.InputID()
	elem := &atomic.Element{
		Key:   utxoID[:],
		Value: utxoBytes,
	}
	if out, ok := utxo.Out.(avax.Addressable); ok {
		elem.Traits = out.Addresses()
	}

	return destinationChainID, &atomic.Requests{
		PutRequests: []*atomic.Element{elem},
	}, nil
}

// exportUTXOEvent includes the non-indexed fields of the ExportUTXO event.
// This allows us to unpack the log data into this struct. Note indexed
// fields appear in the topics array and are not included in the log data.
type exportUTXOEvent struct {
	Amount    uint64
	Locktime  uint64
	Threshold uint64
	Addrs     []common.Address
}

func handleExportUTXO(snowCtx *snow.Context, txHash common.Hash, logIndex int, topics []common.Hash, logData []byte) (ids.ID, *atomic.Requests, error) {
	// Parse the log data into the exportUTXOEvent struct
	ev := &exportUTXOEvent{}
	err := SharedMemoryABI.UnpackIntoInterface(ev, "ExportUTXO", logData)
	if err != nil {
		return ids.ID{}, nil, fmt.Errorf("failed to unpack ExportUTXO event data: %w", err)
	}
	// Parse the topics data.
	// TODO: Improve this by using the ABI to unpack the topics.
	destinationChainID := ids.ID(topics[1])
	assetID := ids.ID(topics[2])

	addrs := make([]ids.ShortID, 0, len(ev.Addrs))
	for _, addr := range ev.Addrs {
		addrs = append(addrs, ids.ShortID(addr))
	}
	utxo := &avax.UTXO{
		// Derive unique UTXOID from txHash and log index
		UTXOID: avax.UTXOID{
			TxID:        ids.ID(txHash),
			OutputIndex: uint32(logIndex),
		},
		Asset: avax.Asset{ID: assetID},
		Out: &secp256k1fx.TransferOutput{
			Amt: ev.Amount,
			OutputOwners: secp256k1fx.OutputOwners{
				Locktime:  ev.Locktime,
				Threshold: uint32(ev.Threshold), // TODO make the actual type uint32 to correspond to this
				Addrs:     addrs,
			},
		},
	}

	utxoBytes, err := codec.Codec.Marshal(codec.CodecVersion, utxo)
	if err != nil {
		return ids.ID{}, nil, err
	}
	utxoID := utxo.InputID()
	elem := &atomic.Element{
		Key:   utxoID[:],
		Value: utxoBytes,
	}
	if out, ok := utxo.Out.(avax.Addressable); ok {
		elem.Traits = out.Addresses()
	}

	return destinationChainID, &atomic.Requests{
		PutRequests: []*atomic.Element{elem},
	}, nil
}

// importAVAXEvent includes the non-indexed fields of the ImportAVAX event.
// This allows us to unpack the log data into this struct. Note indexed
// fields appear in the topics array and are not included in the log data.
type importAVAXEvent struct {
	Amount uint64
	UtxoID [32]byte
}

func handleImportAVAX(snowCtx *snow.Context, txHash common.Hash, logIndex int, topics []common.Hash, logData []byte) (ids.ID, *atomic.Requests, error) {
	// Parse the log data into the importAVAXEvent struct
	ev := &importAVAXEvent{}
	err := SharedMemoryABI.UnpackIntoInterface(ev, "ImportAVAX", logData)
	if err != nil {
		return ids.ID{}, nil, fmt.Errorf("failed to unpack ImportAVAX event data: %w", err)
	}
	// Parse the topics data.
	// TODO: Improve this by using the ABI to unpack the topics.
	sourceChainID := ids.ID(topics[1])

	return sourceChainID, &atomic.Requests{
		RemoveRequests: [][]byte{ev.UtxoID[:]},
	}, nil
}

// importUTXOEvent includes the non-indexed fields of the ImportUTXO event.
// This allows us to unpack the log data into this struct. Note indexed
// fields appear in the topics array and are not included in the log data.
type importUTXOEvent struct {
	Amount uint64
	UtxoID [32]byte
}

func handleImportUTXO(snowCtx *snow.Context, txHash common.Hash, logIndex int, topics []common.Hash, logData []byte) (ids.ID, *atomic.Requests, error) {
	// Parse the log data into the importUTXOEvent struct
	ev := &importUTXOEvent{}
	err := SharedMemoryABI.UnpackIntoInterface(ev, "ImportUTXO", logData)
	if err != nil {
		return ids.ID{}, nil, fmt.Errorf("failed to unpack ImportUTXO event data: %w", err)
	}
	// Parse the topics data.
	// TODO: Improve this by using the ABI to unpack the topics.
	sourceChainID := ids.ID(topics[1])

	// TODO: should we verify the assetID is correct?
	// This is not strictly necessary because the shared memory contract
	// itself has calculated the assetID and included it in the event.

	return sourceChainID, &atomic.Requests{
		RemoveRequests: [][]byte{ev.UtxoID[:]},
	}, nil
}
