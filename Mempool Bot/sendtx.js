
const { ethers } = require('ethers')
const config = require('./config.json')

async function sendTransaction(providerInfo) {
    const provider = new ethers.providers.WebSocketProvider(config["RpcProvider"])
    console.log(await provider.sendTransaction(providerInfo.tx)) // Send the transaction to the network
}

// eslint-disable-next-line no-undef
process.on('message', (providerInfo) => {
    const result = sendTransaction(providerInfo)
    // eslint-disable-next-line no-undef
    process.send(result)
})

