package postgres

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type PgURLVisitRepository struct {
	db      *sqlx.DB
	timeout time.Duration
}

func NewPgURLVisitRepository(db *sqlx.DB, timeout time.Duration) *PgURLVisitRepository {
	return &PgURLVisitRepository{
		db:      db,
		timeout: timeout,
	}
}

func (r *PgURLVisitRepository) AddURLVisit(urlID int64, userIP string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	query := "INSERT INTO url_visits (url_id, ip, created_at) VALUES ($1, $2, $3)"
	_, err := r.db.ExecContext(ctx, query, urlID, userIP, time.Now())
	if err != nil {
		return err
	}

	return nil
}
