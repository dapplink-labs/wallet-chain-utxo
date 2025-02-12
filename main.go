package main

import (
	"flag"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ethereum/go-ethereum/log"

	"github.com/dapplink-labs/wallet-chain-utxo/chaindispatcher"
	"github.com/dapplink-labs/wallet-chain-utxo/config"
	"github.com/dapplink-labs/wallet-chain-utxo/rpc/utxo"
)

func main() {
	var f = flag.String("c", "config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	if err != nil {
		panic(err)
	}
	dispatcher, err := chaindispatcher.New(conf)
	if err != nil {
		log.Error("Setup dispatcher failed", "err", err)
		panic(err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(dispatcher.Interceptor),
	)
	defer grpcServer.GracefulStop()

	utxo.RegisterWalletUtxoServiceServer(grpcServer, dispatcher)

	listen, err := net.Listen("tcp", ":"+conf.Server.Port)
	if err != nil {
		log.Error("Failed to listen", "error", err)
		panic(err)
	}
	reflection.Register(grpcServer)

	log.Info("Starting gRPC server", "port", conf.Server.Port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Error("Failed to serve", "error", err)
		panic(err)
	}
}
