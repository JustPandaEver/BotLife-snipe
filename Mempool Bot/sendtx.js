
const { ethers, BigNumber } = require('ethers')

async function sendTransaction(providerInfo) {
    try {
        const signer = new ethers.Wallet(
            providerInfo.acct,
            providerInfo.provider
        )
        await signer.sendTransaction(providerInfo.tx)
        const balance = await providerInfo.provider.getBalance(providerInfo.contract)
        if (balance.eq(BigNumber(0))) {
            return true
        }
        return false
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
