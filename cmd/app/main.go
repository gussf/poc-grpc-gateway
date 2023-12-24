package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/gussf/poc-grpc-gateway/handlers"
	v1 "github.com/gussf/poc-grpc-gateway/proto/gen/go"
	"github.com/oklog/run"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
	httpServerEndpoint = flag.String("http-server-endpoint", "localhost:8080", "http server endpoint")
	metricsEndpoint    = flag.String("metrics-endpoint", "localhost:9080", "metric scraping endpoint")
)

func newGRPCServer() *grpc.Server {
	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(),
	)

	v1.RegisterEchoerServer(grpcSrv, handlers.NewEchoerHandler())
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

func main() {
	ctx := context.Background()

	// Serve gRPC server
	grpcServer := newGRPCServer()
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
		metricsSrv.GET("/metrics", func(ctx *gin.Context) {
			ctx.Writer.WriteHeader(200)
		})
		return metricsSrv.Run(*metricsEndpoint)
	}, func(error) {
	})

	g.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))

	if err := g.Run(); err != nil {
		log.Fatal(err)
	}

}
