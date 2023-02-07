// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"errors"

	"github.com/ava-labs/avalanchego/vms/components/avax"
)

const (
	// Max number of addresses that can be passed in as argument to GetUTXOs
	maxGetUTXOsAddrs = 1024

	// Max number of items allowed in a page
	maxPageSize uint64 = 1024
)

var (
	errNoAddresses = errors.New("no addresses provided")
)

// SnowmanAPI introduces snowman specific functionality to the evm
type SharedMemoryService struct {
	vm             *VM
	addressManager avax.AddressManager
}

func NewSharedMemoryService(vm *VM) *SharedMemoryService {
	return &SharedMemoryService{
		vm:             vm,
		addressManager: avax.NewAddressManager(vm.ctx),
	}
}

// // GetUTXOs gets all utxos for passed in addresses
// func (s *SharedMemoryService) GetUTXOs(_ *http.Request, args *api.GetUTXOsArgs, reply *api.GetUTXOsReply) error {
// 	s.vm.ctx.Log.Debug("AVM: GetUTXOs called",
// 		logging.UserStrings("addresses", args.Addresses),
// 	)

// 	if len(args.Addresses) == 0 {
// 		return errNoAddresses
// 	}
// 	if len(args.Addresses) > maxGetUTXOsAddrs {
// 		return fmt.Errorf("number of addresses given, %d, exceeds maximum, %d", len(args.Addresses), maxGetUTXOsAddrs)
// 	}

// 	var sourceChain ids.ID
// 	if args.SourceChain == "" {
// 		sourceChain = s.vm.ctx.ChainID
// 	} else {
// 		chainID, err := s.vm.ctx.BCLookup.Lookup(args.SourceChain)
// 		if err != nil {
// 			return fmt.Errorf("problem parsing source chainID %q: %w", args.SourceChain, err)
// 		}
// 		sourceChain = chainID
// 	}

// 	addrSet, err := avax.ParseServiceAddresses(s.addressManager, args.Addresses)
// 	if err != nil {
// 		return err
// 	}

// 	startAddr := ids.ShortEmpty
// 	startUTXO := ids.Empty
// 	if args.StartIndex.Address != "" || args.StartIndex.UTXO != "" {
// 		startAddr, err = avax.ParseServiceAddress(s.addressManager, args.StartIndex.Address)
// 		if err != nil {
// 			return fmt.Errorf("couldn't parse start index address %q: %w", args.StartIndex.Address, err)
// 		}
// 		startUTXO, err = ids.FromString(args.StartIndex.UTXO)
// 		if err != nil {
// 			return fmt.Errorf("couldn't parse start index utxo: %w", err)
// 		}
// 	}

// 	var (
// 		utxos     []*avax.UTXO
// 		endAddr   ids.ShortID
// 		endUTXOID ids.ID
// 	)
// 	limit := int(args.Limit)
// 	if limit <= 0 || int(maxPageSize) < limit {
// 		limit = int(maxPageSize)
// 	}
// 	if sourceChain == s.vm.ctx.ChainID {
// 		utxos, endAddr, endUTXOID, err = avax.GetPaginatedUTXOs(
// 			s.vm.state,
// 			addrSet,
// 			startAddr,
// 			startUTXO,
// 			limit,
// 		)
// 	} else {
// 		utxos, endAddr, endUTXOID, err = s.vm.GetAtomicUTXOs(
// 			sourceChain,
// 			addrSet,
// 			startAddr,
// 			startUTXO,
// 			limit,
// 		)
// 	}
// 	if err != nil {
// 		return fmt.Errorf("problem retrieving UTXOs: %w", err)
// 	}

// 	reply.UTXOs = make([]string, len(utxos))
// 	codec := s.vm.parser.Codec()
// 	for i, utxo := range utxos {
// 		b, err := codec.Marshal(txs.CodecVersion, utxo)
// 		if err != nil {
// 			return fmt.Errorf("problem marshalling UTXO: %w", err)
// 		}
// 		reply.UTXOs[i], err = formatting.Encode(args.Encoding, b)
// 		if err != nil {
// 			return fmt.Errorf("couldn't encode UTXO %s as string: %w", utxo.InputID(), err)
// 		}
// 	}

// 	endAddress := common.Address(endAddr)

// 	reply.EndIndex.Address = endAddress.Hex()
// 	reply.EndIndex.UTXO = endUTXOID.String()
// 	reply.NumFetched = json.Uint64(len(utxos))
// 	reply.Encoding = args.Encoding
// 	return nil
// }
