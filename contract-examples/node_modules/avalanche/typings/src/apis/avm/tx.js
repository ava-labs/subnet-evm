"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Tx = exports.UnsignedTx = exports.SelectTxClass = void 0;
/**
 * @packageDocumentation
 * @module API-AVM-Transactions
 */
const buffer_1 = require("buffer/");
const bintools_1 = __importDefault(require("../../utils/bintools"));
const constants_1 = require("./constants");
const credentials_1 = require("./credentials");
const tx_1 = require("../../common/tx");
const create_hash_1 = __importDefault(require("create-hash"));
const basetx_1 = require("./basetx");
const createassettx_1 = require("./createassettx");
const operationtx_1 = require("./operationtx");
const importtx_1 = require("./importtx");
const exporttx_1 = require("./exporttx");
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
    if (txtype === constants_1.AVMConstants.BASETX) {
        return new basetx_1.BaseTx(...args);
    }
    else if (txtype === constants_1.AVMConstants.CREATEASSETTX) {
        return new createassettx_1.CreateAssetTx(...args);
    }
    else if (txtype === constants_1.AVMConstants.OPERATIONTX) {
        return new operationtx_1.OperationTx(...args);
    }
    else if (txtype === constants_1.AVMConstants.IMPORTTX) {
        return new importtx_1.ImportTx(...args);
    }
    else if (txtype === constants_1.AVMConstants.EXPORTTX) {
        return new exporttx_1.ExportTx(...args);
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
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidHguanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi9zcmMvYXBpcy9hdm0vdHgudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7O0FBQUE7OztHQUdHO0FBQ0gsb0NBQWdDO0FBQ2hDLG9FQUEyQztBQUMzQywyQ0FBMEM7QUFDMUMsK0NBQXFEO0FBR3JELHdDQUFnRTtBQUNoRSw4REFBb0M7QUFDcEMscUNBQWlDO0FBQ2pDLG1EQUErQztBQUMvQywrQ0FBMkM7QUFDM0MseUNBQXFDO0FBQ3JDLHlDQUFxQztBQUVyQywrQ0FBcUQ7QUFFckQ7O0dBRUc7QUFDSCxNQUFNLFFBQVEsR0FBYSxrQkFBUSxDQUFDLFdBQVcsRUFBRSxDQUFBO0FBRWpEOzs7Ozs7R0FNRztBQUNJLE1BQU0sYUFBYSxHQUFHLENBQUMsTUFBYyxFQUFFLEdBQUcsSUFBVyxFQUFVLEVBQUU7SUFDdEUsSUFBSSxNQUFNLEtBQUssd0JBQVksQ0FBQyxNQUFNLEVBQUU7UUFDbEMsT0FBTyxJQUFJLGVBQU0sQ0FBQyxHQUFHLElBQUksQ0FBQyxDQUFBO0tBQzNCO1NBQU0sSUFBSSxNQUFNLEtBQUssd0JBQVksQ0FBQyxhQUFhLEVBQUU7UUFDaEQsT0FBTyxJQUFJLDZCQUFhLENBQUMsR0FBRyxJQUFJLENBQUMsQ0FBQTtLQUNsQztTQUFNLElBQUksTUFBTSxLQUFLLHdCQUFZLENBQUMsV0FBVyxFQUFFO1FBQzlDLE9BQU8sSUFBSSx5QkFBVyxDQUFDLEdBQUcsSUFBSSxDQUFDLENBQUE7S0FDaEM7U0FBTSxJQUFJLE1BQU0sS0FBSyx3QkFBWSxDQUFDLFFBQVEsRUFBRTtRQUMzQyxPQUFPLElBQUksbUJBQVEsQ0FBQyxHQUFHLElBQUksQ0FBQyxDQUFBO0tBQzdCO1NBQU0sSUFBSSxNQUFNLEtBQUssd0JBQVksQ0FBQyxRQUFRLEVBQUU7UUFDM0MsT0FBTyxJQUFJLG1CQUFRLENBQUMsR0FBRyxJQUFJLENBQUMsQ0FBQTtLQUM3QjtJQUNELDBCQUEwQjtJQUMx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QUFJLENBQUMsV0FBVyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtTQUM1QjtRQUNELE9BQU8sTUFBTSxDQUFBO0lBQ2YsQ0FBQztDQUNGO0FBL0NELGdCQStDQyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogQHBhY2thZ2VEb2N1bWVudGF0aW9uXG4gKiBAbW9kdWxlIEFQSS1BVk0tVHJhbnNhY3Rpb25zXG4gKi9cbmltcG9ydCB7IEJ1ZmZlciB9IGZyb20gXCJidWZmZXIvXCJcbmltcG9ydCBCaW5Ub29scyBmcm9tIFwiLi4vLi4vdXRpbHMvYmludG9vbHNcIlxuaW1wb3J0IHsgQVZNQ29uc3RhbnRzIH0gZnJvbSBcIi4vY29uc3RhbnRzXCJcbmltcG9ydCB7IFNlbGVjdENyZWRlbnRpYWxDbGFzcyB9IGZyb20gXCIuL2NyZWRlbnRpYWxzXCJcbmltcG9ydCB7IEtleUNoYWluLCBLZXlQYWlyIH0gZnJvbSBcIi4va2V5Y2hhaW5cIlxuaW1wb3J0IHsgQ3JlZGVudGlhbCB9IGZyb20gXCIuLi8uLi9jb21tb24vY3JlZGVudGlhbHNcIlxuaW1wb3J0IHsgU3RhbmRhcmRUeCwgU3RhbmRhcmRVbnNpZ25lZFR4IH0gZnJvbSBcIi4uLy4uL2NvbW1vbi90eFwiXG5pbXBvcnQgY3JlYXRlSGFzaCBmcm9tIFwiY3JlYXRlLWhhc2hcIlxuaW1wb3J0IHsgQmFzZVR4IH0gZnJvbSBcIi4vYmFzZXR4XCJcbmltcG9ydCB7IENyZWF0ZUFzc2V0VHggfSBmcm9tIFwiLi9jcmVhdGVhc3NldHR4XCJcbmltcG9ydCB7IE9wZXJhdGlvblR4IH0gZnJvbSBcIi4vb3BlcmF0aW9udHhcIlxuaW1wb3J0IHsgSW1wb3J0VHggfSBmcm9tIFwiLi9pbXBvcnR0eFwiXG5pbXBvcnQgeyBFeHBvcnRUeCB9IGZyb20gXCIuL2V4cG9ydHR4XCJcbmltcG9ydCB7IFNlcmlhbGl6ZWRFbmNvZGluZyB9IGZyb20gXCIuLi8uLi91dGlscy9zZXJpYWxpemF0aW9uXCJcbmltcG9ydCB7IFRyYW5zYWN0aW9uRXJyb3IgfSBmcm9tIFwiLi4vLi4vdXRpbHMvZXJyb3JzXCJcblxuLyoqXG4gKiBAaWdub3JlXG4gKi9cbmNvbnN0IGJpbnRvb2xzOiBCaW5Ub29scyA9IEJpblRvb2xzLmdldEluc3RhbmNlKClcblxuLyoqXG4gKiBUYWtlcyBhIGJ1ZmZlciByZXByZXNlbnRpbmcgdGhlIG91dHB1dCBhbmQgcmV0dXJucyB0aGUgcHJvcGVyIFtbQmFzZVR4XV0gaW5zdGFuY2UuXG4gKlxuICogQHBhcmFtIHR4dHlwZSBUaGUgaWQgb2YgdGhlIHRyYW5zYWN0aW9uIHR5cGVcbiAqXG4gKiBAcmV0dXJucyBBbiBpbnN0YW5jZSBvZiBhbiBbW0Jhc2VUeF1dLWV4dGVuZGVkIGNsYXNzLlxuICovXG5leHBvcnQgY29uc3QgU2VsZWN0VHhDbGFzcyA9ICh0eHR5cGU6IG51bWJlciwgLi4uYXJnczogYW55W10pOiBCYXNlVHggPT4ge1xuICBpZiAodHh0eXBlID09PSBBVk1Db25zdGFudHMuQkFTRVRYKSB7XG4gICAgcmV0dXJuIG5ldyBCYXNlVHgoLi4uYXJncylcbiAgfSBlbHNlIGlmICh0eHR5cGUgPT09IEFWTUNvbnN0YW50cy5DUkVBVEVBU1NFVFRYKSB7XG4gICAgcmV0dXJuIG5ldyBDcmVhdGVBc3NldFR4KC4uLmFyZ3MpXG4gIH0gZWxzZSBpZiAodHh0eXBlID09PSBBVk1Db25zdGFudHMuT1BFUkFUSU9OVFgpIHtcbiAgICByZXR1cm4gbmV3IE9wZXJhdGlvblR4KC4uLmFyZ3MpXG4gIH0gZWxzZSBpZiAodHh0eXBlID09PSBBVk1Db25zdGFudHMuSU1QT1JUVFgpIHtcbiAgICByZXR1cm4gbmV3IEltcG9ydFR4KC4uLmFyZ3MpXG4gIH0gZWxzZSBpZiAodHh0eXBlID09PSBBVk1Db25zdGFudHMuRVhQT1JUVFgpIHtcbiAgICByZXR1cm4gbmV3IEV4cG9ydFR4KC4uLmFyZ3MpXG4gIH1cbiAgLyogaXN0YW5idWwgaWdub3JlIG5leHQgKi9cbiAgdGhyb3cgbmV3IFRyYW5zYWN0aW9uRXJyb3IoXCJFcnJvciAtIFNlbGVjdFR4Q2xhc3M6IHVua25vd24gdHh0eXBlXCIpXG59XG5cbmV4cG9ydCBjbGFzcyBVbnNpZ25lZFR4IGV4dGVuZHMgU3RhbmRhcmRVbnNpZ25lZFR4PEtleVBhaXIsIEtleUNoYWluLCBCYXNlVHg+IHtcbiAgcHJvdGVjdGVkIF90eXBlTmFtZSA9IFwiVW5zaWduZWRUeFwiXG4gIHByb3RlY3RlZCBfdHlwZUlEID0gdW5kZWZpbmVkXG5cbiAgLy9zZXJpYWxpemUgaXMgaW5oZXJpdGVkXG5cbiAgZGVzZXJpYWxpemUoZmllbGRzOiBvYmplY3QsIGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKSB7XG4gICAgc3VwZXIuZGVzZXJpYWxpemUoZmllbGRzLCBlbmNvZGluZylcbiAgICB0aGlzLnRyYW5zYWN0aW9uID0gU2VsZWN0VHhDbGFzcyhmaWVsZHNbXCJ0cmFuc2FjdGlvblwiXVtcIl90eXBlSURcIl0pXG4gICAgdGhpcy50cmFuc2FjdGlvbi5kZXNlcmlhbGl6ZShmaWVsZHNbXCJ0cmFuc2FjdGlvblwiXSwgZW5jb2RpbmcpXG4gIH1cblxuICBnZXRUcmFuc2FjdGlvbigpOiBCYXNlVHgge1xuICAgIHJldHVybiB0aGlzLnRyYW5zYWN0aW9uIGFzIEJhc2VUeFxuICB9XG5cbiAgZnJvbUJ1ZmZlcihieXRlczogQnVmZmVyLCBvZmZzZXQ6IG51bWJlciA9IDApOiBudW1iZXIge1xuICAgIHRoaXMuY29kZWNJRCA9IGJpbnRvb2xzLmNvcHlGcm9tKGJ5dGVzLCBvZmZzZXQsIG9mZnNldCArIDIpLnJlYWRVSW50MTZCRSgwKVxuICAgIG9mZnNldCArPSAyXG4gICAgY29uc3QgdHh0eXBlOiBudW1iZXIgPSBiaW50b29sc1xuICAgICAgLmNvcHlGcm9tKGJ5dGVzLCBvZmZzZXQsIG9mZnNldCArIDQpXG4gICAgICAucmVhZFVJbnQzMkJFKDApXG4gICAgb2Zmc2V0ICs9IDRcbiAgICB0aGlzLnRyYW5zYWN0aW9uID0gU2VsZWN0VHhDbGFzcyh0eHR5cGUpXG4gICAgcmV0dXJuIHRoaXMudHJhbnNhY3Rpb24uZnJvbUJ1ZmZlcihieXRlcywgb2Zmc2V0KVxuICB9XG5cbiAgLyoqXG4gICAqIFNpZ25zIHRoaXMgW1tVbnNpZ25lZFR4XV0gYW5kIHJldHVybnMgc2lnbmVkIFtbU3RhbmRhcmRUeF1dXG4gICAqXG4gICAqIEBwYXJhbSBrYyBBbiBbW0tleUNoYWluXV0gdXNlZCBpbiBzaWduaW5nXG4gICAqXG4gICAqIEByZXR1cm5zIEEgc2lnbmVkIFtbU3RhbmRhcmRUeF1dXG4gICAqL1xuICBzaWduKGtjOiBLZXlDaGFpbik6IFR4IHtcbiAgICBjb25zdCB0eGJ1ZmYgPSB0aGlzLnRvQnVmZmVyKClcbiAgICBjb25zdCBtc2c6IEJ1ZmZlciA9IEJ1ZmZlci5mcm9tKFxuICAgICAgY3JlYXRlSGFzaChcInNoYTI1NlwiKS51cGRhdGUodHhidWZmKS5kaWdlc3QoKVxuICAgIClcbiAgICBjb25zdCBjcmVkczogQ3JlZGVudGlhbFtdID0gdGhpcy50cmFuc2FjdGlvbi5zaWduKG1zZywga2MpXG4gICAgcmV0dXJuIG5ldyBUeCh0aGlzLCBjcmVkcylcbiAgfVxufVxuXG5leHBvcnQgY2xhc3MgVHggZXh0ZW5kcyBTdGFuZGFyZFR4PEtleVBhaXIsIEtleUNoYWluLCBVbnNpZ25lZFR4PiB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIlR4XCJcbiAgcHJvdGVjdGVkIF90eXBlSUQgPSB1bmRlZmluZWRcblxuICAvL3NlcmlhbGl6ZSBpcyBpbmhlcml0ZWRcblxuICBkZXNlcmlhbGl6ZShmaWVsZHM6IG9iamVjdCwgZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpIHtcbiAgICBzdXBlci5kZXNlcmlhbGl6ZShmaWVsZHMsIGVuY29kaW5nKVxuICAgIHRoaXMudW5zaWduZWRUeCA9IG5ldyBVbnNpZ25lZFR4KClcbiAgICB0aGlzLnVuc2lnbmVkVHguZGVzZXJpYWxpemUoZmllbGRzW1widW5zaWduZWRUeFwiXSwgZW5jb2RpbmcpXG4gICAgdGhpcy5jcmVkZW50aWFscyA9IFtdXG4gICAgZm9yIChsZXQgaTogbnVtYmVyID0gMDsgaSA8IGZpZWxkc1tcImNyZWRlbnRpYWxzXCJdLmxlbmd0aDsgaSsrKSB7XG4gICAgICBjb25zdCBjcmVkOiBDcmVkZW50aWFsID0gU2VsZWN0Q3JlZGVudGlhbENsYXNzKFxuICAgICAgICBmaWVsZHNbXCJjcmVkZW50aWFsc1wiXVtgJHtpfWBdW1wiX3R5cGVJRFwiXVxuICAgICAgKVxuICAgICAgY3JlZC5kZXNlcmlhbGl6ZShmaWVsZHNbXCJjcmVkZW50aWFsc1wiXVtgJHtpfWBdLCBlbmNvZGluZylcbiAgICAgIHRoaXMuY3JlZGVudGlhbHMucHVzaChjcmVkKVxuICAgIH1cbiAgfVxuXG4gIC8qKlxuICAgKiBUYWtlcyBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IGNvbnRhaW5pbmcgYW4gW1tUeF1dLCBwYXJzZXMgaXQsIHBvcHVsYXRlcyB0aGUgY2xhc3MsIGFuZCByZXR1cm5zIHRoZSBsZW5ndGggb2YgdGhlIFR4IGluIGJ5dGVzLlxuICAgKlxuICAgKiBAcGFyYW0gYnl0ZXMgQSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSBjb250YWluaW5nIGEgcmF3IFtbVHhdXVxuICAgKiBAcGFyYW0gb2Zmc2V0IEEgbnVtYmVyIHJlcHJlc2VudGluZyB0aGUgc3RhcnRpbmcgcG9pbnQgb2YgdGhlIGJ5dGVzIHRvIGJlZ2luIHBhcnNpbmdcbiAgICpcbiAgICogQHJldHVybnMgVGhlIGxlbmd0aCBvZiB0aGUgcmF3IFtbVHhdXVxuICAgKi9cbiAgZnJvbUJ1ZmZlcihieXRlczogQnVmZmVyLCBvZmZzZXQ6IG51bWJlciA9IDApOiBudW1iZXIge1xuICAgIHRoaXMudW5zaWduZWRUeCA9IG5ldyBVbnNpZ25lZFR4KClcbiAgICBvZmZzZXQgPSB0aGlzLnVuc2lnbmVkVHguZnJvbUJ1ZmZlcihieXRlcywgb2Zmc2V0KVxuICAgIGNvbnN0IG51bWNyZWRzOiBudW1iZXIgPSBiaW50b29sc1xuICAgICAgLmNvcHlGcm9tKGJ5dGVzLCBvZmZzZXQsIG9mZnNldCArIDQpXG4gICAgICAucmVhZFVJbnQzMkJFKDApXG4gICAgb2Zmc2V0ICs9IDRcbiAgICB0aGlzLmNyZWRlbnRpYWxzID0gW11cbiAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgbnVtY3JlZHM7IGkrKykge1xuICAgICAgY29uc3QgY3JlZGlkOiBudW1iZXIgPSBiaW50b29sc1xuICAgICAgICAuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgNClcbiAgICAgICAgLnJlYWRVSW50MzJCRSgwKVxuICAgICAgb2Zmc2V0ICs9IDRcbiAgICAgIGNvbnN0IGNyZWQ6IENyZWRlbnRpYWwgPSBTZWxlY3RDcmVkZW50aWFsQ2xhc3MoY3JlZGlkKVxuICAgICAgb2Zmc2V0ID0gY3JlZC5mcm9tQnVmZmVyKGJ5dGVzLCBvZmZzZXQpXG4gICAgICB0aGlzLmNyZWRlbnRpYWxzLnB1c2goY3JlZClcbiAgICB9XG4gICAgcmV0dXJuIG9mZnNldFxuICB9XG59XG4iXX0=