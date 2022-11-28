import { expect } from "chai";
import { ethers } from "hardhat"
import {
    BigNumber,
} from "ethers"

// make sure this is always an admin for minter precompile
const adminAddress: string = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"

describe.only('Order Book', function () {
    let orderBook, alice, bob, order, domain, orderType, signature

    before(async function () {
        // const owner = await ethers.getSigner(adminAddress);
        const OrderBook = await ethers.getContractFactory('OrderBook')
        orderBook = await OrderBook.deploy('Hubble', '2.0')
        const signers = await ethers.getSigners()
            ; ([, alice, bob] = signers)
    })

    it('verify signer', async function () {
        order = {
            trader: alice.address,
            baseAssetQuantity: ethers.utils.parseEther('-5'),
            price: ethers.utils.parseUnits('15', 6),
            salt: 1
        }

        domain = {
            name: 'Hubble',
            version: '2.0',
            chainId: (await ethers.provider.getNetwork()).chainId,
            verifyingContract: orderBook.address
        }

        orderType = {
            Order: [
                // field ordering must be the same as LIMIT_ORDER_TYPEHASH
                { name: "trader", type: "address" },
                { name: "baseAssetQuantity", type: "int256" },
                { name: "price", type: "uint256" },
                { name: "salt", type: "uint256" },
            ]
        }

        signature = await alice._signTypedData(domain, orderType, order)
        const signer = (await orderBook.verifySigner(order, signature))[0]
        expect(signer).to.eq(alice.address)
    })

    it('place an order', async function () {
        const tx = await orderBook.placeOrder(order, signature)
        await expect(tx).to.emit(orderBook, "OrderPlaced").withArgs(
            alice.address,
            order.baseAssetQuantity,
            order.price,
            adminAddress
        )

        let orderHash = await orderBook.getOrderHash(order)
        let status
        status = await orderBook.ordersStatus(orderHash)
        console.log({ status });

        expect(await orderBook.ordersStatus(orderHash)).to.eq(1) // Filled; because evm is fulfilling all orders right now
    })

    it.skip('execute matched orders', async function () {
        const order2 = JSON.parse(JSON.stringify(order))
        order2.baseAssetQuantity = BigNumber.from(order2.baseAssetQuantity).mul(-1)
        order2.trader = bob.address
        const signature2 = await bob._signTypedData(domain, orderType, order2)
        await orderBook.placeOrder(order2, signature2)
        await delay(1000)

        await orderBook.executeMatchedOrders(order, signature, order2, signature2, { gasLimit: 1e6 })
        await delay(1500)

        let position = await orderBook.positions(alice.address)
        expect(position.size).to.eq(order.baseAssetQuantity)
        expect(position.openNotional).to.eq(order.price.mul(order.baseAssetQuantity).abs())

        position = await orderBook.positions(bob.address)
        expect(position.size).to.eq(order2.baseAssetQuantity)
        expect(position.openNotional).to.eq(order2.baseAssetQuantity.mul(order2.price).abs())

        let orderHash = await orderBook.getOrderHash(order)
        expect(await orderBook.ordersStatus(orderHash)).to.eq(1) // Filled
        orderHash = await orderBook.getOrderHash(order2)
        expect(await orderBook.ordersStatus(orderHash)).to.eq(1) // Filled
    })
})

function delay(ms: number) {
    return new Promise(resolve => setTimeout(resolve, ms));
}
