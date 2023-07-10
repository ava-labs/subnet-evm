"use strict";
/**
 * @packageDocumentation
 * @module Utils-PersistanceOptions
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.PersistanceOptions = void 0;
/**
 * A class for defining the persistance behavior of this an API call.
 *
 */
class PersistanceOptions {
    /**
     *
     * @param name The namespace of the database the data
     * @param overwrite True if the data should be completey overwritten
     * @param MergeRule The type of process used to merge with existing data: "intersection", "differenceSelf", "differenceNew", "symDifference", "union", "unionMinusNew", "unionMinusSelf"
     *
     * @remarks
     * The merge rules are as follows:
     *   * "intersection" - the intersection of the set
     *   * "differenceSelf" - the difference between the existing data and new set
     *   * "differenceNew" - the difference between the new data and the existing set
     *   * "symDifference" - the union of the differences between both sets of data
     *   * "union" - the unique set of all elements contained in both sets
     *   * "unionMinusNew" - the unique set of all elements contained in both sets, excluding values only found in the new set
     *   * "unionMinusSelf" - the unique set of all elements contained in both sets, excluding values only found in the existing set
     */
    constructor(name, overwrite = false, mergeRule) {
        this.name = undefined;
        this.overwrite = false;
        this.mergeRule = "union";
        /**
         * Returns the namespace of the instance
         */
        this.getName = () => this.name;
        /**
         * Returns the overwrite rule of the instance
         */
        this.getOverwrite = () => this.overwrite;
        /**
         * Returns the [[MergeRule]] of the instance
         */
        this.getMergeRule = () => this.mergeRule;
        this.name = name;
        this.overwrite = overwrite;
        this.mergeRule = mergeRule;
    }
}
exports.PersistanceOptions = PersistanceOptions;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicGVyc2lzdGVuY2VvcHRpb25zLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL3V0aWxzL3BlcnNpc3RlbmNlb3B0aW9ucy50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiO0FBQUE7OztHQUdHOzs7QUFHSDs7O0dBR0c7QUFDSCxNQUFhLGtCQUFrQjtJQXNCN0I7Ozs7Ozs7Ozs7Ozs7OztPQWVHO0lBQ0gsWUFBWSxJQUFZLEVBQUUsWUFBcUIsS0FBSyxFQUFFLFNBQW9CO1FBckNoRSxTQUFJLEdBQVcsU0FBUyxDQUFBO1FBRXhCLGNBQVMsR0FBWSxLQUFLLENBQUE7UUFFMUIsY0FBUyxHQUFjLE9BQU8sQ0FBQTtRQUV4Qzs7V0FFRztRQUNILFlBQU8sR0FBRyxHQUFXLEVBQUUsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFBO1FBRWpDOztXQUVHO1FBQ0gsaUJBQVksR0FBRyxHQUFZLEVBQUUsQ0FBQyxJQUFJLENBQUMsU0FBUyxDQUFBO1FBRTVDOztXQUVHO1FBQ0gsaUJBQVksR0FBRyxHQUFjLEVBQUUsQ0FBQyxJQUFJLENBQUMsU0FBUyxDQUFBO1FBbUI1QyxJQUFJLENBQUMsSUFBSSxHQUFHLElBQUksQ0FBQTtRQUNoQixJQUFJLENBQUMsU0FBUyxHQUFHLFNBQVMsQ0FBQTtRQUMxQixJQUFJLENBQUMsU0FBUyxHQUFHLFNBQVMsQ0FBQTtJQUM1QixDQUFDO0NBQ0Y7QUEzQ0QsZ0RBMkNDIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgVXRpbHMtUGVyc2lzdGFuY2VPcHRpb25zXG4gKi9cblxuaW1wb3J0IHsgTWVyZ2VSdWxlIH0gZnJvbSBcIi4vY29uc3RhbnRzXCJcbi8qKlxuICogQSBjbGFzcyBmb3IgZGVmaW5pbmcgdGhlIHBlcnNpc3RhbmNlIGJlaGF2aW9yIG9mIHRoaXMgYW4gQVBJIGNhbGwuXG4gKlxuICovXG5leHBvcnQgY2xhc3MgUGVyc2lzdGFuY2VPcHRpb25zIHtcbiAgcHJvdGVjdGVkIG5hbWU6IHN0cmluZyA9IHVuZGVmaW5lZFxuXG4gIHByb3RlY3RlZCBvdmVyd3JpdGU6IGJvb2xlYW4gPSBmYWxzZVxuXG4gIHByb3RlY3RlZCBtZXJnZVJ1bGU6IE1lcmdlUnVsZSA9IFwidW5pb25cIlxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBuYW1lc3BhY2Ugb2YgdGhlIGluc3RhbmNlXG4gICAqL1xuICBnZXROYW1lID0gKCk6IHN0cmluZyA9PiB0aGlzLm5hbWVcblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgb3ZlcndyaXRlIHJ1bGUgb2YgdGhlIGluc3RhbmNlXG4gICAqL1xuICBnZXRPdmVyd3JpdGUgPSAoKTogYm9vbGVhbiA9PiB0aGlzLm92ZXJ3cml0ZVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBbW01lcmdlUnVsZV1dIG9mIHRoZSBpbnN0YW5jZVxuICAgKi9cbiAgZ2V0TWVyZ2VSdWxlID0gKCk6IE1lcmdlUnVsZSA9PiB0aGlzLm1lcmdlUnVsZVxuXG4gIC8qKlxuICAgKlxuICAgKiBAcGFyYW0gbmFtZSBUaGUgbmFtZXNwYWNlIG9mIHRoZSBkYXRhYmFzZSB0aGUgZGF0YVxuICAgKiBAcGFyYW0gb3ZlcndyaXRlIFRydWUgaWYgdGhlIGRhdGEgc2hvdWxkIGJlIGNvbXBsZXRleSBvdmVyd3JpdHRlblxuICAgKiBAcGFyYW0gTWVyZ2VSdWxlIFRoZSB0eXBlIG9mIHByb2Nlc3MgdXNlZCB0byBtZXJnZSB3aXRoIGV4aXN0aW5nIGRhdGE6IFwiaW50ZXJzZWN0aW9uXCIsIFwiZGlmZmVyZW5jZVNlbGZcIiwgXCJkaWZmZXJlbmNlTmV3XCIsIFwic3ltRGlmZmVyZW5jZVwiLCBcInVuaW9uXCIsIFwidW5pb25NaW51c05ld1wiLCBcInVuaW9uTWludXNTZWxmXCJcbiAgICpcbiAgICogQHJlbWFya3NcbiAgICogVGhlIG1lcmdlIHJ1bGVzIGFyZSBhcyBmb2xsb3dzOlxuICAgKiAgICogXCJpbnRlcnNlY3Rpb25cIiAtIHRoZSBpbnRlcnNlY3Rpb24gb2YgdGhlIHNldFxuICAgKiAgICogXCJkaWZmZXJlbmNlU2VsZlwiIC0gdGhlIGRpZmZlcmVuY2UgYmV0d2VlbiB0aGUgZXhpc3RpbmcgZGF0YSBhbmQgbmV3IHNldFxuICAgKiAgICogXCJkaWZmZXJlbmNlTmV3XCIgLSB0aGUgZGlmZmVyZW5jZSBiZXR3ZWVuIHRoZSBuZXcgZGF0YSBhbmQgdGhlIGV4aXN0aW5nIHNldFxuICAgKiAgICogXCJzeW1EaWZmZXJlbmNlXCIgLSB0aGUgdW5pb24gb2YgdGhlIGRpZmZlcmVuY2VzIGJldHdlZW4gYm90aCBzZXRzIG9mIGRhdGFcbiAgICogICAqIFwidW5pb25cIiAtIHRoZSB1bmlxdWUgc2V0IG9mIGFsbCBlbGVtZW50cyBjb250YWluZWQgaW4gYm90aCBzZXRzXG4gICAqICAgKiBcInVuaW9uTWludXNOZXdcIiAtIHRoZSB1bmlxdWUgc2V0IG9mIGFsbCBlbGVtZW50cyBjb250YWluZWQgaW4gYm90aCBzZXRzLCBleGNsdWRpbmcgdmFsdWVzIG9ubHkgZm91bmQgaW4gdGhlIG5ldyBzZXRcbiAgICogICAqIFwidW5pb25NaW51c1NlbGZcIiAtIHRoZSB1bmlxdWUgc2V0IG9mIGFsbCBlbGVtZW50cyBjb250YWluZWQgaW4gYm90aCBzZXRzLCBleGNsdWRpbmcgdmFsdWVzIG9ubHkgZm91bmQgaW4gdGhlIGV4aXN0aW5nIHNldFxuICAgKi9cbiAgY29uc3RydWN0b3IobmFtZTogc3RyaW5nLCBvdmVyd3JpdGU6IGJvb2xlYW4gPSBmYWxzZSwgbWVyZ2VSdWxlOiBNZXJnZVJ1bGUpIHtcbiAgICB0aGlzLm5hbWUgPSBuYW1lXG4gICAgdGhpcy5vdmVyd3JpdGUgPSBvdmVyd3JpdGVcbiAgICB0aGlzLm1lcmdlUnVsZSA9IG1lcmdlUnVsZVxuICB9XG59XG4iXX0=