package transactions

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func SignAndSendLegacyTx(ctx context.Context, tx *types.LegacyTx, client *ethclient.Client, s types.Signer, prv *ecdsa.PrivateKey) common.Hash {

	newTx := types.NewTx(tx)

	signedTx, err := types.SignTx(newTx, s, prv)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatal(err)
	}

	return signedTx.Hash()
}

func ToWei(value float64) *big.Int {
	return big.NewInt(int64(value * (math.Pow(10.0, 18.0))))
}

func PrivateKeyExtrapolate(key string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, common.Address) {
	privateKey, err := crypto.HexToECDSA(key) // Load the private key
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public() // Get the public key from private key

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return privateKey, publicKeyECDSA, address
}

func CurrentBlock(ctx context.Context, client *ethclient.Client) *big.Int { // Get current block
	currentBlock, err := client.BlockNumber(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return big.NewInt(int64(currentBlock))
}

func Unix() int64 { // Get the current epoch time in seconds
	return time.Now().Unix()
}
