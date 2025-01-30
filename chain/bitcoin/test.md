# Bitcoin rpc api test

## 1.get fee
- request
```
grpcurl -plaintext -d '{
  "chain": "Bitcoin",
  "coin": "Btc",
  "network": "mainnet"
}' 127.0.0.1:8289 dapplink.utxo.WalletUtxoService.getFee
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get fee success",
  "best_fee": "0.00001427",
  "best_fee_sat": "1427",
  "slow_fee": "",
  "normal_fee": "",
  "fast_fee": ""
}
```

## 2.get account
- request
```
grpcurl -plaintext -d '{
  "chain": "Bitcoin",
  "network": "mainet",
  "address": "bc1qa2eu6p5rl9255e3xz7fcgm6snn4wl5kdfh7zpt05qp5fad9dmsys0qjg0e"
}' 127.0.0.1:8289 dapplink.utxo.WalletUtxoService.getAccount
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get btc balance success",
  "network": "",
  "balance": "5171031881312"
}
```


## 3.get tx by address
- request
```
grpcurl -plaintext -d '{
  "chain": "Bitcoin",
  "coin": "BTC",
  "network": "mainnet",
  "address": "bc1ql49ydapnjafl5t2cp9zqpjwe6pdgmxy98859v2",
  "page": 1,
  "pagesize": 10
}' 127.0.0.1:8289 dapplink.utxo.WalletUtxoService.getTxByAddress
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get transaction list success",
  "tx": [
    {
      "hash": "110b7ee111ab849682617e7df0a5df137ae3425e434399a9a25adc98de4f3cf5",
      "index": 0,
      "froms": [
        {
          "address": "bc1ql49ydapnjafl5t2cp9zqpjwe6pdgmxy98859v2"
        }
      ],
      "tos": [
        {
          "address": "bc1qhk0ghcywv0mlmcmz408sdaxudxuk9tvng9xx8g"
        },
        {
          "address": "bc1ql49ydapnjafl5t2cp9zqpjwe6pdgmxy98859v2"
        }
      ],
      "fee": "2030",
      "status": "Success",
      "values": [
        {
          "value": "70000000000"
        },
        {
          "value": "37499996955"
        }
      ],
      "type": 0,
      "height": "868100",
      "brc20_address": "",
      "datetime": "1730299612"
    }
  ]
}
```

## 4.get tx by hash
- request
```
grpcurl -plaintext -d '{
  "chain": "Dogecoin",
  "coin": "doge",
  "network": "mainnet",
  "hash": "9b1810bcb18af9916f7ddc2588a9d262fad19cece21936c7a742ae4849ce21b9"
}' 127.0.0.1:8289 dapplink.utxo.WalletUtxoService.getTxByHash
```
- response
```{
  "code": "SUCCESS",
  "msg": "get transaction success",
  "tx": {
    "hash": "9b1810bcb18af9916f7ddc2588a9d262fad19cece21936c7a742ae4849ce21b9",
    "index": 0,
    "froms": [
      {
        "address": "DSWgTjzByQFoGQtGEGq9yzNPr7StEGX7tV"
      }
    ],
    "tos": [
      {
        "address": "DLAGch3DRpf3cDLAXXqDRUffUBvUdCuW9M"
      },
      {
        "address": "DNRkXPjdWfqiy4betAY97y2t9fDMZYdtPb"
      }
    ],
    "fee": "600000000",
    "status": "Success",
    "values": [
      {
        "value": "1437600000000"
      },
      {
        "value": "4457125016418361"
      }
    ],
    "type": 0,
    "height": "5513540",
    "brc20_address": "",
    "datetime": "2024-12-22T15:43:21Z"
  }
}
```

## 5.get address utxo
- request
```
grpcurl -plaintext -d '{
  "chain": "Bitcoin",
  "network": "mainnet",
  "address": "bc1ql49ydapnjafl5t2cp9zqpjwe6pdgmxy98859v2"
}' 127.0.0.1:8289 dapplink.utxo.WalletUtxoService.getUnspentOutputs
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get unspent outputs success",
  "unspent_outputs": [
    {
      "tx_id": "8b57ffaec69a2436ccf5057f92bb9bb6717939eeb5af181081b03f3f60d62061",
      "tx_hash_big_endian": "6120d6603f3fb0811018afb5ee397971b69bbb927f05f5cc36249ac6aeff578b",
      "tx_output_n": "1",
      "script": "0014fd4a46f4339753fa2d58094400c9d9d05a8d9885",
      "height": "",
      "block_time": "",
      "address": "",
      "unspent_amount": "109999998139",
      "value_hex": "",
      "confirmations": "0",
      "index": "4902722177913668"
    },
    {
      "tx_id": "f53c4fde98dc5aa2a99943435e42e37a13dfa5f07d7e61829684ab11e17e0b11",
      "tx_hash_big_endian": "110b7ee111ab849682617e7df0a5df137ae3425e434399a9a25adc98de4f3cf5",
      "tx_output_n": "1",
      "script": "0014fd4a46f4339753fa2d58094400c9d9d05a8d9885",
      "height": "",
      "block_time": "",
      "address": "",
      "unspent_amount": "37499996955",
      "value_hex": "",
      "confirmations": "0",
      "index": "8628460378594187"
    }
  ]
}
```
