package server

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.uber.org/zap"
)

// MetricServer struct
type MetricServer struct {
	server          http.Server
	logger          *zap.SugaredLogger
	port            int
	shutdownTimeout time.Duration
	errors          chan error
}

func NewMetricServer(logger *zap.SugaredLogger, port int, shutdownTimeout time.Duration, errChan chan error) *MetricServer {
	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	return &MetricServer{
		server: http.Server{
			Addr:    net.JoinHostPort("", strconv.Itoa(port)),
			Handler: router,
		},
		logger:          logger,
		port:            port,
		shutdownTimeout: shutdownTimeout,
		errors:          errChan,
	}
}

func (s *MetricServer) Run() {
	go func() {
		s.logger.Infof("Metrics Server is listening on PORT: %d...", s.port)
		err := s.server.ListenAndServe()
		if err != nil {
			s.errors <- err
		}
	}()
}

func (s *MetricServer) Stop() error {
	s.logger.Info("Starting to shutdown the metrics server...")
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
