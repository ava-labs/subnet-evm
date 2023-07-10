"use strict";
/**
 * @packageDocumentation
 * @module API-AVM-MinterSet
 */
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.MinterSet = void 0;
const buffer_1 = require("buffer/");
const bintools_1 = __importDefault(require("../../utils/bintools"));
const serialization_1 = require("../../utils/serialization");
/**
 * @ignore
 */
const bintools = bintools_1.default.getInstance();
const serialization = serialization_1.Serialization.getInstance();
const decimalString = "decimalString";
const cb58 = "cb58";
const num = "number";
const buffer = "Buffer";
/**
 * Class for representing a threshold and set of minting addresses in Avalanche.
 *
 * @typeparam MinterSet including a threshold and array of addresses
 */
class MinterSet extends serialization_1.Serializable {
    /**
     *
     * @param threshold The number of signatures required to mint more of an asset by signing a minting transaction
     * @param minters Array of addresss which are authorized to sign a minting transaction
     */
    constructor(threshold = 1, minters = []) {
        super();
        this._typeName = "MinterSet";
        this._typeID = undefined;
        this.minters = [];
        /**
         * Returns the threshold.
         */
        this.getThreshold = () => {
            return this.threshold;
        };
        /**
         * Returns the minters.
         */
        this.getMinters = () => {
            return this.minters;
        };
        this._cleanAddresses = (addresses) => {
            let addrs = [];
            for (let i = 0; i < addresses.length; i++) {
                if (typeof addresses[`${i}`] === "string") {
                    addrs.push(bintools.stringToAddress(addresses[`${i}`]));
                }
                else if (addresses[`${i}`] instanceof buffer_1.Buffer) {
                    addrs.push(addresses[`${i}`]);
                }
            }
            return addrs;
        };
        this.threshold = threshold;
        this.minters = this._cleanAddresses(minters);
    }
    serialize(encoding = "hex") {
        const fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { threshold: serialization.encoder(this.threshold, encoding, num, decimalString, 4), minters: this.minters.map((m) => serialization.encoder(m, encoding, buffer, cb58, 20)) });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.threshold = serialization.decoder(fields["threshold"], encoding, decimalString, num, 4);
        this.minters = fields["minters"].map((m) => serialization.decoder(m, encoding, cb58, buffer, 20));
    }
}
exports.MinterSet = MinterSet;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibWludGVyc2V0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2FwaXMvYXZtL21pbnRlcnNldC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiO0FBQUE7OztHQUdHOzs7Ozs7QUFFSCxvQ0FBZ0M7QUFDaEMsb0VBQTJDO0FBQzNDLDZEQUtrQztBQUVsQzs7R0FFRztBQUNILE1BQU0sUUFBUSxHQUFhLGtCQUFRLENBQUMsV0FBVyxFQUFFLENBQUE7QUFDakQsTUFBTSxhQUFhLEdBQWtCLDZCQUFhLENBQUMsV0FBVyxFQUFFLENBQUE7QUFDaEUsTUFBTSxhQUFhLEdBQW1CLGVBQWUsQ0FBQTtBQUNyRCxNQUFNLElBQUksR0FBbUIsTUFBTSxDQUFBO0FBQ25DLE1BQU0sR0FBRyxHQUFtQixRQUFRLENBQUE7QUFDcEMsTUFBTSxNQUFNLEdBQW1CLFFBQVEsQ0FBQTtBQUV2Qzs7OztHQUlHO0FBQ0gsTUFBYSxTQUFVLFNBQVEsNEJBQVk7SUErRHpDOzs7O09BSUc7SUFDSCxZQUFZLFlBQW9CLENBQUMsRUFBRSxVQUErQixFQUFFO1FBQ2xFLEtBQUssRUFBRSxDQUFBO1FBcEVDLGNBQVMsR0FBRyxXQUFXLENBQUE7UUFDdkIsWUFBTyxHQUFHLFNBQVMsQ0FBQTtRQWlDbkIsWUFBTyxHQUFhLEVBQUUsQ0FBQTtRQUVoQzs7V0FFRztRQUNILGlCQUFZLEdBQUcsR0FBVyxFQUFFO1lBQzFCLE9BQU8sSUFBSSxDQUFDLFNBQVMsQ0FBQTtRQUN2QixDQUFDLENBQUE7UUFFRDs7V0FFRztRQUNILGVBQVUsR0FBRyxHQUFhLEVBQUU7WUFDMUIsT0FBTyxJQUFJLENBQUMsT0FBTyxDQUFBO1FBQ3JCLENBQUMsQ0FBQTtRQUVTLG9CQUFlLEdBQUcsQ0FBQyxTQUE4QixFQUFZLEVBQUU7WUFDdkUsSUFBSSxLQUFLLEdBQWEsRUFBRSxDQUFBO1lBQ3hCLEtBQUssSUFBSSxDQUFDLEdBQVcsQ0FBQyxFQUFFLENBQUMsR0FBRyxTQUFTLENBQUMsTUFBTSxFQUFFLENBQUMsRUFBRSxFQUFFO2dCQUNqRCxJQUFJLE9BQU8sU0FBUyxDQUFDLEdBQUcsQ0FBQyxFQUFFLENBQUMsS0FBSyxRQUFRLEVBQUU7b0JBQ3pDLEtBQUssQ0FBQyxJQUFJLENBQUMsUUFBUSxDQUFDLGVBQWUsQ0FBQyxTQUFTLENBQUMsR0FBRyxDQUFDLEVBQUUsQ0FBVyxDQUFDLENBQUMsQ0FBQTtpQkFDbEU7cUJBQU0sSUFBSSxTQUFTLENBQUMsR0FBRyxDQUFDLEVBQUUsQ0FBQyxZQUFZLGVBQU0sRUFBRTtvQkFDOUMsS0FBSyxDQUFDLElBQUksQ0FBQyxTQUFTLENBQUMsR0FBRyxDQUFDLEVBQUUsQ0FBVyxDQUFDLENBQUE7aUJBQ3hDO2FBQ0Y7WUFDRCxPQUFPLEtBQUssQ0FBQTtRQUNkLENBQUMsQ0FBQTtRQVNDLElBQUksQ0FBQyxTQUFTLEdBQUcsU0FBUyxDQUFBO1FBQzFCLElBQUksQ0FBQyxPQUFPLEdBQUcsSUFBSSxDQUFDLGVBQWUsQ0FBQyxPQUFPLENBQUMsQ0FBQTtJQUM5QyxDQUFDO0lBcEVELFNBQVMsQ0FBQyxXQUErQixLQUFLO1FBQzVDLE1BQU0sTUFBTSxHQUFXLEtBQUssQ0FBQyxTQUFTLENBQUMsUUFBUSxDQUFDLENBQUE7UUFDaEQsdUNBQ0ssTUFBTSxLQUNULFNBQVMsRUFBRSxhQUFhLENBQUMsT0FBTyxDQUM5QixJQUFJLENBQUMsU0FBUyxFQUNkLFFBQVEsRUFDUixHQUFHLEVBQ0gsYUFBYSxFQUNiLENBQUMsQ0FDRixFQUNELE9BQU8sRUFBRSxJQUFJLENBQUMsT0FBTyxDQUFDLEdBQUcsQ0FBQyxDQUFDLENBQUMsRUFBRSxFQUFFLENBQzlCLGFBQWEsQ0FBQyxPQUFPLENBQUMsQ0FBQyxFQUFFLFFBQVEsRUFBRSxNQUFNLEVBQUUsSUFBSSxFQUFFLEVBQUUsQ0FBQyxDQUNyRCxJQUNGO0lBQ0gsQ0FBQztJQUNELFdBQVcsQ0FBQyxNQUFjLEVBQUUsV0FBK0IsS0FBSztRQUM5RCxLQUFLLENBQUMsV0FBVyxDQUFDLE1BQU0sRUFBRSxRQUFRLENBQUMsQ0FBQTtRQUNuQyxJQUFJLENBQUMsU0FBUyxHQUFHLGFBQWEsQ0FBQyxPQUFPLENBQ3BDLE1BQU0sQ0FBQyxXQUFXLENBQUMsRUFDbkIsUUFBUSxFQUNSLGFBQWEsRUFDYixHQUFHLEVBQ0gsQ0FBQyxDQUNGLENBQUE7UUFDRCxJQUFJLENBQUMsT0FBTyxHQUFHLE1BQU0sQ0FBQyxTQUFTLENBQUMsQ0FBQyxHQUFHLENBQUMsQ0FBQyxDQUFTLEVBQUUsRUFBRSxDQUNqRCxhQUFhLENBQUMsT0FBTyxDQUFDLENBQUMsRUFBRSxRQUFRLEVBQUUsSUFBSSxFQUFFLE1BQU0sRUFBRSxFQUFFLENBQUMsQ0FDckQsQ0FBQTtJQUNILENBQUM7Q0F5Q0Y7QUF6RUQsOEJBeUVDIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgQVBJLUFWTS1NaW50ZXJTZXRcbiAqL1xuXG5pbXBvcnQgeyBCdWZmZXIgfSBmcm9tIFwiYnVmZmVyL1wiXG5pbXBvcnQgQmluVG9vbHMgZnJvbSBcIi4uLy4uL3V0aWxzL2JpbnRvb2xzXCJcbmltcG9ydCB7XG4gIFNlcmlhbGl6YWJsZSxcbiAgU2VyaWFsaXphdGlvbixcbiAgU2VyaWFsaXplZEVuY29kaW5nLFxuICBTZXJpYWxpemVkVHlwZVxufSBmcm9tIFwiLi4vLi4vdXRpbHMvc2VyaWFsaXphdGlvblwiXG5cbi8qKlxuICogQGlnbm9yZVxuICovXG5jb25zdCBiaW50b29sczogQmluVG9vbHMgPSBCaW5Ub29scy5nZXRJbnN0YW5jZSgpXG5jb25zdCBzZXJpYWxpemF0aW9uOiBTZXJpYWxpemF0aW9uID0gU2VyaWFsaXphdGlvbi5nZXRJbnN0YW5jZSgpXG5jb25zdCBkZWNpbWFsU3RyaW5nOiBTZXJpYWxpemVkVHlwZSA9IFwiZGVjaW1hbFN0cmluZ1wiXG5jb25zdCBjYjU4OiBTZXJpYWxpemVkVHlwZSA9IFwiY2I1OFwiXG5jb25zdCBudW06IFNlcmlhbGl6ZWRUeXBlID0gXCJudW1iZXJcIlxuY29uc3QgYnVmZmVyOiBTZXJpYWxpemVkVHlwZSA9IFwiQnVmZmVyXCJcblxuLyoqXG4gKiBDbGFzcyBmb3IgcmVwcmVzZW50aW5nIGEgdGhyZXNob2xkIGFuZCBzZXQgb2YgbWludGluZyBhZGRyZXNzZXMgaW4gQXZhbGFuY2hlLlxuICpcbiAqIEB0eXBlcGFyYW0gTWludGVyU2V0IGluY2x1ZGluZyBhIHRocmVzaG9sZCBhbmQgYXJyYXkgb2YgYWRkcmVzc2VzXG4gKi9cbmV4cG9ydCBjbGFzcyBNaW50ZXJTZXQgZXh0ZW5kcyBTZXJpYWxpemFibGUge1xuICBwcm90ZWN0ZWQgX3R5cGVOYW1lID0gXCJNaW50ZXJTZXRcIlxuICBwcm90ZWN0ZWQgX3R5cGVJRCA9IHVuZGVmaW5lZFxuXG4gIHNlcmlhbGl6ZShlbmNvZGluZzogU2VyaWFsaXplZEVuY29kaW5nID0gXCJoZXhcIik6IG9iamVjdCB7XG4gICAgY29uc3QgZmllbGRzOiBvYmplY3QgPSBzdXBlci5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgcmV0dXJuIHtcbiAgICAgIC4uLmZpZWxkcyxcbiAgICAgIHRocmVzaG9sZDogc2VyaWFsaXphdGlvbi5lbmNvZGVyKFxuICAgICAgICB0aGlzLnRocmVzaG9sZCxcbiAgICAgICAgZW5jb2RpbmcsXG4gICAgICAgIG51bSxcbiAgICAgICAgZGVjaW1hbFN0cmluZyxcbiAgICAgICAgNFxuICAgICAgKSxcbiAgICAgIG1pbnRlcnM6IHRoaXMubWludGVycy5tYXAoKG0pID0+XG4gICAgICAgIHNlcmlhbGl6YXRpb24uZW5jb2RlcihtLCBlbmNvZGluZywgYnVmZmVyLCBjYjU4LCAyMClcbiAgICAgIClcbiAgICB9XG4gIH1cbiAgZGVzZXJpYWxpemUoZmllbGRzOiBvYmplY3QsIGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKSB7XG4gICAgc3VwZXIuZGVzZXJpYWxpemUoZmllbGRzLCBlbmNvZGluZylcbiAgICB0aGlzLnRocmVzaG9sZCA9IHNlcmlhbGl6YXRpb24uZGVjb2RlcihcbiAgICAgIGZpZWxkc1tcInRocmVzaG9sZFwiXSxcbiAgICAgIGVuY29kaW5nLFxuICAgICAgZGVjaW1hbFN0cmluZyxcbiAgICAgIG51bSxcbiAgICAgIDRcbiAgICApXG4gICAgdGhpcy5taW50ZXJzID0gZmllbGRzW1wibWludGVyc1wiXS5tYXAoKG06IHN0cmluZykgPT5cbiAgICAgIHNlcmlhbGl6YXRpb24uZGVjb2RlcihtLCBlbmNvZGluZywgY2I1OCwgYnVmZmVyLCAyMClcbiAgICApXG4gIH1cblxuICBwcm90ZWN0ZWQgdGhyZXNob2xkOiBudW1iZXJcbiAgcHJvdGVjdGVkIG1pbnRlcnM6IEJ1ZmZlcltdID0gW11cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgdGhyZXNob2xkLlxuICAgKi9cbiAgZ2V0VGhyZXNob2xkID0gKCk6IG51bWJlciA9PiB7XG4gICAgcmV0dXJuIHRoaXMudGhyZXNob2xkXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgbWludGVycy5cbiAgICovXG4gIGdldE1pbnRlcnMgPSAoKTogQnVmZmVyW10gPT4ge1xuICAgIHJldHVybiB0aGlzLm1pbnRlcnNcbiAgfVxuXG4gIHByb3RlY3RlZCBfY2xlYW5BZGRyZXNzZXMgPSAoYWRkcmVzc2VzOiBzdHJpbmdbXSB8IEJ1ZmZlcltdKTogQnVmZmVyW10gPT4ge1xuICAgIGxldCBhZGRyczogQnVmZmVyW10gPSBbXVxuICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCBhZGRyZXNzZXMubGVuZ3RoOyBpKyspIHtcbiAgICAgIGlmICh0eXBlb2YgYWRkcmVzc2VzW2Ake2l9YF0gPT09IFwic3RyaW5nXCIpIHtcbiAgICAgICAgYWRkcnMucHVzaChiaW50b29scy5zdHJpbmdUb0FkZHJlc3MoYWRkcmVzc2VzW2Ake2l9YF0gYXMgc3RyaW5nKSlcbiAgICAgIH0gZWxzZSBpZiAoYWRkcmVzc2VzW2Ake2l9YF0gaW5zdGFuY2VvZiBCdWZmZXIpIHtcbiAgICAgICAgYWRkcnMucHVzaChhZGRyZXNzZXNbYCR7aX1gXSBhcyBCdWZmZXIpXG4gICAgICB9XG4gICAgfVxuICAgIHJldHVybiBhZGRyc1xuICB9XG5cbiAgLyoqXG4gICAqXG4gICAqIEBwYXJhbSB0aHJlc2hvbGQgVGhlIG51bWJlciBvZiBzaWduYXR1cmVzIHJlcXVpcmVkIHRvIG1pbnQgbW9yZSBvZiBhbiBhc3NldCBieSBzaWduaW5nIGEgbWludGluZyB0cmFuc2FjdGlvblxuICAgKiBAcGFyYW0gbWludGVycyBBcnJheSBvZiBhZGRyZXNzcyB3aGljaCBhcmUgYXV0aG9yaXplZCB0byBzaWduIGEgbWludGluZyB0cmFuc2FjdGlvblxuICAgKi9cbiAgY29uc3RydWN0b3IodGhyZXNob2xkOiBudW1iZXIgPSAxLCBtaW50ZXJzOiBzdHJpbmdbXSB8IEJ1ZmZlcltdID0gW10pIHtcbiAgICBzdXBlcigpXG4gICAgdGhpcy50aHJlc2hvbGQgPSB0aHJlc2hvbGRcbiAgICB0aGlzLm1pbnRlcnMgPSB0aGlzLl9jbGVhbkFkZHJlc3NlcyhtaW50ZXJzKVxuICB9XG59XG4iXX0=