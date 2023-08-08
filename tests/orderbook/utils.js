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

// Set up contract interface
const OrderBookContractAddress = "0x0300000000000000000000000000000000000000"
const MarginAccountContractAddress = "0x0300000000000000000000000000000000000001"
const ClearingHouseContractAddress = "0x0300000000000000000000000000000000000002"
const HubbleBibliophilePrecompileAddress = "0x0300000000000000000000000000000000000004"
const IOCContractAddress = "0x635c5F96989a4226953FE6361f12B96c5d50289b"

orderBook = new ethers.Contract(OrderBookContractAddress, require('./abi/OrderBook.json'), provider);
clearingHouse = new ethers.Contract(ClearingHouseContractAddress, require('./abi/ClearingHouse.json'), provider);
marginAccount = new ethers.Contract(MarginAccountContractAddress, require('./abi/MarginAccount.json'), provider);
hubblebibliophile = new ethers.Contract(HubbleBibliophilePrecompileAddress, require('./abi/MarginAccount.json'), provider)
ioc = new ethers.Contract(IOCContractAddress, require('./abi/IOC.json'), provider);

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
    return ethers.utils.parseEther(size.toString())
}

function multiplyPrice(price) {
    return ethers.utils.parseUnits(price.toString(), 6)
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
    const order = {
        ammIndex: market,
        trader: trader.address,
        baseAssetQuantity: size,
        price: price,
        salt: salt,
        reduceOnly: reduceOnly,
    }
    const tx = await orderBook.connect(trader).placeOrder(order)
    const txReceipt = await tx.wait()
    return { tx, txReceipt }
}

async function cancelOrder(market, trader, size, price, salt=Date.now(), reduceOnly=false) {
    const order = {
        ammIndex: market,
        trader: trader.address,
        baseAssetQuantity: size,
        price: price,
        salt: salt,
        reduceOnly: reduceOnly,
    }
    const tx = await orderBook.connect(trader).cancelOrder(order)
    const txReceipt = await tx.wait()
    return { tx, txReceipt }
}

async function cancelOrderFromLimitOrder(order, trader) {
    const tx = await orderBook.connect(trader).cancelOrder(order)
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
    if (margin.toNumber() != 0) {
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
    const typedEncodedOrder = ethers.utils.defaultAbiCoder.encode(['uint8', 'bytes'], [0, encodedOrder])
    // console.log({ order, encodedOrder, typedEncodedOrder })
    return typedEncodedOrder
}

function encodeIOCOrder(order) {
    const encodedOrder = ethers.utils.defaultAbiCoder.encode(
        [
          'uint8',
          'uint256',
          'uint256',
          'address',
          'int256',
          'uint256',
          'uint256',
          'bool',
        ],
        [
            order.orderType,
            order.expireAt,
            order.ammIndex,
            order.trader,
            order.baseAssetQuantity,
            order.price,
            order.salt,
            order.reduceOnly,
        ]
    )
    const typedEncodedOrder = ethers.utils.defaultAbiCoder.encode(['uint8', 'bytes'], [1, encodedOrder])
    // console.log({ order, encodedOrder, typedEncodedOrder })
    return typedEncodedOrder
}

async function enableValidatorMatching() {
    const tx = await orderBook.connect(governance).setValidatorStatus(ethers.utils.getAddress('0x4Cf2eD3665F6bFA95cE6A11CFDb7A2EF5FC1C7E4'), true)
    await tx.wait()
}

async function disableValidatorMatching() {
    const tx = await orderBook.connect(governance).setValidatorStatus(ethers.utils.getAddress('0x4Cf2eD3665F6bFA95cE6A11CFDb7A2EF5FC1C7E4'), false)
    await tx.wait()
}

module.exports = {
    _1e6,
    _1e12,
    _1e18,
    addMargin,
    alice,
    bob,
    cancelOrder,
    cancelOrderFromLimitOrder,
    charlie,
    clearingHouse,
    encodeIOCOrder,
    disableValidatorMatching,
    enableValidatorMatching,
    encodeLimitOrder,
    getDomain,
    getOrder,
    getIOCOrder,
    governance,
    hubblebibliophile,
    ioc, 
    marginAccount,
    multiplySize,
    multiplyPrice,
    orderBook,
    orderType,
    provider,
    placeOrder,
    removeAllAvailableMargin,
    removeMargin,
    sleep,
    url,
}
