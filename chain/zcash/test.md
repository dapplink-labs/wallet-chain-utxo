# Zcash rpc api test

## 1.getBlockHeaderByHash
- request
```
grpcurl -plaintext -d '{
  "hash": "00000000014db39226a138002094424894f141ca51388a251fbd970ac712de57",
  "chain": "Zcash"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.getBlockHeaderByHash
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get block header success",
  "parent_hash": "0000000000cf566a30cc144ecc86d91d2dde1367a6610f931f61e02160af388b",
  "block_hash": "00000000014db39226a138002094424894f141ca51388a251fbd970ac712de57",
  "merkle_root": "8deb1a26866788934532c9176fcecc528e3a90db8dd809d136de776ac06cd1e4",
  "number": "2836533"
}
```

## 2. convertAddress
- request
```
grpcurl -plaintext -d '{
  "chain": "Zcash",
  "publicKey": "036e418e6b13e19614d67e281d2635fff7fa5e5d6e10eb6ed03d59dd3fd570ad5c",
  "format": "p2pkh"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.convertAddress
```
- response
```
{
  "code": "SUCCESS",
  "msg": "create address success",
  "address": "t1QzsGFr2iNTxNGAmhw2Nv8P85BG9XU31JJ"
}
```

## 3. validAddress
- request
```
grpcurl -plaintext -d '{
  "chain": "Zcash",
  "address": "t1QzsGFr2iNTxNGAmhw2Nv8P85BG9XU31JJ"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.validAddress
```
- response
```
{
  "code": "SUCCESS",
  "msg": "verify address success",
  "valid": true
}
```

## 4. getFee
- request
```
grpcurl -plaintext -d '{
  "chain": "Zcash"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.getFee
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get fee success",
  "best_fee": "",
  "best_fee_sat": "",
  "slow_fee": "0.00000001",
  "normal_fee": "0.00000002",
  "fast_fee": "0.00000003"
}
```

## 5. getAccount
- request
```
grpcurl -plaintext -d '{
  "chain": "Zcash",
  "address": "t1Qfn52D2TGCAC75aBEKiWkozeJgaQsR638"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.getAccount

```
- response
```
{
  "code": "SUCCESS",
  "msg": "get Zec balance success",
  "network": "",
  "balance": "126.91628771"
}
```

## 6. getUnspentOutput
- request
```
grpcurl -plaintext -d '{
  "chain": "Zcash",
  "address": "t1Qfn52D2TGCAC75aBEKiWkozeJgaQsR638"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.getUnspentOutputs

```
- response
```
{
  "code": "SUCCESS",
  "msg": "get unspent outputs success",
  "unspent_outputs": [
    {
      "tx_id": "9c3f290a78bf8750127c21135c69ccce06a74a72ce77bf028b39ddfcdf58d976",
      "tx_hash_big_endian": "",
      "tx_output_n": "0",
      "script": "",
      "height": "",
      "block_time": "",
      "address": "t1Qfn52D2TGCAC75aBEKiWkozeJgaQsR638",
      "unspent_amount": "0.00001000",
      "value_hex": "",
      "confirmations": "0",
      "index": "0"
    },
    ...
    ...
    ...
  ]
}
```

## 7. getBlockByNumber
- request
```
grpcurl -plaintext -d '{
  "height": "12801",
  "chain": "Zcash"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.getBlockByNumber

```
- response
```
{
  "code": "SUCCESS",
  "msg": "get block by number succcess",
  "height": "12801",
  "hash": "0000000078548f474d22ebdf90081935f32f8feafdf0f47de55acbc45504271f",
  "tx_list": [
    {
      "hash": "80ef9f1a51a4d8ba64361f67c7b51212b8c409fa09326d7f30c099d79132151e",
      ...
      ...
      ...
    },
    ...
    ...
    ...
  ]
} 
```
## 8. getBlockByHash
- request
```
grpcurl -plaintext -d '{
  "chain": "Zcash",
  "hash": "0000000078548f474d22ebdf90081935f32f8feafdf0f47de55acbc45504271f"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.getBlockByHash
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get block by hash succcess",
  "height": "12801",
  "hash": "0000000078548f474d22ebdf90081935f32f8feafdf0f47de55acbc45504271f",
  "tx_list": [
    {
      "hash": "80ef9f1a51a4d8ba64361f67c7b51212b8c409fa09326d7f30c099d79132151e",
      ...
      ...
      ...
    },
    ...
    ...
    ...
  ]
} 
```

## 9. getBlockHeaderByHash
- request
```
grpcurl -plaintext -d '{
  "hash": "0000000078548f474d22ebdf90081935f32f8feafdf0f47de55acbc45504271f",
  "chain": "Zcash"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.getBlockHeaderByHash
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get block header success",
  "parent_hash": "00000000febc373a1da2bd9f887b105ad79ddc26ac26c2b28652d64e5207c5b5",
  "block_hash": "0000000078548f474d22ebdf90081935f32f8feafdf0f47de55acbc45504271f",
  "merkle_root": "192c7a278ee557c8e3ab7d4f5163bd144e1744574570b0a198ac77df2ba5478b",
  "number": "12801"
}
```

## 10. getTransactionByHash
- request
```
grpcurl -plaintext -d '{
  "chain": "Zcash",
  "hash": "d80bf2f2d84071ea8ffef59d16d6f30999a98e468ac5ccb322f7195d1700758a"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.getTxByHash
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get transaction success",
  "tx": {
    "hash": "d80bf2f2d84071ea8ffef59d16d6f30999a98e468ac5ccb322f7195d1700758a",
    "index": 0,
    "froms": [
      {
        "address": "t1Qfn52D2TGCAC75aBEKiWkozeJgaQsR638"
      },
      {
        "address": "t1Qfn52D2TGCAC75aBEKiWkozeJgaQsR638"
      },
      {
        "address": "t1Qfn52D2TGCAC75aBEKiWkozeJgaQsR638"
      },
      {
        "address": "t1Qfn52D2TGCAC75aBEKiWkozeJgaQsR638"
      }
    ],
    "tos": [
      {
        "address": "t1Rz9ytXcKyaPdbM9AwBcUYJePDWaterEUr"
      },
      {
        "address": "t1Qfn52D2TGCAC75aBEKiWkozeJgaQsR638"
      }
    ],
    "fee": "0.00030840",
    "status": "Success",
    "values": [
      {
        "value": "1.73277924"
      },
      {
        "value": "7.02153296"
      }
    ],
    "type": 0,
    "height": "2837800",
    "brc20_address": "",
    "datetime": "1740758065"
  }
}
```

## 11. createUnsignTransaction
- request
```
grpcurl -plaintext -d '{
  "vin": [
    {
      "hash": "5b0191406c2ffa0fc7fe5af68c53c41ab5aa6012a1b619b3c8376911957883da",
      "amount": "4258340",
      "address": "t1QzsGFr2iNTxNGAmhw2Nv8P85BG9XU31JJ"
    }
  ],
  "vout": [
    {
      "address": "t1JM4tZ6anzxXs5bJxGFkA7wYfUkdjoQkCp",
      "amount": "3834556"
    }
  ],
  "chain": "Zcash"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.createUnSignTransaction
```
- response
```
{
  "code": "SUCCESS",
  "msg": "create un sign transaction success",
  "tx_data": "AQAAAAHag3iVEWk3yLMZtqESYKq1GsRTjPZa/scP+i9sQJEBWwAAAAAA/////wG8gjoAAAAAABp2qRW4BTGr/UFyeiFwvRmfNWf9YfSLzWmIrAAAAAA=",
  "sign_hashes": [
    "E0PBA4aH9KnFyfQLQKfG58AjYd9OY4DiWTW3oZp6qXE="
  ]
}
```

## 12. SendTx
- request
```
grpcurl -plaintext -d '{
  "chain": "Zcash",
  "rawTx": "0400008085202f8901d32be6ff4fcd64b7c8def5154000f7ae878f63e485cea24fe8300399f52d0ccf000000006b483045022100aca8dcf8bf7646409830ff60d5c6dce33af9b7e36af8651552e6d3324d0beb7002204fed838250e402adf8a86733572351c2e7b541677f0a66e97ee95742fdb0f6de0121036e418e6b13e19614d67e281d2635fff7fa5e5d6e10eb6ed03d59dd3fd570ad5cffffffff01b2f37000000000001976a914924e352f887f6fdf46d3153fbc5124b5ad3f902988ac00000000000000000000000000000000000000"
}' 127.0.0.1:8389 dapplink.utxo.WalletUtxoService.SendTx
```
- response
```
{
  "code": "SUCCESS",
  "msg": "send tx success",
  "tx_hash": "4f389d3dede396cecf309f232411f4bd426d2383f62a7593b35d726764248625"
}
```



