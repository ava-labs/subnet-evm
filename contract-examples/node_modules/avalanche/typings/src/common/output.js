"use strict";
/**
 * @packageDocumentation
 * @module Common-Output
 */
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.BaseNFTOutput = exports.StandardAmountOutput = exports.StandardTransferableOutput = exports.StandardParseableOutput = exports.Output = exports.OutputOwners = exports.Address = void 0;
const buffer_1 = require("buffer/");
const bn_js_1 = __importDefault(require("bn.js"));
const bintools_1 = __importDefault(require("../utils/bintools"));
const nbytes_1 = require("./nbytes");
const helperfunctions_1 = require("../utils/helperfunctions");
const serialization_1 = require("../utils/serialization");
const errors_1 = require("../utils/errors");
/**
 * @ignore
 */
const bintools = bintools_1.default.getInstance();
const serialization = serialization_1.Serialization.getInstance();
/**
 * Class for representing an address used in [[Output]] types
 */
class Address extends nbytes_1.NBytes {
    /**
     * Class for representing an address used in [[Output]] types
     */
    constructor() {
        super();
        this._typeName = "Address";
        this._typeID = undefined;
        //serialize and deserialize both are inherited
        this.bytes = buffer_1.Buffer.alloc(20);
        this.bsize = 20;
    }
    /**
     * Returns a base-58 representation of the [[Address]].
     */
    toString() {
        return bintools.cb58Encode(this.toBuffer());
    }
    /**
     * Takes a base-58 string containing an [[Address]], parses it, populates the class, and returns the length of the Address in bytes.
     *
     * @param bytes A base-58 string containing a raw [[Address]]
     *
     * @returns The length of the raw [[Address]]
     */
    fromString(addr) {
        const addrbuff = bintools.b58ToBuffer(addr);
        if (addrbuff.length === 24 && bintools.validateChecksum(addrbuff)) {
            const newbuff = bintools.copyFrom(addrbuff, 0, addrbuff.length - 4);
            if (newbuff.length === 20) {
                this.bytes = newbuff;
            }
        }
        else if (addrbuff.length === 24) {
            throw new errors_1.ChecksumError("Error - Address.fromString: invalid checksum on address");
        }
        else if (addrbuff.length === 20) {
            this.bytes = addrbuff;
        }
        else {
            /* istanbul ignore next */
            throw new errors_1.AddressError("Error - Address.fromString: invalid address");
        }
        return this.getSize();
    }
    clone() {
        let newbase = new Address();
        newbase.fromBuffer(this.toBuffer());
        return newbase;
    }
    create(...args) {
        return new Address();
    }
}
exports.Address = Address;
/**
 * Returns a function used to sort an array of [[Address]]es
 */
Address.comparator = () => (a, b) => buffer_1.Buffer.compare(a.toBuffer(), b.toBuffer());
/**
 * Defines the most basic values for output ownership. Mostly inherited from, but can be used in population of NFT Owner data.
 */
class OutputOwners extends serialization_1.Serializable {
    /**
     * An [[Output]] class which contains addresses, locktimes, and thresholds.
     *
     * @param addresses An array of {@link https://github.com/feross/buffer|Buffer}s representing output owner's addresses
     * @param locktime A {@link https://github.com/indutny/bn.js/|BN} representing the locktime
     * @param threshold A number representing the the threshold number of signers required to sign the transaction
     */
    constructor(addresses = undefined, locktime = undefined, threshold = undefined) {
        super();
        this._typeName = "OutputOwners";
        this._typeID = undefined;
        this.locktime = buffer_1.Buffer.alloc(8);
        this.threshold = buffer_1.Buffer.alloc(4);
        this.numaddrs = buffer_1.Buffer.alloc(4);
        this.addresses = [];
        /**
         * Returns the threshold of signers required to spend this output.
         */
        this.getThreshold = () => this.threshold.readUInt32BE(0);
        /**
         * Returns the a {@link https://github.com/indutny/bn.js/|BN} repersenting the UNIX Timestamp when the lock is made available.
         */
        this.getLocktime = () => bintools.fromBufferToBN(this.locktime);
        /**
         * Returns an array of {@link https://github.com/feross/buffer|Buffer}s for the addresses.
         */
        this.getAddresses = () => {
            const result = [];
            for (let i = 0; i < this.addresses.length; i++) {
                result.push(this.addresses[`${i}`].toBuffer());
            }
            return result;
        };
        /**
         * Returns the index of the address.
         *
         * @param address A {@link https://github.com/feross/buffer|Buffer} of the address to look up to return its index.
         *
         * @returns The index of the address.
         */
        this.getAddressIdx = (address) => {
            for (let i = 0; i < this.addresses.length; i++) {
                if (this.addresses[`${i}`].toBuffer().toString("hex") ===
                    address.toString("hex")) {
                    return i;
                }
            }
            /* istanbul ignore next */
            return -1;
        };
        /**
         * Returns the address from the index provided.
         *
         * @param idx The index of the address.
         *
         * @returns Returns the string representing the address.
         */
        this.getAddress = (idx) => {
            if (idx < this.addresses.length) {
                return this.addresses[`${idx}`].toBuffer();
            }
            throw new errors_1.AddressIndexError("Error - Output.getAddress: idx out of range");
        };
        /**
         * Given an array of address {@link https://github.com/feross/buffer|Buffer}s and an optional timestamp, returns true if the addresses meet the threshold required to spend the output.
         */
        this.meetsThreshold = (addresses, asOf = undefined) => {
            let now;
            if (typeof asOf === "undefined") {
                now = (0, helperfunctions_1.UnixNow)();
            }
            else {
                now = asOf;
            }
            const qualified = this.getSpenders(addresses, now);
            const threshold = this.threshold.readUInt32BE(0);
            if (qualified.length >= threshold) {
                return true;
            }
            return false;
        };
        /**
         * Given an array of addresses and an optional timestamp, select an array of address {@link https://github.com/feross/buffer|Buffer}s of qualified spenders for the output.
         */
        this.getSpenders = (addresses, asOf = undefined) => {
            const qualified = [];
            let now;
            if (typeof asOf === "undefined") {
                now = (0, helperfunctions_1.UnixNow)();
            }
            else {
                now = asOf;
            }
            const locktime = bintools.fromBufferToBN(this.locktime);
            if (now.lte(locktime)) {
                // not unlocked, not spendable
                return qualified;
            }
            const threshold = this.threshold.readUInt32BE(0);
            for (let i = 0; i < this.addresses.length && qualified.length < threshold; i++) {
                for (let j = 0; j < addresses.length && qualified.length < threshold; j++) {
                    if (addresses[`${j}`].toString("hex") ===
                        this.addresses[`${i}`].toBuffer().toString("hex")) {
                        qualified.push(addresses[`${j}`]);
                    }
                }
            }
            return qualified;
        };
        if (typeof addresses !== "undefined" && addresses.length) {
            const addrs = [];
            for (let i = 0; i < addresses.length; i++) {
                addrs[`${i}`] = new Address();
                addrs[`${i}`].fromBuffer(addresses[`${i}`]);
            }
            this.addresses = addrs;
            this.addresses.sort(Address.comparator());
            this.numaddrs.writeUInt32BE(this.addresses.length, 0);
        }
        if (typeof threshold !== undefined) {
            this.threshold.writeUInt32BE(threshold || 1, 0);
        }
        if (typeof locktime !== "undefined") {
            this.locktime = bintools.fromBNToBuffer(locktime, 8);
        }
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { locktime: serialization.encoder(this.locktime, encoding, "Buffer", "decimalString", 8), threshold: serialization.encoder(this.threshold, encoding, "Buffer", "decimalString", 4), addresses: this.addresses.map((a) => a.serialize(encoding)) });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.locktime = serialization.decoder(fields["locktime"], encoding, "decimalString", "Buffer", 8);
        this.threshold = serialization.decoder(fields["threshold"], encoding, "decimalString", "Buffer", 4);
        this.addresses = fields["addresses"].map((a) => {
            let addr = new Address();
            addr.deserialize(a, encoding);
            return addr;
        });
        this.numaddrs = buffer_1.Buffer.alloc(4);
        this.numaddrs.writeUInt32BE(this.addresses.length, 0);
    }
    /**
     * Returns a base-58 string representing the [[Output]].
     */
    fromBuffer(bytes, offset = 0) {
        this.locktime = bintools.copyFrom(bytes, offset, offset + 8);
        offset += 8;
        this.threshold = bintools.copyFrom(bytes, offset, offset + 4);
        offset += 4;
        this.numaddrs = bintools.copyFrom(bytes, offset, offset + 4);
        offset += 4;
        const numaddrs = this.numaddrs.readUInt32BE(0);
        this.addresses = [];
        for (let i = 0; i < numaddrs; i++) {
            const addr = new Address();
            offset = addr.fromBuffer(bytes, offset);
            this.addresses.push(addr);
        }
        this.addresses.sort(Address.comparator());
        return offset;
    }
    /**
     * Returns the buffer representing the [[Output]] instance.
     */
    toBuffer() {
        this.addresses.sort(Address.comparator());
        this.numaddrs.writeUInt32BE(this.addresses.length, 0);
        let bsize = this.locktime.length + this.threshold.length + this.numaddrs.length;
        const barr = [this.locktime, this.threshold, this.numaddrs];
        for (let i = 0; i < this.addresses.length; i++) {
            const b = this.addresses[`${i}`].toBuffer();
            barr.push(b);
            bsize += b.length;
        }
        return buffer_1.Buffer.concat(barr, bsize);
    }
    /**
     * Returns a base-58 string representing the [[Output]].
     */
    toString() {
        return bintools.bufferToB58(this.toBuffer());
    }
}
exports.OutputOwners = OutputOwners;
OutputOwners.comparator = () => (a, b) => {
    const aoutid = buffer_1.Buffer.alloc(4);
    aoutid.writeUInt32BE(a.getOutputID(), 0);
    const abuff = a.toBuffer();
    const boutid = buffer_1.Buffer.alloc(4);
    boutid.writeUInt32BE(b.getOutputID(), 0);
    const bbuff = b.toBuffer();
    const asort = buffer_1.Buffer.concat([aoutid, abuff], aoutid.length + abuff.length);
    const bsort = buffer_1.Buffer.concat([boutid, bbuff], boutid.length + bbuff.length);
    return buffer_1.Buffer.compare(asort, bsort);
};
class Output extends OutputOwners {
    constructor() {
        super(...arguments);
        this._typeName = "Output";
        this._typeID = undefined;
    }
}
exports.Output = Output;
class StandardParseableOutput extends serialization_1.Serializable {
    /**
     * Class representing an [[ParseableOutput]] for a transaction.
     *
     * @param output A number representing the InputID of the [[ParseableOutput]]
     */
    constructor(output = undefined) {
        super();
        this._typeName = "StandardParseableOutput";
        this._typeID = undefined;
        this.getOutput = () => this.output;
        if (output instanceof Output) {
            this.output = output;
        }
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { output: this.output.serialize(encoding) });
    }
    toBuffer() {
        const outbuff = this.output.toBuffer();
        const outid = buffer_1.Buffer.alloc(4);
        outid.writeUInt32BE(this.output.getOutputID(), 0);
        const barr = [outid, outbuff];
        return buffer_1.Buffer.concat(barr, outid.length + outbuff.length);
    }
}
exports.StandardParseableOutput = StandardParseableOutput;
/**
 * Returns a function used to sort an array of [[ParseableOutput]]s
 */
StandardParseableOutput.comparator = () => (a, b) => {
    const sorta = a.toBuffer();
    const sortb = b.toBuffer();
    return buffer_1.Buffer.compare(sorta, sortb);
};
class StandardTransferableOutput extends StandardParseableOutput {
    /**
     * Class representing an [[StandardTransferableOutput]] for a transaction.
     *
     * @param assetID A {@link https://github.com/feross/buffer|Buffer} representing the assetID of the [[Output]]
     * @param output A number representing the InputID of the [[StandardTransferableOutput]]
     */
    constructor(assetID = undefined, output = undefined) {
        super(output);
        this._typeName = "StandardTransferableOutput";
        this._typeID = undefined;
        this.assetID = undefined;
        this.getAssetID = () => this.assetID;
        if (typeof assetID !== "undefined") {
            this.assetID = assetID;
        }
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { assetID: serialization.encoder(this.assetID, encoding, "Buffer", "cb58") });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.assetID = serialization.decoder(fields["assetID"], encoding, "cb58", "Buffer", 32);
    }
    toBuffer() {
        const parseableBuff = super.toBuffer();
        const barr = [this.assetID, parseableBuff];
        return buffer_1.Buffer.concat(barr, this.assetID.length + parseableBuff.length);
    }
}
exports.StandardTransferableOutput = StandardTransferableOutput;
/**
 * An [[Output]] class which specifies a token amount .
 */
class StandardAmountOutput extends Output {
    /**
     * A [[StandardAmountOutput]] class which issues a payment on an assetID.
     *
     * @param amount A {@link https://github.com/indutny/bn.js/|BN} representing the amount in the output
     * @param addresses An array of {@link https://github.com/feross/buffer|Buffer}s representing addresses
     * @param locktime A {@link https://github.com/indutny/bn.js/|BN} representing the locktime
     * @param threshold A number representing the the threshold number of signers required to sign the transaction
     */
    constructor(amount = undefined, addresses = undefined, locktime = undefined, threshold = undefined) {
        super(addresses, locktime, threshold);
        this._typeName = "StandardAmountOutput";
        this._typeID = undefined;
        this.amount = buffer_1.Buffer.alloc(8);
        this.amountValue = new bn_js_1.default(0);
        if (typeof amount !== "undefined") {
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
     * Returns the amount as a {@link https://github.com/indutny/bn.js/|BN}.
     */
    getAmount() {
        return this.amountValue.clone();
    }
    /**
     * Popuates the instance from a {@link https://github.com/feross/buffer|Buffer} representing the [[StandardAmountOutput]] and returns the size of the output.
     */
    fromBuffer(outbuff, offset = 0) {
        this.amount = bintools.copyFrom(outbuff, offset, offset + 8);
        this.amountValue = bintools.fromBufferToBN(this.amount);
        offset += 8;
        return super.fromBuffer(outbuff, offset);
    }
    /**
     * Returns the buffer representing the [[StandardAmountOutput]] instance.
     */
    toBuffer() {
        const superbuff = super.toBuffer();
        const bsize = this.amount.length + superbuff.length;
        this.numaddrs.writeUInt32BE(this.addresses.length, 0);
        const barr = [this.amount, superbuff];
        return buffer_1.Buffer.concat(barr, bsize);
    }
}
exports.StandardAmountOutput = StandardAmountOutput;
/**
 * An [[Output]] class which specifies an NFT.
 */
class BaseNFTOutput extends Output {
    constructor() {
        super(...arguments);
        this._typeName = "BaseNFTOutput";
        this._typeID = undefined;
        this.groupID = buffer_1.Buffer.alloc(4);
        /**
         * Returns the groupID as a number.
         */
        this.getGroupID = () => {
            return this.groupID.readUInt32BE(0);
        };
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { groupID: serialization.encoder(this.groupID, encoding, "Buffer", "decimalString", 4) });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.groupID = serialization.decoder(fields["groupID"], encoding, "decimalString", "Buffer", 4);
    }
}
exports.BaseNFTOutput = BaseNFTOutput;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoib3V0cHV0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL2NvbW1vbi9vdXRwdXQudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IjtBQUFBOzs7R0FHRzs7Ozs7O0FBRUgsb0NBQWdDO0FBQ2hDLGtEQUFzQjtBQUN0QixpRUFBd0M7QUFDeEMscUNBQWlDO0FBQ2pDLDhEQUFrRDtBQUNsRCwwREFJK0I7QUFDL0IsNENBQWdGO0FBRWhGOztHQUVHO0FBQ0gsTUFBTSxRQUFRLEdBQWEsa0JBQVEsQ0FBQyxXQUFXLEVBQUUsQ0FBQTtBQUNqRCxNQUFNLGFBQWEsR0FBa0IsNkJBQWEsQ0FBQyxXQUFXLEVBQUUsQ0FBQTtBQUVoRTs7R0FFRztBQUNILE1BQWEsT0FBUSxTQUFRLGVBQU07SUFpRWpDOztPQUVHO0lBQ0g7UUFDRSxLQUFLLEVBQUUsQ0FBQTtRQXBFQyxjQUFTLEdBQUcsU0FBUyxDQUFBO1FBQ3JCLFlBQU8sR0FBRyxTQUFTLENBQUE7UUFFN0IsOENBQThDO1FBRXBDLFVBQUssR0FBRyxlQUFNLENBQUMsS0FBSyxDQUFDLEVBQUUsQ0FBQyxDQUFBO1FBQ3hCLFVBQUssR0FBRyxFQUFFLENBQUE7SUErRHBCLENBQUM7SUFyREQ7O09BRUc7SUFDSCxRQUFRO1FBQ04sT0FBTyxRQUFRLENBQUMsVUFBVSxDQUFDLElBQUksQ0FBQyxRQUFRLEVBQUUsQ0FBQyxDQUFBO0lBQzdDLENBQUM7SUFFRDs7Ozs7O09BTUc7SUFDSCxVQUFVLENBQUMsSUFBWTtRQUNyQixNQUFNLFFBQVEsR0FBVyxRQUFRLENBQUMsV0FBVyxDQUFDLElBQUksQ0FBQyxDQUFBO1FBQ25ELElBQUksUUFBUSxDQUFDLE1BQU0sS0FBSyxFQUFFLElBQUksUUFBUSxDQUFDLGdCQUFnQixDQUFDLFFBQVEsQ0FBQyxFQUFFO1lBQ2pFLE1BQU0sT0FBTyxHQUFXLFFBQVEsQ0FBQyxRQUFRLENBQ3ZDLFFBQVEsRUFDUixDQUFDLEVBQ0QsUUFBUSxDQUFDLE1BQU0sR0FB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QUFjLFNBQVEsTUFBTTtJQUFsRDs7UUFDWSxjQUFTLEdBQUcsZUFBZSxDQUFBO1FBQzNCLFlBQU8sR0FBRyxTQUFTLENBQUE7UUEwQm5CLFlBQU8sR0FBVyxlQUFNLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxDQUFBO1FBRTNDOztXQUVHO1FBQ0gsZUFBVSxHQUFHLEdBQVcsRUFBRTtZQUN4QixPQUFPLElBQUksQ0FBQyxPQUFPLENBQUMsWUFBWSxDQUFDLENBQUMsQ0FBQyxDQUFBO1FBQ3JDLENBQUMsQ0FBQTtJQUNILENBQUM7SUFoQ0MsU0FBUyxDQUFDLFdBQStCLEtBQUs7UUFDNUMsSUFBSSxNQUFNLEdBQVcsS0FBSyxDQUFDLFNBQVMsQ0FBQyxRQUFRLENBQUMsQ0FBQTtRQUM5Qyx1Q0FDSyxNQUFNLEtBQ1QsT0FBTyxFQUFFLGFBQWEsQ0FBQyxPQUFPLENBQzVCLElBQUksQ0FBQyxPQUFPLEVBQ1osUUFBUSxFQUNSLFFBQVEsRUFDUixlQUFlLEVBQ2YsQ0FBQyxDQUNGLElBQ0Y7SUFDSCxDQUFDO0lBQ0QsV0FBVyxDQUFDLE1BQWMsRUFBRSxXQUErQixLQUFLO1FBQzlELEtBQUssQ0FBQyxXQUFXLENBQUMsTUFBTSxFQUFFLFFBQVEsQ0FBQyxDQUFBO1FBQ25DLElBQUksQ0FBQyxPQUFPLEdBQUcsYUFBYSxDQUFDLE9BQU8sQ0FDbEMsTUFBTSxDQUFDLFNBQVMsQ0FBQyxFQUNqQixRQUFRLEVBQ1IsZUFBZSxFQUNmLFFBQVEsRUFDUixDQUFDLENBQ0YsQ0FBQTtJQUNILENBQUM7Q0FVRjtBQXBDRCxzQ0FvQ0MiLCJzb3VyY2VzQ29udGVudCI6WyIvKipcbiAqIEBwYWNrYWdlRG9jdW1lbnRhdGlvblxuICogQG1vZHVsZSBDb21tb24tT3V0cHV0XG4gKi9cblxuaW1wb3J0IHsgQnVmZmVyIH0gZnJvbSBcImJ1ZmZlci9cIlxuaW1wb3J0IEJOIGZyb20gXCJibi5qc1wiXG5pbXBvcnQgQmluVG9vbHMgZnJvbSBcIi4uL3V0aWxzL2JpbnRvb2xzXCJcbmltcG9ydCB7IE5CeXRlcyB9IGZyb20gXCIuL25ieXRlc1wiXG5pbXBvcnQgeyBVbml4Tm93IH0gZnJvbSBcIi4uL3V0aWxzL2hlbHBlcmZ1bmN0aW9uc1wiXG5pbXBvcnQge1xuICBTZXJpYWxpemFibGUsXG4gIFNlcmlhbGl6YXRpb24sXG4gIFNlcmlhbGl6ZWRFbmNvZGluZ1xufSBmcm9tIFwiLi4vdXRpbHMvc2VyaWFsaXphdGlvblwiXG5pbXBvcnQgeyBDaGVja3N1bUVycm9yLCBBZGRyZXNzRXJyb3IsIEFkZHJlc3NJbmRleEVycm9yIH0gZnJvbSBcIi4uL3V0aWxzL2Vycm9yc1wiXG5cbi8qKlxuICogQGlnbm9yZVxuICovXG5jb25zdCBiaW50b29sczogQmluVG9vbHMgPSBCaW5Ub29scy5nZXRJbnN0YW5jZSgpXG5jb25zdCBzZXJpYWxpemF0aW9uOiBTZXJpYWxpemF0aW9uID0gU2VyaWFsaXphdGlvbi5nZXRJbnN0YW5jZSgpXG5cbi8qKlxuICogQ2xhc3MgZm9yIHJlcHJlc2VudGluZyBhbiBhZGRyZXNzIHVzZWQgaW4gW1tPdXRwdXRdXSB0eXBlc1xuICovXG5leHBvcnQgY2xhc3MgQWRkcmVzcyBleHRlbmRzIE5CeXRlcyB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIkFkZHJlc3NcIlxuICBwcm90ZWN0ZWQgX3R5cGVJRCA9IHVuZGVmaW5lZFxuXG4gIC8vc2VyaWFsaXplIGFuZCBkZXNlcmlhbGl6ZSBib3RoIGFyZSBpbmhlcml0ZWRcblxuICBwcm90ZWN0ZWQgYnl0ZXMgPSBCdWZmZXIuYWxsb2MoMjApXG4gIHByb3RlY3RlZCBic2l6ZSA9IDIwXG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSBmdW5jdGlvbiB1c2VkIHRvIHNvcnQgYW4gYXJyYXkgb2YgW1tBZGRyZXNzXV1lc1xuICAgKi9cbiAgc3RhdGljIGNvbXBhcmF0b3IgPVxuICAgICgpOiAoKGE6IEFkZHJlc3MsIGI6IEFkZHJlc3MpID0+IDEgfCAtMSB8IDApID0+XG4gICAgKGE6IEFkZHJlc3MsIGI6IEFkZHJlc3MpOiAxIHwgLTEgfCAwID0+XG4gICAgICBCdWZmZXIuY29tcGFyZShhLnRvQnVmZmVyKCksIGIudG9CdWZmZXIoKSkgYXMgMSB8IC0xIHwgMFxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIGEgYmFzZS01OCByZXByZXNlbnRhdGlvbiBvZiB0aGUgW1tBZGRyZXNzXV0uXG4gICAqL1xuICB0b1N0cmluZygpOiBzdHJpbmcge1xuICAgIHJldHVybiBiaW50b29scy5jYjU4RW5jb2RlKHRoaXMudG9CdWZmZXIoKSlcbiAgfVxuXG4gIC8qKlxuICAgKiBUYWtlcyBhIGJhc2UtNTggc3RyaW5nIGNvbnRhaW5pbmcgYW4gW1tBZGRyZXNzXV0sIHBhcnNlcyBpdCwgcG9wdWxhdGVzIHRoZSBjbGFzcywgYW5kIHJldHVybnMgdGhlIGxlbmd0aCBvZiB0aGUgQWRkcmVzcyBpbiBieXRlcy5cbiAgICpcbiAgICogQHBhcmFtIGJ5dGVzIEEgYmFzZS01OCBzdHJpbmcgY29udGFpbmluZyBhIHJhdyBbW0FkZHJlc3NdXVxuICAgKlxuICAgKiBAcmV0dXJucyBUaGUgbGVuZ3RoIG9mIHRoZSByYXcgW1tBZGRyZXNzXV1cbiAgICovXG4gIGZyb21TdHJpbmcoYWRkcjogc3RyaW5nKTogbnVtYmVyIHtcbiAgICBjb25zdCBhZGRyYnVmZjogQnVmZmVyID0gYmludG9vbHMuYjU4VG9CdWZmZXIoYWRkcilcbiAgICBpZiAoYWRkcmJ1ZmYubGVuZ3RoID09PSAyNCAmJiBiaW50b29scy52YWxpZGF0ZUNoZWNrc3VtKGFkZHJidWZmKSkge1xuICAgICAgY29uc3QgbmV3YnVmZjogQnVmZmVyID0gYmludG9vbHMuY29weUZyb20oXG4gICAgICAgIGFkZHJidWZmLFxuICAgICAgICAwLFxuICAgICAgICBhZGRyYnVmZi5sZW5ndGggLSA0XG4gICAgICApXG4gICAgICBpZiAobmV3YnVmZi5sZW5ndGggPT09IDIwKSB7XG4gICAgICAgIHRoaXMuYnl0ZXMgPSBuZXdidWZmXG4gICAgICB9XG4gICAgfSBlbHNlIGlmIChhZGRyYnVmZi5sZW5ndGggPT09IDI0KSB7XG4gICAgICB0aHJvdyBuZXcgQ2hlY2tzdW1FcnJvcihcbiAgICAgICAgXCJFcnJvciAtIEFkZHJlc3MuZnJvbVN0cmluZzogaW52YWxpZCBjaGVja3N1bSBvbiBhZGRyZXNzXCJcbiAgICAgIClcbiAgICB9IGVsc2UgaWYgKGFkZHJidWZmLmxlbmd0aCA9PT0gMjApIHtcbiAgICAgIHRoaXMuYnl0ZXMgPSBhZGRyYnVmZlxuICAgIH0gZWxzZSB7XG4gICAgICAvKiBpc3RhbmJ1bCBpZ25vcmUgbmV4dCAqL1xuICAgICAgdGhyb3cgbmV3IEFkZHJlc3NFcnJvcihcIkVycm9yIC0gQWRkcmVzcy5mcm9tU3RyaW5nOiBpbnZhbGlkIGFkZHJlc3NcIilcbiAgICB9XG4gICAgcmV0dXJuIHRoaXMuZ2V0U2l6ZSgpXG4gIH1cblxuICBjbG9uZSgpOiB0aGlzIHtcbiAgICBsZXQgbmV3YmFzZTogQWRkcmVzcyA9IG5ldyBBZGRyZXNzKClcbiAgICBuZXdiYXNlLmZyb21CdWZmZXIodGhpcy50b0J1ZmZlcigpKVxuICAgIHJldHVybiBuZXdiYXNlIGFzIHRoaXNcbiAgfVxuXG4gIGNyZWF0ZSguLi5hcmdzOiBhbnlbXSk6IHRoaXMge1xuICAgIHJldHVybiBuZXcgQWRkcmVzcygpIGFzIHRoaXNcbiAgfVxuXG4gIC8qKlxuICAgKiBDbGFzcyBmb3IgcmVwcmVzZW50aW5nIGFuIGFkZHJlc3MgdXNlZCBpbiBbW091dHB1dF1dIHR5cGVzXG4gICAqL1xuICBjb25zdHJ1Y3RvcigpIHtcbiAgICBzdXBlcigpXG4gIH1cbn1cblxuLyoqXG4gKiBEZWZpbmVzIHRoZSBtb3N0IGJhc2ljIHZhbHVlcyBmb3Igb3V0cHV0IG93bmVyc2hpcC4gTW9zdGx5IGluaGVyaXRlZCBmcm9tLCBidXQgY2FuIGJlIHVzZWQgaW4gcG9wdWxhdGlvbiBvZiBORlQgT3duZXIgZGF0YS5cbiAqL1xuZXhwb3J0IGNsYXNzIE91dHB1dE93bmVycyBleHRlbmRzIFNlcmlhbGl6YWJsZSB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIk91dHB1dE93bmVyc1wiXG4gIHByb3RlY3RlZCBfdHlwZUlEID0gdW5kZWZpbmVkXG5cbiAgc2VyaWFsaXplKGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKTogb2JqZWN0IHtcbiAgICBsZXQgZmllbGRzOiBvYmplY3QgPSBzdXBlci5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgcmV0dXJuIHtcbiAgICAgIC4uLmZpZWxkcyxcbiAgICAgIGxvY2t0aW1lOiBzZXJpYWxpemF0aW9uLmVuY29kZXIoXG4gICAgICAgIHRoaXMubG9ja3RpbWUsXG4gICAgICAgIGVuY29kaW5nLFxuICAgICAgICBcIkJ1ZmZlclwiLFxuICAgICAgICBcImRlY2ltYWxTdHJpbmdcIixcbiAgICAgICAgOFxuICAgICAgKSxcbiAgICAgIHRocmVzaG9sZDogc2VyaWFsaXphdGlvbi5lbmNvZGVyKFxuICAgICAgICB0aGlzLnRocmVzaG9sZCxcbiAgICAgICAgZW5jb2RpbmcsXG4gICAgICAgIFwiQnVmZmVyXCIsXG4gICAgICAgIFwiZGVjaW1hbFN0cmluZ1wiLFxuICAgICAgICA0XG4gICAgICApLFxuICAgICAgYWRkcmVzc2VzOiB0aGlzLmFkZHJlc3Nlcy5tYXAoKGE6IEFkZHJlc3MpOiBvYmplY3QgPT5cbiAgICAgICAgYS5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgICApXG4gICAgfVxuICB9XG4gIGRlc2VyaWFsaXplKGZpZWxkczogb2JqZWN0LCBlbmNvZGluZzogU2VyaWFsaXplZEVuY29kaW5nID0gXCJoZXhcIikge1xuICAgIHN1cGVyLmRlc2VyaWFsaXplKGZpZWxkcywgZW5jb2RpbmcpXG4gICAgdGhpcy5sb2NrdGltZSA9IHNlcmlhbGl6YXRpb24uZGVjb2RlcihcbiAgICAgIGZpZWxkc1tcImxvY2t0aW1lXCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBcImRlY2ltYWxTdHJpbmdcIixcbiAgICAgIFwiQnVmZmVyXCIsXG4gICAgICA4XG4gICAgKVxuICAgIHRoaXMudGhyZXNob2xkID0gc2VyaWFsaXphdGlvbi5kZWNvZGVyKFxuICAgICAgZmllbGRzW1widGhyZXNob2xkXCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBcImRlY2ltYWxTdHJpbmdcIixcbiAgICAgIFwiQnVmZmVyXCIsXG4gICAgICA0XG4gICAgKVxuICAgIHRoaXMuYWRkcmVzc2VzID0gZmllbGRzW1wiYWRkcmVzc2VzXCJdLm1hcCgoYTogb2JqZWN0KSA9PiB7XG4gICAgICBsZXQgYWRkcjogQWRkcmVzcyA9IG5ldyBBZGRyZXNzKClcbiAgICAgIGFkZHIuZGVzZXJpYWxpemUoYSwgZW5jb2RpbmcpXG4gICAgICByZXR1cm4gYWRkclxuICAgIH0pXG4gICAgdGhpcy5udW1hZGRycyA9IEJ1ZmZlci5hbGxvYyg0KVxuICAgIHRoaXMubnVtYWRkcnMud3JpdGVVSW50MzJCRSh0aGlzLmFkZHJlc3Nlcy5sZW5ndGgsIDApXG4gIH1cblxuICBwcm90ZWN0ZWQgbG9ja3RpbWU6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYyg4KVxuICBwcm90ZWN0ZWQgdGhyZXNob2xkOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoNClcbiAgcHJvdGVjdGVkIG51bWFkZHJzOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoNClcbiAgcHJvdGVjdGVkIGFkZHJlc3NlczogQWRkcmVzc1tdID0gW11cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgdGhyZXNob2xkIG9mIHNpZ25lcnMgcmVxdWlyZWQgdG8gc3BlbmQgdGhpcyBvdXRwdXQuXG4gICAqL1xuICBnZXRUaHJlc2hvbGQgPSAoKTogbnVtYmVyID0+IHRoaXMudGhyZXNob2xkLnJlYWRVSW50MzJCRSgwKVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vaW5kdXRueS9ibi5qcy98Qk59IHJlcGVyc2VudGluZyB0aGUgVU5JWCBUaW1lc3RhbXAgd2hlbiB0aGUgbG9jayBpcyBtYWRlIGF2YWlsYWJsZS5cbiAgICovXG4gIGdldExvY2t0aW1lID0gKCk6IEJOID0+IGJpbnRvb2xzLmZyb21CdWZmZXJUb0JOKHRoaXMubG9ja3RpbWUpXG5cbiAgLyoqXG4gICAqIFJldHVybnMgYW4gYXJyYXkgb2Yge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn1zIGZvciB0aGUgYWRkcmVzc2VzLlxuICAgKi9cbiAgZ2V0QWRkcmVzc2VzID0gKCk6IEJ1ZmZlcltdID0+IHtcbiAgICBjb25zdCByZXN1bHQ6IEJ1ZmZlcltdID0gW11cbiAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgdGhpcy5hZGRyZXNzZXMubGVuZ3RoOyBpKyspIHtcbiAgICAgIHJlc3VsdC5wdXNoKHRoaXMuYWRkcmVzc2VzW2Ake2l9YF0udG9CdWZmZXIoKSlcbiAgICB9XG4gICAgcmV0dXJuIHJlc3VsdFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIGluZGV4IG9mIHRoZSBhZGRyZXNzLlxuICAgKlxuICAgKiBAcGFyYW0gYWRkcmVzcyBBIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IG9mIHRoZSBhZGRyZXNzIHRvIGxvb2sgdXAgdG8gcmV0dXJuIGl0cyBpbmRleC5cbiAgICpcbiAgICogQHJldHVybnMgVGhlIGluZGV4IG9mIHRoZSBhZGRyZXNzLlxuICAgKi9cbiAgZ2V0QWRkcmVzc0lkeCA9IChhZGRyZXNzOiBCdWZmZXIpOiBudW1iZXIgPT4ge1xuICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCB0aGlzLmFkZHJlc3Nlcy5sZW5ndGg7IGkrKykge1xuICAgICAgaWYgKFxuICAgICAgICB0aGlzLmFkZHJlc3Nlc1tgJHtpfWBdLnRvQnVmZmVyKCkudG9TdHJpbmcoXCJoZXhcIikgPT09XG4gICAgICAgIGFkZHJlc3MudG9TdHJpbmcoXCJoZXhcIilcbiAgICAgICkge1xuICAgICAgICByZXR1cm4gaVxuICAgICAgfVxuICAgIH1cbiAgICAvKiBpc3RhbmJ1bCBpZ25vcmUgbmV4dCAqL1xuICAgIHJldHVybiAtMVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIGFkZHJlc3MgZnJvbSB0aGUgaW5kZXggcHJvdmlkZWQuXG4gICAqXG4gICAqIEBwYXJhbSBpZHggVGhlIGluZGV4IG9mIHRoZSBhZGRyZXNzLlxuICAgKlxuICAgKiBAcmV0dXJucyBSZXR1cm5zIHRoZSBzdHJpbmcgcmVwcmVzZW50aW5nIHRoZSBhZGRyZXNzLlxuICAgKi9cbiAgZ2V0QWRkcmVzcyA9IChpZHg6IG51bWJlcik6IEJ1ZmZlciA9PiB7XG4gICAgaWYgKGlkeCA8IHRoaXMuYWRkcmVzc2VzLmxlbmd0aCkge1xuICAgICAgcmV0dXJuIHRoaXMuYWRkcmVzc2VzW2Ake2lkeH1gXS50b0J1ZmZlcigpXG4gICAgfVxuICAgIHRocm93IG5ldyBBZGRyZXNzSW5kZXhFcnJvcihcIkVycm9yIC0gT3V0cHV0LmdldEFkZHJlc3M6IGlkeCBvdXQgb2YgcmFuZ2VcIilcbiAgfVxuXG4gIC8qKlxuICAgKiBHaXZlbiBhbiBhcnJheSBvZiBhZGRyZXNzIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9cyBhbmQgYW4gb3B0aW9uYWwgdGltZXN0YW1wLCByZXR1cm5zIHRydWUgaWYgdGhlIGFkZHJlc3NlcyBtZWV0IHRoZSB0aHJlc2hvbGQgcmVxdWlyZWQgdG8gc3BlbmQgdGhlIG91dHB1dC5cbiAgICovXG4gIG1lZXRzVGhyZXNob2xkID0gKGFkZHJlc3NlczogQnVmZmVyW10sIGFzT2Y6IEJOID0gdW5kZWZpbmVkKTogYm9vbGVhbiA9PiB7XG4gICAgbGV0IG5vdzogQk5cbiAgICBpZiAodHlwZW9mIGFzT2YgPT09IFwidW5kZWZpbmVkXCIpIHtcbiAgICAgIG5vdyA9IFVuaXhOb3coKVxuICAgIH0gZWxzZSB7XG4gICAgICBub3cgPSBhc09mXG4gICAgfVxuICAgIGNvbnN0IHF1YWxpZmllZDogQnVmZmVyW10gPSB0aGlzLmdldFNwZW5kZXJzKGFkZHJlc3Nlcywgbm93KVxuICAgIGNvbnN0IHRocmVzaG9sZDogbnVtYmVyID0gdGhpcy50aHJlc2hvbGQucmVhZFVJbnQzMkJFKDApXG4gICAgaWYgKHF1YWxpZmllZC5sZW5ndGggPj0gdGhyZXNob2xkKSB7XG4gICAgICByZXR1cm4gdHJ1ZVxuICAgIH1cblxuICAgIHJldHVybiBmYWxzZVxuICB9XG5cbiAgLyoqXG4gICAqIEdpdmVuIGFuIGFycmF5IG9mIGFkZHJlc3NlcyBhbmQgYW4gb3B0aW9uYWwgdGltZXN0YW1wLCBzZWxlY3QgYW4gYXJyYXkgb2YgYWRkcmVzcyB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfXMgb2YgcXVhbGlmaWVkIHNwZW5kZXJzIGZvciB0aGUgb3V0cHV0LlxuICAgKi9cbiAgZ2V0U3BlbmRlcnMgPSAoYWRkcmVzc2VzOiBCdWZmZXJbXSwgYXNPZjogQk4gPSB1bmRlZmluZWQpOiBCdWZmZXJbXSA9PiB7XG4gICAgY29uc3QgcXVhbGlmaWVkOiBCdWZmZXJbXSA9IFtdXG4gICAgbGV0IG5vdzogQk5cbiAgICBpZiAodHlwZW9mIGFzT2YgPT09IFwidW5kZWZpbmVkXCIpIHtcbiAgICAgIG5vdyA9IFVuaXhOb3coKVxuICAgIH0gZWxzZSB7XG4gICAgICBub3cgPSBhc09mXG4gICAgfVxuICAgIGNvbnN0IGxvY2t0aW1lOiBCTiA9IGJpbnRvb2xzLmZyb21CdWZmZXJUb0JOKHRoaXMubG9ja3RpbWUpXG4gICAgaWYgKG5vdy5sdGUobG9ja3RpbWUpKSB7XG4gICAgICAvLyBub3QgdW5sb2NrZWQsIG5vdCBzcGVuZGFibGVcbiAgICAgIHJldHVybiBxdWFsaWZpZWRcbiAgICB9XG5cbiAgICBjb25zdCB0aHJlc2hvbGQ6IG51bWJlciA9IHRoaXMudGhyZXNob2xkLnJlYWRVSW50MzJCRSgwKVxuICAgIGZvciAoXG4gICAgICBsZXQgaTogbnVtYmVyID0gMDtcbiAgICAgIGkgPCB0aGlzLmFkZHJlc3Nlcy5sZW5ndGggJiYgcXVhbGlmaWVkLmxlbmd0aCA8IHRocmVzaG9sZDtcbiAgICAgIGkrK1xuICAgICkge1xuICAgICAgZm9yIChcbiAgICAgICAgbGV0IGo6IG51bWJlciA9IDA7XG4gICAgICAgIGogPCBhZGRyZXNzZXMubGVuZ3RoICYmIHF1YWxpZmllZC5sZW5ndGggPCB0aHJlc2hvbGQ7XG4gICAgICAgIGorK1xuICAgICAgKSB7XG4gICAgICAgIGlmIChcbiAgICAgICAgICBhZGRyZXNzZXNbYCR7an1gXS50b1N0cmluZyhcImhleFwiKSA9PT1cbiAgICAgICAgICB0aGlzLmFkZHJlc3Nlc1tgJHtpfWBdLnRvQnVmZmVyKCkudG9TdHJpbmcoXCJoZXhcIilcbiAgICAgICAgKSB7XG4gICAgICAgICAgcXVhbGlmaWVkLnB1c2goYWRkcmVzc2VzW2Ake2p9YF0pXG4gICAgICAgIH1cbiAgICAgIH1cbiAgICB9XG5cbiAgICByZXR1cm4gcXVhbGlmaWVkXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyBhIGJhc2UtNTggc3RyaW5nIHJlcHJlc2VudGluZyB0aGUgW1tPdXRwdXRdXS5cbiAgICovXG4gIGZyb21CdWZmZXIoYnl0ZXM6IEJ1ZmZlciwgb2Zmc2V0OiBudW1iZXIgPSAwKTogbnVtYmVyIHtcbiAgICB0aGlzLmxvY2t0aW1lID0gYmludG9vbHMuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgOClcbiAgICBvZmZzZXQgKz0gOFxuICAgIHRoaXMudGhyZXNob2xkID0gYmludG9vbHMuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgNClcbiAgICBvZmZzZXQgKz0gNFxuICAgIHRoaXMubnVtYWRkcnMgPSBiaW50b29scy5jb3B5RnJvbShieXRlcywgb2Zmc2V0LCBvZmZzZXQgKyA0KVxuICAgIG9mZnNldCArPSA0XG4gICAgY29uc3QgbnVtYWRkcnM6IG51bWJlciA9IHRoaXMubnVtYWRkcnMucmVhZFVJbnQzMkJFKDApXG4gICAgdGhpcy5hZGRyZXNzZXMgPSBbXVxuICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCBudW1hZGRyczsgaSsrKSB7XG4gICAgICBjb25zdCBhZGRyOiBBZGRyZXNzID0gbmV3IEFkZHJlc3MoKVxuICAgICAgb2Zmc2V0ID0gYWRkci5mcm9tQnVmZmVyKGJ5dGVzLCBvZmZzZXQpXG4gICAgICB0aGlzLmFkZHJlc3Nlcy5wdXNoKGFkZHIpXG4gICAgfVxuICAgIHRoaXMuYWRkcmVzc2VzLnNvcnQoQWRkcmVzcy5jb21wYXJhdG9yKCkpXG4gICAgcmV0dXJuIG9mZnNldFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIGJ1ZmZlciByZXByZXNlbnRpbmcgdGhlIFtbT3V0cHV0XV0gaW5zdGFuY2UuXG4gICAqL1xuICB0b0J1ZmZlcigpOiBCdWZmZXIge1xuICAgIHRoaXMuYWRkcmVzc2VzLnNvcnQoQWRkcmVzcy5jb21wYXJhdG9yKCkpXG4gICAgdGhpcy5udW1hZGRycy53cml0ZVVJbnQzMkJFKHRoaXMuYWRkcmVzc2VzLmxlbmd0aCwgMClcbiAgICBsZXQgYnNpemU6IG51bWJlciA9XG4gICAgICB0aGlzLmxvY2t0aW1lLmxlbmd0aCArIHRoaXMudGhyZXNob2xkLmxlbmd0aCArIHRoaXMubnVtYWRkcnMubGVuZ3RoXG4gICAgY29uc3QgYmFycjogQnVmZmVyW10gPSBbdGhpcy5sb2NrdGltZSwgdGhpcy50aHJlc2hvbGQsIHRoaXMubnVtYWRkcnNdXG4gICAgZm9yIChsZXQgaTogbnVtYmVyID0gMDsgaSA8IHRoaXMuYWRkcmVzc2VzLmxlbmd0aDsgaSsrKSB7XG4gICAgICBjb25zdCBiOiBCdWZmZXIgPSB0aGlzLmFkZHJlc3Nlc1tgJHtpfWBdLnRvQnVmZmVyKClcbiAgICAgIGJhcnIucHVzaChiKVxuICAgICAgYnNpemUgKz0gYi5sZW5ndGhcbiAgICB9XG4gICAgcmV0dXJuIEJ1ZmZlci5jb25jYXQoYmFyciwgYnNpemUpXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyBhIGJhc2UtNTggc3RyaW5nIHJlcHJlc2VudGluZyB0aGUgW1tPdXRwdXRdXS5cbiAgICovXG4gIHRvU3RyaW5nKCk6IHN0cmluZyB7XG4gICAgcmV0dXJuIGJpbnRvb2xzLmJ1ZmZlclRvQjU4KHRoaXMudG9CdWZmZXIoKSlcbiAgfVxuXG4gIHN0YXRpYyBjb21wYXJhdG9yID1cbiAgICAoKTogKChhOiBPdXRwdXQsIGI6IE91dHB1dCkgPT4gMSB8IC0xIHwgMCkgPT5cbiAgICAoYTogT3V0cHV0LCBiOiBPdXRwdXQpOiAxIHwgLTEgfCAwID0+IHtcbiAgICAgIGNvbnN0IGFvdXRpZDogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDQpXG4gICAgICBhb3V0aWQud3JpdGVVSW50MzJCRShhLmdldE91dHB1dElEKCksIDApXG4gICAgICBjb25zdCBhYnVmZjogQnVmZmVyID0gYS50b0J1ZmZlcigpXG5cbiAgICAgIGNvbnN0IGJvdXRpZDogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDQpXG4gICAgICBib3V0aWQud3JpdGVVSW50MzJCRShiLmdldE91dHB1dElEKCksIDApXG4gICAgICBjb25zdCBiYnVmZjogQnVmZmVyID0gYi50b0J1ZmZlcigpXG5cbiAgICAgIGNvbnN0IGFzb3J0OiBCdWZmZXIgPSBCdWZmZXIuY29uY2F0KFxuICAgICAgICBbYW91dGlkLCBhYnVmZl0sXG4gICAgICAgIGFvdXRpZC5sZW5ndGggKyBhYnVmZi5sZW5ndGhcbiAgICAgIClcbiAgICAgIGNvbnN0IGJzb3J0OiBCdWZmZXIgPSBCdWZmZXIuY29uY2F0KFxuICAgICAgICBbYm91dGlkLCBiYnVmZl0sXG4gICAgICAgIGJvdXRpZC5sZW5ndGggKyBiYnVmZi5sZW5ndGhcbiAgICAgIClcbiAgICAgIHJldHVybiBCdWZmZXIuY29tcGFyZShhc29ydCwgYnNvcnQpIGFzIDEgfCAtMSB8IDBcbiAgICB9XG5cbiAgLyoqXG4gICAqIEFuIFtbT3V0cHV0XV0gY2xhc3Mgd2hpY2ggY29udGFpbnMgYWRkcmVzc2VzLCBsb2NrdGltZXMsIGFuZCB0aHJlc2hvbGRzLlxuICAgKlxuICAgKiBAcGFyYW0gYWRkcmVzc2VzIEFuIGFycmF5IG9mIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9cyByZXByZXNlbnRpbmcgb3V0cHV0IG93bmVyJ3MgYWRkcmVzc2VzXG4gICAqIEBwYXJhbSBsb2NrdGltZSBBIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vaW5kdXRueS9ibi5qcy98Qk59IHJlcHJlc2VudGluZyB0aGUgbG9ja3RpbWVcbiAgICogQHBhcmFtIHRocmVzaG9sZCBBIG51bWJlciByZXByZXNlbnRpbmcgdGhlIHRoZSB0aHJlc2hvbGQgbnVtYmVyIG9mIHNpZ25lcnMgcmVxdWlyZWQgdG8gc2lnbiB0aGUgdHJhbnNhY3Rpb25cbiAgICovXG4gIGNvbnN0cnVjdG9yKFxuICAgIGFkZHJlc3NlczogQnVmZmVyW10gPSB1bmRlZmluZWQsXG4gICAgbG9ja3RpbWU6IEJOID0gdW5kZWZpbmVkLFxuICAgIHRocmVzaG9sZDogbnVtYmVyID0gdW5kZWZpbmVkXG4gICkge1xuICAgIHN1cGVyKClcbiAgICBpZiAodHlwZW9mIGFkZHJlc3NlcyAhPT0gXCJ1bmRlZmluZWRcIiAmJiBhZGRyZXNzZXMubGVuZ3RoKSB7XG4gICAgICBjb25zdCBhZGRyczogQWRkcmVzc1tdID0gW11cbiAgICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCBhZGRyZXNzZXMubGVuZ3RoOyBpKyspIHtcbiAgICAgICAgYWRkcnNbYCR7aX1gXSA9IG5ldyBBZGRyZXNzKClcbiAgICAgICAgYWRkcnNbYCR7aX1gXS5mcm9tQnVmZmVyKGFkZHJlc3Nlc1tgJHtpfWBdKVxuICAgICAgfVxuICAgICAgdGhpcy5hZGRyZXNzZXMgPSBhZGRyc1xuICAgICAgdGhpcy5hZGRyZXNzZXMuc29ydChBZGRyZXNzLmNvbXBhcmF0b3IoKSlcbiAgICAgIHRoaXMubnVtYWRkcnMud3JpdGVVSW50MzJCRSh0aGlzLmFkZHJlc3Nlcy5sZW5ndGgsIDApXG4gICAgfVxuICAgIGlmICh0eXBlb2YgdGhyZXNob2xkICE9PSB1bmRlZmluZWQpIHtcbiAgICAgIHRoaXMudGhyZXNob2xkLndyaXRlVUludDMyQkUodGhyZXNob2xkIHx8IDEsIDApXG4gICAgfVxuICAgIGlmICh0eXBlb2YgbG9ja3RpbWUgIT09IFwidW5kZWZpbmVkXCIpIHtcbiAgICAgIHRoaXMubG9ja3RpbWUgPSBiaW50b29scy5mcm9tQk5Ub0J1ZmZlcihsb2NrdGltZSwgOClcbiAgICB9XG4gIH1cbn1cblxuZXhwb3J0IGFic3RyYWN0IGNsYXNzIE91dHB1dCBleHRlbmRzIE91dHB1dE93bmVycyB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIk91dHB1dFwiXG4gIHByb3RlY3RlZCBfdHlwZUlEID0gdW5kZWZpbmVkXG5cbiAgLy9zZXJpYWxpemUgYW5kIGRlc2VyaWFsaXplIGJvdGggYXJlIGluaGVyaXRlZFxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBvdXRwdXRJRCBmb3IgdGhlIG91dHB1dCB3aGljaCB0ZWxscyBwYXJzZXJzIHdoYXQgdHlwZSBpdCBpc1xuICAgKi9cbiAgYWJzdHJhY3QgZ2V0T3V0cHV0SUQoKTogbnVtYmVyXG5cbiAgYWJzdHJhY3QgY2xvbmUoKTogdGhpc1xuXG4gIGFic3RyYWN0IGNyZWF0ZSguLi5hcmdzOiBhbnlbXSk6IHRoaXNcblxuICBhYnN0cmFjdCBzZWxlY3QoaWQ6IG51bWJlciwgLi4uYXJnczogYW55W10pOiBPdXRwdXRcblxuICAvKipcbiAgICpcbiAgICogQHBhcmFtIGFzc2V0SUQgQW4gYXNzZXRJRCB3aGljaCBpcyB3cmFwcGVkIGFyb3VuZCB0aGUgQnVmZmVyIG9mIHRoZSBPdXRwdXRcbiAgICpcbiAgICogTXVzdCBiZSBpbXBsZW1lbnRlZCB0byB1c2UgdGhlIGFwcHJvcHJpYXRlIFRyYW5zZmVyYWJsZU91dHB1dCBmb3IgdGhlIFZNLlxuICAgKi9cbiAgYWJzdHJhY3QgbWFrZVRyYW5zZmVyYWJsZShhc3NldElEOiBCdWZmZXIpOiBTdGFuZGFyZFRyYW5zZmVyYWJsZU91dHB1dFxufVxuXG5leHBvcnQgYWJzdHJhY3QgY2xhc3MgU3RhbmRhcmRQYXJzZWFibGVPdXRwdXQgZXh0ZW5kcyBTZXJpYWxpemFibGUge1xuICBwcm90ZWN0ZWQgX3R5cGVOYW1lID0gXCJTdGFuZGFyZFBhcnNlYWJsZU91dHB1dFwiXG4gIHByb3RlY3RlZCBfdHlwZUlEID0gdW5kZWZpbmVkXG5cbiAgc2VyaWFsaXplKGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKTogb2JqZWN0IHtcbiAgICBsZXQgZmllbGRzOiBvYmplY3QgPSBzdXBlci5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgcmV0dXJuIHtcbiAgICAgIC4uLmZpZWxkcyxcbiAgICAgIG91dHB1dDogdGhpcy5vdXRwdXQuc2VyaWFsaXplKGVuY29kaW5nKVxuICAgIH1cbiAgfVxuXG4gIHByb3RlY3RlZCBvdXRwdXQ6IE91dHB1dFxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIGEgZnVuY3Rpb24gdXNlZCB0byBzb3J0IGFuIGFycmF5IG9mIFtbUGFyc2VhYmxlT3V0cHV0XV1zXG4gICAqL1xuICBzdGF0aWMgY29tcGFyYXRvciA9XG4gICAgKCk6ICgoXG4gICAgICBhOiBTdGFuZGFyZFBhcnNlYWJsZU91dHB1dCxcbiAgICAgIGI6IFN0YW5kYXJkUGFyc2VhYmxlT3V0cHV0XG4gICAgKSA9PiAxIHwgLTEgfCAwKSA9PlxuICAgIChhOiBTdGFuZGFyZFBhcnNlYWJsZU91dHB1dCwgYjogU3RhbmRhcmRQYXJzZWFibGVPdXRwdXQpOiAxIHwgLTEgfCAwID0+IHtcbiAgICAgIGNvbnN0IHNvcnRhID0gYS50b0J1ZmZlcigpXG4gICAgICBjb25zdCBzb3J0YiA9IGIudG9CdWZmZXIoKVxuICAgICAgcmV0dXJuIEJ1ZmZlci5jb21wYXJlKHNvcnRhLCBzb3J0YikgYXMgMSB8IC0xIHwgMFxuICAgIH1cblxuICBnZXRPdXRwdXQgPSAoKTogT3V0cHV0ID0+IHRoaXMub3V0cHV0XG5cbiAgLy8gbXVzdCBiZSBpbXBsZW1lbnRlZCB0byBzZWxlY3Qgb3V0cHV0IHR5cGVzIGZvciB0aGUgVk0gaW4gcXVlc3Rpb25cbiAgYWJzdHJhY3QgZnJvbUJ1ZmZlcihieXRlczogQnVmZmVyLCBvZmZzZXQ/OiBudW1iZXIpOiBudW1iZXJcblxuICB0b0J1ZmZlcigpOiBCdWZmZXIge1xuICAgIGNvbnN0IG91dGJ1ZmY6IEJ1ZmZlciA9IHRoaXMub3V0cHV0LnRvQnVmZmVyKClcbiAgICBjb25zdCBvdXRpZDogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDQpXG4gICAgb3V0aWQud3JpdGVVSW50MzJCRSh0aGlzLm91dHB1dC5nZXRPdXRwdXRJRCgpLCAwKVxuICAgIGNvbnN0IGJhcnI6IEJ1ZmZlcltdID0gW291dGlkLCBvdXRidWZmXVxuICAgIHJldHVybiBCdWZmZXIuY29uY2F0KGJhcnIsIG91dGlkLmxlbmd0aCArIG91dGJ1ZmYubGVuZ3RoKVxuICB9XG5cbiAgLyoqXG4gICAqIENsYXNzIHJlcHJlc2VudGluZyBhbiBbW1BhcnNlYWJsZU91dHB1dF1dIGZvciBhIHRyYW5zYWN0aW9uLlxuICAgKlxuICAgKiBAcGFyYW0gb3V0cHV0IEEgbnVtYmVyIHJlcHJlc2VudGluZyB0aGUgSW5wdXRJRCBvZiB0aGUgW1tQYXJzZWFibGVPdXRwdXRdXVxuICAgKi9cbiAgY29uc3RydWN0b3Iob3V0cHV0OiBPdXRwdXQgPSB1bmRlZmluZWQpIHtcbiAgICBzdXBlcigpXG4gICAgaWYgKG91dHB1dCBpbnN0YW5jZW9mIE91dHB1dCkge1xuICAgICAgdGhpcy5vdXRwdXQgPSBvdXRwdXRcbiAgICB9XG4gIH1cbn1cblxuZXhwb3J0IGFic3RyYWN0IGNsYXNzIFN0YW5kYXJkVHJhbnNmZXJhYmxlT3V0cHV0IGV4dGVuZHMgU3RhbmRhcmRQYXJzZWFibGVPdXRwdXQge1xuICBwcm90ZWN0ZWQgX3R5cGVOYW1lID0gXCJTdGFuZGFyZFRyYW5zZmVyYWJsZU91dHB1dFwiXG4gIHByb3RlY3RlZCBfdHlwZUlEID0gdW5kZWZpbmVkXG5cbiAgc2VyaWFsaXplKGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKTogb2JqZWN0IHtcbiAgICBsZXQgZmllbGRzOiBvYmplY3QgPSBzdXBlci5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgcmV0dXJuIHtcbiAgICAgIC4uLmZpZWxkcyxcbiAgICAgIGFzc2V0SUQ6IHNlcmlhbGl6YXRpb24uZW5jb2Rlcih0aGlzLmFzc2V0SUQsIGVuY29kaW5nLCBcIkJ1ZmZlclwiLCBcImNiNThcIilcbiAgICB9XG4gIH1cbiAgZGVzZXJpYWxpemUoZmllbGRzOiBvYmplY3QsIGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKSB7XG4gICAgc3VwZXIuZGVzZXJpYWxpemUoZmllbGRzLCBlbmNvZGluZylcbiAgICB0aGlzLmFzc2V0SUQgPSBzZXJpYWxpemF0aW9uLmRlY29kZXIoXG4gICAgICBmaWVsZHNbXCJhc3NldElEXCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBcImNiNThcIixcbiAgICAgIFwiQnVmZmVyXCIsXG4gICAgICAzMlxuICAgIClcbiAgfVxuXG4gIHByb3RlY3RlZCBhc3NldElEOiBCdWZmZXIgPSB1bmRlZmluZWRcblxuICBnZXRBc3NldElEID0gKCk6IEJ1ZmZlciA9PiB0aGlzLmFzc2V0SURcblxuICAvLyBtdXN0IGJlIGltcGxlbWVudGVkIHRvIHNlbGVjdCBvdXRwdXQgdHlwZXMgZm9yIHRoZSBWTSBpbiBxdWVzdGlvblxuICBhYnN0cmFjdCBmcm9tQnVmZmVyKGJ5dGVzOiBCdWZmZXIsIG9mZnNldD86IG51bWJlcik6IG51bWJlclxuXG4gIHRvQnVmZmVyKCk6IEJ1ZmZlciB7XG4gICAgY29uc3QgcGFyc2VhYmxlQnVmZjogQnVmZmVyID0gc3VwZXIudG9CdWZmZXIoKVxuICAgIGNvbnN0IGJhcnI6IEJ1ZmZlcltdID0gW3RoaXMuYXNzZXRJRCwgcGFyc2VhYmxlQnVmZl1cbiAgICByZXR1cm4gQnVmZmVyLmNvbmNhdChiYXJyLCB0aGlzLmFzc2V0SUQubGVuZ3RoICsgcGFyc2VhYmxlQnVmZi5sZW5ndGgpXG4gIH1cblxuICAvKipcbiAgICogQ2xhc3MgcmVwcmVzZW50aW5nIGFuIFtbU3RhbmRhcmRUcmFuc2ZlcmFibGVPdXRwdXRdXSBmb3IgYSB0cmFuc2FjdGlvbi5cbiAgICpcbiAgICogQHBhcmFtIGFzc2V0SUQgQSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSByZXByZXNlbnRpbmcgdGhlIGFzc2V0SUQgb2YgdGhlIFtbT3V0cHV0XV1cbiAgICogQHBhcmFtIG91dHB1dCBBIG51bWJlciByZXByZXNlbnRpbmcgdGhlIElucHV0SUQgb2YgdGhlIFtbU3RhbmRhcmRUcmFuc2ZlcmFibGVPdXRwdXRdXVxuICAgKi9cbiAgY29uc3RydWN0b3IoYXNzZXRJRDogQnVmZmVyID0gdW5kZWZpbmVkLCBvdXRwdXQ6IE91dHB1dCA9IHVuZGVmaW5lZCkge1xuICAgIHN1cGVyKG91dHB1dClcbiAgICBpZiAodHlwZW9mIGFzc2V0SUQgIT09IFwidW5kZWZpbmVkXCIpIHtcbiAgICAgIHRoaXMuYXNzZXRJRCA9IGFzc2V0SURcbiAgICB9XG4gIH1cbn1cblxuLyoqXG4gKiBBbiBbW091dHB1dF1dIGNsYXNzIHdoaWNoIHNwZWNpZmllcyBhIHRva2VuIGFtb3VudCAuXG4gKi9cbmV4cG9ydCBhYnN0cmFjdCBjbGFzcyBTdGFuZGFyZEFtb3VudE91dHB1dCBleHRlbmRzIE91dHB1dCB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIlN0YW5kYXJkQW1vdW50T3V0cHV0XCJcbiAgcHJvdGVjdGVkIF90eXBlSUQgPSB1bmRlZmluZWRcblxuICBzZXJpYWxpemUoZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpOiBvYmplY3Qge1xuICAgIGxldCBmaWVsZHM6IG9iamVjdCA9IHN1cGVyLnNlcmlhbGl6ZShlbmNvZGluZylcbiAgICByZXR1cm4ge1xuICAgICAgLi4uZmllbGRzLFxuICAgICAgYW1vdW50OiBzZXJpYWxpemF0aW9uLmVuY29kZXIoXG4gICAgICAgIHRoaXMuYW1vdW50LFxuICAgICAgICBlbmNvZGluZyxcbiAgICAgICAgXCJCdWZmZXJcIixcbiAgICAgICAgXCJkZWNpbWFsU3RyaW5nXCIsXG4gICAgICAgIDhcbiAgICAgIClcbiAgICB9XG4gIH1cbiAgZGVzZXJpYWxpemUoZmllbGRzOiBvYmplY3QsIGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKSB7XG4gICAgc3VwZXIuZGVzZXJpYWxpemUoZmllbGRzLCBlbmNvZGluZylcbiAgICB0aGlzLmFtb3VudCA9IHNlcmlhbGl6YXRpb24uZGVjb2RlcihcbiAgICAgIGZpZWxkc1tcImFtb3VudFwiXSxcbiAgICAgIGVuY29kaW5nLFxuICAgICAgXCJkZWNpbWFsU3RyaW5nXCIsXG4gICAgICBcIkJ1ZmZlclwiLFxuICAgICAgOFxuICAgIClcbiAgICB0aGlzLmFtb3VudFZhbHVlID0gYmludG9vbHMuZnJvbUJ1ZmZlclRvQk4odGhpcy5hbW91bnQpXG4gIH1cblxuICBwcm90ZWN0ZWQgYW1vdW50OiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoOClcbiAgcHJvdGVjdGVkIGFtb3VudFZhbHVlOiBCTiA9IG5ldyBCTigwKVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBhbW91bnQgYXMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2luZHV0bnkvYm4uanMvfEJOfS5cbiAgICovXG4gIGdldEFtb3VudCgpOiBCTiB7XG4gICAgcmV0dXJuIHRoaXMuYW1vdW50VmFsdWUuY2xvbmUoKVxuICB9XG5cbiAgLyoqXG4gICAqIFBvcHVhdGVzIHRoZSBpbnN0YW5jZSBmcm9tIGEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gcmVwcmVzZW50aW5nIHRoZSBbW1N0YW5kYXJkQW1vdW50T3V0cHV0XV0gYW5kIHJldHVybnMgdGhlIHNpemUgb2YgdGhlIG91dHB1dC5cbiAgICovXG4gIGZyb21CdWZmZXIob3V0YnVmZjogQnVmZmVyLCBvZmZzZXQ6IG51bWJlciA9IDApOiBudW1iZXIge1xuICAgIHRoaXMuYW1vdW50ID0gYmludG9vbHMuY29weUZyb20ob3V0YnVmZiwgb2Zmc2V0LCBvZmZzZXQgKyA4KVxuICAgIHRoaXMuYW1vdW50VmFsdWUgPSBiaW50b29scy5mcm9tQnVmZmVyVG9CTih0aGlzLmFtb3VudClcbiAgICBvZmZzZXQgKz0gOFxuICAgIHJldHVybiBzdXBlci5mcm9tQnVmZmVyKG91dGJ1ZmYsIG9mZnNldClcbiAgfVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBidWZmZXIgcmVwcmVzZW50aW5nIHRoZSBbW1N0YW5kYXJkQW1vdW50T3V0cHV0XV0gaW5zdGFuY2UuXG4gICAqL1xuICB0b0J1ZmZlcigpOiBCdWZmZXIge1xuICAgIGNvbnN0IHN1cGVyYnVmZjogQnVmZmVyID0gc3VwZXIudG9CdWZmZXIoKVxuICAgIGNvbnN0IGJzaXplOiBudW1iZXIgPSB0aGlzLmFtb3VudC5sZW5ndGggKyBzdXBlcmJ1ZmYubGVuZ3RoXG4gICAgdGhpcy5udW1hZGRycy53cml0ZVVJbnQzMkJFKHRoaXMuYWRkcmVzc2VzLmxlbmd0aCwgMClcbiAgICBjb25zdCBiYXJyOiBCdWZmZXJbXSA9IFt0aGlzLmFtb3VudCwgc3VwZXJidWZmXVxuICAgIHJldHVybiBCdWZmZXIuY29uY2F0KGJhcnIsIGJzaXplKVxuICB9XG5cbiAgLyoqXG4gICAqIEEgW1tTdGFuZGFyZEFtb3VudE91dHB1dF1dIGNsYXNzIHdoaWNoIGlzc3VlcyBhIHBheW1lbnQgb24gYW4gYXNzZXRJRC5cbiAgICpcbiAgICogQHBhcmFtIGFtb3VudCBBIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vaW5kdXRueS9ibi5qcy98Qk59IHJlcHJlc2VudGluZyB0aGUgYW1vdW50IGluIHRoZSBvdXRwdXRcbiAgICogQHBhcmFtIGFkZHJlc3NlcyBBbiBhcnJheSBvZiB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfXMgcmVwcmVzZW50aW5nIGFkZHJlc3Nlc1xuICAgKiBAcGFyYW0gbG9ja3RpbWUgQSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2luZHV0bnkvYm4uanMvfEJOfSByZXByZXNlbnRpbmcgdGhlIGxvY2t0aW1lXG4gICAqIEBwYXJhbSB0aHJlc2hvbGQgQSBudW1iZXIgcmVwcmVzZW50aW5nIHRoZSB0aGUgdGhyZXNob2xkIG51bWJlciBvZiBzaWduZXJzIHJlcXVpcmVkIHRvIHNpZ24gdGhlIHRyYW5zYWN0aW9uXG4gICAqL1xuICBjb25zdHJ1Y3RvcihcbiAgICBhbW91bnQ6IEJOID0gdW5kZWZpbmVkLFxuICAgIGFkZHJlc3NlczogQnVmZmVyW10gPSB1bmRlZmluZWQsXG4gICAgbG9ja3RpbWU6IEJOID0gdW5kZWZpbmVkLFxuICAgIHRocmVzaG9sZDogbnVtYmVyID0gdW5kZWZpbmVkXG4gICkge1xuICAgIHN1cGVyKGFkZHJlc3NlcywgbG9ja3RpbWUsIHRocmVzaG9sZClcbiAgICBpZiAodHlwZW9mIGFtb3VudCAhPT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgdGhpcy5hbW91bnRWYWx1ZSA9IGFtb3VudC5jbG9uZSgpXG4gICAgICB0aGlzLmFtb3VudCA9IGJpbnRvb2xzLmZyb21CTlRvQnVmZmVyKGFtb3VudCwgOClcbiAgICB9XG4gIH1cbn1cblxuLyoqXG4gKiBBbiBbW091dHB1dF1dIGNsYXNzIHdoaWNoIHNwZWNpZmllcyBhbiBORlQuXG4gKi9cbmV4cG9ydCBhYnN0cmFjdCBjbGFzcyBCYXNlTkZUT3V0cHV0IGV4dGVuZHMgT3V0cHV0IHtcbiAgcHJvdGVjdGVkIF90eXBlTmFtZSA9IFwiQmFzZU5GVE91dHB1dFwiXG4gIHByb3RlY3RlZCBfdHlwZUlEID0gdW5kZWZpbmVkXG5cbiAgc2VyaWFsaXplKGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKTogb2JqZWN0IHtcbiAgICBsZXQgZmllbGRzOiBvYmplY3QgPSBzdXBlci5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgcmV0dXJuIHtcbiAgICAgIC4uLmZpZWxkcyxcbiAgICAgIGdyb3VwSUQ6IHNlcmlhbGl6YXRpb24uZW5jb2RlcihcbiAgICAgICAgdGhpcy5ncm91cElELFxuICAgICAgICBlbmNvZGluZyxcbiAgICAgICAgXCJCdWZmZXJcIixcbiAgICAgICAgXCJkZWNpbWFsU3RyaW5nXCIsXG4gICAgICAgIDRcbiAgICAgIClcbiAgICB9XG4gIH1cbiAgZGVzZXJpYWxpemUoZmllbGRzOiBvYmplY3QsIGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKSB7XG4gICAgc3VwZXIuZGVzZXJpYWxpemUoZmllbGRzLCBlbmNvZGluZylcbiAgICB0aGlzLmdyb3VwSUQgPSBzZXJpYWxpemF0aW9uLmRlY29kZXIoXG4gICAgICBmaWVsZHNbXCJncm91cElEXCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBcImRlY2ltYWxTdHJpbmdcIixcbiAgICAgIFwiQnVmZmVyXCIsXG4gICAgICA0XG4gICAgKVxuICB9XG5cbiAgcHJvdGVjdGVkIGdyb3VwSUQ6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYyg0KVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBncm91cElEIGFzIGEgbnVtYmVyLlxuICAgKi9cbiAgZ2V0R3JvdXBJRCA9ICgpOiBudW1iZXIgPT4ge1xuICAgIHJldHVybiB0aGlzLmdyb3VwSUQucmVhZFVJbnQzMkJFKDApXG4gIH1cbn1cbiJdfQ==