//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
pragma experimental ABIEncoderV2;

import "@openzeppelin/contracts/access/Ownable.sol";
import "ds-test/src/test.sol";
import "./AllowList.sol";
import "./IFeeManager.sol";

address constant FEE_MANAGER_ADDRESS = 0x0200000000000000000000000000000000000003;

uint constant WAGMI_GAS_LIMIT = 20_000_000;
uint constant WAGMI_TARGET_BLOCK_RATE = 2;
uint constant WAGMI_MIN_BASE_FEE = 1_000_000_000;
uint constant WAGMI_TARGET_GAS = 100_000_000;
uint constant WAGMI_BASE_FEE_CHANGE_DENOMINATOR = 48;
uint constant WAGMI_MIN_BLOCK_GAS_COST = 0;
uint constant WAGMI_MAX_BLOCK_GAS_COST = 10_000_000;
uint constant WAGMI_BLOCK_GAS_COST_STEP = 500_000;

uint constant CCHAIN_GAS_LIMIT = 8_000_000;
uint constant CCHAIN_TARGET_BLOCK_RATE = 2;
uint constant CCHAIN_MIN_BASE_FEE = 25_000_000_000;
uint constant CCHAIN_TARGET_GAS = 15_000_000;
uint constant CCHAIN_BASE_FEE_CHANGE_DENOMINATOR = 36;
uint constant CCHAIN_MIN_BLOCK_GAS_COST = 0;
uint constant CCHAIN_MAX_BLOCK_GAS_COST = 1_000_000;
uint constant CCHAIN_BLOCK_GAS_COST_STEP = 100_000;

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

// ExampleFeeManager shows how FeeManager precompile can be used in a smart contract
// All methods of [allowList] can be directly called. There are example calls as tasks in hardhat.config.ts file.
contract ExampleFeeManager is AllowList {
  IFeeManager feeManager = IFeeManager(FEE_MANAGER_ADDRESS);

  constructor() AllowList(FEE_MANAGER_ADDRESS) {}

  function enableWAGMIFees() public onlyEnabled {
    feeManager.setFeeConfig(
      WAGMI_GAS_LIMIT,
      WAGMI_TARGET_BLOCK_RATE,
      WAGMI_MIN_BASE_FEE,
      WAGMI_TARGET_GAS,
      WAGMI_BASE_FEE_CHANGE_DENOMINATOR,
      WAGMI_MIN_BLOCK_GAS_COST,
      WAGMI_MAX_BLOCK_GAS_COST,
      WAGMI_BLOCK_GAS_COST_STEP
    );
  }

  function enableCChainFees() public onlyEnabled {
    feeManager.setFeeConfig(
      CCHAIN_GAS_LIMIT,
      CCHAIN_TARGET_BLOCK_RATE,
      CCHAIN_MIN_BASE_FEE,
      CCHAIN_TARGET_GAS,
      CCHAIN_BASE_FEE_CHANGE_DENOMINATOR,
      CCHAIN_MIN_BLOCK_GAS_COST,
      CCHAIN_MAX_BLOCK_GAS_COST,
      CCHAIN_BLOCK_GAS_COST_STEP
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
    } catch {} // TODO should match on an error to make sure that this is failing in the way that's expected
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

  function test_changeFees() public {
    ExampleFeeManager example = new ExampleFeeManager();
    address exampleAddress = address(example);

    IFeeManager manager = IFeeManager(FEE_MANAGER_ADDRESS);

    manager.setEnabled(exampleAddress);

    FeeConfig memory config = example.getCurrentFeeConfig();

    FeeConfig memory newFeeConfig = FeeConfig({
      gasLimit: CCHAIN_GAS_LIMIT,
      targetBlockRate: CCHAIN_TARGET_BLOCK_RATE,
      minBaseFee: CCHAIN_MIN_BASE_FEE,
      targetGas: CCHAIN_TARGET_GAS,
      baseFeeChangeDenominator: CCHAIN_BASE_FEE_CHANGE_DENOMINATOR,
      minBlockGasCost: CCHAIN_MIN_BLOCK_GAS_COST,
      maxBlockGasCost: CCHAIN_MAX_BLOCK_GAS_COST,
      blockGasCostStep: CCHAIN_BLOCK_GAS_COST_STEP
    });

    assertNotEq(config.gasLimit, newFeeConfig.gasLimit);
    // target block rate is the same for wagmi and cchain
    // assertNotEq(config.targetBlockRate, newFeeConfig.targetBlockRate);
    assertNotEq(config.minBaseFee, newFeeConfig.minBaseFee);
    assertNotEq(config.targetGas, newFeeConfig.targetGas);
    assertNotEq(config.baseFeeChangeDenominator, newFeeConfig.baseFeeChangeDenominator);
    // min block gas cost is the same for wagmi and cchain
    // assertNotEq(config.minBlockGasCost, newFeeConfig.minBlockGasCost);
    assertNotEq(config.maxBlockGasCost, newFeeConfig.maxBlockGasCost);
    assertNotEq(config.blockGasCostStep, newFeeConfig.blockGasCostStep);

    example.enableCChainFees();

    FeeConfig memory changedFeeConfig = example.getCurrentFeeConfig();

    assertEq(changedFeeConfig.gasLimit, newFeeConfig.gasLimit);
    // target block rate is the same for wagmi and cchain
    // assertEq(changedFeeConfig.targetBlockRate, newFeeConfig.targetBlockRate);
    assertEq(changedFeeConfig.minBaseFee, newFeeConfig.minBaseFee);
    assertEq(changedFeeConfig.targetGas, newFeeConfig.targetGas);
    assertEq(changedFeeConfig.baseFeeChangeDenominator, newFeeConfig.baseFeeChangeDenominator);
    // min block gas cost is the same for wagmi and cchain
    // assertEq(changedFeeConfig.minBlockGasCost, newFeeConfig.minBlockGasCost);
    assertEq(changedFeeConfig.maxBlockGasCost, newFeeConfig.maxBlockGasCost);
    assertEq(changedFeeConfig.blockGasCostStep, newFeeConfig.blockGasCostStep);

    assertEq(example.getFeeConfigLastChangedAt(), block.number);

    example.enableCustomFees(config);
  }

  function test_minFeeTransaction() public {
    // noop
  }

  function test_raiseMinFeeByOne() public {
    ExampleFeeManager example = new ExampleFeeManager();
    address exampleAddress = address(example);

    IFeeManager manager = IFeeManager(FEE_MANAGER_ADDRESS);

    manager.setEnabled(exampleAddress);

    FeeConfig memory config = example.getCurrentFeeConfig();
    config.minBaseFee = config.minBaseFee + 1;

    example.enableCustomFees(config);
  }

  function test_lowerMinFeeByOne() public {
    ExampleFeeManager example = new ExampleFeeManager();
    address exampleAddress = address(example);

    IFeeManager manager = IFeeManager(FEE_MANAGER_ADDRESS);

    manager.setEnabled(exampleAddress);

    FeeConfig memory config = example.getCurrentFeeConfig();
    config.minBaseFee = config.minBaseFee - 1;

    example.enableCustomFees(config);
  }
}
