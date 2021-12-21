package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"modules/transactions"
	"os"
	"strings"
	"time"

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

type StorageHex struct { // Initialize struct for storing contract data locations in hex
	startTime    string
	minBuy       string
	maxBuy       string
	hardCap      string
	tokenAddress string
	startBlock   string
}

type FormattedStorage struct { // Struct for storing contract data
	startTime    *big.Int
	minBuy       *big.Int
	maxBuy       *big.Int
	hardCap      *big.Int
	tokenAddress *big.Int
	startBlock   *big.Int
}

// Initialize global variables
var cfg Config
var client *ethclient.Client
var presaleAddress common.Address

func main() {
	/*---------- Initialize -----------*/
	fmt.Println("Retrieving config...")
	f, err := os.Open("./config.yml") // Open the config file
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(f) // Create variables from config options
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	fmt.Println("Done!")

	fmt.Println("Connecting to node...")
	client, err = ethclient.Dial(cfg.Host) // Connect to node
	if err != nil {
		log.Fatal(err)
	}

	rpcClient, err := rpc.DialContext(context.Background(), cfg.Host) // Connect to node for RPC
	if err != nil {
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
		log.Fatal(err)
	}

	chainId, err := client.ChainID(context.Background()) // Get current ChainID
	if err != nil {
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

	/*------------------Storage Addresses----------*/
	pinkSale := StorageHex{ // Set the hex of storage address for the platforms
		startTime:    "0x6A",
		minBuy:       "0x82",
		maxBuy:       "0x83",
		hardCap:      "0x85",
		tokenAddress: "0x69",
	}

	dxSale := StorageHex{
		startTime:    "0x2D",
		minBuy:       "0x1A",
		maxBuy:       "0x1B",
		hardCap:      "0x2C",
		tokenAddress: "0x0",
	}

	unicrypt := StorageHex{
		maxBuy:       "0x5",
		hardCap:      "0x7",
		tokenAddress: "0x2",
		startBlock:   "0xB",
	}
	/*------------------------------------*/

	/*-----------Select Platform----------*/
	var storageHexVals StorageHex

	// Set the hex locations to use and set data of the transaction
	if cfg.Action == "pinksale" {
		fmt.Println("Platform: PinkSale")
		storageHexVals = pinkSale
		tx.Data = common.Hex2Bytes("0x")
	} else if cfg.Action == "dxSale" {
		fmt.Println("Platform: dxSale")
		storageHexVals = dxSale
		tx.Data = common.Hex2Bytes("0x")
	} else if cfg.Action == "unicrypt" {
		fmt.Println("Platform: Unicrypt")
		storageHexVals = unicrypt
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
		log.Fatal(err)
	}
	/*----------------------------------------*/

	/*-----------Get Data-------------*/
	fmt.Println("Getting contract data...")
	bigStorage := formatStorage(storageHexVals) // Retrieve all data from the presale contract
	fmt.Println("Contract data retrieved!")
	/*--------------------------------*/

	pTmp := big.NewInt(1).Mul(bigStorage.hardCap, big.NewInt(cfg.PercentFill))
	buyAt := big.NewInt(1).Div(pTmp, big.NewInt(100)) // Set the value to submit transaction at

	/*-----------------Wait------------*/
	fmt.Println("Waiting for presale to start...")
	if cfg.Action == "pinksale" || cfg.Action == "dxSale" {
		go func() { // Create a goroutine so the sleep is non blocking
			for bigStorage.startTime.Int64()-transactions.Unix() > 4 {
				fmt.Printf("%v Seconds remaining\n", bigStorage.startTime.Int64()-transactions.Unix())
				time.Sleep(time.Second)
			}
		}()
		for bigStorage.startTime.Int64()-transactions.Unix() > 4 { // Block until # seconds before presale start
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

			tx.Gas = pendingTx.Gas()                                   // Set gas limit
			tx.GasPrice = one.Add(pendingTx.GasPrice(), big.NewInt(1)) // Be right in front of the target transaction

			fmt.Printf("Transaction hash: https://bscscan.com/tx/%v", transactions.SignAndSendLegacyTx(context.Background(), tx, client, signer, privateKey))
			os.Exit(0) // Exit the program

		}()
	}
	/*-----------------------------------------------*/

}

func formatStorage(storageHexVals StorageHex) FormattedStorage { // Store the contract data depending on platform
	if cfg.Action == "unicrypt" {
		data := FormattedStorage{
			minBuy:       big.NewInt(1),
			maxBuy:       storages(storageHexVals.maxBuy),
			hardCap:      storages(storageHexVals.hardCap),
			tokenAddress: storages(storageHexVals.tokenAddress),
			startBlock:   storages(storageHexVals.startBlock),
		}
		return data
	} else {
		data := FormattedStorage{
			startTime:    storages(storageHexVals.startTime),
			minBuy:       storages(storageHexVals.minBuy),
			maxBuy:       storages(storageHexVals.maxBuy),
			hardCap:      storages(storageHexVals.hardCap),
			tokenAddress: storages(storageHexVals.tokenAddress),
		}
		return data
	}
}

func storages(hexA string) *big.Int { // Retrieve the contract data
	val, err := client.StorageAt(context.Background(), presaleAddress, common.HexToHash(hexA), nil)
	if err != nil {
		log.Fatal(err)
	}

	data := "0x" + strings.TrimLeft(hex.EncodeToString(val), "0") // Convert from bytes to hex and trim leading zeroes

	big, err := hexutil.DecodeBig(data) // Convert hex data to a big int
	if err != nil {
		fmt.Println(hexA)
		log.Fatal(err)
	}

	return big
}
