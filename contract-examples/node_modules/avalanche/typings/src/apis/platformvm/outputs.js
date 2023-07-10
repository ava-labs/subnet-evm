"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.SECPOwnerOutput = exports.StakeableLockOut = exports.SECPTransferOutput = exports.AmountOutput = exports.ParseableOutput = exports.TransferableOutput = exports.SelectOutputClass = void 0;
/**
 * @packageDocumentation
 * @module API-PlatformVM-Outputs
 */
const buffer_1 = require("buffer/");
const bintools_1 = __importDefault(require("../../utils/bintools"));
const constants_1 = require("./constants");
const output_1 = require("../../common/output");
const serialization_1 = require("../../utils/serialization");
const errors_1 = require("../../utils/errors");
const bintools = bintools_1.default.getInstance();
const serialization = serialization_1.Serialization.getInstance();
/**
 * Takes a buffer representing the output and returns the proper Output instance.
 *
 * @param outputid A number representing the inputID parsed prior to the bytes passed in
 *
 * @returns An instance of an [[Output]]-extended class.
 */
const SelectOutputClass = (outputid, ...args) => {
    if (outputid == constants_1.PlatformVMConstants.SECPXFEROUTPUTID) {
        return new SECPTransferOutput(...args);
    }
    else if (outputid == constants_1.PlatformVMConstants.SECPOWNEROUTPUTID) {
        return new SECPOwnerOutput(...args);
    }
    else if (outputid == constants_1.PlatformVMConstants.STAKEABLELOCKOUTID) {
        return new StakeableLockOut(...args);
    }
    throw new errors_1.OutputIdError("Error - SelectOutputClass: unknown outputid " + outputid);
};
exports.SelectOutputClass = SelectOutputClass;
class TransferableOutput extends output_1.StandardTransferableOutput {
    constructor() {
        super(...arguments);
        this._typeName = "TransferableOutput";
        this._typeID = undefined;
    }
    //serialize is inherited
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.output = (0, exports.SelectOutputClass)(fields["output"]["_typeID"]);
        this.output.deserialize(fields["output"], encoding);
    }
    fromBuffer(bytes, offset = 0) {
        this.assetID = bintools.copyFrom(bytes, offset, offset + constants_1.PlatformVMConstants.ASSETIDLEN);
        offset += constants_1.PlatformVMConstants.ASSETIDLEN;
        const outputid = bintools
            .copyFrom(bytes, offset, offset + 4)
            .readUInt32BE(0);
        offset += 4;
        this.output = (0, exports.SelectOutputClass)(outputid);
        return this.output.fromBuffer(bytes, offset);
    }
}
exports.TransferableOutput = TransferableOutput;
class ParseableOutput extends output_1.StandardParseableOutput {
    constructor() {
        super(...arguments);
        this._typeName = "ParseableOutput";
        this._typeID = undefined;
    }
    //serialize is inherited
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.output = (0, exports.SelectOutputClass)(fields["output"]["_typeID"]);
        this.output.deserialize(fields["output"], encoding);
    }
    fromBuffer(bytes, offset = 0) {
        const outputid = bintools
            .copyFrom(bytes, offset, offset + 4)
            .readUInt32BE(0);
        offset += 4;
        this.output = (0, exports.SelectOutputClass)(outputid);
        return this.output.fromBuffer(bytes, offset);
    }
}
exports.ParseableOutput = ParseableOutput;
class AmountOutput extends output_1.StandardAmountOutput {
    constructor() {
        super(...arguments);
        this._typeName = "AmountOutput";
        this._typeID = undefined;
    }
    //serialize and deserialize both are inherited
    /**
     * @param assetID An assetID which is wrapped around the Buffer of the Output
     */
    makeTransferable(assetID) {
        return new TransferableOutput(assetID, this);
    }
    select(id, ...args) {
        return (0, exports.SelectOutputClass)(id, ...args);
    }
}
exports.AmountOutput = AmountOutput;
/**
 * An [[Output]] class which specifies an Output that carries an ammount for an assetID and uses secp256k1 signature scheme.
 */
class SECPTransferOutput extends AmountOutput {
    constructor() {
        super(...arguments);
        this._typeName = "SECPTransferOutput";
        this._typeID = constants_1.PlatformVMConstants.SECPXFEROUTPUTID;
    }
    //serialize and deserialize both are inherited
    /**
     * Returns the outputID for this output
     */
    getOutputID() {
        return this._typeID;
    }
    create(...args) {
        return new SECPTransferOutput(...args);
    }
    clone() {
        const newout = this.create();
        newout.fromBuffer(this.toBuffer());
        return newout;
    }
}
exports.SECPTransferOutput = SECPTransferOutput;
/**
 * An [[Output]] class which specifies an input that has a locktime which can also enable staking of the value held, preventing transfers but not validation.
 */
class StakeableLockOut extends AmountOutput {
    /**
     * A [[Output]] class which specifies a [[ParseableOutput]] that has a locktime which can also enable staking of the value held, preventing transfers but not validation.
     *
     * @param amount A {@link https://github.com/indutny/bn.js/|BN} representing the amount in the output
     * @param addresses An array of {@link https://github.com/feross/buffer|Buffer}s representing addresses
     * @param locktime A {@link https://github.com/indutny/bn.js/|BN} representing the locktime
     * @param threshold A number representing the the threshold number of signers required to sign the transaction
     * @param stakeableLocktime A {@link https://github.com/indutny/bn.js/|BN} representing the stakeable locktime
     * @param transferableOutput A [[ParseableOutput]] which is embedded into this output.
     */
    constructor(amount = undefined, addresses = undefined, locktime = undefined, threshold = undefined, stakeableLocktime = undefined, transferableOutput = undefined) {
        super(amount, addresses, locktime, threshold);
        this._typeName = "StakeableLockOut";
        this._typeID = constants_1.PlatformVMConstants.STAKEABLELOCKOUTID;
        if (typeof stakeableLocktime !== "undefined") {
            this.stakeableLocktime = bintools.fromBNToBuffer(stakeableLocktime, 8);
        }
        if (typeof transferableOutput !== "undefined") {
            this.transferableOutput = transferableOutput;
            this.synchronize();
        }
    }
    //serialize and deserialize both are inherited
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        let outobj = Object.assign(Object.assign({}, fields), { stakeableLocktime: serialization.encoder(this.stakeableLocktime, encoding, "Buffer", "decimalString", 8), transferableOutput: this.transferableOutput.serialize(encoding) });
        delete outobj["addresses"];
        delete outobj["locktime"];
        delete outobj["threshold"];
        delete outobj["amount"];
        return outobj;
    }
    deserialize(fields, encoding = "hex") {
        fields["addresses"] = [];
        fields["locktime"] = "0";
        fields["threshold"] = "1";
        fields["amount"] = "99";
        super.deserialize(fields, encoding);
        this.stakeableLocktime = serialization.decoder(fields["stakeableLocktime"], encoding, "decimalString", "Buffer", 8);
        this.transferableOutput = new ParseableOutput();
        this.transferableOutput.deserialize(fields["transferableOutput"], encoding);
        this.synchronize();
    }
    //call this every time you load in data
    synchronize() {
        let output = this.transferableOutput.getOutput();
        this.addresses = output.getAddresses().map((a) => {
            let addr = new output_1.Address();
            addr.fromBuffer(a);
            return addr;
        });
        this.numaddrs = buffer_1.Buffer.alloc(4);
        this.numaddrs.writeUInt32BE(this.addresses.length, 0);
        this.locktime = bintools.fromBNToBuffer(output.getLocktime(), 8);
        this.threshold = buffer_1.Buffer.alloc(4);
        this.threshold.writeUInt32BE(output.getThreshold(), 0);
        this.amount = bintools.fromBNToBuffer(output.getAmount(), 8);
        this.amountValue = output.getAmount();
    }
    getStakeableLocktime() {
        return bintools.fromBufferToBN(this.stakeableLocktime);
    }
    getTransferableOutput() {
        return this.transferableOutput;
    }
    /**
     * @param assetID An assetID which is wrapped around the Buffer of the Output
     */
    makeTransferable(assetID) {
        return new TransferableOutput(assetID, this);
    }
    select(id, ...args) {
        return (0, exports.SelectOutputClass)(id, ...args);
    }
    /**
     * Popuates the instance from a {@link https://github.com/feross/buffer|Buffer} representing the [[StakeableLockOut]] and returns the size of the output.
     */
    fromBuffer(outbuff, offset = 0) {
        this.stakeableLocktime = bintools.copyFrom(outbuff, offset, offset + 8);
        offset += 8;
        this.transferableOutput = new ParseableOutput();
        offset = this.transferableOutput.fromBuffer(outbuff, offset);
        this.synchronize();
        return offset;
    }
    /**
     * Returns the buffer representing the [[StakeableLockOut]] instance.
     */
    toBuffer() {
        let xferoutBuff = this.transferableOutput.toBuffer();
        const bsize = this.stakeableLocktime.length + xferoutBuff.length;
        const barr = [this.stakeableLocktime, xferoutBuff];
        return buffer_1.Buffer.concat(barr, bsize);
    }
    /**
     * Returns the outputID for this output
     */
    getOutputID() {
        return this._typeID;
    }
    create(...args) {
        return new StakeableLockOut(...args);
    }
    clone() {
        const newout = this.create();
        newout.fromBuffer(this.toBuffer());
        return newout;
    }
}
exports.StakeableLockOut = StakeableLockOut;
/**
 * An [[Output]] class which only specifies an Output ownership and uses secp256k1 signature scheme.
 */
class SECPOwnerOutput extends output_1.Output {
    constructor() {
        super(...arguments);
        this._typeName = "SECPOwnerOutput";
        this._typeID = constants_1.PlatformVMConstants.SECPOWNEROUTPUTID;
    }
    //serialize and deserialize both are inherited
    /**
     * Returns the outputID for this output
     */
    getOutputID() {
        return this._typeID;
    }
    /**
     *
     * @param assetID An assetID which is wrapped around the Buffer of the Output
     */
    makeTransferable(assetID) {
        return new TransferableOutput(assetID, this);
    }
    create(...args) {
        return new SECPOwnerOutput(...args);
    }
    clone() {
        const newout = this.create();
        newout.fromBuffer(this.toBuffer());
        return newout;
    }
    select(id, ...args) {
        return (0, exports.SelectOutputClass)(id, ...args);
    }
}
exports.SECPOwnerOutput = SECPOwnerOutput;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoib3V0cHV0cy5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3NyYy9hcGlzL3BsYXRmb3Jtdm0vb3V0cHV0cy50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7QUFBQTs7O0dBR0c7QUFDSCxvQ0FBZ0M7QUFDaEMsb0VBQTJDO0FBQzNDLDJDQUFpRDtBQUNqRCxnREFNNEI7QUFDNUIsNkRBQTZFO0FBRTdFLCtDQUFrRDtBQUVsRCxNQUFNLFFBQVEsR0FBYSxrQkFBUSxDQUFDLFdBQVcsRUFBRSxDQUFBO0FBQ2pELE1BQU0sYUFBYSxHQUFrQiw2QkFBYSxDQUFDLFdBQVcsRUFBRSxDQUFBO0FBRWhFOzs7Ozs7R0FNRztBQUNJLE1BQU0saUJBQWlCLEdBQUcsQ0FBQyxRQUFnQixFQUFFLEdBQUcsSUFBVyxFQUFVLEVBQUU7SUFDNUUsSUFBSSxRQUFRLElBQUksK0JBQW1CLENBQUMsZ0JBQWdCLEVBQUU7UUFDcEQsT0FBTyxJQUFJLGtCQUFrQixDQUFDLEdBQUcsSUFBSSxDQUFDLENBQUE7S0FDdkM7U0FBTSxJQUFJLFFBQVEsSUFBSSwrQkFBbUIsQ0FBQyxpQkFBaUIsRUFBRTtRQUM1RCxPQUFPLElBQUksZUFBZSxDQUFDLEdBQUcsSUFBSSxDQUFDLENBQUE7S0FDcEM7U0FBTSxJQUFJLFFBQVEsSUFBSSwrQkFBbUIsQ0FBQyxrQkFBa0IsRUFBRTtRQUM3RCxPQUFPLElBQUksZ0JBQWdCLENBQUMsR0FBRyxJQUFJLENBQUMsQ0FBQTtLQUNyQztJQUNELE1BQU0sSUFBSSxzQkFBYSxDQUNyQiw4Q0FBOEMsR0FBRyxRQUFRLENBQzFELENBQUE7QUFDSCxDQUFDLENBQUE7QUFYWSxRQUFBLGlCQUFpQixxQkFXN0I7QUFFRCxNQUFhLGtCQUFtQixTQUFRLG1DQUEwQjtJQUFsRTs7UUFDWSxjQUFTLEdBQUcsb0JBQW9CLENBQUE7UUFDaEMsWUFBTyxHQUFHLFNBQVMsQ0FBQTtJQXdCL0IsQ0FBQztJQXRCQyx3QkFBd0I7SUFFeEIs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QVEsRUFBRSxDQUFDLENBQUE7UUFDbEMsT0FBTyxNQUFjLENBQUE7SUFDdkIsQ0FBQztJQUVELE1BQU0sQ0FBQyxFQUFVLEVBQUUsR0FBRyxJQUFXO1FBQy9CLE9BQU8sSUFBQSx5QkFBaUIsRUFBQyxFQUFFLEVBQUUsR0FBRyxJQUFJLENBQUMsQ0FBQTtJQUN2QyxDQUFDO0NBQ0Y7QUFsQ0QsMENBa0NDIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgQVBJLVBsYXRmb3JtVk0tT3V0cHV0c1xuICovXG5pbXBvcnQgeyBCdWZmZXIgfSBmcm9tIFwiYnVmZmVyL1wiXG5pbXBvcnQgQmluVG9vbHMgZnJvbSBcIi4uLy4uL3V0aWxzL2JpbnRvb2xzXCJcbmltcG9ydCB7IFBsYXRmb3JtVk1Db25zdGFudHMgfSBmcm9tIFwiLi9jb25zdGFudHNcIlxuaW1wb3J0IHtcbiAgT3V0cHV0LFxuICBTdGFuZGFyZEFtb3VudE91dHB1dCxcbiAgU3RhbmRhcmRUcmFuc2ZlcmFibGVPdXRwdXQsXG4gIFN0YW5kYXJkUGFyc2VhYmxlT3V0cHV0LFxuICBBZGRyZXNzXG59IGZyb20gXCIuLi8uLi9jb21tb24vb3V0cHV0XCJcbmltcG9ydCB7IFNlcmlhbGl6YXRpb24sIFNlcmlhbGl6ZWRFbmNvZGluZyB9IGZyb20gXCIuLi8uLi91dGlscy9zZXJpYWxpemF0aW9uXCJcbmltcG9ydCBCTiBmcm9tIFwiYm4uanNcIlxuaW1wb3J0IHsgT3V0cHV0SWRFcnJvciB9IGZyb20gXCIuLi8uLi91dGlscy9lcnJvcnNcIlxuXG5jb25zdCBiaW50b29sczogQmluVG9vbHMgPSBCaW5Ub29scy5nZXRJbnN0YW5jZSgpXG5jb25zdCBzZXJpYWxpemF0aW9uOiBTZXJpYWxpemF0aW9uID0gU2VyaWFsaXphdGlvbi5nZXRJbnN0YW5jZSgpXG5cbi8qKlxuICogVGFrZXMgYSBidWZmZXIgcmVwcmVzZW50aW5nIHRoZSBvdXRwdXQgYW5kIHJldHVybnMgdGhlIHByb3BlciBPdXRwdXQgaW5zdGFuY2UuXG4gKlxuICogQHBhcmFtIG91dHB1dGlkIEEgbnVtYmVyIHJlcHJlc2VudGluZyB0aGUgaW5wdXRJRCBwYXJzZWQgcHJpb3IgdG8gdGhlIGJ5dGVzIHBhc3NlZCBpblxuICpcbiAqIEByZXR1cm5zIEFuIGluc3RhbmNlIG9mIGFuIFtbT3V0cHV0XV0tZXh0ZW5kZWQgY2xhc3MuXG4gKi9cbmV4cG9ydCBjb25zdCBTZWxlY3RPdXRwdXRDbGFzcyA9IChvdXRwdXRpZDogbnVtYmVyLCAuLi5hcmdzOiBhbnlbXSk6IE91dHB1dCA9PiB7XG4gIGlmIChvdXRwdXRpZCA9PSBQbGF0Zm9ybVZNQ29uc3RhbnRzLlNFQ1BYRkVST1VUUFVUSUQpIHtcbiAgICByZXR1cm4gbmV3IFNFQ1BUcmFuc2Zlck91dHB1dCguLi5hcmdzKVxuICB9IGVsc2UgaWYgKG91dHB1dGlkID09IFBsYXRmb3JtVk1Db25zdGFudHMuU0VDUE9XTkVST1VUUFVUSUQpIHtcbiAgICByZXR1cm4gbmV3IFNFQ1BPd25lck91dHB1dCguLi5hcmdzKVxuICB9IGVsc2UgaWYgKG91dHB1dGlkID09IFBsYXRmb3JtVk1Db25zdGFudHMuU1RBS0VBQkxFTE9DS09VVElEKSB7XG4gICAgcmV0dXJuIG5ldyBTdGFrZWFibGVMb2NrT3V0KC4uLmFyZ3MpXG4gIH1cbiAgdGhyb3cgbmV3IE91dHB1dElkRXJyb3IoXG4gICAgXCJFcnJvciAtIFNlbGVjdE91dHB1dENsYXNzOiB1bmtub3duIG91dHB1dGlkIFwiICsgb3V0cHV0aWRcbiAgKVxufVxuXG5leHBvcnQgY2xhc3MgVHJhbnNmZXJhYmxlT3V0cHV0IGV4dGVuZHMgU3RhbmRhcmRUcmFuc2ZlcmFibGVPdXRwdXQge1xuICBwcm90ZWN0ZWQgX3R5cGVOYW1lID0gXCJUcmFuc2ZlcmFibGVPdXRwdXRcIlxuICBwcm90ZWN0ZWQgX3R5cGVJRCA9IHVuZGVmaW5lZFxuXG4gIC8vc2VyaWFsaXplIGlzIGluaGVyaXRlZFxuXG4gIGRlc2VyaWFsaXplKGZpZWxkczogb2JqZWN0LCBlbmNvZGluZzogU2VyaWFsaXplZEVuY29kaW5nID0gXCJoZXhcIikge1xuICAgIHN1cGVyLmRlc2VyaWFsaXplKGZpZWxkcywgZW5jb2RpbmcpXG4gICAgdGhpcy5vdXRwdXQgPSBTZWxlY3RPdXRwdXRDbGFzcyhmaWVsZHNbXCJvdXRwdXRcIl1bXCJfdHlwZUlEXCJdKVxuICAgIHRoaXMub3V0cHV0LmRlc2VyaWFsaXplKGZpZWxkc1tcIm91dHB1dFwiXSwgZW5jb2RpbmcpXG4gIH1cblxuICBmcm9tQnVmZmVyKGJ5dGVzOiBCdWZmZXIsIG9mZnNldDogbnVtYmVyID0gMCk6IG51bWJlciB7XG4gICAgdGhpcy5hc3NldElEID0gYmludG9vbHMuY29weUZyb20oXG4gICAgICBieXRlcyxcbiAgICAgIG9mZnNldCxcbiAgICAgIG9mZnNldCArIFBsYXRmb3JtVk1Db25zdGFudHMuQVNTRVRJRExFTlxuICAgIClcbiAgICBvZmZzZXQgKz0gUGxhdGZvcm1WTUNvbnN0YW50cy5BU1NFVElETEVOXG4gICAgY29uc3Qgb3V0cHV0aWQ6IG51bWJlciA9IGJpbnRvb2xzXG4gICAgICAuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgNClcbiAgICAgIC5yZWFkVUludDMyQkUoMClcbiAgICBvZmZzZXQgKz0gNFxuICAgIHRoaXMub3V0cHV0ID0gU2VsZWN0T3V0cHV0Q2xhc3Mob3V0cHV0aWQpXG4gICAgcmV0dXJuIHRoaXMub3V0cHV0LmZyb21CdWZmZXIoYnl0ZXMsIG9mZnNldClcbiAgfVxufVxuXG5leHBvcnQgY2xhc3MgUGFyc2VhYmxlT3V0cHV0IGV4dGVuZHMgU3RhbmRhcmRQYXJzZWFibGVPdXRwdXQge1xuICBwcm90ZWN0ZWQgX3R5cGVOYW1lID0gXCJQYXJzZWFibGVPdXRwdXRcIlxuICBwcm90ZWN0ZWQgX3R5cGVJRCA9IHVuZGVmaW5lZFxuXG4gIC8vc2VyaWFsaXplIGlzIGluaGVyaXRlZFxuXG4gIGRlc2VyaWFsaXplKGZpZWxkczogb2JqZWN0LCBlbmNvZGluZzogU2VyaWFsaXplZEVuY29kaW5nID0gXCJoZXhcIikge1xuICAgIHN1cGVyLmRlc2VyaWFsaXplKGZpZWxkcywgZW5jb2RpbmcpXG4gICAgdGhpcy5vdXRwdXQgPSBTZWxlY3RPdXRwdXRDbGFzcyhmaWVsZHNbXCJvdXRwdXRcIl1bXCJfdHlwZUlEXCJdKVxuICAgIHRoaXMub3V0cHV0LmRlc2VyaWFsaXplKGZpZWxkc1tcIm91dHB1dFwiXSwgZW5jb2RpbmcpXG4gIH1cblxuICBmcm9tQnVmZmVyKGJ5dGVzOiBCdWZmZXIsIG9mZnNldDogbnVtYmVyID0gMCk6IG51bWJlciB7XG4gICAgY29uc3Qgb3V0cHV0aWQ6IG51bWJlciA9IGJpbnRvb2xzXG4gICAgICAuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgNClcbiAgICAgIC5yZWFkVUludDMyQkUoMClcbiAgICBvZmZzZXQgKz0gNFxuICAgIHRoaXMub3V0cHV0ID0gU2VsZWN0T3V0cHV0Q2xhc3Mob3V0cHV0aWQpXG4gICAgcmV0dXJuIHRoaXMub3V0cHV0LmZyb21CdWZmZXIoYnl0ZXMsIG9mZnNldClcbiAgfVxufVxuXG5leHBvcnQgYWJzdHJhY3QgY2xhc3MgQW1vdW50T3V0cHV0IGV4dGVuZHMgU3RhbmRhcmRBbW91bnRPdXRwdXQge1xuICBwcm90ZWN0ZWQgX3R5cGVOYW1lID0gXCJBbW91bnRPdXRwdXRcIlxuICBwcm90ZWN0ZWQgX3R5cGVJRCA9IHVuZGVmaW5lZFxuXG4gIC8vc2VyaWFsaXplIGFuZCBkZXNlcmlhbGl6ZSBib3RoIGFyZSBpbmhlcml0ZWRcblxuICAvKipcbiAgICogQHBhcmFtIGFzc2V0SUQgQW4gYXNzZXRJRCB3aGljaCBpcyB3cmFwcGVkIGFyb3VuZCB0aGUgQnVmZmVyIG9mIHRoZSBPdXRwdXRcbiAgICovXG4gIG1ha2VUcmFuc2ZlcmFibGUoYXNzZXRJRDogQnVmZmVyKTogVHJhbnNmZXJhYmxlT3V0cHV0IHtcbiAgICByZXR1cm4gbmV3IFRyYW5zZmVyYWJsZU91dHB1dChhc3NldElELCB0aGlzKVxuICB9XG5cbiAgc2VsZWN0KGlkOiBudW1iZXIsIC4uLmFyZ3M6IGFueVtdKTogT3V0cHV0IHtcbiAgICByZXR1cm4gU2VsZWN0T3V0cHV0Q2xhc3MoaWQsIC4uLmFyZ3MpXG4gIH1cbn1cblxuLyoqXG4gKiBBbiBbW091dHB1dF1dIGNsYXNzIHdoaWNoIHNwZWNpZmllcyBhbiBPdXRwdXQgdGhhdCBjYXJyaWVzIGFuIGFtbW91bnQgZm9yIGFuIGFzc2V0SUQgYW5kIHVzZXMgc2VjcDI1NmsxIHNpZ25hdHVyZSBzY2hlbWUuXG4gKi9cbmV4cG9ydCBjbGFzcyBTRUNQVHJhbnNmZXJPdXRwdXQgZXh0ZW5kcyBBbW91bnRPdXRwdXQge1xuICBwcm90ZWN0ZWQgX3R5cGVOYW1lID0gXCJTRUNQVHJhbnNmZXJPdXRwdXRcIlxuICBwcm90ZWN0ZWQgX3R5cGVJRCA9IFBsYXRmb3JtVk1Db25zdGFudHMuU0VDUFhGRVJPVVRQVVRJRFxuXG4gIC8vc2VyaWFsaXplIGFuZCBkZXNlcmlhbGl6ZSBib3RoIGFyZSBpbmhlcml0ZWRcblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgb3V0cHV0SUQgZm9yIHRoaXMgb3V0cHV0XG4gICAqL1xuICBnZXRPdXRwdXRJRCgpOiBudW1iZXIge1xuICAgIHJldHVybiB0aGlzLl90eXBlSURcbiAgfVxuXG4gIGNyZWF0ZSguLi5hcmdzOiBhbnlbXSk6IHRoaXMge1xuICAgIHJldHVybiBuZXcgU0VDUFRyYW5zZmVyT3V0cHV0KC4uLmFyZ3MpIGFzIHRoaXNcbiAgfVxuXG4gIGNsb25lKCk6IHRoaXMge1xuICAgIGNvbnN0IG5ld291dDogU0VDUFRyYW5zZmVyT3V0cHV0ID0gdGhpcy5jcmVhdGUoKVxuICAgIG5ld291dC5mcm9tQnVmZmVyKHRoaXMudG9CdWZmZXIoKSlcbiAgICByZXR1cm4gbmV3b3V0IGFzIHRoaXNcbiAgfVxufVxuXG4vKipcbiAqIEFuIFtbT3V0cHV0XV0gY2xhc3Mgd2hpY2ggc3BlY2lmaWVzIGFuIGlucHV0IHRoYXQgaGFzIGEgbG9ja3RpbWUgd2hpY2ggY2FuIGFsc28gZW5hYmxlIHN0YWtpbmcgb2YgdGhlIHZhbHVlIGhlbGQsIHByZXZlbnRpbmcgdHJhbnNmZXJzIGJ1dCBub3QgdmFsaWRhdGlvbi5cbiAqL1xuZXhwb3J0IGNsYXNzIFN0YWtlYWJsZUxvY2tPdXQgZXh0ZW5kcyBBbW91bnRPdXRwdXQge1xuICBwcm90ZWN0ZWQgX3R5cGVOYW1lID0gXCJTdGFrZWFibGVMb2NrT3V0XCJcbiAgcHJvdGVjdGVkIF90eXBlSUQgPSBQbGF0Zm9ybVZNQ29uc3RhbnRzLlNUQUtFQUJMRUxPQ0tPVVRJRFxuXG4gIC8vc2VyaWFsaXplIGFuZCBkZXNlcmlhbGl6ZSBib3RoIGFyZSBpbmhlcml0ZWRcblxuICBzZXJpYWxpemUoZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpOiBvYmplY3Qge1xuICAgIGxldCBmaWVsZHM6IG9iamVjdCA9IHN1cGVyLnNlcmlhbGl6ZShlbmNvZGluZylcbiAgICBsZXQgb3V0b2JqOiBvYmplY3QgPSB7XG4gICAgICAuLi5maWVsZHMsIC8vaW5jbHVkZWQgYW55d2F5eXl5Li4uIG5vdCBpZGVhbFxuICAgICAgc3Rha2VhYmxlTG9ja3RpbWU6IHNlcmlhbGl6YXRpb24uZW5jb2RlcihcbiAgICAgICAgdGhpcy5zdGFrZWFibGVMb2NrdGltZSxcbiAgICAgICAgZW5jb2RpbmcsXG4gICAgICAgIFwiQnVmZmVyXCIsXG4gICAgICAgIFwiZGVjaW1hbFN0cmluZ1wiLFxuICAgICAgICA4XG4gICAgICApLFxuICAgICAgdHJhbnNmZXJhYmxlT3V0cHV0OiB0aGlzLnRyYW5zZmVyYWJsZU91dHB1dC5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgfVxuICAgIGRlbGV0ZSBvdXRvYmpbXCJhZGRyZXNzZXNcIl1cbiAgICBkZWxldGUgb3V0b2JqW1wibG9ja3RpbWVcIl1cbiAgICBkZWxldGUgb3V0b2JqW1widGhyZXNob2xkXCJdXG4gICAgZGVsZXRlIG91dG9ialtcImFtb3VudFwiXVxuICAgIHJldHVybiBvdXRvYmpcbiAgfVxuICBkZXNlcmlhbGl6ZShmaWVsZHM6IG9iamVjdCwgZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpIHtcbiAgICBmaWVsZHNbXCJhZGRyZXNzZXNcIl0gPSBbXVxuICAgIGZpZWxkc1tcImxvY2t0aW1lXCJdID0gXCIwXCJcbiAgICBmaWVsZHNbXCJ0aHJlc2hvbGRcIl0gPSBcIjFcIlxuICAgIGZpZWxkc1tcImFtb3VudFwiXSA9IFwiOTlcIlxuICAgIHN1cGVyLmRlc2VyaWFsaXplKGZpZWxkcywgZW5jb2RpbmcpXG4gICAgdGhpcy5zdGFrZWFibGVMb2NrdGltZSA9IHNlcmlhbGl6YXRpb24uZGVjb2RlcihcbiAgICAgIGZpZWxkc1tcInN0YWtlYWJsZUxvY2t0aW1lXCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBcImRlY2ltYWxTdHJpbmdcIixcbiAgICAgIFwiQnVmZmVyXCIsXG4gICAgICA4XG4gICAgKVxuICAgIHRoaXMudHJhbnNmZXJhYmxlT3V0cHV0ID0gbmV3IFBhcnNlYWJsZU91dHB1dCgpXG4gICAgdGhpcy50cmFuc2ZlcmFibGVPdXRwdXQuZGVzZXJpYWxpemUoZmllbGRzW1widHJhbnNmZXJhYmxlT3V0cHV0XCJdLCBlbmNvZGluZylcbiAgICB0aGlzLnN5bmNocm9uaXplKClcbiAgfVxuXG4gIHByb3RlY3RlZCBzdGFrZWFibGVMb2NrdGltZTogQnVmZmVyXG4gIHByb3RlY3RlZCB0cmFuc2ZlcmFibGVPdXRwdXQ6IFBhcnNlYWJsZU91dHB1dFxuXG4gIC8vY2FsbCB0aGlzIGV2ZXJ5IHRpbWUgeW91IGxvYWQgaW4gZGF0YVxuICBwcml2YXRlIHN5bmNocm9uaXplKCkge1xuICAgIGxldCBvdXRwdXQ6IEFtb3VudE91dHB1dCA9XG4gICAgICB0aGlzLnRyYW5zZmVyYWJsZU91dHB1dC5nZXRPdXRwdXQoKSBhcyBBbW91bnRPdXRwdXRcbiAgICB0aGlzLmFkZHJlc3NlcyA9IG91dHB1dC5nZXRBZGRyZXNzZXMoKS5tYXAoKGEpID0+IHtcbiAgICAgIGxldCBhZGRyOiBBZGRyZXNzID0gbmV3IEFkZHJlc3MoKVxuICAgICAgYWRkci5mcm9tQnVmZmVyKGEpXG4gICAgICByZXR1cm4gYWRkclxuICAgIH0pXG4gICAgdGhpcy5udW1hZGRycyA9IEJ1ZmZlci5hbGxvYyg0KVxuICAgIHRoaXMubnVtYWRkcnMud3JpdGVVSW50MzJCRSh0aGlzLmFkZHJlc3Nlcy5sZW5ndGgsIDApXG4gICAgdGhpcy5sb2NrdGltZSA9IGJpbnRvb2xzLmZyb21CTlRvQnVmZmVyKG91dHB1dC5nZXRMb2NrdGltZSgpLCA4KVxuICAgIHRoaXMudGhyZXNob2xkID0gQnVmZmVyLmFsbG9jKDQpXG4gICAgdGhpcy50aHJlc2hvbGQud3JpdGVVSW50MzJCRShvdXRwdXQuZ2V0VGhyZXNob2xkKCksIDApXG4gICAgdGhpcy5hbW91bnQgPSBiaW50b29scy5mcm9tQk5Ub0J1ZmZlcihvdXRwdXQuZ2V0QW1vdW50KCksIDgpXG4gICAgdGhpcy5hbW91bnRWYWx1ZSA9IG91dHB1dC5nZXRBbW91bnQoKVxuICB9XG5cbiAgZ2V0U3Rha2VhYmxlTG9ja3RpbWUoKTogQk4ge1xuICAgIHJldHVybiBiaW50b29scy5mcm9tQnVmZmVyVG9CTih0aGlzLnN0YWtlYWJsZUxvY2t0aW1lKVxuICB9XG5cbiAgZ2V0VHJhbnNmZXJhYmxlT3V0cHV0KCk6IFBhcnNlYWJsZU91dHB1dCB7XG4gICAgcmV0dXJuIHRoaXMudHJhbnNmZXJhYmxlT3V0cHV0XG4gIH1cblxuICAvKipcbiAgICogQHBhcmFtIGFzc2V0SUQgQW4gYXNzZXRJRCB3aGljaCBpcyB3cmFwcGVkIGFyb3VuZCB0aGUgQnVmZmVyIG9mIHRoZSBPdXRwdXRcbiAgICovXG4gIG1ha2VUcmFuc2ZlcmFibGUoYXNzZXRJRDogQnVmZmVyKTogVHJhbnNmZXJhYmxlT3V0cHV0IHtcbiAgICByZXR1cm4gbmV3IFRyYW5zZmVyYWJsZU91dHB1dChhc3NldElELCB0aGlzKVxuICB9XG5cbiAgc2VsZWN0KGlkOiBudW1iZXIsIC4uLmFyZ3M6IGFueVtdKTogT3V0cHV0IHtcbiAgICByZXR1cm4gU2VsZWN0T3V0cHV0Q2xhc3MoaWQsIC4uLmFyZ3MpXG4gIH1cblxuICAvKipcbiAgICogUG9wdWF0ZXMgdGhlIGluc3RhbmNlIGZyb20gYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSByZXByZXNlbnRpbmcgdGhlIFtbU3Rha2VhYmxlTG9ja091dF1dIGFuZCByZXR1cm5zIHRoZSBzaXplIG9mIHRoZSBvdXRwdXQuXG4gICAqL1xuICBmcm9tQnVmZmVyKG91dGJ1ZmY6IEJ1ZmZlciwgb2Zmc2V0OiBudW1iZXIgPSAwKTogbnVtYmVyIHtcbiAgICB0aGlzLnN0YWtlYWJsZUxvY2t0aW1lID0gYmludG9vbHMuY29weUZyb20ob3V0YnVmZiwgb2Zmc2V0LCBvZmZzZXQgKyA4KVxuICAgIG9mZnNldCArPSA4XG4gICAgdGhpcy50cmFuc2ZlcmFibGVPdXRwdXQgPSBuZXcgUGFyc2VhYmxlT3V0cHV0KClcbiAgICBvZmZzZXQgPSB0aGlzLnRyYW5zZmVyYWJsZU91dHB1dC5mcm9tQnVmZmVyKG91dGJ1ZmYsIG9mZnNldClcbiAgICB0aGlzLnN5bmNocm9uaXplKClcbiAgICByZXR1cm4gb2Zmc2V0XG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgYnVmZmVyIHJlcHJlc2VudGluZyB0aGUgW1tTdGFrZWFibGVMb2NrT3V0XV0gaW5zdGFuY2UuXG4gICAqL1xuICB0b0J1ZmZlcigpOiBCdWZmZXIge1xuICAgIGxldCB4ZmVyb3V0QnVmZjogQnVmZmVyID0gdGhpcy50cmFuc2ZlcmFibGVPdXRwdXQudG9CdWZmZXIoKVxuICAgIGNvbnN0IGJzaXplOiBudW1iZXIgPSB0aGlzLnN0YWtlYWJsZUxvY2t0aW1lLmxlbmd0aCArIHhmZXJvdXRCdWZmLmxlbmd0aFxuICAgIGNvbnN0IGJhcnI6IEJ1ZmZlcltdID0gW3RoaXMuc3Rha2VhYmxlTG9ja3RpbWUsIHhmZXJvdXRCdWZmXVxuICAgIHJldHVybiBCdWZmZXIuY29uY2F0KGJhcnIsIGJzaXplKVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIG91dHB1dElEIGZvciB0aGlzIG91dHB1dFxuICAgKi9cbiAgZ2V0T3V0cHV0SUQoKTogbnVtYmVyIHtcbiAgICByZXR1cm4gdGhpcy5fdHlwZUlEXG4gIH1cblxuICBjcmVhdGUoLi4uYXJnczogYW55W10pOiB0aGlzIHtcbiAgICByZXR1cm4gbmV3IFN0YWtlYWJsZUxvY2tPdXQoLi4uYXJncykgYXMgdGhpc1xuICB9XG5cbiAgY2xvbmUoKTogdGhpcyB7XG4gICAgY29uc3QgbmV3b3V0OiBTdGFrZWFibGVMb2NrT3V0ID0gdGhpcy5jcmVhdGUoKVxuICAgIG5ld291dC5mcm9tQnVmZmVyKHRoaXMudG9CdWZmZXIoKSlcbiAgICByZXR1cm4gbmV3b3V0IGFzIHRoaXNcbiAgfVxuXG4gIC8qKlxuICAgKiBBIFtbT3V0cHV0XV0gY2xhc3Mgd2hpY2ggc3BlY2lmaWVzIGEgW1tQYXJzZWFibGVPdXRwdXRdXSB0aGF0IGhhcyBhIGxvY2t0aW1lIHdoaWNoIGNhbiBhbHNvIGVuYWJsZSBzdGFraW5nIG9mIHRoZSB2YWx1ZSBoZWxkLCBwcmV2ZW50aW5nIHRyYW5zZmVycyBidXQgbm90IHZhbGlkYXRpb24uXG4gICAqXG4gICAqIEBwYXJhbSBhbW91bnQgQSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2luZHV0bnkvYm4uanMvfEJOfSByZXByZXNlbnRpbmcgdGhlIGFtb3VudCBpbiB0aGUgb3V0cHV0XG4gICAqIEBwYXJhbSBhZGRyZXNzZXMgQW4gYXJyYXkgb2Yge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn1zIHJlcHJlc2VudGluZyBhZGRyZXNzZXNcbiAgICogQHBhcmFtIGxvY2t0aW1lIEEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9pbmR1dG55L2JuLmpzL3xCTn0gcmVwcmVzZW50aW5nIHRoZSBsb2NrdGltZVxuICAgKiBAcGFyYW0gdGhyZXNob2xkIEEgbnVtYmVyIHJlcHJlc2VudGluZyB0aGUgdGhlIHRocmVzaG9sZCBudW1iZXIgb2Ygc2lnbmVycyByZXF1aXJlZCB0byBzaWduIHRoZSB0cmFuc2FjdGlvblxuICAgKiBAcGFyYW0gc3Rha2VhYmxlTG9ja3RpbWUgQSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2luZHV0bnkvYm4uanMvfEJOfSByZXByZXNlbnRpbmcgdGhlIHN0YWtlYWJsZSBsb2NrdGltZVxuICAgKiBAcGFyYW0gdHJhbnNmZXJhYmxlT3V0cHV0IEEgW1tQYXJzZWFibGVPdXRwdXRdXSB3aGljaCBpcyBlbWJlZGRlZCBpbnRvIHRoaXMgb3V0cHV0LlxuICAgKi9cbiAgY29uc3RydWN0b3IoXG4gICAgYW1vdW50OiBCTiA9IHVuZGVmaW5lZCxcbiAgICBhZGRyZXNzZXM6IEJ1ZmZlcltdID0gdW5kZWZpbmVkLFxuICAgIGxvY2t0aW1lOiBCTiA9IHVuZGVmaW5lZCxcbiAgICB0aHJlc2hvbGQ6IG51bWJlciA9IHVuZGVmaW5lZCxcbiAgICBzdGFrZWFibGVMb2NrdGltZTogQk4gPSB1bmRlZmluZWQsXG4gICAgdHJhbnNmZXJhYmxlT3V0cHV0OiBQYXJzZWFibGVPdXRwdXQgPSB1bmRlZmluZWRcbiAgKSB7XG4gICAgc3VwZXIoYW1vdW50LCBhZGRyZXNzZXMsIGxvY2t0aW1lLCB0aHJlc2hvbGQpXG4gICAgaWYgKHR5cGVvZiBzdGFrZWFibGVMb2NrdGltZSAhPT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgdGhpcy5zdGFrZWFibGVMb2NrdGltZSA9IGJpbnRvb2xzLmZyb21CTlRvQnVmZmVyKHN0YWtlYWJsZUxvY2t0aW1lLCA4KVxuICAgIH1cbiAgICBpZiAodHlwZW9mIHRyYW5zZmVyYWJsZU91dHB1dCAhPT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgdGhpcy50cmFuc2ZlcmFibGVPdXRwdXQgPSB0cmFuc2ZlcmFibGVPdXRwdXRcbiAgICAgIHRoaXMuc3luY2hyb25pemUoKVxuICAgIH1cbiAgfVxufVxuXG4vKipcbiAqIEFuIFtbT3V0cHV0XV0gY2xhc3Mgd2hpY2ggb25seSBzcGVjaWZpZXMgYW4gT3V0cHV0IG93bmVyc2hpcCBhbmQgdXNlcyBzZWNwMjU2azEgc2lnbmF0dXJlIHNjaGVtZS5cbiAqL1xuZXhwb3J0IGNsYXNzIFNFQ1BPd25lck91dHB1dCBleHRlbmRzIE91dHB1dCB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIlNFQ1BPd25lck91dHB1dFwiXG4gIHByb3RlY3RlZCBfdHlwZUlEID0gUGxhdGZvcm1WTUNvbnN0YW50cy5TRUNQT1dORVJPVVRQVVRJRFxuXG4gIC8vc2VyaWFsaXplIGFuZCBkZXNlcmlhbGl6ZSBib3RoIGFyZSBpbmhlcml0ZWRcblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgb3V0cHV0SUQgZm9yIHRoaXMgb3V0cHV0XG4gICAqL1xuICBnZXRPdXRwdXRJRCgpOiBudW1iZXIge1xuICAgIHJldHVybiB0aGlzLl90eXBlSURcbiAgfVxuXG4gIC8qKlxuICAgKlxuICAgKiBAcGFyYW0gYXNzZXRJRCBBbiBhc3NldElEIHdoaWNoIGlzIHdyYXBwZWQgYXJvdW5kIHRoZSBCdWZmZXIgb2YgdGhlIE91dHB1dFxuICAgKi9cbiAgbWFrZVRyYW5zZmVyYWJsZShhc3NldElEOiBCdWZmZXIpOiBUcmFuc2ZlcmFibGVPdXRwdXQge1xuICAgIHJldHVybiBuZXcgVHJhbnNmZXJhYmxlT3V0cHV0KGFzc2V0SUQsIHRoaXMpXG4gIH1cblxuICBjcmVhdGUoLi4uYXJnczogYW55W10pOiB0aGlzIHtcbiAgICByZXR1cm4gbmV3IFNFQ1BPd25lck91dHB1dCguLi5hcmdzKSBhcyB0aGlzXG4gIH1cblxuICBjbG9uZSgpOiB0aGlzIHtcbiAgICBjb25zdCBuZXdvdXQ6IFNFQ1BPd25lck91dHB1dCA9IHRoaXMuY3JlYXRlKClcbiAgICBuZXdvdXQuZnJvbUJ1ZmZlcih0aGlzLnRvQnVmZmVyKCkpXG4gICAgcmV0dXJuIG5ld291dCBhcyB0aGlzXG4gIH1cblxuICBzZWxlY3QoaWQ6IG51bWJlciwgLi4uYXJnczogYW55W10pOiBPdXRwdXQge1xuICAgIHJldHVybiBTZWxlY3RPdXRwdXRDbGFzcyhpZCwgLi4uYXJncylcbiAgfVxufVxuIl19