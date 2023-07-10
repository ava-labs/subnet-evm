"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.StandardAmountInput = exports.StandardTransferableInput = exports.StandardParseableInput = exports.Input = void 0;
/**
 * @packageDocumentation
 * @module Common-Inputs
 */
const buffer_1 = require("buffer/");
const bintools_1 = __importDefault(require("../utils/bintools"));
const bn_js_1 = __importDefault(require("bn.js"));
const credentials_1 = require("./credentials");
const serialization_1 = require("../utils/serialization");
/**
 * @ignore
 */
const bintools = bintools_1.default.getInstance();
const serialization = serialization_1.Serialization.getInstance();
class Input extends serialization_1.Serializable {
    constructor() {
        super(...arguments);
        this._typeName = "Input";
        this._typeID = undefined;
        this.sigCount = buffer_1.Buffer.alloc(4);
        this.sigIdxs = []; // idxs of signers from utxo
        /**
         * Returns the array of [[SigIdx]] for this [[Input]]
         */
        this.getSigIdxs = () => this.sigIdxs;
        /**
         * Creates and adds a [[SigIdx]] to the [[Input]].
         *
         * @param addressIdx The index of the address to reference in the signatures
         * @param address The address of the source of the signature
         */
        this.addSignatureIdx = (addressIdx, address) => {
            const sigidx = new credentials_1.SigIdx();
            const b = buffer_1.Buffer.alloc(4);
            b.writeUInt32BE(addressIdx, 0);
            sigidx.fromBuffer(b);
            sigidx.setSource(address);
            this.sigIdxs.push(sigidx);
            this.sigCount.writeUInt32BE(this.sigIdxs.length, 0);
        };
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { sigIdxs: this.sigIdxs.map((s) => s.serialize(encoding)) });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.sigIdxs = fields["sigIdxs"].map((s) => {
            let sidx = new credentials_1.SigIdx();
            sidx.deserialize(s, encoding);
            return sidx;
        });
        this.sigCount.writeUInt32BE(this.sigIdxs.length, 0);
    }
    fromBuffer(bytes, offset = 0) {
        this.sigCount = bintools.copyFrom(bytes, offset, offset + 4);
        offset += 4;
        const sigCount = this.sigCount.readUInt32BE(0);
        this.sigIdxs = [];
        for (let i = 0; i < sigCount; i++) {
            const sigidx = new credentials_1.SigIdx();
            const sigbuff = bintools.copyFrom(bytes, offset, offset + 4);
            sigidx.fromBuffer(sigbuff);
            offset += 4;
            this.sigIdxs.push(sigidx);
        }
        return offset;
    }
    toBuffer() {
        this.sigCount.writeUInt32BE(this.sigIdxs.length, 0);
        let bsize = this.sigCount.length;
        const barr = [this.sigCount];
        for (let i = 0; i < this.sigIdxs.length; i++) {
            const b = this.sigIdxs[`${i}`].toBuffer();
            barr.push(b);
            bsize += b.length;
        }
        return buffer_1.Buffer.concat(barr, bsize);
    }
    /**
     * Returns a base-58 representation of the [[Input]].
     */
    toString() {
        return bintools.bufferToB58(this.toBuffer());
    }
}
exports.Input = Input;
Input.comparator = () => (a, b) => {
    const aoutid = buffer_1.Buffer.alloc(4);
    aoutid.writeUInt32BE(a.getInputID(), 0);
    const abuff = a.toBuffer();
    const boutid = buffer_1.Buffer.alloc(4);
    boutid.writeUInt32BE(b.getInputID(), 0);
    const bbuff = b.toBuffer();
    const asort = buffer_1.Buffer.concat([aoutid, abuff], aoutid.length + abuff.length);
    const bsort = buffer_1.Buffer.concat([boutid, bbuff], boutid.length + bbuff.length);
    return buffer_1.Buffer.compare(asort, bsort);
};
class StandardParseableInput extends serialization_1.Serializable {
    /**
     * Class representing an [[StandardParseableInput]] for a transaction.
     *
     * @param input A number representing the InputID of the [[StandardParseableInput]]
     */
    constructor(input = undefined) {
        super();
        this._typeName = "StandardParseableInput";
        this._typeID = undefined;
        this.getInput = () => this.input;
        if (input instanceof Input) {
            this.input = input;
        }
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { input: this.input.serialize(encoding) });
    }
    toBuffer() {
        const inbuff = this.input.toBuffer();
        const inid = buffer_1.Buffer.alloc(4);
        inid.writeUInt32BE(this.input.getInputID(), 0);
        const barr = [inid, inbuff];
        return buffer_1.Buffer.concat(barr, inid.length + inbuff.length);
    }
}
exports.StandardParseableInput = StandardParseableInput;
/**
 * Returns a function used to sort an array of [[StandardParseableInput]]s
 */
StandardParseableInput.comparator = () => (a, b) => {
    const sorta = a.toBuffer();
    const sortb = b.toBuffer();
    return buffer_1.Buffer.compare(sorta, sortb);
};
class StandardTransferableInput extends StandardParseableInput {
    /**
     * Class representing an [[StandardTransferableInput]] for a transaction.
     *
     * @param txid A {@link https://github.com/feross/buffer|Buffer} containing the transaction ID of the referenced UTXO
     * @param outputidx A {@link https://github.com/feross/buffer|Buffer} containing the index of the output in the transaction consumed in the [[StandardTransferableInput]]
     * @param assetID A {@link https://github.com/feross/buffer|Buffer} representing the assetID of the [[Input]]
     * @param input An [[Input]] to be made transferable
     */
    constructor(txid = undefined, outputidx = undefined, assetID = undefined, input = undefined) {
        super();
        this._typeName = "StandardTransferableInput";
        this._typeID = undefined;
        this.txid = buffer_1.Buffer.alloc(32);
        this.outputidx = buffer_1.Buffer.alloc(4);
        this.assetID = buffer_1.Buffer.alloc(32);
        /**
         * Returns a {@link https://github.com/feross/buffer|Buffer} of the TxID.
         */
        this.getTxID = () => this.txid;
        /**
         * Returns a {@link https://github.com/feross/buffer|Buffer}  of the OutputIdx.
         */
        this.getOutputIdx = () => this.outputidx;
        /**
         * Returns a base-58 string representation of the UTXOID this [[StandardTransferableInput]] references.
         */
        this.getUTXOID = () => bintools.bufferToB58(buffer_1.Buffer.concat([this.txid, this.outputidx]));
        /**
         * Returns the input.
         */
        this.getInput = () => this.input;
        /**
         * Returns the assetID of the input.
         */
        this.getAssetID = () => this.assetID;
        if (typeof txid !== "undefined" &&
            typeof outputidx !== "undefined" &&
            typeof assetID !== "undefined" &&
            input instanceof Input) {
            this.input = input;
            this.txid = txid;
            this.outputidx = outputidx;
            this.assetID = assetID;
        }
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { txid: serialization.encoder(this.txid, encoding, "Buffer", "cb58"), outputidx: serialization.encoder(this.outputidx, encoding, "Buffer", "decimalString"), assetID: serialization.encoder(this.assetID, encoding, "Buffer", "cb58") });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.txid = serialization.decoder(fields["txid"], encoding, "cb58", "Buffer", 32);
        this.outputidx = serialization.decoder(fields["outputidx"], encoding, "decimalString", "Buffer", 4);
        this.assetID = serialization.decoder(fields["assetID"], encoding, "cb58", "Buffer", 32);
        //input deserialization must be implmented in child classes
    }
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[StandardTransferableInput]].
     */
    toBuffer() {
        const parseableBuff = super.toBuffer();
        const bsize = this.txid.length +
            this.outputidx.length +
            this.assetID.length +
            parseableBuff.length;
        const barr = [
            this.txid,
            this.outputidx,
            this.assetID,
            parseableBuff
        ];
        const buff = buffer_1.Buffer.concat(barr, bsize);
        return buff;
    }
    /**
     * Returns a base-58 representation of the [[StandardTransferableInput]].
     */
    toString() {
        /* istanbul ignore next */
        return bintools.bufferToB58(this.toBuffer());
    }
}
exports.StandardTransferableInput = StandardTransferableInput;
/**
 * An [[Input]] class which specifies a token amount .
 */
class StandardAmountInput extends Input {
    /**
     * An [[AmountInput]] class which issues a payment on an assetID.
     *
     * @param amount A {@link https://github.com/indutny/bn.js/|BN} representing the amount in the input
     */
    constructor(amount = undefined) {
        super();
        this._typeName = "StandardAmountInput";
        this._typeID = undefined;
        this.amount = buffer_1.Buffer.alloc(8);
        this.amountValue = new bn_js_1.default(0);
        /**
         * Returns the amount as a {@link https://github.com/indutny/bn.js/|BN}.
         */
        this.getAmount = () => this.amountValue.clone();
        if (amount) {
            this.amountValue = amount.clone();
            this.amount = bintools.fromBNToBuffer(amount, 8);
        }
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { amount: serialization.encoder(this.amount, encoding, "Buffer", "decimalString", 8) });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.amount = serialization.decoder(fields["amount"], encoding, "decimalString", "Buffer", 8);
        this.amountValue = bintools.fromBufferToBN(this.amount);
    }
    /**
     * Popuates the instance from a {@link https://github.com/feross/buffer|Buffer} representing the [[AmountInput]] and returns the size of the input.
     */
    fromBuffer(bytes, offset = 0) {
        this.amount = bintools.copyFrom(bytes, offset, offset + 8);
        this.amountValue = bintools.fromBufferToBN(this.amount);
        offset += 8;
        return super.fromBuffer(bytes, offset);
    }
    /**
     * Returns the buffer representing the [[AmountInput]] instance.
     */
    toBuffer() {
        const superbuff = super.toBuffer();
        const bsize = this.amount.length + superbuff.length;
        const barr = [this.amount, superbuff];
        return buffer_1.Buffer.concat(barr, bsize);
    }
}
exports.StandardAmountInput = StandardAmountInput;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiaW5wdXQuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi9zcmMvY29tbW9uL2lucHV0LnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7OztBQUFBOzs7R0FHRztBQUNILG9DQUFnQztBQUNoQyxpRUFBd0M7QUFDeEMsa0RBQXNCO0FBQ3RCLCtDQUFzQztBQUN0QywwREFJK0I7QUFFL0I7O0dBRUc7QUFDSCxNQUFNLFFBQVEsR0FBYSxrQkFBUSxDQUFDLFdBQVcsRUFBRSxDQUFBO0FBQ2pELE1BQU0sYUFBYSxHQUFrQiw2QkFBYSxDQUFDLFdBQVcsRUFBRSxDQUFBO0FBRWhFLE1BQXNCLEtBQU0sU0FBUSw0QkFBWTtJQUFoRDs7UUFDWSxjQUFTLEdBQUcsT0FBTyxDQUFBO1FBQ25CLFlBQU8sR0FBRyxTQUFTLENBQUE7UUFtQm5CLGFBQVEsR0FBVyxlQUFNLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxDQUFBO1FBQ2xDLFlBQU8sR0FBYSxFQUFFLENBQUEsQ0FBQyw0QkFBNEI7UUEwQjdEOztXQUVHO1FBQ0gsZUFBVSxHQUFHLEdBQWEsRUFBRSxDQUFDLElBQUksQ0FBQyxPQUFPLENBQUE7UUFJekM7Ozs7O1dBS0c7UUFDSCxvQkFBZSxHQUFHLENBQUMsVUFBa0IsRUFBRSxPQUFlLEVBQUUsRUFBRTtZQUN4RCxNQUFNLE1BQU0sR0FBVyxJQUFJLG9CQUFNLEVBQUUsQ0FBQTtZQUNuQyxNQUFNLENBQUMsR0FBVyxlQUFNLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxDQUFBO1lBQ2pDLENBQUMsQ0FBQyxhQUFhLENBQUMsVUFBVSxFQUFFLENBQUMsQ0FBQyxDQUFBO1lBQzlCLE1BQU0sQ0FBQyxVQUFVLENBQUMsQ0FBQyxDQUFDLENBQUE7WUFDcEIsTUFBTSxDQUFDLFNBQVMsQ0FBQyxPQUFPLENBQUMsQ0FBQTtZQUN6QixJQUFJLENBQUMsT0FBTyxDQUFDLElBQUksQ0FBQyxNQUFNLENBQUMsQ0FBQTtZQUN6QixJQUFJLENBQUMs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UUFBUSxDQUFDLGNBQWMsQ0FBQyxJQUFJLENBQUMsTUFBTSxDQUFDLENBQUE7UUFDdkQsTUFBTSxJQUFJLENBQUMsQ0FBQTtRQUNYLE9BQU8sS0FBSyxDQUFDLFVBQVUsQ0FBQyxLQUFLLEVBQUUsTUFBTSxDQUFDLENBQUE7SUFDeEMsQ0FBQztJQUVEOztPQUVHO0lBQ0gsUUFBUTtRQUNOLE1BQU0sU0FBUyxHQUFXLEtBQUssQ0FBQyxRQUFRLEVBQUUsQ0FBQTtRQUMxQyxNQUFNLEtBQUssR0FBVyxJQUFJLENBQUMsTUFBTSxDQUFDLE1BQU0sR0FBRyxTQUFTLENBQUMsTUFBTSxDQUFBO1FBQzNELE1BQU0sSUFBSSxHQUFhLENBQUMsSUFBSSxDQUFDLE1BQU0sRUFBRSxTQUFTLENBQUMsQ0FBQTtRQUMvQyxPQUFPLGVBQU0sQ0FBQyxNQUFNLENBQUMsSUFBSSxFQUFFLEtBQUssQ0FBQyxDQUFBO0lBQ25DLENBQUM7Q0FjRjtBQXJFRCxrREFxRUMiLCJzb3VyY2VzQ29udGVudCI6WyIvKipcbiAqIEBwYWNrYWdlRG9jdW1lbnRhdGlvblxuICogQG1vZHVsZSBDb21tb24tSW5wdXRzXG4gKi9cbmltcG9ydCB7IEJ1ZmZlciB9IGZyb20gXCJidWZmZXIvXCJcbmltcG9ydCBCaW5Ub29scyBmcm9tIFwiLi4vdXRpbHMvYmludG9vbHNcIlxuaW1wb3J0IEJOIGZyb20gXCJibi5qc1wiXG5pbXBvcnQgeyBTaWdJZHggfSBmcm9tIFwiLi9jcmVkZW50aWFsc1wiXG5pbXBvcnQge1xuICBTZXJpYWxpemFibGUsXG4gIFNlcmlhbGl6YXRpb24sXG4gIFNlcmlhbGl6ZWRFbmNvZGluZ1xufSBmcm9tIFwiLi4vdXRpbHMvc2VyaWFsaXphdGlvblwiXG5cbi8qKlxuICogQGlnbm9yZVxuICovXG5jb25zdCBiaW50b29sczogQmluVG9vbHMgPSBCaW5Ub29scy5nZXRJbnN0YW5jZSgpXG5jb25zdCBzZXJpYWxpemF0aW9uOiBTZXJpYWxpemF0aW9uID0gU2VyaWFsaXphdGlvbi5nZXRJbnN0YW5jZSgpXG5cbmV4cG9ydCBhYnN0cmFjdCBjbGFzcyBJbnB1dCBleHRlbmRzIFNlcmlhbGl6YWJsZSB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIklucHV0XCJcbiAgcHJvdGVjdGVkIF90eXBlSUQgPSB1bmRlZmluZWRcblxuICBzZXJpYWxpemUoZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpOiBvYmplY3Qge1xuICAgIGxldCBmaWVsZHM6IG9iamVjdCA9IHN1cGVyLnNlcmlhbGl6ZShlbmNvZGluZylcbiAgICByZXR1cm4ge1xuICAgICAgLi4uZmllbGRzLFxuICAgICAgc2lnSWR4czogdGhpcy5zaWdJZHhzLm1hcCgocykgPT4gcy5zZXJpYWxpemUoZW5jb2RpbmcpKVxuICAgIH1cbiAgfVxuICBkZXNlcmlhbGl6ZShmaWVsZHM6IG9iamVjdCwgZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpIHtcbiAgICBzdXBlci5kZXNlcmlhbGl6ZShmaWVsZHMsIGVuY29kaW5nKVxuICAgIHRoaXMuc2lnSWR4cyA9IGZpZWxkc1tcInNpZ0lkeHNcIl0ubWFwKChzOiBvYmplY3QpID0+IHtcbiAgICAgIGxldCBzaWR4OiBTaWdJZHggPSBuZXcgU2lnSWR4KClcbiAgICAgIHNpZHguZGVzZXJpYWxpemUocywgZW5jb2RpbmcpXG4gICAgICByZXR1cm4gc2lkeFxuICAgIH0pXG4gICAgdGhpcy5zaWdDb3VudC53cml0ZVVJbnQzMkJFKHRoaXMuc2lnSWR4cy5sZW5ndGgsIDApXG4gIH1cblxuICBwcm90ZWN0ZWQgc2lnQ291bnQ6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYyg0KVxuICBwcm90ZWN0ZWQgc2lnSWR4czogU2lnSWR4W10gPSBbXSAvLyBpZHhzIG9mIHNpZ25lcnMgZnJvbSB1dHhvXG5cbiAgc3RhdGljIGNvbXBhcmF0b3IgPVxuICAgICgpOiAoKGE6IElucHV0LCBiOiBJbnB1dCkgPT4gMSB8IC0xIHwgMCkgPT5cbiAgICAoYTogSW5wdXQsIGI6IElucHV0KTogMSB8IC0xIHwgMCA9PiB7XG4gICAgICBjb25zdCBhb3V0aWQ6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYyg0KVxuICAgICAgYW91dGlkLndyaXRlVUludDMyQkUoYS5nZXRJbnB1dElEKCksIDApXG4gICAgICBjb25zdCBhYnVmZjogQnVmZmVyID0gYS50b0J1ZmZlcigpXG5cbiAgICAgIGNvbnN0IGJvdXRpZDogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDQpXG4gICAgICBib3V0aWQud3JpdGVVSW50MzJCRShiLmdldElucHV0SUQoKSwgMClcbiAgICAgIGNvbnN0IGJidWZmOiBCdWZmZXIgPSBiLnRvQnVmZmVyKClcblxuICAgICAgY29uc3QgYXNvcnQ6IEJ1ZmZlciA9IEJ1ZmZlci5jb25jYXQoXG4gICAgICAgIFthb3V0aWQsIGFidWZmXSxcbiAgICAgICAgYW91dGlkLmxlbmd0aCArIGFidWZmLmxlbmd0aFxuICAgICAgKVxuICAgICAgY29uc3QgYnNvcnQ6IEJ1ZmZlciA9IEJ1ZmZlci5jb25jYXQoXG4gICAgICAgIFtib3V0aWQsIGJidWZmXSxcbiAgICAgICAgYm91dGlkLmxlbmd0aCArIGJidWZmLmxlbmd0aFxuICAgICAgKVxuICAgICAgcmV0dXJuIEJ1ZmZlci5jb21wYXJlKGFzb3J0LCBic29ydCkgYXMgMSB8IC0xIHwgMFxuICAgIH1cblxuICBhYnN0cmFjdCBnZXRJbnB1dElEKCk6IG51bWJlclxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBhcnJheSBvZiBbW1NpZ0lkeF1dIGZvciB0aGlzIFtbSW5wdXRdXVxuICAgKi9cbiAgZ2V0U2lnSWR4cyA9ICgpOiBTaWdJZHhbXSA9PiB0aGlzLnNpZ0lkeHNcblxuICBhYnN0cmFjdCBnZXRDcmVkZW50aWFsSUQoKTogbnVtYmVyXG5cbiAgLyoqXG4gICAqIENyZWF0ZXMgYW5kIGFkZHMgYSBbW1NpZ0lkeF1dIHRvIHRoZSBbW0lucHV0XV0uXG4gICAqXG4gICAqIEBwYXJhbSBhZGRyZXNzSWR4IFRoZSBpbmRleCBvZiB0aGUgYWRkcmVzcyB0byByZWZlcmVuY2UgaW4gdGhlIHNpZ25hdHVyZXNcbiAgICogQHBhcmFtIGFkZHJlc3MgVGhlIGFkZHJlc3Mgb2YgdGhlIHNvdXJjZSBvZiB0aGUgc2lnbmF0dXJlXG4gICAqL1xuICBhZGRTaWduYXR1cmVJZHggPSAoYWRkcmVzc0lkeDogbnVtYmVyLCBhZGRyZXNzOiBCdWZmZXIpID0+IHtcbiAgICBjb25zdCBzaWdpZHg6IFNpZ0lkeCA9IG5ldyBTaWdJZHgoKVxuICAgIGNvbnN0IGI6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYyg0KVxuICAgIGIud3JpdGVVSW50MzJCRShhZGRyZXNzSWR4LCAwKVxuICAgIHNpZ2lkeC5mcm9tQnVmZmVyKGIpXG4gICAgc2lnaWR4LnNldFNvdXJjZShhZGRyZXNzKVxuICAgIHRoaXMuc2lnSWR4cy5wdXNoKHNpZ2lkeClcbiAgICB0aGlzLnNpZ0NvdW50LndyaXRlVUludDMyQkUodGhpcy5zaWdJZHhzLmxlbmd0aCwgMClcbiAgfVxuXG4gIGZyb21CdWZmZXIoYnl0ZXM6IEJ1ZmZlciwgb2Zmc2V0OiBudW1iZXIgPSAwKTogbnVtYmVyIHtcbiAgICB0aGlzLnNpZ0NvdW50ID0gYmludG9vbHMuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgNClcbiAgICBvZmZzZXQgKz0gNFxuICAgIGNvbnN0IHNpZ0NvdW50OiBudW1iZXIgPSB0aGlzLnNpZ0NvdW50LnJlYWRVSW50MzJCRSgwKVxuICAgIHRoaXMuc2lnSWR4cyA9IFtdXG4gICAgZm9yIChsZXQgaTogbnVtYmVyID0gMDsgaSA8IHNpZ0NvdW50OyBpKyspIHtcbiAgICAgIGNvbnN0IHNpZ2lkeCA9IG5ldyBTaWdJZHgoKVxuICAgICAgY29uc3Qgc2lnYnVmZjogQnVmZmVyID0gYmludG9vbHMuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgNClcbiAgICAgIHNpZ2lkeC5mcm9tQnVmZmVyKHNpZ2J1ZmYpXG4gICAgICBvZmZzZXQgKz0gNFxuICAgICAgdGhpcy5zaWdJZHhzLnB1c2goc2lnaWR4KVxuICAgIH1cbiAgICByZXR1cm4gb2Zmc2V0XG4gIH1cblxuICB0b0J1ZmZlcigpOiBCdWZmZXIge1xuICAgIHRoaXMuc2lnQ291bnQud3JpdGVVSW50MzJCRSh0aGlzLnNpZ0lkeHMubGVuZ3RoLCAwKVxuICAgIGxldCBic2l6ZTogbnVtYmVyID0gdGhpcy5zaWdDb3VudC5sZW5ndGhcbiAgICBjb25zdCBiYXJyOiBCdWZmZXJbXSA9IFt0aGlzLnNpZ0NvdW50XVxuICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCB0aGlzLnNpZ0lkeHMubGVuZ3RoOyBpKyspIHtcbiAgICAgIGNvbnN0IGI6IEJ1ZmZlciA9IHRoaXMuc2lnSWR4c1tgJHtpfWBdLnRvQnVmZmVyKClcbiAgICAgIGJhcnIucHVzaChiKVxuICAgICAgYnNpemUgKz0gYi5sZW5ndGhcbiAgICB9XG4gICAgcmV0dXJuIEJ1ZmZlci5jb25jYXQoYmFyciwgYnNpemUpXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyBhIGJhc2UtNTggcmVwcmVzZW50YXRpb24gb2YgdGhlIFtbSW5wdXRdXS5cbiAgICovXG4gIHRvU3RyaW5nKCk6IHN0cmluZyB7XG4gICAgcmV0dXJuIGJpbnRvb2xzLmJ1ZmZlclRvQjU4KHRoaXMudG9CdWZmZXIoKSlcbiAgfVxuXG4gIGFic3RyYWN0IGNsb25lKCk6IHRoaXNcblxuICBhYnN0cmFjdCBjcmVhdGUoLi4uYXJnczogYW55W10pOiB0aGlzXG5cbiAgYWJzdHJhY3Qgc2VsZWN0KGlkOiBudW1iZXIsIC4uLmFyZ3M6IGFueVtdKTogSW5wdXRcbn1cblxuZXhwb3J0IGFic3RyYWN0IGNsYXNzIFN0YW5kYXJkUGFyc2VhYmxlSW5wdXQgZXh0ZW5kcyBTZXJpYWxpemFibGUge1xuICBwcm90ZWN0ZWQgX3R5cGVOYW1lID0gXCJTdGFuZGFyZFBhcnNlYWJsZUlucHV0XCJcbiAgcHJvdGVjdGVkIF90eXBlSUQgPSB1bmRlZmluZWRcblxuICBzZXJpYWxpemUoZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpOiBvYmplY3Qge1xuICAgIGxldCBmaWVsZHM6IG9iamVjdCA9IHN1cGVyLnNlcmlhbGl6ZShlbmNvZGluZylcbiAgICByZXR1cm4ge1xuICAgICAgLi4uZmllbGRzLFxuICAgICAgaW5wdXQ6IHRoaXMuaW5wdXQuc2VyaWFsaXplKGVuY29kaW5nKVxuICAgIH1cbiAgfVxuXG4gIHByb3RlY3RlZCBpbnB1dDogSW5wdXRcblxuICAvKipcbiAgICogUmV0dXJucyBhIGZ1bmN0aW9uIHVzZWQgdG8gc29ydCBhbiBhcnJheSBvZiBbW1N0YW5kYXJkUGFyc2VhYmxlSW5wdXRdXXNcbiAgICovXG4gIHN0YXRpYyBjb21wYXJhdG9yID1cbiAgICAoKTogKChcbiAgICAgIGE6IFN0YW5kYXJkUGFyc2VhYmxlSW5wdXQsXG4gICAgICBiOiBTdGFuZGFyZFBhcnNlYWJsZUlucHV0XG4gICAgKSA9PiAxIHwgLTEgfCAwKSA9PlxuICAgIChhOiBTdGFuZGFyZFBhcnNlYWJsZUlucHV0LCBiOiBTdGFuZGFyZFBhcnNlYWJsZUlucHV0KTogMSB8IC0xIHwgMCA9PiB7XG4gICAgICBjb25zdCBzb3J0YSA9IGEudG9CdWZmZXIoKVxuICAgICAgY29uc3Qgc29ydGIgPSBiLnRvQnVmZmVyKClcbiAgICAgIHJldHVybiBCdWZmZXIuY29tcGFyZShzb3J0YSwgc29ydGIpIGFzIDEgfCAtMSB8IDBcbiAgICB9XG5cbiAgZ2V0SW5wdXQgPSAoKTogSW5wdXQgPT4gdGhpcy5pbnB1dFxuXG4gIC8vIG11c3QgYmUgaW1wbGVtZW50ZWQgdG8gc2VsZWN0IGlucHV0IHR5cGVzIGZvciB0aGUgVk0gaW4gcXVlc3Rpb25cbiAgYWJzdHJhY3QgZnJvbUJ1ZmZlcihieXRlczogQnVmZmVyLCBvZmZzZXQ/OiBudW1iZXIpOiBudW1iZXJcblxuICB0b0J1ZmZlcigpOiBCdWZmZXIge1xuICAgIGNvbnN0IGluYnVmZjogQnVmZmVyID0gdGhpcy5pbnB1dC50b0J1ZmZlcigpXG4gICAgY29uc3QgaW5pZDogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDQpXG4gICAgaW5pZC53cml0ZVVJbnQzMkJFKHRoaXMuaW5wdXQuZ2V0SW5wdXRJRCgpLCAwKVxuICAgIGNvbnN0IGJhcnI6IEJ1ZmZlcltdID0gW2luaWQsIGluYnVmZl1cbiAgICByZXR1cm4gQnVmZmVyLmNvbmNhdChiYXJyLCBpbmlkLmxlbmd0aCArIGluYnVmZi5sZW5ndGgpXG4gIH1cblxuICAvKipcbiAgICogQ2xhc3MgcmVwcmVzZW50aW5nIGFuIFtbU3RhbmRhcmRQYXJzZWFibGVJbnB1dF1dIGZvciBhIHRyYW5zYWN0aW9uLlxuICAgKlxuICAgKiBAcGFyYW0gaW5wdXQgQSBudW1iZXIgcmVwcmVzZW50aW5nIHRoZSBJbnB1dElEIG9mIHRoZSBbW1N0YW5kYXJkUGFyc2VhYmxlSW5wdXRdXVxuICAgKi9cbiAgY29uc3RydWN0b3IoaW5wdXQ6IElucHV0ID0gdW5kZWZpbmVkKSB7XG4gICAgc3VwZXIoKVxuICAgIGlmIChpbnB1dCBpbnN0YW5jZW9mIElucHV0KSB7XG4gICAgICB0aGlzLmlucHV0ID0gaW5wdXRcbiAgICB9XG4gIH1cbn1cblxuZXhwb3J0IGFic3RyYWN0IGNsYXNzIFN0YW5kYXJkVHJhbnNmZXJhYmxlSW5wdXQgZXh0ZW5kcyBTdGFuZGFyZFBhcnNlYWJsZUlucHV0IHtcbiAgcHJvdGVjdGVkIF90eXBlTmFtZSA9IFwiU3RhbmRhcmRUcmFuc2ZlcmFibGVJbnB1dFwiXG4gIHByb3RlY3RlZCBfdHlwZUlEID0gdW5kZWZpbmVkXG5cbiAgc2VyaWFsaXplKGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKTogb2JqZWN0IHtcbiAgICBsZXQgZmllbGRzOiBvYmplY3QgPSBzdXBlci5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgcmV0dXJuIHtcbiAgICAgIC4uLmZpZWxkcyxcbiAgICAgIHR4aWQ6IHNlcmlhbGl6YXRpb24uZW5jb2Rlcih0aGlzLnR4aWQsIGVuY29kaW5nLCBcIkJ1ZmZlclwiLCBcImNiNThcIiksXG4gICAgICBvdXRwdXRpZHg6IHNlcmlhbGl6YXRpb24uZW5jb2RlcihcbiAgICAgICAgdGhpcy5vdXRwdXRpZHgsXG4gICAgICAgIGVuY29kaW5nLFxuICAgICAgICBcIkJ1ZmZlclwiLFxuICAgICAgICBcImRlY2ltYWxTdHJpbmdcIlxuICAgICAgKSxcbiAgICAgIGFzc2V0SUQ6IHNlcmlhbGl6YXRpb24uZW5jb2Rlcih0aGlzLmFzc2V0SUQsIGVuY29kaW5nLCBcIkJ1ZmZlclwiLCBcImNiNThcIilcbiAgICB9XG4gIH1cbiAgZGVzZXJpYWxpemUoZmllbGRzOiBvYmplY3QsIGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKSB7XG4gICAgc3VwZXIuZGVzZXJpYWxpemUoZmllbGRzLCBlbmNvZGluZylcbiAgICB0aGlzLnR4aWQgPSBzZXJpYWxpemF0aW9uLmRlY29kZXIoXG4gICAgICBmaWVsZHNbXCJ0eGlkXCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBcImNiNThcIixcbiAgICAgIFwiQnVmZmVyXCIsXG4gICAgICAzMlxuICAgIClcbiAgICB0aGlzLm91dHB1dGlkeCA9IHNlcmlhbGl6YXRpb24uZGVjb2RlcihcbiAgICAgIGZpZWxkc1tcIm91dHB1dGlkeFwiXSxcbiAgICAgIGVuY29kaW5nLFxuICAgICAgXCJkZWNpbWFsU3RyaW5nXCIsXG4gICAgICBcIkJ1ZmZlclwiLFxuICAgICAgNFxuICAgIClcbiAgICB0aGlzLmFzc2V0SUQgPSBzZXJpYWxpemF0aW9uLmRlY29kZXIoXG4gICAgICBmaWVsZHNbXCJhc3NldElEXCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBcImNiNThcIixcbiAgICAgIFwiQnVmZmVyXCIsXG4gICAgICAzMlxuICAgIClcbiAgICAvL2lucHV0IGRlc2VyaWFsaXphdGlvbiBtdXN0IGJlIGltcGxtZW50ZWQgaW4gY2hpbGQgY2xhc3Nlc1xuICB9XG5cbiAgcHJvdGVjdGVkIHR4aWQ6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYygzMilcbiAgcHJvdGVjdGVkIG91dHB1dGlkeDogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDQpXG4gIHByb3RlY3RlZCBhc3NldElEOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoMzIpXG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSBvZiB0aGUgVHhJRC5cbiAgICovXG4gIGdldFR4SUQgPSAoKTogLyogaXN0YW5idWwgaWdub3JlIG5leHQgKi8gQnVmZmVyID0+IHRoaXMudHhpZFxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIGEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gIG9mIHRoZSBPdXRwdXRJZHguXG4gICAqL1xuICBnZXRPdXRwdXRJZHggPSAoKTogLyogaXN0YW5idWwgaWdub3JlIG5leHQgKi8gQnVmZmVyID0+IHRoaXMub3V0cHV0aWR4XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSBiYXNlLTU4IHN0cmluZyByZXByZXNlbnRhdGlvbiBvZiB0aGUgVVRYT0lEIHRoaXMgW1tTdGFuZGFyZFRyYW5zZmVyYWJsZUlucHV0XV0gcmVmZXJlbmNlcy5cbiAgICovXG4gIGdldFVUWE9JRCA9ICgpOiBzdHJpbmcgPT5cbiAgICBiaW50b29scy5idWZmZXJUb0I1OChCdWZmZXIuY29uY2F0KFt0aGlzLnR4aWQsIHRoaXMub3V0cHV0aWR4XSkpXG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIGlucHV0LlxuICAgKi9cbiAgZ2V0SW5wdXQgPSAoKTogSW5wdXQgPT4gdGhpcy5pbnB1dFxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBhc3NldElEIG9mIHRoZSBpbnB1dC5cbiAgICovXG4gIGdldEFzc2V0SUQgPSAoKTogQnVmZmVyID0+IHRoaXMuYXNzZXRJRFxuXG4gIGFic3RyYWN0IGZyb21CdWZmZXIoYnl0ZXM6IEJ1ZmZlciwgb2Zmc2V0PzogbnVtYmVyKTogbnVtYmVyXG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSByZXByZXNlbnRhdGlvbiBvZiB0aGUgW1tTdGFuZGFyZFRyYW5zZmVyYWJsZUlucHV0XV0uXG4gICAqL1xuICB0b0J1ZmZlcigpOiBCdWZmZXIge1xuICAgIGNvbnN0IHBhcnNlYWJsZUJ1ZmY6IEJ1ZmZlciA9IHN1cGVyLnRvQnVmZmVyKClcbiAgICBjb25zdCBic2l6ZTogbnVtYmVyID1cbiAgICAgIHRoaXMudHhpZC5sZW5ndGggK1xuICAgICAgdGhpcy5vdXRwdXRpZHgubGVuZ3RoICtcbiAgICAgIHRoaXMuYXNzZXRJRC5sZW5ndGggK1xuICAgICAgcGFyc2VhYmxlQnVmZi5sZW5ndGhcbiAgICBjb25zdCBiYXJyOiBCdWZmZXJbXSA9IFtcbiAgICAgIHRoaXMudHhpZCxcbiAgICAgIHRoaXMub3V0cHV0aWR4LFxuICAgICAgdGhpcy5hc3NldElELFxuICAgICAgcGFyc2VhYmxlQnVmZlxuICAgIF1cbiAgICBjb25zdCBidWZmOiBCdWZmZXIgPSBCdWZmZXIuY29uY2F0KGJhcnIsIGJzaXplKVxuICAgIHJldHVybiBidWZmXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyBhIGJhc2UtNTggcmVwcmVzZW50YXRpb24gb2YgdGhlIFtbU3RhbmRhcmRUcmFuc2ZlcmFibGVJbnB1dF1dLlxuICAgKi9cbiAgdG9TdHJpbmcoKTogc3RyaW5nIHtcbiAgICAvKiBpc3RhbmJ1bCBpZ25vcmUgbmV4dCAqL1xuICAgIHJldHVybiBiaW50b29scy5idWZmZXJUb0I1OCh0aGlzLnRvQnVmZmVyKCkpXG4gIH1cblxuICAvKipcbiAgICogQ2xhc3MgcmVwcmVzZW50aW5nIGFuIFtbU3RhbmRhcmRUcmFuc2ZlcmFibGVJbnB1dF1dIGZvciBhIHRyYW5zYWN0aW9uLlxuICAgKlxuICAgKiBAcGFyYW0gdHhpZCBBIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IGNvbnRhaW5pbmcgdGhlIHRyYW5zYWN0aW9uIElEIG9mIHRoZSByZWZlcmVuY2VkIFVUWE9cbiAgICogQHBhcmFtIG91dHB1dGlkeCBBIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IGNvbnRhaW5pbmcgdGhlIGluZGV4IG9mIHRoZSBvdXRwdXQgaW4gdGhlIHRyYW5zYWN0aW9uIGNvbnN1bWVkIGluIHRoZSBbW1N0YW5kYXJkVHJhbnNmZXJhYmxlSW5wdXRdXVxuICAgKiBAcGFyYW0gYXNzZXRJRCBBIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IHJlcHJlc2VudGluZyB0aGUgYXNzZXRJRCBvZiB0aGUgW1tJbnB1dF1dXG4gICAqIEBwYXJhbSBpbnB1dCBBbiBbW0lucHV0XV0gdG8gYmUgbWFkZSB0cmFuc2ZlcmFibGVcbiAgICovXG4gIGNvbnN0cnVjdG9yKFxuICAgIHR4aWQ6IEJ1ZmZlciA9IHVuZGVmaW5lZCxcbiAgICBvdXRwdXRpZHg6IEJ1ZmZlciA9IHVuZGVmaW5lZCxcbiAgICBhc3NldElEOiBCdWZmZXIgPSB1bmRlZmluZWQsXG4gICAgaW5wdXQ6IElucHV0ID0gdW5kZWZpbmVkXG4gICkge1xuICAgIHN1cGVyKClcbiAgICBpZiAoXG4gICAgICB0eXBlb2YgdHhpZCAhPT0gXCJ1bmRlZmluZWRcIiAmJlxuICAgICAgdHlwZW9mIG91dHB1dGlkeCAhPT0gXCJ1bmRlZmluZWRcIiAmJlxuICAgICAgdHlwZW9mIGFzc2V0SUQgIT09IFwidW5kZWZpbmVkXCIgJiZcbiAgICAgIGlucHV0IGluc3RhbmNlb2YgSW5wdXRcbiAgICApIHtcbiAgICAgIHRoaXMuaW5wdXQgPSBpbnB1dFxuICAgICAgdGhpcy50eGlkID0gdHhpZFxuICAgICAgdGhpcy5vdXRwdXRpZHggPSBvdXRwdXRpZHhcbiAgICAgIHRoaXMuYXNzZXRJRCA9IGFzc2V0SURcbiAgICB9XG4gIH1cbn1cblxuLyoqXG4gKiBBbiBbW0lucHV0XV0gY2xhc3Mgd2hpY2ggc3BlY2lmaWVzIGEgdG9rZW4gYW1vdW50IC5cbiAqL1xuZXhwb3J0IGFic3RyYWN0IGNsYXNzIFN0YW5kYXJkQW1vdW50SW5wdXQgZXh0ZW5kcyBJbnB1dCB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIlN0YW5kYXJkQW1vdW50SW5wdXRcIlxuICBwcm90ZWN0ZWQgX3R5cGVJRCA9IHVuZGVmaW5lZFxuXG4gIHNlcmlhbGl6ZShlbmNvZGluZzogU2VyaWFsaXplZEVuY29kaW5nID0gXCJoZXhcIik6IG9iamVjdCB7XG4gICAgbGV0IGZpZWxkczogb2JqZWN0ID0gc3VwZXIuc2VyaWFsaXplKGVuY29kaW5nKVxuICAgIHJldHVybiB7XG4gICAgICAuLi5maWVsZHMsXG4gICAgICBhbW91bnQ6IHNlcmlhbGl6YXRpb24uZW5jb2RlcihcbiAgICAgICAgdGhpcy5hbW91bnQsXG4gICAgICAgIGVuY29kaW5nLFxuICAgICAgICBcIkJ1ZmZlclwiLFxuICAgICAgICBcImRlY2ltYWxTdHJpbmdcIixcbiAgICAgICAgOFxuICAgICAgKVxuICAgIH1cbiAgfVxuICBkZXNlcmlhbGl6ZShmaWVsZHM6IG9iamVjdCwgZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpIHtcbiAgICBzdXBlci5kZXNlcmlhbGl6ZShmaWVsZHMsIGVuY29kaW5nKVxuICAgIHRoaXMuYW1vdW50ID0gc2VyaWFsaXphdGlvbi5kZWNvZGVyKFxuICAgICAgZmllbGRzW1wiYW1vdW50XCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBcImRlY2ltYWxTdHJpbmdcIixcbiAgICAgIFwiQnVmZmVyXCIsXG4gICAgICA4XG4gICAgKVxuICAgIHRoaXMuYW1vdW50VmFsdWUgPSBiaW50b29scy5mcm9tQnVmZmVyVG9CTih0aGlzLmFtb3VudClcbiAgfVxuXG4gIHByb3RlY3RlZCBhbW91bnQ6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYyg4KVxuICBwcm90ZWN0ZWQgYW1vdW50VmFsdWU6IEJOID0gbmV3IEJOKDApXG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIGFtb3VudCBhcyBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vaW5kdXRueS9ibi5qcy98Qk59LlxuICAgKi9cbiAgZ2V0QW1vdW50ID0gKCk6IEJOID0+IHRoaXMuYW1vdW50VmFsdWUuY2xvbmUoKVxuXG4gIC8qKlxuICAgKiBQb3B1YXRlcyB0aGUgaW5zdGFuY2UgZnJvbSBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IHJlcHJlc2VudGluZyB0aGUgW1tBbW91bnRJbnB1dF1dIGFuZCByZXR1cm5zIHRoZSBzaXplIG9mIHRoZSBpbnB1dC5cbiAgICovXG4gIGZyb21CdWZmZXIoYnl0ZXM6IEJ1ZmZlciwgb2Zmc2V0OiBudW1iZXIgPSAwKTogbnVtYmVyIHtcbiAgICB0aGlzLmFtb3VudCA9IGJpbnRvb2xzLmNvcHlGcm9tKGJ5dGVzLCBvZmZzZXQsIG9mZnNldCArIDgpXG4gICAgdGhpcy5hbW91bnRWYWx1ZSA9IGJpbnRvb2xzLmZyb21CdWZmZXJUb0JOKHRoaXMuYW1vdW50KVxuICAgIG9mZnNldCArPSA4XG4gICAgcmV0dXJuIHN1cGVyLmZyb21CdWZmZXIoYnl0ZXMsIG9mZnNldClcbiAgfVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBidWZmZXIgcmVwcmVzZW50aW5nIHRoZSBbW0Ftb3VudElucHV0XV0gaW5zdGFuY2UuXG4gICAqL1xuICB0b0J1ZmZlcigpOiBCdWZmZXIge1xuICAgIGNvbnN0IHN1cGVyYnVmZjogQnVmZmVyID0gc3VwZXIudG9CdWZmZXIoKVxuICAgIGNvbnN0IGJzaXplOiBudW1iZXIgPSB0aGlzLmFtb3VudC5sZW5ndGggKyBzdXBlcmJ1ZmYubGVuZ3RoXG4gICAgY29uc3QgYmFycjogQnVmZmVyW10gPSBbdGhpcy5hbW91bnQsIHN1cGVyYnVmZl1cbiAgICByZXR1cm4gQnVmZmVyLmNvbmNhdChiYXJyLCBic2l6ZSlcbiAgfVxuXG4gIC8qKlxuICAgKiBBbiBbW0Ftb3VudElucHV0XV0gY2xhc3Mgd2hpY2ggaXNzdWVzIGEgcGF5bWVudCBvbiBhbiBhc3NldElELlxuICAgKlxuICAgKiBAcGFyYW0gYW1vdW50IEEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9pbmR1dG55L2JuLmpzL3xCTn0gcmVwcmVzZW50aW5nIHRoZSBhbW91bnQgaW4gdGhlIGlucHV0XG4gICAqL1xuICBjb25zdHJ1Y3RvcihhbW91bnQ6IEJOID0gdW5kZWZpbmVkKSB7XG4gICAgc3VwZXIoKVxuICAgIGlmIChhbW91bnQpIHtcbiAgICAgIHRoaXMuYW1vdW50VmFsdWUgPSBhbW91bnQuY2xvbmUoKVxuICAgICAgdGhpcy5hbW91bnQgPSBiaW50b29scy5mcm9tQk5Ub0J1ZmZlcihhbW91bnQsIDgpXG4gICAgfVxuICB9XG59XG4iXX0=