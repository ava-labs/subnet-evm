# Staking Manager

The Staking Manager is a contract that manages the staking and validators in the subnet. It is capable and responsible for the following:

- Register itself as the manager contract for subnet on P-chain
- Register a validator on the subnet
- Set weight of a validator
- Uptime tracking
- Remove a validator

Validator/Staking manager contract does not directly operate on P-chain. It sends warp messages and receive warp messages from P-Chain. This ultimately makes it asynchronous with P-Chain, all operations should be considered not final until a related configmation warp message is relayed from P-chain. The manager contract should register it's address on P-chain to be able to interact with the P-chain. The P-chain will only accept warp messages from the registered contract address.

## Register Manager Contract

A Manager contract is expected to be registered in P-Chain, before using any staking operations. This contract will be used to manage the staking operations for the Subnet.

Related ACP section can be found [here.](https://github.com/avalanche-foundation/ACPs/tree/main/ACPs/77-reinventing-subnets#setsubnetvalidatormanagertx)

A contract must be deployed on the chain and an address must be obtained. This address will be used to register the validator manager contract for the subnet. For the first time only, the `SetSubnetValidatorManagerTx` will be originated from subnet validators (without a sourceAddress in warp message). A response from P-Chain that the contract is registered should be relayed back to the contract, so that the contract can activate itself (start staking, rewarding etc).

Subsequent `SetSubnetValidatorManagerTx` can be originated from the contract to replace the current manager (itself) with another contract. In this case the received confirmation of `SetSubnetValidatorManagerTx` should be relayed back to both contracts, so that the old manager contract can deactivate itself (stop staking, rewarding etc) and the new manager can activate itself.

## Register Validator

Related ACP-77 section [here](https://github.com/avalanche-foundation/ACPs/tree/main/ACPs/77-reinventing-subnets#step-1-retrieve-a-bls-multisig-from-the-subnet)

### registerValidator

Register validator is the process of adding a validator to the subnet.The contract initiates this with a `registerValidator` function that takes following inputs:

- Subnet ID (bytes32): The ID of the subnet where the validator is being registered
- Node ID (bytes32): The ID of the node that is being registered as a validator
- Amount (uint64): The staked amount that will determine the weight of the validator
- expirationTime (uint64): The time when this register message will be expired
- ed25519 signature (bytes): The signature of the message signed by the node's BLS key

The messageID of the request can be crafted with `sha256(subnetID, nodeID, amount, expirationTime, signature)`. This should be same with the one that will be crafted for P-chain.

The contract should verify these:

- (Optional) the staking amount is greater than the minimum staking amount
- (Optional) The node is not already registered as a validator. This will be handled by the P-Chain though.
- (Optional) Expiry timestamp is not in the past
- The signature has valid length (64 bytes)
- The nodeID is not empty
- The messageID is not already used

The contract should then send a warp message to P-Chain and lock the staked amount in the contract. The manager should not start accruing reward at this time since the validator is not yet accepted by P-Chain and it's not guaranteed. The warp message will be aggregated and signed, then the validator will issue a `RegisterSubnetValidatorTx` transaction on P-Chain with warp message included. More details from ACP-77 [here](https://github.com/avalanche-foundation/ACPs/tree/main/ACPs/77-reinventing-subnets#step-2-issue-a-registersubnetvalidatortx-on-the-p-chain)

### receiveRegisterValidator

Once `RegisterSubnetValidatorTx` is accepted a warp message will be aggregated and relayed to this contract as a result of `RegisterSubnetValidatorTx`. Upon receiving and verifying this message (via `receiveRegisterValidator`), the contract should activate the validator and start accruing rewards. The message will only contain `messageID` which corresponds to the `messageID` in `registerValidator` request. The contract should verify these:

- Verified warp message
- The typecheck on the warp message (== ValidatorRegisteredMessage)
- messageID is not empty
- The message is not already used
- (?) There is a pending request for this messageID
- (?) The nodeID in the pending request is not a validator already
- (?) Message is from P-Chain

TODO: we might not want last 3 checks (with ?) to minimize reverts. The contract should be able to handle these cases without reverting. This is because the timing of the messages is not guaranteed and the contract should be able to handle out-of-order messages.

If `RegisterSubnetValidatorTx` is not accepted, another warp message `InvalidValidatorRegisterMessage` will be relayed to the contract.

## Remove Validator

### removeValidator

The manager contract can remove a validator gracefully by sending a warp message to P-Chain. `removeValidator` function should be called with the following inputs:

- Subnet ID (bytes32): The ID of the subnet where the validator is being removed
- Node ID (bytes32): The ID of the node that is being removed as a validator

The contract should verify these:

- The node is already registered as a validator
- The address that registered the node is the same with the address that is removing the validator
- (Optional) Validator being removed after minimum staking period

After verification contract should perform the following:

- Send a warp message to P-Chain for `SetSubnetValidatorWeightTx{messageID: relatedRegisterValidatorMessageID, nonce: MaxUint64, Weight: 0, }` transaction
- Stop accruing rewards for the validator

### receiveRemoveValidator

Once the `SetSubnetValidatorWeightTx` is accepted, a warp message `InvalidValidatorRegisterMessage` will be relayed to the contract with `messageID` of the request. The message represents a messageID will not possibly become a valid validator. This essentialy means the message was not valid in the first place, or the validator has removed from the set. The InvalidValidatorRegisterMessage can either be a result of `removeValidator` from the contract or it can be originated directly from P-chain via `ExitValidatorSetTx`. In either case this ultimately tells the contract that a validator was evicted. See related ACP-77 section [here](https://github.com/avalanche-foundation/ACPs/tree/main/ACPs/77-reinventing-subnets#exitvalidatorsettx).

The contract should verify these in `receiveRemoveValidator`:

- Verified warp message
- The typecheck on the warp message (== InvalidValidatorRegisterMessage)
- messageID is not empty

At this point the contract should unlock all funds that were locked for the validator and issue locked rewards. If the validator was not removed by the contract itself, the contract should either give partial rewards or no rewards at all. This is because we don't know when the validator was removed.

Rewards should not be unlocked or sent to the staker at this point. This is because we need to first perform a uptime check to determine the rewards.

TODO: should we give a partial reward in case the validator got removed from P-chain (balance drained)? Possibly by timestamping the `InvalidValidatorRegisterMessage`.

TODO: should we add a cooldown period for validator removal/readd?

## Uptime Based Rewards

The manager contract should reward based on uptimes of validators. VMs should be able to track precise uptimes of validators. The manager contract should be able to query the uptime of a validator and reward based on the uptime.

TODO: Design currently is not finalized. We need to discuss how to track uptimes and how to report them to the manager contract. Some ideas:

1- If a validator becomes ineligible permanetly due to low uptime, a warp message should be sent to the manager contract.
2- Once a validator is removed, a warp message that indicates the uptime is enough will be sent to the contract.
3- A message will be aggregated upon a request from the manager contract (e.g via an event), then the message will be sent to the manager contract.

## TODO: Add/Decrease Weight

## TODO: Delegations