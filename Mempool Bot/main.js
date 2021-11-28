#!/usr/bin/env -S node --max-old-space-size=100000
/* eslint-disable no-empty */
//inputs --> config.json


////////////////IMPORTS//////////////////// 
'use strict'
// eslint-disable-next-line no-undef
const ethers = require('ethers')
const spamBot = require('./buy.js')
const approveToken = require('./approve.js')
const sell = require('./sell.js')
const config = require('./config.json')
//Object.assign(process.env, env)
//const env = require('./conf/env.json')
///////////////////////////////////////////


//////////////////CONSTANTS////////////////
const provider = new ethers.providers.WebSocketProvider(config["RpcProvider"]) //Json or WebSocket
const toSnipe = ethers.utils.getAddress(config["target"])
const snipeType = config["snipeType"]
const PCSRouter = config["PCSRouter"]
const pcsAbi = new ethers.utils.Interface(require('./abi.json'))
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
const trading1 = new RegExp('^' + config["tradingMethod"])
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

    console.log('Searching for transaction...\n')

    let pingTimeout = null
    let keepAliveInterval = null

    //Open the WSS
    provider._websocket.on('open', async () => {
        keepAliveInterval = setInterval(() => {
            provider._websocket.ping()
            pingTimeout = setTimeout(() => {
                provider._websocket.terminate()
            }, EXPECTED_PONG_BACK)
        }, KEEP_ALIVE_CHECK_INTERVAL)


        // Approve the token
        approveToken(toSnipe) 

        let multiplier = config["multiplier"]
        let percent = config["percentToSell"];
        const amountIn = await provider.getBalance(config["botContract"])
        let blockDelay = config["blockDelay"]
        provider.on('pending', async (txHash) => { // Look through the mempool
            provider.getTransaction(txHash).then(async (tx) => { // Get the transaction from the hash
                if (tx && tx.to) { // Check the transaction exists and it has a to address
                    switch (snipeType) { // Use the snipe type we select
                        case 'dxsale':
                            if (reDX.test(tx.data)) { // Search the mempool for the selected method
                                console.log("Transaction found: " + tx.hash + "\n");
                                spamBot(toSnipe, tx, blockDelay) // Buy the token and pass on the transaction
                                sell(toSnipe, multiplier, percent, amountIn)
                            }
                            break

                        case 'addliquidity':
                            if (tx.to === PCSRouter) {
                                if (re1.test(tx.data) || re2.test(tx.data)) {
                                    const decodedInput =
                                        pcsAbi.parseTransaction({
                                            data: tx.data,
                                            value: tx.value,
                                        })
                                    if (toSnipe === decodedInput.args[0]) {
                                        console.log("Transaction found: " + tx.hash + "\n");
                                        spamBot(toSnipe, tx, blockDelay)
                                        sell(toSnipe, multiplier, percent, amountIn)
                                    }
                                }
                            }
                            break

                        case 'unicrypt':
                            if (
                                re3.test(tx.data) ||
                                re5.test(tx.data) ||
                                re6.test(tx.data) ||
                                re7.test(tx.data)
                            ) {
                                console.log("Transaction found: " + tx.hash + "\n");
                                spamBot(toSnipe, tx, blockDelay)
                                sell(toSnipe, multiplier, percent, amountIn)
                            }
                            break

                        case 'enabletrading':
                            if (trading1.test(tx.data)) {
                                if (tx.to === toSnipe) {
                                    console.log("Transaction found: " + tx.hash + "\n");
                                    spamBot(toSnipe, tx, blockDelay)
                                    sell(toSnipe, multiplier, percent, amountIn)
                                }
                            }
                            break
                        case 'pinksale':
                            if (re8.test(tx.data)) {
                                console.log("Transaction found: " + tx.hash + "\n");
                                spamBot(toSnipe, tx, blockDelay)
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


////////////////BEGIN BOT////////////////////
startConnection()
/////////////////////////////////////////////