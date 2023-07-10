/**
 * @packageDocumentation
 * @module Common-AssetAmount
 */
import { Buffer } from "buffer/";
import BN from "bn.js";
import { StandardTransferableOutput } from "./output";
import { StandardTransferableInput } from "./input";
/**
 * Class for managing asset amounts in the UTXOSet fee calcuation
 */
export declare class AssetAmount {
    protected assetID: Buffer;
    protected amount: BN;
    protected burn: BN;
    protected spent: BN;
    protected stakeableLockSpent: BN;
    protected change: BN;
    protected stakeableLockChange: boolean;
    protected finished: boolean;
    getAssetID: () => Buffer;
    getAssetIDString: () => string;
    getAmount: () => BN;
    getSpent: () => BN;
    getBurn: () => BN;
    getChange: () => BN;
    getStakeableLockSpent: () => BN;
    getStakeableLockChange: () => boolean;
    isFinished: () => boolean;
    spendAmount: (amt: BN, stakeableLocked?: boolean) => boolean;
    constructor(assetID: Buffer, amount: BN, burn: BN);
}
export declare abstract class StandardAssetAmountDestination<TO extends StandardTransferableOutput, TI extends StandardTransferableInput> {
    protected amounts: AssetAmount[];
    protected destinations: Buffer[];
    protected senders: Buffer[];
    protected changeAddresses: Buffer[];
    protected amountkey: object;
    protected inputs: TI[];
    protected outputs: TO[];
    protected change: TO[];
    addAssetAmount: (assetID: Buffer, amount: BN, burn: BN) => void;
    addInput: (input: TI) => void;
    addOutput: (output: TO) => void;
    addChange: (output: TO) => void;
    getAmounts: () => AssetAmount[];
    getDestinations: () => Buffer[];
    getSenders: () => Buffer[];
    getChangeAddresses: () => Buffer[];
    getAssetAmount: (assetHexStr: string) => AssetAmount;
    assetExists: (assetHexStr: string) => boolean;
    getInputs: () => TI[];
    getOutputs: () => TO[];
    getChangeOutputs: () => TO[];
    getAllOutputs: () => TO[];
    canComplete: () => boolean;
    constructor(destinations: Buffer[], senders: Buffer[], changeAddresses: Buffer[]);
}
//# sourceMappingURL=assetamount.d.ts.map