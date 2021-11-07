package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"

	"github.com/phpCoder88/url-shortener-observable/internal/dto"
	"github.com/phpCoder88/url-shortener-observable/internal/entities"
)

type PgShortURLRepository struct {
	db      *sqlx.DB
	timeout time.Duration
	tracer  opentracing.Tracer
}

func NewPgShortURLRepository(db *sqlx.DB, timeout time.Duration, tracer opentracing.Tracer) *PgShortURLRepository {
	return &PgShortURLRepository{
		db:      db,
		timeout: timeout,
		tracer:  tracer,
	}
}

func (r *PgShortURLRepository) FindAll(ctx context.Context, limit, offset int64) ([]dto.ShortURLReportDto, error) {
	span, newCtx := opentracing.StartSpanFromContextWithTracer(ctx, r.tracer, "PgShortURLRepository.FindAll")
	defer span.Finish()

	var rows []dto.ShortURLReportDto
	timeoutCtx, cancel := context.WithTimeout(newCtx, r.timeout)
	defer cancel()

	query := `SELECT su.*, COUNT(uv.id) AS visits
				FROM short_urls su
				    LEFT OUTER JOIN url_visits uv ON su.id = uv.url_id
				GROUP BY su.id
				LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(timeoutCtx, &rows, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *PgShortURLRepository) Add(ctx context.Context, model *entities.ShortURL) error {
	span, newCtx := opentracing.StartSpanFromContextWithTracer(ctx, r.tracer, "PgShortURLRepository.Add")
	defer span.Finish()

	timeoutCtx, cancel := context.WithTimeout(newCtx, r.timeout)
	defer cancel()

	query := "INSERT INTO short_urls (long_url, token) VALUES ($1, $2)"
	_, err := r.db.ExecContext(timeoutCtx, query, model.LongURL, model.Token)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgShortURLRepository) FindByURL(ctx context.Context, url string) (*entities.ShortURL, error) {
	span, newCtx := opentracing.StartSpanFromContextWithTracer(ctx, r.tracer, "PgShortURLRepository.FindByURL")
	defer span.Finish()

	timeoutCtx, cancel := context.WithTimeout(newCtx, r.timeout)
	defer cancel()

	urlRecord := new(entities.ShortURL)
	err := r.db.GetContext(timeoutCtx, urlRecord, "SELECT * FROM short_urls WHERE long_url = $1", url)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return urlRecord, nil
}

func (r *PgShortURLRepository) FindByToken(ctx context.Context, token string) (*entities.ShortURL, error) {
	span, newCtx := opentracing.StartSpanFromContextWithTracer(ctx, r.tracer, "PgShortURLRepository.FindByToken")
	defer span.Finish()

	timeoutCx, cancel := context.WithTimeout(newCtx, r.timeout)
	defer cancel()

	urlRecord := new(entities.ShortURL)
	err := r.db.GetContext(timeoutCx, urlRecord, "SELECT * FROM short_urls WHERE token = $1", token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return urlRecord, nil
}
