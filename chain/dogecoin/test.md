# Bitcoin rpc api test

## 1.get fee
- request
```
grpcurl -plaintext -d '{
  "chain": "Dogecoin",
  "coin": "doge",
  "network": "miannet"
}' 127.0.0.1:8289 dapplink.utxo.WalletUtxoService.getFee
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get fee success",
  "best_fee": "1.08205959",
  "best_fee_sat": "108205959",
  "slow_fee": "0.17836453",
  "normal_fee": "1.08205959",
  "fast_fee": "2.73674534"
}
```

## 2.get account
- request
```
grpcurl -plaintext -d '{
  "chain": "Dogecoin",
  "network": "mainnet",
  "address": "DLAGch3DRpf3cDLAXXqDRUffUBvUdCuW9M"
}' 127.0.0.1:8289 dapplink.utxo.WalletUtxoService.getAccount
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get account success",
  "network": "",
  "balance": "1444600000000"
}
```


## 3.get tx by address
- request
```

```

- response
```

```

## 4.get tx by hash
- request
```
grpcurl -plaintext -d '{
  "chain": "Bitcoin",
  "coin": "BTC",
  "network": "mainnet",
  "hash": "6120d6603f3fb0811018afb5ee397971b69bbb927f05f5cc36249ac6aeff578b"
}' 127.0.0.1:8289 dapplink.utxo.WalletUtxoService.getTxByHash
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get transaction success",
  "tx": {
    "hash": "6120d6603f3fb0811018afb5ee397971b69bbb927f05f5cc36249ac6aeff578b",
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
    "fee": "1861",
    "status": "Success",
    "values": [
      {
        "value": "70000000000"
      },
      {
        "value": "109999998139"
      }
    ],
    "type": 0,
    "height": "868100",
    "brc20_address": "",
    "datetime": "1730299686"
  }
}
```

## 5.get address utxo
- request
```
grpcurl -plaintext -d '{
  "chain": "Dogecoin",
  "network": "doge",
  "address": "DLAGch3DRpf3cDLAXXqDRUffUBvUdCuW9M"
}' 127.0.0.1:8289 dapplink.utxo.WalletUtxoService.getUnspentOutputs
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get unspent outputs success",
  "unspent_outputs": [
    {
      "tx_id": "9b1810bcb18af9916f7ddc2588a9d262fad19cece21936c7a742ae4849ce21b9",
      "tx_hash_big_endian": "",
      "tx_output_n": "0",
      "script": "",
      "height": "",
      "block_time": "",
      "address": "",
      "unspent_amount": "1437600000000",
      "value_hex": "",
      "confirmations": "3855",
      "index": "0"
    },
    {
      "tx_id": "b4e3622ea1f151b05ddd2736f59fa1420fd54e43625fbd8ee05b719c81b2c3ce",
      "tx_hash_big_endian": "",
      "tx_output_n": "0",
      "script": "",
      "height": "",
      "block_time": "",
      "address": "",
      "unspent_amount": "7000000000",
      "value_hex": "",
      "confirmations": "3866",
      "index": "0"
    }
  ]
}
```
