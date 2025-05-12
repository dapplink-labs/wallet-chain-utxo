<!--
parent:
  order: false
-->

<div align="center">
  <h1> wallet-chain-utxo repo </h1>
</div>

<div align="center">
  <a href="https://github.com/dapplink-labs/wallet-chain-utxo/releases/latest">
    <img alt="Version" src="https://img.shields.io/github/tag/dapplink-labs/wallet-chain-utxo.svg" />
  </a>
  <a href="https://github.com/dapplink-labs/wallet-chain-utxo/blob/main/LICENSE">
    <img alt="License: Apache-2.0" src="https://img.shields.io/github/license/dapplink-labs/wallet-chain-utxo.svg" />
  </a>
  <a href="https://pkg.go.dev/github.com/dapplink-labs/wallet-chain-utxo">
    <img alt="GoDoc" src="https://godoc.org/github.com/dapplink-labs/wallet-chain-utxo?status.svg" />
  </a>
</div>

This repo is utxo chains rpc service gateway. currently support `Bitcoin`, `Bitcoincash`, `Dash`, `Dogecoin`, `Litecoin`, written in golang, provides grpc interface for upper-layer service access

**Tips**: need [Go 1.22+](https://golang.org/dl/)

## Install

### Install dependencies
```bash
go mod tidy
```
### build
```bash
go build or go install wallet-chain-utxo
```

### start
```bash
./wallet-chain-utxo -c ./config.yml
```

### Start the RPC interface test interface

```bash
grpcui -plaintext 127.0.0.1:8389
```

## Contribute

### 1.fork repo

fork wallet-chain-utxo to your github

### 2.clone repo

```bash
git@github.com:guoshijiang/wallet-chain-utxo.git
```

### 3. create new branch and commit code

```bash
git branch -C xxx
git checkout xxx

coding

git add .
git commit -m "xxx"
git push origin xxx
```

### 4.commit PR

Have a pr on your github and submit it to the wallet-chain-utxo repository

### 5.review

After the wallet-chain-utxo code maintainer has passed the review, the code will be merged into the wallet-chain-utxo library. At this point, your PR submission is complete

### 6.Disclaimer

This code has not yet been audited, and should not be used in any production systems.
