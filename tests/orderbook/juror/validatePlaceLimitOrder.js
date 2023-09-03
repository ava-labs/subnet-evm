const { expect } = require("chai");
const { BigNumber } = require("ethers");
const utils = require("../utils")

const {
    _1e6,
    _1e18,
    addMargin,
    alice,
    bob,
    cancelOrderFromLimitOrderV2,
    clearingHouse,
    getMinSizeRequirement,
    getOrderV2,
    getRandomSalt,
    juror,
    multiplyPrice,
    multiplySize,
    orderBook,
    placeOrderFromLimitOrderV2,
    removeAllAvailableMargin,
    waitForOrdersToMatch,
} = utils

describe("Test validatePlaceLimitOrder", async function () {
    market = BigNumber.from(0)
    longBaseAssetQuantity = multiplySize(0.1)
    shortBaseAssetQuantity = multiplySize(-0.1)
    price = multiplyPrice(1800)
    initialMargin = multiplyPrice(600000)

    context("when order's baseAssetQuantity is 0", async function () {
        it("returns error", async function () {
            longOrder = getOrderV2(market, alice.address, 0, price, getRandomSalt())
            response = await juror.validatePlaceLimitOrder(longOrder, alice.address)
            expect(response.err).to.eq("baseAssetQuantity is zero")
            longOrderHash = await orderBook.getOrderHashV2(longOrder)
            expect(response.orderHash).to.eq(longOrderHash)
            // expect(response.res.reserveAmount.toNumber()).to.eq(0)
            // expect(response.res.amm.toString()).to.eq("0x0000000000000000000000000000000000000000")
        })
    })

    context("when order's baseAssetQuantity is not 0", async function () {
        context("when order's baseAssetQuantity is not a multiple of minSizeRequirement", async function () {
            it("returns error when order's baseAssetQuantity.abs() is > 0 but < minSizeRequirement", async function () {
                minSizeRequirement = await getMinSizeRequirement(market)
                let invalidLongBaseAssetQuantity = minSizeRequirement.sub(1)

                //longOrder
                longOrder = getOrderV2(market, alice.address, invalidLongBaseAssetQuantity, price, getRandomSalt())
                response = await juror.validatePlaceLimitOrder(longOrder, alice.address)
                expect(response.err).to.eq("not multiple")
                // longOrderHash = await orderBook.getOrderHashV2(longOrder)
                // expect(response.orderHash).to.eq(longOrderHash)
                // expect(response.res.reserveAmount.toNumber()).to.eq(0)
                // expect(response.res.amm.toString()).to.eq("0x0000000000000000000000000000000000000000")

                //shortOrder
                shortOrder = getOrderV2(market, alice.address, invalidLongBaseAssetQuantity.mul("-1"), price, getRandomSalt())
                response = await juror.validatePlaceLimitOrder(shortOrder, alice.address)
                expect(response.err).to.eq("not multiple")
                // shortOrderHash = await orderBook.getOrderHashV2(longOrder)
                // expect(response.orderHash).to.eq(shortOrderHash)
                // expect(response.res.reserveAmount.toNumber()).to.eq(0)
                // expect(response.res.amm.toString()).to.eq("0x0000000000000000000000000000000000000000")
            })
            it("returns error when order's baseAssetQuantity.abs() is > minSizeRequirement", async function () {
                minSizeRequirement = await getMinSizeRequirement(market)
                let invalidLongBaseAssetQuantity = minSizeRequirement.mul(3).div(2)

                //longOrder
                longOrder = getOrderV2(market, alice.address, invalidLongBaseAssetQuantity, price, getRandomSalt())
                response = await juror.validatePlaceLimitOrder(longOrder, alice.address)
                expect(response.err).to.eq("not multiple")
                // longOrderHash = await orderBook.getOrderHashV2(longOrder)
                // expect(response.orderHash).to.eq(longOrderHash)
                // expect(response.res.reserveAmount.toNumber()).to.eq(0)
                // expect(response.res.amm.toString()).to.eq("0x0000000000000000000000000000000000000000")

                //shortOrder
                shortOrder = getOrderV2(market, alice.address, invalidLongBaseAssetQuantity.mul("-1"), price, getRandomSalt())
                response = await juror.validatePlaceLimitOrder(shortOrder, alice.address)
                expect(response.err).to.eq("not multiple")
                // shortOrderHash = await orderBook.getOrderHashV2(longOrder)
                // expect(response.orderHash).to.eq(shortOrderHash)
                // expect(response.res.reserveAmount.toNumber()).to.eq(0)
                // expect(response.res.amm.toString()).to.eq("0x0000000000000000000000000000000000000000")
            })
        })
        context("when order's quoteAssetQuantity is a multiple of minSizeRequirement", async function () {
            context("when order was already placed", async function () {
                this.beforeAll(async function() {
                    await addMargin(alice, initialMargin)
                    await addMargin(bob, initialMargin)
                })
                this.afterAll(async function() {
                    await removeAllAvailableMargin(alice)
                    await removeAllAvailableMargin(bob)
                })
                context("when order's status is placed", function() {
                    it("returns error for a longOrder", async function() {
                        let longOrder = getOrderV2(market, alice.address, longBaseAssetQuantity, price, getRandomSalt())
                        console.log("placing order")
                        response = await placeOrderFromLimitOrderV2(longOrder, alice)
                        response = await juror.validatePlaceLimitOrder(longOrder, alice.address)
                        //cleanup
                        await cancelOrderFromLimitOrderV2(longOrder, alice)

                        expect(response.err).to.eq("order already exists")
                        longOrderHash = await orderBook.getOrderHashV2(longOrder)
                        expect(response.orderHash).to.eq(longOrderHash)
                        // expect(response.res.reserveAmount.toNumber()).to.eq(0)
                        // expect(response.res.amm.toString()).to.eq("0x0000000000000000000000000000000000000000")
                    })
                    it("returns error for a shortOrder", async function() {
                        let shortOrder = getOrderV2(market, bob.address, shortBaseAssetQuantity, price, getRandomSalt())
                        await placeOrderFromLimitOrderV2(shortOrder, bob)
                        response = await juror.validatePlaceLimitOrder(shortOrder, bob.address)
                        //cleanup
                        await cancelOrderFromLimitOrderV2(shortOrder, bob)

                        expect(response.err).to.eq("order already exists")
                        shortOrderHash = await orderBook.getOrderHashV2(shortOrder)
                        expect(response.orderHash).to.eq(shortOrderHash)
                        // expect(response.res.reserveAmount.toNumber()).to.eq(0)
                        // expect(response.res.amm.toString()).to.eq("0x0000000000000000000000000000000000000000")
                    })
                })
                context.skip("when order status is filled", async function () {
                    it("returns error", async function() {
                        await utils.enableValidatorMatching()
                        let longOrder = getOrderV2(market, alice.address, longBaseAssetQuantity, price, getRandomSalt())
                        let shortOrder = getOrderV2(market, bob.address, shortBaseAssetQuantity, price, getRandomSalt())
                        await placeOrderFromLimitOrderV2(longOrder, alice)
                        await placeOrderFromLimitOrderV2(shortOrder, bob)
                        await waitForOrdersToMatch()
                        responseLong = await juror.validatePlaceLimitOrder(longOrder, alice.address)
                        responseShort = await juror.validatePlaceLimitOrder(shortOrder, alice.address)
                        //cleanup
                        await placeOrderFromLimitOrderV2(longOrder, bob)
                        await placeOrderFromLimitOrderV2(shortOrder, alice)
                        await waitForOrdersToMatch()

                        expect(responseLong.err).to.eq("order already exists")
                        longOrderHash = await orderBook.getOrderHashV2(longOrder)
                        expect(responseLong.orderHash).to.eq(longOrderHash)
                        // expect(response.res.reserveAmount.toNumber()).to.eq(0)
                        // expect(response.res.amm.toString()).to.eq("0x0000000000000000000000000000000000000000")

                        // expect(responseShort.err).to.eq("order already exists")
                        // shortOrderHash = await orderBook.getOrderHashV2(shortOrder)
                        // expect(responseShort.orderHash).to.eq(shortOrderHash)
                        // expect(response.res.reserveAmount.toNumber()).to.eq(0)
                        // expect(response.res.amm.toString()).to.eq("0x0000000000000000000000000000000000000000")
                    })
                })
                context("when order status is cancelled", async function () {
                    it("returns error for a longOrder", async function() {
                        let longOrder = getOrderV2(market, alice.address, longBaseAssetQuantity, price, getRandomSalt())
                        await placeOrderFromLimitOrderV2(longOrder, alice)
                        await cancelOrderFromLimitOrderV2(longOrder, alice)

                        response = await juror.validatePlaceLimitOrder(longOrder, alice.address)
                        expect(response.err).to.eq("order already exists")
                        longOrderHash = await orderBook.getOrderHashV2(longOrder)
                        expect(response.orderHash).to.eq(longOrderHash)
                        // expect(response.res.reserveAmount.toNumber()).to.eq(0)
                        // expect(response.res.amm.toString()).to.eq("0x0000000000000000000000000000000000000000")
                    })
                    it.skip("returns error for a shortOrder", async function() {
                        let shortOrder = getOrderV2(market, bob.address, shortBaseAssetQuantity, price, getRandomSalt())
                        await addMargin(bob, initialMargin)
                        await placeOrderFromLimitOrderV2(shortOrder, bob)
                        await cancelOrderFromLimitOrderV2(shortOrder, bob)
                        response = await juror.validatePlaceLimitOrder(shortOrder, bob.address)
                        //cleanup
                        await removeAllAvailableMargin(bob)

                        expect(response.err).to.eq("order already exists")
                        shortOrderHash = await orderBook.getOrderHashV2(shortOrder)
                        expect(response.orderHash).to.eq(shortOrderHash)
                        // expect(response.res.reserveAmount.toNumber()).to.eq(0)
                        // expect(response.res.amm.toString()).to.eq("0x0000000000000000000000000000000000000000")
                    })
                })
            })
            context("when order was never placed", async function () {
                context("when order is not reduceOnly", async function () {
                    context.skip("when order is in opposite direction to currentPosition and trader has unfilled reduceOnly Orders", async function() {
                        this.beforeEach(async function() {
                            await addMargin(alice, initialMargin)
                            await addMargin(bob, initialMargin)
                        })
                        this.afterEach(async function() {
                            await removeAllAvailableMargin(alice)
                            await removeAllAvailableMargin(bob)
                        })
                        it("returns error", async function () {
                            let longOrder = getOrderV2(market, alice.address, longBaseAssetQuantity, price, getRandomSalt())
                            let shortOrder = getOrderV2(market, bob.address, shortBaseAssetQuantity, price, getRandomSalt())
                            await placeOrderFromLimitOrderV2(longOrder, alice)
                            await placeOrderFromLimitOrderV2(shortOrder, bob)
                            await waitForOrdersToMatch()

                            longReduceOnlyOrder = getOrderV2(market, bob.address, longBaseAssetQuantity.div(2), price, getRandomSalt())
                            await placeOrderFromLimitOrderV2(longReduceOnlyOrder, bob)
                            response = await juror.validatePlaceLimitOrder(longOrder, bob.address)
                            //cleanup
                            await cancelOrderFromLimitOrderV2(longReduceOnlyOrder, bob)
                            await placeOrderFromLimitOrderV2(longOrder, bob)
                            await placeOrderFromLimitOrderV2(shortOrder, alice)
                            await waitForOrdersToMatch()

                            expect(responseLong.err).to.eq("")
                            longOrderHash = await orderBook.getOrderHashV2(longOrder)
                            expect(responseLong.orderHash).to.eq(longOrderHash)
                        })
                    })
                    context("when order is not in opposite direction to currentPostion if trader has unfilled reduceOnly orders", async function() {
                        context("when trader does not have sufficient margin", async function() {
                            it("returns error", async function() {
                                await removeAllAvailableMargin(alice)
                                let longOrder = getOrderV2(market, alice.address, longBaseAssetQuantity, price, getRandomSalt())
                                longOrderHash = await orderBook.getOrderHashV2(longOrder)
                                response = await juror.validatePlaceLimitOrder(longOrder, alice.address)
                                expect(response.err).to.eq("insufficient margin")
                                expect(response.orderHash).to.eq(longOrderHash)
                            })
                        })
                        context("when trader has sufficient margin", async function () {
                            context("when order is not postOnly", async function() {
                                this.beforeAll(async function() {
                                    await addMargin(alice, initialMargin)
                                })
                                this.afterAll(async function() {
                                    await removeAllAvailableMargin(alice)
                                })
                                it("returns success", async function () {
                                    let longOrder = getOrderV2(market, alice.address, longBaseAssetQuantity, price, getRandomSalt())
                                    response = await juror.validatePlaceLimitOrder(longOrder, alice.address)

                                    minAllowableMargin = await clearingHouse.minAllowableMargin()
                                    takerFee = await clearingHouse.takerFee()
                                    expect(response.err).to.eq("")
                                    quoteAsset = longBaseAssetQuantity.mul(price).div(_1e18)
                                    expectedReserveAmount = quoteAsset.mul(minAllowableMargin).div(_1e6)
                                    expectedTakerFee = quoteAsset.mul(takerFee).div(_1e6)
                                    expect(response.res.reserveAmount.toString()).to.eq(expectedReserveAmount.add(expectedTakerFee).toString())
                                    longOrderHash = await orderBook.getOrderHashV2(longOrder)
                                    expect(response.orderHash).to.eq(longOrderHash)
                                })
                            })
                            context.skip("when order is postOnly", async function () {
                                this.beforeAll(async function() {
                                    await addMargin(alice, initialMargin)
                                })
                                this.afterAll(async function() {
                                    await removeAllAvailableMargin(alice)
                                })
                                context.skip("for a long order", async function() {
                                    context("when there is no asks in orderbook", async function() {
                                        it("returns success", async function() {
                                            longPostOnlyOrder = getOrderV2(market, alice.address, longBaseAssetQuantity, price, getRandomSalt(), false, true)
                                            response = await juror.validatePlaceLimitOrder(longPostOnlyOrder, alice.address)

                                            minAllowableMargin = await clearingHouse.minAllowableMargin()
                                            takerFee = await clearingHouse.takerFee()
                                            longPostOnlyOrderHash = await orderBook.getOrderHashV2(longPostOnlyOrder, alice.address)
                                            quoteAsset = longBaseAssetQuantity.mul(price).div(_1e18)
                                            expectedReserveAmount = quoteAsset.mul(minAllowableMargin).div(_1e6)
                                            expectedTakerFee = quoteAsset.mul(takerFee).div(_1e6)
                                            expect(response.err).to.eq("")
                                            expect(response.res.reserveAmount.toString()).to.eq(expectedReserveAmount.add(expectedTakerFee).toString())
                                            expect(response.orderHash).to.eq(longPostOnlyOrderHash)
                                        })
                                    })
                                    context("when there are asks in orderbook", async function() {
                                        let shortOrder = getOrderV2(market, bob.address, shortBaseAssetQuantity, price, getRandomSalt(), false, false)
                                        let minAllowableMargin, takerFee

                                        this.beforeAll(async function(){
                                            console.log("inner beforeall")
                                            minAllowableMargin = await clearingHouse.minAllowableMargin()
                                            takerFee = await clearingHouse.takerFee()
                                            await addMargin(bob, initialMargin)
                                            await placeOrderFromLimitOrderV2(shortOrder, bob)
                                        })
                                        this.afterAll(async function() {
                                            await cancelOrderFromLimitOrderV2(shortOrder, bob)
                                            await removeAllAvailableMargin(bob)
                                        })

                                        context("when order's price < asksHead price", async function(){
                                            it("returns success", async function() {
                                                newPrice = price.sub(1)
                                                longPostOnlyOrder = getOrderV2(market, alice.address, longBaseAssetQuantity, newPrice, getRandomSalt(), false, true)
                                                response = await juror.validatePlaceLimitOrder(longPostOnlyOrder, alice.address)

                                                longPostOnlyOrderHash = await orderBook.getOrderHashV2(longPostOnlyOrder, alice.address)
                                                quoteAsset = longBaseAssetQuantity.mul(newPrice).div(_1e18)
                                                expectedReserveAmount = quoteAsset.mul(minAllowableMargin).div(_1e6)
                                                expectedTakerFee = quoteAsset.mul(takerFee).div(_1e6)
                                                expect(response.err).to.eq("")
                                                expect(response.res.reserveAmount.toString()).to.eq(expectedReserveAmount.add(expectedTakerFee).toString())
                                                expect(response.orderHash).to.eq(longPostOnlyOrderHash)
                                            })
                                        })
                                        context("when order's price >= asksHead price", async function(){
                                            it("returns error if price == asksHead", async function() {
                                                shortOrder = getOrderV2(market, bob.address, shortBaseAssetQuantity, price, getRandomSalt(), false, false)
                                                await placeOrderFromLimitOrderV2(shortOrder, bob)

                                                longPostOnlyOrder = getOrderV2(market, alice.address, longBaseAssetQuantity, price, getRandomSalt(), false, true)
                                                response = await juror.validatePlaceLimitOrder(longPostOnlyOrder, alice.address)
                                                longPostOnlyOrderHash = await orderBook.getOrderHashV2(longPostOnlyOrder, alice.address)
                                                expect(response.err).to.eq("crossing market")
                                                expect(response.orderHash).to.eq(longPostOnlyOrderHash)
                                            })
                                            it("returns error if price > asksHead", async function() {
                                                newPrice = price.add(1)
                                                shortOrder = getOrderV2(market, bob.address, shortBaseAssetQuantity, price, getRandomSalt(), false, false)
                                                await placeOrderFromLimitOrderV2(shortOrder, bob)

                                                longPostOnlyOrder = getOrderV2(market, alice.address, longBaseAssetQuantity, newPrice, getRandomSalt(), false, true)
                                                response = await juror.validatePlaceLimitOrder(longPostOnlyOrder, alice.address)
                                                longPostOnlyOrderHash = await orderBook.getOrderHashV2(longPostOnlyOrder, alice.address)
                                                expect(response.err).to.eq("crossing market")
                                                expect(response.orderHash).to.eq(longPostOnlyOrderHash)
                                            })
                                        })
                                    })
                                })
                                context("for a short order", async function() {
                                    context("when there is no bids in orderbook", async function() {
                                        it("returns success", async function(){
                                        })
                                    })
                                    context("when there are bids in orderbook", async function() {
                                        context("when order's price < asksHead price", async function(){
                                            it("returns error", async function() {
                                            })
                                        })
                                        context("when order's price > asksHead price", async function(){
                                            it("returns success", async function(){
                                            })
                                        })
                                    })
                                })
                            })
                        })
                    })
                })
                context("when order is reduceOnly", async function () {
                    this.beforeEach(async function () {
                    })
                    this.afterEach(async function () {
                    })

                    context("when order is not opposite of currentPosition", async function () {
                        it("returns error", async function () {
                            await removeAllAvailableMargin(alice)
                            await removeAllAvailableMargin(bob)
                            await addMargin(alice, initialMargin)
                            await addMargin(bob, initialMargin)
                            longOrder = getOrderV2(market, alice.address, longBaseAssetQuantity, price, getRandomSalt())
                            shortOrder = getOrderV2(market, bob.address, shortBaseAssetQuantity, price, getRandomSalt())
                            await placeOrderFromLimitOrderV2(longOrder, alice)
                            await placeOrderFromLimitOrderV2(shortOrder, bob)
                            await waitForOrdersToMatch()

                            longOrder.trader = bob.address
                            await placeOrderFromLimitOrderV2(longOrder, bob)
                            shortOrder.trader = bob.address
                            await placeOrderFromLimitOrderV2(shortOrder, alice)
                            await waitForOrdersToMatch()
                            await removeAllAvailableMargin(alice)
                            await removeAllAvailableMargin(bob)
                            return
                            longReduceOnlyOrder = getOrderV2(market, alice.address, longBaseAssetQuantity, price, getRandomSalt(), true)
                            shortReduceOnlyOrder = getOrderV2(market, shortBaseAssetQuantity, price, getRandomSalt(), true)
                            //longOrder
                            response = await juror.validatePlaceLimitOrder(longReduceOnlyOrder, alice.address)
                            expect(response.err).to.eq("reduce only order must reduce position")
                            // longOrderHash = await orderBook.getOrderHashV2(longOrder)
                            // expect(response.orderHash).to.eq(longOrderHash)
                            // expect(response.res.reserveAmount.toNumber()).to.eq(0)
                            // expect(response.res.amm.toString()).to.eq("0x0000000000000000000000000000000000000000")

                            //shortOrder
                            shortOrder = getOrderV2(market, alice.address, invalidLongBaseAssetQuantity.mul("-1"), price, getRandomSalt())
                            response = await juror.validatePlaceLimitOrder(shortOrder, alice.address)
                            expect(response.err).to.eq("reduce only order must reduce position")
                            // shortOrderHash = await orderBook.getOrderHashV2(longOrder)
                            // expect(response.orderHash).to.eq(shortOrderHash)
                            // expect(response.res.reserveAmount.toNumber()).to.eq(0)
                            // expect(response.res.amm.toString()).to.eq("0x0000000000000000000000000000000000000000")
                        })
                    })
                    context("when order is opposite of currentPosition", async function () {
                        context("when order is longOrder", async function () {
                            context("when trader already has open longOrders", async function () {
                                it("returns error", async function () {
                                })
                            })
                            context("when trader does not have open longOrders", async function () {
                                context("when order's baseAssetQuantity + reduceOnlyAmount > trader's longPosition", async function () {
                                    it("returns error", async function () {
                                    })
                                })
                            })
                        })
                        context("when order is shortOrder", async function () {
                            context("when trader already has open shortOrders", async function () {
                                it("returns error", async function () {
                                })
                            })
                            context("when trader does not have open shortOrders", async function () {
                                context("when order's baseAssetQuantity + reduceOnlyAmount > trader's shortPosition", async function () {
                                    it("returns error", async function () {
                                    })
                                })
                            })
                        })
                    })
                })
            })
        })
    })
})
