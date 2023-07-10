"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.StandardUTXOSet = exports.StandardUTXO = void 0;
/**
 * @packageDocumentation
 * @module Common-UTXOs
 */
const buffer_1 = require("buffer/");
const bintools_1 = __importDefault(require("../utils/bintools"));
const bn_js_1 = __importDefault(require("bn.js"));
const output_1 = require("./output");
const helperfunctions_1 = require("../utils/helperfunctions");
const serialization_1 = require("../utils/serialization");
const errors_1 = require("../utils/errors");
/**
 * @ignore
 */
const bintools = bintools_1.default.getInstance();
const serialization = serialization_1.Serialization.getInstance();
/**
 * Class for representing a single StandardUTXO.
 */
class StandardUTXO extends serialization_1.Serializable {
    /**
     * Class for representing a single StandardUTXO.
     *
     * @param codecID Optional number which specifies the codeID of the UTXO. Default 0
     * @param txID Optional {@link https://github.com/feross/buffer|Buffer} of transaction ID for the StandardUTXO
     * @param txidx Optional {@link https://github.com/feross/buffer|Buffer} or number for the index of the transaction's [[Output]]
     * @param assetID Optional {@link https://github.com/feross/buffer|Buffer} of the asset ID for the StandardUTXO
     * @param outputid Optional {@link https://github.com/feross/buffer|Buffer} or number of the output ID for the StandardUTXO
     */
    constructor(codecID = 0, txID = undefined, outputidx = undefined, assetID = undefined, output = undefined) {
        super();
        this._typeName = "StandardUTXO";
        this._typeID = undefined;
        this.codecID = buffer_1.Buffer.alloc(2);
        this.txid = buffer_1.Buffer.alloc(32);
        this.outputidx = buffer_1.Buffer.alloc(4);
        this.assetID = buffer_1.Buffer.alloc(32);
        this.output = undefined;
        /**
         * Returns the numeric representation of the CodecID.
         */
        this.getCodecID = () => this.codecID.readUInt8(0);
        /**
         * Returns the {@link https://github.com/feross/buffer|Buffer} representation of the CodecID
         */
        this.getCodecIDBuffer = () => this.codecID;
        /**
         * Returns a {@link https://github.com/feross/buffer|Buffer} of the TxID.
         */
        this.getTxID = () => this.txid;
        /**
         * Returns a {@link https://github.com/feross/buffer|Buffer}  of the OutputIdx.
         */
        this.getOutputIdx = () => this.outputidx;
        /**
         * Returns the assetID as a {@link https://github.com/feross/buffer|Buffer}.
         */
        this.getAssetID = () => this.assetID;
        /**
         * Returns the UTXOID as a base-58 string (UTXOID is a string )
         */
        this.getUTXOID = () => bintools.bufferToB58(buffer_1.Buffer.concat([this.getTxID(), this.getOutputIdx()]));
        /**
         * Returns a reference to the output
         */
        this.getOutput = () => this.output;
        if (typeof codecID !== "undefined") {
            this.codecID.writeUInt8(codecID, 0);
        }
        if (typeof txID !== "undefined") {
            this.txid = txID;
        }
        if (typeof outputidx === "number") {
            this.outputidx.writeUInt32BE(outputidx, 0);
        }
        else if (outputidx instanceof buffer_1.Buffer) {
            this.outputidx = outputidx;
        }
        if (typeof assetID !== "undefined") {
            this.assetID = assetID;
        }
        if (typeof output !== "undefined") {
            this.output = output;
        }
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { codecID: serialization.encoder(this.codecID, encoding, "Buffer", "decimalString"), txid: serialization.encoder(this.txid, encoding, "Buffer", "cb58"), outputidx: serialization.encoder(this.outputidx, encoding, "Buffer", "decimalString"), assetID: serialization.encoder(this.assetID, encoding, "Buffer", "cb58"), output: this.output.serialize(encoding) });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.codecID = serialization.decoder(fields["codecID"], encoding, "decimalString", "Buffer", 2);
        this.txid = serialization.decoder(fields["txid"], encoding, "cb58", "Buffer", 32);
        this.outputidx = serialization.decoder(fields["outputidx"], encoding, "decimalString", "Buffer", 4);
        this.assetID = serialization.decoder(fields["assetID"], encoding, "cb58", "Buffer", 32);
    }
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[StandardUTXO]].
     */
    toBuffer() {
        const outbuff = this.output.toBuffer();
        const outputidbuffer = buffer_1.Buffer.alloc(4);
        outputidbuffer.writeUInt32BE(this.output.getOutputID(), 0);
        const barr = [
            this.codecID,
            this.txid,
            this.outputidx,
            this.assetID,
            outputidbuffer,
            outbuff
        ];
        return buffer_1.Buffer.concat(barr, this.codecID.length +
            this.txid.length +
            this.outputidx.length +
            this.assetID.length +
            outputidbuffer.length +
            outbuff.length);
    }
}
exports.StandardUTXO = StandardUTXO;
/**
 * Class representing a set of [[StandardUTXO]]s.
 */
class StandardUTXOSet extends serialization_1.Serializable {
    constructor() {
        super(...arguments);
        this._typeName = "StandardUTXOSet";
        this._typeID = undefined;
        this.utxos = {};
        this.addressUTXOs = {}; // maps address to utxoids:locktime
        /**
         * Returns true if the [[StandardUTXO]] is in the StandardUTXOSet.
         *
         * @param utxo Either a [[StandardUTXO]] a cb58 serialized string representing a StandardUTXO
         */
        this.includes = (utxo) => {
            let utxoX = undefined;
            let utxoid = undefined;
            try {
                utxoX = this.parseUTXO(utxo);
                utxoid = utxoX.getUTXOID();
            }
            catch (e) {
                if (e instanceof Error) {
                    console.log(e.message);
                }
                else {
                    console.log(e);
                }
                return false;
            }
            return utxoid in this.utxos;
        };
        /**
         * Removes a [[StandardUTXO]] from the [[StandardUTXOSet]] if it exists.
         *
         * @param utxo Either a [[StandardUTXO]] an cb58 serialized string representing a StandardUTXO
         *
         * @returns A [[StandardUTXO]] if it was removed and undefined if nothing was removed.
         */
        this.remove = (utxo) => {
            let utxovar = undefined;
            try {
                utxovar = this.parseUTXO(utxo);
            }
            catch (e) {
                if (e instanceof Error) {
                    console.log(e.message);
                }
                else {
                    console.log(e);
                }
                return undefined;
            }
            const utxoid = utxovar.getUTXOID();
            if (!(utxoid in this.utxos)) {
                return undefined;
            }
            delete this.utxos[`${utxoid}`];
            const addresses = Object.keys(this.addressUTXOs);
            for (let i = 0; i < addresses.length; i++) {
                if (utxoid in this.addressUTXOs[addresses[`${i}`]]) {
                    delete this.addressUTXOs[addresses[`${i}`]][`${utxoid}`];
                }
            }
            return utxovar;
        };
        /**
         * Removes an array of [[StandardUTXO]]s to the [[StandardUTXOSet]].
         *
         * @param utxo Either a [[StandardUTXO]] an cb58 serialized string representing a StandardUTXO
         * @param overwrite If true, if the UTXOID already exists, overwrite it... default false
         *
         * @returns An array of UTXOs which were removed.
         */
        this.removeArray = (utxos) => {
            const removed = [];
            for (let i = 0; i < utxos.length; i++) {
                const result = this.remove(utxos[`${i}`]);
                if (typeof result !== "undefined") {
                    removed.push(result);
                }
            }
            return removed;
        };
        /**
         * Gets a [[StandardUTXO]] from the [[StandardUTXOSet]] by its UTXOID.
         *
         * @param utxoid String representing the UTXOID
         *
         * @returns A [[StandardUTXO]] if it exists in the set.
         */
        this.getUTXO = (utxoid) => this.utxos[`${utxoid}`];
        /**
         * Gets all the [[StandardUTXO]]s, optionally that match with UTXOIDs in an array
         *
         * @param utxoids An optional array of UTXOIDs, returns all [[StandardUTXO]]s if not provided
         *
         * @returns An array of [[StandardUTXO]]s.
         */
        this.getAllUTXOs = (utxoids = undefined) => {
            let results = [];
            if (typeof utxoids !== "undefined" && Array.isArray(utxoids)) {
                results = utxoids
                    .filter((utxoid) => this.utxos[`${utxoid}`])
                    .map((utxoid) => this.utxos[`${utxoid}`]);
            }
            else {
                results = Object.values(this.utxos);
            }
            return results;
        };
        /**
         * Gets all the [[StandardUTXO]]s as strings, optionally that match with UTXOIDs in an array.
         *
         * @param utxoids An optional array of UTXOIDs, returns all [[StandardUTXO]]s if not provided
         *
         * @returns An array of [[StandardUTXO]]s as cb58 serialized strings.
         */
        this.getAllUTXOStrings = (utxoids = undefined) => {
            const results = [];
            const utxos = Object.keys(this.utxos);
            if (typeof utxoids !== "undefined" && Array.isArray(utxoids)) {
                for (let i = 0; i < utxoids.length; i++) {
                    if (utxoids[`${i}`] in this.utxos) {
                        results.push(this.utxos[utxoids[`${i}`]].toString());
                    }
                }
            }
            else {
                for (const u of utxos) {
                    results.push(this.utxos[`${u}`].toString());
                }
            }
            return results;
        };
        /**
         * Given an address or array of addresses, returns all the UTXOIDs for those addresses
         *
         * @param address An array of address {@link https://github.com/feross/buffer|Buffer}s
         * @param spendable If true, only retrieves UTXOIDs whose locktime has passed
         *
         * @returns An array of addresses.
         */
        this.getUTXOIDs = (addresses = undefined, spendable = true) => {
            if (typeof addresses !== "undefined") {
                const results = [];
                const now = (0, helperfunctions_1.UnixNow)();
                for (let i = 0; i < addresses.length; i++) {
                    if (addresses[`${i}`].toString("hex") in this.addressUTXOs) {
                        const entries = Object.entries(this.addressUTXOs[addresses[`${i}`].toString("hex")]);
                        for (const [utxoid, locktime] of entries) {
                            if ((results.indexOf(utxoid) === -1 &&
                                spendable &&
                                locktime.lte(now)) ||
                                !spendable) {
                                results.push(utxoid);
                            }
                        }
                    }
                }
                return results;
            }
            return Object.keys(this.utxos);
        };
        /**
         * Gets the addresses in the [[StandardUTXOSet]] and returns an array of {@link https://github.com/feross/buffer|Buffer}.
         */
        this.getAddresses = () => Object.keys(this.addressUTXOs).map((k) => buffer_1.Buffer.from(k, "hex"));
        /**
         * Returns the balance of a set of addresses in the StandardUTXOSet.
         *
         * @param addresses An array of addresses
         * @param assetID Either a {@link https://github.com/feross/buffer|Buffer} or an cb58 serialized representation of an AssetID
         * @param asOf The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
         *
         * @returns Returns the total balance as a {@link https://github.com/indutny/bn.js/|BN}.
         */
        this.getBalance = (addresses, assetID, asOf = undefined) => {
            const utxoids = this.getUTXOIDs(addresses);
            const utxos = this.getAllUTXOs(utxoids);
            let spend = new bn_js_1.default(0);
            let asset;
            if (typeof assetID === "string") {
                asset = bintools.cb58Decode(assetID);
            }
            else {
                asset = assetID;
            }
            for (let i = 0; i < utxos.length; i++) {
                if (utxos[`${i}`].getOutput() instanceof output_1.StandardAmountOutput &&
                    utxos[`${i}`].getAssetID().toString("hex") === asset.toString("hex") &&
                    utxos[`${i}`].getOutput().meetsThreshold(addresses, asOf)) {
                    spend = spend.add(utxos[`${i}`].getOutput().getAmount());
                }
            }
            return spend;
        };
        /**
         * Gets all the Asset IDs, optionally that match with Asset IDs in an array
         *
         * @param utxoids An optional array of Addresses as string or Buffer, returns all Asset IDs if not provided
         *
         * @returns An array of {@link https://github.com/feross/buffer|Buffer} representing the Asset IDs.
         */
        this.getAssetIDs = (addresses = undefined) => {
            const results = new Set();
            let utxoids = [];
            if (typeof addresses !== "undefined") {
                utxoids = this.getUTXOIDs(addresses);
            }
            else {
                utxoids = this.getUTXOIDs();
            }
            for (let i = 0; i < utxoids.length; i++) {
                if (utxoids[`${i}`] in this.utxos && !(utxoids[`${i}`] in results)) {
                    results.add(this.utxos[utxoids[`${i}`]].getAssetID());
                }
            }
            return [...results];
        };
        /**
         * Returns a new set with copy of UTXOs in this and set parameter.
         *
         * @param utxoset The [[StandardUTXOSet]] to merge with this one
         * @param hasUTXOIDs Will subselect a set of [[StandardUTXO]]s which have the UTXOIDs provided in this array, defults to all UTXOs
         *
         * @returns A new StandardUTXOSet that contains all the filtered elements.
         */
        this.merge = (utxoset, hasUTXOIDs = undefined) => {
            const results = this.create();
            const utxos1 = this.getAllUTXOs(hasUTXOIDs);
            const utxos2 = utxoset.getAllUTXOs(hasUTXOIDs);
            const process = (utxo) => {
                results.add(utxo);
            };
            utxos1.forEach(process);
            utxos2.forEach(process);
            return results;
        };
        /**
         * Set intersetion between this set and a parameter.
         *
         * @param utxoset The set to intersect
         *
         * @returns A new StandardUTXOSet containing the intersection
         */
        this.intersection = (utxoset) => {
            const us1 = this.getUTXOIDs();
            const us2 = utxoset.getUTXOIDs();
            const results = us1.filter((utxoid) => us2.includes(utxoid));
            return this.merge(utxoset, results);
        };
        /**
         * Set difference between this set and a parameter.
         *
         * @param utxoset The set to difference
         *
         * @returns A new StandardUTXOSet containing the difference
         */
        this.difference = (utxoset) => {
            const us1 = this.getUTXOIDs();
            const us2 = utxoset.getUTXOIDs();
            const results = us1.filter((utxoid) => !us2.includes(utxoid));
            return this.merge(utxoset, results);
        };
        /**
         * Set symmetrical difference between this set and a parameter.
         *
         * @param utxoset The set to symmetrical difference
         *
         * @returns A new StandardUTXOSet containing the symmetrical difference
         */
        this.symDifference = (utxoset) => {
            const us1 = this.getUTXOIDs();
            const us2 = utxoset.getUTXOIDs();
            const results = us1
                .filter((utxoid) => !us2.includes(utxoid))
                .concat(us2.filter((utxoid) => !us1.includes(utxoid)));
            return this.merge(utxoset, results);
        };
        /**
         * Set union between this set and a parameter.
         *
         * @param utxoset The set to union
         *
         * @returns A new StandardUTXOSet containing the union
         */
        this.union = (utxoset) => this.merge(utxoset);
        /**
         * Merges a set by the rule provided.
         *
         * @param utxoset The set to merge by the MergeRule
         * @param mergeRule The [[MergeRule]] to apply
         *
         * @returns A new StandardUTXOSet containing the merged data
         *
         * @remarks
         * The merge rules are as follows:
         *   * "intersection" - the intersection of the set
         *   * "differenceSelf" - the difference between the existing data and new set
         *   * "differenceNew" - the difference between the new data and the existing set
         *   * "symDifference" - the union of the differences between both sets of data
         *   * "union" - the unique set of all elements contained in both sets
         *   * "unionMinusNew" - the unique set of all elements contained in both sets, excluding values only found in the new set
         *   * "unionMinusSelf" - the unique set of all elements contained in both sets, excluding values only found in the existing set
         */
        this.mergeByRule = (utxoset, mergeRule) => {
            let uSet;
            switch (mergeRule) {
                case "intersection":
                    return this.intersection(utxoset);
                case "differenceSelf":
                    return this.difference(utxoset);
                case "differenceNew":
                    return utxoset.difference(this);
                case "symDifference":
                    return this.symDifference(utxoset);
                case "union":
                    return this.union(utxoset);
                case "unionMinusNew":
                    uSet = this.union(utxoset);
                    return uSet.difference(utxoset);
                case "unionMinusSelf":
                    uSet = this.union(utxoset);
                    return uSet.difference(this);
                default:
                    throw new errors_1.MergeRuleError("Error - StandardUTXOSet.mergeByRule: bad MergeRule");
            }
        };
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        let utxos = {};
        for (let utxoid in this.utxos) {
            let utxoidCleaned = serialization.encoder(utxoid, encoding, "base58", "base58");
            utxos[`${utxoidCleaned}`] = this.utxos[`${utxoid}`].serialize(encoding);
        }
        let addressUTXOs = {};
        for (let address in this.addressUTXOs) {
            let addressCleaned = serialization.encoder(address, encoding, "hex", "cb58");
            let utxobalance = {};
            for (let utxoid in this.addressUTXOs[`${address}`]) {
                let utxoidCleaned = serialization.encoder(utxoid, encoding, "base58", "base58");
                utxobalance[`${utxoidCleaned}`] = serialization.encoder(this.addressUTXOs[`${address}`][`${utxoid}`], encoding, "BN", "decimalString");
            }
            addressUTXOs[`${addressCleaned}`] = utxobalance;
        }
        return Object.assign(Object.assign({}, fields), { utxos,
            addressUTXOs });
    }
    /**
     * Adds a [[StandardUTXO]] to the StandardUTXOSet.
     *
     * @param utxo Either a [[StandardUTXO]] an cb58 serialized string representing a StandardUTXO
     * @param overwrite If true, if the UTXOID already exists, overwrite it... default false
     *
     * @returns A [[StandardUTXO]] if one was added and undefined if nothing was added.
     */
    add(utxo, overwrite = false) {
        let utxovar = undefined;
        try {
            utxovar = this.parseUTXO(utxo);
        }
        catch (e) {
            if (e instanceof Error) {
                console.log(e.message);
            }
            else {
                console.log(e);
            }
            return undefined;
        }
        const utxoid = utxovar.getUTXOID();
        if (!(utxoid in this.utxos) || overwrite === true) {
            this.utxos[`${utxoid}`] = utxovar;
            const addresses = utxovar.getOutput().getAddresses();
            const locktime = utxovar.getOutput().getLocktime();
            for (let i = 0; i < addresses.length; i++) {
                const address = addresses[`${i}`].toString("hex");
                if (!(address in this.addressUTXOs)) {
                    this.addressUTXOs[`${address}`] = {};
                }
                this.addressUTXOs[`${address}`][`${utxoid}`] = locktime;
            }
            return utxovar;
        }
        return undefined;
    }
    /**
     * Adds an array of [[StandardUTXO]]s to the [[StandardUTXOSet]].
     *
     * @param utxo Either a [[StandardUTXO]] an cb58 serialized string representing a StandardUTXO
     * @param overwrite If true, if the UTXOID already exists, overwrite it... default false
     *
     * @returns An array of StandardUTXOs which were added.
     */
    addArray(utxos, overwrite = false) {
        const added = [];
        for (let i = 0; i < utxos.length; i++) {
            let result = this.add(utxos[`${i}`], overwrite);
            if (typeof result !== "undefined") {
                added.push(result);
            }
        }
        return added;
    }
    filter(args, lambda) {
        let newset = this.clone();
        let utxos = this.getAllUTXOs();
        for (let i = 0; i < utxos.length; i++) {
            if (lambda(utxos[`${i}`], ...args) === false) {
                newset.remove(utxos[`${i}`]);
            }
        }
        return newset;
    }
}
exports.StandardUTXOSet = StandardUTXOSet;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidXR4b3MuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi9zcmMvY29tbW9uL3V0eG9zLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7OztBQUFBOzs7R0FHRztBQUNILG9DQUFnQztBQUNoQyxpRUFBd0M7QUFDeEMsa0RBQXNCO0FBQ3RCLHFDQUF1RDtBQUN2RCw4REFBa0Q7QUFFbEQsMERBSStCO0FBQy9CLDRDQUFnRDtBQUVoRDs7R0FFRztBQUNILE1BQU0sUUFBUSxHQUFhLGtCQUFRLENBQUMsV0FBVyxFQUFFLENBQUE7QUFDakQsTUFBTSxhQUFhLEdBQWtCLDZCQUFhLENBQUMsV0FBVyxFQUFFLENBQUE7QUFFaEU7O0dBRUc7QUFDSCxNQUFzQixZQUFhLFNBQVEsNEJBQVk7SUFtSnJEOzs7Ozs7OztPQVFHO0lBQ0gsWUFDRSxVQUFrQixDQUFDLEVBQ25CLE9BQWUsU0FBUyxFQUN4QixZQUE2QixTQUFTLEVBQ3RDLFVBQWtCLFNBQVMsRUFDM0IsU0FBaUIsU0FBUztRQUUxQixLQUFLLEVBQUUsQ0FBQTtRQWxLQyxjQUFTLEdBQUcsY0FBYyxDQUFBO1FBQzFCLFlBQU8sR0FBRyxTQUFTLENBQUE7UUF1RG5CLFlBQU8sR0FBVyxlQUFNLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxDQUFBO1FBQ2pDLFNBQUksR0FBVyxlQUFNLENBQUMsS0FBSyxDQUFDLEVBQUUsQ0FBQyxDQUFBO1FBQy9CLGNBQVMsR0FBVyxlQUFNLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxDQUFBO1FBQ25DLFlBQU8sR0FBVyxlQUFNLENBQUMsS0FBSyxDQUFDLEVBQUUsQ0FBQyxDQUFBO1FBQ2xDLFdBQU0sR0FBVyxTQUFTLENBQUE7UUFFcEM7O1dBRUc7UUFDSCxlQUFVLEdBQUcsR0FBc0MsRUFBRSxDQUNuRCxJQUFJLENBQUMsT0FBTyxDQUFDLFNBQVMsQ0FBQyxDQUFDLENBQUMsQ0FBQTtRQUUzQjs7V0FFRztRQUNILHFCQUFnQixHQUFHLEdBQVcs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LENBQUMsS0FBSyxDQUFDLEdBQUcsQ0FBQyxFQUFFLENBQUMsQ0FBQyxDQUFBO2FBQzdCO1NBQ0Y7UUFDRCxPQUFPLE1BQU0sQ0FBQTtJQUNmLENBQUM7Q0FzSEY7QUFqZUQsMENBaWVDIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgQ29tbW9uLVVUWE9zXG4gKi9cbmltcG9ydCB7IEJ1ZmZlciB9IGZyb20gXCJidWZmZXIvXCJcbmltcG9ydCBCaW5Ub29scyBmcm9tIFwiLi4vdXRpbHMvYmludG9vbHNcIlxuaW1wb3J0IEJOIGZyb20gXCJibi5qc1wiXG5pbXBvcnQgeyBPdXRwdXQsIFN0YW5kYXJkQW1vdW50T3V0cHV0IH0gZnJvbSBcIi4vb3V0cHV0XCJcbmltcG9ydCB7IFVuaXhOb3cgfSBmcm9tIFwiLi4vdXRpbHMvaGVscGVyZnVuY3Rpb25zXCJcbmltcG9ydCB7IE1lcmdlUnVsZSB9IGZyb20gXCIuLi91dGlscy9jb25zdGFudHNcIlxuaW1wb3J0IHtcbiAgU2VyaWFsaXphYmxlLFxuICBTZXJpYWxpemF0aW9uLFxuICBTZXJpYWxpemVkRW5jb2Rpbmdcbn0gZnJvbSBcIi4uL3V0aWxzL3NlcmlhbGl6YXRpb25cIlxuaW1wb3J0IHsgTWVyZ2VSdWxlRXJyb3IgfSBmcm9tIFwiLi4vdXRpbHMvZXJyb3JzXCJcblxuLyoqXG4gKiBAaWdub3JlXG4gKi9cbmNvbnN0IGJpbnRvb2xzOiBCaW5Ub29scyA9IEJpblRvb2xzLmdldEluc3RhbmNlKClcbmNvbnN0IHNlcmlhbGl6YXRpb246IFNlcmlhbGl6YXRpb24gPSBTZXJpYWxpemF0aW9uLmdldEluc3RhbmNlKClcblxuLyoqXG4gKiBDbGFzcyBmb3IgcmVwcmVzZW50aW5nIGEgc2luZ2xlIFN0YW5kYXJkVVRYTy5cbiAqL1xuZXhwb3J0IGFic3RyYWN0IGNsYXNzIFN0YW5kYXJkVVRYTyBleHRlbmRzIFNlcmlhbGl6YWJsZSB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIlN0YW5kYXJkVVRYT1wiXG4gIHByb3RlY3RlZCBfdHlwZUlEID0gdW5kZWZpbmVkXG5cbiAgc2VyaWFsaXplKGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKTogb2JqZWN0IHtcbiAgICBsZXQgZmllbGRzOiBvYmplY3QgPSBzdXBlci5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgcmV0dXJuIHtcbiAgICAgIC4uLmZpZWxkcyxcbiAgICAgIGNvZGVjSUQ6IHNlcmlhbGl6YXRpb24uZW5jb2RlcihcbiAgICAgICAgdGhpcy5jb2RlY0lELFxuICAgICAgICBlbmNvZGluZyxcbiAgICAgICAgXCJCdWZmZXJcIixcbiAgICAgICAgXCJkZWNpbWFsU3RyaW5nXCJcbiAgICAgICksXG4gICAgICB0eGlkOiBzZXJpYWxpemF0aW9uLmVuY29kZXIodGhpcy50eGlkLCBlbmNvZGluZywgXCJCdWZmZXJcIiwgXCJjYjU4XCIpLFxuICAgICAgb3V0cHV0aWR4OiBzZXJpYWxpemF0aW9uLmVuY29kZXIoXG4gICAgICAgIHRoaXMub3V0cHV0aWR4LFxuICAgICAgICBlbmNvZGluZyxcbiAgICAgICAgXCJCdWZmZXJcIixcbiAgICAgICAgXCJkZWNpbWFsU3RyaW5nXCJcbiAgICAgICksXG4gICAgICBhc3NldElEOiBzZXJpYWxpemF0aW9uLmVuY29kZXIodGhpcy5hc3NldElELCBlbmNvZGluZywgXCJCdWZmZXJcIiwgXCJjYjU4XCIpLFxuICAgICAgb3V0cHV0OiB0aGlzLm91dHB1dC5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgfVxuICB9XG4gIGRlc2VyaWFsaXplKGZpZWxkczogb2JqZWN0LCBlbmNvZGluZzogU2VyaWFsaXplZEVuY29kaW5nID0gXCJoZXhcIikge1xuICAgIHN1cGVyLmRlc2VyaWFsaXplKGZpZWxkcywgZW5jb2RpbmcpXG4gICAgdGhpcy5jb2RlY0lEID0gc2VyaWFsaXphdGlvbi5kZWNvZGVyKFxuICAgICAgZmllbGRzW1wiY29kZWNJRFwiXSxcbiAgICAgIGVuY29kaW5nLFxuICAgICAgXCJkZWNpbWFsU3RyaW5nXCIsXG4gICAgICBcIkJ1ZmZlclwiLFxuICAgICAgMlxuICAgIClcbiAgICB0aGlzLnR4aWQgPSBzZXJpYWxpemF0aW9uLmRlY29kZXIoXG4gICAgICBmaWVsZHNbXCJ0eGlkXCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBcImNiNThcIixcbiAgICAgIFwiQnVmZmVyXCIsXG4gICAgICAzMlxuICAgIClcbiAgICB0aGlzLm91dHB1dGlkeCA9IHNlcmlhbGl6YXRpb24uZGVjb2RlcihcbiAgICAgIGZpZWxkc1tcIm91dHB1dGlkeFwiXSxcbiAgICAgIGVuY29kaW5nLFxuICAgICAgXCJkZWNpbWFsU3RyaW5nXCIsXG4gICAgICBcIkJ1ZmZlclwiLFxuICAgICAgNFxuICAgIClcbiAgICB0aGlzLmFzc2V0SUQgPSBzZXJpYWxpemF0aW9uLmRlY29kZXIoXG4gICAgICBmaWVsZHNbXCJhc3NldElEXCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBcImNiNThcIixcbiAgICAgIFwiQnVmZmVyXCIsXG4gICAgICAzMlxuICAgIClcbiAgfVxuXG4gIHByb3RlY3RlZCBjb2RlY0lEOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoMilcbiAgcHJvdGVjdGVkIHR4aWQ6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYygzMilcbiAgcHJvdGVjdGVkIG91dHB1dGlkeDogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDQpXG4gIHByb3RlY3RlZCBhc3NldElEOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoMzIpXG4gIHByb3RlY3RlZCBvdXRwdXQ6IE91dHB1dCA9IHVuZGVmaW5lZFxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBudW1lcmljIHJlcHJlc2VudGF0aW9uIG9mIHRoZSBDb2RlY0lELlxuICAgKi9cbiAgZ2V0Q29kZWNJRCA9ICgpOiAvKiBpc3RhbmJ1bCBpZ25vcmUgbmV4dCAqLyBudW1iZXIgPT5cbiAgICB0aGlzLmNvZGVjSUQucmVhZFVJbnQ4KDApXG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IHJlcHJlc2VudGF0aW9uIG9mIHRoZSBDb2RlY0lEXG4gICAqL1xuICBnZXRDb2RlY0lEQnVmZmVyID0gKCk6IEJ1ZmZlciA9PiB0aGlzLmNvZGVjSURcblxuICAvKipcbiAgICogUmV0dXJucyBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IG9mIHRoZSBUeElELlxuICAgKi9cbiAgZ2V0VHhJRCA9ICgpOiAvKiBpc3RhbmJ1bCBpZ25vcmUgbmV4dCAqLyBCdWZmZXIgPT4gdGhpcy50eGlkXG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSAgb2YgdGhlIE91dHB1dElkeC5cbiAgICovXG4gIGdldE91dHB1dElkeCA9ICgpOiAvKiBpc3RhbmJ1bCBpZ25vcmUgbmV4dCAqLyBCdWZmZXIgPT4gdGhpcy5vdXRwdXRpZHhcblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgYXNzZXRJRCBhcyBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9LlxuICAgKi9cbiAgZ2V0QXNzZXRJRCA9ICgpOiBCdWZmZXIgPT4gdGhpcy5hc3NldElEXG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIFVUWE9JRCBhcyBhIGJhc2UtNTggc3RyaW5nIChVVFhPSUQgaXMgYSBzdHJpbmcgKVxuICAgKi9cbiAgZ2V0VVRYT0lEID0gKCk6IC8qIGlzdGFuYnVsIGlnbm9yZSBuZXh0ICovIHN0cmluZyA9PlxuICAgIGJpbnRvb2xzLmJ1ZmZlclRvQjU4KEJ1ZmZlci5jb25jYXQoW3RoaXMuZ2V0VHhJRCgpLCB0aGlzLmdldE91dHB1dElkeCgpXSkpXG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSByZWZlcmVuY2UgdG8gdGhlIG91dHB1dFxuICAgKi9cbiAgZ2V0T3V0cHV0ID0gKCk6IE91dHB1dCA9PiB0aGlzLm91dHB1dFxuXG4gIC8qKlxuICAgKiBUYWtlcyBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IGNvbnRhaW5pbmcgYW4gW1tTdGFuZGFyZFVUWE9dXSwgcGFyc2VzIGl0LCBwb3B1bGF0ZXMgdGhlIGNsYXNzLCBhbmQgcmV0dXJucyB0aGUgbGVuZ3RoIG9mIHRoZSBTdGFuZGFyZFVUWE8gaW4gYnl0ZXMuXG4gICAqXG4gICAqIEBwYXJhbSBieXRlcyBBIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IGNvbnRhaW5pbmcgYSByYXcgW1tTdGFuZGFyZFVUWE9dXVxuICAgKi9cbiAgYWJzdHJhY3QgZnJvbUJ1ZmZlcihieXRlczogQnVmZmVyLCBvZmZzZXQ/OiBudW1iZXIpOiBudW1iZXJcblxuICAvKipcbiAgICogUmV0dXJucyBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IHJlcHJlc2VudGF0aW9uIG9mIHRoZSBbW1N0YW5kYXJkVVRYT11dLlxuICAgKi9cbiAgdG9CdWZmZXIoKTogQnVmZmVyIHtcbiAgICBjb25zdCBvdXRidWZmOiBCdWZmZXIgPSB0aGlzLm91dHB1dC50b0J1ZmZlcigpXG4gICAgY29uc3Qgb3V0cHV0aWRidWZmZXI6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYyg0KVxuICAgIG91dHB1dGlkYnVmZmVyLndyaXRlVUludDMyQkUodGhpcy5vdXRwdXQuZ2V0T3V0cHV0SUQoKSwgMClcbiAgICBjb25zdCBiYXJyOiBCdWZmZXJbXSA9IFtcbiAgICAgIHRoaXMuY29kZWNJRCxcbiAgICAgIHRoaXMudHhpZCxcbiAgICAgIHRoaXMub3V0cHV0aWR4LFxuICAgICAgdGhpcy5hc3NldElELFxuICAgICAgb3V0cHV0aWRidWZmZXIsXG4gICAgICBvdXRidWZmXG4gICAgXVxuICAgIHJldHVybiBCdWZmZXIuY29uY2F0KFxuICAgICAgYmFycixcbiAgICAgIHRoaXMuY29kZWNJRC5sZW5ndGggK1xuICAgICAgICB0aGlzLnR4aWQubGVuZ3RoICtcbiAgICAgICAgdGhpcy5vdXRwdXRpZHgubGVuZ3RoICtcbiAgICAgICAgdGhpcy5hc3NldElELmxlbmd0aCArXG4gICAgICAgIG91dHB1dGlkYnVmZmVyLmxlbmd0aCArXG4gICAgICAgIG91dGJ1ZmYubGVuZ3RoXG4gICAgKVxuICB9XG5cbiAgYWJzdHJhY3QgZnJvbVN0cmluZyhzZXJpYWxpemVkOiBzdHJpbmcpOiBudW1iZXJcblxuICBhYnN0cmFjdCB0b1N0cmluZygpOiBzdHJpbmdcblxuICBhYnN0cmFjdCBjbG9uZSgpOiB0aGlzXG5cbiAgYWJzdHJhY3QgY3JlYXRlKFxuICAgIGNvZGVjSUQ/OiBudW1iZXIsXG4gICAgdHhpZD86IEJ1ZmZlcixcbiAgICBvdXRwdXRpZHg/OiBCdWZmZXIgfCBudW1iZXIsXG4gICAgYXNzZXRJRD86IEJ1ZmZlcixcbiAgICBvdXRwdXQ/OiBPdXRwdXRcbiAgKTogdGhpc1xuXG4gIC8qKlxuICAgKiBDbGFzcyBmb3IgcmVwcmVzZW50aW5nIGEgc2luZ2xlIFN0YW5kYXJkVVRYTy5cbiAgICpcbiAgICogQHBhcmFtIGNvZGVjSUQgT3B0aW9uYWwgbnVtYmVyIHdoaWNoIHNwZWNpZmllcyB0aGUgY29kZUlEIG9mIHRoZSBVVFhPLiBEZWZhdWx0IDBcbiAgICogQHBhcmFtIHR4SUQgT3B0aW9uYWwge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gb2YgdHJhbnNhY3Rpb24gSUQgZm9yIHRoZSBTdGFuZGFyZFVUWE9cbiAgICogQHBhcmFtIHR4aWR4IE9wdGlvbmFsIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IG9yIG51bWJlciBmb3IgdGhlIGluZGV4IG9mIHRoZSB0cmFuc2FjdGlvbidzIFtbT3V0cHV0XV1cbiAgICogQHBhcmFtIGFzc2V0SUQgT3B0aW9uYWwge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gb2YgdGhlIGFzc2V0IElEIGZvciB0aGUgU3RhbmRhcmRVVFhPXG4gICAqIEBwYXJhbSBvdXRwdXRpZCBPcHRpb25hbCB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSBvciBudW1iZXIgb2YgdGhlIG91dHB1dCBJRCBmb3IgdGhlIFN0YW5kYXJkVVRYT1xuICAgKi9cbiAgY29uc3RydWN0b3IoXG4gICAgY29kZWNJRDogbnVtYmVyID0gMCxcbiAgICB0eElEOiBCdWZmZXIgPSB1bmRlZmluZWQsXG4gICAgb3V0cHV0aWR4OiBCdWZmZXIgfCBudW1iZXIgPSB1bmRlZmluZWQsXG4gICAgYXNzZXRJRDogQnVmZmVyID0gdW5kZWZpbmVkLFxuICAgIG91dHB1dDogT3V0cHV0ID0gdW5kZWZpbmVkXG4gICkge1xuICAgIHN1cGVyKClcbiAgICBpZiAodHlwZW9mIGNvZGVjSUQgIT09IFwidW5kZWZpbmVkXCIpIHtcbiAgICAgIHRoaXMuY29kZWNJRC53cml0ZVVJbnQ4KGNvZGVjSUQsIDApXG4gICAgfVxuICAgIGlmICh0eXBlb2YgdHhJRCAhPT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgdGhpcy50eGlkID0gdHhJRFxuICAgIH1cbiAgICBpZiAodHlwZW9mIG91dHB1dGlkeCA9PT0gXCJudW1iZXJcIikge1xuICAgICAgdGhpcy5vdXRwdXRpZHgud3JpdGVVSW50MzJCRShvdXRwdXRpZHgsIDApXG4gICAgfSBlbHNlIGlmIChvdXRwdXRpZHggaW5zdGFuY2VvZiBCdWZmZXIpIHtcbiAgICAgIHRoaXMub3V0cHV0aWR4ID0gb3V0cHV0aWR4XG4gICAgfVxuXG4gICAgaWYgKHR5cGVvZiBhc3NldElEICE9PSBcInVuZGVmaW5lZFwiKSB7XG4gICAgICB0aGlzLmFzc2V0SUQgPSBhc3NldElEXG4gICAgfVxuICAgIGlmICh0eXBlb2Ygb3V0cHV0ICE9PSBcInVuZGVmaW5lZFwiKSB7XG4gICAgICB0aGlzLm91dHB1dCA9IG91dHB1dFxuICAgIH1cbiAgfVxufVxuLyoqXG4gKiBDbGFzcyByZXByZXNlbnRpbmcgYSBzZXQgb2YgW1tTdGFuZGFyZFVUWE9dXXMuXG4gKi9cbmV4cG9ydCBhYnN0cmFjdCBjbGFzcyBTdGFuZGFyZFVUWE9TZXQ8XG4gIFVUWE9DbGFzcyBleHRlbmRzIFN0YW5kYXJkVVRYT1xuPiBleHRlbmRzIFNlcmlhbGl6YWJsZSB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIlN0YW5kYXJkVVRYT1NldFwiXG4gIHByb3RlY3RlZCBfdHlwZUlEID0gdW5kZWZpbmVkXG5cbiAgc2VyaWFsaXplKGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKTogb2JqZWN0IHtcbiAgICBsZXQgZmllbGRzOiBvYmplY3QgPSBzdXBlci5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgbGV0IHV0eG9zID0ge31cbiAgICBmb3IgKGxldCB1dHhvaWQgaW4gdGhpcy51dHhvcykge1xuICAgICAgbGV0IHV0eG9pZENsZWFuZWQ6IHN0cmluZyA9IHNlcmlhbGl6YXRpb24uZW5jb2RlcihcbiAgICAgICAgdXR4b2lkLFxuICAgICAgICBlbmNvZGluZyxcbiAgICAgICAgXCJiYXNlNThcIixcbiAgICAgICAgXCJiYXNlNThcIlxuICAgICAgKVxuICAgICAgdXR4b3NbYCR7dXR4b2lkQ2xlYW5lZH1gXSA9IHRoaXMudXR4b3NbYCR7dXR4b2lkfWBdLnNlcmlhbGl6ZShlbmNvZGluZylcbiAgICB9XG4gICAgbGV0IGFkZHJlc3NVVFhPcyA9IHt9XG4gICAgZm9yIChsZXQgYWRkcmVzcyBpbiB0aGlzLmFkZHJlc3NVVFhPcykge1xuICAgICAgbGV0IGFkZHJlc3NDbGVhbmVkOiBzdHJpbmcgPSBzZXJpYWxpemF0aW9uLmVuY29kZXIoXG4gICAgICAgIGFkZHJlc3MsXG4gICAgICAgIGVuY29kaW5nLFxuICAgICAgICBcImhleFwiLFxuICAgICAgICBcImNiNThcIlxuICAgICAgKVxuICAgICAgbGV0IHV0eG9iYWxhbmNlID0ge31cbiAgICAgIGZvciAobGV0IHV0eG9pZCBpbiB0aGlzLmFkZHJlc3NVVFhPc1tgJHthZGRyZXNzfWBdKSB7XG4gICAgICAgIGxldCB1dHhvaWRDbGVhbmVkOiBzdHJpbmcgPSBzZXJpYWxpemF0aW9uLmVuY29kZXIoXG4gICAgICAgICAgdXR4b2lkLFxuICAgICAgICAgIGVuY29kaW5nLFxuICAgICAgICAgIFwiYmFzZTU4XCIsXG4gICAgICAgICAgXCJiYXNlNThcIlxuICAgICAgICApXG4gICAgICAgIHV0eG9iYWxhbmNlW2Ake3V0eG9pZENsZWFuZWR9YF0gPSBzZXJpYWxpemF0aW9uLmVuY29kZXIoXG4gICAgICAgICAgdGhpcy5hZGRyZXNzVVRYT3NbYCR7YWRkcmVzc31gXVtgJHt1dHhvaWR9YF0sXG4gICAgICAgICAgZW5jb2RpbmcsXG4gICAgICAgICAgXCJCTlwiLFxuICAgICAgICAgIFwiZGVjaW1hbFN0cmluZ1wiXG4gICAgICAgIClcbiAgICAgIH1cbiAgICAgIGFkZHJlc3NVVFhPc1tgJHthZGRyZXNzQ2xlYW5lZH1gXSA9IHV0eG9iYWxhbmNlXG4gICAgfVxuICAgIHJldHVybiB7XG4gICAgICAuLi5maWVsZHMsXG4gICAgICB1dHhvcyxcbiAgICAgIGFkZHJlc3NVVFhPc1xuICAgIH1cbiAgfVxuXG4gIHByb3RlY3RlZCB1dHhvczogeyBbdXR4b2lkOiBzdHJpbmddOiBVVFhPQ2xhc3MgfSA9IHt9XG4gIHByb3RlY3RlZCBhZGRyZXNzVVRYT3M6IHsgW2FkZHJlc3M6IHN0cmluZ106IHsgW3V0eG9pZDogc3RyaW5nXTogQk4gfSB9ID0ge30gLy8gbWFwcyBhZGRyZXNzIHRvIHV0eG9pZHM6bG9ja3RpbWVcblxuICBhYnN0cmFjdCBwYXJzZVVUWE8odXR4bzogVVRYT0NsYXNzIHwgc3RyaW5nKTogVVRYT0NsYXNzXG5cbiAgLyoqXG4gICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgW1tTdGFuZGFyZFVUWE9dXSBpcyBpbiB0aGUgU3RhbmRhcmRVVFhPU2V0LlxuICAgKlxuICAgKiBAcGFyYW0gdXR4byBFaXRoZXIgYSBbW1N0YW5kYXJkVVRYT11dIGEgY2I1OCBzZXJpYWxpemVkIHN0cmluZyByZXByZXNlbnRpbmcgYSBTdGFuZGFyZFVUWE9cbiAgICovXG4gIGluY2x1ZGVzID0gKHV0eG86IFVUWE9DbGFzcyB8IHN0cmluZyk6IGJvb2xlYW4gPT4ge1xuICAgIGxldCB1dHhvWDogVVRYT0NsYXNzID0gdW5kZWZpbmVkXG4gICAgbGV0IHV0eG9pZDogc3RyaW5nID0gdW5kZWZpbmVkXG4gICAgdHJ5IHtcbiAgICAgIHV0eG9YID0gdGhpcy5wYXJzZVVUWE8odXR4bylcbiAgICAgIHV0eG9pZCA9IHV0eG9YLmdldFVUWE9JRCgpXG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgaWYgKGUgaW5zdGFuY2VvZiBFcnJvcikge1xuICAgICAgICBjb25zb2xlLmxvZyhlLm1lc3NhZ2UpXG4gICAgICB9IGVsc2Uge1xuICAgICAgICBjb25zb2xlLmxvZyhlKVxuICAgICAgfVxuICAgICAgcmV0dXJuIGZhbHNlXG4gICAgfVxuICAgIHJldHVybiB1dHhvaWQgaW4gdGhpcy51dHhvc1xuICB9XG5cbiAgLyoqXG4gICAqIEFkZHMgYSBbW1N0YW5kYXJkVVRYT11dIHRvIHRoZSBTdGFuZGFyZFVUWE9TZXQuXG4gICAqXG4gICAqIEBwYXJhbSB1dHhvIEVpdGhlciBhIFtbU3RhbmRhcmRVVFhPXV0gYW4gY2I1OCBzZXJpYWxpemVkIHN0cmluZyByZXByZXNlbnRpbmcgYSBTdGFuZGFyZFVUWE9cbiAgICogQHBhcmFtIG92ZXJ3cml0ZSBJZiB0cnVlLCBpZiB0aGUgVVRYT0lEIGFscmVhZHkgZXhpc3RzLCBvdmVyd3JpdGUgaXQuLi4gZGVmYXVsdCBmYWxzZVxuICAgKlxuICAgKiBAcmV0dXJucyBBIFtbU3RhbmRhcmRVVFhPXV0gaWYgb25lIHdhcyBhZGRlZCBhbmQgdW5kZWZpbmVkIGlmIG5vdGhpbmcgd2FzIGFkZGVkLlxuICAgKi9cbiAgYWRkKHV0eG86IFVUWE9DbGFzcyB8IHN0cmluZywgb3ZlcndyaXRlOiBib29sZWFuID0gZmFsc2UpOiBVVFhPQ2xhc3Mge1xuICAgIGxldCB1dHhvdmFyOiBVVFhPQ2xhc3MgPSB1bmRlZmluZWRcbiAgICB0cnkge1xuICAgICAgdXR4b3ZhciA9IHRoaXMucGFyc2VVVFhPKHV0eG8pXG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgaWYgKGUgaW5zdGFuY2VvZiBFcnJvcikge1xuICAgICAgICBjb25zb2xlLmxvZyhlLm1lc3NhZ2UpXG4gICAgICB9IGVsc2Uge1xuICAgICAgICBjb25zb2xlLmxvZyhlKVxuICAgICAgfVxuICAgICAgcmV0dXJuIHVuZGVmaW5lZFxuICAgIH1cblxuICAgIGNvbnN0IHV0eG9pZDogc3RyaW5nID0gdXR4b3Zhci5nZXRVVFhPSUQoKVxuICAgIGlmICghKHV0eG9pZCBpbiB0aGlzLnV0eG9zKSB8fCBvdmVyd3JpdGUgPT09IHRydWUpIHtcbiAgICAgIHRoaXMudXR4b3NbYCR7dXR4b2lkfWBdID0gdXR4b3ZhclxuICAgICAgY29uc3QgYWRkcmVzc2VzOiBCdWZmZXJbXSA9IHV0eG92YXIuZ2V0T3V0cHV0KCkuZ2V0QWRkcmVzc2VzKClcbiAgICAgIGNvbnN0IGxvY2t0aW1lOiBCTiA9IHV0eG92YXIuZ2V0T3V0cHV0KCkuZ2V0TG9ja3RpbWUoKVxuICAgICAgZm9yIChsZXQgaTogbnVtYmVyID0gMDsgaSA8IGFkZHJlc3Nlcy5sZW5ndGg7IGkrKykge1xuICAgICAgICBjb25zdCBhZGRyZXNzOiBzdHJpbmcgPSBhZGRyZXNzZXNbYCR7aX1gXS50b1N0cmluZyhcImhleFwiKVxuICAgICAgICBpZiAoIShhZGRyZXNzIGluIHRoaXMuYWRkcmVzc1VUWE9zKSkge1xuICAgICAgICAgIHRoaXMuYWRkcmVzc1VUWE9zW2Ake2FkZHJlc3N9YF0gPSB7fVxuICAgICAgICB9XG4gICAgICAgIHRoaXMuYWRkcmVzc1VUWE9zW2Ake2FkZHJlc3N9YF1bYCR7dXR4b2lkfWBdID0gbG9ja3RpbWVcbiAgICAgIH1cbiAgICAgIHJldHVybiB1dHhvdmFyXG4gICAgfVxuICAgIHJldHVybiB1bmRlZmluZWRcbiAgfVxuXG4gIC8qKlxuICAgKiBBZGRzIGFuIGFycmF5IG9mIFtbU3RhbmRhcmRVVFhPXV1zIHRvIHRoZSBbW1N0YW5kYXJkVVRYT1NldF1dLlxuICAgKlxuICAgKiBAcGFyYW0gdXR4byBFaXRoZXIgYSBbW1N0YW5kYXJkVVRYT11dIGFuIGNiNTggc2VyaWFsaXplZCBzdHJpbmcgcmVwcmVzZW50aW5nIGEgU3RhbmRhcmRVVFhPXG4gICAqIEBwYXJhbSBvdmVyd3JpdGUgSWYgdHJ1ZSwgaWYgdGhlIFVUWE9JRCBhbHJlYWR5IGV4aXN0cywgb3ZlcndyaXRlIGl0Li4uIGRlZmF1bHQgZmFsc2VcbiAgICpcbiAgICogQHJldHVybnMgQW4gYXJyYXkgb2YgU3RhbmRhcmRVVFhPcyB3aGljaCB3ZXJlIGFkZGVkLlxuICAgKi9cbiAgYWRkQXJyYXkoXG4gICAgdXR4b3M6IHN0cmluZ1tdIHwgVVRYT0NsYXNzW10sXG4gICAgb3ZlcndyaXRlOiBib29sZWFuID0gZmFsc2VcbiAgKTogU3RhbmRhcmRVVFhPW10ge1xuICAgIGNvbnN0IGFkZGVkOiBVVFhPQ2xhc3NbXSA9IFtdXG4gICAgZm9yIChsZXQgaTogbnVtYmVyID0gMDsgaSA8IHV0eG9zLmxlbmd0aDsgaSsrKSB7XG4gICAgICBsZXQgcmVzdWx0OiBVVFhPQ2xhc3MgPSB0aGlzLmFkZCh1dHhvc1tgJHtpfWBdLCBvdmVyd3JpdGUpXG4gICAgICBpZiAodHlwZW9mIHJlc3VsdCAhPT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgICBhZGRlZC5wdXNoKHJlc3VsdClcbiAgICAgIH1cbiAgICB9XG4gICAgcmV0dXJuIGFkZGVkXG4gIH1cblxuICAvKipcbiAgICogUmVtb3ZlcyBhIFtbU3RhbmRhcmRVVFhPXV0gZnJvbSB0aGUgW1tTdGFuZGFyZFVUWE9TZXRdXSBpZiBpdCBleGlzdHMuXG4gICAqXG4gICAqIEBwYXJhbSB1dHhvIEVpdGhlciBhIFtbU3RhbmRhcmRVVFhPXV0gYW4gY2I1OCBzZXJpYWxpemVkIHN0cmluZyByZXByZXNlbnRpbmcgYSBTdGFuZGFyZFVUWE9cbiAgICpcbiAgICogQHJldHVybnMgQSBbW1N0YW5kYXJkVVRYT11dIGlmIGl0IHdhcyByZW1vdmVkIGFuZCB1bmRlZmluZWQgaWYgbm90aGluZyB3YXMgcmVtb3ZlZC5cbiAgICovXG4gIHJlbW92ZSA9ICh1dHhvOiBVVFhPQ2xhc3MgfCBzdHJpbmcpOiBVVFhPQ2xhc3MgPT4ge1xuICAgIGxldCB1dHhvdmFyOiBVVFhPQ2xhc3MgPSB1bmRlZmluZWRcbiAgICB0cnkge1xuICAgICAgdXR4b3ZhciA9IHRoaXMucGFyc2VVVFhPKHV0eG8pXG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgaWYgKGUgaW5zdGFuY2VvZiBFcnJvcikge1xuICAgICAgICBjb25zb2xlLmxvZyhlLm1lc3NhZ2UpXG4gICAgICB9IGVsc2Uge1xuICAgICAgICBjb25zb2xlLmxvZyhlKVxuICAgICAgfVxuICAgICAgcmV0dXJuIHVuZGVmaW5lZFxuICAgIH1cblxuICAgIGNvbnN0IHV0eG9pZDogc3RyaW5nID0gdXR4b3Zhci5nZXRVVFhPSUQoKVxuICAgIGlmICghKHV0eG9pZCBpbiB0aGlzLnV0eG9zKSkge1xuICAgICAgcmV0dXJuIHVuZGVmaW5lZFxuICAgIH1cbiAgICBkZWxldGUgdGhpcy51dHhvc1tgJHt1dHhvaWR9YF1cbiAgICBjb25zdCBhZGRyZXNzZXMgPSBPYmplY3Qua2V5cyh0aGlzLmFkZHJlc3NVVFhPcylcbiAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgYWRkcmVzc2VzLmxlbmd0aDsgaSsrKSB7XG4gICAgICBpZiAodXR4b2lkIGluIHRoaXMuYWRkcmVzc1VUWE9zW2FkZHJlc3Nlc1tgJHtpfWBdXSkge1xuICAgICAgICBkZWxldGUgdGhpcy5hZGRyZXNzVVRYT3NbYWRkcmVzc2VzW2Ake2l9YF1dW2Ake3V0eG9pZH1gXVxuICAgICAgfVxuICAgIH1cbiAgICByZXR1cm4gdXR4b3ZhclxuICB9XG5cbiAgLyoqXG4gICAqIFJlbW92ZXMgYW4gYXJyYXkgb2YgW1tTdGFuZGFyZFVUWE9dXXMgdG8gdGhlIFtbU3RhbmRhcmRVVFhPU2V0XV0uXG4gICAqXG4gICAqIEBwYXJhbSB1dHhvIEVpdGhlciBhIFtbU3RhbmRhcmRVVFhPXV0gYW4gY2I1OCBzZXJpYWxpemVkIHN0cmluZyByZXByZXNlbnRpbmcgYSBTdGFuZGFyZFVUWE9cbiAgICogQHBhcmFtIG92ZXJ3cml0ZSBJZiB0cnVlLCBpZiB0aGUgVVRYT0lEIGFscmVhZHkgZXhpc3RzLCBvdmVyd3JpdGUgaXQuLi4gZGVmYXVsdCBmYWxzZVxuICAgKlxuICAgKiBAcmV0dXJucyBBbiBhcnJheSBvZiBVVFhPcyB3aGljaCB3ZXJlIHJlbW92ZWQuXG4gICAqL1xuICByZW1vdmVBcnJheSA9ICh1dHhvczogc3RyaW5nW10gfCBVVFhPQ2xhc3NbXSk6IFVUWE9DbGFzc1tdID0+IHtcbiAgICBjb25zdCByZW1vdmVkOiBVVFhPQ2xhc3NbXSA9IFtdXG4gICAgZm9yIChsZXQgaTogbnVtYmVyID0gMDsgaSA8IHV0eG9zLmxlbmd0aDsgaSsrKSB7XG4gICAgICBjb25zdCByZXN1bHQ6IFVUWE9DbGFzcyA9IHRoaXMucmVtb3ZlKHV0eG9zW2Ake2l9YF0pXG4gICAgICBpZiAodHlwZW9mIHJlc3VsdCAhPT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgICByZW1vdmVkLnB1c2gocmVzdWx0KVxuICAgICAgfVxuICAgIH1cbiAgICByZXR1cm4gcmVtb3ZlZFxuICB9XG5cbiAgLyoqXG4gICAqIEdldHMgYSBbW1N0YW5kYXJkVVRYT11dIGZyb20gdGhlIFtbU3RhbmRhcmRVVFhPU2V0XV0gYnkgaXRzIFVUWE9JRC5cbiAgICpcbiAgICogQHBhcmFtIHV0eG9pZCBTdHJpbmcgcmVwcmVzZW50aW5nIHRoZSBVVFhPSURcbiAgICpcbiAgICogQHJldHVybnMgQSBbW1N0YW5kYXJkVVRYT11dIGlmIGl0IGV4aXN0cyBpbiB0aGUgc2V0LlxuICAgKi9cbiAgZ2V0VVRYTyA9ICh1dHhvaWQ6IHN0cmluZyk6IFVUWE9DbGFzcyA9PiB0aGlzLnV0eG9zW2Ake3V0eG9pZH1gXVxuXG4gIC8qKlxuICAgKiBHZXRzIGFsbCB0aGUgW1tTdGFuZGFyZFVUWE9dXXMsIG9wdGlvbmFsbHkgdGhhdCBtYXRjaCB3aXRoIFVUWE9JRHMgaW4gYW4gYXJyYXlcbiAgICpcbiAgICogQHBhcmFtIHV0eG9pZHMgQW4gb3B0aW9uYWwgYXJyYXkgb2YgVVRYT0lEcywgcmV0dXJucyBhbGwgW1tTdGFuZGFyZFVUWE9dXXMgaWYgbm90IHByb3ZpZGVkXG4gICAqXG4gICAqIEByZXR1cm5zIEFuIGFycmF5IG9mIFtbU3RhbmRhcmRVVFhPXV1zLlxuICAgKi9cbiAgZ2V0QWxsVVRYT3MgPSAodXR4b2lkczogc3RyaW5nW10gPSB1bmRlZmluZWQpOiBVVFhPQ2xhc3NbXSA9PiB7XG4gICAgbGV0IHJlc3VsdHM6IFVUWE9DbGFzc1tdID0gW11cbiAgICBpZiAodHlwZW9mIHV0eG9pZHMgIT09IFwidW5kZWZpbmVkXCIgJiYgQXJyYXkuaXNBcnJheSh1dHhvaWRzKSkge1xuICAgICAgcmVzdWx0cyA9IHV0eG9pZHNcbiAgICAgICAgLmZpbHRlcigodXR4b2lkKSA9PiB0aGlzLnV0eG9zW2Ake3V0eG9pZH1gXSlcbiAgICAgICAgLm1hcCgodXR4b2lkKSA9PiB0aGlzLnV0eG9zW2Ake3V0eG9pZH1gXSlcbiAgICB9IGVsc2Uge1xuICAgICAgcmVzdWx0cyA9IE9iamVjdC52YWx1ZXModGhpcy51dHhvcylcbiAgICB9XG4gICAgcmV0dXJuIHJlc3VsdHNcbiAgfVxuXG4gIC8qKlxuICAgKiBHZXRzIGFsbCB0aGUgW1tTdGFuZGFyZFVUWE9dXXMgYXMgc3RyaW5ncywgb3B0aW9uYWxseSB0aGF0IG1hdGNoIHdpdGggVVRYT0lEcyBpbiBhbiBhcnJheS5cbiAgICpcbiAgICogQHBhcmFtIHV0eG9pZHMgQW4gb3B0aW9uYWwgYXJyYXkgb2YgVVRYT0lEcywgcmV0dXJucyBhbGwgW1tTdGFuZGFyZFVUWE9dXXMgaWYgbm90IHByb3ZpZGVkXG4gICAqXG4gICAqIEByZXR1cm5zIEFuIGFycmF5IG9mIFtbU3RhbmRhcmRVVFhPXV1zIGFzIGNiNTggc2VyaWFsaXplZCBzdHJpbmdzLlxuICAgKi9cbiAgZ2V0QWxsVVRYT1N0cmluZ3MgPSAodXR4b2lkczogc3RyaW5nW10gPSB1bmRlZmluZWQpOiBzdHJpbmdbXSA9PiB7XG4gICAgY29uc3QgcmVzdWx0czogc3RyaW5nW10gPSBbXVxuICAgIGNvbnN0IHV0eG9zID0gT2JqZWN0LmtleXModGhpcy51dHhvcylcbiAgICBpZiAodHlwZW9mIHV0eG9pZHMgIT09IFwidW5kZWZpbmVkXCIgJiYgQXJyYXkuaXNBcnJheSh1dHhvaWRzKSkge1xuICAgICAgZm9yIChsZXQgaTogbnVtYmVyID0gMDsgaSA8IHV0eG9pZHMubGVuZ3RoOyBpKyspIHtcbiAgICAgICAgaWYgKHV0eG9pZHNbYCR7aX1gXSBpbiB0aGlzLnV0eG9zKSB7XG4gICAgICAgICAgcmVzdWx0cy5wdXNoKHRoaXMudXR4b3NbdXR4b2lkc1tgJHtpfWBdXS50b1N0cmluZygpKVxuICAgICAgICB9XG4gICAgICB9XG4gICAgfSBlbHNlIHtcbiAgICAgIGZvciAoY29uc3QgdSBvZiB1dHhvcykge1xuICAgICAgICByZXN1bHRzLnB1c2godGhpcy51dHhvc1tgJHt1fWBdLnRvU3RyaW5nKCkpXG4gICAgICB9XG4gICAgfVxuICAgIHJldHVybiByZXN1bHRzXG4gIH1cblxuICAvKipcbiAgICogR2l2ZW4gYW4gYWRkcmVzcyBvciBhcnJheSBvZiBhZGRyZXNzZXMsIHJldHVybnMgYWxsIHRoZSBVVFhPSURzIGZvciB0aG9zZSBhZGRyZXNzZXNcbiAgICpcbiAgICogQHBhcmFtIGFkZHJlc3MgQW4gYXJyYXkgb2YgYWRkcmVzcyB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfXNcbiAgICogQHBhcmFtIHNwZW5kYWJsZSBJZiB0cnVlLCBvbmx5IHJldHJpZXZlcyBVVFhPSURzIHdob3NlIGxvY2t0aW1lIGhhcyBwYXNzZWRcbiAgICpcbiAgICogQHJldHVybnMgQW4gYXJyYXkgb2YgYWRkcmVzc2VzLlxuICAgKi9cbiAgZ2V0VVRYT0lEcyA9IChcbiAgICBhZGRyZXNzZXM6IEJ1ZmZlcltdID0gdW5kZWZpbmVkLFxuICAgIHNwZW5kYWJsZTogYm9vbGVhbiA9IHRydWVcbiAgKTogc3RyaW5nW10gPT4ge1xuICAgIGlmICh0eXBlb2YgYWRkcmVzc2VzICE9PSBcInVuZGVmaW5lZFwiKSB7XG4gICAgICBjb25zdCByZXN1bHRzOiBzdHJpbmdbXSA9IFtdXG4gICAgICBjb25zdCBub3c6IEJOID0gVW5peE5vdygpXG4gICAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgYWRkcmVzc2VzLmxlbmd0aDsgaSsrKSB7XG4gICAgICAgIGlmIChhZGRyZXNzZXNbYCR7aX1gXS50b1N0cmluZyhcImhleFwiKSBpbiB0aGlzLmFkZHJlc3NVVFhPcykge1xuICAgICAgICAgIGNvbnN0IGVudHJpZXMgPSBPYmplY3QuZW50cmllcyhcbiAgICAgICAgICAgIHRoaXMuYWRkcmVzc1VUWE9zW2FkZHJlc3Nlc1tgJHtpfWBdLnRvU3RyaW5nKFwiaGV4XCIpXVxuICAgICAgICAgIClcbiAgICAgICAgICBmb3IgKGNvbnN0IFt1dHhvaWQsIGxvY2t0aW1lXSBvZiBlbnRyaWVzKSB7XG4gICAgICAgICAgICBpZiAoXG4gICAgICAgICAgICAgIChyZXN1bHRzLmluZGV4T2YodXR4b2lkKSA9PT0gLTEgJiZcbiAgICAgICAgICAgICAgICBzcGVuZGFibGUgJiZcbiAgICAgICAgICAgICAgICBsb2NrdGltZS5sdGUobm93KSkgfHxcbiAgICAgICAgICAgICAgIXNwZW5kYWJsZVxuICAgICAgICAgICAgKSB7XG4gICAgICAgICAgICAgIHJlc3VsdHMucHVzaCh1dHhvaWQpXG4gICAgICAgICAgICB9XG4gICAgICAgICAgfVxuICAgICAgICB9XG4gICAgICB9XG4gICAgICByZXR1cm4gcmVzdWx0c1xuICAgIH1cbiAgICByZXR1cm4gT2JqZWN0LmtleXModGhpcy51dHhvcylcbiAgfVxuXG4gIC8qKlxuICAgKiBHZXRzIHRoZSBhZGRyZXNzZXMgaW4gdGhlIFtbU3RhbmRhcmRVVFhPU2V0XV0gYW5kIHJldHVybnMgYW4gYXJyYXkgb2Yge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0uXG4gICAqL1xuICBnZXRBZGRyZXNzZXMgPSAoKTogQnVmZmVyW10gPT5cbiAgICBPYmplY3Qua2V5cyh0aGlzLmFkZHJlc3NVVFhPcykubWFwKChrKSA9PiBCdWZmZXIuZnJvbShrLCBcImhleFwiKSlcblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgYmFsYW5jZSBvZiBhIHNldCBvZiBhZGRyZXNzZXMgaW4gdGhlIFN0YW5kYXJkVVRYT1NldC5cbiAgICpcbiAgICogQHBhcmFtIGFkZHJlc3NlcyBBbiBhcnJheSBvZiBhZGRyZXNzZXNcbiAgICogQHBhcmFtIGFzc2V0SUQgRWl0aGVyIGEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gb3IgYW4gY2I1OCBzZXJpYWxpemVkIHJlcHJlc2VudGF0aW9uIG9mIGFuIEFzc2V0SURcbiAgICogQHBhcmFtIGFzT2YgVGhlIHRpbWVzdGFtcCB0byB2ZXJpZnkgdGhlIHRyYW5zYWN0aW9uIGFnYWluc3QgYXMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2luZHV0bnkvYm4uanMvfEJOfVxuICAgKlxuICAgKiBAcmV0dXJucyBSZXR1cm5zIHRoZSB0b3RhbCBiYWxhbmNlIGFzIGEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9pbmR1dG55L2JuLmpzL3xCTn0uXG4gICAqL1xuICBnZXRCYWxhbmNlID0gKFxuICAgIGFkZHJlc3NlczogQnVmZmVyW10sXG4gICAgYXNzZXRJRDogQnVmZmVyIHwgc3RyaW5nLFxuICAgIGFzT2Y6IEJOID0gdW5kZWZpbmVkXG4gICk6IEJOID0+IHtcbiAgICBjb25zdCB1dHhvaWRzOiBzdHJpbmdbXSA9IHRoaXMuZ2V0VVRYT0lEcyhhZGRyZXNzZXMpXG4gICAgY29uc3QgdXR4b3M6IFN0YW5kYXJkVVRYT1tdID0gdGhpcy5nZXRBbGxVVFhPcyh1dHhvaWRzKVxuICAgIGxldCBzcGVuZDogQk4gPSBuZXcgQk4oMClcbiAgICBsZXQgYXNzZXQ6IEJ1ZmZlclxuICAgIGlmICh0eXBlb2YgYXNzZXRJRCA9PT0gXCJzdHJpbmdcIikge1xuICAgICAgYXNzZXQgPSBiaW50b29scy5jYjU4RGVjb2RlKGFzc2V0SUQpXG4gICAgfSBlbHNlIHtcbiAgICAgIGFzc2V0ID0gYXNzZXRJRFxuICAgIH1cbiAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgdXR4b3MubGVuZ3RoOyBpKyspIHtcbiAgICAgIGlmIChcbiAgICAgICAgdXR4b3NbYCR7aX1gXS5nZXRPdXRwdXQoKSBpbnN0YW5jZW9mIFN0YW5kYXJkQW1vdW50T3V0cHV0ICYmXG4gICAgICAgIHV0eG9zW2Ake2l9YF0uZ2V0QXNzZXRJRCgpLnRvU3RyaW5nKFwiaGV4XCIpID09PSBhc3NldC50b1N0cmluZyhcImhleFwiKSAmJlxuICAgICAgICB1dHhvc1tgJHtpfWBdLmdldE91dHB1dCgpLm1lZXRzVGhyZXNob2xkKGFkZHJlc3NlcywgYXNPZilcbiAgICAgICkge1xuICAgICAgICBzcGVuZCA9IHNwZW5kLmFkZChcbiAgICAgICAgICAodXR4b3NbYCR7aX1gXS5nZXRPdXRwdXQoKSBhcyBTdGFuZGFyZEFtb3VudE91dHB1dCkuZ2V0QW1vdW50KClcbiAgICAgICAgKVxuICAgICAgfVxuICAgIH1cbiAgICByZXR1cm4gc3BlbmRcbiAgfVxuXG4gIC8qKlxuICAgKiBHZXRzIGFsbCB0aGUgQXNzZXQgSURzLCBvcHRpb25hbGx5IHRoYXQgbWF0Y2ggd2l0aCBBc3NldCBJRHMgaW4gYW4gYXJyYXlcbiAgICpcbiAgICogQHBhcmFtIHV0eG9pZHMgQW4gb3B0aW9uYWwgYXJyYXkgb2YgQWRkcmVzc2VzIGFzIHN0cmluZyBvciBCdWZmZXIsIHJldHVybnMgYWxsIEFzc2V0IElEcyBpZiBub3QgcHJvdmlkZWRcbiAgICpcbiAgICogQHJldHVybnMgQW4gYXJyYXkgb2Yge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gcmVwcmVzZW50aW5nIHRoZSBBc3NldCBJRHMuXG4gICAqL1xuICBnZXRBc3NldElEcyA9IChhZGRyZXNzZXM6IEJ1ZmZlcltdID0gdW5kZWZpbmVkKTogQnVmZmVyW10gPT4ge1xuICAgIGNvbnN0IHJlc3VsdHM6IFNldDxCdWZmZXI+ID0gbmV3IFNldCgpXG4gICAgbGV0IHV0eG9pZHM6IHN0cmluZ1tdID0gW11cbiAgICBpZiAodHlwZW9mIGFkZHJlc3NlcyAhPT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgdXR4b2lkcyA9IHRoaXMuZ2V0VVRYT0lEcyhhZGRyZXNzZXMpXG4gICAgfSBlbHNlIHtcbiAgICAgIHV0eG9pZHMgPSB0aGlzLmdldFVUWE9JRHMoKVxuICAgIH1cblxuICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCB1dHhvaWRzLmxlbmd0aDsgaSsrKSB7XG4gICAgICBpZiAodXR4b2lkc1tgJHtpfWBdIGluIHRoaXMudXR4b3MgJiYgISh1dHhvaWRzW2Ake2l9YF0gaW4gcmVzdWx0cykpIHtcbiAgICAgICAgcmVzdWx0cy5hZGQodGhpcy51dHhvc1t1dHhvaWRzW2Ake2l9YF1dLmdldEFzc2V0SUQoKSlcbiAgICAgIH1cbiAgICB9XG5cbiAgICByZXR1cm4gWy4uLnJlc3VsdHNdXG4gIH1cblxuICBhYnN0cmFjdCBjbG9uZSgpOiB0aGlzXG5cbiAgYWJzdHJhY3QgY3JlYXRlKC4uLmFyZ3M6IGFueVtdKTogdGhpc1xuXG4gIGZpbHRlcihcbiAgICBhcmdzOiBhbnlbXSxcbiAgICBsYW1iZGE6ICh1dHhvOiBVVFhPQ2xhc3MsIC4uLmxhcmdzOiBhbnlbXSkgPT4gYm9vbGVhblxuICApOiB0aGlzIHtcbiAgICBsZXQgbmV3c2V0OiB0aGlzID0gdGhpcy5jbG9uZSgpXG4gICAgbGV0IHV0eG9zOiBVVFhPQ2xhc3NbXSA9IHRoaXMuZ2V0QWxsVVRYT3MoKVxuICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCB1dHhvcy5sZW5ndGg7IGkrKykge1xuICAgICAgaWYgKGxhbWJkYSh1dHhvc1tgJHtpfWBdLCAuLi5hcmdzKSA9PT0gZmFsc2UpIHtcbiAgICAgICAgbmV3c2V0LnJlbW92ZSh1dHhvc1tgJHtpfWBdKVxuICAgICAgfVxuICAgIH1cbiAgICByZXR1cm4gbmV3c2V0XG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyBhIG5ldyBzZXQgd2l0aCBjb3B5IG9mIFVUWE9zIGluIHRoaXMgYW5kIHNldCBwYXJhbWV0ZXIuXG4gICAqXG4gICAqIEBwYXJhbSB1dHhvc2V0IFRoZSBbW1N0YW5kYXJkVVRYT1NldF1dIHRvIG1lcmdlIHdpdGggdGhpcyBvbmVcbiAgICogQHBhcmFtIGhhc1VUWE9JRHMgV2lsbCBzdWJzZWxlY3QgYSBzZXQgb2YgW1tTdGFuZGFyZFVUWE9dXXMgd2hpY2ggaGF2ZSB0aGUgVVRYT0lEcyBwcm92aWRlZCBpbiB0aGlzIGFycmF5LCBkZWZ1bHRzIHRvIGFsbCBVVFhPc1xuICAgKlxuICAgKiBAcmV0dXJucyBBIG5ldyBTdGFuZGFyZFVUWE9TZXQgdGhhdCBjb250YWlucyBhbGwgdGhlIGZpbHRlcmVkIGVsZW1lbnRzLlxuICAgKi9cbiAgbWVyZ2UgPSAodXR4b3NldDogdGhpcywgaGFzVVRYT0lEczogc3RyaW5nW10gPSB1bmRlZmluZWQpOiB0aGlzID0+IHtcbiAgICBjb25zdCByZXN1bHRzOiB0aGlzID0gdGhpcy5jcmVhdGUoKVxuICAgIGNvbnN0IHV0eG9zMTogVVRYT0NsYXNzW10gPSB0aGlzLmdldEFsbFVUWE9zKGhhc1VUWE9JRHMpXG4gICAgY29uc3QgdXR4b3MyOiBVVFhPQ2xhc3NbXSA9IHV0eG9zZXQuZ2V0QWxsVVRYT3MoaGFzVVRYT0lEcylcbiAgICBjb25zdCBwcm9jZXNzID0gKHV0eG86IFVUWE9DbGFzcykgPT4ge1xuICAgICAgcmVzdWx0cy5hZGQodXR4bylcbiAgICB9XG4gICAgdXR4b3MxLmZvckVhY2gocHJvY2VzcylcbiAgICB1dHhvczIuZm9yRWFjaChwcm9jZXNzKVxuICAgIHJldHVybiByZXN1bHRzIGFzIHRoaXNcbiAgfVxuXG4gIC8qKlxuICAgKiBTZXQgaW50ZXJzZXRpb24gYmV0d2VlbiB0aGlzIHNldCBhbmQgYSBwYXJhbWV0ZXIuXG4gICAqXG4gICAqIEBwYXJhbSB1dHhvc2V0IFRoZSBzZXQgdG8gaW50ZXJzZWN0XG4gICAqXG4gICAqIEByZXR1cm5zIEEgbmV3IFN0YW5kYXJkVVRYT1NldCBjb250YWluaW5nIHRoZSBpbnRlcnNlY3Rpb25cbiAgICovXG4gIGludGVyc2VjdGlvbiA9ICh1dHhvc2V0OiB0aGlzKTogdGhpcyA9PiB7XG4gICAgY29uc3QgdXMxOiBzdHJpbmdbXSA9IHRoaXMuZ2V0VVRYT0lEcygpXG4gICAgY29uc3QgdXMyOiBzdHJpbmdbXSA9IHV0eG9zZXQuZ2V0VVRYT0lEcygpXG4gICAgY29uc3QgcmVzdWx0czogc3RyaW5nW10gPSB1czEuZmlsdGVyKCh1dHhvaWQpID0+IHVzMi5pbmNsdWRlcyh1dHhvaWQpKVxuICAgIHJldHVybiB0aGlzLm1lcmdlKHV0eG9zZXQsIHJlc3VsdHMpIGFzIHRoaXNcbiAgfVxuXG4gIC8qKlxuICAgKiBTZXQgZGlmZmVyZW5jZSBiZXR3ZWVuIHRoaXMgc2V0IGFuZCBhIHBhcmFtZXRlci5cbiAgICpcbiAgICogQHBhcmFtIHV0eG9zZXQgVGhlIHNldCB0byBkaWZmZXJlbmNlXG4gICAqXG4gICAqIEByZXR1cm5zIEEgbmV3IFN0YW5kYXJkVVRYT1NldCBjb250YWluaW5nIHRoZSBkaWZmZXJlbmNlXG4gICAqL1xuICBkaWZmZXJlbmNlID0gKHV0eG9zZXQ6IHRoaXMpOiB0aGlzID0+IHtcbiAgICBjb25zdCB1czE6IHN0cmluZ1tdID0gdGhpcy5nZXRVVFhPSURzKClcbiAgICBjb25zdCB1czI6IHN0cmluZ1tdID0gdXR4b3NldC5nZXRVVFhPSURzKClcbiAgICBjb25zdCByZXN1bHRzOiBzdHJpbmdbXSA9IHVzMS5maWx0ZXIoKHV0eG9pZCkgPT4gIXVzMi5pbmNsdWRlcyh1dHhvaWQpKVxuICAgIHJldHVybiB0aGlzLm1lcmdlKHV0eG9zZXQsIHJlc3VsdHMpIGFzIHRoaXNcbiAgfVxuXG4gIC8qKlxuICAgKiBTZXQgc3ltbWV0cmljYWwgZGlmZmVyZW5jZSBiZXR3ZWVuIHRoaXMgc2V0IGFuZCBhIHBhcmFtZXRlci5cbiAgICpcbiAgICogQHBhcmFtIHV0eG9zZXQgVGhlIHNldCB0byBzeW1tZXRyaWNhbCBkaWZmZXJlbmNlXG4gICAqXG4gICAqIEByZXR1cm5zIEEgbmV3IFN0YW5kYXJkVVRYT1NldCBjb250YWluaW5nIHRoZSBzeW1tZXRyaWNhbCBkaWZmZXJlbmNlXG4gICAqL1xuICBzeW1EaWZmZXJlbmNlID0gKHV0eG9zZXQ6IHRoaXMpOiB0aGlzID0+IHtcbiAgICBjb25zdCB1czE6IHN0cmluZ1tdID0gdGhpcy5nZXRVVFhPSURzKClcbiAgICBjb25zdCB1czI6IHN0cmluZ1tdID0gdXR4b3NldC5nZXRVVFhPSURzKClcbiAgICBjb25zdCByZXN1bHRzOiBzdHJpbmdbXSA9IHVzMVxuICAgICAgLmZpbHRlcigodXR4b2lkKSA9PiAhdXMyLmluY2x1ZGVzKHV0eG9pZCkpXG4gICAgICAuY29uY2F0KHVzMi5maWx0ZXIoKHV0eG9pZCkgPT4gIXVzMS5pbmNsdWRlcyh1dHhvaWQpKSlcbiAgICByZXR1cm4gdGhpcy5tZXJnZSh1dHhvc2V0LCByZXN1bHRzKSBhcyB0aGlzXG4gIH1cblxuICAvKipcbiAgICogU2V0IHVuaW9uIGJldHdlZW4gdGhpcyBzZXQgYW5kIGEgcGFyYW1ldGVyLlxuICAgKlxuICAgKiBAcGFyYW0gdXR4b3NldCBUaGUgc2V0IHRvIHVuaW9uXG4gICAqXG4gICAqIEByZXR1cm5zIEEgbmV3IFN0YW5kYXJkVVRYT1NldCBjb250YWluaW5nIHRoZSB1bmlvblxuICAgKi9cbiAgdW5pb24gPSAodXR4b3NldDogdGhpcyk6IHRoaXMgPT4gdGhpcy5tZXJnZSh1dHhvc2V0KSBhcyB0aGlzXG5cbiAgLyoqXG4gICAqIE1lcmdlcyBhIHNldCBieSB0aGUgcnVsZSBwcm92aWRlZC5cbiAgICpcbiAgICogQHBhcmFtIHV0eG9zZXQgVGhlIHNldCB0byBtZXJnZSBieSB0aGUgTWVyZ2VSdWxlXG4gICAqIEBwYXJhbSBtZXJnZVJ1bGUgVGhlIFtbTWVyZ2VSdWxlXV0gdG8gYXBwbHlcbiAgICpcbiAgICogQHJldHVybnMgQSBuZXcgU3RhbmRhcmRVVFhPU2V0IGNvbnRhaW5pbmcgdGhlIG1lcmdlZCBkYXRhXG4gICAqXG4gICAqIEByZW1hcmtzXG4gICAqIFRoZSBtZXJnZSBydWxlcyBhcmUgYXMgZm9sbG93czpcbiAgICogICAqIFwiaW50ZXJzZWN0aW9uXCIgLSB0aGUgaW50ZXJzZWN0aW9uIG9mIHRoZSBzZXRcbiAgICogICAqIFwiZGlmZmVyZW5jZVNlbGZcIiAtIHRoZSBkaWZmZXJlbmNlIGJldHdlZW4gdGhlIGV4aXN0aW5nIGRhdGEgYW5kIG5ldyBzZXRcbiAgICogICAqIFwiZGlmZmVyZW5jZU5ld1wiIC0gdGhlIGRpZmZlcmVuY2UgYmV0d2VlbiB0aGUgbmV3IGRhdGEgYW5kIHRoZSBleGlzdGluZyBzZXRcbiAgICogICAqIFwic3ltRGlmZmVyZW5jZVwiIC0gdGhlIHVuaW9uIG9mIHRoZSBkaWZmZXJlbmNlcyBiZXR3ZWVuIGJvdGggc2V0cyBvZiBkYXRhXG4gICAqICAgKiBcInVuaW9uXCIgLSB0aGUgdW5pcXVlIHNldCBvZiBhbGwgZWxlbWVudHMgY29udGFpbmVkIGluIGJvdGggc2V0c1xuICAgKiAgICogXCJ1bmlvbk1pbnVzTmV3XCIgLSB0aGUgdW5pcXVlIHNldCBvZiBhbGwgZWxlbWVudHMgY29udGFpbmVkIGluIGJvdGggc2V0cywgZXhjbHVkaW5nIHZhbHVlcyBvbmx5IGZvdW5kIGluIHRoZSBuZXcgc2V0XG4gICAqICAgKiBcInVuaW9uTWludXNTZWxmXCIgLSB0aGUgdW5pcXVlIHNldCBvZiBhbGwgZWxlbWVudHMgY29udGFpbmVkIGluIGJvdGggc2V0cywgZXhjbHVkaW5nIHZhbHVlcyBvbmx5IGZvdW5kIGluIHRoZSBleGlzdGluZyBzZXRcbiAgICovXG4gIG1lcmdlQnlSdWxlID0gKHV0eG9zZXQ6IHRoaXMsIG1lcmdlUnVsZTogTWVyZ2VSdWxlKTogdGhpcyA9PiB7XG4gICAgbGV0IHVTZXQ6IHRoaXNcbiAgICBzd2l0Y2ggKG1lcmdlUnVsZSkge1xuICAgICAgY2FzZSBcImludGVyc2VjdGlvblwiOlxuICAgICAgICByZXR1cm4gdGhpcy5pbnRlcnNlY3Rpb24odXR4b3NldClcbiAgICAgIGNhc2UgXCJkaWZmZXJlbmNlU2VsZlwiOlxuICAgICAgICByZXR1cm4gdGhpcy5kaWZmZXJlbmNlKHV0eG9zZXQpXG4gICAgICBjYXNlIFwiZGlmZmVyZW5jZU5ld1wiOlxuICAgICAgICByZXR1cm4gdXR4b3NldC5kaWZmZXJlbmNlKHRoaXMpIGFzIHRoaXNcbiAgICAgIGNhc2UgXCJzeW1EaWZmZXJlbmNlXCI6XG4gICAgICAgIHJldHVybiB0aGlzLnN5bURpZmZlcmVuY2UodXR4b3NldClcbiAgICAgIGNhc2UgXCJ1bmlvblwiOlxuICAgICAgICByZXR1cm4gdGhpcy51bmlvbih1dHhvc2V0KVxuICAgICAgY2FzZSBcInVuaW9uTWludXNOZXdcIjpcbiAgICAgICAgdVNldCA9IHRoaXMudW5pb24odXR4b3NldClcbiAgICAgICAgcmV0dXJuIHVTZXQuZGlmZmVyZW5jZSh1dHhvc2V0KSBhcyB0aGlzXG4gICAgICBjYXNlIFwidW5pb25NaW51c1NlbGZcIjpcbiAgICAgICAgdVNldCA9IHRoaXMudW5pb24odXR4b3NldClcbiAgICAgICAgcmV0dXJuIHVTZXQuZGlmZmVyZW5jZSh0aGlzKSBhcyB0aGlzXG4gICAgICBkZWZhdWx0OlxuICAgICAgICB0aHJvdyBuZXcgTWVyZ2VSdWxlRXJyb3IoXG4gICAgICAgICAgXCJFcnJvciAtIFN0YW5kYXJkVVRYT1NldC5tZXJnZUJ5UnVsZTogYmFkIE1lcmdlUnVsZVwiXG4gICAgICAgIClcbiAgICB9XG4gIH1cbn1cbiJdfQ==