/* eslint-disable no-undef */
const { ethers } = require('ethers')
const axios = require('axios').default
var _ = require('lodash')
var argv = require('minimist')(process.argv.slice(2))
const { rjust } = require('justify-text')
let keys = require('./keys.js')
const ethUtils = require('ethereumjs-util')
const { fork } = require('child_process')

let smartContract = '0xa34FE6296e6AFB7d8Ba7433465C9e4289A95786C'
let txValue = 1
let txGasPrice = 5000000000
let gasLimit = 5000000

// args --> method, target, amountOut, tries

const provider = new ethers.providers.JsonRpcProvider(
    'wss://bsc-ws-node.nariox.org:443'
)
let host = 'https://bsc-dataseed1.binance.org/'

let txMethodId
if (argv.method == 'buy') {
    txMethodId =
        '0xf088d547000000000000000000000000' + argv.target.replace('0x', '')
}
if (argv.method == 'buy2') {
    txMethodId =
        '0x25b3283e' +
        rjust(argv.target.replace('0x', ''), 64, '0') +
        rjust(parseInt(argv.amountOut).toString(16).replace('x', ''), 64, '0') +
        rjust(argv.tries, 64, '0')
}

async function getNonces(acctInfo) {
    let noncePayloadList = []
    let nonceFetchId = 1
    let mapping = {}

    for (i in acctInfo) {
        staticAddress = acctInfo[i]['Address']
        noncePayloadList.push({
            jsonrpc: '2.0',
            method: 'eth_getTransactionCount',
            params: [staticAddress, 'latest'],
            id: nonceFetchId,
        })

        mapping[nonceFetchId] = staticAddress
        nonceFetchId += 1
    }

    noncesRaw = await RsvFetchCall(host, noncePayloadList)

    newAcctInfo = _.cloneDeep(acctInfo)

    for (result in noncesRaw) {
        address = mapping[result['id']]

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
    resultJSON = res.data

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
})()

;(async () => {
    // eslint-disable-next-line no-constant-condition
    while (true) {
        for (i in acctInfo) {
            tx = {
                value: txValue,
                nonce: acctInfo[i]['Nonce'],
                gasPrice: txGasPrice,
                gasLimit: gasLimit,
                chainId: 56,
                to: smartContract,
                data: txMethodId,
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
                    process.exit()
                }
            })
        }
        acctInfo[i]['Nonce'] = acctInfo[i]['Nonce'] + 1
    }
})()
