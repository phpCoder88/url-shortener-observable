package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/suite"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type PgURLVisitRepositoryTestSuite struct {
	suite.Suite
	mDB    *sql.DB
	mock   sqlmock.Sqlmock
	db     *sqlx.DB
	tracer *mocktracer.MockTracer
}

func (s *PgURLVisitRepositoryTestSuite) SetupTest() {
	var err error
	s.mDB, s.mock, err = sqlmock.New()
	s.Require().NoError(err)
	s.db = sqlx.NewDb(s.mDB, "sqlmock")
	s.tracer = mocktracer.New()
}

func (s *PgURLVisitRepositoryTestSuite) TearDownTest() {
	s.db.Close()
	s.mDB.Close()
}

func (s *PgURLVisitRepositoryTestSuite) TestPgURLVisitRepository_AddURLVisit() {
	type args struct {
		urlID int64
		ip    string
	}
	type mockBehavior func(args args)
	testTable := []struct {
		name         string
		args         args
		timeout      time.Duration
		wantError    bool
		mockBehavior mockBehavior
	}{
		{
			name: "ok",
			args: args{
				urlID: 1,
				ip:    "127.0.0.1",
			},
			timeout:   time.Second,
			wantError: false,
			mockBehavior: func(args args) {
				s.mock.ExpectExec("INSERT INTO url_visits").
					WithArgs(args.urlID, args.ip, AnyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "with timeout error",
			args: args{
				urlID: 1,
				ip:    "127.0.0.1",
			},
			timeout:   500 * time.Microsecond,
			wantError: true,
			mockBehavior: func(args args) {
				s.mock.ExpectExec("INSERT INTO url_visits").
					WithArgs(args.urlID, args.ip, AnyTime{}).
					WillReturnError(ErrContextDeadlineExceeded)
			},
		},
	}

	for _, tt := range testTable {
		s.Run(tt.name, func() {
			tt.mockBehavior(tt.args)

			repository := NewPgURLVisitRepository(s.db, tt.timeout, s.tracer)
			err := repository.AddURLVisit(context.Background(), tt.args.urlID, tt.args.ip)
			if tt.wantError {
				s.Error(err)
			} else {
				s.NoError(err)
			}

			err = s.mock.ExpectationsWereMet()
			s.NoError(err)
		})
	}
}

func TestPgURLVisitRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PgURLVisitRepositoryTestSuite))
}
