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

type Config struct {
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

type StorageHex struct {
	startTime    string
	minBuy       string
	maxBuy       string
	hardCap      string
	tokenAddress string
	startBlock   string
}

var cfg Config

var client *ethclient.Client
var presaleAddress common.Address

func main() {

	f, err := os.Open("./config.yml")
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	client, err = ethclient.Dial(cfg.Parameters.Host)
	if err != nil {
		log.Fatal(err)
	}
	rpcClient, err := rpc.DialContext(context.Background(), cfg.Parameters.Host)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(cfg.Parameters.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	presaleAddress = common.HexToAddress(cfg.Parameters.PresaleAddress)
	presaleAddressPointer := &presaleAddress
	value := big.NewInt(int64(cfg.Parameters.AmountIn * (math.Pow(10.0, 18.0))))
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	var data []byte
	tx := &types.LegacyTx{
		Nonce: nonce,
		To:    presaleAddressPointer,
		Value: value,
		Data:  data,
	}
	pinkSale := StorageHex{
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

	var storageHexVals StorageHex
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
	hashes := make(chan common.Hash)
	gethClient := gethclient.New(rpcClient)
	_, err = gethClient.SubscribePendingTransactions(context.Background(), hashes)
	if err != nil {
		log.Fatal(err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	signer := types.NewEIP155Signer(chainId)
	one := big.NewInt(1)
	var contributedAmount *big.Int
	bigStorage := formatStorage(storageHexVals)
	pTmp := big.NewInt(1).Mul(bigStorage.hardCap, big.NewInt(cfg.Parameters.PercentFill))
	buyAt := big.NewInt(1).Div(pTmp, big.NewInt(100))
	if cfg.Parameters.Action == "pinksale" || cfg.Parameters.Action == "dxSale" {
		go func() {
			for bigStorage.startTime.Int64()-unix() > 4 {
				fmt.Printf("%v Seconds remaining\n", bigStorage.startTime.Int64()-unix())
				time.Sleep(time.Second)
			}
		}()
		for bigStorage.startTime.Int64()-unix() > 4 {
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
		for bigStorage.startBlock.Uint64()-curBlock > 2 {
		}
	}
	// Be efficient now
	for {
		go func() {

			hash := <-hashes
			pendingTx, _, err := client.TransactionByHash(context.Background(), hash)
			if err != nil {
				return
			}

			if pendingTx.Value().Cmp(bigStorage.minBuy) == -1 && pendingTx.Value().Cmp(bigStorage.maxBuy) == 1 {
				return
			}

			if pendingTx.To() != presaleAddressPointer {
				return
			}
			contributedAmount.Add(pendingTx.Value(), contributedAmount)
			if contributedAmount.Cmp(buyAt) == -1 {
				return
			}

			tx.Gas = pendingTx.Gas() + 1
			tx.GasPrice = one.Add(pendingTx.GasPrice(), big.NewInt(1))
			newTx := types.NewTx(tx)
			signedTx, err := types.SignTx(newTx, signer, privateKey)
			if err != nil {
				log.Fatal(err)
			}
			err = client.SendTransaction(context.Background(), signedTx)
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)

		}()
	}

}

func formatStorage(storageHexVals StorageHex) FormattedStorage {
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

type FormattedStorage struct {
	startTime    *big.Int
	minBuy       *big.Int
	maxBuy       *big.Int
	hardCap      *big.Int
	tokenAddress *big.Int
	startBlock   *big.Int
}

func storages(hexA string) *big.Int {
	val, err := client.StorageAt(context.Background(), presaleAddress, common.HexToHash(hexA), nil)
	if err != nil {
		log.Fatal(err)
	}
	data := "0x" + strings.TrimLeft(hex.EncodeToString(val), "0")

	big, err := hexutil.DecodeBig(data)
	if err != nil {
		log.Fatal(err)
	}
	return big

}

func unix() int64 {
	return time.Now().Unix()
}
