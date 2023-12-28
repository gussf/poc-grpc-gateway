package handlers

import (
	"context"
	"time"

	"github.com/gussf/poc-grpc-gateway/metrics"
	v1 "github.com/gussf/poc-grpc-gateway/proto/gen/go"
	"go.mongodb.org/mongo-driver/mongo"
)

type Echoer struct {
	mongo *mongo.Collection
	m     metrics.MongoDBMetrics
	v1.UnimplementedEchoerServer
}

func NewEchoerHandler(mongo *mongo.Client) Echoer {
	return Echoer{
		m:     metrics.NewMongoDBMetrics(),
		mongo: mongo.Database("echoer").Collection("words"),
	}
}

func (e Echoer) PostEcho(ctx context.Context, msg *v1.StringMessage) (*v1.StringMessage, error) {
	method := "PostEcho"
	start := time.Now()

	_, err := e.mongo.InsertOne(ctx, msg)
	e.m.Request.WithLabelValues(method).Observe(float64(time.Since(start).Seconds()))
	if err != nil {
		e.m.Fail.WithLabelValues(method).Add(1)
		return nil, err
	}

	e.m.Success.WithLabelValues(method).Add(1)
	return &v1.StringMessage{Value: "OK"}, nil
}

func (e Echoer) GetEcho(ctx context.Context, msg *v1.StringMessage) (*v1.StringMessage, error) {
	method := "GetEcho"
	start := time.Now()

	item := e.mongo.FindOne(ctx, msg)
	e.m.Request.WithLabelValues(method).Observe(float64(time.Since(start).Seconds()))

	b, err := item.Raw()
	if err != nil {
		e.m.Fail.WithLabelValues(method).Add(1)
		return nil, err
	}

	e.m.Success.WithLabelValues(method).Add(1)
	return &v1.StringMessage{
		Value: b.String(),
	}, nil
}
