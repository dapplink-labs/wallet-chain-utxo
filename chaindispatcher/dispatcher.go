package chaindispatcher

import (
	"context"
	"runtime/debug"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ethereum/go-ethereum/log"

	"github.com/dapplink-labs/wallet-chain-utxo/chain"
	"github.com/dapplink-labs/wallet-chain-utxo/chain/bitcoin"
	"github.com/dapplink-labs/wallet-chain-utxo/chain/bitcoincash"
	"github.com/dapplink-labs/wallet-chain-utxo/chain/dash"
	"github.com/dapplink-labs/wallet-chain-utxo/chain/litecoin"
	"github.com/dapplink-labs/wallet-chain-utxo/chain/zen"
	"github.com/dapplink-labs/wallet-chain-utxo/config"
	"github.com/dapplink-labs/wallet-chain-utxo/rpc/common"
	"github.com/dapplink-labs/wallet-chain-utxo/rpc/utxo"
)

type CommonRequest interface {
	GetChain() string
}

type CommonReply = utxo.SupportChainsResponse

type ChainType = string

type ChainDispatcher struct {
	registry map[ChainType]chain.IChainAdaptor
}

func New(conf *config.Config) (*ChainDispatcher, error) {
	dispatcher := ChainDispatcher{
		registry: make(map[ChainType]chain.IChainAdaptor),
	}
	chainAdaptorFactoryMap := map[string]func(conf *config.Config) (chain.IChainAdaptor, error){
		strings.ToLower(bitcoin.ChainName):     bitcoin.NewChainAdaptor,
		strings.ToLower(bitcoincash.ChainName): bitcoincash.NewChainAdaptor,
		strings.ToLower(dash.ChainName):        dash.NewChainAdaptor,
		strings.ToLower(litecoin.ChainName):    litecoin.NewChainAdaptor,
		strings.ToLower(zen.ChainName):         zen.NewChainAdaptor,
	}
	supportedChains := []string{
		strings.ToLower(bitcoin.ChainName),
		strings.ToLower(bitcoincash.ChainName),
		strings.ToLower(dash.ChainName),
		strings.ToLower(litecoin.ChainName),
		strings.ToLower(zen.ChainName),
	}
	for _, c := range conf.Chains {
		chainName := strings.ToLower(c)
		if factory, ok := chainAdaptorFactoryMap[chainName]; ok {
			adaptor, err := factory(conf)
			if err != nil {
				log.Crit("failed to setup chain", "chain", chainName, "error", err)
			}
			dispatcher.registry[chainName] = adaptor
		} else {
			log.Error("unsupported chain", "chain", chainName, "supportedChains", supportedChains)
		}
	}
	return &dispatcher, nil
}

func (d *ChainDispatcher) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error("panic error", "msg", e)
			log.Debug(string(debug.Stack()))
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()

	pos := strings.LastIndex(info.FullMethod, "/")
	method := info.FullMethod[pos+1:]

	chainName := strings.ToLower(req.(CommonRequest).GetChain())
	log.Info(method, "chain", chainName, "req", req)

	resp, err = handler(ctx, req)
	log.Debug("Finish handling", "resp", resp, "err", err)
	return
}

func (d *ChainDispatcher) preHandler(req interface{}) (resp *CommonReply, chainName string) {
	chainName = strings.ToLower(req.(CommonRequest).GetChain())
	if _, ok := d.registry[chainName]; !ok {
		return &CommonReply{
			Code:    common.ReturnCode_ERROR,
			Msg:     config.UnsupportedOperation,
			Support: false,
		}, chainName
	}
	return nil, chainName
}

func (d *ChainDispatcher) GetSupportChains(ctx context.Context, request *utxo.SupportChainsRequest) (*utxo.SupportChainsResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &utxo.SupportChainsResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedOperation,
		}, nil
	}
	return d.registry[chainName].GetSupportChains(request)
}

func (d *ChainDispatcher) ConvertAddress(ctx context.Context, request *utxo.ConvertAddressRequest) (*utxo.ConvertAddressResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &utxo.ConvertAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "covert address fail at pre handle",
		}, nil
	}
	return d.registry[chainName].ConvertAddress(request)
}

func (d *ChainDispatcher) ValidAddress(ctx context.Context, request *utxo.ValidAddressRequest) (*utxo.ValidAddressResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &utxo.ValidAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "valid address error at pre handle",
		}, nil
	}
	return d.registry[chainName].ValidAddress(request)
}

func (d *ChainDispatcher) GetFee(ctx context.Context, request *utxo.FeeRequest) (*utxo.FeeResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &utxo.FeeResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get fee fail at pre handle",
		}, nil
	}
	return d.registry[chainName].GetFee(request)
}

func (d *ChainDispatcher) GetAccount(ctx context.Context, request *utxo.AccountRequest) (*utxo.AccountResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &utxo.AccountResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get account information fail at pre handle",
		}, nil
	}
	return d.registry[chainName].GetAccount(request)
}

func (d *ChainDispatcher) GetUnspentOutputs(ctx context.Context, request *utxo.UnspentOutputsRequest) (*utxo.UnspentOutputsResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &utxo.UnspentOutputsResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get un spend out fail at pre handle",
		}, nil
	}
	return d.registry[chainName].GetUnspentOutputs(request)
}

func (d *ChainDispatcher) GetBlockByNumber(ctx context.Context, request *utxo.BlockNumberRequest) (*utxo.BlockResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &utxo.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block by number fail at pre handle",
		}, nil
	}
	return d.registry[chainName].GetBlockByNumber(request)
}

func (d *ChainDispatcher) GetBlockByHash(ctx context.Context, request *utxo.BlockHashRequest) (*utxo.BlockResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &utxo.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block by hash fail at pre handle",
		}, nil
	}
	return d.registry[chainName].GetBlockByHash(request)
}

func (d *ChainDispatcher) GetBlockHeaderByHash(ctx context.Context, request *utxo.BlockHeaderHashRequest) (*utxo.BlockHeaderResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &utxo.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block header by hash fail at pre handle",
		}, nil
	}
	return d.registry[chainName].GetBlockHeaderByHash(request)
}

func (d *ChainDispatcher) GetBlockHeaderByNumber(ctx context.Context, request *utxo.BlockHeaderNumberRequest) (*utxo.BlockHeaderResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &utxo.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block header by number fail at pre handle",
		}, nil
	}
	return d.registry[chainName].GetBlockHeaderByNumber(request)
}

func (d *ChainDispatcher) SendTx(ctx context.Context, request *utxo.SendTxRequest) (*utxo.SendTxResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &utxo.SendTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "send tx fail at pre handle",
		}, nil
	}
	return d.registry[chainName].SendTx(request)
}

func (d *ChainDispatcher) GetTxByAddress(ctx context.Context, request *utxo.TxAddressRequest) (*utxo.TxAddressResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &utxo.TxAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get tx by address fail pre handle",
		}, nil
	}
	return d.registry[chainName].GetTxByAddress(request)
}

func (d *ChainDispatcher) GetTxByHash(ctx context.Context, request *utxo.TxHashRequest) (*utxo.TxHashResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &utxo.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get tx by hash fail at pre handle",
		}, nil
	}
	return d.registry[chainName].GetTxByHash(request)
}

func (d *ChainDispatcher) CreateUnSignTransaction(ctx context.Context, request *utxo.UnSignTransactionRequest) (*utxo.UnSignTransactionResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &utxo.UnSignTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get un sign tx fail at pre handle",
		}, nil
	}
	return d.registry[chainName].CreateUnSignTransaction(request)
}

func (d *ChainDispatcher) BuildSignedTransaction(ctx context.Context, request *utxo.SignedTransactionRequest) (*utxo.SignedTransactionResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &utxo.SignedTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "signed tx fail at pre handle",
		}, nil
	}
	return d.registry[chainName].BuildSignedTransaction(request)
}

func (d *ChainDispatcher) DecodeTransaction(ctx context.Context, request *utxo.DecodeTransactionRequest) (*utxo.DecodeTransactionResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &utxo.DecodeTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "decode tx fail at pre handle",
		}, nil
	}
	return d.registry[chainName].DecodeTransaction(request)
}

func (d *ChainDispatcher) VerifySignedTransaction(ctx context.Context, request *utxo.VerifyTransactionRequest) (*utxo.VerifyTransactionResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &utxo.VerifyTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "verify tx fail at pre handle",
		}, nil
	}
	return d.registry[chainName].VerifySignedTransaction(request)
}
