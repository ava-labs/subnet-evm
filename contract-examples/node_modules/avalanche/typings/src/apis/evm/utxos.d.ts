/**
 * @packageDocumentation
 * @module API-EVM-UTXOs
 */
import { Buffer } from "buffer/";
import BN from "bn.js";
import { TransferableOutput } from "./outputs";
import { TransferableInput } from "./inputs";
import { Output } from "../../common/output";
import { StandardUTXO, StandardUTXOSet } from "../../common/utxos";
import { StandardAssetAmountDestination } from "../../common/assetamount";
import { SerializedEncoding } from "../../utils/serialization";
import { UnsignedTx } from "./tx";
/**
 * Class for representing a single UTXO.
 */
export declare class UTXO extends StandardUTXO {
    protected _typeName: string;
    protected _typeID: any;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Takes a base-58 string containing a [[UTXO]], parses it, populates the class, and returns the length of the StandardUTXO in bytes.
     *
     * @param serialized A base-58 string containing a raw [[UTXO]]
     *
     * @returns The length of the raw [[UTXO]]
     *
     * @remarks
     * unlike most fromStrings, it expects the string to be serialized in cb58 format
     */
    fromString(serialized: string): number;
    /**
     * Returns a base-58 representation of the [[UTXO]].
     *
     * @remarks
     * unlike most toStrings, this returns in cb58 serialization format
     */
    toString(): string;
    clone(): this;
    create(codecID?: number, txID?: Buffer, outputidx?: Buffer | number, assetID?: Buffer, output?: Output): this;
}
export declare class AssetAmountDestination extends StandardAssetAmountDestination<TransferableOutput, TransferableInput> {
}
/**
 * Class representing a set of [[UTXO]]s.
 */
export declare class UTXOSet extends StandardUTXOSet<UTXO> {
    protected _typeName: string;
    protected _typeID: any;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    parseUTXO(utxo: UTXO | string): UTXO;
    create(): this;
    clone(): this;
    _feeCheck(fee: BN, feeAssetID: Buffer): boolean;
    getMinimumSpendable: (aad: AssetAmountDestination, asOf?: BN, locktime?: BN, threshold?: number) => Error;
    /**
     * Creates an unsigned ImportTx transaction.
     *
     * @param networkID The number representing NetworkID of the node
     * @param blockchainID The {@link https://github.com/feross/buffer|Buffer} representing the BlockchainID for the transaction
     * @param toAddress The address to send the funds
     * @param importIns An array of [[TransferableInput]]s being imported
     * @param sourceChain A {@link https://github.com/feross/buffer|Buffer} for the chainid where the imports are coming from.
     * @param fee Optional. The amount of fees to burn in its smallest denomination, represented as {@link https://github.com/indutny/bn.js/|BN}. Fee will come from the inputs first, if they can.
     * @param feeAssetID Optional. The assetID of the fees being burned.
     * @returns An unsigned transaction created from the passed in parameters.
     *
     */
    buildImportTx: (networkID: number, blockchainID: Buffer, toAddress: string, atomics: UTXO[], sourceChain?: Buffer, fee?: BN, feeAssetID?: Buffer) => UnsignedTx;
    /**
     * Creates an unsigned ExportTx transaction.
     *
     * @param networkID The number representing NetworkID of the node
     * @param blockchainID The {@link https://github.com/feross/buffer|Buffer} representing the BlockchainID for the transaction
     * @param amount The amount being exported as a {@link https://github.com/indutny/bn.js/|BN}
     * @param avaxAssetID {@link https://github.com/feross/buffer|Buffer} of the AssetID for AVAX
     * @param toAddresses An array of addresses as {@link https://github.com/feross/buffer|Buffer} who recieves the AVAX
     * @param fromAddresses An array of addresses as {@link https://github.com/feross/buffer|Buffer} who owns the AVAX
     * @param changeAddresses Optional. The addresses that can spend the change remaining from the spent UTXOs.
     * @param destinationChain Optional. A {@link https://github.com/feross/buffer|Buffer} for the chainid where to send the asset.
     * @param fee Optional. The amount of fees to burn in its smallest denomination, represented as {@link https://github.com/indutny/bn.js/|BN}
     * @param feeAssetID Optional. The assetID of the fees being burned.
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     * @param locktime Optional. The locktime field created in the resulting outputs
     * @param threshold Optional. The number of signatures required to spend the funds in the resultant UTXO
     * @returns An unsigned transaction created from the passed in parameters.
     *
     */
    buildExportTx: (networkID: number, blockchainID: Buffer, amount: BN, avaxAssetID: Buffer, toAddresses: Buffer[], fromAddresses: Buffer[], changeAddresses?: Buffer[], destinationChain?: Buffer, fee?: BN, feeAssetID?: Buffer, asOf?: BN, locktime?: BN, threshold?: number) => UnsignedTx;
}
//# sourceMappingURL=utxos.d.ts.map