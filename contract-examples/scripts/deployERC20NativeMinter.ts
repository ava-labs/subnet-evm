import {
  Contract,
  ContractFactory
} from "ethers"
import { ethers } from "hardhat"

const main = async (): Promise<any> => {
  const Token: ContractFactory = await ethers.getContractFactory("ERC20NativeMinter")
  const token: Contract = await Token.deploy(5000)

  await token.deployed()
}

main()
  .then(() => process.exit(0))
  .catch(error => {
    console.error(error)
    process.exit(1)
  })
