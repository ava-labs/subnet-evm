import {
  Contract,
  ContractFactory
} from "ethers"
import { ethers } from "hardhat"

const rewardRate = ethers.utils.parseEther("0.00001")

const main = async (): Promise<any> => {
  const Contract: ContractFactory = await ethers.getContractFactory("ExampleRewardDistributor")
  const contract: Contract = await Contract.deploy(rewardRate)

  await contract.deployed()
  console.log(`Contract deployed to: ${contract.address}`)
}

main()
  .then(() => process.exit(0))
  .catch(error => {
    console.error(error)
    process.exit(1)
  })
