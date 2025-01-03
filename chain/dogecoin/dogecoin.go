package dogecoin

import (
	"bytes"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	"github.com/dapplink-labs/wallet-chain-utxo/chain"
	"github.com/dapplink-labs/wallet-chain-utxo/chain/dogecoin/base"
	"github.com/dapplink-labs/wallet-chain-utxo/config"
	common2 "github.com/dapplink-labs/wallet-chain-utxo/rpc/common"
	"github.com/dapplink-labs/wallet-chain-utxo/rpc/utxo"
)

const ChainName = "Dogecoin"

type ChainAdaptor struct {
	dogeClient *base.DogeClient
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	baseClient, err := base.NewDogeClient(conf.WalletNode.Doge.DataApiUrl, conf.WalletNode.Doge.DataApiKey)
	if err != nil {
		log.Error("new dogecoin data client fail", "err", err)
		return nil, err
	}
	return &ChainAdaptor{
		dogeClient: baseClient,
	}, nil
}
func (c *ChainAdaptor) GetSupportChains(req *utxo.SupportChainsRequest) (*utxo.SupportChainsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetFee(req *utxo.FeeRequest) (*utxo.FeeResponse, error) {
	feeRate, err := c.dogeClient.GetFeeRate()
	if err != nil {
		return &utxo.FeeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	// 将 satoshi/kb 转换为 DOGE/kb
	normalFee := float64(feeRate.MediumFeePerKb) / 1e8
	slowFee := float64(feeRate.LowFeePerKb) / 1e8
	fastFee := float64(feeRate.HighFeePerKb) / 1e8

	return &utxo.FeeResponse{
		Code:       common2.ReturnCode_SUCCESS,
		Msg:        "get fee success",
		BestFee:    strconv.FormatFloat(normalFee, 'f', -1, 64),
		BestFeeSat: strconv.FormatInt(feeRate.MediumFeePerKb, 10),
		SlowFee:    strconv.FormatFloat(slowFee, 'f', -1, 64),
		NormalFee:  strconv.FormatFloat(normalFee, 'f', -1, 64),
		FastFee:    strconv.FormatFloat(fastFee, 'f', -1, 64),
	}, nil

}

func (c *ChainAdaptor) GetAccount(req *utxo.AccountRequest) (*utxo.AccountResponse, error) {
	account, err := c.dogeClient.GetAccountBalance(req.Address)
	if err != nil {
		return &utxo.AccountResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	return &utxo.AccountResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "get account success",
		Balance: strconv.FormatInt(account.Balance, 10),
	}, nil
}

func (c *ChainAdaptor) GetUnspentOutputs(req *utxo.UnspentOutputsRequest) (*utxo.UnspentOutputsResponse, error) {
	utxos, err := c.dogeClient.GetUnspentOutputs(req.Address)
	if err != nil {
		return &utxo.UnspentOutputsResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	var unspentOutputList []*utxo.UnspentOutput
	for _, value := range utxos.TxRefs {
		if value.TxInputN == -1 {
			unspentOutput := &utxo.UnspentOutput{
				TxHashBigEndian: "",
				TxId:            value.TxHash,
				TxOutputN:       uint64(value.TxOutputN),
				Script:          "",
				UnspentAmount:   strconv.FormatInt(value.Value, 10),
				Index:           0,
				Confirmations:   uint64(value.Confirmations),
			}
			unspentOutputList = append(unspentOutputList, unspentOutput)
		}
	}
	return &utxo.UnspentOutputsResponse{
		Code:           common2.ReturnCode_SUCCESS,
		Msg:            "get unspent outputs success",
		UnspentOutputs: unspentOutputList,
	}, nil
}

func (c *ChainAdaptor) GetBlockByNumber(req *utxo.BlockNumberRequest) (*utxo.BlockResponse, error) {
	block, err := c.dogeClient.GetBlockByNumber(req.GetHeight())
	if err != nil {
		return &utxo.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	// 构造交易列表
	var txList []*utxo.TransactionList
	for _, txid := range block.TxIds {
		// 获取每个交易的详细信息
		tx, err := c.dogeClient.GetTransaction(txid)
		if err != nil {
			continue
		}

		// 构造输入
		var vins []*utxo.Vin
		for _, input := range tx.Inputs {
			address := ""
			if len(input.Addresses) > 0 {
				address = input.Addresses[0]
			}
			vins = append(vins, &utxo.Vin{
				Hash:    input.PrevHash,
				Index:   uint32(input.OutputIndex),
				Amount:  input.OutputValue,
				Address: address,
			})
		}

		// 构造输出
		var vouts []*utxo.Vout
		for i, output := range tx.Outputs {
			address := ""
			if len(output.Addresses) > 0 {
				address = output.Addresses[0]
			}
			vouts = append(vouts, &utxo.Vout{
				Address: address,
				Amount:  output.Value,
				Index:   uint32(i),
			})
		}

		txList = append(txList, &utxo.TransactionList{
			Hash: txid,
			Fee:  strconv.FormatInt(tx.Fees, 10),
			Vin:  vins,
			Vout: vouts,
		})
	}

	return &utxo.BlockResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "get block success",
		Height: uint64(block.Height),
		Hash:   block.Hash,
		TxList: txList,
	}, nil
}

func (c *ChainAdaptor) GetBlockByHash(req *utxo.BlockHashRequest) (*utxo.BlockResponse, error) {
	block, err := c.dogeClient.GetBlockByHash(req.Hash)
	if err != nil {
		return &utxo.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	// 构造交易列表
	var txList []*utxo.TransactionList
	for _, txid := range block.TxIds {
		// 获取每个交易的详细信息
		tx, err := c.dogeClient.GetTransaction(txid)
		if err != nil {
			continue
		}

		// 构造输入
		var vins []*utxo.Vin
		for _, input := range tx.Inputs {
			address := ""
			if len(input.Addresses) > 0 {
				address = input.Addresses[0]
			}
			vins = append(vins, &utxo.Vin{
				Hash:    input.PrevHash,
				Index:   uint32(input.OutputIndex),
				Amount:  input.OutputValue,
				Address: address,
			})
		}

		// 构造输出
		var vouts []*utxo.Vout
		for i, output := range tx.Outputs {
			address := ""
			if len(output.Addresses) > 0 {
				address = output.Addresses[0]
			}
			vouts = append(vouts, &utxo.Vout{
				Address: address,
				Amount:  output.Value,
				Index:   uint32(i),
			})
		}

		txList = append(txList, &utxo.TransactionList{
			Hash: txid,
			Fee:  strconv.FormatInt(tx.Fees, 10),
			Vin:  vins,
			Vout: vouts,
		})
	}

	return &utxo.BlockResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "get block success",
		Height: uint64(block.Height),
		Hash:   block.Hash,
		TxList: txList,
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByHash(req *utxo.BlockHeaderHashRequest) (*utxo.BlockHeaderResponse, error) {
	block, err := c.dogeClient.GetBlockByHash(req.Hash)
	if err != nil {
		return &utxo.BlockHeaderResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get block header fail",
		}, err
	}

	return &utxo.BlockHeaderResponse{
		Code:       common2.ReturnCode_SUCCESS,
		Msg:        "get block header success",
		ParentHash: block.PrevBlock,
		Number:     string(block.Ver),
		BlockHash:  req.Hash,
		MerkleRoot: block.MrklRoot,
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByNumber(req *utxo.BlockHeaderNumberRequest) (*utxo.BlockHeaderResponse, error) {
	block, err := c.dogeClient.GetBlockByNumber(req.Height)
	if err != nil {
		return &utxo.BlockHeaderResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get block header fail",
		}, err
	}

	return &utxo.BlockHeaderResponse{
		Code:       common2.ReturnCode_SUCCESS,
		Msg:        "get block header success",
		ParentHash: block.PrevBlock,
		Number:     string(block.Ver),
		BlockHash:  block.Hash,
		MerkleRoot: block.MrklRoot,
	}, nil
}

func (c *ChainAdaptor) GetTxByAddress(req *utxo.TxAddressRequest) (*utxo.TxAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetTxByHash(req *utxo.TxHashRequest) (*utxo.TxHashResponse, error) {
	tx, err := c.dogeClient.GetTransaction(req.Hash)
	if err != nil {
		return &utxo.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	// 构造 froms (输入地址)
	var froms []*utxo.Address
	for _, input := range tx.Inputs {
		if len(input.Addresses) > 0 {
			froms = append(froms, &utxo.Address{
				Address: input.Addresses[0],
			})
		}
	}

	// 构造 tos (输出地址)
	var tos []*utxo.Address
	// 构造 values (对应的金额)
	var values []*utxo.Value
	for _, output := range tx.Outputs {
		if len(output.Addresses) > 0 {
			tos = append(tos, &utxo.Address{
				Address: output.Addresses[0],
			})
			values = append(values, &utxo.Value{
				Value: strconv.FormatInt(output.Value, 10),
			})
		}
	}

	txMsg := &utxo.TxMessage{
		Hash:     tx.Hash,                               // 交易哈希
		Index:    uint32(0),                             // 索引
		Froms:    froms,                                 // 输入地址列表
		Tos:      tos,                                   // 输出地址列表
		Values:   values,                                // 对应的金额列表
		Fee:      strconv.FormatInt(tx.Fees, 10),        // 手续费
		Status:   utxo.TxStatus_Success,                 // 交易状态
		Type:     0,                                     // 交易类型
		Height:   strconv.FormatInt(tx.BlockHeight, 10), // 区块高度
		Datetime: tx.Confirmed.Format(time.RFC3339),     // 确认时间
	}

	return &utxo.TxHashResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get transaction success",
		Tx:   txMsg,
	}, nil
}

// 主要修改的部分是地址格式和网络参数
func (c *ChainAdaptor) ConvertAddress(req *utxo.ConvertAddressRequest) (*utxo.ConvertAddressResponse, error) {
	var address string
	compressedPubKeyBytes, _ := hex.DecodeString(req.PublicKey)
	pubKeyHash := btcutil.Hash160(compressedPubKeyBytes)

	// 使用 chaincfg.MainNetParams，但需要修改狗狗币的参数
	params := &chaincfg.MainNetParams
	params.PubKeyHashAddrID = 0x1e // Dogecoin P2PKH地址前缀
	params.ScriptHashAddrID = 0x16 // Dogecoin P2SH地址前缀

	switch req.Format {
	case "p2pkh":
		// 生成普通地址（D开头）
		p2pkhAddr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, params)
		if err != nil {
			return nil, err
		}
		address = p2pkhAddr.EncodeAddress()
	case "p2sh":
		// 生成多签地址（A开头）
		p2shAddr, err := btcutil.NewAddressScriptHash(pubKeyHash, params)
		if err != nil {
			return nil, err
		}
		address = p2shAddr.EncodeAddress()
	default:
		return nil, errors.New("unsupported address format for Dogecoin")
	}

	return &utxo.ConvertAddressResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "create address success",
		Address: address,
	}, nil
}

// ValidAddress 验证地址
func (c *ChainAdaptor) ValidAddress(req *utxo.ValidAddressRequest) (*utxo.ValidAddressResponse, error) {
	// 使用Dogecoin的网络参数
	params := &chaincfg.MainNetParams
	params.PubKeyHashAddrID = 0x1e
	params.ScriptHashAddrID = 0x16

	address, err := btcutil.DecodeAddress(req.Address, params)
	if err != nil {
		return &utxo.ValidAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}

	if !address.IsForNet(params) {
		return &utxo.ValidAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "address is not valid for Dogecoin network",
		}, nil
	}

	return &utxo.ValidAddressResponse{
		Code:  common2.ReturnCode_SUCCESS,
		Msg:   "verify address success",
		Valid: true,
	}, nil
}

//
//// 其他方法基本保持不变,只需要修改一些具体的参数
//// ... 其他方法的实现
//
//func (c *ChainAdaptor) EstimateSmartFee() (float64, error) {
//	resp, err := c.dogeClient.Call()
//	if err != nil {
//		return 0, fmt.Errorf("failed to send request: %v", err)
//	}
//
//	var result struct {
//		Result struct {
//			FeeRate float64 `json:"feerate"`
//			Blocks  int64   `json:"blocks"`
//		} `json:"result"`
//		Error interface{} `json:"error"`
//		ID    string      `json:"id"`
//	}
//
//	if err := json.Unmarshal(resp, &result); err != nil {
//		return 0, fmt.Errorf("failed to unmarshal response: %v", err)
//	}
//
//	if result.Error != nil {
//		return 0, fmt.Errorf("RPC error: %v", result.Error)
//	}
//
//	if result.Result.FeeRate < 0 {
//		return MinFeeRate, nil
//	}
//
//	return result.Result.FeeRate, nil
//}

func (c *ChainAdaptor) CreateUnSignTransaction(req *utxo.UnSignTransactionRequest) (*utxo.UnSignTransactionResponse, error) {
	// 1. 参数验证
	if len(req.Vin) == 0 || len(req.Vout) == 0 {
		return &utxo.UnSignTransactionResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "invalid transaction inputs or outputs",
		}, errors.New("invalid transaction inputs or outputs")
	}

	// 2. 使用狗狗币的网络参数
	params := &chaincfg.MainNetParams
	params.PubKeyHashAddrID = 0x1e // Dogecoin P2PKH地址前缀
	params.ScriptHashAddrID = 0x16 // Dogecoin P2SH地址前缀

	// 3. 创建交易
	rawTx := wire.NewMsgTx(wire.TxVersion)

	// 4. 添加输入
	for _, in := range req.Vin {
		utxoHash, err := chainhash.NewHashFromStr(in.Hash)
		if err != nil {
			return &utxo.UnSignTransactionResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  "invalid input hash: " + err.Error(),
			}, err
		}
		txIn := wire.NewTxIn(wire.NewOutPoint(utxoHash, in.Index), nil, nil)
		rawTx.AddTxIn(txIn)
	}

	// 5. 添加输出
	for _, out := range req.Vout {
		toAddress, err := btcutil.DecodeAddress(out.Address, params)
		if err != nil {
			return &utxo.UnSignTransactionResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  "invalid output address: " + err.Error(),
			}, err
		}
		toPkScript, err := txscript.PayToAddrScript(toAddress)
		if err != nil {
			return &utxo.UnSignTransactionResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  "create output script failed: " + err.Error(),
			}, err
		}
		rawTx.AddTxOut(wire.NewTxOut(out.Amount, toPkScript))
	}

	// 6. 计算签名哈希
	signHashes := make([][]byte, len(req.Vin))
	for i, in := range req.Vin {
		fromAddr, err := btcutil.DecodeAddress(in.Address, params)
		if err != nil {
			return &utxo.UnSignTransactionResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  "invalid input address: " + err.Error(),
			}, err
		}
		fromPkScript, err := txscript.PayToAddrScript(fromAddr)
		if err != nil {
			return &utxo.UnSignTransactionResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  "create input script failed: " + err.Error(),
			}, err
		}
		signHash, err := txscript.CalcSignatureHash(fromPkScript, txscript.SigHashAll, rawTx, i)
		if err != nil {
			return &utxo.UnSignTransactionResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  "calculate sign hash failed: " + err.Error(),
			}, err
		}
		signHashes[i] = signHash
	}

	// 7. 序列化交易
	var txBuf bytes.Buffer
	if err := rawTx.Serialize(&txBuf); err != nil {
		return &utxo.UnSignTransactionResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "serialize transaction failed: " + err.Error(),
		}, err
	}

	// 将二进制数据转换为十六进制字符串
	txHex := hex.EncodeToString(txBuf.Bytes())

	// 签名哈希转换为十六进制字符串数组
	signHashesHex := make([][]byte, 0)
	for _, hash := range signHashes {
		// 将每个哈希转换为十六进制字符串
		hashHex := hex.EncodeToString(hash)
		signHashesHex = append(signHashesHex, []byte(hashHex))
	}

	log.Info("Transaction created",
		"txHex", txHex,
		"signHashes", signHashesHex)

	// 8. 返回结果
	return &utxo.UnSignTransactionResponse{
		Code:       common2.ReturnCode_SUCCESS,
		Msg:        "create unsigned transaction success",
		TxData:     []byte(txHex),
		SignHashes: signHashesHex,
	}, nil
}

// CalcSignHashes 计算签名哈希
func (c *ChainAdaptor) CalcSignHashes(Vins []*utxo.Vin, Vouts []*utxo.Vout) ([][]byte, []byte, error) {
	if len(Vins) == 0 || len(Vouts) == 0 {
		return nil, nil, errors.New("invalid len in or out")
	}

	// 使用狗狗币的网络参数
	params := &chaincfg.MainNetParams
	params.PubKeyHashAddrID = 0x1e // Dogecoin P2PKH地址前缀
	params.ScriptHashAddrID = 0x16 // Dogecoin P2SH地址前缀

	rawTx := wire.NewMsgTx(wire.TxVersion)

	// 添加输入
	for _, in := range Vins {
		utxoHash, err := chainhash.NewHashFromStr(in.Hash)
		if err != nil {
			return nil, nil, err
		}
		txIn := wire.NewTxIn(wire.NewOutPoint(utxoHash, in.Index), nil, nil)
		rawTx.AddTxIn(txIn)
	}

	// 添加输出
	for _, out := range Vouts {
		toAddress, err := btcutil.DecodeAddress(out.Address, params)
		if err != nil {
			return nil, nil, err
		}
		toPkScript, err := txscript.PayToAddrScript(toAddress)
		if err != nil {
			return nil, nil, err
		}
		rawTx.AddTxOut(wire.NewTxOut(out.Amount, toPkScript))
	}

	// 计算每个输入的签名哈希
	signHashes := make([][]byte, len(Vins))
	for i, in := range Vins {
		fromAddr, err := btcutil.DecodeAddress(in.Address, params)
		if err != nil {
			return nil, nil, err
		}
		fromPkScript, err := txscript.PayToAddrScript(fromAddr)
		if err != nil {
			return nil, nil, err
		}
		signHash, err := txscript.CalcSignatureHash(fromPkScript, txscript.SigHashAll, rawTx, i)
		if err != nil {
			return nil, nil, err
		}
		signHashes[i] = signHash
	}

	// 序列化交易
	buf := bytes.NewBuffer(make([]byte, 0, rawTx.SerializeSize()))
	err := rawTx.Serialize(buf)
	if err != nil {
		return nil, nil, err
	}

	return signHashes, buf.Bytes(), nil
}

func (c *ChainAdaptor) BuildSignedTransaction(req *utxo.SignedTransactionRequest) (*utxo.SignedTransactionResponse, error) {
	// 1. 将十六进制字符串转换回字节数组
	txHex := string(req.TxData)
	log.Info("Transaction hex", "hex", txHex)

	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		log.Error("Failed to decode transaction hex",
			"error", err,
			"txHex", txHex)
		return nil, err
	}

	// 2. 反序列化交易
	var rawTx wire.MsgTx
	if err := rawTx.Deserialize(bytes.NewReader(txBytes)); err != nil {
		log.Error("Failed to deserialize transaction",
			"error", err,
			"txBytes", hex.EncodeToString(txBytes))
		return nil, err
	}

	// 3. 检查签名和公钥
	if len(req.Signatures) != len(req.PublicKeys) || len(req.Signatures) == 0 {
		log.Error("Invalid signatures or public keys",
			"sigCount", len(req.Signatures),
			"pubKeyCount", len(req.PublicKeys))
		return nil, errors.New("invalid signatures or public keys")
	}

	// 4. 构建签名脚本
	for i := 0; i < len(rawTx.TxIn); i++ {
		builder := txscript.NewScriptBuilder()
		builder.AddData(req.Signatures[0])
		builder.AddData(req.PublicKeys[0])
		signScript, err := builder.Script()
		if err != nil {
			log.Error("Failed to build signature script", "error", err)
			return nil, err
		}
		rawTx.TxIn[i].SignatureScript = signScript
	}

	// 5. 序列化签名后的交易
	var signedTxBuf bytes.Buffer
	if err := rawTx.Serialize(&signedTxBuf); err != nil {
		log.Error("Failed to serialize signed transaction", "error", err)
		return nil, err
	}

	// 6. 计算交易哈希
	txHash := rawTx.TxHash()

	return &utxo.SignedTransactionResponse{
		Code:         common2.ReturnCode_SUCCESS,
		Msg:          "Transaction signed successfully",
		SignedTxData: signedTxBuf.Bytes(),
		Hash:         txHash[:],
	}, nil
}

func (c *ChainAdaptor) DecodeTransaction(req *utxo.DecodeTransactionRequest) (*utxo.DecodeTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) VerifySignedTransaction(req *utxo.VerifyTransactionRequest) (*utxo.VerifyTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}
func (c *ChainAdaptor) SendTx(req *utxo.SendTxRequest) (*utxo.SendTxResponse, error) {
	//TODO implement me
	panic("implement me")
}
