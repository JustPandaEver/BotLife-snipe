// inputs --> target, amountOut, tries

const { ethers } = require('ethers')
const axios = require('axios').default
var _ = require('lodash')
//var argv = require('minimist')(process.argv.slice(2))
const { rjust } = require('justify-text')
let keys = require('./keys.js')
const ethUtils = require('ethereumjs-util')
const { fork } = require('child_process')

let smartContract = '0xa34FE6296e6AFB7d8Ba7433465C9e4289A95786C'
let txValue = 1
let txGasPrice = 5000000000
let gasLimit = 5000000

const provider = new ethers.providers.JsonRpcProvider(
    'wss://bsc-ws-node.nariox.org:443'
)
let host = 'https://bsc-dataseed1.binance.org/'


function txMethodId(target, amountOut) {
    let method = '0x25b3283e' +
    rjust(target.replace('0x', ''), 64, '0') +
    rjust(parseInt(amountOut).toString(16).replace('x', ''), 64, '0') +
    rjust(1, 64, '0');
    return method;
}
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

let acctInfo = []
//keys are array in keys.keys

var j = 0
var arrayLength = keys.keys.length
for (var i = 0; i < arrayLength; i++) {
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
    acctInfo = await getNonces(acctInfo)
})();

module.exports = async function startBot(target, amountOut, tries, tx, blockDelay) {
    let methodData = txMethodId(target, amountOut);
    if (blockDelay > 0) { // Not sure this delay will work
        const receiptLP = await tx.wait()
        const buyBlock = receiptLP.blockNumber + (blockDelay - 2) // What block to buy on. Why - 2?
        // eslint-disable-next-line no-empty
        while ((await provider.getBlockNumber()) < buyBlock) {}
    } 
    // eslint-disable-next-line no-constant-condition
    for (let j = 0; j <= tries; j++) {
        for (i in acctInfo) {
            let tx = {
                value: txValue,
                nonce: acctInfo[i]['Nonce'],
                gasPrice: txGasPrice,
                gasLimit: gasLimit,
                chainId: 56,
                to: smartContract,
                data: methodData,
            }
            const signer = new ethers.Wallet(acctInfo[i]['Skey'], provider)
            let signedTx = await signer.signTransaction(tx)
            const sendTx = fork('sendtx.js')
            let providerInfo = {
                tx: signedTx,
                acct: acctInfo[i]['Skey'],
                contract: smartContract,
                provider: provider,
            }
            sendTx.send(providerInfo)
            sendTx.on('message', (result) => {
                if (result) {
                    // eslint-disable-next-line no-undef
                    process.exit()
                }
            })
        }
        acctInfo[i]['Nonce'] = acctInfo[i]['Nonce'] + 1
    }
}
