package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/fullstorydev/grpcui/standalone"
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

	// Load configuration
	conf, err := config.New(*f)
	if err != nil {
		log.Error("Failed to load configuration", "error", err)
		panic(err)
	}

	// Create gRPC dispatcher
	dispatcher, err := chaindispatcher.New(conf)
	if err != nil {
		log.Error("Setup dispatcher failed", "err", err)
		panic(err)
	}

	// Initialize gRPC server
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(dispatcher.Interceptor))
	defer grpcServer.GracefulStop()

	// Register services
	utxo.RegisterWalletUtxoServiceServer(grpcServer, dispatcher)
	reflection.Register(grpcServer)

	// Start gRPC server
	grpcAddr := fmt.Sprintf(":%s", conf.Server.GrpcPort)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Error("Failed to create gRPC listener", "error", err)
		panic(err)
	}

	// Run gRPC server in background
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("gRPC server failed", "error", err)
			panic(err)
		}
	}()

	// Wait for gRPC server to start
	time.Sleep(time.Second)

	// Create gRPC connection
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("localhost%s", grpcAddr),
		grpc.WithInsecure(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		log.Error("Failed to create gRPC connection", "error", err)
		panic(err)
	}
	defer conn.Close()

	// Create gRPC Web UI handler
	handler, err := standalone.HandlerViaReflection(context.Background(), conn, "")
	if err != nil {
		log.Error("Failed to create Web UI handler", "error", err)
		panic(err)
	}

	// Start Web UI server
	webAddr := fmt.Sprintf(":%s", conf.Server.WebPort)
	log.Info("Service started",
		"grpc", fmt.Sprintf("localhost%s", grpcAddr),
		"web", fmt.Sprintf("http://localhost%s", webAddr))

	if err := http.ListenAndServe(webAddr, handler); err != nil {
		log.Error("Failed to start Web UI server", "error", err)
		panic(err)
	}
}
