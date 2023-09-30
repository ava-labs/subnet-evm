import {
    Contract,
    ContractFactory
  } from "ethers"
  import { ethers } from "hardhat"
  
  const main = async (): Promise<any> => {
    const contractFactory: ContractFactory = await ethers.getContractFactory("OrderBook")
    contractFactory.attach
    const contract: Contract = await contractFactory.deploy('orderBook', 1)
  
    await contract.deployed()
    console.log(`OrderBook Contract deployed to: ${contract.address}`)
  }
  
  main()
    .then(() => process.exit(0))
    .catch(error => {
      console.error(error)
      process.exit(1)
    })
  