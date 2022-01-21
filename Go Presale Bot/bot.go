package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"modules/licensing"
	"modules/transactions"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"gopkg.in/yaml.v2"
)

type Config struct { // Initialize struct for config options
	Host           string  `yaml:"host"`
	PrivateKey     string  `yaml:"privatekey"`
	AmountIn       float64 `yaml:"amountin"`
	PresaleAddress string  `yaml:"presaleaddress"`
	PercentFill    int64   `yaml:"percentfill"`
	Action         string  `yaml:"action"`
}

type FormattedStorage struct { // Struct for storing contract data
	startTime  int64
	minBuy     *big.Int
	maxBuy     *big.Int
	hardCap    *big.Int
	startBlock *big.Int
}

// Initialize global variables
var cfg Config
var client *ethclient.Client
var presaleAddress common.Address

func main() {
	/*---------- Initialize -----------*/
	licensing.CheckPresaleBotLicense("http://wisdom-bots.com:3002", false, false)
	fmt.Println("Retrieving config...")
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
	fmt.Println("Done!")

	fmt.Println("Connecting to node...")
	client, err = ethclient.Dial(cfg.Host) // Connect to node
	if err != nil {
		fmt.Println("Unable to connect to node.")
		log.Fatal(err)
	}

	rpcClient, err := rpc.DialContext(context.Background(), cfg.Host) // Connect to node for RPC
	if err != nil {
		fmt.Println("Unable to connect to node.")
		log.Fatal(err)
	}

	gethClient := gethclient.New(rpcClient) // Initialize the GETH client
	fmt.Println("Connected!")

	privateKey, _, fromAddress := transactions.PrivateKeyExtrapolate(cfg.PrivateKey)
	/*------------------------*/

	/*----------Initialize Transaction------------*/
	presaleAddress = common.HexToAddress(cfg.PresaleAddress) // Init presale address
	presaleAddressPointer := &presaleAddress                 // Create a pointer type
	presaleAddressString := presaleAddressPointer.String()   // Create string of presale address

	value := transactions.ToWei(cfg.AmountIn) // Convert to wei

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress) // Get next available nonce
	if err != nil {
		fmt.Println("Unable to retrieve nonce.")
		log.Fatal(err)
	}

	chainId, err := client.ChainID(context.Background()) // Get current ChainID
	if err != nil {
		fmt.Println("Unable to retrieve ChainID.")
		log.Fatal(err)
	}

	signer := types.NewEIP155Signer(chainId) // Init signer

	one := big.NewInt(1) // Init big int of 1

	// Init some vars
	contributedAmount := big.NewInt(0)
	var data []byte

	// Start forming the tx
	tx := &types.LegacyTx{
		Nonce: nonce,
		To:    presaleAddressPointer,
		Value: value,
		Data:  data,
	}
	/*---------------------------------------------*/

	/*-----------Select Platform----------*/

	// Set data of the transaction
	if cfg.Action == "pinksale" {
		fmt.Println("Platform: PinkSale")
		tx.Data = common.Hex2Bytes("0xd7bb99ba")
	} else if cfg.Action == "unicrypt" {
		fmt.Println("Platform: Unicrypt")
		tx.Data, err = hexutil.Decode("0xf868e7660000000000000000000000000000000000000000000000000000000000000000")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("Unsupported action: %v", cfg.Action)
	}
	/*-------------------------------------*/

	/*----------Subscribe TxPool-----------*/
	hashes := make(chan common.Hash)

	_, err = gethClient.SubscribePendingTransactions(context.Background(), hashes) // Subscribe to the TxPool over WS
	if err != nil {
		fmt.Println("Unable to subscribe to the mempool. Ensure your node supports txpool.")
		log.Fatal(err)
	}
	/*----------------------------------------*/

	/*-----------Get Data-------------*/
	fmt.Println("Getting contract data...")
	bigStorage := getData() // Retrieve all data from the presale contract
	fmt.Println("Contract data retrieved!")
	/*--------------------------------*/

	pTmp := big.NewInt(1).Mul(bigStorage.hardCap, big.NewInt(cfg.PercentFill))
	buyAt := big.NewInt(1).Div(pTmp, big.NewInt(100)) // Set the value to submit transaction at

	/*-----------------Wait------------*/
	fmt.Println("Waiting for presale to start...")
	if cfg.Action == "pinksale" || cfg.Action == "dxSale" {
		go func() { // Create a goroutine so the sleep is non blocking
			for bigStorage.startTime-transactions.Unix() > 3 {
				fmt.Printf("%v Seconds remaining\n", bigStorage.startTime-transactions.Unix())
				time.Sleep(time.Second)
			}
		}()
		for bigStorage.startTime-transactions.Unix() > 4 { // Block until # seconds before presale start
		}
	} else if cfg.Action == "unicrypt" {

		go func() {
			for big.NewInt(0).Sub(bigStorage.startBlock, transactions.CurrentBlock(context.Background(), client)).Cmp(big.NewInt(2)) == 1 {
				fmt.Printf("%v Blocks remaining\n", big.NewInt(0).Sub(bigStorage.startBlock, transactions.CurrentBlock(context.Background(), client)))
				time.Sleep(time.Second * 3)
			}
		}()
		for big.NewInt(0).Sub(bigStorage.startBlock, transactions.CurrentBlock(context.Background(), client)).Cmp(big.NewInt(2)) == 1 { // Block until # blocks before presale start
		}
	}
	fmt.Println("Presale started... Searching the mempool")
	/*---------------------------------------*/
	// Code after this is time-sensitive
	/*-----------Analyze and Buy-------------*/
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

			// Instead of doing everything in the if statement, return if a condition is not met
			// Start with easiest, this avoids unnecessary computations
			// This should help with speed or CPU load
			if pendingTxTo.String() != presaleAddressString { // Check if the target address is the presale contract
				return
			}
			if pendingTx.Value().Cmp(bigStorage.minBuy) == -1 && pendingTx.Value().Cmp(bigStorage.maxBuy) == 1 { // Check if value is between the min and max buy
				return
			}

			contributedAmount.Add(pendingTx.Value(), contributedAmount) // Add the value of the pending transaction

			if contributedAmount.Cmp(buyAt) == -1 { // Check if we have reached the target contributed amount to buy at
				return
			}

			tx.Gas = pendingTx.Gas() * 2                               // Set gas limit
			tx.GasPrice = one.Add(pendingTx.GasPrice(), big.NewInt(1)) // Be right in front of the target transaction

			fmt.Printf("Transaction hash: https://bscscan.com/tx/%v", transactions.SignAndSendLegacyTx(context.Background(), tx, client, signer, privateKey))
			os.Exit(0) // Exit the program

		}()
	}
	/*-----------------------------------------------*/

}

var data FormattedStorage

func getData() FormattedStorage {
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	// navigate to a page, wait for an element, click
	if cfg.Action == "pinksale" {
		var start string
		var cap string
		var min string
		var max string
		err := chromedp.Run(ctx,
			emulation.SetUserAgentOverride("WebScraper 1.0"),
			chromedp.Navigate("https://www.pinksale.finance/#/launchpad/"+cfg.PresaleAddress+"?chain=BSC"),
			// wait for footer element is visible (ie, page is loaded)
			chromedp.ScrollIntoView(`footer`),
			chromedp.WaitVisible(`div.ant-card-body > div.ant-alert`, chromedp.ByQuery),
			chromedp.Text(`//tr[contains(., 'Presale Start Time')]/td[2]`, &start, chromedp.NodeVisible, chromedp.BySearch),
			chromedp.Text(`//tr[contains(., 'Hard Cap')]/td[2]`, &cap, chromedp.NodeVisible, chromedp.BySearch),
			chromedp.Text(`//tr[contains(., 'Minimum Buy')]/td[2]`, &min, chromedp.NodeVisible, chromedp.BySearch),
			chromedp.Text(`//tr[contains(., 'Maximum Buy')]/td[2]`, &max, chromedp.NodeVisible, chromedp.BySearch),
		)
		if err != nil {
			fmt.Println("Unable to retrieve contract information. Please make sure the contract address and action are entered correctly. Then run the bot again. Sometimes this error can happen multiple times in a row.")
			log.Fatal(err)
		}
		start = strings.TrimRight(start, "(UTC) ")
		t, err := time.Parse("2006.01.02 15:04", start)
		if err != nil {
			panic(err)
		}
		cap = strings.TrimRight(cap, " BN")
		min = strings.TrimRight(min, " BN")
		max = strings.TrimRight(max, " BN")

		cap64, err := strconv.ParseFloat(cap, 64)
		if err != nil {
			panic(err)
		}

		min64, err := strconv.ParseFloat(min, 64)
		if err != nil {
			panic(err)
		}

		max64, err := strconv.ParseFloat(max, 64)
		if err != nil {
			panic(err)
		}
		data = FormattedStorage{
			startTime: t.Unix(),
			minBuy:    transactions.ToWei(min64),
			maxBuy:    transactions.ToWei(max64),
			hardCap:   transactions.ToWei(cap64),
		}
	} else {
		var start string
		var cap string
		var max string
		err := chromedp.Run(ctx,
			emulation.SetUserAgentOverride("WebScraper 1.0"),
			chromedp.Navigate("https://app.unicrypt.network/amm/pancake-v2/ilo/"+cfg.PresaleAddress),
			// wait for footer element is visible (ie, page is loaded)
			chromedp.ScrollIntoView(`footer`),
			chromedp.WaitVisible(`div.v-tabs`, chromedp.ByQuery),
			chromedp.EvaluateAsDevTools(`$x("/html/body/div/div[1]/main/div/div[2]/div/div[2]/div[2]/div[3]/div[2]/div[3]/div[1]/div/div[2]/div/div[3]")[0].click()`, nil),

			chromedp.Text(`/html/body/div/div[1]/main/div/div[2]/div/div[2]/div[2]/div[3]/div[2]/div[3]/div[2]/div[2]/div/div[1]/div[3]/div[1]`, &start, chromedp.NodeVisible, chromedp.BySearch),
			chromedp.Text(`/html/body/div/div[1]/main/div/div[2]/div/div[2]/div[2]/div[3]/div[2]/div[3]/div[2]/div[2]/div/div[4]/div/div[1]`, &cap, chromedp.NodeVisible, chromedp.BySearch),
			chromedp.Text(`/html/body/div/div[1]/main/div/div[2]/div/div[2]/div[2]/div[3]/div[2]/div[3]/div[2]/div[2]/div/div[5]/div/div[1]`, &max, chromedp.NodeVisible, chromedp.BySearch),
		)
		if err != nil {
			fmt.Println("Unable to retrieve contract information. Please make sure the contract address and action are entered correctly. Then run the bot again. Sometimes this error can happen multiple times in a row.")
			log.Fatal(err)
		}

		parts := strings.Split(start, " ")
		start = parts[2]
		cap = strings.TrimRight(cap, " BN")
		max = strings.TrimRight(max, " BN")

		cap64, err := strconv.ParseFloat(cap, 64)
		if err != nil {
			panic(err)
		}

		max64, err := strconv.ParseFloat(max, 64)
		if err != nil {
			panic(err)
		}

		startInt, err := strconv.ParseInt(start, 10, 64)
		if err != nil {
			panic(err)
		}

		data = FormattedStorage{
			startBlock: big.NewInt(startInt),
			maxBuy:     transactions.ToWei(max64),
			hardCap:    transactions.ToWei(cap64),
		}
	}
	return data
}
