package postgres

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

type PgURLVisitRepository struct {
	db      *sqlx.DB
	timeout time.Duration
	tracer  opentracing.Tracer
}

func NewPgURLVisitRepository(db *sqlx.DB, timeout time.Duration, tracer opentracing.Tracer) *PgURLVisitRepository {
	return &PgURLVisitRepository{
		db:      db,
		timeout: timeout,
		tracer:  tracer,
	}
}

func (r *PgURLVisitRepository) AddURLVisit(ctx context.Context, urlID int64, userIP string) error {
	span, newCtx := opentracing.StartSpanFromContextWithTracer(ctx, r.tracer, "PgURLVisitRepository.AddURLVisit")
	defer span.Finish()

	timeoutCtx, cancel := context.WithTimeout(newCtx, r.timeout)
	defer cancel()

	query := "INSERT INTO url_visits (url_id, ip, created_at) VALUES ($1, $2, $3)"
	_, err := r.db.ExecContext(timeoutCtx, query, urlID, userIP, time.Now())
	if err != nil {
		return err
	}

	return nil
}
