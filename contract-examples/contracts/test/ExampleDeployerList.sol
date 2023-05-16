//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../ExampleDeployerList.sol";
import "../IAllowList.sol";
import "../AllowList.sol";
import "./AllowList.sol";

// ExampleDeployerListTest defines transactions that are used to test
// the DeployerAllowList precompile by instantiating and calling the
// ExampleDeployerList and making assertions.
// The transactions are put together as steps of a complete test in contract_deployer_allow_list.ts.
// TODO: a bunch of these tests have repeated code that should be combined
contract ExampleDeployerListTest is AllowListTest {
  ExampleDeployerList private example;

  function setUp() public {
    example = new ExampleDeployerList();
    IAllowList allowList = IAllowList(DEPLOYER_LIST);
    allowList.setNone(OTHER_ADDRESS);
  }

  function test_verifySenderIsAdmin() public {
    IAllowList allowList = IAllowList(DEPLOYER_LIST);
    assertRole(allowList.readAllowList(msg.sender), AllowList.Role.Admin);
  }

  function test_newAddressHasNoRole() public {
    address exampleAddress = address(example);
    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertRole(allowList.readAllowList(exampleAddress), AllowList.Role.None);
  }

  function test_noRoleIsNotAdmin() public {
    address exampleAddress = address(example);

    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertRole(allowList.readAllowList(exampleAddress), AllowList.Role.None);
    assertTrue(!example.isAdmin(exampleAddress));
  }

  function test_ownerIsAdmin() public {
    address exampleAddress = address(example);

    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertRole(allowList.readAllowList(exampleAddress), AllowList.Role.None);
    assertTrue(example.isAdmin(address(this)));
  }

  function test_noRoleCannotDeploy() public {
    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertRole(allowList.readAllowList(tx.origin), AllowList.Role.None);

    try example.deployContract() {
      assertTrue(false, "deployContract should fail");
    } catch {} // TODO should match on an error to make sure that this is failing in the way that's expected
  }

  function test_adminAddContractAsAdmin() public {
    address exampleAddress = address(example);

    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertRole(allowList.readAllowList(exampleAddress), AllowList.Role.None);

    allowList.setAdmin(exampleAddress);

    assertRole(allowList.readAllowList(exampleAddress), AllowList.Role.Admin);

    assertTrue(example.isAdmin(exampleAddress));
  }

  function test_addDeployerThroughContract() public {
    ExampleDeployerList other = new ExampleDeployerList();
    address exampleAddress = address(example);
    address otherAddress = address(other);

    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertRole(allowList.readAllowList(exampleAddress), AllowList.Role.None);

    allowList.setAdmin(exampleAddress);

    assertRole(allowList.readAllowList(exampleAddress), AllowList.Role.Admin);

    example.setEnabled(otherAddress);

    assertTrue(example.isEnabled(otherAddress));
  }

  function test_deployerCanDeploy() public {
    ExampleDeployerList deployer = new ExampleDeployerList();
    address exampleAddress = address(example);
    address deployerAddress = address(deployer);

    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertRole(allowList.readAllowList(exampleAddress), AllowList.Role.None);

    allowList.setAdmin(exampleAddress);

    assertRole(allowList.readAllowList(exampleAddress), AllowList.Role.Admin);

    example.setEnabled(deployerAddress);

    assertTrue(example.isEnabled(deployerAddress));

    deployer.deployContract();
  }

  function test_adminCanRevokeDeployer() public {
    ExampleDeployerList deployer = new ExampleDeployerList();
    address exampleAddress = address(example);
    address deployerAddress = address(deployer);

    IAllowList allowList = IAllowList(DEPLOYER_LIST);

    assertRole(allowList.readAllowList(exampleAddress), AllowList.Role.None);

    allowList.setAdmin(exampleAddress);

    assertRole(allowList.readAllowList(exampleAddress), AllowList.Role.Admin);

    example.setEnabled(deployerAddress);

    assertTrue(example.isEnabled(deployerAddress));

    example.revoke(deployerAddress);

    assertRole(allowList.readAllowList(deployerAddress), AllowList.Role.None);
  }
}
