package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"modules/transactions"
	"os"
	"strings"
	"time"

	erc20 "BotLife/AntirugBot/erc20"
	pcsRouter "BotLife/AntirugBot/pcsrouter"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"gopkg.in/yaml.v2"
)

type Config struct { // Initialize struct for config options
	Host       string `yaml:"host"`
	PrivateKey string `yaml:"privateKey"`
	Target     string `yaml:"target"`
	PcsRouter  string `yaml:"pcsRouter"`
	BnbAddress string `yaml:"bnbAddress"`
}

type Client struct {
	c *rpc.Client
}

var cfg Config
var client *ethclient.Client
var privKey *ecdsa.PrivateKey
var fromAddress common.Address
var pcsRouterAddress common.Address
var targetAddress common.Address
var bnbAddress common.Address
var pcsInstance *pcsRouter.PcsRouter
var ercInstance *erc20.ERC20Abi

func main() {

	f, err := os.Open("./config.yml") // Open the config file
	if err != nil {
		fmt.Println("Ensure config is present.")
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(f) // Create variables from config options
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Println("Config file is likely formatted incorrectly.")
		log.Fatal(err)
	}
	f.Close()

	client, err = ethclient.Dial(cfg.Host)
	if err != nil {
		panic(err)
	}

	rpcClient, err := rpc.Dial(cfg.Host)
	if err != nil {
		panic(err)
	}

	pcsRouterAddress = common.HexToAddress(cfg.PcsRouter)

	gethClient := gethclient.New(rpcClient)

	privKey, _, fromAddress = transactions.PrivateKeyExtrapolate(cfg.PrivateKey)

	pcsInstance, err = pcsRouter.NewPcsRouter(pcsRouterAddress, client)
	if err != nil {
		fmt.Println("Unable to connect to PCS Router")
	}

	bnbAddress = common.HexToAddress(cfg.BnbAddress)
	targetAddress = common.HexToAddress(cfg.Target)
	ercInstance, err = erc20.NewERC20Abi(targetAddress, client)
	if err != nil {
		panic(err)
	}

	hashes := make(chan common.Hash)

	_, err = gethClient.SubscribePendingTransactions(context.Background(), hashes) // Subscribe to the TxPool over WS
	if err != nil {
		fmt.Println("Unable to subscribe to the mempool. Ensure your node supports txpool.")
		log.Fatal(err)
	}

	fmt.Println("Go!")
	for { // Create a new thread for every pending transaction
		go func() {

			hash := <-hashes // Receive the hash

			pendingTx, _, err := client.TransactionByHash(context.Background(), hash) // Get transaction from hash
			if err != nil {
				return
			}

			pendingTxTo := pendingTx.To()

			if pendingTxTo == nil {
				return
			}

			if *pendingTx.To() == pcsRouterAddress {
				if decodeLiquidityInput(pendingTx) {
					sellCoin(pendingTx)
				}
			} else {
				return
			}

			if checkRug(pendingTx) {
				sellCoin(pendingTx)
			} else {
				return
			}

		}()
	}

}

func checkRug(pendingTx *types.Transaction) bool {

	if checkHoneypot(pendingTx) || checkLiquidity(pendingTx) || devSell(pendingTx) || mint(pendingTx) {
		return true
	}
	return false
}

func sellCoin(pendingTx *types.Transaction) {

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privKey, chainId)
	if err != nil {
		log.Fatal(err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println("Unable to retrieve pending nonce.")
		log.Fatal(err)
	}

	txNonce := big.NewInt(int64(nonce))

	auth.GasLimit = 500000
	auth.GasPrice = pendingTx.GasPrice().Add(pendingTx.GasPrice(), big.NewInt(10000))

	transactOps := &bind.TransactOpts{
		From:     fromAddress,
		Nonce:    txNonce,
		Signer:   auth.Signer,
		Value:    big.NewInt(0),
		GasPrice: auth.GasPrice,
		GasLimit: auth.GasLimit,
		Context:  context.Background(),
	}

	balance := getBalance(false)

	deadline := big.NewInt(time.Now().Unix() + 1200)
	path := []common.Address{targetAddress, bnbAddress}
	tx, err := pcsInstance.SwapExactTokensForETHSupportingFeeOnTransferTokens(transactOps, balance, big.NewInt(0), path, fromAddress, deadline)
	if err != nil {
		fmt.Printf("Transaction failed.\n%q", err)
	}
	fmt.Printf("Transaction receipt: %v", tx.Hash())
	os.Exit(0)
}

func getBalance(pending bool) *big.Int {
	opts := &bind.CallOpts{
		Pending: pending,
		Context: context.Background(),
	}

	balance, err := ercInstance.BalanceOf(opts, fromAddress)
	if err != nil {
		fmt.Println(err)
	}

	return balance
}

func checkHoneypot(pendingTx *types.Transaction) bool {
	return false
}

func checkLiquidity(pendingTx *types.Transaction) bool {

	if *pendingTx.To() == pcsRouterAddress {
		return decodeLiquidityInput(pendingTx)
	}
	return false
}

func devSell(pendingTx *types.Transaction) bool {
	return false
}

func mint(pendingTx *types.Transaction) bool {
	return false
}

func checkAmountOut() bool {
	opts := &bind.CallOpts{
		Pending: true,
		Context: context.Background(),
	}
	path := []common.Address{targetAddress, bnbAddress}
	amountsOut, err := pcsInstance.GetAmountsOut(opts, getBalance(true), path)
	if err != nil {
		panic(err)
	}
	// return getBalance(true).Cmp(big.NewInt(100000000000000)) == -1
	return amountsOut[1].Cmp(big.NewInt(100000000000000)) == -1
}

func decodeLiquidityInput(pendingTx *types.Transaction) bool {

	data := common.Bytes2Hex(pendingTx.Data())
	var method string
	if len(data) >= 10 {
		method = data[:8]
	} else {
		return false
	}

	removeLiquidity := "baa2abde"
	removeLiquidityEth := "02751cec"
	removeLiquidityETHSupportingFeeOnTransferTokens := "af2979eb"
	removeLiquidityETHWithPermit := "ded9382a"
	removeLiquidityETHWithPermitSupportingFeeOnTransferTokens := "5b0d5984"
	removeLiquidityWithPermit := "2195995c"

	if method == removeLiquidity ||
		method == removeLiquidityEth ||
		method == removeLiquidityETHSupportingFeeOnTransferTokens ||
		method == removeLiquidityETHWithPermit ||
		method == removeLiquidityETHWithPermitSupportingFeeOnTransferTokens ||
		method == removeLiquidityWithPermit {

		if strings.Contains(data, strings.ToLower(cfg.Target[2:])) {
			return true
		}

	}
	return false
}
