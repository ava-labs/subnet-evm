"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Tx = exports.UnsignedTx = exports.SelectTxClass = void 0;
/**
 * @packageDocumentation
 * @module API-PlatformVM-Transactions
 */
const buffer_1 = require("buffer/");
const bintools_1 = __importDefault(require("../../utils/bintools"));
const constants_1 = require("./constants");
const credentials_1 = require("./credentials");
const tx_1 = require("../../common/tx");
const create_hash_1 = __importDefault(require("create-hash"));
const basetx_1 = require("./basetx");
const importtx_1 = require("./importtx");
const exporttx_1 = require("./exporttx");
const validationtx_1 = require("./validationtx");
const createsubnettx_1 = require("./createsubnettx");
const errors_1 = require("../../utils/errors");
/**
 * @ignore
 */
const bintools = bintools_1.default.getInstance();
/**
 * Takes a buffer representing the output and returns the proper [[BaseTx]] instance.
 *
 * @param txtype The id of the transaction type
 *
 * @returns An instance of an [[BaseTx]]-extended class.
 */
const SelectTxClass = (txtype, ...args) => {
    if (txtype === constants_1.PlatformVMConstants.BASETX) {
        return new basetx_1.BaseTx(...args);
    }
    else if (txtype === constants_1.PlatformVMConstants.IMPORTTX) {
        return new importtx_1.ImportTx(...args);
    }
    else if (txtype === constants_1.PlatformVMConstants.EXPORTTX) {
        return new exporttx_1.ExportTx(...args);
    }
    else if (txtype === constants_1.PlatformVMConstants.ADDDELEGATORTX) {
        return new validationtx_1.AddDelegatorTx(...args);
    }
    else if (txtype === constants_1.PlatformVMConstants.ADDVALIDATORTX) {
        return new validationtx_1.AddValidatorTx(...args);
    }
    else if (txtype === constants_1.PlatformVMConstants.CREATESUBNETTX) {
        return new createsubnettx_1.CreateSubnetTx(...args);
    }
    /* istanbul ignore next */
    throw new errors_1.TransactionError("Error - SelectTxClass: unknown txtype");
};
exports.SelectTxClass = SelectTxClass;
class UnsignedTx extends tx_1.StandardUnsignedTx {
    constructor() {
        super(...arguments);
        this._typeName = "UnsignedTx";
        this._typeID = undefined;
    }
    //serialize is inherited
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.transaction = (0, exports.SelectTxClass)(fields["transaction"]["_typeID"]);
        this.transaction.deserialize(fields["transaction"], encoding);
    }
    getTransaction() {
        return this.transaction;
    }
    fromBuffer(bytes, offset = 0) {
        this.codecID = bintools.copyFrom(bytes, offset, offset + 2).readUInt16BE(0);
        offset += 2;
        const txtype = bintools
            .copyFrom(bytes, offset, offset + 4)
            .readUInt32BE(0);
        offset += 4;
        this.transaction = (0, exports.SelectTxClass)(txtype);
        return this.transaction.fromBuffer(bytes, offset);
    }
    /**
     * Signs this [[UnsignedTx]] and returns signed [[StandardTx]]
     *
     * @param kc An [[KeyChain]] used in signing
     *
     * @returns A signed [[StandardTx]]
     */
    sign(kc) {
        const txbuff = this.toBuffer();
        const msg = buffer_1.Buffer.from((0, create_hash_1.default)("sha256").update(txbuff).digest());
        const creds = this.transaction.sign(msg, kc);
        return new Tx(this, creds);
    }
}
exports.UnsignedTx = UnsignedTx;
class Tx extends tx_1.StandardTx {
    constructor() {
        super(...arguments);
        this._typeName = "Tx";
        this._typeID = undefined;
    }
    //serialize is inherited
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.unsignedTx = new UnsignedTx();
        this.unsignedTx.deserialize(fields["unsignedTx"], encoding);
        this.credentials = [];
        for (let i = 0; i < fields["credentials"].length; i++) {
            const cred = (0, credentials_1.SelectCredentialClass)(fields["credentials"][`${i}`]["_typeID"]);
            cred.deserialize(fields["credentials"][`${i}`], encoding);
            this.credentials.push(cred);
        }
    }
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} containing an [[Tx]], parses it, populates the class, and returns the length of the Tx in bytes.
     *
     * @param bytes A {@link https://github.com/feross/buffer|Buffer} containing a raw [[Tx]]
     * @param offset A number representing the starting point of the bytes to begin parsing
     *
     * @returns The length of the raw [[Tx]]
     */
    fromBuffer(bytes, offset = 0) {
        this.unsignedTx = new UnsignedTx();
        offset = this.unsignedTx.fromBuffer(bytes, offset);
        const numcreds = bintools
            .copyFrom(bytes, offset, offset + 4)
            .readUInt32BE(0);
        offset += 4;
        this.credentials = [];
        for (let i = 0; i < numcreds; i++) {
            const credid = bintools
                .copyFrom(bytes, offset, offset + 4)
                .readUInt32BE(0);
            offset += 4;
            const cred = (0, credentials_1.SelectCredentialClass)(credid);
            offset = cred.fromBuffer(bytes, offset);
            this.credentials.push(cred);
        }
        return offset;
    }
}
exports.Tx = Tx;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidHguanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi9zcmMvYXBpcy9wbGF0Zm9ybXZtL3R4LnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7OztBQUFBOzs7R0FHRztBQUNILG9DQUFnQztBQUNoQyxvRUFBMkM7QUFDM0MsMkNBQWlEO0FBQ2pELCtDQUFxRDtBQUVyRCx3Q0FBZ0U7QUFFaEUsOERBQW9DO0FBQ3BDLHFDQUFpQztBQUNqQyx5Q0FBcUM7QUFDckMseUNBQXFDO0FBRXJDLGlEQUErRDtBQUMvRCxxREFBaUQ7QUFDakQsK0NBQXFEO0FBRXJEOztHQUVHO0FBQ0gsTUFBTSxRQUFRLEdBQWEsa0JBQVEsQ0FBQyxXQUFXLEVBQUUsQ0FBQTtBQUVqRDs7Ozs7O0dBTUc7QUFDSSxNQUFNLGFBQWEsR0FBRyxDQUFDLE1BQWMsRUFBRSxHQUFHLElBQVcsRUFBVSxFQUFFO0lBQ3RFLElBQUksTUFBTSxLQUFLLCtCQUFtQixDQUFDLE1BQU0sRUFBRTtRQUN6QyxPQUFPLElBQUksZUFBTSxDQUFDLEdBQUcsSUFBSSxDQUFDLENBQUE7S0FDM0I7U0FBTSxJQUFJLE1BQU0sS0FBSywrQkFBbUIsQ0FBQyxRQUFRLEVBQUU7UUFDbEQsT0FBTyxJQUFJLG1CQUFRLENBQUMsR0FBRyxJQUFJLENBQUMsQ0FBQTtLQUM3QjtTQUFNLElBQUksTUFBTSxLQUFLLCtCQUFtQixDQUFDLFFBQVEsRUFBRTtRQUNsRCxPQUFPLElBQUksbUJBQVEsQ0FBQyxHQUFHLElBQUksQ0FBQyxDQUFBO0tBQzdCO1NBQU0sSUFBSSxNQUFNLEtBQUssK0JBQW1CLENBQUMsY0FBYyxFQUFFO1FBQ3hELE9BQU8sSUFBSSw2QkFBYyxDQUFDLEdBQUcsSUFBSSxDQUFDLENBQUE7S0FDbkM7U0FBTSxJQUFJLE1BQU0sS0FBSywrQkFBbUIsQ0FBQyxjQUFjLEVBQUU7UUFDeEQsT0FBTyxJQUFJLDZCQUFjLENBQUMsR0FBRyxJQUFJLENBQUMsQ0FBQTtLQUNuQztTQUFN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QUFJLEdBQWUsSUFBQSxtQ0FBcUIsRUFBQyxNQUFNLENBQUMsQ0FBQTtZQUN0RCxNQUFNLEdBQUcsSUFBSSxDQUFDLFVBQVUsQ0FBQyxLQUFLLEVBQUUsTUFBTSxDQUFDLENBQUE7WUFDdkMsSUFBSSxDQUFDLFdBQVcsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7U0FDNUI7UUFDRCxPQUFPLE1BQU0sQ0FBQTtJQUNmLENBQUM7Q0FDRjtBQS9DRCxnQkErQ0MiLCJzb3VyY2VzQ29udGVudCI6WyIvKipcbiAqIEBwYWNrYWdlRG9jdW1lbnRhdGlvblxuICogQG1vZHVsZSBBUEktUGxhdGZvcm1WTS1UcmFuc2FjdGlvbnNcbiAqL1xuaW1wb3J0IHsgQnVmZmVyIH0gZnJvbSBcImJ1ZmZlci9cIlxuaW1wb3J0IEJpblRvb2xzIGZyb20gXCIuLi8uLi91dGlscy9iaW50b29sc1wiXG5pbXBvcnQgeyBQbGF0Zm9ybVZNQ29uc3RhbnRzIH0gZnJvbSBcIi4vY29uc3RhbnRzXCJcbmltcG9ydCB7IFNlbGVjdENyZWRlbnRpYWxDbGFzcyB9IGZyb20gXCIuL2NyZWRlbnRpYWxzXCJcbmltcG9ydCB7IEtleUNoYWluLCBLZXlQYWlyIH0gZnJvbSBcIi4va2V5Y2hhaW5cIlxuaW1wb3J0IHsgU3RhbmRhcmRUeCwgU3RhbmRhcmRVbnNpZ25lZFR4IH0gZnJvbSBcIi4uLy4uL2NvbW1vbi90eFwiXG5pbXBvcnQgeyBDcmVkZW50aWFsIH0gZnJvbSBcIi4uLy4uL2NvbW1vbi9jcmVkZW50aWFsc1wiXG5pbXBvcnQgY3JlYXRlSGFzaCBmcm9tIFwiY3JlYXRlLWhhc2hcIlxuaW1wb3J0IHsgQmFzZVR4IH0gZnJvbSBcIi4vYmFzZXR4XCJcbmltcG9ydCB7IEltcG9ydFR4IH0gZnJvbSBcIi4vaW1wb3J0dHhcIlxuaW1wb3J0IHsgRXhwb3J0VHggfSBmcm9tIFwiLi9leHBvcnR0eFwiXG5pbXBvcnQgeyBTZXJpYWxpemVkRW5jb2RpbmcgfSBmcm9tIFwiLi4vLi4vdXRpbHMvc2VyaWFsaXphdGlvblwiXG5pbXBvcnQgeyBBZGREZWxlZ2F0b3JUeCwgQWRkVmFsaWRhdG9yVHggfSBmcm9tIFwiLi92YWxpZGF0aW9udHhcIlxuaW1wb3J0IHsgQ3JlYXRlU3VibmV0VHggfSBmcm9tIFwiLi9jcmVhdGVzdWJuZXR0eFwiXG5pbXBvcnQgeyBUcmFuc2FjdGlvbkVycm9yIH0gZnJvbSBcIi4uLy4uL3V0aWxzL2Vycm9yc1wiXG5cbi8qKlxuICogQGlnbm9yZVxuICovXG5jb25zdCBiaW50b29sczogQmluVG9vbHMgPSBCaW5Ub29scy5nZXRJbnN0YW5jZSgpXG5cbi8qKlxuICogVGFrZXMgYSBidWZmZXIgcmVwcmVzZW50aW5nIHRoZSBvdXRwdXQgYW5kIHJldHVybnMgdGhlIHByb3BlciBbW0Jhc2VUeF1dIGluc3RhbmNlLlxuICpcbiAqIEBwYXJhbSB0eHR5cGUgVGhlIGlkIG9mIHRoZSB0cmFuc2FjdGlvbiB0eXBlXG4gKlxuICogQHJldHVybnMgQW4gaW5zdGFuY2Ugb2YgYW4gW1tCYXNlVHhdXS1leHRlbmRlZCBjbGFzcy5cbiAqL1xuZXhwb3J0IGNvbnN0IFNlbGVjdFR4Q2xhc3MgPSAodHh0eXBlOiBudW1iZXIsIC4uLmFyZ3M6IGFueVtdKTogQmFzZVR4ID0+IHtcbiAgaWYgKHR4dHlwZSA9PT0gUGxhdGZvcm1WTUNvbnN0YW50cy5CQVNFVFgpIHtcbiAgICByZXR1cm4gbmV3IEJhc2VUeCguLi5hcmdzKVxuICB9IGVsc2UgaWYgKHR4dHlwZSA9PT0gUGxhdGZvcm1WTUNvbnN0YW50cy5JTVBPUlRUWCkge1xuICAgIHJldHVybiBuZXcgSW1wb3J0VHgoLi4uYXJncylcbiAgfSBlbHNlIGlmICh0eHR5cGUgPT09IFBsYXRmb3JtVk1Db25zdGFudHMuRVhQT1JUVFgpIHtcbiAgICByZXR1cm4gbmV3IEV4cG9ydFR4KC4uLmFyZ3MpXG4gIH0gZWxzZSBpZiAodHh0eXBlID09PSBQbGF0Zm9ybVZNQ29uc3RhbnRzLkFERERFTEVHQVRPUlRYKSB7XG4gICAgcmV0dXJuIG5ldyBBZGREZWxlZ2F0b3JUeCguLi5hcmdzKVxuICB9IGVsc2UgaWYgKHR4dHlwZSA9PT0gUGxhdGZvcm1WTUNvbnN0YW50cy5BRERWQUxJREFUT1JUWCkge1xuICAgIHJldHVybiBuZXcgQWRkVmFsaWRhdG9yVHgoLi4uYXJncylcbiAgfSBlbHNlIGlmICh0eHR5cGUgPT09IFBsYXRmb3JtVk1Db25zdGFudHMuQ1JFQVRFU1VCTkVUVFgpIHtcbiAgICByZXR1cm4gbmV3IENyZWF0ZVN1Ym5ldFR4KC4uLmFyZ3MpXG4gIH1cbiAgLyogaXN0YW5idWwgaWdub3JlIG5leHQgKi9cbiAgdGhyb3cgbmV3IFRyYW5zYWN0aW9uRXJyb3IoXCJFcnJvciAtIFNlbGVjdFR4Q2xhc3M6IHVua25vd24gdHh0eXBlXCIpXG59XG5cbmV4cG9ydCBjbGFzcyBVbnNpZ25lZFR4IGV4dGVuZHMgU3RhbmRhcmRVbnNpZ25lZFR4PEtleVBhaXIsIEtleUNoYWluLCBCYXNlVHg+IHtcbiAgcHJvdGVjdGVkIF90eXBlTmFtZSA9IFwiVW5zaWduZWRUeFwiXG4gIHByb3RlY3RlZCBfdHlwZUlEID0gdW5kZWZpbmVkXG5cbiAgLy9zZXJpYWxpemUgaXMgaW5oZXJpdGVkXG5cbiAgZGVzZXJpYWxpemUoZmllbGRzOiBvYmplY3QsIGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKSB7XG4gICAgc3VwZXIuZGVzZXJpYWxpemUoZmllbGRzLCBlbmNvZGluZylcbiAgICB0aGlzLnRyYW5zYWN0aW9uID0gU2VsZWN0VHhDbGFzcyhmaWVsZHNbXCJ0cmFuc2FjdGlvblwiXVtcIl90eXBlSURcIl0pXG4gICAgdGhpcy50cmFuc2FjdGlvbi5kZXNlcmlhbGl6ZShmaWVsZHNbXCJ0cmFuc2FjdGlvblwiXSwgZW5jb2RpbmcpXG4gIH1cblxuICBnZXRUcmFuc2FjdGlvbigpOiBCYXNlVHgge1xuICAgIHJldHVybiB0aGlzLnRyYW5zYWN0aW9uIGFzIEJhc2VUeFxuICB9XG5cbiAgZnJvbUJ1ZmZlcihieXRlczogQnVmZmVyLCBvZmZzZXQ6IG51bWJlciA9IDApOiBudW1iZXIge1xuICAgIHRoaXMuY29kZWNJRCA9IGJpbnRvb2xzLmNvcHlGcm9tKGJ5dGVzLCBvZmZzZXQsIG9mZnNldCArIDIpLnJlYWRVSW50MTZCRSgwKVxuICAgIG9mZnNldCArPSAyXG4gICAgY29uc3QgdHh0eXBlOiBudW1iZXIgPSBiaW50b29sc1xuICAgICAgLmNvcHlGcm9tKGJ5dGVzLCBvZmZzZXQsIG9mZnNldCArIDQpXG4gICAgICAucmVhZFVJbnQzMkJFKDApXG4gICAgb2Zmc2V0ICs9IDRcbiAgICB0aGlzLnRyYW5zYWN0aW9uID0gU2VsZWN0VHhDbGFzcyh0eHR5cGUpXG4gICAgcmV0dXJuIHRoaXMudHJhbnNhY3Rpb24uZnJvbUJ1ZmZlcihieXRlcywgb2Zmc2V0KVxuICB9XG5cbiAgLyoqXG4gICAqIFNpZ25zIHRoaXMgW1tVbnNpZ25lZFR4XV0gYW5kIHJldHVybnMgc2lnbmVkIFtbU3RhbmRhcmRUeF1dXG4gICAqXG4gICAqIEBwYXJhbSBrYyBBbiBbW0tleUNoYWluXV0gdXNlZCBpbiBzaWduaW5nXG4gICAqXG4gICAqIEByZXR1cm5zIEEgc2lnbmVkIFtbU3RhbmRhcmRUeF1dXG4gICAqL1xuICBzaWduKGtjOiBLZXlDaGFpbik6IFR4IHtcbiAgICBjb25zdCB0eGJ1ZmYgPSB0aGlzLnRvQnVmZmVyKClcbiAgICBjb25zdCBtc2c6IEJ1ZmZlciA9IEJ1ZmZlci5mcm9tKFxuICAgICAgY3JlYXRlSGFzaChcInNoYTI1NlwiKS51cGRhdGUodHhidWZmKS5kaWdlc3QoKVxuICAgIClcbiAgICBjb25zdCBjcmVkczogQ3JlZGVudGlhbFtdID0gdGhpcy50cmFuc2FjdGlvbi5zaWduKG1zZywga2MpXG4gICAgcmV0dXJuIG5ldyBUeCh0aGlzLCBjcmVkcylcbiAgfVxufVxuXG5leHBvcnQgY2xhc3MgVHggZXh0ZW5kcyBTdGFuZGFyZFR4PEtleVBhaXIsIEtleUNoYWluLCBVbnNpZ25lZFR4PiB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIlR4XCJcbiAgcHJvdGVjdGVkIF90eXBlSUQgPSB1bmRlZmluZWRcblxuICAvL3NlcmlhbGl6ZSBpcyBpbmhlcml0ZWRcblxuICBkZXNlcmlhbGl6ZShmaWVsZHM6IG9iamVjdCwgZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpIHtcbiAgICBzdXBlci5kZXNlcmlhbGl6ZShmaWVsZHMsIGVuY29kaW5nKVxuICAgIHRoaXMudW5zaWduZWRUeCA9IG5ldyBVbnNpZ25lZFR4KClcbiAgICB0aGlzLnVuc2lnbmVkVHguZGVzZXJpYWxpemUoZmllbGRzW1widW5zaWduZWRUeFwiXSwgZW5jb2RpbmcpXG4gICAgdGhpcy5jcmVkZW50aWFscyA9IFtdXG4gICAgZm9yIChsZXQgaTogbnVtYmVyID0gMDsgaSA8IGZpZWxkc1tcImNyZWRlbnRpYWxzXCJdLmxlbmd0aDsgaSsrKSB7XG4gICAgICBjb25zdCBjcmVkOiBDcmVkZW50aWFsID0gU2VsZWN0Q3JlZGVudGlhbENsYXNzKFxuICAgICAgICBmaWVsZHNbXCJjcmVkZW50aWFsc1wiXVtgJHtpfWBdW1wiX3R5cGVJRFwiXVxuICAgICAgKVxuICAgICAgY3JlZC5kZXNlcmlhbGl6ZShmaWVsZHNbXCJjcmVkZW50aWFsc1wiXVtgJHtpfWBdLCBlbmNvZGluZylcbiAgICAgIHRoaXMuY3JlZGVudGlhbHMucHVzaChjcmVkKVxuICAgIH1cbiAgfVxuXG4gIC8qKlxuICAgKiBUYWtlcyBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IGNvbnRhaW5pbmcgYW4gW1tUeF1dLCBwYXJzZXMgaXQsIHBvcHVsYXRlcyB0aGUgY2xhc3MsIGFuZCByZXR1cm5zIHRoZSBsZW5ndGggb2YgdGhlIFR4IGluIGJ5dGVzLlxuICAgKlxuICAgKiBAcGFyYW0gYnl0ZXMgQSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSBjb250YWluaW5nIGEgcmF3IFtbVHhdXVxuICAgKiBAcGFyYW0gb2Zmc2V0IEEgbnVtYmVyIHJlcHJlc2VudGluZyB0aGUgc3RhcnRpbmcgcG9pbnQgb2YgdGhlIGJ5dGVzIHRvIGJlZ2luIHBhcnNpbmdcbiAgICpcbiAgICogQHJldHVybnMgVGhlIGxlbmd0aCBvZiB0aGUgcmF3IFtbVHhdXVxuICAgKi9cbiAgZnJvbUJ1ZmZlcihieXRlczogQnVmZmVyLCBvZmZzZXQ6IG51bWJlciA9IDApOiBudW1iZXIge1xuICAgIHRoaXMudW5zaWduZWRUeCA9IG5ldyBVbnNpZ25lZFR4KClcbiAgICBvZmZzZXQgPSB0aGlzLnVuc2lnbmVkVHguZnJvbUJ1ZmZlcihieXRlcywgb2Zmc2V0KVxuICAgIGNvbnN0IG51bWNyZWRzOiBudW1iZXIgPSBiaW50b29sc1xuICAgICAgLmNvcHlGcm9tKGJ5dGVzLCBvZmZzZXQsIG9mZnNldCArIDQpXG4gICAgICAucmVhZFVJbnQzMkJFKDApXG4gICAgb2Zmc2V0ICs9IDRcbiAgICB0aGlzLmNyZWRlbnRpYWxzID0gW11cbiAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgbnVtY3JlZHM7IGkrKykge1xuICAgICAgY29uc3QgY3JlZGlkOiBudW1iZXIgPSBiaW50b29sc1xuICAgICAgICAuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgNClcbiAgICAgICAgLnJlYWRVSW50MzJCRSgwKVxuICAgICAgb2Zmc2V0ICs9IDRcbiAgICAgIGNvbnN0IGNyZWQ6IENyZWRlbnRpYWwgPSBTZWxlY3RDcmVkZW50aWFsQ2xhc3MoY3JlZGlkKVxuICAgICAgb2Zmc2V0ID0gY3JlZC5mcm9tQnVmZmVyKGJ5dGVzLCBvZmZzZXQpXG4gICAgICB0aGlzLmNyZWRlbnRpYWxzLnB1c2goY3JlZClcbiAgICB9XG4gICAgcmV0dXJuIG9mZnNldFxuICB9XG59XG4iXX0=