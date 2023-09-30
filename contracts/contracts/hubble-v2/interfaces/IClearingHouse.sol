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

    event PositionModified(address indexed trader, uint indexed idx, int256 baseAsset, uint price, int256 realizedPnl, int256 size, uint256 openNotional, int256 fee, OrderExecutionMode mode, uint256 timestamp);
    event PositionLiquidated(address indexed trader, uint indexed idx, int256 baseAsset, uint256 price, int256 realizedPnl, int256 size, uint256 openNotional, int256 fee, uint256 timestamp);
    event MarketAdded(uint indexed idx, address indexed amm);
    event ReferralBonusAdded(address indexed referrer, uint referralBonus);
    event FundingPaid(address indexed trader, uint indexed idx, int256 takerFundingPayment, int256 cumulativePremiumFraction);
    event FundingRateUpdated(uint indexed idx, int256 premiumFraction, uint256 underlyingPrice, int256 cumulativePremiumFraction, uint256 nextFundingTime, uint256 timestamp, uint256 blockNumber);

    function openComplementaryPositions(
        Instruction[2] memory orders,
        int256 fillAmount,
        uint fulfillPrice
    )  external returns (uint256 openInterest);

    function settleFunding() external;
    function getTotalNotionalPositionAndUnrealizedPnl(address trader, int256 margin, Mode mode)
        external
        view
        returns(uint256 notionalPosition, int256 unrealizedPnl);
    function isAboveMaintenanceMargin(address trader) external view returns(bool);
    function assertMarginRequirement(address trader) external view;
    function updatePositions(address trader) external;
    function getTotalFunding(address trader) external view returns(int256 totalFunding);
    function getAmmsLength() external view returns(uint);
    // function amms(uint idx) external view returns(IAMM);
    function maintenanceMargin() external view returns(int256);
    function minAllowableMargin() external view returns(int256);
    function takerFee() external view returns(int256);
    function makerFee() external view returns(int256);
    function liquidationPenalty() external view returns(uint256);
    function getNotionalPositionAndMargin(address trader, bool includeFundingPayments, Mode mode)
        external
        view
        returns(uint256 notionalPosition, int256 margin);
    function getNotionalPositionAndMarginVanilla(address trader, bool includeFundingPayments, Mode mode)
        external
        view
        returns(uint256 notionalPosition, int256 margin);
    function liquidate(
        Instruction calldata instruction,
        int256 liquidationAmount,
        uint price,
        address trader
    ) external returns (uint256 openInterest);
    function feeSink() external view returns(address);
    function calcMarginFraction(address trader, bool includeFundingPayments, Mode mode) external view returns(int256);
    function getUnderlyingPrice() external view returns(uint[] memory prices);
    function orderBook() external view returns(address);
}
