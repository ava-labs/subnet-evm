"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.PubSub = void 0;
class PubSub {
    newSet() {
        return { newSet: {} };
    }
    newBloom(maxElements = 1000, collisionProb = 0.01) {
        return {
            newBloom: {
                maxElements: maxElements,
                collisionProb: collisionProb
            }
        };
    }
    addAddresses(addresses) {
        return {
            addAddresses: {
                addresses: addresses
            }
        };
    }
}
exports.PubSub = PubSub;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicHVic3ViLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2FwaXMvcHVic3ViL3B1YnN1Yi50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFBQSxNQUFhLE1BQU07SUFDakIsTUFBTTtRQUNKLE9BQU8sRUFBQyxNQUFNLEVBQUMsRUFBRSxFQUFDLENBQUM7SUFDckIsQ0FBQztJQUNELFFBQVEsQ0FBQyxjQUFzQixJQUFJLEVBQUUsZ0JBQXdCLElBQUk7UUFDL0QsT0FBTztZQUNMLFFBQVEsRUFBRTtnQkFDUixXQUFXLEVBQUUsV0FBVztnQkFDeEIsYUFBYSxFQUFFLGFBQWE7YUFDN0I7U0FDRixDQUFBO0lBQ0gsQ0FBQztJQUNELFlBQVksQ0FBQyxTQUFtQjtRQUM5QixPQUFPO1lBQ0wsWUFBWSxFQUFFO2dCQUNaLFNBQVMsRUFBRSxTQUFTO2FBQ3JCO1NBQ0YsQ0FBQztJQUNKLENBQUM7Q0FDRjtBQW5CRCx3QkFtQkMiLCJzb3VyY2VzQ29udGVudCI6WyJleHBvcnQgY2xhc3MgUHViU3ViIHtcbiAgbmV3U2V0KCkge1xuICAgIHJldHVybiB7bmV3U2V0Ont9fTtcbiAgfVxuICBuZXdCbG9vbShtYXhFbGVtZW50czogbnVtYmVyID0gMTAwMCwgY29sbGlzaW9uUHJvYjogbnVtYmVyID0gMC4wMSkge1xuICAgIHJldHVybiB7XG4gICAgICBuZXdCbG9vbToge1xuICAgICAgICBtYXhFbGVtZW50czogbWF4RWxlbWVudHMsXG4gICAgICAgIGNvbGxpc2lvblByb2I6IGNvbGxpc2lvblByb2JcbiAgICAgIH1cbiAgICB9XG4gIH1cbiAgYWRkQWRkcmVzc2VzKGFkZHJlc3Nlczogc3RyaW5nW10pIHtcbiAgICByZXR1cm4ge1xuICAgICAgYWRkQWRkcmVzc2VzOiB7XG4gICAgICAgIGFkZHJlc3NlczogYWRkcmVzc2VzXG4gICAgICB9XG4gICAgfTtcbiAgfVxufSJdfQ==