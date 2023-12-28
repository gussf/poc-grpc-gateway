package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MongoDBMetrics struct {
	Success *prometheus.CounterVec
	Fail    *prometheus.CounterVec
	Request *prometheus.HistogramVec
}

func NewMongoDBMetrics() MongoDBMetrics {
	return MongoDBMetrics{
		Success: promauto.NewCounterVec(prometheus.CounterOpts{Name: "mongodb_success"}, []string{"method"}),
		Fail:    promauto.NewCounterVec(prometheus.CounterOpts{Name: "mongodb_failure"}, []string{"method"}),
		Request: promauto.NewHistogramVec(prometheus.HistogramOpts{Name: "mongodb_request"}, []string{"method"}),
	}
}
