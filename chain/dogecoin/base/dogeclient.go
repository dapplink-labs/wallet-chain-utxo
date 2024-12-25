package base

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type DogeClient struct {
	client  *resty.Client
	baseURL string
	token   string
}

func NewDogeClient(baseURL, token string) (*DogeClient, error) {
	client := resty.New()

	return &DogeClient{
		client:  client,
		baseURL: baseURL,
		token:   token,
	}, nil
}

// GetFeeRate 获取费率
func (c *DogeClient) GetFeeRate() (*FeeRateResponse, error) {
	var result FeeRateResponse
	_, err := c.client.R().
		SetResult(&result).
		Get(c.baseURL)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTransaction 获取交易详情
func (c *DogeClient) GetTransaction(txHash string) (*TransactionResponse, error) {
	var result TransactionResponse
	_, err := c.client.R().
		SetQueryParam("token", c.token).
		SetResult(&result).
		Get(fmt.Sprintf("%s/txs/%s", c.baseURL, txHash))

	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAccountBalance 获取账户余额
func (c *DogeClient) GetAccountBalance(address string) (*AccountResponse, error) {
	var result AccountResponse
	_, err := c.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("%s/addrs/%s/balance", c.baseURL, address))

	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetUnspentOutputs 获取UTXO
func (c *DogeClient) GetUnspentOutputs(address string) (*UnspentOutputResponse, error) {
	var result UnspentOutputResponse
	_, err := c.client.R().
		SetQueryParam("unspentOnly", "true").
		SetQueryParam("token", c.token).
		SetResult(&result).
		Get(fmt.Sprintf("%s/addrs/%s", c.baseURL, address))

	if err != nil {
		return nil, err
	}
	return &result, nil
}
