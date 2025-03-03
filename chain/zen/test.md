# zen rpc api test



## 1. - getAccount
- request
```
grpcurl -plaintext -d '{
  "chain": "Zen",
  "network": "mainnet",
  "address": "znVmcBG35teueHJTuB1dQuc94XY733Z1Hec"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.getAccount
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get Zen balance success",
  "network": "mainnet",
  "balance": "970000"
}
```

## 2.- getFee
- request
```
grpcurl -plaintext -d '{
  "chain": "Zen",
  "coin": "Zen",
  "network": "mainnet"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.getFee
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get fee success",
  "best_fee": "0.0001",
  "best_fee_sat": "10000",
  "slow_fee": "",
  "normal_fee": "",
  "fast_fee": ""
}
```


## 3.- getTxByAddress
- request
```
grpcurl -plaintext -d '{
  "chain": "Zen",
  "coin": "Zen",
  "network": "mainnet",
  "address": "znVmcBG35teueHJTuB1dQuc94XY733Z1Hec",
  "page": 1,
  "pagesize": 10
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.getTxByAddress
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get transaction list success",
  "tx": [
    {
      "hash": "3ab4283b9d90aaee261fbe7065c88d7f20b2dae349c8039dbe55f55ee970df34",
      "index": 0,
      "froms": [
        {
          "address": "znp5aYNr9f848b4pFH3pbX7wtRWrTaqb4cg"
        },
        {
          "address": "znp5aYNr9f848b4pFH3pbX7wtRWrTaqb4cg"
        }
      ],
      "tos": [
        {
          "address": "znVmcBG35teueHJTuB1dQuc94XY733Z1Hec"
        }
      ],
      "fee": "0.0001",
      "status": "Success",
      "values": [
        {
          "value": "0.00970000"
        }
      ],
      "type": 1,
      "height": "1721899",
      "brc20_address": "",
      "datetime": "1740362359"
    },
    {
      "hash": "20fb67ff4f88c191505392faf11541653a2689e5273e4235c448f810cbe81488",
      "index": 0,
      "froms": [
        {
          "address": "znVmcBG35teueHJTuB1dQuc94XY733Z1Hec"
        }
      ],
      "tos": [
        {
          "address": "znp5aYNr9f848b4pFH3pbX7wtRWrTaqb4cg"
        }
      ],
      "fee": "0.0001",
      "status": "Success",
      "values": [
        {
          "value": "0.00890000"
        }
      ],
      "type": 0,
      "height": "1721896",
      "brc20_address": "",
      "datetime": "1740361825"
    },
    {
      "hash": "40e2468bd53c9afef28f66367d89fc739269c3aa57fba7880c2428ee11e1fafe",
      "index": 0,
      "froms": [
        {
          "address": "znp5aYNr9f848b4pFH3pbX7wtRWrTaqb4cg"
        },
        {
          "address": "znp5aYNr9f848b4pFH3pbX7wtRWrTaqb4cg"
        }
      ],
      "tos": [
        {
          "address": "znVmcBG35teueHJTuB1dQuc94XY733Z1Hec"
        },
        {
          "address": "znp5aYNr9f848b4pFH3pbX7wtRWrTaqb4cg"
        }
      ],
      "fee": "0.0001",
      "status": "Success",
      "values": [
        {
          "value": "0.00900000"
        },
        {
          "value": "0.00090000"
        }
      ],
      "type": 1,
      "height": "1721651",
      "brc20_address": "",
      "datetime": "1740325316"
    }
  ]
}
```

## 4.- getTxByHash
- request
```
grpcurl -plaintext -d '{
  "chain": "Zen",
  "coin": "Zen",
  "network": "mainnet",
  "hash": "3ab4283b9d90aaee261fbe7065c88d7f20b2dae349c8039dbe55f55ee970df34"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.getTxByHash
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get transaction success",
  "tx": {
    "hash": "3ab4283b9d90aaee261fbe7065c88d7f20b2dae349c8039dbe55f55ee970df34",
    "index": 0,
    "froms": [
      {
        "address": "znp5aYNr9f848b4pFH3pbX7wtRWrTaqb4cg"
      },
      {
        "address": "znp5aYNr9f848b4pFH3pbX7wtRWrTaqb4cg"
      }
    ],
    "tos": [
      {
        "address": "znVmcBG35teueHJTuB1dQuc94XY733Z1Hec"
      }
    ],
    "fee": "0.0001",
    "status": "Success",
    "values": [
      {
        "value": "0.00970000"
      }
    ],
    "type": 0,
    "height": "1721899",
    "brc20_address": "",
    "datetime": "1740362359"
  }
}
```

## 5.- SendTx
- request
```
grpcurl -plaintext -d '{
  "chain": "Zen",
  "coin": "Zen",
  "network": "mainnet",
  "rawTx": "0100000001a9b2524c650acf2e82ffa4b9ab38ee968a9d9c1e02ca7ec710d1257950102f22000000006a473044022048ba2f3c136b67ede9c17590698a89e4a09a54f2c1ba74a00ea37ca44776eb5102204ec00f3bb1638bfb72cb307c365c44e4ea93c73056dc810db14590a82b10aff0012103d95d0cbe7e72b6fc79e64394e0356c6ab8c662c71f06bbe35ff895bb03bdbaa2ffffffff01e0570e00000000003f76a9143362459b23db8f89d61e4224e2ba1b28cd10b88c88ac20f06a300e50d0b29b4218ec9018302362294b98734a3edf0da88f0c01000000000321511ab400000000"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.SendTx
```

- response
```
{
  "code": "SUCCESS",
  "msg": "send tx success ",
  "tx_hash": "42bdaf444c0ca92b01b13a4c810768904533473f02de17991870677068230037"
}
```

## 6. - convertAddress
- request
```
grpcurl -plaintext -d '{
  "chain": "Zen",
  "network": "mainnet",
  "publicKey": "03d95d0cbe7e72b6fc79e64394e0356c6ab8c662c71f06bbe35ff895bb03bdbaa2"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.convertAddress
```

- response
```
{
  "code": "SUCCESS",
  "msg": "create address success",
  "address": "znp5aYNr9f848b4pFH3pbX7wtRWrTaqb4cg"
}
```

## 7. - validAddress
- request
```
grpcurl -plaintext -d '{
  "chain": "Zen",
  "network": "mainnet",
  "address": "znp5aYNr9f848b4pFH3pbX7wtRWrTaqb4cg"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.validAddress
```

- response
```
{
  "code": "SUCCESS",
  "msg": "verify address success, It is p2pkh address",
  "valid": true
}
```

## 8. - CreateUnSignTransaction
- request
```
grpcurl -plaintext -d '{
  "vin": [
    {
      "hash": "42bdaf444c0ca92b01b13a4c810768904533473f02de17991870677068230037",
      "amount": "940000",
      "address": "znVmcBG35teueHJTuB1dQuc94XY733Z1Hec"
    }
  ],
  "vout": [
    {
      "address": "znp5aYNr9f848b4pFH3pbX7wtRWrTaqb4cg",
      "amount": "930000"
    }
  ],
  "chain": "Zen"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.createUnSignTransaction
```

- response
```
{
  "code": "SUCCESS",
  "msg": "create un sign transaction success",
  "tx_data": "QX8DAQEIVFhPYmplY3QB/4AAAQQBCExvY2tUaW1lAQYAAQdWZXJzaW9uAQYAAQNJbnMB/4YAAQRPdXRzAf+KAAAAGv+FAgEBC1tddHlwZXMuSW5zAf+GAAH/ggAASv+BAwEBA0lucwH/ggABBAEGT3V0cHV0Af+EAAEGU2NyaXB0AQwAARBQcmV2U2NyaXB0UHViS2V5AQwAAQhTZXF1ZW5jZQEMAAAAJv+DAwEBBk91dHB1dAH/hAABAgEESGFzaAEMAAEEVm91dAEGAAAAG/+JAgEBDFtddHlwZXMuT3V0cwH/igAB/4gAACr/hwMBAQRPdXRzAf+IAAECAQhTYXRvc2hpcwEGAAEGU2NyaXB0AQwAAAD+AV7/gAIBAQEBAUA0MmJkYWY0NDRjMGNhOTJiMDFiMTNhNGM4MTA3Njg5MDQ1MzM0NzNmMDJkZTE3OTkxODcwNjc3MDY4MjMwMDM3AAJ+NzZhOTE0MzM2MjQ1OWIyM2RiOGY4OWQ2MWU0MjI0ZTJiYTFiMjhjZDEwYjg4Yzg4YWMyMGYwNmEzMDBlNTBkMGIyOWI0MjE4ZWM5MDE4MzAyMzYyMjk0Yjk4NzM0YTNlZGYwZGE4OGYwYzAxMDAwMDAwMDAwMzIxNTExYWI0AQhmZmZmZmZmZgABAQH9DjDQAX43NmE5MTRmYzNhYTg2MTJmZWIyOGZlYjAzZWZhNTA2NWQ1NzZkYmI2N2M5ZjhjODhhYzIwYmFiOTIyOThlY2I2YjIwOGVkYTU0NjkxZTAyZjE0NjUwZTk0ZjdmNGEzOTU4YTc3YTY5OTM2MDIwMDAwMDAwMDAzYWE1NjFhYjQAAA==",
  "sign_hashes": [
    "rvH6w3RvzDvsrS8elT34Lzypy30FiPjvvoRzR1NFebw="
  ]
}
```

## 9. buildSignedTransaction
- request
```
grpcurl -plaintext -d '{
  "signatures": [
    "MEQCIGUYSPcElDGXF7UiPciFmwF70TPtUCpi+lDSsvOkHZbLAiAZ69i0jkqqerR/yPQESQ38sxZrqzEjUc3D1hk7bAGrogE="
  ],
  "publicKeys": [
    "AifHfY2zwfnRpBL+ZRVRo3n6OyWjflGir816ly7eeUSi"
  ],
  "chain": "Zen",
  "txData": "QX8DAQEIVFhPYmplY3QB/4AAAQQBCExvY2tUaW1lAQYAAQdWZXJzaW9uAQYAAQNJbnMB/4YAAQRPdXRzAf+KAAAAGv+FAgEBC1tddHlwZXMuSW5zAf+GAAH/ggAASv+BAwEBA0lucwH/ggABBAEGT3V0cHV0Af+EAAEGU2NyaXB0AQwAARBQcmV2U2NyaXB0UHViS2V5AQwAAQhTZXF1ZW5jZQEMAAAAJv+DAwEBBk91dHB1dAH/hAABAgEESGFzaAEMAAEEVm91dAEGAAAAG/+JAgEBDFtddHlwZXMuT3V0cwH/igAB/4gAACr/hwMBAQRPdXRzAf+IAAECAQhTYXRvc2hpcwEGAAEGU2NyaXB0AQwAAAD+AV7/gAIBAQEBAUA0MmJkYWY0NDRjMGNhOTJiMDFiMTNhNGM4MTA3Njg5MDQ1MzM0NzNmMDJkZTE3OTkxODcwNjc3MDY4MjMwMDM3AAJ+NzZhOTE0MzM2MjQ1OWIyM2RiOGY4OWQ2MWU0MjI0ZTJiYTFiMjhjZDEwYjg4Yzg4YWMyMGYwNmEzMDBlNTBkMGIyOWI0MjE4ZWM5MDE4MzAyMzYyMjk0Yjk4NzM0YTNlZGYwZGE4OGYwYzAxMDAwMDAwMDAwMzIxNTExYWI0AQhmZmZmZmZmZgABAQH9DjDQAX43NmE5MTRmYzNhYTg2MTJmZWIyOGZlYjAzZWZhNTA2NWQ1NzZkYmI2N2M5ZjhjODhhYzIwYmFiOTIyOThlY2I2YjIwOGVkYTU0NjkxZTAyZjE0NjUwZTk0ZjdmNGEzOTU4YTc3YTY5OTM2MDIwMDAwMDAwMDAzYWE1NjFhYjQAAA=="
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.buildSignedTransaction
```

- response
```
{
  "code": "SUCCESS",
  "msg": "",
  "signed_tx_data": "AQAAAAE3ACNocGdwGJkX3gI/RzNFkGgHgUw6sQErqQxMRK+9QgAAAABqRzBEAiBlGEj3BJQxlxe1Ij3IhZsBe9Ez7VAqYvpQ0rLzpB2WywIgGevYtI5Kqnq0f8j0BEkN/LMWa6sxI1HNw9YZO2wBq6IBIQInx32Ns8H50aQS/mUVUaN5+jslo35Roq/Nepcu3nlEov////8B0DAOAAAAAAA/dqkU/DqoYS/rKP6wPvpQZdV227Z8n4yIrCC6uSKY7LayCO2lRpHgLxRlDpT39KOVinemmTYCAAAAAAOqVhq0AAAAAA==",
  "hash": "B0L3aUp3TKOjVStklz7IcEe4a5usXQeuFflRFlvqg4o="
}
```



## 10. - GetBlockByHash
- request
```
grpcurl -plaintext -d '{
  "chain": "Zen",
  "hash": "0000000003469fb58f05f2723405d5eca4003a485430ea9bfe22ac38f40e7b7c"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.getBlockByHash
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get block by hash succcess",
  "height": "1724946",
  "hash": "0000000003469fb58f05f2723405d5eca4003a485430ea9bfe22ac38f40e7b7c",
  "tx_list": [
    {
      "hash": "51b16705a65966fbb2a3421f05e0b7628685c907b498697564d00607363ab2d5",
      "fee": "",
      "vin": [],
      "vout": []
    }
  ]
}
```

## 11. - GetBlockByNumber
- request
```
grpcurl -plaintext -d '{
  "chain": "Zen"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.getBlockByNumber
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get block by number succcess",
  "height": "1724982",
  "hash": "0000000000712b6c5fe6ce4bfe7110de60bc4e63bb6404f75e3fd7016d3bd158",
  "tx_list": [
    {
      "hash": "8c308f0804c5f3bc5c1da171979d7d642ac8560fb137830441288aee01098de4",
      "fee": "",
      "vin": [],
      "vout": []
    }
  ]
}
```

- request
```
grpcurl -plaintext -d '{
  "chain": "Zen",
  "height": "1724946"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.getBlockByNumber
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get block by number succcess",
  "height": "1724946",
  "hash": "0000000003469fb58f05f2723405d5eca4003a485430ea9bfe22ac38f40e7b7c",
  "tx_list": [
    {
      "hash": "51b16705a65966fbb2a3421f05e0b7628685c907b498697564d00607363ab2d5",
      "fee": "",
      "vin": [],
      "vout": []
    }
  ]
}
```


## 12.get address utxo
- request
```
grpcurl -plaintext -d '{
  "chain": "Zen",
  "network": "mainnet",
  "address": "znVmcBG35teueHJTuB1dQuc94XY733Z1Hec"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.getUnspentOutputs
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get unspent outputs success",
  "unspent_outputs": [
    {
      "tx_id": "3ab4283b9d90aaee261fbe7065c88d7f20b2dae349c8039dbe55f55ee970df34",
      "tx_hash_big_endian": "",
      "tx_output_n": "0",
      "script": "76a9143362459b23db8f89d61e4224e2ba1b28cd10b88c88ac20abca01b259dcb3bf17c4e8a4b15602867c1ebe9fe0e90933990b58000000000003fc441ab4",
      "height": "1721899",
      "block_time": "",
      "address": "znVmcBG35teueHJTuB1dQuc94XY733Z1Hec",
      "unspent_amount": "970000",
      "value_hex": "",
      "confirmations": "2044",
      "index": "0"
    }
  ]
}
```
