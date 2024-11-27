# Validators

The Validators package is a collection of structs and functions to manage the state and uptime of validators in the Subnet-EVM. It consists of the following components:

- State package : The state package stores the validator state and uptime information.
- Uptime package: The uptime package manages the uptime tracking of the validators.
- Manager struct: The manager struct is responsible for managing the state and uptime of the validators.

## State Package

The state package stores the validator state and uptime information. The state package implements a CRUD interface for validators. The implementation tracks validators by their validationIDs and assumes they're unique per node and their validation period. The state implementation also assumes NodeIDs are unique in the tracked set. The state implementation only allows existing validator's `weight` and `IsActive` fields to be updated; all other fields should be constant and if any other field changes, the state manager errors and does not update the validator.

For L1 validators, an `active` status implies the validator balance on the P-Chain is sufficient to cover the continuous validation fee. When an L1 validator balance is depleted, it is marked as `inactive` on the P-Chain and this information is passed to the Subnet-EVM's state.

The State interface allows a listener to register to the state changes including validator addition, removal and active status change. The listener always receives the full state when it first subscribes.

The package defines how to serialize the data according to the codec. It can read and write the validator state and uptime information within the database.

## Uptime Package

The uptime package manages the uptime tracking of the L1 validators. It wraps AvalancheGo's uptime tracking manager under the hood and additionally introduces pausable uptime manager interface. The pausable uptime manager interface allows the manager to pause and resume the uptime tracking for a specific validator. 

The uptime package must be run on at least one L1 node, referred to in this document as the "tracker node".

Uptime tracking works as follows:

1- StartTracking: Nodes can start uptime tracking with `StartTracking` method when they're bootstrapped. This method updates the uptime of up-to-date validators by adding the duration between their last updated and tracker node's initializing time to their uptimes. This effectively adds the offline duration of tracker node's to the uptime of validators and optimistically assumes that the validators are online during this period. Subnet-EVM's Pausable manager does not directly modifies this behaviour and it also updates validators that were paused/inactive before the node initialized. Pausable Uptime Manager assumes peers are online and active (has enough fees) when tracker nodes are offline.

2- Connected: Avalanche Uptime manager records the time when a peer is connected to the tracker node (the node running the uptime tracking). When a paused (inactive) validator is connected, pausable uptime manager does not directly invokes the `Connected` on Avalanche Uptime manager, thus the connection time is not directly recorded. Instead, pausable uptime manager waits for the validator to be resumed (top-up fee balance). When the validator is resumed, the tracker records the resumed time and starts tracking the uptime of the validator. Note: Uptime manager does not check if the connected peer is a validator or not. It records the connection time assuming that a non-validator peer can become a validator whilst they're connected to the uptime manager.

3- Disconnected: When a peer validator is disconnected, Avalanche Uptime manager updates the uptime of the validator by adding the duration between the connection time and the disconnection time to the uptime of the validator. The pausable uptime manager handles the inactive peers as if they were disconnected when they are paused, thus it assumes that no paused peers can be disconnected again from the pausable uptime manager.

4- Pause: Pausable Uptime Manager can listen the validator status change via subscribing to the state. When state invokes the `OnValidatorStatusChange` method, pausable uptime manager pauses the uptime tracking of the validator if the validator is inactive. When a validator is paused, it is treated as if it is disconnected from the tracker node; thus it's uptime is updated from the connection time to the pause time and uptime manager stops tracking the uptime of the validator.

5- Resume: When a paused validator peer resumes (status become active), pausable uptime manager resumes the uptime tracking of the validator. It basically treat the peer as if it is connected to the tracker node. Note: Pausable uptime manager holds the set of connected peers that does track the connected peer in p2p layer. The set is used to start tracking the uptime of the paused validators when they resume; this is because the inner AvalancheGo manager thinks that the peer is completely disconnected when it is paused. Pausable uptime manager is able to re-connect them to the inner manager by using this additional connected set.

6- CalculateUptime: The CalculateUptime method calculates a node's updated uptime based on its connection status, connected time and the current time. It first retrieves the node's current uptime and last update time from the state, returning an error if retrieval fails. If tracking hasn’t started, it assumes the node has been online since the last update, adding this duration to its uptime. If the node is not connected and tracking is active, uptime remains unchanged and returned. For connected nodes, the method ensures the connection time does not predate the last update to avoid double-counting. Finally, it adds the duration since the connection time to the node's uptime and returns the updated values.

## Validator Manager Struct

`ValidatorManager` struct is responsible for managing the state of the validators by fetching the information from P-Chain state (via `GetCurrentValidatorSet` in chain context) and updating the state accordingly. It dispatches a `goroutine` to sync the validator state every 60 seconds. The manager fetches the up-to-date validator set from P-Chain and performs the sync operation. The sync operation first performs removing the validators from the state that are not in the P-Chain validator set. Then it adds new validators and updates the existing validators in the state. This order off operations ensures that the uptimes of validators being removed and re-added under same nodeIDs are updated in the same sync operation despite having different validationIDs.

P-Chain's `GetCurrentValidatorSet` can report both L1 and Subnet validators. Subnet-EVM's uptime manager also tracks both of these validator types. So even if a the Subnet has not yet been converted to an L1, the uptime and validator state tracking is still performed by Subnet-EVM.

Validator Manager persists the state to disk at the end of every sync operation. The VM also persists the validator database when the node is shutting down.