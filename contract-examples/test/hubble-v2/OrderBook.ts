import { expect } from "chai";
import { ethers } from "hardhat"
import { BigNumber } from "ethers"

// make sure this is always an admin for minter precompile
const adminAddress: string = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const GENESIS_ORDERBOOK_ADDRESS = '0x0300000000000000000000000000000000000069'

describe.only('Order Book', function () {
    let orderBook, alice, bob, order, domain, orderType, signature

    before(async function () {
        const signers = await ethers.getSigners()
        ;([, alice, bob] = signers)

        // 1. set proxyAdmin
        const genesisTUP = await ethers.getContractAt('GenesisTUP', GENESIS_ORDERBOOK_ADDRESS)
        let _admin = await ethers.provider.getStorageAt(GENESIS_ORDERBOOK_ADDRESS, '0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103')
        // console.log({ _admin })
        let proxyAdmin
        if (_admin == '0x' + '0'.repeat(64)) { // because we don't run a fresh subnet everytime
            const ProxyAdmin = await ethers.getContractFactory('ProxyAdmin')
            proxyAdmin = await ProxyAdmin.deploy()
            await genesisTUP.init(proxyAdmin.address)
            console.log('genesisTUP.init done...')
            await delay(2000)
        } else {
            proxyAdmin = await ethers.getContractAt('ProxyAdmin', '0x' + _admin.slice(26))
        }
        // _admin = await ethers.provider.getStorageAt(GENESIS_ORDERBOOK_ADDRESS, '0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103')
        // console.log({ _admin })

        // 2. set implementation
        const OrderBook = await ethers.getContractFactory('OrderBook')
        const orderBookImpl = await OrderBook.deploy()

        orderBook = await ethers.getContractAt('OrderBook', GENESIS_ORDERBOOK_ADDRESS)
        let _impl = await ethers.provider.getStorageAt(GENESIS_ORDERBOOK_ADDRESS, '0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc')

        let isInitialized = false
        if (_impl != '0x' + '0'.repeat(64)) {
            isInitialized = await orderBook.isInitialized()
        }

        if (isInitialized) {
            await proxyAdmin.upgrade(GENESIS_ORDERBOOK_ADDRESS, orderBookImpl.address)
        } else {
            await proxyAdmin.upgradeAndCall(
                GENESIS_ORDERBOOK_ADDRESS,
                orderBookImpl.address,
                orderBookImpl.interface.encodeFunctionData('initialize', ['Hubble', '2.0'])
            )
        }
        await delay(2000)

        _impl = await ethers.provider.getStorageAt(GENESIS_ORDERBOOK_ADDRESS, '0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc')
        // console.log({ _impl })
        expect(ethers.utils.getAddress('0x' + _impl.slice(26))).to.eq(orderBookImpl.address)
    })

    it.only('verify signer', async function() {
        order = {
            trader: alice.address,
            baseAssetQuantity: ethers.utils.parseEther('-5'),
            price: ethers.utils.parseUnits('15', 6),
            salt: Date.now()
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

    it('place an order', async function() {
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

    it('execute matched orders', async function() {
        const order2 = {
            trader: bob.address,
            baseAssetQuantity: BigNumber.from(order.baseAssetQuantity).mul(-1),
            price: ethers.utils.parseUnits('15', 6),
            salt: Date.now()
        }

        const signature2 = await bob._signTypedData(domain, orderType, order2)
        await orderBook.placeOrder(order2, signature2)
        await delay(1000)

        await orderBook.executeMatchedOrders(order, signature, order2, signature2, {gasLimit: 1e6})
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
