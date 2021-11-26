/* eslint-disable no-undef */
const { ethers, BigNumber } = require('ethers')

async function sendTransaction(providerInfo) {
    try {
        const signer = new ethers.Wallet(
            providerInfo.acct,
            providerInfo.provider
        )
        await signer.sendTransaction(providerInfo.tx)
        console.log(object)
        const balance = await provider.getBalance(providerInfo.contract)
        if (balance.eq(BigNumber(0))) {
            return true
        }
        return false
    } catch (error) {
        console.log('Exception occurred in sendTransaction')
    }
}

process.on('message', (providerInfo) => {
    const result = sendTransaction(providerInfo)
    process.send(result)
})
