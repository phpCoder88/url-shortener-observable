package handlers

import (
	"go.uber.org/zap"

	"github.com/phpCoder88/url-shortener/internal/ioc"
	"github.com/phpCoder88/url-shortener/internal/metrics"
)

type Handler struct {
	logger    *zap.SugaredLogger
	container *ioc.Container
	metrics   *metrics.Metrics
}

func NewHandler(logger *zap.SugaredLogger, container *ioc.Container) *Handler {
	return &Handler{
		logger:    logger,
		container: container,
		metrics:   metrics.New(),
	}
}
