package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/gussf/poc-grpc-gateway/handlers"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func startAPI(ctx context.Context, handler handlers.Echoer) {
	// Serve gRPC server
	grpcServer := newGRPCServer(handler)
	g := &run.Group{}
	g.Add(func() error {
		l, err := net.Listen("tcp", *grpcServerEndpoint)
		if err != nil {
			return err
		}
		log.Println("starting grpc on", *grpcServerEndpoint)
		return grpcServer.Serve(l)
	}, func(err error) {
		if err != nil {
			log.Fatal(err)
		}
		grpcServer.GracefulStop()
		grpcServer.Stop()
	})

	// Serve gRPC Gateway
	grpcGateway := &http.Server{Addr: *httpServerEndpoint}
	g.Add(func() error {
		mux := runtime.NewServeMux()
		log.Println("starting grpc-gateway on", *httpServerEndpoint)
		registerGRPCGateway(ctx, mux)
		return grpcGateway.ListenAndServe()
	}, func(error) {
		if err := grpcGateway.Close(); err != nil {
			log.Fatal(err)
		}
	})

	// Serve Metrics
	metricsSrv := gin.New()
	g.Add(func() error {
		log.Println("starting metrics on", *metricsEndpoint)
		metricsSrv.GET("/metrics", gin.WrapH(promhttp.Handler()))
		return metricsSrv.Run(*metricsEndpoint)
	}, func(error) {
	})

	g.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))

	if err := g.Run(); err != nil {
		log.Fatal(err)
	}

}
