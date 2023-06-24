// (c) 2022-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// SPDX-License-Identifier: MIT

pragma solidity >=0.8.0;
import "./IAllowList.sol";

interface ICounter is IAllowList {

    //state changing functions
    function IncrementByOne() external;
    function IncrementByX(uint64 x) external;

    // reads the state without changing
    function getCounter() external view returns (uint64);
}