"use strict";
/**
 * @packageDocumentation
 * @module Utils-Payload
 */
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.MAGNETPayload = exports.ONIONPayload = exports.IPFSPayload = exports.URLPayload = exports.EMAILPayload = exports.YAMLPayload = exports.JSONPayload = exports.CSVPayload = exports.SVGPayload = exports.ICOPayload = exports.BMPPayload = exports.PNGPayload = exports.JPEGPayload = exports.SECPENCPayload = exports.SECPSIGPayload = exports.NODEIDPayload = exports.CHAINIDPayload = exports.SUBNETIDPayload = exports.NFTIDPayload = exports.UTXOIDPayload = exports.ASSETIDPayload = exports.TXIDPayload = exports.cb58EncodedPayload = exports.CCHAINADDRPayload = exports.PCHAINADDRPayload = exports.XCHAINADDRPayload = exports.ChainAddressPayload = exports.BIGNUMPayload = exports.B64STRPayload = exports.B58STRPayload = exports.HEXSTRPayload = exports.UTF8Payload = exports.BINPayload = exports.PayloadBase = exports.PayloadTypes = void 0;
const buffer_1 = require("buffer/");
const bintools_1 = __importDefault(require("./bintools"));
const bn_js_1 = __importDefault(require("bn.js"));
const errors_1 = require("../utils/errors");
const serialization_1 = require("../utils/serialization");
/**
 * @ignore
 */
const bintools = bintools_1.default.getInstance();
const serialization = serialization_1.Serialization.getInstance();
/**
 * Class for determining payload types and managing the lookup table.
 */
class PayloadTypes {
    constructor() {
        this.types = [];
        this.types = [
            "BIN",
            "UTF8",
            "HEXSTR",
            "B58STR",
            "B64STR",
            "BIGNUM",
            "XCHAINADDR",
            "PCHAINADDR",
            "CCHAINADDR",
            "TXID",
            "ASSETID",
            "UTXOID",
            "NFTID",
            "SUBNETID",
            "CHAINID",
            "NODEID",
            "SECPSIG",
            "SECPENC",
            "JPEG",
            "PNG",
            "BMP",
            "ICO",
            "SVG",
            "CSV",
            "JSON",
            "YAML",
            "EMAIL",
            "URL",
            "IPFS",
            "ONION",
            "MAGNET"
        ];
    }
    /**
     * Given an encoded payload buffer returns the payload content (minus typeID).
     */
    getContent(payload) {
        const pl = bintools.copyFrom(payload, 5);
        return pl;
    }
    /**
     * Given an encoded payload buffer returns the payload (with typeID).
     */
    getPayload(payload) {
        const pl = bintools.copyFrom(payload, 4);
        return pl;
    }
    /**
     * Given a payload buffer returns the proper TypeID.
     */
    getTypeID(payload) {
        const offset = 4;
        const typeID = bintools
            .copyFrom(payload, offset, offset + 1)
            .readUInt8(0);
        return typeID;
    }
    /**
     * Given a type string returns the proper TypeID.
     */
    lookupID(typestr) {
        return this.types.indexOf(typestr);
    }
    /**
     * Given a TypeID returns a string describing the payload type.
     */
    lookupType(value) {
        return this.types[`${value}`];
    }
    /**
     * Given a TypeID returns the proper [[PayloadBase]].
     */
    select(typeID, ...args) {
        switch (typeID) {
            case 0:
                return new BINPayload(...args);
            case 1:
                return new UTF8Payload(...args);
            case 2:
                return new HEXSTRPayload(...args);
            case 3:
                return new B58STRPayload(...args);
            case 4:
                return new B64STRPayload(...args);
            case 5:
                return new BIGNUMPayload(...args);
            case 6:
                return new XCHAINADDRPayload(...args);
            case 7:
                return new PCHAINADDRPayload(...args);
            case 8:
                return new CCHAINADDRPayload(...args);
            case 9:
                return new TXIDPayload(...args);
            case 10:
                return new ASSETIDPayload(...args);
            case 11:
                return new UTXOIDPayload(...args);
            case 12:
                return new NFTIDPayload(...args);
            case 13:
                return new SUBNETIDPayload(...args);
            case 14:
                return new CHAINIDPayload(...args);
            case 15:
                return new NODEIDPayload(...args);
            case 16:
                return new SECPSIGPayload(...args);
            case 17:
                return new SECPENCPayload(...args);
            case 18:
                return new JPEGPayload(...args);
            case 19:
                return new PNGPayload(...args);
            case 20:
                return new BMPPayload(...args);
            case 21:
                return new ICOPayload(...args);
            case 22:
                return new SVGPayload(...args);
            case 23:
                return new CSVPayload(...args);
            case 24:
                return new JSONPayload(...args);
            case 25:
                return new YAMLPayload(...args);
            case 26:
                return new EMAILPayload(...args);
            case 27:
                return new URLPayload(...args);
            case 28:
                return new IPFSPayload(...args);
            case 29:
                return new ONIONPayload(...args);
            case 30:
                return new MAGNETPayload(...args);
        }
        throw new errors_1.TypeIdError(`Error - PayloadTypes.select: unknown typeid ${typeID}`);
    }
    /**
     * Given a [[PayloadBase]] which may not be cast properly, returns a properly cast [[PayloadBase]].
     */
    recast(unknowPayload) {
        return this.select(unknowPayload.typeID(), unknowPayload.returnType());
    }
    /**
     * Returns the [[PayloadTypes]] singleton.
     */
    static getInstance() {
        if (!PayloadTypes.instance) {
            PayloadTypes.instance = new PayloadTypes();
        }
        return PayloadTypes.instance;
    }
}
exports.PayloadTypes = PayloadTypes;
/**
 * Base class for payloads.
 */
class PayloadBase {
    constructor() {
        this.payload = buffer_1.Buffer.alloc(0);
        this.typeid = undefined;
    }
    /**
     * Returns the TypeID for the payload.
     */
    typeID() {
        return this.typeid;
    }
    /**
     * Returns the string name for the payload's type.
     */
    typeName() {
        return PayloadTypes.getInstance().lookupType(this.typeid);
    }
    /**
     * Returns the payload content (minus typeID).
     */
    getContent() {
        const pl = bintools.copyFrom(this.payload);
        return pl;
    }
    /**
     * Returns the payload (with typeID).
     */
    getPayload() {
        const typeID = buffer_1.Buffer.alloc(1);
        typeID.writeUInt8(this.typeid, 0);
        const pl = buffer_1.Buffer.concat([typeID, bintools.copyFrom(this.payload)]);
        return pl;
    }
    /**
     * Decodes the payload as a {@link https://github.com/feross/buffer|Buffer} including 4 bytes for the length and TypeID.
     */
    fromBuffer(bytes, offset = 0) {
        const size = bintools
            .copyFrom(bytes, offset, offset + 4)
            .readUInt32BE(0);
        offset += 4;
        this.typeid = bintools.copyFrom(bytes, offset, offset + 1).readUInt8(0);
        offset += 1;
        this.payload = bintools.copyFrom(bytes, offset, offset + size - 1);
        offset += size - 1;
        return offset;
    }
    /**
     * Encodes the payload as a {@link https://github.com/feross/buffer|Buffer} including 4 bytes for the length and TypeID.
     */
    toBuffer() {
        const sizebuff = buffer_1.Buffer.alloc(4);
        sizebuff.writeUInt32BE(this.payload.length + 1, 0);
        const typebuff = buffer_1.Buffer.alloc(1);
        typebuff.writeUInt8(this.typeid, 0);
        return buffer_1.Buffer.concat([sizebuff, typebuff, this.payload]);
    }
}
exports.PayloadBase = PayloadBase;
/**
 * Class for payloads representing simple binary blobs.
 */
class BINPayload extends PayloadBase {
    /**
     * @param payload Buffer only
     */
    constructor(payload = undefined) {
        super();
        this.typeid = 0;
        if (payload instanceof buffer_1.Buffer) {
            this.payload = payload;
        }
        else if (typeof payload === "string") {
            this.payload = bintools.b58ToBuffer(payload);
        }
    }
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} for the payload.
     */
    returnType() {
        return this.payload;
    }
}
exports.BINPayload = BINPayload;
/**
 * Class for payloads representing UTF8 encoding.
 */
class UTF8Payload extends PayloadBase {
    /**
     * @param payload Buffer utf8 string
     */
    constructor(payload = undefined) {
        super();
        this.typeid = 1;
        if (payload instanceof buffer_1.Buffer) {
            this.payload = payload;
        }
        else if (typeof payload === "string") {
            this.payload = buffer_1.Buffer.from(payload, "utf8");
        }
    }
    /**
     * Returns a string for the payload.
     */
    returnType() {
        return this.payload.toString("utf8");
    }
}
exports.UTF8Payload = UTF8Payload;
/**
 * Class for payloads representing Hexadecimal encoding.
 */
class HEXSTRPayload extends PayloadBase {
    /**
     * @param payload Buffer or hex string
     */
    constructor(payload = undefined) {
        super();
        this.typeid = 2;
        if (payload instanceof buffer_1.Buffer) {
            this.payload = payload;
        }
        else if (typeof payload === "string") {
            if (payload.startsWith("0x") || !payload.match(/^[0-9A-Fa-f]+$/)) {
                throw new errors_1.HexError("HEXSTRPayload.constructor -- hex string may not start with 0x and must be in /^[0-9A-Fa-f]+$/: " +
                    payload);
            }
            this.payload = buffer_1.Buffer.from(payload, "hex");
        }
    }
    /**
     * Returns a hex string for the payload.
     */
    returnType() {
        return this.payload.toString("hex");
    }
}
exports.HEXSTRPayload = HEXSTRPayload;
/**
 * Class for payloads representing Base58 encoding.
 */
class B58STRPayload extends PayloadBase {
    /**
     * @param payload Buffer or cb58 encoded string
     */
    constructor(payload = undefined) {
        super();
        this.typeid = 3;
        if (payload instanceof buffer_1.Buffer) {
            this.payload = payload;
        }
        else if (typeof payload === "string") {
            this.payload = bintools.b58ToBuffer(payload);
        }
    }
    /**
     * Returns a base58 string for the payload.
     */
    returnType() {
        return bintools.bufferToB58(this.payload);
    }
}
exports.B58STRPayload = B58STRPayload;
/**
 * Class for payloads representing Base64 encoding.
 */
class B64STRPayload extends PayloadBase {
    /**
     * @param payload Buffer of base64 string
     */
    constructor(payload = undefined) {
        super();
        this.typeid = 4;
        if (payload instanceof buffer_1.Buffer) {
            this.payload = payload;
        }
        else if (typeof payload === "string") {
            this.payload = buffer_1.Buffer.from(payload, "base64");
        }
    }
    /**
     * Returns a base64 string for the payload.
     */
    returnType() {
        return this.payload.toString("base64");
    }
}
exports.B64STRPayload = B64STRPayload;
/**
 * Class for payloads representing Big Numbers.
 *
 * @param payload Accepts a Buffer, BN, or base64 string
 */
class BIGNUMPayload extends PayloadBase {
    /**
     * @param payload Buffer, BN, or base64 string
     */
    constructor(payload = undefined) {
        super();
        this.typeid = 5;
        if (payload instanceof buffer_1.Buffer) {
            this.payload = payload;
        }
        else if (payload instanceof bn_js_1.default) {
            this.payload = bintools.fromBNToBuffer(payload);
        }
        else if (typeof payload === "string") {
            this.payload = buffer_1.Buffer.from(payload, "hex");
        }
    }
    /**
     * Returns a {@link https://github.com/indutny/bn.js/|BN} for the payload.
     */
    returnType() {
        return bintools.fromBufferToBN(this.payload);
    }
}
exports.BIGNUMPayload = BIGNUMPayload;
/**
 * Class for payloads representing chain addresses.
 *
 */
class ChainAddressPayload extends PayloadBase {
    /**
     * @param payload Buffer or address string
     */
    constructor(payload = undefined, hrp) {
        super();
        this.typeid = 6;
        this.chainid = "";
        if (payload instanceof buffer_1.Buffer) {
            this.payload = payload;
        }
        else if (typeof payload === "string") {
            if (hrp != undefined) {
                this.payload = bintools.stringToAddress(payload, hrp);
            }
            else {
                this.payload = bintools.stringToAddress(payload);
            }
        }
    }
    /**
     * Returns the chainid.
     */
    returnChainID() {
        return this.chainid;
    }
    /**
     * Returns an address string for the payload.
     */
    returnType(hrp) {
        const type = "bech32";
        return serialization.bufferToType(this.payload, type, hrp, this.chainid);
    }
}
exports.ChainAddressPayload = ChainAddressPayload;
/**
 * Class for payloads representing X-Chin addresses.
 */
class XCHAINADDRPayload extends ChainAddressPayload {
    constructor() {
        super(...arguments);
        this.typeid = 6;
        this.chainid = "X";
    }
}
exports.XCHAINADDRPayload = XCHAINADDRPayload;
/**
 * Class for payloads representing P-Chain addresses.
 */
class PCHAINADDRPayload extends ChainAddressPayload {
    constructor() {
        super(...arguments);
        this.typeid = 7;
        this.chainid = "P";
    }
}
exports.PCHAINADDRPayload = PCHAINADDRPayload;
/**
 * Class for payloads representing C-Chain addresses.
 */
class CCHAINADDRPayload extends ChainAddressPayload {
    constructor() {
        super(...arguments);
        this.typeid = 8;
        this.chainid = "C";
    }
}
exports.CCHAINADDRPayload = CCHAINADDRPayload;
/**
 * Class for payloads representing data serialized by bintools.cb58Encode().
 */
class cb58EncodedPayload extends PayloadBase {
    /**
     * Returns a bintools.cb58Encoded string for the payload.
     */
    returnType() {
        return bintools.cb58Encode(this.payload);
    }
    /**
     * @param payload Buffer or cb58 encoded string
     */
    constructor(payload = undefined) {
        super();
        if (payload instanceof buffer_1.Buffer) {
            this.payload = payload;
        }
        else if (typeof payload === "string") {
            this.payload = bintools.cb58Decode(payload);
        }
    }
}
exports.cb58EncodedPayload = cb58EncodedPayload;
/**
 * Class for payloads representing TxIDs.
 */
class TXIDPayload extends cb58EncodedPayload {
    constructor() {
        super(...arguments);
        this.typeid = 9;
    }
}
exports.TXIDPayload = TXIDPayload;
/**
 * Class for payloads representing AssetIDs.
 */
class ASSETIDPayload extends cb58EncodedPayload {
    constructor() {
        super(...arguments);
        this.typeid = 10;
    }
}
exports.ASSETIDPayload = ASSETIDPayload;
/**
 * Class for payloads representing NODEIDs.
 */
class UTXOIDPayload extends cb58EncodedPayload {
    constructor() {
        super(...arguments);
        this.typeid = 11;
    }
}
exports.UTXOIDPayload = UTXOIDPayload;
/**
 * Class for payloads representing NFTIDs (UTXOIDs in an NFT context).
 */
class NFTIDPayload extends UTXOIDPayload {
    constructor() {
        super(...arguments);
        this.typeid = 12;
    }
}
exports.NFTIDPayload = NFTIDPayload;
/**
 * Class for payloads representing SubnetIDs.
 */
class SUBNETIDPayload extends cb58EncodedPayload {
    constructor() {
        super(...arguments);
        this.typeid = 13;
    }
}
exports.SUBNETIDPayload = SUBNETIDPayload;
/**
 * Class for payloads representing ChainIDs.
 */
class CHAINIDPayload extends cb58EncodedPayload {
    constructor() {
        super(...arguments);
        this.typeid = 14;
    }
}
exports.CHAINIDPayload = CHAINIDPayload;
/**
 * Class for payloads representing NodeIDs.
 */
class NODEIDPayload extends cb58EncodedPayload {
    constructor() {
        super(...arguments);
        this.typeid = 15;
    }
}
exports.NODEIDPayload = NODEIDPayload;
/**
 * Class for payloads representing secp256k1 signatures.
 * convention: secp256k1 signature (130 bytes)
 */
class SECPSIGPayload extends B58STRPayload {
    constructor() {
        super(...arguments);
        this.typeid = 16;
    }
}
exports.SECPSIGPayload = SECPSIGPayload;
/**
 * Class for payloads representing secp256k1 encrypted messages.
 * convention: public key (65 bytes) + secp256k1 encrypted message for that public key
 */
class SECPENCPayload extends B58STRPayload {
    constructor() {
        super(...arguments);
        this.typeid = 17;
    }
}
exports.SECPENCPayload = SECPENCPayload;
/**
 * Class for payloads representing JPEG images.
 */
class JPEGPayload extends BINPayload {
    constructor() {
        super(...arguments);
        this.typeid = 18;
    }
}
exports.JPEGPayload = JPEGPayload;
class PNGPayload extends BINPayload {
    constructor() {
        super(...arguments);
        this.typeid = 19;
    }
}
exports.PNGPayload = PNGPayload;
/**
 * Class for payloads representing BMP images.
 */
class BMPPayload extends BINPayload {
    constructor() {
        super(...arguments);
        this.typeid = 20;
    }
}
exports.BMPPayload = BMPPayload;
/**
 * Class for payloads representing ICO images.
 */
class ICOPayload extends BINPayload {
    constructor() {
        super(...arguments);
        this.typeid = 21;
    }
}
exports.ICOPayload = ICOPayload;
/**
 * Class for payloads representing SVG images.
 */
class SVGPayload extends UTF8Payload {
    constructor() {
        super(...arguments);
        this.typeid = 22;
    }
}
exports.SVGPayload = SVGPayload;
/**
 * Class for payloads representing CSV files.
 */
class CSVPayload extends UTF8Payload {
    constructor() {
        super(...arguments);
        this.typeid = 23;
    }
}
exports.CSVPayload = CSVPayload;
/**
 * Class for payloads representing JSON strings.
 */
class JSONPayload extends PayloadBase {
    constructor(payload = undefined) {
        super();
        this.typeid = 24;
        if (payload instanceof buffer_1.Buffer) {
            this.payload = payload;
        }
        else if (typeof payload === "string") {
            this.payload = buffer_1.Buffer.from(payload, "utf8");
        }
        else if (payload) {
            let jsonstr = JSON.stringify(payload);
            this.payload = buffer_1.Buffer.from(jsonstr, "utf8");
        }
    }
    /**
     * Returns a JSON-decoded object for the payload.
     */
    returnType() {
        return JSON.parse(this.payload.toString("utf8"));
    }
}
exports.JSONPayload = JSONPayload;
/**
 * Class for payloads representing YAML definitions.
 */
class YAMLPayload extends UTF8Payload {
    constructor() {
        super(...arguments);
        this.typeid = 25;
    }
}
exports.YAMLPayload = YAMLPayload;
/**
 * Class for payloads representing email addresses.
 */
class EMAILPayload extends UTF8Payload {
    constructor() {
        super(...arguments);
        this.typeid = 26;
    }
}
exports.EMAILPayload = EMAILPayload;
/**
 * Class for payloads representing URL strings.
 */
class URLPayload extends UTF8Payload {
    constructor() {
        super(...arguments);
        this.typeid = 27;
    }
}
exports.URLPayload = URLPayload;
/**
 * Class for payloads representing IPFS addresses.
 */
class IPFSPayload extends B58STRPayload {
    constructor() {
        super(...arguments);
        this.typeid = 28;
    }
}
exports.IPFSPayload = IPFSPayload;
/**
 * Class for payloads representing onion URLs.
 */
class ONIONPayload extends UTF8Payload {
    constructor() {
        super(...arguments);
        this.typeid = 29;
    }
}
exports.ONIONPayload = ONIONPayload;
/**
 * Class for payloads representing torrent magnet links.
 */
class MAGNETPayload extends UTF8Payload {
    constructor() {
        super(...arguments);
        this.typeid = 30;
    }
}
exports.MAGNETPayload = MAGNETPayload;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicGF5bG9hZC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3NyYy91dGlscy9wYXlsb2FkLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7QUFBQTs7O0dBR0c7Ozs7OztBQUVILG9DQUFnQztBQUNoQywwREFBaUM7QUFDakMsa0RBQXNCO0FBQ3RCLDRDQUF1RDtBQUN2RCwwREFBc0U7QUFFdEU7O0dBRUc7QUFDSCxNQUFNLFFBQVEsR0FBYSxrQkFBUSxDQUFDLFdBQVcsRUFBRSxDQUFBO0FBQ2pELE1BQU0sYUFBYSxHQUFrQiw2QkFBYSxDQUFDLFdBQVcsRUFBRSxDQUFBO0FBRWhFOztHQUVHO0FBQ0gsTUFBYSxZQUFZO0lBd0l2QjtRQXRJVSxVQUFLLEdBQWEsRUFBRSxDQUFBO1FBdUk1QixJQUFJLENBQUMsS0FBSyxHQUFHO1lBQ1gsS0FBSztZQUNMLE1BQU07WUFDTixRQUFRO1lBQ1IsUUFBUTtZQUNSLFFBQVE7WUFDUixRQUFRO1lBQ1IsWUFBWTtZQUNaLFlBQVk7WUFDWixZQUFZO1lBQ1osTUFBTTtZQUNOLFNBQVM7WUFDVCxRQUFRO1lBQ1IsT0FBTztZQUNQLFVBQVU7WUFDVixTQUFTO1lBQ1QsUUFBUTtZQUNSLFNBQVM7WUFDVCxTQUFTO1lBQ1QsTUFBTTtZQUNOLEtBQUs7WUFDTCxLQUFLO1lBQ0wsS0FBSztZQUNMLEtBQUs7WUFDTCxLQUFLO1lBQ0wsTUFBTTtZQUNOLE1BQU07WUFDTixPQUFPO1lBQ1AsS0FBSztZQUNMLE1BQU07WUFDTixPQUFPO1lBQ1AsUUFBUTtTQUNULENBQUE7SUFDSCxDQUFDO0lBdEtEOztPQUVHO0lBQ0gsVUFBVSxDQUFDLE9BQWU7UUFDeEIsTUFBTSxFQUFFLEdBQVcsUUFBUSxDQUFDLFFBQVEsQ0FBQyxPQUFPLEVBQUUsQ0FBQyxDQUFDLENBQUE7UUFDaEQsT0FBTyxFQUFFLENBQUE7SUFDWCxDQUFDO0lBRUQ7O09BRUc7SUFDSCxVQUFVLENBQUMsT0FBZTtRQUN4QixNQUFNLEVBQUUsR0FBVyxRQUFRLENBQUMsUUFBUSxDQUFDLE9BQU8sRUFBRSxDQUFDLENBQUMsQ0FBQTtRQUNoRCxPQUFPLEVBQUUsQ0FBQTtJQUNYLENBQUM7SUFFRDs7T0FFRztJQUNILFNBQVMsQ0FBQyxPQUFlO1FBQ3ZCLE1BQU0sTUFBTSxHQUFXLENBQUMsQ0FBQTtRQUN4QixNQUFNLE1BQU0sR0FBVyxRQUFRO2FBQzVCLFFBQVEsQ0FBQyxPQUFPLEVBQUUsTUFBTSxFQUFFLE1BQU0sR0FBRyxDQUFDLENBQUM7YUFDckMsU0FBUyxDQUFDLENBQUMsQ0FBQyxDQUFBO1FBQ2YsT0FBTyxNQUFNLENBQUE7SUFDZixDQUFDO0lBRUQ7O09BRUc7SUFDSCxRQUFRLENBQUMsT0FBZTtRQUN0QixPQUFPLElBQUksQ0FBQyxLQUFLLENBQUMsT0FBTyxDQUFDLE9BQU8sQ0FBQyxDQUFBO0lBQ3BDLENBQUM7SUFFRDs7T0FFRztJQUNILFVBQVUsQ0FBQyxLQUFhO1FBQ3RCLE9BQU8sSUFBSSxDQUFDLEtBQUssQ0FBQyxHQUFHLEtBQUssRUFBRSxDQUFDLENBQUE7SUFDL0IsQ0FBQztJQUVEOztPQUVHO0lBQ0gsTUFBTSxDQUFDLE1BQWMsRUFBRSxHQUFHLElBQVc7UUFDbkMsUUFBUSxNQUFNLEVBQUU7WUFDZCxLQUFLLENBQUM7Z0JBQ0osT0FBTyxJQUFJLFVBQVUsQ0FBQyxHQUFHLElBQUksQ0FBQyxDQUFBO1lBQ2hDLEtBQUssQ0FBQztnQkFDSixPQUFPLElBQUksV0FBVyxDQUFDLEdBQUcsSUFBSSxDQUFDLENBQUE7WUFDakMsS0FBSyxDQUFDO2dCQUNKLE9BQU8sSUFBSSxhQUFhLENBQUMsR0FBRyxJQUFJLENBQUMsQ0FBQTtZQUNuQyxLQUFLLENBQUM7Z0JBQ0osT0FBTyxJQUFJLGFBQWEsQ0FBQyxHQUFHLElBQUksQ0FBQyxDQUFBO1lBQ25DLEtBQUssQ0FBQztnQkFDSixPQUFPLElBQUksYUFBYSxDQUFDLEdBQUcsSUFBSSxDQUFDLENBQUE7WUFDbkMsS0FBSyxDQUFDO2dCQUNKLE9BQU8sSUFBSSxhQUFhLENBQUMsR0FBRyxJQUFJLENBQUMsQ0FBQTtZQUNuQyxLQUFLLENBQUM7Z0JBQ0osT0FBTyxJQUFJLGlCQUFpQixDQUFDLEdBQUcsSUFBSSxDQUFDLENBQUE7WUFDdkMsS0FBSyxDQUFDO2dCQUNKLE9BQU8sSUFBSSxpQkFBaUIsQ0FBQyxHQUFHLElBQUksQ0FBQyxDQUFBO1lBQ3ZDLEtBQUssQ0FBQztnQkFDSixPQUFPLElBQUksaUJBQWlCLENBQUMsR0FBRyxJQUFJLENBQUMsQ0FBQTtZQUN2QyxLQUFLLENBQUM7Z0JBQ0osT0FBTyxJQUFJLFdBQVcsQ0FBQyxHQUFHLElBQUksQ0FBQyxDQUFBO1lBQ2pDLEtBQUssRUFBRTtnQkFDTCxPQUFPLElBQUksY0FBYyxDQUFDLEdBQUcsSUFBSSxDQUFDLENBQUE7WUFDcEMsS0FBSyxFQUFFO2dCQUNMLE9BQU8sSUFBSSxhQUFhLENBQUMsR0FBRyxJQUFJLENBQUMsQ0FBQTtZQUNuQyxLQUFLLEVBQUU7Z0JBQ0wsT0FBTyxJQUFJLFlBQVksQ0FBQyxHQUFHLElBQUksQ0FBQyxDQUFBO1lBQ2xDLEtBQUssRUFBRTtnQkFDTCxPQUFPLElBQUksZUFBZSxDQUFDLEdBQUcsSUFBSSxDQUFDLENBQUE7WUFDckMsS0FBSyxFQUFFO2dCQUNMLE9BQU8sSUFBSSxjQUFjLENBQUMsR0FBRyxJQUFJLENBQUMsQ0FBQTtZQUNwQyxLQUFLLEVBQUU7Z0JBQ0wsT0FBTyxJQUFJLGFBQWEsQ0FBQyxHQUFHLElBQUksQ0FBQyxDQUFBO1lBQ25DLEtBQUssRUFBRTtnQkFDTCxPQUFPLElBQUksY0FBYyxDQUFDLEdBQUcsSUFBSSxDQUFDLENBQUE7WUFDcEMsS0FBSyxFQUFFO2dCQUNMLE9BQU8sSUFBSSxjQUFjLENBQUMsR0FBRyxJQUFJLENBQUMsQ0FBQTtZQUNwQyxLQUFLLEVBQUU7Z0JBQ0wsT0FBTyxJQUFJLFdBQVcsQ0FBQyxHQUFHLElBQUksQ0FBQyxDQUFBO1lBQ2pDLEtBQUssRUFBRTtnQkFDTCxPQUFPLElBQUksVUFBVSxDQUFDLEdBQUcsSUFBSSxDQUFDLENBQUE7WUFDaEMsS0FBSyxFQUFFO2dCQUNMLE9BQU8sSUFBSSxVQUFVLENBQUMsR0FBRyxJQUFJLENBQUMsQ0FBQTtZQUNoQyxLQUFLLEVBQUU7Z0JBQ0wsT0FBTyxJQUFJLFVBQVUsQ0FBQyxHQUFHLElBQUksQ0FBQyxDQUFBO1lBQ2hDLEtBQUssRUFBRTtnQkFDTCxPQUFPLElBQUksVUFBVSxDQUFDLEdBQUcsSUFBSSxDQUFDLENBQUE7WUFDaEMsS0FBSyxFQUFFO2dCQUNMLE9BQU8sSUFBSSxVQUFVLENBQUMsR0FBRyxJQUFJLENBQUMsQ0FBQTtZQUNoQyxLQUFLLEVBQUU7Z0JBQ0wsT0FBTyxJQUFJLFdBQVcsQ0FBQyxHQUFHLElBQUksQ0FBQyxDQUFBO1lBQ2pDLEtBQUssRUFBRTtnQkFDTCxPQUFPLElBQUksV0FBVyxDQUFDLEdBQUcsSUFBSSxDQUFDLENBQUE7WUFDakMsS0FBSyxFQUFFO2dCQUNMLE9BQU8sSUFBSSxZQUFZLENBQUMsR0FBRyxJQUFJLENBQUMsQ0FBQTtZQUNsQyxLQUFLLEVBQUU7Z0JBQ0wsT0FBTyxJQUFJLFVBQVUsQ0FBQyxHQUFHLElBQUksQ0FBQyxDQUFBO1lBQ2hDLEtBQUssRUFBRTtnQkFDTCxPQUFPLElBQUksV0FBVyxDQUFDLEdBQUcsSUFBSSxDQUFDLENBQUE7WUFDakMsS0FBSyxFQUFFO2dCQUNMLE9BQU8sSUFBSSxZQUFZLENBQUMsR0FBRyxJQUFJLENBQUMsQ0FBQTtZQUNsQyxLQUFLLEVBQUU7Z0JBQ0wsT0FBTyxJQUFJLGFBQWEsQ0FBQyxHQUFHLElBQUksQ0FBQyxDQUFBO1NBQ3BDO1FBQ0QsTUFBTSxJQUFJLG9CQUFXLENBQ25CLCtDQUErQyxNQUFNLEVBQUUsQ0FDeEQsQ0FBQTtJQUNILENBQUM7SUFFRDs7T0FFRztJQUNILE1BQU0sQ0FBQyxhQUEwQjtRQUMvQixPQUFPLElBQUksQ0FBQyxNQUFNLENBQUMsYUFBYSxDQUFDLE1BQU0sRUFBRSxFQUFFLGFBQWEsQ0FBQyxVQUFVLEVBQUUsQ0FBQyxDQUFBO0lBQ3hFLENBQUM7SUFFRDs7T0FFRztJQUNILE1BQU0sQ0FBQyxXQUFXO1FBQ2hCLElBQUksQ0FBQyxZQUFZLENBQUMsUUFBUSxFQUFFO1lBQzFCLFlBQVksQ0FBQyxRQUFRLEdBQUcsSUFBSSxZQUFZLEVBQUUsQ0FBQTtTQUMzQztRQUVELE9BQU8sWUFBWSxDQUFDLFFBQVEsQ0FBQTtJQUM5QixDQUFDO0NBcUNGO0FBM0tELG9DQTJLQztBQUVEOztHQUVHO0FBQ0gsTUFBc0IsV0FBVztJQW1FL0I7UUFsRVUsWUFBTyxHQUFXLGVBQU0sQ0FBQyxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUE7UUFDakMsV0FBTSxHQUFXLFNBQVMsQ0FBQTtJQWlFckIsQ0FBQztJQS9EaEI7O09BRUc7SUFDSCxNQUFNO1FBQ0osT0FBTyxJQUFJLENBQUMsTUFBTSxDQUFBO0lBQ3BCLENBQUM7SUFFRDs7T0FFRztJQUNILFFBQVE7UUFDTixPQUFPLFlBQVksQ0FBQyxXQUFXLEVBQUUsQ0FBQyxVQUFVLENBQUMsSUFBSSxDQUFDLE1BQU0sQ0FBQyxDQUFBO0lBQzNELENBQUM7SUFFRDs7T0FFRztJQUNILFVBQVU7UUFDUixNQUFNLEVBQUUsR0FBVyxRQUFRLENBQUMsUUFBUSxDQUFDLElBQUksQ0FBQyxPQUFPLENBQUMsQ0FBQTtRQUNsRCxPQUFPLEVBQUUsQ0FBQTtJQUNYLENBQUM7SUFFRDs7T0FFRztJQUNILFVBQVU7UUFDUixNQUFNLE1BQU0sR0FBVyxlQUFNLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxDQUFBO1FBQ3RDLE1BQU0sQ0FBQyxVQUFVLENBQUMsSUFBSSxDQUFDLE1BQU0sRUFBRSxDQUFDLENBQUMsQ0FBQTtRQUNqQyxNQUFNLEVBQUUsR0FBVyxlQUFNLENBQUMsTUFBTSxDQUFDLENBQUMsTUFBTSxFQUFFLFFBQVEsQ0FBQyxRQUFRLENBQUMsSUFBSSxDQUFDLE9BQU8sQ0FBQyxDQUFDLENBQUMsQ0FBQTtRQUMzRSxPQUFPLEVBQUUsQ0FBQTtJQUNYLENBQUM7SUFFRDs7T0FFRztJQUNILFVBQVUsQ0FBQyxLQUFhLEVBQUUsU0FBaUIsQ0FBQztRQUMxQyxNQUFNLElBQUksR0FBVyxRQUFRO2FBQzFCLFFBQVEsQ0FBQyxLQUFLLEVBQUUsTUFBTSxFQUFFLE1BQU0sR0FBRyxDQUFDLENBQUM7YUFDbkMsWUFBWSxDQUFDLENBQUMsQ0FBQyxDQUFBO1FBQ2xCLE1BQU0sSUFBSSxDQUFDLENBQUE7UUFDWCxJQUFJLENBQUMsTUFBTSxHQUFHLFFBQVEsQ0FBQyxRQUFRLENBQUMsS0FBSyxFQUFFLE1BQU0sRUFBRSxNQUFNLEdBQUcsQ0FBQyxDQUFDLENBQUMsU0FBUyxDQUFDLENBQUMsQ0FBQyxDQUFBO1FBQ3ZFLE1BQU0sSUFBSSxDQUFDLENBQUE7UUFDWCxJQUFJLENBQUMsT0FBTyxHQUFHLFFBQVEsQ0FBQyxRQUFRLENBQUMsS0FBSyxFQUFFLE1BQU0sRUFBRSxNQUFNLEdBQUcsSUFBSSxHQUFHLENBQUMsQ0FBQyxDQUFBO1FBQ2xFLE1BQU0sSUFBSSxJQUFJLEdBQUcsQ0FBQyxDQUFBO1FBQ2xCLE9BQU8sTUFBTSxDQUFBO0lBQ2YsQ0FBQztJQUVEOztPQUVHO0lBQ0gsUUFBUTtRQUNOLE1BQU0sUUFBUSxHQUFXLGVBQU0sQ0FBQyxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUE7UUFDeEMsUUFBUSxDQUFDLGFBQWEsQ0FBQyxJQUFJLENBQUMsT0FBTyxDQUFDLE1BQU0sR0FBRyxDQUFDLEVBQUUsQ0FBQyxDQUFDLENBQUE7UUFDbEQsTUFBTSxRQUFRLEdBQVcsZUFBTSxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQTtRQUN4QyxRQUFRLENBQUMsVUFBVSxDQUFDLElBQUksQ0FBQyxNQUFNLEVBQUUsQ0FBQyxDQUFDLENBQUE7UUFDbkMsT0FBTyxlQUFNLENBQUMsTUFBTSxDQUFDLENBQUMsUUFBUSxFQUFFLFFBQVEsRUFBRSxJQUFJLENBQUMsT0FBTyxDQUFDLENBQUMsQ0FBQTtJQUMxRCxDQUFDO0NBUUY7QUFwRUQsa0NBb0VDO0FBRUQ7O0dBRUc7QUFDSCxNQUFhLFVBQVcsU0FBUSxXQUFXO0lBU3pDOztPQUVHO0lBQ0gsWUFBWSxVQUFlLFNBQVM7UUFDbEMsS0FBSyxFQUFFLENBQUE7UUFaQyxXQUFNLEdBQUcsQ0FBQyxDQUFBO1FBYWxCLElBQUksT0FBTyxZQUFZLGVBQU0sRUFBRTtZQUM3QixJQUFJLENBQUMsT0FBTyxHQUFHLE9BQU8sQ0FBQTtTQUN2QjthQUFNLElBQUksT0FBTyxPQUFPLEtBQUssUUFBUSxFQUFFO1lBQ3RDLElBQUksQ0FBQyxPQUFPLEdBQUcsUUFBUSxDQUFDLFdBQVcsQ0FBQyxPQUFPLENBQUMsQ0FBQTtTQUM3QztJQUNILENBQUM7SUFoQkQ7O09BRUc7SUFDSCxVQUFVO1FBQ1IsT0FBTyxJQUFJLENBQUMsT0FBTyxDQUFBO0lBQ3JCLENBQUM7Q0FZRjtBQXBCRCxnQ0FvQkM7QUFFRDs7R0FFRztBQUNILE1BQWEsV0FBWSxTQUFRLFdBQVc7SUFTMUM7O09BRUc7SUFDSCxZQUFZLFVBQWUsU0FBUztRQUNsQyxLQUFLLEVBQUUsQ0FBQTtRQVpDLFdBQU0sR0FBRyxDQUFDLENBQUE7UUFhbEIsSUFBSSxPQUFPLFlBQVksZUFBTSxFQUFFO1lBQzdCLElBQUksQ0FBQyxPQUFPLEdBQUcsT0FBTyxDQUFBO1NBQ3ZCO2FBQU0sSUFBSSxPQUFPLE9BQU8sS0FBSyxRQUFRLEVBQUU7WUFDdEMsSUFBSSxDQUFDLE9BQU8sR0FBRyxlQUFNLENBQUMsSUFBSSxDQUFDLE9BQU8sRUFBRSxNQUFNLENBQUMsQ0FBQTtTQUM1QztJQUNILENBQUM7SUFoQkQ7O09BRUc7SUFDSCxVQUFVO1FBQ1IsT0FBTyxJQUFJLENBQUMsT0FBTyxDQUFDLFFBQVEsQ0FBQyxNQUFNLENBQUMsQ0FBQTtJQUN0QyxDQUFDO0NBWUY7QUFwQkQsa0NBb0JDO0FBRUQ7O0dBRUc7QUFDSCxNQUFhLGFBQWMsU0FBUSxXQUFXO0lBUzVDOztPQUVHO0lBQ0gsWUFBWSxVQUFlLFNBQVM7UUFDbEMsS0FBSyxFQUFFLENBQUE7UUFaQyxXQUFNLEdBQUcsQ0FBQyxDQUFBO1FBYWxCLElBQUksT0FBTyxZQUFZLGVBQU0sRUFBRTtZQUM3QixJQUFJLENBQUMsT0FBTyxHQUFHLE9BQU8sQ0FBQTtTQUN2QjthQUFNLElBQUksT0FBTyxPQUFPLEtBQUssUUFBUSxFQUFFO1lBQ3RDLElBQUksT0FBTyxDQUFDLFVBQVUsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLE9BQU8sQ0FBQyxLQUFLLENBQUMsZ0JBQWdCLENBQUMsRUFBRTtnQkFDaEUsTUFBTSxJQUFJLGlCQUFRLENBQ2hCLGlHQUFpRztvQkFDL0YsT0FBTyxDQUNWLENBQUE7YUFDRjtZQUNELElBQUksQ0FBQyxPQUFPLEdBQUcsZUFBTSxDQUFDLElBQUksQ0FBQyxPQUFPLEVBQUUsS0FBSyxDQUFDLENBQUE7U0FDM0M7SUFDSCxDQUFDO0lBdEJEOztPQUVHO0lBQ0gsVUFBVTtRQUNSLE9BQU8sSUFBSSxDQUFDLE9BQU8sQ0FBQyxRQUFRLENBQUMsS0FBSyxDQUFDLENBQUE7SUFDckMsQ0FBQztDQWtCRjtBQTFCRCxzQ0EwQkM7QUFFRDs7R0FFRztBQUNILE1BQWEsYUFBYyxTQUFRLFdBQVc7SUFTNUM7O09BRUc7SUFDSCxZQUFZLFVBQWUsU0FBUztRQUNsQyxLQUFLLEVBQUUsQ0FBQTtRQVpDLFdBQU0sR0FBRyxDQUFDLENBQUE7UUFhbEIsSUFBSSxPQUFPLFlBQVksZUFBTSxFQUFFO1lBQzdCLElBQUksQ0FBQyxPQUFPLEdBQUcsT0FBTyxDQUFBO1NBQ3ZCO2FBQU0sSUFBSSxPQUFPLE9BQU8sS0FBSyxRQUFRLEVBQUU7WUFDdEMsSUFBSSxDQUFDLE9BQU8sR0FBRyxRQUFRLENBQUMsV0FBVyxDQUFDLE9BQU8sQ0FBQyxDQUFBO1NBQzdDO0lBQ0gsQ0FBQztJQWhCRDs7T0FFRztJQUNILFVBQVU7UUFDUixPQUFPLFFBQVEsQ0FBQyxXQUFXLENBQUMsSUFBSSxDQUFDLE9BQU8sQ0FBQyxDQUFBO0lBQzNDLENBQUM7Q0FZRjtBQXBCRCxzQ0FvQkM7QUFFRDs7R0FFRztBQUNILE1BQWEsYUFBYyxTQUFRLFdBQVc7SUFTNUM7O09BRUc7SUFDSCxZQUFZLFVBQWUsU0FBUztRQUNsQyxLQUFLLEVBQUUsQ0FBQTtRQVpDLFdBQU0sR0FBRyxDQUFDLENBQUE7UUFhbEIsSUFBSSxPQUFPLFlBQVksZUFBTSxFQUFFO1lBQzdCLElBQUksQ0FBQyxPQUFPLEdBQUcsT0FBTyxDQUFBO1NBQ3ZCO2FBQU0sSUFBSSxPQUFPLE9BQU8sS0FBSyxRQUFRLEVBQUU7WUFDdEMsSUFBSSxDQUFDLE9BQU8sR0FBRyxlQUFNLENBQUMsSUFBSSxDQUFDLE9BQU8sRUFBRSxRQUFRLENBQUMsQ0FBQTtTQUM5QztJQUNILENBQUM7SUFoQkQ7O09BRUc7SUFDSCxVQUFVO1FBQ1IsT0FBTyxJQUFJLENBQUMsT0FBTyxDQUFDLFFBQVEsQ0FBQyxRQUFRLENBQUMsQ0FBQTtJQUN4QyxDQUFDO0NBWUY7QUFwQkQsc0NBb0JDO0FBRUQ7Ozs7R0FJRztBQUNILE1BQWEsYUFBYyxTQUFRLFdBQVc7SUFTNUM7O09BRUc7SUFDSCxZQUFZLFVBQWUsU0FBUztRQUNsQyxLQUFLLEVBQUUsQ0FBQTtRQVpDLFdBQU0sR0FBRyxDQUFDLENBQUE7UUFhbEIsSUFBSSxPQUFPLFlBQVksZUFBTSxFQUFFO1lBQzdCLElBQUksQ0FBQyxPQUFPLEdBQUcsT0FBTyxDQUFBO1NBQ3ZCO2FBQU0sSUFBSSxPQUFPLFlBQVksZUFBRSxFQUFFO1lBQ2hDLElBQUksQ0FBQyxPQUFPLEdBQUcsUUFBUSxDQUFDLGNBQWMsQ0FBQyxPQUFPLENBQUMsQ0FBQTtTQUNoRDthQUFNLElBQUksT0FBTyxPQUFPLEtBQUssUUFBUSxFQUFFO1lBQ3RDLElBQUksQ0FBQyxPQUFPLEdBQUcsZUFBTSxDQUFDLElBQUksQ0FBQyxPQUFPLEVBQUUsS0FBSyxDQUFDLENBQUE7U0FDM0M7SUFDSCxDQUFDO0lBbEJEOztPQUVHO0lBQ0gsVUFBVTtRQUNSLE9BQU8sUUFBUSxDQUFDLGNBQWMsQ0FBQyxJQUFJLENBQUMsT0FBTyxDQUFDLENBQUE7SUFDOUMsQ0FBQztDQWNGO0FBdEJELHNDQXNCQztBQUVEOzs7R0FHRztBQUNILE1BQXNCLG1CQUFvQixTQUFRLFdBQVc7SUFrQjNEOztPQUVHO0lBQ0gsWUFBWSxVQUFlLFNBQVMsRUFBRSxHQUFZO1FBQ2hELEtBQUssRUFBRSxDQUFBO1FBckJDLFdBQU0sR0FBRyxDQUFDLENBQUE7UUFDVixZQUFPLEdBQVcsRUFBRSxDQUFBO1FBcUI1QixJQUFJLE9BQU8sWUFBWSxlQUFNLEVBQUU7WUFDN0IsSUFBSSxDQUFDLE9BQU8sR0FBRyxPQUFPLENBQUE7U0FDdkI7YUFBTSxJQUFJLE9BQU8sT0FBTyxLQUFLLFFBQVEsRUFBRTtZQUN0QyxJQUFJLEdBQUcsSUFBSSxTQUFTLEVBQUU7Z0JBQ3BCLElBQUksQ0FBQyxPQUFPLEdBQUcsUUFBUSxDQUFDLGVBQWUsQ0FBQyxPQUFPLEVBQUUsR0FBRyxDQUFDLENBQUE7YUFDdEQ7aUJBQU07Z0JBQ0wsSUFBSSxDQUFDLE9BQU8sR0FBRyxRQUFRLENBQUMsZUFBZSxDQUFDLE9BQU8sQ0FBQyxDQUFBO2FBQ2pEO1NBQ0Y7SUFDSCxDQUFDO0lBNUJEOztPQUVHO0lBQ0gsYUFBYTtRQUNYLE9BQU8sSUFBSSxDQUFDLE9BQU8sQ0FBQTtJQUNyQixDQUFDO0lBRUQ7O09BRUc7SUFDSCxVQUFVLENBQUMsR0FBVztRQUNwQixNQUFNLElBQUksR0FBbUIsUUFBUSxDQUFBO1FBQ3JDLE9BQU8sYUFBYSxDQUFDLFlBQVksQ0FBQyxJQUFJLENBQUMsT0FBTyxFQUFFLElBQUksRUFBRSxHQUFHLEVBQUUsSUFBSSxDQUFDLE9BQU8sQ0FBQyxDQUFBO0lBQzFFLENBQUM7Q0FnQkY7QUFqQ0Qsa0RBaUNDO0FBRUQ7O0dBRUc7QUFDSCxNQUFhLGlCQUFrQixTQUFRLG1CQUFtQjtJQUExRDs7UUFDWSxXQUFNLEdBQUcsQ0FBQyxDQUFBO1FBQ1YsWUFBTyxHQUFHLEdBQUcsQ0FBQTtJQUN6QixDQUFDO0NBQUE7QUFIRCw4Q0FHQztBQUVEOztHQUVHO0FBQ0gsTUFBYSxpQkFBa0IsU0FBUSxtQkFBbUI7SUFBMUQ7O1FBQ1ksV0FBTSxHQUFHLENBQUMsQ0FBQTtRQUNWLFlBQU8sR0FBRyxHQUFHLENBQUE7SUFDekIsQ0FBQztDQUFBO0FBSEQsOENBR0M7QUFFRDs7R0FFRztBQUNILE1BQWEsaUJBQWtCLFNBQVEsbUJBQW1CO0lBQTFEOztRQUNZLFdBQU0sR0FBRyxDQUFDLENBQUE7UUFDVixZQUFPLEdBQUcsR0FBRyxDQUFBO0lBQ3pCLENBQUM7Q0FBQTtBQUhELDhDQUdDO0FBRUQ7O0dBRUc7QUFDSCxNQUFzQixrQkFBbUIsU0FBUSxXQUFXO0lBQzFEOztPQUVHO0lBQ0gsVUFBVTtRQUNSLE9BQU8sUUFBUSxDQUFDLFVBQVUsQ0FBQyxJQUFJLENBQUMsT0FBTyxDQUFDLENBQUE7SUFDMUMsQ0FBQztJQUNEOztPQUVHO0lBQ0gsWUFBWSxVQUFlLFNBQVM7UUFDbEMsS0FBSyxFQUFFLENBQUE7UUFDUCxJQUFJLE9BQU8sWUFBWSxlQUFNLEVBQUU7WUFDN0IsSUFBSSxDQUFDLE9BQU8sR0FBRyxPQUFPLENBQUE7U0FDdkI7YUFBTSxJQUFJLE9BQU8sT0FBTyxLQUFLLFFBQVEsRUFBRTtZQUN0QyxJQUFJLENBQUMsT0FBTyxHQUFHLFFBQVEsQ0FBQyxVQUFVLENBQUMsT0FBTyxDQUFDLENBQUE7U0FDNUM7SUFDSCxDQUFDO0NBQ0Y7QUFsQkQsZ0RBa0JDO0FBRUQ7O0dBRUc7QUFDSCxNQUFhLFdBQVksU0FBUSxrQkFBa0I7SUFBbkQ7O1FBQ1ksV0FBTSxHQUFHLENBQUMsQ0FBQTtJQUN0QixDQUFDO0NBQUE7QUFGRCxrQ0FFQztBQUVEOztHQUVHO0FBQ0gsTUFBYSxjQUFlLFNBQVEsa0JBQWtCO0lBQXREOztRQUNZLFdBQU0sR0FBRyxFQUFFLENBQUE7SUFDdkIsQ0FBQztDQUFBO0FBRkQsd0NBRUM7QUFFRDs7R0FFRztBQUNILE1BQWEsYUFBYyxTQUFRLGtCQUFrQjtJQUFyRDs7UUFDWSxXQUFNLEdBQUcsRUFBRSxDQUFBO0lBQ3ZCLENBQUM7Q0FBQTtBQUZELHNDQUVDO0FBRUQ7O0dBRUc7QUFDSCxNQUFhLFlBQWEsU0FBUSxhQUFhO0lBQS9DOztRQUNZLFdBQU0sR0FBRyxFQUFFLENBQUE7SUFDdkIsQ0FBQztDQUFBO0FBRkQsb0NBRUM7QUFFRDs7R0FFRztBQUNILE1BQWEsZUFBZ0IsU0FBUSxrQkFBa0I7SUFBdkQ7O1FBQ1ksV0FBTSxHQUFHLEVBQUUsQ0FBQTtJQUN2QixDQUFDO0NBQUE7QUFGRCwwQ0FFQztBQUVEOztHQUVHO0FBQ0gsTUFBYSxjQUFlLFNBQVEsa0JBQWtCO0lBQXREOztRQUNZLFdBQU0sR0FBRyxFQUFFLENBQUE7SUFDdkIsQ0FBQztDQUFBO0FBRkQsd0NBRUM7QUFFRDs7R0FFRztBQUNILE1BQWEsYUFBYyxTQUFRLGtCQUFrQjtJQUFyRDs7UUFDWSxXQUFNLEdBQUcsRUFBRSxDQUFBO0lBQ3ZCLENBQUM7Q0FBQTtBQUZELHNDQUVDO0FBRUQ7OztHQUdHO0FBQ0gsTUFBYSxjQUFlLFNBQVEsYUFBYTtJQUFqRDs7UUFDWSxXQUFNLEdBQUcsRUFBRSxDQUFBO0lBQ3ZCLENBQUM7Q0FBQTtBQUZELHdDQUVDO0FBRUQ7OztHQUdHO0FBQ0gsTUFBYSxjQUFlLFNBQVEsYUFBYTtJQUFqRDs7UUFDWSxXQUFNLEdBQUcsRUFBRSxDQUFBO0lBQ3ZCLENBQUM7Q0FBQTtBQUZELHdDQUVDO0FBRUQ7O0dBRUc7QUFDSCxNQUFhLFdBQVksU0FBUSxVQUFVO0lBQTNDOztRQUNZLFdBQU0sR0FBRyxFQUFFLENBQUE7SUFDdkIsQ0FBQztDQUFBO0FBRkQsa0NBRUM7QUFFRCxNQUFhLFVBQVcsU0FBUSxVQUFVO0lBQTFDOztRQUNZLFdBQU0sR0FBRyxFQUFFLENBQUE7SUFDdkIsQ0FBQztDQUFBO0FBRkQsZ0NBRUM7QUFFRDs7R0FFRztBQUNILE1BQWEsVUFBVyxTQUFRLFVBQVU7SUFBMUM7O1FBQ1ksV0FBTSxHQUFHLEVBQUUsQ0FBQTtJQUN2QixDQUFDO0NBQUE7QUFGRCxnQ0FFQztBQUVEOztHQUVHO0FBQ0gsTUFBYSxVQUFXLFNBQVEsVUFBVTtJQUExQzs7UUFDWSxXQUFNLEdBQUcsRUFBRSxDQUFBO0lBQ3ZCLENBQUM7Q0FBQTtBQUZELGdDQUVDO0FBRUQ7O0dBRUc7QUFDSCxNQUFhLFVBQVcsU0FBUSxXQUFXO0lBQTNDOztRQUNZLFdBQU0sR0FBRyxFQUFFLENBQUE7SUFDdkIsQ0FBQztDQUFBO0FBRkQsZ0NBRUM7QUFFRDs7R0FFRztBQUNILE1BQWEsVUFBVyxTQUFRLFdBQVc7SUFBM0M7O1FBQ1ksV0FBTSxHQUFHLEVBQUUsQ0FBQTtJQUN2QixDQUFDO0NBQUE7QUFGRCxnQ0FFQztBQUVEOztHQUVHO0FBQ0gsTUFBYSxXQUFZLFNBQVEsV0FBVztJQVUxQyxZQUFZLFVBQWUsU0FBUztRQUNsQyxLQUFLLEVBQUUsQ0FBQTtRQVZDLFdBQU0sR0FBRyxFQUFFLENBQUE7UUFXbkIsSUFBSSxPQUFPLFlBQVksZUFBTSxFQUFFO1lBQzdCLElBQUksQ0FBQyxPQUFPLEdBQUcsT0FBTyxDQUFBO1NBQ3ZCO2FBQU0sSUFBSSxPQUFPLE9BQU8sS0FBSyxRQUFRLEVBQUU7WUFDdEMsSUFBSSxDQUFDLE9BQU8sR0FBRyxlQUFNLENBQUMsSUFBSSxDQUFDLE9BQU8sRUFBRSxNQUFNLENBQUMsQ0FBQTtTQUM1QzthQUFNLElBQUksT0FBTyxFQUFFO1lBQ2xCLElBQUksT0FBTyxHQUFXLElBQUksQ0FBQyxTQUFTLENBQUMsT0FBTyxDQUFDLENBQUE7WUFDN0MsSUFBSSxDQUFDLE9BQU8sR0FBRyxlQUFNLENBQUMsSUFBSSxDQUFDLE9BQU8sRUFBRSxNQUFNLENBQUMsQ0FBQTtTQUM1QztJQUNILENBQUM7SUFqQkQ7O09BRUc7SUFDSCxVQUFVO1FBQ1IsT0FBTyxJQUFJLENBQUMsS0FBSyxDQUFDLElBQUksQ0FBQyxPQUFPLENBQUMsUUFBUSxDQUFDLE1BQU0sQ0FBQyxDQUFDLENBQUE7SUFDbEQsQ0FBQztDQWFGO0FBckJELGtDQXFCQztBQUVEOztHQUVHO0FBQ0gsTUFBYSxXQUFZLFNBQVEsV0FBVztJQUE1Qzs7UUFDWSxXQUFNLEdBQUcsRUFBRSxDQUFBO0lBQ3ZCLENBQUM7Q0FBQTtBQUZELGtDQUVDO0FBRUQ7O0dBRUc7QUFDSCxNQUFhLFlBQWEsU0FBUSxXQUFXO0lBQTdDOztRQUNZLFdBQU0sR0FBRyxFQUFFLENBQUE7SUFDdkIsQ0FBQztDQUFBO0FBRkQsb0NBRUM7QUFFRDs7R0FFRztBQUNILE1BQWEsVUFBVyxTQUFRLFdBQVc7SUFBM0M7O1FBQ1ksV0FBTSxHQUFHLEVBQUUsQ0FBQTtJQUN2QixDQUFDO0NBQUE7QUFGRCxnQ0FFQztBQUVEOztHQUVHO0FBQ0gsTUFBYSxXQUFZLFNBQVEsYUFBYTtJQUE5Qzs7UUFDWSxXQUFNLEdBQUcsRUFBRSxDQUFBO0lBQ3ZCLENBQUM7Q0FBQTtBQUZELGtDQUVDO0FBRUQ7O0dBRUc7QUFDSCxNQUFhLFlBQWEsU0FBUSxXQUFXO0lBQTdDOztRQUNZLFdBQU0sR0FBRyxFQUFFLENBQUE7SUFDdkIsQ0FBQztDQUFBO0FBRkQsb0NBRUM7QUFFRDs7R0FFRztBQUNILE1BQWEsYUFBYyxTQUFRLFdBQVc7SUFBOUM7O1FBQ1ksV0FBTSxHQUFHLEVBQUUsQ0FBQTtJQUN2QixDQUFDO0NBQUE7QUFGRCxzQ0FFQyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogQHBhY2thZ2VEb2N1bWVudGF0aW9uXG4gKiBAbW9kdWxlIFV0aWxzLVBheWxvYWRcbiAqL1xuXG5pbXBvcnQgeyBCdWZmZXIgfSBmcm9tIFwiYnVmZmVyL1wiXG5pbXBvcnQgQmluVG9vbHMgZnJvbSBcIi4vYmludG9vbHNcIlxuaW1wb3J0IEJOIGZyb20gXCJibi5qc1wiXG5pbXBvcnQgeyBUeXBlSWRFcnJvciwgSGV4RXJyb3IgfSBmcm9tIFwiLi4vdXRpbHMvZXJyb3JzXCJcbmltcG9ydCB7IFNlcmlhbGl6YXRpb24sIFNlcmlhbGl6ZWRUeXBlIH0gZnJvbSBcIi4uL3V0aWxzL3NlcmlhbGl6YXRpb25cIlxuXG4vKipcbiAqIEBpZ25vcmVcbiAqL1xuY29uc3QgYmludG9vbHM6IEJpblRvb2xzID0gQmluVG9vbHMuZ2V0SW5zdGFuY2UoKVxuY29uc3Qgc2VyaWFsaXphdGlvbjogU2VyaWFsaXphdGlvbiA9IFNlcmlhbGl6YXRpb24uZ2V0SW5zdGFuY2UoKVxuXG4vKipcbiAqIENsYXNzIGZvciBkZXRlcm1pbmluZyBwYXlsb2FkIHR5cGVzIGFuZCBtYW5hZ2luZyB0aGUgbG9va3VwIHRhYmxlLlxuICovXG5leHBvcnQgY2xhc3MgUGF5bG9hZFR5cGVzIHtcbiAgcHJpdmF0ZSBzdGF0aWMgaW5zdGFuY2U6IFBheWxvYWRUeXBlc1xuICBwcm90ZWN0ZWQgdHlwZXM6IHN0cmluZ1tdID0gW11cblxuICAvKipcbiAgICogR2l2ZW4gYW4gZW5jb2RlZCBwYXlsb2FkIGJ1ZmZlciByZXR1cm5zIHRoZSBwYXlsb2FkIGNvbnRlbnQgKG1pbnVzIHR5cGVJRCkuXG4gICAqL1xuICBnZXRDb250ZW50KHBheWxvYWQ6IEJ1ZmZlcik6IEJ1ZmZlciB7XG4gICAgY29uc3QgcGw6IEJ1ZmZlciA9IGJpbnRvb2xzLmNvcHlGcm9tKHBheWxvYWQsIDUpXG4gICAgcmV0dXJuIHBsXG4gIH1cblxuICAvKipcbiAgICogR2l2ZW4gYW4gZW5jb2RlZCBwYXlsb2FkIGJ1ZmZlciByZXR1cm5zIHRoZSBwYXlsb2FkICh3aXRoIHR5cGVJRCkuXG4gICAqL1xuICBnZXRQYXlsb2FkKHBheWxvYWQ6IEJ1ZmZlcik6IEJ1ZmZlciB7XG4gICAgY29uc3QgcGw6IEJ1ZmZlciA9IGJpbnRvb2xzLmNvcHlGcm9tKHBheWxvYWQsIDQpXG4gICAgcmV0dXJuIHBsXG4gIH1cblxuICAvKipcbiAgICogR2l2ZW4gYSBwYXlsb2FkIGJ1ZmZlciByZXR1cm5zIHRoZSBwcm9wZXIgVHlwZUlELlxuICAgKi9cbiAgZ2V0VHlwZUlEKHBheWxvYWQ6IEJ1ZmZlcik6IG51bWJlciB7XG4gICAgY29uc3Qgb2Zmc2V0OiBudW1iZXIgPSA0XG4gICAgY29uc3QgdHlwZUlEOiBudW1iZXIgPSBiaW50b29sc1xuICAgICAgLmNvcHlGcm9tKHBheWxvYWQsIG9mZnNldCwgb2Zmc2V0ICsgMSlcbiAgICAgIC5yZWFkVUludDgoMClcbiAgICByZXR1cm4gdHlwZUlEXG4gIH1cblxuICAvKipcbiAgICogR2l2ZW4gYSB0eXBlIHN0cmluZyByZXR1cm5zIHRoZSBwcm9wZXIgVHlwZUlELlxuICAgKi9cbiAgbG9va3VwSUQodHlwZXN0cjogc3RyaW5nKTogbnVtYmVyIHtcbiAgICByZXR1cm4gdGhpcy50eXBlcy5pbmRleE9mKHR5cGVzdHIpXG4gIH1cblxuICAvKipcbiAgICogR2l2ZW4gYSBUeXBlSUQgcmV0dXJucyBhIHN0cmluZyBkZXNjcmliaW5nIHRoZSBwYXlsb2FkIHR5cGUuXG4gICAqL1xuICBsb29rdXBUeXBlKHZhbHVlOiBudW1iZXIpOiBzdHJpbmcge1xuICAgIHJldHVybiB0aGlzLnR5cGVzW2Ake3ZhbHVlfWBdXG4gIH1cblxuICAvKipcbiAgICogR2l2ZW4gYSBUeXBlSUQgcmV0dXJucyB0aGUgcHJvcGVyIFtbUGF5bG9hZEJhc2VdXS5cbiAgICovXG4gIHNlbGVjdCh0eXBlSUQ6IG51bWJlciwgLi4uYXJnczogYW55W10pOiBQYXlsb2FkQmFzZSB7XG4gICAgc3dpdGNoICh0eXBlSUQpIHtcbiAgICAgIGNhc2UgMDpcbiAgICAgICAgcmV0dXJuIG5ldyBCSU5QYXlsb2FkKC4uLmFyZ3MpXG4gICAgICBjYXNlIDE6XG4gICAgICAgIHJldHVybiBuZXcgVVRGOFBheWxvYWQoLi4uYXJncylcbiAgICAgIGNhc2UgMjpcbiAgICAgICAgcmV0dXJuIG5ldyBIRVhTVFJQYXlsb2FkKC4uLmFyZ3MpXG4gICAgICBjYXNlIDM6XG4gICAgICAgIHJldHVybiBuZXcgQjU4U1RSUGF5bG9hZCguLi5hcmdzKVxuICAgICAgY2FzZSA0OlxuICAgICAgICByZXR1cm4gbmV3IEI2NFNUUlBheWxvYWQoLi4uYXJncylcbiAgICAgIGNhc2UgNTpcbiAgICAgICAgcmV0dXJuIG5ldyBCSUdOVU1QYXlsb2FkKC4uLmFyZ3MpXG4gICAgICBjYXNlIDY6XG4gICAgICAgIHJldHVybiBuZXcgWENIQUlOQUREUlBheWxvYWQoLi4uYXJncylcbiAgICAgIGNhc2UgNzpcbiAgICAgICAgcmV0dXJuIG5ldyBQQ0hBSU5BRERSUGF5bG9hZCguLi5hcmdzKVxuICAgICAgY2FzZSA4OlxuICAgICAgICByZXR1cm4gbmV3IENDSEFJTkFERFJQYXlsb2FkKC4uLmFyZ3MpXG4gICAgICBjYXNlIDk6XG4gICAgICAgIHJldHVybiBuZXcgVFhJRFBheWxvYWQoLi4uYXJncylcbiAgICAgIGNhc2UgMTA6XG4gICAgICAgIHJldHVybiBuZXcgQVNTRVRJRFBheWxvYWQoLi4uYXJncylcbiAgICAgIGNhc2UgMTE6XG4gICAgICAgIHJldHVybiBuZXcgVVRYT0lEUGF5bG9hZCguLi5hcmdzKVxuICAgICAgY2FzZSAxMjpcbiAgICAgICAgcmV0dXJuIG5ldyBORlRJRFBheWxvYWQoLi4uYXJncylcbiAgICAgIGNhc2UgMTM6XG4gICAgICAgIHJldHVybiBuZXcgU1VCTkVUSURQYXlsb2FkKC4uLmFyZ3MpXG4gICAgICBjYXNlIDE0OlxuICAgICAgICByZXR1cm4gbmV3IENIQUlOSURQYXlsb2FkKC4uLmFyZ3MpXG4gICAgICBjYXNlIDE1OlxuICAgICAgICByZXR1cm4gbmV3IE5PREVJRFBheWxvYWQoLi4uYXJncylcbiAgICAgIGNhc2UgMTY6XG4gICAgICAgIHJldHVybiBuZXcgU0VDUFNJR1BheWxvYWQoLi4uYXJncylcbiAgICAgIGNhc2UgMTc6XG4gICAgICAgIHJldHVybiBuZXcgU0VDUEVOQ1BheWxvYWQoLi4uYXJncylcbiAgICAgIGNhc2UgMTg6XG4gICAgICAgIHJldHVybiBuZXcgSlBFR1BheWxvYWQoLi4uYXJncylcbiAgICAgIGNhc2UgMTk6XG4gICAgICAgIHJldHVybiBuZXcgUE5HUGF5bG9hZCguLi5hcmdzKVxuICAgICAgY2FzZSAyMDpcbiAgICAgICAgcmV0dXJuIG5ldyBCTVBQYXlsb2FkKC4uLmFyZ3MpXG4gICAgICBjYXNlIDIxOlxuICAgICAgICByZXR1cm4gbmV3IElDT1BheWxvYWQoLi4uYXJncylcbiAgICAgIGNhc2UgMjI6XG4gICAgICAgIHJldHVybiBuZXcgU1ZHUGF5bG9hZCguLi5hcmdzKVxuICAgICAgY2FzZSAyMzpcbiAgICAgICAgcmV0dXJuIG5ldyBDU1ZQYXlsb2FkKC4uLmFyZ3MpXG4gICAgICBjYXNlIDI0OlxuICAgICAgICByZXR1cm4gbmV3IEpTT05QYXlsb2FkKC4uLmFyZ3MpXG4gICAgICBjYXNlIDI1OlxuICAgICAgICByZXR1cm4gbmV3IFlBTUxQYXlsb2FkKC4uLmFyZ3MpXG4gICAgICBjYXNlIDI2OlxuICAgICAgICByZXR1cm4gbmV3IEVNQUlMUGF5bG9hZCguLi5hcmdzKVxuICAgICAgY2FzZSAyNzpcbiAgICAgICAgcmV0dXJuIG5ldyBVUkxQYXlsb2FkKC4uLmFyZ3MpXG4gICAgICBjYXNlIDI4OlxuICAgICAgICByZXR1cm4gbmV3IElQRlNQYXlsb2FkKC4uLmFyZ3MpXG4gICAgICBjYXNlIDI5OlxuICAgICAgICByZXR1cm4gbmV3IE9OSU9OUGF5bG9hZCguLi5hcmdzKVxuICAgICAgY2FzZSAzMDpcbiAgICAgICAgcmV0dXJuIG5ldyBNQUdORVRQYXlsb2FkKC4uLmFyZ3MpXG4gICAgfVxuICAgIHRocm93IG5ldyBUeXBlSWRFcnJvcihcbiAgICAgIGBFcnJvciAtIFBheWxvYWRUeXBlcy5zZWxlY3Q6IHVua25vd24gdHlwZWlkICR7dHlwZUlEfWBcbiAgICApXG4gIH1cblxuICAvKipcbiAgICogR2l2ZW4gYSBbW1BheWxvYWRCYXNlXV0gd2hpY2ggbWF5IG5vdCBiZSBjYXN0IHByb3Blcmx5LCByZXR1cm5zIGEgcHJvcGVybHkgY2FzdCBbW1BheWxvYWRCYXNlXV0uXG4gICAqL1xuICByZWNhc3QodW5rbm93UGF5bG9hZDogUGF5bG9hZEJhc2UpOiBQYXlsb2FkQmFzZSB7XG4gICAgcmV0dXJuIHRoaXMuc2VsZWN0KHVua25vd1BheWxvYWQudHlwZUlEKCksIHVua25vd1BheWxvYWQucmV0dXJuVHlwZSgpKVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIFtbUGF5bG9hZFR5cGVzXV0gc2luZ2xldG9uLlxuICAgKi9cbiAgc3RhdGljIGdldEluc3RhbmNlKCk6IFBheWxvYWRUeXBlcyB7XG4gICAgaWYgKCFQYXlsb2FkVHlwZXMuaW5zdGFuY2UpIHtcbiAgICAgIFBheWxvYWRUeXBlcy5pbnN0YW5jZSA9IG5ldyBQYXlsb2FkVHlwZXMoKVxuICAgIH1cblxuICAgIHJldHVybiBQYXlsb2FkVHlwZXMuaW5zdGFuY2VcbiAgfVxuXG4gIHByaXZhdGUgY29uc3RydWN0b3IoKSB7XG4gICAgdGhpcy50eXBlcyA9IFtcbiAgICAgIFwiQklOXCIsXG4gICAgICBcIlVURjhcIixcbiAgICAgIFwiSEVYU1RSXCIsXG4gICAgICBcIkI1OFNUUlwiLFxuICAgICAgXCJCNjRTVFJcIixcbiAgICAgIFwiQklHTlVNXCIsXG4gICAgICBcIlhDSEFJTkFERFJcIixcbiAgICAgIFwiUENIQUlOQUREUlwiLFxuICAgICAgXCJDQ0hBSU5BRERSXCIsXG4gICAgICBcIlRYSURcIixcbiAgICAgIFwiQVNTRVRJRFwiLFxuICAgICAgXCJVVFhPSURcIixcbiAgICAgIFwiTkZUSURcIixcbiAgICAgIFwiU1VCTkVUSURcIixcbiAgICAgIFwiQ0hBSU5JRFwiLFxuICAgICAgXCJOT0RFSURcIixcbiAgICAgIFwiU0VDUFNJR1wiLFxuICAgICAgXCJTRUNQRU5DXCIsXG4gICAgICBcIkpQRUdcIixcbiAgICAgIFwiUE5HXCIsXG4gICAgICBcIkJNUFwiLFxuICAgICAgXCJJQ09cIixcbiAgICAgIFwiU1ZHXCIsXG4gICAgICBcIkNTVlwiLFxuICAgICAgXCJKU09OXCIsXG4gICAgICBcIllBTUxcIixcbiAgICAgIFwiRU1BSUxcIixcbiAgICAgIFwiVVJMXCIsXG4gICAgICBcIklQRlNcIixcbiAgICAgIFwiT05JT05cIixcbiAgICAgIFwiTUFHTkVUXCJcbiAgICBdXG4gIH1cbn1cblxuLyoqXG4gKiBCYXNlIGNsYXNzIGZvciBwYXlsb2Fkcy5cbiAqL1xuZXhwb3J0IGFic3RyYWN0IGNsYXNzIFBheWxvYWRCYXNlIHtcbiAgcHJvdGVjdGVkIHBheWxvYWQ6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYygwKVxuICBwcm90ZWN0ZWQgdHlwZWlkOiBudW1iZXIgPSB1bmRlZmluZWRcblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgVHlwZUlEIGZvciB0aGUgcGF5bG9hZC5cbiAgICovXG4gIHR5cGVJRCgpOiBudW1iZXIge1xuICAgIHJldHVybiB0aGlzLnR5cGVpZFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIHN0cmluZyBuYW1lIGZvciB0aGUgcGF5bG9hZCdzIHR5cGUuXG4gICAqL1xuICB0eXBlTmFtZSgpOiBzdHJpbmcge1xuICAgIHJldHVybiBQYXlsb2FkVHlwZXMuZ2V0SW5zdGFuY2UoKS5sb29rdXBUeXBlKHRoaXMudHlwZWlkKVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIHBheWxvYWQgY29udGVudCAobWludXMgdHlwZUlEKS5cbiAgICovXG4gIGdldENvbnRlbnQoKTogQnVmZmVyIHtcbiAgICBjb25zdCBwbDogQnVmZmVyID0gYmludG9vbHMuY29weUZyb20odGhpcy5wYXlsb2FkKVxuICAgIHJldHVybiBwbFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIHBheWxvYWQgKHdpdGggdHlwZUlEKS5cbiAgICovXG4gIGdldFBheWxvYWQoKTogQnVmZmVyIHtcbiAgICBjb25zdCB0eXBlSUQ6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYygxKVxuICAgIHR5cGVJRC53cml0ZVVJbnQ4KHRoaXMudHlwZWlkLCAwKVxuICAgIGNvbnN0IHBsOiBCdWZmZXIgPSBCdWZmZXIuY29uY2F0KFt0eXBlSUQsIGJpbnRvb2xzLmNvcHlGcm9tKHRoaXMucGF5bG9hZCldKVxuICAgIHJldHVybiBwbFxuICB9XG5cbiAgLyoqXG4gICAqIERlY29kZXMgdGhlIHBheWxvYWQgYXMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSBpbmNsdWRpbmcgNCBieXRlcyBmb3IgdGhlIGxlbmd0aCBhbmQgVHlwZUlELlxuICAgKi9cbiAgZnJvbUJ1ZmZlcihieXRlczogQnVmZmVyLCBvZmZzZXQ6IG51bWJlciA9IDApOiBudW1iZXIge1xuICAgIGNvbnN0IHNpemU6IG51bWJlciA9IGJpbnRvb2xzXG4gICAgICAuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgNClcbiAgICAgIC5yZWFkVUludDMyQkUoMClcbiAgICBvZmZzZXQgKz0gNFxuICAgIHRoaXMudHlwZWlkID0gYmludG9vbHMuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgMSkucmVhZFVJbnQ4KDApXG4gICAgb2Zmc2V0ICs9IDFcbiAgICB0aGlzLnBheWxvYWQgPSBiaW50b29scy5jb3B5RnJvbShieXRlcywgb2Zmc2V0LCBvZmZzZXQgKyBzaXplIC0gMSlcbiAgICBvZmZzZXQgKz0gc2l6ZSAtIDFcbiAgICByZXR1cm4gb2Zmc2V0XG4gIH1cblxuICAvKipcbiAgICogRW5jb2RlcyB0aGUgcGF5bG9hZCBhcyBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IGluY2x1ZGluZyA0IGJ5dGVzIGZvciB0aGUgbGVuZ3RoIGFuZCBUeXBlSUQuXG4gICAqL1xuICB0b0J1ZmZlcigpOiBCdWZmZXIge1xuICAgIGNvbnN0IHNpemVidWZmOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoNClcbiAgICBzaXplYnVmZi53cml0ZVVJbnQzMkJFKHRoaXMucGF5bG9hZC5sZW5ndGggKyAxLCAwKVxuICAgIGNvbnN0IHR5cGVidWZmOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoMSlcbiAgICB0eXBlYnVmZi53cml0ZVVJbnQ4KHRoaXMudHlwZWlkLCAwKVxuICAgIHJldHVybiBCdWZmZXIuY29uY2F0KFtzaXplYnVmZiwgdHlwZWJ1ZmYsIHRoaXMucGF5bG9hZF0pXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgZXhwZWN0ZWQgdHlwZSBmb3IgdGhlIHBheWxvYWQuXG4gICAqL1xuICBhYnN0cmFjdCByZXR1cm5UeXBlKC4uLmFyZ3M6IGFueSk6IGFueVxuXG4gIGNvbnN0cnVjdG9yKCkge31cbn1cblxuLyoqXG4gKiBDbGFzcyBmb3IgcGF5bG9hZHMgcmVwcmVzZW50aW5nIHNpbXBsZSBiaW5hcnkgYmxvYnMuXG4gKi9cbmV4cG9ydCBjbGFzcyBCSU5QYXlsb2FkIGV4dGVuZHMgUGF5bG9hZEJhc2Uge1xuICBwcm90ZWN0ZWQgdHlwZWlkID0gMFxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIGEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gZm9yIHRoZSBwYXlsb2FkLlxuICAgKi9cbiAgcmV0dXJuVHlwZSgpOiBCdWZmZXIge1xuICAgIHJldHVybiB0aGlzLnBheWxvYWRcbiAgfVxuICAvKipcbiAgICogQHBhcmFtIHBheWxvYWQgQnVmZmVyIG9ubHlcbiAgICovXG4gIGNvbnN0cnVjdG9yKHBheWxvYWQ6IGFueSA9IHVuZGVmaW5lZCkge1xuICAgIHN1cGVyKClcbiAgICBpZiAocGF5bG9hZCBpbnN0YW5jZW9mIEJ1ZmZlcikge1xuICAgICAgdGhpcy5wYXlsb2FkID0gcGF5bG9hZFxuICAgIH0gZWxzZSBpZiAodHlwZW9mIHBheWxvYWQgPT09IFwic3RyaW5nXCIpIHtcbiAgICAgIHRoaXMucGF5bG9hZCA9IGJpbnRvb2xzLmI1OFRvQnVmZmVyKHBheWxvYWQpXG4gICAgfVxuICB9XG59XG5cbi8qKlxuICogQ2xhc3MgZm9yIHBheWxvYWRzIHJlcHJlc2VudGluZyBVVEY4IGVuY29kaW5nLlxuICovXG5leHBvcnQgY2xhc3MgVVRGOFBheWxvYWQgZXh0ZW5kcyBQYXlsb2FkQmFzZSB7XG4gIHByb3RlY3RlZCB0eXBlaWQgPSAxXG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSBzdHJpbmcgZm9yIHRoZSBwYXlsb2FkLlxuICAgKi9cbiAgcmV0dXJuVHlwZSgpOiBzdHJpbmcge1xuICAgIHJldHVybiB0aGlzLnBheWxvYWQudG9TdHJpbmcoXCJ1dGY4XCIpXG4gIH1cbiAgLyoqXG4gICAqIEBwYXJhbSBwYXlsb2FkIEJ1ZmZlciB1dGY4IHN0cmluZ1xuICAgKi9cbiAgY29uc3RydWN0b3IocGF5bG9hZDogYW55ID0gdW5kZWZpbmVkKSB7XG4gICAgc3VwZXIoKVxuICAgIGlmIChwYXlsb2FkIGluc3RhbmNlb2YgQnVmZmVyKSB7XG4gICAgICB0aGlzLnBheWxvYWQgPSBwYXlsb2FkXG4gICAgfSBlbHNlIGlmICh0eXBlb2YgcGF5bG9hZCA9PT0gXCJzdHJpbmdcIikge1xuICAgICAgdGhpcy5wYXlsb2FkID0gQnVmZmVyLmZyb20ocGF5bG9hZCwgXCJ1dGY4XCIpXG4gICAgfVxuICB9XG59XG5cbi8qKlxuICogQ2xhc3MgZm9yIHBheWxvYWRzIHJlcHJlc2VudGluZyBIZXhhZGVjaW1hbCBlbmNvZGluZy5cbiAqL1xuZXhwb3J0IGNsYXNzIEhFWFNUUlBheWxvYWQgZXh0ZW5kcyBQYXlsb2FkQmFzZSB7XG4gIHByb3RlY3RlZCB0eXBlaWQgPSAyXG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSBoZXggc3RyaW5nIGZvciB0aGUgcGF5bG9hZC5cbiAgICovXG4gIHJldHVyblR5cGUoKTogc3RyaW5nIHtcbiAgICByZXR1cm4gdGhpcy5wYXlsb2FkLnRvU3RyaW5nKFwiaGV4XCIpXG4gIH1cbiAgLyoqXG4gICAqIEBwYXJhbSBwYXlsb2FkIEJ1ZmZlciBvciBoZXggc3RyaW5nXG4gICAqL1xuICBjb25zdHJ1Y3RvcihwYXlsb2FkOiBhbnkgPSB1bmRlZmluZWQpIHtcbiAgICBzdXBlcigpXG4gICAgaWYgKHBheWxvYWQgaW5zdGFuY2VvZiBCdWZmZXIpIHtcbiAgICAgIHRoaXMucGF5bG9hZCA9IHBheWxvYWRcbiAgICB9IGVsc2UgaWYgKHR5cGVvZiBwYXlsb2FkID09PSBcInN0cmluZ1wiKSB7XG4gICAgICBpZiAocGF5bG9hZC5zdGFydHNXaXRoKFwiMHhcIikgfHwgIXBheWxvYWQubWF0Y2goL15bMC05QS1GYS1mXSskLykpIHtcbiAgICAgICAgdGhyb3cgbmV3IEhleEVycm9yKFxuICAgICAgICAgIFwiSEVYU1RSUGF5bG9hZC5jb25zdHJ1Y3RvciAtLSBoZXggc3RyaW5nIG1heSBub3Qgc3RhcnQgd2l0aCAweCBhbmQgbXVzdCBiZSBpbiAvXlswLTlBLUZhLWZdKyQvOiBcIiArXG4gICAgICAgICAgICBwYXlsb2FkXG4gICAgICAgIClcbiAgICAgIH1cbiAgICAgIHRoaXMucGF5bG9hZCA9IEJ1ZmZlci5mcm9tKHBheWxvYWQsIFwiaGV4XCIpXG4gICAgfVxuICB9XG59XG5cbi8qKlxuICogQ2xhc3MgZm9yIHBheWxvYWRzIHJlcHJlc2VudGluZyBCYXNlNTggZW5jb2RpbmcuXG4gKi9cbmV4cG9ydCBjbGFzcyBCNThTVFJQYXlsb2FkIGV4dGVuZHMgUGF5bG9hZEJhc2Uge1xuICBwcm90ZWN0ZWQgdHlwZWlkID0gM1xuXG4gIC8qKlxuICAgKiBSZXR1cm5zIGEgYmFzZTU4IHN0cmluZyBmb3IgdGhlIHBheWxvYWQuXG4gICAqL1xuICByZXR1cm5UeXBlKCk6IHN0cmluZyB7XG4gICAgcmV0dXJuIGJpbnRvb2xzLmJ1ZmZlclRvQjU4KHRoaXMucGF5bG9hZClcbiAgfVxuICAvKipcbiAgICogQHBhcmFtIHBheWxvYWQgQnVmZmVyIG9yIGNiNTggZW5jb2RlZCBzdHJpbmdcbiAgICovXG4gIGNvbnN0cnVjdG9yKHBheWxvYWQ6IGFueSA9IHVuZGVmaW5lZCkge1xuICAgIHN1cGVyKClcbiAgICBpZiAocGF5bG9hZCBpbnN0YW5jZW9mIEJ1ZmZlcikge1xuICAgICAgdGhpcy5wYXlsb2FkID0gcGF5bG9hZFxuICAgIH0gZWxzZSBpZiAodHlwZW9mIHBheWxvYWQgPT09IFwic3RyaW5nXCIpIHtcbiAgICAgIHRoaXMucGF5bG9hZCA9IGJpbnRvb2xzLmI1OFRvQnVmZmVyKHBheWxvYWQpXG4gICAgfVxuICB9XG59XG5cbi8qKlxuICogQ2xhc3MgZm9yIHBheWxvYWRzIHJlcHJlc2VudGluZyBCYXNlNjQgZW5jb2RpbmcuXG4gKi9cbmV4cG9ydCBjbGFzcyBCNjRTVFJQYXlsb2FkIGV4dGVuZHMgUGF5bG9hZEJhc2Uge1xuICBwcm90ZWN0ZWQgdHlwZWlkID0gNFxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIGEgYmFzZTY0IHN0cmluZyBmb3IgdGhlIHBheWxvYWQuXG4gICAqL1xuICByZXR1cm5UeXBlKCk6IHN0cmluZyB7XG4gICAgcmV0dXJuIHRoaXMucGF5bG9hZC50b1N0cmluZyhcImJhc2U2NFwiKVxuICB9XG4gIC8qKlxuICAgKiBAcGFyYW0gcGF5bG9hZCBCdWZmZXIgb2YgYmFzZTY0IHN0cmluZ1xuICAgKi9cbiAgY29uc3RydWN0b3IocGF5bG9hZDogYW55ID0gdW5kZWZpbmVkKSB7XG4gICAgc3VwZXIoKVxuICAgIGlmIChwYXlsb2FkIGluc3RhbmNlb2YgQnVmZmVyKSB7XG4gICAgICB0aGlzLnBheWxvYWQgPSBwYXlsb2FkXG4gICAgfSBlbHNlIGlmICh0eXBlb2YgcGF5bG9hZCA9PT0gXCJzdHJpbmdcIikge1xuICAgICAgdGhpcy5wYXlsb2FkID0gQnVmZmVyLmZyb20ocGF5bG9hZCwgXCJiYXNlNjRcIilcbiAgICB9XG4gIH1cbn1cblxuLyoqXG4gKiBDbGFzcyBmb3IgcGF5bG9hZHMgcmVwcmVzZW50aW5nIEJpZyBOdW1iZXJzLlxuICpcbiAqIEBwYXJhbSBwYXlsb2FkIEFjY2VwdHMgYSBCdWZmZXIsIEJOLCBvciBiYXNlNjQgc3RyaW5nXG4gKi9cbmV4cG9ydCBjbGFzcyBCSUdOVU1QYXlsb2FkIGV4dGVuZHMgUGF5bG9hZEJhc2Uge1xuICBwcm90ZWN0ZWQgdHlwZWlkID0gNVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIGEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9pbmR1dG55L2JuLmpzL3xCTn0gZm9yIHRoZSBwYXlsb2FkLlxuICAgKi9cbiAgcmV0dXJuVHlwZSgpOiBCTiB7XG4gICAgcmV0dXJuIGJpbnRvb2xzLmZyb21CdWZmZXJUb0JOKHRoaXMucGF5bG9hZClcbiAgfVxuICAvKipcbiAgICogQHBhcmFtIHBheWxvYWQgQnVmZmVyLCBCTiwgb3IgYmFzZTY0IHN0cmluZ1xuICAgKi9cbiAgY29uc3RydWN0b3IocGF5bG9hZDogYW55ID0gdW5kZWZpbmVkKSB7XG4gICAgc3VwZXIoKVxuICAgIGlmIChwYXlsb2FkIGluc3RhbmNlb2YgQnVmZmVyKSB7XG4gICAgICB0aGlzLnBheWxvYWQgPSBwYXlsb2FkXG4gICAgfSBlbHNlIGlmIChwYXlsb2FkIGluc3RhbmNlb2YgQk4pIHtcbiAgICAgIHRoaXMucGF5bG9hZCA9IGJpbnRvb2xzLmZyb21CTlRvQnVmZmVyKHBheWxvYWQpXG4gICAgfSBlbHNlIGlmICh0eXBlb2YgcGF5bG9hZCA9PT0gXCJzdHJpbmdcIikge1xuICAgICAgdGhpcy5wYXlsb2FkID0gQnVmZmVyLmZyb20ocGF5bG9hZCwgXCJoZXhcIilcbiAgICB9XG4gIH1cbn1cblxuLyoqXG4gKiBDbGFzcyBmb3IgcGF5bG9hZHMgcmVwcmVzZW50aW5nIGNoYWluIGFkZHJlc3Nlcy5cbiAqXG4gKi9cbmV4cG9ydCBhYnN0cmFjdCBjbGFzcyBDaGFpbkFkZHJlc3NQYXlsb2FkIGV4dGVuZHMgUGF5bG9hZEJhc2Uge1xuICBwcm90ZWN0ZWQgdHlwZWlkID0gNlxuICBwcm90ZWN0ZWQgY2hhaW5pZDogc3RyaW5nID0gXCJcIlxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBjaGFpbmlkLlxuICAgKi9cbiAgcmV0dXJuQ2hhaW5JRCgpOiBzdHJpbmcge1xuICAgIHJldHVybiB0aGlzLmNoYWluaWRcbiAgfVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIGFuIGFkZHJlc3Mgc3RyaW5nIGZvciB0aGUgcGF5bG9hZC5cbiAgICovXG4gIHJldHVyblR5cGUoaHJwOiBzdHJpbmcpOiBzdHJpbmcge1xuICAgIGNvbnN0IHR5cGU6IFNlcmlhbGl6ZWRUeXBlID0gXCJiZWNoMzJcIlxuICAgIHJldHVybiBzZXJpYWxpemF0aW9uLmJ1ZmZlclRvVHlwZSh0aGlzLnBheWxvYWQsIHR5cGUsIGhycCwgdGhpcy5jaGFpbmlkKVxuICB9XG4gIC8qKlxuICAgKiBAcGFyYW0gcGF5bG9hZCBCdWZmZXIgb3IgYWRkcmVzcyBzdHJpbmdcbiAgICovXG4gIGNvbnN0cnVjdG9yKHBheWxvYWQ6IGFueSA9IHVuZGVmaW5lZCwgaHJwPzogc3RyaW5nKSB7XG4gICAgc3VwZXIoKVxuICAgIGlmIChwYXlsb2FkIGluc3RhbmNlb2YgQnVmZmVyKSB7XG4gICAgICB0aGlzLnBheWxvYWQgPSBwYXlsb2FkXG4gICAgfSBlbHNlIGlmICh0eXBlb2YgcGF5bG9hZCA9PT0gXCJzdHJpbmdcIikge1xuICAgICAgaWYgKGhycCAhPSB1bmRlZmluZWQpIHtcbiAgICAgICAgdGhpcy5wYXlsb2FkID0gYmludG9vbHMuc3RyaW5nVG9BZGRyZXNzKHBheWxvYWQsIGhycClcbiAgICAgIH0gZWxzZSB7XG4gICAgICAgIHRoaXMucGF5bG9hZCA9IGJpbnRvb2xzLnN0cmluZ1RvQWRkcmVzcyhwYXlsb2FkKVxuICAgICAgfVxuICAgIH1cbiAgfVxufVxuXG4vKipcbiAqIENsYXNzIGZvciBwYXlsb2FkcyByZXByZXNlbnRpbmcgWC1DaGluIGFkZHJlc3Nlcy5cbiAqL1xuZXhwb3J0IGNsYXNzIFhDSEFJTkFERFJQYXlsb2FkIGV4dGVuZHMgQ2hhaW5BZGRyZXNzUGF5bG9hZCB7XG4gIHByb3RlY3RlZCB0eXBlaWQgPSA2XG4gIHByb3RlY3RlZCBjaGFpbmlkID0gXCJYXCJcbn1cblxuLyoqXG4gKiBDbGFzcyBmb3IgcGF5bG9hZHMgcmVwcmVzZW50aW5nIFAtQ2hhaW4gYWRkcmVzc2VzLlxuICovXG5leHBvcnQgY2xhc3MgUENIQUlOQUREUlBheWxvYWQgZXh0ZW5kcyBDaGFpbkFkZHJlc3NQYXlsb2FkIHtcbiAgcHJvdGVjdGVkIHR5cGVpZCA9IDdcbiAgcHJvdGVjdGVkIGNoYWluaWQgPSBcIlBcIlxufVxuXG4vKipcbiAqIENsYXNzIGZvciBwYXlsb2FkcyByZXByZXNlbnRpbmcgQy1DaGFpbiBhZGRyZXNzZXMuXG4gKi9cbmV4cG9ydCBjbGFzcyBDQ0hBSU5BRERSUGF5bG9hZCBleHRlbmRzIENoYWluQWRkcmVzc1BheWxvYWQge1xuICBwcm90ZWN0ZWQgdHlwZWlkID0gOFxuICBwcm90ZWN0ZWQgY2hhaW5pZCA9IFwiQ1wiXG59XG5cbi8qKlxuICogQ2xhc3MgZm9yIHBheWxvYWRzIHJlcHJlc2VudGluZyBkYXRhIHNlcmlhbGl6ZWQgYnkgYmludG9vbHMuY2I1OEVuY29kZSgpLlxuICovXG5leHBvcnQgYWJzdHJhY3QgY2xhc3MgY2I1OEVuY29kZWRQYXlsb2FkIGV4dGVuZHMgUGF5bG9hZEJhc2Uge1xuICAvKipcbiAgICogUmV0dXJucyBhIGJpbnRvb2xzLmNiNThFbmNvZGVkIHN0cmluZyBmb3IgdGhlIHBheWxvYWQuXG4gICAqL1xuICByZXR1cm5UeXBlKCk6IHN0cmluZyB7XG4gICAgcmV0dXJuIGJpbnRvb2xzLmNiNThFbmNvZGUodGhpcy5wYXlsb2FkKVxuICB9XG4gIC8qKlxuICAgKiBAcGFyYW0gcGF5bG9hZCBCdWZmZXIgb3IgY2I1OCBlbmNvZGVkIHN0cmluZ1xuICAgKi9cbiAgY29uc3RydWN0b3IocGF5bG9hZDogYW55ID0gdW5kZWZpbmVkKSB7XG4gICAgc3VwZXIoKVxuICAgIGlmIChwYXlsb2FkIGluc3RhbmNlb2YgQnVmZmVyKSB7XG4gICAgICB0aGlzLnBheWxvYWQgPSBwYXlsb2FkXG4gICAgfSBlbHNlIGlmICh0eXBlb2YgcGF5bG9hZCA9PT0gXCJzdHJpbmdcIikge1xuICAgICAgdGhpcy5wYXlsb2FkID0gYmludG9vbHMuY2I1OERlY29kZShwYXlsb2FkKVxuICAgIH1cbiAgfVxufVxuXG4vKipcbiAqIENsYXNzIGZvciBwYXlsb2FkcyByZXByZXNlbnRpbmcgVHhJRHMuXG4gKi9cbmV4cG9ydCBjbGFzcyBUWElEUGF5bG9hZCBleHRlbmRzIGNiNThFbmNvZGVkUGF5bG9hZCB7XG4gIHByb3RlY3RlZCB0eXBlaWQgPSA5XG59XG5cbi8qKlxuICogQ2xhc3MgZm9yIHBheWxvYWRzIHJlcHJlc2VudGluZyBBc3NldElEcy5cbiAqL1xuZXhwb3J0IGNsYXNzIEFTU0VUSURQYXlsb2FkIGV4dGVuZHMgY2I1OEVuY29kZWRQYXlsb2FkIHtcbiAgcHJvdGVjdGVkIHR5cGVpZCA9IDEwXG59XG5cbi8qKlxuICogQ2xhc3MgZm9yIHBheWxvYWRzIHJlcHJlc2VudGluZyBOT0RFSURzLlxuICovXG5leHBvcnQgY2xhc3MgVVRYT0lEUGF5bG9hZCBleHRlbmRzIGNiNThFbmNvZGVkUGF5bG9hZCB7XG4gIHByb3RlY3RlZCB0eXBlaWQgPSAxMVxufVxuXG4vKipcbiAqIENsYXNzIGZvciBwYXlsb2FkcyByZXByZXNlbnRpbmcgTkZUSURzIChVVFhPSURzIGluIGFuIE5GVCBjb250ZXh0KS5cbiAqL1xuZXhwb3J0IGNsYXNzIE5GVElEUGF5bG9hZCBleHRlbmRzIFVUWE9JRFBheWxvYWQge1xuICBwcm90ZWN0ZWQgdHlwZWlkID0gMTJcbn1cblxuLyoqXG4gKiBDbGFzcyBmb3IgcGF5bG9hZHMgcmVwcmVzZW50aW5nIFN1Ym5ldElEcy5cbiAqL1xuZXhwb3J0IGNsYXNzIFNVQk5FVElEUGF5bG9hZCBleHRlbmRzIGNiNThFbmNvZGVkUGF5bG9hZCB7XG4gIHByb3RlY3RlZCB0eXBlaWQgPSAxM1xufVxuXG4vKipcbiAqIENsYXNzIGZvciBwYXlsb2FkcyByZXByZXNlbnRpbmcgQ2hhaW5JRHMuXG4gKi9cbmV4cG9ydCBjbGFzcyBDSEFJTklEUGF5bG9hZCBleHRlbmRzIGNiNThFbmNvZGVkUGF5bG9hZCB7XG4gIHByb3RlY3RlZCB0eXBlaWQgPSAxNFxufVxuXG4vKipcbiAqIENsYXNzIGZvciBwYXlsb2FkcyByZXByZXNlbnRpbmcgTm9kZUlEcy5cbiAqL1xuZXhwb3J0IGNsYXNzIE5PREVJRFBheWxvYWQgZXh0ZW5kcyBjYjU4RW5jb2RlZFBheWxvYWQge1xuICBwcm90ZWN0ZWQgdHlwZWlkID0gMTVcbn1cblxuLyoqXG4gKiBDbGFzcyBmb3IgcGF5bG9hZHMgcmVwcmVzZW50aW5nIHNlY3AyNTZrMSBzaWduYXR1cmVzLlxuICogY29udmVudGlvbjogc2VjcDI1NmsxIHNpZ25hdHVyZSAoMTMwIGJ5dGVzKVxuICovXG5leHBvcnQgY2xhc3MgU0VDUFNJR1BheWxvYWQgZXh0ZW5kcyBCNThTVFJQYXlsb2FkIHtcbiAgcHJvdGVjdGVkIHR5cGVpZCA9IDE2XG59XG5cbi8qKlxuICogQ2xhc3MgZm9yIHBheWxvYWRzIHJlcHJlc2VudGluZyBzZWNwMjU2azEgZW5jcnlwdGVkIG1lc3NhZ2VzLlxuICogY29udmVudGlvbjogcHVibGljIGtleSAoNjUgYnl0ZXMpICsgc2VjcDI1NmsxIGVuY3J5cHRlZCBtZXNzYWdlIGZvciB0aGF0IHB1YmxpYyBrZXlcbiAqL1xuZXhwb3J0IGNsYXNzIFNFQ1BFTkNQYXlsb2FkIGV4dGVuZHMgQjU4U1RSUGF5bG9hZCB7XG4gIHByb3RlY3RlZCB0eXBlaWQgPSAxN1xufVxuXG4vKipcbiAqIENsYXNzIGZvciBwYXlsb2FkcyByZXByZXNlbnRpbmcgSlBFRyBpbWFnZXMuXG4gKi9cbmV4cG9ydCBjbGFzcyBKUEVHUGF5bG9hZCBleHRlbmRzIEJJTlBheWxvYWQge1xuICBwcm90ZWN0ZWQgdHlwZWlkID0gMThcbn1cblxuZXhwb3J0IGNsYXNzIFBOR1BheWxvYWQgZXh0ZW5kcyBCSU5QYXlsb2FkIHtcbiAgcHJvdGVjdGVkIHR5cGVpZCA9IDE5XG59XG5cbi8qKlxuICogQ2xhc3MgZm9yIHBheWxvYWRzIHJlcHJlc2VudGluZyBCTVAgaW1hZ2VzLlxuICovXG5leHBvcnQgY2xhc3MgQk1QUGF5bG9hZCBleHRlbmRzIEJJTlBheWxvYWQge1xuICBwcm90ZWN0ZWQgdHlwZWlkID0gMjBcbn1cblxuLyoqXG4gKiBDbGFzcyBmb3IgcGF5bG9hZHMgcmVwcmVzZW50aW5nIElDTyBpbWFnZXMuXG4gKi9cbmV4cG9ydCBjbGFzcyBJQ09QYXlsb2FkIGV4dGVuZHMgQklOUGF5bG9hZCB7XG4gIHByb3RlY3RlZCB0eXBlaWQgPSAyMVxufVxuXG4vKipcbiAqIENsYXNzIGZvciBwYXlsb2FkcyByZXByZXNlbnRpbmcgU1ZHIGltYWdlcy5cbiAqL1xuZXhwb3J0IGNsYXNzIFNWR1BheWxvYWQgZXh0ZW5kcyBVVEY4UGF5bG9hZCB7XG4gIHByb3RlY3RlZCB0eXBlaWQgPSAyMlxufVxuXG4vKipcbiAqIENsYXNzIGZvciBwYXlsb2FkcyByZXByZXNlbnRpbmcgQ1NWIGZpbGVzLlxuICovXG5leHBvcnQgY2xhc3MgQ1NWUGF5bG9hZCBleHRlbmRzIFVURjhQYXlsb2FkIHtcbiAgcHJvdGVjdGVkIHR5cGVpZCA9IDIzXG59XG5cbi8qKlxuICogQ2xhc3MgZm9yIHBheWxvYWRzIHJlcHJlc2VudGluZyBKU09OIHN0cmluZ3MuXG4gKi9cbmV4cG9ydCBjbGFzcyBKU09OUGF5bG9hZCBleHRlbmRzIFBheWxvYWRCYXNlIHtcbiAgcHJvdGVjdGVkIHR5cGVpZCA9IDI0XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSBKU09OLWRlY29kZWQgb2JqZWN0IGZvciB0aGUgcGF5bG9hZC5cbiAgICovXG4gIHJldHVyblR5cGUoKTogYW55IHtcbiAgICByZXR1cm4gSlNPTi5wYXJzZSh0aGlzLnBheWxvYWQudG9TdHJpbmcoXCJ1dGY4XCIpKVxuICB9XG5cbiAgY29uc3RydWN0b3IocGF5bG9hZDogYW55ID0gdW5kZWZpbmVkKSB7XG4gICAgc3VwZXIoKVxuICAgIGlmIChwYXlsb2FkIGluc3RhbmNlb2YgQnVmZmVyKSB7XG4gICAgICB0aGlzLnBheWxvYWQgPSBwYXlsb2FkXG4gICAgfSBlbHNlIGlmICh0eXBlb2YgcGF5bG9hZCA9PT0gXCJzdHJpbmdcIikge1xuICAgICAgdGhpcy5wYXlsb2FkID0gQnVmZmVyLmZyb20ocGF5bG9hZCwgXCJ1dGY4XCIpXG4gICAgfSBlbHNlIGlmIChwYXlsb2FkKSB7XG4gICAgICBsZXQganNvbnN0cjogc3RyaW5nID0gSlNPTi5zdHJpbmdpZnkocGF5bG9hZClcbiAgICAgIHRoaXMucGF5bG9hZCA9IEJ1ZmZlci5mcm9tKGpzb25zdHIsIFwidXRmOFwiKVxuICAgIH1cbiAgfVxufVxuXG4vKipcbiAqIENsYXNzIGZvciBwYXlsb2FkcyByZXByZXNlbnRpbmcgWUFNTCBkZWZpbml0aW9ucy5cbiAqL1xuZXhwb3J0IGNsYXNzIFlBTUxQYXlsb2FkIGV4dGVuZHMgVVRGOFBheWxvYWQge1xuICBwcm90ZWN0ZWQgdHlwZWlkID0gMjVcbn1cblxuLyoqXG4gKiBDbGFzcyBmb3IgcGF5bG9hZHMgcmVwcmVzZW50aW5nIGVtYWlsIGFkZHJlc3Nlcy5cbiAqL1xuZXhwb3J0IGNsYXNzIEVNQUlMUGF5bG9hZCBleHRlbmRzIFVURjhQYXlsb2FkIHtcbiAgcHJvdGVjdGVkIHR5cGVpZCA9IDI2XG59XG5cbi8qKlxuICogQ2xhc3MgZm9yIHBheWxvYWRzIHJlcHJlc2VudGluZyBVUkwgc3RyaW5ncy5cbiAqL1xuZXhwb3J0IGNsYXNzIFVSTFBheWxvYWQgZXh0ZW5kcyBVVEY4UGF5bG9hZCB7XG4gIHByb3RlY3RlZCB0eXBlaWQgPSAyN1xufVxuXG4vKipcbiAqIENsYXNzIGZvciBwYXlsb2FkcyByZXByZXNlbnRpbmcgSVBGUyBhZGRyZXNzZXMuXG4gKi9cbmV4cG9ydCBjbGFzcyBJUEZTUGF5bG9hZCBleHRlbmRzIEI1OFNUUlBheWxvYWQge1xuICBwcm90ZWN0ZWQgdHlwZWlkID0gMjhcbn1cblxuLyoqXG4gKiBDbGFzcyBmb3IgcGF5bG9hZHMgcmVwcmVzZW50aW5nIG9uaW9uIFVSTHMuXG4gKi9cbmV4cG9ydCBjbGFzcyBPTklPTlBheWxvYWQgZXh0ZW5kcyBVVEY4UGF5bG9hZCB7XG4gIHByb3RlY3RlZCB0eXBlaWQgPSAyOVxufVxuXG4vKipcbiAqIENsYXNzIGZvciBwYXlsb2FkcyByZXByZXNlbnRpbmcgdG9ycmVudCBtYWduZXQgbGlua3MuXG4gKi9cbmV4cG9ydCBjbGFzcyBNQUdORVRQYXlsb2FkIGV4dGVuZHMgVVRGOFBheWxvYWQge1xuICBwcm90ZWN0ZWQgdHlwZWlkID0gMzBcbn1cbiJdfQ==