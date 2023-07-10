"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.GenesisAsset = void 0;
/**
 * @packageDocumentation
 * @module API-AVM-GenesisAsset
 */
const buffer_1 = require("buffer/");
const bintools_1 = __importDefault(require("../../utils/bintools"));
const initialstates_1 = require("./initialstates");
const constants_1 = require("../../utils/constants");
const serialization_1 = require("../../utils/serialization");
const createassettx_1 = require("./createassettx");
const bn_js_1 = __importDefault(require("bn.js"));
/**
 * @ignore
 */
const serialization = serialization_1.Serialization.getInstance();
const bintools = bintools_1.default.getInstance();
const utf8 = "utf8";
const buffer = "Buffer";
const decimalString = "decimalString";
class GenesisAsset extends createassettx_1.CreateAssetTx {
    /**
     * Class representing a GenesisAsset
     *
     * @param assetAlias Optional String for the asset alias
     * @param name Optional String for the descriptive name of the asset
     * @param symbol Optional String for the ticker symbol of the asset
     * @param denomination Optional number for the denomination which is 10^D. D must be >= 0 and <= 32. Ex: $1 AVAX = 10^9 $nAVAX
     * @param initialState Optional [[InitialStates]] that represent the intial state of a created asset
     * @param memo Optional {@link https://github.com/feross/buffer|Buffer} for the memo field
     */
    constructor(assetAlias = undefined, name = undefined, symbol = undefined, denomination = undefined, initialState = undefined, memo = undefined) {
        super(constants_1.DefaultNetworkID, buffer_1.Buffer.alloc(32), [], [], memo);
        this._typeName = "GenesisAsset";
        this._codecID = undefined;
        this._typeID = undefined;
        this.assetAlias = "";
        /**
         * Returns the string representation of the assetAlias
         */
        this.getAssetAlias = () => this.assetAlias;
        if (typeof assetAlias === "string" &&
            typeof name === "string" &&
            typeof symbol === "string" &&
            typeof denomination === "number" &&
            denomination >= 0 &&
            denomination <= 32 &&
            typeof initialState !== "undefined") {
            this.assetAlias = assetAlias;
            this.name = name;
            this.symbol = symbol;
            this.denomination.writeUInt8(denomination, 0);
            this.initialState = initialState;
        }
    }
    serialize(encoding = "hex") {
        const fields = super.serialize(encoding);
        delete fields["blockchainID"];
        delete fields["outs"];
        delete fields["ins"];
        return Object.assign(Object.assign({}, fields), { assetAlias: serialization.encoder(this.assetAlias, encoding, utf8, utf8), name: serialization.encoder(this.name, encoding, utf8, utf8), symbol: serialization.encoder(this.symbol, encoding, utf8, utf8), denomination: serialization.encoder(this.denomination, encoding, buffer, decimalString, 1), initialState: this.initialState.serialize(encoding) });
    }
    deserialize(fields, encoding = "hex") {
        fields["blockchainID"] = buffer_1.Buffer.alloc(32, 16).toString("hex");
        fields["outs"] = [];
        fields["ins"] = [];
        super.deserialize(fields, encoding);
        this.assetAlias = serialization.decoder(fields["assetAlias"], encoding, utf8, utf8);
        this.name = serialization.decoder(fields["name"], encoding, utf8, utf8);
        this.symbol = serialization.decoder(fields["symbol"], encoding, utf8, utf8);
        this.denomination = serialization.decoder(fields["denomination"], encoding, decimalString, buffer, 1);
        this.initialState = new initialstates_1.InitialStates();
        this.initialState.deserialize(fields["initialState"], encoding);
    }
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} containing an [[GenesisAsset]], parses it, populates the class, and returns the length of the [[GenesisAsset]] in bytes.
     *
     * @param bytes A {@link https://github.com/feross/buffer|Buffer} containing a raw [[GenesisAsset]]
     *
     * @returns The length of the raw [[GenesisAsset]]
     *
     * @remarks assume not-checksummed
     */
    fromBuffer(bytes, offset = 0) {
        const assetAliasSize = bintools
            .copyFrom(bytes, offset, offset + 2)
            .readUInt16BE(0);
        offset += 2;
        this.assetAlias = bintools
            .copyFrom(bytes, offset, offset + assetAliasSize)
            .toString("utf8");
        offset += assetAliasSize;
        offset += super.fromBuffer(bytes, offset);
        return offset;
    }
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[GenesisAsset]].
     */
    toBuffer(networkID = constants_1.DefaultNetworkID) {
        // asset alias
        const assetAlias = this.getAssetAlias();
        const assetAliasbuffSize = buffer_1.Buffer.alloc(2);
        assetAliasbuffSize.writeUInt16BE(assetAlias.length, 0);
        let bsize = assetAliasbuffSize.length;
        let barr = [assetAliasbuffSize];
        const assetAliasbuff = buffer_1.Buffer.alloc(assetAlias.length);
        assetAliasbuff.write(assetAlias, 0, assetAlias.length, utf8);
        bsize += assetAliasbuff.length;
        barr.push(assetAliasbuff);
        const networkIDBuff = buffer_1.Buffer.alloc(4);
        networkIDBuff.writeUInt32BE(new bn_js_1.default(networkID).toNumber(), 0);
        bsize += networkIDBuff.length;
        barr.push(networkIDBuff);
        // Blockchain ID
        bsize += 32;
        barr.push(buffer_1.Buffer.alloc(32));
        // num Outputs
        bsize += 4;
        barr.push(buffer_1.Buffer.alloc(4));
        // num Inputs
        bsize += 4;
        barr.push(buffer_1.Buffer.alloc(4));
        // memo
        const memo = this.getMemo();
        const memobuffSize = buffer_1.Buffer.alloc(4);
        memobuffSize.writeUInt32BE(memo.length, 0);
        bsize += memobuffSize.length;
        barr.push(memobuffSize);
        bsize += memo.length;
        barr.push(memo);
        // asset name
        const name = this.getName();
        const namebuffSize = buffer_1.Buffer.alloc(2);
        namebuffSize.writeUInt16BE(name.length, 0);
        bsize += namebuffSize.length;
        barr.push(namebuffSize);
        const namebuff = buffer_1.Buffer.alloc(name.length);
        namebuff.write(name, 0, name.length, utf8);
        bsize += namebuff.length;
        barr.push(namebuff);
        // symbol
        const symbol = this.getSymbol();
        const symbolbuffSize = buffer_1.Buffer.alloc(2);
        symbolbuffSize.writeUInt16BE(symbol.length, 0);
        bsize += symbolbuffSize.length;
        barr.push(symbolbuffSize);
        const symbolbuff = buffer_1.Buffer.alloc(symbol.length);
        symbolbuff.write(symbol, 0, symbol.length, utf8);
        bsize += symbolbuff.length;
        barr.push(symbolbuff);
        // denomination
        const denomination = this.getDenomination();
        const denominationbuffSize = buffer_1.Buffer.alloc(1);
        denominationbuffSize.writeUInt8(denomination, 0);
        bsize += denominationbuffSize.length;
        barr.push(denominationbuffSize);
        bsize += this.initialState.toBuffer().length;
        barr.push(this.initialState.toBuffer());
        return buffer_1.Buffer.concat(barr, bsize);
    }
}
exports.GenesisAsset = GenesisAsset;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZ2VuZXNpc2Fzc2V0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2FwaXMvYXZtL2dlbmVzaXNhc3NldC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7QUFBQTs7O0dBR0c7QUFDSCxvQ0FBZ0M7QUFDaEMsb0VBQTJDO0FBQzNDLG1EQUErQztBQUMvQyxxREFBd0Q7QUFDeEQsNkRBSWtDO0FBQ2xDLG1EQUErQztBQUMvQyxrREFBc0I7QUFFdEI7O0dBRUc7QUFDSCxNQUFNLGFBQWEsR0FBa0IsNkJBQWEsQ0FBQyxXQUFXLEVBQUUsQ0FBQTtBQUNoRSxNQUFNLFFBQVEsR0FBYSxrQkFBUSxDQUFDLFdBQVcsRUFBRSxDQUFBO0FBQ2pELE1BQU0sSUFBSSxHQUFtQixNQUFNLENBQUE7QUFDbkMsTUFBTSxNQUFNLEdBQW1CLFFBQVEsQ0FBQTtBQUN2QyxNQUFNLGFBQWEsR0FBbUIsZUFBZSxDQUFBO0FBRXJELE1BQWEsWUFBYSxTQUFRLDZCQUFhO0lBNEo3Qzs7Ozs7Ozs7O09BU0c7SUFDSCxZQUNFLGFBQXFCLFNBQVMsRUFDOUIsT0FBZSxTQUFTLEVBQ3hCLFNBQWlCLFNBQVMsRUFDMUIsZUFBdUIsU0FBUyxFQUNoQyxlQUE4QixTQUFTLEVBQ3ZDLE9BQWUsU0FBUztRQUV4QixLQUFLLENBQUMsNEJBQWdCLEVBQUUsZUFBTSxDQUFDLEtBQUssQ0FBQyxFQUFFLENBQUMsRUFBRSxFQUFFLEVBQUUsRUFBRSxFQUFFLElBQUksQ0FBQyxDQUFBO1FBN0svQyxjQUFTLEdBQUcsY0FBYyxDQUFBO1FBQzFCLGFBQVEsR0FBRyxTQUFTLENBQUE7UUFDcEIsWUFBTyxHQUFHLFNBQVMsQ0FBQTtRQStDbkIsZUFBVSxHQUFXLEVBQUUsQ0FBQTtRQUVqQzs7V0FFRztRQUNILGtCQUFhLEdBQUcsR0FBVyxFQUFFLENBQUMsSUFBSSxDQUFDLFVBQVUsQ0FBQTtRQXdIM0MsSUFDRSxPQUFPLFVBQVUsS0FBSyxRQUFRO1lBQzlCLE9B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QUksQ0FBQyxDQUFBO1FBQ2hELEtBQUssSUFBSSxVQUFVLENBQUMsTUFBTSxDQUFBO1FBQzFCLElBQUksQ0FBQyxJQUFJLENBQUMsVUFBVSxDQUFDLENBQUE7UUFFckIsZUFBZTtRQUNmLE1BQU0sWUFBWSxHQUFXLElBQUksQ0FBQyxlQUFlLEVBQUUsQ0FBQTtRQUNuRCxNQUFNLG9CQUFvQixHQUFXLGVBQU0sQ0FBQyxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUE7UUFDcEQsb0JBQW9CLENBQUMsVUFBVSxDQUFDLFlBQVksRUFBRSxDQUFDLENBQUMsQ0FBQTtRQUNoRCxLQUFLLElBQUksb0JBQW9CLENBQUMsTUFBTSxDQUFBO1FBQ3BDLElBQUksQ0FBQyxJQUFJLENBQUMsb0JBQW9CLENBQUMsQ0FBQTtRQUUvQixLQUFLLElBQUksSUFBSSxDQUFDLFlBQVksQ0FBQyxRQUFRLEVBQUUsQ0FBQyxNQUFNLENBQUE7UUFDNUMsSUFBSSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsWUFBWSxDQUFDLFFBQVEsRUFBRSxDQUFDLENBQUE7UUFDdkMsT0FBTyxlQUFNLENBQUMsTUFBTSxDQUFDLElBQUksRUFBRSxLQUFLLENBQUMsQ0FBQTtJQUNuQyxDQUFDO0NBcUNGO0FBL0xELG9DQStMQyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogQHBhY2thZ2VEb2N1bWVudGF0aW9uXG4gKiBAbW9kdWxlIEFQSS1BVk0tR2VuZXNpc0Fzc2V0XG4gKi9cbmltcG9ydCB7IEJ1ZmZlciB9IGZyb20gXCJidWZmZXIvXCJcbmltcG9ydCBCaW5Ub29scyBmcm9tIFwiLi4vLi4vdXRpbHMvYmludG9vbHNcIlxuaW1wb3J0IHsgSW5pdGlhbFN0YXRlcyB9IGZyb20gXCIuL2luaXRpYWxzdGF0ZXNcIlxuaW1wb3J0IHsgRGVmYXVsdE5ldHdvcmtJRCB9IGZyb20gXCIuLi8uLi91dGlscy9jb25zdGFudHNcIlxuaW1wb3J0IHtcbiAgU2VyaWFsaXphdGlvbixcbiAgU2VyaWFsaXplZEVuY29kaW5nLFxuICBTZXJpYWxpemVkVHlwZVxufSBmcm9tIFwiLi4vLi4vdXRpbHMvc2VyaWFsaXphdGlvblwiXG5pbXBvcnQgeyBDcmVhdGVBc3NldFR4IH0gZnJvbSBcIi4vY3JlYXRlYXNzZXR0eFwiXG5pbXBvcnQgQk4gZnJvbSBcImJuLmpzXCJcblxuLyoqXG4gKiBAaWdub3JlXG4gKi9cbmNvbnN0IHNlcmlhbGl6YXRpb246IFNlcmlhbGl6YXRpb24gPSBTZXJpYWxpemF0aW9uLmdldEluc3RhbmNlKClcbmNvbnN0IGJpbnRvb2xzOiBCaW5Ub29scyA9IEJpblRvb2xzLmdldEluc3RhbmNlKClcbmNvbnN0IHV0Zjg6IFNlcmlhbGl6ZWRUeXBlID0gXCJ1dGY4XCJcbmNvbnN0IGJ1ZmZlcjogU2VyaWFsaXplZFR5cGUgPSBcIkJ1ZmZlclwiXG5jb25zdCBkZWNpbWFsU3RyaW5nOiBTZXJpYWxpemVkVHlwZSA9IFwiZGVjaW1hbFN0cmluZ1wiXG5cbmV4cG9ydCBjbGFzcyBHZW5lc2lzQXNzZXQgZXh0ZW5kcyBDcmVhdGVBc3NldFR4IHtcbiAgcHJvdGVjdGVkIF90eXBlTmFtZSA9IFwiR2VuZXNpc0Fzc2V0XCJcbiAgcHJvdGVjdGVkIF9jb2RlY0lEID0gdW5kZWZpbmVkXG4gIHByb3RlY3RlZCBfdHlwZUlEID0gdW5kZWZpbmVkXG5cbiAgc2VyaWFsaXplKGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKTogb2JqZWN0IHtcbiAgICBjb25zdCBmaWVsZHM6IG9iamVjdCA9IHN1cGVyLnNlcmlhbGl6ZShlbmNvZGluZylcbiAgICBkZWxldGUgZmllbGRzW1wiYmxvY2tjaGFpbklEXCJdXG4gICAgZGVsZXRlIGZpZWxkc1tcIm91dHNcIl1cbiAgICBkZWxldGUgZmllbGRzW1wiaW5zXCJdXG4gICAgcmV0dXJuIHtcbiAgICAgIC4uLmZpZWxkcyxcbiAgICAgIGFzc2V0QWxpYXM6IHNlcmlhbGl6YXRpb24uZW5jb2Rlcih0aGlzLmFzc2V0QWxpYXMsIGVuY29kaW5nLCB1dGY4LCB1dGY4KSxcbiAgICAgIG5hbWU6IHNlcmlhbGl6YXRpb24uZW5jb2Rlcih0aGlzLm5hbWUsIGVuY29kaW5nLCB1dGY4LCB1dGY4KSxcbiAgICAgIHN5bWJvbDogc2VyaWFsaXphdGlvbi5lbmNvZGVyKHRoaXMuc3ltYm9sLCBlbmNvZGluZywgdXRmOCwgdXRmOCksXG4gICAgICBkZW5vbWluYXRpb246IHNlcmlhbGl6YXRpb24uZW5jb2RlcihcbiAgICAgICAgdGhpcy5kZW5vbWluYXRpb24sXG4gICAgICAgIGVuY29kaW5nLFxuICAgICAgICBidWZmZXIsXG4gICAgICAgIGRlY2ltYWxTdHJpbmcsXG4gICAgICAgIDFcbiAgICAgICksXG4gICAgICBpbml0aWFsU3RhdGU6IHRoaXMuaW5pdGlhbFN0YXRlLnNlcmlhbGl6ZShlbmNvZGluZylcbiAgICB9XG4gIH1cblxuICBkZXNlcmlhbGl6ZShmaWVsZHM6IG9iamVjdCwgZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpIHtcbiAgICBmaWVsZHNbXCJibG9ja2NoYWluSURcIl0gPSBCdWZmZXIuYWxsb2MoMzIsIDE2KS50b1N0cmluZyhcImhleFwiKVxuICAgIGZpZWxkc1tcIm91dHNcIl0gPSBbXVxuICAgIGZpZWxkc1tcImluc1wiXSA9IFtdXG4gICAgc3VwZXIuZGVzZXJpYWxpemUoZmllbGRzLCBlbmNvZGluZylcbiAgICB0aGlzLmFzc2V0QWxpYXMgPSBzZXJpYWxpemF0aW9uLmRlY29kZXIoXG4gICAgICBmaWVsZHNbXCJhc3NldEFsaWFzXCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICB1dGY4LFxuICAgICAgdXRmOFxuICAgIClcbiAgICB0aGlzLm5hbWUgPSBzZXJpYWxpemF0aW9uLmRlY29kZXIoZmllbGRzW1wibmFtZVwiXSwgZW5jb2RpbmcsIHV0ZjgsIHV0ZjgpXG4gICAgdGhpcy5zeW1ib2wgPSBzZXJpYWxpemF0aW9uLmRlY29kZXIoZmllbGRzW1wic3ltYm9sXCJdLCBlbmNvZGluZywgdXRmOCwgdXRmOClcbiAgICB0aGlzLmRlbm9taW5hdGlvbiA9IHNlcmlhbGl6YXRpb24uZGVjb2RlcihcbiAgICAgIGZpZWxkc1tcImRlbm9taW5hdGlvblwiXSxcbiAgICAgIGVuY29kaW5nLFxuICAgICAgZGVjaW1hbFN0cmluZyxcbiAgICAgIGJ1ZmZlcixcbiAgICAgIDFcbiAgICApXG4gICAgdGhpcy5pbml0aWFsU3RhdGUgPSBuZXcgSW5pdGlhbFN0YXRlcygpXG4gICAgdGhpcy5pbml0aWFsU3RhdGUuZGVzZXJpYWxpemUoZmllbGRzW1wiaW5pdGlhbFN0YXRlXCJdLCBlbmNvZGluZylcbiAgfVxuXG4gIHByb3RlY3RlZCBhc3NldEFsaWFzOiBzdHJpbmcgPSBcIlwiXG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIHN0cmluZyByZXByZXNlbnRhdGlvbiBvZiB0aGUgYXNzZXRBbGlhc1xuICAgKi9cbiAgZ2V0QXNzZXRBbGlhcyA9ICgpOiBzdHJpbmcgPT4gdGhpcy5hc3NldEFsaWFzXG5cbiAgLyoqXG4gICAqIFRha2VzIGEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gY29udGFpbmluZyBhbiBbW0dlbmVzaXNBc3NldF1dLCBwYXJzZXMgaXQsIHBvcHVsYXRlcyB0aGUgY2xhc3MsIGFuZCByZXR1cm5zIHRoZSBsZW5ndGggb2YgdGhlIFtbR2VuZXNpc0Fzc2V0XV0gaW4gYnl0ZXMuXG4gICAqXG4gICAqIEBwYXJhbSBieXRlcyBBIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IGNvbnRhaW5pbmcgYSByYXcgW1tHZW5lc2lzQXNzZXRdXVxuICAgKlxuICAgKiBAcmV0dXJucyBUaGUgbGVuZ3RoIG9mIHRoZSByYXcgW1tHZW5lc2lzQXNzZXRdXVxuICAgKlxuICAgKiBAcmVtYXJrcyBhc3N1bWUgbm90LWNoZWNrc3VtbWVkXG4gICAqL1xuICBmcm9tQnVmZmVyKGJ5dGVzOiBCdWZmZXIsIG9mZnNldDogbnVtYmVyID0gMCk6IG51bWJlciB7XG4gICAgY29uc3QgYXNzZXRBbGlhc1NpemU6IG51bWJlciA9IGJpbnRvb2xzXG4gICAgICAuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgMilcbiAgICAgIC5yZWFkVUludDE2QkUoMClcbiAgICBvZmZzZXQgKz0gMlxuICAgIHRoaXMuYXNzZXRBbGlhcyA9IGJpbnRvb2xzXG4gICAgICAuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgYXNzZXRBbGlhc1NpemUpXG4gICAgICAudG9TdHJpbmcoXCJ1dGY4XCIpXG4gICAgb2Zmc2V0ICs9IGFzc2V0QWxpYXNTaXplXG4gICAgb2Zmc2V0ICs9IHN1cGVyLmZyb21CdWZmZXIoYnl0ZXMsIG9mZnNldClcbiAgICByZXR1cm4gb2Zmc2V0XG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IHJlcHJlc2VudGF0aW9uIG9mIHRoZSBbW0dlbmVzaXNBc3NldF1dLlxuICAgKi9cbiAgdG9CdWZmZXIobmV0d29ya0lEOiBudW1iZXIgPSBEZWZhdWx0TmV0d29ya0lEKTogQnVmZmVyIHtcbiAgICAvLyBhc3NldCBhbGlhc1xuICAgIGNvbnN0IGFzc2V0QWxpYXM6IHN0cmluZyA9IHRoaXMuZ2V0QXNzZXRBbGlhcygpXG4gICAgY29uc3QgYXNzZXRBbGlhc2J1ZmZTaXplOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoMilcbiAgICBhc3NldEFsaWFzYnVmZlNpemUud3JpdGVVSW50MTZCRShhc3NldEFsaWFzLmxlbmd0aCwgMClcbiAgICBsZXQgYnNpemU6IG51bWJlciA9IGFzc2V0QWxpYXNidWZmU2l6ZS5sZW5ndGhcbiAgICBsZXQgYmFycjogQnVmZmVyW10gPSBbYXNzZXRBbGlhc2J1ZmZTaXplXVxuICAgIGNvbnN0IGFzc2V0QWxpYXNidWZmOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoYXNzZXRBbGlhcy5sZW5ndGgpXG4gICAgYXNzZXRBbGlhc2J1ZmYud3JpdGUoYXNzZXRBbGlhcywgMCwgYXNzZXRBbGlhcy5sZW5ndGgsIHV0ZjgpXG4gICAgYnNpemUgKz0gYXNzZXRBbGlhc2J1ZmYubGVuZ3RoXG4gICAgYmFyci5wdXNoKGFzc2V0QWxpYXNidWZmKVxuXG4gICAgY29uc3QgbmV0d29ya0lEQnVmZjogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDQpXG4gICAgbmV0d29ya0lEQnVmZi53cml0ZVVJbnQzMkJFKG5ldyBCTihuZXR3b3JrSUQpLnRvTnVtYmVyKCksIDApXG4gICAgYnNpemUgKz0gbmV0d29ya0lEQnVmZi5sZW5ndGhcbiAgICBiYXJyLnB1c2gobmV0d29ya0lEQnVmZilcblxuICAgIC8vIEJsb2NrY2hhaW4gSURcbiAgICBic2l6ZSArPSAzMlxuICAgIGJhcnIucHVzaChCdWZmZXIuYWxsb2MoMzIpKVxuXG4gICAgLy8gbnVtIE91dHB1dHNcbiAgICBic2l6ZSArPSA0XG4gICAgYmFyci5wdXNoKEJ1ZmZlci5hbGxvYyg0KSlcblxuICAgIC8vIG51bSBJbnB1dHNcbiAgICBic2l6ZSArPSA0XG4gICAgYmFyci5wdXNoKEJ1ZmZlci5hbGxvYyg0KSlcblxuICAgIC8vIG1lbW9cbiAgICBjb25zdCBtZW1vOiBCdWZmZXIgPSB0aGlzLmdldE1lbW8oKVxuICAgIGNvbnN0IG1lbW9idWZmU2l6ZTogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDQpXG4gICAgbWVtb2J1ZmZTaXplLndyaXRlVUludDMyQkUobWVtby5sZW5ndGgsIDApXG4gICAgYnNpemUgKz0gbWVtb2J1ZmZTaXplLmxlbmd0aFxuICAgIGJhcnIucHVzaChtZW1vYnVmZlNpemUpXG5cbiAgICBic2l6ZSArPSBtZW1vLmxlbmd0aFxuICAgIGJhcnIucHVzaChtZW1vKVxuXG4gICAgLy8gYXNzZXQgbmFtZVxuICAgIGNvbnN0IG5hbWU6IHN0cmluZyA9IHRoaXMuZ2V0TmFtZSgpXG4gICAgY29uc3QgbmFtZWJ1ZmZTaXplOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoMilcbiAgICBuYW1lYnVmZlNpemUud3JpdGVVSW50MTZCRShuYW1lLmxlbmd0aCwgMClcbiAgICBic2l6ZSArPSBuYW1lYnVmZlNpemUubGVuZ3RoXG4gICAgYmFyci5wdXNoKG5hbWVidWZmU2l6ZSlcbiAgICBjb25zdCBuYW1lYnVmZjogQnVmZmVyID0gQnVmZmVyLmFsbG9jKG5hbWUubGVuZ3RoKVxuICAgIG5hbWVidWZmLndyaXRlKG5hbWUsIDAsIG5hbWUubGVuZ3RoLCB1dGY4KVxuICAgIGJzaXplICs9IG5hbWVidWZmLmxlbmd0aFxuICAgIGJhcnIucHVzaChuYW1lYnVmZilcblxuICAgIC8vIHN5bWJvbFxuICAgIGNvbnN0IHN5bWJvbDogc3RyaW5nID0gdGhpcy5nZXRTeW1ib2woKVxuICAgIGNvbnN0IHN5bWJvbGJ1ZmZTaXplOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoMilcbiAgICBzeW1ib2xidWZmU2l6ZS53cml0ZVVJbnQxNkJFKHN5bWJvbC5sZW5ndGgsIDApXG4gICAgYnNpemUgKz0gc3ltYm9sYnVmZlNpemUubGVuZ3RoXG4gICAgYmFyci5wdXNoKHN5bWJvbGJ1ZmZTaXplKVxuXG4gICAgY29uc3Qgc3ltYm9sYnVmZjogQnVmZmVyID0gQnVmZmVyLmFsbG9jKHN5bWJvbC5sZW5ndGgpXG4gICAgc3ltYm9sYnVmZi53cml0ZShzeW1ib2wsIDAsIHN5bWJvbC5sZW5ndGgsIHV0ZjgpXG4gICAgYnNpemUgKz0gc3ltYm9sYnVmZi5sZW5ndGhcbiAgICBiYXJyLnB1c2goc3ltYm9sYnVmZilcblxuICAgIC8vIGRlbm9taW5hdGlvblxuICAgIGNvbnN0IGRlbm9taW5hdGlvbjogbnVtYmVyID0gdGhpcy5nZXREZW5vbWluYXRpb24oKVxuICAgIGNvbnN0IGRlbm9taW5hdGlvbmJ1ZmZTaXplOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoMSlcbiAgICBkZW5vbWluYXRpb25idWZmU2l6ZS53cml0ZVVJbnQ4KGRlbm9taW5hdGlvbiwgMClcbiAgICBic2l6ZSArPSBkZW5vbWluYXRpb25idWZmU2l6ZS5sZW5ndGhcbiAgICBiYXJyLnB1c2goZGVub21pbmF0aW9uYnVmZlNpemUpXG5cbiAgICBic2l6ZSArPSB0aGlzLmluaXRpYWxTdGF0ZS50b0J1ZmZlcigpLmxlbmd0aFxuICAgIGJhcnIucHVzaCh0aGlzLmluaXRpYWxTdGF0ZS50b0J1ZmZlcigpKVxuICAgIHJldHVybiBCdWZmZXIuY29uY2F0KGJhcnIsIGJzaXplKVxuICB9XG5cbiAgLyoqXG4gICAqIENsYXNzIHJlcHJlc2VudGluZyBhIEdlbmVzaXNBc3NldFxuICAgKlxuICAgKiBAcGFyYW0gYXNzZXRBbGlhcyBPcHRpb25hbCBTdHJpbmcgZm9yIHRoZSBhc3NldCBhbGlhc1xuICAgKiBAcGFyYW0gbmFtZSBPcHRpb25hbCBTdHJpbmcgZm9yIHRoZSBkZXNjcmlwdGl2ZSBuYW1lIG9mIHRoZSBhc3NldFxuICAgKiBAcGFyYW0gc3ltYm9sIE9wdGlvbmFsIFN0cmluZyBmb3IgdGhlIHRpY2tlciBzeW1ib2wgb2YgdGhlIGFzc2V0XG4gICAqIEBwYXJhbSBkZW5vbWluYXRpb24gT3B0aW9uYWwgbnVtYmVyIGZvciB0aGUgZGVub21pbmF0aW9uIHdoaWNoIGlzIDEwXkQuIEQgbXVzdCBiZSA+PSAwIGFuZCA8PSAzMi4gRXg6ICQxIEFWQVggPSAxMF45ICRuQVZBWFxuICAgKiBAcGFyYW0gaW5pdGlhbFN0YXRlIE9wdGlvbmFsIFtbSW5pdGlhbFN0YXRlc11dIHRoYXQgcmVwcmVzZW50IHRoZSBpbnRpYWwgc3RhdGUgb2YgYSBjcmVhdGVkIGFzc2V0XG4gICAqIEBwYXJhbSBtZW1vIE9wdGlvbmFsIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IGZvciB0aGUgbWVtbyBmaWVsZFxuICAgKi9cbiAgY29uc3RydWN0b3IoXG4gICAgYXNzZXRBbGlhczogc3RyaW5nID0gdW5kZWZpbmVkLFxuICAgIG5hbWU6IHN0cmluZyA9IHVuZGVmaW5lZCxcbiAgICBzeW1ib2w6IHN0cmluZyA9IHVuZGVmaW5lZCxcbiAgICBkZW5vbWluYXRpb246IG51bWJlciA9IHVuZGVmaW5lZCxcbiAgICBpbml0aWFsU3RhdGU6IEluaXRpYWxTdGF0ZXMgPSB1bmRlZmluZWQsXG4gICAgbWVtbzogQnVmZmVyID0gdW5kZWZpbmVkXG4gICkge1xuICAgIHN1cGVyKERlZmF1bHROZXR3b3JrSUQsIEJ1ZmZlci5hbGxvYygzMiksIFtdLCBbXSwgbWVtbylcbiAgICBpZiAoXG4gICAgICB0eXBlb2YgYXNzZXRBbGlhcyA9PT0gXCJzdHJpbmdcIiAmJlxuICAgICAgdHlwZW9mIG5hbWUgPT09IFwic3RyaW5nXCIgJiZcbiAgICAgIHR5cGVvZiBzeW1ib2wgPT09IFwic3RyaW5nXCIgJiZcbiAgICAgIHR5cGVvZiBkZW5vbWluYXRpb24gPT09IFwibnVtYmVyXCIgJiZcbiAgICAgIGRlbm9taW5hdGlvbiA+PSAwICYmXG4gICAgICBkZW5vbWluYXRpb24gPD0gMzIgJiZcbiAgICAgIHR5cGVvZiBpbml0aWFsU3RhdGUgIT09IFwidW5kZWZpbmVkXCJcbiAgICApIHtcbiAgICAgIHRoaXMuYXNzZXRBbGlhcyA9IGFzc2V0QWxpYXNcbiAgICAgIHRoaXMubmFtZSA9IG5hbWVcbiAgICAgIHRoaXMuc3ltYm9sID0gc3ltYm9sXG4gICAgICB0aGlzLmRlbm9taW5hdGlvbi53cml0ZVVJbnQ4KGRlbm9taW5hdGlvbiwgMClcbiAgICAgIHRoaXMuaW5pdGlhbFN0YXRlID0gaW5pdGlhbFN0YXRlXG4gICAgfVxuICB9XG59XG4iXX0=