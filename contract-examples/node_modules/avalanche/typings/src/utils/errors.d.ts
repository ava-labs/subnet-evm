export declare class AvalancheError extends Error {
    errorCode: string;
    constructor(m: string, code: string);
    getCode(): string;
}
export declare class AddressError extends AvalancheError {
    constructor(m: string);
}
export declare class GooseEggCheckError extends AvalancheError {
    constructor(m: string);
}
export declare class ChainIdError extends AvalancheError {
    constructor(m: string);
}
export declare class NoAtomicUTXOsError extends AvalancheError {
    constructor(m: string);
}
export declare class SymbolError extends AvalancheError {
    constructor(m: string);
}
export declare class NameError extends AvalancheError {
    constructor(m: string);
}
export declare class TransactionError extends AvalancheError {
    constructor(m: string);
}
export declare class CodecIdError extends AvalancheError {
    constructor(m: string);
}
export declare class CredIdError extends AvalancheError {
    constructor(m: string);
}
export declare class TransferableOutputError extends AvalancheError {
    constructor(m: string);
}
export declare class TransferableInputError extends AvalancheError {
    constructor(m: string);
}
export declare class InputIdError extends AvalancheError {
    constructor(m: string);
}
export declare class OperationError extends AvalancheError {
    constructor(m: string);
}
export declare class InvalidOperationIdError extends AvalancheError {
    constructor(m: string);
}
export declare class ChecksumError extends AvalancheError {
    constructor(m: string);
}
export declare class OutputIdError extends AvalancheError {
    constructor(m: string);
}
export declare class UTXOError extends AvalancheError {
    constructor(m: string);
}
export declare class InsufficientFundsError extends AvalancheError {
    constructor(m: string);
}
export declare class ThresholdError extends AvalancheError {
    constructor(m: string);
}
export declare class SECPMintOutputError extends AvalancheError {
    constructor(m: string);
}
export declare class EVMInputError extends AvalancheError {
    constructor(m: string);
}
export declare class EVMOutputError extends AvalancheError {
    constructor(m: string);
}
export declare class FeeAssetError extends AvalancheError {
    constructor(m: string);
}
export declare class StakeError extends AvalancheError {
    constructor(m: string);
}
export declare class TimeError extends AvalancheError {
    constructor(m: string);
}
export declare class DelegationFeeError extends AvalancheError {
    constructor(m: string);
}
export declare class SubnetOwnerError extends AvalancheError {
    constructor(m: string);
}
export declare class BufferSizeError extends AvalancheError {
    constructor(m: string);
}
export declare class AddressIndexError extends AvalancheError {
    constructor(m: string);
}
export declare class PublicKeyError extends AvalancheError {
    constructor(m: string);
}
export declare class MergeRuleError extends AvalancheError {
    constructor(m: string);
}
export declare class Base58Error extends AvalancheError {
    constructor(m: string);
}
export declare class PrivateKeyError extends AvalancheError {
    constructor(m: string);
}
export declare class NodeIdError extends AvalancheError {
    constructor(m: string);
}
export declare class HexError extends AvalancheError {
    constructor(m: string);
}
export declare class TypeIdError extends AvalancheError {
    constructor(m: string);
}
export declare class TypeNameError extends AvalancheError {
    constructor(m: string);
}
export declare class UnknownTypeError extends AvalancheError {
    constructor(m: string);
}
export declare class Bech32Error extends AvalancheError {
    constructor(m: string);
}
export declare class EVMFeeError extends AvalancheError {
    constructor(m: string);
}
export declare class InvalidEntropy extends AvalancheError {
    constructor(m: string);
}
export declare class ProtocolError extends AvalancheError {
    constructor(m: string);
}
export declare class SubnetIdError extends AvalancheError {
    constructor(m: string);
}
export declare class SubnetThresholdError extends AvalancheError {
    constructor(m: string);
}
export declare class SubnetAddressError extends AvalancheError {
    constructor(m: string);
}
export interface ErrorResponseObject {
    code: number;
    message: string;
    data?: null;
}
//# sourceMappingURL=errors.d.ts.map