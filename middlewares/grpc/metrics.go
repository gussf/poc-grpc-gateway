package grpc

import (
	"context"
	"time"

	"github.com/gussf/poc-grpc-gateway/metrics"
	"google.golang.org/grpc"
)

func HttpMetricsInterceptor() grpc.UnaryServerInterceptor {
	metrics := metrics.NewAPIMetrics()
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		metrics.Requests.WithLabelValues("obama", info.FullMethod).Observe(time.Since(start).Seconds())
		return resp, err
	}
}
