package zcash

import (
	"github.com/go-resty/resty/v2"
)

type BaseClient struct {
	*resty.Client
	RequestUrl string
}

func NewBaseClient(RpcUrl, RpcUser, RpcPass string) (*BaseClient, error) {
	client := resty.New()
	header := map[string]string{
		"Content-Type": "application/json",
	}
	client.SetHeaders(header)
	return &BaseClient{
		Client:     client,
		RequestUrl: RpcUrl,
	}, nil
}

type BaseDataClient struct {
	*resty.Client
	RequestUrl string
}

func NewBaseDataClient(RpcUrl, ApiKey, RpcUser, RpcPass string) (*BaseDataClient, error) {
	client := resty.New()
	header := map[string]string{
		"accept":           "application/json, text/plain, */*",
		"accept-encoding":  "gzip, compress, deflate, br",
		"x-api-key":        ApiKey,
		"x-api-version":    "2024-12-12",
		"x-envoy-internal": "true",
	}
	client.SetHeaders(header)
	return &BaseDataClient{
		Client:     client,
		RequestUrl: RpcUrl,
	}, nil
}
