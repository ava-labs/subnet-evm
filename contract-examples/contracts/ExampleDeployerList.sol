//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "./IAllowList.sol";
import "./AllowList.sol";
import "ds-test/src/test.sol";

address constant DEPLOYER_LIST = 0x0200000000000000000000000000000000000000;
address constant OTHER_ADDRESS = 0x0Fa8EA536Be85F32724D57A37758761B86416123;

// ExampleDeployerList shows how ContractDeployerAllowList precompile can be used in a smart contract
// All methods of [allowList] can be directly called. There are example calls as tasks in hardhat.config.ts file.
contract ExampleDeployerList is AllowList {
  // Precompiled Allow List Contract Address
  constructor() AllowList(DEPLOYER_LIST) {}

  function deployContract() public {
    new Example();
  }
}

contract Example {}

// TODO: a bunch of these tests have repeated code that should be combined
contract ExampleDeployerListTest is DSTest {
  ExampleDeployerList private example;

  function setUp() public {
    example = new ExampleDeployerList();
    IAllowList allowList = IAllowList(DEPLOYER_LIST);
    allowList.setNone(OTHER_ADDRESS);
  }

  function test_verifySenderIsAdmin() public {
    IAllowList allowList = IAllowList(DEPLOYER_LIST);
    assertEq(allowList.readAllowList(msg.sender), 2);
  }

  function test_newAddressHasNoRole() public {
    address exampleAddress = address(example);
    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertEq(allowList.readAllowList(exampleAddress), 0);
  }

  function test_noRoleIsNotAdmin() public {
    address exampleAddress = address(example);

    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertEq(allowList.readAllowList(exampleAddress), 0);
    assertTrue(!example.isAdmin(exampleAddress));
  }

  function test_ownerIsAdmin() public {
    address exampleAddress = address(example);

    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertEq(allowList.readAllowList(exampleAddress), 0);
    assertTrue(example.isAdmin(address(this)));
  }

  function test_noRoleCannotDeploy() public {
    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertEq(allowList.readAllowList(tx.origin), 0);

    try example.deployContract() {
      assertTrue(false, "deployContract should fail");
    } catch {} // TODO should match on an error to make sure that this is failing in the way that's expected
  }

  function test_adminAddContractAsAdmin() public {
    address exampleAddress = address(example);

    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertEq(allowList.readAllowList(exampleAddress), 0);

    allowList.setAdmin(exampleAddress);

    assertEq(allowList.readAllowList(exampleAddress), 2);

    assertTrue(example.isAdmin(exampleAddress));
  }

  function test_addDeployerThroughContract() public {
    ExampleDeployerList other = new ExampleDeployerList();
    address exampleAddress = address(example);
    address otherAddress = address(other);

    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertEq(allowList.readAllowList(exampleAddress), 0);

    allowList.setAdmin(exampleAddress);

    assertEq(allowList.readAllowList(exampleAddress), 2);

    example.setEnabled(otherAddress);

    assertTrue(example.isEnabled(otherAddress));
  }

  function test_deployerCanDeploy() public {
    ExampleDeployerList deployer = new ExampleDeployerList();
    address exampleAddress = address(example);
    address deployerAddress = address(deployer);

    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertEq(allowList.readAllowList(exampleAddress), 0);

    allowList.setAdmin(exampleAddress);

    assertEq(allowList.readAllowList(exampleAddress), 2);

    example.setEnabled(deployerAddress);

    assertTrue(example.isEnabled(deployerAddress));

    deployer.deployContract();
  }

  function test_adminCanRevokeDeployer() public {
    ExampleDeployerList deployer = new ExampleDeployerList();
    address exampleAddress = address(example);
    address deployerAddress = address(deployer);

    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertEq(allowList.readAllowList(exampleAddress), 0);

    allowList.setAdmin(exampleAddress);

    assertEq(allowList.readAllowList(exampleAddress), 2);

    example.setEnabled(deployerAddress);

    assertTrue(example.isEnabled(deployerAddress));

    example.revoke(deployerAddress);

    assertEq(allowList.readAllowList(deployerAddress), 0);
  }
}
