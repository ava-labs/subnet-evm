/**
 * @packageDocumentation
 * @module Utils-Payload
 */
import { Buffer } from "buffer/";
import BN from "bn.js";
/**
 * Class for determining payload types and managing the lookup table.
 */
export declare class PayloadTypes {
    private static instance;
    protected types: string[];
    /**
     * Given an encoded payload buffer returns the payload content (minus typeID).
     */
    getContent(payload: Buffer): Buffer;
    /**
     * Given an encoded payload buffer returns the payload (with typeID).
     */
    getPayload(payload: Buffer): Buffer;
    /**
     * Given a payload buffer returns the proper TypeID.
     */
    getTypeID(payload: Buffer): number;
    /**
     * Given a type string returns the proper TypeID.
     */
    lookupID(typestr: string): number;
    /**
     * Given a TypeID returns a string describing the payload type.
     */
    lookupType(value: number): string;
    /**
     * Given a TypeID returns the proper [[PayloadBase]].
     */
    select(typeID: number, ...args: any[]): PayloadBase;
    /**
     * Given a [[PayloadBase]] which may not be cast properly, returns a properly cast [[PayloadBase]].
     */
    recast(unknowPayload: PayloadBase): PayloadBase;
    /**
     * Returns the [[PayloadTypes]] singleton.
     */
    static getInstance(): PayloadTypes;
    private constructor();
}
/**
 * Base class for payloads.
 */
export declare abstract class PayloadBase {
    protected payload: Buffer;
    protected typeid: number;
    /**
     * Returns the TypeID for the payload.
     */
    typeID(): number;
    /**
     * Returns the string name for the payload's type.
     */
    typeName(): string;
    /**
     * Returns the payload content (minus typeID).
     */
    getContent(): Buffer;
    /**
     * Returns the payload (with typeID).
     */
    getPayload(): Buffer;
    /**
     * Decodes the payload as a {@link https://github.com/feross/buffer|Buffer} including 4 bytes for the length and TypeID.
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Encodes the payload as a {@link https://github.com/feross/buffer|Buffer} including 4 bytes for the length and TypeID.
     */
    toBuffer(): Buffer;
    /**
     * Returns the expected type for the payload.
     */
    abstract returnType(...args: any): any;
    constructor();
}
/**
 * Class for payloads representing simple binary blobs.
 */
export declare class BINPayload extends PayloadBase {
    protected typeid: number;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} for the payload.
     */
    returnType(): Buffer;
    /**
     * @param payload Buffer only
     */
    constructor(payload?: any);
}
/**
 * Class for payloads representing UTF8 encoding.
 */
export declare class UTF8Payload extends PayloadBase {
    protected typeid: number;
    /**
     * Returns a string for the payload.
     */
    returnType(): string;
    /**
     * @param payload Buffer utf8 string
     */
    constructor(payload?: any);
}
/**
 * Class for payloads representing Hexadecimal encoding.
 */
export declare class HEXSTRPayload extends PayloadBase {
    protected typeid: number;
    /**
     * Returns a hex string for the payload.
     */
    returnType(): string;
    /**
     * @param payload Buffer or hex string
     */
    constructor(payload?: any);
}
/**
 * Class for payloads representing Base58 encoding.
 */
export declare class B58STRPayload extends PayloadBase {
    protected typeid: number;
    /**
     * Returns a base58 string for the payload.
     */
    returnType(): string;
    /**
     * @param payload Buffer or cb58 encoded string
     */
    constructor(payload?: any);
}
/**
 * Class for payloads representing Base64 encoding.
 */
export declare class B64STRPayload extends PayloadBase {
    protected typeid: number;
    /**
     * Returns a base64 string for the payload.
     */
    returnType(): string;
    /**
     * @param payload Buffer of base64 string
     */
    constructor(payload?: any);
}
/**
 * Class for payloads representing Big Numbers.
 *
 * @param payload Accepts a Buffer, BN, or base64 string
 */
export declare class BIGNUMPayload extends PayloadBase {
    protected typeid: number;
    /**
     * Returns a {@link https://github.com/indutny/bn.js/|BN} for the payload.
     */
    returnType(): BN;
    /**
     * @param payload Buffer, BN, or base64 string
     */
    constructor(payload?: any);
}
/**
 * Class for payloads representing chain addresses.
 *
 */
export declare abstract class ChainAddressPayload extends PayloadBase {
    protected typeid: number;
    protected chainid: string;
    /**
     * Returns the chainid.
     */
    returnChainID(): string;
    /**
     * Returns an address string for the payload.
     */
    returnType(hrp: string): string;
    /**
     * @param payload Buffer or address string
     */
    constructor(payload?: any, hrp?: string);
}
/**
 * Class for payloads representing X-Chin addresses.
 */
export declare class XCHAINADDRPayload extends ChainAddressPayload {
    protected typeid: number;
    protected chainid: string;
}
/**
 * Class for payloads representing P-Chain addresses.
 */
export declare class PCHAINADDRPayload extends ChainAddressPayload {
    protected typeid: number;
    protected chainid: string;
}
/**
 * Class for payloads representing C-Chain addresses.
 */
export declare class CCHAINADDRPayload extends ChainAddressPayload {
    protected typeid: number;
    protected chainid: string;
}
/**
 * Class for payloads representing data serialized by bintools.cb58Encode().
 */
export declare abstract class cb58EncodedPayload extends PayloadBase {
    /**
     * Returns a bintools.cb58Encoded string for the payload.
     */
    returnType(): string;
    /**
     * @param payload Buffer or cb58 encoded string
     */
    constructor(payload?: any);
}
/**
 * Class for payloads representing TxIDs.
 */
export declare class TXIDPayload extends cb58EncodedPayload {
    protected typeid: number;
}
/**
 * Class for payloads representing AssetIDs.
 */
export declare class ASSETIDPayload extends cb58EncodedPayload {
    protected typeid: number;
}
/**
 * Class for payloads representing NODEIDs.
 */
export declare class UTXOIDPayload extends cb58EncodedPayload {
    protected typeid: number;
}
/**
 * Class for payloads representing NFTIDs (UTXOIDs in an NFT context).
 */
export declare class NFTIDPayload extends UTXOIDPayload {
    protected typeid: number;
}
/**
 * Class for payloads representing SubnetIDs.
 */
export declare class SUBNETIDPayload extends cb58EncodedPayload {
    protected typeid: number;
}
/**
 * Class for payloads representing ChainIDs.
 */
export declare class CHAINIDPayload extends cb58EncodedPayload {
    protected typeid: number;
}
/**
 * Class for payloads representing NodeIDs.
 */
export declare class NODEIDPayload extends cb58EncodedPayload {
    protected typeid: number;
}
/**
 * Class for payloads representing secp256k1 signatures.
 * convention: secp256k1 signature (130 bytes)
 */
export declare class SECPSIGPayload extends B58STRPayload {
    protected typeid: number;
}
/**
 * Class for payloads representing secp256k1 encrypted messages.
 * convention: public key (65 bytes) + secp256k1 encrypted message for that public key
 */
export declare class SECPENCPayload extends B58STRPayload {
    protected typeid: number;
}
/**
 * Class for payloads representing JPEG images.
 */
export declare class JPEGPayload extends BINPayload {
    protected typeid: number;
}
export declare class PNGPayload extends BINPayload {
    protected typeid: number;
}
/**
 * Class for payloads representing BMP images.
 */
export declare class BMPPayload extends BINPayload {
    protected typeid: number;
}
/**
 * Class for payloads representing ICO images.
 */
export declare class ICOPayload extends BINPayload {
    protected typeid: number;
}
/**
 * Class for payloads representing SVG images.
 */
export declare class SVGPayload extends UTF8Payload {
    protected typeid: number;
}
/**
 * Class for payloads representing CSV files.
 */
export declare class CSVPayload extends UTF8Payload {
    protected typeid: number;
}
/**
 * Class for payloads representing JSON strings.
 */
export declare class JSONPayload extends PayloadBase {
    protected typeid: number;
    /**
     * Returns a JSON-decoded object for the payload.
     */
    returnType(): any;
    constructor(payload?: any);
}
/**
 * Class for payloads representing YAML definitions.
 */
export declare class YAMLPayload extends UTF8Payload {
    protected typeid: number;
}
/**
 * Class for payloads representing email addresses.
 */
export declare class EMAILPayload extends UTF8Payload {
    protected typeid: number;
}
/**
 * Class for payloads representing URL strings.
 */
export declare class URLPayload extends UTF8Payload {
    protected typeid: number;
}
/**
 * Class for payloads representing IPFS addresses.
 */
export declare class IPFSPayload extends B58STRPayload {
    protected typeid: number;
}
/**
 * Class for payloads representing onion URLs.
 */
export declare class ONIONPayload extends UTF8Payload {
    protected typeid: number;
}
/**
 * Class for payloads representing torrent magnet links.
 */
export declare class MAGNETPayload extends UTF8Payload {
    protected typeid: number;
}
//# sourceMappingURL=payload.d.ts.map