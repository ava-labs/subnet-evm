"use strict";
/**
 * @packageDocumentation
 * @module API-PlatformVM-ValidationTx
 */
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.AddValidatorTx = exports.AddDelegatorTx = exports.WeightedValidatorTx = exports.ValidatorTx = void 0;
const bn_js_1 = __importDefault(require("bn.js"));
const bintools_1 = __importDefault(require("../../utils/bintools"));
const basetx_1 = require("./basetx");
const outputs_1 = require("../platformvm/outputs");
const buffer_1 = require("buffer/");
const constants_1 = require("./constants");
const constants_2 = require("../../utils/constants");
const helperfunctions_1 = require("../../utils/helperfunctions");
const outputs_2 = require("./outputs");
const serialization_1 = require("../../utils/serialization");
const errors_1 = require("../../utils/errors");
/**
 * @ignore
 */
const bintools = bintools_1.default.getInstance();
const serialization = serialization_1.Serialization.getInstance();
/**
 * Abstract class representing an transactions with validation information.
 */
class ValidatorTx extends basetx_1.BaseTx {
    constructor(networkID, blockchainID, outs, ins, memo, nodeID, startTime, endTime) {
        super(networkID, blockchainID, outs, ins, memo);
        this._typeName = "ValidatorTx";
        this._typeID = undefined;
        this.nodeID = buffer_1.Buffer.alloc(20);
        this.startTime = buffer_1.Buffer.alloc(8);
        this.endTime = buffer_1.Buffer.alloc(8);
        this.nodeID = nodeID;
        this.startTime = bintools.fromBNToBuffer(startTime, 8);
        this.endTime = bintools.fromBNToBuffer(endTime, 8);
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { nodeID: serialization.encoder(this.nodeID, encoding, "Buffer", "nodeID"), startTime: serialization.encoder(this.startTime, encoding, "Buffer", "decimalString"), endTime: serialization.encoder(this.endTime, encoding, "Buffer", "decimalString") });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.nodeID = serialization.decoder(fields["nodeID"], encoding, "nodeID", "Buffer", 20);
        this.startTime = serialization.decoder(fields["startTime"], encoding, "decimalString", "Buffer", 8);
        this.endTime = serialization.decoder(fields["endTime"], encoding, "decimalString", "Buffer", 8);
    }
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} for the stake amount.
     */
    getNodeID() {
        return this.nodeID;
    }
    /**
     * Returns a string for the nodeID amount.
     */
    getNodeIDString() {
        return (0, helperfunctions_1.bufferToNodeIDString)(this.nodeID);
    }
    /**
     * Returns a {@link https://github.com/indutny/bn.js/|BN} for the stake amount.
     */
    getStartTime() {
        return bintools.fromBufferToBN(this.startTime);
    }
    /**
     * Returns a {@link https://github.com/indutny/bn.js/|BN} for the stake amount.
     */
    getEndTime() {
        return bintools.fromBufferToBN(this.endTime);
    }
    fromBuffer(bytes, offset = 0) {
        offset = super.fromBuffer(bytes, offset);
        this.nodeID = bintools.copyFrom(bytes, offset, offset + 20);
        offset += 20;
        this.startTime = bintools.copyFrom(bytes, offset, offset + 8);
        offset += 8;
        this.endTime = bintools.copyFrom(bytes, offset, offset + 8);
        offset += 8;
        return offset;
    }
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[ValidatorTx]].
     */
    toBuffer() {
        const superbuff = super.toBuffer();
        const bsize = superbuff.length +
            this.nodeID.length +
            this.startTime.length +
            this.endTime.length;
        return buffer_1.Buffer.concat([superbuff, this.nodeID, this.startTime, this.endTime], bsize);
    }
}
exports.ValidatorTx = ValidatorTx;
class WeightedValidatorTx extends ValidatorTx {
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
    constructor(networkID = constants_2.DefaultNetworkID, blockchainID = buffer_1.Buffer.alloc(32, 16), outs = undefined, ins = undefined, memo = undefined, nodeID = undefined, startTime = undefined, endTime = undefined, weight = undefined) {
        super(networkID, blockchainID, outs, ins, memo, nodeID, startTime, endTime);
        this._typeName = "WeightedValidatorTx";
        this._typeID = undefined;
        this.weight = buffer_1.Buffer.alloc(8);
        if (typeof weight !== undefined) {
            this.weight = bintools.fromBNToBuffer(weight, 8);
        }
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { weight: serialization.encoder(this.weight, encoding, "Buffer", "decimalString") });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.weight = serialization.decoder(fields["weight"], encoding, "decimalString", "Buffer", 8);
    }
    /**
     * Returns a {@link https://github.com/indutny/bn.js/|BN} for the stake amount.
     */
    getWeight() {
        return bintools.fromBufferToBN(this.weight);
    }
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} for the stake amount.
     */
    getWeightBuffer() {
        return this.weight;
    }
    fromBuffer(bytes, offset = 0) {
        offset = super.fromBuffer(bytes, offset);
        this.weight = bintools.copyFrom(bytes, offset, offset + 8);
        offset += 8;
        return offset;
    }
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[AddSubnetValidatorTx]].
     */
    toBuffer() {
        const superbuff = super.toBuffer();
        return buffer_1.Buffer.concat([superbuff, this.weight]);
    }
}
exports.WeightedValidatorTx = WeightedValidatorTx;
/**
 * Class representing an unsigned AddDelegatorTx transaction.
 */
class AddDelegatorTx extends WeightedValidatorTx {
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
    constructor(networkID = constants_2.DefaultNetworkID, blockchainID = buffer_1.Buffer.alloc(32, 16), outs = undefined, ins = undefined, memo = undefined, nodeID = undefined, startTime = undefined, endTime = undefined, stakeAmount = undefined, stakeOuts = undefined, rewardOwners = undefined) {
        super(networkID, blockchainID, outs, ins, memo, nodeID, startTime, endTime, stakeAmount);
        this._typeName = "AddDelegatorTx";
        this._typeID = constants_1.PlatformVMConstants.ADDDELEGATORTX;
        this.stakeOuts = [];
        this.rewardOwners = undefined;
        if (typeof stakeOuts !== undefined) {
            this.stakeOuts = stakeOuts;
        }
        this.rewardOwners = rewardOwners;
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { stakeOuts: this.stakeOuts.map((s) => s.serialize(encoding)), rewardOwners: this.rewardOwners.serialize(encoding) });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.stakeOuts = fields["stakeOuts"].map((s) => {
            let xferout = new outputs_1.TransferableOutput();
            xferout.deserialize(s, encoding);
            return xferout;
        });
        this.rewardOwners = new outputs_2.ParseableOutput();
        this.rewardOwners.deserialize(fields["rewardOwners"], encoding);
    }
    /**
     * Returns the id of the [[AddDelegatorTx]]
     */
    getTxType() {
        return this._typeID;
    }
    /**
     * Returns a {@link https://github.com/indutny/bn.js/|BN} for the stake amount.
     */
    getStakeAmount() {
        return this.getWeight();
    }
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} for the stake amount.
     */
    getStakeAmountBuffer() {
        return this.weight;
    }
    /**
     * Returns the array of outputs being staked.
     */
    getStakeOuts() {
        return this.stakeOuts;
    }
    /**
     * Should match stakeAmount. Used in sanity checking.
     */
    getStakeOutsTotal() {
        let val = new bn_js_1.default(0);
        for (let i = 0; i < this.stakeOuts.length; i++) {
            val = val.add(this.stakeOuts[`${i}`].getOutput().getAmount());
        }
        return val;
    }
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} for the reward address.
     */
    getRewardOwners() {
        return this.rewardOwners;
    }
    getTotalOuts() {
        return [...this.getOuts(), ...this.getStakeOuts()];
    }
    fromBuffer(bytes, offset = 0) {
        offset = super.fromBuffer(bytes, offset);
        const numstakeouts = bintools.copyFrom(bytes, offset, offset + 4);
        offset += 4;
        const outcount = numstakeouts.readUInt32BE(0);
        this.stakeOuts = [];
        for (let i = 0; i < outcount; i++) {
            const xferout = new outputs_1.TransferableOutput();
            offset = xferout.fromBuffer(bytes, offset);
            this.stakeOuts.push(xferout);
        }
        this.rewardOwners = new outputs_2.ParseableOutput();
        offset = this.rewardOwners.fromBuffer(bytes, offset);
        return offset;
    }
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[AddDelegatorTx]].
     */
    toBuffer() {
        const superbuff = super.toBuffer();
        let bsize = superbuff.length;
        const numouts = buffer_1.Buffer.alloc(4);
        numouts.writeUInt32BE(this.stakeOuts.length, 0);
        let barr = [super.toBuffer(), numouts];
        bsize += numouts.length;
        this.stakeOuts = this.stakeOuts.sort(outputs_1.TransferableOutput.comparator());
        for (let i = 0; i < this.stakeOuts.length; i++) {
            let out = this.stakeOuts[`${i}`].toBuffer();
            barr.push(out);
            bsize += out.length;
        }
        let ro = this.rewardOwners.toBuffer();
        barr.push(ro);
        bsize += ro.length;
        return buffer_1.Buffer.concat(barr, bsize);
    }
    clone() {
        let newbase = new AddDelegatorTx();
        newbase.fromBuffer(this.toBuffer());
        return newbase;
    }
    create(...args) {
        return new AddDelegatorTx(...args);
    }
}
exports.AddDelegatorTx = AddDelegatorTx;
class AddValidatorTx extends AddDelegatorTx {
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
    constructor(networkID = constants_2.DefaultNetworkID, blockchainID = buffer_1.Buffer.alloc(32, 16), outs = undefined, ins = undefined, memo = undefined, nodeID = undefined, startTime = undefined, endTime = undefined, stakeAmount = undefined, stakeOuts = undefined, rewardOwners = undefined, delegationFee = undefined) {
        super(networkID, blockchainID, outs, ins, memo, nodeID, startTime, endTime, stakeAmount, stakeOuts, rewardOwners);
        this._typeName = "AddValidatorTx";
        this._typeID = constants_1.PlatformVMConstants.ADDVALIDATORTX;
        this.delegationFee = 0;
        if (typeof delegationFee === "number") {
            if (delegationFee >= 0 && delegationFee <= 100) {
                this.delegationFee = parseFloat(delegationFee.toFixed(4));
            }
            else {
                throw new errors_1.DelegationFeeError("AddValidatorTx.constructor -- delegationFee must be in the range of 0 and 100, inclusively.");
            }
        }
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { delegationFee: serialization.encoder(this.getDelegationFeeBuffer(), encoding, "Buffer", "decimalString", 4) });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        let dbuff = serialization.decoder(fields["delegationFee"], encoding, "decimalString", "Buffer", 4);
        this.delegationFee =
            dbuff.readUInt32BE(0) / AddValidatorTx.delegatorMultiplier;
    }
    /**
     * Returns the id of the [[AddValidatorTx]]
     */
    getTxType() {
        return this._typeID;
    }
    /**
     * Returns the delegation fee (represents a percentage from 0 to 100);
     */
    getDelegationFee() {
        return this.delegationFee;
    }
    /**
     * Returns the binary representation of the delegation fee as a {@link https://github.com/feross/buffer|Buffer}.
     */
    getDelegationFeeBuffer() {
        let dBuff = buffer_1.Buffer.alloc(4);
        let buffnum = parseFloat(this.delegationFee.toFixed(4)) *
            AddValidatorTx.delegatorMultiplier;
        dBuff.writeUInt32BE(buffnum, 0);
        return dBuff;
    }
    fromBuffer(bytes, offset = 0) {
        offset = super.fromBuffer(bytes, offset);
        let dbuff = bintools.copyFrom(bytes, offset, offset + 4);
        offset += 4;
        this.delegationFee =
            dbuff.readUInt32BE(0) / AddValidatorTx.delegatorMultiplier;
        return offset;
    }
    toBuffer() {
        let superBuff = super.toBuffer();
        let feeBuff = this.getDelegationFeeBuffer();
        return buffer_1.Buffer.concat([superBuff, feeBuff]);
    }
}
exports.AddValidatorTx = AddValidatorTx;
AddValidatorTx.delegatorMultiplier = 10000;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidmFsaWRhdGlvbnR4LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2FwaXMvcGxhdGZvcm12bS92YWxpZGF0aW9udHgudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IjtBQUFBOzs7R0FHRzs7Ozs7O0FBRUgsa0RBQXNCO0FBQ3RCLG9FQUEyQztBQUMzQyxxQ0FBaUM7QUFDakMsbURBQTBEO0FBRTFELG9DQUFnQztBQUNoQywyQ0FBaUQ7QUFDakQscURBQXdEO0FBQ3hELGlFQUFrRTtBQUNsRSx1Q0FBeUQ7QUFDekQsNkRBQTZFO0FBQzdFLCtDQUF1RDtBQUV2RDs7R0FFRztBQUNILE1BQU0sUUFBUSxHQUFhLGtCQUFRLENBQUMsV0FBVyxFQUFFLENBQUE7QUFDakQsTUFBTSxhQUFhLEdBQWtCLDZCQUFhLENBQUMsV0FBVyxFQUFFLENBQUE7QUFFaEU7O0dBRUc7QUFDSCxNQUFzQixXQUFZLFNBQVEsZUFBTTtJQTBHOUMsWUFDRSxTQUFpQixFQUNqQixZQUFvQixFQUNwQixJQUEwQixFQUMxQixHQUF3QixFQUN4QixJQUFhLEVBQ2IsTUFBZSxFQUNmLFNBQWMsRUFDZCxPQUFZO1FBRVosS0FBSyxDQUFDLFNBQVMsRUFBRSxZQUFZLEVBQUUsSUFBSSxFQUFFLEdBQUcsRUFBRSxJQUFJLENBQUMsQ0FBQTtRQW5IdkMsY0FBUyxHQUFHLGFBQWEsQ0FBQTtRQUN6QixZQUFPLEdBQUcsU0FBUyxDQUFBO1FBOENuQixXQUFNLEdBQVcsZUFBTSxDQUFDLEtBQUssQ0FBQyxFQUFFLENBQUMsQ0FBQTtRQUNqQyxjQUFTLEdBQVcsZUFBTSxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQTtRQUNuQyxZQUFPLEdBQVcsZUFBTSxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQTtRQW1FekMsSUFBSSxDQUFDLE1BQU0sR0FBRyxNQUFNLENBQUE7UUFDcEIsSUFBSSxDQUFDLFNBQVMsR0FBRyxRQUFRLENBQUMsY0FBYyxDQUFDLFNB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Q0FBQyxDQUFDLENBQUMsQ0FBQTtRQUNuQyxJQUFJLE9BQU8sR0FDVCxVQUFVLENBQUMsSUFBSSxDQUFDLGFBQWEsQ0FBQyxPQUFPLENBQUMsQ0FBQyxDQUFDLENBQUM7WUFDekMsY0FBYyxDQUFDLG1CQUFtQixDQUFBO1FBQ3BDLEtBQUssQ0FBQyxhQUFhLENBQUMsT0FBTyxFQUFFLENBQUMsQ0FBQyxDQUFBO1FBQy9CLE9BQU8sS0FBSyxDQUFBO0lBQ2QsQ0FBQztJQUVELFVBQVUsQ0FBQyxLQUFhLEVBQUUsU0FBaUIsQ0FBQztRQUMxQyxNQUFNLEdBQUcsS0FBSyxDQUFDLFVBQVUsQ0FBQyxLQUFLLEVBQUUsTUFBTSxDQUFDLENBQUE7UUFDeEMsSUFBSSxLQUFLLEdBQVcsUUFBUSxDQUFDLFFBQVEsQ0FBQyxLQUFLLEVBQUUsTUFBTSxFQUFFLE1BQU0sR0FBRyxDQUFDLENBQUMsQ0FBQTtRQUNoRSxNQUFNLElBQUksQ0FBQyxDQUFBO1FBQ1gsSUFBSSxDQUFDLGFBQWE7WUFDaEIsS0FBSyxDQUFDLFlBQVksQ0FBQyxDQUFDLENBQUMsR0FBRyxjQUFjLENBQUMsbUJBQW1CLENBQUE7UUFDNUQsT0FBTyxNQUFNLENBQUE7SUFDZixDQUFDO0lBRUQsUUFBUTtRQUNOLElBQUksU0FBUyxHQUFXLEtBQUssQ0FBQyxRQUFRLEVBQUUsQ0FBQTtRQUN4QyxJQUFJLE9BQU8sR0FBVyxJQUFJLENBQUMsc0JBQXNCLEVBQUUsQ0FBQTtRQUNuRCxPQUFPLGVBQU0sQ0FBQyxNQUFNLENBQUMsQ0FBQyxTQUFTLEVBQUUsT0FBTyxDQUFDLENBQUMsQ0FBQTtJQUM1QyxDQUFDOztBQXhFSCx3Q0FrSUM7QUFuR2dCLGtDQUFtQixHQUFXLEtBQUssQ0FBQSIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogQHBhY2thZ2VEb2N1bWVudGF0aW9uXG4gKiBAbW9kdWxlIEFQSS1QbGF0Zm9ybVZNLVZhbGlkYXRpb25UeFxuICovXG5cbmltcG9ydCBCTiBmcm9tIFwiYm4uanNcIlxuaW1wb3J0IEJpblRvb2xzIGZyb20gXCIuLi8uLi91dGlscy9iaW50b29sc1wiXG5pbXBvcnQgeyBCYXNlVHggfSBmcm9tIFwiLi9iYXNldHhcIlxuaW1wb3J0IHsgVHJhbnNmZXJhYmxlT3V0cHV0IH0gZnJvbSBcIi4uL3BsYXRmb3Jtdm0vb3V0cHV0c1wiXG5pbXBvcnQgeyBUcmFuc2ZlcmFibGVJbnB1dCB9IGZyb20gXCIuLi9wbGF0Zm9ybXZtL2lucHV0c1wiXG5pbXBvcnQgeyBCdWZmZXIgfSBmcm9tIFwiYnVmZmVyL1wiXG5pbXBvcnQgeyBQbGF0Zm9ybVZNQ29uc3RhbnRzIH0gZnJvbSBcIi4vY29uc3RhbnRzXCJcbmltcG9ydCB7IERlZmF1bHROZXR3b3JrSUQgfSBmcm9tIFwiLi4vLi4vdXRpbHMvY29uc3RhbnRzXCJcbmltcG9ydCB7IGJ1ZmZlclRvTm9kZUlEU3RyaW5nIH0gZnJvbSBcIi4uLy4uL3V0aWxzL2hlbHBlcmZ1bmN0aW9uc1wiXG5pbXBvcnQgeyBBbW91bnRPdXRwdXQsIFBhcnNlYWJsZU91dHB1dCB9IGZyb20gXCIuL291dHB1dHNcIlxuaW1wb3J0IHsgU2VyaWFsaXphdGlvbiwgU2VyaWFsaXplZEVuY29kaW5nIH0gZnJvbSBcIi4uLy4uL3V0aWxzL3NlcmlhbGl6YXRpb25cIlxuaW1wb3J0IHsgRGVsZWdhdGlvbkZlZUVycm9yIH0gZnJvbSBcIi4uLy4uL3V0aWxzL2Vycm9yc1wiXG5cbi8qKlxuICogQGlnbm9yZVxuICovXG5jb25zdCBiaW50b29sczogQmluVG9vbHMgPSBCaW5Ub29scy5nZXRJbnN0YW5jZSgpXG5jb25zdCBzZXJpYWxpemF0aW9uOiBTZXJpYWxpemF0aW9uID0gU2VyaWFsaXphdGlvbi5nZXRJbnN0YW5jZSgpXG5cbi8qKlxuICogQWJzdHJhY3QgY2xhc3MgcmVwcmVzZW50aW5nIGFuIHRyYW5zYWN0aW9ucyB3aXRoIHZhbGlkYXRpb24gaW5mb3JtYXRpb24uXG4gKi9cbmV4cG9ydCBhYnN0cmFjdCBjbGFzcyBWYWxpZGF0b3JUeCBleHRlbmRzIEJhc2VUeCB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIlZhbGlkYXRvclR4XCJcbiAgcHJvdGVjdGVkIF90eXBlSUQgPSB1bmRlZmluZWRcblxuICBzZXJpYWxpemUoZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpOiBvYmplY3Qge1xuICAgIGxldCBmaWVsZHM6IG9iamVjdCA9IHN1cGVyLnNlcmlhbGl6ZShlbmNvZGluZylcbiAgICByZXR1cm4ge1xuICAgICAgLi4uZmllbGRzLFxuICAgICAgbm9kZUlEOiBzZXJpYWxpemF0aW9uLmVuY29kZXIodGhpcy5ub2RlSUQsIGVuY29kaW5nLCBcIkJ1ZmZlclwiLCBcIm5vZGVJRFwiKSxcbiAgICAgIHN0YXJ0VGltZTogc2VyaWFsaXphdGlvbi5lbmNvZGVyKFxuICAgICAgICB0aGlzLnN0YXJ0VGltZSxcbiAgICAgICAgZW5jb2RpbmcsXG4gICAgICAgIFwiQnVmZmVyXCIsXG4gICAgICAgIFwiZGVjaW1hbFN0cmluZ1wiXG4gICAgICApLFxuICAgICAgZW5kVGltZTogc2VyaWFsaXphdGlvbi5lbmNvZGVyKFxuICAgICAgICB0aGlzLmVuZFRpbWUsXG4gICAgICAgIGVuY29kaW5nLFxuICAgICAgICBcIkJ1ZmZlclwiLFxuICAgICAgICBcImRlY2ltYWxTdHJpbmdcIlxuICAgICAgKVxuICAgIH1cbiAgfVxuICBkZXNlcmlhbGl6ZShmaWVsZHM6IG9iamVjdCwgZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpIHtcbiAgICBzdXBlci5kZXNlcmlhbGl6ZShmaWVsZHMsIGVuY29kaW5nKVxuICAgIHRoaXMubm9kZUlEID0gc2VyaWFsaXphdGlvbi5kZWNvZGVyKFxuICAgICAgZmllbGRzW1wibm9kZUlEXCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBcIm5vZGVJRFwiLFxuICAgICAgXCJCdWZmZXJcIixcbiAgICAgIDIwXG4gICAgKVxuICAgIHRoaXMuc3RhcnRUaW1lID0gc2VyaWFsaXphdGlvbi5kZWNvZGVyKFxuICAgICAgZmllbGRzW1wic3RhcnRUaW1lXCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBcImRlY2ltYWxTdHJpbmdcIixcbiAgICAgIFwiQnVmZmVyXCIsXG4gICAgICA4XG4gICAgKVxuICAgIHRoaXMuZW5kVGltZSA9IHNlcmlhbGl6YXRpb24uZGVjb2RlcihcbiAgICAgIGZpZWxkc1tcImVuZFRpbWVcIl0sXG4gICAgICBlbmNvZGluZyxcbiAgICAgIFwiZGVjaW1hbFN0cmluZ1wiLFxuICAgICAgXCJCdWZmZXJcIixcbiAgICAgIDhcbiAgICApXG4gIH1cblxuICBwcm90ZWN0ZWQgbm9kZUlEOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoMjApXG4gIHByb3RlY3RlZCBzdGFydFRpbWU6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYyg4KVxuICBwcm90ZWN0ZWQgZW5kVGltZTogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDgpXG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSBmb3IgdGhlIHN0YWtlIGFtb3VudC5cbiAgICovXG4gIGdldE5vZGVJRCgpOiBCdWZmZXIge1xuICAgIHJldHVybiB0aGlzLm5vZGVJRFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSBzdHJpbmcgZm9yIHRoZSBub2RlSUQgYW1vdW50LlxuICAgKi9cbiAgZ2V0Tm9kZUlEU3RyaW5nKCk6IHN0cmluZyB7XG4gICAgcmV0dXJuIGJ1ZmZlclRvTm9kZUlEU3RyaW5nKHRoaXMubm9kZUlEKVxuICB9XG4gIC8qKlxuICAgKiBSZXR1cm5zIGEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9pbmR1dG55L2JuLmpzL3xCTn0gZm9yIHRoZSBzdGFrZSBhbW91bnQuXG4gICAqL1xuICBnZXRTdGFydFRpbWUoKSB7XG4gICAgcmV0dXJuIGJpbnRvb2xzLmZyb21CdWZmZXJUb0JOKHRoaXMuc3RhcnRUaW1lKVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2luZHV0bnkvYm4uanMvfEJOfSBmb3IgdGhlIHN0YWtlIGFtb3VudC5cbiAgICovXG4gIGdldEVuZFRpbWUoKSB7XG4gICAgcmV0dXJuIGJpbnRvb2xzLmZyb21CdWZmZXJUb0JOKHRoaXMuZW5kVGltZSlcbiAgfVxuXG4gIGZyb21CdWZmZXIoYnl0ZXM6IEJ1ZmZlciwgb2Zmc2V0OiBudW1iZXIgPSAwKTogbnVtYmVyIHtcbiAgICBvZmZzZXQgPSBzdXBlci5mcm9tQnVmZmVyKGJ5dGVzLCBvZmZzZXQpXG4gICAgdGhpcy5ub2RlSUQgPSBiaW50b29scy5jb3B5RnJvbShieXRlcywgb2Zmc2V0LCBvZmZzZXQgKyAyMClcbiAgICBvZmZzZXQgKz0gMjBcbiAgICB0aGlzLnN0YXJ0VGltZSA9IGJpbnRvb2xzLmNvcHlGcm9tKGJ5dGVzLCBvZmZzZXQsIG9mZnNldCArIDgpXG4gICAgb2Zmc2V0ICs9IDhcbiAgICB0aGlzLmVuZFRpbWUgPSBiaW50b29scy5jb3B5RnJvbShieXRlcywgb2Zmc2V0LCBvZmZzZXQgKyA4KVxuICAgIG9mZnNldCArPSA4XG4gICAgcmV0dXJuIG9mZnNldFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSByZXByZXNlbnRhdGlvbiBvZiB0aGUgW1tWYWxpZGF0b3JUeF1dLlxuICAgKi9cbiAgdG9CdWZmZXIoKTogQnVmZmVyIHtcbiAgICBjb25zdCBzdXBlcmJ1ZmY6IEJ1ZmZlciA9IHN1cGVyLnRvQnVmZmVyKClcbiAgICBjb25zdCBic2l6ZTogbnVtYmVyID1cbiAgICAgIHN1cGVyYnVmZi5sZW5ndGggK1xuICAgICAgdGhpcy5ub2RlSUQubGVuZ3RoICtcbiAgICAgIHRoaXMuc3RhcnRUaW1lLmxlbmd0aCArXG4gICAgICB0aGlzLmVuZFRpbWUubGVuZ3RoXG4gICAgcmV0dXJuIEJ1ZmZlci5jb25jYXQoXG4gICAgICBbc3VwZXJidWZmLCB0aGlzLm5vZGVJRCwgdGhpcy5zdGFydFRpbWUsIHRoaXMuZW5kVGltZV0sXG4gICAgICBic2l6ZVxuICAgIClcbiAgfVxuXG4gIGNvbnN0cnVjdG9yKFxuICAgIG5ldHdvcmtJRDogbnVtYmVyLFxuICAgIGJsb2NrY2hhaW5JRDogQnVmZmVyLFxuICAgIG91dHM6IFRyYW5zZmVyYWJsZU91dHB1dFtdLFxuICAgIGluczogVHJhbnNmZXJhYmxlSW5wdXRbXSxcbiAgICBtZW1vPzogQnVmZmVyLFxuICAgIG5vZGVJRD86IEJ1ZmZlcixcbiAgICBzdGFydFRpbWU/OiBCTixcbiAgICBlbmRUaW1lPzogQk5cbiAgKSB7XG4gICAgc3VwZXIobmV0d29ya0lELCBibG9ja2NoYWluSUQsIG91dHMsIGlucywgbWVtbylcbiAgICB0aGlzLm5vZGVJRCA9IG5vZGVJRFxuICAgIHRoaXMuc3RhcnRUaW1lID0gYmludG9vbHMuZnJvbUJOVG9CdWZmZXIoc3RhcnRUaW1lLCA4KVxuICAgIHRoaXMuZW5kVGltZSA9IGJpbnRvb2xzLmZyb21CTlRvQnVmZmVyKGVuZFRpbWUsIDgpXG4gIH1cbn1cblxuZXhwb3J0IGFic3RyYWN0IGNsYXNzIFdlaWdodGVkVmFsaWRhdG9yVHggZXh0ZW5kcyBWYWxpZGF0b3JUeCB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIldlaWdodGVkVmFsaWRhdG9yVHhcIlxuICBwcm90ZWN0ZWQgX3R5cGVJRCA9IHVuZGVmaW5lZFxuXG4gIHNlcmlhbGl6ZShlbmNvZGluZzogU2VyaWFsaXplZEVuY29kaW5nID0gXCJoZXhcIik6IG9iamVjdCB7XG4gICAgbGV0IGZpZWxkczogb2JqZWN0ID0gc3VwZXIuc2VyaWFsaXplKGVuY29kaW5nKVxuICAgIHJldHVybiB7XG4gICAgICAuLi5maWVsZHMsXG4gICAgICB3ZWlnaHQ6IHNlcmlhbGl6YXRpb24uZW5jb2RlcihcbiAgICAgICAgdGhpcy53ZWlnaHQsXG4gICAgICAgIGVuY29kaW5nLFxuICAgICAgICBcIkJ1ZmZlclwiLFxuICAgICAgICBcImRlY2ltYWxTdHJpbmdcIlxuICAgICAgKVxuICAgIH1cbiAgfVxuICBkZXNlcmlhbGl6ZShmaWVsZHM6IG9iamVjdCwgZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpIHtcbiAgICBzdXBlci5kZXNlcmlhbGl6ZShmaWVsZHMsIGVuY29kaW5nKVxuICAgIHRoaXMud2VpZ2h0ID0gc2VyaWFsaXphdGlvbi5kZWNvZGVyKFxuICAgICAgZmllbGRzW1wid2VpZ2h0XCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBcImRlY2ltYWxTdHJpbmdcIixcbiAgICAgIFwiQnVmZmVyXCIsXG4gICAgICA4XG4gICAgKVxuICB9XG5cbiAgcHJvdGVjdGVkIHdlaWdodDogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDgpXG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2luZHV0bnkvYm4uanMvfEJOfSBmb3IgdGhlIHN0YWtlIGFtb3VudC5cbiAgICovXG4gIGdldFdlaWdodCgpOiBCTiB7XG4gICAgcmV0dXJuIGJpbnRvb2xzLmZyb21CdWZmZXJUb0JOKHRoaXMud2VpZ2h0KVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSBmb3IgdGhlIHN0YWtlIGFtb3VudC5cbiAgICovXG4gIGdldFdlaWdodEJ1ZmZlcigpOiBCdWZmZXIge1xuICAgIHJldHVybiB0aGlzLndlaWdodFxuICB9XG5cbiAgZnJvbUJ1ZmZlcihieXRlczogQnVmZmVyLCBvZmZzZXQ6IG51bWJlciA9IDApOiBudW1iZXIge1xuICAgIG9mZnNldCA9IHN1cGVyLmZyb21CdWZmZXIoYnl0ZXMsIG9mZnNldClcbiAgICB0aGlzLndlaWdodCA9IGJpbnRvb2xzLmNvcHlGcm9tKGJ5dGVzLCBvZmZzZXQsIG9mZnNldCArIDgpXG4gICAgb2Zmc2V0ICs9IDhcbiAgICByZXR1cm4gb2Zmc2V0XG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IHJlcHJlc2VudGF0aW9uIG9mIHRoZSBbW0FkZFN1Ym5ldFZhbGlkYXRvclR4XV0uXG4gICAqL1xuICB0b0J1ZmZlcigpOiBCdWZmZXIge1xuICAgIGNvbnN0IHN1cGVyYnVmZjogQnVmZmVyID0gc3VwZXIudG9CdWZmZXIoKVxuICAgIHJldHVybiBCdWZmZXIuY29uY2F0KFtzdXBlcmJ1ZmYsIHRoaXMud2VpZ2h0XSlcbiAgfVxuXG4gIC8qKlxuICAgKiBDbGFzcyByZXByZXNlbnRpbmcgYW4gdW5zaWduZWQgQWRkU3VibmV0VmFsaWRhdG9yVHggdHJhbnNhY3Rpb24uXG4gICAqXG4gICAqIEBwYXJhbSBuZXR3b3JrSUQgT3B0aW9uYWwuIE5ldHdvcmtpZCwgW1tEZWZhdWx0TmV0d29ya0lEXV1cbiAgICogQHBhcmFtIGJsb2NrY2hhaW5JRCBPcHRpb25hbC4gQmxvY2tjaGFpbmlkLCBkZWZhdWx0IEJ1ZmZlci5hbGxvYygzMiwgMTYpXG4gICAqIEBwYXJhbSBvdXRzIE9wdGlvbmFsLiBBcnJheSBvZiB0aGUgW1tUcmFuc2ZlcmFibGVPdXRwdXRdXXNcbiAgICogQHBhcmFtIGlucyBPcHRpb25hbC4gQXJyYXkgb2YgdGhlIFtbVHJhbnNmZXJhYmxlSW5wdXRdXXNcbiAgICogQHBhcmFtIG1lbW8gT3B0aW9uYWwuIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IGZvciB0aGUgbWVtbyBmaWVsZFxuICAgKiBAcGFyYW0gbm9kZUlEIE9wdGlvbmFsLiBUaGUgbm9kZSBJRCBvZiB0aGUgdmFsaWRhdG9yIGJlaW5nIGFkZGVkLlxuICAgKiBAcGFyYW0gc3RhcnRUaW1lIE9wdGlvbmFsLiBUaGUgVW5peCB0aW1lIHdoZW4gdGhlIHZhbGlkYXRvciBzdGFydHMgdmFsaWRhdGluZyB0aGUgUHJpbWFyeSBOZXR3b3JrLlxuICAgKiBAcGFyYW0gZW5kVGltZSBPcHRpb25hbC4gVGhlIFVuaXggdGltZSB3aGVuIHRoZSB2YWxpZGF0b3Igc3RvcHMgdmFsaWRhdGluZyB0aGUgUHJpbWFyeSBOZXR3b3JrIChhbmQgc3Rha2VkIEFWQVggaXMgcmV0dXJuZWQpLlxuICAgKiBAcGFyYW0gd2VpZ2h0IE9wdGlvbmFsLiBUaGUgYW1vdW50IG9mIG5BVkFYIHRoZSB2YWxpZGF0b3IgaXMgc3Rha2luZy5cbiAgICovXG4gIGNvbnN0cnVjdG9yKFxuICAgIG5ldHdvcmtJRDogbnVtYmVyID0gRGVmYXVsdE5ldHdvcmtJRCxcbiAgICBibG9ja2NoYWluSUQ6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYygzMiwgMTYpLFxuICAgIG91dHM6IFRyYW5zZmVyYWJsZU91dHB1dFtdID0gdW5kZWZpbmVkLFxuICAgIGluczogVHJhbnNmZXJhYmxlSW5wdXRbXSA9IHVuZGVmaW5lZCxcbiAgICBtZW1vOiBCdWZmZXIgPSB1bmRlZmluZWQsXG4gICAgbm9kZUlEOiBCdWZmZXIgPSB1bmRlZmluZWQsXG4gICAgc3RhcnRUaW1lOiBCTiA9IHVuZGVmaW5lZCxcbiAgICBlbmRUaW1lOiBCTiA9IHVuZGVmaW5lZCxcbiAgICB3ZWlnaHQ6IEJOID0gdW5kZWZpbmVkXG4gICkge1xuICAgIHN1cGVyKG5ldHdvcmtJRCwgYmxvY2tjaGFpbklELCBvdXRzLCBpbnMsIG1lbW8sIG5vZGVJRCwgc3RhcnRUaW1lLCBlbmRUaW1lKVxuICAgIGlmICh0eXBlb2Ygd2VpZ2h0ICE9PSB1bmRlZmluZWQpIHtcbiAgICAgIHRoaXMud2VpZ2h0ID0gYmludG9vbHMuZnJvbUJOVG9CdWZmZXIod2VpZ2h0LCA4KVxuICAgIH1cbiAgfVxufVxuXG4vKipcbiAqIENsYXNzIHJlcHJlc2VudGluZyBhbiB1bnNpZ25lZCBBZGREZWxlZ2F0b3JUeCB0cmFuc2FjdGlvbi5cbiAqL1xuZXhwb3J0IGNsYXNzIEFkZERlbGVnYXRvclR4IGV4dGVuZHMgV2VpZ2h0ZWRWYWxpZGF0b3JUeCB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIkFkZERlbGVnYXRvclR4XCJcbiAgcHJvdGVjdGVkIF90eXBlSUQgPSBQbGF0Zm9ybVZNQ29uc3RhbnRzLkFERERFTEVHQVRPUlRYXG5cbiAgc2VyaWFsaXplKGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKTogb2JqZWN0IHtcbiAgICBsZXQgZmllbGRzOiBvYmplY3QgPSBzdXBlci5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgcmV0dXJuIHtcbiAgICAgIC4uLmZpZWxkcyxcbiAgICAgIHN0YWtlT3V0czogdGhpcy5zdGFrZU91dHMubWFwKChzKSA9PiBzLnNlcmlhbGl6ZShlbmNvZGluZykpLFxuICAgICAgcmV3YXJkT3duZXJzOiB0aGlzLnJld2FyZE93bmVycy5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgfVxuICB9XG4gIGRlc2VyaWFsaXplKGZpZWxkczogb2JqZWN0LCBlbmNvZGluZzogU2VyaWFsaXplZEVuY29kaW5nID0gXCJoZXhcIikge1xuICAgIHN1cGVyLmRlc2VyaWFsaXplKGZpZWxkcywgZW5jb2RpbmcpXG4gICAgdGhpcy5zdGFrZU91dHMgPSBmaWVsZHNbXCJzdGFrZU91dHNcIl0ubWFwKChzOiBvYmplY3QpID0+IHtcbiAgICAgIGxldCB4ZmVyb3V0OiBUcmFuc2ZlcmFibGVPdXRwdXQgPSBuZXcgVHJhbnNmZXJhYmxlT3V0cHV0KClcbiAgICAgIHhmZXJvdXQuZGVzZXJpYWxpemUocywgZW5jb2RpbmcpXG4gICAgICByZXR1cm4geGZlcm91dFxuICAgIH0pXG4gICAgdGhpcy5yZXdhcmRPd25lcnMgPSBuZXcgUGFyc2VhYmxlT3V0cHV0KClcbiAgICB0aGlzLnJld2FyZE93bmVycy5kZXNlcmlhbGl6ZShmaWVsZHNbXCJyZXdhcmRPd25lcnNcIl0sIGVuY29kaW5nKVxuICB9XG5cbiAgcHJvdGVjdGVkIHN0YWtlT3V0czogVHJhbnNmZXJhYmxlT3V0cHV0W10gPSBbXVxuICBwcm90ZWN0ZWQgcmV3YXJkT3duZXJzOiBQYXJzZWFibGVPdXRwdXQgPSB1bmRlZmluZWRcblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgaWQgb2YgdGhlIFtbQWRkRGVsZWdhdG9yVHhdXVxuICAgKi9cbiAgZ2V0VHhUeXBlKCk6IG51bWJlciB7XG4gICAgcmV0dXJuIHRoaXMuX3R5cGVJRFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2luZHV0bnkvYm4uanMvfEJOfSBmb3IgdGhlIHN0YWtlIGFtb3VudC5cbiAgICovXG4gIGdldFN0YWtlQW1vdW50KCk6IEJOIHtcbiAgICByZXR1cm4gdGhpcy5nZXRXZWlnaHQoKVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSBmb3IgdGhlIHN0YWtlIGFtb3VudC5cbiAgICovXG4gIGdldFN0YWtlQW1vdW50QnVmZmVyKCk6IEJ1ZmZlciB7XG4gICAgcmV0dXJuIHRoaXMud2VpZ2h0XG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgYXJyYXkgb2Ygb3V0cHV0cyBiZWluZyBzdGFrZWQuXG4gICAqL1xuICBnZXRTdGFrZU91dHMoKTogVHJhbnNmZXJhYmxlT3V0cHV0W10ge1xuICAgIHJldHVybiB0aGlzLnN0YWtlT3V0c1xuICB9XG5cbiAgLyoqXG4gICAqIFNob3VsZCBtYXRjaCBzdGFrZUFtb3VudC4gVXNlZCBpbiBzYW5pdHkgY2hlY2tpbmcuXG4gICAqL1xuICBnZXRTdGFrZU91dHNUb3RhbCgpOiBCTiB7XG4gICAgbGV0IHZhbDogQk4gPSBuZXcgQk4oMClcbiAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgdGhpcy5zdGFrZU91dHMubGVuZ3RoOyBpKyspIHtcbiAgICAgIHZhbCA9IHZhbC5hZGQoXG4gICAgICAgICh0aGlzLnN0YWtlT3V0c1tgJHtpfWBdLmdldE91dHB1dCgpIGFzIEFtb3VudE91dHB1dCkuZ2V0QW1vdW50KClcbiAgICAgIClcbiAgICB9XG4gICAgcmV0dXJuIHZhbFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSBmb3IgdGhlIHJld2FyZCBhZGRyZXNzLlxuICAgKi9cbiAgZ2V0UmV3YXJkT3duZXJzKCk6IFBhcnNlYWJsZU91dHB1dCB7XG4gICAgcmV0dXJuIHRoaXMucmV3YXJkT3duZXJzXG4gIH1cblxuICBnZXRUb3RhbE91dHMoKTogVHJhbnNmZXJhYmxlT3V0cHV0W10ge1xuICAgIHJldHVybiBbLi4uKHRoaXMuZ2V0T3V0cygpIGFzIFRyYW5zZmVyYWJsZU91dHB1dFtdKSwgLi4udGhpcy5nZXRTdGFrZU91dHMoKV1cbiAgfVxuXG4gIGZyb21CdWZmZXIoYnl0ZXM6IEJ1ZmZlciwgb2Zmc2V0OiBudW1iZXIgPSAwKTogbnVtYmVyIHtcbiAgICBvZmZzZXQgPSBzdXBlci5mcm9tQnVmZmVyKGJ5dGVzLCBvZmZzZXQpXG4gICAgY29uc3QgbnVtc3Rha2VvdXRzID0gYmludG9vbHMuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgNClcbiAgICBvZmZzZXQgKz0gNFxuICAgIGNvbnN0IG91dGNvdW50OiBudW1iZXIgPSBudW1zdGFrZW91dHMucmVhZFVJbnQzMkJFKDApXG4gICAgdGhpcy5zdGFrZU91dHMgPSBbXVxuICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCBvdXRjb3VudDsgaSsrKSB7XG4gICAgICBjb25zdCB4ZmVyb3V0OiBUcmFuc2ZlcmFibGVPdXRwdXQgPSBuZXcgVHJhbnNmZXJhYmxlT3V0cHV0KClcbiAgICAgIG9mZnNldCA9IHhmZXJvdXQuZnJvbUJ1ZmZlcihieXRlcywgb2Zmc2V0KVxuICAgICAgdGhpcy5zdGFrZU91dHMucHVzaCh4ZmVyb3V0KVxuICAgIH1cbiAgICB0aGlzLnJld2FyZE93bmVycyA9IG5ldyBQYXJzZWFibGVPdXRwdXQoKVxuICAgIG9mZnNldCA9IHRoaXMucmV3YXJkT3duZXJzLmZyb21CdWZmZXIoYnl0ZXMsIG9mZnNldClcbiAgICByZXR1cm4gb2Zmc2V0XG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IHJlcHJlc2VudGF0aW9uIG9mIHRoZSBbW0FkZERlbGVnYXRvclR4XV0uXG4gICAqL1xuICB0b0J1ZmZlcigpOiBCdWZmZXIge1xuICAgIGNvbnN0IHN1cGVyYnVmZjogQnVmZmVyID0gc3VwZXIudG9CdWZmZXIoKVxuICAgIGxldCBic2l6ZTogbnVtYmVyID0gc3VwZXJidWZmLmxlbmd0aFxuICAgIGNvbnN0IG51bW91dHM6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYyg0KVxuICAgIG51bW91dHMud3JpdGVVSW50MzJCRSh0aGlzLnN0YWtlT3V0cy5sZW5ndGgsIDApXG4gICAgbGV0IGJhcnI6IEJ1ZmZlcltdID0gW3N1cGVyLnRvQnVmZmVyKCksIG51bW91dHNdXG4gICAgYnNpemUgKz0gbnVtb3V0cy5sZW5ndGhcbiAgICB0aGlzLnN0YWtlT3V0cyA9IHRoaXMuc3Rha2VPdXRzLnNvcnQoVHJhbnNmZXJhYmxlT3V0cHV0LmNvbXBhcmF0b3IoKSlcbiAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgdGhpcy5zdGFrZU91dHMubGVuZ3RoOyBpKyspIHtcbiAgICAgIGxldCBvdXQ6IEJ1ZmZlciA9IHRoaXMuc3Rha2VPdXRzW2Ake2l9YF0udG9CdWZmZXIoKVxuICAgICAgYmFyci5wdXNoKG91dClcbiAgICAgIGJzaXplICs9IG91dC5sZW5ndGhcbiAgICB9XG4gICAgbGV0IHJvOiBCdWZmZXIgPSB0aGlzLnJld2FyZE93bmVycy50b0J1ZmZlcigpXG4gICAgYmFyci5wdXNoKHJvKVxuICAgIGJzaXplICs9IHJvLmxlbmd0aFxuICAgIHJldHVybiBCdWZmZXIuY29uY2F0KGJhcnIsIGJzaXplKVxuICB9XG5cbiAgY2xvbmUoKTogdGhpcyB7XG4gICAgbGV0IG5ld2Jhc2U6IEFkZERlbGVnYXRvclR4ID0gbmV3IEFkZERlbGVnYXRvclR4KClcbiAgICBuZXdiYXNlLmZyb21CdWZmZXIodGhpcy50b0J1ZmZlcigpKVxuICAgIHJldHVybiBuZXdiYXNlIGFzIHRoaXNcbiAgfVxuXG4gIGNyZWF0ZSguLi5hcmdzOiBhbnlbXSk6IHRoaXMge1xuICAgIHJldHVybiBuZXcgQWRkRGVsZWdhdG9yVHgoLi4uYXJncykgYXMgdGhpc1xuICB9XG5cbiAgLyoqXG4gICAqIENsYXNzIHJlcHJlc2VudGluZyBhbiB1bnNpZ25lZCBBZGREZWxlZ2F0b3JUeCB0cmFuc2FjdGlvbi5cbiAgICpcbiAgICogQHBhcmFtIG5ldHdvcmtJRCBPcHRpb25hbC4gTmV0d29ya2lkLCBbW0RlZmF1bHROZXR3b3JrSURdXVxuICAgKiBAcGFyYW0gYmxvY2tjaGFpbklEIE9wdGlvbmFsLiBCbG9ja2NoYWluaWQsIGRlZmF1bHQgQnVmZmVyLmFsbG9jKDMyLCAxNilcbiAgICogQHBhcmFtIG91dHMgT3B0aW9uYWwuIEFycmF5IG9mIHRoZSBbW1RyYW5zZmVyYWJsZU91dHB1dF1dc1xuICAgKiBAcGFyYW0gaW5zIE9wdGlvbmFsLiBBcnJheSBvZiB0aGUgW1tUcmFuc2ZlcmFibGVJbnB1dF1dc1xuICAgKiBAcGFyYW0gbWVtbyBPcHRpb25hbC4ge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gZm9yIHRoZSBtZW1vIGZpZWxkXG4gICAqIEBwYXJhbSBub2RlSUQgT3B0aW9uYWwuIFRoZSBub2RlIElEIG9mIHRoZSB2YWxpZGF0b3IgYmVpbmcgYWRkZWQuXG4gICAqIEBwYXJhbSBzdGFydFRpbWUgT3B0aW9uYWwuIFRoZSBVbml4IHRpbWUgd2hlbiB0aGUgdmFsaWRhdG9yIHN0YXJ0cyB2YWxpZGF0aW5nIHRoZSBQcmltYXJ5IE5ldHdvcmsuXG4gICAqIEBwYXJhbSBlbmRUaW1lIE9wdGlvbmFsLiBUaGUgVW5peCB0aW1lIHdoZW4gdGhlIHZhbGlkYXRvciBzdG9wcyB2YWxpZGF0aW5nIHRoZSBQcmltYXJ5IE5ldHdvcmsgKGFuZCBzdGFrZWQgQVZBWCBpcyByZXR1cm5lZCkuXG4gICAqIEBwYXJhbSBzdGFrZUFtb3VudCBPcHRpb25hbC4gVGhlIGFtb3VudCBvZiBuQVZBWCB0aGUgdmFsaWRhdG9yIGlzIHN0YWtpbmcuXG4gICAqIEBwYXJhbSBzdGFrZU91dHMgT3B0aW9uYWwuIFRoZSBvdXRwdXRzIHVzZWQgaW4gcGF5aW5nIHRoZSBzdGFrZS5cbiAgICogQHBhcmFtIHJld2FyZE93bmVycyBPcHRpb25hbC4gVGhlIFtbUGFyc2VhYmxlT3V0cHV0XV0gY29udGFpbmluZyBhIFtbU0VDUE93bmVyT3V0cHV0XV0gZm9yIHRoZSByZXdhcmRzLlxuICAgKi9cbiAgY29uc3RydWN0b3IoXG4gICAgbmV0d29ya0lEOiBudW1iZXIgPSBEZWZhdWx0TmV0d29ya0lELFxuICAgIGJsb2NrY2hhaW5JRDogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDMyLCAxNiksXG4gICAgb3V0czogVHJhbnNmZXJhYmxlT3V0cHV0W10gPSB1bmRlZmluZWQsXG4gICAgaW5zOiBUcmFuc2ZlcmFibGVJbnB1dFtdID0gdW5kZWZpbmVkLFxuICAgIG1lbW86IEJ1ZmZlciA9IHVuZGVmaW5lZCxcbiAgICBub2RlSUQ6IEJ1ZmZlciA9IHVuZGVmaW5lZCxcbiAgICBzdGFydFRpbWU6IEJOID0gdW5kZWZpbmVkLFxuICAgIGVuZFRpbWU6IEJOID0gdW5kZWZpbmVkLFxuICAgIHN0YWtlQW1vdW50OiBCTiA9IHVuZGVmaW5lZCxcbiAgICBzdGFrZU91dHM6IFRyYW5zZmVyYWJsZU91dHB1dFtdID0gdW5kZWZpbmVkLFxuICAgIHJld2FyZE93bmVyczogUGFyc2VhYmxlT3V0cHV0ID0gdW5kZWZpbmVkXG4gICkge1xuICAgIHN1cGVyKFxuICAgICAgbmV0d29ya0lELFxuICAgICAgYmxvY2tjaGFpbklELFxuICAgICAgb3V0cyxcbiAgICAgIGlucyxcbiAgICAgIG1lbW8sXG4gICAgICBub2RlSUQsXG4gICAgICBzdGFydFRpbWUsXG4gICAgICBlbmRUaW1lLFxuICAgICAgc3Rha2VBbW91bnRcbiAgICApXG4gICAgaWYgKHR5cGVvZiBzdGFrZU91dHMgIT09IHVuZGVmaW5lZCkge1xuICAgICAgdGhpcy5zdGFrZU91dHMgPSBzdGFrZU91dHNcbiAgICB9XG4gICAgdGhpcy5yZXdhcmRPd25lcnMgPSByZXdhcmRPd25lcnNcbiAgfVxufVxuXG5leHBvcnQgY2xhc3MgQWRkVmFsaWRhdG9yVHggZXh0ZW5kcyBBZGREZWxlZ2F0b3JUeCB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIkFkZFZhbGlkYXRvclR4XCJcbiAgcHJvdGVjdGVkIF90eXBlSUQgPSBQbGF0Zm9ybVZNQ29uc3RhbnRzLkFERFZBTElEQVRPUlRYXG5cbiAgc2VyaWFsaXplKGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKTogb2JqZWN0IHtcbiAgICBsZXQgZmllbGRzOiBvYmplY3QgPSBzdXBlci5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgcmV0dXJuIHtcbiAgICAgIC4uLmZpZWxkcyxcbiAgICAgIGRlbGVnYXRpb25GZWU6IHNlcmlhbGl6YXRpb24uZW5jb2RlcihcbiAgICAgICAgdGhpcy5nZXREZWxlZ2F0aW9uRmVlQnVmZmVyKCksXG4gICAgICAgIGVuY29kaW5nLFxuICAgICAgICBcIkJ1ZmZlclwiLFxuICAgICAgICBcImRlY2ltYWxTdHJpbmdcIixcbiAgICAgICAgNFxuICAgICAgKVxuICAgIH1cbiAgfVxuICBkZXNlcmlhbGl6ZShmaWVsZHM6IG9iamVjdCwgZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpIHtcbiAgICBzdXBlci5kZXNlcmlhbGl6ZShmaWVsZHMsIGVuY29kaW5nKVxuICAgIGxldCBkYnVmZjogQnVmZmVyID0gc2VyaWFsaXphdGlvbi5kZWNvZGVyKFxuICAgICAgZmllbGRzW1wiZGVsZWdhdGlvbkZlZVwiXSxcbiAgICAgIGVuY29kaW5nLFxuICAgICAgXCJkZWNpbWFsU3RyaW5nXCIsXG4gICAgICBcIkJ1ZmZlclwiLFxuICAgICAgNFxuICAgIClcbiAgICB0aGlzLmRlbGVnYXRpb25GZWUgPVxuICAgICAgZGJ1ZmYucmVhZFVJbnQzMkJFKDApIC8gQWRkVmFsaWRhdG9yVHguZGVsZWdhdG9yTXVsdGlwbGllclxuICB9XG5cbiAgcHJvdGVjdGVkIGRlbGVnYXRpb25GZWU6IG51bWJlciA9IDBcbiAgcHJpdmF0ZSBzdGF0aWMgZGVsZWdhdG9yTXVsdGlwbGllcjogbnVtYmVyID0gMTAwMDBcblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgaWQgb2YgdGhlIFtbQWRkVmFsaWRhdG9yVHhdXVxuICAgKi9cbiAgZ2V0VHhUeXBlKCk6IG51bWJlciB7XG4gICAgcmV0dXJuIHRoaXMuX3R5cGVJRFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIGRlbGVnYXRpb24gZmVlIChyZXByZXNlbnRzIGEgcGVyY2VudGFnZSBmcm9tIDAgdG8gMTAwKTtcbiAgICovXG4gIGdldERlbGVnYXRpb25GZWUoKTogbnVtYmVyIHtcbiAgICByZXR1cm4gdGhpcy5kZWxlZ2F0aW9uRmVlXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgYmluYXJ5IHJlcHJlc2VudGF0aW9uIG9mIHRoZSBkZWxlZ2F0aW9uIGZlZSBhcyBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9LlxuICAgKi9cbiAgZ2V0RGVsZWdhdGlvbkZlZUJ1ZmZlcigpOiBCdWZmZXIge1xuICAgIGxldCBkQnVmZjogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDQpXG4gICAgbGV0IGJ1ZmZudW06IG51bWJlciA9XG4gICAgICBwYXJzZUZsb2F0KHRoaXMuZGVsZWdhdGlvbkZlZS50b0ZpeGVkKDQpKSAqXG4gICAgICBBZGRWYWxpZGF0b3JUeC5kZWxlZ2F0b3JNdWx0aXBsaWVyXG4gICAgZEJ1ZmYud3JpdGVVSW50MzJCRShidWZmbnVtLCAwKVxuICAgIHJldHVybiBkQnVmZlxuICB9XG5cbiAgZnJvbUJ1ZmZlcihieXRlczogQnVmZmVyLCBvZmZzZXQ6IG51bWJlciA9IDApOiBudW1iZXIge1xuICAgIG9mZnNldCA9IHN1cGVyLmZyb21CdWZmZXIoYnl0ZXMsIG9mZnNldClcbiAgICBsZXQgZGJ1ZmY6IEJ1ZmZlciA9IGJpbnRvb2xzLmNvcHlGcm9tKGJ5dGVzLCBvZmZzZXQsIG9mZnNldCArIDQpXG4gICAgb2Zmc2V0ICs9IDRcbiAgICB0aGlzLmRlbGVnYXRpb25GZWUgPVxuICAgICAgZGJ1ZmYucmVhZFVJbnQzMkJFKDApIC8gQWRkVmFsaWRhdG9yVHguZGVsZWdhdG9yTXVsdGlwbGllclxuICAgIHJldHVybiBvZmZzZXRcbiAgfVxuXG4gIHRvQnVmZmVyKCk6IEJ1ZmZlciB7XG4gICAgbGV0IHN1cGVyQnVmZjogQnVmZmVyID0gc3VwZXIudG9CdWZmZXIoKVxuICAgIGxldCBmZWVCdWZmOiBCdWZmZXIgPSB0aGlzLmdldERlbGVnYXRpb25GZWVCdWZmZXIoKVxuICAgIHJldHVybiBCdWZmZXIuY29uY2F0KFtzdXBlckJ1ZmYsIGZlZUJ1ZmZdKVxuICB9XG5cbiAgLyoqXG4gICAqIENsYXNzIHJlcHJlc2VudGluZyBhbiB1bnNpZ25lZCBBZGRWYWxpZGF0b3JUeCB0cmFuc2FjdGlvbi5cbiAgICpcbiAgICogQHBhcmFtIG5ldHdvcmtJRCBPcHRpb25hbC4gTmV0d29ya2lkLCBbW0RlZmF1bHROZXR3b3JrSURdXVxuICAgKiBAcGFyYW0gYmxvY2tjaGFpbklEIE9wdGlvbmFsLiBCbG9ja2NoYWluaWQsIGRlZmF1bHQgQnVmZmVyLmFsbG9jKDMyLCAxNilcbiAgICogQHBhcmFtIG91dHMgT3B0aW9uYWwuIEFycmF5IG9mIHRoZSBbW1RyYW5zZmVyYWJsZU91dHB1dF1dc1xuICAgKiBAcGFyYW0gaW5zIE9wdGlvbmFsLiBBcnJheSBvZiB0aGUgW1tUcmFuc2ZlcmFibGVJbnB1dF1dc1xuICAgKiBAcGFyYW0gbWVtbyBPcHRpb25hbC4ge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gZm9yIHRoZSBtZW1vIGZpZWxkXG4gICAqIEBwYXJhbSBub2RlSUQgT3B0aW9uYWwuIFRoZSBub2RlIElEIG9mIHRoZSB2YWxpZGF0b3IgYmVpbmcgYWRkZWQuXG4gICAqIEBwYXJhbSBzdGFydFRpbWUgT3B0aW9uYWwuIFRoZSBVbml4IHRpbWUgd2hlbiB0aGUgdmFsaWRhdG9yIHN0YXJ0cyB2YWxpZGF0aW5nIHRoZSBQcmltYXJ5IE5ldHdvcmsuXG4gICAqIEBwYXJhbSBlbmRUaW1lIE9wdGlvbmFsLiBUaGUgVW5peCB0aW1lIHdoZW4gdGhlIHZhbGlkYXRvciBzdG9wcyB2YWxpZGF0aW5nIHRoZSBQcmltYXJ5IE5ldHdvcmsgKGFuZCBzdGFrZWQgQVZBWCBpcyByZXR1cm5lZCkuXG4gICAqIEBwYXJhbSBzdGFrZUFtb3VudCBPcHRpb25hbC4gVGhlIGFtb3VudCBvZiBuQVZBWCB0aGUgdmFsaWRhdG9yIGlzIHN0YWtpbmcuXG4gICAqIEBwYXJhbSBzdGFrZU91dHMgT3B0aW9uYWwuIFRoZSBvdXRwdXRzIHVzZWQgaW4gcGF5aW5nIHRoZSBzdGFrZS5cbiAgICogQHBhcmFtIHJld2FyZE93bmVycyBPcHRpb25hbC4gVGhlIFtbUGFyc2VhYmxlT3V0cHV0XV0gY29udGFpbmluZyB0aGUgW1tTRUNQT3duZXJPdXRwdXRdXSBmb3IgdGhlIHJld2FyZHMuXG4gICAqIEBwYXJhbSBkZWxlZ2F0aW9uRmVlIE9wdGlvbmFsLiBUaGUgcGVyY2VudCBmZWUgdGhpcyB2YWxpZGF0b3IgY2hhcmdlcyB3aGVuIG90aGVycyBkZWxlZ2F0ZSBzdGFrZSB0byB0aGVtLlxuICAgKiBVcCB0byA0IGRlY2ltYWwgcGxhY2VzIGFsbG93ZWQ7IGFkZGl0aW9uYWwgZGVjaW1hbCBwbGFjZXMgYXJlIGlnbm9yZWQuIE11c3QgYmUgYmV0d2VlbiAwIGFuZCAxMDAsIGluY2x1c2l2ZS5cbiAgICogRm9yIGV4YW1wbGUsIGlmIGRlbGVnYXRpb25GZWVSYXRlIGlzIDEuMjM0NSBhbmQgc29tZW9uZSBkZWxlZ2F0ZXMgdG8gdGhpcyB2YWxpZGF0b3IsIHRoZW4gd2hlbiB0aGUgZGVsZWdhdGlvblxuICAgKiBwZXJpb2QgaXMgb3ZlciwgMS4yMzQ1JSBvZiB0aGUgcmV3YXJkIGdvZXMgdG8gdGhlIHZhbGlkYXRvciBhbmQgdGhlIHJlc3QgZ29lcyB0byB0aGUgZGVsZWdhdG9yLlxuICAgKi9cbiAgY29uc3RydWN0b3IoXG4gICAgbmV0d29ya0lEOiBudW1iZXIgPSBEZWZhdWx0TmV0d29ya0lELFxuICAgIGJsb2NrY2hhaW5JRDogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDMyLCAxNiksXG4gICAgb3V0czogVHJhbnNmZXJhYmxlT3V0cHV0W10gPSB1bmRlZmluZWQsXG4gICAgaW5zOiBUcmFuc2ZlcmFibGVJbnB1dFtdID0gdW5kZWZpbmVkLFxuICAgIG1lbW86IEJ1ZmZlciA9IHVuZGVmaW5lZCxcbiAgICBub2RlSUQ6IEJ1ZmZlciA9IHVuZGVmaW5lZCxcbiAgICBzdGFydFRpbWU6IEJOID0gdW5kZWZpbmVkLFxuICAgIGVuZFRpbWU6IEJOID0gdW5kZWZpbmVkLFxuICAgIHN0YWtlQW1vdW50OiBCTiA9IHVuZGVmaW5lZCxcbiAgICBzdGFrZU91dHM6IFRyYW5zZmVyYWJsZU91dHB1dFtdID0gdW5kZWZpbmVkLFxuICAgIHJld2FyZE93bmVyczogUGFyc2VhYmxlT3V0cHV0ID0gdW5kZWZpbmVkLFxuICAgIGRlbGVnYXRpb25GZWU6IG51bWJlciA9IHVuZGVmaW5lZFxuICApIHtcbiAgICBzdXBlcihcbiAgICAgIG5ldHdvcmtJRCxcbiAgICAgIGJsb2NrY2hhaW5JRCxcbiAgICAgIG91dHMsXG4gICAgICBpbnMsXG4gICAgICBtZW1vLFxuICAgICAgbm9kZUlELFxuICAgICAgc3RhcnRUaW1lLFxuICAgICAgZW5kVGltZSxcbiAgICAgIHN0YWtlQW1vdW50LFxuICAgICAgc3Rha2VPdXRzLFxuICAgICAgcmV3YXJkT3duZXJzXG4gICAgKVxuICAgIGlmICh0eXBlb2YgZGVsZWdhdGlvbkZlZSA9PT0gXCJudW1iZXJcIikge1xuICAgICAgaWYgKGRlbGVnYXRpb25GZWUgPj0gMCAmJiBkZWxlZ2F0aW9uRmVlIDw9IDEwMCkge1xuICAgICAgICB0aGlzLmRlbGVnYXRpb25GZWUgPSBwYXJzZUZsb2F0KGRlbGVnYXRpb25GZWUudG9GaXhlZCg0KSlcbiAgICAgIH0gZWxzZSB7XG4gICAgICAgIHRocm93IG5ldyBEZWxlZ2F0aW9uRmVlRXJyb3IoXG4gICAgICAgICAgXCJBZGRWYWxpZGF0b3JUeC5jb25zdHJ1Y3RvciAtLSBkZWxlZ2F0aW9uRmVlIG11c3QgYmUgaW4gdGhlIHJhbmdlIG9mIDAgYW5kIDEwMCwgaW5jbHVzaXZlbHkuXCJcbiAgICAgICAgKVxuICAgICAgfVxuICAgIH1cbiAgfVxufVxuIl19