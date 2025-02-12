package dogecoin

import (
	"math"
	"math/big"
	"strconv"

	"github.com/dapplink-labs/wallet-chain-utxo/rpc/utxo"
	"github.com/shopspring/decimal"
)

const (
	dogeDecimals = 8      // Dogecoin也是8位小数
	MinFeeRate   = 0.01   // Dogecoin最小费率 (DOGE/kb)
	DogePerKb    = 100000 // 每kb的聪数
)

type DecodeTxRes struct {
	Hash       string
	SignHashes [][]byte
	Vins       []*utxo.Vin
	Vouts      []*utxo.Vout
	CostFee    *big.Int
}

func dogeToSatoshi(dogeCount float64) *big.Int {
	amount := strconv.FormatFloat(dogeCount, 'f', -1, 64)
	amountDm, _ := decimal.NewFromString(amount)
	tenDm := decimal.NewFromFloat(math.Pow(10, float64(dogeDecimals)))
	satoshiDm, _ := big.NewInt(0).SetString(amountDm.Mul(tenDm).String(), 10)
	return satoshiDm
}

// 计算交易大小的预估费用
func calculateFee(txSize int64, feeRate float64) *big.Int {
	// 将kb转换为字节的费率
	feePerByte := feeRate / 1024.0

	// 计算总费用（DOGE）
	totalFee := feePerByte * float64(txSize)

	// 确保不低于最小费率
	if totalFee < MinFeeRate {
		totalFee = MinFeeRate
	}

	// 转换为聪
	return dogeToSatoshi(totalFee)
}

// 估算交易大小
func estimateTxSize(numInputs, numOutputs int) int64 {
	// 基础交易大小
	baseSize := int64(10)
	// 每个输入大小
	inputSize := int64(148)
	// 每个输出大小
	outputSize := int64(34)

	return baseSize + (inputSize * int64(numInputs)) + (outputSize * int64(numOutputs))
}
