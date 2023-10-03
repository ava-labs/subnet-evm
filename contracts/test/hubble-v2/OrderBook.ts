import { expect } from "chai";
import { ethers } from "hardhat"
import { BigNumber } from "ethers"
// import * as _ from "lodash";

// make sure this is always an admin for minter precompile
const adminAddress: string = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const GENESIS_ORDERBOOK_ADDRESS = '0x03000000000000000000000000000000000000b0'

describe.only('Order Book', function () {
    let orderBook, alice, bob, longOrder, shortOrder, domain, orderType, signature

    before(async function () {
        const signers = await ethers.getSigners()
        ;([, alice, bob] = signers)

        console.log({alice: alice.address, bob: bob.address})
        // 1. set proxyAdmin
        // const genesisTUP = await ethers.getContractAt('GenesisTUP', GENESIS_ORDERBOOK_ADDRESS)
        // let _admin = await ethers.provider.getStorageAt(GENESIS_ORDERBOOK_ADDRESS, '0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103')
        // // console.log({ _admin })
        // let proxyAdmin
        // if (_admin == '0x' + '0'.repeat(64)) { // because we don't run a fresh subnet everytime
        //     const ProxyAdmin = await ethers.getContractFactory('ProxyAdmin')
        //     proxyAdmin = await ProxyAdmin.deploy()
        //     await genesisTUP.init(proxyAdmin.address)
        //     console.log('genesisTUP.init done...')
        //     await delay(2000)
        // } else {
        //     proxyAdmin = await ethers.getContractAt('ProxyAdmin', '0x' + _admin.slice(26))
        // }
        // // _admin = await ethers.provider.getStorageAt(GENESIS_ORDERBOOK_ADDRESS, '0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103')
        // // console.log({ _admin })

        // // 2. set implementation
        // const OrderBook = await ethers.getContractFactory('OrderBook')
        // const orderBookImpl = await OrderBook.deploy()

        // await delay(2000)
        // orderBook = await ethers.getContractAt('OrderBook', GENESIS_ORDERBOOK_ADDRESS)
        // let _impl = await ethers.provider.getStorageAt(GENESIS_ORDERBOOK_ADDRESS, '0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc')

        // if (_impl != '0x' + '0'.repeat(64)) {
        //     await proxyAdmin.upgrade(GENESIS_ORDERBOOK_ADDRESS, orderBookImpl.address)
        // } else {
        //     await proxyAdmin.upgradeAndCall(
        //         GENESIS_ORDERBOOK_ADDRESS,
        //         orderBookImpl.address,
        //         orderBookImpl.interface.encodeFunctionData('initialize', ['Hubble', '2.0'])
        //     )
        // }
        // await delay(2000)

        // _impl = await ethers.provider.getStorageAt(GENESIS_ORDERBOOK_ADDRESS, '0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc')
        // // console.log({ _impl })
        // expect(ethers.utils.getAddress('0x' + _impl.slice(26))).to.eq(orderBookImpl.address)
    })

    it('verify signer', async function() {

        orderBook = await ethers.getContractAt('OrderBook', GENESIS_ORDERBOOK_ADDRESS)
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
        shortOrder = {
              trader: alice.address,
              baseAssetQuantity: ethers.utils.parseEther('-5'),
              price: ethers.utils.parseUnits('15', 6),
              salt: Date.now()
        }

        signature = await alice._signTypedData(domain, orderType, shortOrder)
        const signer = (await orderBook.verifySigner(shortOrder, signature))[0]
        expect(signer).to.eq(alice.address)
    })

    it('place an order', async function() {
        const tx = await orderBook.placeOrder(shortOrder, signature)
        await expect(tx).to.emit(orderBook, "OrderPlaced").withArgs(
            shortOrder.trader,
            shortOrder,
            signature
        )
    })

    it('matches orders with same price and opposite base asset quantity', async function() {
      // long order with same price and baseAssetQuantity
        longOrder = {
            trader: bob.address,
            baseAssetQuantity: ethers.utils.parseEther('5'),
            price: ethers.utils.parseUnits('15', 6),
            salt: Date.now()
        }
        let signature = await bob._signTypedData(domain, orderType, longOrder)
        const tx = await orderBook.placeOrder(longOrder, signature)

        await delay(6000)

        const filter = orderBook.filters
        let events = await orderBook.queryFilter(filter)
        console.log({events});

        let matchedOrderEvent = events[events.length -1]
        // expect(matchedOrderEvent.event).to.eq('OrderMatched')
    })

    it.skip('make lots of orders', async function() {
        const signers = await ethers.getSigners()

      // long order with same price and baseAssetQuantity
        longOrder = {
            trader: _.sample(signers).address,
            baseAssetQuantity: ethers.utils.parseEther('5'),
            price: ethers.utils.parseUnits('15', 6),
            salt: Date.now()
        }
        shortOrder = {
            trader: _.sample(signers).address,
            baseAssetQuantity: ethers.utils.parseEther('-5'),
            price: ethers.utils.parseUnits('15', 6),
            salt: Date.now()
        }
        let signature = await bob._signTypedData(domain, orderType, longOrder)
        const tx = await orderBook.placeOrder(longOrder, signature)

        signature = await bob._signTypedData(domain, orderType, shortOrder)
        const tx = await orderBook.placeOrder(longOrder, signature)

        await delay(6000)

        const filter = orderBook.filters
        let events = await orderBook.queryFilter(filter)
        let matchedOrderEvent = events[events.length -1]
        // expect(matchedOrderEvent.event).to.eq('OrderMatched')
    })

    it.skip('matches multiple long orders with same price and opposite base asset quantity with short orders', async function() {
        longOrder.salt = Date.now()
        signature = await bob._signTypedData(domain, orderType, longOrder)
        const longOrderTx1 = await orderBook.placeOrder(longOrder, signature)

        longOrder.salt = Date.now()
        signature = await bob._signTypedData(domain, orderType, longOrder)
        const longOrderTx2 = await orderBook.placeOrder(longOrder, signature)

        shortOrder.salt = Date.now()
        signature = await alice._signTypedData(domain, orderType, shortOrder)
        let shortOrderTx1 = await orderBook.placeOrder(shortOrder, signature)

        shortOrder.salt = Date.now()
        signature = await alice._signTypedData(domain, orderType, shortOrder)
        let shortOrderTx2 = await orderBook.placeOrder(shortOrder, signature)

        // waiting for next buildblock call
        await delay(6000)
        const filter = orderBook.filters
        let events = await orderBook.queryFilter(filter)

        expect(events[events.length - 1].event).to.eq('OrderMatched')
        expect(events[events.length - 2].event).to.eq('OrderMatched')
        expect(events[events.length - 3].event).to.eq('OrderPlaced')
        expect(events[events.length - 4].event).to.eq('OrderPlaced')
    })
})

function delay(ms: number) {
    return new Promise(resolve => setTimeout(resolve, ms));
}
