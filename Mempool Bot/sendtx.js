
const { ethers } = require('ethers')

async function sendTransaction(providerInfo) {
    try {
        const signer = new ethers.Wallet( // Create our signer
            providerInfo.acct,
            providerInfo.provider
        )
        await signer.sendTransaction(providerInfo.tx) // Send the transaction to the network
    } catch (error) {
        console.log('Exception occurred in sendTransaction')
    }
}

// eslint-disable-next-line no-undef
process.on('message', (providerInfo) => {
    const result = sendTransaction(providerInfo)
    // eslint-disable-next-line no-undef
    process.send(result)
})
