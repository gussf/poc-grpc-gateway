package main

import (
	"context"
	"flag"
	"log"

	"github.com/gussf/poc-grpc-gateway/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
	httpServerEndpoint = flag.String("http-server-endpoint", "localhost:8080", "http server endpoint")
	metricsEndpoint    = flag.String("metrics-endpoint", "localhost:9080", "metric scraping endpoint")
)

// I know it's ugly but I don't wanna do a DAL :(
func NewMongoDBClient(ctx context.Context) *mongo.Client {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:root@0.0.0.0:27017"))
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func main() {
	ctx := context.Background()

	client := NewMongoDBClient(ctx)
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	handler := handlers.NewEchoerHandler(client)
	startAPI(ctx, handler)
}
