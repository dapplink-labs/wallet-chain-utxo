package base

import "time"

// FeeRateResponse 费率响应
type FeeRateResponse struct {
	Name             string    `json:"name"`
	Height           int64     `json:"height"`
	Hash             string    `json:"hash"`
	Time             time.Time `json:"time"`
	LatestURL        string    `json:"latest_url"`
	PreviousHash     string    `json:"previous_hash"`
	PreviousURL      string    `json:"previous_url"`
	PeerCount        int       `json:"peer_count"`
	UnconfirmedCount int       `json:"unconfirmed_count"`
	HighFeePerKb     int64     `json:"high_fee_per_kb"`
	MediumFeePerKb   int64     `json:"medium_fee_per_kb"`
	LowFeePerKb      int64     `json:"low_fee_per_kb"`
	LastForkHeight   int64     `json:"last_fork_height"`
	LastForkHash     string    `json:"last_fork_hash"`
}

// TransactionResponse 交易响应
type TransactionResponse struct {
	BlockHash     string    `json:"block_hash"`
	BlockHeight   int64     `json:"block_height"`
	BlockIndex    int       `json:"block_index"`
	Hash          string    `json:"hash"`
	Addresses     []string  `json:"addresses"`
	Total         int64     `json:"total"`
	Fees          int64     `json:"fees"`
	Size          int       `json:"size"`
	Preference    string    `json:"preference"`
	RelayedBy     string    `json:"relayed_by"`
	Confirmed     time.Time `json:"confirmed"`
	Received      time.Time `json:"received"`
	Ver           int       `json:"ver"`
	DoubleSpend   bool      `json:"double_spend"`
	VinSz         int       `json:"vin_sz"`
	VoutSz        int       `json:"vout_sz"`
	Confirmations int64     `json:"confirmations"`
	Confidence    float64   `json:"confidence"`
	Inputs        []Input   `json:"inputs"`
	Outputs       []Output  `json:"outputs"`
}

type Input struct {
	PrevHash    string   `json:"prev_hash"`
	OutputIndex int      `json:"output_index"`
	Script      string   `json:"script"`
	OutputValue int64    `json:"output_value"`
	Sequence    int64    `json:"sequence"`
	Addresses   []string `json:"addresses"`
	ScriptType  string   `json:"script_type"`
	Age         int64    `json:"age"`
}

type Output struct {
	Value      int64    `json:"value"`
	Script     string   `json:"script"`
	Addresses  []string `json:"addresses"`
	ScriptType string   `json:"script_type"`
}

// AccountResponse 账户余额响应
type AccountResponse struct {
	Address            string `json:"address"`
	TotalReceived      int64  `json:"total_received"`
	TotalSent          int64  `json:"total_sent"`
	Balance            int64  `json:"balance"`
	UnconfirmedBalance int64  `json:"unconfirmed_balance"`
	FinalBalance       int64  `json:"final_balance"`
	NTx                int    `json:"n_tx"`
	UnconfirmedNTx     int    `json:"unconfirmed_n_tx"`
	FinalNTx           int    `json:"final_n_tx"`
}

// UnspentOutputResponse UTXO响应
type UnspentOutputResponse struct {
	Address            string  `json:"address"`
	TotalReceived      int64   `json:"total_received"`
	TotalSent          int64   `json:"total_sent"`
	Balance            int64   `json:"balance"`
	UnconfirmedBalance int64   `json:"unconfirmed_balance"`
	FinalBalance       int64   `json:"final_balance"`
	NTx                int     `json:"n_tx"`
	UnconfirmedNTx     int     `json:"unconfirmed_n_tx"`
	FinalNTx           int     `json:"final_n_tx"`
	TxRefs             []TxRef `json:"txrefs"`
	TxURL              string  `json:"tx_url"`
}

type TxRef struct {
	TxHash        string    `json:"tx_hash"`
	BlockHeight   int64     `json:"block_height"`
	TxInputN      int       `json:"tx_input_n"`
	TxOutputN     int       `json:"tx_output_n"`
	Value         int64     `json:"value"`
	RefBalance    int64     `json:"ref_balance"`
	Spent         bool      `json:"spent"`
	Confirmations int64     `json:"confirmations"`
	Confirmed     time.Time `json:"confirmed"`
	DoubleSpend   bool      `json:"double_spend"`
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
}

type ScriptPubKey struct {
	Asm     string `json:"asm"`
	Hex     string `json:"hex"`
	Desc    string `json:"desc"`
	Address string `json:"addresses"`
	Type    string `json:"type"`
}

type Vout struct {
	Value        interface{}  `json:"value"`
	N            uint64       `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptpubkey"`
}

// ResultRawData 交易原始数据
type ResultRawData struct {
	TxId          string `json:"txid"`
	Hash          string `json:"hash"`
	Version       uint64 `json:"version"`
	Size          uint64 `json:"size"`
	VSize         uint64 `json:"vsize"`
	Weight        uint64 `json:"weight"`
	LockTime      uint64 `json:"locktime"`
	Vin           []Vin  `json:"vin"`
	Vout          []Vout `json:"vout"`
	Hex           string `json:"hex"`
	Blockhash     string `json:"blockhash"`
	Confirmations uint64 `json:"confirmations"`
	BlockTime     uint64 `json:"blocktime"`
	Time          uint64 `json:"time"`
}

// BlockResponse 区块响应
type BlockResponse struct {
	Hash         string    `json:"hash"`
	Height       int64     `json:"height"`
	Chain        string    `json:"chain"`
	Total        int64     `json:"total"`
	Fees         int64     `json:"fees"`
	Size         int       `json:"size"`
	VSize        int       `json:"vsize"`
	Ver          int       `json:"ver"`
	Time         time.Time `json:"time"`
	ReceivedTime time.Time `json:"received_time"`
	CoinbaseAddr string    `json:"coinbase_addr"`
	RelayedBy    string    `json:"relayed_by"`
	Bits         int       `json:"bits"`
	Nonce        int64     `json:"nonce"`
	NTx          int       `json:"n_tx"`
	PrevBlock    string    `json:"prev_block"`
	MrklRoot     string    `json:"mrkl_root"`
	TxIds        []string  `json:"txids"`
	Depth        int       `json:"depth"`
	PrevBlockURL string    `json:"prev_block_url"`
	TxURL        string    `json:"tx_url"`
	NextTxIds    string    `json:"next_txids"`
}
