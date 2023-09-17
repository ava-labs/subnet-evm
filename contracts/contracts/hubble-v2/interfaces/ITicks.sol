// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface ITicks {
    function getPrevTick(address amm, bool isBid, uint tick) external view returns (uint prevTick);
    function sampleImpactBid(address amm) external view returns (uint impactBid);
    function sampleImpactAsk(address amm) external view returns (uint impactAsk);
    function getQuote(address amm, int256 baseAssetQuantity) external view returns (uint256 rate);
    function getBaseQuote(address amm, int256 quoteQuantity) external view returns (uint256 rate);
}
