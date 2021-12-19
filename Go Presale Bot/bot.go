package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://bsc-dataseed1.defibit.io/")
	if err != nil {
		log.Fatal(err)
	}
	address := common.HexToAddress("0x9aeCccc1e5E4172c849895CA6351e5DD359e1C72")
	num := []byte("131")
	key := common.BytesToHash(num)
	val, err := client.StorageAt(context.Background(), address, key, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(val)
}
