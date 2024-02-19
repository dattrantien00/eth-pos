package randomtransfer

import (
	"context"
	"crypto/ecdsa"
	"deposit/utils"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
	"gitlab.123xe.vn/cryptocurrency/smart-contract-gen/NFT/ERC20"
)

type TransferInfo struct {
	Address            string
	ReceiverPrivateKey string
	Value              string
	Token              string
}

var (
	mu           sync.Mutex
	erc20        *ERC20.ERC20
	tokenDecimal uint8
)

func genKeyToFile() *xlsx.File {
	cli, err := ethclient.Dial(rpc)
	_ = cli
	if err != nil {
		logrus.Error(err)
		return nil
	}

	file := createRandomTransferFile()
	if file == nil {
		return nil
	}

	sender := addressByPrivKey(privateKey)

	nonceStart, _ := cli.NonceAt(context.Background(), common.HexToAddress(sender), nil)
	fmt.Println(nonceStart)
	if token != "" {
		erc20, err = ERC20.NewERC20(common.HexToAddress(token), cli)
		if err != nil {
			logrus.Errorln(err)
			return nil
		}
		tokenDecimal, err = erc20.Decimals(nil)
		if err != nil {
			logrus.Errorln(err)
			return nil
		}
	}

	for i := 0; i < int(receiverNumber); i++ {
		// wg.Add(1)
		nonce := nonceStart + uint64(i)
		// go func() {
		randomAndTransfer(cli, file.Sheets[0], nonce)
		// wg.Done()
		// }()
		// time.Sleep(2*time.Second)

	}

	err = file.Save("./tmp/" + fileName + ".xlsx")
	if err != nil {
		fmt.Println("Lỗi khi lưu file Excel:", err)
		return nil
	}
	return file
}

func transferRandom() {
	cli, err := ethclient.Dial(rpc)
	_ = cli
	if err != nil {
		logrus.Error(err)
		return
	}

	file := genKeyToFile()
	//cross transfer user in file
	fmt.Println("======================Cross Transfer==================")
	time.Sleep(5 * time.Second)

	loopCrossTransfer(cli, file.Sheets[0])
}

func createRandomTransferFile() *xlsx.File {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		logrus.Error(err)
		return nil
	}
	row := sheet.AddRow()
	row.WriteStruct(&TransferInfo{"Address", "Receiver PrivateKey", "Value", "Token"}, -1)
	return file
}

func importFileAndTransfer(){
	cli, err := ethclient.Dial(rpc)
	if err != nil{
		logrus.Errorln(err)
		return 
	}
	file,err := xlsx.OpenFile("./tmp/" + fileName + ".xlsx")
	if err != nil{
		logrus.Errorln(err)
	}
	if len(file.Sheets)==0{
		logrus.Errorln("file hasn't any sheet")
	}
	loopCrossTransfer(cli,file.Sheets[0])
}

func loopCrossTransfer(client *ethclient.Client, sheet *xlsx.Sheet){
	go func() {
		for {
			crossTransfer(client, sheet)
			time.Sleep(time.Duration(timeBetweenTrans) * time.Millisecond)
		}
	}()
	time.Sleep(time.Duration(timeCrossTransfer) * time.Second)
}
func randomAndTransfer(client *ethclient.Client, sheet *xlsx.Sheet, nonce uint64) {
	receiverPrvK := genPrivateKey()
	address := addressByPrivKey(receiverPrvK)
	amount := rand.Float64() * maxValuePerTransaction
	data := TransferInfo{}
	if token == "" {
		err := transferCoinbase(client, privateKey, address, amount, nonce)
		if err != nil {
			return
		}
		data = TransferInfo{
			Address:            address,
			ReceiverPrivateKey: receiverPrvK,
			Value:              fmt.Sprintf("%f", amount),
			Token:              "0x0",
		}
	} else {
		err := transferErc20(client, address, amount, nonce)
		if err != nil {
			return
		}
		data = TransferInfo{
			Address:            address,
			ReceiverPrivateKey: receiverPrvK,
			Value:              fmt.Sprintf("%f", amount),
			Token:              token,
		}
	}

	writeToXlsx(sheet, data)

}

func writeToXlsx(sheet *xlsx.Sheet, data TransferInfo) {
	mu.Lock()
	defer mu.Unlock()
	row := sheet.AddRow()
	row.WriteStruct(&data, -1)

}

func transferCoinbase(client *ethclient.Client, prvKey string, receiver string, amount float64, nonce uint64) error {

	privK, err := crypto.HexToECDSA(prvKey)
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	value := utils.FloatStringToBigInt(amount, 18)
	gasLimit := uint64(21000)         // in units
	tip := big.NewInt(2000000000)     // maxPriorityFeePerGas = 2 Gwei
	feeCap := big.NewInt(20000000000) // maxFeePerGas = 20 Gwei
	if err != nil {
		log.Fatal(err)
		return err
	}

	toAddress := common.HexToAddress(receiver)
	var data []byte

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
		return err
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: feeCap,
		GasTipCap: tip,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
		Data:      data,
	})

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privK)

	if err != nil {
		log.Fatal(err)
		return err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	fmt.Printf("Transaction hash: %s\n", signedTx.Hash().Hex())
	return nil
}

func transferErc20(client *ethclient.Client, receiver string, amount float64, nonce uint64) error {
	auth, err := utils.GetAuth(client, privateKey)
	if err != nil {
		return err
	}
	auth.Nonce = new(big.Int).SetUint64(nonce)

	_, err = erc20.Transfer(auth, common.HexToAddress(receiver), utils.FloatStringToBigInt(amount, int(tokenDecimal)))
	if err != nil {
		logrus.Errorln(err)
		return err
	}
	return nil
}
func genPrivateKey() string {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	return hexutil.Encode(privateKeyBytes)[2:]
}

func addressByPrivKey(key string) string {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return fromAddress.Hex()
}

func randomIndex(min, max int) (senderIndex, toIndex int) {
	senderIndex = rand.Intn(max-min) + min
	toIndex = rand.Intn(max-min) + min
	return
}

func crossTransfer(client *ethclient.Client, sheet *xlsx.Sheet) {

	senderIndex, toIndex := randomIndex(1, sheet.MaxRow)
	senderKey := sheet.Rows[senderIndex].Cells[1].Value
	senderAddr := sheet.Rows[senderIndex].Cells[0].Value
	toAddr := sheet.Rows[toIndex].Cells[0].Value

	sBal, err := client.PendingBalanceAt(context.Background(), common.HexToAddress(senderAddr))
	if err != nil {
		return
	}
	sBalF, _ := utils.WeiToEther(sBal).Float64()
	amount := rand.Float64() * sBalF

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(senderAddr))
	transferCoinbase(client, senderKey, toAddr, amount, nonce)
}
