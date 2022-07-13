// (c) 2022-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// SPDX-License-Identifier: MIT

pragma solidity >=0.8.0;

interface HelloWorldInterface {
    function sayHello() external;

    // setRecipient
    function setReceipient(string calldata recipient) external;
}

// TODO add it to ChainConfig in the same way as for the other precompiles
// if the state of the stateful precompile is intended to impact the EVM or the node's behavior, you can also add a getter function to retrieve the state
// and add this wherever necessary in the codebase ie. TxAllowList in core/tx_pool.go and core/state_transition.go. This enables users to modify the behavior
// of the EVM and the entire node at a deeper level if needed and they have considered the tradeoffs.

// add four simple unit tests in core/stateful_precompile_test.go
// add a test in plugin/evm/vm_test.go if it modifies any core behavior beyond the simple state transition ex. TxAllowList
// set recipient read only
// set recipient no read only
// sayHello normal
// sayHello after setting the recipient
// add hardhat task tests in smart contract examples (Connor has written these so he's better to ask than I am on this)