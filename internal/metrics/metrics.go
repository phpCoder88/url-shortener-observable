package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	InfoCounter      prometheus.Counter
	RedirectCounter  prometheus.Counter
	LatencyHistogram *prometheus.HistogramVec
}

const (
	Namespace    = "shortener"
	LabelHandler = "handler"
)

func New() *Metrics {
	metric := &Metrics{
		InfoCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "info_requests",
			Help:      "The number of build info requests",
		}),
		RedirectCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "redirect_counts",
			Help:      "The number of redirects",
		}),
		LatencyHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: Namespace,
			Name:      "latency",
			Help:      "The distribution of the latencies",
			Buckets:   []float64{0, 25, 50, 75, 100, 200, 400, 600, 800, 1000, 2000, 4000, 6000},
		}, []string{LabelHandler}),
	}

	prometheus.MustRegister(metric.InfoCounter)
	prometheus.MustRegister(metric.RedirectCounter)
	prometheus.MustRegister(metric.LatencyHistogram)

	return metric
}
