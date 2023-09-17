// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

interface IClearingHouse {
    enum OrderExecutionMode {
        Taker,
        Maker,
        SameBlock, // not used
        Liquidation
    }

    /**
     * @param ammIndex Market id to place the order. In Hubble, market ids are sequential and start from 0
     * @param trader Address of the trader
     * @param mode Whether to be executed as a Maker, Taker or Liquidation
    */
    struct Instruction {
        uint256 ammIndex;
        address trader;
        bytes32 orderHash;
        OrderExecutionMode mode;
    }

    enum Mode { Maintenance_Margin, Min_Allowable_Margin }
}
