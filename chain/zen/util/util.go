package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"math"
	"strconv"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/dapplink-labs/wallet-chain-utxo/chain/zen/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

const messagePrefix = "ZENCash main net"

var pubKeyHash = []byte{0x20, 0x89}
var scriptHash = []byte{0x20, 0x96}
var zcPaymentAddressHash = 0x169a
var zcSpendingKeyHash = 0xab36
var WIF = 0x80
var bip32PublicKey = 0x0488b21e
var bip32PrivateKey = 0x0488ade4

const (
	OP_0                  = "00"
	OP_1                  = "51"
	OP_2                  = "52"
	OP_3                  = "53"
	OP_4                  = "54"
	OP_DUP                = "76"
	OP_NIP                = "77"
	OP_OVER               = "78"
	OP_HASH160            = "a9"
	OP_EQUALVERIFY        = "88"
	OP_CHECKSIG           = "ac"
	OP_CHECKBLOCKATHEIGHT = "b4"
	OP_EQUAL              = "87"
	OP_REVERSED           = "89"
	OP_CHECKMULTISIG      = "ae"
	OP_PUSHDATA1          = "4c"
	OP_RETURN             = "6a"

	SIGHASH_ALL          = 1
	SIGHASH_NONE         = 2
	SIGHASH_SINGLE       = 3
	SIGHASH_ANYONECANPAY = 0x80
)

func OPGetNum(op string) int {
	num, _ := strconv.Atoi(op)
	return num
}

func ReadVarInt(buf []byte, offset int) (uint64, int) {
	if offset >= len(buf) {
		return 0, 0
	}

	first := buf[offset]
	offset++

	switch {
	case first < 0xfd:
		return uint64(first), 1
	case first == 0xfd:
		if offset+2 > len(buf) {
			return 0, 0
		}
		return uint64(binary.LittleEndian.Uint16(buf[offset:])), 3
	case first == 0xfe:
		if offset+4 > len(buf) {
			return 0, 0
		}
		return uint64(binary.LittleEndian.Uint32(buf[offset:])), 5
	case first == 0xff:
		if offset+8 > len(buf) {
			return 0, 0
		}
		return binary.LittleEndian.Uint64(buf[offset:]), 9
	}

	return 0, 0
}

func ReverseBytes(data []byte) []byte {
	length := len(data)
	reversed := make([]byte, length)

	for i := 0; i < length; i++ {
		reversed[i] = data[length-1-i]
	}

	return reversed
}

func ConvertToSatoshi(valueStr string) int64 {
	valueDecimal, err := decimal.NewFromString(valueStr)
	if err != nil {
		return 0
	}

	multiplier := decimal.NewFromFloat(1e8)
	result := valueDecimal.Mul(multiplier)

	intValue := result.IntPart()
	return intValue
}

func ValidAddress(address string) (validAddress bool, msg string, err error) {
	addressBytes := base58.Decode(address)
	validAddress = true

	// check prefix
	if len(addressBytes) >= 7 && bytes.HasPrefix(addressBytes, pubKeyHash) {
		msg = "verify address success, It is p2pkh address"
	} else if len(addressBytes) >= 7 && bytes.HasPrefix(addressBytes, scriptHash) {
		msg = "verify address success, It is p2sh address"
	} else {
		msg = "verify address fail, It is not a valid address"
		validAddress = false
		err = errors.New(msg)
		return
	}

	// check checksum
	if len(addressBytes) < 7 {
		msg = "verify address fail, address length is too short"
		validAddress = false
		err = errors.New(msg)
		return
	} else {
		splitHash := addressBytes[2 : len(addressBytes)-4]
		splitChecksum := addressBytes[len(addressBytes)-4:]
		computeChecksum := checksum(append(pubKeyHash, splitHash...))
		if !bytes.Equal(splitChecksum, computeChecksum) {
			msg = "verify address fail, checksum error"
			validAddress = false
			err = errors.New(msg)
			return
		}
	}
	return
}

func DoubleSha256(data []byte) []byte {
	hash := sha256.Sum256(data)
	hash = sha256.Sum256(hash[:])
	return hash[:]
}

func checksum(payload []byte) []byte {
	return DoubleSha256(payload)[:4]
}

func PublicKeyToAddress(publicKey string) (string, error) {
	compressedPubKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return "", err
	}
	hash160 := btcutil.Hash160(compressedPubKeyBytes)
	payload := append(pubKeyHash, hash160...)
	address := base58.Encode(append(payload, checksum(payload)...))
	return address, nil
}

func CreateRawTX(txObjectInputParam *types.TXObjectInputParam) (*types.TXObject, error) {
	hisotory := txObjectInputParam.ParamIns
	recipients := txObjectInputParam.ParamOuts
	blockHeight := txObjectInputParam.BlockHeight
	blockHash := txObjectInputParam.BlockHash

	txObject := &types.TXObject{
		LockTime: 0,
		Version:  1,
		Ins:      []types.Ins{},
		Outs:     []types.Outs{},
	}
	for _, h := range hisotory {
		txObject.Ins = append(txObject.Ins, types.Ins{
			Output: types.Output{
				Hash: h.TXID,
				Vout: h.Vout,
			},
			Script:           "",
			PrevScriptPubKey: h.ScriptPubKey,
			Sequence:         "ffffffff",
		})
	}

	for _, o := range recipients {
		script, err := addressToScript(o.Address, blockHeight, blockHash)
		if err != nil {
			return nil, err
		}
		txObject.Outs = append(txObject.Outs, types.Outs{
			Script:   script,
			Satoshis: o.Satoshis,
		})
	}
	return txObject, nil
}

func addressToScript(address string, blockHeight int64, blockHash string) (string, error) {
	addressBytes := base58.Decode(address)

	if bytes.HasPrefix(addressBytes, scriptHash) {
		return mkScriptHashReplayScript(address, blockHeight, blockHash)
	} else if bytes.HasPrefix(addressBytes, pubKeyHash) {
		return mkPubkeyHashReplayScript(address, blockHeight, blockHash)
	}
	return "", errors.New("addressToScript error: invalid address")
}

func mkPubkeyHashReplayScript(address string, blockHeight int64, blockHash string) (string, error) {
	data, version, _ := base58.CheckDecode(address)
	addrHex := hex.EncodeToString(append([]byte{version}, data...))
	subAddrHex := addrHex[len(hex.EncodeToString(pubKeyHash)):]
	blockHashHex, err := hex.DecodeString(blockHash)
	if err != nil {
		return "", err
	}
	RevertedBlockHashHex := hex.EncodeToString(ReverseBytes(blockHashHex))
	subAddrHexLength, err := GetPushDataLength(subAddrHex)
	if err != nil {
		return "", err
	}

	blockHashHexLength, err := GetPushDataLength(RevertedBlockHashHex)
	if err != nil {
		return "", err
	}
	s, err := serializeScriptBlockHeight(blockHeight)
	if err != nil {
		return "", err
	}

	return OP_DUP + OP_HASH160 + subAddrHexLength + subAddrHex + OP_EQUALVERIFY + OP_CHECKSIG + blockHashHexLength + RevertedBlockHashHex + s + OP_CHECKBLOCKATHEIGHT, nil
}

func mkScriptHashReplayScript(address string, blockHeight int64, blockHash string) (string, error) {
	data, version, _ := base58.CheckDecode(address)
	addrHex := hex.EncodeToString(append([]byte{version}, data...))

	subAddrHex := addrHex[len(hex.EncodeToString(scriptHash)):]

	blockHashHex, err := hex.DecodeString(blockHash)
	if err != nil {
		return "", err
	}
	sah, err := GetPushDataLength(subAddrHex)
	if err != nil {
		return "", err
	}
	RevertedBlockHashHex := hex.EncodeToString(ReverseBytes(blockHashHex))
	bhh, err := GetPushDataLength(RevertedBlockHashHex)
	if err != nil {
		return "", err
	}
	sbh, err := serializeScriptBlockHeight(blockHeight)
	if err != nil {
		return "", err
	}

	return OP_HASH160 + sah + subAddrHex + OP_EQUAL + bhh + RevertedBlockHashHex + sbh + OP_CHECKBLOCKATHEIGHT, nil
}

func GetPushDataLength(s string) (string, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return "", errors.Errorf("hex decode error: %v", err)
	}

	hexLength := len(data)

	return NumToVarInt(uint64(hexLength)), nil
}

func NumToVarInt(num uint64) string {
	var buf bytes.Buffer

	switch {
	case num < 0xfd:
		buf.WriteByte(byte(num))
	case num <= 0xffff:
		buf.WriteByte(0xfd)
		binary.Write(&buf, binary.LittleEndian, uint16(num))
	case num <= 0xffffffff:
		buf.WriteByte(0xfe)
		binary.Write(&buf, binary.LittleEndian, uint32(num))
	default:
		buf.WriteByte(0xff)
		binary.Write(&buf, binary.LittleEndian, num)
	}

	return hex.EncodeToString(buf.Bytes())
}

func serializeScriptBlockHeight(blockHeight int64) (string, error) {
	// check for scriptNum special case values
	if blockHeight >= -1 && blockHeight <= 16 {
		var res = 0

		if blockHeight == 0 || blockHeight >= 1 && blockHeight <= 16 {

			res = int(blockHeight) + OPGetNum(OP_1) - 1
		} else if blockHeight == 0 {
			res = OPGetNum(OP_0)
		}

		return strconv.Itoa(res), nil
	} else {
		// Minimal encoding
		var blockHeightBuffer = scriptNumEncode(blockHeight)
		var blockHeightHex = hex.EncodeToString(blockHeightBuffer)
		blockHeightHexLength, err := GetPushDataLength(blockHeightHex)
		if err != nil {
			return "", err
		}
		return blockHeightHexLength + blockHeightHex, nil
	}
}

func scriptNumEncode(_number int64) []byte {
	value := _number
	var negative = value < 0
	if negative {
		value = -value
	}

	var size = scriptNumSize(value)
	buffer := make([]byte, size)

	for i := 0; i < size; i++ {
		buffer[i] = byte(value & 0xff)
		value >>= 8
	}

	if buffer[size-1]&0x80 != 0 {
		if negative {
			buffer[size-1] = 0x80
		} else {
			buffer[size-1] = 0x00
		}
	} else if negative {
		buffer[size-1] |= 0x80
	}

	return buffer
}

func scriptNumSize(_number int64) int {
	if _number > 0x7fffffff {
		return 5
	} else if _number > 0x7fffff {
		return 4
	} else if _number > 0x7fff {
		return 3
	} else if _number > 0x7f {
		return 2
	} else if _number >= 0x00 {
		return 1
	} else {
		return 0
	}
}

func SignatureForm(txObject *types.TXObject, i int, script string, hashcode uint64) (*types.TXObject, error) {
	var newTx types.TXObject
	err := DeepCopy(&txObject, &newTx)
	if err != nil {
		return nil, err
	}

	for j := 0; j < len(newTx.Ins); j++ {
		newTx.Ins[j].Script = ""
	}

	if hashcode == SIGHASH_NONE {
		newTx.Outs = []types.Outs{}
	} else if hashcode == SIGHASH_SINGLE {
		newTx.Outs = newTx.Outs[0:len(newTx.Ins)]

		for _j := 0; _j < len(newTx.Ins)-1; _j++ {
			newTx.Outs[_j].Satoshis = math.MaxUint64
			newTx.Outs[_j].Script = ""
		}
	} else if hashcode == SIGHASH_ANYONECANPAY {
		newTx.Ins = newTx.Ins[i : i+1]
	}

	newTx.Ins[i].Script = script
	return &newTx, nil
}

func DeepCopy(src, dst interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

func GetScriptSignature(txObject *types.TXObject, hashcode uint64) ([]byte, error) {
	// Buffer
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint16(buf, uint16(hashcode))
	hashcodeHex := hex.EncodeToString(buf)

	signingTxHex, err := SerializeTX(txObject, false)
	if err != nil {
		return nil, err
	}

	signingTxWithHashcode, err := hex.DecodeString(hex.EncodeToString(signingTxHex) + hashcodeHex)
	if err != nil {
		return nil, err
	}

	msg := DoubleSha256(signingTxWithHashcode)

	return msg[:], nil
}

func SerializeTX(txObj *types.TXObject, withPrevScriptPubKey bool) ([]byte, error) {
	var buffer bytes.Buffer

	versionBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(versionBuf, txObj.Version)
	buffer.Write(versionBuf)

	insLenBytes, err := hex.DecodeString(NumToVarInt(uint64(len(txObj.Ins))))
	if err != nil {
		return nil, err
	}
	buffer.Write(insLenBytes)

	for _, input := range txObj.Ins {
		hashBytes, err := hex.DecodeString(input.Output.Hash)
		if err != nil {
			return nil, err
		}
		for i, j := 0, len(hashBytes)-1; i < j; i, j = i+1, j-1 {
			hashBytes[i], hashBytes[j] = hashBytes[j], hashBytes[i]
		}
		buffer.Write(hashBytes)

		voutBuf := make([]byte, 4)
		binary.LittleEndian.PutUint32(voutBuf, uint32(input.Output.Vout))
		buffer.Write(voutBuf)

		if withPrevScriptPubKey {
			pushDataLen, err := GetPushDataLength(input.PrevScriptPubKey)
			if err != nil {
				return nil, err
			}
			pushDataLenBytes, err := hex.DecodeString(pushDataLen)
			if err != nil {
				return nil, err
			}
			buffer.Write(pushDataLenBytes)

			prevScriptBytes, err := hex.DecodeString(input.PrevScriptPubKey)
			if err != nil {
				return nil, err
			}
			buffer.Write(prevScriptBytes)
		}

		scriptPushDataLen, err := GetPushDataLength(input.Script)
		if err != nil {
			return nil, err
		}
		scriptPushDataLenBytes, err := hex.DecodeString(scriptPushDataLen)
		if err != nil {
			return nil, err
		}
		buffer.Write(scriptPushDataLenBytes)

		scriptBytes, err := hex.DecodeString(input.Script)
		if err != nil {
			return nil, err
		}
		buffer.Write(scriptBytes)

		sequenceBytes, err := hex.DecodeString(input.Sequence)
		if err != nil {
			return nil, err
		}
		buffer.Write(sequenceBytes)
	}

	outsLenBytes, err := hex.DecodeString(NumToVarInt(uint64(len(txObj.Outs))))
	if err != nil {
		return nil, err
	}
	buffer.Write(outsLenBytes)

	for _, out := range txObj.Outs {
		satoshisBuf := make([]byte, 8)
		binary.LittleEndian.PutUint32(satoshisBuf, uint32(out.Satoshis&0xFFFFFFFF))
		binary.LittleEndian.PutUint32(satoshisBuf[4:], uint32(out.Satoshis>>32))
		buffer.Write(satoshisBuf)

		scriptPushDataLen, err := GetPushDataLength(out.Script)
		if err != nil {
			return nil, err
		}
		scriptPushDataLenBytes, err := hex.DecodeString(scriptPushDataLen)
		if err != nil {
			return nil, err
		}
		buffer.Write(scriptPushDataLenBytes)

		scriptBytes, err := hex.DecodeString(out.Script)
		if err != nil {
			return nil, err
		}
		buffer.Write(scriptBytes)
	}

	locktimeBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(locktimeBuf, uint32(txObj.LockTime))
	buffer.Write(locktimeBuf)

	return buffer.Bytes(), nil
}
