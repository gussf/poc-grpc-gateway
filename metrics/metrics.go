package metrics

type Metrics struct {
	API APIMetrics
}

func New() Metrics {
	return Metrics{
		API: newAPIMetrics(),
	}
}
