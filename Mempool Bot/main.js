#!/usr/bin/env -S node --max-old-space-size=100000
/* eslint-disable no-empty */
//inputs (args) --> snipeType, target, method, amountOut, X, %, tries, amountIn, blockDelay


////////////////IMPORTS//////////////////// 
'use strict'
const { exec } = require('child_process')
// eslint-disable-next-line no-undef
var myArgs = process.argv.slice(2)
const ethers = require('ethers')
const spamBot = require('./buy.js')
const approveToken = require('./approve.js')
const sell = require('./sell.js')
//const yargs = require('yargs/yargs')
//const { hideBin } = require('yargs/helpers')
//const argv = yargs(hideBin(process.argv)).argv
//Object.assign(process.env, env)
//const env = require('./conf/env.json')
///////////////////////////////////////////


//////////////////CONSTANTS////////////////
const provider = new ethers.providers.WebSocketProvider('ws://localhost:8546')
const toSnipe = ethers.utils.getAddress(String(myArgs[1]))
const snipeType = String(myArgs[0])
const PCSRouter = '0x10ED43C718714eb63d5aA57B78B54704E256024E'
const pcsAbi = new ethers.utils.Interface(require('./conf/abi.json'))
//const blockDelay = parseInt(String(myArgs[2]))
//const contractAdd = ethers.utils.getAddress('0xFDFb632bE900583a0fBa25b72cE6A6763a70Fbe8')
//const snipeData = '0xf088d547' + ethers.utils.hexZeroPad(toSnipe, 32).replace('0x', '')
///////////////////////////////////////////



///////////////////REGEX/////////////////// --> Choose contract method to look for
//DXSale
const reDX = new RegExp('^0x267dd102')
//addLiquidity
const re1 = new RegExp('^0xf305d719')
const re2 = new RegExp('^0xe8e33700')
//UniCrypt
const re3 = new RegExp('^0xe8078d94')
const re5 = new RegExp('^0x4d536fe3')
const re6 = new RegExp('^0xa6f2ae3a')
const re7 = new RegExp('^0xd98a9469')
//tradingEnabled
const trading1 = new RegExp('^' + myArgs[2])
//PinkSale
const re8 = new RegExp('^0x4bb278f3')
//const re4 = new RegExp("^0x8119c065");
////////////////////////////////////////////


//////////////WATCHDOG////////////////////// --> Make sure we stay connected
const EXPECTED_PONG_BACK = 30000
const KEEP_ALIVE_CHECK_INTERVAL = 15000
////////////////////////////////////////////


////////////////////CHOOSE METHOD////////////
const startConnection = () => {


    console.log(
        '1 = DX \n2 = AddLiq \n3 = Uni \n4 = TradingEnabled SnipeID Address MethodIDofTradingEnabled'
    )
    console.log('--blockDelay # blocks to delay buy after liquidity add')
    console.log('Waiting for Transaction Data...\n')


    // eslint-disable-next-line no-unused-vars
    exec('./clean', (err, stdout, stderr) => {
        console.log(stdout)
    })
    let pingTimeout = null
    let keepAliveInterval = null


    //Open the WSS
    provider._websocket.on('open', () => {
        keepAliveInterval = setInterval(() => {
            provider._websocket.ping()
            pingTimeout = setTimeout(() => {
                provider._websocket.terminate()
            }, EXPECTED_PONG_BACK)
        }, KEEP_ALIVE_CHECK_INTERVAL)


        // Approve the token
        approveToken(toSnipe) 

        let amountOut = myArgs[3]
        let tries = myArgs[6]
        let multiplier = myArgs[4]
        let percent = myArgs[5]
        let amountIn = myArgs[7]
        let blockDelay = myArgs[8]
        provider.on('pending', async (txHash) => { // Look through the mempool
            provider.getTransaction(txHash).then(async (tx) => { // Get the transaction from the hash
                if (tx && tx.to) { // Check the transaction exists and it has a to address
                    switch (snipeType) { // Use the snipe type we select
                        case '1':
                            if (reDX.test(tx.data)) { // Search the mempool for the selected method
                                spamBot(toSnipe, amountOut, tries, tx, blockDelay) // Buy the token and pass on the transaction
                                sell(toSnipe, multiplier, percent, amountIn)
                            }
                            break

                        case '2':
                            if (tx.to === PCSRouter) {
                                if (re1.test(tx.data) || re2.test(tx.data)) {
                                    const decodedInput =
                                        pcsAbi.parseTransaction({
                                            data: tx.data,
                                            value: tx.value,
                                        })
                                    if (toSnipe === decodedInput.args[0]) {
                                        spamBot(toSnipe, amountOut, tries, tx, blockDelay)
                                        sell(toSnipe, multiplier, percent, amountIn)
                                    }
                                }
                            }
                            break

                        case '3':
                            if (
                                re3.test(tx.data) ||
                                re5.test(tx.data) ||
                                re6.test(tx.data) ||
                                re7.test(tx.data)
                            ) {
                                spamBot(toSnipe, amountOut, tries, tx, blockDelay)
                                sell(toSnipe, multiplier, percent, amountIn)
                            }
                            break

                        case '4':
                            if (trading1.test(tx.data)) {
                                if (tx.to === toSnipe) {
                                    spamBot(toSnipe, amountOut, tries, tx, blockDelay)
                                    sell(toSnipe, multiplier, percent, amountIn)
                                }
                            }
                            break
                        case '5':
                            if (re8.test(tx.data)) {
                                spamBot(toSnipe, amountOut, tries, tx, blockDelay)
                                sell(toSnipe, multiplier, percent, amountIn)
                            }
                            break
                    }
                }
            })
        })
    })


    // Check connection of WSS
    provider._websocket.on('close', () => {
        console.log('WebSocket Closed...Reconnecting...')
        clearInterval(keepAliveInterval)
        clearTimeout(pingTimeout)
        startConnection()
    })


    provider._websocket.on('error', () => {
        console.log('Error. Attempting to Reconnect...')
        clearInterval(keepAliveInterval)
        clearTimeout(pingTimeout)
        startConnection()
    })


    provider._websocket.on('pong', () => {
        clearInterval(pingTimeout)
    })
}
/////////////////////////////////////////////


/////////////////BUY/////////////////////////
// const BuyToken = async (txLP) => { // --> Receive the transaction from mempool
//     try {
//         if (blockDelay > 0) { // Do we want to delay before buying
//             console.log('Waiting ' + blockDelay + ' Blocks before buying...')
//             const receiptLP = await txLP.wait() 

//             const buyBlock = receiptLP.blockNumber + (blockDelay - 2) // What block to buy on
//             while ((await provider.getBlockNumber()) < buyBlock) {} // Pause until we are at the buy block

//             console.log('Waiting for Transaction receipt...')
//             const tx = await walletBot.sendTransaction({ // Send the buy transaction; we could pre sign this
//                 to: contractAdd,
//                 from: walletBot.address,
//                 gasLimit: 5000000,
//                 data: snipeData,
//                 gasPrice: 250000000000,
//             })
//             console.log(await provider.getBlockNumber())

//             const receipt = await tx.wait() // Get the transaction receipt
//             console.log('Token Purchase Complete')
//             console.log('Liquidity Addition Transaction Hash: ' + txLP.hash)
//             console.log('Your Transaction Hash: ' + receipt.transactionHash)

//             if (receipt.blockNumber === receiptLP.blockNumber) { // If we bought in the same block
//                 console.log(
//                     `There are ${
//                         receipt.transactionIndex - receiptLP.transactionIndex // Number of transactions between liquidity add and buy
//                     } transactions in the block between your snipe and the LP transaction.`
//                 )
//             } else {
//                 console.log(
//                     `There are ${
//                         receipt.blockNumber - receiptLP.blockNumber // Number of blocks between liquidity add and buy
//                     } block(s) between your snipe and the LP transaction.`
//                 )
//             }
//         } else { // Buy with no delay
//             console.log('Waiting for (NODELAY) Transaction receipt...')
//             const tx = await walletBot.sendTransaction({ // Send the buy transaction
//                 to: contractAdd,
//                 from: walletBot.address,
//                 gasLimit: 5000000,
//                 data: snipeData,
//                 gasPrice: txLP.gasPrice,
//             })

//             const receipt = await tx.wait()
//             const receiptLP = await txLP.wait()
//             console.log('Token Purchase Complete')
//             console.log('Liquidity Addition Transaction Hash: ' + txLP.hash)
//             console.log('Your Transaction Hash: ' + receipt.transactionHash)

//             if (receipt.blockNumber === receiptLP.blockNumber) {
//                 console.log(
//                     `There are ${
//                         receipt.transactionIndex - receiptLP.transactionIndex
//                     } transactions in the block between your snipe and the LP transaction.`
//                 )
//             } else {
//                 console.log(
//                     `There are ${
//                         receipt.blockNumber - receiptLP.blockNumber
//                     } block(s) between your snipe and the LP transaction.`
//                 )
//             }
//         }

//         if (argv.sell) { // If we chose to sell
//             await getTokenBalanceAndTrade(toSnipe) // Find how many tokens we got and sell them
//         } else {
//             process.exit() // Exit the bot
//         }
//     } catch (e) {}
// }
/////////////////////////////////////////////


////////////////BEGIN BOT////////////////////
startConnection()
/////////////////////////////////////////////