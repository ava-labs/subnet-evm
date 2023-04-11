//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
pragma experimental ABIEncoderV2;

import "@openzeppelin/contracts/access/Ownable.sol";
import "ds-test/src/test.sol";
import "./AllowList.sol";
import "./IFeeManager.sol";

address constant FEE_MANAGER_ADDRESS = 0x0200000000000000000000000000000000000003;

// ExampleFeeManager shows how FeeManager precompile can be used in a smart contract
// All methods of [allowList] can be directly called. There are example calls as tasks in hardhat.config.ts file.
contract ExampleFeeManager is AllowList {
  // Precompiled Fee Manager Contract Address
  IFeeManager feeManager = IFeeManager(FEE_MANAGER_ADDRESS);

  struct FeeConfig {
    uint256 gasLimit;
    uint256 targetBlockRate;
    uint256 minBaseFee;
    uint256 targetGas;
    uint256 baseFeeChangeDenominator;
    uint256 minBlockGasCost;
    uint256 maxBlockGasCost;
    uint256 blockGasCostStep;
  }

  constructor() AllowList(FEE_MANAGER_ADDRESS) {}

  function enableWAGMIFees() public onlyEnabled {
    feeManager.setFeeConfig(
      20_000_000, // gasLimit
      2, // targetBlockRate
      1_000_000_000, // minBaseFee
      100_000_000, // targetGas
      48, // baseFeeChangeDenominator
      0, // minBlockGasCost
      10_000_000, // maxBlockGasCost
      500_000 // blockGasCostStep
    );
  }

  function enableCChainFees() public onlyEnabled {
    feeManager.setFeeConfig(
      8_000_000, // gasLimit
      2, // targetBlockRate
      25_000_000_000, // minBaseFee
      15_000_000, // targetGas
      36, // baseFeeChangeDenominator
      0, // minBlockGasCost
      1_000_000, // maxBlockGasCost
      200_000 // blockGasCostStep
    );
  }

  function enableCustomFees(FeeConfig memory config) public onlyEnabled {
    feeManager.setFeeConfig(
      config.gasLimit,
      config.targetBlockRate,
      config.minBaseFee,
      config.targetGas,
      config.baseFeeChangeDenominator,
      config.minBlockGasCost,
      config.maxBlockGasCost,
      config.blockGasCostStep
    );
  }

  function getCurrentFeeConfig() public view returns (FeeConfig memory) {
    FeeConfig memory config;
    (
      config.gasLimit,
      config.targetBlockRate,
      config.minBaseFee,
      config.targetGas,
      config.baseFeeChangeDenominator,
      config.minBlockGasCost,
      config.maxBlockGasCost,
      config.blockGasCostStep
    ) = feeManager.getFeeConfig();
    return config;
  }

  function getFeeConfigLastChangedAt() public view returns (uint256) {
    return feeManager.getFeeConfigLastChangedAt();
  }
}

contract ExampleFeeManagerTest is DSTest {
  uint256 testNumber;

  function setUp() public {
    // noop
  }

  function test_addContractDeployerAsOwner() public {
    ExampleFeeManager manager = new ExampleFeeManager();
    assertEq(address(this), manager.owner());
  }

  function test_enableWAGMIFeesFailure() public {
    ExampleFeeManager example = new ExampleFeeManager();

    IFeeManager manager = IFeeManager(FEE_MANAGER_ADDRESS);

    // TODO: make roles const (role: None = 0)
    assertEq(manager.readAllowList(address(example)), 0);

    try example.enableWAGMIFees() {
      assertTrue(false, "enableWAGMIFees should fail");
    } catch {}
  }

  function test_addContractToManagerList() public {
    ExampleFeeManager example = new ExampleFeeManager();

    address exampleAddress = address(example);
    address thisAddress = address(this);

    IFeeManager manager = IFeeManager(FEE_MANAGER_ADDRESS);

    // TODO: make this a const (role: ADMIN = 2)
    assertEq(manager.readAllowList(thisAddress), 2);
    assertEq(manager.readAllowList(exampleAddress), 0);

    manager.setEnabled(exampleAddress);

    // TODO: make this a const (role: ENABLED = 1)
    assertEq(manager.readAllowList(exampleAddress), 1);
  }

  function test_enableCustomFees() public {
    ExampleFeeManager example = new ExampleFeeManager();
    address exampleAddress = address(example);

    IFeeManager manager = IFeeManager(FEE_MANAGER_ADDRESS);

    manager.setEnabled(exampleAddress);

    ExampleFeeManager.FeeConfig memory config = example.getCurrentFeeConfig();

    uint256 newGasLimit = config.gasLimit + 10;
    uint256 newMinBaseFee = config.minBaseFee + 10;

    config.gasLimit = newGasLimit;
    config.minBaseFee = newMinBaseFee;

    example.enableCustomFees(config);

    ExampleFeeManager.FeeConfig memory newConfig = example.getCurrentFeeConfig();

    assertEq(newConfig.gasLimit, newGasLimit);
    assertEq(newConfig.minBaseFee, newMinBaseFee);

    assertEq(example.getFeeConfigLastChangedAt(), block.number);
  }
}
