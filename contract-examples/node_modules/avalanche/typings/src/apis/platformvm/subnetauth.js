"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.SubnetAuth = void 0;
/**
 * @packageDocumentation
 * @module API-PlatformVM-SubnetAuth
 */
const buffer_1 = require("buffer/");
const bintools_1 = __importDefault(require("../../utils/bintools"));
const utils_1 = require("../../utils");
const _1 = require(".");
/**
 * @ignore
 */
const bintools = bintools_1.default.getInstance();
class SubnetAuth extends utils_1.Serializable {
    constructor() {
        super(...arguments);
        this._typeName = "SubnetAuth";
        this._typeID = _1.PlatformVMConstants.SUBNETAUTH;
        this.addressIndices = [];
        this.numAddressIndices = buffer_1.Buffer.alloc(4);
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign({}, fields);
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
    }
    /**
     * Add an address index for Subnet Auth signing
     *
     * @param index the Buffer of the address index to add
     */
    addAddressIndex(index) {
        const numAddrIndices = this.getNumAddressIndices();
        this.numAddressIndices.writeUIntBE(numAddrIndices + 1, 0, 4);
        this.addressIndices.push(index);
    }
    /**
     * Returns the number of address indices as a number
     */
    getNumAddressIndices() {
        return this.numAddressIndices.readUIntBE(0, 4);
    }
    /**
     * Returns an array of AddressIndices as Buffers
     */
    getAddressIndices() {
        return this.addressIndices;
    }
    fromBuffer(bytes, offset = 0) {
        // increase offset for type id
        offset += 4;
        this.numAddressIndices = bintools.copyFrom(bytes, offset, offset + 4);
        offset += 4;
        for (let i = 0; i < this.getNumAddressIndices(); i++) {
            this.addressIndices.push(bintools.copyFrom(bytes, offset, offset + 4));
            offset += 4;
        }
        return offset;
    }
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[SubnetAuth]].
     */
    toBuffer() {
        const typeIDBuf = buffer_1.Buffer.alloc(4);
        typeIDBuf.writeUIntBE(this._typeID, 0, 4);
        const numAddressIndices = buffer_1.Buffer.alloc(4);
        numAddressIndices.writeIntBE(this.addressIndices.length, 0, 4);
        const barr = [typeIDBuf, numAddressIndices];
        let bsize = typeIDBuf.length + numAddressIndices.length;
        this.addressIndices.forEach((addressIndex, i) => {
            bsize += 4;
            barr.push(this.addressIndices[`${i}`]);
        });
        return buffer_1.Buffer.concat(barr, bsize);
    }
}
exports.SubnetAuth = SubnetAuth;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoic3VibmV0YXV0aC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3NyYy9hcGlzL3BsYXRmb3Jtdm0vc3VibmV0YXV0aC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7QUFBQTs7O0dBR0c7QUFDSCxvQ0FBZ0M7QUFDaEMsb0VBQTJDO0FBQzNDLHVDQUE4RDtBQUM5RCx3QkFBdUM7QUFFdkM7O0dBRUc7QUFDSCxNQUFNLFFBQVEsR0FBYSxrQkFBUSxDQUFDLFdBQVcsRUFBRSxDQUFBO0FBRWpELE1BQWEsVUFBVyxTQUFRLG9CQUFZO0lBQTVDOztRQUNZLGNBQVMsR0FBRyxZQUFZLENBQUE7UUFDeEIsWUFBTyxHQUFHLHNCQUFtQixDQUFDLFVBQVUsQ0FBQTtRQXFDeEMsbUJBQWMsR0FBYSxFQUFFLENBQUE7UUFDN0Isc0JBQWlCLEdBQVcsZUFBTSxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQTtJQThCdkQsQ0FBQztJQWxFQyxTQUFTLENBQUMsV0FBK0IsS0FBSztRQUM1QyxJQUFJLE1BQU0sR0FBVyxLQUFLLENBQUMsU0FBUyxDQUFDLFFBQVEsQ0FBQyxDQUFBO1FBQzlDLHlCQUNLLE1BQU0sRUFDVjtJQUNILENBQUM7SUFDRCxXQUFXLENBQUMsTUFBYyxFQUFFLFdBQStCLEtBQUs7UUFDOUQsS0FBSyxDQUFDLFdBQVcsQ0FBQyxNQUFNLEVBQUUsUUFBUSxDQUFDLENBQUE7SUFDckMsQ0FBQztJQUVEOzs7O09BSUc7SUFDSCxlQUFlLENBQUMsS0FBYTtRQUMzQixNQUFNLGNBQWMsR0FBVyxJQUFJLENBQUMsb0JBQW9CLEVBQUUsQ0FBQTtRQUMxRCxJQUFJLENBQUMsaUJBQWlCLENBQUMsV0FBVyxDQUFDLGNBQWMsR0FBRyxDQUFDLEVBQUUsQ0FBQyxFQUFFLENBQUMsQ0FBQyxDQUFBO1FBQzVELElBQUksQ0FBQyxjQUFjLENBQUMsSUFBSSxDQUFDLEtBQUssQ0FBQyxDQUFBO0lBQ2pDLENBQUM7SUFFRDs7T0FFRztJQUNILG9CQUFvQjtRQUNsQixPQUFPLElBQUksQ0FBQyxpQkFBaUIsQ0FBQyxVQUFVLENBQUMsQ0FBQyxFQUFFLENBQUMsQ0FBQyxDQUFBO0lBQ2hELENBQUM7SUFFRDs7T0FFRztJQUNILGlCQUFpQjtRQUNmLE9BQU8sSUFBSSxDQUFDLGNBQWMsQ0FBQTtJQUM1QixDQUFDO0lBS0QsVUFBVSxDQUFDLEtBQWEsRUFBRSxTQUFpQixDQUFDO1FBQzFDLDhCQUE4QjtRQUM5QixNQUFNLElBQUksQ0FBQyxDQUFBO1FBQ1gsSUFBSSxDQUFDLGlCQUFpQixHQUFHLFFBQVEsQ0FBQyxRQUFRLENBQUMsS0FBSyxFQUFFLE1BQU0sRUFBRSxNQUFNLEdBQUcsQ0FBQyxDQUFDLENBQUE7UUFDckUsTUFBTSxJQUFJLENBQUMsQ0FBQTtRQUNYLEtBQUssSUFBSSxDQUFDLEdBQVcsQ0FBQyxFQUFFLENBQUMsR0FBRyxJQUFJLENBQUMsb0JBQW9CLEVBQUUsRUFBRSxDQUFDLEVBQUUsRUFBRTtZQUM1RCxJQUFJLENBQUMsY0FBYyxDQUFDLElBQUksQ0FBQyxRQUFRLENBQUMsUUFBUSxDQUFDLEtBQUssRUFBRSxNQUFNLEVBQUUsTUFBTSxHQUFHLENBQUMsQ0FBQyxDQUFDLENBQUE7WUFDdEUsTUFBTSxJQUFJLENBQUMsQ0FBQTtTQUNaO1FBQ0QsT0FBTyxNQUFNLENBQUE7SUFDZixDQUFDO0lBRUQ7O09BRUc7SUFDSCxRQUFRO1FBQ04sTUFBTSxTQUFTLEdBQVcsZUFBTSxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQTtRQUN6QyxTQUFTLENBQUMsV0FBVyxDQUFDLElBQUksQ0FBQyxPQUFPLEVBQUUsQ0FBQyxFQUFFLENBQUMsQ0FBQyxDQUFBO1FBQ3pDLE1BQU0saUJBQWlCLEdBQVcsZUFBTSxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQTtRQUNqRCxpQkFBaUIsQ0FBQyxVQUFVLENBQUMsSUFBSSxDQUFDLGNBQWMsQ0FBQyxNQUFNLEVBQUUsQ0FBQyxFQUFFLENBQUMsQ0FBQyxDQUFBO1FBQzlELE1BQU0sSUFBSSxHQUFhLENBQUMsU0FBUyxFQUFFLGlCQUFpQixDQUFDLENBQUE7UUFDckQsSUFBSSxLQUFLLEdBQVcsU0FBUyxDQUFDLE1BQU0sR0FBRyxpQkFBaUIsQ0FBQyxNQUFNLENBQUE7UUFDL0QsSUFBSSxDQUFDLGNBQWMsQ0FBQyxPQUFPLENBQUMsQ0FBQyxZQUFvQixFQUFFLENBQVMsRUFBUSxFQUFFO1lBQ3BFLEtBQUssSUFBSSxDQUFDLENBQUE7WUFDVixJQUFJLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxjQUFjLENBQUMsR0FBRyxDQUFDLEVBQUUsQ0FBQyxDQUFDLENBQUE7UUFDeEMsQ0FBQyxDQUFDLENBQUE7UUFDRixPQUFPLGVBQU0sQ0FBQyxNQUFNLENBQUMsSUFBSSxFQUFFLEtBQUssQ0FBQyxDQUFBO0lBQ25DLENBQUM7Q0FDRjtBQXRFRCxnQ0FzRUMiLCJzb3VyY2VzQ29udGVudCI6WyIvKipcbiAqIEBwYWNrYWdlRG9jdW1lbnRhdGlvblxuICogQG1vZHVsZSBBUEktUGxhdGZvcm1WTS1TdWJuZXRBdXRoXG4gKi9cbmltcG9ydCB7IEJ1ZmZlciB9IGZyb20gXCJidWZmZXIvXCJcbmltcG9ydCBCaW5Ub29scyBmcm9tIFwiLi4vLi4vdXRpbHMvYmludG9vbHNcIlxuaW1wb3J0IHsgU2VyaWFsaXphYmxlLCBTZXJpYWxpemVkRW5jb2RpbmcgfSBmcm9tIFwiLi4vLi4vdXRpbHNcIlxuaW1wb3J0IHsgUGxhdGZvcm1WTUNvbnN0YW50cyB9IGZyb20gXCIuXCJcblxuLyoqXG4gKiBAaWdub3JlXG4gKi9cbmNvbnN0IGJpbnRvb2xzOiBCaW5Ub29scyA9IEJpblRvb2xzLmdldEluc3RhbmNlKClcblxuZXhwb3J0IGNsYXNzIFN1Ym5ldEF1dGggZXh0ZW5kcyBTZXJpYWxpemFibGUge1xuICBwcm90ZWN0ZWQgX3R5cGVOYW1lID0gXCJTdWJuZXRBdXRoXCJcbiAgcHJvdGVjdGVkIF90eXBlSUQgPSBQbGF0Zm9ybVZNQ29uc3RhbnRzLlNVQk5FVEFVVEhcblxuICBzZXJpYWxpemUoZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpOiBvYmplY3Qge1xuICAgIGxldCBmaWVsZHM6IG9iamVjdCA9IHN1cGVyLnNlcmlhbGl6ZShlbmNvZGluZylcbiAgICByZXR1cm4ge1xuICAgICAgLi4uZmllbGRzXG4gICAgfVxuICB9XG4gIGRlc2VyaWFsaXplKGZpZWxkczogb2JqZWN0LCBlbmNvZGluZzogU2VyaWFsaXplZEVuY29kaW5nID0gXCJoZXhcIikge1xuICAgIHN1cGVyLmRlc2VyaWFsaXplKGZpZWxkcywgZW5jb2RpbmcpXG4gIH1cblxuICAvKipcbiAgICogQWRkIGFuIGFkZHJlc3MgaW5kZXggZm9yIFN1Ym5ldCBBdXRoIHNpZ25pbmdcbiAgICpcbiAgICogQHBhcmFtIGluZGV4IHRoZSBCdWZmZXIgb2YgdGhlIGFkZHJlc3MgaW5kZXggdG8gYWRkXG4gICAqL1xuICBhZGRBZGRyZXNzSW5kZXgoaW5kZXg6IEJ1ZmZlcik6IHZvaWQge1xuICAgIGNvbnN0IG51bUFkZHJJbmRpY2VzOiBudW1iZXIgPSB0aGlzLmdldE51bUFkZHJlc3NJbmRpY2VzKClcbiAgICB0aGlzLm51bUFkZHJlc3NJbmRpY2VzLndyaXRlVUludEJFKG51bUFkZHJJbmRpY2VzICsgMSwgMCwgNClcbiAgICB0aGlzLmFkZHJlc3NJbmRpY2VzLnB1c2goaW5kZXgpXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgbnVtYmVyIG9mIGFkZHJlc3MgaW5kaWNlcyBhcyBhIG51bWJlclxuICAgKi9cbiAgZ2V0TnVtQWRkcmVzc0luZGljZXMoKTogbnVtYmVyIHtcbiAgICByZXR1cm4gdGhpcy5udW1BZGRyZXNzSW5kaWNlcy5yZWFkVUludEJFKDAsIDQpXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyBhbiBhcnJheSBvZiBBZGRyZXNzSW5kaWNlcyBhcyBCdWZmZXJzXG4gICAqL1xuICBnZXRBZGRyZXNzSW5kaWNlcygpOiBCdWZmZXJbXSB7XG4gICAgcmV0dXJuIHRoaXMuYWRkcmVzc0luZGljZXNcbiAgfVxuXG4gIHByb3RlY3RlZCBhZGRyZXNzSW5kaWNlczogQnVmZmVyW10gPSBbXVxuICBwcm90ZWN0ZWQgbnVtQWRkcmVzc0luZGljZXM6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYyg0KVxuXG4gIGZyb21CdWZmZXIoYnl0ZXM6IEJ1ZmZlciwgb2Zmc2V0OiBudW1iZXIgPSAwKTogbnVtYmVyIHtcbiAgICAvLyBpbmNyZWFzZSBvZmZzZXQgZm9yIHR5cGUgaWRcbiAgICBvZmZzZXQgKz0gNFxuICAgIHRoaXMubnVtQWRkcmVzc0luZGljZXMgPSBiaW50b29scy5jb3B5RnJvbShieXRlcywgb2Zmc2V0LCBvZmZzZXQgKyA0KVxuICAgIG9mZnNldCArPSA0XG4gICAgZm9yIChsZXQgaTogbnVtYmVyID0gMDsgaSA8IHRoaXMuZ2V0TnVtQWRkcmVzc0luZGljZXMoKTsgaSsrKSB7XG4gICAgICB0aGlzLmFkZHJlc3NJbmRpY2VzLnB1c2goYmludG9vbHMuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgNCkpXG4gICAgICBvZmZzZXQgKz0gNFxuICAgIH1cbiAgICByZXR1cm4gb2Zmc2V0XG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IHJlcHJlc2VudGF0aW9uIG9mIHRoZSBbW1N1Ym5ldEF1dGhdXS5cbiAgICovXG4gIHRvQnVmZmVyKCk6IEJ1ZmZlciB7XG4gICAgY29uc3QgdHlwZUlEQnVmOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoNClcbiAgICB0eXBlSURCdWYud3JpdGVVSW50QkUodGhpcy5fdHlwZUlELCAwLCA0KVxuICAgIGNvbnN0IG51bUFkZHJlc3NJbmRpY2VzOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoNClcbiAgICBudW1BZGRyZXNzSW5kaWNlcy53cml0ZUludEJFKHRoaXMuYWRkcmVzc0luZGljZXMubGVuZ3RoLCAwLCA0KVxuICAgIGNvbnN0IGJhcnI6IEJ1ZmZlcltdID0gW3R5cGVJREJ1ZiwgbnVtQWRkcmVzc0luZGljZXNdXG4gICAgbGV0IGJzaXplOiBudW1iZXIgPSB0eXBlSURCdWYubGVuZ3RoICsgbnVtQWRkcmVzc0luZGljZXMubGVuZ3RoXG4gICAgdGhpcy5hZGRyZXNzSW5kaWNlcy5mb3JFYWNoKChhZGRyZXNzSW5kZXg6IEJ1ZmZlciwgaTogbnVtYmVyKTogdm9pZCA9PiB7XG4gICAgICBic2l6ZSArPSA0XG4gICAgICBiYXJyLnB1c2godGhpcy5hZGRyZXNzSW5kaWNlc1tgJHtpfWBdKVxuICAgIH0pXG4gICAgcmV0dXJuIEJ1ZmZlci5jb25jYXQoYmFyciwgYnNpemUpXG4gIH1cbn1cbiJdfQ==