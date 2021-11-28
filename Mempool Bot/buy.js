// inputs --> target, amountOut

//////////////////IMPORTS//////////////////////
const { ethers } = require('ethers')
const { rjust } = require('justify-text')
let keys = require('./keys.js')
const ethUtils = require('ethereumjs-util')
const { fork } = require('child_process')
const config = require('./config.json')
///////////////////////////////////////////////


////////////////CONFIG/////////////////////////
let smartContract = config["botContract"]
let buyGasLimit = config["buyGasLimit"]
let minAmtOut = config["minAmountOut"]
let target = config["target"]
let blockDelay = config["blockDelay"]
const provider = new ethers.providers.WebSocketProvider(config["RpcProvider"])
///////////////////////////////////////////////


//////////////BUY METHOD///////////////////////
function txMethodId(target) { // Set the data for the transaction
    let method = config['buyMethod'] +
    rjust(target.replace('0x', ''), 64, '0') +
    rjust(parseInt(minAmtOut).toString(16).replace('x', ''), 64, '0')
    return method;
}
///////////////////////////////////////////////


///////////////NONCE///////////////////////////
async function getNonces(j) {
    let nonce = await provider.getTransactionCount(acctInfo[j]["Address"])
    acctInfo[j]["Nonce"] = nonce
}
///////////////////////////////////////////////


////////////////////WALLETS////////////////////
let acctInfo = []
var j = 0
var arrayLength = keys.keys.length
for (var i = 0; i < arrayLength; i++) { // Go through each key, add it to acctInfo and get its address
    if (keys.keys[i].length > 1) {
        // eslint-disable-next-line no-undef
        let address = '0x' + ethUtils.privateToAddress(Buffer.from(keys.keys[i].trim().toLowerCase(), 'hex')).toString('hex')
        acctInfo[j] = {
            Skey: '0x' + keys.keys[i].trim(),
            Address: address,
            Nonce: 0,
        }
        getNonces(j)
    }
    j += 1
}
///////////////////////////////////////////////


//////////////////////BUY//////////////////////
let methodData = txMethodId(target); // Set the transaction data
module.exports = async function startBot(tx) {
    let goalGas = tx.gasPrice
    if (blockDelay > 0) { // Not sure this delay will work
        const receiptLP = await tx.wait()
        const buyBlock = receiptLP.blockNumber + (blockDelay - 2) // What block to buy on. Why - 2?
        // eslint-disable-next-line no-empty
        while ((await provider.getBlockNumber()) < buyBlock) {}
    }
    console.log("Beginning buys...\n");
    // eslint-disable-next-line no-constant-condition
    for (i in acctInfo) { // Buy with each wallet
        let tx = { // Create the transaction
            value: 0,
            nonce: acctInfo[i]['Nonce'],
            gasPrice: goalGas,
            gasLimit: ethers.BigNumber.from(buyGasLimit),
            chainId: parseInt(config["chainId"], 10),
            to: smartContract,
            data: methodData,
        }
        const signer = new ethers.Wallet(acctInfo[i]['Skey'], provider) // Create the signer from the wallet
        let signedTx = await signer.signTransaction(tx) // Sign the transaction
        const sendTx = fork('sendtx.js') // Use multiprocessing to send the transaction
        let providerInfo = { // Set the info the child process needs
            tx: signedTx,
            acct: acctInfo[i]['Skey']
        }
        sendTx.send(providerInfo) // Send the info
        sendTx.on('message', () => {})
    }
}
///////////////////////////////////////////////