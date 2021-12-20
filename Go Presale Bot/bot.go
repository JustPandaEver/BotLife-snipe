package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"gopkg.in/yaml.v2"
)

type Config struct { // Initialize struct for config options
	Parameters struct {
		Host           string  `yaml:"host"`
		PCSAddress     string  `yaml:"pcsaddress"`
		PrivateKey     string  `yaml:"privatekey"`
		AmountIn       float64 `yaml:"amountin"`
		PresaleAddress string  `yaml:"presaleaddress"`
		PercentFill    int64   `yaml:"percentfill"`
		Action         string  `yaml:"action"`
	} `yaml:"parameters"`
}

type StorageHex struct { // Initialize struct for storing contract data locations in hex
	startTime    string
	minBuy       string
	maxBuy       string
	hardCap      string
	tokenAddress string
	startBlock   string
}

// Initialize global variables
var cfg Config
var client *ethclient.Client
var presaleAddress common.Address

func main() {
	/*---------- Initialize -----------*/
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

	client, err = ethclient.Dial(cfg.Parameters.Host) // Connect to node
	if err != nil {
		log.Fatal(err)
	}

	rpcClient, err := rpc.DialContext(context.Background(), cfg.Parameters.Host) // Connect to node for RPC
	if err != nil {
		log.Fatal(err)
	}

	gethClient := gethclient.New(rpcClient) // Initialize the GETH client

	privateKey, err := crypto.HexToECDSA(cfg.Parameters.PrivateKey) // Load the private key
	if err != nil {
		log.Fatal(err)
	}
	/*------------------------*/

	/*--------Initialize Key----------*/
	publicKey := privateKey.Public() // Get the public key from private key

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	/*---------------------------------*/

	/*----------Initialize Transaction------------*/
	presaleAddress = common.HexToAddress(cfg.Parameters.PresaleAddress) // Init presale address
	presaleAddressPointer := &presaleAddress                            // Create a pointer type

	value := big.NewInt(int64(cfg.Parameters.AmountIn * (math.Pow(10.0, 18.0)))) // Convert to wei

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA) // Get the wallet address from public key

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
	var contributedAmount *big.Int
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
	if cfg.Parameters.Action == "pinksale" {
		storageHexVals = pinkSale
		tx.Data = common.Hex2Bytes("0x")
	} else if cfg.Parameters.Action == "dxSale" {
		storageHexVals = dxSale
		tx.Data = common.Hex2Bytes("0x")
	} else if cfg.Parameters.Action == "unicrypt" {
		storageHexVals = unicrypt
		tx.Data = common.Hex2Bytes("0xf868e7660000000000000000000000000000000000000000000000000000000000000000")
	} else {
		log.Fatalf("Unsupported action: %v", cfg.Parameters.Action)
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
	bigStorage := formatStorage(storageHexVals) // Retrieve all data from the presale contract
	/*--------------------------------*/

	pTmp := big.NewInt(1).Mul(bigStorage.hardCap, big.NewInt(cfg.Parameters.PercentFill))
	buyAt := big.NewInt(1).Div(pTmp, big.NewInt(100)) // Set the value to submit transaction at

	/*-----------------Wait------------*/
	if cfg.Parameters.Action == "pinksale" || cfg.Parameters.Action == "dxSale" {
		go func() { // Create a goroutine so the sleep is non blocking
			for bigStorage.startTime.Int64()-unix() > 4 {
				fmt.Printf("%v Seconds remaining\n", bigStorage.startTime.Int64()-unix())
				time.Sleep(time.Second)
			}
		}()
		for bigStorage.startTime.Int64()-unix() > 4 { // Block until # seconds before presale start
		}
	} else if cfg.Parameters.Action == "unicrypt" {
		curBlock, err := client.BlockNumber(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			for bigStorage.startBlock.Uint64()-curBlock > 2 {
				fmt.Printf("%v Blocks remaining\n", bigStorage.startBlock.Uint64()-curBlock)
				time.Sleep(time.Second * 3)
			}
		}()
		for bigStorage.startBlock.Uint64()-curBlock > 2 { // Block until # blocks before presale start
		}
	}
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

			// Instead of doing everything in the if statement, return if a condition is not met
			// Start with easiest, this avoids unnecessary computations
			// This should help with speed or CPU load
			if pendingTx.To() != presaleAddressPointer { // Check if the target address is the presale contract
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
			newTx := types.NewTx(tx)                                   // Create the transaction

			signedTx, err := types.SignTx(newTx, signer, privateKey) // Sign the transaction
			if err != nil {
				log.Fatal(err)
			}

			err = client.SendTransaction(context.Background(), signedTx) // Send the transaction
			if err != nil {
				log.Fatal(err)
			}

			os.Exit(0) // Exit the program

		}()
	}
	/*-----------------------------------------------*/

}

type FormattedStorage struct { // Struct for storing contract data
	startTime    *big.Int
	minBuy       *big.Int
	maxBuy       *big.Int
	hardCap      *big.Int
	tokenAddress *big.Int
	startBlock   *big.Int
}

func formatStorage(storageHexVals StorageHex) FormattedStorage { // Store the contract data depending on platform
	if cfg.Parameters.Action == "unicrypt" {
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
		log.Fatal(err)
	}

	return big
}

func unix() int64 { // Get the current epoch time in seconds
	return time.Now().Unix()
}
