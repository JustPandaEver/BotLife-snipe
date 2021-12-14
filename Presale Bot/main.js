const ethers = require("ethers");
const config = require('./config.json')
const ethUtils = require('ethereumjs-util')

const provider = new ethers.providers.WebSocketProvider(config["provider"]);
const wallet = new ethers.Wallet(config["privateKey"], provider)

//////////////WATCHDOG////////////////////// --> Make sure we stay connected
const EXPECTED_PONG_BACK = 30000
const KEEP_ALIVE_CHECK_INTERVAL = 15000
////////////////////////////////////////////

async function storage(storageInt) {
    return parseInt(await provider.getStorageAt(config["presaleAddress"], storageInt), 16);
}

const wait = async () => {
  const startTime = await storage(106)
  const minBuy = (await storage(130)).toString()
  const maxBuy = (await storage(131)).toString()
  const hardCap = ethers.BigNumber.from((await storage(133)).toString())
  let contributedAmount = ethers.BigNumber.from(0);
  let nonce = (await provider.getTransactionCount('0x' + ethUtils.privateToAddress(Buffer.from(config["privateKey"].trim().toLowerCase(), 'hex')).toString('hex'))) // quite janky just make sure to use a fresh wallet
  while (startTime - unix() > 2) { 
    process.stdout.write(startTime - unix() + " Seconds remaining\n");
    await new Promise(r => setTimeout(r, 500));
  }
  provider.on('pending', async (txHash) => { // Look through the mempool
      provider.getTransaction(txHash).then(async (tx) => { // Get the transaction from the hash
          if (tx && tx.to) {
            if (tx.to === config["presaleAddress"] && tx.value.lte(maxBuy) && tx.value.gte(minBuy)) {
            contributedAmount = tx.value.add(contributedAmount)
            if (contributedAmount.gte(hardCap.mul(config["percentFill"]).div(100))) {
              let sendTx = { // Create the transaction
                value: ethers.BigNumber.from(config["amtBuy"]).mul(ethers.BigNumber.from(10).pow(18)).div(1000),
                nonce: nonce,
                gasPrice: (tx.gasPrice).add(1),
                gasLimit: (tx.gasLimit),
                chainId: parseInt(config["chainId"], 10),
                to: config["presaleAddress"],
                data: "0x"
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

function unix() {
    return parseInt(Math.floor(new Date().getTime() / 1000));
  }

wait()