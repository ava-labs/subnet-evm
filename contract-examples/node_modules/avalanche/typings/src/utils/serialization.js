"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Serialization = exports.Serializable = exports.SERIALIZATIONVERSION = void 0;
/**
 * @packageDocumentation
 * @module Utils-Serialization
 */
const bintools_1 = __importDefault(require("../utils/bintools"));
const bn_js_1 = __importDefault(require("bn.js"));
const buffer_1 = require("buffer/");
const xss_1 = __importDefault(require("xss"));
const helperfunctions_1 = require("./helperfunctions");
const errors_1 = require("../utils/errors");
exports.SERIALIZATIONVERSION = 0;
class Serializable {
    constructor() {
        this._typeName = undefined;
        this._typeID = undefined;
        this._codecID = undefined;
    }
    /**
     * Used in serialization. TypeName is a string name for the type of object being output.
     */
    getTypeName() {
        return this._typeName;
    }
    /**
     * Used in serialization. Optional. TypeID is a number for the typeID of object being output.
     */
    getTypeID() {
        return this._typeID;
    }
    /**
     * Used in serialization. Optional. TypeID is a number for the typeID of object being output.
     */
    getCodecID() {
        return this._codecID;
    }
    /**
     * Sanitize to prevent cross scripting attacks.
     */
    sanitizeObject(obj) {
        for (const k in obj) {
            if (typeof obj[`${k}`] === "object" && obj[`${k}`] !== null) {
                this.sanitizeObject(obj[`${k}`]);
            }
            else if (typeof obj[`${k}`] === "string") {
                obj[`${k}`] = (0, xss_1.default)(obj[`${k}`]);
            }
        }
        return obj;
    }
    //sometimes the parent class manages the fields
    //these are so you can say super.serialize(encoding)
    serialize(encoding) {
        return {
            _typeName: (0, xss_1.default)(this._typeName),
            _typeID: typeof this._typeID === "undefined" ? null : this._typeID,
            _codecID: typeof this._codecID === "undefined" ? null : this._codecID
        };
    }
    deserialize(fields, encoding) {
        fields = this.sanitizeObject(fields);
        if (typeof fields["_typeName"] !== "string") {
            throw new errors_1.TypeNameError("Error - Serializable.deserialize: _typeName must be a string, found: " +
                typeof fields["_typeName"]);
        }
        if (fields["_typeName"] !== this._typeName) {
            throw new errors_1.TypeNameError("Error - Serializable.deserialize: _typeName mismatch -- expected: " +
                this._typeName +
                " -- received: " +
                fields["_typeName"]);
        }
        if (typeof fields["_typeID"] !== "undefined" &&
            fields["_typeID"] !== null) {
            if (typeof fields["_typeID"] !== "number") {
                throw new errors_1.TypeIdError("Error - Serializable.deserialize: _typeID must be a number, found: " +
                    typeof fields["_typeID"]);
            }
            if (fields["_typeID"] !== this._typeID) {
                throw new errors_1.TypeIdError("Error - Serializable.deserialize: _typeID mismatch -- expected: " +
                    this._typeID +
                    " -- received: " +
                    fields["_typeID"]);
            }
        }
        if (typeof fields["_codecID"] !== "undefined" &&
            fields["_codecID"] !== null) {
            if (typeof fields["_codecID"] !== "number") {
                throw new errors_1.CodecIdError("Error - Serializable.deserialize: _codecID must be a number, found: " +
                    typeof fields["_codecID"]);
            }
            if (fields["_codecID"] !== this._codecID) {
                throw new errors_1.CodecIdError("Error - Serializable.deserialize: _codecID mismatch -- expected: " +
                    this._codecID +
                    " -- received: " +
                    fields["_codecID"]);
            }
        }
    }
}
exports.Serializable = Serializable;
class Serialization {
    constructor() {
        this.bintools = bintools_1.default.getInstance();
    }
    /**
     * Retrieves the Serialization singleton.
     */
    static getInstance() {
        if (!Serialization.instance) {
            Serialization.instance = new Serialization();
        }
        return Serialization.instance;
    }
    /**
     * Convert {@link https://github.com/feross/buffer|Buffer} to [[SerializedType]]
     *
     * @param vb {@link https://github.com/feross/buffer|Buffer}
     * @param type [[SerializedType]]
     * @param ...args remaining arguments
     * @returns type of [[SerializedType]]
     */
    bufferToType(vb, type, ...args) {
        if (type === "BN") {
            return new bn_js_1.default(vb.toString("hex"), "hex");
        }
        else if (type === "Buffer") {
            if (args.length == 1 && typeof args[0] === "number") {
                vb = buffer_1.Buffer.from(vb.toString("hex").padStart(args[0] * 2, "0"), "hex");
            }
            return vb;
        }
        else if (type === "bech32") {
            return this.bintools.addressToString(args[0], args[1], vb);
        }
        else if (type === "nodeID") {
            return (0, helperfunctions_1.bufferToNodeIDString)(vb);
        }
        else if (type === "privateKey") {
            return (0, helperfunctions_1.bufferToPrivateKeyString)(vb);
        }
        else if (type === "cb58") {
            return this.bintools.cb58Encode(vb);
        }
        else if (type === "base58") {
            return this.bintools.bufferToB58(vb);
        }
        else if (type === "base64") {
            return vb.toString("base64");
        }
        else if (type === "hex") {
            return vb.toString("hex");
        }
        else if (type === "decimalString") {
            return new bn_js_1.default(vb.toString("hex"), "hex").toString(10);
        }
        else if (type === "number") {
            return new bn_js_1.default(vb.toString("hex"), "hex").toNumber();
        }
        else if (type === "utf8") {
            return vb.toString("utf8");
        }
        return undefined;
    }
    /**
     * Convert [[SerializedType]] to {@link https://github.com/feross/buffer|Buffer}
     *
     * @param v type of [[SerializedType]]
     * @param type [[SerializedType]]
     * @param ...args remaining arguments
     * @returns {@link https://github.com/feross/buffer|Buffer}
     */
    typeToBuffer(v, type, ...args) {
        if (type === "BN") {
            let str = v.toString("hex");
            if (args.length == 1 && typeof args[0] === "number") {
                return buffer_1.Buffer.from(str.padStart(args[0] * 2, "0"), "hex");
            }
            return buffer_1.Buffer.from(str, "hex");
        }
        else if (type === "Buffer") {
            return v;
        }
        else if (type === "bech32") {
            return this.bintools.stringToAddress(v, ...args);
        }
        else if (type === "nodeID") {
            return (0, helperfunctions_1.NodeIDStringToBuffer)(v);
        }
        else if (type === "privateKey") {
            return (0, helperfunctions_1.privateKeyStringToBuffer)(v);
        }
        else if (type === "cb58") {
            return this.bintools.cb58Decode(v);
        }
        else if (type === "base58") {
            return this.bintools.b58ToBuffer(v);
        }
        else if (type === "base64") {
            return buffer_1.Buffer.from(v, "base64");
        }
        else if (type === "hex") {
            if (v.startsWith("0x")) {
                v = v.slice(2);
            }
            return buffer_1.Buffer.from(v, "hex");
        }
        else if (type === "decimalString") {
            let str = new bn_js_1.default(v, 10).toString("hex");
            if (args.length == 1 && typeof args[0] === "number") {
                return buffer_1.Buffer.from(str.padStart(args[0] * 2, "0"), "hex");
            }
            return buffer_1.Buffer.from(str, "hex");
        }
        else if (type === "number") {
            let str = new bn_js_1.default(v, 10).toString("hex");
            if (args.length == 1 && typeof args[0] === "number") {
                return buffer_1.Buffer.from(str.padStart(args[0] * 2, "0"), "hex");
            }
            return buffer_1.Buffer.from(str, "hex");
        }
        else if (type === "utf8") {
            if (args.length == 1 && typeof args[0] === "number") {
                let b = buffer_1.Buffer.alloc(args[0]);
                b.write(v);
                return b;
            }
            return buffer_1.Buffer.from(v, "utf8");
        }
        return undefined;
    }
    /**
     * Convert value to type of [[SerializedType]] or [[SerializedEncoding]]
     *
     * @param value
     * @param encoding [[SerializedEncoding]]
     * @param intype [[SerializedType]]
     * @param outtype [[SerializedType]]
     * @param ...args remaining arguments
     * @returns type of [[SerializedType]] or [[SerializedEncoding]]
     */
    encoder(value, encoding, intype, outtype, ...args) {
        if (typeof value === "undefined") {
            throw new errors_1.UnknownTypeError("Error - Serializable.encoder: value passed is undefined");
        }
        if (encoding !== "display") {
            outtype = encoding;
        }
        const vb = this.typeToBuffer(value, intype, ...args);
        return this.bufferToType(vb, outtype, ...args);
    }
    /**
     * Convert value to type of [[SerializedType]] or [[SerializedEncoding]]
     *
     * @param value
     * @param encoding [[SerializedEncoding]]
     * @param intype [[SerializedType]]
     * @param outtype [[SerializedType]]
     * @param ...args remaining arguments
     * @returns type of [[SerializedType]] or [[SerializedEncoding]]
     */
    decoder(value, encoding, intype, outtype, ...args) {
        if (typeof value === "undefined") {
            throw new errors_1.UnknownTypeError("Error - Serializable.decoder: value passed is undefined");
        }
        if (encoding !== "display") {
            intype = encoding;
        }
        const vb = this.typeToBuffer(value, intype, ...args);
        return this.bufferToType(vb, outtype, ...args);
    }
    serialize(serialize, vm, encoding = "display", notes = undefined) {
        if (typeof notes === "undefined") {
            notes = serialize.getTypeName();
        }
        return {
            vm,
            encoding,
            version: exports.SERIALIZATIONVERSION,
            notes,
            fields: serialize.serialize(encoding)
        };
    }
    deserialize(input, output) {
        output.deserialize(input.fields, input.encoding);
    }
}
exports.Serialization = Serialization;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoic2VyaWFsaXphdGlvbi5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3NyYy91dGlscy9zZXJpYWxpemF0aW9uLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7OztBQUFBOzs7R0FHRztBQUNILGlFQUF3QztBQUN4QyxrREFBc0I7QUFDdEIsb0NBQWdDO0FBQ2hDLDhDQUFxQjtBQUNyQix1REFLMEI7QUFDMUIsNENBS3dCO0FBR1gsUUFBQSxvQkFBb0IsR0FBVyxDQUFDLENBQUE7QUF5QjdDLE1BQXNCLFlBQVk7SUFBbEM7UUFDWSxjQUFTLEdBQVcsU0FBUyxDQUFBO1FBQzdCLFlBQU8sR0FBVyxTQUFTLENBQUE7UUFDM0IsYUFBUSxHQUFXLFNBQVMsQ0FBQTtJQXFHeEMsQ0FBQztJQW5HQzs7T0FFRztJQUNILFdBQVc7UUFDVCxPQUFPLElBQUksQ0FBQyxTQUFTLENBQUE7SUFDdkIsQ0FBQztJQUVEOztPQUVHO0lBQ0gsU0FBUztRQUNQLE9BQU8sSUFBSSxDQUFDLE9BQU8sQ0FBQTtJQUNyQixDQUFDO0lBRUQ7O09BRUc7SUFDSCxVQUFVO1FBQ1IsT0FBTyxJQUFJLENBQUMsUUFBUSxDQUFBO0lBQ3RCLENBQUM7SUFFRDs7T0FFRztJQUNILGNBQWMsQ0FBQyxHQUFXO1FBQ3hCLEtBQUssTUFBTSxDQUFDLElBQUksR0FBRyxFQUFFO1lBQ25CLElBQUksT0FBTyxHQUFHLENBQUMsR0FBRyxDQUFDLEVBQUUsQ0FBQyxLQUFLLFFBQVEsSUFBSSxHQUFHLENBQUMsR0FBRyxDQUFDLEVBQUUsQ0FBQyxLQUFLLElBQUksRUFBRTtnQkFDM0QsSUFBSSxDQUFDLGNBQWMsQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFDLEVBQUUsQ0FBQyxDQUFDLENBQUE7YUFDakM7aUJBQU0sSUFBSSxPQUFPLEdBQUcsQ0FBQyxHQUFHLENBQUMsRUFBRSxDQUFDLEtBQUssUUFBUSxFQUFFO2dCQUMxQyxHQUFHLENBQUMsR0FBRyxDQUFDLEVBQUUsQ0FBQyxH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Q2xCO1FBQ0QsTUFBTSxFQUFFLEdBQVcsSUFBSSxDQUFDLFlBQVksQ0FBQyxLQUFLLEVBQUUsTUFBTSxFQUFFLEdBQUcsSUFBSSxDQUFDLENBQUE7UUFDNUQsT0FBTyxJQUFJLENBQUMsWUFBWSxDQUFDLEVBQUUsRUFBRSxPQUFPLEVBQUUsR0FBRyxJQUFJLENBQUMsQ0FBQTtJQUNoRCxDQUFDO0lBRUQsU0FBUyxDQUNQLFNBQXVCLEVBQ3ZCLEVBQVUsRUFDVixXQUErQixTQUFTLEVBQ3hDLFFBQWdCLFNBQVM7UUFFekIsSUFBSSxPQUFPLEtBQUssS0FBSyxXQUFXLEVBQUU7WUFDaEMsS0FBSyxHQUFHLFNBQVMsQ0FBQyxXQUFXLEVBQUUsQ0FBQTtTQUNoQztRQUNELE9BQU87WUFDTCxFQUFFO1lBQ0YsUUFBUTtZQUNSLE9BQU8sRUFBRSw0QkFBb0I7WUFDN0IsS0FBSztZQUNMLE1BQU0sRUFBRSxTQUFTLENBQUMsU0FBUyxDQUFDLFFBQVEsQ0FBQztTQUN0QyxDQUFBO0lBQ0gsQ0FBQztJQUVELFdBQVcsQ0FBQyxLQUFpQixFQUFFLE1BQW9CO1FBQ2pELE1BQU0sQ0FBQyxXQUFXLENBQUMsS0FBSyxDQUFDLE1BQU0sRUFBRSxLQUFLLENBQUMsUUFBUSxDQUFDLENBQUE7SUFDbEQsQ0FBQztDQUNGO0FBbE1ELHNDQWtNQyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogQHBhY2thZ2VEb2N1bWVudGF0aW9uXG4gKiBAbW9kdWxlIFV0aWxzLVNlcmlhbGl6YXRpb25cbiAqL1xuaW1wb3J0IEJpblRvb2xzIGZyb20gXCIuLi91dGlscy9iaW50b29sc1wiXG5pbXBvcnQgQk4gZnJvbSBcImJuLmpzXCJcbmltcG9ydCB7IEJ1ZmZlciB9IGZyb20gXCJidWZmZXIvXCJcbmltcG9ydCB4c3MgZnJvbSBcInhzc1wiXG5pbXBvcnQge1xuICBOb2RlSURTdHJpbmdUb0J1ZmZlcixcbiAgcHJpdmF0ZUtleVN0cmluZ1RvQnVmZmVyLFxuICBidWZmZXJUb05vZGVJRFN0cmluZyxcbiAgYnVmZmVyVG9Qcml2YXRlS2V5U3RyaW5nXG59IGZyb20gXCIuL2hlbHBlcmZ1bmN0aW9uc1wiXG5pbXBvcnQge1xuICBDb2RlY0lkRXJyb3IsXG4gIFR5cGVJZEVycm9yLFxuICBUeXBlTmFtZUVycm9yLFxuICBVbmtub3duVHlwZUVycm9yXG59IGZyb20gXCIuLi91dGlscy9lcnJvcnNcIlxuaW1wb3J0IHsgU2VyaWFsaXplZCB9IGZyb20gXCIuLi9jb21tb25cIlxuXG5leHBvcnQgY29uc3QgU0VSSUFMSVpBVElPTlZFUlNJT046IG51bWJlciA9IDBcbmV4cG9ydCB0eXBlIFNlcmlhbGl6ZWRUeXBlID1cbiAgfCBcImhleFwiXG4gIHwgXCJCTlwiXG4gIHwgXCJCdWZmZXJcIlxuICB8IFwiYmVjaDMyXCJcbiAgfCBcIm5vZGVJRFwiXG4gIHwgXCJwcml2YXRlS2V5XCJcbiAgfCBcImNiNThcIlxuICB8IFwiYmFzZTU4XCJcbiAgfCBcImJhc2U2NFwiXG4gIHwgXCJkZWNpbWFsU3RyaW5nXCJcbiAgfCBcIm51bWJlclwiXG4gIHwgXCJ1dGY4XCJcblxuZXhwb3J0IHR5cGUgU2VyaWFsaXplZEVuY29kaW5nID1cbiAgfCBcImhleFwiXG4gIHwgXCJjYjU4XCJcbiAgfCBcImJhc2U1OFwiXG4gIHwgXCJiYXNlNjRcIlxuICB8IFwiZGVjaW1hbFN0cmluZ1wiXG4gIHwgXCJudW1iZXJcIlxuICB8IFwidXRmOFwiXG4gIHwgXCJkaXNwbGF5XCJcblxuZXhwb3J0IGFic3RyYWN0IGNsYXNzIFNlcmlhbGl6YWJsZSB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWU6IHN0cmluZyA9IHVuZGVmaW5lZFxuICBwcm90ZWN0ZWQgX3R5cGVJRDogbnVtYmVyID0gdW5kZWZpbmVkXG4gIHByb3RlY3RlZCBfY29kZWNJRDogbnVtYmVyID0gdW5kZWZpbmVkXG5cbiAgLyoqXG4gICAqIFVzZWQgaW4gc2VyaWFsaXphdGlvbi4gVHlwZU5hbWUgaXMgYSBzdHJpbmcgbmFtZSBmb3IgdGhlIHR5cGUgb2Ygb2JqZWN0IGJlaW5nIG91dHB1dC5cbiAgICovXG4gIGdldFR5cGVOYW1lKCk6IHN0cmluZyB7XG4gICAgcmV0dXJuIHRoaXMuX3R5cGVOYW1lXG4gIH1cblxuICAvKipcbiAgICogVXNlZCBpbiBzZXJpYWxpemF0aW9uLiBPcHRpb25hbC4gVHlwZUlEIGlzIGEgbnVtYmVyIGZvciB0aGUgdHlwZUlEIG9mIG9iamVjdCBiZWluZyBvdXRwdXQuXG4gICAqL1xuICBnZXRUeXBlSUQoKTogbnVtYmVyIHtcbiAgICByZXR1cm4gdGhpcy5fdHlwZUlEXG4gIH1cblxuICAvKipcbiAgICogVXNlZCBpbiBzZXJpYWxpemF0aW9uLiBPcHRpb25hbC4gVHlwZUlEIGlzIGEgbnVtYmVyIGZvciB0aGUgdHlwZUlEIG9mIG9iamVjdCBiZWluZyBvdXRwdXQuXG4gICAqL1xuICBnZXRDb2RlY0lEKCk6IG51bWJlciB7XG4gICAgcmV0dXJuIHRoaXMuX2NvZGVjSURcbiAgfVxuXG4gIC8qKlxuICAgKiBTYW5pdGl6ZSB0byBwcmV2ZW50IGNyb3NzIHNjcmlwdGluZyBhdHRhY2tzLlxuICAgKi9cbiAgc2FuaXRpemVPYmplY3Qob2JqOiBvYmplY3QpOiBvYmplY3Qge1xuICAgIGZvciAoY29uc3QgayBpbiBvYmopIHtcbiAgICAgIGlmICh0eXBlb2Ygb2JqW2Ake2t9YF0gPT09IFwib2JqZWN0XCIgJiYgb2JqW2Ake2t9YF0gIT09IG51bGwpIHtcbiAgICAgICAgdGhpcy5zYW5pdGl6ZU9iamVjdChvYmpbYCR7a31gXSlcbiAgICAgIH0gZWxzZSBpZiAodHlwZW9mIG9ialtgJHtrfWBdID09PSBcInN0cmluZ1wiKSB7XG4gICAgICAgIG9ialtgJHtrfWBdID0geHNzKG9ialtgJHtrfWBdKVxuICAgICAgfVxuICAgIH1cbiAgICByZXR1cm4gb2JqXG4gIH1cblxuICAvL3NvbWV0aW1lcyB0aGUgcGFyZW50IGNsYXNzIG1hbmFnZXMgdGhlIGZpZWxkc1xuICAvL3RoZXNlIGFyZSBzbyB5b3UgY2FuIHNheSBzdXBlci5zZXJpYWxpemUoZW5jb2RpbmcpXG4gIHNlcmlhbGl6ZShlbmNvZGluZz86IFNlcmlhbGl6ZWRFbmNvZGluZyk6IG9iamVjdCB7XG4gICAgcmV0dXJuIHtcbiAgICAgIF90eXBlTmFtZTogeHNzKHRoaXMuX3R5cGVOYW1lKSxcbiAgICAgIF90eXBlSUQ6IHR5cGVvZiB0aGlzLl90eXBlSUQgPT09IFwidW5kZWZpbmVkXCIgPyBudWxsIDogdGhpcy5fdHlwZUlELFxuICAgICAgX2NvZGVjSUQ6IHR5cGVvZiB0aGlzLl9jb2RlY0lEID09PSBcInVuZGVmaW5lZFwiID8gbnVsbCA6IHRoaXMuX2NvZGVjSURcbiAgICB9XG4gIH1cbiAgZGVzZXJpYWxpemUoZmllbGRzOiBvYmplY3QsIGVuY29kaW5nPzogU2VyaWFsaXplZEVuY29kaW5nKTogdm9pZCB7XG4gICAgZmllbGRzID0gdGhpcy5zYW5pdGl6ZU9iamVjdChmaWVsZHMpXG4gICAgaWYgKHR5cGVvZiBmaWVsZHNbXCJfdHlwZU5hbWVcIl0gIT09IFwic3RyaW5nXCIpIHtcbiAgICAgIHRocm93IG5ldyBUeXBlTmFtZUVycm9yKFxuICAgICAgICBcIkVycm9yIC0gU2VyaWFsaXphYmxlLmRlc2VyaWFsaXplOiBfdHlwZU5hbWUgbXVzdCBiZSBhIHN0cmluZywgZm91bmQ6IFwiICtcbiAgICAgICAgICB0eXBlb2YgZmllbGRzW1wiX3R5cGVOYW1lXCJdXG4gICAgICApXG4gICAgfVxuICAgIGlmIChmaWVsZHNbXCJfdHlwZU5hbWVcIl0gIT09IHRoaXMuX3R5cGVOYW1lKSB7XG4gICAgICB0aHJvdyBuZXcgVHlwZU5hbWVFcnJvcihcbiAgICAgICAgXCJFcnJvciAtIFNlcmlhbGl6YWJsZS5kZXNlcmlhbGl6ZTogX3R5cGVOYW1lIG1pc21hdGNoIC0tIGV4cGVjdGVkOiBcIiArXG4gICAgICAgICAgdGhpcy5fdHlwZU5hbWUgK1xuICAgICAgICAgIFwiIC0tIHJlY2VpdmVkOiBcIiArXG4gICAgICAgICAgZmllbGRzW1wiX3R5cGVOYW1lXCJdXG4gICAgICApXG4gICAgfVxuICAgIGlmIChcbiAgICAgIHR5cGVvZiBmaWVsZHNbXCJfdHlwZUlEXCJdICE9PSBcInVuZGVmaW5lZFwiICYmXG4gICAgICBmaWVsZHNbXCJfdHlwZUlEXCJdICE9PSBudWxsXG4gICAgKSB7XG4gICAgICBpZiAodHlwZW9mIGZpZWxkc1tcIl90eXBlSURcIl0gIT09IFwibnVtYmVyXCIpIHtcbiAgICAgICAgdGhyb3cgbmV3IFR5cGVJZEVycm9yKFxuICAgICAgICAgIFwiRXJyb3IgLSBTZXJpYWxpemFibGUuZGVzZXJpYWxpemU6IF90eXBlSUQgbXVzdCBiZSBhIG51bWJlciwgZm91bmQ6IFwiICtcbiAgICAgICAgICAgIHR5cGVvZiBmaWVsZHNbXCJfdHlwZUlEXCJdXG4gICAgICAgIClcbiAgICAgIH1cbiAgICAgIGlmIChmaWVsZHNbXCJfdHlwZUlEXCJdICE9PSB0aGlzLl90eXBlSUQpIHtcbiAgICAgICAgdGhyb3cgbmV3IFR5cGVJZEVycm9yKFxuICAgICAgICAgIFwiRXJyb3IgLSBTZXJpYWxpemFibGUuZGVzZXJpYWxpemU6IF90eXBlSUQgbWlzbWF0Y2ggLS0gZXhwZWN0ZWQ6IFwiICtcbiAgICAgICAgICAgIHRoaXMuX3R5cGVJRCArXG4gICAgICAgICAgICBcIiAtLSByZWNlaXZlZDogXCIgK1xuICAgICAgICAgICAgZmllbGRzW1wiX3R5cGVJRFwiXVxuICAgICAgICApXG4gICAgICB9XG4gICAgfVxuICAgIGlmIChcbiAgICAgIHR5cGVvZiBmaWVsZHNbXCJfY29kZWNJRFwiXSAhPT0gXCJ1bmRlZmluZWRcIiAmJlxuICAgICAgZmllbGRzW1wiX2NvZGVjSURcIl0gIT09IG51bGxcbiAgICApIHtcbiAgICAgIGlmICh0eXBlb2YgZmllbGRzW1wiX2NvZGVjSURcIl0gIT09IFwibnVtYmVyXCIpIHtcbiAgICAgICAgdGhyb3cgbmV3IENvZGVjSWRFcnJvcihcbiAgICAgICAgICBcIkVycm9yIC0gU2VyaWFsaXphYmxlLmRlc2VyaWFsaXplOiBfY29kZWNJRCBtdXN0IGJlIGEgbnVtYmVyLCBmb3VuZDogXCIgK1xuICAgICAgICAgICAgdHlwZW9mIGZpZWxkc1tcIl9jb2RlY0lEXCJdXG4gICAgICAgIClcbiAgICAgIH1cbiAgICAgIGlmIChmaWVsZHNbXCJfY29kZWNJRFwiXSAhPT0gdGhpcy5fY29kZWNJRCkge1xuICAgICAgICB0aHJvdyBuZXcgQ29kZWNJZEVycm9yKFxuICAgICAgICAgIFwiRXJyb3IgLSBTZXJpYWxpemFibGUuZGVzZXJpYWxpemU6IF9jb2RlY0lEIG1pc21hdGNoIC0tIGV4cGVjdGVkOiBcIiArXG4gICAgICAgICAgICB0aGlzLl9jb2RlY0lEICtcbiAgICAgICAgICAgIFwiIC0tIHJlY2VpdmVkOiBcIiArXG4gICAgICAgICAgICBmaWVsZHNbXCJfY29kZWNJRFwiXVxuICAgICAgICApXG4gICAgICB9XG4gICAgfVxuICB9XG59XG5cbmV4cG9ydCBjbGFzcyBTZXJpYWxpemF0aW9uIHtcbiAgcHJpdmF0ZSBzdGF0aWMgaW5zdGFuY2U6IFNlcmlhbGl6YXRpb25cblxuICBwcml2YXRlIGNvbnN0cnVjdG9yKCkge1xuICAgIHRoaXMuYmludG9vbHMgPSBCaW5Ub29scy5nZXRJbnN0YW5jZSgpXG4gIH1cbiAgcHJpdmF0ZSBiaW50b29sczogQmluVG9vbHNcblxuICAvKipcbiAgICogUmV0cmlldmVzIHRoZSBTZXJpYWxpemF0aW9uIHNpbmdsZXRvbi5cbiAgICovXG4gIHN0YXRpYyBnZXRJbnN0YW5jZSgpOiBTZXJpYWxpemF0aW9uIHtcbiAgICBpZiAoIVNlcmlhbGl6YXRpb24uaW5zdGFuY2UpIHtcbiAgICAgIFNlcmlhbGl6YXRpb24uaW5zdGFuY2UgPSBuZXcgU2VyaWFsaXphdGlvbigpXG4gICAgfVxuICAgIHJldHVybiBTZXJpYWxpemF0aW9uLmluc3RhbmNlXG4gIH1cblxuICAvKipcbiAgICogQ29udmVydCB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSB0byBbW1NlcmlhbGl6ZWRUeXBlXV1cbiAgICpcbiAgICogQHBhcmFtIHZiIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9XG4gICAqIEBwYXJhbSB0eXBlIFtbU2VyaWFsaXplZFR5cGVdXVxuICAgKiBAcGFyYW0gLi4uYXJncyByZW1haW5pbmcgYXJndW1lbnRzXG4gICAqIEByZXR1cm5zIHR5cGUgb2YgW1tTZXJpYWxpemVkVHlwZV1dXG4gICAqL1xuICBidWZmZXJUb1R5cGUodmI6IEJ1ZmZlciwgdHlwZTogU2VyaWFsaXplZFR5cGUsIC4uLmFyZ3M6IGFueVtdKTogYW55IHtcbiAgICBpZiAodHlwZSA9PT0gXCJCTlwiKSB7XG4gICAgICByZXR1cm4gbmV3IEJOKHZiLnRvU3RyaW5nKFwiaGV4XCIpLCBcImhleFwiKVxuICAgIH0gZWxzZSBpZiAodHlwZSA9PT0gXCJCdWZmZXJcIikge1xuICAgICAgaWYgKGFyZ3MubGVuZ3RoID09IDEgJiYgdHlwZW9mIGFyZ3NbMF0gPT09IFwibnVtYmVyXCIpIHtcbiAgICAgICAgdmIgPSBCdWZmZXIuZnJvbSh2Yi50b1N0cmluZyhcImhleFwiKS5wYWRTdGFydChhcmdzWzBdICogMiwgXCIwXCIpLCBcImhleFwiKVxuICAgICAgfVxuICAgICAgcmV0dXJuIHZiXG4gICAgfSBlbHNlIGlmICh0eXBlID09PSBcImJlY2gzMlwiKSB7XG4gICAgICByZXR1cm4gdGhpcy5iaW50b29scy5hZGRyZXNzVG9TdHJpbmcoYXJnc1swXSwgYXJnc1sxXSwgdmIpXG4gICAgfSBlbHNlIGlmICh0eXBlID09PSBcIm5vZGVJRFwiKSB7XG4gICAgICByZXR1cm4gYnVmZmVyVG9Ob2RlSURTdHJpbmcodmIpXG4gICAgfSBlbHNlIGlmICh0eXBlID09PSBcInByaXZhdGVLZXlcIikge1xuICAgICAgcmV0dXJuIGJ1ZmZlclRvUHJpdmF0ZUtleVN0cmluZyh2YilcbiAgICB9IGVsc2UgaWYgKHR5cGUgPT09IFwiY2I1OFwiKSB7XG4gICAgICByZXR1cm4gdGhpcy5iaW50b29scy5jYjU4RW5jb2RlKHZiKVxuICAgIH0gZWxzZSBpZiAodHlwZSA9PT0gXCJiYXNlNThcIikge1xuICAgICAgcmV0dXJuIHRoaXMuYmludG9vbHMuYnVmZmVyVG9CNTgodmIpXG4gICAgfSBlbHNlIGlmICh0eXBlID09PSBcImJhc2U2NFwiKSB7XG4gICAgICByZXR1cm4gdmIudG9TdHJpbmcoXCJiYXNlNjRcIilcbiAgICB9IGVsc2UgaWYgKHR5cGUgPT09IFwiaGV4XCIpIHtcbiAgICAgIHJldHVybiB2Yi50b1N0cmluZyhcImhleFwiKVxuICAgIH0gZWxzZSBpZiAodHlwZSA9PT0gXCJkZWNpbWFsU3RyaW5nXCIpIHtcbiAgICAgIHJldHVybiBuZXcgQk4odmIudG9TdHJpbmcoXCJoZXhcIiksIFwiaGV4XCIpLnRvU3RyaW5nKDEwKVxuICAgIH0gZWxzZSBpZiAodHlwZSA9PT0gXCJudW1iZXJcIikge1xuICAgICAgcmV0dXJuIG5ldyBCTih2Yi50b1N0cmluZyhcImhleFwiKSwgXCJoZXhcIikudG9OdW1iZXIoKVxuICAgIH0gZWxzZSBpZiAodHlwZSA9PT0gXCJ1dGY4XCIpIHtcbiAgICAgIHJldHVybiB2Yi50b1N0cmluZyhcInV0ZjhcIilcbiAgICB9XG4gICAgcmV0dXJuIHVuZGVmaW5lZFxuICB9XG5cbiAgLyoqXG4gICAqIENvbnZlcnQgW1tTZXJpYWxpemVkVHlwZV1dIHRvIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9XG4gICAqXG4gICAqIEBwYXJhbSB2IHR5cGUgb2YgW1tTZXJpYWxpemVkVHlwZV1dXG4gICAqIEBwYXJhbSB0eXBlIFtbU2VyaWFsaXplZFR5cGVdXVxuICAgKiBAcGFyYW0gLi4uYXJncyByZW1haW5pbmcgYXJndW1lbnRzXG4gICAqIEByZXR1cm5zIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9XG4gICAqL1xuICB0eXBlVG9CdWZmZXIodjogYW55LCB0eXBlOiBTZXJpYWxpemVkVHlwZSwgLi4uYXJnczogYW55W10pOiBCdWZmZXIge1xuICAgIGlmICh0eXBlID09PSBcIkJOXCIpIHtcbiAgICAgIGxldCBzdHI6IHN0cmluZyA9ICh2IGFzIEJOKS50b1N0cmluZyhcImhleFwiKVxuICAgICAgaWYgKGFyZ3MubGVuZ3RoID09IDEgJiYgdHlwZW9mIGFyZ3NbMF0gPT09IFwibnVtYmVyXCIpIHtcbiAgICAgICAgcmV0dXJuIEJ1ZmZlci5mcm9tKHN0ci5wYWRTdGFydChhcmdzWzBdICogMiwgXCIwXCIpLCBcImhleFwiKVxuICAgICAgfVxuICAgICAgcmV0dXJuIEJ1ZmZlci5mcm9tKHN0ciwgXCJoZXhcIilcbiAgICB9IGVsc2UgaWYgKHR5cGUgPT09IFwiQnVmZmVyXCIpIHtcbiAgICAgIHJldHVybiB2XG4gICAgfSBlbHNlIGlmICh0eXBlID09PSBcImJlY2gzMlwiKSB7XG4gICAgICByZXR1cm4gdGhpcy5iaW50b29scy5zdHJpbmdUb0FkZHJlc3ModiwgLi4uYXJncylcbiAgICB9IGVsc2UgaWYgKHR5cGUgPT09IFwibm9kZUlEXCIpIHtcbiAgICAgIHJldHVybiBOb2RlSURTdHJpbmdUb0J1ZmZlcih2KVxuICAgIH0gZWxzZSBpZiAodHlwZSA9PT0gXCJwcml2YXRlS2V5XCIpIHtcbiAgICAgIHJldHVybiBwcml2YXRlS2V5U3RyaW5nVG9CdWZmZXIodilcbiAgICB9IGVsc2UgaWYgKHR5cGUgPT09IFwiY2I1OFwiKSB7XG4gICAgICByZXR1cm4gdGhpcy5iaW50b29scy5jYjU4RGVjb2RlKHYpXG4gICAgfSBlbHNlIGlmICh0eXBlID09PSBcImJhc2U1OFwiKSB7XG4gICAgICByZXR1cm4gdGhpcy5iaW50b29scy5iNThUb0J1ZmZlcih2KVxuICAgIH0gZWxzZSBpZiAodHlwZSA9PT0gXCJiYXNlNjRcIikge1xuICAgICAgcmV0dXJuIEJ1ZmZlci5mcm9tKHYgYXMgc3RyaW5nLCBcImJhc2U2NFwiKVxuICAgIH0gZWxzZSBpZiAodHlwZSA9PT0gXCJoZXhcIikge1xuICAgICAgaWYgKCh2IGFzIHN0cmluZykuc3RhcnRzV2l0aChcIjB4XCIpKSB7XG4gICAgICAgIHYgPSAodiBhcyBzdHJpbmcpLnNsaWNlKDIpXG4gICAgICB9XG4gICAgICByZXR1cm4gQnVmZmVyLmZyb20odiBhcyBzdHJpbmcsIFwiaGV4XCIpXG4gICAgfSBlbHNlIGlmICh0eXBlID09PSBcImRlY2ltYWxTdHJpbmdcIikge1xuICAgICAgbGV0IHN0cjogc3RyaW5nID0gbmV3IEJOKHYgYXMgc3RyaW5nLCAxMCkudG9TdHJpbmcoXCJoZXhcIilcbiAgICAgIGlmIChhcmdzLmxlbmd0aCA9PSAxICYmIHR5cGVvZiBhcmdzWzBdID09PSBcIm51bWJlclwiKSB7XG4gICAgICAgIHJldHVybiBCdWZmZXIuZnJvbShzdHIucGFkU3RhcnQoYXJnc1swXSAqIDIsIFwiMFwiKSwgXCJoZXhcIilcbiAgICAgIH1cbiAgICAgIHJldHVybiBCdWZmZXIuZnJvbShzdHIsIFwiaGV4XCIpXG4gICAgfSBlbHNlIGlmICh0eXBlID09PSBcIm51bWJlclwiKSB7XG4gICAgICBsZXQgc3RyOiBzdHJpbmcgPSBuZXcgQk4odiwgMTApLnRvU3RyaW5nKFwiaGV4XCIpXG4gICAgICBpZiAoYXJncy5sZW5ndGggPT0gMSAmJiB0eXBlb2YgYXJnc1swXSA9PT0gXCJudW1iZXJcIikge1xuICAgICAgICByZXR1cm4gQnVmZmVyLmZyb20oc3RyLnBhZFN0YXJ0KGFyZ3NbMF0gKiAyLCBcIjBcIiksIFwiaGV4XCIpXG4gICAgICB9XG4gICAgICByZXR1cm4gQnVmZmVyLmZyb20oc3RyLCBcImhleFwiKVxuICAgIH0gZWxzZSBpZiAodHlwZSA9PT0gXCJ1dGY4XCIpIHtcbiAgICAgIGlmIChhcmdzLmxlbmd0aCA9PSAxICYmIHR5cGVvZiBhcmdzWzBdID09PSBcIm51bWJlclwiKSB7XG4gICAgICAgIGxldCBiOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoYXJnc1swXSlcbiAgICAgICAgYi53cml0ZSh2KVxuICAgICAgICByZXR1cm4gYlxuICAgICAgfVxuICAgICAgcmV0dXJuIEJ1ZmZlci5mcm9tKHYsIFwidXRmOFwiKVxuICAgIH1cbiAgICByZXR1cm4gdW5kZWZpbmVkXG4gIH1cblxuICAvKipcbiAgICogQ29udmVydCB2YWx1ZSB0byB0eXBlIG9mIFtbU2VyaWFsaXplZFR5cGVdXSBvciBbW1NlcmlhbGl6ZWRFbmNvZGluZ11dXG4gICAqXG4gICAqIEBwYXJhbSB2YWx1ZVxuICAgKiBAcGFyYW0gZW5jb2RpbmcgW1tTZXJpYWxpemVkRW5jb2RpbmddXVxuICAgKiBAcGFyYW0gaW50eXBlIFtbU2VyaWFsaXplZFR5cGVdXVxuICAgKiBAcGFyYW0gb3V0dHlwZSBbW1NlcmlhbGl6ZWRUeXBlXV1cbiAgICogQHBhcmFtIC4uLmFyZ3MgcmVtYWluaW5nIGFyZ3VtZW50c1xuICAgKiBAcmV0dXJucyB0eXBlIG9mIFtbU2VyaWFsaXplZFR5cGVdXSBvciBbW1NlcmlhbGl6ZWRFbmNvZGluZ11dXG4gICAqL1xuICBlbmNvZGVyKFxuICAgIHZhbHVlOiBhbnksXG4gICAgZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyxcbiAgICBpbnR5cGU6IFNlcmlhbGl6ZWRUeXBlLFxuICAgIG91dHR5cGU6IFNlcmlhbGl6ZWRUeXBlLFxuICAgIC4uLmFyZ3M6IGFueVtdXG4gICk6IGFueSB7XG4gICAgaWYgKHR5cGVvZiB2YWx1ZSA9PT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgdGhyb3cgbmV3IFVua25vd25UeXBlRXJyb3IoXG4gICAgICAgIFwiRXJyb3IgLSBTZXJpYWxpemFibGUuZW5jb2RlcjogdmFsdWUgcGFzc2VkIGlzIHVuZGVmaW5lZFwiXG4gICAgICApXG4gICAgfVxuICAgIGlmIChlbmNvZGluZyAhPT0gXCJkaXNwbGF5XCIpIHtcbiAgICAgIG91dHR5cGUgPSBlbmNvZGluZ1xuICAgIH1cbiAgICBjb25zdCB2YjogQnVmZmVyID0gdGhpcy50eXBlVG9CdWZmZXIodmFsdWUsIGludHlwZSwgLi4uYXJncylcbiAgICByZXR1cm4gdGhpcy5idWZmZXJUb1R5cGUodmIsIG91dHR5cGUsIC4uLmFyZ3MpXG4gIH1cblxuICAvKipcbiAgICogQ29udmVydCB2YWx1ZSB0byB0eXBlIG9mIFtbU2VyaWFsaXplZFR5cGVdXSBvciBbW1NlcmlhbGl6ZWRFbmNvZGluZ11dXG4gICAqXG4gICAqIEBwYXJhbSB2YWx1ZVxuICAgKiBAcGFyYW0gZW5jb2RpbmcgW1tTZXJpYWxpemVkRW5jb2RpbmddXVxuICAgKiBAcGFyYW0gaW50eXBlIFtbU2VyaWFsaXplZFR5cGVdXVxuICAgKiBAcGFyYW0gb3V0dHlwZSBbW1NlcmlhbGl6ZWRUeXBlXV1cbiAgICogQHBhcmFtIC4uLmFyZ3MgcmVtYWluaW5nIGFyZ3VtZW50c1xuICAgKiBAcmV0dXJucyB0eXBlIG9mIFtbU2VyaWFsaXplZFR5cGVdXSBvciBbW1NlcmlhbGl6ZWRFbmNvZGluZ11dXG4gICAqL1xuICBkZWNvZGVyKFxuICAgIHZhbHVlOiBzdHJpbmcsXG4gICAgZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyxcbiAgICBpbnR5cGU6IFNlcmlhbGl6ZWRUeXBlLFxuICAgIG91dHR5cGU6IFNlcmlhbGl6ZWRUeXBlLFxuICAgIC4uLmFyZ3M6IGFueVtdXG4gICk6IGFueSB7XG4gICAgaWYgKHR5cGVvZiB2YWx1ZSA9PT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgdGhyb3cgbmV3IFVua25vd25UeXBlRXJyb3IoXG4gICAgICAgIFwiRXJyb3IgLSBTZXJpYWxpemFibGUuZGVjb2RlcjogdmFsdWUgcGFzc2VkIGlzIHVuZGVmaW5lZFwiXG4gICAgICApXG4gICAgfVxuICAgIGlmIChlbmNvZGluZyAhPT0gXCJkaXNwbGF5XCIpIHtcbiAgICAgIGludHlwZSA9IGVuY29kaW5nXG4gICAgfVxuICAgIGNvbnN0IHZiOiBCdWZmZXIgPSB0aGlzLnR5cGVUb0J1ZmZlcih2YWx1ZSwgaW50eXBlLCAuLi5hcmdzKVxuICAgIHJldHVybiB0aGlzLmJ1ZmZlclRvVHlwZSh2Yiwgb3V0dHlwZSwgLi4uYXJncylcbiAgfVxuXG4gIHNlcmlhbGl6ZShcbiAgICBzZXJpYWxpemU6IFNlcmlhbGl6YWJsZSxcbiAgICB2bTogc3RyaW5nLFxuICAgIGVuY29kaW5nOiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImRpc3BsYXlcIixcbiAgICBub3Rlczogc3RyaW5nID0gdW5kZWZpbmVkXG4gICk6IFNlcmlhbGl6ZWQge1xuICAgIGlmICh0eXBlb2Ygbm90ZXMgPT09IFwidW5kZWZpbmVkXCIpIHtcbiAgICAgIG5vdGVzID0gc2VyaWFsaXplLmdldFR5cGVOYW1lKClcbiAgICB9XG4gICAgcmV0dXJuIHtcbiAgICAgIHZtLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICB2ZXJzaW9uOiBTRVJJQUxJWkFUSU9OVkVSU0lPTixcbiAgICAgIG5vdGVzLFxuICAgICAgZmllbGRzOiBzZXJpYWxpemUuc2VyaWFsaXplKGVuY29kaW5nKVxuICAgIH1cbiAgfVxuXG4gIGRlc2VyaWFsaXplKGlucHV0OiBTZXJpYWxpemVkLCBvdXRwdXQ6IFNlcmlhbGl6YWJsZSkge1xuICAgIG91dHB1dC5kZXNlcmlhbGl6ZShpbnB1dC5maWVsZHMsIGlucHV0LmVuY29kaW5nKVxuICB9XG59XG4iXX0=