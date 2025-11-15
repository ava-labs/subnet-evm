//SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "precompile/allowlist/allowlisttest/IAllowList.sol";
import "precompile/allowlist/allowlisttest/AllowList.sol";
import "precompile/allowlist/allowlisttest/AllowListTest.sol";

// DeployerListTest defines transactions that are used to test
// the DeployerAllowList precompile by instantiating and calling the
// DeployerList and making assertions.
// The transactions are put together as steps of a complete test in contract_deployer_allow_list.ts.
// TODO: a bunch of these tests have repeated code that should be combined
contract DeployerListTest is AllowListTestHelper {
    address constant DEPLOYER_LIST = 0x0200000000000000000000000000000000000000;
    address constant OTHER_ADDRESS = 0x0Fa8EA536Be85F32724D57A37758761B86416123;

    IAllowList allowListPrecompile = IAllowList(DEPLOYER_LIST);
    AllowListTest private allowListTest;

    function setUp() public {
        allowListTest = new AllowListTest(DEPLOYER_LIST);
        allowListPrecompile.setNone(OTHER_ADDRESS);
    }

    function step_verifySenderIsAdmin() public {
        assertRole(
            allowListPrecompile.readAllowList(msg.sender),
            AllowList.Role.Admin
        );
    }

    function step_newAddressHasNoRole() public {
        address allowListTestAddress = address(allowListTest);

        assertRole(
            allowListPrecompile.readAllowList(allowListTestAddress),
            AllowList.Role.None
        );
    }

    function step_noRoleIsNotAdmin() public {
        address allowListTestAddress = address(allowListTest);

        assertRole(
            allowListPrecompile.readAllowList(allowListTestAddress),
            AllowList.Role.None
        );
        assertTrue(!allowListTest.isAdmin(allowListTestAddress));
    }

    function step_noRoleCannotDeploy() public {
        assertRole(
            allowListPrecompile.readAllowList(tx.origin),
            AllowList.Role.None
        );

        try allowListTest.deployContract() {
            assertTrue(false, "deployContract should fail");
        } catch {} // TODO should match on an error to make sure that this is failing in the way that's expected
    }

    function step_adminAddContractAsAdmin() public {
        address allowListTestAddress = address(allowListTest);

        assertRole(
            allowListPrecompile.readAllowList(allowListTestAddress),
            AllowList.Role.None
        );

        allowListPrecompile.setAdmin(allowListTestAddress);

        assertRole(
            allowListPrecompile.readAllowList(allowListTestAddress),
            AllowList.Role.Admin
        );

        assertTrue(allowListTest.isAdmin(allowListTestAddress));
    }

    function step_addDeployerThroughContract() public {
        AllowListTest other = new AllowListTest(DEPLOYER_LIST);
        address allowListTestAddress = address(allowListTest);
        address otherAddress = address(other);

        assertRole(
            allowListPrecompile.readAllowList(allowListTestAddress),
            AllowList.Role.None
        );

        allowListPrecompile.setAdmin(allowListTestAddress);

        assertRole(
            allowListPrecompile.readAllowList(allowListTestAddress),
            AllowList.Role.Admin
        );

        allowListTest.setEnabled(otherAddress);

        assertTrue(allowListTest.isEnabled(otherAddress));
    }

    function step_deployerCanDeploy() public {
        AllowListTest deployer = new AllowListTest(DEPLOYER_LIST);
        address allowListTestAddress = address(allowListTest);
        address deployerAddress = address(deployer);

        assertRole(
            allowListPrecompile.readAllowList(allowListTestAddress),
            AllowList.Role.None
        );

        allowListPrecompile.setAdmin(allowListTestAddress);

        assertRole(
            allowListPrecompile.readAllowList(allowListTestAddress),
            AllowList.Role.Admin
        );

        allowListTest.setEnabled(deployerAddress);

        assertTrue(allowListTest.isEnabled(deployerAddress));

        deployer.deployContract();
    }

    function step_adminCanRevokeDeployer() public {
        AllowListTest deployer = new AllowListTest(DEPLOYER_LIST);
        address allowListTestAddress = address(allowListTest);
        address deployerAddress = address(deployer);

        assertRole(
            allowListPrecompile.readAllowList(allowListTestAddress),
            AllowList.Role.None
        );

        allowListPrecompile.setAdmin(allowListTestAddress);

        assertRole(
            allowListPrecompile.readAllowList(allowListTestAddress),
            AllowList.Role.Admin
        );

        allowListTest.setEnabled(deployerAddress);

        assertTrue(allowListTest.isEnabled(deployerAddress));

        allowListTest.revoke(deployerAddress);

        assertRole(
            allowListPrecompile.readAllowList(deployerAddress),
            AllowList.Role.None
        );
    }
}
