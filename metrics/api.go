package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type APIMetrics struct {
	Requests *prometheus.HistogramVec
}

func NewAPIMetrics() APIMetrics {
	return APIMetrics{
		Requests: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "api_requests",
			},
			[]string{"client", "method"},
		),
	}
}
