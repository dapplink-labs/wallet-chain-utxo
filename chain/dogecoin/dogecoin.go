package dogecoin

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
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
	// 1. 记录输入参数
	log.Info("Building signed transaction",
		"txDataLen", len(req.TxData),
		"signaturesLen", len(req.Signatures),
		"publicKeysLen", len(req.PublicKeys))

	// 2. 创建一个新的交易对象
	rawTx := wire.NewMsgTx(wire.TxVersion)

	// 3. 手动解析交易数据
	data := req.TxData
	if len(data) < 5 {
		return nil, errors.New("transaction data too short")
	}

	// 4. 解析版本号 (4 bytes)
	rawTx.Version = int32(binary.LittleEndian.Uint32(data[0:4]))
	pos := 4

	// 5. 解析输入数量 (1 byte)
	numInputs := int(data[pos])
	pos++

	log.Info("Parsing transaction header",
		"version", rawTx.Version,
		"numInputs", numInputs,
		"pos", pos,
		"data", fmt.Sprintf("%x", data[pos:pos+32]))

	// 6. 解析输入
	for i := 0; i < numInputs; i++ {
		if pos+36 > len(data) {
			return nil, fmt.Errorf("invalid input data at position %d", pos)
		}

		// 解析前一个输出哈希 (32 bytes)
		var hash chainhash.Hash
		copy(hash[:], data[pos:pos+32])
		pos += 32

		// 解析前一个输出索引 (4 bytes)
		index := binary.LittleEndian.Uint32(data[pos : pos+4])
		pos += 4

		// 创建输入
		txIn := wire.NewTxIn(&wire.OutPoint{Hash: hash, Index: index}, nil, nil)

		// 解析脚本长度
		if pos >= len(data) {
			return nil, fmt.Errorf("invalid script length position at %d", pos)
		}
		scriptLen := int(data[pos])
		pos++

		log.Info("Parsing input",
			"index", i,
			"hash", hash.String(),
			"outIndex", index,
			"scriptLen", scriptLen,
			"pos", pos)

		// 跳过脚本数据
		if pos+scriptLen > len(data) {
			return nil, fmt.Errorf("invalid script data at position %d, length %d", pos, scriptLen)
		}
		pos += scriptLen

		// 解析 sequence (4 bytes)
		if pos+4 > len(data) {
			return nil, fmt.Errorf("invalid sequence data at position %d", pos)
		}
		txIn.Sequence = binary.LittleEndian.Uint32(data[pos : pos+4])
		pos += 4

		log.Info("Input parsed",
			"index", i,
			"sequence", txIn.Sequence,
			"pos", pos,
			"nextBytes", fmt.Sprintf("%x", data[pos:pos+8]))

		rawTx.AddTxIn(txIn)
	}

	// 7. 解析输出数量 (1 byte)
	if pos >= len(data) {
		return nil, fmt.Errorf("invalid output count position at %d", pos)
	}

	// 检查是否需要跳过 sequence 的剩余部分
	if data[pos] == 0xff && pos+1 < len(data) && data[pos+1] == 0xff {
		pos += 2 // 跳过 sequence 的剩余部分 (ffff)
	}

	// 现在读取真正的输出数量
	numOutputs := int(data[pos])
	pos++

	log.Info("Parsing outputs",
		"numOutputs", numOutputs,
		"pos", pos,
		"nextBytes", fmt.Sprintf("%x", data[pos:pos+16]))

	// 8. 解析输出
	for i := 0; i < numOutputs; i++ {
		if pos+8 > len(data) {
			return nil, fmt.Errorf("invalid output value at position %d", pos)
		}

		// 解析金额 (8 bytes)
		value := int64(binary.LittleEndian.Uint64(data[pos : pos+8]))
		pos += 8

		// 解析公钥脚本长度 (1 byte)
		if pos >= len(data) {
			return nil, fmt.Errorf("invalid script length position at %d", pos)
		}

		// 读取脚本长度并进行验证
		scriptLen := int(data[pos])
		pos++

		// 调试信息
		log.Info("Script length analysis",
			"index", i,
			"rawByte", fmt.Sprintf("%x", data[pos-1]),
			"scriptLen", scriptLen,
			"pos", pos,
			"nextBytes", fmt.Sprintf("%x", data[pos:min(pos+32, len(data))]))

		// 特殊处理：如果脚本长度看起来不合理，尝试使用标准长度
		if scriptLen > 0x9a {
			scriptLen = 25 // 标准 P2PKH 脚本长度
		}

		// 验证脚本长度的合理性
		if scriptLen > 100 {
			return nil, fmt.Errorf("unreasonable script length: %d at position %d", scriptLen, pos-1)
		}

		// 确保有足够的数据
		if pos+scriptLen > len(data) {
			return nil, fmt.Errorf("invalid script data at position %d, length %d, remaining %d",
				pos, scriptLen, len(data)-pos)
		}

		// 复制脚本数据
		pkScript := make([]byte, scriptLen)
		copy(pkScript, data[pos:pos+scriptLen])
		pos += scriptLen

		// 创建输出
		txOut := wire.NewTxOut(value, pkScript)
		rawTx.AddTxOut(txOut)

		log.Info("Output added",
			"index", i,
			"value", value,
			"scriptLen", scriptLen,
			"script", fmt.Sprintf("%x", pkScript))
	}

	// 9. 添加签名
	if len(req.Signatures) == 0 || len(req.PublicKeys) == 0 {
		return nil, errors.New("missing signatures or public keys")
	}

	// 10. 构建签名脚本
	builder := txscript.NewScriptBuilder()
	builder.AddData(req.Signatures[0])
	builder.AddData(req.PublicKeys[0])

	signScript, err := builder.Script()
	if err != nil {
		return nil, fmt.Errorf("failed to build signature script: %v", err)
	}

	// 11. 添加签名脚本到输入
	for i := range rawTx.TxIn {
		rawTx.TxIn[i].SignatureScript = signScript
	}

	// 12. 序列化签名后的交易
	var signedTxBuf bytes.Buffer
	if err := rawTx.Serialize(&signedTxBuf); err != nil {
		return nil, fmt.Errorf("failed to serialize transaction: %v", err)
	}

	// 13. 计算交易哈希
	txHash := rawTx.TxHash()

	log.Info("Transaction successfully signed",
		"hash", txHash.String(),
		"size", signedTxBuf.Len(),
		"numInputs", len(rawTx.TxIn),
		"numOutputs", len(rawTx.TxOut))

	return &utxo.SignedTransactionResponse{
		Code:         common2.ReturnCode_SUCCESS,
		Msg:          "Transaction signed successfully",
		SignedTxData: signedTxBuf.Bytes(),
		Hash:         txHash[:],
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
