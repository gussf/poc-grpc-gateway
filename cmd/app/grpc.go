package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/gussf/poc-grpc-gateway/handlers"
	middlewares "github.com/gussf/poc-grpc-gateway/middlewares/grpc"
	v1 "github.com/gussf/poc-grpc-gateway/proto/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newGRPCServer(handler handlers.Echoer) *grpc.Server {
	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middlewares.HttpMetricsInterceptor(),
		),
	)

	v1.RegisterEchoerServer(grpcSrv, handler)
	return grpcSrv
}

func registerGRPCGateway(ctx context.Context, mux *runtime.ServeMux) {
	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := v1.RegisterEchoerHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		log.Fatal(err)
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	http.ListenAndServe(*httpServerEndpoint, mux)
}
