package types

import "math/big"

// 大部分类型定义可以复用比特币的
type AccountBalance struct {
	FinalBalance  big.Int `json:"final_balance"`
	NTx           big.Int `json:"n_tx"`
	TotalReceived big.Int `json:"total_received"`
}

type UnspentOutput struct {
	TxHashBigEndian string `json:"tx_hash_big_endian"`
	TxHash          string `json:"tx_hash"`
	TxOutputN       uint64 `json:"tx_output_n"`
	Script          string `json:"script"`
	Value           uint64 `json:"value"`
	ValueHex        string `json:"value_hex"`
	Confirmations   uint64 `json:"confirmations"`
	TxIndex         uint64 `json:"tx_index"`
}

type UnspentOutputList struct {
	Notice         string          `json:"notice"`
	UnspentOutputs []UnspentOutput `json:"unspent_outputs"`
}

type Transaction struct {
	Hash160       string    `json:"hash160"`
	Address       string    `json:"address"`
	NTx           uint64    `json:"n_tx"`
	TotalReceived big.Int   `json:"total_received"`
	TotalSent     big.Int   `json:"total_sent"`
	FinalBalance  big.Int   `json:"final_balance"`
	Txs           []TxsItem `json:"txs"`
}

type TxsItem struct {
	Hash        string      `json:"hash"`
	Version     uint64      `json:"ver"`
	Fee         big.Int     `json:"fee"`
	Time        big.Int     `json:"time"`
	BlockHeight big.Int     `json:"block_height"`
	Inputs      []InputItem `json:"inputs"`
	Out         []OutItem   `json:"out"`
}

type InputItem struct {
	PrevOut PrevOut `json:"prev_out"`
}

type PrevOut struct {
	Addr  string  `json:"addr"`
	Value big.Int `json:"value"`
}

type OutItem struct {
	Value big.Int `json:"value"`
	Addr  string  `json:"addr"`
}

// ... 其他类型定义
