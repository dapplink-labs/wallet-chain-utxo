FROM golang:1.22.11-alpine3.21 as builder

RUN apk add --no-cache make ca-certificates gcc musl-dev linux-headers git jq bash

COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

WORKDIR /app

RUN go mod download

ARG CONFIG=config.yml

# build wallet-chain-utxo with the shared go.mod & go.sum files
COPY . /app/wallet-chain-utxo

WORKDIR /app/wallet-chain-utxo

RUN make

FROM alpine:3.18

COPY --from=builder /app/wallet-chain-utxo/wallet-chain-utxo /usr/local/bin
COPY --from=builder /app/wallet-chain-utxo/${CONFIG} /etc/wallet-chain-utxo/

WORKDIR /app

ENTRYPOINT ["wallet-chain-utxo"]
CMD ["-c", "/etc/wallet-chain-utxo/config.yml"]
