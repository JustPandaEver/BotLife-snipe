// inputs --> target, multiplier, percent, amountIn
/* eslint-disable no-empty */
const ethers = require('ethers')
const config = require('./config.json')

//const account = new ethers.Wallet(process.env.MNEMONIC, provider);
const masterKey = config["masterKey"]
const provider = new ethers.providers.WebSocketProvider(config["RpcProvider"])
const walletBot = new ethers.Wallet(masterKey, provider) // masterKey needs to be the wallet coins are sent to from the contract
const PCSRouter = config["PCSRouter"]
const wBNB = config["wBNB"]
const pcsAbi = new ethers.utils.Interface(require('./abi.json'))
const router = new ethers.Contract(PCSRouter, pcsAbi, walletBot)


///////////////FIND BALANCE////////////////// --> Find how many tokens we received
async function getTokenBalance(sellContract) {
    return (await sellContract.balanceOf(walletBot.address))
}
/////////////////////////////////////////////

/////////////////////SELL//////////////////// --> Sell our tokens
async function tradeTokensForExactETHWithSupportingFee(
    tokenToSell, // --> Address of token to sell
    bigNumberToSell, // --> Number of tokens to sell
    percent // --> Percent of tokens to sell
) {
    try {
        console.log(`Selling Token: ${tokenToSell}`)
        const amountOutMin = ethers.BigNumber.from(0)
        let a = ethers.BigNumber.from(bigNumberToSell)
        let sellGasLimit = config["sellGasLimit"]
        let sellGasPrice = config["sellGasPrice"]
        const tx = // --> Create the transaction
            await router.swapExactTokensForETHSupportingFeeOnTransferTokens(
                a.mul(percent).div(100), // Amount of tokens to sell
                amountOutMin,
                [tokenToSell, wBNB], // Path
                walletBot.address, // To
                Date.now() + 1000 * 60 * 10, // Deadline --> 10 minutes
                {
                    gasLimit: sellGasLimit,
                    gasPrice: sellGasPrice,
                }
            )
        const receipt = await tx.wait() // --> Get the transaction hash
        console.log('Sell Transaction Receipt:')
        console.log(receipt.transactionHash)
        // eslint-disable-next-line no-undef
        process.exit() // Kill the bot
    } catch (e) {
        console.log("The sale failed!");
    }
}
/////////////////////////////////////////////


/////////////////////MAIN////////////////////
module.exports = async function sell(target, multiplier, percent, amountIn) {
    const toSnipe = ethers.utils.getAddress(target)
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
    let balance = getTokenBalance(sellContract)
    // eslint-disable-next-line no-constant-condition
    while (true) {
        let amountOut = router.getAmountsOut(
            balance,
            [toSnipe, wBNB]
        )
        let goal = amountIn.mul(multiplier)
        if (amountOut.gt(goal)) {
            await tradeTokensForExactETHWithSupportingFee(toSnipe, balance, percent)
        }
    }
}
/////////////////////////////////////////////