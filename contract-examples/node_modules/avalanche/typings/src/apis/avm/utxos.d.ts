/**
 * @packageDocumentation
 * @module API-AVM-UTXOs
 */
import { Buffer } from "buffer/";
import BN from "bn.js";
import { TransferableOutput, SECPMintOutput, SECPTransferOutput } from "./outputs";
import { UnsignedTx } from "./tx";
import { TransferableInput } from "./inputs";
import { Output, OutputOwners } from "../../common/output";
import { InitialStates } from "./initialstates";
import { MinterSet } from "./minterset";
import { StandardUTXO, StandardUTXOSet } from "../../common/utxos";
import { StandardAssetAmountDestination } from "../../common/assetamount";
import { SerializedEncoding } from "../../utils/serialization";
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
    create(codecID?: number, txid?: Buffer, outputidx?: Buffer | number, assetID?: Buffer, output?: Output): this;
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
    create(...args: any[]): this;
    clone(): this;
    _feeCheck(fee: BN, feeAssetID: Buffer): boolean;
    getMinimumSpendable: (aad: AssetAmountDestination, asOf?: BN, locktime?: BN, threshold?: number) => Error;
    /**
     * Creates an [[UnsignedTx]] wrapping a [[BaseTx]]. For more granular control, you may create your own
     * [[UnsignedTx]] wrapping a [[BaseTx]] manually (with their corresponding [[TransferableInput]]s and [[TransferableOutput]]s).
     *
     * @param networkID The number representing NetworkID of the node
     * @param blockchainID The {@link https://github.com/feross/buffer|Buffer} representing the BlockchainID for the transaction
     * @param amount The amount of the asset to be spent in its smallest denomination, represented as {@link https://github.com/indutny/bn.js/|BN}.
     * @param assetID {@link https://github.com/feross/buffer|Buffer} of the asset ID for the UTXO
     * @param toAddresses The addresses to send the funds
     * @param fromAddresses The addresses being used to send the funds from the UTXOs {@link https://github.com/feross/buffer|Buffer}
     * @param changeAddresses Optional. The addresses that can spend the change remaining from the spent UTXOs. Default: toAddresses
     * @param fee Optional. The amount of fees to burn in its smallest denomination, represented as {@link https://github.com/indutny/bn.js/|BN}
     * @param feeAssetID Optional. The assetID of the fees being burned. Default: assetID
     * @param memo Optional. Contains arbitrary data, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     * @param locktime Optional. The locktime field created in the resulting outputs
     * @param threshold Optional. The number of signatures required to spend the funds in the resultant UTXO
     *
     * @returns An unsigned transaction created from the passed in parameters.
     *
     */
    buildBaseTx: (networkID: number, blockchainID: Buffer, amount: BN, assetID: Buffer, toAddresses: Buffer[], fromAddresses: Buffer[], changeAddresses?: Buffer[], fee?: BN, feeAssetID?: Buffer, memo?: Buffer, asOf?: BN, locktime?: BN, threshold?: number) => UnsignedTx;
    /**
     * Creates an unsigned Create Asset transaction. For more granular control, you may create your own
     * [[CreateAssetTX]] manually (with their corresponding [[TransferableInput]]s, [[TransferableOutput]]s).
     *
     * @param networkID The number representing NetworkID of the node
     * @param blockchainID The {@link https://github.com/feross/buffer|Buffer} representing the BlockchainID for the transaction
     * @param fromAddresses The addresses being used to send the funds from the UTXOs {@link https://github.com/feross/buffer|Buffer}
     * @param changeAddresses Optional. The addresses that can spend the change remaining from the spent UTXOs
     * @param initialState The [[InitialStates]] that represent the intial state of a created asset
     * @param name String for the descriptive name of the asset
     * @param symbol String for the ticker symbol of the asset
     * @param denomination Optional number for the denomination which is 10^D. D must be >= 0 and <= 32. Ex: $1 AVAX = 10^9 $nAVAX
     * @param mintOutputs Optional. Array of [[SECPMintOutput]]s to be included in the transaction. These outputs can be spent to mint more tokens.
     * @param fee Optional. The amount of fees to burn in its smallest denomination, represented as {@link https://github.com/indutny/bn.js/|BN}
     * @param feeAssetID Optional. The assetID of the fees being burned.
     * @param memo Optional contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     *
     * @returns An unsigned transaction created from the passed in parameters.
     *
     */
    buildCreateAssetTx: (networkID: number, blockchainID: Buffer, fromAddresses: Buffer[], changeAddresses: Buffer[], initialState: InitialStates, name: string, symbol: string, denomination: number, mintOutputs?: SECPMintOutput[], fee?: BN, feeAssetID?: Buffer, memo?: Buffer, asOf?: BN) => UnsignedTx;
    /**
     * Creates an unsigned Secp mint transaction. For more granular control, you may create your own
     * [[OperationTx]] manually (with their corresponding [[TransferableInput]]s, [[TransferableOutput]]s, and [[TransferOperation]]s).
     *
     * @param networkID The number representing NetworkID of the node
     * @param blockchainID The {@link https://github.com/feross/buffer|Buffer} representing the BlockchainID for the transaction
     * @param mintOwner A [[SECPMintOutput]] which specifies the new set of minters
     * @param transferOwner A [[SECPTransferOutput]] which specifies where the minted tokens will go
     * @param fromAddresses The addresses being used to send the funds from the UTXOs {@link https://github.com/feross/buffer|Buffer}
     * @param changeAddresses The addresses that can spend the change remaining from the spent UTXOs
     * @param mintUTXOID The UTXOID for the [[SCPMintOutput]] being spent to produce more tokens
     * @param fee Optional. The amount of fees to burn in its smallest denomination, represented as {@link https://github.com/indutny/bn.js/|BN}
     * @param feeAssetID Optional. The assetID of the fees being burned.
     * @param memo Optional contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     */
    buildSECPMintTx: (networkID: number, blockchainID: Buffer, mintOwner: SECPMintOutput, transferOwner: SECPTransferOutput, fromAddresses: Buffer[], changeAddresses: Buffer[], mintUTXOID: string, fee?: BN, feeAssetID?: Buffer, memo?: Buffer, asOf?: BN) => UnsignedTx;
    /**
     * Creates an unsigned Create Asset transaction. For more granular control, you may create your own
     * [[CreateAssetTX]] manually (with their corresponding [[TransferableInput]]s, [[TransferableOutput]]s).
     *
     * @param networkID The number representing NetworkID of the node
     * @param blockchainID The {@link https://github.com/feross/buffer|Buffer} representing the BlockchainID for the transaction
     * @param fromAddresses The addresses being used to send the funds from the UTXOs {@link https://github.com/feross/buffer|Buffer}
     * @param changeAddresses Optional. The addresses that can spend the change remaining from the spent UTXOs.
     * @param minterSets The minters and thresholds required to mint this nft asset
     * @param name String for the descriptive name of the nft asset
     * @param symbol String for the ticker symbol of the nft asset
     * @param fee Optional. The amount of fees to burn in its smallest denomination, represented as {@link https://github.com/indutny/bn.js/|BN}
     * @param feeAssetID Optional. The assetID of the fees being burned.
     * @param memo Optional contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     * @param locktime Optional. The locktime field created in the resulting mint output
     *
     * @returns An unsigned transaction created from the passed in parameters.
     *
     */
    buildCreateNFTAssetTx: (networkID: number, blockchainID: Buffer, fromAddresses: Buffer[], changeAddresses: Buffer[], minterSets: MinterSet[], name: string, symbol: string, fee?: BN, feeAssetID?: Buffer, memo?: Buffer, asOf?: BN, locktime?: BN) => UnsignedTx;
    /**
     * Creates an unsigned NFT mint transaction. For more granular control, you may create your own
     * [[OperationTx]] manually (with their corresponding [[TransferableInput]]s, [[TransferableOutput]]s, and [[TransferOperation]]s).
     *
     * @param networkID The number representing NetworkID of the node
     * @param blockchainID The {@link https://github.com/feross/buffer|Buffer} representing the BlockchainID for the transaction
     * @param owners An array of [[OutputOwners]] who will be given the NFTs.
     * @param fromAddresses The addresses being used to send the funds from the UTXOs
     * @param changeAddresses Optional. The addresses that can spend the change remaining from the spent UTXOs.
     * @param utxoids An array of strings for the NFTs being transferred
     * @param groupID Optional. The group this NFT is issued to.
     * @param payload Optional. Data for NFT Payload.
     * @param fee Optional. The amount of fees to burn in its smallest denomination, represented as {@link https://github.com/indutny/bn.js/|BN}
     * @param feeAssetID Optional. The assetID of the fees being burned.
     * @param memo Optional contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     *
     * @returns An unsigned transaction created from the passed in parameters.
     *
     */
    buildCreateNFTMintTx: (networkID: number, blockchainID: Buffer, owners: OutputOwners[], fromAddresses: Buffer[], changeAddresses: Buffer[], utxoids: string[], groupID?: number, payload?: Buffer, fee?: BN, feeAssetID?: Buffer, memo?: Buffer, asOf?: BN) => UnsignedTx;
    /**
     * Creates an unsigned NFT transfer transaction. For more granular control, you may create your own
     * [[OperationTx]] manually (with their corresponding [[TransferableInput]]s, [[TransferableOutput]]s, and [[TransferOperation]]s).
     *
     * @param networkID The number representing NetworkID of the node
     * @param blockchainID The {@link https://github.com/feross/buffer|Buffer} representing the BlockchainID for the transaction
     * @param toAddresses An array of {@link https://github.com/feross/buffer|Buffer}s which indicate who recieves the NFT
     * @param fromAddresses An array for {@link https://github.com/feross/buffer|Buffer} who owns the NFT
     * @param changeAddresses Optional. The addresses that can spend the change remaining from the spent UTXOs.
     * @param utxoids An array of strings for the NFTs being transferred
     * @param fee Optional. The amount of fees to burn in its smallest denomination, represented as {@link https://github.com/indutny/bn.js/|BN}
     * @param feeAssetID Optional. The assetID of the fees being burned.
     * @param memo Optional contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     * @param locktime Optional. The locktime field created in the resulting outputs
     * @param threshold Optional. The number of signatures required to spend the funds in the resultant UTXO
     *
     * @returns An unsigned transaction created from the passed in parameters.
     *
     */
    buildNFTTransferTx: (networkID: number, blockchainID: Buffer, toAddresses: Buffer[], fromAddresses: Buffer[], changeAddresses: Buffer[], utxoids: string[], fee?: BN, feeAssetID?: Buffer, memo?: Buffer, asOf?: BN, locktime?: BN, threshold?: number) => UnsignedTx;
    /**
     * Creates an unsigned ImportTx transaction.
     *
     * @param networkID The number representing NetworkID of the node
     * @param blockchainID The {@link https://github.com/feross/buffer|Buffer} representing the BlockchainID for the transaction
     * @param toAddresses The addresses to send the funds
     * @param fromAddresses The addresses being used to send the funds from the UTXOs {@link https://github.com/feross/buffer|Buffer}
     * @param changeAddresses Optional. The addresses that can spend the change remaining from the spent UTXOs.
     * @param importIns An array of [[TransferableInput]]s being imported
     * @param sourceChain A {@link https://github.com/feross/buffer|Buffer} for the chainid where the imports are coming from.
     * @param fee Optional. The amount of fees to burn in its smallest denomination, represented as {@link https://github.com/indutny/bn.js/|BN}. Fee will come from the inputs first, if they can.
     * @param feeAssetID Optional. The assetID of the fees being burned.
     * @param memo Optional contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     * @param locktime Optional. The locktime field created in the resulting outputs
     * @param threshold Optional. The number of signatures required to spend the funds in the resultant UTXO
     * @returns An unsigned transaction created from the passed in parameters.
     *
     */
    buildImportTx: (networkID: number, blockchainID: Buffer, toAddresses: Buffer[], fromAddresses: Buffer[], changeAddresses: Buffer[], atomics: UTXO[], sourceChain?: Buffer, fee?: BN, feeAssetID?: Buffer, memo?: Buffer, asOf?: BN, locktime?: BN, threshold?: number) => UnsignedTx;
    /**
     * Creates an unsigned ExportTx transaction.
     *
     * @param networkID The number representing NetworkID of the node
     * @param blockchainID The {@link https://github.com/feross/buffer|Buffer} representing the BlockchainID for the transaction
     * @param amount The amount being exported as a {@link https://github.com/indutny/bn.js/|BN}
     * @param avaxAssetID {@link https://github.com/feross/buffer|Buffer} of the asset ID for AVAX
     * @param toAddresses An array of addresses as {@link https://github.com/feross/buffer|Buffer} who recieves the AVAX
     * @param fromAddresses An array of addresses as {@link https://github.com/feross/buffer|Buffer} who owns the AVAX
     * @param changeAddresses Optional. The addresses that can spend the change remaining from the spent UTXOs.
     * @param fee Optional. The amount of fees to burn in its smallest denomination, represented as {@link https://github.com/indutny/bn.js/|BN}
     * @param destinationChain Optional. A {@link https://github.com/feross/buffer|Buffer} for the chainid where to send the asset.
     * @param feeAssetID Optional. The assetID of the fees being burned.
     * @param memo Optional contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     * @param locktime Optional. The locktime field created in the resulting outputs
     * @param threshold Optional. The number of signatures required to spend the funds in the resultant UTXO
     * @returns An unsigned transaction created from the passed in parameters.
     *
     */
    buildExportTx: (networkID: number, blockchainID: Buffer, amount: BN, assetID: Buffer, toAddresses: Buffer[], fromAddresses: Buffer[], changeAddresses?: Buffer[], destinationChain?: Buffer, fee?: BN, feeAssetID?: Buffer, memo?: Buffer, asOf?: BN, locktime?: BN, threshold?: number) => UnsignedTx;
}
//# sourceMappingURL=utxos.d.ts.map