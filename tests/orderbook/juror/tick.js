const { expect } = require("chai");
const { BigNumber } = require("ethers");
const utils = require("../utils")

const {
    addMargin,
    alice,
    bnToFloat,
    cancelV2Orders,
    getAMMContract,
    getOrderV2,
    getRandomSalt,
    limitOrderBook,
    multiplyPrice,
    multiplySize,
    placeOrderFromLimitOrderV2,
    placeV2Orders,
    removeAllAvailableMargin,
} = utils

describe("Testing Tick methods", async function() {
    market = BigNumber.from(0)
    initialMargin = multiplyPrice(500000)
    let amm

    this.beforeAll(async function() {
        amm = await getAMMContract(0)
    })
    this.afterEach(async function() {
        // get all OrderAccepted events
        let filter = limitOrderBook.filters.OrderAccepted(alice.address)
        let orderAcceptedEvents = await limitOrderBook.queryFilter(filter)
        // console.log(orderAcceptedEvents)

        // get all OrderCancelAccepted events
        filter = limitOrderBook.filters.OrderCancelAccepted(alice.address)
        let orderCancelAccepted = await limitOrderBook.queryFilter(filter)
        // console.log(orderCancelAccepted)
        const openOrders = orderAcceptedEvents.filter(e => {
            return orderCancelAccepted.filter(e2 => e2.args.orderHash == e.args.orderHash).length == 0
        }).map(e => e.args.order)

        console.log('openOrders', openOrders.length)
        if (openOrders.length) {
            const { txReceipt } = await cancelV2Orders(openOrders, alice)
            // const orderRejected = txReceipt.events.filter(l => l.event == 'OrderCancelRejected')
            // console.log(orderRejected.map(l => l.args))
        }
        await removeAllAvailableMargin(alice)
    })

    // these 2 tests when run together have a problem that they dont account for live matching
    it("bids", async function() {
        expect((await amm.bidsHead()).toNumber()).to.equal(0)
        let orderData = generateRandomArray(15)
        orderData = orderData.map(a => {
            return { price: multiplyPrice(a.price), size: multiplySize(a.size) }
        })


        const orders = []
        let requiredMargin = BigNumber.from(0)
        for (let i = 0; i < orderData.length; i++) {
            let longOrder = getOrderV2(market, alice.address, orderData[i].size, orderData[i].price, getRandomSalt())
            requiredMargin = requiredMargin.add(await utils.getRequiredMarginForLongOrder(longOrder))
            orders.push(longOrder)
        }
        await addMargin(alice, requiredMargin)

        const { txReceipt } = await placeV2Orders(orders, alice)
        txReceipt.events.forEach(e => prettyPrintEvents(e))

        // sort orderData based on descending price
        orderData = orderData
        .reduce((accumulator, order) => {
            // Find an existing order in the accumulator with the same price
            const existingOrder = accumulator.find(item => item.price.eq(order.price));

            if (existingOrder) {
                // If the order exists, add the size
                existingOrder.size = existingOrder.size.add(order.size);
            } else {
                // If the order doesn't exist, push it to the accumulator
                accumulator.push(order);
            }

            return accumulator;
        }, [])
        .sort((a, b) => (a.price.lt(b.price) ? 1 : -1))
        expect((await amm.bidsHead()).toString()).to.equal(orderData[0].price.toString())

        for (let i = 0; i < orderData.length; i++) {
            const { nextTick, amount } = await amm.bids(orderData[i].price)
            expect(amount.toString()).to.equal(orderData[i].size.toString())
            expect(nextTick.toString()).to.equal(i == orderData.length-1 ? '0' : orderData[i+1].price.toString())
        }
    })

    it("asks", async function() {
        expect((await amm.asksHead()).toNumber()).to.equal(0)
        // let orderData = generateRandomArray(15)
        let orderData = [
            { price: 2056, size: 0.5 },
            { price: 2075, size: 0.5 },
            { price: 2022, size: 0.2 },
            { price: 2040, size: 0.5 },
            { price: 2045, size: 0.4 },
            { price: 1955, size: 0.7 },
            { price: 2069, size: 0.7 },
            { price: 2050, size: 0.4 },
            { price: 1978, size: 0.3 },
            { price: 2044, size: 0.5 },
            { price: 2028, size: 0.4 },
            { price: 1993, size: 0.5 },
            { price: 2063, size: 1 },
            { price: 1943, size: 0.4 },
            { price: 2018, size: 0.5 }
          ]
        console.log(orderData)
        orderData = orderData.map(a => {
            return { price: multiplyPrice(a.price), size: multiplySize(a.size * -1) }
        })

        const orders = []
        let requiredMargin = BigNumber.from(0)
        for (let i = 0; i < orderData.length; i++) {
            let shortOrder = getOrderV2(market, alice.address, orderData[i].size, orderData[i].price, getRandomSalt())
            requiredMargin = requiredMargin.add(await utils.getRequiredMarginForShortOrder(shortOrder))
            orders.push(shortOrder)
        }
        await addMargin(alice, requiredMargin)
        const { txReceipt } = await placeV2Orders(orders, alice)
        txReceipt.events.forEach(e => prettyPrintEvents(e))

        orderData = orderData
        .reduce((accumulator, order) => {
            // Find an existing order in the accumulator with the same price
            const existingOrder = accumulator.find(item => item.price.eq(order.price));

            if (existingOrder) {
                // If the order exists, add the size
                existingOrder.size = existingOrder.size.add(order.size);
            } else {
                // If the order doesn't exist, push it to the accumulator
                accumulator.push(order);
            }

            return accumulator;
        }, [])
        .sort((a, b) => (a.price.lt(b.price) ? -1 : 1))
        console.log(orderData.map(a => { return { price: bnToFloat(a.price), size: bnToFloat(a.size, 18) }}))

        console.log('asksHead', (await amm.asksHead()).toString())
        expect((await amm.asksHead()).toString()).to.equal(orderData[0].price.toString())

        for (let i = 0; i < orderData.length; i++) {
            const { nextTick, amount } = await amm.asks(orderData[i].price)
            console.log({
                tick: bnToFloat(orderData[i].price),
                storage: {
                    nextTick: bnToFloat(nextTick),
                    amount: bnToFloat(amount, 18),
                },
                actual: {
                    nextTick: i == orderData.length-1 ? 0 : bnToFloat(orderData[i+1].price),
                    amount: bnToFloat(orderData[i].size.mul(-1), 18),
                }
            })
            expect(amount.toString()).to.equal(orderData[i].size.mul(-1).toString())
            expect(nextTick.toString()).to.equal(i == orderData.length-1 ? '0' : orderData[i+1].price.toString())
        }
    })
})

function prettyPrintEvents(event) {
    // console.log(event.event)
    if (event.event != 'OrderAccepted' && event.event != 'OrderRejected') return
    const res = {
        event: event.event,
        args: {
            order: {
                price: bnToFloat(event.args.order.price),
                size: bnToFloat(event.args.order.baseAssetQuantity, 18),
            }
        }
    }
    if (event.event == 'OrderRejected') {
        res.args.err = event.args.err
    }
    console.log(res)
}

function getRandomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

function getRandomFloat(min, max, decimalPlaces) {
    let rand = Math.random() * (max - min) + min;
    let power = Math.pow(10, decimalPlaces);
    return Math.round(rand * power) / power;
}

function generateRandomArray(n) {
    let arr = [];
    for (let i = 0; i < n; i++) {
        let price = getRandomInt(1900, 2100);
        let size = getRandomFloat(0.1, 1, 1); // 1 decimal place

        // arr.push({
        //     price: multiplyPrice(price),
        //     size: multiplySize(size)
        // });

        arr.push({ price, size });
    }
    return arr;
}
