package handlers

import (
	"go.uber.org/zap"

	"github.com/phpCoder88/url-shortener/internal/ioc"
	"github.com/prometheus/client_golang/prometheus"
)

type Handler struct {
	logger    *zap.SugaredLogger
	container *ioc.Container

	infoCounter      prometheus.Counter
	redirectCounter  prometheus.Counter
	latencyHistogram *prometheus.HistogramVec
}

func NewHandler(logger *zap.SugaredLogger, container *ioc.Container) *Handler {
	handler := &Handler{
		logger:    logger,
		container: container,
		infoCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "shortener",
			Name:      "info_requests",
			Help:      "The number of build info requests",
		}),
		redirectCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "shortener",
			Name:      "redirect_counts",
			Help:      "The number of redirects",
		}),
		latencyHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "shortener",
			Name:      "latency",
			Help:      "The distribution of the latencies",
			Buckets:   []float64{0, 25, 50, 75, 100, 200, 400, 600, 800, 1000, 2000, 4000, 6000},
		}, []string{"handler"}),
	}

	prometheus.MustRegister(handler.infoCounter)
	prometheus.MustRegister(handler.redirectCounter)
	prometheus.MustRegister(handler.latencyHistogram)

	return handler
}
