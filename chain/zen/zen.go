package zen

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"

	"github.com/dapplink-labs/wallet-chain-utxo/chain"
	"github.com/dapplink-labs/wallet-chain-utxo/chain/zen/types"
	"github.com/dapplink-labs/wallet-chain-utxo/chain/zen/util"
	"github.com/dapplink-labs/wallet-chain-utxo/config"
	common2 "github.com/dapplink-labs/wallet-chain-utxo/rpc/common"
	"github.com/dapplink-labs/wallet-chain-utxo/rpc/utxo"
)

const ChainName = "Zen"

type ChainAdaptor struct {
	thirdPartClient *zenTPClient
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	zenTPClient, err := NewBlockChainClient(conf.WalletNode.Zen.TpApiUrl)
	if err != nil {
		log.Error("new blockchain client fail", "err", err)
		return nil, err
	}
	return &ChainAdaptor{
		// zenDataClient:   baseDataClient,
		thirdPartClient: zenTPClient,
	}, nil
}

func (c *ChainAdaptor) GetAccount(req *utxo.AccountRequest) (*utxo.AccountResponse, error) {
	balance, err := c.thirdPartClient.GetAccountBalance(req.Address)
	if err != nil {
		return &utxo.AccountResponse{
			Code:    common2.ReturnCode_ERROR,
			Msg:     "get Zen balance fail",
			Balance: "0",
			Network: "",
		}, err
	}
	return &utxo.AccountResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "get Zen balance success",
		Balance: balance,
		Network: req.Network,
	}, nil
}

/*
*
TODO
horizen api provide a estimate fee api with {nil} as response, need to implement gas fee
estimate by calculate the past serval block average fee.
*/
func (c *ChainAdaptor) GetFee(req *utxo.FeeRequest) (*utxo.FeeResponse, error) {
	gasFeeResp, err := c.thirdPartClient.GetFee()
	if err != nil {
		return &utxo.FeeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	return &utxo.FeeResponse{
		Code:       common2.ReturnCode_SUCCESS,
		Msg:        "get fee success",
		BestFee:    gasFeeResp.BestTransactionFee,
		BestFeeSat: gasFeeResp.BestTransactionFeeSat,
		SlowFee:    gasFeeResp.SlowGasPrice,
		NormalFee:  gasFeeResp.StandardGasPrice,
		FastFee:    gasFeeResp.RapidGasPrice,
	}, nil
}

func (c *ChainAdaptor) GetTxByAddress(req *utxo.TxAddressRequest) (*utxo.TxAddressResponse, error) {
	transaction, err := c.thirdPartClient.GetTransactionsByAddress(req.Address, req.Pagesize, req.Page)
	if err != nil {
		return &utxo.TxAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get transaction list fail",
			Tx:   nil,
		}, err

	}

	if err != nil {
		return &utxo.TxAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get transaction list fail",
			Tx:   nil,
		}, err
	}
	var tx_list []*utxo.TxMessage
	for _, ttxs := range transaction.Items {
		var from_addrs []*utxo.Address
		var to_addrs []*utxo.Address
		var value_list []*utxo.Value
		var direction int32
		for _, inputs := range ttxs.Vin {
			from_addrs = append(from_addrs, &utxo.Address{Address: inputs.Address})
		}
		tx_fee := ttxs.Fees
		for _, out := range ttxs.Vout {
			for _, addr := range out.ScriptPubKey.Address {
				to_addrs = append(to_addrs, &utxo.Address{Address: addr})
			}
			value_list = append(value_list, &utxo.Value{Value: out.Value})
		}
		datetime := strconv.FormatUint(ttxs.Time, 10)
		if strings.EqualFold(req.Address, from_addrs[0].Address) {
			direction = 0
		} else {
			direction = 1
		}
		tx := &utxo.TxMessage{
			Hash:     ttxs.TXID,
			Froms:    from_addrs,
			Tos:      to_addrs,
			Values:   value_list,
			Fee:      strconv.FormatFloat(tx_fee, 'f', -1, 64),
			Status:   utxo.TxStatus_Success,
			Type:     direction,
			Height:   strconv.FormatUint(ttxs.BlockHeight, 10),
			Datetime: datetime,
		}
		tx_list = append(tx_list, tx)
	}
	return &utxo.TxAddressResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get transaction list success",
		Tx:   tx_list,
	}, nil
}

func (c *ChainAdaptor) GetTxByHash(req *utxo.TxHashRequest) (*utxo.TxHashResponse, error) {
	transaction, err := c.thirdPartClient.GetTransactionsByHash(req.Hash)
	if err != nil {
		return &utxo.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get transaction list fail",
			Tx:   nil,
		}, err
	}
	var from_addrs []*utxo.Address
	var to_addrs []*utxo.Address
	var value_list []*utxo.Value
	for _, inputs := range transaction.Vin {
		from_addrs = append(from_addrs, &utxo.Address{Address: inputs.Address})
	}
	tx_fee := transaction.Fees
	for _, out := range transaction.Vout {
		for _, addr := range out.ScriptPubKey.Address {
			to_addrs = append(to_addrs, &utxo.Address{Address: addr})
		}
		value_list = append(value_list, &utxo.Value{Value: out.Value})
	}
	datetime := strconv.FormatUint(transaction.Time, 10)

	tx := &utxo.TxMessage{
		Hash:     transaction.TXID,
		Froms:    from_addrs,
		Tos:      to_addrs,
		Values:   value_list,
		Fee:      strconv.FormatFloat(tx_fee, 'f', -1, 64),
		Status:   utxo.TxStatus_Success,
		Type:     0,
		Height:   strconv.FormatUint(transaction.BlockHeight, 10),
		Datetime: datetime,
	}
	return &utxo.TxHashResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get transaction success",
		Tx:   tx,
		// txMsg,
	}, nil
}

func (c *ChainAdaptor) ConvertAddress(req *utxo.ConvertAddressRequest) (*utxo.ConvertAddressResponse, error) {
	address, err := util.PublicKeyToAddress(req.PublicKey)
	if err != nil {
		return &utxo.ConvertAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	return &utxo.ConvertAddressResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "create address success",
		Address: address,
	}, nil
}

func (c *ChainAdaptor) ValidAddress(req *utxo.ValidAddressRequest) (*utxo.ValidAddressResponse, error) {
	addressValid, msg, err := util.ValidAddress(req.Address)
	if err != nil {
		return &utxo.ValidAddressResponse{
			Code:  common2.ReturnCode_ERROR,
			Msg:   err.Error(),
			Valid: false,
		}, err
	}

	return &utxo.ValidAddressResponse{
		Code:  common2.ReturnCode_SUCCESS,
		Msg:   msg,
		Valid: addressValid,
	}, nil
}

func (c *ChainAdaptor) GetUnspentOutputs(req *utxo.UnspentOutputsRequest) (*utxo.UnspentOutputsResponse, error) {
	utxoList, err := c.thirdPartClient.GetAccountUtxo(req.Address)
	if err != nil {
		return &utxo.UnspentOutputsResponse{
			Code:           common2.ReturnCode_ERROR,
			Msg:            err.Error(),
			UnspentOutputs: nil,
		}, err
	}
	var unspentOutputList []*utxo.UnspentOutput
	for _, value := range utxoList {
		unspentOutput := &utxo.UnspentOutput{
			Address:       value.Address,
			TxId:          value.Txid,
			Script:        value.ScriptPubKey,
			UnspentAmount: strconv.FormatUint(value.Satoshis, 10),
			Height:        strconv.FormatUint(value.Height, 10),
			Confirmations: value.Confirmations,
		}
		unspentOutputList = append(unspentOutputList, unspentOutput)
	}
	return &utxo.UnspentOutputsResponse{
		Code:           common2.ReturnCode_SUCCESS,
		Msg:            "get unspent outputs success",
		UnspentOutputs: unspentOutputList,
	}, nil
}

func (c *ChainAdaptor) SendTx(req *utxo.SendTxRequest) (*utxo.SendTxResponse, error) {
	_, err := DeserializeTx(req.RawTx, false)
	if err != nil {
		return &utxo.SendTxResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	response, err := c.thirdPartClient.SendTransaction(req.RawTx)
	if err != nil {
		return &utxo.SendTxResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	return &utxo.SendTxResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "send tx success ",
		TxHash: response,
	}, nil
}

func (c *ChainAdaptor) CreateUnSignTransaction(req *utxo.UnSignTransactionRequest) (*utxo.UnSignTransactionResponse, error) {
	txHash, buf, err := c.CalcSignHashes(req.Vin, req.Vout)
	if err != nil {
		log.Error("calc sign hashes fail", "err", err)
		return nil, err
	}
	return &utxo.UnSignTransactionResponse{
		Code:       common2.ReturnCode_SUCCESS,
		Msg:        "create un sign transaction success",
		TxData:     buf,
		SignHashes: txHash,
	}, nil
}

func (c *ChainAdaptor) CalcSignHashes(Vins []*utxo.Vin, Vouts []*utxo.Vout) (signHash [][]byte, rawTx []byte, err error) {
	// TODO use go routing to get vin's  scriptPubKey
	// TODO need optimize the compute, use string less than byte
	block, err := c.thirdPartClient.GetBlockByHeight(0)
	if err != nil {
		return nil, nil, errors.Errorf("get block by height fail, height: %d, err: %v", 0, err)
	}
	txObjectInputParam := types.TXObjectInputParam{
		BlockHeight: block.Height,
		BlockHash:   block.Hash,
		ParamIns:    []types.ParamIn{},
		ParamOuts:   []types.ParamOut{},
	}

	var scriptPubKey string
	for _, vin := range Vins {
		tx, err := c.thirdPartClient.GetTransactionsByHash(vin.Hash)
		if err != nil {
			return nil, nil, errors.Errorf("get transaction by hash fail, hash: %s, err: %v", vin.Hash, err)
		}

		for _, vout := range tx.Vout {
			for _, address := range vout.ScriptPubKey.Address {
				if address == vin.Address {
					scriptPubKey = vout.ScriptPubKey.Hex
					vinAmount := util.ConvertToSatoshi(vout.Value)
					// TODO generate another txoutput to change for address who signed tx
					if vinAmount != vin.Amount {
						return nil, nil, errors.Errorf("vin amount that utxo have not equal to vin input, utxo vin amount: %d, vin input amount: %d", vinAmount, vin.Amount)
					}
					goto getScript
				}
			}
		}
		return nil, nil, errors.Errorf("cant find address in given tx, txhash: %s, address: %s, err: %v", vin.Hash, vin.Address, err)

	getScript:
		tmt := types.ParamIn{
			TXID:         vin.Hash,
			Vout:         uint64(vin.Index),
			ScriptPubKey: scriptPubKey,
		}
		txObjectInputParam.ParamIns = append(txObjectInputParam.ParamIns, tmt)
	}

	for _, vout := range Vouts {
		txObjectInputParam.ParamOuts = append(txObjectInputParam.ParamOuts, types.ParamOut{
			Address:  vout.Address,
			Satoshis: uint64(vout.Amount),
		})
	}

	txObject, err := util.CreateRawTX(&txObjectInputParam)
	if err != nil {
		return nil, nil, errors.Errorf("failed to create raw transaction: %v", err)
	}

	if err != nil {
		return nil, nil, errors.Errorf("failed to create raw transaction: %v", err)
	}
	for index, ins := range txObject.Ins {
		script := ins.PrevScriptPubKey
		newTx, err := util.SignatureForm(txObject, index, script, util.SIGHASH_ALL)
		if err != nil {
			return nil, nil, errors.Errorf("sign tx fail, err: %v", err)
		}

		scriptSig, err := util.GetScriptSignature(newTx, util.SIGHASH_ALL)
		if err != nil {
			return nil, nil, errors.Errorf("get script signature fail, err: %v", err)
		}
		signHash = append(signHash, scriptSig)
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(txObject)
	if err != nil {
		return nil, nil, errors.Errorf("gob encode fail, err: %v", err)
	}
	return signHash, buf.Bytes(), nil
}

func (c *ChainAdaptor) BuildSignedTransaction(req *utxo.SignedTransactionRequest) (*utxo.SignedTransactionResponse, error) {
	if len(req.Signatures) != len(req.PublicKeys) {
		return &utxo.SignedTransactionResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "signatures length not equal to public keys length",
		}, nil
	}
	var txObj types.TXObject
	dec := gob.NewDecoder(bytes.NewReader(req.TxData))
	err := dec.Decode(&txObj)
	if err != nil {
		return &utxo.SignedTransactionResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	if len(req.Signatures) != len(txObj.Ins) {
		return &utxo.SignedTransactionResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "signatures length not equal to decoded txObj ins length",
		}, nil
	}

	for i, _ := range txObj.Ins {
		scriptSig := hex.EncodeToString(req.Signatures[i])
		slen, err := util.GetPushDataLength(scriptSig)
		if err != nil {
			return &utxo.SignedTransactionResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  "get push data length fail",
			}, err
		}
		pubKey := hex.EncodeToString(req.PublicKeys[i])
		pklen, err := util.GetPushDataLength(pubKey)
		if err != nil {
			return &utxo.SignedTransactionResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  "get push data length fail",
			}, err
		}
		txObj.Ins[i].Script = slen + scriptSig + pklen + pubKey
	}

	rawTx, err := util.SerializeTX(&txObj, false)

	return &utxo.SignedTransactionResponse{
		Code:         common2.ReturnCode_SUCCESS,
		SignedTxData: rawTx,
		Hash:         util.DoubleSha256(rawTx),
	}, nil
}

func (c *ChainAdaptor) GetBlockByNumber(req *utxo.BlockNumberRequest) (*utxo.BlockResponse, error) {
	if req.Height < 0 {
		return &utxo.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "block height must be greater than 0",
		}, nil
	}

	block, err := c.thirdPartClient.GetBlockByHeight(req.Height)
	if err != nil {
		return &utxo.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	var txList []*utxo.TransactionList
	for _, txHash := range block.Tx {
		tx, err := c.thirdPartClient.GetTransactionsByHash(txHash)
		if err != nil {
			return &utxo.BlockResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  "txHash:" + txHash + err.Error(),
			}, err
		}
		var vinList []*utxo.Vin
		var voutList []*utxo.Vout
		for _, vin := range tx.Vin {
			if vin.Coinbase != "" {
				vinList = append(vinList, &utxo.Vin{
					//TODO cant modify .proto file to add some fields，use hash to store coinbase
					Hash:    "coinbase:" + vin.Coinbase,
					Address: vin.Address,
				})
			} else {
				vinList = append(vinList, &utxo.Vin{
					Hash:    vin.TxId,
					Address: vin.Address,
				})
			}
		}
		for _, vout := range tx.Vout {
			voutList = append(voutList, &utxo.Vout{
				Amount:  util.ConvertToSatoshi(vout.Value),
				Address: vout.ScriptPubKey.Address[0],
			})
		}

		txList = append(txList, &utxo.TransactionList{
			Hash: txHash,
			Fee:  strconv.FormatFloat(tx.Fees, 'f', -1, 64),
			Vin:  vinList,
			Vout: voutList,
		})
	}
	return &utxo.BlockResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "get block by number succcess",
		Height: uint64(block.Height),
		Hash:   block.Hash,
		TxList: txList,
	}, nil
}

func (c *ChainAdaptor) GetBlockByHash(req *utxo.BlockHashRequest) (*utxo.BlockResponse, error) {
	block, err := c.thirdPartClient.GetBlockByHash(req.Hash)
	if err != nil {
		return &utxo.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	var txList []*utxo.TransactionList
	for _, txHash := range block.Tx {
		txList = append(txList, &utxo.TransactionList{
			Hash: txHash,
		})
	}
	return &utxo.BlockResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "get block by hash succcess",
		Height: uint64(block.Height),
		Hash:   block.Hash,
		TxList: txList,
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByHash(req *utxo.BlockHeaderHashRequest) (*utxo.BlockHeaderResponse, error) {
	return &utxo.BlockHeaderResponse{
		Code: common2.ReturnCode_ERROR,
		Msg:  "method not support",
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByNumber(req *utxo.BlockHeaderNumberRequest) (*utxo.BlockHeaderResponse, error) {
	return &utxo.BlockHeaderResponse{
		Code: common2.ReturnCode_ERROR,
		Msg:  "method not support",
	}, nil
}

func (c *ChainAdaptor) DecodeTransaction(req *utxo.DecodeTransactionRequest) (*utxo.DecodeTransactionResponse, error) {
	return &utxo.DecodeTransactionResponse{
		Code: common2.ReturnCode_ERROR,
		Msg:  "method not support",
	}, nil
}

func (c *ChainAdaptor) VerifySignedTransaction(req *utxo.VerifyTransactionRequest) (*utxo.VerifyTransactionResponse, error) {
	return &utxo.VerifyTransactionResponse{
		Code: common2.ReturnCode_ERROR,
		Msg:  "method not support",
	}, nil
}

func (c *ChainAdaptor) GetSupportChains(req *utxo.SupportChainsRequest) (*utxo.SupportChainsResponse, error) {
	return &utxo.SupportChainsResponse{
		Code:    common2.ReturnCode_ERROR,
		Msg:     "method not support",
		Support: false,
	}, nil
}
