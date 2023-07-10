/**
 * @packageDocumentation
 * @module Utils-HelperFunctions
 */
import BN from "bn.js";
import { Buffer } from "buffer/";
import { UnsignedTx } from "../apis/evm";
export declare function getPreferredHRP(networkID?: number): string;
export declare function MaxWeightFormula(staked: BN, cap: BN): BN;
/**
 * Function providing the current UNIX time using a {@link https://github.com/indutny/bn.js/|BN}.
 */
export declare function UnixNow(): BN;
/**
 * Takes a private key buffer and produces a private key string with prefix.
 *
 * @param pk A {@link https://github.com/feross/buffer|Buffer} for the private key.
 */
export declare function bufferToPrivateKeyString(pk: Buffer): string;
/**
 * Takes a private key string and produces a private key {@link https://github.com/feross/buffer|Buffer}.
 *
 * @param pk A string for the private key.
 */
export declare function privateKeyStringToBuffer(pk: string): Buffer;
/**
 * Takes a nodeID buffer and produces a nodeID string with prefix.
 *
 * @param pk A {@link https://github.com/feross/buffer|Buffer} for the nodeID.
 */
export declare function bufferToNodeIDString(pk: Buffer): string;
/**
 * Takes a nodeID string and produces a nodeID {@link https://github.com/feross/buffer|Buffer}.
 *
 * @param pk A string for the nodeID.
 */
export declare function NodeIDStringToBuffer(pk: string): Buffer;
export declare function costImportTx(tx: UnsignedTx): number;
export declare function calcBytesCost(len: number): number;
export declare function costExportTx(tx: UnsignedTx): number;
//# sourceMappingURL=helperfunctions.d.ts.map