"use strict";
/**
 * @packageDocumentation
 * @module API-EVM-ExportTx
 */
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.ExportTx = void 0;
const buffer_1 = require("buffer/");
const bintools_1 = __importDefault(require("../../utils/bintools"));
const constants_1 = require("./constants");
const basetx_1 = require("./basetx");
const credentials_1 = require("./credentials");
const credentials_2 = require("../../common/credentials");
const inputs_1 = require("./inputs");
const serialization_1 = require("../../utils/serialization");
const outputs_1 = require("./outputs");
const errors_1 = require("../../utils/errors");
/**
 * @ignore
 */
const bintools = bintools_1.default.getInstance();
const serializer = serialization_1.Serialization.getInstance();
class ExportTx extends basetx_1.EVMBaseTx {
    /**
     * Class representing a ExportTx.
     *
     * @param networkID Optional networkID
     * @param blockchainID Optional blockchainID, default Buffer.alloc(32, 16)
     * @param destinationChain Optional destinationChain, default Buffer.alloc(32, 16)
     * @param inputs Optional array of the [[EVMInputs]]s
     * @param exportedOutputs Optional array of the [[EVMOutputs]]s
     */
    constructor(networkID = undefined, blockchainID = buffer_1.Buffer.alloc(32, 16), destinationChain = buffer_1.Buffer.alloc(32, 16), inputs = undefined, exportedOutputs = undefined) {
        super(networkID, blockchainID);
        this._typeName = "ExportTx";
        this._typeID = constants_1.EVMConstants.EXPORTTX;
        this.destinationChain = buffer_1.Buffer.alloc(32);
        this.numInputs = buffer_1.Buffer.alloc(4);
        this.inputs = [];
        this.numExportedOutputs = buffer_1.Buffer.alloc(4);
        this.exportedOutputs = [];
        this.destinationChain = destinationChain;
        if (typeof inputs !== "undefined" && Array.isArray(inputs)) {
            inputs.forEach((input) => {
                if (!(input instanceof inputs_1.EVMInput)) {
                    throw new errors_1.EVMInputError("Error - ExportTx.constructor: invalid EVMInput in array parameter 'inputs'");
                }
            });
            if (inputs.length > 1) {
                inputs = inputs.sort(inputs_1.EVMInput.comparator());
            }
            this.inputs = inputs;
        }
        if (typeof exportedOutputs !== "undefined" &&
            Array.isArray(exportedOutputs)) {
            exportedOutputs.forEach((exportedOutput) => {
                if (!(exportedOutput instanceof outputs_1.TransferableOutput)) {
                    throw new errors_1.TransferableOutputError("Error - ExportTx.constructor: TransferableOutput EVMInput in array parameter 'exportedOutputs'");
                }
            });
            this.exportedOutputs = exportedOutputs;
        }
    }
    serialize(encoding = "hex") {
        let fields = super.serialize(encoding);
        return Object.assign(Object.assign({}, fields), { destinationChain: serializer.encoder(this.destinationChain, encoding, "Buffer", "cb58"), exportedOutputs: this.exportedOutputs.map((i) => i.serialize(encoding)) });
    }
    deserialize(fields, encoding = "hex") {
        super.deserialize(fields, encoding);
        this.destinationChain = serializer.decoder(fields["destinationChain"], encoding, "cb58", "Buffer", 32);
        this.exportedOutputs = fields["exportedOutputs"].map((i) => {
            let eo = new outputs_1.TransferableOutput();
            eo.deserialize(i, encoding);
            return eo;
        });
        this.numExportedOutputs = buffer_1.Buffer.alloc(4);
        this.numExportedOutputs.writeUInt32BE(this.exportedOutputs.length, 0);
    }
    /**
     * Returns the destinationChain as a {@link https://github.com/feross/buffer|Buffer}
     */
    getDestinationChain() {
        return this.destinationChain;
    }
    /**
     * Returns the inputs as an array of [[EVMInputs]]
     */
    getInputs() {
        return this.inputs;
    }
    /**
     * Returns the outs as an array of [[EVMOutputs]]
     */
    getExportedOutputs() {
        return this.exportedOutputs;
    }
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[ExportTx]].
     */
    toBuffer() {
        if (typeof this.destinationChain === "undefined") {
            throw new errors_1.ChainIdError("ExportTx.toBuffer -- this.destinationChain is undefined");
        }
        this.numInputs.writeUInt32BE(this.inputs.length, 0);
        this.numExportedOutputs.writeUInt32BE(this.exportedOutputs.length, 0);
        let barr = [
            super.toBuffer(),
            this.destinationChain,
            this.numInputs
        ];
        let bsize = super.toBuffer().length +
            this.destinationChain.length +
            this.numInputs.length;
        this.inputs.forEach((importIn) => {
            bsize += importIn.toBuffer().length;
            barr.push(importIn.toBuffer());
        });
        bsize += this.numExportedOutputs.length;
        barr.push(this.numExportedOutputs);
        this.exportedOutputs.forEach((out) => {
            bsize += out.toBuffer().length;
            barr.push(out.toBuffer());
        });
        return buffer_1.Buffer.concat(barr, bsize);
    }
    /**
     * Decodes the [[ExportTx]] as a {@link https://github.com/feross/buffer|Buffer} and returns the size.
     */
    fromBuffer(bytes, offset = 0) {
        offset = super.fromBuffer(bytes, offset);
        this.destinationChain = bintools.copyFrom(bytes, offset, offset + 32);
        offset += 32;
        this.numInputs = bintools.copyFrom(bytes, offset, offset + 4);
        offset += 4;
        const numInputs = this.numInputs.readUInt32BE(0);
        for (let i = 0; i < numInputs; i++) {
            const anIn = new inputs_1.EVMInput();
            offset = anIn.fromBuffer(bytes, offset);
            this.inputs.push(anIn);
        }
        this.numExportedOutputs = bintools.copyFrom(bytes, offset, offset + 4);
        offset += 4;
        const numExportedOutputs = this.numExportedOutputs.readUInt32BE(0);
        for (let i = 0; i < numExportedOutputs; i++) {
            const anOut = new outputs_1.TransferableOutput();
            offset = anOut.fromBuffer(bytes, offset);
            this.exportedOutputs.push(anOut);
        }
        return offset;
    }
    /**
     * Returns a base-58 representation of the [[ExportTx]].
     */
    toString() {
        return bintools.bufferToB58(this.toBuffer());
    }
    /**
     * Takes the bytes of an [[UnsignedTx]] and returns an array of [[Credential]]s
     *
     * @param msg A Buffer for the [[UnsignedTx]]
     * @param kc An [[KeyChain]] used in signing
     *
     * @returns An array of [[Credential]]s
     */
    sign(msg, kc) {
        const creds = super.sign(msg, kc);
        this.inputs.forEach((input) => {
            const cred = (0, credentials_1.SelectCredentialClass)(input.getCredentialID());
            const sigidxs = input.getSigIdxs();
            sigidxs.forEach((sigidx) => {
                const keypair = kc.getKey(sigidx.getSource());
                const signval = keypair.sign(msg);
                const sig = new credentials_2.Signature();
                sig.fromBuffer(signval);
                cred.addSignature(sig);
            });
            creds.push(cred);
        });
        return creds;
    }
}
exports.ExportTx = ExportTx;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZXhwb3J0dHguanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi9zcmMvYXBpcy9ldm0vZXhwb3J0dHgudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IjtBQUFBOzs7R0FHRzs7Ozs7O0FBRUgsb0NBQWdDO0FBQ2hDLG9FQUEyQztBQUMzQywyQ0FBMEM7QUFFMUMscUNBQW9DO0FBQ3BDLCtDQUFxRDtBQUNyRCwwREFBd0U7QUFDeEUscUNBQW1DO0FBQ25DLDZEQUE2RTtBQUM3RSx1Q0FBOEM7QUFDOUMsK0NBSTJCO0FBRTNCOztHQUVHO0FBQ0gsTUFBTSxRQUFRLEdBQWEsa0JBQVEsQ0FBQyxXQUFXLEVBQUUsQ0FBQTtBQUNqRCxNQUFNLFVBQVUsR0FBa0IsNkJBQWEsQ0FBQyxXQUFXLEVBQUUsQ0FBQTtBQUU3RCxNQUFhLFFBQVMsU0FBUSxrQkFBUztJQXlKckM7Ozs7Ozs7O09BUUc7SUFDSCxZQUNFLFlBQW9CLFNBQVMsRUFDN0IsZUFBdUIsZUFBTSxDQUFDLEtBQUssQ0FBQyxFQUFFLEVBQUUsRUFBRSxDQUFDLEVBQzNDLG1CQUEyQixlQUFNLENBQUMsS0FBSyxDQUFDLEVBQUUsRUFBRSxFQUFFLENBQUMsRUFDL0MsU0FBcUIsU0FBUyxFQUM5QixrQkFBd0MsU0FBUztRQUVqRCxLQUFLLENBQUMsU0FBUyxFQUFFLFlBQVksQ0FBQyxDQUFBO1FBeEt0QixjQUFTLEdBQUcsVUFBVSxDQUFBO1FBQ3RCLFlBQU8sR0FBRyx3QkFBWSxDQUFDLFFBQVEsQ0FBQTtRQWlDL0IscUJBQWdCLEdBQVcsZUFBTSxDQUFDLEtBQUssQ0FBQyxFQUFFLENBQUMsQ0FBQTtRQUMzQyxjQUFTLEdBQVcsZUFBTSxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQTtRQUNuQyxXQUFNLEdBQWUsRUFBRSxDQUFBO1FBQ3ZCLHVCQUFrQixHQUFXLGVBQU0sQ0FBQyxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUE7UUFDNUMsb0JBQWUsR0FBeUIsRUFBRSxDQUFB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QyxNQUFNLEdBQUcsR0FBYyxJQUFJLHVCQUFTLEVBQUUsQ0FBQTtnQkFDdEMsR0FBRyxDQUFDLFVBQVUsQ0FBQyxPQUFPLENBQUMsQ0FBQTtnQkFDdkIsSUFBSSxDQUFDLFlBQVksQ0FBQyxHQUFHLENBQUMsQ0FBQTtZQUN4QixDQUFDLENBQUMsQ0FBQTtZQUNGLEtBQUssQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7UUFDbEIsQ0FBQyxDQUFDLENBQUE7UUFDRixPQUFPLEtBQUssQ0FBQTtJQUNkLENBQUM7Q0ErQ0Y7QUF0TUQsNEJBc01DIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgQVBJLUVWTS1FeHBvcnRUeFxuICovXG5cbmltcG9ydCB7IEJ1ZmZlciB9IGZyb20gXCJidWZmZXIvXCJcbmltcG9ydCBCaW5Ub29scyBmcm9tIFwiLi4vLi4vdXRpbHMvYmludG9vbHNcIlxuaW1wb3J0IHsgRVZNQ29uc3RhbnRzIH0gZnJvbSBcIi4vY29uc3RhbnRzXCJcbmltcG9ydCB7IEtleUNoYWluLCBLZXlQYWlyIH0gZnJvbSBcIi4va2V5Y2hhaW5cIlxuaW1wb3J0IHsgRVZNQmFzZVR4IH0gZnJvbSBcIi4vYmFzZXR4XCJcbmltcG9ydCB7IFNlbGVjdENyZWRlbnRpYWxDbGFzcyB9IGZyb20gXCIuL2NyZWRlbnRpYWxzXCJcbmltcG9ydCB7IFNpZ25hdHVyZSwgU2lnSWR4LCBDcmVkZW50aWFsIH0gZnJvbSBcIi4uLy4uL2NvbW1vbi9jcmVkZW50aWFsc1wiXG5pbXBvcnQgeyBFVk1JbnB1dCB9IGZyb20gXCIuL2lucHV0c1wiXG5pbXBvcnQgeyBTZXJpYWxpemF0aW9uLCBTZXJpYWxpemVkRW5jb2RpbmcgfSBmcm9tIFwiLi4vLi4vdXRpbHMvc2VyaWFsaXphdGlvblwiXG5pbXBvcnQgeyBUcmFuc2ZlcmFibGVPdXRwdXQgfSBmcm9tIFwiLi9vdXRwdXRzXCJcbmltcG9ydCB7XG4gIENoYWluSWRFcnJvcixcbiAgRVZNSW5wdXRFcnJvcixcbiAgVHJhbnNmZXJhYmxlT3V0cHV0RXJyb3Jcbn0gZnJvbSBcIi4uLy4uL3V0aWxzL2Vycm9yc1wiXG5cbi8qKlxuICogQGlnbm9yZVxuICovXG5jb25zdCBiaW50b29sczogQmluVG9vbHMgPSBCaW5Ub29scy5nZXRJbnN0YW5jZSgpXG5jb25zdCBzZXJpYWxpemVyOiBTZXJpYWxpemF0aW9uID0gU2VyaWFsaXphdGlvbi5nZXRJbnN0YW5jZSgpXG5cbmV4cG9ydCBjbGFzcyBFeHBvcnRUeCBleHRlbmRzIEVWTUJhc2VUeCB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIkV4cG9ydFR4XCJcbiAgcHJvdGVjdGVkIF90eXBlSUQgPSBFVk1Db25zdGFudHMuRVhQT1JUVFhcblxuICBzZXJpYWxpemUoZW5jb2Rpbmc6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiaGV4XCIpOiBvYmplY3Qge1xuICAgIGxldCBmaWVsZHM6IG9iamVjdCA9IHN1cGVyLnNlcmlhbGl6ZShlbmNvZGluZylcbiAgICByZXR1cm4ge1xuICAgICAgLi4uZmllbGRzLFxuICAgICAgZGVzdGluYXRpb25DaGFpbjogc2VyaWFsaXplci5lbmNvZGVyKFxuICAgICAgICB0aGlzLmRlc3RpbmF0aW9uQ2hhaW4sXG4gICAgICAgIGVuY29kaW5nLFxuICAgICAgICBcIkJ1ZmZlclwiLFxuICAgICAgICBcImNiNThcIlxuICAgICAgKSxcbiAgICAgIGV4cG9ydGVkT3V0cHV0czogdGhpcy5leHBvcnRlZE91dHB1dHMubWFwKChpKSA9PiBpLnNlcmlhbGl6ZShlbmNvZGluZykpXG4gICAgfVxuICB9XG4gIGRlc2VyaWFsaXplKGZpZWxkczogb2JqZWN0LCBlbmNvZGluZzogU2VyaWFsaXplZEVuY29kaW5nID0gXCJoZXhcIikge1xuICAgIHN1cGVyLmRlc2VyaWFsaXplKGZpZWxkcywgZW5jb2RpbmcpXG4gICAgdGhpcy5kZXN0aW5hdGlvbkNoYWluID0gc2VyaWFsaXplci5kZWNvZGVyKFxuICAgICAgZmllbGRzW1wiZGVzdGluYXRpb25DaGFpblwiXSxcbiAgICAgIGVuY29kaW5nLFxuICAgICAgXCJjYjU4XCIsXG4gICAgICBcIkJ1ZmZlclwiLFxuICAgICAgMzJcbiAgICApXG4gICAgdGhpcy5leHBvcnRlZE91dHB1dHMgPSBmaWVsZHNbXCJleHBvcnRlZE91dHB1dHNcIl0ubWFwKChpOiBvYmplY3QpID0+IHtcbiAgICAgIGxldCBlbzogVHJhbnNmZXJhYmxlT3V0cHV0ID0gbmV3IFRyYW5zZmVyYWJsZU91dHB1dCgpXG4gICAgICBlby5kZXNlcmlhbGl6ZShpLCBlbmNvZGluZylcbiAgICAgIHJldHVybiBlb1xuICAgIH0pXG4gICAgdGhpcy5udW1FeHBvcnRlZE91dHB1dHMgPSBCdWZmZXIuYWxsb2MoNClcbiAgICB0aGlzLm51bUV4cG9ydGVkT3V0cHV0cy53cml0ZVVJbnQzMkJFKHRoaXMuZXhwb3J0ZWRPdXRwdXRzLmxlbmd0aCwgMClcbiAgfVxuXG4gIHByb3RlY3RlZCBkZXN0aW5hdGlvbkNoYWluOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoMzIpXG4gIHByb3RlY3RlZCBudW1JbnB1dHM6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYyg0KVxuICBwcm90ZWN0ZWQgaW5wdXRzOiBFVk1JbnB1dFtdID0gW11cbiAgcHJvdGVjdGVkIG51bUV4cG9ydGVkT3V0cHV0czogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDQpXG4gIHByb3RlY3RlZCBleHBvcnRlZE91dHB1dHM6IFRyYW5zZmVyYWJsZU91dHB1dFtdID0gW11cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgZGVzdGluYXRpb25DaGFpbiBhcyBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9XG4gICAqL1xuICBnZXREZXN0aW5hdGlvbkNoYWluKCk6IEJ1ZmZlciB7XG4gICAgcmV0dXJuIHRoaXMuZGVzdGluYXRpb25DaGFpblxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIGlucHV0cyBhcyBhbiBhcnJheSBvZiBbW0VWTUlucHV0c11dXG4gICAqL1xuICBnZXRJbnB1dHMoKTogRVZNSW5wdXRbXSB7XG4gICAgcmV0dXJuIHRoaXMuaW5wdXRzXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgb3V0cyBhcyBhbiBhcnJheSBvZiBbW0VWTU91dHB1dHNdXVxuICAgKi9cbiAgZ2V0RXhwb3J0ZWRPdXRwdXRzKCk6IFRyYW5zZmVyYWJsZU91dHB1dFtdIHtcbiAgICByZXR1cm4gdGhpcy5leHBvcnRlZE91dHB1dHNcbiAgfVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIGEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gcmVwcmVzZW50YXRpb24gb2YgdGhlIFtbRXhwb3J0VHhdXS5cbiAgICovXG4gIHRvQnVmZmVyKCk6IEJ1ZmZlciB7XG4gICAgaWYgKHR5cGVvZiB0aGlzLmRlc3RpbmF0aW9uQ2hhaW4gPT09IFwidW5kZWZpbmVkXCIpIHtcbiAgICAgIHRocm93IG5ldyBDaGFpbklkRXJyb3IoXG4gICAgICAgIFwiRXhwb3J0VHgudG9CdWZmZXIgLS0gdGhpcy5kZXN0aW5hdGlvbkNoYWluIGlzIHVuZGVmaW5lZFwiXG4gICAgICApXG4gICAgfVxuICAgIHRoaXMubnVtSW5wdXRzLndyaXRlVUludDMyQkUodGhpcy5pbnB1dHMubGVuZ3RoLCAwKVxuICAgIHRoaXMubnVtRXhwb3J0ZWRPdXRwdXRzLndyaXRlVUludDMyQkUodGhpcy5leHBvcnRlZE91dHB1dHMubGVuZ3RoLCAwKVxuICAgIGxldCBiYXJyOiBCdWZmZXJbXSA9IFtcbiAgICAgIHN1cGVyLnRvQnVmZmVyKCksXG4gICAgICB0aGlzLmRlc3RpbmF0aW9uQ2hhaW4sXG4gICAgICB0aGlzLm51bUlucHV0c1xuICAgIF1cbiAgICBsZXQgYnNpemU6IG51bWJlciA9XG4gICAgICBzdXBlci50b0J1ZmZlcigpLmxlbmd0aCArXG4gICAgICB0aGlzLmRlc3RpbmF0aW9uQ2hhaW4ubGVuZ3RoICtcbiAgICAgIHRoaXMubnVtSW5wdXRzLmxlbmd0aFxuICAgIHRoaXMuaW5wdXRzLmZvckVhY2goKGltcG9ydEluOiBFVk1JbnB1dCkgPT4ge1xuICAgICAgYnNpemUgKz0gaW1wb3J0SW4udG9CdWZmZXIoKS5sZW5ndGhcbiAgICAgIGJhcnIucHVzaChpbXBvcnRJbi50b0J1ZmZlcigpKVxuICAgIH0pXG4gICAgYnNpemUgKz0gdGhpcy5udW1FeHBvcnRlZE91dHB1dHMubGVuZ3RoXG4gICAgYmFyci5wdXNoKHRoaXMubnVtRXhwb3J0ZWRPdXRwdXRzKVxuICAgIHRoaXMuZXhwb3J0ZWRPdXRwdXRzLmZvckVhY2goKG91dDogVHJhbnNmZXJhYmxlT3V0cHV0KSA9PiB7XG4gICAgICBic2l6ZSArPSBvdXQudG9CdWZmZXIoKS5sZW5ndGhcbiAgICAgIGJhcnIucHVzaChvdXQudG9CdWZmZXIoKSlcbiAgICB9KVxuICAgIHJldHVybiBCdWZmZXIuY29uY2F0KGJhcnIsIGJzaXplKVxuICB9XG5cbiAgLyoqXG4gICAqIERlY29kZXMgdGhlIFtbRXhwb3J0VHhdXSBhcyBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IGFuZCByZXR1cm5zIHRoZSBzaXplLlxuICAgKi9cbiAgZnJvbUJ1ZmZlcihieXRlczogQnVmZmVyLCBvZmZzZXQ6IG51bWJlciA9IDApOiBudW1iZXIge1xuICAgIG9mZnNldCA9IHN1cGVyLmZyb21CdWZmZXIoYnl0ZXMsIG9mZnNldClcbiAgICB0aGlzLmRlc3RpbmF0aW9uQ2hhaW4gPSBiaW50b29scy5jb3B5RnJvbShieXRlcywgb2Zmc2V0LCBvZmZzZXQgKyAzMilcbiAgICBvZmZzZXQgKz0gMzJcbiAgICB0aGlzLm51bUlucHV0cyA9IGJpbnRvb2xzLmNvcHlGcm9tKGJ5dGVzLCBvZmZzZXQsIG9mZnNldCArIDQpXG4gICAgb2Zmc2V0ICs9IDRcbiAgICBjb25zdCBudW1JbnB1dHM6IG51bWJlciA9IHRoaXMubnVtSW5wdXRzLnJlYWRVSW50MzJCRSgwKVxuICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCBudW1JbnB1dHM7IGkrKykge1xuICAgICAgY29uc3QgYW5JbjogRVZNSW5wdXQgPSBuZXcgRVZNSW5wdXQoKVxuICAgICAgb2Zmc2V0ID0gYW5Jbi5mcm9tQnVmZmVyKGJ5dGVzLCBvZmZzZXQpXG4gICAgICB0aGlzLmlucHV0cy5wdXNoKGFuSW4pXG4gICAgfVxuICAgIHRoaXMubnVtRXhwb3J0ZWRPdXRwdXRzID0gYmludG9vbHMuY29weUZyb20oYnl0ZXMsIG9mZnNldCwgb2Zmc2V0ICsgNClcbiAgICBvZmZzZXQgKz0gNFxuICAgIGNvbnN0IG51bUV4cG9ydGVkT3V0cHV0czogbnVtYmVyID0gdGhpcy5udW1FeHBvcnRlZE91dHB1dHMucmVhZFVJbnQzMkJFKDApXG4gICAgZm9yIChsZXQgaTogbnVtYmVyID0gMDsgaSA8IG51bUV4cG9ydGVkT3V0cHV0czsgaSsrKSB7XG4gICAgICBjb25zdCBhbk91dDogVHJhbnNmZXJhYmxlT3V0cHV0ID0gbmV3IFRyYW5zZmVyYWJsZU91dHB1dCgpXG4gICAgICBvZmZzZXQgPSBhbk91dC5mcm9tQnVmZmVyKGJ5dGVzLCBvZmZzZXQpXG4gICAgICB0aGlzLmV4cG9ydGVkT3V0cHV0cy5wdXNoKGFuT3V0KVxuICAgIH1cbiAgICByZXR1cm4gb2Zmc2V0XG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyBhIGJhc2UtNTggcmVwcmVzZW50YXRpb24gb2YgdGhlIFtbRXhwb3J0VHhdXS5cbiAgICovXG4gIHRvU3RyaW5nKCk6IHN0cmluZyB7XG4gICAgcmV0dXJuIGJpbnRvb2xzLmJ1ZmZlclRvQjU4KHRoaXMudG9CdWZmZXIoKSlcbiAgfVxuXG4gIC8qKlxuICAgKiBUYWtlcyB0aGUgYnl0ZXMgb2YgYW4gW1tVbnNpZ25lZFR4XV0gYW5kIHJldHVybnMgYW4gYXJyYXkgb2YgW1tDcmVkZW50aWFsXV1zXG4gICAqXG4gICAqIEBwYXJhbSBtc2cgQSBCdWZmZXIgZm9yIHRoZSBbW1Vuc2lnbmVkVHhdXVxuICAgKiBAcGFyYW0ga2MgQW4gW1tLZXlDaGFpbl1dIHVzZWQgaW4gc2lnbmluZ1xuICAgKlxuICAgKiBAcmV0dXJucyBBbiBhcnJheSBvZiBbW0NyZWRlbnRpYWxdXXNcbiAgICovXG4gIHNpZ24obXNnOiBCdWZmZXIsIGtjOiBLZXlDaGFpbik6IENyZWRlbnRpYWxbXSB7XG4gICAgY29uc3QgY3JlZHM6IENyZWRlbnRpYWxbXSA9IHN1cGVyLnNpZ24obXNnLCBrYylcbiAgICB0aGlzLmlucHV0cy5mb3JFYWNoKChpbnB1dDogRVZNSW5wdXQpID0+IHtcbiAgICAgIGNvbnN0IGNyZWQ6IENyZWRlbnRpYWwgPSBTZWxlY3RDcmVkZW50aWFsQ2xhc3MoaW5wdXQuZ2V0Q3JlZGVudGlhbElEKCkpXG4gICAgICBjb25zdCBzaWdpZHhzOiBTaWdJZHhbXSA9IGlucHV0LmdldFNpZ0lkeHMoKVxuICAgICAgc2lnaWR4cy5mb3JFYWNoKChzaWdpZHg6IFNpZ0lkeCkgPT4ge1xuICAgICAgICBjb25zdCBrZXlwYWlyOiBLZXlQYWlyID0ga2MuZ2V0S2V5KHNpZ2lkeC5nZXRTb3VyY2UoKSlcbiAgICAgICAgY29uc3Qgc2lnbnZhbDogQnVmZmVyID0ga2V5cGFpci5zaWduKG1zZylcbiAgICAgICAgY29uc3Qgc2lnOiBTaWduYXR1cmUgPSBuZXcgU2lnbmF0dXJlKClcbiAgICAgICAgc2lnLmZyb21CdWZmZXIoc2lnbnZhbClcbiAgICAgICAgY3JlZC5hZGRTaWduYXR1cmUoc2lnKVxuICAgICAgfSlcbiAgICAgIGNyZWRzLnB1c2goY3JlZClcbiAgICB9KVxuICAgIHJldHVybiBjcmVkc1xuICB9XG5cbiAgLyoqXG4gICAqIENsYXNzIHJlcHJlc2VudGluZyBhIEV4cG9ydFR4LlxuICAgKlxuICAgKiBAcGFyYW0gbmV0d29ya0lEIE9wdGlvbmFsIG5ldHdvcmtJRFxuICAgKiBAcGFyYW0gYmxvY2tjaGFpbklEIE9wdGlvbmFsIGJsb2NrY2hhaW5JRCwgZGVmYXVsdCBCdWZmZXIuYWxsb2MoMzIsIDE2KVxuICAgKiBAcGFyYW0gZGVzdGluYXRpb25DaGFpbiBPcHRpb25hbCBkZXN0aW5hdGlvbkNoYWluLCBkZWZhdWx0IEJ1ZmZlci5hbGxvYygzMiwgMTYpXG4gICAqIEBwYXJhbSBpbnB1dHMgT3B0aW9uYWwgYXJyYXkgb2YgdGhlIFtbRVZNSW5wdXRzXV1zXG4gICAqIEBwYXJhbSBleHBvcnRlZE91dHB1dHMgT3B0aW9uYWwgYXJyYXkgb2YgdGhlIFtbRVZNT3V0cHV0c11dc1xuICAgKi9cbiAgY29uc3RydWN0b3IoXG4gICAgbmV0d29ya0lEOiBudW1iZXIgPSB1bmRlZmluZWQsXG4gICAgYmxvY2tjaGFpbklEOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoMzIsIDE2KSxcbiAgICBkZXN0aW5hdGlvbkNoYWluOiBCdWZmZXIgPSBCdWZmZXIuYWxsb2MoMzIsIDE2KSxcbiAgICBpbnB1dHM6IEVWTUlucHV0W10gPSB1bmRlZmluZWQsXG4gICAgZXhwb3J0ZWRPdXRwdXRzOiBUcmFuc2ZlcmFibGVPdXRwdXRbXSA9IHVuZGVmaW5lZFxuICApIHtcbiAgICBzdXBlcihuZXR3b3JrSUQsIGJsb2NrY2hhaW5JRClcbiAgICB0aGlzLmRlc3RpbmF0aW9uQ2hhaW4gPSBkZXN0aW5hdGlvbkNoYWluXG4gICAgaWYgKHR5cGVvZiBpbnB1dHMgIT09IFwidW5kZWZpbmVkXCIgJiYgQXJyYXkuaXNBcnJheShpbnB1dHMpKSB7XG4gICAgICBpbnB1dHMuZm9yRWFjaCgoaW5wdXQ6IEVWTUlucHV0KSA9PiB7XG4gICAgICAgIGlmICghKGlucHV0IGluc3RhbmNlb2YgRVZNSW5wdXQpKSB7XG4gICAgICAgICAgdGhyb3cgbmV3IEVWTUlucHV0RXJyb3IoXG4gICAgICAgICAgICBcIkVycm9yIC0gRXhwb3J0VHguY29uc3RydWN0b3I6IGludmFsaWQgRVZNSW5wdXQgaW4gYXJyYXkgcGFyYW1ldGVyICdpbnB1dHMnXCJcbiAgICAgICAgICApXG4gICAgICAgIH1cbiAgICAgIH0pXG4gICAgICBpZiAoaW5wdXRzLmxlbmd0aCA+IDEpIHtcbiAgICAgICAgaW5wdXRzID0gaW5wdXRzLnNvcnQoRVZNSW5wdXQuY29tcGFyYXRvcigpKVxuICAgICAgfVxuICAgICAgdGhpcy5pbnB1dHMgPSBpbnB1dHNcbiAgICB9XG4gICAgaWYgKFxuICAgICAgdHlwZW9mIGV4cG9ydGVkT3V0cHV0cyAhPT0gXCJ1bmRlZmluZWRcIiAmJlxuICAgICAgQXJyYXkuaXNBcnJheShleHBvcnRlZE91dHB1dHMpXG4gICAgKSB7XG4gICAgICBleHBvcnRlZE91dHB1dHMuZm9yRWFjaCgoZXhwb3J0ZWRPdXRwdXQ6IFRyYW5zZmVyYWJsZU91dHB1dCkgPT4ge1xuICAgICAgICBpZiAoIShleHBvcnRlZE91dHB1dCBpbnN0YW5jZW9mIFRyYW5zZmVyYWJsZU91dHB1dCkpIHtcbiAgICAgICAgICB0aHJvdyBuZXcgVHJhbnNmZXJhYmxlT3V0cHV0RXJyb3IoXG4gICAgICAgICAgICBcIkVycm9yIC0gRXhwb3J0VHguY29uc3RydWN0b3I6IFRyYW5zZmVyYWJsZU91dHB1dCBFVk1JbnB1dCBpbiBhcnJheSBwYXJhbWV0ZXIgJ2V4cG9ydGVkT3V0cHV0cydcIlxuICAgICAgICAgIClcbiAgICAgICAgfVxuICAgICAgfSlcbiAgICAgIHRoaXMuZXhwb3J0ZWRPdXRwdXRzID0gZXhwb3J0ZWRPdXRwdXRzXG4gICAgfVxuICB9XG59XG4iXX0=