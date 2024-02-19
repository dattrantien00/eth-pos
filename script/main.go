package main

import (
	"context"
	"crypto/ecdsa"
	"deposit/contract"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

const (
	PRIV_KEY         = "2e0834786285daccd064ca17f1654f67b4aef298acbb82cef9ec422fb4975622"
	DEPOSIT_CONTRACT = "0x4242424242424242424242424242424242424242"

	VALIDATOR_PUBKEY       = "a9f5e7a3e9481a2ed971f1f5fbff3981a9eb16f2a18647b10d5c0ca29401115dcf1a23cd1f59f994e1c955f935151095"
	WITHDRAWAL_CREDENTIALS = "00b343f6a5717be63801cf770b921002d84a9e0b6230316602d6aca161cc5185"
	SIGNATURE              = "8ce6de7fd012a607cc96a698348b3a97f45410a215e5b1b1a9ec6b5c510e7c34d5ae771f5a42cd1ab0d4a801401691d70f79dc114abc378b18559452cda44ceb317ae8abf4b71b2c1a682dac0a14abdfdf6dd31c99816df1c37c4e01a91ac5aa"
	DEPOSIT_DATA_ROOT      = "09fcce618b44239763fa37e376b38b43b1f9d02b39ec922f6b15c46ce4c24e72"
)

func main() {

	cli, _ := ethclient.Dial("http://127.0.0.1:8000")
	fmt.Println(cli.BlockNumber(context.Background()))
	// fmt.Println(cli.BlockByNumber(context.Background(), big.NewInt(-3)))
	ins, _ := contract.NewStore(common.HexToAddress(DEPOSIT_CONTRACT), cli)

	root, _ := ins.GetDepositRoot(nil)
	fmt.Println(hex.EncodeToString(root[:]))
	auth, _ := GetAuth(cli, PRIV_KEY)

	pubk, _ := hex.DecodeString(VALIDATOR_PUBKEY)
	cridential, _ := hex.DecodeString(WITHDRAWAL_CREDENTIALS)
	sig, _ := hex.DecodeString(SIGNATURE)
	depositRoot, _ := hex.DecodeString(DEPOSIT_DATA_ROOT)
	de := [32]byte{}
	copy(de[:], depositRoot)

	tx, err := ins.Deposit(auth, pubk, cridential, sig, de)
	fmt.Println(tx.Hash().Hex())
	fmt.Println(err)
}
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
	auth.Value = FloatStringToBigInt("16", 18)
	fmt.Println(auth.Value)
	return auth, nil
}

func FloatStringToBigInt(amount string, decimals int) *big.Int {
	fAmount, _ := new(big.Float).SetString(amount)
	fi, _ := new(big.Float).Mul(fAmount, big.NewFloat(math.Pow10(decimals))).Int(nil)
	return fi
}
