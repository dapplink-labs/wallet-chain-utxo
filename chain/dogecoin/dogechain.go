package dogecoin

import (
	"fmt"

	"github.com/dapplink-labs/wallet-chain-utxo/chain/dogecoin/types"

	gresty "github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

var errDogeChainHTTPError = errors.New("dogechain http error")

type DogeClient struct {
	client *gresty.Client
}

func NewDogeChainClient(url string) (*DogeClient, error) {
	if url == "" {
		return nil, fmt.Errorf("dogechain URL cannot be empty")
	}

	client := gresty.New()
	client.SetHostURL(url)
	client.OnAfterResponse(func(c *gresty.Client, r *gresty.Response) error {
		statusCode := r.StatusCode()
		if statusCode >= 400 {
			method := r.Request.Method
			url := r.Request.URL
			return fmt.Errorf("%d cannot %s %s: %w", statusCode, method, url, errDogeChainHTTPError)
		}
		return nil
	})
	return &DogeClient{
		client: client,
	}, nil
}

// 实现获取余额等方法
func (c *DogeClient) GetAccountBalance(address string) (string, error) {
	var accountBalance map[string]*types.AccountBalance
	response, err := c.client.R().
		SetResult(&accountBalance).
		Get("/balance?active=" + address)
	if err != nil {
		return "", fmt.Errorf("cannot get account balance: %w", err)
	}
	if response.StatusCode() != 200 {
		return "", errors.New("get account balance fail")
	}
	return accountBalance[address].FinalBalance.String(), nil
}

func (c *DogeClient) GetAccountUtxo(address string) ([]types.UnspentOutput, error) {
	var utxoList types.UnspentOutputList
	response, err := c.client.R().
		SetResult(&utxoList).
		Get("/unspent?active=" + address)
	if err != nil {
		return nil, fmt.Errorf("cannot get utxo: %w", err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get account utxo fail")
	}
	return utxoList.UnspentOutputs, nil
}

func (c *DogeClient) GetTransactionsByAddress(address, page, pageSize string) (*types.Transaction, error) {
	var txList types.Transaction
	response, err := c.client.R().
		SetResult(&txList).
		Get(fmt.Sprintf("/address/%s/transactions?page=%s&pageSize=%s", address, page, pageSize))
	if err != nil {
		return nil, fmt.Errorf("cannot get transactions: %w", err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get transactions fail")
	}
	return &txList, nil
}

func (c *DogeClient) GetTransactionByHash(txHash string) (*types.TxsItem, error) {
	var tx types.TxsItem
	response, err := c.client.R().
		SetResult(&tx).
		Get("/transaction/" + txHash)
	if err != nil {
		return nil, fmt.Errorf("cannot get transaction: %w", err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get transaction fail")
	}
	return &tx, nil
}

// ... 其他方法实现
