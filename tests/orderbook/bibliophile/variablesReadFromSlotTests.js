const { ethers, BigNumber } = require('ethers');
const utils = require('../utils.js');
const chai = require('chai');
const { assert, expect } = chai;
let chaiHttp = require('chai-http');

chai.use(chaiHttp);

const {
    _1e6,
    _1e18,
    addMargin,
    alice,
    bob,
    cancelOrderFromLimitOrderV2,
    charlie,
    clearingHouse,
    getAMMContract,
    getIOCOrder,
    getOrderV2,
    getRandomSalt,
    getRequiredMarginForLongOrder,
    getRequiredMarginForShortOrder,
    governance,
    ioc,
    limitOrderBook,
    multiplyPrice,
    multiplySize,
    orderBook,
    placeOrderFromLimitOrderV2,
    placeIOCOrder,
    provider,
    removeAllAvailableMargin,
    url,
} = utils;



describe('Testing variables read from slots by precompile', function () {
    context("Clearing house contract variables", function () {
        // vars read from slot
        // minAllowableMargin, maintenanceMargin, takerFee, amms
        it("should read the correct value from contracts", async function () {
            method = "testing_getClearingHouseVars"
            params =[ charlie.address ]
            response = await makehttpCall(method, params)
            result = response.body.result

            actualMaintenanceMargin = await clearingHouse.maintenanceMargin()
            actualMinAllowableMargin = await clearingHouse.minAllowableMargin()
            actualTakerFee = await clearingHouse.takerFee()
            actualAmms = await clearingHouse.getAMMs()

            expect(result.maintenance_margin).to.equal(actualMaintenanceMargin.toNumber())
            expect(result.min_allowable_margin).to.equal(actualMinAllowableMargin.toNumber())
            expect(result.taker_fee).to.equal(actualTakerFee.toNumber())
            expect(result.amms.length).to.equal(actualAmms.length)
            for(let i = 0; i < result.amms.length; i++) {
                expect(result.amms[i].toLowerCase()).to.equal(actualAmms[i].toLowerCase())
            }
            newMaintenanceMargin = BigNumber.from(20000)
            newMinAllowableMargin = BigNumber.from(40000)
            newTakerFee = BigNumber.from(10000)
            makerFee = await clearingHouse.makerFee()
            referralShare = await clearingHouse.referralShare()
            tradingFeeDiscount = await clearingHouse.tradingFeeDiscount()
            liquidationPenalty = await clearingHouse.liquidationPenalty()
            tx = await clearingHouse.connect(governance).setParams(
                newMaintenanceMargin,
                newMinAllowableMargin,
                newTakerFee,
                makerFee,
                referralShare,
                tradingFeeDiscount,
                liquidationPenalty
            )
            await tx.wait()

            response = await makehttpCall(method, params)
            result = response.body.result

            expect(result.maintenance_margin).to.equal(newMaintenanceMargin.toNumber())
            expect(result.min_allowable_margin).to.equal(newMinAllowableMargin.toNumber())
            expect(result.taker_fee).to.equal(newTakerFee.toNumber())

            // revert config
            tx = await clearingHouse.connect(governance).setParams(
                actualMaintenanceMargin,
                actualMinAllowableMargin,
                actualTakerFee,
                makerFee,
                referralShare,
                tradingFeeDiscount,
                liquidationPenalty
            )
            await tx.wait()
        })
    })
    context("Margin account contract variables", function () {
        // vars read from slot
        // margin, reservedMargin
        it("should read the correct value from contracts", async function () {
            // zero balance
            method ="testing_getMarginAccountVars"
            params =[ 0, charlie.address ]
            response = await makehttpCall(method, params)
            expect(response.body.result.margin).to.equal(0)
            expect(response.body.result.reserved_margin).to.equal(0)

            // add balance for order and then place
            longOrder = getOrderV2(0, charlie.address, multiplySize(0.1), multiplyPrice(2000), BigNumber.from(Date.now()))
            requiredMargin = await getRequiredMarginForLongOrder(longOrder)
            await addMargin(charlie, requiredMargin)
            await placeOrderFromLimitOrderV2(longOrder, charlie)

            method ="testing_getMarginAccountVars"
            params =[ 0, charlie.address ]
            response = await makehttpCall(method, params)

            //cleanup
            await cancelOrderFromLimitOrderV2(longOrder, charlie)
            await removeAllAvailableMargin(charlie)

            expect(response.body.result.margin).to.equal(requiredMargin.toNumber())
            expect(response.body.result.reserved_margin).to.equal(requiredMargin.toNumber())
        })
    })
    context("AMM contract variables", function () {
        // vars read from slot
        // positions, cumulativePremiumFraction, maxOracleSpreadRatio, maxLiquidationRatio, minSizeRequirement, oracle, underlyingAsset,
        // maxLiquidationPriceSpread, redStoneAdapter, redStoneFeedId, impactMarginNotional, lastTradePrice, bids, asks, bidsHead, asksHead
        let ammIndex = 0
        let method ="testing_getAMMVars"
        let ammAddress

        this.beforeAll(async function () {
            amms = await clearingHouse.getAMMs()
            ammAddress = amms[ammIndex]
        })
        context("when variables have default value after setup", async function () {
            it("should read the correct value of variables from contracts", async function () {
                // maxOracleSpreadRatio, maxLiquidationRatio, minSizeRequirement, oracle, underlyingAsset, maxLiquidationPriceSpread
                params =[ ammAddress, ammIndex, charlie.address ]
                response = await makehttpCall(method, params)

                amm = new ethers.Contract(ammAddress, require('../abi/AMM.json'), provider)
                actualMaxOracleSpreadRatio = await amm.maxOracleSpreadRatio()
                actualOracleAddress = await amm.oracle()
                actualMaxLiquidationRatio = await amm.maxLiquidationRatio()
                actualMinSizeRequirement = await amm.minSizeRequirement()
                actualUnderlyingAssetAddress = await amm.underlyingAsset()
                actualMaxLiquidationPriceSpread = await amm.maxLiquidationPriceSpread()

                result = response.body.result
                expect(result.max_oracle_spread_ratio).to.equal(actualMaxOracleSpreadRatio.toNumber())
                expect(result.oracle_address.toLowerCase()).to.equal(actualOracleAddress.toString().toLowerCase())
                expect(result.max_liquidation_ratio).to.equal(actualMaxLiquidationRatio.toNumber())
                expect(String(result.min_size_requirement)).to.equal(actualMinSizeRequirement.toString())
                expect(result.underlying_asset_address.toLowerCase()).to.equal(actualUnderlyingAssetAddress.toString().toLowerCase())
                expect(result.max_liquidation_price_spread).to.equal(actualMaxLiquidationPriceSpread.toNumber())
            })
        })
        context("when variables dont have default value after setup", async function () {
            // positions, cumulativePremiumFraction, redStoneAdapter, redStoneFeedId, impactMarginNotional, lastTradePrice, bids, asks, bidsHead, asksHead
            context("variables which need set config before reading", async function () {
                let amm, oracleAddress, redStoneAdapterAddress, impactMarginNotional
                this.beforeAll(async function () {
                    amm = await getAMMContract(ammIndex)
                    oracleAddress = await amm.oracle()
                    oracle = new ethers.Contract(oracleAddress, require("../abi/Oracle.json"), provider);
                    marginAccount = new ethers.Contract(await amm.marginAccount(), require("../abi/MarginAccount.json"), provider);

                    redStoneAdapterAddress = await oracle.redStoneAdapter()
                    impactMarginNotional = await amm.impactMarginNotional()
                })
                this.afterAll(async function () {
                    await oracle.connect(governance).setRedStoneAdapterAddress(redStoneAdapterAddress)
                    await marginAccount.connect(governance).setOracle(oracleAddress)
                    await amm.connect(governance).setImpactMarginNotional(impactMarginNotional)
                })
                it("should read the correct value from contracts", async function () {
                    newOracleAddress = alice.address
                    newRedStoneAdapterAddress = bob.address
                    newImpactMarginNotional = BigNumber.from(100000)

                    tx = await oracle.connect(governance).setRedStoneAdapterAddress(newRedStoneAdapterAddress)
                    tx = await amm.connect(governance).setImpactMarginNotional(newImpactMarginNotional)
                    await tx.wait()

                    params =[ ammAddress, ammIndex, charlie.address ]
                    response = await makehttpCall(method, params)
                    result = response.body.result

                    expect(result.red_stone_adapter_address.toLowerCase()).to.equal(newRedStoneAdapterAddress.toLowerCase())
                    expect(result.impact_margin_notional).to.equal(newImpactMarginNotional.toNumber())

                    // setOracle
                    tx = await marginAccount.connect(governance).setOracle(newOracleAddress)
                    await tx.wait()
                    response = await makehttpCall(method, params)
                    result = response.body.result
                    expect(result.oracle_address.toLowerCase()).to.equal(newOracleAddress.toLowerCase())
                    expect(result.red_stone_adapter_address.toLowerCase()).to.equal('0x' + '0'.repeat(40)) // red stone adapter should be zero in new oracle
                    expect(result.impact_margin_notional).to.equal(newImpactMarginNotional.toNumber())
                })
            })
            context("variables which need place order before reading", async function () {
                //bids, asks, bidsHead, asksHead
                let longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
                let shortOrderBaseAssetQuantity = multiplySize(-0.1) // 0.1 ether
                let longOrderPrice = multiplyPrice(1799)
                let shortOrderPrice = multiplyPrice(1801)
                let longOrder = getOrderV2(ammIndex, alice.address, longOrderBaseAssetQuantity, longOrderPrice, BigNumber.from(Date.now()), false)
                let shortOrder = getOrderV2(ammIndex, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, BigNumber.from(Date.now()), false)

                this.beforeAll(async function () {
                    requiredMarginAlice = await getRequiredMarginForLongOrder(longOrder)
                    await addMargin(alice, requiredMarginAlice)
                    await placeOrderFromLimitOrderV2(longOrder, alice)
                    requiredMarginBob = await getRequiredMarginForShortOrder(shortOrder)
                    await addMargin(bob, requiredMarginBob)
                    await placeOrderFromLimitOrderV2(shortOrder, bob)
                })

                this.afterAll(async function () {
                    await cancelOrderFromLimitOrderV2(longOrder, alice)
                    await cancelOrderFromLimitOrderV2(shortOrder, bob)
                    await removeAllAvailableMargin(alice)
                    await removeAllAvailableMargin(bob)
                })

                it("should read the correct values from contract", async function () {
                    params =[ ammAddress, ammIndex, alice.address ]
                    response = await makehttpCall(method, params)
                    result = response.body.result
                    expect(result.asks_head).to.equal(shortOrderPrice.toNumber())
                    expect(result.bids_head).to.equal(longOrderPrice.toNumber())
                    expect(String(result.bids_head_size)).to.equal(longOrderBaseAssetQuantity.toString())
                    expect(String(result.asks_head_size)).to.equal(shortOrderBaseAssetQuantity.abs().toString())
                })
            })
            context("variables which need position before reading", async function () {
                let longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
                let shortOrderBaseAssetQuantity = multiplySize(-0.1) // 0.1 ether
                let orderPrice = multiplyPrice(2000)
                let longOrder = getOrderV2(ammIndex, alice.address, longOrderBaseAssetQuantity, orderPrice, BigNumber.from(Date.now()), false)
                let shortOrder = getOrderV2(ammIndex, bob.address, shortOrderBaseAssetQuantity, orderPrice, BigNumber.from(Date.now()), false)

                this.beforeAll(async function () {
                    requiredMarginAlice = await getRequiredMarginForLongOrder(longOrder)
                    await addMargin(alice, requiredMarginAlice)
                    await placeOrderFromLimitOrderV2(longOrder, alice)
                    requiredMarginBob = await getRequiredMarginForShortOrder(shortOrder)
                    await addMargin(bob, requiredMarginBob)
                    await placeOrderFromLimitOrderV2(shortOrder, bob)
                })

                this.afterAll(async function () {
                    oppositeLongOrder = getOrderV2(ammIndex, bob.address, longOrderBaseAssetQuantity, orderPrice, BigNumber.from(Date.now()), true)
                    oppositeShortOrder = getOrderV2(ammIndex, alice.address, shortOrderBaseAssetQuantity, orderPrice, BigNumber.from(Date.now()), true)
                    await placeOrderFromLimitOrderV2(oppositeLongOrder, bob)
                    await placeOrderFromLimitOrderV2(oppositeShortOrder, alice)
                    await utils.waitForOrdersToMatch()
                    await removeAllAvailableMargin(alice)
                    await removeAllAvailableMargin(bob)
                })

                it("should read the correct values from contract", async function () {
                    params =[ ammAddress, ammIndex, alice.address ]
                    resultAlice = (await makehttpCall(method, params)).body.result
                    params =[ ammAddress, ammIndex, bob.address ]
                    resultBob = (await makehttpCall(method, params)).body.result

                    expect(String(resultAlice.position.size)).to.equal(longOrderBaseAssetQuantity.toString())
                    expect(String(resultAlice.position.open_notional)).to.equal(longOrderBaseAssetQuantity.mul(orderPrice).div(_1e18).toString())
                    expect(String(resultBob.position.size)).to.equal(shortOrderBaseAssetQuantity.toString())
                    expect(String(resultBob.position.open_notional)).to.equal(shortOrderBaseAssetQuantity.mul(orderPrice).abs().div(_1e18).toString())
                    expect(resultAlice.last_price).to.equal(orderPrice.toNumber())
                    expect(resultBob.last_price).to.equal(orderPrice.toNumber())
                })
            })
        })
    })
    context("IOC order contract variables", function () {
        let method ="testing_getIOCOrdersVars"
        let longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
        let shortOrderBaseAssetQuantity = multiplySize(-0.1) // 0.1 ether
        let orderPrice = multiplyPrice(2000)
        let market = BigNumber.from(0)


        context("variable which have default value after setup", async function () {
            //ioc expiration cap
            it("should read the correct value from contracts", async function () {
                params = [ "0xe97a0702264091714ea19b481c1fd12d9686cb4602efbfbec41ec5ea5410da84"]

                result = (await makehttpCall(method, params)).body.result
                actualExpirationCap = await ioc.expirationCap()
                expect(result.ioc_expiration_cap).to.eq(actualExpirationCap.toNumber())
            })
        })
        context("variable which need place order before reading", async function () {
            //blockPlaced, filledAmount, orderStatus
            context("for a long IOC order", async function () {
                it("returns correct value when order is not placed", async function () {
                    latestBlockNumber = await provider.getBlockNumber()
                    lastTimestamp = (await provider.getBlock(latestBlockNumber)).timestamp
                    expireAt = lastTimestamp + 6
                    longIOCOrder = getIOCOrder(expireAt, market, charlie.address, longOrderBaseAssetQuantity, orderPrice, getRandomSalt(), false)
                    orderHash = await ioc.getOrderHash(longIOCOrder)
                    params = [ orderHash ]
                    result = (await makehttpCall(method, params)).body.result
                    expect(result.order_details.block_placed).to.eq(0)
                    expect(result.order_details.filled_amount).to.eq(0)
                    expect(result.order_details.order_status).to.eq(0)
                })
                it("returns correct value when order is placed", async function () {
                    let charlieBalance = _1e6.mul(150)
                    await addMargin(charlie, charlieBalance)

                    //placing order
                    latestBlockNumber = await provider.getBlockNumber()
                    lastTimestamp = (await provider.getBlock(latestBlockNumber)).timestamp
                    expireAt = lastTimestamp + 6
                    longIOCOrder = getIOCOrder(expireAt, market, charlie.address, longOrderBaseAssetQuantity, orderPrice, getRandomSalt(), false)
                    orderHash = await ioc.getOrderHash(longIOCOrder)
                    params = [ orderHash ]
                    txDetails = await placeIOCOrder(longIOCOrder, charlie)
                    result = (await makehttpCall(method, params)).body.result

                    //cleanup
                    await removeAllAvailableMargin(charlie)

                    actualBlockPlaced = txDetails.txReceipt.blockNumber
                    expect(result.order_details.block_placed).to.eq(actualBlockPlaced)
                    expect(result.order_details.filled_amount).to.eq(0)
                    expect(result.order_details.order_status).to.eq(1)

                })
            })
            context("for a short IOC order", async function () {
                it("returns correct value when order is not placed", async function () {
                    latestBlockNumber = await provider.getBlockNumber()
                    lastTimestamp = (await provider.getBlock(latestBlockNumber)).timestamp
                    expireAt = lastTimestamp + 6
                    shortIOCOrder = getIOCOrder(expireAt, market, charlie.address, shortOrderBaseAssetQuantity, orderPrice, getRandomSalt(), false)
                    orderHash = await ioc.getOrderHash(shortIOCOrder)
                    params = [ orderHash ]
                    result = (await makehttpCall(method, params)).body.result
                    expect(result.order_details.block_placed).to.eq(0)
                    expect(result.order_details.filled_amount).to.eq(0)
                    expect(result.order_details.order_status).to.eq(0)
                })
                it("returns correct value when order is placed", async function () {
                    let charlieBalance = _1e6.mul(150)
                    await addMargin(charlie, charlieBalance)

                    //placing order
                    latestBlockNumber = await provider.getBlockNumber()
                    lastTimestamp = (await provider.getBlock(latestBlockNumber)).timestamp
                    expireAt = lastTimestamp + 6
                    shortIOCOrder = getIOCOrder(expireAt, market, charlie.address, shortOrderBaseAssetQuantity, orderPrice, getRandomSalt(), false)
                    orderHash = await ioc.getOrderHash(shortIOCOrder)
                    params = [ orderHash ]
                    txDetails = await placeIOCOrder(shortIOCOrder, charlie)
                    result = (await makehttpCall(method, params)).body.result

                    //cleanup
                    await removeAllAvailableMargin(charlie)

                    actualBlockPlaced = txDetails.txReceipt.blockNumber
                    expect(result.order_details.block_placed).to.eq(actualBlockPlaced)
                    expect(result.order_details.filled_amount).to.eq(0)
                    expect(result.order_details.order_status).to.eq(1)
                })
            })
        })
    })
    context("order book contract variables", function () {
        let method ="testing_getOrderBookVars"
        let traderAddress = alice.address
        let senderAddress = charlie.address
        let longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
        let shortOrderBaseAssetQuantity = multiplySize(-0.1) // 0.1 ether
        let orderPrice = multiplyPrice(2000)
        let market = BigNumber.from(0)

        context("variables which dont need place order before reading", async function () {
            let params = [ traderAddress, senderAddress, "0xe97a0702264091714ea19b481c1fd12d9686cb4602efbfbec41ec5ea5410da84" ]
            //isTradingAuthority
            it("should return false when sender is not a tradingAuthority for an address", async function () {
                result = (await makehttpCall(method, params)).body.result
                expect(result.is_trading_authority).to.eq(false)
            })
            // need to implement adding trading authority for an address
            it.skip("should return true when sender is a tradingAuthority for an address", async function () {
                await orderBook.connect(alice).setTradingAuthority(traderAddress, senderAddress)
                result = (await makehttpCall(method, params)).body.result
                expect(result.is_trading_authority).to.eq(true)
            })
        })
        context("variables which need place order before reading", async function () {
            context("for a long limit order", async function () {
                let longOrder = getOrderV2(market, traderAddress, longOrderBaseAssetQuantity, orderPrice, getRandomSalt())
                it("returns correct value when order is not placed", async function () {
                    orderHash = await limitOrderBook.getOrderHash(longOrder)
                    params = [ traderAddress, senderAddress, orderHash ]
                    result = (await makehttpCall(method, params)).body.result
                    expect(result.order_details.block_placed).to.eq(0)
                    expect(result.order_details.filled_amount).to.eq(0)
                    expect(result.order_details.order_status).to.eq(0)
                })
                it("returns correct value when order is placed", async function () {
                    requiredMargin = await getRequiredMarginForLongOrder(longOrder)
                    await addMargin(alice, requiredMargin)
                    orderHash = await limitOrderBook.getOrderHash(longOrder)
                    const {txReceipt} = await placeOrderFromLimitOrderV2(longOrder, alice)
                    params = [ traderAddress, traderAddress, orderHash ]
                    result = (await makehttpCall(method, params)).body.result
                    // cleanup
                    await cancelOrderFromLimitOrderV2(longOrder, alice)
                    await removeAllAvailableMargin(alice)

                    expectedBlockPlaced = txReceipt.blockNumber
                    expect(result.order_details.block_placed).to.eq(expectedBlockPlaced)
                    expect(result.order_details.filled_amount).to.eq(0)
                    expect(result.order_details.order_status).to.eq(1)
                })
            })
            context("for a short limit order", async function () {
                let shortOrder = getOrderV2(market, traderAddress, shortOrderBaseAssetQuantity, orderPrice, getRandomSalt())
                it("returns correct value when order is not placed", async function () {
                    orderHash = await limitOrderBook.getOrderHash(shortOrder)
                    params = [ traderAddress, traderAddress, orderHash ]
                    result = (await makehttpCall(method, params)).body.result
                    expect(result.order_details.block_placed).to.eq(0)
                    expect(result.order_details.filled_amount).to.eq(0)
                    expect(result.order_details.order_status).to.eq(0)
                })
                it("returns correct value when order is placed", async function () {
                    requiredMargin = await getRequiredMarginForShortOrder(shortOrder)
                    await addMargin(alice, requiredMargin)

                    orderHash = await limitOrderBook.getOrderHash(shortOrder)
                    const { txReceipt } = await placeOrderFromLimitOrderV2(shortOrder, alice)
                    params = [ traderAddress, traderAddress, orderHash ]
                    result = (await makehttpCall(method, params)).body.result
                    // cleanup
                    await cancelOrderFromLimitOrderV2(shortOrder, alice)
                    await removeAllAvailableMargin(alice)

                    expectedBlockPlaced = txReceipt.blockNumber
                    expect(result.order_details.block_placed).to.eq(expectedBlockPlaced)
                    expect(result.order_details.filled_amount).to.eq(0)
                    expect(result.order_details.order_status).to.eq(1)
                })
            })
        })
    })
})

async function makehttpCall(method, params=[]) {
    body = {
        "jsonrpc":"2.0",
        "id" :1,
        "method" : method,
        "params" : params
    }

    const serverUrl = url.split("/").slice(0, 3).join("/")
    path = "/".concat(url.split("/").slice(3).join("/"))
    return chai.request(serverUrl)
        .post(path)
        .send(body)
}
