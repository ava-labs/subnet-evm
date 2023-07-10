"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
class PubSub {
    newSet() {
        return JSON.stringify({ newSet: {} });
    }
    newBloom(maxElements = 1000, collisionProb = 0.01) {
        return JSON.stringify({
            newBloom: {
                maxElements: maxElements,
                collisionProb: collisionProb
            }
        });
    }
    addAddresses(addresses) {
        return JSON.stringify({
            addAddresses: {
                addresses: addresses
            }
        });
    }
}
exports.default = PubSub;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicHVic3ViLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL3V0aWxzL3B1YnN1Yi50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOztBQUFBLE1BQXFCLE1BQU07SUFDekIsTUFBTTtRQUNKLE9BQU8sSUFBSSxDQUFDLFNBQVMsQ0FBQyxFQUFFLE1BQU0sRUFBRSxFQUFFLEVBQUUsQ0FBQyxDQUFBO0lBQ3ZDLENBQUM7SUFDRCxRQUFRLENBQUMsY0FBc0IsSUFBSSxFQUFFLGdCQUF3QixJQUFJO1FBQy9ELE9BQU8sSUFBSSxDQUFDLFNBQVMsQ0FBQztZQUNwQixRQUFRLEVBQUU7Z0JBQ1IsV0FBVyxFQUFFLFdBQVc7Z0JBQ3hCLGFBQWEsRUFBRSxhQUFhO2FBQzdCO1NBQ0YsQ0FBQyxDQUFBO0lBQ0osQ0FBQztJQUNELFlBQVksQ0FBQyxTQUFtQjtRQUM5QixPQUFPLElBQUksQ0FBQyxTQUFTLENBQUM7WUFDcEIsWUFBWSxFQUFFO2dCQUNaLFNBQVMsRUFBRSxTQUFTO2FBQ3JCO1NBQ0YsQ0FBQyxDQUFBO0lBQ0osQ0FBQztDQUNGO0FBbkJELHlCQW1CQyIsInNvdXJjZXNDb250ZW50IjpbImV4cG9ydCBkZWZhdWx0IGNsYXNzIFB1YlN1YiB7XG4gIG5ld1NldCgpIHtcbiAgICByZXR1cm4gSlNPTi5zdHJpbmdpZnkoeyBuZXdTZXQ6IHt9IH0pXG4gIH1cbiAgbmV3Qmxvb20obWF4RWxlbWVudHM6IG51bWJlciA9IDEwMDAsIGNvbGxpc2lvblByb2I6IG51bWJlciA9IDAuMDEpIHtcbiAgICByZXR1cm4gSlNPTi5zdHJpbmdpZnkoe1xuICAgICAgbmV3Qmxvb206IHtcbiAgICAgICAgbWF4RWxlbWVudHM6IG1heEVsZW1lbnRzLFxuICAgICAgICBjb2xsaXNpb25Qcm9iOiBjb2xsaXNpb25Qcm9iXG4gICAgICB9XG4gICAgfSlcbiAgfVxuICBhZGRBZGRyZXNzZXMoYWRkcmVzc2VzOiBzdHJpbmdbXSkge1xuICAgIHJldHVybiBKU09OLnN0cmluZ2lmeSh7XG4gICAgICBhZGRBZGRyZXNzZXM6IHtcbiAgICAgICAgYWRkcmVzc2VzOiBhZGRyZXNzZXNcbiAgICAgIH1cbiAgICB9KVxuICB9XG59XG4iXX0=