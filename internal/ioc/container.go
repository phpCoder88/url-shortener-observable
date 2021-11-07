package ioc

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"

	"github.com/phpCoder88/url-shortener-observable/internal/repositories/postgres"
	"github.com/phpCoder88/url-shortener-observable/internal/services/shortener"
)

type Container struct {
	ShortenerService *shortener.Service
}

func NewContainer(db *sqlx.DB, queryTimeout time.Duration, tracer opentracing.Tracer) *Container {
	shortURLRepo := postgres.NewPgShortURLRepository(db, queryTimeout, tracer)
	urlVisitRepo := postgres.NewPgURLVisitRepository(db, queryTimeout, tracer)

	return &Container{
		ShortenerService: shortener.NewService(shortURLRepo, urlVisitRepo, tracer),
	}
}
