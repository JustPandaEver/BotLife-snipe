const ethers = require("ethers");
const config = require('./config.json')
const ethUtils = require('ethereumjs-util')
const approveToken = require('./approve.js')


//////////////////INIT ETHERS//////////////////////////
const provider = new ethers.providers.WebSocketProvider(config["provider"]);
const wallet = new ethers.Wallet(config["privateKey"], provider)
///////////////////////////////////////////////////////

////////////////////RETRIEVE DATA//////////////////////
async function storage(storageInt) {
    return await provider.getStorageAt(config["presaleAddress"], storageInt);
}
///////////////////////////////////////////////////////

//////////////////LOCATIONS////////////////////////////
const storageAddress = { // Set location in storage for the data
  pinkSale : {
    startTime : 106,
    minBuy : 130,
    maxBuy : 131,
    hardCap : 133,
    tokenAddress : 105,
    startBlock: 0
  },
  dxSale : {
    startTime : 45,
    minBuy : 26,
    maxBuy : 27,
    hardCap : 44,
    tokenAddress : 0,
    startBlock: 0
  },
  unicrypt : {
    startTime : 0,
    minBuy : 0,
    maxBuy : 5,
    hardCap : 7,
    tokenAddress : 2,
    startBlock: 11
  }
}
///////////////////////////////////////////////////////

/////////////////////MAIN//////////////////////////////
const wait = async () => {


  async function data(data) { // Get the important data from contract
    startTime = parseInt(await storage(data.startTime), 16)
    minBuy = ethers.BigNumber.from(parseInt((await storage(data.minBuy)), 16).toString())
    maxBuy = ethers.BigNumber.from((await storage(data.maxBuy)).toString())
    hardCap = ethers.BigNumber.from(parseInt((await storage(data.hardCap)), 16).toString())
    tokenAddress = "0x" + (await storage(data.tokenAddress)).toString().slice(26, 66)
    startBlock = parseInt(await storage(data.startBlock), 16)
  }

  let zero = ethers.BigNumber.from(0)
  let contributedAmount = zero; // Initialize contributed amount at 0
  let nonce = (await provider.getTransactionCount('0x' + ethUtils.privateToAddress(Buffer.from(config["privateKey"].trim().toLowerCase(), 'hex')).toString('hex'))) // quite janky just make sure to use a fresh wallet
  

  ///////////////CHOOSE PLATFORM////////////////////////
  if (config["action"] == "pinkSale") {
    await data(storageAddress.pinkSale)
  } else if (config["action"] == "dxSale") {
    await data(storageAddress.dxSale)
    hardCap = hardCap.mul(ethers.BigNumber.from(10).pow(19))
  } else if (config["action"] == "unicrypt") {
    await data(storageAddress.unicrypt)
    minBuy = zero.add(1)
  }
  ///////////////////////////////////////////////////////

  approveToken(tokenAddress) // Approve

  /////////////////////WAIT//////////////////////////////
  if (config["action"] == "pinkSale" || config["action"] == "dxSale") {
    while (startTime - unix() > 4) { 
      process.stdout.write(startTime - unix() + " Seconds remaining\n");
      await new Promise(r => setTimeout(r, 500));
    }
  }
  if (config["action"] == "unicrypt") {
    while (await provider.getBlockNumber() < (startBlock - 2)) {
      console.log("Current Block:" + await provider.getBlockNumber());
      await new Promise(r => setTimeout(r, 1000));
    }
  }
  ///////////////////////////////////////////////////////

  ////////////////////////BUY////////////////////////////
  provider.on('pending', async (txHash) => { // Look through the mempool
      provider.getTransaction(txHash).then(async (tx) => { // Get the transaction from the hash
          if (tx && tx.to) {
            if (tx.to == config["presaleAddress"] && tx.value.lte(maxBuy) && tx.value.gte(minBuy)) {
            contributedAmount = tx.value.add(contributedAmount)
            if (contributedAmount.gte(hardCap.mul(config["percentFill"]).div(100))) {
              if (config["action"] == "unicrypt") {
                txData = "0xf868e7660000000000000000000000000000000000000000000000000000000000000000"
              } else {
                txData = "0x"
              }
              let sendTx = { // Create the transaction
                value: ethers.BigNumber.from(config["amtBuy"]).mul(ethers.BigNumber.from(10).pow(18)),
                nonce: nonce,
                gasPrice: (tx.gasPrice).add(1),
                gasLimit: (tx.gasLimit),
                chainId: parseInt(config["chainId"], 10),
                to: config["presaleAddress"],
                data: txData
              }
              let signedTx = await wallet.signTransaction(sendTx) // Sign the transaction
              let receipt = await provider.sendTransaction(signedTx) // Send the transaction to the network
              console.log("Buy receipt: " + receipt.hash + "\n");
              process.exit(1)
            }
          }
          }
      })
  })
}
///////////////////////////////////////////////////////

/////////////////////TIME//////////////////////////////
function unix() {
    return parseInt(Math.floor(new Date().getTime() / 1000));
  }
///////////////////////////////////////////////////////

wait()