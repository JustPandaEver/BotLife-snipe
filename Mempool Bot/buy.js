// inputs --> target, amountOut, tries

//////////////////IMPORTS//////////////////////
const { ethers } = require('ethers')
const axios = require('axios').default
var _ = require('lodash')
const { rjust } = require('justify-text')
let keys = require('./keys.js')
const ethUtils = require('ethereumjs-util')
const { fork } = require('child_process')
const config = require('./config.json')
///////////////////////////////////////////////


////////////////CONFIG/////////////////////////
let smartContract = config["botContract"]
let txValue = 1
let buyTxGasPrice = config["buyTxGasPrice"] // do these need to be converted to numbers
let buyGasLimit = config["buyGasLimit"]
const provider = new ethers.providers.JsonRpcProvider( //Json or WebSocket
    config["RpcProvider"]
)
let host = config["host"]
///////////////////////////////////////////////


//////////////BUY METHOD///////////////////////
function txMethodId(target, amountOut) { // Set the data for the transaction
    let method = config['buyMethod'] +
    rjust(target.replace('0x', ''), 64, '0') +
    rjust(parseInt(amountOut).toString(16).replace('x', ''), 64, '0') +
    rjust(1, 64, '0');
    return method;
}
///////////////////////////////////////////////


///////////////NONCE///////////////////////////
// What is all this nonce stuff? Is it necessary? Can we just start at 1 and increase
async function getNonces(acctInfo) {
    let noncePayloadList = []
    let nonceFetchId = 1
    let mapping = {}

    for (i in acctInfo) {
        let staticAddress = acctInfo[i]['Address']
        noncePayloadList.push({
            jsonrpc: '2.0',
            method: 'eth_getTransactionCount',
            params: [staticAddress, 'latest'],
            id: nonceFetchId,
        })

        mapping[nonceFetchId] = staticAddress
        nonceFetchId += 1
    }

    let noncesRaw = await RsvFetchCall(host, noncePayloadList)

    let newAcctInfo = _.cloneDeep(acctInfo)

    for (let result in noncesRaw) {
        let address = mapping[result['id']]

        for (i in acctInfo) {
            if (acctInfo[i]['Address'].toLowerCase() == address) {
                newAcctInfo[i]['Nonce'] = parseInt(result['result'], 16)
            }
        }
    }
    console.log(newAcctInfo)
    return newAcctInfo
}
async function RsvFetchCall(url, payload) {
    let axiosHeaders = {
        headers: {
            'Content-Type': 'application/json',
        },
    }

    let res = await axios.post(url, payload, axiosHeaders)
    let resultJSON = res.data

    return resultJSON
}
///////////////////////////////////////////////


////////////////////WALLETS////////////////////
let acctInfo = []
//keys are array in keys.keys

var j = 0
var arrayLength = keys.keys.length
for (var i = 0; i < arrayLength; i++) { // Go through each key, add it to acctInfo and get its address
    if (keys.keys[i].length > 1) {
        acctInfo[j] = {
            Skey: '0x' + keys.keys[i].trim(),
            Address:
                '0x' +
                ethUtils
                    .privateToAddress(
                        // eslint-disable-next-line no-undef
                        Buffer.from(keys.keys[i].trim().toLowerCase(), 'hex')
                    )
                    .toString('hex'),
            Nonce: 0,
        }
    }
    j += 1
}

(async () => {
    acctInfo = await getNonces(acctInfo) // Get nonces for each account?
})();
///////////////////////////////////////////////


//////////////////////BUY//////////////////////
module.exports = async function startBot(target, amountOut, tx, blockDelay) {
    let methodData = txMethodId(target, amountOut); // Set the transaction data
    if (blockDelay > 0) { // Not sure this delay will work
        const receiptLP = await tx.wait()
        const buyBlock = receiptLP.blockNumber + (blockDelay - 2) // What block to buy on. Why - 2?
        // eslint-disable-next-line no-empty
        while ((await provider.getBlockNumber()) < buyBlock) {}
    } 
    // eslint-disable-next-line no-constant-condition
        for (i in acctInfo) { // Buy with each wallet
            let tx = { // Create the transaction
                value: txValue,
                nonce: acctInfo[i]['Nonce'],
                gasPrice: buyTxGasPrice,
                gasLimit: buyGasLimit,
                chainId: config["chainId"],
                to: smartContract,
                data: methodData,
            }
            const signer = new ethers.Wallet(acctInfo[i]['Skey'], provider) // Create the signer from the wallet
            let signedTx = await signer.signTransaction(tx) // Sign the transaction
            const sendTx = fork('sendtx.js') // Use multiprocessing to send the transaction
            let providerInfo = { // Set the info the child process needs
                tx: signedTx,
                acct: acctInfo[i]['Skey'],
                contract: smartContract,
                provider: provider,
            }
            sendTx.send(providerInfo) // Send the info
            sendTx.on('message', () => {})
        }
        acctInfo[i]['Nonce'] = acctInfo[i]['Nonce'] + 1 // Increment the nonce
}
///////////////////////////////////////////////