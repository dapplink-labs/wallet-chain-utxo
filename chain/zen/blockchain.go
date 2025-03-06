package zen

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/log"
	gresty "github.com/go-resty/resty/v2"
	"github.com/pkg/errors"

	"github.com/dapplink-labs/chain-explorer-api/common/gas_fee"
	"github.com/dapplink-labs/wallet-chain-utxo/chain/zen/types"
	"github.com/dapplink-labs/wallet-chain-utxo/chain/zen/util"
)

var errBlockChainHTTPError = errors.New("blockchain http error")

type zenTPClient struct {
	client *gresty.Client
}

func NewBlockChainClient(url string) (*zenTPClient, error) {
	// validate if blockchain url is provided or not
	if url == "" {
		return nil, fmt.Errorf("blockchain URL cannot be empty")
	}

	client := gresty.New()
	log.Info("url", url)
	client.SetHostURL(url)
	client.OnAfterResponse(func(c *gresty.Client, r *gresty.Response) error {
		statusCode := r.StatusCode()
		if statusCode >= 400 {
			method := r.Request.Method
			url := r.Request.URL
			return fmt.Errorf("%d cannot %s %s: %w", statusCode, method, url, errBlockChainHTTPError)
		}
		return nil
	})
	return &zenTPClient{
		client: client,
	}, nil
}

func (c *zenTPClient) GetAccountBalance(address string) (string, error) {
	var accountBalance big.Int
	response, err := c.client.R().
		SetResult(&accountBalance).
		Get("/addr/" + address + "/balance")
	if err != nil {
		return "", fmt.Errorf("cannot get account balance: %w", err)
	}
	if response.StatusCode() != 200 {
		return "", errors.New("get account balance fail")
	}
	return accountBalance.String(), nil
}

func (c *zenTPClient) GetFee() (*gas_fee.GasEstimateFeeResponse, error) {
	return &gas_fee.GasEstimateFeeResponse{
		BestTransactionFee:    "0.0001",
		BestTransactionFeeSat: "10000",
		SlowGasPrice:          "",
		StandardGasPrice:      "",
		RapidGasPrice:         "",
	}, nil
}

func (c *zenTPClient) GetAccountUtxo(address string) ([]types.UnspentOutput, error) {
	var utxoUnspentList []types.UnspentOutput
	response, err := c.client.R().
		SetResult(&utxoUnspentList).
		Get("/addr/" + address + "/utxo")
	if err != nil {
		return nil, fmt.Errorf("cannot utxo fail: %w", err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get account counter fail")
	}
	return utxoUnspentList, nil
}

func (c *zenTPClient) GetTransactionsByAddress(address string, pageSize, page uint32) (*types.Transaction, error) {
	var transactionList types.Transaction
	response, err := c.client.R().
		SetResult(&transactionList).
		SetBody(map[string]interface{}{
			"addrs":       address,
			"from":        pageSize * (page - 1),
			"to":          pageSize * page,
			"noAsm":       1,
			"noScriptSig": 1,
			"noSpent":     1,
		}).
		Post("/addrs/txs")

	if err != nil {
		return nil, fmt.Errorf("cannot utxo fail: %w", err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get account counter fail")
	}
	return &transactionList, nil
}

func (c *zenTPClient) GetTransactionsByHash(txHash string) (*types.Items, error) {
	var transaction types.Items
	response, err := c.client.R().
		SetResult(&transaction).
		Get("/tx/" + txHash)

	if err != nil {
		return nil, fmt.Errorf("cannot utxo fail: %w", err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get account counter fail")
	}
	return &transaction, nil
}

func (c *zenTPClient) GetBlockByHeight(blockHeight int64) (*types.Block, error) {
	if blockHeight == 0 {
		blockInfo, err := c.GetLatestBlock()
		if err != nil {
			return nil, err
		}
		blockHeight = int64(blockInfo.Blocks)
	}
	txBlockHash, err := c.GetBlockHashByHeight(blockHeight)
	if err != nil {
		return nil, err
	}
	block, err := c.GetBlockByHash(txBlockHash)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (c *zenTPClient) GetBlockByHash(blockHash string) (*types.Block, error) {
	var block types.Block
	response, err := c.client.R().
		SetResult(&block).
		Get("/block/" + blockHash)
	if err != nil {
		return nil, fmt.Errorf("cannot get block by height: %w", err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get block by height fail")
	}
	return &block, nil
}

func (c *zenTPClient) GetLatestBlock() (*types.BlockInfo, error) {
	var block types.BlockGetInfo
	response, err := c.client.R().
		SetResult(&block).
		Get("/status?q=getInfo")

	if err != nil {
		return nil, fmt.Errorf("cannot get latest block: %w", err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get latest block fail")
	}
	return &block.BlockInfo, nil
}

func (c *zenTPClient) GetBlockHashByHeight(blockHeight int64) (string, error) {
	var block types.BlockHash
	response, err := c.client.R().
		SetResult(&block).
		Get("/block-index/" + strconv.FormatInt(blockHeight, 10))
	if err != nil {
		return "", fmt.Errorf("cannot get block hash by height: %w", err)
	}
	if response.StatusCode() != 200 {
		return "", errors.New("get block hash by height fail")
	}
	return block.BlockHash, nil
}

func (c *zenTPClient) GetBlockAndHashByOffset(offset int64) (int64, string, error) {
	blockInfo, err := c.GetLatestBlock()
	if err != nil {
		log.Error("get block and hash by offset  fail", "err", err)
		return 0, "", err
	}

	txBlock := blockInfo.Blocks - offset
	txBlockHash, err := c.GetBlockHashByHeight(txBlock)
	if err != nil {
		return 0, "", err
	}

	return txBlock, txBlockHash, nil
}

func (c *zenTPClient) SendTransaction(tx string) (string, error) {
	var result types.SendRawTransactionResponse
	response, err := c.client.R().
		SetResult(&result).
		SetBody(types.SendRawTransaction{
			RawTx: tx,
		}).
		Post("/tx/send")

	if err != nil {
		return "", fmt.Errorf("cannot send raw tx: %w, server response: %v", err, response)
	}
	if response.StatusCode() != 200 {
		return "", errors.New("send raw tx fail")
	}
	return result.TxId, nil
}

func DeserializeTx(hexStr string, withPrevScriptPubKey bool) (*types.RawTransaction, error) {
	buf, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("cant decode hex str: %v", err)
	}

	offset := 0

	tx := &types.RawTransaction{
		Ins:  []types.TxInput{},
		Outs: []types.TxOut{},
	}

	if offset+4 > len(buf) {
		return nil, fmt.Errorf("no enough data to read version")
	}
	tx.Version = binary.LittleEndian.Uint32(buf[offset:])
	offset += 4

	vinLen, bytesRead := util.ReadVarInt(buf, offset)
	if bytesRead == 0 {
		return nil, fmt.Errorf("cant read input length")
	}
	offset += bytesRead

	for i := uint64(0); i < vinLen; i++ {
		if offset+32 > len(buf) {
			return nil, fmt.Errorf("cant read prev tx hash")
		}
		hash := util.ReverseBytes(buf[offset : offset+32])
		offset += 32

		if offset+4 > len(buf) {
			return nil, fmt.Errorf("cant read output index")
		}
		vout := binary.LittleEndian.Uint32(buf[offset:])
		offset += 4

		prevScriptPubKey := ""
		if withPrevScriptPubKey {
			prevScriptPubKeyLen, bytesRead := util.ReadVarInt(buf, offset)
			if bytesRead == 0 {
				return nil, fmt.Errorf("cant read prev script pubkey length")
			}
			offset += bytesRead

			if offset+int(prevScriptPubKeyLen) > len(buf) {
				return nil, fmt.Errorf("cant read prev script pubkey")
			}
			prevScriptPubKey = hex.EncodeToString(buf[offset : offset+int(prevScriptPubKeyLen)])
			offset += int(prevScriptPubKeyLen)
		}

		scriptLen, bytesRead := util.ReadVarInt(buf, offset)
		if bytesRead == 0 {
			return nil, fmt.Errorf("cant read script length")
		}
		offset += bytesRead

		if offset+int(scriptLen) > len(buf) {
			return nil, fmt.Errorf("cant read script")
		}
		script := hex.EncodeToString(buf[offset : offset+int(scriptLen)])
		offset += int(scriptLen)

		if offset+4 > len(buf) {
			return nil, fmt.Errorf("cant read sequence")
		}
		sequence := hex.EncodeToString(buf[offset : offset+4])
		offset += 4

		tx.Ins = append(tx.Ins, types.TxInput{
			Output: types.TxOutput{
				Hash: hex.EncodeToString(hash),
				Vout: vout,
			},
			Script:           script,
			Sequence:         sequence,
			PrevScriptPubKey: prevScriptPubKey,
		})
	}

	voutLen, bytesRead := util.ReadVarInt(buf, offset)
	if bytesRead == 0 {
		return nil, fmt.Errorf("cant read output length")
	}
	offset += bytesRead

	for i := uint64(0); i < voutLen; i++ {
		if offset+8 > len(buf) {
			return nil, fmt.Errorf("cant read satoshis")
		}
		satoshis := binary.LittleEndian.Uint64(buf[offset:])
		offset += 8

		scriptLen, bytesRead := util.ReadVarInt(buf, offset)
		if bytesRead == 0 {
			return nil, fmt.Errorf("cant read script length")
		}
		offset += bytesRead

		if offset+int(scriptLen) > len(buf) {
			return nil, fmt.Errorf("cant read script")
		}
		script := hex.EncodeToString(buf[offset : offset+int(scriptLen)])
		offset += int(scriptLen)

		tx.Outs = append(tx.Outs, types.TxOut{
			Satoshis: satoshis,
			Script:   script,
		})
	}

	if offset+4 > len(buf) {
		return nil, fmt.Errorf("cant read locktime")
	}
	tx.Locktime = binary.LittleEndian.Uint32(buf[offset:])

	return tx, nil
}
