#!/usr/bin/env -S node

"use strict";
const ethers = require("ethers");
var myArgs = process.argv.slice(2);
var fs = require('fs');
var array = fs.readFileSync('./conf/presale.txt').toString().split("\n");
const {
  Worker
} = require('worker_threads');
console.log('Commands: dx, approve, dxclaim, uni, uniclaim, sell, withdraw, send, bnb, uncx, unidump, pink, pclaim, psnipe');

const masterAcc = "";
const masterKey = "";
const bscNode = "ws://localhost:8546";
const gasPrice = 250
const gasLimit = 500000
const claimLimit = 2000000
const claimPrice = 25

const provider = new ethers.providers.WebSocketProvider(bscNode);
const walletCount = Number(myArgs[0]);
const actionArg = String(myArgs[1]);
const txAmount = String(myArgs[2]);
const presaleAddy = ethers.utils.getAddress(String(myArgs[3]));
const tokenAddy = ethers.utils.getAddress(String(myArgs[4]));

if (actionArg == "uni") {
  uniWait();
} else if (actionArg == "pink") {
  pinkWait();
}  else if (actionArg == "psnipe") {
  pinkWhitelist();
}  else if (actionArg == "pclaim") {
  pinkClaim();
} else if (actionArg == "storage") {
  storage();
} else {
  start();
  console.log("Initializing workers...");
}

async function pinkWait() {
  //you may have to adjust the 3 if ur server is slow, to like 4 or 5
  while ((parseInt(await provider.getStorageAt(presaleAddy, 106)) - unix()) > 3) {
    process.stdout.write(parseInt(await provider.getStorageAt(presaleAddy, 106)) - unix() + " ");
    await new Promise(r => setTimeout(r, 500));
  }
  console.log("Onwards");
  start();
  console.log("Initializing workers...");
}

async function pinkWhitelist() {
  while (parseInt(await provider.getStorageAt(presaleAddy, 142)) < unix()) {
    console.log("Waiting for Whitelist end to be called")
    await new Promise(r => setTimeout(r, 1000));
  }
  //console.log(startTime);
  //you may have to adjust the 3 if ur server is slow, to like 4 or 5
  while ((parseInt(await provider.getStorageAt(presaleAddy, 142)) - unix()) > 3) {
    process.stdout.write(parseInt(await provider.getStorageAt(presaleAddy, 142)) - unix() + " ");
    await new Promise(r => setTimeout(r, 500));
  }
  console.log("Onwards");
  start();
  console.log("Initializing workers...");
}

async function pinkClaim() {
  console.log("Waiting for Finalize to be called")
  // can sometimes be 117 i guess
  while (parseInt(await provider.getStorageAt(presaleAddy, 118)) == 0) {
    await new Promise(r => setTimeout(r, 1000));
  }
  console.log("Onwards");
  start();
  console.log("Initializing workers...");
  console.log(parseInt(await provider.getStorageAt(presaleAddy, 118)) + " " + unix());
}

async function uniWait() {
  const startBlock = parseInt(await provider.getStorageAt(presaleAddy, 11))
  //console.log(startBlock)
  if (await provider.getBlockNumber() < (startBlock - 2)) {
    console.log("Waiting until Block " + (startBlock - 2) +
      " before submitting transactions.");
  }

  while (await provider.getBlockNumber() < (startBlock - 2)) {
    console.log("Current Block:" + await provider.getBlockNumber());
    await new Promise(r => setTimeout(r, 1000));
  }
  console.log("Onwards");
  start();
  console.log("Initializing workers...");
}

async function storage() {
   let index = 0
   for (index = 0; index < 100000; index++) {
     if (parseInt(await provider.getStorageAt(presaleAddy, index)) != 0){
     console.log(`[${index}]` + parseInt(await provider.getStorageAt(presaleAddy, index)));
   }
   }
}

function unix() {
  return parseInt(Math.floor(new Date().getTime() / 1000));
}

async function start() {
  let i = 1;
  while (walletCount >= i) {
    const temp = new ethers.Wallet(array[i - 1]);
    const myMessage = {
      masterKey: masterKey,
      workerCount: i,
      fillGasLimit: gasLimit,
      fillGasPrice: gasPrice,
      claimSellGasLimit: claimLimit,
      claimSellGasPrice: claimPrice,
      presaleAddy: presaleAddy,
      tokenAddy: tokenAddy,
      txAmount: txAmount,
      action: actionArg,
      privateKey: array[i - 1],
      walletAddy: temp.address,
      masterAcc: masterAcc,
      bscNode: bscNode
    }
    const worker = new Worker("./conf/worker.js");
    worker.postMessage(myMessage);
    i++
  }
}
