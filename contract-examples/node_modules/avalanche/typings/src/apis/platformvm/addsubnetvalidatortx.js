"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.AddSubnetValidatorTx = void 0;
/**
 * @packageDocumentation
 * @module API-PlatformVM-AddSubnetValidatorTx
 */
const buffer_1 = require("buffer/");
const bintools_1 = __importDefault(require("../../utils/bintools"));
const constants_1 = require("./constants");
const credentials_1 = require("../../common/credentials");
const basetx_1 = require("./basetx");
const constants_2 = require("../../utils/constants");
const serialization_1 = require("../../utils/serialization");
const _1 = require(".");
const utils_1 = require("../../utils");
/**
 * @ignore
 */
const bintools = bintools_1.default.getInstance();
const serialization = serialization_1.Serialization.getInstance();
/**
 * Class representing an unsigned AddSubnetValidatorTx transaction.
 */
class AddSubnetValidatorTx extends basetx_1.BaseTx {
    /**
     * Class representing an unsigned AddSubnetValidator transaction.
     *
     * @param networkID Optional networkID, [[DefaultNetworkID]]
     * @param blockchainID Optional blockchainID, default Buffer.alloc(32, 16)
     * @param outs Optional array of the [[TransferableOutput]]s
     * @param ins Optional array of the [[TransferableInput]]s
     * @param memo Optional {@link https://github.com/feross/buffer|Buffer} for the memo field
     * @param nodeID Optional. The node ID of the validator being added.
     * @param startTime Optional. The Unix time when the validator starts validating the Primary Network.
     * @param endTime Optional. The Unix time when the validator stops validating the Primary Network (and staked AVAX is returned).
     * @param weight Optional. Weight of this validator used when sampling
     * @param subnetID Optional. ID of the subnet this validator is validating
     */
    constructor(networkID = constants_2.DefaultNetworkID, blockchainID = buffer_1.Buffer.alloc(32, 16), outs = undefined, ins = undefined, memo = undefined, nodeID = undefined, startTime = undefined, endTime = undefined, weight = undefined, subnetID = undefined) {
        super(networkID, blockchainID, outs, ins, memo);
        this._typeName = "AddSubnetValidatorTx";
        this._typeID = constants_1.PlatformVMConstants.ADDSUBNETVALIDATORTX;
        this.nodeID = buffer_1.Buffer.alloc(20);
        this.startTime = buffer_1.Buffer.alloc(8);
        this.endTime = buffer_1.Buffer.alloc(8);
        this.weight = buffer_1.Buffer.alloc(8);
        this.subnetID = buffer_1.Buffer.alloc(32);
        this.sigCount = buffer_1.Buffer.alloc(4);
        this.sigIdxs = []; // idxs of subnet auth signers
        if (typeof subnetID != "undefined") {
            if (typeof subnetID === "string") {
                this.subnetID = bintools.cb58Decode(subnetID);
            }
            else {
                this.subnetID = subnetID;
            }
        }
        if (typeof nodeID != "undefined") {
            this.nodeID = nodeID;
        }
        if (typeof startTime != "undefined") {
            this.startTime = bintools.fromBNToBuffer(startTime, 8);
        }
        if (typeof endTime != "undefined") {
            this.endTime = bintools.fromBNToBuffer(endTime, 8);
        }
        if (typeof weight != "undefined") {
            this.weight = bintools.fromBNToBuffer(weight, 8);
        }
        const subnetAuth = new _1.SubnetAuth();
        this.subnetAuth = subnetAuth;
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { subnetID: serialization.encoder(this.subnetID, encoding, "Buffer", "cb58") });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.subnetID = serialization.decoder(fields["subnetID"], encoding, "cb58", "Buffer", 32);
        // this.exportOuts = fields["exportOuts"].map((e: object) => {
        //   let eo: TransferableOutput = new TransferableOutput()
        //   eo.deserialize(e, encoding)
        //   return eo
        // })
    }
    /**
     * Returns the id of the [[AddSubnetValidatorTx]]
     */
    getTxType() {
        return constants_1.PlatformVMConstants.ADDSUBNETVALIDATORTX;
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
        return (0, utils_1.bufferToNodeIDString)(this.nodeID);
    }
    /**
     * Returns a {@link https://github.com/indutny/bn.js/|BN} for the startTime.
     */
    getStartTime() {
        return bintools.fromBufferToBN(this.startTime);
    }
    /**
     * Returns a {@link https://github.com/indutny/bn.js/|BN} for the endTime.
     */
    getEndTime() {
        return bintools.fromBufferToBN(this.endTime);
    }
    /**
     * Returns a {@link https://github.com/indutny/bn.js/|BN} for the weight
     */
    getWeight() {
        return bintools.fromBufferToBN(this.weight);
    }
    /**
     * Returns the subnetID as a string
     */
    getSubnetID() {
        return bintools.cb58Encode(this.subnetID);
    }
    /**
     * Returns the subnetAuth
     */
    getSubnetAuth() {
        return this.subnetAuth;
    }
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} containing an [[AddSubnetValidatorTx]], parses it, populates the class, and returns the length of the [[CreateChainTx]] in bytes.
     *
     * @param bytes A {@link https://github.com/feross/buffer|Buffer} containing a raw [[AddSubnetValidatorTx]]
     *
     * @returns The length of the raw [[AddSubnetValidatorTx]]
     *
     * @remarks assume not-checksummed
     */
    fromBuffer(bytes, offset = 0) {
        offset = super.fromBuffer(bytes, offset);
        this.nodeID = bintools.copyFrom(bytes, offset, offset + 20);
        offset += 20;
        this.startTime = bintools.copyFrom(bytes, offset, offset + 8);
        offset += 8;
        this.endTime = bintools.copyFrom(bytes, offset, offset + 8);
        offset += 8;
        this.weight = bintools.copyFrom(bytes, offset, offset + 8);
        offset += 8;
        this.subnetID = bintools.copyFrom(bytes, offset, offset + 32);
        offset += 32;
        const sa = new _1.SubnetAuth();
        offset += sa.fromBuffer(bintools.copyFrom(bytes, offset));
        this.subnetAuth = sa;
        return offset;
    }
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[CreateChainTx]].
     */
    toBuffer() {
        const superbuff = super.toBuffer();
        const bsize = superbuff.length +
            this.nodeID.length +
            this.startTime.length +
            this.endTime.length +
            this.weight.length +
            this.subnetID.length +
            this.subnetAuth.toBuffer().length;
        const barr = [
            superbuff,
            this.nodeID,
            this.startTime,
            this.endTime,
            this.weight,
            this.subnetID,
            this.subnetAuth.toBuffer()
        ];
        return buffer_1.Buffer.concat(barr, bsize);
    }
    clone() {
        const newAddSubnetValidatorTx = new AddSubnetValidatorTx();
        newAddSubnetValidatorTx.fromBuffer(this.toBuffer());
        return newAddSubnetValidatorTx;
    }
    create(...args) {
        return new AddSubnetValidatorTx(...args);
    }
    /**
     * Creates and adds a [[SigIdx]] to the [[AddSubnetValidatorTx]].
     *
     * @param addressIdx The index of the address to reference in the signatures
     * @param address The address of the source of the signature
     */
    addSignatureIdx(addressIdx, address) {
        const addressIndex = buffer_1.Buffer.alloc(4);
        addressIndex.writeUIntBE(addressIdx, 0, 4);
        this.subnetAuth.addAddressIndex(addressIndex);
        const sigidx = new credentials_1.SigIdx();
        const b = buffer_1.Buffer.alloc(4);
        b.writeUInt32BE(addressIdx, 0);
        sigidx.fromBuffer(b);
        sigidx.setSource(address);
        this.sigIdxs.push(sigidx);
        this.sigCount.writeUInt32BE(this.sigIdxs.length, 0);
    }
    /**
     * Returns the array of [[SigIdx]] for this [[Input]]
     */
    getSigIdxs() {
        return this.sigIdxs;
    }
    getCredentialID() {
        return constants_1.PlatformVMConstants.SECPCREDENTIAL;
    }
    /**
     * Takes the bytes of an [[UnsignedTx]] and returns an array of [[Credential]]s
     *
     * @param msg A Buffer for the [[UnsignedTx]]
     * @param kc An [[KeyChain]] used in signing
     *
     * @returns An array of [[Credential]]s
     */
    sign(msg, kc) {
        const creds = super.sign(msg, kc);
        const sigidxs = this.getSigIdxs();
        const cred = (0, _1.SelectCredentialClass)(this.getCredentialID());
        for (let i = 0; i < sigidxs.length; i++) {
            const keypair = kc.getKey(sigidxs[`${i}`].getSource());
            const signval = keypair.sign(msg);
            const sig = new credentials_1.Signature();
            sig.fromBuffer(signval);
            cred.addSignature(sig);
        }
        creds.push(cred);
        return creds;
    }
}
exports.AddSubnetValidatorTx = AddSubnetValidatorTx;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYWRkc3VibmV0dmFsaWRhdG9ydHguanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi9zcmMvYXBpcy9wbGF0Zm9ybXZtL2FkZHN1Ym5ldHZhbGlkYXRvcnR4LnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7OztBQUFBOzs7R0FHRztBQUNILG9DQUFnQztBQUNoQyxvRUFBMkM7QUFDM0MsMkNBQWlEO0FBR2pELDBEQUF3RTtBQUN4RSxxQ0FBaUM7QUFDakMscURBQXdEO0FBQ3hELDZEQUE2RTtBQUM3RSx3QkFBcUQ7QUFHckQsdUNBQWtEO0FBRWxEOztHQUVHO0FBQ0gsTUFBTSxRQUFRLEdBQWEsa0JBQVEsQ0FBQyxXQUFXLEVBQUUsQ0FBQTtBQUNqRCxNQUFNLGFBQWEsR0FBa0IsNkJBQWEsQ0FBQyxXQUFXLEVBQUUsQ0FBQTtBQUVoRTs7R0FFRztBQUNILE1BQWEsb0JBQXFCLFNBQVEsZUFBTTtJQTBOOUM7Ozs7Ozs7Ozs7Ozs7T0FhRztJQUNILFlBQ0UsWUFBb0IsNEJBQWdCLEVBQ3BDLGVBQXVCLGVBQU0sQ0FBQyxLQUFLLENBQUMsRUFBRSxFQUFFLEVBQUUsQ0FBQyxFQUMzQyxPQUE2QixTQUFTLEVBQ3RDLE1BQTJCLFNBQVMsRUFDcEMsT0FBZSxTQUFTLEVBQ3hCLFNBQWlCLFNBQVMsRUFDMUIsWUFBZ0IsU0FBUyxFQUN6QixVQUFjLFNBQVMsRUFDdkIsU0FBYSxTQUFTLEVBQ3RCLFdBQTRCLFNBQVM7UUFFckMsS0FBSyxDQUFDLFNBQVMsRUFBRSxZQUFZLEVBQUUsSUFBSSxFQUFFLEdBQUcsRUFBRSxJQUFJLENBQUMsQ0FBQTtRQW5QdkMsY0FBUyxHQUFHLHNCQUFzQixDQUFBO1FBQ2xDLFlBQU8sR0FBRywrQkFBbUIsQ0FBQyxvQkFBb0IsQ0FBQTtRQTBCbEQsV0FBTSxHQUFXLGVBQU0sQ0FBQyxLQUFLLENBQUMsRUFBRSxDQUFDLENBQUE7UUFDakMsY0FBUyxHQUFXLGVBQU0sQ0FBQyxLQUFLLENBQUMsQ0FBQyxD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QUMsR0FBRyxFQUFFLEVBQUUsQ0FBQyxDQUFBO1FBQy9DLE1BQU0sT0FBTyxHQUFhLElBQUksQ0FBQyxVQUFVLEVBQUUsQ0FBQTtRQUMzQyxNQUFNLElBQUksR0FBZSxJQUFBLHdCQUFxQixFQUFDLElBQUksQ0FBQyxlQUFlLEVBQUUsQ0FBQyxDQUFBO1FBQ3RFLEtBQUssSUFBSSxDQUFDLEdBQVcsQ0FBQyxFQUFFLENBQUMsR0FBRyxPQUFPLENBQUMsTUFBTSxFQUFFLENBQUMsRUFBRSxFQUFFO1lBQy9DLE1BQU0sT0FBTyxHQUFZLEVBQUUsQ0FBQyxNQUFNLENBQUMsT0FBTyxDQUFDLEdBQUcsQ0FBQyxFQUFFLENBQUMsQ0FBQyxTQUFTLEVBQUUsQ0FBQyxDQUFBO1lBQy9ELE1BQU0sT0FBTyxHQUFXLE9BQU8sQ0FBQyxJQUFJLENBQUMsR0FBRyxDQUFDLENBQUE7WUFDekMsTUFBTSxHQUFHLEdBQWMsSUFBSSx1QkFBUyxFQUFFLENBQUE7WUFDdEMsR0FBRyxDQUFDLFVBQVUsQ0FBQyxPQUFPLENBQUMsQ0FBQTtZQUN2QixJQUFJLENBQUMsWUFBWSxDQUFDLEdBQUcsQ0FBQyxDQUFBO1NBQ3ZCO1FBQ0QsS0FBSyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtRQUNoQixPQUFPLEtBQUssQ0FBQTtJQUNkLENBQUM7Q0FvREY7QUE1UUQsb0RBNFFDIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgQVBJLVBsYXRmb3JtVk0tQWRkU3VibmV0VmFsaWRhdG9yVHhcbiAqL1xuaW1wb3J0IHsgQnVmZmVyIH0gZnJvbSBcImJ1ZmZlci9cIlxuaW1wb3J0IEJpblRvb2xzIGZyb20gXCIuLi8uLi91dGlscy9iaW50b29sc1wiXG5pbXBvcnQgeyBQbGF0Zm9ybVZNQ29uc3RhbnRzIH0gZnJvbSBcIi4vY29uc3RhbnRzXCJcbmltcG9ydCB7IFRyYW5zZmVyYWJsZU91dHB1dCB9IGZyb20gXCIuL291dHB1dHNcIlxuaW1wb3J0IHsgVHJhbnNmZXJhYmxlSW5wdXQgfSBmcm9tIFwiLi9pbnB1dHNcIlxuaW1wb3J0IHsgQ3JlZGVudGlhbCwgU2lnSWR4LCBTaWduYXR1cmUgfSBmcm9tIFwiLi4vLi4vY29tbW9uL2NyZWRlbnRpYWxzXCJcbmltcG9ydCB7IEJhc2VUeCB9IGZyb20gXCIuL2Jhc2V0eFwiXG5pbXBvcnQgeyBEZWZhdWx0TmV0d29ya0lEIH0gZnJvbSBcIi4uLy4uL3V0aWxzL2NvbnN0YW50c1wiXG5pbXBvcnQgeyBTZXJpYWxpemF0aW9uLCBTZXJpYWxpemVkRW5jb2RpbmcgfSBmcm9tIFwiLi4vLi4vdXRpbHMvc2VyaWFsaXphdGlvblwiXG5pbXBvcnQgeyBTZWxlY3RDcmVkZW50aWFsQ2xhc3MsIFN1Ym5ldEF1dGggfSBmcm9tIFwiLlwiXG5pbXBvcnQgeyBLZXlDaGFpbiwgS2V5UGFpciB9IGZyb20gXCIuL2tleWNoYWluXCJcbmltcG9ydCBCTiBmcm9tIFwiYm4uanNcIlxuaW1wb3J0IHsgYnVmZmVyVG9Ob2RlSURTdHJpbmcgfSBmcm9tIFwiLi4vLi4vdXRpbHNcIlxuXG4vKipcbiAqIEBpZ25vcmVcbiAqL1xuY29uc3QgYmludG9vbHM6IEJpblRvb2xzID0gQmluVG9vbHMuZ2V0SW5zdGFuY2UoKVxuY29uc3Qgc2VyaWFsaXphdGlvbjogU2VyaWFsaXphdGlvbiA9IFNlcmlhbGl6YXRpb24uZ2V0SW5zdGFuY2UoKVxuXG4vKipcbiAqIENsYXNzIHJlcHJlc2VudGluZyBhbiB1bnNpZ25lZCBBZGRTdWJuZXRWYWxpZGF0b3JUeCB0cmFuc2FjdGlvbi5cbiAqL1xuZXhwb3J0IGNsYXNzIEFkZFN1Ym5ldFZhbGlkYXRvclR4IGV4dGVuZHMgQmFzZVR4IHtcbiAgcHJvdGVjdGVkIF90eXBlTmFtZSA9IFwiQWRkU3VibmV0VmFsaWRhdG9yVHhcIlxuICBwcm90ZWN0ZWQgX3R5cGVJRCA9IFBsYXRmb3JtVk1Db25zdGFudHMuQUREU1VCTkVUVkFMSURBVE9SVFhcblxuICBzZXJpYWxpemUoZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpOiBvYmplY3Qge1xuICAgIGxldCBmaWVsZHM6IG9iamVjdCA9IHN1cGVyLnNlcmlhbGl6ZShlbmNvZGluZylcbiAgICByZXR1cm4ge1xuICAgICAgLi4uZmllbGRzLFxuICAgICAgc3VibmV0SUQ6IHNlcmlhbGl6YXRpb24uZW5jb2Rlcih0aGlzLnN1Ym5ldElELCBlbmNvZGluZywgXCJCdWZmZXJcIiwgXCJjYjU4XCIpXG4gICAgICAvLyBleHBvcnRPdXRzOiB0aGlzLmV4cG9ydE91dHMubWFwKChlKSA9PiBlLnNlcmlhbGl6ZShlbmNvZGluZykpXG4gICAgfVxuICB9XG4gIGRlc2VyaWFsaXplKGZpZWxkczogb2JqZWN0LCBlbmNvZGluZzogU2VyaWFsaXplZEVuY29kaW5nID0gXCJoZXhcIikge1xuICAgIHN1cGVyLmRlc2VyaWFsaXplKGZpZWxkcywgZW5jb2RpbmcpXG4gICAgdGhpcy5zdWJuZXRJRCA9IHNlcmlhbGl6YXRpb24uZGVjb2RlcihcbiAgICAgIGZpZWxkc1tcInN1Ym5ldElEXCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBcImNiNThcIixcbiAgICAgIFwiQnVmZmVyXCIsXG4gICAgICAzMlxuICAgIClcbiAgICAvLyB0aGlzLmV4cG9ydE91dHMgPSBmaWVsZHNbXCJleHBvcnRPdXRzXCJdLm1hcCgoZTogb2JqZWN0KSA9PiB7XG4gICAgLy8gICBsZXQgZW86IFRyYW5zZmVyYWJsZU91dHB1dCA9IG5ldyBUcmFuc2ZlcmFibGVPdXRwdXQoKVxuICAgIC8vICAgZW8uZGVzZXJpYWxpemUoZSwgZW5jb2RpbmcpXG4gICAgLy8gICByZXR1cm4gZW9cbiAgICAvLyB9KVxuICB9XG5cbiAgcHJvdGVjdGVkIG5vZGVJRDogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDIwKVxuICBwcm90ZWN0ZWQgc3RhcnRUaW1lOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoOClcbiAgcHJvdGVjdGVkIGVuZFRpbWU6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYyg4KVxuICBwcm90ZWN0ZWQgd2VpZ2h0OiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoOClcbiAgcHJvdGVjdGVkIHN1Ym5ldElEOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoMzIpXG4gIHByb3RlY3RlZCBzdWJuZXRBdXRoOiBTdWJuZXRBdXRoXG4gIHByb3RlY3RlZCBzaWdDb3VudDogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDQpXG4gIHByb3RlY3RlZCBzaWdJZHhzOiBTaWdJZHhbXSA9IFtdIC8vIGlkeHMgb2Ygc3VibmV0IGF1dGggc2lnbmVyc1xuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBpZCBvZiB0aGUgW1tBZGRTdWJuZXRWYWxpZGF0b3JUeF1dXG4gICAqL1xuICBnZXRUeFR5cGUoKTogbnVtYmVyIHtcbiAgICByZXR1cm4gUGxhdGZvcm1WTUNvbnN0YW50cy5BRERTVUJORVRWQUxJREFUT1JUWFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSBmb3IgdGhlIHN0YWtlIGFtb3VudC5cbiAgICovXG4gIGdldE5vZGVJRCgpOiBCdWZmZXIge1xuICAgIHJldHVybiB0aGlzLm5vZGVJRFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSBzdHJpbmcgZm9yIHRoZSBub2RlSUQgYW1vdW50LlxuICAgKi9cbiAgZ2V0Tm9kZUlEU3RyaW5nKCk6IHN0cmluZyB7XG4gICAgcmV0dXJuIGJ1ZmZlclRvTm9kZUlEU3RyaW5nKHRoaXMubm9kZUlEKVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2luZHV0bnkvYm4uanMvfEJOfSBmb3IgdGhlIHN0YXJ0VGltZS5cbiAgICovXG4gIGdldFN0YXJ0VGltZSgpOiBCTiB7XG4gICAgcmV0dXJuIGJpbnRvb2xzLmZyb21CdWZmZXJUb0JOKHRoaXMuc3RhcnRUaW1lKVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2luZHV0bnkvYm4uanMvfEJOfSBmb3IgdGhlIGVuZFRpbWUuXG4gICAqL1xuICBnZXRFbmRUaW1lKCk6IEJOIHtcbiAgICByZXR1cm4gYmludG9vbHMuZnJvbUJ1ZmZlclRvQk4odGhpcy5lbmRUaW1lKVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2luZHV0bnkvYm4uanMvfEJOfSBmb3IgdGhlIHdlaWdodFxuICAgKi9cbiAgZ2V0V2VpZ2h0KCk6IEJOIHtcbiAgICByZXR1cm4gYmludG9vbHMuZnJvbUJ1ZmZlclRvQk4odGhpcy53ZWlnaHQpXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgc3VibmV0SUQgYXMgYSBzdHJpbmdcbiAgICovXG4gIGdldFN1Ym5ldElEKCk6IHN0cmluZyB7XG4gICAgcmV0dXJuIGJpbnRvb2xzLmNiNThFbmNvZGUodGhpcy5zdWJuZXRJRClcbiAgfVxuICAvKipcbiAgICogUmV0dXJucyB0aGUgc3VibmV0QXV0aFxuICAgKi9cbiAgZ2V0U3VibmV0QXV0aCgpOiBTdWJuZXRBdXRoIHtcbiAgICByZXR1cm4gdGhpcy5zdWJuZXRBdXRoXG4gIH1cblxuICAvKipcbiAgICogVGFrZXMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSBjb250YWluaW5nIGFuIFtbQWRkU3VibmV0VmFsaWRhdG9yVHhdXSwgcGFyc2VzIGl0LCBwb3B1bGF0ZXMgdGhlIGNsYXNzLCBhbmQgcmV0dXJucyB0aGUgbGVuZ3RoIG9mIHRoZSBbW0NyZWF0ZUNoYWluVHhdXSBpbiBieXRlcy5cbiAgICpcbiAgICogQHBhcmFtIGJ5dGVzIEEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gY29udGFpbmluZyBhIHJhdyBbW0FkZFN1Ym5ldFZhbGlkYXRvclR4XV1cbiAgICpcbiAgICogQHJldHVybnMgVGhlIGxlbmd0aCBvZiB0aGUgcmF3IFtbQWRkU3VibmV0VmFsaWRhdG9yVHhdXVxuICAgKlxuICAgKiBAcmVtYXJrcyBhc3N1bWUgbm90LWNoZWNrc3VtbWVkXG4gICAqL1xuICBmcm9tQnVmZmVyKGJ5dGVzOiBCdWZmZXIsIG9mZnNldDogbnVtYmVyID0gMCk6IG51bWJlciB7XG4gICAgb2Zmc2V0ID0gc3VwZXIuZnJvbUJ1ZmZlcihieXRlcywgb2Zmc2V0KVxuXG4gICAgdGhpcy5ub2RlSUQgPSBiaW50b29scy5jb3B5RnJvbShieXRlcywgb2Zmc2V0LCBvZmZzZXQgKyAyMClcbiAgICBvZmZzZXQgKz0gMjBcblxuICAgIHRoaXMuc3RhcnRUaW1lID0gYmludG9vbHMuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgOClcbiAgICBvZmZzZXQgKz0gOFxuXG4gICAgdGhpcy5lbmRUaW1lID0gYmludG9vbHMuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgOClcbiAgICBvZmZzZXQgKz0gOFxuXG4gICAgdGhpcy53ZWlnaHQgPSBiaW50b29scy5jb3B5RnJvbShieXRlcywgb2Zmc2V0LCBvZmZzZXQgKyA4KVxuICAgIG9mZnNldCArPSA4XG5cbiAgICB0aGlzLnN1Ym5ldElEID0gYmludG9vbHMuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgMzIpXG4gICAgb2Zmc2V0ICs9IDMyXG5cbiAgICBjb25zdCBzYTogU3VibmV0QXV0aCA9IG5ldyBTdWJuZXRBdXRoKClcbiAgICBvZmZzZXQgKz0gc2EuZnJvbUJ1ZmZlcihiaW50b29scy5jb3B5RnJvbShieXRlcywgb2Zmc2V0KSlcbiAgICB0aGlzLnN1Ym5ldEF1dGggPSBzYVxuXG4gICAgcmV0dXJuIG9mZnNldFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSByZXByZXNlbnRhdGlvbiBvZiB0aGUgW1tDcmVhdGVDaGFpblR4XV0uXG4gICAqL1xuICB0b0J1ZmZlcigpOiBCdWZmZXIge1xuICAgIGNvbnN0IHN1cGVyYnVmZjogQnVmZmVyID0gc3VwZXIudG9CdWZmZXIoKVxuXG4gICAgY29uc3QgYnNpemU6IG51bWJlciA9XG4gICAgICBzdXBlcmJ1ZmYubGVuZ3RoICtcbiAgICAgIHRoaXMubm9kZUlELmxlbmd0aCArXG4gICAgICB0aGlzLnN0YXJ0VGltZS5sZW5ndGggK1xuICAgICAgdGhpcy5lbmRUaW1lLmxlbmd0aCArXG4gICAgICB0aGlzLndlaWdodC5sZW5ndGggK1xuICAgICAgdGhpcy5zdWJuZXRJRC5sZW5ndGggK1xuICAgICAgdGhpcy5zdWJuZXRBdXRoLnRvQnVmZmVyKCkubGVuZ3RoXG5cbiAgICBjb25zdCBiYXJyOiBCdWZmZXJbXSA9IFtcbiAgICAgIHN1cGVyYnVmZixcbiAgICAgIHRoaXMubm9kZUlELFxuICAgICAgdGhpcy5zdGFydFRpbWUsXG4gICAgICB0aGlzLmVuZFRpbWUsXG4gICAgICB0aGlzLndlaWdodCxcbiAgICAgIHRoaXMuc3VibmV0SUQsXG4gICAgICB0aGlzLnN1Ym5ldEF1dGgudG9CdWZmZXIoKVxuICAgIF1cbiAgICByZXR1cm4gQnVmZmVyLmNvbmNhdChiYXJyLCBic2l6ZSlcbiAgfVxuXG4gIGNsb25lKCk6IHRoaXMge1xuICAgIGNvbnN0IG5ld0FkZFN1Ym5ldFZhbGlkYXRvclR4OiBBZGRTdWJuZXRWYWxpZGF0b3JUeCA9XG4gICAgICBuZXcgQWRkU3VibmV0VmFsaWRhdG9yVHgoKVxuICAgIG5ld0FkZFN1Ym5ldFZhbGlkYXRvclR4LmZyb21CdWZmZXIodGhpcy50b0J1ZmZlcigpKVxuICAgIHJldHVybiBuZXdBZGRTdWJuZXRWYWxpZGF0b3JUeCBhcyB0aGlzXG4gIH1cblxuICBjcmVhdGUoLi4uYXJnczogYW55W10pOiB0aGlzIHtcbiAgICByZXR1cm4gbmV3IEFkZFN1Ym5ldFZhbGlkYXRvclR4KC4uLmFyZ3MpIGFzIHRoaXNcbiAgfVxuXG4gIC8qKlxuICAgKiBDcmVhdGVzIGFuZCBhZGRzIGEgW1tTaWdJZHhdXSB0byB0aGUgW1tBZGRTdWJuZXRWYWxpZGF0b3JUeF1dLlxuICAgKlxuICAgKiBAcGFyYW0gYWRkcmVzc0lkeCBUaGUgaW5kZXggb2YgdGhlIGFkZHJlc3MgdG8gcmVmZXJlbmNlIGluIHRoZSBzaWduYXR1cmVzXG4gICAqIEBwYXJhbSBhZGRyZXNzIFRoZSBhZGRyZXNzIG9mIHRoZSBzb3VyY2Ugb2YgdGhlIHNpZ25hdHVyZVxuICAgKi9cbiAgYWRkU2lnbmF0dXJlSWR4KGFkZHJlc3NJZHg6IG51bWJlciwgYWRkcmVzczogQnVmZmVyKTogdm9pZCB7XG4gICAgY29uc3QgYWRkcmVzc0luZGV4OiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoNClcbiAgICBhZGRyZXNzSW5kZXgud3JpdGVVSW50QkUoYWRkcmVzc0lkeCwgMCwgNClcbiAgICB0aGlzLnN1Ym5ldEF1dGguYWRkQWRkcmVzc0luZGV4KGFkZHJlc3NJbmRleClcblxuICAgIGNvbnN0IHNpZ2lkeDogU2lnSWR4ID0gbmV3IFNpZ0lkeCgpXG4gICAgY29uc3QgYjogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDQpXG4gICAgYi53cml0ZVVJbnQzMkJFKGFkZHJlc3NJZHgsIDApXG4gICAgc2lnaWR4LmZyb21CdWZmZXIoYilcbiAgICBzaWdpZHguc2V0U291cmNlKGFkZHJlc3MpXG4gICAgdGhpcy5zaWdJZHhzLnB1c2goc2lnaWR4KVxuICAgIHRoaXMuc2lnQ291bnQud3JpdGVVSW50MzJCRSh0aGlzLnNpZ0lkeHMubGVuZ3RoLCAwKVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIGFycmF5IG9mIFtbU2lnSWR4XV0gZm9yIHRoaXMgW1tJbnB1dF1dXG4gICAqL1xuICBnZXRTaWdJZHhzKCk6IFNpZ0lkeFtdIHtcbiAgICByZXR1cm4gdGhpcy5zaWdJZHhzXG4gIH1cblxuICBnZXRDcmVkZW50aWFsSUQoKTogbnVtYmVyIHtcbiAgICByZXR1cm4gUGxhdGZvcm1WTUNvbnN0YW50cy5TRUNQQ1JFREVOVElBTFxuICB9XG5cbiAgLyoqXG4gICAqIFRha2VzIHRoZSBieXRlcyBvZiBhbiBbW1Vuc2lnbmVkVHhdXSBhbmQgcmV0dXJucyBhbiBhcnJheSBvZiBbW0NyZWRlbnRpYWxdXXNcbiAgICpcbiAgICogQHBhcmFtIG1zZyBBIEJ1ZmZlciBmb3IgdGhlIFtbVW5zaWduZWRUeF1dXG4gICAqIEBwYXJhbSBrYyBBbiBbW0tleUNoYWluXV0gdXNlZCBpbiBzaWduaW5nXG4gICAqXG4gICAqIEByZXR1cm5zIEFuIGFycmF5IG9mIFtbQ3JlZGVudGlhbF1dc1xuICAgKi9cbiAgc2lnbihtc2c6IEJ1ZmZlciwga2M6IEtleUNoYWluKTogQ3JlZGVudGlhbFtdIHtcbiAgICBjb25zdCBjcmVkczogQ3JlZGVudGlhbFtdID0gc3VwZXIuc2lnbihtc2csIGtjKVxuICAgIGNvbnN0IHNpZ2lkeHM6IFNpZ0lkeFtdID0gdGhpcy5nZXRTaWdJZHhzKClcbiAgICBjb25zdCBjcmVkOiBDcmVkZW50aWFsID0gU2VsZWN0Q3JlZGVudGlhbENsYXNzKHRoaXMuZ2V0Q3JlZGVudGlhbElEKCkpXG4gICAgZm9yIChsZXQgaTogbnVtYmVyID0gMDsgaSA8IHNpZ2lkeHMubGVuZ3RoOyBpKyspIHtcbiAgICAgIGNvbnN0IGtleXBhaXI6IEtleVBhaXIgPSBrYy5nZXRLZXkoc2lnaWR4c1tgJHtpfWBdLmdldFNvdXJjZSgpKVxuICAgICAgY29uc3Qgc2lnbnZhbDogQnVmZmVyID0ga2V5cGFpci5zaWduKG1zZylcbiAgICAgIGNvbnN0IHNpZzogU2lnbmF0dXJlID0gbmV3IFNpZ25hdHVyZSgpXG4gICAgICBzaWcuZnJvbUJ1ZmZlcihzaWdudmFsKVxuICAgICAgY3JlZC5hZGRTaWduYXR1cmUoc2lnKVxuICAgIH1cbiAgICBjcmVkcy5wdXNoKGNyZWQpXG4gICAgcmV0dXJuIGNyZWRzXG4gIH1cblxuICAvKipcbiAgICogQ2xhc3MgcmVwcmVzZW50aW5nIGFuIHVuc2lnbmVkIEFkZFN1Ym5ldFZhbGlkYXRvciB0cmFuc2FjdGlvbi5cbiAgICpcbiAgICogQHBhcmFtIG5ldHdvcmtJRCBPcHRpb25hbCBuZXR3b3JrSUQsIFtbRGVmYXVsdE5ldHdvcmtJRF1dXG4gICAqIEBwYXJhbSBibG9ja2NoYWluSUQgT3B0aW9uYWwgYmxvY2tjaGFpbklELCBkZWZhdWx0IEJ1ZmZlci5hbGxvYygzMiwgMTYpXG4gICAqIEBwYXJhbSBvdXRzIE9wdGlvbmFsIGFycmF5IG9mIHRoZSBbW1RyYW5zZmVyYWJsZU91dHB1dF1dc1xuICAgKiBAcGFyYW0gaW5zIE9wdGlvbmFsIGFycmF5IG9mIHRoZSBbW1RyYW5zZmVyYWJsZUlucHV0XV1zXG4gICAqIEBwYXJhbSBtZW1vIE9wdGlvbmFsIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IGZvciB0aGUgbWVtbyBmaWVsZFxuICAgKiBAcGFyYW0gbm9kZUlEIE9wdGlvbmFsLiBUaGUgbm9kZSBJRCBvZiB0aGUgdmFsaWRhdG9yIGJlaW5nIGFkZGVkLlxuICAgKiBAcGFyYW0gc3RhcnRUaW1lIE9wdGlvbmFsLiBUaGUgVW5peCB0aW1lIHdoZW4gdGhlIHZhbGlkYXRvciBzdGFydHMgdmFsaWRhdGluZyB0aGUgUHJpbWFyeSBOZXR3b3JrLlxuICAgKiBAcGFyYW0gZW5kVGltZSBPcHRpb25hbC4gVGhlIFVuaXggdGltZSB3aGVuIHRoZSB2YWxpZGF0b3Igc3RvcHMgdmFsaWRhdGluZyB0aGUgUHJpbWFyeSBOZXR3b3JrIChhbmQgc3Rha2VkIEFWQVggaXMgcmV0dXJuZWQpLlxuICAgKiBAcGFyYW0gd2VpZ2h0IE9wdGlvbmFsLiBXZWlnaHQgb2YgdGhpcyB2YWxpZGF0b3IgdXNlZCB3aGVuIHNhbXBsaW5nXG4gICAqIEBwYXJhbSBzdWJuZXRJRCBPcHRpb25hbC4gSUQgb2YgdGhlIHN1Ym5ldCB0aGlzIHZhbGlkYXRvciBpcyB2YWxpZGF0aW5nXG4gICAqL1xuICBjb25zdHJ1Y3RvcihcbiAgICBuZXR3b3JrSUQ6IG51bWJlciA9IERlZmF1bHROZXR3b3JrSUQsXG4gICAgYmxvY2tjaGFpbklEOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoMzIsIDE2KSxcbiAgICBvdXRzOiBUcmFuc2ZlcmFibGVPdXRwdXRbXSA9IHVuZGVmaW5lZCxcbiAgICBpbnM6IFRyYW5zZmVyYWJsZUlucHV0W10gPSB1bmRlZmluZWQsXG4gICAgbWVtbzogQnVmZmVyID0gdW5kZWZpbmVkLFxuICAgIG5vZGVJRDogQnVmZmVyID0gdW5kZWZpbmVkLFxuICAgIHN0YXJ0VGltZTogQk4gPSB1bmRlZmluZWQsXG4gICAgZW5kVGltZTogQk4gPSB1bmRlZmluZWQsXG4gICAgd2VpZ2h0OiBCTiA9IHVuZGVmaW5lZCxcbiAgICBzdWJuZXRJRDogc3RyaW5nIHwgQnVmZmVyID0gdW5kZWZpbmVkXG4gICkge1xuICAgIHN1cGVyKG5ldHdvcmtJRCwgYmxvY2tjaGFpbklELCBvdXRzLCBpbnMsIG1lbW8pXG4gICAgaWYgKHR5cGVvZiBzdWJuZXRJRCAhPSBcInVuZGVmaW5lZFwiKSB7XG4gICAgICBpZiAodHlwZW9mIHN1Ym5ldElEID09PSBcInN0cmluZ1wiKSB7XG4gICAgICAgIHRoaXMuc3VibmV0SUQgPSBiaW50b29scy5jYjU4RGVjb2RlKHN1Ym5ldElEKVxuICAgICAgfSBlbHNlIHtcbiAgICAgICAgdGhpcy5zdWJuZXRJRCA9IHN1Ym5ldElEXG4gICAgICB9XG4gICAgfVxuICAgIGlmICh0eXBlb2Ygbm9kZUlEICE9IFwidW5kZWZpbmVkXCIpIHtcbiAgICAgIHRoaXMubm9kZUlEID0gbm9kZUlEXG4gICAgfVxuICAgIGlmICh0eXBlb2Ygc3RhcnRUaW1lICE9IFwidW5kZWZpbmVkXCIpIHtcbiAgICAgIHRoaXMuc3RhcnRUaW1lID0gYmludG9vbHMuZnJvbUJOVG9CdWZmZXIoc3RhcnRUaW1lLCA4KVxuICAgIH1cbiAgICBpZiAodHlwZW9mIGVuZFRpbWUgIT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgdGhpcy5lbmRUaW1lID0gYmludG9vbHMuZnJvbUJOVG9CdWZmZXIoZW5kVGltZSwgOClcbiAgICB9XG4gICAgaWYgKHR5cGVvZiB3ZWlnaHQgIT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgdGhpcy53ZWlnaHQgPSBiaW50b29scy5mcm9tQk5Ub0J1ZmZlcih3ZWlnaHQsIDgpXG4gICAgfVxuXG4gICAgY29uc3Qgc3VibmV0QXV0aDogU3VibmV0QXV0aCA9IG5ldyBTdWJuZXRBdXRoKClcbiAgICB0aGlzLnN1Ym5ldEF1dGggPSBzdWJuZXRBdXRoXG4gIH1cbn1cbiJdfQ==