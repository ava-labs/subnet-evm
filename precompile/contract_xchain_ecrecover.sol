// (c) 2022-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// SPDX-License-Identifier: MIT

pragma solidity >=0.8.0;

interface XChainECRecover {
    
    function xChainECRecover(string memory input) external view returns(string memory);

    function getXChainECRecover(string memory input) external view returns (string memory);
}