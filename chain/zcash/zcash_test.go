package zcash

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/dapplink-labs/wallet-chain-utxo/chain"
	"github.com/dapplink-labs/wallet-chain-utxo/config"
	"github.com/dapplink-labs/wallet-chain-utxo/rpc/utxo"
	"github.com/ethereum/go-ethereum/log"

	_ "github.com/stretchr/testify/assert"
)

func setup() (adaptor chain.IChainAdaptor, err error) {
	conf, err := config.New("../../config.yml")
	if err != nil {
		log.Error("load config failed, error:", err)
		return nil, err
	}
	adaptor, err = NewChainAdaptor(conf)
	if err != nil {
		log.Error("create chain adaptor failed, error:", err)
		return nil, err
	}
	return adaptor, nil
}

func TestChainAdaptor_GetSupportChains(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	adaptor.GetSupportChains(&utxo.SupportChainsRequest{
		Chain: ChainName,
	})
}

func TestChainAdaptor_ConvertAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := adaptor.ConvertAddress(&utxo.ConvertAddressRequest{
		Chain:     ChainName,
		Format:    "p2pkh",
		PublicKey: "036e418e6b13e19614d67e281d2635fff7fa5e5d6e10eb6ed03d59dd3fd570ad5c",
	})

	log.Info("address: ", rsp.Address)
	assert.Equal(t, "t1QzsGFr2iNTxNGAmhw2Nv8P85BG9XU31JJ", rsp.Address, "Expected this address")

	js, _ := json.Marshal(rsp)

	log.Info(string(js))

}

func TestChainAdaptor_VarifyAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := adaptor.ValidAddress(&utxo.ValidAddressRequest{
		Chain:   ChainName,
		Address: "t1QzsGFr2iNTxNGAmhw2Nv8P85BG9XU31JJ",
	})

	log.Info("address is valid?: ", rsp.Valid)
	assert.Equal(t, true, rsp.Valid, "Expected correct address")

	js, _ := json.Marshal(rsp)

	log.Info(string(js))

}
