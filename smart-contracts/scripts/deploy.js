const { ethers } = require("hardhat");

async function main() {
  console.log("Deploying LaunchpadFactory contract...");

  // Get the ContractFactory and Signers here.
  const LaunchpadFactory = await ethers.getContractFactory("LaunchpadFactory");
  
  // Deploy the contract
  const factory = await LaunchpadFactory.deploy();
  await factory.waitForDeployment();

  const factoryAddress = await factory.getAddress();
  console.log("LaunchpadFactory deployed to:", factoryAddress);

  // Save the deployment info
  const deploymentInfo = {
    network: network.name,
    factoryAddress: factoryAddress,
    deployedAt: new Date().toISOString(),
    blockNumber: await ethers.provider.getBlockNumber(),
  };

  console.log("Deployment info:", deploymentInfo);

  return deploymentInfo;
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
if (require.main === module) {
  main()
    .then(() => process.exit(0))
    .catch((error) => {
      console.error(error);
      process.exit(1);
    });
}

module.exports = main;