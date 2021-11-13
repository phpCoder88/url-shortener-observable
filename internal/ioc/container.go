package ioc

import (
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/phpCoder88/url-shortener-observable/internal/repositories/postgres"
	"github.com/phpCoder88/url-shortener-observable/internal/services/shortener"
)

type Container struct {
	ShortenerService *shortener.Service
}

func NewContainer(db *sqlx.DB, queryTimeout time.Duration) *Container {
	shortURLRepo := postgres.NewPgCachedShortURLRepository(db, queryTimeout, 100)
	urlVisitRepo := postgres.NewPgURLVisitRepository(db, queryTimeout)

	return &Container{
		ShortenerService: shortener.NewService(shortURLRepo, urlVisitRepo),
	}
}
