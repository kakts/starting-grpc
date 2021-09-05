package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"pancake.maker/gen/api"
	"pancake.maker/handler"
	"pancake.maker/interceptor"
)

func main() {
	port := 50051
	// tcpで50051でlistenする
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// NewServerの可変超引数でインターセプタを追加できる
	// 単項RPCの場合 		// grpc.UnaryInterceptor(myOriginalInterceptor),
	// これだと引数1つしか渡せないので、ChainUnaryServerを使って逐次的にインターセプタ実行できる
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			interceptor.FirstRequestInterceptor(),
			interceptor.SecondRequestInterceptor(),
		)),
	)

	// 重要
	// リフレクションを有効にしておく
	api.RegisterPancakeBakerServiceServer(server,
		handler.NewBakerHandler(),
	)
	reflection.Register(server)

	go func() {
		log.Printf("Start gRPC server port: %v", port)
		server.Serve(lis)
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	server.GracefulStop()

}
