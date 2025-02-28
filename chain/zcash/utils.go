package zcash

import (
	"math"
	"math/big"
	"strconv"

	"github.com/shopspring/decimal"

	"github.com/dapplink-labs/wallet-chain-utxo/rpc/utxo"
)

const (
	zecDecimals = 8
)

type DecodeTxRes struct {
	Hash       string
	SignHashes [][]byte
	Vins       []*utxo.Vin
	Vouts      []*utxo.Vout
	CostFee    *big.Int
}

func zecToSatoshi(zecCount float64) *big.Int {
	amount := strconv.FormatFloat(zecCount, 'f', -1, 64)
	amountDm, _ := decimal.NewFromString(amount)
	tenDm := decimal.NewFromFloat(math.Pow(10, float64(zecDecimals)))
	satoshiDm, _ := big.NewInt(0).SetString(amountDm.Mul(tenDm).String(), 10)
	return satoshiDm
}
