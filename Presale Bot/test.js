const ethers = require("ethers");
const config = require('./config.json')

const provider = new ethers.providers.WebSocketProvider(config["provider"]);

// async function storage() {
//     let index = 0
//     for (index = 0; index < 100000; index++) {
//       if (parseInt(await provider.getStorageAt("0x11F0768b946E0447CD2b81AAfD67E9273031D3b2", index), 16) != 0){
//       console.log(`[${index}]` + parseInt(await provider.getStorageAt("0x11F0768b946E0447CD2b81AAfD67E9273031D3b2", index), 16));
//     }
//     }
//  }

//storage()
txHash = "0xcd0d984e0d3c4daa9ef3822e02a9bbe87c725446f96733a10ef0e459bd78af1c"
async function test () {
    hash = await provider.getTransaction(txHash)
    val = hash.data.slice(0, 10)
console.log(await val);
}
test()