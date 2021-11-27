/* eslint-disable no-empty */
// inputs --> target
const ethers = require('ethers')
const config = require('./config.json')

const masterKey = config["masterKey"]
const provider = new ethers.providers.WebSocketProvider(config["RpcProvider"])
const walletBot = new ethers.Wallet(masterKey, provider)
const PCSRouter = config["PCSRouter"]



///////////////////INIT//////////////////////
module.exports = async function init(tokenToPurchase) {
    const toSnipe = ethers.utils.getAddress(tokenToPurchase)
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
    await approveToken(toSnipe, sellContract);
}
/////////////////////////////////////////////

///////////////////APPROVE/////////////////// --> Approve token for selling
async function approveToken(toSnipe, sellContract) {
    try {
        const allowance = await sellContract.allowance(
            walletBot.address,
            PCSRouter
        )
        if (allowance._hex === '0x00') {
            console.log(`Approving Token: ${toSnipe}`)
            const tx = await sellContract.approve(
                PCSRouter,
                ethers.constants.MaxUint256
            )
            await tx.wait()
            console.log('Token Approved')
        }
    } catch (e) {}
}
///////////////////////////////////////////////