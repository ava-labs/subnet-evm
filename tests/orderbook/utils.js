const { ethers, BigNumber } = require('ethers');

const _1e6 = BigNumber.from(10).pow(6)
const _1e12 = BigNumber.from(10).pow(12)
const _1e18 = BigNumber.from(10).pow(18)
const homedir = require('os').homedir()
let conf = require(`${homedir}/.hubblenet.json`)
const url = `http://127.0.0.1:9650/ext/bc/${conf.chain_id}/rpc`
provider = new ethers.providers.JsonRpcProvider(url);

// Set up signer
governance = new ethers.Wallet('0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80', provider) // governance
alice = new ethers.Wallet('0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d', provider); // 0x70997970c51812dc3a010c7d01b50e0d17dc79c8
bob = new ethers.Wallet('0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a', provider); // 0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc
charlie = new ethers.Wallet('15614556be13730e9e8d6eacc1603143e7b96987429df8726384c2ec4502ef6e', provider); // 0x55ee05df718f1a5c1441e76190eb1a19ee2c9430

const OrderBookContractAddress = "0x03000000000000000000000000000000000000b0"
const MarginAccountContractAddress = "0x03000000000000000000000000000000000000b1"
const ClearingHouseContractAddress = "0x03000000000000000000000000000000000000b2"
const JurorPrecompileAddress = "0x03000000000000000000000000000000000000a0"
const TicksPrecompileAddress = "0x03000000000000000000000000000000000000a1"
const LimitOrderBookContractAddress = "0x03000000000000000000000000000000000000b3"
const IOCContractAddress = "0x03000000000000000000000000000000000000b4"

orderBook = new ethers.Contract(OrderBookContractAddress, require('./abi/OrderBook.json'), provider);
limitOrderBook = new ethers.Contract(LimitOrderBookContractAddress, require('./abi/LimitOrderBook.json'), provider);
clearingHouse = new ethers.Contract(ClearingHouseContractAddress, require('./abi/ClearingHouse.json'), provider);
marginAccount = new ethers.Contract(MarginAccountContractAddress, require('./abi/MarginAccount.json'), provider);
ioc = new ethers.Contract(IOCContractAddress, require('./abi/IOC.json'), provider);
juror = new ethers.Contract(JurorPrecompileAddress, require('./abi/Juror.json'), provider);
ticks = new ethers.Contract(TicksPrecompileAddress, require('./abi/Ticks.json'), provider);

orderType = {
    Order: [
        // field ordering must be the same as LIMIT_ORDER_TYPEHASH
        { name: "trader", type: "address" },
        { name: "baseAssetQuantity", type: "int256" },
        { name: "price", type: "uint256" },
        { name: "salt", type: "uint256" },
    ]
}

function getOrder(market, traderAddress, baseAssetQuantity, price, salt, reduceOnly=false) {
    return {
        ammIndex: market,
        trader: traderAddress,
        baseAssetQuantity: baseAssetQuantity,
        price: price,
        salt: BigNumber.from(salt),
        reduceOnly: reduceOnly,
    }
}

function getOrderV2(ammIndex, trader, baseAssetQuantity, price, salt, reduceOnly=false, postOnly=false) {
    return {
        ammIndex,
        trader,
        baseAssetQuantity,
        price,
        salt: BigNumber.from(salt),
        reduceOnly,
        postOnly
    }
}

function getIOCOrder(expireAt, ammIndex, trader, baseAssetQuantity, price, salt, reduceOnly=false) {
    return {
        orderType: 1,
        expireAt: expireAt,
        ammIndex: ammIndex,
        trader: trader,
        baseAssetQuantity: baseAssetQuantity,
        price: price,
        salt: salt,
        reduceOnly: false
    }
}

//Convert to wei units to support 18 decimals
function multiplySize(size) {
    // return _1e18.mul(size)
    return ethers.utils.parseEther(size.toString())
}

function multiplyPrice(price) {
    return _1e6.mul(price)
    // return ethers.utils.parseUnits(price.toString(), 6)
}

async function getDomain() {
    domain = {
        name: "Hubble",
        version: "2.0",
        chainId: (await provider.getNetwork()).chainId,
        verifyingContract: orderBook.address
    }
    return domain
}

async function placeOrder(market, trader, size, price, salt=Date.now(), reduceOnly=false) {
    order = getOrder(market, trader.address, size, price, salt, reduceOnly)
    return placeOrderFromLimitOrder(order, trader)
}

async function placeOrderFromLimitOrder(order, trader) {
    const tx = await limitOrderBook.connect(trader).placeOrders([order])
    const txReceipt = await tx.wait()
    return { tx, txReceipt }
}

async function placeOrderFromLimitOrderV2(order, trader) {
    // console.log({ placeOrderEstimateGas: (await limitOrderBook.connect(trader).estimateGas.placeOrders([order])).toNumber() })
    // return limitOrderBook.connect(trader).placeOrders([order])
    const tx = await limitOrderBook.connect(trader).placeOrders([order])
    const txReceipt = await tx.wait()
    return { tx, txReceipt }
}

async function placeV2Orders(orders, trader) {
    console.log({ placeOrdersEstimateGas: (await limitOrderBook.connect(trader).estimateGas.placeOrders(orders)).toNumber() })
    const tx = await limitOrderBook.connect(trader).placeOrders(orders)
    const txReceipt = await tx.wait()
    return { tx, txReceipt }
}

async function placeIOCOrder(order, trader) {
    const tx = await ioc.connect(trader).placeOrders([order])
    const txReceipt = await tx.wait()
    return { tx, txReceipt }
}

async function cancelOrderFromLimitOrder(order, trader) {
    const tx = await limitOrderBook.connect(trader).cancelOrder(order)
    const txReceipt = await tx.wait()
    return { tx, txReceipt }
}

async function cancelOrderFromLimitOrderV2(order, trader) {
    // console.log({ estimateGas: (await limitOrderBook.connect(trader).estimateGas.cancelOrders([order])).toNumber() })
    // return limitOrderBook.connect(trader).cancelOrders([order])
    const tx = await limitOrderBook.connect(trader).cancelOrders([order])
    const txReceipt = await tx.wait()
    return { tx, txReceipt }
}

async function cancelV2Orders(orders, trader) {
    console.log({ cancelV2OrdersEstimateGas: (await limitOrderBook.connect(trader).estimateGas.cancelOrders(orders)).toNumber() })
    const tx = await limitOrderBook.connect(trader).cancelOrders(orders)
    const txReceipt = await tx.wait()
    return { tx, txReceipt }
}

function sleep(s) {
    return new Promise(resolve => setTimeout(resolve, s * 1000));
}

async function addMargin(trader, amount, txOpts={}) {
    const hgtAmount = _1e12.mul(amount)
    marginAccountHelper = await getMarginAccountHelper()
    const tx = await marginAccountHelper.connect(trader).addVUSDMarginWithReserve(amount, trader.address, Object.assign(txOpts, { value: hgtAmount }))
    const result = await marginAccount.marginAccountHelper()
    const txReceipt = await tx.wait()
    return { tx, txReceipt }
}

async function removeMargin(trader, amount) {
    const hgtAmount = _1e12.mul(amount)
    marginAccountHelper = await getMarginAccountHelper()
    const tx = await marginAccountHelper.connect(trader).removeMarginInUSD(hgtAmount)
    const txReceipt = await tx.wait()
    return { tx, txReceipt }
}

async function removeAllAvailableMargin(trader) {
    margin = await marginAccount.getAvailableMargin(trader.address)
    marginAccountHelper = await getMarginAccountHelper()
    if (margin.toNumber() > 0) {
        // const tx = await marginAccountHelper.connect(trader).removeMarginInUSD(5e11)
        const tx = await marginAccountHelper.connect(trader).removeMarginInUSD(margin.toNumber())
        await tx.wait()
    }
    return
}

async function getMarginAccountHelper() {
    marginAccountHelperAddress = await marginAccount.marginAccountHelper()
    return new ethers.Contract(marginAccountHelperAddress, require('./abi/MarginAccountHelper.json'), provider)
}

function encodeLimitOrder(order) {
    const encodedOrder = ethers.utils.defaultAbiCoder.encode(
        [
          'uint256',
          'address',
          'int256',
          'uint256',
          'uint256',
          'bool',
        ],
        [
            order.ammIndex,
            order.trader,
            order.baseAssetQuantity,
            order.price,
            order.salt,
            order.reduceOnly,
        ]
    )
    return encodedOrder
}

function encodeLimitOrderWithType(order) {
    encodedOrder = encodeLimitOrder(order)
    const typedEncodedOrder = ethers.utils.defaultAbiCoder.encode(['uint8', 'bytes'], [0, encodedOrder])
    return typedEncodedOrder
}

function encodeLimitOrderV2(order) {
    const encodedOrder = ethers.utils.defaultAbiCoder.encode(
        [
          'uint256',
          'address',
          'int256',
          'uint256',
          'uint256',
          'bool',
          'bool',
        ],
        [
            order.ammIndex,
            order.trader,
            order.baseAssetQuantity,
            order.price,
            order.salt,
            order.reduceOnly,
            order.postOnly,
        ]
    )
    return encodedOrder
}

function encodeLimitOrderV2WithType(order) {
    encodedOrder = encodeLimitOrderV2(order)
    const typedEncodedOrder = ethers.utils.defaultAbiCoder.encode(['uint8', 'bytes'], [0, encodedOrder])
    return typedEncodedOrder
}

// async function cleanUpPositionsAndRemoveMargin(market, trader1, trader2) {
//     position1 = await amm.positions(trader1.address)
//     position2 = await amm.positions(trader2.address)
//     if (position1.size.toString() != "0" && position2.size.toString() != "0") {
//         if (position1.size.toString() != positionSize2.size.toString()) {
//             console.log("Position sizes are not equal")
//             return
//         }
//         price = BigNumber.from(position1.notionalPosition.toString()).div(BigNumber.from(position1.size.toString()))
//         console.log("placing opposite orders to close positions")
//         await placeOrder(market, trader1, positionSize1, price)
//         await placeOrder(market, trader2, positionSize2, price)
//         await sleep(10)
//     }

//     console.log("removing margin for both traders")
//     await removeAllAvailableMargin(trader1)
//     await removeAllAvailableMargin(trader2)
// }

function getRandomSalt() {
    // return date and add random number to generate unique salts even when called concurrently(hopefully)
    // Return a random number between 0 and 1000:
    randomNumber = Date.now() + Math.floor(Math.random()*1000)
    return BigNumber.from(randomNumber)
}

async function waitForOrdersToMatch() {
    await sleep(5)
}

async function enableValidatorMatching() {
    const tx = await orderBook.connect(governance).setValidatorStatus(ethers.utils.getAddress('0x4Cf2eD3665F6bFA95cE6A11CFDb7A2EF5FC1C7E4'), true)
    await tx.wait()
}

async function disableValidatorMatching() {
    const tx = await orderBook.connect(governance).setValidatorStatus(ethers.utils.getAddress('0x4Cf2eD3665F6bFA95cE6A11CFDb7A2EF5FC1C7E4'), false)
    await tx.wait()
}

async function getAMMContract(market) {
    const ammAddress = await clearingHouse.amms(market)
    amm =  new ethers.Contract(ammAddress, require("./abi/AMM.json"), provider);
    return amm
}

async function getMinSizeRequirement(market) {
    const amm = await getAMMContract(market)
    return await amm.minSizeRequirement()
}

async function getMakerFee() {
    return await clearingHouse.makerFee()
}

async function getTakerFee() {
    return await clearingHouse.takerFee()
}

async function getOrderBookEvents(fromBlock=0) {
    block = await provider.getBlock("latest")
    events = await orderBook.queryFilter("*",fromBlock,block.number)
    return events
}

async function getLimitOrderBookEvents(fromBlock=0) {
    block = await provider.getBlock("latest")
    events = await limitOrderBook.queryFilter("*",fromBlock,block.number)
    return events
}

function bnToFloat(num, decimals = 6) {
    return parseFloat(ethers.utils.formatUnits(num.toString(), decimals))
}

async function getRequiredMarginForLongOrder(longOrder) {
    price = longOrder.price
    baseAssetQuantity = longOrder.baseAssetQuantity

    minAllowableMargin = await clearingHouse.minAllowableMargin()
    takerFee = await clearingHouse.takerFee()

    quoteAsset = baseAssetQuantity.mul(price).div(_1e18).abs()
    requiredMargin = quoteAsset.mul(minAllowableMargin).div(_1e6)
    requiredTakerFee = quoteAsset.mul(takerFee).div(_1e6)
    totalRequiredMargin = requiredMargin.add(requiredTakerFee)
    return totalRequiredMargin
}

async function getRequiredMarginForShortOrder(shortOrder) {
    price = shortOrder.price
    baseAssetQuantity = shortOrder.baseAssetQuantity

    minAllowableMargin = await clearingHouse.minAllowableMargin()
    takerFee = await clearingHouse.takerFee()
    amm = await getAMMContract(shortOrder.ammIndex)
    oraclePrice = await amm.getUnderlyingPrice()
    maxOracleSpreadRatio = await amm.maxOracleSpreadRatio()
    upperBound = oraclePrice.mul(maxOracleSpreadRatio.add(_1e6)).div(_1e6)
    if (price < upperBound) {
        price = upperBound
    }
    // for shortOrder we use upperBound to reservePrice as it is the worst price
    quoteAsset = baseAssetQuantity.mul(price).div(_1e18).abs()
    requiredMargin = quoteAsset.mul(minAllowableMargin).div(_1e6)
    requiredTakerFee = quoteAsset.mul(takerFee).div(_1e6)
    return requiredMargin.add(requiredTakerFee)
}

async function getEventsFromOrderBookTx(transactionHash) {
    tx = await provider.getTransaction(transactionHash)
    events = await getOrderBookEvents(tx.blockNumber)
    var orderBookLogsWithEvent = []
    for(i = 0; i < events.length; i++) {
        if(events[i].transactionHash == transactionHash) {
            orderBookLogsWithEvent.push(events[i])
            break
        }
    }
    return orderBookLogsWithEvent
}

async function getEventsFromLimitOrderBookTx(transactionHash) {
    tx = await provider.getTransaction(transactionHash)
    events = await getLimitOrderBookEvents(tx.blockNumber)
    var limitOrderBookLogsWithEvent = []
    for(i = 0; i < events.length; i++) {
        if(events[i].transactionHash == transactionHash) {
            limitOrderBookLogsWithEvent.push(events[i])
            break
        }
    }
    return limitOrderBookLogsWithEvent
}

module.exports = {
    _1e6,
    _1e12,
    _1e18,
    addMargin,
    alice,
    bnToFloat,
    bob,
    cancelOrderFromLimitOrder,
    cancelOrderFromLimitOrderV2,
    cancelV2Orders,
    charlie,
    clearingHouse,
    disableValidatorMatching,
    enableValidatorMatching,
    encodeLimitOrder,
    encodeLimitOrderWithType,
    encodeLimitOrderV2,
    encodeLimitOrderV2WithType,
    getAMMContract,
    getDomain,
    getEventsFromOrderBookTx,
    getEventsFromLimitOrderBookTx,
    getIOCOrder,
    getOrder,
    getOrderV2,
    getMakerFee,
    getMinSizeRequirement,
    getOrderBookEvents,
    getLimitOrderBookEvents,
    getRandomSalt,
    getRequiredMarginForLongOrder,
    getRequiredMarginForShortOrder,
    getTakerFee,
    governance,
    ioc,
    juror,
    limitOrderBook,
    marginAccount,
    multiplySize,
    multiplyPrice,
    orderBook,
    orderType,
    provider,
    placeOrder,
    placeV2Orders,
    placeOrderFromLimitOrder,
    placeOrderFromLimitOrderV2,
    placeIOCOrder,
    removeAllAvailableMargin,
    removeMargin,
    sleep,
    ticks,
    url,
    waitForOrdersToMatch,
}
