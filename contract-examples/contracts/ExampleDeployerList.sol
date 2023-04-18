//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "./IAllowList.sol";
import "./AllowList.sol";
import "ds-test/src/test.sol";

address constant DEPLOYER_LIST = 0x0200000000000000000000000000000000000000;

// ExampleDeployerList shows how ContractDeployerAllowList precompile can be used in a smart contract
// All methods of [allowList] can be directly called. There are example calls as tasks in hardhat.config.ts file.
contract ExampleDeployerList is AllowList {
  // Precompiled Allow List Contract Address
  constructor() AllowList(DEPLOYER_LIST) {}

  function deployToken() public {
    new AllowList(DEPLOYER_LIST);
  }
}

// TODO: a bunch of these tests have repeated code that should be combined
contract ExampleDeployerListTest is DSTest {
  function setUp() public {
    // noop
  }

  function test_verifySenderIsAdmin() public {
    IAllowList allowList = IAllowList(DEPLOYER_LIST);
    assertEq(allowList.readAllowList(msg.sender), 2);
  }

  function test_newAddressHasNoRole() public {
    ExampleDeployerList example = new ExampleDeployerList();
    address exampleAddress = address(example);
    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertEq(allowList.readAllowList(exampleAddress), 0);
  }

  function test_noRoleIsNotAdmin() public {
    ExampleDeployerList example = new ExampleDeployerList();
    address exampleAddress = address(example);

    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertEq(allowList.readAllowList(exampleAddress), 0);
    assertTrue(!example.isAdmin(exampleAddress));
  }

  function test_ownerIsAdmin() public {
    ExampleDeployerList example = new ExampleDeployerList();
    address exampleAddress = address(example);

    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertEq(allowList.readAllowList(exampleAddress), 0);
    assertTrue(example.isAdmin(address(this)));
  }

  function test_noRoleCannotDeploy() public {
    ExampleDeployerList example = new ExampleDeployerList();
    address exampleAddress = address(example);

    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertEq(allowList.readAllowList(exampleAddress), 0);

    try example.deployToken() {
      assertTrue(false, "deployToken should fail");
    } catch {} // TODO should match on an error to make sure that this is failing in the way that's expected
  }

  function test_adminAddContractAsAdmin() public {
    ExampleDeployerList example = new ExampleDeployerList();
    address exampleAddress = address(example);

    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertEq(allowList.readAllowList(exampleAddress), 0);

    allowList.setAdmin(exampleAddress);

    assertEq(allowList.readAllowList(exampleAddress), 2);

    assertTrue(example.isAdmin(exampleAddress));
  }

  function test_addDeployerThroughContract() public {
    ExampleDeployerList example = new ExampleDeployerList();
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
    ExampleDeployerList example = new ExampleDeployerList();
    ExampleDeployerList deployer = new ExampleDeployerList();
    address exampleAddress = address(example);
    address deployerAddress = address(deployer);

    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertEq(allowList.readAllowList(exampleAddress), 0);

    allowList.setAdmin(exampleAddress);

    assertEq(allowList.readAllowList(exampleAddress), 2);

    example.setEnabled(deployerAddress);

    assertTrue(example.isEnabled(deployerAddress));

    deployer.deployToken();
  }

  function test_adminCanRevokeDeployer() public {
    ExampleDeployerList example = new ExampleDeployerList();
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
