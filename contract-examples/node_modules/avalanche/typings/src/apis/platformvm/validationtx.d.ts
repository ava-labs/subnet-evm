/**
 * @packageDocumentation
 * @module API-PlatformVM-ValidationTx
 */
import BN from "bn.js";
import { BaseTx } from "./basetx";
import { TransferableOutput } from "../platformvm/outputs";
import { TransferableInput } from "../platformvm/inputs";
import { Buffer } from "buffer/";
import { ParseableOutput } from "./outputs";
import { SerializedEncoding } from "../../utils/serialization";
/**
 * Abstract class representing an transactions with validation information.
 */
export declare abstract class ValidatorTx extends BaseTx {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected nodeID: Buffer;
    protected startTime: Buffer;
    protected endTime: Buffer;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} for the stake amount.
     */
    getNodeID(): Buffer;
    /**
     * Returns a string for the nodeID amount.
     */
    getNodeIDString(): string;
    /**
     * Returns a {@link https://github.com/indutny/bn.js/|BN} for the stake amount.
     */
    getStartTime(): BN;
    /**
     * Returns a {@link https://github.com/indutny/bn.js/|BN} for the stake amount.
     */
    getEndTime(): BN;
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[ValidatorTx]].
     */
    toBuffer(): Buffer;
    constructor(networkID: number, blockchainID: Buffer, outs: TransferableOutput[], ins: TransferableInput[], memo?: Buffer, nodeID?: Buffer, startTime?: BN, endTime?: BN);
}
export declare abstract class WeightedValidatorTx extends ValidatorTx {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected weight: Buffer;
    /**
     * Returns a {@link https://github.com/indutny/bn.js/|BN} for the stake amount.
     */
    getWeight(): BN;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} for the stake amount.
     */
    getWeightBuffer(): Buffer;
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[AddSubnetValidatorTx]].
     */
    toBuffer(): Buffer;
    /**
     * Class representing an unsigned AddSubnetValidatorTx transaction.
     *
     * @param networkID Optional. Networkid, [[DefaultNetworkID]]
     * @param blockchainID Optional. Blockchainid, default Buffer.alloc(32, 16)
     * @param outs Optional. Array of the [[TransferableOutput]]s
     * @param ins Optional. Array of the [[TransferableInput]]s
     * @param memo Optional. {@link https://github.com/feross/buffer|Buffer} for the memo field
     * @param nodeID Optional. The node ID of the validator being added.
     * @param startTime Optional. The Unix time when the validator starts validating the Primary Network.
     * @param endTime Optional. The Unix time when the validator stops validating the Primary Network (and staked AVAX is returned).
     * @param weight Optional. The amount of nAVAX the validator is staking.
     */
    constructor(networkID?: number, blockchainID?: Buffer, outs?: TransferableOutput[], ins?: TransferableInput[], memo?: Buffer, nodeID?: Buffer, startTime?: BN, endTime?: BN, weight?: BN);
}
/**
 * Class representing an unsigned AddDelegatorTx transaction.
 */
export declare class AddDelegatorTx extends WeightedValidatorTx {
    protected _typeName: string;
    protected _typeID: number;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected stakeOuts: TransferableOutput[];
    protected rewardOwners: ParseableOutput;
    /**
     * Returns the id of the [[AddDelegatorTx]]
     */
    getTxType(): number;
    /**
     * Returns a {@link https://github.com/indutny/bn.js/|BN} for the stake amount.
     */
    getStakeAmount(): BN;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} for the stake amount.
     */
    getStakeAmountBuffer(): Buffer;
    /**
     * Returns the array of outputs being staked.
     */
    getStakeOuts(): TransferableOutput[];
    /**
     * Should match stakeAmount. Used in sanity checking.
     */
    getStakeOutsTotal(): BN;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} for the reward address.
     */
    getRewardOwners(): ParseableOutput;
    getTotalOuts(): TransferableOutput[];
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[AddDelegatorTx]].
     */
    toBuffer(): Buffer;
    clone(): this;
    create(...args: any[]): this;
    /**
     * Class representing an unsigned AddDelegatorTx transaction.
     *
     * @param networkID Optional. Networkid, [[DefaultNetworkID]]
     * @param blockchainID Optional. Blockchainid, default Buffer.alloc(32, 16)
     * @param outs Optional. Array of the [[TransferableOutput]]s
     * @param ins Optional. Array of the [[TransferableInput]]s
     * @param memo Optional. {@link https://github.com/feross/buffer|Buffer} for the memo field
     * @param nodeID Optional. The node ID of the validator being added.
     * @param startTime Optional. The Unix time when the validator starts validating the Primary Network.
     * @param endTime Optional. The Unix time when the validator stops validating the Primary Network (and staked AVAX is returned).
     * @param stakeAmount Optional. The amount of nAVAX the validator is staking.
     * @param stakeOuts Optional. The outputs used in paying the stake.
     * @param rewardOwners Optional. The [[ParseableOutput]] containing a [[SECPOwnerOutput]] for the rewards.
     */
    constructor(networkID?: number, blockchainID?: Buffer, outs?: TransferableOutput[], ins?: TransferableInput[], memo?: Buffer, nodeID?: Buffer, startTime?: BN, endTime?: BN, stakeAmount?: BN, stakeOuts?: TransferableOutput[], rewardOwners?: ParseableOutput);
}
export declare class AddValidatorTx extends AddDelegatorTx {
    protected _typeName: string;
    protected _typeID: number;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected delegationFee: number;
    private static delegatorMultiplier;
    /**
     * Returns the id of the [[AddValidatorTx]]
     */
    getTxType(): number;
    /**
     * Returns the delegation fee (represents a percentage from 0 to 100);
     */
    getDelegationFee(): number;
    /**
     * Returns the binary representation of the delegation fee as a {@link https://github.com/feross/buffer|Buffer}.
     */
    getDelegationFeeBuffer(): Buffer;
    fromBuffer(bytes: Buffer, offset?: number): number;
    toBuffer(): Buffer;
    /**
     * Class representing an unsigned AddValidatorTx transaction.
     *
     * @param networkID Optional. Networkid, [[DefaultNetworkID]]
     * @param blockchainID Optional. Blockchainid, default Buffer.alloc(32, 16)
     * @param outs Optional. Array of the [[TransferableOutput]]s
     * @param ins Optional. Array of the [[TransferableInput]]s
     * @param memo Optional. {@link https://github.com/feross/buffer|Buffer} for the memo field
     * @param nodeID Optional. The node ID of the validator being added.
     * @param startTime Optional. The Unix time when the validator starts validating the Primary Network.
     * @param endTime Optional. The Unix time when the validator stops validating the Primary Network (and staked AVAX is returned).
     * @param stakeAmount Optional. The amount of nAVAX the validator is staking.
     * @param stakeOuts Optional. The outputs used in paying the stake.
     * @param rewardOwners Optional. The [[ParseableOutput]] containing the [[SECPOwnerOutput]] for the rewards.
     * @param delegationFee Optional. The percent fee this validator charges when others delegate stake to them.
     * Up to 4 decimal places allowed; additional decimal places are ignored. Must be between 0 and 100, inclusive.
     * For example, if delegationFeeRate is 1.2345 and someone delegates to this validator, then when the delegation
     * period is over, 1.2345% of the reward goes to the validator and the rest goes to the delegator.
     */
    constructor(networkID?: number, blockchainID?: Buffer, outs?: TransferableOutput[], ins?: TransferableInput[], memo?: Buffer, nodeID?: Buffer, startTime?: BN, endTime?: BN, stakeAmount?: BN, stakeOuts?: TransferableOutput[], rewardOwners?: ParseableOutput, delegationFee?: number);
}
//# sourceMappingURL=validationtx.d.ts.map