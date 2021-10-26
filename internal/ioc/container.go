package ioc

import (
	"time"

	"github.com/phpCoder88/url-shortener/internal/repositories/postgres"

	"github.com/jmoiron/sqlx"

	"github.com/phpCoder88/url-shortener/internal/services/shortener"
)

type Container struct {
	ShortenerService *shortener.Service
}

func NewContainer(db *sqlx.DB, queryTimeout time.Duration) *Container {
	shortURLRepo := postgres.NewPgShortURLRepository(db, queryTimeout)
	urlVisitRepo := postgres.NewPgURLVisitRepository(db, queryTimeout)

	return &Container{
		ShortenerService: shortener.NewService(shortURLRepo, urlVisitRepo),
	}
}
