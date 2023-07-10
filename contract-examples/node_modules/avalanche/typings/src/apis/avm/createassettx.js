"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateAssetTx = void 0;
/**
 * @packageDocumentation
 * @module API-AVM-CreateAssetTx
 */
const buffer_1 = require("buffer/");
const bintools_1 = __importDefault(require("../../utils/bintools"));
const constants_1 = require("./constants");
const initialstates_1 = require("./initialstates");
const basetx_1 = require("./basetx");
const constants_2 = require("../../utils/constants");
const serialization_1 = require("../../utils/serialization");
const errors_1 = require("../../utils/errors");
/**
 * @ignore
 */
const bintools = bintools_1.default.getInstance();
const serialization = serialization_1.Serialization.getInstance();
const utf8 = "utf8";
const decimalString = "decimalString";
const buffer = "Buffer";
class CreateAssetTx extends basetx_1.BaseTx {
    /**
     * Class representing an unsigned Create Asset transaction.
     *
     * @param networkID Optional networkID, [[DefaultNetworkID]]
     * @param blockchainID Optional blockchainID, default Buffer.alloc(32, 16)
     * @param outs Optional array of the [[TransferableOutput]]s
     * @param ins Optional array of the [[TransferableInput]]s
     * @param memo Optional {@link https://github.com/feross/buffer|Buffer} for the memo field
     * @param name String for the descriptive name of the asset
     * @param symbol String for the ticker symbol of the asset
     * @param denomination Optional number for the denomination which is 10^D. D must be >= 0 and <= 32. Ex: $1 AVAX = 10^9 $nAVAX
     * @param initialState Optional [[InitialStates]] that represent the intial state of a created asset
     */
    constructor(networkID = constants_2.DefaultNetworkID, blockchainID = buffer_1.Buffer.alloc(32, 16), outs = undefined, ins = undefined, memo = undefined, name = undefined, symbol = undefined, denomination = undefined, initialState = undefined) {
        super(networkID, blockchainID, outs, ins, memo);
        this._typeName = "CreateAssetTx";
        this._codecID = constants_1.AVMConstants.LATESTCODEC;
        this._typeID = this._codecID === 0
            ? constants_1.AVMConstants.CREATEASSETTX
            : constants_1.AVMConstants.CREATEASSETTX_CODECONE;
        this.name = "";
        this.symbol = "";
        this.denomination = buffer_1.Buffer.alloc(1);
        this.initialState = new initialstates_1.InitialStates();
        if (typeof name === "string" &&
            typeof symbol === "string" &&
            typeof denomination === "number" &&
            denomination >= 0 &&
            denomination <= 32 &&
            typeof initialState !== "undefined") {
            this.initialState = initialState;
            this.name = name;
            this.symbol = symbol;
            this.denomination.writeUInt8(denomination, 0);
        }
    }
    serialize(encoding = "hex") {
        const fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { name: serialization.encoder(this.name, encoding, utf8, utf8), symbol: serialization.encoder(this.symbol, encoding, utf8, utf8), denomination: serialization.encoder(this.denomination, encoding, buffer, decimalString, 1), initialState: this.initialState.serialize(encoding) });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.name = serialization.decoder(fields["name"], encoding, utf8, utf8);
        this.symbol = serialization.decoder(fields["symbol"], encoding, utf8, utf8);
        this.denomination = serialization.decoder(fields["denomination"], encoding, decimalString, buffer, 1);
        this.initialState = new initialstates_1.InitialStates();
        this.initialState.deserialize(fields["initialState"], encoding);
    }
    /**
     * Set the codecID
     *
     * @param codecID The codecID to set
     */
    setCodecID(codecID) {
        if (codecID !== 0 && codecID !== 1) {
            /* istanbul ignore next */
            throw new errors_1.CodecIdError("Error - CreateAssetTx.setCodecID: invalid codecID. Valid codecIDs are 0 and 1.");
        }
        this._codecID = codecID;
        this._typeID =
            this._codecID === 0
                ? constants_1.AVMConstants.CREATEASSETTX
                : constants_1.AVMConstants.CREATEASSETTX_CODECONE;
    }
    /**
     * Returns the id of the [[CreateAssetTx]]
     */
    getTxType() {
        return this._typeID;
    }
    /**
     * Returns the array of array of [[Output]]s for the initial state
     */
    getInitialStates() {
        return this.initialState;
    }
    /**
     * Returns the string representation of the name
     */
    getName() {
        return this.name;
    }
    /**
     * Returns the string representation of the symbol
     */
    getSymbol() {
        return this.symbol;
    }
    /**
     * Returns the numeric representation of the denomination
     */
    getDenomination() {
        return this.denomination.readUInt8(0);
    }
    /**
     * Returns the {@link https://github.com/feross/buffer|Buffer} representation of the denomination
     */
    getDenominationBuffer() {
        return this.denomination;
    }
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} containing an [[CreateAssetTx]], parses it, populates the class, and returns the length of the [[CreateAssetTx]] in bytes.
     *
     * @param bytes A {@link https://github.com/feross/buffer|Buffer} containing a raw [[CreateAssetTx]]
     *
     * @returns The length of the raw [[CreateAssetTx]]
     *
     * @remarks assume not-checksummed
     */
    fromBuffer(bytes, offset = 0) {
        offset = super.fromBuffer(bytes, offset);
        const namesize = bintools
            .copyFrom(bytes, offset, offset + 2)
            .readUInt16BE(0);
        offset += 2;
        this.name = bintools
            .copyFrom(bytes, offset, offset + namesize)
            .toString("utf8");
        offset += namesize;
        const symsize = bintools
            .copyFrom(bytes, offset, offset + 2)
            .readUInt16BE(0);
        offset += 2;
        this.symbol = bintools
            .copyFrom(bytes, offset, offset + symsize)
            .toString("utf8");
        offset += symsize;
        this.denomination = bintools.copyFrom(bytes, offset, offset + 1);
        offset += 1;
        const inits = new initialstates_1.InitialStates();
        offset = inits.fromBuffer(bytes, offset);
        this.initialState = inits;
        return offset;
    }
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[CreateAssetTx]].
     */
    toBuffer() {
        const superbuff = super.toBuffer();
        const initstatebuff = this.initialState.toBuffer();
        const namebuff = buffer_1.Buffer.alloc(this.name.length);
        namebuff.write(this.name, 0, this.name.length, utf8);
        const namesize = buffer_1.Buffer.alloc(2);
        namesize.writeUInt16BE(this.name.length, 0);
        const symbuff = buffer_1.Buffer.alloc(this.symbol.length);
        symbuff.write(this.symbol, 0, this.symbol.length, utf8);
        const symsize = buffer_1.Buffer.alloc(2);
        symsize.writeUInt16BE(this.symbol.length, 0);
        const bsize = superbuff.length +
            namesize.length +
            namebuff.length +
            symsize.length +
            symbuff.length +
            this.denomination.length +
            initstatebuff.length;
        const barr = [
            superbuff,
            namesize,
            namebuff,
            symsize,
            symbuff,
            this.denomination,
            initstatebuff
        ];
        return buffer_1.Buffer.concat(barr, bsize);
    }
    clone() {
        let newbase = new CreateAssetTx();
        newbase.fromBuffer(this.toBuffer());
        return newbase;
    }
    create(...args) {
        return new CreateAssetTx(...args);
    }
}
exports.CreateAssetTx = CreateAssetTx;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY3JlYXRlYXNzZXR0eC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3NyYy9hcGlzL2F2bS9jcmVhdGVhc3NldHR4LnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7OztBQUFBOzs7R0FHRztBQUNILG9DQUFnQztBQUNoQyxvRUFBMkM7QUFDM0MsMkNBQTBDO0FBRzFDLG1EQUErQztBQUMvQyxxQ0FBaUM7QUFDakMscURBQXdEO0FBQ3hELDZEQUlrQztBQUNsQywrQ0FBaUQ7QUFFakQ7O0dBRUc7QUFDSCxNQUFNLFFBQVEsR0FBYSxrQkFBUSxDQUFDLFdBQVcsRUFBRSxDQUFBO0FBQ2pELE1BQU0sYUFBYSxHQUFrQiw2QkFBYSxDQUFDLFdBQVcsRUFBRSxDQUFBO0FBQ2hFLE1BQU0sSUFBSSxHQUFtQixNQUFNLENBQUE7QUFDbkMsTUFBTSxhQUFhLEdBQW1CLGVBQWUsQ0FBQTtBQUNyRCxNQUFNLE1BQU0sR0FBbUIsUUFBUSxDQUFBO0FBRXZDLE1BQWEsYUFBYyxTQUFRLGVBQU07SUFnTXZDOzs7Ozs7Ozs7Ozs7T0FZRztJQUNILFlBQ0UsWUFBb0IsNEJBQWdCLEVBQ3BDLGVBQXVCLGVBQU0sQ0FBQyxLQUFLLENBQUMsRUFBRSxFQUFFLEVBQUUsQ0FBQyxFQUMzQyxPQUE2QixTQUFTLEVBQ3RDLE1BQTJCLFNBQVMsRUFDcEMsT0FBZSxTQUFTLEVBQ3hCLE9BQWUsU0FBUyxFQUN4QixTQUFpQixTQUFTLEVBQzFCLGVBQXVCLFNBQVMsRUFDaEMsZUFBOEIsU0FBUztRQUV2QyxLQUFLLENBQUMsU0FBUyxFQUFFLFlBQVksRUFBRSxJQUFJLEVBQUUsR0FBRyxFQUFFLElBQUksQ0FBQyxDQUFBO1FBdk52QyxjQUFTLEdBQUcsZUFBZSxDQUFBO1FBQzNCLGFBQVEsR0FBRyx3QkFBWSxDQUFDLFdBQVcsQ0FBQTtRQUNuQyxZQUFPLEdBQ2YsSUFBSSxDQUFDLFFBQVEsS0FBSyxDQUFDO1lBQ2pCLENBQUMsQ0FBQyx3QkFBWSxDQUFD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QzNDLENBQUM7Q0F5Q0Y7QUF2T0Qsc0NBdU9DIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgQVBJLUFWTS1DcmVhdGVBc3NldFR4XG4gKi9cbmltcG9ydCB7IEJ1ZmZlciB9IGZyb20gXCJidWZmZXIvXCJcbmltcG9ydCBCaW5Ub29scyBmcm9tIFwiLi4vLi4vdXRpbHMvYmludG9vbHNcIlxuaW1wb3J0IHsgQVZNQ29uc3RhbnRzIH0gZnJvbSBcIi4vY29uc3RhbnRzXCJcbmltcG9ydCB7IFRyYW5zZmVyYWJsZU91dHB1dCB9IGZyb20gXCIuL291dHB1dHNcIlxuaW1wb3J0IHsgVHJhbnNmZXJhYmxlSW5wdXQgfSBmcm9tIFwiLi9pbnB1dHNcIlxuaW1wb3J0IHsgSW5pdGlhbFN0YXRlcyB9IGZyb20gXCIuL2luaXRpYWxzdGF0ZXNcIlxuaW1wb3J0IHsgQmFzZVR4IH0gZnJvbSBcIi4vYmFzZXR4XCJcbmltcG9ydCB7IERlZmF1bHROZXR3b3JrSUQgfSBmcm9tIFwiLi4vLi4vdXRpbHMvY29uc3RhbnRzXCJcbmltcG9ydCB7XG4gIFNlcmlhbGl6YXRpb24sXG4gIFNlcmlhbGl6ZWRFbmNvZGluZyxcbiAgU2VyaWFsaXplZFR5cGVcbn0gZnJvbSBcIi4uLy4uL3V0aWxzL3NlcmlhbGl6YXRpb25cIlxuaW1wb3J0IHsgQ29kZWNJZEVycm9yIH0gZnJvbSBcIi4uLy4uL3V0aWxzL2Vycm9yc1wiXG5cbi8qKlxuICogQGlnbm9yZVxuICovXG5jb25zdCBiaW50b29sczogQmluVG9vbHMgPSBCaW5Ub29scy5nZXRJbnN0YW5jZSgpXG5jb25zdCBzZXJpYWxpemF0aW9uOiBTZXJpYWxpemF0aW9uID0gU2VyaWFsaXphdGlvbi5nZXRJbnN0YW5jZSgpXG5jb25zdCB1dGY4OiBTZXJpYWxpemVkVHlwZSA9IFwidXRmOFwiXG5jb25zdCBkZWNpbWFsU3RyaW5nOiBTZXJpYWxpemVkVHlwZSA9IFwiZGVjaW1hbFN0cmluZ1wiXG5jb25zdCBidWZmZXI6IFNlcmlhbGl6ZWRUeXBlID0gXCJCdWZmZXJcIlxuXG5leHBvcnQgY2xhc3MgQ3JlYXRlQXNzZXRUeCBleHRlbmRzIEJhc2VUeCB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIkNyZWF0ZUFzc2V0VHhcIlxuICBwcm90ZWN0ZWQgX2NvZGVjSUQgPSBBVk1Db25zdGFudHMuTEFURVNUQ09ERUNcbiAgcHJvdGVjdGVkIF90eXBlSUQgPVxuICAgIHRoaXMuX2NvZGVjSUQgPT09IDBcbiAgICAgID8gQVZNQ29uc3RhbnRzLkNSRUFURUFTU0VUVFhcbiAgICAgIDogQVZNQ29uc3RhbnRzLkNSRUFURUFTU0VUVFhfQ09ERUNPTkVcblxuICBzZXJpYWxpemUoZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpOiBvYmplY3Qge1xuICAgIGNvbnN0IGZpZWxkczogb2JqZWN0ID0gc3VwZXIuc2VyaWFsaXplKGVuY29kaW5nKVxuICAgIHJldHVybiB7XG4gICAgICAuLi5maWVsZHMsXG4gICAgICBuYW1lOiBzZXJpYWxpemF0aW9uLmVuY29kZXIodGhpcy5uYW1lLCBlbmNvZGluZywgdXRmOCwgdXRmOCksXG4gICAgICBzeW1ib2w6IHNlcmlhbGl6YXRpb24uZW5jb2Rlcih0aGlzLnN5bWJvbCwgZW5jb2RpbmcsIHV0ZjgsIHV0ZjgpLFxuICAgICAgZGVub21pbmF0aW9uOiBzZXJpYWxpemF0aW9uLmVuY29kZXIoXG4gICAgICAgIHRoaXMuZGVub21pbmF0aW9uLFxuICAgICAgICBlbmNvZGluZyxcbiAgICAgICAgYnVmZmVyLFxuICAgICAgICBkZWNpbWFsU3RyaW5nLFxuICAgICAgICAxXG4gICAgICApLFxuICAgICAgaW5pdGlhbFN0YXRlOiB0aGlzLmluaXRpYWxTdGF0ZS5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgfVxuICB9XG4gIGRlc2VyaWFsaXplKGZpZWxkczogb2JqZWN0LCBlbmNvZGluZzogU2VyaWFsaXplZEVuY29kaW5nID0gXCJoZXhcIikge1xuICAgIHN1cGVyLmRlc2VyaWFsaXplKGZpZWxkcywgZW5jb2RpbmcpXG4gICAgdGhpcy5uYW1lID0gc2VyaWFsaXphdGlvbi5kZWNvZGVyKGZpZWxkc1tcIm5hbWVcIl0sIGVuY29kaW5nLCB1dGY4LCB1dGY4KVxuICAgIHRoaXMuc3ltYm9sID0gc2VyaWFsaXphdGlvbi5kZWNvZGVyKGZpZWxkc1tcInN5bWJvbFwiXSwgZW5jb2RpbmcsIHV0ZjgsIHV0ZjgpXG4gICAgdGhpcy5kZW5vbWluYXRpb24gPSBzZXJpYWxpemF0aW9uLmRlY29kZXIoXG4gICAgICBmaWVsZHNbXCJkZW5vbWluYXRpb25cIl0sXG4gICAgICBlbmNvZGluZyxcbiAgICAgIGRlY2ltYWxTdHJpbmcsXG4gICAgICBidWZmZXIsXG4gICAgICAxXG4gICAgKVxuICAgIHRoaXMuaW5pdGlhbFN0YXRlID0gbmV3IEluaXRpYWxTdGF0ZXMoKVxuICAgIHRoaXMuaW5pdGlhbFN0YXRlLmRlc2VyaWFsaXplKGZpZWxkc1tcImluaXRpYWxTdGF0ZVwiXSwgZW5jb2RpbmcpXG4gIH1cblxuICBwcm90ZWN0ZWQgbmFtZTogc3RyaW5nID0gXCJcIlxuICBwcm90ZWN0ZWQgc3ltYm9sOiBzdHJpbmcgPSBcIlwiXG4gIHByb3RlY3RlZCBkZW5vbWluYXRpb246IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYygxKVxuICBwcm90ZWN0ZWQgaW5pdGlhbFN0YXRlOiBJbml0aWFsU3RhdGVzID0gbmV3IEluaXRpYWxTdGF0ZXMoKVxuXG4gIC8qKlxuICAgKiBTZXQgdGhlIGNvZGVjSURcbiAgICpcbiAgICogQHBhcmFtIGNvZGVjSUQgVGhlIGNvZGVjSUQgdG8gc2V0XG4gICAqL1xuICBzZXRDb2RlY0lEKGNvZGVjSUQ6IG51bWJlcik6IHZvaWQge1xuICAgIGlmIChjb2RlY0lEICE9PSAwICYmIGNvZGVjSUQgIT09IDEpIHtcbiAgICAgIC8qIGlzdGFuYnVsIGlnbm9yZSBuZXh0ICovXG4gICAgICB0aHJvdyBuZXcgQ29kZWNJZEVycm9yKFxuICAgICAgICBcIkVycm9yIC0gQ3JlYXRlQXNzZXRUeC5zZXRDb2RlY0lEOiBpbnZhbGlkIGNvZGVjSUQuIFZhbGlkIGNvZGVjSURzIGFyZSAwIGFuZCAxLlwiXG4gICAgICApXG4gICAgfVxuICAgIHRoaXMuX2NvZGVjSUQgPSBjb2RlY0lEXG4gICAgdGhpcy5fdHlwZUlEID1cbiAgICAgIHRoaXMuX2NvZGVjSUQgPT09IDBcbiAgICAgICAgPyBBVk1Db25zdGFudHMuQ1JFQVRFQVNTRVRUWFxuICAgICAgICA6IEFWTUNvbnN0YW50cy5DUkVBVEVBU1NFVFRYX0NPREVDT05FXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgaWQgb2YgdGhlIFtbQ3JlYXRlQXNzZXRUeF1dXG4gICAqL1xuICBnZXRUeFR5cGUoKTogbnVtYmVyIHtcbiAgICByZXR1cm4gdGhpcy5fdHlwZUlEXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgYXJyYXkgb2YgYXJyYXkgb2YgW1tPdXRwdXRdXXMgZm9yIHRoZSBpbml0aWFsIHN0YXRlXG4gICAqL1xuICBnZXRJbml0aWFsU3RhdGVzKCk6IEluaXRpYWxTdGF0ZXMge1xuICAgIHJldHVybiB0aGlzLmluaXRpYWxTdGF0ZVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIHN0cmluZyByZXByZXNlbnRhdGlvbiBvZiB0aGUgbmFtZVxuICAgKi9cbiAgZ2V0TmFtZSgpOiBzdHJpbmcge1xuICAgIHJldHVybiB0aGlzLm5hbWVcbiAgfVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBzdHJpbmcgcmVwcmVzZW50YXRpb24gb2YgdGhlIHN5bWJvbFxuICAgKi9cbiAgZ2V0U3ltYm9sKCk6IHN0cmluZyB7XG4gICAgcmV0dXJuIHRoaXMuc3ltYm9sXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgbnVtZXJpYyByZXByZXNlbnRhdGlvbiBvZiB0aGUgZGVub21pbmF0aW9uXG4gICAqL1xuICBnZXREZW5vbWluYXRpb24oKTogbnVtYmVyIHtcbiAgICByZXR1cm4gdGhpcy5kZW5vbWluYXRpb24ucmVhZFVJbnQ4KDApXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gcmVwcmVzZW50YXRpb24gb2YgdGhlIGRlbm9taW5hdGlvblxuICAgKi9cbiAgZ2V0RGVub21pbmF0aW9uQnVmZmVyKCk6IEJ1ZmZlciB7XG4gICAgcmV0dXJuIHRoaXMuZGVub21pbmF0aW9uXG4gIH1cblxuICAvKipcbiAgICogVGFrZXMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSBjb250YWluaW5nIGFuIFtbQ3JlYXRlQXNzZXRUeF1dLCBwYXJzZXMgaXQsIHBvcHVsYXRlcyB0aGUgY2xhc3MsIGFuZCByZXR1cm5zIHRoZSBsZW5ndGggb2YgdGhlIFtbQ3JlYXRlQXNzZXRUeF1dIGluIGJ5dGVzLlxuICAgKlxuICAgKiBAcGFyYW0gYnl0ZXMgQSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSBjb250YWluaW5nIGEgcmF3IFtbQ3JlYXRlQXNzZXRUeF1dXG4gICAqXG4gICAqIEByZXR1cm5zIFRoZSBsZW5ndGggb2YgdGhlIHJhdyBbW0NyZWF0ZUFzc2V0VHhdXVxuICAgKlxuICAgKiBAcmVtYXJrcyBhc3N1bWUgbm90LWNoZWNrc3VtbWVkXG4gICAqL1xuICBmcm9tQnVmZmVyKGJ5dGVzOiBCdWZmZXIsIG9mZnNldDogbnVtYmVyID0gMCk6IG51bWJlciB7XG4gICAgb2Zmc2V0ID0gc3VwZXIuZnJvbUJ1ZmZlcihieXRlcywgb2Zmc2V0KVxuXG4gICAgY29uc3QgbmFtZXNpemU6IG51bWJlciA9IGJpbnRvb2xzXG4gICAgICAuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgMilcbiAgICAgIC5yZWFkVUludDE2QkUoMClcbiAgICBvZmZzZXQgKz0gMlxuICAgIHRoaXMubmFtZSA9IGJpbnRvb2xzXG4gICAgICAuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgbmFtZXNpemUpXG4gICAgICAudG9TdHJpbmcoXCJ1dGY4XCIpXG4gICAgb2Zmc2V0ICs9IG5hbWVzaXplXG5cbiAgICBjb25zdCBzeW1zaXplOiBudW1iZXIgPSBiaW50b29sc1xuICAgICAgLmNvcHlGcm9tKGJ5dGVzLCBvZmZzZXQsIG9mZnNldCArIDIpXG4gICAgICAucmVhZFVJbnQxNkJFKDApXG4gICAgb2Zmc2V0ICs9IDJcbiAgICB0aGlzLnN5bWJvbCA9IGJpbnRvb2xzXG4gICAgICAuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgc3ltc2l6ZSlcbiAgICAgIC50b1N0cmluZyhcInV0ZjhcIilcbiAgICBvZmZzZXQgKz0gc3ltc2l6ZVxuXG4gICAgdGhpcy5kZW5vbWluYXRpb24gPSBiaW50b29scy5jb3B5RnJvbShieXRlcywgb2Zmc2V0LCBvZmZzZXQgKyAxKVxuICAgIG9mZnNldCArPSAxXG5cbiAgICBjb25zdCBpbml0czogSW5pdGlhbFN0YXRlcyA9IG5ldyBJbml0aWFsU3RhdGVzKClcbiAgICBvZmZzZXQgPSBpbml0cy5mcm9tQnVmZmVyKGJ5dGVzLCBvZmZzZXQpXG4gICAgdGhpcy5pbml0aWFsU3RhdGUgPSBpbml0c1xuXG4gICAgcmV0dXJuIG9mZnNldFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSByZXByZXNlbnRhdGlvbiBvZiB0aGUgW1tDcmVhdGVBc3NldFR4XV0uXG4gICAqL1xuICB0b0J1ZmZlcigpOiBCdWZmZXIge1xuICAgIGNvbnN0IHN1cGVyYnVmZjogQnVmZmVyID0gc3VwZXIudG9CdWZmZXIoKVxuICAgIGNvbnN0IGluaXRzdGF0ZWJ1ZmY6IEJ1ZmZlciA9IHRoaXMuaW5pdGlhbFN0YXRlLnRvQnVmZmVyKClcblxuICAgIGNvbnN0IG5hbWVidWZmOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2ModGhpcy5uYW1lLmxlbmd0aClcbiAgICBuYW1lYnVmZi53cml0ZSh0aGlzLm5hbWUsIDAsIHRoaXMubmFtZS5sZW5ndGgsIHV0ZjgpXG4gICAgY29uc3QgbmFtZXNpemU6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYygyKVxuICAgIG5hbWVzaXplLndyaXRlVUludDE2QkUodGhpcy5uYW1lLmxlbmd0aCwgMClcblxuICAgIGNvbnN0IHN5bWJ1ZmY6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYyh0aGlzLnN5bWJvbC5sZW5ndGgpXG4gICAgc3ltYnVmZi53cml0ZSh0aGlzLnN5bWJvbCwgMCwgdGhpcy5zeW1ib2wubGVuZ3RoLCB1dGY4KVxuICAgIGNvbnN0IHN5bXNpemU6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYygyKVxuICAgIHN5bXNpemUud3JpdGVVSW50MTZCRSh0aGlzLnN5bWJvbC5sZW5ndGgsIDApXG5cbiAgICBjb25zdCBic2l6ZTogbnVtYmVyID1cbiAgICAgIHN1cGVyYnVmZi5sZW5ndGggK1xuICAgICAgbmFtZXNpemUubGVuZ3RoICtcbiAgICAgIG5hbWVidWZmLmxlbmd0aCArXG4gICAgICBzeW1zaXplLmxlbmd0aCArXG4gICAgICBzeW1idWZmLmxlbmd0aCArXG4gICAgICB0aGlzLmRlbm9taW5hdGlvbi5sZW5ndGggK1xuICAgICAgaW5pdHN0YXRlYnVmZi5sZW5ndGhcbiAgICBjb25zdCBiYXJyOiBCdWZmZXJbXSA9IFtcbiAgICAgIHN1cGVyYnVmZixcbiAgICAgIG5hbWVzaXplLFxuICAgICAgbmFtZWJ1ZmYsXG4gICAgICBzeW1zaXplLFxuICAgICAgc3ltYnVmZixcbiAgICAgIHRoaXMuZGVub21pbmF0aW9uLFxuICAgICAgaW5pdHN0YXRlYnVmZlxuICAgIF1cbiAgICByZXR1cm4gQnVmZmVyLmNvbmNhdChiYXJyLCBic2l6ZSlcbiAgfVxuXG4gIGNsb25lKCk6IHRoaXMge1xuICAgIGxldCBuZXdiYXNlOiBDcmVhdGVBc3NldFR4ID0gbmV3IENyZWF0ZUFzc2V0VHgoKVxuICAgIG5ld2Jhc2UuZnJvbUJ1ZmZlcih0aGlzLnRvQnVmZmVyKCkpXG4gICAgcmV0dXJuIG5ld2Jhc2UgYXMgdGhpc1xuICB9XG5cbiAgY3JlYXRlKC4uLmFyZ3M6IGFueVtdKTogdGhpcyB7XG4gICAgcmV0dXJuIG5ldyBDcmVhdGVBc3NldFR4KC4uLmFyZ3MpIGFzIHRoaXNcbiAgfVxuXG4gIC8qKlxuICAgKiBDbGFzcyByZXByZXNlbnRpbmcgYW4gdW5zaWduZWQgQ3JlYXRlIEFzc2V0IHRyYW5zYWN0aW9uLlxuICAgKlxuICAgKiBAcGFyYW0gbmV0d29ya0lEIE9wdGlvbmFsIG5ldHdvcmtJRCwgW1tEZWZhdWx0TmV0d29ya0lEXV1cbiAgICogQHBhcmFtIGJsb2NrY2hhaW5JRCBPcHRpb25hbCBibG9ja2NoYWluSUQsIGRlZmF1bHQgQnVmZmVyLmFsbG9jKDMyLCAxNilcbiAgICogQHBhcmFtIG91dHMgT3B0aW9uYWwgYXJyYXkgb2YgdGhlIFtbVHJhbnNmZXJhYmxlT3V0cHV0XV1zXG4gICAqIEBwYXJhbSBpbnMgT3B0aW9uYWwgYXJyYXkgb2YgdGhlIFtbVHJhbnNmZXJhYmxlSW5wdXRdXXNcbiAgICogQHBhcmFtIG1lbW8gT3B0aW9uYWwge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gZm9yIHRoZSBtZW1vIGZpZWxkXG4gICAqIEBwYXJhbSBuYW1lIFN0cmluZyBmb3IgdGhlIGRlc2NyaXB0aXZlIG5hbWUgb2YgdGhlIGFzc2V0XG4gICAqIEBwYXJhbSBzeW1ib2wgU3RyaW5nIGZvciB0aGUgdGlja2VyIHN5bWJvbCBvZiB0aGUgYXNzZXRcbiAgICogQHBhcmFtIGRlbm9taW5hdGlvbiBPcHRpb25hbCBudW1iZXIgZm9yIHRoZSBkZW5vbWluYXRpb24gd2hpY2ggaXMgMTBeRC4gRCBtdXN0IGJlID49IDAgYW5kIDw9IDMyLiBFeDogJDEgQVZBWCA9IDEwXjkgJG5BVkFYXG4gICAqIEBwYXJhbSBpbml0aWFsU3RhdGUgT3B0aW9uYWwgW1tJbml0aWFsU3RhdGVzXV0gdGhhdCByZXByZXNlbnQgdGhlIGludGlhbCBzdGF0ZSBvZiBhIGNyZWF0ZWQgYXNzZXRcbiAgICovXG4gIGNvbnN0cnVjdG9yKFxuICAgIG5ldHdvcmtJRDogbnVtYmVyID0gRGVmYXVsdE5ldHdvcmtJRCxcbiAgICBibG9ja2NoYWluSUQ6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYygzMiwgMTYpLFxuICAgIG91dHM6IFRyYW5zZmVyYWJsZU91dHB1dFtdID0gdW5kZWZpbmVkLFxuICAgIGluczogVHJhbnNmZXJhYmxlSW5wdXRbXSA9IHVuZGVmaW5lZCxcbiAgICBtZW1vOiBCdWZmZXIgPSB1bmRlZmluZWQsXG4gICAgbmFtZTogc3RyaW5nID0gdW5kZWZpbmVkLFxuICAgIHN5bWJvbDogc3RyaW5nID0gdW5kZWZpbmVkLFxuICAgIGRlbm9taW5hdGlvbjogbnVtYmVyID0gdW5kZWZpbmVkLFxuICAgIGluaXRpYWxTdGF0ZTogSW5pdGlhbFN0YXRlcyA9IHVuZGVmaW5lZFxuICApIHtcbiAgICBzdXBlcihuZXR3b3JrSUQsIGJsb2NrY2hhaW5JRCwgb3V0cywgaW5zLCBtZW1vKVxuICAgIGlmIChcbiAgICAgIHR5cGVvZiBuYW1lID09PSBcInN0cmluZ1wiICYmXG4gICAgICB0eXBlb2Ygc3ltYm9sID09PSBcInN0cmluZ1wiICYmXG4gICAgICB0eXBlb2YgZGVub21pbmF0aW9uID09PSBcIm51bWJlclwiICYmXG4gICAgICBkZW5vbWluYXRpb24gPj0gMCAmJlxuICAgICAgZGVub21pbmF0aW9uIDw9IDMyICYmXG4gICAgICB0eXBlb2YgaW5pdGlhbFN0YXRlICE9PSBcInVuZGVmaW5lZFwiXG4gICAgKSB7XG4gICAgICB0aGlzLmluaXRpYWxTdGF0ZSA9IGluaXRpYWxTdGF0ZVxuICAgICAgdGhpcy5uYW1lID0gbmFtZVxuICAgICAgdGhpcy5zeW1ib2wgPSBzeW1ib2xcbiAgICAgIHRoaXMuZGVub21pbmF0aW9uLndyaXRlVUludDgoZGVub21pbmF0aW9uLCAwKVxuICAgIH1cbiAgfVxufVxuIl19