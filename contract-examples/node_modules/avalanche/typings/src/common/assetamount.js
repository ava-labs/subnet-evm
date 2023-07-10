"use strict";
/**
 * @packageDocumentation
 * @module Common-AssetAmount
 */
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.StandardAssetAmountDestination = exports.AssetAmount = void 0;
const buffer_1 = require("buffer/");
const bn_js_1 = __importDefault(require("bn.js"));
const errors_1 = require("../utils/errors");
/**
 * Class for managing asset amounts in the UTXOSet fee calcuation
 */
class AssetAmount {
    constructor(assetID, amount, burn) {
        // assetID that is amount is managing.
        this.assetID = buffer_1.Buffer.alloc(32);
        // amount of this asset that should be sent.
        this.amount = new bn_js_1.default(0);
        // burn is the amount of this asset that should be burned.
        this.burn = new bn_js_1.default(0);
        // spent is the total amount of this asset that has been consumed.
        this.spent = new bn_js_1.default(0);
        // stakeableLockSpent is the amount of this asset that has been consumed that
        // was locked.
        this.stakeableLockSpent = new bn_js_1.default(0);
        // change is the excess amount of this asset that was consumed over the amount
        // requested to be consumed(amount + burn).
        this.change = new bn_js_1.default(0);
        // stakeableLockChange is a flag to mark if the input that generated the
        // change was locked.
        this.stakeableLockChange = false;
        // finished is a convenience flag to track "spent >= amount + burn"
        this.finished = false;
        this.getAssetID = () => {
            return this.assetID;
        };
        this.getAssetIDString = () => {
            return this.assetID.toString("hex");
        };
        this.getAmount = () => {
            return this.amount;
        };
        this.getSpent = () => {
            return this.spent;
        };
        this.getBurn = () => {
            return this.burn;
        };
        this.getChange = () => {
            return this.change;
        };
        this.getStakeableLockSpent = () => {
            return this.stakeableLockSpent;
        };
        this.getStakeableLockChange = () => {
            return this.stakeableLockChange;
        };
        this.isFinished = () => {
            return this.finished;
        };
        // spendAmount should only be called if this asset is still awaiting more
        // funds to consume.
        this.spendAmount = (amt, stakeableLocked = false) => {
            if (this.finished) {
                /* istanbul ignore next */
                throw new errors_1.InsufficientFundsError("Error - AssetAmount.spendAmount: attempted to spend " + "excess funds");
            }
            this.spent = this.spent.add(amt);
            if (stakeableLocked) {
                this.stakeableLockSpent = this.stakeableLockSpent.add(amt);
            }
            const total = this.amount.add(this.burn);
            if (this.spent.gte(total)) {
                this.change = this.spent.sub(total);
                if (stakeableLocked) {
                    this.stakeableLockChange = true;
                }
                this.finished = true;
            }
            return this.finished;
        };
        this.assetID = assetID;
        this.amount = typeof amount === "undefined" ? new bn_js_1.default(0) : amount;
        this.burn = typeof burn === "undefined" ? new bn_js_1.default(0) : burn;
        this.spent = new bn_js_1.default(0);
        this.stakeableLockSpent = new bn_js_1.default(0);
        this.stakeableLockChange = false;
    }
}
exports.AssetAmount = AssetAmount;
class StandardAssetAmountDestination {
    constructor(destinations, senders, changeAddresses) {
        this.amounts = [];
        this.destinations = [];
        this.senders = [];
        this.changeAddresses = [];
        this.amountkey = {};
        this.inputs = [];
        this.outputs = [];
        this.change = [];
        // TODO: should this function allow for repeated calls with the same
        //       assetID?
        this.addAssetAmount = (assetID, amount, burn) => {
            let aa = new AssetAmount(assetID, amount, burn);
            this.amounts.push(aa);
            this.amountkey[aa.getAssetIDString()] = aa;
        };
        this.addInput = (input) => {
            this.inputs.push(input);
        };
        this.addOutput = (output) => {
            this.outputs.push(output);
        };
        this.addChange = (output) => {
            this.change.push(output);
        };
        this.getAmounts = () => {
            return this.amounts;
        };
        this.getDestinations = () => {
            return this.destinations;
        };
        this.getSenders = () => {
            return this.senders;
        };
        this.getChangeAddresses = () => {
            return this.changeAddresses;
        };
        this.getAssetAmount = (assetHexStr) => {
            return this.amountkey[`${assetHexStr}`];
        };
        this.assetExists = (assetHexStr) => {
            return assetHexStr in this.amountkey;
        };
        this.getInputs = () => {
            return this.inputs;
        };
        this.getOutputs = () => {
            return this.outputs;
        };
        this.getChangeOutputs = () => {
            return this.change;
        };
        this.getAllOutputs = () => {
            return this.outputs.concat(this.change);
        };
        this.canComplete = () => {
            for (let i = 0; i < this.amounts.length; i++) {
                if (!this.amounts[`${i}`].isFinished()) {
                    return false;
                }
            }
            return true;
        };
        this.destinations = destinations;
        this.changeAddresses = changeAddresses;
        this.senders = senders;
    }
}
exports.StandardAssetAmountDestination = StandardAssetAmountDestination;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYXNzZXRhbW91bnQuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi9zcmMvY29tbW9uL2Fzc2V0YW1vdW50LnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7QUFBQTs7O0dBR0c7Ozs7OztBQUVILG9DQUFnQztBQUNoQyxrREFBc0I7QUFHdEIsNENBQXdEO0FBRXhEOztHQUVHO0FBQ0gsTUFBYSxXQUFXO0lBcUZ0QixZQUFZLE9BQWUsRUFBRSxNQUFVLEVBQUUsSUFBUTtRQXBGakQsc0NBQXNDO1FBQzVCLFlBQU8sR0FBVyxlQUFNLENBQUMsS0FBSyxDQUFDLEVBQUUsQ0FBQyxDQUFBO1FBQzVDLDRDQUE0QztRQUNsQyxXQUFNLEdBQU8sSUFBSSxlQUFFLENBQUMsQ0FBQyxDQUFDLENBQUE7UUFDaEMsMERBQTBEO1FBQ2hELFNBQUksR0FBTyxJQUFJLGVBQUUsQ0FBQyxDQUFDLENBQUMsQ0FBQTtRQUU5QixrRUFBa0U7UUFDeEQsVUFBSyxHQUFPLElBQUksZUFBRSxDQUFDLENBQUMsQ0FBQyxDQUFBO1FBQy9CLDZFQUE2RTtRQUM3RSxjQUFjO1FBQ0osdUJBQWtCLEdBQU8sSUFBSSxlQUFFLENBQUMsQ0FBQyxDQUFDLENBQUE7UUFFNUMsOEVBQThFO1FBQzlFLDJDQUEyQztRQUNqQyxXQUFNLEdBQU8sSUFBSSxlQUFFLENBQUMsQ0FBQyxDQUFDLENBQUE7UUFDaEMsd0VBQXdFO1FBQ3hFLHFCQUFxQjtRQUNYLHdCQUFtQixHQUFZLEtBQUssQ0FBQTtRQUU5QyxtRUFBbUU7UUFDekQsYUFBUSxHQUFZLEtBQUssQ0FBQTtRQUVuQyxlQUFVLEdBQUcsR0FBVyxFQUFFO1lBQ3hCLE9BQU8sSUFBSSxDQUFDLE9BQU8sQ0FBQTtRQUNyQixDQUFDLENBQUE7UUFFRCxxQkFBZ0IsR0FBRyxHQUFXLEVBQUU7WUFDOUIsT0FBTyxJQUFJLENBQUMsT0FBTyxDQUFDLFFBQVEsQ0FBQyxLQUFLLENBQUMsQ0FBQTtRQUNyQyxDQUFDLENBQUE7UUFFRCxj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QUNiLENBQUMsQ0FBQTtRQU9DLElBQUksQ0FBQyxZQUFZLEdBQUcsWUFBWSxDQUFBO1FBQ2hDLElBQUksQ0FBQyxlQUFlLEdBQUcsZUFBZSxDQUFBO1FBQ3RDLElBQUksQ0FBQyxPQUFPLEdBQUcsT0FBTyxDQUFBO0lBQ3hCLENBQUM7Q0FDRjtBQTNGRCx3RUEyRkMiLCJzb3VyY2VzQ29udGVudCI6WyIvKipcbiAqIEBwYWNrYWdlRG9jdW1lbnRhdGlvblxuICogQG1vZHVsZSBDb21tb24tQXNzZXRBbW91bnRcbiAqL1xuXG5pbXBvcnQgeyBCdWZmZXIgfSBmcm9tIFwiYnVmZmVyL1wiXG5pbXBvcnQgQk4gZnJvbSBcImJuLmpzXCJcbmltcG9ydCB7IFN0YW5kYXJkVHJhbnNmZXJhYmxlT3V0cHV0IH0gZnJvbSBcIi4vb3V0cHV0XCJcbmltcG9ydCB7IFN0YW5kYXJkVHJhbnNmZXJhYmxlSW5wdXQgfSBmcm9tIFwiLi9pbnB1dFwiXG5pbXBvcnQgeyBJbnN1ZmZpY2llbnRGdW5kc0Vycm9yIH0gZnJvbSBcIi4uL3V0aWxzL2Vycm9yc1wiXG5cbi8qKlxuICogQ2xhc3MgZm9yIG1hbmFnaW5nIGFzc2V0IGFtb3VudHMgaW4gdGhlIFVUWE9TZXQgZmVlIGNhbGN1YXRpb25cbiAqL1xuZXhwb3J0IGNsYXNzIEFzc2V0QW1vdW50IHtcbiAgLy8gYXNzZXRJRCB0aGF0IGlzIGFtb3VudCBpcyBtYW5hZ2luZy5cbiAgcHJvdGVjdGVkIGFzc2V0SUQ6IEJ1ZmZlciA9IEJ1ZmZlci5hbGxvYygzMilcbiAgLy8gYW1vdW50IG9mIHRoaXMgYXNzZXQgdGhhdCBzaG91bGQgYmUgc2VudC5cbiAgcHJvdGVjdGVkIGFtb3VudDogQk4gPSBuZXcgQk4oMClcbiAgLy8gYnVybiBpcyB0aGUgYW1vdW50IG9mIHRoaXMgYXNzZXQgdGhhdCBzaG91bGQgYmUgYnVybmVkLlxuICBwcm90ZWN0ZWQgYnVybjogQk4gPSBuZXcgQk4oMClcblxuICAvLyBzcGVudCBpcyB0aGUgdG90YWwgYW1vdW50IG9mIHRoaXMgYXNzZXQgdGhhdCBoYXMgYmVlbiBjb25zdW1lZC5cbiAgcHJvdGVjdGVkIHNwZW50OiBCTiA9IG5ldyBCTigwKVxuICAvLyBzdGFrZWFibGVMb2NrU3BlbnQgaXMgdGhlIGFtb3VudCBvZiB0aGlzIGFzc2V0IHRoYXQgaGFzIGJlZW4gY29uc3VtZWQgdGhhdFxuICAvLyB3YXMgbG9ja2VkLlxuICBwcm90ZWN0ZWQgc3Rha2VhYmxlTG9ja1NwZW50OiBCTiA9IG5ldyBCTigwKVxuXG4gIC8vIGNoYW5nZSBpcyB0aGUgZXhjZXNzIGFtb3VudCBvZiB0aGlzIGFzc2V0IHRoYXQgd2FzIGNvbnN1bWVkIG92ZXIgdGhlIGFtb3VudFxuICAvLyByZXF1ZXN0ZWQgdG8gYmUgY29uc3VtZWQoYW1vdW50ICsgYnVybikuXG4gIHByb3RlY3RlZCBjaGFuZ2U6IEJOID0gbmV3IEJOKDApXG4gIC8vIHN0YWtlYWJsZUxvY2tDaGFuZ2UgaXMgYSBmbGFnIHRvIG1hcmsgaWYgdGhlIGlucHV0IHRoYXQgZ2VuZXJhdGVkIHRoZVxuICAvLyBjaGFuZ2Ugd2FzIGxvY2tlZC5cbiAgcHJvdGVjdGVkIHN0YWtlYWJsZUxvY2tDaGFuZ2U6IGJvb2xlYW4gPSBmYWxzZVxuXG4gIC8vIGZpbmlzaGVkIGlzIGEgY29udmVuaWVuY2UgZmxhZyB0byB0cmFjayBcInNwZW50ID49IGFtb3VudCArIGJ1cm5cIlxuICBwcm90ZWN0ZWQgZmluaXNoZWQ6IGJvb2xlYW4gPSBmYWxzZVxuXG4gIGdldEFzc2V0SUQgPSAoKTogQnVmZmVyID0+IHtcbiAgICByZXR1cm4gdGhpcy5hc3NldElEXG4gIH1cblxuICBnZXRBc3NldElEU3RyaW5nID0gKCk6IHN0cmluZyA9PiB7XG4gICAgcmV0dXJuIHRoaXMuYXNzZXRJRC50b1N0cmluZyhcImhleFwiKVxuICB9XG5cbiAgZ2V0QW1vdW50ID0gKCk6IEJOID0+IHtcbiAgICByZXR1cm4gdGhpcy5hbW91bnRcbiAgfVxuXG4gIGdldFNwZW50ID0gKCk6IEJOID0+IHtcbiAgICByZXR1cm4gdGhpcy5zcGVudFxuICB9XG5cbiAgZ2V0QnVybiA9ICgpOiBCTiA9PiB7XG4gICAgcmV0dXJuIHRoaXMuYnVyblxuICB9XG5cbiAgZ2V0Q2hhbmdlID0gKCk6IEJOID0+IHtcbiAgICByZXR1cm4gdGhpcy5jaGFuZ2VcbiAgfVxuXG4gIGdldFN0YWtlYWJsZUxvY2tTcGVudCA9ICgpOiBCTiA9PiB7XG4gICAgcmV0dXJuIHRoaXMuc3Rha2VhYmxlTG9ja1NwZW50XG4gIH1cblxuICBnZXRTdGFrZWFibGVMb2NrQ2hhbmdlID0gKCk6IGJvb2xlYW4gPT4ge1xuICAgIHJldHVybiB0aGlzLnN0YWtlYWJsZUxvY2tDaGFuZ2VcbiAgfVxuXG4gIGlzRmluaXNoZWQgPSAoKTogYm9vbGVhbiA9PiB7XG4gICAgcmV0dXJuIHRoaXMuZmluaXNoZWRcbiAgfVxuXG4gIC8vIHNwZW5kQW1vdW50IHNob3VsZCBvbmx5IGJlIGNhbGxlZCBpZiB0aGlzIGFzc2V0IGlzIHN0aWxsIGF3YWl0aW5nIG1vcmVcbiAgLy8gZnVuZHMgdG8gY29uc3VtZS5cbiAgc3BlbmRBbW91bnQgPSAoYW10OiBCTiwgc3Rha2VhYmxlTG9ja2VkOiBib29sZWFuID0gZmFsc2UpOiBib29sZWFuID0+IHtcbiAgICBpZiAodGhpcy5maW5pc2hlZCkge1xuICAgICAgLyogaXN0YW5idWwgaWdub3JlIG5leHQgKi9cbiAgICAgIHRocm93IG5ldyBJbnN1ZmZpY2llbnRGdW5kc0Vycm9yKFxuICAgICAgICBcIkVycm9yIC0gQXNzZXRBbW91bnQuc3BlbmRBbW91bnQ6IGF0dGVtcHRlZCB0byBzcGVuZCBcIiArIFwiZXhjZXNzIGZ1bmRzXCJcbiAgICAgIClcbiAgICB9XG4gICAgdGhpcy5zcGVudCA9IHRoaXMuc3BlbnQuYWRkKGFtdClcbiAgICBpZiAoc3Rha2VhYmxlTG9ja2VkKSB7XG4gICAgICB0aGlzLnN0YWtlYWJsZUxvY2tTcGVudCA9IHRoaXMuc3Rha2VhYmxlTG9ja1NwZW50LmFkZChhbXQpXG4gICAgfVxuXG4gICAgY29uc3QgdG90YWw6IEJOID0gdGhpcy5hbW91bnQuYWRkKHRoaXMuYnVybilcbiAgICBpZiAodGhpcy5zcGVudC5ndGUodG90YWwpKSB7XG4gICAgICB0aGlzLmNoYW5nZSA9IHRoaXMuc3BlbnQuc3ViKHRvdGFsKVxuICAgICAgaWYgKHN0YWtlYWJsZUxvY2tlZCkge1xuICAgICAgICB0aGlzLnN0YWtlYWJsZUxvY2tDaGFuZ2UgPSB0cnVlXG4gICAgICB9XG4gICAgICB0aGlzLmZpbmlzaGVkID0gdHJ1ZVxuICAgIH1cbiAgICByZXR1cm4gdGhpcy5maW5pc2hlZFxuICB9XG5cbiAgY29uc3RydWN0b3IoYXNzZXRJRDogQnVmZmVyLCBhbW91bnQ6IEJOLCBidXJuOiBCTikge1xuICAgIHRoaXMuYXNzZXRJRCA9IGFzc2V0SURcbiAgICB0aGlzLmFtb3VudCA9IHR5cGVvZiBhbW91bnQgPT09IFwidW5kZWZpbmVkXCIgPyBuZXcgQk4oMCkgOiBhbW91bnRcbiAgICB0aGlzLmJ1cm4gPSB0eXBlb2YgYnVybiA9PT0gXCJ1bmRlZmluZWRcIiA/IG5ldyBCTigwKSA6IGJ1cm5cbiAgICB0aGlzLnNwZW50ID0gbmV3IEJOKDApXG4gICAgdGhpcy5zdGFrZWFibGVMb2NrU3BlbnQgPSBuZXcgQk4oMClcbiAgICB0aGlzLnN0YWtlYWJsZUxvY2tDaGFuZ2UgPSBmYWxzZVxuICB9XG59XG5cbmV4cG9ydCBhYnN0cmFjdCBjbGFzcyBTdGFuZGFyZEFzc2V0QW1vdW50RGVzdGluYXRpb248XG4gIFRPIGV4dGVuZHMgU3RhbmRhcmRUcmFuc2ZlcmFibGVPdXRwdXQsXG4gIFRJIGV4dGVuZHMgU3RhbmRhcmRUcmFuc2ZlcmFibGVJbnB1dFxuPiB7XG4gIHByb3RlY3RlZCBhbW91bnRzOiBBc3NldEFtb3VudFtdID0gW11cbiAgcHJvdGVjdGVkIGRlc3RpbmF0aW9uczogQnVmZmVyW10gPSBbXVxuICBwcm90ZWN0ZWQgc2VuZGVyczogQnVmZmVyW10gPSBbXVxuICBwcm90ZWN0ZWQgY2hhbmdlQWRkcmVzc2VzOiBCdWZmZXJbXSA9IFtdXG4gIHByb3RlY3RlZCBhbW91bnRrZXk6IG9iamVjdCA9IHt9XG4gIHByb3RlY3RlZCBpbnB1dHM6IFRJW10gPSBbXVxuICBwcm90ZWN0ZWQgb3V0cHV0czogVE9bXSA9IFtdXG4gIHByb3RlY3RlZCBjaGFuZ2U6IFRPW10gPSBbXVxuXG4gIC8vIFRPRE86IHNob3VsZCB0aGlzIGZ1bmN0aW9uIGFsbG93IGZvciByZXBlYXRlZCBjYWxscyB3aXRoIHRoZSBzYW1lXG4gIC8vICAgICAgIGFzc2V0SUQ/XG4gIGFkZEFzc2V0QW1vdW50ID0gKGFzc2V0SUQ6IEJ1ZmZlciwgYW1vdW50OiBCTiwgYnVybjogQk4pID0+IHtcbiAgICBsZXQgYWE6IEFzc2V0QW1vdW50ID0gbmV3IEFzc2V0QW1vdW50KGFzc2V0SUQsIGFtb3VudCwgYnVybilcbiAgICB0aGlzLmFtb3VudHMucHVzaChhYSlcbiAgICB0aGlzLmFtb3VudGtleVthYS5nZXRBc3NldElEU3RyaW5nKCldID0gYWFcbiAgfVxuXG4gIGFkZElucHV0ID0gKGlucHV0OiBUSSkgPT4ge1xuICAgIHRoaXMuaW5wdXRzLnB1c2goaW5wdXQpXG4gIH1cblxuICBhZGRPdXRwdXQgPSAob3V0cHV0OiBUTykgPT4ge1xuICAgIHRoaXMub3V0cHV0cy5wdXNoKG91dHB1dClcbiAgfVxuXG4gIGFkZENoYW5nZSA9IChvdXRwdXQ6IFRPKSA9PiB7XG4gICAgdGhpcy5jaGFuZ2UucHVzaChvdXRwdXQpXG4gIH1cblxuICBnZXRBbW91bnRzID0gKCk6IEFzc2V0QW1vdW50W10gPT4ge1xuICAgIHJldHVybiB0aGlzLmFtb3VudHNcbiAgfVxuXG4gIGdldERlc3RpbmF0aW9ucyA9ICgpOiBCdWZmZXJbXSA9PiB7XG4gICAgcmV0dXJuIHRoaXMuZGVzdGluYXRpb25zXG4gIH1cblxuICBnZXRTZW5kZXJzID0gKCk6IEJ1ZmZlcltdID0+IHtcbiAgICByZXR1cm4gdGhpcy5zZW5kZXJzXG4gIH1cblxuICBnZXRDaGFuZ2VBZGRyZXNzZXMgPSAoKTogQnVmZmVyW10gPT4ge1xuICAgIHJldHVybiB0aGlzLmNoYW5nZUFkZHJlc3Nlc1xuICB9XG5cbiAgZ2V0QXNzZXRBbW91bnQgPSAoYXNzZXRIZXhTdHI6IHN0cmluZyk6IEFzc2V0QW1vdW50ID0+IHtcbiAgICByZXR1cm4gdGhpcy5hbW91bnRrZXlbYCR7YXNzZXRIZXhTdHJ9YF1cbiAgfVxuXG4gIGFzc2V0RXhpc3RzID0gKGFzc2V0SGV4U3RyOiBzdHJpbmcpOiBib29sZWFuID0+IHtcbiAgICByZXR1cm4gYXNzZXRIZXhTdHIgaW4gdGhpcy5hbW91bnRrZXlcbiAgfVxuXG4gIGdldElucHV0cyA9ICgpOiBUSVtdID0+IHtcbiAgICByZXR1cm4gdGhpcy5pbnB1dHNcbiAgfVxuXG4gIGdldE91dHB1dHMgPSAoKTogVE9bXSA9PiB7XG4gICAgcmV0dXJuIHRoaXMub3V0cHV0c1xuICB9XG5cbiAgZ2V0Q2hhbmdlT3V0cHV0cyA9ICgpOiBUT1tdID0+IHtcbiAgICByZXR1cm4gdGhpcy5jaGFuZ2VcbiAgfVxuXG4gIGdldEFsbE91dHB1dHMgPSAoKTogVE9bXSA9PiB7XG4gICAgcmV0dXJuIHRoaXMub3V0cHV0cy5jb25jYXQodGhpcy5jaGFuZ2UpXG4gIH1cblxuICBjYW5Db21wbGV0ZSA9ICgpOiBib29sZWFuID0+IHtcbiAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgdGhpcy5hbW91bnRzLmxlbmd0aDsgaSsrKSB7XG4gICAgICBpZiAoIXRoaXMuYW1vdW50c1tgJHtpfWBdLmlzRmluaXNoZWQoKSkge1xuICAgICAgICByZXR1cm4gZmFsc2VcbiAgICAgIH1cbiAgICB9XG4gICAgcmV0dXJuIHRydWVcbiAgfVxuXG4gIGNvbnN0cnVjdG9yKFxuICAgIGRlc3RpbmF0aW9uczogQnVmZmVyW10sXG4gICAgc2VuZGVyczogQnVmZmVyW10sXG4gICAgY2hhbmdlQWRkcmVzc2VzOiBCdWZmZXJbXVxuICApIHtcbiAgICB0aGlzLmRlc3RpbmF0aW9ucyA9IGRlc3RpbmF0aW9uc1xuICAgIHRoaXMuY2hhbmdlQWRkcmVzc2VzID0gY2hhbmdlQWRkcmVzc2VzXG4gICAgdGhpcy5zZW5kZXJzID0gc2VuZGVyc1xuICB9XG59XG4iXX0=