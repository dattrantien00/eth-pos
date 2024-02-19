package utils

import (
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

func FloatStringToBigInt(amount float64, decimals int) *big.Int {
	fAmount := new(big.Float).SetFloat64(amount)
	fi, _ := new(big.Float).Mul(fAmount, big.NewFloat(math.Pow10(decimals))).Int(nil)
	return fi
}

func WeiToEther(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.Ether))
}
