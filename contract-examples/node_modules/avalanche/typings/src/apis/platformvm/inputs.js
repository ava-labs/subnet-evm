"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.StakeableLockIn = exports.SECPTransferInput = exports.AmountInput = exports.TransferableInput = exports.ParseableInput = exports.SelectInputClass = void 0;
/**
 * @packageDocumentation
 * @module API-PlatformVM-Inputs
 */
const buffer_1 = require("buffer/");
const bintools_1 = __importDefault(require("../../utils/bintools"));
const constants_1 = require("./constants");
const input_1 = require("../../common/input");
const serialization_1 = require("../../utils/serialization");
const errors_1 = require("../../utils/errors");
/**
 * @ignore
 */
const bintools = bintools_1.default.getInstance();
const serialization = serialization_1.Serialization.getInstance();
/**
 * Takes a buffer representing the output and returns the proper [[Input]] instance.
 *
 * @param inputid A number representing the inputID parsed prior to the bytes passed in
 *
 * @returns An instance of an [[Input]]-extended class.
 */
const SelectInputClass = (inputid, ...args) => {
    if (inputid === constants_1.PlatformVMConstants.SECPINPUTID) {
        return new SECPTransferInput(...args);
    }
    else if (inputid === constants_1.PlatformVMConstants.STAKEABLELOCKINID) {
        return new StakeableLockIn(...args);
    }
    /* istanbul ignore next */
    throw new errors_1.InputIdError("Error - SelectInputClass: unknown inputid");
};
exports.SelectInputClass = SelectInputClass;
class ParseableInput extends input_1.StandardParseableInput {
    constructor() {
        super(...arguments);
        this._typeName = "ParseableInput";
        this._typeID = undefined;
    }
    //serialize is inherited
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.input = (0, exports.SelectInputClass)(fields["input"]["_typeID"]);
        this.input.deserialize(fields["input"], encoding);
    }
    fromBuffer(bytes, offset = 0) {
        const inputid = bintools
            .copyFrom(bytes, offset, offset + 4)
            .readUInt32BE(0);
        offset += 4;
        this.input = (0, exports.SelectInputClass)(inputid);
        return this.input.fromBuffer(bytes, offset);
    }
}
exports.ParseableInput = ParseableInput;
class TransferableInput extends input_1.StandardTransferableInput {
    constructor() {
        super(...arguments);
        this._typeName = "TransferableInput";
        this._typeID = undefined;
    }
    //serialize is inherited
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.input = (0, exports.SelectInputClass)(fields["input"]["_typeID"]);
        this.input.deserialize(fields["input"], encoding);
    }
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} containing a [[TransferableInput]], parses it, populates the class, and returns the length of the [[TransferableInput]] in bytes.
     *
     * @param bytes A {@link https://github.com/feross/buffer|Buffer} containing a raw [[TransferableInput]]
     *
     * @returns The length of the raw [[TransferableInput]]
     */
    fromBuffer(bytes, offset = 0) {
        this.txid = bintools.copyFrom(bytes, offset, offset + 32);
        offset += 32;
        this.outputidx = bintools.copyFrom(bytes, offset, offset + 4);
        offset += 4;
        this.assetID = bintools.copyFrom(bytes, offset, offset + constants_1.PlatformVMConstants.ASSETIDLEN);
        offset += 32;
        const inputid = bintools
            .copyFrom(bytes, offset, offset + 4)
            .readUInt32BE(0);
        offset += 4;
        this.input = (0, exports.SelectInputClass)(inputid);
        return this.input.fromBuffer(bytes, offset);
    }
}
exports.TransferableInput = TransferableInput;
class AmountInput extends input_1.StandardAmountInput {
    constructor() {
        super(...arguments);
        this._typeName = "AmountInput";
        this._typeID = undefined;
    }
    //serialize and deserialize both are inherited
    select(id, ...args) {
        return (0, exports.SelectInputClass)(id, ...args);
    }
}
exports.AmountInput = AmountInput;
class SECPTransferInput extends AmountInput {
    constructor() {
        super(...arguments);
        this._typeName = "SECPTransferInput";
        this._typeID = constants_1.PlatformVMConstants.SECPINPUTID;
        this.getCredentialID = () => constants_1.PlatformVMConstants.SECPCREDENTIAL;
    }
    //serialize and deserialize both are inherited
    /**
     * Returns the inputID for this input
     */
    getInputID() {
        return this._typeID;
    }
    create(...args) {
        return new SECPTransferInput(...args);
    }
    clone() {
        const newout = this.create();
        newout.fromBuffer(this.toBuffer());
        return newout;
    }
}
exports.SECPTransferInput = SECPTransferInput;
/**
 * An [[Input]] class which specifies an input that has a locktime which can also enable staking of the value held, preventing transfers but not validation.
 */
class StakeableLockIn extends AmountInput {
    /**
     * A [[Output]] class which specifies an [[Input]] that has a locktime which can also enable staking of the value held, preventing transfers but not validation.
     *
     * @param amount A {@link https://github.com/indutny/bn.js/|BN} representing the amount in the input
     * @param stakeableLocktime A {@link https://github.com/indutny/bn.js/|BN} representing the stakeable locktime
     * @param transferableInput A [[ParseableInput]] which is embedded into this input.
     */
    constructor(amount = undefined, stakeableLocktime = undefined, transferableInput = undefined) {
        super(amount);
        this._typeName = "StakeableLockIn";
        this._typeID = constants_1.PlatformVMConstants.STAKEABLELOCKINID;
        this.getCredentialID = () => constants_1.PlatformVMConstants.SECPCREDENTIAL;
        if (typeof stakeableLocktime !== "undefined") {
            this.stakeableLocktime = bintools.fromBNToBuffer(stakeableLocktime, 8);
        }
        if (typeof transferableInput !== "undefined") {
            this.transferableInput = transferableInput;
            this.synchronize();
        }
    }
    //serialize and deserialize both are inherited
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        let outobj = Object.assign(Object.assign({}, fields), { stakeableLocktime: serialization.encoder(this.stakeableLocktime, encoding, "Buffer", "decimalString", 8), transferableInput: this.transferableInput.serialize(encoding) });
        delete outobj["sigIdxs"];
        delete outobj["sigCount"];
        delete outobj["amount"];
        return outobj;
    }
    deserialize(fields, encoding = "hex") {
        fields["sigIdxs"] = [];
        fields["sigCount"] = "0";
        fields["amount"] = "98";
        super.deserialize(fields, encoding);
        this.stakeableLocktime = serialization.decoder(fields["stakeableLocktime"], encoding, "decimalString", "Buffer", 8);
        this.transferableInput = new ParseableInput();
        this.transferableInput.deserialize(fields["transferableInput"], encoding);
        this.synchronize();
    }
    synchronize() {
        let input = this.transferableInput.getInput();
        this.sigIdxs = input.getSigIdxs();
        this.sigCount = buffer_1.Buffer.alloc(4);
        this.sigCount.writeUInt32BE(this.sigIdxs.length, 0);
        this.amount = bintools.fromBNToBuffer(input.getAmount(), 8);
        this.amountValue = input.getAmount();
    }
    getStakeableLocktime() {
        return bintools.fromBufferToBN(this.stakeableLocktime);
    }
    getTransferablInput() {
        return this.transferableInput;
    }
    /**
     * Returns the inputID for this input
     */
    getInputID() {
        return this._typeID;
    }
    /**
     * Popuates the instance from a {@link https://github.com/feross/buffer|Buffer} representing the [[StakeableLockIn]] and returns the size of the output.
     */
    fromBuffer(bytes, offset = 0) {
        this.stakeableLocktime = bintools.copyFrom(bytes, offset, offset + 8);
        offset += 8;
        this.transferableInput = new ParseableInput();
        offset = this.transferableInput.fromBuffer(bytes, offset);
        this.synchronize();
        return offset;
    }
    /**
     * Returns the buffer representing the [[StakeableLockIn]] instance.
     */
    toBuffer() {
        const xferinBuff = this.transferableInput.toBuffer();
        const bsize = this.stakeableLocktime.length + xferinBuff.length;
        const barr = [this.stakeableLocktime, xferinBuff];
        return buffer_1.Buffer.concat(barr, bsize);
    }
    create(...args) {
        return new StakeableLockIn(...args);
    }
    clone() {
        const newout = this.create();
        newout.fromBuffer(this.toBuffer());
        return newout;
    }
    select(id, ...args) {
        return (0, exports.SelectInputClass)(id, ...args);
    }
}
exports.StakeableLockIn = StakeableLockIn;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiaW5wdXRzLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2FwaXMvcGxhdGZvcm12bS9pbnB1dHMudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7O0FBQUE7OztHQUdHO0FBQ0gsb0NBQWdDO0FBQ2hDLG9FQUEyQztBQUMzQywyQ0FBaUQ7QUFDakQsOENBSzJCO0FBQzNCLDZEQUE2RTtBQUU3RSwrQ0FBaUQ7QUFFakQ7O0dBRUc7QUFDSCxNQUFNLFFBQVEsR0FBYSxrQkFBUSxDQUFDLFdBQVcsRUFBRSxDQUFBO0FBQ2pELE1BQU0sYUFBYSxHQUFrQiw2QkFBYSxDQUFDLFdBQVcsRUFBRSxDQUFBO0FBRWhFOzs7Ozs7R0FNRztBQUNJLE1BQU0sZ0JBQWdCLEdBQUcsQ0FBQyxPQUFlLEVBQUUsR0FBRyxJQUFXLEVBQVMsRUFBRTtJQUN6RSxJQUFJLE9BQU8sS0FBSywrQkFBbUIsQ0FBQyxXQUFXLEVBQUU7UUFDL0MsT0FBTyxJQUFJLGlCQUFpQixDQUFDLEdBQUcsSUFBSSxDQUFDLENBQUE7S0FDdEM7U0FBTSxJQUFJLE9BQU8sS0FBSywrQkFBbUIsQ0FBQyxpQkFBaUIsRUFBRTtRQUM1RCxPQUFPLElBQUksZUFBZSxDQUFDLEdBQUcsSUFBSSxDQUFDLENBQUE7S0FDcEM7SUFDRCwwQkFBMEI7SUFDMUIsTUFBTSxJQUFJLHFCQUFZLENBQUMsMkNBQTJDLENBQUMsQ0FBQTtBQUNyRSxDQUFDLENBQUE7QUFSWSxRQUFBLGdCQUFnQixvQkFRNUI7QUFFRCxNQUFhLGNBQWUsU0FBUSw4QkFBc0I7SUFBMUQ7O1FBQ1ksY0FBUyxHQUFHLGdCQUFnQixDQUFBO1FBQzVCLFlBQU8sR0FBRyxTQUFTLENBQUE7SUFrQi9CLENBQUM7SUFoQkMsd0JBQXdCO0lBRXhCLFdBQVcsQ0FBQyxNQUFjLEVBQUUsV0FBK0IsS0FBSztRQUM5RCxLQUFLLENBQUMsV0FBVyxDQUFDLE1BQU0sRUFBRSxRQUFRLENBQUMsQ0FBQTtRQUNuQyxJQUFJ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Q0FBQTtRQUNYLElBQUksQ0FBQyxpQkFBaUIsR0FBRyxJQUFJLGNBQWMsRUFBRSxDQUFBO1FBQzdDLE1BQU0sR0FBRyxJQUFJLENBQUMsaUJBQWlCLENBQUMsVUFBVSxDQUFDLEtBQUssRUFBRSxNQUFNLENBQUMsQ0FBQTtRQUN6RCxJQUFJLENBQUMsV0FBVyxFQUFFLENBQUE7UUFDbEIsT0FBTyxNQUFNLENBQUE7SUFDZixDQUFDO0lBRUQ7O09BRUc7SUFDSCxRQUFRO1FBQ04sTUFBTSxVQUFVLEdBQVcsSUFBSSxDQUFDLGlCQUFpQixDQUFDLFFBQVEsRUFBRSxDQUFBO1FBQzVELE1BQU0sS0FBSyxHQUFXLElBQUksQ0FBQyxpQkFBaUIsQ0FBQyxNQUFNLEdBQUcsVUFBVSxDQUFDLE1BQU0sQ0FBQTtRQUN2RSxNQUFNLElBQUksR0FBYSxDQUFDLElBQUksQ0FBQyxpQkFBaUIsRUFBRSxVQUFVLENBQUMsQ0FBQTtRQUMzRCxPQUFPLGVBQU0sQ0FBQyxNQUFNLENBQUMsSUFBSSxFQUFFLEtBQUssQ0FBQyxDQUFBO0lBQ25DLENBQUM7SUFFRCxNQUFNLENBQUMsR0FBRyxJQUFXO1FBQ25CLE9BQU8sSUFBSSxlQUFlLENBQUMsR0FBRyxJQUFJLENBQVMsQ0FBQTtJQUM3QyxDQUFDO0lBRUQsS0FBSztRQUNILE1BQU0sTUFBTSxHQUFvQixJQUFJLENBQUMsTUFBTSxFQUFFLENBQUE7UUFDN0MsTUFBTSxDQUFDLFVBQVUsQ0FBQyxJQUFJLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQTtRQUNsQyxPQUFPLE1BQWMsQ0FBQTtJQUN2QixDQUFDO0lBRUQsTUFBTSxDQUFDLEVBQVUsRUFBRSxHQUFHLElBQVc7UUFDL0IsT0FBTyxJQUFBLHdCQUFnQixFQUFDLEVBQUUsRUFBRSxHQUFHLElBQUksQ0FBQyxDQUFBO0lBQ3RDLENBQUM7Q0F1QkY7QUE5SEQsMENBOEhDIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgQVBJLVBsYXRmb3JtVk0tSW5wdXRzXG4gKi9cbmltcG9ydCB7IEJ1ZmZlciB9IGZyb20gXCJidWZmZXIvXCJcbmltcG9ydCBCaW5Ub29scyBmcm9tIFwiLi4vLi4vdXRpbHMvYmludG9vbHNcIlxuaW1wb3J0IHsgUGxhdGZvcm1WTUNvbnN0YW50cyB9IGZyb20gXCIuL2NvbnN0YW50c1wiXG5pbXBvcnQge1xuICBJbnB1dCxcbiAgU3RhbmRhcmRUcmFuc2ZlcmFibGVJbnB1dCxcbiAgU3RhbmRhcmRBbW91bnRJbnB1dCxcbiAgU3RhbmRhcmRQYXJzZWFibGVJbnB1dFxufSBmcm9tIFwiLi4vLi4vY29tbW9uL2lucHV0XCJcbmltcG9ydCB7IFNlcmlhbGl6YXRpb24sIFNlcmlhbGl6ZWRFbmNvZGluZyB9IGZyb20gXCIuLi8uLi91dGlscy9zZXJpYWxpemF0aW9uXCJcbmltcG9ydCBCTiBmcm9tIFwiYm4uanNcIlxuaW1wb3J0IHsgSW5wdXRJZEVycm9yIH0gZnJvbSBcIi4uLy4uL3V0aWxzL2Vycm9yc1wiXG5cbi8qKlxuICogQGlnbm9yZVxuICovXG5jb25zdCBiaW50b29sczogQmluVG9vbHMgPSBCaW5Ub29scy5nZXRJbnN0YW5jZSgpXG5jb25zdCBzZXJpYWxpemF0aW9uOiBTZXJpYWxpemF0aW9uID0gU2VyaWFsaXphdGlvbi5nZXRJbnN0YW5jZSgpXG5cbi8qKlxuICogVGFrZXMgYSBidWZmZXIgcmVwcmVzZW50aW5nIHRoZSBvdXRwdXQgYW5kIHJldHVybnMgdGhlIHByb3BlciBbW0lucHV0XV0gaW5zdGFuY2UuXG4gKlxuICogQHBhcmFtIGlucHV0aWQgQSBudW1iZXIgcmVwcmVzZW50aW5nIHRoZSBpbnB1dElEIHBhcnNlZCBwcmlvciB0byB0aGUgYnl0ZXMgcGFzc2VkIGluXG4gKlxuICogQHJldHVybnMgQW4gaW5zdGFuY2Ugb2YgYW4gW1tJbnB1dF1dLWV4dGVuZGVkIGNsYXNzLlxuICovXG5leHBvcnQgY29uc3QgU2VsZWN0SW5wdXRDbGFzcyA9IChpbnB1dGlkOiBudW1iZXIsIC4uLmFyZ3M6IGFueVtdKTogSW5wdXQgPT4ge1xuICBpZiAoaW5wdXRpZCA9PT0gUGxhdGZvcm1WTUNvbnN0YW50cy5TRUNQSU5QVVRJRCkge1xuICAgIHJldHVybiBuZXcgU0VDUFRyYW5zZmVySW5wdXQoLi4uYXJncylcbiAgfSBlbHNlIGlmIChpbnB1dGlkID09PSBQbGF0Zm9ybVZNQ29uc3RhbnRzLlNUQUtFQUJMRUxPQ0tJTklEKSB7XG4gICAgcmV0dXJuIG5ldyBTdGFrZWFibGVMb2NrSW4oLi4uYXJncylcbiAgfVxuICAvKiBpc3RhbmJ1bCBpZ25vcmUgbmV4dCAqL1xuICB0aHJvdyBuZXcgSW5wdXRJZEVycm9yKFwiRXJyb3IgLSBTZWxlY3RJbnB1dENsYXNzOiB1bmtub3duIGlucHV0aWRcIilcbn1cblxuZXhwb3J0IGNsYXNzIFBhcnNlYWJsZUlucHV0IGV4dGVuZHMgU3RhbmRhcmRQYXJzZWFibGVJbnB1dCB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIlBhcnNlYWJsZUlucHV0XCJcbiAgcHJvdGVjdGVkIF90eXBlSUQgPSB1bmRlZmluZWRcblxuICAvL3NlcmlhbGl6ZSBpcyBpbmhlcml0ZWRcblxuICBkZXNlcmlhbGl6ZShmaWVsZHM6IG9iamVjdCwgZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpIHtcbiAgICBzdXBlci5kZXNlcmlhbGl6ZShmaWVsZHMsIGVuY29kaW5nKVxuICAgIHRoaXMuaW5wdXQgPSBTZWxlY3RJbnB1dENsYXNzKGZpZWxkc1tcImlucHV0XCJdW1wiX3R5cGVJRFwiXSlcbiAgICB0aGlzLmlucHV0LmRlc2VyaWFsaXplKGZpZWxkc1tcImlucHV0XCJdLCBlbmNvZGluZylcbiAgfVxuXG4gIGZyb21CdWZmZXIoYnl0ZXM6IEJ1ZmZlciwgb2Zmc2V0OiBudW1iZXIgPSAwKTogbnVtYmVyIHtcbiAgICBjb25zdCBpbnB1dGlkOiBudW1iZXIgPSBiaW50b29sc1xuICAgICAgLmNvcHlGcm9tKGJ5dGVzLCBvZmZzZXQsIG9mZnNldCArIDQpXG4gICAgICAucmVhZFVJbnQzMkJFKDApXG4gICAgb2Zmc2V0ICs9IDRcbiAgICB0aGlzLmlucHV0ID0gU2VsZWN0SW5wdXRDbGFzcyhpbnB1dGlkKVxuICAgIHJldHVybiB0aGlzLmlucHV0LmZyb21CdWZmZXIoYnl0ZXMsIG9mZnNldClcbiAgfVxufVxuXG5leHBvcnQgY2xhc3MgVHJhbnNmZXJhYmxlSW5wdXQgZXh0ZW5kcyBTdGFuZGFyZFRyYW5zZmVyYWJsZUlucHV0IHtcbiAgcHJvdGVjdGVkIF90eXBlTmFtZSA9IFwiVHJhbnNmZXJhYmxlSW5wdXRcIlxuICBwcm90ZWN0ZWQgX3R5cGVJRCA9IHVuZGVmaW5lZFxuXG4gIC8vc2VyaWFsaXplIGlzIGluaGVyaXRlZFxuXG4gIGRlc2VyaWFsaXplKGZpZWxkczogb2JqZWN0LCBlbmNvZGluZzogU2VyaWFsaXplZEVuY29kaW5nID0gXCJoZXhcIikge1xuICAgIHN1cGVyLmRlc2VyaWFsaXplKGZpZWxkcywgZW5jb2RpbmcpXG4gICAgdGhpcy5pbnB1dCA9IFNlbGVjdElucHV0Q2xhc3MoZmllbGRzW1wiaW5wdXRcIl1bXCJfdHlwZUlEXCJdKVxuICAgIHRoaXMuaW5wdXQuZGVzZXJpYWxpemUoZmllbGRzW1wiaW5wdXRcIl0sIGVuY29kaW5nKVxuICB9XG5cbiAgLyoqXG4gICAqIFRha2VzIGEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gY29udGFpbmluZyBhIFtbVHJhbnNmZXJhYmxlSW5wdXRdXSwgcGFyc2VzIGl0LCBwb3B1bGF0ZXMgdGhlIGNsYXNzLCBhbmQgcmV0dXJucyB0aGUgbGVuZ3RoIG9mIHRoZSBbW1RyYW5zZmVyYWJsZUlucHV0XV0gaW4gYnl0ZXMuXG4gICAqXG4gICAqIEBwYXJhbSBieXRlcyBBIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IGNvbnRhaW5pbmcgYSByYXcgW1tUcmFuc2ZlcmFibGVJbnB1dF1dXG4gICAqXG4gICAqIEByZXR1cm5zIFRoZSBsZW5ndGggb2YgdGhlIHJhdyBbW1RyYW5zZmVyYWJsZUlucHV0XV1cbiAgICovXG4gIGZyb21CdWZmZXIoYnl0ZXM6IEJ1ZmZlciwgb2Zmc2V0OiBudW1iZXIgPSAwKTogbnVtYmVyIHtcbiAgICB0aGlzLnR4aWQgPSBiaW50b29scy5jb3B5RnJvbShieXRlcywgb2Zmc2V0LCBvZmZzZXQgKyAzMilcbiAgICBvZmZzZXQgKz0gMzJcbiAgICB0aGlzLm91dHB1dGlkeCA9IGJpbnRvb2xzLmNvcHlGcm9tKGJ5dGVzLCBvZmZzZXQsIG9mZnNldCArIDQpXG4gICAgb2Zmc2V0ICs9IDRcbiAgICB0aGlzLmFzc2V0SUQgPSBiaW50b29scy5jb3B5RnJvbShcbiAgICAgIGJ5dGVzLFxuICAgICAgb2Zmc2V0LFxuICAgICAgb2Zmc2V0ICsgUGxhdGZvcm1WTUNvbnN0YW50cy5BU1NFVElETEVOXG4gICAgKVxuICAgIG9mZnNldCArPSAzMlxuICAgIGNvbnN0IGlucHV0aWQ6IG51bWJlciA9IGJpbnRvb2xzXG4gICAgICAuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgNClcbiAgICAgIC5yZWFkVUludDMyQkUoMClcbiAgICBvZmZzZXQgKz0gNFxuICAgIHRoaXMuaW5wdXQgPSBTZWxlY3RJbnB1dENsYXNzKGlucHV0aWQpXG4gICAgcmV0dXJuIHRoaXMuaW5wdXQuZnJvbUJ1ZmZlcihieXRlcywgb2Zmc2V0KVxuICB9XG59XG5cbmV4cG9ydCBhYnN0cmFjdCBjbGFzcyBBbW91bnRJbnB1dCBleHRlbmRzIFN0YW5kYXJkQW1vdW50SW5wdXQge1xuICBwcm90ZWN0ZWQgX3R5cGVOYW1lID0gXCJBbW91bnRJbnB1dFwiXG4gIHByb3RlY3RlZCBfdHlwZUlEID0gdW5kZWZpbmVkXG5cbiAgLy9zZXJpYWxpemUgYW5kIGRlc2VyaWFsaXplIGJvdGggYXJlIGluaGVyaXRlZFxuXG4gIHNlbGVjdChpZDogbnVtYmVyLCAuLi5hcmdzOiBhbnlbXSk6IElucHV0IHtcbiAgICByZXR1cm4gU2VsZWN0SW5wdXRDbGFzcyhpZCwgLi4uYXJncylcbiAgfVxufVxuXG5leHBvcnQgY2xhc3MgU0VDUFRyYW5zZmVySW5wdXQgZXh0ZW5kcyBBbW91bnRJbnB1dCB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIlNFQ1BUcmFuc2ZlcklucHV0XCJcbiAgcHJvdGVjdGVkIF90eXBlSUQgPSBQbGF0Zm9ybVZNQ29uc3RhbnRzLlNFQ1BJTlBVVElEXG5cbiAgLy9zZXJpYWxpemUgYW5kIGRlc2VyaWFsaXplIGJvdGggYXJlIGluaGVyaXRlZFxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBpbnB1dElEIGZvciB0aGlzIGlucHV0XG4gICAqL1xuICBnZXRJbnB1dElEKCk6IG51bWJlciB7XG4gICAgcmV0dXJuIHRoaXMuX3R5cGVJRFxuICB9XG5cbiAgZ2V0Q3JlZGVudGlhbElEID0gKCk6IG51bWJlciA9PiBQbGF0Zm9ybVZNQ29uc3RhbnRzLlNFQ1BDUkVERU5USUFMXG5cbiAgY3JlYXRlKC4uLmFyZ3M6IGFueVtdKTogdGhpcyB7XG4gICAgcmV0dXJuIG5ldyBTRUNQVHJhbnNmZXJJbnB1dCguLi5hcmdzKSBhcyB0aGlzXG4gIH1cblxuICBjbG9uZSgpOiB0aGlzIHtcbiAgICBjb25zdCBuZXdvdXQ6IFNFQ1BUcmFuc2ZlcklucHV0ID0gdGhpcy5jcmVhdGUoKVxuICAgIG5ld291dC5mcm9tQnVmZmVyKHRoaXMudG9CdWZmZXIoKSlcbiAgICByZXR1cm4gbmV3b3V0IGFzIHRoaXNcbiAgfVxufVxuXG4vKipcbiAqIEFuIFtbSW5wdXRdXSBjbGFzcyB3aGljaCBzcGVjaWZpZXMgYW4gaW5wdXQgdGhhdCBoYXMgYSBsb2NrdGltZSB3aGljaCBjYW4gYWxzbyBlbmFibGUgc3Rha2luZyBvZiB0aGUgdmFsdWUgaGVsZCwgcHJldmVudGluZyB0cmFuc2ZlcnMgYnV0IG5vdCB2YWxpZGF0aW9uLlxuICovXG5leHBvcnQgY2xhc3MgU3Rha2VhYmxlTG9ja0luIGV4dGVuZHMgQW1vdW50SW5wdXQge1xuICBwcm90ZWN0ZWQgX3R5cGVOYW1lID0gXCJTdGFrZWFibGVMb2NrSW5cIlxuICBwcm90ZWN0ZWQgX3R5cGVJRCA9IFBsYXRmb3JtVk1Db25zdGFudHMuU1RBS0VBQkxFTE9DS0lOSURcblxuICAvL3NlcmlhbGl6ZSBhbmQgZGVzZXJpYWxpemUgYm90aCBhcmUgaW5oZXJpdGVkXG5cbiAgc2VyaWFsaXplKGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKTogb2JqZWN0IHtcbiAgICBsZXQgZmllbGRzOiBvYmplY3QgPSBzdXBlci5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgbGV0IG91dG9iajogb2JqZWN0ID0ge1xuICAgICAgLi4uZmllbGRzLFxuICAgICAgc3Rha2VhYmxlTG9ja3RpbWU6IHNlcmlhbGl6YXRpb24uZW5jb2RlcihcbiAgICAgICAgdGhpcy5zdGFrZWFibGVMb2NrdGltZSxcbiAgICAgICAgZW5jb2RpbmcsXG4gICAgICAgIFwiQnVmZmVyXCIsXG4gICAgICAgIFwiZGVjaW1hbFN0cmluZ1wiLFxuICAgICAgICA4XG4gICAgICApLFxuICAgICAgdHJhbnNmZXJhYmxlSW5wdXQ6IHRoaXMudHJhbnNmZXJhYmxlSW5wdXQuc2VyaWFsaXplKGVuY29kaW5nKVxuICAgIH1cbiAgICBkZWxldGUgb3V0b2JqW1wic2lnSWR4c1wiXVxuICAgIGRlbGV0ZSBvdXRvYmpbXCJzaWdDb3VudFwiXVxuICAgIGRlbGV0ZSBvdXRvYmpbXCJhbW91bnRcIl1cbiAgICByZXR1cm4gb3V0b2JqXG4gIH1cbiAgZGVzZXJpYWxpemUoZmllbGRzOiBvYmplY3QsIGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKSB7XG4gICAgZmllbGRzW1wic2lnSWR4c1wiXSA9IFtdXG4gICAgZmllbGRzW1wic2lnQ291bnRcIl0gPSBcIjBcIlxuICAgIGZpZWxkc1tcImFtb3VudFwiXSA9IFwiOThcIlxuICAgIHN1cGVyLmRlc2VyaWFsaXplKGZpZWxkcywgZW5jb2RpbmcpXG4gICAgdGhpcy5zdGFrZWFibGVMb2NrdGltZSA9IHNlcmlhbGl6YXRpb24uZGVjb2RlcihcbiAgICAgIGZpZWxkc1tcInN0YWtlYWJsZUxvY2t0aW1lXCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBcImRlY2ltYWxTdHJpbmdcIixcbiAgICAgIFwiQnVmZmVyXCIsXG4gICAgICA4XG4gICAgKVxuICAgIHRoaXMudHJhbnNmZXJhYmxlSW5wdXQgPSBuZXcgUGFyc2VhYmxlSW5wdXQoKVxuICAgIHRoaXMudHJhbnNmZXJhYmxlSW5wdXQuZGVzZXJpYWxpemUoZmllbGRzW1widHJhbnNmZXJhYmxlSW5wdXRcIl0sIGVuY29kaW5nKVxuICAgIHRoaXMuc3luY2hyb25pemUoKVxuICB9XG5cbiAgcHJvdGVjdGVkIHN0YWtlYWJsZUxvY2t0aW1lOiBCdWZmZXJcbiAgcHJvdGVjdGVkIHRyYW5zZmVyYWJsZUlucHV0OiBQYXJzZWFibGVJbnB1dFxuXG4gIHByaXZhdGUgc3luY2hyb25pemUoKSB7XG4gICAgbGV0IGlucHV0OiBBbW91bnRJbnB1dCA9IHRoaXMudHJhbnNmZXJhYmxlSW5wdXQuZ2V0SW5wdXQoKSBhcyBBbW91bnRJbnB1dFxuICAgIHRoaXMuc2lnSWR4cyA9IGlucHV0LmdldFNpZ0lkeHMoKVxuICAgIHRoaXMuc2lnQ291bnQgPSBCdWZmZXIuYWxsb2MoNClcbiAgICB0aGlzLnNpZ0NvdW50LndyaXRlVUludDMyQkUodGhpcy5zaWdJZHhzLmxlbmd0aCwgMClcbiAgICB0aGlzLmFtb3VudCA9IGJpbnRvb2xzLmZyb21CTlRvQnVmZmVyKGlucHV0LmdldEFtb3VudCgpLCA4KVxuICAgIHRoaXMuYW1vdW50VmFsdWUgPSBpbnB1dC5nZXRBbW91bnQoKVxuICB9XG5cbiAgZ2V0U3Rha2VhYmxlTG9ja3RpbWUoKTogQk4ge1xuICAgIHJldHVybiBiaW50b29scy5mcm9tQnVmZmVyVG9CTih0aGlzLnN0YWtlYWJsZUxvY2t0aW1lKVxuICB9XG5cbiAgZ2V0VHJhbnNmZXJhYmxJbnB1dCgpOiBQYXJzZWFibGVJbnB1dCB7XG4gICAgcmV0dXJuIHRoaXMudHJhbnNmZXJhYmxlSW5wdXRcbiAgfVxuICAvKipcbiAgICogUmV0dXJucyB0aGUgaW5wdXRJRCBmb3IgdGhpcyBpbnB1dFxuICAgKi9cbiAgZ2V0SW5wdXRJRCgpOiBudW1iZXIge1xuICAgIHJldHVybiB0aGlzLl90eXBlSURcbiAgfVxuXG4gIGdldENyZWRlbnRpYWxJRCA9ICgpOiBudW1iZXIgPT4gUGxhdGZvcm1WTUNvbnN0YW50cy5TRUNQQ1JFREVOVElBTFxuXG4gIC8qKlxuICAgKiBQb3B1YXRlcyB0aGUgaW5zdGFuY2UgZnJvbSBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IHJlcHJlc2VudGluZyB0aGUgW1tTdGFrZWFibGVMb2NrSW5dXSBhbmQgcmV0dXJucyB0aGUgc2l6ZSBvZiB0aGUgb3V0cHV0LlxuICAgKi9cbiAgZnJvbUJ1ZmZlcihieXRlczogQnVmZmVyLCBvZmZzZXQ6IG51bWJlciA9IDApOiBudW1iZXIge1xuICAgIHRoaXMuc3Rha2VhYmxlTG9ja3RpbWUgPSBiaW50b29scy5jb3B5RnJvbShieXRlcywgb2Zmc2V0LCBvZmZzZXQgKyA4KVxuICAgIG9mZnNldCArPSA4XG4gICAgdGhpcy50cmFuc2ZlcmFibGVJbnB1dCA9IG5ldyBQYXJzZWFibGVJbnB1dCgpXG4gICAgb2Zmc2V0ID0gdGhpcy50cmFuc2ZlcmFibGVJbnB1dC5mcm9tQnVmZmVyKGJ5dGVzLCBvZmZzZXQpXG4gICAgdGhpcy5zeW5jaHJvbml6ZSgpXG4gICAgcmV0dXJuIG9mZnNldFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIGJ1ZmZlciByZXByZXNlbnRpbmcgdGhlIFtbU3Rha2VhYmxlTG9ja0luXV0gaW5zdGFuY2UuXG4gICAqL1xuICB0b0J1ZmZlcigpOiBCdWZmZXIge1xuICAgIGNvbnN0IHhmZXJpbkJ1ZmY6IEJ1ZmZlciA9IHRoaXMudHJhbnNmZXJhYmxlSW5wdXQudG9CdWZmZXIoKVxuICAgIGNvbnN0IGJzaXplOiBudW1iZXIgPSB0aGlzLnN0YWtlYWJsZUxvY2t0aW1lLmxlbmd0aCArIHhmZXJpbkJ1ZmYubGVuZ3RoXG4gICAgY29uc3QgYmFycjogQnVmZmVyW10gPSBbdGhpcy5zdGFrZWFibGVMb2NrdGltZSwgeGZlcmluQnVmZl1cbiAgICByZXR1cm4gQnVmZmVyLmNvbmNhdChiYXJyLCBic2l6ZSlcbiAgfVxuXG4gIGNyZWF0ZSguLi5hcmdzOiBhbnlbXSk6IHRoaXMge1xuICAgIHJldHVybiBuZXcgU3Rha2VhYmxlTG9ja0luKC4uLmFyZ3MpIGFzIHRoaXNcbiAgfVxuXG4gIGNsb25lKCk6IHRoaXMge1xuICAgIGNvbnN0IG5ld291dDogU3Rha2VhYmxlTG9ja0luID0gdGhpcy5jcmVhdGUoKVxuICAgIG5ld291dC5mcm9tQnVmZmVyKHRoaXMudG9CdWZmZXIoKSlcbiAgICByZXR1cm4gbmV3b3V0IGFzIHRoaXNcbiAgfVxuXG4gIHNlbGVjdChpZDogbnVtYmVyLCAuLi5hcmdzOiBhbnlbXSk6IElucHV0IHtcbiAgICByZXR1cm4gU2VsZWN0SW5wdXRDbGFzcyhpZCwgLi4uYXJncylcbiAgfVxuXG4gIC8qKlxuICAgKiBBIFtbT3V0cHV0XV0gY2xhc3Mgd2hpY2ggc3BlY2lmaWVzIGFuIFtbSW5wdXRdXSB0aGF0IGhhcyBhIGxvY2t0aW1lIHdoaWNoIGNhbiBhbHNvIGVuYWJsZSBzdGFraW5nIG9mIHRoZSB2YWx1ZSBoZWxkLCBwcmV2ZW50aW5nIHRyYW5zZmVycyBidXQgbm90IHZhbGlkYXRpb24uXG4gICAqXG4gICAqIEBwYXJhbSBhbW91bnQgQSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2luZHV0bnkvYm4uanMvfEJOfSByZXByZXNlbnRpbmcgdGhlIGFtb3VudCBpbiB0aGUgaW5wdXRcbiAgICogQHBhcmFtIHN0YWtlYWJsZUxvY2t0aW1lIEEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9pbmR1dG55L2JuLmpzL3xCTn0gcmVwcmVzZW50aW5nIHRoZSBzdGFrZWFibGUgbG9ja3RpbWVcbiAgICogQHBhcmFtIHRyYW5zZmVyYWJsZUlucHV0IEEgW1tQYXJzZWFibGVJbnB1dF1dIHdoaWNoIGlzIGVtYmVkZGVkIGludG8gdGhpcyBpbnB1dC5cbiAgICovXG4gIGNvbnN0cnVjdG9yKFxuICAgIGFtb3VudDogQk4gPSB1bmRlZmluZWQsXG4gICAgc3Rha2VhYmxlTG9ja3RpbWU6IEJOID0gdW5kZWZpbmVkLFxuICAgIHRyYW5zZmVyYWJsZUlucHV0OiBQYXJzZWFibGVJbnB1dCA9IHVuZGVmaW5lZFxuICApIHtcbiAgICBzdXBlcihhbW91bnQpXG4gICAgaWYgKHR5cGVvZiBzdGFrZWFibGVMb2NrdGltZSAhPT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgdGhpcy5zdGFrZWFibGVMb2NrdGltZSA9IGJpbnRvb2xzLmZyb21CTlRvQnVmZmVyKHN0YWtlYWJsZUxvY2t0aW1lLCA4KVxuICAgIH1cbiAgICBpZiAodHlwZW9mIHRyYW5zZmVyYWJsZUlucHV0ICE9PSBcInVuZGVmaW5lZFwiKSB7XG4gICAgICB0aGlzLnRyYW5zZmVyYWJsZUlucHV0ID0gdHJhbnNmZXJhYmxlSW5wdXRcbiAgICAgIHRoaXMuc3luY2hyb25pemUoKVxuICAgIH1cbiAgfVxufVxuIl19