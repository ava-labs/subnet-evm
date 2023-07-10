"use strict";
/**
 * @packageDocumentation
 * @module Common-NBytes
 */
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.NBytes = void 0;
const bintools_1 = __importDefault(require("../utils/bintools"));
const serialization_1 = require("../utils/serialization");
const errors_1 = require("../utils/errors");
/**
 * @ignore
 */
const bintools = bintools_1.default.getInstance();
const serialization = serialization_1.Serialization.getInstance();
/**
 * Abstract class that implements basic functionality for managing a
 * {@link https://github.com/feross/buffer|Buffer} of an exact length.
 *
 * Create a class that extends this one and override bsize to make it validate for exactly
 * the correct length.
 */
class NBytes extends serialization_1.Serializable {
    constructor() {
        super(...arguments);
        this._typeName = "NBytes";
        this._typeID = undefined;
        /**
         * Returns the length of the {@link https://github.com/feross/buffer|Buffer}.
         *
         * @returns The exact length requirement of this class
         */
        this.getSize = () => this.bsize;
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { bsize: serialization.encoder(this.bsize, encoding, "number", "decimalString", 4), bytes: serialization.encoder(this.bytes, encoding, "Buffer", "hex", this.bsize) });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.bsize = serialization.decoder(fields["bsize"], encoding, "decimalString", "number", 4);
        this.bytes = serialization.decoder(fields["bytes"], encoding, "hex", "Buffer", this.bsize);
    }
    /**
     * Takes a base-58 encoded string, verifies its length, and stores it.
     *
     * @returns The size of the {@link https://github.com/feross/buffer|Buffer}
     */
    fromString(b58str) {
        try {
            this.fromBuffer(bintools.b58ToBuffer(b58str));
        }
        catch (e) {
            /* istanbul ignore next */
            const emsg = `Error - NBytes.fromString: ${e}`;
            /* istanbul ignore next */
            throw new Error(emsg);
        }
        return this.bsize;
    }
    /**
     * Takes a [[Buffer]], verifies its length, and stores it.
     *
     * @returns The size of the {@link https://github.com/feross/buffer|Buffer}
     */
    fromBuffer(buff, offset = 0) {
        try {
            if (buff.length - offset < this.bsize) {
                /* istanbul ignore next */
                throw new errors_1.BufferSizeError("Error - NBytes.fromBuffer: not enough space available in buffer.");
            }
            this.bytes = bintools.copyFrom(buff, offset, offset + this.bsize);
        }
        catch (e) {
            /* istanbul ignore next */
            const emsg = `Error - NBytes.fromBuffer: ${e}`;
            /* istanbul ignore next */
            throw new Error(emsg);
        }
        return offset + this.bsize;
    }
    /**
     * @returns A reference to the stored {@link https://github.com/feross/buffer|Buffer}
     */
    toBuffer() {
        return this.bytes;
    }
    /**
     * @returns A base-58 string of the stored {@link https://github.com/feross/buffer|Buffer}
     */
    toString() {
        return bintools.bufferToB58(this.toBuffer());
    }
}
exports.NBytes = NBytes;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibmJ5dGVzLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL2NvbW1vbi9uYnl0ZXMudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IjtBQUFBOzs7R0FHRzs7Ozs7O0FBR0gsaUVBQXdDO0FBQ3hDLDBEQUkrQjtBQUMvQiw0Q0FBaUQ7QUFFakQ7O0dBRUc7QUFDSCxNQUFNLFFBQVEsR0FBYSxrQkFBUSxDQUFDLFdBQVcsRUFBRSxDQUFBO0FBQ2pELE1BQU0sYUFBYSxHQUFrQiw2QkFBYSxDQUFDLFdBQVcsRUFBRSxDQUFBO0FBRWhFOzs7Ozs7R0FNRztBQUNILE1BQXNCLE1BQU8sU0FBUSw0QkFBWTtJQUFqRDs7UUFDWSxjQUFTLEdBQUcsUUFBUSxDQUFBO1FBQ3BCLFlBQU8sR0FBRyxTQUFTLENBQUE7UUEyQzdCOzs7O1dBSUc7UUFDSCxZQUFPLEdBQUcsR0FBRyxFQUFFLENBQUMsSUFBSSxDQUFDLEtBQUssQ0FBQTtJQTJENUIsQ0FBQztJQXpHQyxTQUFTLENBQUMsV0FBK0IsS0FBSztRQUM1QyxJQUFJLE1BQU0sR0FBVyxLQUFLLENBQUMsU0FBUyxDQUFDLFFBQVEsQ0FBQyxDQUFBO1FBQzlDLHVDQUNLLE1BQU0sS0FDVCxLQUFLLEVBQUUsYUFBYSxDQUFDLE9BQU8sQ0FDMUIsSUFBSSxDQUFDLEtBQUssRUFDVixRQUFRLEVBQ1IsUUFBUSxFQUNSLGVBQWUsRUFDZixDQUFDLENBQ0YsRUFDRCxLQUFLLEVBQUUsYUFBYSxDQUFDLE9BQU8sQ0FDMUIsSUFBSSxDQUFDLEtBQUssRUFDVixRQUFRLEVBQ1IsUUFBUSxFQUNSLEtBQUssRUFDTCxJQUFJLENBQUMsS0FBSyxDQUNYLElBQ0Y7SUFDSCxDQUFDO0lBQ0QsV0FBVyxDQUFDLE1BQWMsRUFBRSxXQUErQixLQUFLO1FBQzlELEtBQUssQ0FBQyxXQUFXLENBQUMsTUFBTSxFQUFFLFFBQVEsQ0FBQyxDQUFBO1FBQ25DLElBQUksQ0FBQyxLQUFLLEdBQUcsYUFBYSxDQUFDLE9BQU8sQ0FDaEMsTUFBTSxDQUFDLE9BQU8sQ0FBQyxFQUNmLFFBQVEsRUFDUixlQUFlLEVBQ2YsUUFBUSxFQUNSLENBQUMsQ0FDRixDQUFBO1FBQ0QsSUFBSSxDQUFDLEtBQUssR0FBRyxhQUFhLENBQUMsT0FBTyxDQUNoQyxNQUFNLENBQUMsT0FBTyxDQUFDLEVBQ2YsUUFBUSxFQUNSLEtBQUssRUFDTCxRQUFRLEVBQ1IsSUFBSSxDQUFDLEtBQUssQ0FDWCxDQUFBO0lBQ0gsQ0FBQztJQVlEOzs7O09BSUc7SUFDSCxVQUFVLENBQUMsTUFBYztRQUN2QixJQUFJO1lBQ0YsSUFBSSxDQUFDLFVBQVUsQ0FBQyxRQUFRLENBQUMsV0FBVyxDQUFDLE1BQU0sQ0FBQyxDQUFDLENBQUE7U0FDOUM7UUFBQyxPQUFPLENBQUMsRUFBRTtZQUNWLDBCQUEwQjtZQUMxQixNQUFNLElBQUksR0FBVyw4QkFBOEIsQ0FBQyxFQUFFLENBQUE7WUFDdEQsMEJBQTBCO1lBQzFCLE1BQU0sSUFBSSxLQUFLLENBQUMsSUFBSSxDQUFDLENBQUE7U0FDdEI7UUFDRCxPQUFPLElBQUksQ0FBQyxLQUFLLENBQUE7SUFDbkIsQ0FBQztJQUVEOzs7O09BSUc7SUFDSCxVQUFVLENBQUMsSUFBWSxFQUFFLFNBQWlCLENBQUM7UUFDekMsSUFBSTtZQUNGLElBQUksSUFBSSxDQUFDLE1BQU0sR0FBRyxNQUFNLEdBQUcsSUFBSSxDQUFDLEtBQUssRUFBRTtnQkFDckMsMEJBQTBCO2dCQUMxQixNQUFNLElBQUksd0JBQWUsQ0FDdkIsa0VBQWtFLENBQ25FLENBQUE7YUFDRjtZQUVELElBQUksQ0FBQyxLQUFLLEdBQUcsUUFBUSxDQUFDLFFBQVEsQ0FBQyxJQUFJLEVBQUUsTUFBTSxFQUFFLE1BQU0sR0FBRyxJQUFJLENBQUMsS0FBSyxDQUFDLENBQUE7U0FDbEU7UUFBQyxPQUFPLENBQUMsRUFBRTtZQUNWLDBCQUEwQjtZQUMxQixNQUFNLElBQUksR0FBVyw4QkFBOEIsQ0FBQyxFQUFFLENBQUE7WUFDdEQsMEJBQTBCO1lBQzFCLE1BQU0sSUFBSSxLQUFLLENBQUMsSUFBSSxDQUFDLENBQUE7U0FDdEI7UUFDRCxPQUFPLE1BQU0sR0FBRyxJQUFJLENBQUMsS0FBSyxDQUFBO0lBQzVCLENBQUM7SUFFRDs7T0FFRztJQUNILFFBQVE7UUFDTixPQUFPLElBQUksQ0FBQyxLQUFLLENBQUE7SUFDbkIsQ0FBQztJQUVEOztPQUVHO0lBQ0gsUUFBUTtRQUNOLE9BQU8sUUFBUSxDQUFDLFdBQVcsQ0FBQyxJQUFJLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQTtJQUM5QyxDQUFDO0NBSUY7QUE3R0Qsd0JBNkdDIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgQ29tbW9uLU5CeXRlc1xuICovXG5cbmltcG9ydCB7IEJ1ZmZlciB9IGZyb20gXCJidWZmZXIvXCJcbmltcG9ydCBCaW5Ub29scyBmcm9tIFwiLi4vdXRpbHMvYmludG9vbHNcIlxuaW1wb3J0IHtcbiAgU2VyaWFsaXphYmxlLFxuICBTZXJpYWxpemF0aW9uLFxuICBTZXJpYWxpemVkRW5jb2Rpbmdcbn0gZnJvbSBcIi4uL3V0aWxzL3NlcmlhbGl6YXRpb25cIlxuaW1wb3J0IHsgQnVmZmVyU2l6ZUVycm9yIH0gZnJvbSBcIi4uL3V0aWxzL2Vycm9yc1wiXG5cbi8qKlxuICogQGlnbm9yZVxuICovXG5jb25zdCBiaW50b29sczogQmluVG9vbHMgPSBCaW5Ub29scy5nZXRJbnN0YW5jZSgpXG5jb25zdCBzZXJpYWxpemF0aW9uOiBTZXJpYWxpemF0aW9uID0gU2VyaWFsaXphdGlvbi5nZXRJbnN0YW5jZSgpXG5cbi8qKlxuICogQWJzdHJhY3QgY2xhc3MgdGhhdCBpbXBsZW1lbnRzIGJhc2ljIGZ1bmN0aW9uYWxpdHkgZm9yIG1hbmFnaW5nIGFcbiAqIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IG9mIGFuIGV4YWN0IGxlbmd0aC5cbiAqXG4gKiBDcmVhdGUgYSBjbGFzcyB0aGF0IGV4dGVuZHMgdGhpcyBvbmUgYW5kIG92ZXJyaWRlIGJzaXplIHRvIG1ha2UgaXQgdmFsaWRhdGUgZm9yIGV4YWN0bHlcbiAqIHRoZSBjb3JyZWN0IGxlbmd0aC5cbiAqL1xuZXhwb3J0IGFic3RyYWN0IGNsYXNzIE5CeXRlcyBleHRlbmRzIFNlcmlhbGl6YWJsZSB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIk5CeXRlc1wiXG4gIHByb3RlY3RlZCBfdHlwZUlEID0gdW5kZWZpbmVkXG5cbiAgc2VyaWFsaXplKGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImhleFwiKTogb2JqZWN0IHtcbiAgICBsZXQgZmllbGRzOiBvYmplY3QgPSBzdXBlci5zZXJpYWxpemUoZW5jb2RpbmcpXG4gICAgcmV0dXJuIHtcbiAgICAgIC4uLmZpZWxkcyxcbiAgICAgIGJzaXplOiBzZXJpYWxpemF0aW9uLmVuY29kZXIoXG4gICAgICAgIHRoaXMuYnNpemUsXG4gICAgICAgIGVuY29kaW5nLFxuICAgICAgICBcIm51bWJlclwiLFxuICAgICAgICBcImRlY2ltYWxTdHJpbmdcIixcbiAgICAgICAgNFxuICAgICAgKSxcbiAgICAgIGJ5dGVzOiBzZXJpYWxpemF0aW9uLmVuY29kZXIoXG4gICAgICAgIHRoaXMuYnl0ZXMsXG4gICAgICAgIGVuY29kaW5nLFxuICAgICAgICBcIkJ1ZmZlclwiLFxuICAgICAgICBcImhleFwiLFxuICAgICAgICB0aGlzLmJzaXplXG4gICAgICApXG4gICAgfVxuICB9XG4gIGRlc2VyaWFsaXplKGZpZWxkczogb2JqZWN0LCBlbmNvZGluZzogU2VyaWFsaXplZEVuY29kaW5nID0gXCJoZXhcIikge1xuICAgIHN1cGVyLmRlc2VyaWFsaXplKGZpZWxkcywgZW5jb2RpbmcpXG4gICAgdGhpcy5ic2l6ZSA9IHNlcmlhbGl6YXRpb24uZGVjb2RlcihcbiAgICAgIGZpZWxkc1tcImJzaXplXCJdLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBcImRlY2ltYWxTdHJpbmdcIixcbiAgICAgIFwibnVtYmVyXCIsXG4gICAgICA0XG4gICAgKVxuICAgIHRoaXMuYnl0ZXMgPSBzZXJpYWxpemF0aW9uLmRlY29kZXIoXG4gICAgICBmaWVsZHNbXCJieXRlc1wiXSxcbiAgICAgIGVuY29kaW5nLFxuICAgICAgXCJoZXhcIixcbiAgICAgIFwiQnVmZmVyXCIsXG4gICAgICB0aGlzLmJzaXplXG4gICAgKVxuICB9XG5cbiAgcHJvdGVjdGVkIGJ5dGVzOiBCdWZmZXJcbiAgcHJvdGVjdGVkIGJzaXplOiBudW1iZXJcblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgbGVuZ3RoIG9mIHRoZSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfS5cbiAgICpcbiAgICogQHJldHVybnMgVGhlIGV4YWN0IGxlbmd0aCByZXF1aXJlbWVudCBvZiB0aGlzIGNsYXNzXG4gICAqL1xuICBnZXRTaXplID0gKCkgPT4gdGhpcy5ic2l6ZVxuXG4gIC8qKlxuICAgKiBUYWtlcyBhIGJhc2UtNTggZW5jb2RlZCBzdHJpbmcsIHZlcmlmaWVzIGl0cyBsZW5ndGgsIGFuZCBzdG9yZXMgaXQuXG4gICAqXG4gICAqIEByZXR1cm5zIFRoZSBzaXplIG9mIHRoZSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfVxuICAgKi9cbiAgZnJvbVN0cmluZyhiNThzdHI6IHN0cmluZyk6IG51bWJlciB7XG4gICAgdHJ5IHtcbiAgICAgIHRoaXMuZnJvbUJ1ZmZlcihiaW50b29scy5iNThUb0J1ZmZlcihiNThzdHIpKVxuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgIC8qIGlzdGFuYnVsIGlnbm9yZSBuZXh0ICovXG4gICAgICBjb25zdCBlbXNnOiBzdHJpbmcgPSBgRXJyb3IgLSBOQnl0ZXMuZnJvbVN0cmluZzogJHtlfWBcbiAgICAgIC8qIGlzdGFuYnVsIGlnbm9yZSBuZXh0ICovXG4gICAgICB0aHJvdyBuZXcgRXJyb3IoZW1zZylcbiAgICB9XG4gICAgcmV0dXJuIHRoaXMuYnNpemVcbiAgfVxuXG4gIC8qKlxuICAgKiBUYWtlcyBhIFtbQnVmZmVyXV0sIHZlcmlmaWVzIGl0cyBsZW5ndGgsIGFuZCBzdG9yZXMgaXQuXG4gICAqXG4gICAqIEByZXR1cm5zIFRoZSBzaXplIG9mIHRoZSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfVxuICAgKi9cbiAgZnJvbUJ1ZmZlcihidWZmOiBCdWZmZXIsIG9mZnNldDogbnVtYmVyID0gMCk6IG51bWJlciB7XG4gICAgdHJ5IHtcbiAgICAgIGlmIChidWZmLmxlbmd0aCAtIG9mZnNldCA8IHRoaXMuYnNpemUpIHtcbiAgICAgICAgLyogaXN0YW5idWwgaWdub3JlIG5leHQgKi9cbiAgICAgICAgdGhyb3cgbmV3IEJ1ZmZlclNpemVFcnJvcihcbiAgICAgICAgICBcIkVycm9yIC0gTkJ5dGVzLmZyb21CdWZmZXI6IG5vdCBlbm91Z2ggc3BhY2UgYXZhaWxhYmxlIGluIGJ1ZmZlci5cIlxuICAgICAgICApXG4gICAgICB9XG5cbiAgICAgIHRoaXMuYnl0ZXMgPSBiaW50b29scy5jb3B5RnJvbShidWZmLCBvZmZzZXQsIG9mZnNldCArIHRoaXMuYnNpemUpXG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgLyogaXN0YW5idWwgaWdub3JlIG5leHQgKi9cbiAgICAgIGNvbnN0IGVtc2c6IHN0cmluZyA9IGBFcnJvciAtIE5CeXRlcy5mcm9tQnVmZmVyOiAke2V9YFxuICAgICAgLyogaXN0YW5idWwgaWdub3JlIG5leHQgKi9cbiAgICAgIHRocm93IG5ldyBFcnJvcihlbXNnKVxuICAgIH1cbiAgICByZXR1cm4gb2Zmc2V0ICsgdGhpcy5ic2l6ZVxuICB9XG5cbiAgLyoqXG4gICAqIEByZXR1cm5zIEEgcmVmZXJlbmNlIHRvIHRoZSBzdG9yZWQge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn1cbiAgICovXG4gIHRvQnVmZmVyKCk6IEJ1ZmZlciB7XG4gICAgcmV0dXJuIHRoaXMuYnl0ZXNcbiAgfVxuXG4gIC8qKlxuICAgKiBAcmV0dXJucyBBIGJhc2UtNTggc3RyaW5nIG9mIHRoZSBzdG9yZWQge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn1cbiAgICovXG4gIHRvU3RyaW5nKCk6IHN0cmluZyB7XG4gICAgcmV0dXJuIGJpbnRvb2xzLmJ1ZmZlclRvQjU4KHRoaXMudG9CdWZmZXIoKSlcbiAgfVxuXG4gIGFic3RyYWN0IGNsb25lKCk6IHRoaXNcbiAgYWJzdHJhY3QgY3JlYXRlKC4uLmFyZ3M6IGFueVtdKTogdGhpc1xufVxuIl19