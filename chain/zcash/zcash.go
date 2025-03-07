package zcash

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/dapplink-labs/wallet-chain-utxo/chain/zcash/common"
	common2 "github.com/dapplink-labs/wallet-chain-utxo/rpc/common"
	"strconv"
	"strings"

	"github.com/dapplink-labs/wallet-chain-utxo/chain"
	"github.com/dapplink-labs/wallet-chain-utxo/config"
	"github.com/dapplink-labs/wallet-chain-utxo/rpc/utxo"
	"github.com/ethereum/go-ethereum/log"
)

const ChainName = "Zcash"

type ChainAdaptor struct {
	zcashClient           *BaseClient
	zcashDataClientClient *BaseDataClient
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	baseClient, err := NewBaseClient(conf.WalletNode.Zec.RpcUrl, "", "")
	if err != nil {
		log.Error("new zcash rpc client fail", "err", err)
		return nil, err
	}
	baseDataClient, err := NewBaseDataClient(conf.WalletNode.Zec.DataApiUrl, conf.WalletNode.Zec.DataApiKey, "", "")
	if err != nil {
		log.Error("new zcash data client fail", "err", err)
		return nil, err
	}
	return &ChainAdaptor{
		zcashClient:           baseClient,
		zcashDataClientClient: baseDataClient,
	}, nil
}

func (c *ChainAdaptor) GetSupportChains(req *utxo.SupportChainsRequest) (*utxo.SupportChainsResponse, error) {
	return &utxo.SupportChainsResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

func (c *ChainAdaptor) ConvertAddress(req *utxo.ConvertAddressRequest) (*utxo.ConvertAddressResponse, error) {
	var address string
	compressedPubKeyBytes, _ := hex.DecodeString(req.PublicKey)
	pubKeyHash := btcutil.Hash160(compressedPubKeyBytes)
	switch req.Format {
	case "p2pkh":
		// the version of zcash is 7352(0x1cb8)
		versionedPayload := append([]byte{0x1c, 0xb8}, pubKeyHash...)
		hash1 := sha256.Sum256(versionedPayload)
		hash2 := sha256.Sum256(hash1[:])
		checksum := hash2[:4]
		payloadWithChecksum := append(versionedPayload, checksum...)
		address = base58.Encode(payloadWithChecksum)
		break
	default:
		return nil, errors.New("Do not support address type")
	}
	return &utxo.ConvertAddressResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "create address success",
		Address: address,
	}, nil
}

func (c *ChainAdaptor) ValidAddress(req *utxo.ValidAddressRequest) (*utxo.ValidAddressResponse, error) {
	address := req.GetAddress()
	isValidateAddress := strings.HasPrefix(address, "t1")
	if !isValidateAddress {
		return &utxo.ValidAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "address is not valid",
		}, nil
	} else {
		return &utxo.ValidAddressResponse{
			Code:  common2.ReturnCode_SUCCESS,
			Msg:   "verify address success",
			Valid: true,
		}, nil
	}
}

// get fee
func (c *ChainAdaptor) GetFee(req *utxo.FeeRequest) (*utxo.FeeResponse, error) {
	client := c.zcashDataClientClient
	resp, err := client.R().Get(client.RequestUrl + "blockchain-fees/utxo/zcash/mainnet/mempool")
	if err != nil {
		return &utxo.FeeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	log.Debug("get mempool fees", "get", resp.String())
	resMap := make(map[string]interface{}, 3)
	err = json.Unmarshal([]byte(resp.String()), &resMap)
	itemMap := resMap["data"].(map[string]interface{})
	responseMap := itemMap["item"].(map[string]interface{})
	if err != nil {
		log.Error("Response is not a map[string]interface{}")
		return &utxo.FeeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	return &utxo.FeeResponse{
		Code:       common2.ReturnCode_SUCCESS,
		Msg:        "get fee success",
		BestFee:    "",
		BestFeeSat: "",
		SlowFee:    responseMap["slow"].(string),
		NormalFee:  responseMap["standard"].(string),
		FastFee:    responseMap["fast"].(string),
	}, nil
}

func (c *ChainAdaptor) GetAccount(req *utxo.AccountRequest) (*utxo.AccountResponse, error) {
	address := req.GetAddress()
	resp, err := c.zcashDataClientClient.R().Get(c.zcashDataClientClient.RequestUrl + "addresses-latest/utxo/zcash/mainnet/" + address + "/balance")

	if err != nil {
		return &utxo.AccountResponse{
			Code:    common2.ReturnCode_ERROR,
			Msg:     "get Zec balance fail",
			Balance: "0",
		}, err
	}
	resMap := make(map[string]interface{}, 3)
	err = json.Unmarshal([]byte(resp.String()), &resMap)
	itemMap := resMap["data"].(map[string]interface{})
	responseMap := itemMap["item"].(map[string]interface{})
	comfirmMap := responseMap["confirmedBalance"].(map[string]interface{})
	return &utxo.AccountResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "get Zec balance success",
		Balance: comfirmMap["amount"].(string),
	}, nil
}

func (c *ChainAdaptor) GetUnspentOutputs(req *utxo.UnspentOutputsRequest) (*utxo.UnspentOutputsResponse, error) {
	address := req.GetAddress()
	resp, err := c.zcashDataClientClient.R().Get(c.zcashDataClientClient.RequestUrl + "addresses-latest/utxo/zcash/mainnet/" + address + "/unspent-outputs")
	if err != nil {
		return &utxo.UnspentOutputsResponse{
			Code:           common2.ReturnCode_ERROR,
			Msg:            "get unspend utxo balance fail",
			UnspentOutputs: []*utxo.UnspentOutput{},
		}, err
	}
	resMap := make(map[string]interface{}, 3)
	err = json.Unmarshal([]byte(resp.String()), &resMap)
	itemMap := resMap["data"].(map[string]interface{})
	utxoSlice := itemMap["items"].([]any)
	outputs := make([]*utxo.UnspentOutput, 0)
	for i := 0; i < len(utxoSlice); i++ {
		utxoItem := utxoSlice[i].(map[string]interface{})
		value := utxoItem["value"].(map[string]interface{})
		unspentOutput := &utxo.UnspentOutput{
			TxHashBigEndian: "",
			TxId:            utxoItem["transactionId"].(string),
			TxOutputN:       0,
			Script:          "",
			Address:         utxoItem["address"].(string),
			UnspentAmount:   value["amount"].(string),
			Index:           uint64(int(utxoItem["index"].(float64))),
		}
		outputs = append(outputs, unspentOutput)
	}
	return &utxo.UnspentOutputsResponse{
		Code:           common2.ReturnCode_SUCCESS,
		Msg:            "get unspent outputs success",
		UnspentOutputs: outputs,
	}, nil

}

func (c *ChainAdaptor) GetBlockByNumber(req *utxo.BlockNumberRequest) (*utxo.BlockResponse, error) {
	height := req.GetHeight()
	requestBody := common.RequestData
	requestBody[common.METHOD] = "getblock"
	requestBody[common.PARAMS] = []any{strconv.Itoa(int(height)), 2}
	client := c.zcashClient

	block, err := client.R().SetBody(requestBody).Post(client.RequestUrl)
	if err != nil {
		return &utxo.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get block  fail",
		}, err
	}
	log.Info(block.String())

	allMap := make(map[string]interface{})
	json.Unmarshal([]byte(block.String()), &allMap)
	resultMap := allMap["result"].(map[string]interface{})
	blockHash := resultMap["hash"].(string)
	txSlice := resultMap["tx"].([]any)

	var txList []*utxo.TransactionList

	for i := 0; i < len(txSlice); i++ {
		txMap := txSlice[i].(map[string]interface{})

		if err != nil {
			log.Error("json unmarshal fail", "err", err)
			return nil, err
		}
		var vinSlice = txMap["vin"].([]interface{})
		var vinList []*utxo.Vin
		for i := 0; i < len(vinSlice); i++ {
			vin := vinSlice[i].(map[string]interface{})
			secriptMap, ok := vin["scriptSig"].(map[string]interface{})
			if !ok {
				continue
			}
			vinItem := &utxo.Vin{
				Hash:    vin["txid"].(string),
				Index:   uint32(vin["sequence"].(float64)),
				Amount:  10,
				Address: secriptMap["asm"].(string),
			}
			vinList = append(vinList, vinItem)
		}
		var voutSlice = txMap["vout"].([]interface{})

		var voutList []*utxo.Vout
		for i := 0; i < len(voutSlice); i++ {
			vout := voutSlice[i].(map[string]interface{})
			secriptMap, ok := vout["scriptPubKey"].(map[string]interface{})
			if !ok {
				continue
			}
			voutItem := &utxo.Vout{
				Address: secriptMap["addresses"].([]interface{})[0].(string),
				Amount:  int64(vout["value"].(float64)),
			}
			voutList = append(voutList, voutItem)
		}
		txItem := &utxo.TransactionList{
			Hash: txMap["txid"].(string),
			Vin:  vinList,
			Vout: voutList,
		}
		txList = append(txList, txItem)

	}

	//var resultBlock types.BlockData
	//err = json.Unmarshal(block, &resultBlock)
	//if err != nil {
	//	log.Error("Unmarshal json fail", "err", err)
	//}

	return &utxo.BlockResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "get block by number succcess",
		Height: uint64(req.Height),
		Hash:   blockHash,
		TxList: txList,
	}, nil
	return nil, err
}

func (c *ChainAdaptor) GetBlockByHash(req *utxo.BlockHashRequest) (*utxo.BlockResponse, error) {
	blockHash := req.GetHash()
	requestBody := common.RequestData
	requestBody[common.METHOD] = "getblock"
	requestBody[common.PARAMS] = []any{blockHash, 2}
	client := c.zcashClient

	block, err := client.R().SetBody(requestBody).Post(client.RequestUrl)
	if err != nil {
		return &utxo.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get block  fail",
		}, err
	}
	log.Info(block.String())

	allMap := make(map[string]interface{})
	json.Unmarshal([]byte(block.String()), &allMap)
	resultMap := allMap["result"].(map[string]interface{})
	//blockHash := resultMap["hash"].(string)
	blockHeight := resultMap["height"].(float64)
	txSlice := resultMap["tx"].([]any)

	var txList []*utxo.TransactionList

	for i := 0; i < len(txSlice); i++ {
		txMap := txSlice[i].(map[string]interface{})

		if err != nil {
			log.Error("json unmarshal fail", "err", err)
			return nil, err
		}
		var vinSlice = txMap["vin"].([]interface{})
		var vinList []*utxo.Vin
		for i := 0; i < len(vinSlice); i++ {
			vin := vinSlice[i].(map[string]interface{})
			secriptMap, ok := vin["scriptSig"].(map[string]interface{})
			if !ok {
				continue
			}
			vinItem := &utxo.Vin{
				Hash:    vin["txid"].(string),
				Index:   uint32(vin["sequence"].(float64)),
				Amount:  10,
				Address: secriptMap["asm"].(string),
			}
			vinList = append(vinList, vinItem)
		}
		var voutSlice = txMap["vout"].([]interface{})

		var voutList []*utxo.Vout
		for i := 0; i < len(voutSlice); i++ {
			vout := voutSlice[i].(map[string]interface{})
			secriptMap, ok := vout["scriptPubKey"].(map[string]interface{})
			if !ok {
				continue
			}
			voutItem := &utxo.Vout{
				Address: secriptMap["addresses"].([]interface{})[0].(string),
				Amount:  int64(vout["value"].(float64)),
			}
			voutList = append(voutList, voutItem)
		}
		txItem := &utxo.TransactionList{
			Hash: txMap["txid"].(string),
			Vin:  vinList,
			Vout: voutList,
		}
		txList = append(txList, txItem)

	}

	//var resultBlock types.BlockData
	//err = json.Unmarshal(block, &resultBlock)
	//if err != nil {
	//	log.Error("Unmarshal json fail", "err", err)
	//}

	return &utxo.BlockResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "get block by number succcess",
		Height: uint64(blockHeight),
		Hash:   blockHash,
		TxList: txList,
	}, nil
	return nil, err
}

func (c *ChainAdaptor) GetBlockHeaderByHash(req *utxo.BlockHeaderHashRequest) (*utxo.BlockHeaderResponse, error) {
	hash, err := chainhash.NewHashFromStr(req.Hash)
	if err != nil {
		log.Error("format string to hash fail", "err", err)
	}
	requestBody := common.RequestData
	requestBody[common.METHOD] = "getblockheader"
	requestBody[common.PARAMS] = []any{hash.String()}
	client := c.zcashClient
	resp, err := client.R().SetBody(requestBody).Post(client.RequestUrl)
	if err != nil {
		return &utxo.BlockHeaderResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get block header fail",
		}, err
	}
	resMap := make(map[string]interface{}, 3)
	err = json.Unmarshal([]byte(resp.String()), &resMap)
	responseMap := resMap["result"].(map[string]interface{})
	if err != nil {
		log.Error("Response is not a map[string]interface{}")
		return &utxo.BlockHeaderResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get block header fail",
		}, err
	}
	log.Info("response: " + resp.String())
	return &utxo.BlockHeaderResponse{
		Code:       common2.ReturnCode_SUCCESS,
		Msg:        "get block header success",
		ParentHash: responseMap["previousblockhash"].(string),
		Number:     fmt.Sprintf("%d", int(responseMap["height"].(float64))),
		BlockHash:  req.Hash,
		MerkleRoot: responseMap["merkleroot"].(string),
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByNumber(req *utxo.BlockHeaderNumberRequest) (*utxo.BlockHeaderResponse, error) {

	panic("Dot support this method")
}

func (c *ChainAdaptor) SendTx(req *utxo.SendTxRequest) (*utxo.SendTxResponse, error) {

	requestBody := common.RequestData
	requestBody[common.METHOD] = "sendrawtransaction"
	requestBody[common.PARAMS] = []any{req.RawTx}
	client := c.zcashClient
	resp, err := client.R().SetBody(requestBody).Post(client.RequestUrl)
	if err != nil {
		return &utxo.SendTxResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "send TX fail",
		}, err
	}
	resMap := make(map[string]interface{}, 3)
	err = json.Unmarshal([]byte(resp.String()), &resMap)
	if err != nil {
		return &utxo.SendTxResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	txHash, ok := resMap["result"].(string)
	if !ok {
		errorMap := resMap["error"].(map[string]interface{})
		errMsg := errorMap["message"].(string)
		return &utxo.SendTxResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  errMsg,
		}, err
	}

	return &utxo.SendTxResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "send tx success",
		TxHash: txHash,
	}, nil
}

func (c *ChainAdaptor) GetTxByAddress(req *utxo.TxAddressRequest) (*utxo.TxAddressResponse, error) {
	panic("Not supported now")
}

func (c *ChainAdaptor) GetTxByHash(req *utxo.TxHashRequest) (*utxo.TxHashResponse, error) {
	txHash := req.GetHash()
	resp, err := c.zcashDataClientClient.R().Get(c.zcashDataClientClient.RequestUrl + "transactions/utxo/zcash/mainnet/" + txHash)
	if err != nil {
		return &utxo.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get tx balance fail",
			Tx:   nil,
		}, err
	}
	resMap := make(map[string]interface{}, 3)
	err = json.Unmarshal([]byte(resp.String()), &resMap)
	dataMap := resMap["data"].(map[string]interface{})
	itemMap := dataMap["item"].(map[string]interface{})
	feeMap := itemMap["fee"].(map[string]interface{})
	inputs := itemMap["inputs"].([]interface{})
	outputs := itemMap["outputs"].([]interface{})
	mineMap := itemMap["minedInBlock"].(map[string]interface{})
	height := strconv.FormatFloat(mineMap["height"].(float64), 'f', -1, 64)
	timestamp := strconv.FormatFloat(itemMap["timestamp"].(float64), 'f', -1, 64)
	from_addrs := make([]*utxo.Address, 0)
	fee := feeMap["amount"].(string)
	for _, input := range inputs {
		inputMap := input.(map[string]interface{})
		addressesSlice := inputMap["addresses"].([]interface{})
		for _, address := range addressesSlice {
			from_addrs = append(from_addrs, &utxo.Address{Address: address.(string)})
		}
	}
	to_addrs := make([]*utxo.Address, 0)
	value_list := make([]*utxo.Value, 0)
	for _, output := range outputs {
		outputMap := output.(map[string]interface{})
		addressesSlice := outputMap["addresses"].([]interface{})
		valueMap := outputMap["value"].(map[string]interface{})
		for _, address := range addressesSlice {
			to_addrs = append(to_addrs, &utxo.Address{Address: address.(string)})
		}
		value_list = append(value_list, &utxo.Value{Value: valueMap["amount"].(string)})
	}
	transactionHash := itemMap["hash"].(string)

	txMsg := &utxo.TxMessage{
		Hash:     transactionHash,
		Froms:    from_addrs,
		Tos:      to_addrs,
		Values:   value_list,
		Fee:      fee,
		Status:   utxo.TxStatus_Success,
		Type:     0,
		Height:   height,
		Datetime: timestamp,
	}
	return &utxo.TxHashResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get transaction success",
		Tx:   txMsg,
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

func (c *ChainAdaptor) BuildSignedTransaction(req *utxo.SignedTransactionRequest) (*utxo.SignedTransactionResponse, error) {
	r := bytes.NewReader(req.TxData)
	var msgTx wire.MsgTx
	err := msgTx.Deserialize(r)
	if err != nil {
		log.Error("Create signed transaction msg tx deserialize", "err", err)
		return &utxo.SignedTransactionResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	if len(req.Signatures) != len(msgTx.TxIn) {
		log.Error("CreateSignedTransaction invalid params", "err", "Signature number mismatch Txin number")
		err = errors.New("Signature number != Txin number")
		return &utxo.SignedTransactionResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	if len(req.PublicKeys) != len(msgTx.TxIn) {
		log.Error("CreateSignedTransaction invalid params", "err", "Pubkey number mismatch Txin number")
		err = errors.New("Pubkey number != Txin number")
		return &utxo.SignedTransactionResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	for i, in := range msgTx.TxIn {
		txHashRequest := &utxo.TxHashRequest{
			Chain: ChainName,
			Hash:  in.PreviousOutPoint.Hash.String(),
		}

		preTx, err2 := c.GetTxByHash(txHashRequest)
		if err2 != nil {
			log.Error("CreateSignedTransaction GetRawTransactionVerbose", "err", err2)
			return &utxo.SignedTransactionResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  err2.Error(),
			}, err2
		}
		fromAddress := preTx.Tx.Tos[in.PreviousOutPoint.Index].Address
		log.Info("CreateSignedTransaction ", "from address", fromAddress)

		ripemd160Result, _ := extractPubKeyHashFromAddress(fromAddress)

		fromPkScript, err2 := txscript.NewScriptBuilder().AddOp(txscript.OP_DUP).AddOp(txscript.OP_HASH160).
			AddData(ripemd160Result).AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG).
			Script()

		if err2 != nil {
			log.Error("CreateSignedTransaction PayToAddrScript", "err", err2)
			return &utxo.SignedTransactionResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  err2.Error(),
			}, err2
		}

		if len(req.Signatures[i]) < 64 {
			err2 = errors.New("Invalid signature length")
			return &utxo.SignedTransactionResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  err2.Error(),
			}, err2
		}

		// 初始化 r 和 s
		r := new(btcec.ModNScalar)
		s := new(btcec.ModNScalar)

		// 从签名中提取前 32 字节作为 R 和后 32 字节作为 S
		r.SetBytes((*[32]byte)(req.Signatures[i][0:32]))
		s.SetBytes((*[32]byte)(req.Signatures[i][32:64]))

		btcecSig := ecdsa.NewSignature(r, s)
		if err != nil {
			return nil, fmt.Errorf("无法从签名中恢复公钥: %v", err)
		}
		sig := append(btcecSig.Serialize(), byte(txscript.SigHashAll))

		btcecPub, err2 := btcec.ParsePubKey(req.PublicKeys[i])
		if err2 != nil {
			log.Error("CreateSignedTransaction ParsePubKey", "err", err2)
			return &utxo.SignedTransactionResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  err2.Error(),
			}, err2
		}
		var pkData []byte
		if btcec.IsCompressedPubKey(req.PublicKeys[i]) {
			pkData = btcecPub.SerializeCompressed()
		} else {
			pkData = btcecPub.SerializeUncompressed()
		}

		sigScript, err2 := txscript.NewScriptBuilder().AddData(sig).AddData(pkData).Script()
		if err2 != nil {
			log.Error("create signed transaction new script builder", "err", err2)
			return &utxo.SignedTransactionResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  err2.Error(),
			}, err2
		}
		//fmt.Println("txHash:", txHashRequest.Hash)
		//fmt.Println("signature:", hex.EncodeToString(sig))
		//fmt.Println("fromPKScript:", hex.EncodeToString(fromPkScript))
		//fmt.Println("sigScript:", hex.EncodeToString(sigScript))
		//fmt.Println("publicKey:", hex.EncodeToString(req.PublicKeys[i]))
		//fmt.Println("pkData", hex.EncodeToString(pkData))
		publicKeyHash, _ := extractPubKeyHashFromAddress(preTx.Tx.Tos[i].Address)
		fmt.Println("publicKeyHash", hex.EncodeToString(publicKeyHash))
		msgTx.TxIn[i].SignatureScript = sigScript
		amountStr, _ := strconv.ParseFloat(preTx.Tx.Values[in.PreviousOutPoint.Index].Value, 64)
		amount := zecToSatoshi(amountStr)
		log.Info("CreateSignedTransaction ", "amount", preTx.Tx.Values[in.PreviousOutPoint.Index].Value, "int amount", amount)

		vm, err2 := txscript.NewEngine(fromPkScript, &msgTx, i, txscript.StandardVerifyFlags, nil, nil, amount.Int64(), nil)
		if err2 != nil {
			log.Error("create signed transaction newEngine", "err", err2)
			return &utxo.SignedTransactionResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  err2.Error(),
			}, err2
		}
		if err3 := vm.Execute(); err3 != nil {
			log.Error("CreateSignedTransaction NewEngine Execute", "err", err3)
			return &utxo.SignedTransactionResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  err3.Error(),
			}, err3
		}
	}
	// serialize tx
	buf := bytes.NewBuffer(make([]byte, 0, msgTx.SerializeSize()))
	err = msgTx.Serialize(buf)
	if err != nil {
		log.Error("CreateSignedTransaction tx Serialize", "err", err)
		return &utxo.SignedTransactionResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	hash := msgTx.TxHash()
	return &utxo.SignedTransactionResponse{
		Code:         common2.ReturnCode_SUCCESS,
		SignedTxData: buf.Bytes(),
		Hash:         (&hash).CloneBytes(),
	}, nil
}

func (c *ChainAdaptor) DecodeTransaction(req *utxo.DecodeTransactionRequest) (*utxo.DecodeTransactionResponse, error) {
	panic("Not supported now")
}

func (c *ChainAdaptor) VerifySignedTransaction(req *utxo.VerifyTransactionRequest) (*utxo.VerifyTransactionResponse, error) {
	panic("Not supported now")
}

func (c *ChainAdaptor) CalcSignHashes(Vins []*utxo.Vin, Vouts []*utxo.Vout) ([][]byte, []byte, error) {
	if len(Vins) == 0 || len(Vouts) == 0 {
		return nil, nil, errors.New("invalid len in or out")
	}
	rawTx := wire.NewMsgTx(wire.TxVersion)
	for _, in := range Vins {
		utxoHash, err := chainhash.NewHashFromStr(in.Hash)
		if err != nil {
			return nil, nil, err
		}
		txIn := wire.NewTxIn(wire.NewOutPoint(utxoHash, in.Index), nil, nil)
		rawTx.AddTxIn(txIn)
	}
	for _, out := range Vouts {
		//toAddress, err := btcutil.DecodeAddress(out.Address, &chaincfg.MainNetParams)
		//if err != nil {
		//	return nil, nil, err
		//}
		// 提取公钥哈希
		pubKeyHash, err := extractPubKeyHashFromAddress(out.Address)
		if err != nil {
			fmt.Println("Error:", err)
		}

		toPkScript, err := txscript.NewScriptBuilder().AddOp(txscript.OP_DUP).AddOp(txscript.OP_HASH160).
			AddData(pubKeyHash).AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG).
			Script()

		//toPkScript, err := txscript.PayToAddrScript(toAddress)
		if err != nil {
			return nil, nil, err
		}
		rawTx.AddTxOut(wire.NewTxOut(out.Amount, toPkScript))
	}
	signHashes := make([][]byte, len(Vins))
	for i, in := range Vins {
		from := in.Address
		//fromAddr, err := btcutil.DecodeAddress(from, &chaincfg.MainNetParams)
		//if err != nil {
		//	log.Info("decode address error", "from", from, "err", err)
		//	return nil, nil, err
		//}
		//fromPkScript, err := txscript.PayToAddrScript(fromAddr)
		// 提取公钥哈希
		pubKeyHash, err := extractPubKeyHashFromAddress(from)
		if err != nil {
			fmt.Println("Error:", err)
		}

		fromPkScript, err := txscript.NewScriptBuilder().AddOp(txscript.OP_DUP).AddOp(txscript.OP_HASH160).
			AddData(pubKeyHash).AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG).
			Script()
		if err != nil {
			log.Info("pay to addr script err", "err", err)
			return nil, nil, err
		}
		signHash, err := txscript.CalcSignatureHash(fromPkScript, txscript.SigHashAll, rawTx, i)
		if err != nil {
			log.Info("Calc signature hash error", "err", err)
			return nil, nil, err
		}
		signHashes[i] = signHash
	}
	var buffer bytes.Buffer
	err := rawTx.Serialize(&buffer)
	if err != nil {
		log.Info("serialize fail")
	}
	//buf := bytes.NewBuffer(make([]byte, 0, rawTx.SerializeSize()))
	log.Info("rawTx:", hex.EncodeToString(buffer.Bytes()))
	return signHashes, buffer.Bytes(), nil
}

func extractPubKeyHashFromAddress(address string) ([]byte, error) {
	// Step 1: Base58 解码地址
	addressBytes := base58.Decode(address)
	if len(addressBytes) < 4 {
		return nil, fmt.Errorf("invalid address length")
	}

	// Step 2: 提取版本字节、公钥哈希和校验和
	versionByte := addressBytes[0]
	pubKeyHash := addressBytes[1 : len(addressBytes)-4] // 公钥哈希 20 字节
	checksum := addressBytes[len(addressBytes)-4:]      // 校验和 4 字节

	// Step 3: 验证校验和
	payload := append([]byte{versionByte}, pubKeyHash...)
	// 计算校验和
	checksumCalc := sha256.Sum256(payload)
	checksumCalc = sha256.Sum256(checksumCalc[:])

	// 校验计算出的校验和是否和提取的校验和一致
	if !bytes.Equal(checksumCalc[:4], checksum) {
		return nil, fmt.Errorf("invalid checksum")
	}

	// 返回公钥哈希
	return pubKeyHash, nil
}
