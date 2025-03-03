package types

type Items struct {
	TXID          string  `json:"txid"`
	Version       uint64  `json:"version"`
	LockTime      uint64  `json:"locktime"`
	Vin           []Vin   `json:"vin"`
	Vout          []Vout  `json:"vout"`
	BlockHash     string  `json:"blockhash"`
	BlockHeight   uint64  `json:"blockheight"`
	Confirmations uint64  `json:"confirmations"`
	Time          uint64  `json:"time"`
	BlockTime     uint64  `json:"blocktime"`
	IsCoinBase    bool    `json:"isCoinBase"`
	ValueOut      float64 `json:"valueOut"`
	Size          uint64  `json:"size"`
	Fees          float64 `json:"fees"`
}

type Transaction struct {
	TotalItems uint64  `json:"totalItems"`
	From       uint64  `json:"from"`
	To         uint64  `json:"to"`
	Items      []Items `json:"items"`
}

type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

type Vin struct {
	TxId        string    `json:"txid"`
	Vout        uint64    `json:"vout"`
	ScriptSig   ScriptSig `json:"scriptSig"`
	Sequence    uint64    `json:"sequence"`
	TxInWitness []string  `json:"txinwitness"`
	Address     string    `json:"addr"`
	Coinbase    string    `json:"coinbase"`
}

type ScriptPubKey struct {
	Asm     string   `json:"asm"`
	Hex     string   `json:"hex"`
	Desc    string   `json:"desc"`
	Address []string `json:"addresses"`
	Type    string   `json:"type"`
}

type Vout struct {
	Value        string       `json:"value"`
	N            uint64       `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptpubkey"`
}

type UnspentOutput struct {
	Address       string  `json:"address"`
	Txid          string  `json:"txid"`
	Vout          uint64  `json:"vout"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	Amount        float64 `json:"amount"`
	Satoshis      uint64  `json:"satoshis"`
	Height        uint64  `json:"height"`
	Confirmations uint64  `json:"confirmations"`
}

type BlockGetInfo struct {
	BlockInfo BlockInfo `json:"info"`
}

type BlockInfo struct {
	Version         uint64  `json:"version"`
	Protocolversion uint64  `json:"protocolversion"`
	Blocks          int64   `json:"blocks"`
	Timeoffset      uint64  `json:"timeoffset"`
	Connections     uint64  `json:"connections"`
	Proxy           string  `json:"proxy"`
	Difficulty      float64 `json:"difficulty"`
	Testnet         bool    `json:"testnet"`
	Relayfee        float64 `json:"relayfee"`
	Errors          string  `json:"errors"`
	Network         string  `json:"network"`
}

type BlockHash struct {
	BlockHash string `json:"blockHash"`
}

type Block struct {
	Hash              string   `json:"hash"`
	Size              uint64   `json:"size"`
	Height            int64    `json:"height"`
	Version           uint64   `json:"version"`
	Merkleroot        string   `json:"merkleroot"`
	Tx                []string `json:"tx"`
	Time              uint64   `json:"time"`
	Nonce             string   `json:"nonce"`
	Solution          string   `json:"solution"`
	Bits              string   `json:"bits"`
	Difficulty        float64  `json:"difficulty"`
	Chainwork         string   `json:"chainwork"`
	Confirmations     uint64   `json:"confirmations"`
	Previousblockhash string   `json:"previousblockhash"`
	Nextblockhash     string   `json:"nextblockhash"`
	Reward            float64  `json:"reward"`
	IsMainChain       bool     `json:"isMainChain"`
	PoolInfo          PoolInfo `json:"poolInfo"`
	// Cert            []Cert `json:"cert"`
	ScTxsCommitment string `json:"scTxsCommitment"`
}

type PoolInfo struct {
	PoolName string `json:"poolName"`
	Url      string `json:"url"`
}

type TxInput struct {
	Output           TxOutput `json:"output"`
	Script           string   `json:"script"`
	Sequence         string   `json:"sequence"`
	PrevScriptPubKey string   `json:"prevScriptPubKey"`
}

type TxOutput struct {
	Hash string `json:"hash"`
	Vout uint32 `json:"vout"`
}

type TxOut struct {
	Satoshis uint64 `json:"satoshis"`
	Script   string `json:"script"`
}

type RawTransaction struct {
	Version  uint32    `json:"version"`
	Locktime uint32    `json:"locktime"`
	Ins      []TxInput `json:"ins"`
	Outs     []TxOut   `json:"outs"`
}

type SendRawTransaction struct {
	RawTx string `json:"rawtx"`
}

type SendRawTransactionResponse struct {
	TxId string `json:"txid"`
}

type TXObjectInputParam struct {
	BlockHeight int64      `json:"blockHeight"`
	BlockHash   string     `json:"blockHash"`
	ParamIns    []ParamIn  `json:"paramIn"`
	ParamOuts   []ParamOut `json:"paramOut"`
}

type ParamIn struct {
	TXID         string `json:"txid"`
	Vout         uint64 `json:"vout"`
	ScriptPubKey string `json:"scriptPubKey"`
}

type ParamOut struct {
	Address  string `json:"address"`
	Satoshis uint64 `json:"satoshis"`
}

type TXObject struct {
	LockTime uint64 `json:"locktime"`
	Version  uint32 `json:"version"`
	Ins      []Ins  `json:"ins"`
	Outs     []Outs `json:"outs"`
}

type Ins struct {
	Output           Output `json:"output"`
	Script           string `json:"script"`
	PrevScriptPubKey string `json:"prevScriptPubKey"`
	Sequence         string `json:"sequence"`
}

type Outs struct {
	Satoshis uint64 `json:"satoshis"`
	Script   string `json:"script"`
}

type Output struct {
	Hash string `json:"hash"`
	Vout uint64 `json:"vout"`
}
