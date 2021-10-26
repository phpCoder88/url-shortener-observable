package postgres

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"

	"github.com/phpCoder88/url-shortener/internal/dto"
	"github.com/phpCoder88/url-shortener/internal/entities"
)

var ErrContextDeadlineExceeded = errors.New("context deadline exceeded")

type PgShortURLRepositoryTestSuite struct {
	suite.Suite
	mDB  *sql.DB
	mock sqlmock.Sqlmock
	db   *sqlx.DB
}

func (s *PgShortURLRepositoryTestSuite) SetupTest() {
	var err error
	s.mDB, s.mock, err = sqlmock.New()
	s.Require().NoError(err)
	s.db = sqlx.NewDb(s.mDB, "sqlmock")
}

func (s *PgShortURLRepositoryTestSuite) TearDownTest() {
	s.db.Close()
	s.mDB.Close()
}

func (s *PgShortURLRepositoryTestSuite) TestPgShortURLRepository_Add() {
	type mockBehavior func(model *entities.ShortURL)
	tableTests := []struct {
		name         string
		model        *entities.ShortURL
		timeout      time.Duration
		wantError    bool
		mockBehavior mockBehavior
	}{
		{
			name: "ok",
			model: &entities.ShortURL{
				LongURL: "https://rxjs.dev/guide/overview",
				Token:   "G6X5g",
			},
			timeout:   time.Second,
			wantError: false,
			mockBehavior: func(model *entities.ShortURL) {
				s.mock.ExpectExec("INSERT INTO short_urls").
					WithArgs(model.LongURL, model.Token).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "with timeout error",
			model: &entities.ShortURL{
				LongURL: "https://rxjs.dev/guide/overview",
				Token:   "G6X5g",
			},
			timeout:   500 * time.Microsecond,
			wantError: true,
			mockBehavior: func(model *entities.ShortURL) {
				s.mock.ExpectExec("INSERT INTO short_urls").
					WithArgs(model.LongURL, model.Token).
					WillReturnError(ErrContextDeadlineExceeded)
			},
		},
	}

	for _, tt := range tableTests {
		s.Run(tt.name, func() {
			tt.mockBehavior(tt.model)

			repository := NewPgShortURLRepository(s.db, tt.timeout)
			err := repository.Add(tt.model)
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

func (s *PgShortURLRepositoryTestSuite) TestPgShortURLRepository_FindAll() {
	sqlQuery := `^SELECT (.+) FROM short_urls su LEFT OUTER JOIN url_visits uv ON (.+) GROUP BY su.id LIMIT \$1 OFFSET \$2$`
	type args struct {
		limit  int64
		offset int64
	}
	type mockBehavior func(args args, urls []dto.ShortURLReportDto)
	testTable := []struct {
		name         string
		args         args
		timeout      time.Duration
		wantError    bool
		want         []dto.ShortURLReportDto
		mockBehavior mockBehavior
	}{
		{
			name:      "ok",
			args:      args{limit: 10, offset: 0},
			timeout:   time.Second,
			wantError: false,
			want: []dto.ShortURLReportDto{
				{
					ID:        1,
					LongURL:   "https://rxjs.dev/guide/overview",
					Token:     "DkZ9P",
					Enabled:   true,
					CreatedAt: time.Now(),
					Visits:    7,
				},
				{
					ID:        13,
					LongURL:   "https://laravel.com/docs/8.x/validation",
					Token:     "G6X5",
					Enabled:   true,
					CreatedAt: time.Now().Add(time.Hour),
					Visits:    2,
				},
			},
			mockBehavior: func(args args, urls []dto.ShortURLReportDto) {
				rows := sqlmock.NewRows([]string{"id", "long_url", "token", "enabled", "created_at", "visits"})
				for _, url := range urls {
					rows.AddRow(url.ID, url.LongURL, url.Token, url.Enabled, url.CreatedAt, url.Visits)
				}
				s.mock.ExpectQuery(sqlQuery).WithArgs(args.limit, args.offset).WillReturnRows(rows)
			},
		},
		{
			name:      "with timeout error",
			args:      args{limit: 10, offset: 0},
			timeout:   500 * time.Microsecond,
			wantError: true,
			want: []dto.ShortURLReportDto{
				{
					ID:        1,
					LongURL:   "https://rxjs.dev/guide/overview",
					Token:     "DkZ9P",
					Enabled:   true,
					CreatedAt: time.Now(),
					Visits:    7,
				},
			},
			mockBehavior: func(args args, urls []dto.ShortURLReportDto) {
				rows := sqlmock.NewRows([]string{"id", "long_url", "token", "enabled", "created_at", "visits"})
				for _, url := range urls {
					rows.AddRow(url.ID, url.LongURL, url.Token, url.Enabled, url.CreatedAt, url.Visits)
				}
				s.mock.ExpectQuery(sqlQuery).WithArgs(args.limit, args.offset).WillReturnError(ErrContextDeadlineExceeded)
			},
		},
		{
			name:      "with empty result",
			args:      args{limit: 10, offset: 0},
			timeout:   time.Second,
			wantError: false,
			want:      nil,
			mockBehavior: func(args args, urls []dto.ShortURLReportDto) {
				rows := sqlmock.NewRows([]string{"id", "long_url", "token", "enabled", "created_at", "visits"})
				s.mock.ExpectQuery(sqlQuery).WithArgs(args.limit, args.offset).WillReturnRows(rows)
			},
		},
	}

	for _, tt := range testTable {
		s.Run(tt.name, func() {
			tt.mockBehavior(tt.args, tt.want)

			repository := NewPgShortURLRepository(s.db, tt.timeout)
			urls, err := repository.FindAll(tt.args.limit, tt.args.offset)
			if tt.wantError {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tt.want, urls)
			}

			err = s.mock.ExpectationsWereMet()
			s.NoError(err)
		})
	}
}

func (s *PgShortURLRepositoryTestSuite) TestPgShortURLRepository_FindByURL() {
	sqlQuery := `SELECT \* FROM short_urls WHERE long_url = \$1`
	tableTests := getTableTestDataForFindingByTokenOrURL(sqlQuery, "https://rxjs.dev/guide/overview", s.mock)

	for _, tt := range tableTests {
		s.Run(tt.name, func() {
			tt.mockBehavior(tt.colValue, tt.want)

			repository := NewPgShortURLRepository(s.db, tt.timeout)
			url, err := repository.FindByURL(tt.colValue)
			checkSearchByTokenOrURLResults(s, tt.wantError, tt.want, url, err)
		})
	}
}

func (s *PgShortURLRepositoryTestSuite) TestPgShortURLRepository_FindByToken() {
	sqlQuery := `SELECT \* FROM short_urls WHERE token = \$1`
	tableTests := getTableTestDataForFindingByTokenOrURL(sqlQuery, "jRgqL", s.mock)

	for _, tt := range tableTests {
		s.Run(tt.name, func() {
			tt.mockBehavior(tt.colValue, tt.want)

			repository := NewPgShortURLRepository(s.db, tt.timeout)
			url, err := repository.FindByToken(tt.colValue)
			checkSearchByTokenOrURLResults(s, tt.wantError, tt.want, url, err)
		})
	}
}

func checkSearchByTokenOrURLResults(s *PgShortURLRepositoryTestSuite, wantError bool, want, actual *entities.ShortURL, err error) {
	if wantError {
		s.Error(err)
	} else {
		s.NoError(err)
		s.Equal(want, actual)
	}

	err = s.mock.ExpectationsWereMet()
	s.NoError(err)
}

type mockSearchingByTokenOrURLBehavior func(column string, urlRec *entities.ShortURL)

type searchingByTokenOrURLTableTests struct {
	name         string
	colValue     string
	timeout      time.Duration
	wantError    bool
	want         *entities.ShortURL
	mockBehavior mockSearchingByTokenOrURLBehavior
}

func getTableTestDataForFindingByTokenOrURL(sqlQuery, colValue string, mock sqlmock.Sqlmock) []searchingByTokenOrURLTableTests {
	return []searchingByTokenOrURLTableTests{
		{
			name:      "ok",
			colValue:  colValue,
			timeout:   time.Second,
			wantError: false,
			want: &entities.ShortURL{
				ID:        17,
				LongURL:   "https://rxjs.dev/guide/overview",
				Token:     "G6X5g",
				Enabled:   true,
				CreatedAt: time.Now(),
			},
			mockBehavior: func(url string, urlRec *entities.ShortURL) {
				rows := sqlmock.NewRows([]string{"id", "long_url", "token", "enabled", "created_at"}).
					AddRow(urlRec.ID, urlRec.LongURL, urlRec.Token, urlRec.Enabled, urlRec.CreatedAt)
				mock.ExpectQuery(sqlQuery).WithArgs(url).WillReturnRows(rows)
			},
		},
		{
			name:      "not found",
			colValue:  colValue,
			timeout:   time.Second,
			wantError: false,
			want:      nil,
			mockBehavior: func(url string, urlRec *entities.ShortURL) {
				rows := sqlmock.NewRows([]string{"id", "long_url", "token", "enabled", "created_at"})
				mock.ExpectQuery(sqlQuery).WithArgs(url).WillReturnRows(rows)
			},
		},
		{
			name:      "with timeout error",
			colValue:  colValue,
			timeout:   500 * time.Microsecond,
			wantError: true,
			want:      nil,
			mockBehavior: func(url string, urlRec *entities.ShortURL) {
				mock.ExpectQuery(sqlQuery).WithArgs(url).WillReturnError(ErrContextDeadlineExceeded)
			},
		},
	}
}

func TestPgShortURLRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PgShortURLRepositoryTestSuite))
}
