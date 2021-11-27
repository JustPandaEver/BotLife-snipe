/* eslint-disable no-empty */
const ethers = require('ethers')

const walletBot = new ethers.Wallet(masterKey, provider)
const masterKey = ''
const provider = new ethers.providers.WebSocketProvider('ws://localhost:8546')
const PCSRouter = '0x10ED43C718714eb63d5aA57B78B54704E256024E'
const toSnipe = ethers.utils.getAddress(addressToSell)
let addressToSell;

///////////////////APPROVE/////////////////// --> Approve token for selling
module.exports = async function approveToken(tokenToPurchase) {
    try {
        const allowance = await sellContract.allowance(
            walletBot.address,
            PCSRouter
        )
        if (allowance._hex === '0x00') {
            console.log(`Approving Token: ${tokenToPurchase}`)
            const tx = await sellContract.approve(
                PCSRouter,
                ethers.constants.MaxUint256
            )
            await tx.wait()
            console.log('Token Approved')
        }
    } catch (e) {}
}
/////////////////////////////////////////////

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