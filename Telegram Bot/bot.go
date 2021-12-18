package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	pcsRouter "tgbot/contracts"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	Parameters struct {
		Host       string  `yaml:"host"`
		Contracts  string  `yaml:"contracts"`
		Bought     string  `yaml:"bought"`
		PCSAddress string  `yaml:"pcsaddress"`
		BNBAddress string  `yaml:"bnbaddress"`
		PrivateKey string  `yaml:"privatekey"`
		AmountIn   float64 `yaml:"amountin"`
		GasLimit   uint64  `yaml:"gaslimit"`
		GasPrice   int64   `yaml:"gasprice"`
	} `yaml:"parameters"`
}

func main() {
	fmt.Println("Starting...")
	bot, err := tgbotapi.NewBotAPI("5057473003:AAGZj29_0TfCHqQHyd3oMxNNMRQg4doIV-A")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open("./config.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	filePath := cfg.Parameters.Contracts
	updates := bot.GetUpdatesChan(u)
	fmt.Println("Listening for updates now!")
	for update := range updates {
		f, _ := os.OpenFile(filePath, os.O_RDWR, 0777)
		writer := csv.NewWriter(f)
		contract := getContract(update)
		if !exists(contract, f) {
			writer.Write([]string{contract})
			writer.Flush()
			f.Close()
			buy(contract, cfg)
		}
	}

}

func exists(contract string, f *os.File) bool {
	reader := csv.NewReader(f)
	contracts, _ := reader.ReadAll()
	if len(contracts) != 0 {
		for _, existingContract := range contracts {
			if runtime.GOOS == "windows" {
				existingContract[0] = strings.TrimRight(existingContract[0], "\r\n")
			} else {
				existingContract[0] = strings.TrimRight(existingContract[0], "\n")
			}
			if contract == existingContract[0] {
				return true
			}
		}
		return false
	} else {
		return true
	}
}

func getContract(update tgbotapi.Update) string {
	regex := regexp.MustCompile(`0x.*`)
	contractExt := regex.FindAllString(update.Message.Text, -1)[0]
	runes := []rune(contractExt)
	contract := string(runes[0:42])
	return contract
}

func buy(contract string, cfg Config) bool {
	client, err := ethclient.Dial(cfg.Parameters.Host)
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(contract)
	pcsRouterAddress := common.HexToAddress(cfg.Parameters.PCSAddress)
	bnbAddress := common.HexToAddress(cfg.Parameters.BNBAddress)
	privateKey, err := crypto.HexToECDSA(cfg.Parameters.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, (*big.Int)(chainID))
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	txNonce := big.NewInt(int64(nonce))
	amountIn := big.NewInt(int64(cfg.Parameters.AmountIn * (math.Pow(10.0, 18.0))))
	auth.GasLimit = cfg.Parameters.GasLimit // in units
	auth.GasPrice = big.NewInt(cfg.Parameters.GasPrice * (1000000000))
	instance, err := pcsRouter.NewPcsRouter(pcsRouterAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	opts := &bind.CallOpts{
		Pending: false,
		Context: context.Background(),
	}
	path := []common.Address{bnbAddress, contractAddress}
	out, err := instance.GetAmountsOut(opts, amountIn, path)
	if err != nil {
		log.Fatal(err)
	}
	minAmtOut := out[1].Mul(out[1], big.NewInt(85)).Div(out[1], big.NewInt(100))
	transactOps := &bind.TransactOpts{
		From:     fromAddress,
		Nonce:    txNonce,
		Signer:   auth.Signer,
		Value:    amountIn,
		GasPrice: auth.GasPrice,
		GasLimit: auth.GasLimit,
		Context:  context.Background(),
	}
	deadline := big.NewInt(time.Now().Unix() + 1200)
	tx, err := instance.SwapExactETHForTokensSupportingFeeOnTransferTokens(transactOps, minAmtOut, path, fromAddress, deadline)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transaction receipt: %v", tx.Hash())
	f, _ := os.OpenFile(cfg.Parameters.Bought, os.O_RDWR, 0777)
	writer := csv.NewWriter(f)
	writer.Write([]string{fmt.Sprintf("Transaction receipt: %v", tx.Hash())})
	writer.Flush()
	f.Close()
	return true
}
