package config

import (
	"os"

	"gopkg.in/yaml.v2"

	"github.com/ethereum/go-ethereum/log"
)

type Server struct {
	GrpcPort string `yaml:"grpc_port"`
	WebPort  string `yaml:"web_port"`
}

type RPC struct {
	RpcUrl  string `yaml:"rpc_url"`
	RpcUser string `yaml:"rpc_user"`
	RpcPass string `yaml:"rpc_pass"`
}

type Node struct {
	RpcUrl       string `yaml:"rpc_url"`
	RpcUser      string `yaml:"rpc_user"`
	RpcPass      string `yaml:"rpc_pass"`
	DataApiUrl   string `yaml:"data_api_url"`
	DataApiKey   string `yaml:"data_api_key"`
	DataApiToken string `yaml:"data_api_token"`
	TpApiUrl     string `yaml:"tp_api_url"`
	TimeOut      uint64 `yaml:"time_out"`
}

type WalletNode struct {
	Btc  Node `yaml:"btc"`
	Bch  Node `yaml:"bch"`
	Dash Node `yaml:"dash"`
	Doge Node `yaml:"doge"`
	Ltc  Node `yaml:"ltc"`
}

type Config struct {
	Server     Server     `yaml:"server"`
	WalletNode WalletNode `yaml:"walletnode"`
	NetWork    string     `yaml:"network"`
	Chains     []string   `yaml:"chains"`
}

type NetWorkType int

const (
	MainNet NetWorkType = iota
	TestNet
	RegTest
)

func New(path string) (*Config, error) {
	var config = new(Config)
	h := log.NewTerminalHandler(os.Stdout, true)
	log.SetDefault(log.NewLogger(h))

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

const UnsupportedChain = "Unsupport chain"
const UnsupportedOperation = UnsupportedChain
