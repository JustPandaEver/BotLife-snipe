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
txHash = "0xa8ab3931c3476c97f4bc3443d388a44bb3d6c163a49e87c878956b875f3e9a17"
async function test () {
    hash = await provider.getTransaction(txHash)
    val = hash.value
console.log(await val);
}
test()