package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
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

	cli, _ := ethclient.Dial("https://l1testnet.trustkeys.network")
	fmt.Println(cli.BlockNumber(context.Background()))

	number := big.NewInt(657400)
	for {
		block,_ := cli.BlockByNumber(context.Background(),number)
		fmt.Println(len(block.Body().Transactions))
		number.Add(number,big.NewInt(1))
	}
	
}
