/* eslint-disable no-empty */
const ethers = require('ethers')

let addressToSell;
//const account = new ethers.Wallet(process.env.MNEMONIC, provider);
const walletBot = new ethers.Wallet(masterKey, provider)
const provider = new ethers.providers.WebSocketProvider('ws://localhost:8546')
const toSnipe = ethers.utils.getAddress(addressToSell)
const PCSRouter = '0x10ED43C718714eb63d5aA57B78B54704E256024E'
const wBNB = '0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c'
const pcsAbi = new ethers.utils.Interface(require('./conf/abi.json'))
const router = new ethers.Contract(PCSRouter, pcsAbi, walletBot)
const masterKey = ''

// Sell at X gain?
/////////////SELL CONTRACT////////////////// --> Connect to a contract
const sellContract = new ethers.Contract( 
    toSnipe, // Address to connect to
    [ // ABI of contract
        'function approve(address _spender, uint256 _value) public returns (bool success)',
        'function balanceOf(address account) external view returns (uint256)',
        'function decimals() view returns (uint8)',
        'function allowance(address owner, address spender) external view returns (uint)',
    ],
    provider // Provider or Signer
)
/////////////////////////////////////////////


//////////////SLEEP////////////////////////// --> Pause before selling
function sleep(ms) {
    return new Promise((resolve) => {
        setTimeout(resolve, ms)
    })
}
/////////////////////////////////////////////

///////////////FIND BALANCE////////////////// --> Find how many tokens we received and sell
async function getTokenBalanceAndTrade(tokenToSell) {
    sellContract.balanceOf(walletBot.address).then(async (balanceBigNumber) => { // --> Check balance of bot
        tradeTokensForExactETHWithSupportingFee(tokenToSell, balanceBigNumber) // --> Sell
    })
}
/////////////////////////////////////////////

/////////////////////SELL//////////////////// --> Sell our tokens
module.exports = async function tradeTokensForExactETHWithSupportingFee(
    tokenToSell, // --> Address of token to sell
    bigNumberToSell, // --> Number of tokens to sell
    sleepTime
) {
    try {
        sleep(sleepTime * 1000) // Wait x time before selling
        console.log(`Selling Token: ${tokenToSell}`)
        const amountOutMin = ethers.BigNumber.from(0)
        let a = ethers.BigNumber.from(bigNumberToSell)
        const tx = // --> Create the transaction
            await router.swapExactTokensForETHSupportingFeeOnTransferTokens(
                a, // Amount of tokens to sell. What was sellTrue?
                amountOutMin,
                [tokenToSell, wBNB], // Path
                walletBot.address, // To
                Date.now() + 1000 * 60 * 1, // Deadline --> 10 minutes
                {
                    gasLimit: 10000000,
                    gasPrice: 15000000000,
                }
            )
        const receipt = await tx.wait() // --> Get the transaction hash
        console.log('Sell Transaction Receipt:')
        console.log(receipt.transactionHash)
        // eslint-disable-next-line no-undef
        process.exit() // Kill the bot
    } catch (e) {
        getTokenBalanceAndTrade(tokenToSell) // If the sale fails, retry
    }
}
/////////////////////////////////////////////