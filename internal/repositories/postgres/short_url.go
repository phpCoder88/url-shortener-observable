package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/phpCoder88/url-shortener/internal/dto"
	"github.com/phpCoder88/url-shortener/internal/entities"
)

type PgShortURLRepository struct {
	db      *sqlx.DB
	timeout time.Duration
}

func NewPgShortURLRepository(db *sqlx.DB, timeout time.Duration) *PgShortURLRepository {
	return &PgShortURLRepository{
		db:      db,
		timeout: timeout,
	}
}

func (r *PgShortURLRepository) FindAll(limit, offset int64) ([]dto.ShortURLReportDto, error) {
	var rows []dto.ShortURLReportDto
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	query := `SELECT su.*, COUNT(uv.id) AS visits
				FROM short_urls su
				    LEFT OUTER JOIN url_visits uv ON su.id = uv.url_id
				GROUP BY su.id
				LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &rows, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *PgShortURLRepository) Add(model *entities.ShortURL) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	query := "INSERT INTO short_urls (long_url, token) VALUES ($1, $2)"
	_, err := r.db.ExecContext(ctx, query, model.LongURL, model.Token)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgShortURLRepository) FindByURL(url string) (*entities.ShortURL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	urlRecord := new(entities.ShortURL)
	err := r.db.GetContext(ctx, urlRecord, "SELECT * FROM short_urls WHERE long_url = $1", url)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return urlRecord, nil
}

func (r *PgShortURLRepository) FindByToken(token string) (*entities.ShortURL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	urlRecord := new(entities.ShortURL)
	err := r.db.GetContext(ctx, urlRecord, "SELECT * FROM short_urls WHERE token = $1", token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return urlRecord, nil
}
