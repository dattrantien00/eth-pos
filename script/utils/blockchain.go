package utils

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func GetAuth(client *ethclient.Client, privateKey string) (*bind.TransactOpts, error) {

	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privKey, chainID)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	pubkey := privKey.Public()
	pubkeyECDSA := pubkey.(*ecdsa.PublicKey)
	addr := crypto.PubkeyToAddress(*pubkeyECDSA)

	nonce, _ := client.PendingNonceAt(context.Background(), addr)
	//gasPrice, err := cli.SuggestGasPrice(context.Background())
	// gasPrice := big.NewInt(consts.BSC_GAS_PRICE)
	auth.GasLimit = 8 * 1e6
	auth.GasPrice = nil
	auth.Nonce = big.NewInt(int64(nonce))
	// auth.Value = big.NewInt(32000000000)
	auth.Value = big.NewInt(0)
	fmt.Println(auth.Value)
	return auth, nil
}
