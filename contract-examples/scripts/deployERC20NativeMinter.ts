import {
  Contract,
  ContractFactory
} from "ethers"
import { ethers } from "hardhat"

const main = async (): Promise<any> => {
  const Token: ContractFactory = await ethers.getContractFactory("ERC20NativeMinter")
  const token: Contract = await Token.deploy(5000)

  await token.deployed()
  console.log(`Token deployed to: ${token.address}`)
  const name: string = await token.name()
  console.log(`Name: ${name}`)

  const symbol: string = await token.symbol()
  console.log(`Symbol: ${symbol}`)

  const decimals: string = await token.decimals()
  console.log(`Decimals: ${decimals}`)
}

main()
  .then(() => process.exit(0))
  .catch(error => {
    console.error(error)
    process.exit(1)
  })
