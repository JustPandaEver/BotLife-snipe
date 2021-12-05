"use strict";
const ethers = require("ethers");
const {
  Worker,
  parentPort
} = require('worker_threads');

const getInfo = parentPort.on('message', (message) => {
  const workerCount = message.workerCount
  const actionArg = message.action
  const walletAddy = message.walletAddy
  const privateKey = message.privateKey
  const buyGasLimit = message.fillGasLimit
  const buyGasPrice = message.fillGasPrice
  const presaleAddy = message.presaleAddy
  const tokenAddy = message.tokenAddy
  const txAmount = message.txAmount
  const claimSellGasLimit = message.claimSellGasLimit
  const claimSellGasPrice = message.claimSellGasPrice
  const masterKey = message.masterKey
  const masterAccount = message.masterAcc
  const bscProvider = message.bscNode

  const getRandomGasPrice = (min, max) => {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min) + min);
  }
  const pcsrouter = "0x10ED43C718714eb63d5aA57B78B54704E256024E";
  const bnb = "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c";
  const txAmountFormat = ethers.utils.parseUnits(String(txAmount), "ether");
  const buyGasLimitFormat = ethers.BigNumber.from(buyGasLimit)
  const gasPriceFormat = ethers.utils.parseUnits(`${buyGasPrice}`, 'gwei')
  const realBuyGasPrice = gasPriceFormat
  const claimSellGasLimitFormat = ethers.BigNumber.from(claimSellGasLimit)
  const claimSellGasPriceFormat = ethers.utils.parseUnits(`${claimSellGasPrice}`, 'gwei')
  const realClaimSellGasPrice = claimSellGasPriceFormat
  const EXPECTED_PONG_BACK = 30000;
  const KEEP_ALIVE_CHECK_INTERVAL = 15000;
  const provider = new ethers.providers.WebSocketProvider(
    bscProvider
  );
  const pcsAbi = new ethers.utils.Interface(require("./abi.json"));
  const wallet = new ethers.Wallet(privateKey);
  const account = wallet.connect(provider);
  const masterWallet = new ethers.Wallet(masterKey);
  const masterConnect = masterWallet.connect(provider);
  const router = new ethers.Contract(pcsrouter, pcsAbi, account);
  const startConnection = () => {
    let pingTimeout = null;
    let keepAliveInterval = null;
    provider._websocket.on("open", () => {
      keepAliveInterval = setInterval(() => {
        provider._websocket.ping();
        pingTimeout = setTimeout(() => {
          provider._websocket.terminate();
        }, EXPECTED_PONG_BACK);
      }, KEEP_ALIVE_CHECK_INTERVAL);

      if (actionArg === "dx") {

        dxSaleBuy();

      } else if (actionArg === "approve") {

        getApprove();

      } else if (actionArg === "dxclaim") {

        dxSaleClaimTokens();

      } else if (actionArg === "uni") {

        uniOld();

      } else if (actionArg === "uniclaim") {

        uniClaim();

      } else if (actionArg === "sell") {

        sellToken();

      } else if (actionArg === "withdraw") {

        withDrawAllTokens();

      } else if (actionArg === "send") {

        sendEnMasse();

      } else if (actionArg === "bnb") {

        bnbb();
      } else if (actionArg === "uncx") {

        uncx();
      } else if (actionArg === "unidump") {

        uniDump();
      } else if (actionArg === "pink") {

        pink();
      } else if (actionArg === "psnipe") {

        pink();
      } else if (actionArg === "pclaim") {

        pinkClaim();
      } else {

        console.log("Invalid action passed...")
      }

    });

    provider._websocket.on("close", () => {
      console.log("WebSocket Closed...Reconnecting...");
      clearInterval(keepAliveInterval);
      clearTimeout(pingTimeout);
      startConnection();
    });

    provider._websocket.on("error", () => {
      console.log("Error. Attemptiing to Reconnect...");
      clearInterval(keepAliveInterval);
      clearTimeout(pingTimeout);
      startConnection();
    });

    provider._websocket.on("pong", () => {
      clearInterval(pingTimeout);
    });
  };
  //Buy DXSALE
  const dxSaleBuy = async () => {
    try {
      const tx =
        await account.sendTransaction({
          to: presaleAddy,
          from: walletAddy,
          gasLimit: buyGasLimitFormat,
          gasPrice: realBuyGasPrice,
          value: txAmountFormat,
        });
      const receipt = await tx.wait();
      console.log("Worker " + workerCount + " bought " + txAmount + " BNB of " + presaleAddy + " successfully.  Receipt: " + receipt.transactionHash)
    } catch (error) {
      if (error.reason, !error.transactionHash) {
        console.log("Worker " + workerCount + " Error reason: " + error.reason)
      } else if (error.transactionHash) {
        console.log("Worker " + workerCount + " Error transaction hash: " + error.transactionHash)
      }
    }
  };
  //unicrypt old
  const uniOld = async () => {
    try {
      const tx =
        await account.sendTransaction({
          to: presaleAddy,
          from: walletAddy,
          gasLimit: buyGasLimitFormat,
          gasPrice: realBuyGasPrice,
          value: txAmountFormat,
          data: "0xf868e7660000000000000000000000000000000000000000000000000000000000000000"
        });
      const receipt = await tx.wait();
      console.log("Worker " + workerCount + " bought Unicrypt " + presaleAddy + " successfully.  Receipt: " + receipt.transactionHash)
    } catch (error) {
      if (error.reason, !error.transactionHash) {
        console.log("Worker " + workerCount + " Error reason: " + error.reason)
      } else if (error.transactionHash) {
        console.log("Worker " + workerCount + " Error transaction hash: " + error.transactionHash)
      }
    }
  };
  const pink = async () => {
    try {
      const tx =
        await account.sendTransaction({
          to: presaleAddy,
          from: walletAddy,
          gasLimit: buyGasLimitFormat,
          gasPrice: realBuyGasPrice,
          value: txAmountFormat,
          data: "0xd7bb99ba"
        });
      const receipt = await tx.wait();
      console.log("Worker " + workerCount + " bought " + txAmount + " BNB of " + presaleAddy + " successfully.  Receipt: " + receipt.transactionHash)
    } catch (error) {
      if (error.reason, !error.transactionHash) {
        console.log("Worker " + workerCount + " Error reason: " + error.reason)
      } else if (error.transactionHash) {
        console.log("Worker " + workerCount + " Error transaction hash: " + error.transactionHash)
      }
    }
  };
  const pinkClaim = async () => {
    try {
      const tx =
        await account.sendTransaction({
          to: presaleAddy,
          from: walletAddy,
          gasLimit: claimSellGasLimitFormat,
          gasPrice: realClaimSellGasPrice,
          data: "0x4e71d92d"
        });
      const receipt = await tx.wait();
      console.log("Worker " + workerCount + " claimed tokens from " + presaleAddy + " successfully.  Receipt: " + receipt.transactionHash)
    } catch (error) {
      if (error.reason, !error.transactionHash) {
        console.log("Worker " + workerCount + " Error reason: " + error.reason)
      } else if (error.transactionHash) {
        console.log("Worker " + workerCount + " Error transaction hash: " + error.transactionHash)
      }
    }
  };
  //unicrypt claim
  const uniClaim = async () => {
    try {
      const tx =
        await account.sendTransaction({
          to: presaleAddy,
          from: walletAddy,
          gasLimit: claimSellGasLimitFormat,
          gasPrice: realClaimSellGasPrice,
          data: "0xfe8121de"
        });
      const receipt = await tx.wait();
      console.log("Worker " + workerCount + " claimed tokens from " + presaleAddy + " successfully.  Receipt: " + receipt.transactionHash)
    } catch (error) {
      if (error.reason, !error.transactionHash) {
        console.log("Worker " + workerCount + " Error reason: " + error.reason)
      } else if (error.transactionHash) {
        console.log("Worker " + workerCount + " Error transaction hash: " + error.transactionHash)
      }
    }
  };
  //unidump
  const uniDump = async () => {
    try {
      const tx =
        await account.sendTransaction({
          to: presaleAddy,
          from: walletAddy,
          gasLimit: claimSellGasLimitFormat,
          gasPrice: realClaimSellGasPrice,
          data: "0xfe8121de"
        });
      const receipt = await tx.wait();
      console.log("Worker " + workerCount + " claimed tokens from " + presaleAddy + " successfully.  Receipt: " + receipt.transactionHash)
      await sellToken();
    } catch (error) {
      if (error.reason, !error.transactionHash) {
        console.log("Worker " + workerCount + " Error reason: " + error.reason)
      } else if (error.transactionHash) {
        console.log("Worker " + workerCount + " Error transaction hash: " + error.transactionHash)
      }
    }
  };
  //Approve
  const getApprove = async () => {
    try {
      const sellAmount = ethers.constants.MaxUint256;
      const sellContract = new ethers.Contract(
        tokenAddy,
        [
          "function approve(address _spender, uint256 _value) public returns (bool success)",
        ],
        account
      );
      console.log(((new Date()).toISOString()).replace('Z', '') + ": Worker " + workerCount + " Before Approval Sent!");
      const tx = await sellContract.approve(pcsrouter, sellAmount);
      console.log(((new Date()).toISOString()).replace('Z', '') + ": Worker " + workerCount + " Approval Sent!");
      const receipt = await tx.wait();
      console.log(((new Date()).toISOString()).replace('Z', '') + ": Worker " + workerCount + " Approved! Receipt: " + receipt.transactionHash);
    } catch (error) {
      if (error.reason, !error.transactionHash) {
        console.log("Worker " + workerCount + " Error reason: " + error.reason)
      } else if (error.transactionHash) {
        console.log("Worker " + workerCount + " Error transaction hash: " + error.transactionHash)
      }
    }
  };
  //Claim DXSALE
  const dxSaleClaimTokens = async () => {
    try {
      const tx =
        await account.sendTransaction({
          to: presaleAddy,
          from: walletAddy,
          gasLimit: claimSellGasLimitFormat,
          gasPrice: realClaimSellGasPrice,
          data: "0x48c54b9d"
        });
      const receipt = await tx.wait();
      console.log("Worker " + workerCount + " claimed tokens from " + presaleAddy + " successfully.  Receipt: " + receipt.transactionHash)
    } catch (error) {
      if (error.reason, !error.transactionHash) {
        console.log("Worker " + workerCount + " Error reason: " + error.reason)
      } else if (error.transactionHash) {
        console.log("Worker " + workerCount + " Error transaction hash: " + error.transactionHash)
      }
    }
  };
  //Sell Tokens
  const sellToken = async () => {
    try {
      const tokenTransfer = new ethers.Contract(
        tokenAddy,
        [
          "function transfer(address recipient, uint256 amount) external returns (bool)",
          'function balanceOf(address owner) view returns (uint)',
        ],
        account
      );
      let a = ethers.BigNumber.from(await tokenTransfer.balanceOf(walletAddy));

      const tx =
        await router.swapExactTokensForETHSupportingFeeOnTransferTokens(
          a.mul(txAmount).div(100),
          ethers.BigNumber.from(0),
          [tokenAddy, bnb],
          walletAddy,
          Date.now() + 1000 * 60 * 5, //5 minutes
          {
            gasLimit: claimSellGasLimitFormat,
            gasPrice: realClaimSellGasPrice,
          }
        );
      const receipt = await tx.wait();
      console.log("Worker " + workerCount + " sold " + txAmount + " tokens of " + tokenAddy + " successfully.  Receipt: " + receipt.transactionHash)
    } catch (error) {
      if (error.reason, !error.transactionHash) {
        console.log("Worker " + workerCount + " Error reason: " + error.reason)
      } else if (error.transactionHash) {
        console.log("Worker " + workerCount + " Error transaction hash: " + error.transactionHash)
      }
    }
  }
  //withdraw
  const withDrawAllTokens = async () => {
    try {
      const tokenTransfer = new ethers.Contract(
        tokenAddy,
        [
          "function transfer(address recipient, uint256 amount) external returns (bool)",
          'function balanceOf(address owner) view returns (uint)',
        ],
        account
      );
      const tx = await tokenTransfer.transfer(masterAccount, tokenTransfer.balanceOf(account.address), {
        gasLimit: ethers.BigNumber.from(250000),
        gasPrice: ethers.utils.parseUnits('15', 'gwei'),
      });
      const receipt = await tx.wait();
      console.log("Worker " + workerCount + " has transferred tokens to master! Receipt: " + receipt.transactionHash);
    } catch (error) {
      if (error.reason, !error.transactionHash) {
        console.log("Worker " + workerCount + " Error reason: " + error.reason)
      } else if (error.transactionHash) {
        console.log("Worker " + workerCount + " Error transaction hash: " + error.transactionHash)
      }
    }
  }
  //withdraw BNB
  const bnbb = async () => {
    try {
      let balance = await account.getBalance();
      let gasPrice = await provider.getGasPrice();
      let gasLimit = 21000;
      let value = balance.sub(gasPrice.mul(gasLimit))

      let tx = await account.sendTransaction({
        gasLimit: gasLimit,
        gasPrice: gasPrice,
        to: masterAccount,
        value: value
      });
      const receipt = await tx.wait();
      console.log("Worker " + workerCount + " has transferred all BNB to master! Receipt: " + receipt.transactionHash);
    } catch (error) {
      if (error.reason, !error.transactionHash) {
        console.log("Worker " + workerCount + " Error reason: " + error.reason)
      } else if (error.transactionHash) {
        console.log("Worker " + workerCount + " Error transaction hash: " + error.transactionHash)
      }
    }
  }
  //deposit bnb to wallets
  const sendEnMasse = async () => {
    try {
      const tx =
        await masterConnect.sendTransaction({
          to: walletAddy,
          from: masterAccount,
          gasLimit: ethers.BigNumber.from(110000),
          gasPrice: ethers.utils.parseUnits(`6`, 'gwei'),
          value: txAmountFormat
        });
      const receipt = await tx.wait();
      console.log("Master sent " + txAmount + "BNB successfully to worker " + workerCount + " .  Receipt: " + receipt.transactionHash)
    } catch (error) {
      if (error.reason, !error.transactionHash) {
        console.log("Worker " + workerCount + " Error reason: " + error.reason)
        sendEnMasse();
      } else if (error.transactionHash) {
        console.log("Worker " + workerCount + " Error transaction hash: " + error.transactionHash)
        sendEnMasse();
      }
    }
  }
  const uncx = async () => {
    try {
      const tokenTransfer = new ethers.Contract(
        '0x09a6c44c3947b69e2b45f4d51b67e6a39acfb506',
        [
          "function transfer(address recipient, uint256 amount) external returns (bool)",
          'function balanceOf(address owner) view returns (uint)',
        ],
        masterConnect
      );
      const tx = await tokenTransfer.transfer(walletAddy, ethers.utils.parseUnits(`3000000000000000000`, 'wei'), {
        gasLimit: ethers.BigNumber.from(100000),
        gasPrice: ethers.utils.parseUnits('5', 'gwei'),
      });
      const receipt = await tx.wait();
      console.log("Sent UNCX to " + workerCount + "! Receipt: " + receipt.transactionHash);
    } catch (error) {
      if (error.reason, !error.transactionHash) {
        console.log("Worker " + workerCount + " Error reason: " + error.reason)
        uncx();
      } else if (error.transactionHash) {
        console.log("Worker " + workerCount + " Error transaction hash: " + error.transactionHash)
        uncx();
      }
    }
  }
  startConnection();

});
