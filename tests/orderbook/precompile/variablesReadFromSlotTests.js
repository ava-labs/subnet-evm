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
    cancelOrderFromLimitOrder,
    charlie,
    clearingHouse,
    getIOCOrder,
    getOrder,
    governance,
    ioc,
    marginAccount,
    multiplyPrice,
    multiplySize,
    orderBook,
    provider,
    removeAllAvailableMargin,
    sleep,
    url,
} = utils;



describe('Testing variables read from slots by precompile', function () {
    context("Clearing house contract variables", function () {
        it("should read the correct value from contracts", async function () {
            method = "testing_getClearingHouseVars"
            params =[ charlie.address ]
            response = await makehttpCall(method, params)
            result = response.body.result

            actualMaintenanceMargin = await clearingHouse.maintenanceMargin()
            actualMinAllowableMargin = await clearingHouse.minAllowableMargin()
            actualAmms = await clearingHouse.getAMMs()

            expect(result.maintenance_margin).to.equal(actualMaintenanceMargin.toNumber())
            expect(result.min_allowable_margin).to.equal(actualMinAllowableMargin.toNumber())
            expect(result.amms.length).to.equal(actualAmms.length)
            for(let i = 0; i < result.amms.length; i++) {
                expect(result.amms[i].toLowerCase()).to.equal(actualAmms[i].toLowerCase())
            }
            newMaintenanceMargin = BigNumber.from(20000)
            newMinAllowableMargin = BigNumber.from(40000)
            takerFee = await clearingHouse.takerFee()
            makerFee = await clearingHouse.makerFee()
            referralShare = await clearingHouse.referralShare()
            tradingFeeDiscount = await clearingHouse.tradingFeeDiscount()
            liquidationPenalty = await clearingHouse.liquidationPenalty()
            tx = await clearingHouse.connect(governance).setParams(
                newMaintenanceMargin,
                newMinAllowableMargin,
                takerFee,
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

            // revert config
            tx = await clearingHouse.connect(governance).setParams(
                actualMaintenanceMargin,
                actualMinAllowableMargin,
                takerFee,
                makerFee,
                referralShare,
                tradingFeeDiscount,
                liquidationPenalty
            )
            await tx.wait()
        })
    })

    context("Margin account contract variables", function () {
        it("should read the correct value from contracts", async function () {
            await removeAllAvailableMargin(charlie)

            let charlieBalance = _1e6.mul(150)
            await addMargin(charlie, charlieBalance)

            vusdIdx = BigNumber.from(0)
            method ="testing_getMarginAccountVars"
            params =[ 0, charlie.address ]
            response = await makehttpCall(method, params)

            actualMargin = await marginAccount.getAvailableMargin(charlie.address)
            expect(response.body.result.margin).to.equal(charlieBalance.toNumber())
            expect(actualMargin.toNumber()).to.equal(charlieBalance.toNumber())
        })
    })

    context("AMM contract variables", function () {
        it("should read the correct value from contracts", async function () {
            amms = await clearingHouse.getAMMs()
            ammIndex = 0
            ammAddress = amms[ammIndex]
            amm = new ethers.Contract(ammAddress, require('../abi/AMM.json'), provider)

            actualLastPrice = await amm.lastPrice()
            actualCumulativePremiumFraction = await amm.cumulativePremiumFraction()
            actualMaxOracleSpreadRatio = await amm.maxOracleSpreadRatio()
            actualOracleAddress = await amm.oracle()
            actualMaxLiquidationRatio = await amm.maxLiquidationRatio()
            actualMinSizeRequirement = await amm.minSizeRequirement()
            actualUnderlyingAssetAddress = await amm.underlyingAsset()
            actualMaxLiquidationPriceSpread = await amm.maxLiquidationPriceSpread()
            actualRedStoneAdapterAddress = await amm.redStoneAdapter()
            actualRedStoneFeedId = await amm.redStoneFeedId()
            actualPosition = await amm.positions(charlie.address)

            // testing for amms[0]
            method ="testing_getAMMVars"
            params =[ammAddress, ammIndex, charlie.address]
            response = await makehttpCall(method, params)
            result = response.body.result
            // expect(result.last_price).to.equal(actualLastPrice.toNumber())
            expect(result.cumulative_premium_fraction).to.equal(actualCumulativePremiumFraction.toNumber())
            expect(result.max_oracle_spread_ratio).to.equal(actualMaxOracleSpreadRatio.toNumber())
            expect(result.oracle_address.toLowerCase()).to.equal(actualOracleAddress.toString().toLowerCase())
            expect(result.max_liquidation_ratio).to.equal(actualMaxLiquidationRatio.toNumber())
            expect(String(result.min_size_requirement)).to.equal(actualMinSizeRequirement.toString())
            expect(result.underlying_asset_address.toLowerCase()).to.equal(actualUnderlyingAssetAddress.toString().toLowerCase())
            expect(result.max_liquidation_price_spread).to.equal(actualMaxLiquidationPriceSpread.toNumber())
            expect(result.red_stone_adapter_address).to.equal(actualRedStoneAdapterAddress)
            expect(result.red_stone_feed_id).to.equal(actualRedStoneFeedId)
            expect(String(result.position.size)).to.equal(actualPosition.size.toString())
            expect(result.position.open_notional).to.equal(actualPosition.openNotional.toNumber())
            expect(result.position.last_premium_fraction).to.equal(actualPosition.lastPremiumFraction.toNumber())


            // creating positions
            await removeAllAvailableMargin(charlie)
            await removeAllAvailableMargin(alice)
            let charlieBalance = _1e6.mul(150)
            await addMargin(charlie, charlieBalance)
            await addMargin(alice, charlieBalance)

            longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
            shortOrderBaseAssetQuantity = multiplySize(-0.1) // 0.1 ether
            orderPrice = multiplyPrice(1800)
            salt = BigNumber.from(Date.now())
            market = BigNumber.from(0)

            longOrder = getOrder(market, charlie.address, longOrderBaseAssetQuantity, orderPrice, salt, false)
            shortOrder = getOrder(market, alice.address, shortOrderBaseAssetQuantity, orderPrice, salt, false)
            await orderBook.connect(charlie).placeOrders([longOrder])
            await orderBook.connect(alice).placeOrders([shortOrder])
            await sleep(10)

            //testing for charlie
            response = await makehttpCall(method, params)
            result = response.body.result
            actualPosition = await amm.positions(charlie.address)
            expect(String(result.position.size)).to.equal(longOrderBaseAssetQuantity.toString())
            expect(result.position.open_notional).to.equal(longOrderBaseAssetQuantity.mul(orderPrice).div(_1e18).toNumber())
            expect(result.position.last_premium_fraction).to.equal(actualPosition.lastPremiumFraction.toNumber())

            // testing for alice
            params =[ammAddress, ammIndex, alice.address]
            response = await makehttpCall(method, params)
            actualPosition = await amm.positions(alice.address)
            expect(String(result.position.size)).to.equal(shortOrderBaseAssetQuantity.abs().toString())
            expect(result.position.open_notional).to.equal(shortOrderBaseAssetQuantity.mul(orderPrice).abs().div(_1e18).toNumber())
            expect(result.position.last_premium_fraction).to.equal(actualPosition.lastPremiumFraction.toNumber())

            //cleanup
            longOrder = getOrder(market, alice.address, longOrderBaseAssetQuantity, orderPrice, salt, false)
            shortOrder = getOrder(market, charlie.address, shortOrderBaseAssetQuantity, orderPrice, salt, false)
            await orderBook.connect(charlie).placeOrders([shortOrder])
            await orderBook.connect(alice).placeOrders([longOrder])
            await sleep(10)
            await removeAllAvailableMargin(charlie)
            await removeAllAvailableMargin(alice)
        })
    })

    context("IOC order contract variables", function () {
        it("should read the correct value from contracts", async function () {
            let charlieBalance = _1e6.mul(150)
            await addMargin(charlie, charlieBalance)

            longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
            orderPrice = multiplyPrice(1800)
            salt = BigNumber.from(Date.now())
            market = BigNumber.from(0)

            latestBlockNumber = await provider.getBlockNumber()
            lastTimestamp = (await provider.getBlock(latestBlockNumber)).timestamp
            expireAt = lastTimestamp + 6
            IOCOrder = getIOCOrder(expireAt, market, charlie.address, longOrderBaseAssetQuantity, orderPrice, salt, false)
            orderHash = await ioc.getOrderHash(IOCOrder)
            params = [ orderHash ]
            method ="testing_getIOCOrdersVars"

            // before placing order
            result = (await makehttpCall(method, params)).body.result

            actualExpirationCap = await ioc.expirationCap()
            expectedExpirationCap = result.ioc_expiration_cap

            expect(expectedExpirationCap).to.equal(actualExpirationCap.toNumber())
            expect(result.order_details.block_placed).to.eq(0)
            expect(result.order_details.filled_amount).to.eq(0)
            expect(result.order_details.order_status).to.eq(0)

            //placing order
            const tx = await ioc.connect(charlie).placeOrders([IOCOrder])
            const txReceipt = await tx.wait()
            result = (await makehttpCall(method, params)).body.result

            actualBlockPlaced = txReceipt.blockNumber
            expect(result.order_details.block_placed).to.eq(actualBlockPlaced)
            expect(result.order_details.filled_amount).to.eq(0)
            expect(result.order_details.order_status).to.eq(1)

            //cleanup
            await removeAllAvailableMargin(charlie)
        })
    })
    context("order book contract variables", function () {
        it("should read the correct value from contracts", async function () {

            await removeAllAvailableMargin(charlie)
            let charlieBalance = _1e6.mul(150)
            await addMargin(charlie, charlieBalance)

            longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
            orderPrice = multiplyPrice(1800)
            salt = BigNumber.from(Date.now())
            market = BigNumber.from(0)

            latestBlockNumber = await provider.getBlockNumber()
            lastTimestamp = (await provider.getBlock(latestBlockNumber)).timestamp
            expireAt = lastTimestamp + 6
            order = getOrder(market, charlie.address, longOrderBaseAssetQuantity, orderPrice, salt, false)
            orderHash = await orderBook.getOrderHash(order)
            params=[charlie.address, alice.address, orderHash]
            method ="testing_getOrderBookVars"

            // before placing order
            result = (await makehttpCall(method, params)).body.result

            actualResult = await orderBook.isTradingAuthority(charlie.address, alice.address)
            expect(result.is_trading_authority).to.equal(actualResult)

            expect(result.order_details.block_placed).to.eq(0)
            expect(result.order_details.filled_amount).to.eq(0)
            expect(result.order_details.order_status).to.eq(0)

            //placing order
            const tx = await orderBook.connect(charlie).placeOrders([order])
            const txReceipt = await tx.wait()
            result = (await makehttpCall(method, params)).body.result

            actualBlockPlaced = txReceipt.blockNumber
            expect(result.order_details.block_placed).to.eq(actualBlockPlaced)
            expect(result.order_details.filled_amount).to.eq(0)
            expect(result.order_details.order_status).to.eq(1)

            // cleanup
            await cancelOrderFromLimitOrder(order, charlie)
            await removeAllAvailableMargin(charlie)
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
