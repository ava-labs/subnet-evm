"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.GenesisState = void 0;
/**
 * @packageDocumentation
 * @module API-AVM-CreateAssetTx
 */
const buffer_1 = require("buffer/");
const bintools_1 = __importDefault(require("../../utils/bintools"));
const constants_1 = require("./constants");
const constants_2 = require("../../utils/constants");
const serialization_1 = require("../../utils/serialization");
/**
 * @ignore
 */
const bintools = bintools_1.default.getInstance();
const serializer = serialization_1.Serialization.getInstance();
class GenesisState extends serialization_1.Serializable {
    /**
    * Class representing a GenesisState
    *
    * @param networkid Optional networkid, [[DefaultNetworkID]]
    * @param blockchainid Optional blockchainid, default Buffer.alloc(32, 16)
    */
    constructor(networkid = constants_2.DefaultNetworkID, blockchainid = buffer_1.Buffer.alloc(32)) {
        super();
        this._typeName = "GenesisState";
        this._codecID = constants_1.AVMConstants.LATESTCODEC;
        this.networkid = buffer_1.Buffer.alloc(4);
        this.blockchainid = buffer_1.Buffer.alloc(32);
        this.networkid.writeUInt32BE(networkid, 0);
        this.blockchainid = blockchainid;
    }
    serialize(encoding = "utf8") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { "networkid": serializer.encoder(this.networkid, encoding, "Buffer", "decimalString"), "blockchainid": serializer.encoder(this.blockchainid, encoding, "Buffer", "cb58") });
    }
    ;
}
exports.GenesisState = GenesisState;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZ2VuZXNpc3N0YXRlLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2FwaXMvYXZtL2dlbmVzaXNzdGF0ZS50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7QUFBQTs7O0dBR0c7QUFDSCxvQ0FBZ0M7QUFDaEMsb0VBQTJDO0FBQzNDLDJDQUEwQztBQUsxQyxxREFBd0Q7QUFDeEQsNkRBQTJGO0FBRzNGOztHQUVHO0FBQ0gsTUFBTSxRQUFRLEdBQUcsa0JBQVEsQ0FBQyxXQUFXLEVBQUUsQ0FBQTtBQUN2QyxNQUFNLFVBQVUsR0FBRyw2QkFBYSxDQUFDLFdBQVcsRUFBRSxDQUFBO0FBRTlDLE1BQWEsWUFBYSxTQUFRLDRCQUFZO0lBZTVDOzs7OztNQUtFO0lBQ0YsWUFBWSxZQUFvQiw0QkFBZ0IsRUFBRSxlQUF1QixlQUFNLENBQUMsS0FBSyxDQUFDLEVBQUUsQ0FBQztRQUN2RixLQUFLLEVBQUUsQ0FBQTtRQXJCQyxjQUFTLEdBQUcsY0FBYyxDQUFBO1FBQzFCLGFBQVEsR0FBRyx3QkFBWSxDQUFDLFdBQVcsQ0FBQTtRQUVuQyxjQUFTLEdBQVcsZUFBTSxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQTtRQUNuQyxpQkFBWSxHQUFXLGVBQU0sQ0FBQyxLQUFLLENBQUMsRUFBRSxDQUFDLENBQUE7UUFrQi9DLElBQUksQ0FBQyxTQUFTLENBQUMsYUFBYSxDQUFDLFNBQVMsRUFBRSxDQUFDLENBQUMsQ0FBQTtRQUMxQyxJQUFJLENBQUMsWUFBWSxHQUFHLFlBQVksQ0FBQTtJQUNsQyxDQUFDO0lBbEJELFNBQVMsQ0FBQyxXQUErQixNQUFNO1FBQzdDLElBQUksTUFBTSxHQUFXLEtBQUssQ0FBQyxTQUFTLENBQUMsUUFBUSxDQUFDLENBQUE7UUFDOUMsdUNBQ0ssTUFBTSxLQUNULFdBQVcsRUFBRSxVQUFVLENBQUMsT0FBTyxDQUFDLElBQUksQ0FBQyxTQUFTLEVBQUUsUUFBUSxFQUFFLFFBQVEsRUFBRSxlQUFlLENBQUMsRUFDcEYsY0FBYyxFQUFFLFVBQVUsQ0FBQyxPQUFPLENBQUMsSUFBSSxDQUFDLFlBQVksRUFBRSxRQUFRLEVBQUUsUUFBUSxFQUFFLE1BQU0sQ0FBQyxJQUNsRjtJQUNILENBQUM7SUFBQSxDQUFDO0NBWUg7QUExQkQsb0NBMEJDIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgQVBJLUFWTS1DcmVhdGVBc3NldFR4XG4gKi9cbmltcG9ydCB7IEJ1ZmZlciB9IGZyb20gJ2J1ZmZlci8nXG5pbXBvcnQgQmluVG9vbHMgZnJvbSAnLi4vLi4vdXRpbHMvYmludG9vbHMnXG5pbXBvcnQgeyBBVk1Db25zdGFudHMgfSBmcm9tICcuL2NvbnN0YW50cydcbmltcG9ydCB7IFRyYW5zZmVyYWJsZU91dHB1dCB9IGZyb20gJy4vb3V0cHV0cydcbmltcG9ydCB7IFRyYW5zZmVyYWJsZUlucHV0IH0gZnJvbSAnLi9pbnB1dHMnXG5pbXBvcnQgeyBJbml0aWFsU3RhdGVzIH0gZnJvbSAnLi9pbml0aWFsc3RhdGVzJ1xuaW1wb3J0IHsgQmFzZVR4IH0gZnJvbSAnLi9iYXNldHgnXG5pbXBvcnQgeyBEZWZhdWx0TmV0d29ya0lEIH0gZnJvbSAnLi4vLi4vdXRpbHMvY29uc3RhbnRzJ1xuaW1wb3J0IHsgU2VyaWFsaXphYmxlLCBTZXJpYWxpemF0aW9uLCBTZXJpYWxpemVkRW5jb2RpbmcgfSBmcm9tICcuLi8uLi91dGlscy9zZXJpYWxpemF0aW9uJ1xuaW1wb3J0IHsgQ29kZWNJZEVycm9yIH0gZnJvbSAnLi4vLi4vdXRpbHMvZXJyb3JzJ1xuXG4vKipcbiAqIEBpZ25vcmVcbiAqL1xuY29uc3QgYmludG9vbHMgPSBCaW5Ub29scy5nZXRJbnN0YW5jZSgpXG5jb25zdCBzZXJpYWxpemVyID0gU2VyaWFsaXphdGlvbi5nZXRJbnN0YW5jZSgpXG5cbmV4cG9ydCBjbGFzcyBHZW5lc2lzU3RhdGUgZXh0ZW5kcyBTZXJpYWxpemFibGUge1xuICBwcm90ZWN0ZWQgX3R5cGVOYW1lID0gXCJHZW5lc2lzU3RhdGVcIlxuICBwcm90ZWN0ZWQgX2NvZGVjSUQgPSBBVk1Db25zdGFudHMuTEFURVNUQ09ERUNcblxuICBwcm90ZWN0ZWQgbmV0d29ya2lkOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoNClcbiAgcHJvdGVjdGVkIGJsb2NrY2hhaW5pZDogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDMyKVxuXG4gIHNlcmlhbGl6ZShlbmNvZGluZzogU2VyaWFsaXplZEVuY29kaW5nID0gXCJ1dGY4XCIpOiBvYmplY3Qge1xuICAgIGxldCBmaWVsZHM6IG9iamVjdCA9IHN1cGVyLnNlcmlhbGl6ZShlbmNvZGluZylcbiAgICByZXR1cm4ge1xuICAgICAgLi4uZmllbGRzLFxuICAgICAgXCJuZXR3b3JraWRcIjogc2VyaWFsaXplci5lbmNvZGVyKHRoaXMubmV0d29ya2lkLCBlbmNvZGluZywgXCJCdWZmZXJcIiwgXCJkZWNpbWFsU3RyaW5nXCIpLFxuICAgICAgXCJibG9ja2NoYWluaWRcIjogc2VyaWFsaXplci5lbmNvZGVyKHRoaXMuYmxvY2tjaGFpbmlkLCBlbmNvZGluZywgXCJCdWZmZXJcIiwgXCJjYjU4XCIpXG4gICAgfVxuICB9O1xuICAvKipcbiAgKiBDbGFzcyByZXByZXNlbnRpbmcgYSBHZW5lc2lzU3RhdGVcbiAgKlxuICAqIEBwYXJhbSBuZXR3b3JraWQgT3B0aW9uYWwgbmV0d29ya2lkLCBbW0RlZmF1bHROZXR3b3JrSURdXVxuICAqIEBwYXJhbSBibG9ja2NoYWluaWQgT3B0aW9uYWwgYmxvY2tjaGFpbmlkLCBkZWZhdWx0IEJ1ZmZlci5hbGxvYygzMiwgMTYpXG4gICovXG4gIGNvbnN0cnVjdG9yKG5ldHdvcmtpZDogbnVtYmVyID0gRGVmYXVsdE5ldHdvcmtJRCwgYmxvY2tjaGFpbmlkOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoMzIpKSB7XG4gICAgc3VwZXIoKVxuICAgIHRoaXMubmV0d29ya2lkLndyaXRlVUludDMyQkUobmV0d29ya2lkLCAwKVxuICAgIHRoaXMuYmxvY2tjaGFpbmlkID0gYmxvY2tjaGFpbmlkXG4gIH1cbn0iXX0=