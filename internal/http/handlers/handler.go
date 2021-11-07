package handlers

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/phpCoder88/url-shortener-observable/internal/ioc"
)

type Handler struct {
	logger    *zap.SugaredLogger
	tracer    opentracing.Tracer
	container *ioc.Container
}

func NewHandler(logger *zap.SugaredLogger, container *ioc.Container, tracer opentracing.Tracer) *Handler {
	return &Handler{
		logger:    logger,
		tracer:    tracer,
		container: container,
	}
}
