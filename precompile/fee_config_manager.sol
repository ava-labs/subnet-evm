// (c) 2022-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// SPDX-License-Identifier: MIT

pragma solidity >=0.8.0;

interface FeeConfigManagerInterface {
    // Set [addr] to have the admin role over the fee config manager list
    function setAdmin(address addr) external;

    // Set [addr] to be enabled on the fee config manager list
    function setEnabled(address addr) external;

    // Set [addr] to have no role over the fee config manager list
    function setNone(address addr) external;

    // Set fee config fields to contract storage
    function setFeeConfig(
        uint256 gasLimit,
        uint256 targetBlockRate,
        uint256 minBaseFee,
        uint256 targetGas,
        uint256 baseFeeChangeDenominator,
        uint256 minBlockGasCost,
        uint256 maxBlockGasCost,
        uint256 blockGasCostStep
    ) external;

    // Get fee config from the contract storage
    function getFeeConfig()
        external
        view
        returns (
            uint256 gasLimit,
            uint256 targetBlockRate,
            uint256 minBaseFee,
            uint256 targetGas,
            uint256 baseFeeChangeDenominator,
            uint256 minBlockGasCost,
            uint256 maxBlockGasCost,
            uint256 blockGasCostStep
        );
}
