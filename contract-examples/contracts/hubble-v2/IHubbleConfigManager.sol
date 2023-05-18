pragma solidity 0.8.9;

import "../IAllowList.sol";

interface IHubbleConfigManager is IAllowList{
  //getSpreadRatioThreshold returns the spreadRatioThreshold stored in evm state
  function getSpreadRatioThreshold() external view returns (uint256 spreadRatioThreshold);

  //setSpreadRatioThreshold stores the spreadRatioThreshold in evm state
  function setSpreadRatioThreshold(uint256 spreadRatioThreshold) external;

  //getMinAllowableMargin returns the minAllowableMargin stored in evm state
  function getMinAllowableMargin() external view returns (uint256 minAllowableMargin);

  //setMinAllowableMargin stores the minAllowableMargin in evm state
  function setMinAllowableMargin(uint256 minAllowableMargin) external;

  //getMaintenanceMargin return the maintenanceMargin stored in evm state
  function getMaintenanceMargin() external view returns(uint256 maintenanceMargin);

  //setMaintenanceMargin stores the maintenanceMargin stored in evm state
  function setMaintenanceMargin(uint256 maintenanceMargin) external;

  //getMaxLiquidationRatio return the maxLiquidationRatio stored in evm state
  function getMaxLiquidationRatio() external view returns(uint256 maxLiquidationRatio);

  //setMaxLiquidationRatio stores the maxLiquidationRatio in evm state
  function setMaxLiquidationRatio(uint256 maxLiquidationRatio) external;

  //getMinSizeRequirement returns minSizeRequirement stored in evm state
  function getMinSizeRequirement() external view returns(uint256 minSizeRequirement);

  //setMinSizeRequirement stores the minSizeRequirement in evm state
  function setMinSizeRequirement(uint256 minSizeRequirement) external;
}