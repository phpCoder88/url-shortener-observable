package handlers

import (
	"github.com/phpCoder88/url-shortener-observable/internal/ioc"
	"go.uber.org/zap"
)

type Handler struct {
	logger    *zap.SugaredLogger
	container *ioc.Container
}

func NewHandler(logger *zap.SugaredLogger, container *ioc.Container) *Handler {
	return &Handler{
		logger:    logger,
		container: container,
	}
}
