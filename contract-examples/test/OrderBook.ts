
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { expect } from "chai";
import exp = require("constants");
import {
  BigNumber,
  Contract,
  ContractFactory,
} from "ethers"
import { ethers } from "hardhat"
import { EmitHint } from "typescript";

// make sure this is always an admin for minter precompile
const adminAddress: string = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"

describe.only('Order Book', function() {
    let orderBook, alice, order, domain, orderType, signature

    before(async function () {
        // const owner = await ethers.getSigner(adminAddress);
        const OrderBook = await ethers.getContractFactory('OrderBook')
        orderBook = await OrderBook.deploy('Hubble', '2.0')
        const signers = await ethers.getSigners()
        alice = signers[1]
    })

    it('verify signer', async function() {
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

        orderType =  {
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

        expect(await orderBook.getOrdersLen()).to.eq(1)
        const contractOrder = await orderBook.orders(0)
        expect(contractOrder[0]).to.eq(order.trader)
        expect(contractOrder[1]).to.eq(order.baseAssetQuantity)
        expect(contractOrder[2]).to.eq(order.price)
        expect(contractOrder[3]).to.eq(order.salt)
    })
})
