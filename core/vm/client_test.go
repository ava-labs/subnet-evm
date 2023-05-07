package vm

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/params"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/ethereum/go-ethereum/common"

	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
)

func (suite *KeeperTestSuite) TestCreateClient() {
	cases := []struct {
		msg            string
		clientState    exported.ClientState
		consensusState exported.ConsensusState
		expPass        bool
	}{
		{
			"success: 07-tendermint client type supported",
			ibctm.NewClientState(testChainID, ibctm.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath),
			suite.consensusState,
			true,
		},
	}

	statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	statedb.Finalise(true)
	vmctx := BlockContext{
		CanTransfer: func(StateDB, common.Address, *big.Int) bool { return true },
		Transfer:    func(StateDB, common.Address, common.Address, *big.Int) {},
	}
	vmenv := NewEVM(vmctx, TxContext{}, statedb, params.TestChainConfig, Config{ExtraEips: []int{2200}})

	for _, tc := range cases {
		test_precompiles := &createClient{}
		var input []byte

		// clientStateLen     - clientState
		clientState, ok := tc.clientState.(*ibctm.ClientState)
		if !ok {
			suite.Require().NoError(errors.New("convert to proto.Message failer"))
		}
		clientStateByte, err := clientState.Marshal()

		suite.Require().NoError(err)

		// 8 byte             - clientStateLen
		clientStateLen := make([]byte, 8)
		binary.BigEndian.PutUint64(clientStateLen, uint64(len(clientStateByte)))

		input = append(input, clientStateLen...)
		input = append(input, clientStateByte...)

		// consensusStateLen  - consensusState

		consensusState, ok := tc.consensusState.(*ibctm.ConsensusState)
		if !ok {
			suite.Require().NoError(errors.New("convert to proto.Message failer"))
		}
		consensusStateByte, err := consensusState.Marshal()
		suite.Require().NoError(err)

		// 8 byte             - consensusStateLen
		consensusStateLen := make([]byte, 8)
		binary.BigEndian.PutUint64(consensusStateLen, uint64(len(consensusStateByte)))

		input = append(input, consensusStateLen...)
		input = append(input, consensusStateByte...)

		output, err := test_precompiles.Run(vmenv, input)
		suite.Require().NoError(err)
		suite.Equal(string(output), fmt.Sprintf("%s-%d", tc.clientState.ClientType(), 0), "clientID bad formatting")
	}
}

func (suite *KeeperTestSuite) TestUpdateClientTendermint() {
	var (
		path         *ibctesting.Path
		updateHeader *ibctm.Header
	)

	// Must create header creation functions since suite.header gets recreated on each test case
	createFutureUpdateFn := func(trustedHeight clienttypes.Height) *ibctm.Header {
		header, err := suite.chainA.ConstructUpdateTMClientHeaderWithTrustedHeight(path.EndpointB.Chain, path.EndpointA.ClientID, trustedHeight)
		suite.Require().NoError(err)
		return header
	}
	createPastUpdateFn := func(fillHeight, trustedHeight clienttypes.Height) *ibctm.Header {
		consState, found := suite.chainA.App.GetIBCKeeper().ClientKeeper.GetClientConsensusState(suite.chainA.GetContext(), path.EndpointA.ClientID, trustedHeight)
		suite.Require().True(found)

		return suite.chainB.CreateTMClientHeader(suite.chainB.ChainID, int64(fillHeight.RevisionHeight), trustedHeight, consState.(*ibctm.ConsensusState).Timestamp.Add(time.Second*5),
			suite.chainB.Vals, suite.chainB.Vals, suite.chainB.Vals, suite.chainB.Signers)
	}

	statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	statedb.Finalise(true)
	vmctx := BlockContext{
		CanTransfer: func(StateDB, common.Address, *big.Int) bool { return true },
		Transfer:    func(StateDB, common.Address, common.Address, *big.Int) {},
	}
	vmenv := NewEVM(vmctx, TxContext{}, statedb, params.TestChainConfig, Config{ExtraEips: []int{2200}})

	cases := []struct {
		name      string
		malleate  func()
		expPass   bool
		expFreeze bool
	}{
		{"valid past update", func() {
			statedb, _ = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
			statedb.Finalise(true)
			vmenv = NewEVM(vmctx, TxContext{}, statedb, params.TestChainConfig, Config{ExtraEips: []int{2200}})

			clientState := path.EndpointA.GetClientState()
			trustedHeight := clientState.GetLatestHeight().(clienttypes.Height)

			currHeight := suite.chainB.CurrentHeader.Height
			fillHeight := clienttypes.NewHeight(clientState.GetLatestHeight().GetRevisionNumber(), uint64(currHeight))

			// commit a couple blocks to allow client to fill in gaps
			suite.coordinator.CommitBlock(suite.chainB) // this height is not filled in yet
			suite.coordinator.CommitBlock(suite.chainB) // this height is filled in by the update below

			err := path.EndpointA.UpdateClient()
			suite.Require().NoError(err)

			clientStateByte, _ := clientState.(*ibctm.ClientState).Marshal()
			clientStatePath := fmt.Sprintf("clients/%s/clientState", path.EndpointA.ClientID)
			vmenv.StateDB.SetPrecompileState(
				common.BytesToAddress([]byte(clientStatePath)),
				clientStateByte,
			)
			// store previous consensus state
			prevConsState := &ibctm.ConsensusState{
				Timestamp:          suite.past,
				NextValidatorsHash: suite.chainB.Vals.Hash(),
			}
			prevConsStateByte, _ := prevConsState.Marshal()
			consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", path.EndpointA.ClientID, clientState.GetLatestHeight())
			vmenv.StateDB.SetPrecompileState(
				common.BytesToAddress([]byte(consensusStatePath)),
				prevConsStateByte,
			)
			// ensure fill height not set
			_, found := suite.chainA.App.GetIBCKeeper().ClientKeeper.GetClientConsensusState(suite.chainA.GetContext(), path.EndpointA.ClientID, fillHeight)
			suite.Require().False(found)

			// updateHeader will fill in consensus state between prevConsState and suite.consState
			// clientState should not be updated
			updateHeader = createPastUpdateFn(fillHeight, trustedHeight)
		}, true, false},
		{"misbehaviour detection: conflicting header", func() {
			// statedb, _ = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
			// statedb.Finalise(true)
			// vmenv = NewEVM(vmctx, TxContext{}, statedb, params.TestChainConfig, Config{ExtraEips: []int{2200}})

			clientID := path.EndpointA.ClientID

			height1 := clienttypes.NewHeight(1, 1)
			// store previous consensus state
			prevConsState := &ibctm.ConsensusState{
				Timestamp:          suite.past,
				NextValidatorsHash: suite.chainB.Vals.Hash(),
			}
			suite.chainA.App.GetIBCKeeper().ClientKeeper.SetClientConsensusState(suite.chainA.GetContext(), clientID, height1, prevConsState)

			height5 := clienttypes.NewHeight(1, 5)
			// store next consensus state to check that trustedHeight does not need to be hightest consensus state before header height
			nextConsState := &ibctm.ConsensusState{
				Timestamp:          suite.past.Add(time.Minute),
				NextValidatorsHash: suite.chainB.Vals.Hash(),
			}
			suite.chainA.App.GetIBCKeeper().ClientKeeper.SetClientConsensusState(suite.chainA.GetContext(), clientID, height5, nextConsState)

			height3 := clienttypes.NewHeight(1, 3)
			// updateHeader will fill in consensus state between prevConsState and suite.consState
			// clientState should not be updated
			updateHeader = createPastUpdateFn(height3, height1)
			// set conflicting consensus state in store to create misbehaviour scenario
			conflictConsState := updateHeader.ConsensusState()
			conflictConsState.Root = commitmenttypes.NewMerkleRoot([]byte("conflicting apphash"))
			suite.chainA.App.GetIBCKeeper().ClientKeeper.SetClientConsensusState(suite.chainA.GetContext(), clientID, updateHeader.GetHeight(), conflictConsState)
		}, true, true},
		{"client state not found", func() {
			updateHeader = createFutureUpdateFn(path.EndpointA.GetClientState().GetLatestHeight().(clienttypes.Height))
			path.EndpointA.ClientID = ibctesting.InvalidID
		}, false, false},
		{"consensus state not found", func() {
			statedb, _ = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
			statedb.Finalise(true)
			vmenv = NewEVM(vmctx, TxContext{}, statedb, params.TestChainConfig, Config{ExtraEips: []int{2200}})

			clientState := path.EndpointA.GetClientState()
			tmClient, ok := clientState.(*ibctm.ClientState)
			suite.Require().True(ok)
			tmClient.LatestHeight = tmClient.LatestHeight.Increment().(clienttypes.Height)
			clientStateByte, _ := clientState.(*ibctm.ClientState).Marshal()
			clientStatePath := fmt.Sprintf("clients/%s/clientState", path.EndpointA.ClientID)
			vmenv.StateDB.SetPrecompileState(
				common.BytesToAddress([]byte(clientStatePath)),
				clientStateByte,
			)
			updateHeader = createFutureUpdateFn(clientState.GetLatestHeight().(clienttypes.Height))
		}, false, false},
		{"client is not active", func() {
			statedb, _ = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
			statedb.Finalise(true)
			vmenv = NewEVM(vmctx, TxContext{}, statedb, params.TestChainConfig, Config{ExtraEips: []int{2200}})

			clientState := path.EndpointA.GetClientState().(*ibctm.ClientState)
			clientState.FrozenHeight = clienttypes.NewHeight(1, 1)
			clientStateByte, _ := clientState.Marshal()
			clientStatePath := fmt.Sprintf("clients/%s/clientState", path.EndpointA.ClientID)
			vmenv.StateDB.SetPrecompileState(
				common.BytesToAddress([]byte(clientStatePath)),
				clientStateByte,
			)
			updateHeader = createFutureUpdateFn(clientState.GetLatestHeight().(clienttypes.Height))
		}, false, false},
		{"invalid header", func() {
			updateHeader = createFutureUpdateFn(path.EndpointA.GetClientState().GetLatestHeight().(clienttypes.Height))
			updateHeader.TrustedHeight = updateHeader.TrustedHeight.Increment().(clienttypes.Height)
		}, false, false},
	}

	for _, tc := range cases {
		tc := tc
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest()
			path = ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			tc.malleate()

			var input []byte

			ClientIDByte := []byte(path.EndpointA.ClientID)
			ClientIDByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(ClientIDByteLen, uint64(len(ClientIDByte)))

			input = append(input, ClientIDByteLen...)
			input = append(input, ClientIDByte...)

			clientMessageByte, err := updateHeader.Marshal()
			suite.Require().NoError(err)
			clientMessageLen := make([]byte, 8)
			binary.BigEndian.PutUint64(clientMessageLen, uint64(len(clientMessageByte)))

			input = append(input, clientMessageLen...)
			input = append(input, clientMessageByte...)

			test_precompiles := &updateClient{}

			_, err = test_precompiles.Run(vmenv, input)

			if tc.expPass {
				suite.Require().NoError(err, err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestUpgradeClient() {
	var (
		path                                        *ibctesting.Path
		upgradedClient                              exported.ClientState
		upgradedConsState                           exported.ConsensusState
		lastHeight                                  exported.Height
		proofUpgradedClient, proofUpgradedConsState []byte
	)

	statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	statedb.Finalise(true)
	vmctx := BlockContext{
		CanTransfer: func(StateDB, common.Address, *big.Int) bool { return true },
		Transfer:    func(StateDB, common.Address, common.Address, *big.Int) {},
	}
	vmenv := NewEVM(vmctx, TxContext{}, statedb, params.TestChainConfig, Config{ExtraEips: []int{2200}})

	testCases := []struct {
		name    string
		setup   func()
		expPass bool
	}{
		{
			name: "successful upgrade",
			setup: func() {
				// last Height is at next block
				lastHeight = clienttypes.NewHeight(1, uint64(suite.chainB.GetContext().BlockHeight()+1))

				// commit upgrade store changes and update clients
				suite.coordinator.CommitBlock(suite.chainB)
				err := path.EndpointA.UpdateClient()
				suite.Require().NoError(err)

				cs, found := suite.chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(suite.chainA.GetContext(), path.EndpointA.ClientID)
				suite.Require().True(found)

				proofUpgradedClient, _ = suite.chainB.QueryUpgradeProof(upgradetypes.UpgradedClientKey(int64(lastHeight.GetRevisionHeight())), cs.GetLatestHeight().GetRevisionHeight())
				proofUpgradedConsState, _ = suite.chainB.QueryUpgradeProof(upgradetypes.UpgradedConsStateKey(int64(lastHeight.GetRevisionHeight())), cs.GetLatestHeight().GetRevisionHeight())
			},
			expPass: true,
		},
		{
			name: "client state not found",
			setup: func() {
				// last Height is at next block
				lastHeight = clienttypes.NewHeight(1, uint64(suite.chainB.GetContext().BlockHeight()+1))

				// commit upgrade store changes and update clients
				suite.coordinator.CommitBlock(suite.chainB)
				err := path.EndpointA.UpdateClient()
				suite.Require().NoError(err)

				cs, found := suite.chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(suite.chainA.GetContext(), path.EndpointA.ClientID)
				suite.Require().True(found)

				proofUpgradedClient, _ = suite.chainB.QueryUpgradeProof(upgradetypes.UpgradedClientKey(int64(lastHeight.GetRevisionHeight())), cs.GetLatestHeight().GetRevisionHeight())
				proofUpgradedConsState, _ = suite.chainB.QueryUpgradeProof(upgradetypes.UpgradedConsStateKey(int64(lastHeight.GetRevisionHeight())), cs.GetLatestHeight().GetRevisionHeight())

				path.EndpointA.ClientID = "wrongclientid"
			},
			expPass: false,
		},
		{
			name: "client state is not active",
			setup: func() {
				// client is frozen

				// last Height is at next block
				lastHeight = clienttypes.NewHeight(1, uint64(suite.chainB.GetContext().BlockHeight()+1))

				// commit upgrade store changes and update clients
				suite.coordinator.CommitBlock(suite.chainB)
				err := path.EndpointA.UpdateClient()
				suite.Require().NoError(err)

				cs, found := suite.chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(suite.chainA.GetContext(), path.EndpointA.ClientID)
				suite.Require().True(found)

				proofUpgradedClient, _ = suite.chainB.QueryUpgradeProof(upgradetypes.UpgradedClientKey(int64(lastHeight.GetRevisionHeight())), cs.GetLatestHeight().GetRevisionHeight())
				proofUpgradedConsState, _ = suite.chainB.QueryUpgradeProof(upgradetypes.UpgradedConsStateKey(int64(lastHeight.GetRevisionHeight())), cs.GetLatestHeight().GetRevisionHeight())

				// set frozen client in store
				tmClient, ok := cs.(*ibctm.ClientState)
				suite.Require().True(ok)
				tmClient.FrozenHeight = clienttypes.NewHeight(1, 1)
				suite.chainA.App.GetIBCKeeper().ClientKeeper.SetClientState(suite.chainA.GetContext(), path.EndpointA.ClientID, tmClient)
			},
			expPass: false,
		},
		{
			name: "tendermint client VerifyUpgrade fails",
			setup: func() {
				// last Height is at next block
				lastHeight = clienttypes.NewHeight(1, uint64(suite.chainB.GetContext().BlockHeight()+1))

				// change upgradedClient client-specified parameters
				tmClient := upgradedClient.(*ibctm.ClientState)
				tmClient.ChainId = "wrongchainID"
				upgradedClient = tmClient

				suite.coordinator.CommitBlock(suite.chainB)
				err := path.EndpointA.UpdateClient()
				suite.Require().NoError(err)

				cs, found := suite.chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(suite.chainA.GetContext(), path.EndpointA.ClientID)
				suite.Require().True(found)

				proofUpgradedClient, _ = suite.chainB.QueryUpgradeProof(upgradetypes.UpgradedClientKey(int64(lastHeight.GetRevisionHeight())), cs.GetLatestHeight().GetRevisionHeight())
				proofUpgradedConsState, _ = suite.chainB.QueryUpgradeProof(upgradetypes.UpgradedConsStateKey(int64(lastHeight.GetRevisionHeight())), cs.GetLatestHeight().GetRevisionHeight())
			},
			expPass: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		path = ibctesting.NewPath(suite.chainA, suite.chainB)
		suite.coordinator.SetupClients(path)

		clientState := path.EndpointA.GetClientState().(*ibctm.ClientState)
		revisionNumber := clienttypes.ParseChainID(clientState.ChainId)

		newChainID, err := clienttypes.SetRevisionNumber(clientState.ChainId, revisionNumber+1)
		suite.Require().NoError(err)

		upgradedClient = ibctm.NewClientState(newChainID, ibctm.DefaultTrustLevel, trustingPeriod, ubdPeriod+trustingPeriod, maxClockDrift, clienttypes.NewHeight(revisionNumber+1, clientState.GetLatestHeight().GetRevisionHeight()+1), commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath)
		upgradedClient = upgradedClient.ZeroCustomFields()

		upgradedConsState = &ibctm.ConsensusState{
			NextValidatorsHash: []byte("nextValsHash"),
		}

		tc.setup()

		clientStateByte, _ := upgradedClient.(*ibctm.ClientState).Marshal()
		clientStatePath := fmt.Sprintf("clients/%s/clientState", clientState.ChainId)
		vmenv.StateDB.SetPrecompileState(
			common.BytesToAddress([]byte(clientStatePath)),
			clientStateByte,
		)
		consStateByte, _ := upgradedConsState.(*ibctm.ConsensusState).Marshal()
		consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", clientState.ChainId, lastHeight)
		vmenv.StateDB.SetPrecompileState(
			common.BytesToAddress([]byte(consensusStatePath)),
			consStateByte,
		)

		var input []byte

		ClientIDByte := []byte(path.EndpointA.ClientID)
		ClientIDByteLen := make([]byte, 8)
		binary.BigEndian.PutUint64(ClientIDByteLen, uint64(len(ClientIDByte)))

		input = append(input, ClientIDByteLen...)
		input = append(input, ClientIDByte...)

		upgradedClientByte, err := upgradedClient.(*ibctm.ClientState).Marshal()
		suite.Require().NoError(err)
		upgradedClientLen := make([]byte, 8)
		binary.BigEndian.PutUint64(upgradedClientLen, uint64(len(upgradedClientByte)))

		input = append(input, upgradedClientLen...)
		input = append(input, upgradedClientByte...)

		upgradedConsStateByte, err := upgradedConsState.(*ibctm.ConsensusState).Marshal()
		suite.Require().NoError(err)
		upgradedConsStateLen := make([]byte, 8)
		binary.BigEndian.PutUint64(upgradedConsStateLen, uint64(len(upgradedConsStateByte)))

		input = append(input, upgradedConsStateLen...)
		input = append(input, upgradedConsStateByte...)

		proofUpgradedClientLen := make([]byte, 8)
		binary.BigEndian.PutUint64(proofUpgradedClientLen, uint64(len(proofUpgradedClient)))

		input = append(input, proofUpgradedClientLen...)
		input = append(input, proofUpgradedClient...)

		proofUpgradedConsStateLen := make([]byte, 8)
		binary.BigEndian.PutUint64(proofUpgradedConsStateLen, uint64(len(proofUpgradedConsState)))

		input = append(input, proofUpgradedConsStateLen...)
		input = append(input, proofUpgradedConsState...)

		test_precompiles := &upgradeClient{}

		_, err = test_precompiles.Run(vmenv, input)

		if tc.expPass {
			suite.Require().NoError(err, "verify upgrade failed on valid case: %s", tc.name)
		} else {
			suite.Require().Error(err, "verify upgrade passed on invalid case: %s", tc.name)
		}
	}
}
