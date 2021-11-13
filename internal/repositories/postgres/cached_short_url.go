package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/phpCoder88/url-shortener-observable/internal/dto"
	"github.com/phpCoder88/url-shortener-observable/internal/entities"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type PgCachedShortURLRepository struct {
	cache        *cache.Cache
	cacheTimeout time.Duration
	repo         *PgShortURLRepository
}

func NewPgCachedShortURLRepository(db *sqlx.DB, timeout time.Duration, cacheSize int) *PgCachedShortURLRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	rCache := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(cacheSize, time.Minute),
	})

	return &PgCachedShortURLRepository{
		cache:        rCache,
		cacheTimeout: time.Second,
		repo:         NewPgShortURLRepository(db, timeout),
	}
}

func (r *PgCachedShortURLRepository) FindAll(limit, offset int64) ([]dto.ShortURLReportDto, error) {
	return r.repo.FindAll(limit, offset)
}

func (r *PgCachedShortURLRepository) Add(model *entities.ShortURL) error {
	return r.repo.Add(model)
}

func (r *PgCachedShortURLRepository) FindByURL(url string) (*entities.ShortURL, error) {
	return r.repo.FindByURL(url)
}

func (r *PgCachedShortURLRepository) FindByToken(token string) (*entities.ShortURL, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), r.cacheTimeout)
	defer cancelFunc()

	key := fmt.Sprintf("token:%s", token)
	var cachedShortURL entities.ShortURL

	err := r.cache.Get(ctx, key, &cachedShortURL)

	switch err {
	case nil:
		return &cachedShortURL, nil

	case cache.ErrCacheMiss:
		var shortURL *entities.ShortURL
		shortURL, err = r.repo.FindByToken(token)
		if err != nil {
			return nil, err
		}

		err = r.cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   key,
			Value: shortURL,
			TTL:   time.Hour,
		})

		if err != nil {
			return nil, err
		}

		return shortURL, nil
	}

	return nil, err
}
