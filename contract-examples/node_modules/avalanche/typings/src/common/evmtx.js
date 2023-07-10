"use strict";
/**
 * @packageDocumentation
 * @module Common-Transactions
 */
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.EVMStandardTx = exports.EVMStandardUnsignedTx = exports.EVMStandardBaseTx = void 0;
const buffer_1 = require("buffer/");
const bintools_1 = __importDefault(require("../utils/bintools"));
const bn_js_1 = __importDefault(require("bn.js"));
const input_1 = require("./input");
const output_1 = require("./output");
const constants_1 = require("../utils/constants");
const serialization_1 = require("../utils/serialization");
/**
 * @ignore
 */
const bintools = bintools_1.default.getInstance();
const serializer = serialization_1.Serialization.getInstance();
/**
 * Class representing a base for all transactions.
 */
class EVMStandardBaseTx extends serialization_1.Serializable {
    /**
     * Class representing a StandardBaseTx which is the foundation for all transactions.
     *
     * @param networkID Optional networkID, [[DefaultNetworkID]]
     * @param blockchainID Optional blockchainID, default Buffer.alloc(32, 16)
     * @param outs Optional array of the [[TransferableOutput]]s
     * @param ins Optional array of the [[TransferableInput]]s
     */
    constructor(networkID = constants_1.DefaultNetworkID, blockchainID = buffer_1.Buffer.alloc(32, 16)) {
        super();
        this._typeName = "EVMStandardBaseTx";
        this._typeID = undefined;
        this.networkID = buffer_1.Buffer.alloc(4);
        this.blockchainID = buffer_1.Buffer.alloc(32);
        this.networkID.writeUInt32BE(networkID, 0);
        this.blockchainID = blockchainID;
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { networkID: serializer.encoder(this.networkID, encoding, "Buffer", "decimalString"), blockchainID: serializer.encoder(this.blockchainID, encoding, "Buffer", "cb58") });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.networkID = serializer.decoder(fields["networkID"], encoding, "decimalString", "Buffer", 4);
        this.blockchainID = serializer.decoder(fields["blockchainID"], encoding, "cb58", "Buffer", 32);
    }
    /**
     * Returns the NetworkID as a number
     */
    getNetworkID() {
        return this.networkID.readUInt32BE(0);
    }
    /**
     * Returns the Buffer representation of the BlockchainID
     */
    getBlockchainID() {
        return this.blockchainID;
    }
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[StandardBaseTx]].
     */
    toBuffer() {
        let bsize = this.networkID.length + this.blockchainID.length;
        const barr = [this.networkID, this.blockchainID];
        const buff = buffer_1.Buffer.concat(barr, bsize);
        return buff;
    }
    /**
     * Returns a base-58 representation of the [[StandardBaseTx]].
     */
    toString() {
        return bintools.bufferToB58(this.toBuffer());
    }
}
exports.EVMStandardBaseTx = EVMStandardBaseTx;
/**
 * Class representing an unsigned transaction.
 */
class EVMStandardUnsignedTx extends serialization_1.Serializable {
    constructor(transaction = undefined, codecID = 0) {
        super();
        this._typeName = "StandardUnsignedTx";
        this._typeID = undefined;
        this.codecID = 0;
        this.codecID = codecID;
        this.transaction = transaction;
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { codecID: serializer.encoder(this.codecID, encoding, "number", "decimalString", 2), transaction: this.transaction.serialize(encoding) });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.codecID = serializer.decoder(fields["codecID"], encoding, "decimalString", "number");
    }
    /**
     * Returns the CodecID as a number
     */
    getCodecID() {
        return this.codecID;
    }
    /**
     * Returns the {@link https://github.com/feross/buffer|Buffer} representation of the CodecID
     */
    getCodecIDBuffer() {
        let codecBuf = buffer_1.Buffer.alloc(2);
        codecBuf.writeUInt16BE(this.codecID, 0);
        return codecBuf;
    }
    /**
     * Returns the inputTotal as a BN
     */
    getInputTotal(assetID) {
        const ins = [];
        const aIDHex = assetID.toString("hex");
        let total = new bn_js_1.default(0);
        ins.forEach((input) => {
            // only check StandardAmountInputs
            if (input.getInput() instanceof input_1.StandardAmountInput &&
                aIDHex === input.getAssetID().toString("hex")) {
                const i = input.getInput();
                total = total.add(i.getAmount());
            }
        });
        return total;
    }
    /**
     * Returns the outputTotal as a BN
     */
    getOutputTotal(assetID) {
        const outs = [];
        const aIDHex = assetID.toString("hex");
        let total = new bn_js_1.default(0);
        outs.forEach((out) => {
            // only check StandardAmountOutput
            if (out.getOutput() instanceof output_1.StandardAmountOutput &&
                aIDHex === out.getAssetID().toString("hex")) {
                const output = out.getOutput();
                total = total.add(output.getAmount());
            }
        });
        return total;
    }
    /**
     * Returns the number of burned tokens as a BN
     */
    getBurn(assetID) {
        return this.getInputTotal(assetID).sub(this.getOutputTotal(assetID));
    }
    toBuffer() {
        const codecID = this.getCodecIDBuffer();
        const txtype = buffer_1.Buffer.alloc(4);
        txtype.writeUInt32BE(this.transaction.getTxType(), 0);
        const basebuff = this.transaction.toBuffer();
        return buffer_1.Buffer.concat([codecID, txtype, basebuff], codecID.length + txtype.length + basebuff.length);
    }
}
exports.EVMStandardUnsignedTx = EVMStandardUnsignedTx;
/**
 * Class representing a signed transaction.
 */
class EVMStandardTx extends serialization_1.Serializable {
    /**
     * Class representing a signed transaction.
     *
     * @param unsignedTx Optional [[StandardUnsignedTx]]
     * @param signatures Optional array of [[Credential]]s
     */
    constructor(unsignedTx = undefined, credentials = undefined) {
        super();
        this._typeName = "StandardTx";
        this._typeID = undefined;
        this.unsignedTx = undefined;
        this.credentials = [];
        if (typeof unsignedTx !== "undefined") {
            this.unsignedTx = unsignedTx;
            if (typeof credentials !== "undefined") {
                this.credentials = credentials;
            }
        }
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { unsignedTx: this.unsignedTx.serialize(encoding), credentials: this.credentials.map((c) => c.serialize(encoding)) });
    }
    /**
     * Returns the [[StandardUnsignedTx]]
     */
    getUnsignedTx() {
        return this.unsignedTx;
    }
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[StandardTx]].
     */
    toBuffer() {
        const txbuff = this.unsignedTx.toBuffer();
        let bsize = txbuff.length;
        const credlen = buffer_1.Buffer.alloc(4);
        credlen.writeUInt32BE(this.credentials.length, 0);
        const barr = [txbuff, credlen];
        bsize += credlen.length;
        this.credentials.forEach((credential) => {
            const credid = buffer_1.Buffer.alloc(4);
            credid.writeUInt32BE(credential.getCredentialID(), 0);
            barr.push(credid);
            bsize += credid.length;
            const credbuff = credential.toBuffer();
            bsize += credbuff.length;
            barr.push(credbuff);
        });
        const buff = buffer_1.Buffer.concat(barr, bsize);
        return buff;
    }
    /**
     * Takes a base-58 string containing an [[StandardTx]], parses it, populates the class, and returns the length of the Tx in bytes.
     *
     * @param serialized A base-58 string containing a raw [[StandardTx]]
     *
     * @returns The length of the raw [[StandardTx]]
     *
     * @remarks
     * unlike most fromStrings, it expects the string to be serialized in cb58 format
     */
    fromString(serialized) {
        return this.fromBuffer(bintools.cb58Decode(serialized));
    }
    /**
     * Returns a cb58 representation of the [[StandardTx]].
     *
     * @remarks
     * unlike most toStrings, this returns in cb58 serialization format
     */
    toString() {
        return bintools.cb58Encode(this.toBuffer());
    }
}
exports.EVMStandardTx = EVMStandardTx;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZXZtdHguanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi9zcmMvY29tbW9uL2V2bXR4LnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7QUFBQTs7O0dBR0c7Ozs7OztBQUVILG9DQUFnQztBQUNoQyxpRUFBd0M7QUFFeEMsa0RBQXNCO0FBRXRCLG1DQUF3RTtBQUN4RSxxQ0FBMkU7QUFDM0Usa0RBQXFEO0FBQ3JELDBEQUkrQjtBQUUvQjs7R0FFRztBQUNILE1BQU0sUUFBUSxHQUFhLGtCQUFRLENBQUMsV0FBVyxFQUFFLENBQUE7QUFDakQsTUFBTSxVQUFVLEdBQWtCLDZCQUFhLENBQUMsV0FBVyxFQUFFLENBQUE7QUFFN0Q7O0dBRUc7QUFDSCxNQUFzQixpQkFHcEIsU0FBUSw0QkFBWTtJQXNGcEI7Ozs7Ozs7T0FPRztJQUNILFlBQ0UsWUFBb0IsNEJBQWdCLEVBQ3BDLGVBQXVCLGVBQU0sQ0FBQyxLQUFLLENBQUMsRUFBRSxFQUFFLEVBQUUsQ0FBQztRQUUzQyxLQUFLLEVBQUUsQ0FBQTtRQWpHQyxjQUFTLEdBQUcsbUJBQW1CLENBQUE7UUFDL0IsWUFBTyxHQUFHLFNBQVMsQ0FBQTtRQXVDbkIsY0FBUyxHQUFXLGVBQU0sQ0FBQyxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUE7UUFDbkMsaUJBQVksR0FBVyxlQUFNLENBQUMsS0FBSyxDQUFDLEVBQUUsQ0FBQyxDQUFBO1FBeUQvQyxJQUFJLENBQUMsU0FBUyxDQUFDLGFBQWEsQ0FBQyxTQUFTLEVBQUUsQ0FBQyxDQUFDLENBQUE7UUFDMUMsSUFBSSxDQUFDLFlBQVksR0FBRyxZQUFZLENBQUE7SUFDbEMsQ0FBQztJQWpHRCxTQUFTLENBQUMsV0FBK0IsS0FBSztRQUM1QyxJQUFJLE1BQU0sR0FBVyxLQUFLLENBQUMsU0FBUyxDQUFDLFFBQVEsQ0FBQyxDQUFBO1FBQzlDLHVDQUNLLE1BQU0sS0FDVCxTQUFTLEVBQUUsVUFBVSxDQUFDLE9BQU8sQ0FDM0IsSUFBSSxDQUFDLFNBQVMs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QyxDQUFDLENBQUE7UUFDakQsTUFBTSxJQUFJLEdBQWEsQ0FBQyxNQUFNLEVBQUUsT0FBTyxDQUFDLENBQUE7UUFDeEMsS0FBSyxJQUFJLE9BQU8sQ0FBQyxNQUFNLENBQUE7UUFDdkIsSUFBSSxDQUFDLFdBQVcsQ0FBQyxPQUFPLENBQUMsQ0FBQyxVQUFzQixFQUFFLEVBQUU7WUFDbEQsTUFBTSxNQUFNLEdBQVcsZUFBTSxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQTtZQUN0QyxNQUFNLENBQUMsYUFBYSxDQUFDLFVBQVUsQ0FBQyxlQUFlLEVBQUUsRUFBRSxDQUFDLENBQUMsQ0FBQTtZQUNyRCxJQUFJLENBQUMsSUFBSSxDQUFDLE1BQU0sQ0FBQyxDQUFBO1lBQ2pCLEtBQUssSUFBSSxNQUFNLENBQUMsTUFBTSxDQUFBO1lBQ3RCLE1BQU0sUUFBUSxHQUFXLFVBQVUsQ0FBQyxRQUFRLEVBQUUsQ0FBQTtZQUM5QyxLQUFLLElBQUksUUFBUSxDQUFDLE1BQU0sQ0FBQTtZQUN4QixJQUFJLENBQUMsSUFBSSxDQUFDLFFBQVEsQ0FBQyxDQUFBO1FBQ3JCLENBQUMsQ0FBQyxDQUFBO1FBQ0YsTUFBTSxJQUFJLEdBQVcsZUFBTSxDQUFDLE1BQU0sQ0FBQyxJQUFJLEVBQUUsS0FBSyxDQUFDLENBQUE7UUFDL0MsT0FBTyxJQUFJLENBQUE7SUFDYixDQUFDO0lBRUQ7Ozs7Ozs7OztPQVNHO0lBQ0gsVUFBVSxDQUFDLFVBQWtCO1FBQzNCLE9BQU8sSUFBSSxDQUFDLFVBQVUsQ0FBQyxRQUFRLENBQUMsVUFBVSxDQUFDLFVBQVUsQ0FBQyxDQUFDLENBQUE7SUFDekQsQ0FBQztJQUVEOzs7OztPQUtHO0lBQ0gsUUFBUTtRQUNOLE9BQU8sUUFBUSxDQUFDLFVBQVUsQ0FBQyxJQUFJLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQTtJQUM3QyxDQUFDO0NBb0JGO0FBbEdELHNDQWtHQyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogQHBhY2thZ2VEb2N1bWVudGF0aW9uXG4gKiBAbW9kdWxlIENvbW1vbi1UcmFuc2FjdGlvbnNcbiAqL1xuXG5pbXBvcnQgeyBCdWZmZXIgfSBmcm9tIFwiYnVmZmVyL1wiXG5pbXBvcnQgQmluVG9vbHMgZnJvbSBcIi4uL3V0aWxzL2JpbnRvb2xzXCJcbmltcG9ydCB7IENyZWRlbnRpYWwgfSBmcm9tIFwiLi9jcmVkZW50aWFsc1wiXG5pbXBvcnQgQk4gZnJvbSBcImJuLmpzXCJcbmltcG9ydCB7IFN0YW5kYXJkS2V5Q2hhaW4sIFN0YW5kYXJkS2V5UGFpciB9IGZyb20gXCIuL2tleWNoYWluXCJcbmltcG9ydCB7IFN0YW5kYXJkQW1vdW50SW5wdXQsIFN0YW5kYXJkVHJhbnNmZXJhYmxlSW5wdXQgfSBmcm9tIFwiLi9pbnB1dFwiXG5pbXBvcnQgeyBTdGFuZGFyZEFtb3VudE91dHB1dCwgU3RhbmRhcmRUcmFuc2ZlcmFibGVPdXRwdXQgfSBmcm9tIFwiLi9vdXRwdXRcIlxuaW1wb3J0IHsgRGVmYXVsdE5ldHdvcmtJRCB9IGZyb20gXCIuLi91dGlscy9jb25zdGFudHNcIlxuaW1wb3J0IHtcbiAgU2VyaWFsaXphYmxlLFxuICBTZXJpYWxpemF0aW9uLFxuICBTZXJpYWxpemVkRW5jb2Rpbmdcbn0gZnJvbSBcIi4uL3V0aWxzL3NlcmlhbGl6YXRpb25cIlxuXG4vKipcbiAqIEBpZ25vcmVcbiAqL1xuY29uc3QgYmludG9vbHM6IEJpblRvb2xzID0gQmluVG9vbHMuZ2V0SW5zdGFuY2UoKVxuY29uc3Qgc2VyaWFsaXplcjogU2VyaWFsaXphdGlvbiA9IFNlcmlhbGl6YXRpb24uZ2V0SW5zdGFuY2UoKVxuXG4vKipcbiAqIENsYXNzIHJlcHJlc2VudGluZyBhIGJhc2UgZm9yIGFsbCB0cmFuc2FjdGlvbnMuXG4gKi9cbmV4cG9ydCBhYnN0cmFjdCBjbGFzcyBFVk1TdGFuZGFyZEJhc2VUeDxcbiAgS1BDbGFzcyBleHRlbmRzIFN0YW5kYXJkS2V5UGFpcixcbiAgS0NDbGFzcyBleHRlbmRzIFN0YW5kYXJkS2V5Q2hhaW48S1BDbGFzcz5cbj4gZXh0ZW5kcyBTZXJpYWxpemFibGUge1xuICBwcm90ZWN0ZWQgX3R5cGVOYW1lID0gXCJFVk1TdGFuZGFyZEJhc2VUeFwiXG4gIHByb3RlY3RlZCBfdHlwZUlEID0gdW5kZWZpbmVkXG5cbiAgc2VyaWFsaXplKGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKTogb2JqZWN0IHtcbiAgICBsZXQgZmllbGRzOiBvYmplY3QgPSBzdXBlci5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgcmV0dXJuIHtcbiAgICAgIC4uLmZpZWxkcyxcbiAgICAgIG5ldHdvcmtJRDogc2VyaWFsaXplci5lbmNvZGVyKFxuICAgICAgICB0aGlzLm5ldHdvcmtJRCxcbiAgICAgICAgZW5jb2RpbmcsXG4gICAgICAgIFwiQnVmZmVyXCIsXG4gICAgICAgIFwiZGVjaW1hbFN0cmluZ1wiXG4gICAgICApLFxuICAgICAgYmxvY2tjaGFpbklEOiBzZXJpYWxpemVyLmVuY29kZXIoXG4gICAgICAgIHRoaXMuYmxvY2tjaGFpbklELFxuICAgICAgICBlbmNvZGluZyxcbiAgICAgICAgXCJCdWZmZXJcIixcbiAgICAgICAgXCJjYjU4XCJcbiAgICAgIClcbiAgICB9XG4gIH1cblxuICBkZXNlcmlhbGl6ZShmaWVsZHM6IG9iamVjdCwgZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpIHtcbiAgICBzdXBlci5kZXNlcmlhbGl6ZShmaWVsZHMsIGVuY29kaW5nKVxuICAgIHRoaXMubmV0d29ya0lEID0gc2VyaWFsaXplci5kZWNvZGVyKFxuICAgICAgZmllbGRzW1wibmV0d29ya0lEXCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBcImRlY2ltYWxTdHJpbmdcIixcbiAgICAgIFwiQnVmZmVyXCIsXG4gICAgICA0XG4gICAgKVxuICAgIHRoaXMuYmxvY2tjaGFpbklEID0gc2VyaWFsaXplci5kZWNvZGVyKFxuICAgICAgZmllbGRzW1wiYmxvY2tjaGFpbklEXCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBcImNiNThcIixcbiAgICAgIFwiQnVmZmVyXCIsXG4gICAgICAzMlxuICAgIClcbiAgfVxuXG4gIHByb3RlY3RlZCBuZXR3b3JrSUQ6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYyg0KVxuICBwcm90ZWN0ZWQgYmxvY2tjaGFpbklEOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoMzIpXG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIGlkIG9mIHRoZSBbW1N0YW5kYXJkQmFzZVR4XV1cbiAgICovXG4gIGFic3RyYWN0IGdldFR4VHlwZSgpOiBudW1iZXJcblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgTmV0d29ya0lEIGFzIGEgbnVtYmVyXG4gICAqL1xuICBnZXROZXR3b3JrSUQoKTogbnVtYmVyIHtcbiAgICByZXR1cm4gdGhpcy5uZXR3b3JrSUQucmVhZFVJbnQzMkJFKDApXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgQnVmZmVyIHJlcHJlc2VudGF0aW9uIG9mIHRoZSBCbG9ja2NoYWluSURcbiAgICovXG4gIGdldEJsb2NrY2hhaW5JRCgpOiBCdWZmZXIge1xuICAgIHJldHVybiB0aGlzLmJsb2NrY2hhaW5JRFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSByZXByZXNlbnRhdGlvbiBvZiB0aGUgW1tTdGFuZGFyZEJhc2VUeF1dLlxuICAgKi9cbiAgdG9CdWZmZXIoKTogQnVmZmVyIHtcbiAgICBsZXQgYnNpemU6IG51bWJlciA9IHRoaXMubmV0d29ya0lELmxlbmd0aCArIHRoaXMuYmxvY2tjaGFpbklELmxlbmd0aFxuICAgIGNvbnN0IGJhcnI6IEJ1ZmZlcltdID0gW3RoaXMubmV0d29ya0lELCB0aGlzLmJsb2NrY2hhaW5JRF1cbiAgICBjb25zdCBidWZmOiBCdWZmZXIgPSBCdWZmZXIuY29uY2F0KGJhcnIsIGJzaXplKVxuICAgIHJldHVybiBidWZmXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyBhIGJhc2UtNTggcmVwcmVzZW50YXRpb24gb2YgdGhlIFtbU3RhbmRhcmRCYXNlVHhdXS5cbiAgICovXG4gIHRvU3RyaW5nKCk6IHN0cmluZyB7XG4gICAgcmV0dXJuIGJpbnRvb2xzLmJ1ZmZlclRvQjU4KHRoaXMudG9CdWZmZXIoKSlcbiAgfVxuXG4gIGFic3RyYWN0IGNsb25lKCk6IHRoaXNcblxuICBhYnN0cmFjdCBjcmVhdGUoLi4uYXJnczogYW55W10pOiB0aGlzXG5cbiAgYWJzdHJhY3Qgc2VsZWN0KGlkOiBudW1iZXIsIC4uLmFyZ3M6IGFueVtdKTogdGhpc1xuXG4gIC8qKlxuICAgKiBDbGFzcyByZXByZXNlbnRpbmcgYSBTdGFuZGFyZEJhc2VUeCB3aGljaCBpcyB0aGUgZm91bmRhdGlvbiBmb3IgYWxsIHRyYW5zYWN0aW9ucy5cbiAgICpcbiAgICogQHBhcmFtIG5ldHdvcmtJRCBPcHRpb25hbCBuZXR3b3JrSUQsIFtbRGVmYXVsdE5ldHdvcmtJRF1dXG4gICAqIEBwYXJhbSBibG9ja2NoYWluSUQgT3B0aW9uYWwgYmxvY2tjaGFpbklELCBkZWZhdWx0IEJ1ZmZlci5hbGxvYygzMiwgMTYpXG4gICAqIEBwYXJhbSBvdXRzIE9wdGlvbmFsIGFycmF5IG9mIHRoZSBbW1RyYW5zZmVyYWJsZU91dHB1dF1dc1xuICAgKiBAcGFyYW0gaW5zIE9wdGlvbmFsIGFycmF5IG9mIHRoZSBbW1RyYW5zZmVyYWJsZUlucHV0XV1zXG4gICAqL1xuICBjb25zdHJ1Y3RvcihcbiAgICBuZXR3b3JrSUQ6IG51bWJlciA9IERlZmF1bHROZXR3b3JrSUQsXG4gICAgYmxvY2tjaGFpbklEOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoMzIsIDE2KVxuICApIHtcbiAgICBzdXBlcigpXG4gICAgdGhpcy5uZXR3b3JrSUQud3JpdGVVSW50MzJCRShuZXR3b3JrSUQsIDApXG4gICAgdGhpcy5ibG9ja2NoYWluSUQgPSBibG9ja2NoYWluSURcbiAgfVxufVxuXG4vKipcbiAqIENsYXNzIHJlcHJlc2VudGluZyBhbiB1bnNpZ25lZCB0cmFuc2FjdGlvbi5cbiAqL1xuZXhwb3J0IGFic3RyYWN0IGNsYXNzIEVWTVN0YW5kYXJkVW5zaWduZWRUeDxcbiAgS1BDbGFzcyBleHRlbmRzIFN0YW5kYXJkS2V5UGFpcixcbiAgS0NDbGFzcyBleHRlbmRzIFN0YW5kYXJkS2V5Q2hhaW48S1BDbGFzcz4sXG4gIFNCVHggZXh0ZW5kcyBFVk1TdGFuZGFyZEJhc2VUeDxLUENsYXNzLCBLQ0NsYXNzPlxuPiBleHRlbmRzIFNlcmlhbGl6YWJsZSB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIlN0YW5kYXJkVW5zaWduZWRUeFwiXG4gIHByb3RlY3RlZCBfdHlwZUlEID0gdW5kZWZpbmVkXG5cbiAgc2VyaWFsaXplKGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKTogb2JqZWN0IHtcbiAgICBsZXQgZmllbGRzOiBvYmplY3QgPSBzdXBlci5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgcmV0dXJuIHtcbiAgICAgIC4uLmZpZWxkcyxcbiAgICAgIGNvZGVjSUQ6IHNlcmlhbGl6ZXIuZW5jb2RlcihcbiAgICAgICAgdGhpcy5jb2RlY0lELFxuICAgICAgICBlbmNvZGluZyxcbiAgICAgICAgXCJudW1iZXJcIixcbiAgICAgICAgXCJkZWNpbWFsU3RyaW5nXCIsXG4gICAgICAgIDJcbiAgICAgICksXG4gICAgICB0cmFuc2FjdGlvbjogdGhpcy50cmFuc2FjdGlvbi5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgfVxuICB9XG5cbiAgZGVzZXJpYWxpemUoZmllbGRzOiBvYmplY3QsIGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKSB7XG4gICAgc3VwZXIuZGVzZXJpYWxpemUoZmllbGRzLCBlbmNvZGluZylcbiAgICB0aGlzLmNvZGVjSUQgPSBzZXJpYWxpemVyLmRlY29kZXIoXG4gICAgICBmaWVsZHNbXCJjb2RlY0lEXCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBcImRlY2ltYWxTdHJpbmdcIixcbiAgICAgIFwibnVtYmVyXCJcbiAgICApXG4gIH1cblxuICBwcm90ZWN0ZWQgY29kZWNJRDogbnVtYmVyID0gMFxuICBwcm90ZWN0ZWQgdHJhbnNhY3Rpb246IFNCVHhcblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgQ29kZWNJRCBhcyBhIG51bWJlclxuICAgKi9cbiAgZ2V0Q29kZWNJRCgpOiBudW1iZXIge1xuICAgIHJldHVybiB0aGlzLmNvZGVjSURcbiAgfVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSByZXByZXNlbnRhdGlvbiBvZiB0aGUgQ29kZWNJRFxuICAgKi9cbiAgZ2V0Q29kZWNJREJ1ZmZlcigpOiBCdWZmZXIge1xuICAgIGxldCBjb2RlY0J1ZjogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDIpXG4gICAgY29kZWNCdWYud3JpdGVVSW50MTZCRSh0aGlzLmNvZGVjSUQsIDApXG4gICAgcmV0dXJuIGNvZGVjQnVmXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgaW5wdXRUb3RhbCBhcyBhIEJOXG4gICAqL1xuICBnZXRJbnB1dFRvdGFsKGFzc2V0SUQ6IEJ1ZmZlcik6IEJOIHtcbiAgICBjb25zdCBpbnM6IFN0YW5kYXJkVHJhbnNmZXJhYmxlSW5wdXRbXSA9IFtdXG4gICAgY29uc3QgYUlESGV4OiBzdHJpbmcgPSBhc3NldElELnRvU3RyaW5nKFwiaGV4XCIpXG4gICAgbGV0IHRvdGFsOiBCTiA9IG5ldyBCTigwKVxuICAgIGlucy5mb3JFYWNoKChpbnB1dDogU3RhbmRhcmRUcmFuc2ZlcmFibGVJbnB1dCkgPT4ge1xuICAgICAgLy8gb25seSBjaGVjayBTdGFuZGFyZEFtb3VudElucHV0c1xuICAgICAgaWYgKFxuICAgICAgICBpbnB1dC5nZXRJbnB1dCgpIGluc3RhbmNlb2YgU3RhbmRhcmRBbW91bnRJbnB1dCAmJlxuICAgICAgICBhSURIZXggPT09IGlucHV0LmdldEFzc2V0SUQoKS50b1N0cmluZyhcImhleFwiKVxuICAgICAgKSB7XG4gICAgICAgIGNvbnN0IGkgPSBpbnB1dC5nZXRJbnB1dCgpIGFzIFN0YW5kYXJkQW1vdW50SW5wdXRcbiAgICAgICAgdG90YWwgPSB0b3RhbC5hZGQoaS5nZXRBbW91bnQoKSlcbiAgICAgIH1cbiAgICB9KVxuICAgIHJldHVybiB0b3RhbFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIG91dHB1dFRvdGFsIGFzIGEgQk5cbiAgICovXG4gIGdldE91dHB1dFRvdGFsKGFzc2V0SUQ6IEJ1ZmZlcik6IEJOIHtcbiAgICBjb25zdCBvdXRzOiBTdGFuZGFyZFRyYW5zZmVyYWJsZU91dHB1dFtdID0gW11cbiAgICBjb25zdCBhSURIZXg6IHN0cmluZyA9IGFzc2V0SUQudG9TdHJpbmcoXCJoZXhcIilcbiAgICBsZXQgdG90YWw6IEJOID0gbmV3IEJOKDApXG5cbiAgICBvdXRzLmZvckVhY2goKG91dDogU3RhbmRhcmRUcmFuc2ZlcmFibGVPdXRwdXQpID0+IHtcbiAgICAgIC8vIG9ubHkgY2hlY2sgU3RhbmRhcmRBbW91bnRPdXRwdXRcbiAgICAgIGlmIChcbiAgICAgICAgb3V0LmdldE91dHB1dCgpIGluc3RhbmNlb2YgU3RhbmRhcmRBbW91bnRPdXRwdXQgJiZcbiAgICAgICAgYUlESGV4ID09PSBvdXQuZ2V0QXNzZXRJRCgpLnRvU3RyaW5nKFwiaGV4XCIpXG4gICAgICApIHtcbiAgICAgICAgY29uc3Qgb3V0cHV0OiBTdGFuZGFyZEFtb3VudE91dHB1dCA9XG4gICAgICAgICAgb3V0LmdldE91dHB1dCgpIGFzIFN0YW5kYXJkQW1vdW50T3V0cHV0XG4gICAgICAgIHRvdGFsID0gdG90YWwuYWRkKG91dHB1dC5nZXRBbW91bnQoKSlcbiAgICAgIH1cbiAgICB9KVxuICAgIHJldHVybiB0b3RhbFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIG51bWJlciBvZiBidXJuZWQgdG9rZW5zIGFzIGEgQk5cbiAgICovXG4gIGdldEJ1cm4oYXNzZXRJRDogQnVmZmVyKTogQk4ge1xuICAgIHJldHVybiB0aGlzLmdldElucHV0VG90YWwoYXNzZXRJRCkuc3ViKHRoaXMuZ2V0T3V0cHV0VG90YWwoYXNzZXRJRCkpXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgVHJhbnNhY3Rpb25cbiAgICovXG4gIGFic3RyYWN0IGdldFRyYW5zYWN0aW9uKCk6IFNCVHhcblxuICBhYnN0cmFjdCBmcm9tQnVmZmVyKGJ5dGVzOiBCdWZmZXIsIG9mZnNldD86IG51bWJlcik6IG51bWJlclxuXG4gIHRvQnVmZmVyKCk6IEJ1ZmZlciB7XG4gICAgY29uc3QgY29kZWNJRDogQnVmZmVyID0gdGhpcy5nZXRDb2RlY0lEQnVmZmVyKClcbiAgICBjb25zdCB0eHR5cGU6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYyg0KVxuICAgIHR4dHlwZS53cml0ZVVJbnQzMkJFKHRoaXMudHJhbnNhY3Rpb24uZ2V0VHhUeXBlKCksIDApXG4gICAgY29uc3QgYmFzZWJ1ZmY6IEJ1ZmZlciA9IHRoaXMudHJhbnNhY3Rpb24udG9CdWZmZXIoKVxuICAgIHJldHVybiBCdWZmZXIuY29uY2F0KFxuICAgICAgW2NvZGVjSUQsIHR4dHlwZSwgYmFzZWJ1ZmZdLFxuICAgICAgY29kZWNJRC5sZW5ndGggKyB0eHR5cGUubGVuZ3RoICsgYmFzZWJ1ZmYubGVuZ3RoXG4gICAgKVxuICB9XG5cbiAgLyoqXG4gICAqIFNpZ25zIHRoaXMgW1tVbnNpZ25lZFR4XV0gYW5kIHJldHVybnMgc2lnbmVkIFtbU3RhbmRhcmRUeF1dXG4gICAqXG4gICAqIEBwYXJhbSBrYyBBbiBbW0tleUNoYWluXV0gdXNlZCBpbiBzaWduaW5nXG4gICAqXG4gICAqIEByZXR1cm5zIEEgc2lnbmVkIFtbU3RhbmRhcmRUeF1dXG4gICAqL1xuICBhYnN0cmFjdCBzaWduKFxuICAgIGtjOiBLQ0NsYXNzXG4gICk6IEVWTVN0YW5kYXJkVHg8XG4gICAgS1BDbGFzcyxcbiAgICBLQ0NsYXNzLFxuICAgIEVWTVN0YW5kYXJkVW5zaWduZWRUeDxLUENsYXNzLCBLQ0NsYXNzLCBTQlR4PlxuICA+XG5cbiAgY29uc3RydWN0b3IodHJhbnNhY3Rpb246IFNCVHggPSB1bmRlZmluZWQsIGNvZGVjSUQ6IG51bWJlciA9IDApIHtcbiAgICBzdXBlcigpXG4gICAgdGhpcy5jb2RlY0lEID0gY29kZWNJRFxuICAgIHRoaXMudHJhbnNhY3Rpb24gPSB0cmFuc2FjdGlvblxuICB9XG59XG5cbi8qKlxuICogQ2xhc3MgcmVwcmVzZW50aW5nIGEgc2lnbmVkIHRyYW5zYWN0aW9uLlxuICovXG5leHBvcnQgYWJzdHJhY3QgY2xhc3MgRVZNU3RhbmRhcmRUeDxcbiAgS1BDbGFzcyBleHRlbmRzIFN0YW5kYXJkS2V5UGFpcixcbiAgS0NDbGFzcyBleHRlbmRzIFN0YW5kYXJkS2V5Q2hhaW48S1BDbGFzcz4sXG4gIFNVQlR4IGV4dGVuZHMgRVZNU3RhbmRhcmRVbnNpZ25lZFR4PFxuICAgIEtQQ2xhc3MsXG4gICAgS0NDbGFzcyxcbiAgICBFVk1TdGFuZGFyZEJhc2VUeDxLUENsYXNzLCBLQ0NsYXNzPlxuICA+XG4+IGV4dGVuZHMgU2VyaWFsaXphYmxlIHtcbiAgcHJvdGVjdGVkIF90eXBlTmFtZSA9IFwiU3RhbmRhcmRUeFwiXG4gIHByb3RlY3RlZCBfdHlwZUlEID0gdW5kZWZpbmVkXG5cbiAgc2VyaWFsaXplKGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKTogb2JqZWN0IHtcbiAgICBsZXQgZmllbGRzOiBvYmplY3QgPSBzdXBlci5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgcmV0dXJuIHtcbiAgICAgIC4uLmZpZWxkcyxcbiAgICAgIHVuc2lnbmVkVHg6IHRoaXMudW5zaWduZWRUeC5zZXJpYWxpemUoZW5jb2RpbmcpLFxuICAgICAgY3JlZGVudGlhbHM6IHRoaXMuY3JlZGVudGlhbHMubWFwKChjKSA9PiBjLnNlcmlhbGl6ZShlbmNvZGluZykpXG4gICAgfVxuICB9XG5cbiAgcHJvdGVjdGVkIHVuc2lnbmVkVHg6IFNVQlR4ID0gdW5kZWZpbmVkXG4gIHByb3RlY3RlZCBjcmVkZW50aWFsczogQ3JlZGVudGlhbFtdID0gW11cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgW1tTdGFuZGFyZFVuc2lnbmVkVHhdXVxuICAgKi9cbiAgZ2V0VW5zaWduZWRUeCgpOiBTVUJUeCB7XG4gICAgcmV0dXJuIHRoaXMudW5zaWduZWRUeFxuICB9XG5cbiAgYWJzdHJhY3QgZnJvbUJ1ZmZlcihieXRlczogQnVmZmVyLCBvZmZzZXQ/OiBudW1iZXIpOiBudW1iZXJcblxuICAvKipcbiAgICogUmV0dXJucyBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IHJlcHJlc2VudGF0aW9uIG9mIHRoZSBbW1N0YW5kYXJkVHhdXS5cbiAgICovXG4gIHRvQnVmZmVyKCk6IEJ1ZmZlciB7XG4gICAgY29uc3QgdHhidWZmOiBCdWZmZXIgPSB0aGlzLnVuc2lnbmVkVHgudG9CdWZmZXIoKVxuICAgIGxldCBic2l6ZTogbnVtYmVyID0gdHhidWZmLmxlbmd0aFxuICAgIGNvbnN0IGNyZWRsZW46IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYyg0KVxuICAgIGNyZWRsZW4ud3JpdGVVSW50MzJCRSh0aGlzLmNyZWRlbnRpYWxzLmxlbmd0aCwgMClcbiAgICBjb25zdCBiYXJyOiBCdWZmZXJbXSA9IFt0eGJ1ZmYsIGNyZWRsZW5dXG4gICAgYnNpemUgKz0gY3JlZGxlbi5sZW5ndGhcbiAgICB0aGlzLmNyZWRlbnRpYWxzLmZvckVhY2goKGNyZWRlbnRpYWw6IENyZWRlbnRpYWwpID0+IHtcbiAgICAgIGNvbnN0IGNyZWRpZDogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDQpXG4gICAgICBjcmVkaWQud3JpdGVVSW50MzJCRShjcmVkZW50aWFsLmdldENyZWRlbnRpYWxJRCgpLCAwKVxuICAgICAgYmFyci5wdXNoKGNyZWRpZClcbiAgICAgIGJzaXplICs9IGNyZWRpZC5sZW5ndGhcbiAgICAgIGNvbnN0IGNyZWRidWZmOiBCdWZmZXIgPSBjcmVkZW50aWFsLnRvQnVmZmVyKClcbiAgICAgIGJzaXplICs9IGNyZWRidWZmLmxlbmd0aFxuICAgICAgYmFyci5wdXNoKGNyZWRidWZmKVxuICAgIH0pXG4gICAgY29uc3QgYnVmZjogQnVmZmVyID0gQnVmZmVyLmNvbmNhdChiYXJyLCBic2l6ZSlcbiAgICByZXR1cm4gYnVmZlxuICB9XG5cbiAgLyoqXG4gICAqIFRha2VzIGEgYmFzZS01OCBzdHJpbmcgY29udGFpbmluZyBhbiBbW1N0YW5kYXJkVHhdXSwgcGFyc2VzIGl0LCBwb3B1bGF0ZXMgdGhlIGNsYXNzLCBhbmQgcmV0dXJucyB0aGUgbGVuZ3RoIG9mIHRoZSBUeCBpbiBieXRlcy5cbiAgICpcbiAgICogQHBhcmFtIHNlcmlhbGl6ZWQgQSBiYXNlLTU4IHN0cmluZyBjb250YWluaW5nIGEgcmF3IFtbU3RhbmRhcmRUeF1dXG4gICAqXG4gICAqIEByZXR1cm5zIFRoZSBsZW5ndGggb2YgdGhlIHJhdyBbW1N0YW5kYXJkVHhdXVxuICAgKlxuICAgKiBAcmVtYXJrc1xuICAgKiB1bmxpa2UgbW9zdCBmcm9tU3RyaW5ncywgaXQgZXhwZWN0cyB0aGUgc3RyaW5nIHRvIGJlIHNlcmlhbGl6ZWQgaW4gY2I1OCBmb3JtYXRcbiAgICovXG4gIGZyb21TdHJpbmcoc2VyaWFsaXplZDogc3RyaW5nKTogbnVtYmVyIHtcbiAgICByZXR1cm4gdGhpcy5mcm9tQnVmZmVyKGJpbnRvb2xzLmNiNThEZWNvZGUoc2VyaWFsaXplZCkpXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyBhIGNiNTggcmVwcmVzZW50YXRpb24gb2YgdGhlIFtbU3RhbmRhcmRUeF1dLlxuICAgKlxuICAgKiBAcmVtYXJrc1xuICAgKiB1bmxpa2UgbW9zdCB0b1N0cmluZ3MsIHRoaXMgcmV0dXJucyBpbiBjYjU4IHNlcmlhbGl6YXRpb24gZm9ybWF0XG4gICAqL1xuICB0b1N0cmluZygpOiBzdHJpbmcge1xuICAgIHJldHVybiBiaW50b29scy5jYjU4RW5jb2RlKHRoaXMudG9CdWZmZXIoKSlcbiAgfVxuXG4gIC8qKlxuICAgKiBDbGFzcyByZXByZXNlbnRpbmcgYSBzaWduZWQgdHJhbnNhY3Rpb24uXG4gICAqXG4gICAqIEBwYXJhbSB1bnNpZ25lZFR4IE9wdGlvbmFsIFtbU3RhbmRhcmRVbnNpZ25lZFR4XV1cbiAgICogQHBhcmFtIHNpZ25hdHVyZXMgT3B0aW9uYWwgYXJyYXkgb2YgW1tDcmVkZW50aWFsXV1zXG4gICAqL1xuICBjb25zdHJ1Y3RvcihcbiAgICB1bnNpZ25lZFR4OiBTVUJUeCA9IHVuZGVmaW5lZCxcbiAgICBjcmVkZW50aWFsczogQ3JlZGVudGlhbFtdID0gdW5kZWZpbmVkXG4gICkge1xuICAgIHN1cGVyKClcbiAgICBpZiAodHlwZW9mIHVuc2lnbmVkVHggIT09IFwidW5kZWZpbmVkXCIpIHtcbiAgICAgIHRoaXMudW5zaWduZWRUeCA9IHVuc2lnbmVkVHhcbiAgICAgIGlmICh0eXBlb2YgY3JlZGVudGlhbHMgIT09IFwidW5kZWZpbmVkXCIpIHtcbiAgICAgICAgdGhpcy5jcmVkZW50aWFscyA9IGNyZWRlbnRpYWxzXG4gICAgICB9XG4gICAgfVxuICB9XG59XG4iXX0=