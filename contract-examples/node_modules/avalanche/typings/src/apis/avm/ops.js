"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.UTXOID = exports.NFTTransferOperation = exports.NFTMintOperation = exports.SECPMintOperation = exports.TransferableOperation = exports.Operation = exports.SelectOperationClass = void 0;
/**
 * @packageDocumentation
 * @module API-AVM-Operations
 */
const buffer_1 = require("buffer/");
const bintools_1 = __importDefault(require("../../utils/bintools"));
const constants_1 = require("./constants");
const outputs_1 = require("./outputs");
const nbytes_1 = require("../../common/nbytes");
const credentials_1 = require("../../common/credentials");
const output_1 = require("../../common/output");
const serialization_1 = require("../../utils/serialization");
const errors_1 = require("../../utils/errors");
const bintools = bintools_1.default.getInstance();
const serialization = serialization_1.Serialization.getInstance();
const cb58 = "cb58";
const buffer = "Buffer";
const hex = "hex";
const decimalString = "decimalString";
/**
 * Takes a buffer representing the output and returns the proper [[Operation]] instance.
 *
 * @param opid A number representing the operation ID parsed prior to the bytes passed in
 *
 * @returns An instance of an [[Operation]]-extended class.
 */
const SelectOperationClass = (opid, ...args) => {
    if (opid === constants_1.AVMConstants.SECPMINTOPID ||
        opid === constants_1.AVMConstants.SECPMINTOPID_CODECONE) {
        return new SECPMintOperation(...args);
    }
    else if (opid === constants_1.AVMConstants.NFTMINTOPID ||
        opid === constants_1.AVMConstants.NFTMINTOPID_CODECONE) {
        return new NFTMintOperation(...args);
    }
    else if (opid === constants_1.AVMConstants.NFTXFEROPID ||
        opid === constants_1.AVMConstants.NFTXFEROPID_CODECONE) {
        return new NFTTransferOperation(...args);
    }
    /* istanbul ignore next */
    throw new errors_1.InvalidOperationIdError(`Error - SelectOperationClass: unknown opid ${opid}`);
};
exports.SelectOperationClass = SelectOperationClass;
/**
 * A class representing an operation. All operation types must extend on this class.
 */
class Operation extends serialization_1.Serializable {
    constructor() {
        super(...arguments);
        this._typeName = "Operation";
        this._typeID = undefined;
        this.sigCount = buffer_1.Buffer.alloc(4);
        this.sigIdxs = []; // idxs of signers from utxo
        /**
         * Returns the array of [[SigIdx]] for this [[Operation]]
         */
        this.getSigIdxs = () => this.sigIdxs;
        /**
         * Creates and adds a [[SigIdx]] to the [[Operation]].
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
     * Returns a base-58 string representing the [[NFTMintOperation]].
     */
    toString() {
        return bintools.bufferToB58(this.toBuffer());
    }
}
exports.Operation = Operation;
Operation.comparator = () => (a, b) => {
    const aoutid = buffer_1.Buffer.alloc(4);
    aoutid.writeUInt32BE(a.getOperationID(), 0);
    const abuff = a.toBuffer();
    const boutid = buffer_1.Buffer.alloc(4);
    boutid.writeUInt32BE(b.getOperationID(), 0);
    const bbuff = b.toBuffer();
    const asort = buffer_1.Buffer.concat([aoutid, abuff], aoutid.length + abuff.length);
    const bsort = buffer_1.Buffer.concat([boutid, bbuff], boutid.length + bbuff.length);
    return buffer_1.Buffer.compare(asort, bsort);
};
/**
 * A class which contains an [[Operation]] for transfers.
 *
 */
class TransferableOperation extends serialization_1.Serializable {
    constructor(assetID = undefined, utxoids = undefined, operation = undefined) {
        super();
        this._typeName = "TransferableOperation";
        this._typeID = undefined;
        this.assetID = buffer_1.Buffer.alloc(32);
        this.utxoIDs = [];
        /**
         * Returns the assetID as a {@link https://github.com/feross/buffer|Buffer}.
         */
        this.getAssetID = () => this.assetID;
        /**
         * Returns an array of UTXOIDs in this operation.
         */
        this.getUTXOIDs = () => this.utxoIDs;
        /**
         * Returns the operation
         */
        this.getOperation = () => this.operation;
        if (typeof assetID !== "undefined" &&
            assetID.length === constants_1.AVMConstants.ASSETIDLEN &&
            operation instanceof Operation &&
            typeof utxoids !== "undefined" &&
            Array.isArray(utxoids)) {
            this.assetID = assetID;
            this.operation = operation;
            for (let i = 0; i < utxoids.length; i++) {
                const utxoid = new UTXOID();
                if (typeof utxoids[`${i}`] === "string") {
                    utxoid.fromString(utxoids[`${i}`]);
                }
                else if (utxoids[`${i}`] instanceof buffer_1.Buffer) {
                    utxoid.fromBuffer(utxoids[`${i}`]);
                }
                else if (utxoids[`${i}`] instanceof UTXOID) {
                    utxoid.fromString(utxoids[`${i}`].toString()); // clone
                }
                this.utxoIDs.push(utxoid);
            }
        }
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { assetID: serialization.encoder(this.assetID, encoding, buffer, cb58, 32), utxoIDs: this.utxoIDs.map((u) => u.serialize(encoding)), operation: this.operation.serialize(encoding) });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.assetID = serialization.decoder(fields["assetID"], encoding, cb58, buffer, 32);
        this.utxoIDs = fields["utxoIDs"].map((u) => {
            let utxoid = new UTXOID();
            utxoid.deserialize(u, encoding);
            return utxoid;
        });
        this.operation = (0, exports.SelectOperationClass)(fields["operation"]["_typeID"]);
        this.operation.deserialize(fields["operation"], encoding);
    }
    fromBuffer(bytes, offset = 0) {
        this.assetID = bintools.copyFrom(bytes, offset, offset + 32);
        offset += 32;
        const numutxoIDs = bintools
            .copyFrom(bytes, offset, offset + 4)
            .readUInt32BE(0);
        offset += 4;
        this.utxoIDs = [];
        for (let i = 0; i < numutxoIDs; i++) {
            const utxoid = new UTXOID();
            offset = utxoid.fromBuffer(bytes, offset);
            this.utxoIDs.push(utxoid);
        }
        const opid = bintools
            .copyFrom(bytes, offset, offset + 4)
            .readUInt32BE(0);
        offset += 4;
        this.operation = (0, exports.SelectOperationClass)(opid);
        return this.operation.fromBuffer(bytes, offset);
    }
    toBuffer() {
        const numutxoIDs = buffer_1.Buffer.alloc(4);
        numutxoIDs.writeUInt32BE(this.utxoIDs.length, 0);
        let bsize = this.assetID.length + numutxoIDs.length;
        const barr = [this.assetID, numutxoIDs];
        this.utxoIDs = this.utxoIDs.sort(UTXOID.comparator());
        for (let i = 0; i < this.utxoIDs.length; i++) {
            const b = this.utxoIDs[`${i}`].toBuffer();
            barr.push(b);
            bsize += b.length;
        }
        const opid = buffer_1.Buffer.alloc(4);
        opid.writeUInt32BE(this.operation.getOperationID(), 0);
        barr.push(opid);
        bsize += opid.length;
        const b = this.operation.toBuffer();
        bsize += b.length;
        barr.push(b);
        return buffer_1.Buffer.concat(barr, bsize);
    }
}
exports.TransferableOperation = TransferableOperation;
/**
 * Returns a function used to sort an array of [[TransferableOperation]]s
 */
TransferableOperation.comparator = () => {
    return function (a, b) {
        return buffer_1.Buffer.compare(a.toBuffer(), b.toBuffer());
    };
};
/**
 * An [[Operation]] class which specifies a SECP256k1 Mint Op.
 */
class SECPMintOperation extends Operation {
    /**
     * An [[Operation]] class which mints new tokens on an assetID.
     *
     * @param mintOutput The [[SECPMintOutput]] that will be produced by this transaction.
     * @param transferOutput A [[SECPTransferOutput]] that will be produced from this minting operation.
     */
    constructor(mintOutput = undefined, transferOutput = undefined) {
        super();
        this._typeName = "SECPMintOperation";
        this._codecID = constants_1.AVMConstants.LATESTCODEC;
        this._typeID = this._codecID === 0
            ? constants_1.AVMConstants.SECPMINTOPID
            : constants_1.AVMConstants.SECPMINTOPID_CODECONE;
        this.mintOutput = undefined;
        this.transferOutput = undefined;
        if (typeof mintOutput !== "undefined") {
            this.mintOutput = mintOutput;
        }
        if (typeof transferOutput !== "undefined") {
            this.transferOutput = transferOutput;
        }
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { mintOutput: this.mintOutput.serialize(encoding), transferOutputs: this.transferOutput.serialize(encoding) });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.mintOutput = new outputs_1.SECPMintOutput();
        this.mintOutput.deserialize(fields["mintOutput"], encoding);
        this.transferOutput = new outputs_1.SECPTransferOutput();
        this.transferOutput.deserialize(fields["transferOutputs"], encoding);
    }
    /**
     * Set the codecID
     *
     * @param codecID The codecID to set
     */
    setCodecID(codecID) {
        if (codecID !== 0 && codecID !== 1) {
            /* istanbul ignore next */
            throw new errors_1.CodecIdError("Error - SECPMintOperation.setCodecID: invalid codecID. Valid codecIDs are 0 and 1.");
        }
        this._codecID = codecID;
        this._typeID =
            this._codecID === 0
                ? constants_1.AVMConstants.SECPMINTOPID
                : constants_1.AVMConstants.SECPMINTOPID_CODECONE;
    }
    /**
     * Returns the operation ID.
     */
    getOperationID() {
        return this._typeID;
    }
    /**
     * Returns the credential ID.
     */
    getCredentialID() {
        if (this._codecID === 0) {
            return constants_1.AVMConstants.SECPCREDENTIAL;
        }
        else if (this._codecID === 1) {
            return constants_1.AVMConstants.SECPCREDENTIAL_CODECONE;
        }
    }
    /**
     * Returns the [[SECPMintOutput]] to be produced by this operation.
     */
    getMintOutput() {
        return this.mintOutput;
    }
    /**
     * Returns [[SECPTransferOutput]] to be produced by this operation.
     */
    getTransferOutput() {
        return this.transferOutput;
    }
    /**
     * Popuates the instance from a {@link https://github.com/feross/buffer|Buffer} representing the [[SECPMintOperation]] and returns the updated offset.
     */
    fromBuffer(bytes, offset = 0) {
        offset = super.fromBuffer(bytes, offset);
        this.mintOutput = new outputs_1.SECPMintOutput();
        offset = this.mintOutput.fromBuffer(bytes, offset);
        this.transferOutput = new outputs_1.SECPTransferOutput();
        offset = this.transferOutput.fromBuffer(bytes, offset);
        return offset;
    }
    /**
     * Returns the buffer representing the [[SECPMintOperation]] instance.
     */
    toBuffer() {
        const superbuff = super.toBuffer();
        const mintoutBuff = this.mintOutput.toBuffer();
        const transferOutBuff = this.transferOutput.toBuffer();
        const bsize = superbuff.length + mintoutBuff.length + transferOutBuff.length;
        const barr = [superbuff, mintoutBuff, transferOutBuff];
        return buffer_1.Buffer.concat(barr, bsize);
    }
}
exports.SECPMintOperation = SECPMintOperation;
/**
 * An [[Operation]] class which specifies a NFT Mint Op.
 */
class NFTMintOperation extends Operation {
    /**
     * An [[Operation]] class which contains an NFT on an assetID.
     *
     * @param groupID The group to which to issue the NFT Output
     * @param payload A {@link https://github.com/feross/buffer|Buffer} of the NFT payload
     * @param outputOwners An array of outputOwners
     */
    constructor(groupID = undefined, payload = undefined, outputOwners = undefined) {
        super();
        this._typeName = "NFTMintOperation";
        this._codecID = constants_1.AVMConstants.LATESTCODEC;
        this._typeID = this._codecID === 0
            ? constants_1.AVMConstants.NFTMINTOPID
            : constants_1.AVMConstants.NFTMINTOPID_CODECONE;
        this.groupID = buffer_1.Buffer.alloc(4);
        this.outputOwners = [];
        /**
         * Returns the credential ID.
         */
        this.getCredentialID = () => {
            if (this._codecID === 0) {
                return constants_1.AVMConstants.NFTCREDENTIAL;
            }
            else if (this._codecID === 1) {
                return constants_1.AVMConstants.NFTCREDENTIAL_CODECONE;
            }
        };
        /**
         * Returns the payload.
         */
        this.getGroupID = () => {
            return bintools.copyFrom(this.groupID, 0);
        };
        /**
         * Returns the payload.
         */
        this.getPayload = () => {
            return bintools.copyFrom(this.payload, 0);
        };
        /**
         * Returns the payload's raw {@link https://github.com/feross/buffer|Buffer} with length prepended, for use with [[PayloadBase]]'s fromBuffer
         */
        this.getPayloadBuffer = () => {
            let payloadlen = buffer_1.Buffer.alloc(4);
            payloadlen.writeUInt32BE(this.payload.length, 0);
            return buffer_1.Buffer.concat([payloadlen, bintools.copyFrom(this.payload, 0)]);
        };
        /**
         * Returns the outputOwners.
         */
        this.getOutputOwners = () => {
            return this.outputOwners;
        };
        if (typeof groupID !== "undefined" &&
            typeof payload !== "undefined" &&
            outputOwners.length) {
            this.groupID.writeUInt32BE(groupID ? groupID : 0, 0);
            this.payload = payload;
            this.outputOwners = outputOwners;
        }
    }
    serialize(encoding = "hex") {
        const fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { groupID: serialization.encoder(this.groupID, encoding, buffer, decimalString, 4), payload: serialization.encoder(this.payload, encoding, buffer, hex), outputOwners: this.outputOwners.map((o) => o.serialize(encoding)) });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.groupID = serialization.decoder(fields["groupID"], encoding, decimalString, buffer, 4);
        this.payload = serialization.decoder(fields["payload"], encoding, hex, buffer);
        // this.outputOwners = fields["outputOwners"].map((o: NFTMintOutput) => {
        //   let oo: NFTMintOutput = new NFTMintOutput()
        //   oo.deserialize(o, encoding)
        //   return oo
        // })
        this.outputOwners = fields["outputOwners"].map((o) => {
            let oo = new output_1.OutputOwners();
            oo.deserialize(o, encoding);
            return oo;
        });
    }
    /**
     * Set the codecID
     *
     * @param codecID The codecID to set
     */
    setCodecID(codecID) {
        if (codecID !== 0 && codecID !== 1) {
            /* istanbul ignore next */
            throw new errors_1.CodecIdError("Error - NFTMintOperation.setCodecID: invalid codecID. Valid codecIDs are 0 and 1.");
        }
        this._codecID = codecID;
        this._typeID =
            this._codecID === 0
                ? constants_1.AVMConstants.NFTMINTOPID
                : constants_1.AVMConstants.NFTMINTOPID_CODECONE;
    }
    /**
     * Returns the operation ID.
     */
    getOperationID() {
        return this._typeID;
    }
    /**
     * Popuates the instance from a {@link https://github.com/feross/buffer|Buffer} representing the [[NFTMintOperation]] and returns the updated offset.
     */
    fromBuffer(bytes, offset = 0) {
        offset = super.fromBuffer(bytes, offset);
        this.groupID = bintools.copyFrom(bytes, offset, offset + 4);
        offset += 4;
        let payloadLen = bintools
            .copyFrom(bytes, offset, offset + 4)
            .readUInt32BE(0);
        offset += 4;
        this.payload = bintools.copyFrom(bytes, offset, offset + payloadLen);
        offset += payloadLen;
        let numoutputs = bintools
            .copyFrom(bytes, offset, offset + 4)
            .readUInt32BE(0);
        offset += 4;
        this.outputOwners = [];
        for (let i = 0; i < numoutputs; i++) {
            let outputOwner = new output_1.OutputOwners();
            offset = outputOwner.fromBuffer(bytes, offset);
            this.outputOwners.push(outputOwner);
        }
        return offset;
    }
    /**
     * Returns the buffer representing the [[NFTMintOperation]] instance.
     */
    toBuffer() {
        const superbuff = super.toBuffer();
        const payloadlen = buffer_1.Buffer.alloc(4);
        payloadlen.writeUInt32BE(this.payload.length, 0);
        const outputownerslen = buffer_1.Buffer.alloc(4);
        outputownerslen.writeUInt32BE(this.outputOwners.length, 0);
        let bsize = superbuff.length +
            this.groupID.length +
            payloadlen.length +
            this.payload.length +
            outputownerslen.length;
        const barr = [
            superbuff,
            this.groupID,
            payloadlen,
            this.payload,
            outputownerslen
        ];
        for (let i = 0; i < this.outputOwners.length; i++) {
            let b = this.outputOwners[`${i}`].toBuffer();
            barr.push(b);
            bsize += b.length;
        }
        return buffer_1.Buffer.concat(barr, bsize);
    }
    /**
     * Returns a base-58 string representing the [[NFTMintOperation]].
     */
    toString() {
        return bintools.bufferToB58(this.toBuffer());
    }
}
exports.NFTMintOperation = NFTMintOperation;
/**
 * A [[Operation]] class which specifies a NFT Transfer Op.
 */
class NFTTransferOperation extends Operation {
    /**
     * An [[Operation]] class which contains an NFT on an assetID.
     *
     * @param output An [[NFTTransferOutput]]
     */
    constructor(output = undefined) {
        super();
        this._typeName = "NFTTransferOperation";
        this._codecID = constants_1.AVMConstants.LATESTCODEC;
        this._typeID = this._codecID === 0
            ? constants_1.AVMConstants.NFTXFEROPID
            : constants_1.AVMConstants.NFTXFEROPID_CODECONE;
        this.getOutput = () => this.output;
        if (typeof output !== "undefined") {
            this.output = output;
        }
    }
    serialize(encoding = "hex") {
        const fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { output: this.output.serialize(encoding) });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.output = new outputs_1.NFTTransferOutput();
        this.output.deserialize(fields["output"], encoding);
    }
    /**
     * Set the codecID
     *
     * @param codecID The codecID to set
     */
    setCodecID(codecID) {
        if (codecID !== 0 && codecID !== 1) {
            /* istanbul ignore next */
            throw new errors_1.CodecIdError("Error - NFTTransferOperation.setCodecID: invalid codecID. Valid codecIDs are 0 and 1.");
        }
        this._codecID = codecID;
        this._typeID =
            this._codecID === 0
                ? constants_1.AVMConstants.NFTXFEROPID
                : constants_1.AVMConstants.NFTXFEROPID_CODECONE;
    }
    /**
     * Returns the operation ID.
     */
    getOperationID() {
        return this._typeID;
    }
    /**
     * Returns the credential ID.
     */
    getCredentialID() {
        if (this._codecID === 0) {
            return constants_1.AVMConstants.NFTCREDENTIAL;
        }
        else if (this._codecID === 1) {
            return constants_1.AVMConstants.NFTCREDENTIAL_CODECONE;
        }
    }
    /**
     * Popuates the instance from a {@link https://github.com/feross/buffer|Buffer} representing the [[NFTTransferOperation]] and returns the updated offset.
     */
    fromBuffer(bytes, offset = 0) {
        offset = super.fromBuffer(bytes, offset);
        this.output = new outputs_1.NFTTransferOutput();
        return this.output.fromBuffer(bytes, offset);
    }
    /**
     * Returns the buffer representing the [[NFTTransferOperation]] instance.
     */
    toBuffer() {
        const superbuff = super.toBuffer();
        const outbuff = this.output.toBuffer();
        const bsize = superbuff.length + outbuff.length;
        const barr = [superbuff, outbuff];
        return buffer_1.Buffer.concat(barr, bsize);
    }
    /**
     * Returns a base-58 string representing the [[NFTTransferOperation]].
     */
    toString() {
        return bintools.bufferToB58(this.toBuffer());
    }
}
exports.NFTTransferOperation = NFTTransferOperation;
/**
 * Class for representing a UTXOID used in [[TransferableOp]] types
 */
class UTXOID extends nbytes_1.NBytes {
    /**
     * Class for representing a UTXOID used in [[TransferableOp]] types
     */
    constructor() {
        super();
        this._typeName = "UTXOID";
        this._typeID = undefined;
        //serialize and deserialize both are inherited
        this.bytes = buffer_1.Buffer.alloc(36);
        this.bsize = 36;
    }
    /**
     * Returns a base-58 representation of the [[UTXOID]].
     */
    toString() {
        return bintools.cb58Encode(this.toBuffer());
    }
    /**
     * Takes a base-58 string containing an [[UTXOID]], parses it, populates the class, and returns the length of the UTXOID in bytes.
     *
     * @param bytes A base-58 string containing a raw [[UTXOID]]
     *
     * @returns The length of the raw [[UTXOID]]
     */
    fromString(utxoid) {
        const utxoidbuff = bintools.b58ToBuffer(utxoid);
        if (utxoidbuff.length === 40 && bintools.validateChecksum(utxoidbuff)) {
            const newbuff = bintools.copyFrom(utxoidbuff, 0, utxoidbuff.length - 4);
            if (newbuff.length === 36) {
                this.bytes = newbuff;
            }
        }
        else if (utxoidbuff.length === 40) {
            throw new errors_1.ChecksumError("Error - UTXOID.fromString: invalid checksum on address");
        }
        else if (utxoidbuff.length === 36) {
            this.bytes = utxoidbuff;
        }
        else {
            /* istanbul ignore next */
            throw new errors_1.AddressError("Error - UTXOID.fromString: invalid address");
        }
        return this.getSize();
    }
    clone() {
        const newbase = new UTXOID();
        newbase.fromBuffer(this.toBuffer());
        return newbase;
    }
    create(...args) {
        return new UTXOID();
    }
}
exports.UTXOID = UTXOID;
/**
 * Returns a function used to sort an array of [[UTXOID]]s
 */
UTXOID.comparator = () => (a, b) => buffer_1.Buffer.compare(a.toBuffer(), b.toBuffer());
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoib3BzLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2FwaXMvYXZtL29wcy50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7QUFBQTs7O0dBR0c7QUFDSCxvQ0FBZ0M7QUFDaEMsb0VBQTJDO0FBQzNDLDJDQUEwQztBQUMxQyx1Q0FJa0I7QUFDbEIsZ0RBQTRDO0FBQzVDLDBEQUFpRDtBQUNqRCxnREFBa0Q7QUFDbEQsNkRBS2tDO0FBQ2xDLCtDQUsyQjtBQUUzQixNQUFNLFFBQVEsR0FBYSxrQkFBUSxDQUFDLFdBQVcsRUFBRSxDQUFBO0FBQ2pELE1BQU0sYUFBYSxHQUFrQiw2QkFBYSxDQUFDLFdBQVcsRUFBRSxDQUFBO0FBQ2hFLE1BQU0sSUFBSSxHQUFtQixNQUFNLENBQUE7QUFDbkMsTUFBTSxNQUFNLEdBQW1CLFFBQVEsQ0FBQTtBQUN2QyxNQUFNLEdBQUcsR0FBbUIsS0FBSyxDQUFBO0FBQ2pDLE1BQU0sYUFBYSxHQUFtQixlQUFlLENBQUE7QUFFckQ7Ozs7OztHQU1HO0FBQ0ksTUFBTSxvQkFBb0IsR0FBRyxDQUNsQyxJQUFZLEVBQ1osR0FBRyxJQUFXLEVBQ0gsRUFBRTtJQUNiLElBQ0UsSUFBSSxLQUFLLHdCQUFZLENBQUMsWUFBWTtRQUNsQyxJQUFJLEtBQUssd0JBQVksQ0FBQyxxQkFBcUIsRUFDM0M7UUFDQSxPQUFPLElBQUksaUJBQWlCLENBQUMsR0FBRyxJQUFJLENBQUMsQ0FBQTtLQUN0QztTQUFNLElBQ0wsSUFBSSxLQUFLLHdCQUFZLENBQUMsV0FBVztRQUNqQyxJQUFJLEtBQUssd0JBQVksQ0FBQyxvQkFBb0IsRUFDMUM7UUFDQSxPQUFPLElBQUksZ0JBQWdCLENBQUMsR0FBRyxJQUFJLENBQUMsQ0FBQTtLQUNyQztTQUFNLElBQ0wsSUFBSSxLQUFLLHdCQUFZLENBQUMsV0FBVztRQUNqQyxJQUFJLEtBQUssd0JBQVksQ0FBQyxvQkFBb0IsRUFDMUM7UUFDQSxPQUFPLElBQUksb0JBQW9C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QUUsQ0FBQTtJQUN2QixDQUFDO0lBRUQsS0FBSztRQUNILE1BQU0sT0FBTyxHQUFXLElBQUksTUFBTSxFQUFFLENBQUE7UUFDcEMsT0FBTyxDQUFDLFVBQVUsQ0FBQyxJQUFJLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQTtRQUNuQyxPQUFPLE9BQWUsQ0FBQTtJQUN4QixDQUFDO0lBRUQsTUFBTSxDQUFDLEdBQUcsSUFBVztRQUNuQixPQUFPLElBQUksTUFBTSxFQUFVLENBQUE7SUFDN0IsQ0FBQzs7QUEvREgsd0JBdUVDO0FBOURDOztHQUVHO0FBQ0ksaUJBQVUsR0FDZixHQUEyQyxFQUFFLENBQzdDLENBQUMsQ0FBUyxFQUFFLENBQVMsRUFBYyxFQUFFLENBQ25DLGVBQU0sQ0FBQyxPQUFPLENBQUMsQ0FBQyxDQUFDLFFBQVEsRUFBRSxFQUFFLENBQUMsQ0FBQyxRQUFRLEVBQUUsQ0FBZSxDQUFBIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgQVBJLUFWTS1PcGVyYXRpb25zXG4gKi9cbmltcG9ydCB7IEJ1ZmZlciB9IGZyb20gXCJidWZmZXIvXCJcbmltcG9ydCBCaW5Ub29scyBmcm9tIFwiLi4vLi4vdXRpbHMvYmludG9vbHNcIlxuaW1wb3J0IHsgQVZNQ29uc3RhbnRzIH0gZnJvbSBcIi4vY29uc3RhbnRzXCJcbmltcG9ydCB7XG4gIE5GVFRyYW5zZmVyT3V0cHV0LFxuICBTRUNQTWludE91dHB1dCxcbiAgU0VDUFRyYW5zZmVyT3V0cHV0XG59IGZyb20gXCIuL291dHB1dHNcIlxuaW1wb3J0IHsgTkJ5dGVzIH0gZnJvbSBcIi4uLy4uL2NvbW1vbi9uYnl0ZXNcIlxuaW1wb3J0IHsgU2lnSWR4IH0gZnJvbSBcIi4uLy4uL2NvbW1vbi9jcmVkZW50aWFsc1wiXG5pbXBvcnQgeyBPdXRwdXRPd25lcnMgfSBmcm9tIFwiLi4vLi4vY29tbW9uL291dHB1dFwiXG5pbXBvcnQge1xuICBTZXJpYWxpemFibGUsXG4gIFNlcmlhbGl6YXRpb24sXG4gIFNlcmlhbGl6ZWRFbmNvZGluZyxcbiAgU2VyaWFsaXplZFR5cGVcbn0gZnJvbSBcIi4uLy4uL3V0aWxzL3NlcmlhbGl6YXRpb25cIlxuaW1wb3J0IHtcbiAgSW52YWxpZE9wZXJhdGlvbklkRXJyb3IsXG4gIENvZGVjSWRFcnJvcixcbiAgQ2hlY2tzdW1FcnJvcixcbiAgQWRkcmVzc0Vycm9yXG59IGZyb20gXCIuLi8uLi91dGlscy9lcnJvcnNcIlxuXG5jb25zdCBiaW50b29sczogQmluVG9vbHMgPSBCaW5Ub29scy5nZXRJbnN0YW5jZSgpXG5jb25zdCBzZXJpYWxpemF0aW9uOiBTZXJpYWxpemF0aW9uID0gU2VyaWFsaXphdGlvbi5nZXRJbnN0YW5jZSgpXG5jb25zdCBjYjU4OiBTZXJpYWxpemVkVHlwZSA9IFwiY2I1OFwiXG5jb25zdCBidWZmZXI6IFNlcmlhbGl6ZWRUeXBlID0gXCJCdWZmZXJcIlxuY29uc3QgaGV4OiBTZXJpYWxpemVkVHlwZSA9IFwiaGV4XCJcbmNvbnN0IGRlY2ltYWxTdHJpbmc6IFNlcmlhbGl6ZWRUeXBlID0gXCJkZWNpbWFsU3RyaW5nXCJcblxuLyoqXG4gKiBUYWtlcyBhIGJ1ZmZlciByZXByZXNlbnRpbmcgdGhlIG91dHB1dCBhbmQgcmV0dXJucyB0aGUgcHJvcGVyIFtbT3BlcmF0aW9uXV0gaW5zdGFuY2UuXG4gKlxuICogQHBhcmFtIG9waWQgQSBudW1iZXIgcmVwcmVzZW50aW5nIHRoZSBvcGVyYXRpb24gSUQgcGFyc2VkIHByaW9yIHRvIHRoZSBieXRlcyBwYXNzZWQgaW5cbiAqXG4gKiBAcmV0dXJucyBBbiBpbnN0YW5jZSBvZiBhbiBbW09wZXJhdGlvbl1dLWV4dGVuZGVkIGNsYXNzLlxuICovXG5leHBvcnQgY29uc3QgU2VsZWN0T3BlcmF0aW9uQ2xhc3MgPSAoXG4gIG9waWQ6IG51bWJlcixcbiAgLi4uYXJnczogYW55W11cbik6IE9wZXJhdGlvbiA9PiB7XG4gIGlmIChcbiAgICBvcGlkID09PSBBVk1Db25zdGFudHMuU0VDUE1JTlRPUElEIHx8XG4gICAgb3BpZCA9PT0gQVZNQ29uc3RhbnRzLlNFQ1BNSU5UT1BJRF9DT0RFQ09ORVxuICApIHtcbiAgICByZXR1cm4gbmV3IFNFQ1BNaW50T3BlcmF0aW9uKC4uLmFyZ3MpXG4gIH0gZWxzZSBpZiAoXG4gICAgb3BpZCA9PT0gQVZNQ29uc3RhbnRzLk5GVE1JTlRPUElEIHx8XG4gICAgb3BpZCA9PT0gQVZNQ29uc3RhbnRzLk5GVE1JTlRPUElEX0NPREVDT05FXG4gICkge1xuICAgIHJldHVybiBuZXcgTkZUTWludE9wZXJhdGlvbiguLi5hcmdzKVxuICB9IGVsc2UgaWYgKFxuICAgIG9waWQgPT09IEFWTUNvbnN0YW50cy5ORlRYRkVST1BJRCB8fFxuICAgIG9waWQgPT09IEFWTUNvbnN0YW50cy5ORlRYRkVST1BJRF9DT0RFQ09ORVxuICApIHtcbiAgICByZXR1cm4gbmV3IE5GVFRyYW5zZmVyT3BlcmF0aW9uKC4uLmFyZ3MpXG4gIH1cbiAgLyogaXN0YW5idWwgaWdub3JlIG5leHQgKi9cbiAgdGhyb3cgbmV3IEludmFsaWRPcGVyYXRpb25JZEVycm9yKFxuICAgIGBFcnJvciAtIFNlbGVjdE9wZXJhdGlvbkNsYXNzOiB1bmtub3duIG9waWQgJHtvcGlkfWBcbiAgKVxufVxuXG4vKipcbiAqIEEgY2xhc3MgcmVwcmVzZW50aW5nIGFuIG9wZXJhdGlvbi4gQWxsIG9wZXJhdGlvbiB0eXBlcyBtdXN0IGV4dGVuZCBvbiB0aGlzIGNsYXNzLlxuICovXG5leHBvcnQgYWJzdHJhY3QgY2xhc3MgT3BlcmF0aW9uIGV4dGVuZHMgU2VyaWFsaXphYmxlIHtcbiAgcHJvdGVjdGVkIF90eXBlTmFtZSA9IFwiT3BlcmF0aW9uXCJcbiAgcHJvdGVjdGVkIF90eXBlSUQgPSB1bmRlZmluZWRcblxuICBzZXJpYWxpemUoZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpOiBvYmplY3Qge1xuICAgIGxldCBmaWVsZHM6IG9iamVjdCA9IHN1cGVyLnNlcmlhbGl6ZShlbmNvZGluZylcbiAgICByZXR1cm4ge1xuICAgICAgLi4uZmllbGRzLFxuICAgICAgc2lnSWR4czogdGhpcy5zaWdJZHhzLm1hcCgoczogU2lnSWR4KTogb2JqZWN0ID0+IHMuc2VyaWFsaXplKGVuY29kaW5nKSlcbiAgICB9XG4gIH1cbiAgZGVzZXJpYWxpemUoZmllbGRzOiBvYmplY3QsIGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKSB7XG4gICAgc3VwZXIuZGVzZXJpYWxpemUoZmllbGRzLCBlbmNvZGluZylcbiAgICB0aGlzLnNpZ0lkeHMgPSBmaWVsZHNbXCJzaWdJZHhzXCJdLm1hcCgoczogb2JqZWN0KTogU2lnSWR4ID0+IHtcbiAgICAgIGxldCBzaWR4OiBTaWdJZHggPSBuZXcgU2lnSWR4KClcbiAgICAgIHNpZHguZGVzZXJpYWxpemUocywgZW5jb2RpbmcpXG4gICAgICByZXR1cm4gc2lkeFxuICAgIH0pXG4gICAgdGhpcy5zaWdDb3VudC53cml0ZVVJbnQzMkJFKHRoaXMuc2lnSWR4cy5sZW5ndGgsIDApXG4gIH1cblxuICBwcm90ZWN0ZWQgc2lnQ291bnQ6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYyg0KVxuICBwcm90ZWN0ZWQgc2lnSWR4czogU2lnSWR4W10gPSBbXSAvLyBpZHhzIG9mIHNpZ25lcnMgZnJvbSB1dHhvXG5cbiAgc3RhdGljIGNvbXBhcmF0b3IgPVxuICAgICgpOiAoKGE6IE9wZXJhdGlvbiwgYjogT3BlcmF0aW9uKSA9PiAxIHwgLTEgfCAwKSA9PlxuICAgIChhOiBPcGVyYXRpb24sIGI6IE9wZXJhdGlvbik6IDEgfCAtMSB8IDAgPT4ge1xuICAgICAgY29uc3QgYW91dGlkOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoNClcbiAgICAgIGFvdXRpZC53cml0ZVVJbnQzMkJFKGEuZ2V0T3BlcmF0aW9uSUQoKSwgMClcbiAgICAgIGNvbnN0IGFidWZmOiBCdWZmZXIgPSBhLnRvQnVmZmVyKClcblxuICAgICAgY29uc3QgYm91dGlkOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoNClcbiAgICAgIGJvdXRpZC53cml0ZVVJbnQzMkJFKGIuZ2V0T3BlcmF0aW9uSUQoKSwgMClcbiAgICAgIGNvbnN0IGJidWZmOiBCdWZmZXIgPSBiLnRvQnVmZmVyKClcblxuICAgICAgY29uc3QgYXNvcnQ6IEJ1ZmZlciA9IEJ1ZmZlci5jb25jYXQoXG4gICAgICAgIFthb3V0aWQsIGFidWZmXSxcbiAgICAgICAgYW91dGlkLmxlbmd0aCArIGFidWZmLmxlbmd0aFxuICAgICAgKVxuICAgICAgY29uc3QgYnNvcnQ6IEJ1ZmZlciA9IEJ1ZmZlci5jb25jYXQoXG4gICAgICAgIFtib3V0aWQsIGJidWZmXSxcbiAgICAgICAgYm91dGlkLmxlbmd0aCArIGJidWZmLmxlbmd0aFxuICAgICAgKVxuICAgICAgcmV0dXJuIEJ1ZmZlci5jb21wYXJlKGFzb3J0LCBic29ydCkgYXMgMSB8IC0xIHwgMFxuICAgIH1cblxuICBhYnN0cmFjdCBnZXRPcGVyYXRpb25JRCgpOiBudW1iZXJcblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgYXJyYXkgb2YgW1tTaWdJZHhdXSBmb3IgdGhpcyBbW09wZXJhdGlvbl1dXG4gICAqL1xuICBnZXRTaWdJZHhzID0gKCk6IFNpZ0lkeFtdID0+IHRoaXMuc2lnSWR4c1xuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBjcmVkZW50aWFsIElELlxuICAgKi9cbiAgYWJzdHJhY3QgZ2V0Q3JlZGVudGlhbElEKCk6IG51bWJlclxuXG4gIC8qKlxuICAgKiBDcmVhdGVzIGFuZCBhZGRzIGEgW1tTaWdJZHhdXSB0byB0aGUgW1tPcGVyYXRpb25dXS5cbiAgICpcbiAgICogQHBhcmFtIGFkZHJlc3NJZHggVGhlIGluZGV4IG9mIHRoZSBhZGRyZXNzIHRvIHJlZmVyZW5jZSBpbiB0aGUgc2lnbmF0dXJlc1xuICAgKiBAcGFyYW0gYWRkcmVzcyBUaGUgYWRkcmVzcyBvZiB0aGUgc291cmNlIG9mIHRoZSBzaWduYXR1cmVcbiAgICovXG4gIGFkZFNpZ25hdHVyZUlkeCA9IChhZGRyZXNzSWR4OiBudW1iZXIsIGFkZHJlc3M6IEJ1ZmZlcikgPT4ge1xuICAgIGNvbnN0IHNpZ2lkeDogU2lnSWR4ID0gbmV3IFNpZ0lkeCgpXG4gICAgY29uc3QgYjogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDQpXG4gICAgYi53cml0ZVVJbnQzMkJFKGFkZHJlc3NJZHgsIDApXG4gICAgc2lnaWR4LmZyb21CdWZmZXIoYilcbiAgICBzaWdpZHguc2V0U291cmNlKGFkZHJlc3MpXG4gICAgdGhpcy5zaWdJZHhzLnB1c2goc2lnaWR4KVxuICAgIHRoaXMuc2lnQ291bnQud3JpdGVVSW50MzJCRSh0aGlzLnNpZ0lkeHMubGVuZ3RoLCAwKVxuICB9XG5cbiAgZnJvbUJ1ZmZlcihieXRlczogQnVmZmVyLCBvZmZzZXQ6IG51bWJlciA9IDApOiBudW1iZXIge1xuICAgIHRoaXMuc2lnQ291bnQgPSBiaW50b29scy5jb3B5RnJvbShieXRlcywgb2Zmc2V0LCBvZmZzZXQgKyA0KVxuICAgIG9mZnNldCArPSA0XG4gICAgY29uc3Qgc2lnQ291bnQ6IG51bWJlciA9IHRoaXMuc2lnQ291bnQucmVhZFVJbnQzMkJFKDApXG4gICAgdGhpcy5zaWdJZHhzID0gW11cbiAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgc2lnQ291bnQ7IGkrKykge1xuICAgICAgY29uc3Qgc2lnaWR4OiBTaWdJZHggPSBuZXcgU2lnSWR4KClcbiAgICAgIGNvbnN0IHNpZ2J1ZmY6IEJ1ZmZlciA9IGJpbnRvb2xzLmNvcHlGcm9tKGJ5dGVzLCBvZmZzZXQsIG9mZnNldCArIDQpXG4gICAgICBzaWdpZHguZnJvbUJ1ZmZlcihzaWdidWZmKVxuICAgICAgb2Zmc2V0ICs9IDRcbiAgICAgIHRoaXMuc2lnSWR4cy5wdXNoKHNpZ2lkeClcbiAgICB9XG4gICAgcmV0dXJuIG9mZnNldFxuICB9XG5cbiAgdG9CdWZmZXIoKTogQnVmZmVyIHtcbiAgICB0aGlzLnNpZ0NvdW50LndyaXRlVUludDMyQkUodGhpcy5zaWdJZHhzLmxlbmd0aCwgMClcbiAgICBsZXQgYnNpemU6IG51bWJlciA9IHRoaXMuc2lnQ291bnQubGVuZ3RoXG4gICAgY29uc3QgYmFycjogQnVmZmVyW10gPSBbdGhpcy5zaWdDb3VudF1cbiAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgdGhpcy5zaWdJZHhzLmxlbmd0aDsgaSsrKSB7XG4gICAgICBjb25zdCBiOiBCdWZmZXIgPSB0aGlzLnNpZ0lkeHNbYCR7aX1gXS50b0J1ZmZlcigpXG4gICAgICBiYXJyLnB1c2goYilcbiAgICAgIGJzaXplICs9IGIubGVuZ3RoXG4gICAgfVxuICAgIHJldHVybiBCdWZmZXIuY29uY2F0KGJhcnIsIGJzaXplKVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSBiYXNlLTU4IHN0cmluZyByZXByZXNlbnRpbmcgdGhlIFtbTkZUTWludE9wZXJhdGlvbl1dLlxuICAgKi9cbiAgdG9TdHJpbmcoKTogc3RyaW5nIHtcbiAgICByZXR1cm4gYmludG9vbHMuYnVmZmVyVG9CNTgodGhpcy50b0J1ZmZlcigpKVxuICB9XG59XG5cbi8qKlxuICogQSBjbGFzcyB3aGljaCBjb250YWlucyBhbiBbW09wZXJhdGlvbl1dIGZvciB0cmFuc2ZlcnMuXG4gKlxuICovXG5leHBvcnQgY2xhc3MgVHJhbnNmZXJhYmxlT3BlcmF0aW9uIGV4dGVuZHMgU2VyaWFsaXphYmxlIHtcbiAgcHJvdGVjdGVkIF90eXBlTmFtZSA9IFwiVHJhbnNmZXJhYmxlT3BlcmF0aW9uXCJcbiAgcHJvdGVjdGVkIF90eXBlSUQgPSB1bmRlZmluZWRcblxuICBzZXJpYWxpemUoZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpOiBvYmplY3Qge1xuICAgIGxldCBmaWVsZHM6IG9iamVjdCA9IHN1cGVyLnNlcmlhbGl6ZShlbmNvZGluZylcbiAgICByZXR1cm4ge1xuICAgICAgLi4uZmllbGRzLFxuICAgICAgYXNzZXRJRDogc2VyaWFsaXphdGlvbi5lbmNvZGVyKHRoaXMuYXNzZXRJRCwgZW5jb2RpbmcsIGJ1ZmZlciwgY2I1OCwgMzIpLFxuICAgICAgdXR4b0lEczogdGhpcy51dHhvSURzLm1hcCgodSkgPT4gdS5zZXJpYWxpemUoZW5jb2RpbmcpKSxcbiAgICAgIG9wZXJhdGlvbjogdGhpcy5vcGVyYXRpb24uc2VyaWFsaXplKGVuY29kaW5nKVxuICAgIH1cbiAgfVxuICBkZXNlcmlhbGl6ZShmaWVsZHM6IG9iamVjdCwgZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpIHtcbiAgICBzdXBlci5kZXNlcmlhbGl6ZShmaWVsZHMsIGVuY29kaW5nKVxuICAgIHRoaXMuYXNzZXRJRCA9IHNlcmlhbGl6YXRpb24uZGVjb2RlcihcbiAgICAgIGZpZWxkc1tcImFzc2V0SURcIl0sXG4gICAgICBlbmNvZGluZyxcbiAgICAgIGNiNTgsXG4gICAgICBidWZmZXIsXG4gICAgICAzMlxuICAgIClcbiAgICB0aGlzLnV0eG9JRHMgPSBmaWVsZHNbXCJ1dHhvSURzXCJdLm1hcCgodTogb2JqZWN0KSA9PiB7XG4gICAgICBsZXQgdXR4b2lkOiBVVFhPSUQgPSBuZXcgVVRYT0lEKClcbiAgICAgIHV0eG9pZC5kZXNlcmlhbGl6ZSh1LCBlbmNvZGluZylcbiAgICAgIHJldHVybiB1dHhvaWRcbiAgICB9KVxuICAgIHRoaXMub3BlcmF0aW9uID0gU2VsZWN0T3BlcmF0aW9uQ2xhc3MoZmllbGRzW1wib3BlcmF0aW9uXCJdW1wiX3R5cGVJRFwiXSlcbiAgICB0aGlzLm9wZXJhdGlvbi5kZXNlcmlhbGl6ZShmaWVsZHNbXCJvcGVyYXRpb25cIl0sIGVuY29kaW5nKVxuICB9XG5cbiAgcHJvdGVjdGVkIGFzc2V0SUQ6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYygzMilcbiAgcHJvdGVjdGVkIHV0eG9JRHM6IFVUWE9JRFtdID0gW11cbiAgcHJvdGVjdGVkIG9wZXJhdGlvbjogT3BlcmF0aW9uXG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSBmdW5jdGlvbiB1c2VkIHRvIHNvcnQgYW4gYXJyYXkgb2YgW1tUcmFuc2ZlcmFibGVPcGVyYXRpb25dXXNcbiAgICovXG4gIHN0YXRpYyBjb21wYXJhdG9yID0gKCk6ICgoXG4gICAgYTogVHJhbnNmZXJhYmxlT3BlcmF0aW9uLFxuICAgIGI6IFRyYW5zZmVyYWJsZU9wZXJhdGlvblxuICApID0+IDEgfCAtMSB8IDApID0+IHtcbiAgICByZXR1cm4gZnVuY3Rpb24gKFxuICAgICAgYTogVHJhbnNmZXJhYmxlT3BlcmF0aW9uLFxuICAgICAgYjogVHJhbnNmZXJhYmxlT3BlcmF0aW9uXG4gICAgKTogMSB8IC0xIHwgMCB7XG4gICAgICByZXR1cm4gQnVmZmVyLmNvbXBhcmUoYS50b0J1ZmZlcigpLCBiLnRvQnVmZmVyKCkpIGFzIDEgfCAtMSB8IDBcbiAgICB9XG4gIH1cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIGFzc2V0SUQgYXMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfS5cbiAgICovXG4gIGdldEFzc2V0SUQgPSAoKTogQnVmZmVyID0+IHRoaXMuYXNzZXRJRFxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIGFuIGFycmF5IG9mIFVUWE9JRHMgaW4gdGhpcyBvcGVyYXRpb24uXG4gICAqL1xuICBnZXRVVFhPSURzID0gKCk6IFVUWE9JRFtdID0+IHRoaXMudXR4b0lEc1xuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBvcGVyYXRpb25cbiAgICovXG4gIGdldE9wZXJhdGlvbiA9ICgpOiBPcGVyYXRpb24gPT4gdGhpcy5vcGVyYXRpb25cblxuICBmcm9tQnVmZmVyKGJ5dGVzOiBCdWZmZXIsIG9mZnNldDogbnVtYmVyID0gMCk6IG51bWJlciB7XG4gICAgdGhpcy5hc3NldElEID0gYmludG9vbHMuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgMzIpXG4gICAgb2Zmc2V0ICs9IDMyXG4gICAgY29uc3QgbnVtdXR4b0lEczogbnVtYmVyID0gYmludG9vbHNcbiAgICAgIC5jb3B5RnJvbShieXRlcywgb2Zmc2V0LCBvZmZzZXQgKyA0KVxuICAgICAgLnJlYWRVSW50MzJCRSgwKVxuICAgIG9mZnNldCArPSA0XG4gICAgdGhpcy51dHhvSURzID0gW11cbiAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgbnVtdXR4b0lEczsgaSsrKSB7XG4gICAgICBjb25zdCB1dHhvaWQ6IFVUWE9JRCA9IG5ldyBVVFhPSUQoKVxuICAgICAgb2Zmc2V0ID0gdXR4b2lkLmZyb21CdWZmZXIoYnl0ZXMsIG9mZnNldClcbiAgICAgIHRoaXMudXR4b0lEcy5wdXNoKHV0eG9pZClcbiAgICB9XG4gICAgY29uc3Qgb3BpZDogbnVtYmVyID0gYmludG9vbHNcbiAgICAgIC5jb3B5RnJvbShieXRlcywgb2Zmc2V0LCBvZmZzZXQgKyA0KVxuICAgICAgLnJlYWRVSW50MzJCRSgwKVxuICAgIG9mZnNldCArPSA0XG4gICAgdGhpcy5vcGVyYXRpb24gPSBTZWxlY3RPcGVyYXRpb25DbGFzcyhvcGlkKVxuICAgIHJldHVybiB0aGlzLm9wZXJhdGlvbi5mcm9tQnVmZmVyKGJ5dGVzLCBvZmZzZXQpXG4gIH1cblxuICB0b0J1ZmZlcigpOiBCdWZmZXIge1xuICAgIGNvbnN0IG51bXV0eG9JRHMgPSBCdWZmZXIuYWxsb2MoNClcbiAgICBudW11dHhvSURzLndyaXRlVUludDMyQkUodGhpcy51dHhvSURzLmxlbmd0aCwgMClcbiAgICBsZXQgYnNpemU6IG51bWJlciA9IHRoaXMuYXNzZXRJRC5sZW5ndGggKyBudW11dHhvSURzLmxlbmd0aFxuICAgIGNvbnN0IGJhcnI6IEJ1ZmZlcltdID0gW3RoaXMuYXNzZXRJRCwgbnVtdXR4b0lEc11cbiAgICB0aGlzLnV0eG9JRHMgPSB0aGlzLnV0eG9JRHMuc29ydChVVFhPSUQuY29tcGFyYXRvcigpKVxuICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCB0aGlzLnV0eG9JRHMubGVuZ3RoOyBpKyspIHtcbiAgICAgIGNvbnN0IGI6IEJ1ZmZlciA9IHRoaXMudXR4b0lEc1tgJHtpfWBdLnRvQnVmZmVyKClcbiAgICAgIGJhcnIucHVzaChiKVxuICAgICAgYnNpemUgKz0gYi5sZW5ndGhcbiAgICB9XG4gICAgY29uc3Qgb3BpZDogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDQpXG4gICAgb3BpZC53cml0ZVVJbnQzMkJFKHRoaXMub3BlcmF0aW9uLmdldE9wZXJhdGlvbklEKCksIDApXG4gICAgYmFyci5wdXNoKG9waWQpXG4gICAgYnNpemUgKz0gb3BpZC5sZW5ndGhcbiAgICBjb25zdCBiOiBCdWZmZXIgPSB0aGlzLm9wZXJhdGlvbi50b0J1ZmZlcigpXG4gICAgYnNpemUgKz0gYi5sZW5ndGhcbiAgICBiYXJyLnB1c2goYilcbiAgICByZXR1cm4gQnVmZmVyLmNvbmNhdChiYXJyLCBic2l6ZSlcbiAgfVxuXG4gIGNvbnN0cnVjdG9yKFxuICAgIGFzc2V0SUQ6IEJ1ZmZlciA9IHVuZGVmaW5lZCxcbiAgICB1dHhvaWRzOiBVVFhPSURbXSB8IHN0cmluZ1tdIHwgQnVmZmVyW10gPSB1bmRlZmluZWQsXG4gICAgb3BlcmF0aW9uOiBPcGVyYXRpb24gPSB1bmRlZmluZWRcbiAgKSB7XG4gICAgc3VwZXIoKVxuICAgIGlmIChcbiAgICAgIHR5cGVvZiBhc3NldElEICE9PSBcInVuZGVmaW5lZFwiICYmXG4gICAgICBhc3NldElELmxlbmd0aCA9PT0gQVZNQ29uc3RhbnRzLkFTU0VUSURMRU4gJiZcbiAgICAgIG9wZXJhdGlvbiBpbnN0YW5jZW9mIE9wZXJhdGlvbiAmJlxuICAgICAgdHlwZW9mIHV0eG9pZHMgIT09IFwidW5kZWZpbmVkXCIgJiZcbiAgICAgIEFycmF5LmlzQXJyYXkodXR4b2lkcylcbiAgICApIHtcbiAgICAgIHRoaXMuYXNzZXRJRCA9IGFzc2V0SURcbiAgICAgIHRoaXMub3BlcmF0aW9uID0gb3BlcmF0aW9uXG4gICAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgdXR4b2lkcy5sZW5ndGg7IGkrKykge1xuICAgICAgICBjb25zdCB1dHhvaWQ6IFVUWE9JRCA9IG5ldyBVVFhPSUQoKVxuICAgICAgICBpZiAodHlwZW9mIHV0eG9pZHNbYCR7aX1gXSA9PT0gXCJzdHJpbmdcIikge1xuICAgICAgICAgIHV0eG9pZC5mcm9tU3RyaW5nKHV0eG9pZHNbYCR7aX1gXSBhcyBzdHJpbmcpXG4gICAgICAgIH0gZWxzZSBpZiAodXR4b2lkc1tgJHtpfWBdIGluc3RhbmNlb2YgQnVmZmVyKSB7XG4gICAgICAgICAgdXR4b2lkLmZyb21CdWZmZXIodXR4b2lkc1tgJHtpfWBdIGFzIEJ1ZmZlcilcbiAgICAgICAgfSBlbHNlIGlmICh1dHhvaWRzW2Ake2l9YF0gaW5zdGFuY2VvZiBVVFhPSUQpIHtcbiAgICAgICAgICB1dHhvaWQuZnJvbVN0cmluZyh1dHhvaWRzW2Ake2l9YF0udG9TdHJpbmcoKSkgLy8gY2xvbmVcbiAgICAgICAgfVxuICAgICAgICB0aGlzLnV0eG9JRHMucHVzaCh1dHhvaWQpXG4gICAgICB9XG4gICAgfVxuICB9XG59XG5cbi8qKlxuICogQW4gW1tPcGVyYXRpb25dXSBjbGFzcyB3aGljaCBzcGVjaWZpZXMgYSBTRUNQMjU2azEgTWludCBPcC5cbiAqL1xuZXhwb3J0IGNsYXNzIFNFQ1BNaW50T3BlcmF0aW9uIGV4dGVuZHMgT3BlcmF0aW9uIHtcbiAgcHJvdGVjdGVkIF90eXBlTmFtZSA9IFwiU0VDUE1pbnRPcGVyYXRpb25cIlxuICBwcm90ZWN0ZWQgX2NvZGVjSUQgPSBBVk1Db25zdGFudHMuTEFURVNUQ09ERUNcbiAgcHJvdGVjdGVkIF90eXBlSUQgPVxuICAgIHRoaXMuX2NvZGVjSUQgPT09IDBcbiAgICAgID8gQVZNQ29uc3RhbnRzLlNFQ1BNSU5UT1BJRFxuICAgICAgOiBBVk1Db25zdGFudHMuU0VDUE1JTlRPUElEX0NPREVDT05FXG5cbiAgc2VyaWFsaXplKGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKTogb2JqZWN0IHtcbiAgICBsZXQgZmllbGRzOiBvYmplY3QgPSBzdXBlci5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgcmV0dXJuIHtcbiAgICAgIC4uLmZpZWxkcyxcbiAgICAgIG1pbnRPdXRwdXQ6IHRoaXMubWludE91dHB1dC5zZXJpYWxpemUoZW5jb2RpbmcpLFxuICAgICAgdHJhbnNmZXJPdXRwdXRzOiB0aGlzLnRyYW5zZmVyT3V0cHV0LnNlcmlhbGl6ZShlbmNvZGluZylcbiAgICB9XG4gIH1cbiAgZGVzZXJpYWxpemUoZmllbGRzOiBvYmplY3QsIGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKSB7XG4gICAgc3VwZXIuZGVzZXJpYWxpemUoZmllbGRzLCBlbmNvZGluZylcbiAgICB0aGlzLm1pbnRPdXRwdXQgPSBuZXcgU0VDUE1pbnRPdXRwdXQoKVxuICAgIHRoaXMubWludE91dHB1dC5kZXNlcmlhbGl6ZShmaWVsZHNbXCJtaW50T3V0cHV0XCJdLCBlbmNvZGluZylcbiAgICB0aGlzLnRyYW5zZmVyT3V0cHV0ID0gbmV3IFNFQ1BUcmFuc2Zlck91dHB1dCgpXG4gICAgdGhpcy50cmFuc2Zlck91dHB1dC5kZXNlcmlhbGl6ZShmaWVsZHNbXCJ0cmFuc2Zlck91dHB1dHNcIl0sIGVuY29kaW5nKVxuICB9XG5cbiAgcHJvdGVjdGVkIG1pbnRPdXRwdXQ6IFNFQ1BNaW50T3V0cHV0ID0gdW5kZWZpbmVkXG4gIHByb3RlY3RlZCB0cmFuc2Zlck91dHB1dDogU0VDUFRyYW5zZmVyT3V0cHV0ID0gdW5kZWZpbmVkXG5cbiAgLyoqXG4gICAqIFNldCB0aGUgY29kZWNJRFxuICAgKlxuICAgKiBAcGFyYW0gY29kZWNJRCBUaGUgY29kZWNJRCB0byBzZXRcbiAgICovXG4gIHNldENvZGVjSUQoY29kZWNJRDogbnVtYmVyKTogdm9pZCB7XG4gICAgaWYgKGNvZGVjSUQgIT09IDAgJiYgY29kZWNJRCAhPT0gMSkge1xuICAgICAgLyogaXN0YW5idWwgaWdub3JlIG5leHQgKi9cbiAgICAgIHRocm93IG5ldyBDb2RlY0lkRXJyb3IoXG4gICAgICAgIFwiRXJyb3IgLSBTRUNQTWludE9wZXJhdGlvbi5zZXRDb2RlY0lEOiBpbnZhbGlkIGNvZGVjSUQuIFZhbGlkIGNvZGVjSURzIGFyZSAwIGFuZCAxLlwiXG4gICAgICApXG4gICAgfVxuICAgIHRoaXMuX2NvZGVjSUQgPSBjb2RlY0lEXG4gICAgdGhpcy5fdHlwZUlEID1cbiAgICAgIHRoaXMuX2NvZGVjSUQgPT09IDBcbiAgICAgICAgPyBBVk1Db25zdGFudHMuU0VDUE1JTlRPUElEXG4gICAgICAgIDogQVZNQ29uc3RhbnRzLlNFQ1BNSU5UT1BJRF9DT0RFQ09ORVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIG9wZXJhdGlvbiBJRC5cbiAgICovXG4gIGdldE9wZXJhdGlvbklEKCk6IG51bWJlciB7XG4gICAgcmV0dXJuIHRoaXMuX3R5cGVJRFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIGNyZWRlbnRpYWwgSUQuXG4gICAqL1xuICBnZXRDcmVkZW50aWFsSUQoKTogbnVtYmVyIHtcbiAgICBpZiAodGhpcy5fY29kZWNJRCA9PT0gMCkge1xuICAgICAgcmV0dXJuIEFWTUNvbnN0YW50cy5TRUNQQ1JFREVOVElBTFxuICAgIH0gZWxzZSBpZiAodGhpcy5fY29kZWNJRCA9PT0gMSkge1xuICAgICAgcmV0dXJuIEFWTUNvbnN0YW50cy5TRUNQQ1JFREVOVElBTF9DT0RFQ09ORVxuICAgIH1cbiAgfVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBbW1NFQ1BNaW50T3V0cHV0XV0gdG8gYmUgcHJvZHVjZWQgYnkgdGhpcyBvcGVyYXRpb24uXG4gICAqL1xuICBnZXRNaW50T3V0cHV0KCk6IFNFQ1BNaW50T3V0cHV0IHtcbiAgICByZXR1cm4gdGhpcy5taW50T3V0cHV0XG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyBbW1NFQ1BUcmFuc2Zlck91dHB1dF1dIHRvIGJlIHByb2R1Y2VkIGJ5IHRoaXMgb3BlcmF0aW9uLlxuICAgKi9cbiAgZ2V0VHJhbnNmZXJPdXRwdXQoKTogU0VDUFRyYW5zZmVyT3V0cHV0IHtcbiAgICByZXR1cm4gdGhpcy50cmFuc2Zlck91dHB1dFxuICB9XG5cbiAgLyoqXG4gICAqIFBvcHVhdGVzIHRoZSBpbnN0YW5jZSBmcm9tIGEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gcmVwcmVzZW50aW5nIHRoZSBbW1NFQ1BNaW50T3BlcmF0aW9uXV0gYW5kIHJldHVybnMgdGhlIHVwZGF0ZWQgb2Zmc2V0LlxuICAgKi9cbiAgZnJvbUJ1ZmZlcihieXRlczogQnVmZmVyLCBvZmZzZXQ6IG51bWJlciA9IDApOiBudW1iZXIge1xuICAgIG9mZnNldCA9IHN1cGVyLmZyb21CdWZmZXIoYnl0ZXMsIG9mZnNldClcbiAgICB0aGlzLm1pbnRPdXRwdXQgPSBuZXcgU0VDUE1pbnRPdXRwdXQoKVxuICAgIG9mZnNldCA9IHRoaXMubWludE91dHB1dC5mcm9tQnVmZmVyKGJ5dGVzLCBvZmZzZXQpXG4gICAgdGhpcy50cmFuc2Zlck91dHB1dCA9IG5ldyBTRUNQVHJhbnNmZXJPdXRwdXQoKVxuICAgIG9mZnNldCA9IHRoaXMudHJhbnNmZXJPdXRwdXQuZnJvbUJ1ZmZlcihieXRlcywgb2Zmc2V0KVxuICAgIHJldHVybiBvZmZzZXRcbiAgfVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBidWZmZXIgcmVwcmVzZW50aW5nIHRoZSBbW1NFQ1BNaW50T3BlcmF0aW9uXV0gaW5zdGFuY2UuXG4gICAqL1xuICB0b0J1ZmZlcigpOiBCdWZmZXIge1xuICAgIGNvbnN0IHN1cGVyYnVmZjogQnVmZmVyID0gc3VwZXIudG9CdWZmZXIoKVxuICAgIGNvbnN0IG1pbnRvdXRCdWZmOiBCdWZmZXIgPSB0aGlzLm1pbnRPdXRwdXQudG9CdWZmZXIoKVxuICAgIGNvbnN0IHRyYW5zZmVyT3V0QnVmZjogQnVmZmVyID0gdGhpcy50cmFuc2Zlck91dHB1dC50b0J1ZmZlcigpXG4gICAgY29uc3QgYnNpemU6IG51bWJlciA9XG4gICAgICBzdXBlcmJ1ZmYubGVuZ3RoICsgbWludG91dEJ1ZmYubGVuZ3RoICsgdHJhbnNmZXJPdXRCdWZmLmxlbmd0aFxuXG4gICAgY29uc3QgYmFycjogQnVmZmVyW10gPSBbc3VwZXJidWZmLCBtaW50b3V0QnVmZiwgdHJhbnNmZXJPdXRCdWZmXVxuXG4gICAgcmV0dXJuIEJ1ZmZlci5jb25jYXQoYmFyciwgYnNpemUpXG4gIH1cblxuICAvKipcbiAgICogQW4gW1tPcGVyYXRpb25dXSBjbGFzcyB3aGljaCBtaW50cyBuZXcgdG9rZW5zIG9uIGFuIGFzc2V0SUQuXG4gICAqXG4gICAqIEBwYXJhbSBtaW50T3V0cHV0IFRoZSBbW1NFQ1BNaW50T3V0cHV0XV0gdGhhdCB3aWxsIGJlIHByb2R1Y2VkIGJ5IHRoaXMgdHJhbnNhY3Rpb24uXG4gICAqIEBwYXJhbSB0cmFuc2Zlck91dHB1dCBBIFtbU0VDUFRyYW5zZmVyT3V0cHV0XV0gdGhhdCB3aWxsIGJlIHByb2R1Y2VkIGZyb20gdGhpcyBtaW50aW5nIG9wZXJhdGlvbi5cbiAgICovXG4gIGNvbnN0cnVjdG9yKFxuICAgIG1pbnRPdXRwdXQ6IFNFQ1BNaW50T3V0cHV0ID0gdW5kZWZpbmVkLFxuICAgIHRyYW5zZmVyT3V0cHV0OiBTRUNQVHJhbnNmZXJPdXRwdXQgPSB1bmRlZmluZWRcbiAgKSB7XG4gICAgc3VwZXIoKVxuICAgIGlmICh0eXBlb2YgbWludE91dHB1dCAhPT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgdGhpcy5taW50T3V0cHV0ID0gbWludE91dHB1dFxuICAgIH1cbiAgICBpZiAodHlwZW9mIHRyYW5zZmVyT3V0cHV0ICE9PSBcInVuZGVmaW5lZFwiKSB7XG4gICAgICB0aGlzLnRyYW5zZmVyT3V0cHV0ID0gdHJhbnNmZXJPdXRwdXRcbiAgICB9XG4gIH1cbn1cblxuLyoqXG4gKiBBbiBbW09wZXJhdGlvbl1dIGNsYXNzIHdoaWNoIHNwZWNpZmllcyBhIE5GVCBNaW50IE9wLlxuICovXG5leHBvcnQgY2xhc3MgTkZUTWludE9wZXJhdGlvbiBleHRlbmRzIE9wZXJhdGlvbiB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIk5GVE1pbnRPcGVyYXRpb25cIlxuICBwcm90ZWN0ZWQgX2NvZGVjSUQgPSBBVk1Db25zdGFudHMuTEFURVNUQ09ERUNcbiAgcHJvdGVjdGVkIF90eXBlSUQgPVxuICAgIHRoaXMuX2NvZGVjSUQgPT09IDBcbiAgICAgID8gQVZNQ29uc3RhbnRzLk5GVE1JTlRPUElEXG4gICAgICA6IEFWTUNvbnN0YW50cy5ORlRNSU5UT1BJRF9DT0RFQ09ORVxuXG4gIHNlcmlhbGl6ZShlbmNvZGluZzogU2VyaWFsaXplZEVuY29kaW5nID0gXCJoZXhcIik6IG9iamVjdCB7XG4gICAgY29uc3QgZmllbGRzOiBvYmplY3QgPSBzdXBlci5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgcmV0dXJuIHtcbiAgICAgIC4uLmZpZWxkcyxcbiAgICAgIGdyb3VwSUQ6IHNlcmlhbGl6YXRpb24uZW5jb2RlcihcbiAgICAgICAgdGhpcy5ncm91cElELFxuICAgICAgICBlbmNvZGluZyxcbiAgICAgICAgYnVmZmVyLFxuICAgICAgICBkZWNpbWFsU3RyaW5nLFxuICAgICAgICA0XG4gICAgICApLFxuICAgICAgcGF5bG9hZDogc2VyaWFsaXphdGlvbi5lbmNvZGVyKHRoaXMucGF5bG9hZCwgZW5jb2RpbmcsIGJ1ZmZlciwgaGV4KSxcbiAgICAgIG91dHB1dE93bmVyczogdGhpcy5vdXRwdXRPd25lcnMubWFwKChvKSA9PiBvLnNlcmlhbGl6ZShlbmNvZGluZykpXG4gICAgfVxuICB9XG4gIGRlc2VyaWFsaXplKGZpZWxkczogb2JqZWN0LCBlbmNvZGluZzogU2VyaWFsaXplZEVuY29kaW5nID0gXCJoZXhcIikge1xuICAgIHN1cGVyLmRlc2VyaWFsaXplKGZpZWxkcywgZW5jb2RpbmcpXG4gICAgdGhpcy5ncm91cElEID0gc2VyaWFsaXphdGlvbi5kZWNvZGVyKFxuICAgICAgZmllbGRzW1wiZ3JvdXBJRFwiXSxcbiAgICAgIGVuY29kaW5nLFxuICAgICAgZGVjaW1hbFN0cmluZyxcbiAgICAgIGJ1ZmZlcixcbiAgICAgIDRcbiAgICApXG4gICAgdGhpcy5wYXlsb2FkID0gc2VyaWFsaXphdGlvbi5kZWNvZGVyKFxuICAgICAgZmllbGRzW1wicGF5bG9hZFwiXSxcbiAgICAgIGVuY29kaW5nLFxuICAgICAgaGV4LFxuICAgICAgYnVmZmVyXG4gICAgKVxuICAgIC8vIHRoaXMub3V0cHV0T3duZXJzID0gZmllbGRzW1wib3V0cHV0T3duZXJzXCJdLm1hcCgobzogTkZUTWludE91dHB1dCkgPT4ge1xuICAgIC8vICAgbGV0IG9vOiBORlRNaW50T3V0cHV0ID0gbmV3IE5GVE1pbnRPdXRwdXQoKVxuICAgIC8vICAgb28uZGVzZXJpYWxpemUobywgZW5jb2RpbmcpXG4gICAgLy8gICByZXR1cm4gb29cbiAgICAvLyB9KVxuICAgIHRoaXMub3V0cHV0T3duZXJzID0gZmllbGRzW1wib3V0cHV0T3duZXJzXCJdLm1hcChcbiAgICAgIChvOiBvYmplY3QpOiBPdXRwdXRPd25lcnMgPT4ge1xuICAgICAgICBsZXQgb286IE91dHB1dE93bmVycyA9IG5ldyBPdXRwdXRPd25lcnMoKVxuICAgICAgICBvby5kZXNlcmlhbGl6ZShvLCBlbmNvZGluZylcbiAgICAgICAgcmV0dXJuIG9vXG4gICAgICB9XG4gICAgKVxuICB9XG5cbiAgcHJvdGVjdGVkIGdyb3VwSUQ6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYyg0KVxuICBwcm90ZWN0ZWQgcGF5bG9hZDogQnVmZmVyXG4gIHByb3RlY3RlZCBvdXRwdXRPd25lcnM6IE91dHB1dE93bmVyc1tdID0gW11cblxuICAvKipcbiAgICogU2V0IHRoZSBjb2RlY0lEXG4gICAqXG4gICAqIEBwYXJhbSBjb2RlY0lEIFRoZSBjb2RlY0lEIHRvIHNldFxuICAgKi9cbiAgc2V0Q29kZWNJRChjb2RlY0lEOiBudW1iZXIpOiB2b2lkIHtcbiAgICBpZiAoY29kZWNJRCAhPT0gMCAmJiBjb2RlY0lEICE9PSAxKSB7XG4gICAgICAvKiBpc3RhbmJ1bCBpZ25vcmUgbmV4dCAqL1xuICAgICAgdGhyb3cgbmV3IENvZGVjSWRFcnJvcihcbiAgICAgICAgXCJFcnJvciAtIE5GVE1pbnRPcGVyYXRpb24uc2V0Q29kZWNJRDogaW52YWxpZCBjb2RlY0lELiBWYWxpZCBjb2RlY0lEcyBhcmUgMCBhbmQgMS5cIlxuICAgICAgKVxuICAgIH1cbiAgICB0aGlzLl9jb2RlY0lEID0gY29kZWNJRFxuICAgIHRoaXMuX3R5cGVJRCA9XG4gICAgICB0aGlzLl9jb2RlY0lEID09PSAwXG4gICAgICAgID8gQVZNQ29uc3RhbnRzLk5GVE1JTlRPUElEXG4gICAgICAgIDogQVZNQ29uc3RhbnRzLk5GVE1JTlRPUElEX0NPREVDT05FXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgb3BlcmF0aW9uIElELlxuICAgKi9cbiAgZ2V0T3BlcmF0aW9uSUQoKTogbnVtYmVyIHtcbiAgICByZXR1cm4gdGhpcy5fdHlwZUlEXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgY3JlZGVudGlhbCBJRC5cbiAgICovXG4gIGdldENyZWRlbnRpYWxJRCA9ICgpOiBudW1iZXIgPT4ge1xuICAgIGlmICh0aGlzLl9jb2RlY0lEID09PSAwKSB7XG4gICAgICByZXR1cm4gQVZNQ29uc3RhbnRzLk5GVENSRURFTlRJQUxcbiAgICB9IGVsc2UgaWYgKHRoaXMuX2NvZGVjSUQgPT09IDEpIHtcbiAgICAgIHJldHVybiBBVk1Db25zdGFudHMuTkZUQ1JFREVOVElBTF9DT0RFQ09ORVxuICAgIH1cbiAgfVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBwYXlsb2FkLlxuICAgKi9cbiAgZ2V0R3JvdXBJRCA9ICgpOiBCdWZmZXIgPT4ge1xuICAgIHJldHVybiBiaW50b29scy5jb3B5RnJvbSh0aGlzLmdyb3VwSUQsIDApXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgcGF5bG9hZC5cbiAgICovXG4gIGdldFBheWxvYWQgPSAoKTogQnVmZmVyID0+IHtcbiAgICByZXR1cm4gYmludG9vbHMuY29weUZyb20odGhpcy5wYXlsb2FkLCAwKVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIHBheWxvYWQncyByYXcge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gd2l0aCBsZW5ndGggcHJlcGVuZGVkLCBmb3IgdXNlIHdpdGggW1tQYXlsb2FkQmFzZV1dJ3MgZnJvbUJ1ZmZlclxuICAgKi9cbiAgZ2V0UGF5bG9hZEJ1ZmZlciA9ICgpOiBCdWZmZXIgPT4ge1xuICAgIGxldCBwYXlsb2FkbGVuOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoNClcbiAgICBwYXlsb2FkbGVuLndyaXRlVUludDMyQkUodGhpcy5wYXlsb2FkLmxlbmd0aCwgMClcbiAgICByZXR1cm4gQnVmZmVyLmNvbmNhdChbcGF5bG9hZGxlbiwgYmludG9vbHMuY29weUZyb20odGhpcy5wYXlsb2FkLCAwKV0pXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgb3V0cHV0T3duZXJzLlxuICAgKi9cbiAgZ2V0T3V0cHV0T3duZXJzID0gKCk6IE91dHB1dE93bmVyc1tdID0+IHtcbiAgICByZXR1cm4gdGhpcy5vdXRwdXRPd25lcnNcbiAgfVxuXG4gIC8qKlxuICAgKiBQb3B1YXRlcyB0aGUgaW5zdGFuY2UgZnJvbSBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IHJlcHJlc2VudGluZyB0aGUgW1tORlRNaW50T3BlcmF0aW9uXV0gYW5kIHJldHVybnMgdGhlIHVwZGF0ZWQgb2Zmc2V0LlxuICAgKi9cbiAgZnJvbUJ1ZmZlcihieXRlczogQnVmZmVyLCBvZmZzZXQ6IG51bWJlciA9IDApOiBudW1iZXIge1xuICAgIG9mZnNldCA9IHN1cGVyLmZyb21CdWZmZXIoYnl0ZXMsIG9mZnNldClcbiAgICB0aGlzLmdyb3VwSUQgPSBiaW50b29scy5jb3B5RnJvbShieXRlcywgb2Zmc2V0LCBvZmZzZXQgKyA0KVxuICAgIG9mZnNldCArPSA0XG4gICAgbGV0IHBheWxvYWRMZW46IG51bWJlciA9IGJpbnRvb2xzXG4gICAgICAuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgNClcbiAgICAgIC5yZWFkVUludDMyQkUoMClcbiAgICBvZmZzZXQgKz0gNFxuICAgIHRoaXMucGF5bG9hZCA9IGJpbnRvb2xzLmNvcHlGcm9tKGJ5dGVzLCBvZmZzZXQsIG9mZnNldCArIHBheWxvYWRMZW4pXG4gICAgb2Zmc2V0ICs9IHBheWxvYWRMZW5cbiAgICBsZXQgbnVtb3V0cHV0czogbnVtYmVyID0gYmludG9vbHNcbiAgICAgIC5jb3B5RnJvbShieXRlcywgb2Zmc2V0LCBvZmZzZXQgKyA0KVxuICAgICAgLnJlYWRVSW50MzJCRSgwKVxuICAgIG9mZnNldCArPSA0XG4gICAgdGhpcy5vdXRwdXRPd25lcnMgPSBbXVxuICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCBudW1vdXRwdXRzOyBpKyspIHtcbiAgICAgIGxldCBvdXRwdXRPd25lcjogT3V0cHV0T3duZXJzID0gbmV3IE91dHB1dE93bmVycygpXG4gICAgICBvZmZzZXQgPSBvdXRwdXRPd25lci5mcm9tQnVmZmVyKGJ5dGVzLCBvZmZzZXQpXG4gICAgICB0aGlzLm91dHB1dE93bmVycy5wdXNoKG91dHB1dE93bmVyKVxuICAgIH1cbiAgICByZXR1cm4gb2Zmc2V0XG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgYnVmZmVyIHJlcHJlc2VudGluZyB0aGUgW1tORlRNaW50T3BlcmF0aW9uXV0gaW5zdGFuY2UuXG4gICAqL1xuICB0b0J1ZmZlcigpOiBCdWZmZXIge1xuICAgIGNvbnN0IHN1cGVyYnVmZjogQnVmZmVyID0gc3VwZXIudG9CdWZmZXIoKVxuICAgIGNvbnN0IHBheWxvYWRsZW46IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYyg0KVxuICAgIHBheWxvYWRsZW4ud3JpdGVVSW50MzJCRSh0aGlzLnBheWxvYWQubGVuZ3RoLCAwKVxuXG4gICAgY29uc3Qgb3V0cHV0b3duZXJzbGVuOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoNClcbiAgICBvdXRwdXRvd25lcnNsZW4ud3JpdGVVSW50MzJCRSh0aGlzLm91dHB1dE93bmVycy5sZW5ndGgsIDApXG5cbiAgICBsZXQgYnNpemU6IG51bWJlciA9XG4gICAgICBzdXBlcmJ1ZmYubGVuZ3RoICtcbiAgICAgIHRoaXMuZ3JvdXBJRC5sZW5ndGggK1xuICAgICAgcGF5bG9hZGxlbi5sZW5ndGggK1xuICAgICAgdGhpcy5wYXlsb2FkLmxlbmd0aCArXG4gICAgICBvdXRwdXRvd25lcnNsZW4ubGVuZ3RoXG5cbiAgICBjb25zdCBiYXJyOiBCdWZmZXJbXSA9IFtcbiAgICAgIHN1cGVyYnVmZixcbiAgICAgIHRoaXMuZ3JvdXBJRCxcbiAgICAgIHBheWxvYWRsZW4sXG4gICAgICB0aGlzLnBheWxvYWQsXG4gICAgICBvdXRwdXRvd25lcnNsZW5cbiAgICBdXG5cbiAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgdGhpcy5vdXRwdXRPd25lcnMubGVuZ3RoOyBpKyspIHtcbiAgICAgIGxldCBiOiBCdWZmZXIgPSB0aGlzLm91dHB1dE93bmVyc1tgJHtpfWBdLnRvQnVmZmVyKClcbiAgICAgIGJhcnIucHVzaChiKVxuICAgICAgYnNpemUgKz0gYi5sZW5ndGhcbiAgICB9XG5cbiAgICByZXR1cm4gQnVmZmVyLmNvbmNhdChiYXJyLCBic2l6ZSlcbiAgfVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIGEgYmFzZS01OCBzdHJpbmcgcmVwcmVzZW50aW5nIHRoZSBbW05GVE1pbnRPcGVyYXRpb25dXS5cbiAgICovXG4gIHRvU3RyaW5nKCk6IHN0cmluZyB7XG4gICAgcmV0dXJuIGJpbnRvb2xzLmJ1ZmZlclRvQjU4KHRoaXMudG9CdWZmZXIoKSlcbiAgfVxuXG4gIC8qKlxuICAgKiBBbiBbW09wZXJhdGlvbl1dIGNsYXNzIHdoaWNoIGNvbnRhaW5zIGFuIE5GVCBvbiBhbiBhc3NldElELlxuICAgKlxuICAgKiBAcGFyYW0gZ3JvdXBJRCBUaGUgZ3JvdXAgdG8gd2hpY2ggdG8gaXNzdWUgdGhlIE5GVCBPdXRwdXRcbiAgICogQHBhcmFtIHBheWxvYWQgQSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSBvZiB0aGUgTkZUIHBheWxvYWRcbiAgICogQHBhcmFtIG91dHB1dE93bmVycyBBbiBhcnJheSBvZiBvdXRwdXRPd25lcnNcbiAgICovXG4gIGNvbnN0cnVjdG9yKFxuICAgIGdyb3VwSUQ6IG51bWJlciA9IHVuZGVmaW5lZCxcbiAgICBwYXlsb2FkOiBCdWZmZXIgPSB1bmRlZmluZWQsXG4gICAgb3V0cHV0T3duZXJzOiBPdXRwdXRPd25lcnNbXSA9IHVuZGVmaW5lZFxuICApIHtcbiAgICBzdXBlcigpXG4gICAgaWYgKFxuICAgICAgdHlwZW9mIGdyb3VwSUQgIT09IFwidW5kZWZpbmVkXCIgJiZcbiAgICAgIHR5cGVvZiBwYXlsb2FkICE9PSBcInVuZGVmaW5lZFwiICYmXG4gICAgICBvdXRwdXRPd25lcnMubGVuZ3RoXG4gICAgKSB7XG4gICAgICB0aGlzLmdyb3VwSUQud3JpdGVVSW50MzJCRShncm91cElEID8gZ3JvdXBJRCA6IDAsIDApXG4gICAgICB0aGlzLnBheWxvYWQgPSBwYXlsb2FkXG4gICAgICB0aGlzLm91dHB1dE93bmVycyA9IG91dHB1dE93bmVyc1xuICAgIH1cbiAgfVxufVxuXG4vKipcbiAqIEEgW1tPcGVyYXRpb25dXSBjbGFzcyB3aGljaCBzcGVjaWZpZXMgYSBORlQgVHJhbnNmZXIgT3AuXG4gKi9cbmV4cG9ydCBjbGFzcyBORlRUcmFuc2Zlck9wZXJhdGlvbiBleHRlbmRzIE9wZXJhdGlvbiB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIk5GVFRyYW5zZmVyT3BlcmF0aW9uXCJcbiAgcHJvdGVjdGVkIF9jb2RlY0lEID0gQVZNQ29uc3RhbnRzLkxBVEVTVENPREVDXG4gIHByb3RlY3RlZCBfdHlwZUlEID1cbiAgICB0aGlzLl9jb2RlY0lEID09PSAwXG4gICAgICA/IEFWTUNvbnN0YW50cy5ORlRYRkVST1BJRFxuICAgICAgOiBBVk1Db25zdGFudHMuTkZUWEZFUk9QSURfQ09ERUNPTkVcblxuICBzZXJpYWxpemUoZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpOiBvYmplY3Qge1xuICAgIGNvbnN0IGZpZWxkczogb2JqZWN0ID0gc3VwZXIuc2VyaWFsaXplKGVuY29kaW5nKVxuICAgIHJldHVybiB7XG4gICAgICAuLi5maWVsZHMsXG4gICAgICBvdXRwdXQ6IHRoaXMub3V0cHV0LnNlcmlhbGl6ZShlbmNvZGluZylcbiAgICB9XG4gIH1cbiAgZGVzZXJpYWxpemUoZmllbGRzOiBvYmplY3QsIGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKSB7XG4gICAgc3VwZXIuZGVzZXJpYWxpemUoZmllbGRzLCBlbmNvZGluZylcbiAgICB0aGlzLm91dHB1dCA9IG5ldyBORlRUcmFuc2Zlck91dHB1dCgpXG4gICAgdGhpcy5vdXRwdXQuZGVzZXJpYWxpemUoZmllbGRzW1wib3V0cHV0XCJdLCBlbmNvZGluZylcbiAgfVxuXG4gIHByb3RlY3RlZCBvdXRwdXQ6IE5GVFRyYW5zZmVyT3V0cHV0XG5cbiAgLyoqXG4gICAqIFNldCB0aGUgY29kZWNJRFxuICAgKlxuICAgKiBAcGFyYW0gY29kZWNJRCBUaGUgY29kZWNJRCB0byBzZXRcbiAgICovXG4gIHNldENvZGVjSUQoY29kZWNJRDogbnVtYmVyKTogdm9pZCB7XG4gICAgaWYgKGNvZGVjSUQgIT09IDAgJiYgY29kZWNJRCAhPT0gMSkge1xuICAgICAgLyogaXN0YW5idWwgaWdub3JlIG5leHQgKi9cbiAgICAgIHRocm93IG5ldyBDb2RlY0lkRXJyb3IoXG4gICAgICAgIFwiRXJyb3IgLSBORlRUcmFuc2Zlck9wZXJhdGlvbi5zZXRDb2RlY0lEOiBpbnZhbGlkIGNvZGVjSUQuIFZhbGlkIGNvZGVjSURzIGFyZSAwIGFuZCAxLlwiXG4gICAgICApXG4gICAgfVxuICAgIHRoaXMuX2NvZGVjSUQgPSBjb2RlY0lEXG4gICAgdGhpcy5fdHlwZUlEID1cbiAgICAgIHRoaXMuX2NvZGVjSUQgPT09IDBcbiAgICAgICAgPyBBVk1Db25zdGFudHMuTkZUWEZFUk9QSURcbiAgICAgICAgOiBBVk1Db25zdGFudHMuTkZUWEZFUk9QSURfQ09ERUNPTkVcbiAgfVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBvcGVyYXRpb24gSUQuXG4gICAqL1xuICBnZXRPcGVyYXRpb25JRCgpOiBudW1iZXIge1xuICAgIHJldHVybiB0aGlzLl90eXBlSURcbiAgfVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBjcmVkZW50aWFsIElELlxuICAgKi9cbiAgZ2V0Q3JlZGVudGlhbElEKCk6IG51bWJlciB7XG4gICAgaWYgKHRoaXMuX2NvZGVjSUQgPT09IDApIHtcbiAgICAgIHJldHVybiBBVk1Db25zdGFudHMuTkZUQ1JFREVOVElBTFxuICAgIH0gZWxzZSBpZiAodGhpcy5fY29kZWNJRCA9PT0gMSkge1xuICAgICAgcmV0dXJuIEFWTUNvbnN0YW50cy5ORlRDUkVERU5USUFMX0NPREVDT05FXG4gICAgfVxuICB9XG5cbiAgZ2V0T3V0cHV0ID0gKCk6IE5GVFRyYW5zZmVyT3V0cHV0ID0+IHRoaXMub3V0cHV0XG5cbiAgLyoqXG4gICAqIFBvcHVhdGVzIHRoZSBpbnN0YW5jZSBmcm9tIGEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gcmVwcmVzZW50aW5nIHRoZSBbW05GVFRyYW5zZmVyT3BlcmF0aW9uXV0gYW5kIHJldHVybnMgdGhlIHVwZGF0ZWQgb2Zmc2V0LlxuICAgKi9cbiAgZnJvbUJ1ZmZlcihieXRlczogQnVmZmVyLCBvZmZzZXQ6IG51bWJlciA9IDApOiBudW1iZXIge1xuICAgIG9mZnNldCA9IHN1cGVyLmZyb21CdWZmZXIoYnl0ZXMsIG9mZnNldClcbiAgICB0aGlzLm91dHB1dCA9IG5ldyBORlRUcmFuc2Zlck91dHB1dCgpXG4gICAgcmV0dXJuIHRoaXMub3V0cHV0LmZyb21CdWZmZXIoYnl0ZXMsIG9mZnNldClcbiAgfVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBidWZmZXIgcmVwcmVzZW50aW5nIHRoZSBbW05GVFRyYW5zZmVyT3BlcmF0aW9uXV0gaW5zdGFuY2UuXG4gICAqL1xuICB0b0J1ZmZlcigpOiBCdWZmZXIge1xuICAgIGNvbnN0IHN1cGVyYnVmZjogQnVmZmVyID0gc3VwZXIudG9CdWZmZXIoKVxuICAgIGNvbnN0IG91dGJ1ZmY6IEJ1ZmZlciA9IHRoaXMub3V0cHV0LnRvQnVmZmVyKClcbiAgICBjb25zdCBic2l6ZTogbnVtYmVyID0gc3VwZXJidWZmLmxlbmd0aCArIG91dGJ1ZmYubGVuZ3RoXG4gICAgY29uc3QgYmFycjogQnVmZmVyW10gPSBbc3VwZXJidWZmLCBvdXRidWZmXVxuICAgIHJldHVybiBCdWZmZXIuY29uY2F0KGJhcnIsIGJzaXplKVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSBiYXNlLTU4IHN0cmluZyByZXByZXNlbnRpbmcgdGhlIFtbTkZUVHJhbnNmZXJPcGVyYXRpb25dXS5cbiAgICovXG4gIHRvU3RyaW5nKCk6IHN0cmluZyB7XG4gICAgcmV0dXJuIGJpbnRvb2xzLmJ1ZmZlclRvQjU4KHRoaXMudG9CdWZmZXIoKSlcbiAgfVxuXG4gIC8qKlxuICAgKiBBbiBbW09wZXJhdGlvbl1dIGNsYXNzIHdoaWNoIGNvbnRhaW5zIGFuIE5GVCBvbiBhbiBhc3NldElELlxuICAgKlxuICAgKiBAcGFyYW0gb3V0cHV0IEFuIFtbTkZUVHJhbnNmZXJPdXRwdXRdXVxuICAgKi9cbiAgY29uc3RydWN0b3Iob3V0cHV0OiBORlRUcmFuc2Zlck91dHB1dCA9IHVuZGVmaW5lZCkge1xuICAgIHN1cGVyKClcbiAgICBpZiAodHlwZW9mIG91dHB1dCAhPT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgdGhpcy5vdXRwdXQgPSBvdXRwdXRcbiAgICB9XG4gIH1cbn1cblxuLyoqXG4gKiBDbGFzcyBmb3IgcmVwcmVzZW50aW5nIGEgVVRYT0lEIHVzZWQgaW4gW1tUcmFuc2ZlcmFibGVPcF1dIHR5cGVzXG4gKi9cbmV4cG9ydCBjbGFzcyBVVFhPSUQgZXh0ZW5kcyBOQnl0ZXMge1xuICBwcm90ZWN0ZWQgX3R5cGVOYW1lID0gXCJVVFhPSURcIlxuICBwcm90ZWN0ZWQgX3R5cGVJRCA9IHVuZGVmaW5lZFxuXG4gIC8vc2VyaWFsaXplIGFuZCBkZXNlcmlhbGl6ZSBib3RoIGFyZSBpbmhlcml0ZWRcblxuICBwcm90ZWN0ZWQgYnl0ZXMgPSBCdWZmZXIuYWxsb2MoMzYpXG4gIHByb3RlY3RlZCBic2l6ZSA9IDM2XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSBmdW5jdGlvbiB1c2VkIHRvIHNvcnQgYW4gYXJyYXkgb2YgW1tVVFhPSURdXXNcbiAgICovXG4gIHN0YXRpYyBjb21wYXJhdG9yID1cbiAgICAoKTogKChhOiBVVFhPSUQsIGI6IFVUWE9JRCkgPT4gMSB8IC0xIHwgMCkgPT5cbiAgICAoYTogVVRYT0lELCBiOiBVVFhPSUQpOiAxIHwgLTEgfCAwID0+XG4gICAgICBCdWZmZXIuY29tcGFyZShhLnRvQnVmZmVyKCksIGIudG9CdWZmZXIoKSkgYXMgMSB8IC0xIHwgMFxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIGEgYmFzZS01OCByZXByZXNlbnRhdGlvbiBvZiB0aGUgW1tVVFhPSURdXS5cbiAgICovXG4gIHRvU3RyaW5nKCk6IHN0cmluZyB7XG4gICAgcmV0dXJuIGJpbnRvb2xzLmNiNThFbmNvZGUodGhpcy50b0J1ZmZlcigpKVxuICB9XG5cbiAgLyoqXG4gICAqIFRha2VzIGEgYmFzZS01OCBzdHJpbmcgY29udGFpbmluZyBhbiBbW1VUWE9JRF1dLCBwYXJzZXMgaXQsIHBvcHVsYXRlcyB0aGUgY2xhc3MsIGFuZCByZXR1cm5zIHRoZSBsZW5ndGggb2YgdGhlIFVUWE9JRCBpbiBieXRlcy5cbiAgICpcbiAgICogQHBhcmFtIGJ5dGVzIEEgYmFzZS01OCBzdHJpbmcgY29udGFpbmluZyBhIHJhdyBbW1VUWE9JRF1dXG4gICAqXG4gICAqIEByZXR1cm5zIFRoZSBsZW5ndGggb2YgdGhlIHJhdyBbW1VUWE9JRF1dXG4gICAqL1xuICBmcm9tU3RyaW5nKHV0eG9pZDogc3RyaW5nKTogbnVtYmVyIHtcbiAgICBjb25zdCB1dHhvaWRidWZmOiBCdWZmZXIgPSBiaW50b29scy5iNThUb0J1ZmZlcih1dHhvaWQpXG4gICAgaWYgKHV0eG9pZGJ1ZmYubGVuZ3RoID09PSA0MCAmJiBiaW50b29scy52YWxpZGF0ZUNoZWNrc3VtKHV0eG9pZGJ1ZmYpKSB7XG4gICAgICBjb25zdCBuZXdidWZmOiBCdWZmZXIgPSBiaW50b29scy5jb3B5RnJvbShcbiAgICAgICAgdXR4b2lkYnVmZixcbiAgICAgICAgMCxcbiAgICAgICAgdXR4b2lkYnVmZi5sZW5ndGggLSA0XG4gICAgICApXG4gICAgICBpZiAobmV3YnVmZi5sZW5ndGggPT09IDM2KSB7XG4gICAgICAgIHRoaXMuYnl0ZXMgPSBuZXdidWZmXG4gICAgICB9XG4gICAgfSBlbHNlIGlmICh1dHhvaWRidWZmLmxlbmd0aCA9PT0gNDApIHtcbiAgICAgIHRocm93IG5ldyBDaGVja3N1bUVycm9yKFxuICAgICAgICBcIkVycm9yIC0gVVRYT0lELmZyb21TdHJpbmc6IGludmFsaWQgY2hlY2tzdW0gb24gYWRkcmVzc1wiXG4gICAgICApXG4gICAgfSBlbHNlIGlmICh1dHhvaWRidWZmLmxlbmd0aCA9PT0gMzYpIHtcbiAgICAgIHRoaXMuYnl0ZXMgPSB1dHhvaWRidWZmXG4gICAgfSBlbHNlIHtcbiAgICAgIC8qIGlzdGFuYnVsIGlnbm9yZSBuZXh0ICovXG4gICAgICB0aHJvdyBuZXcgQWRkcmVzc0Vycm9yKFwiRXJyb3IgLSBVVFhPSUQuZnJvbVN0cmluZzogaW52YWxpZCBhZGRyZXNzXCIpXG4gICAgfVxuICAgIHJldHVybiB0aGlzLmdldFNpemUoKVxuICB9XG5cbiAgY2xvbmUoKTogdGhpcyB7XG4gICAgY29uc3QgbmV3YmFzZTogVVRYT0lEID0gbmV3IFVUWE9JRCgpXG4gICAgbmV3YmFzZS5mcm9tQnVmZmVyKHRoaXMudG9CdWZmZXIoKSlcbiAgICByZXR1cm4gbmV3YmFzZSBhcyB0aGlzXG4gIH1cblxuICBjcmVhdGUoLi4uYXJnczogYW55W10pOiB0aGlzIHtcbiAgICByZXR1cm4gbmV3IFVUWE9JRCgpIGFzIHRoaXNcbiAgfVxuXG4gIC8qKlxuICAgKiBDbGFzcyBmb3IgcmVwcmVzZW50aW5nIGEgVVRYT0lEIHVzZWQgaW4gW1tUcmFuc2ZlcmFibGVPcF1dIHR5cGVzXG4gICAqL1xuICBjb25zdHJ1Y3RvcigpIHtcbiAgICBzdXBlcigpXG4gIH1cbn1cbiJdfQ==