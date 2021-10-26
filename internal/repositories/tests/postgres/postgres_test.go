//go:build integration
// +build integration

package postgres

import (
	"testing"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"

	"github.com/phpCoder88/url-shortener/internal/config"
	pgStorage "github.com/phpCoder88/url-shortener/internal/storages/postgres"
)

type PgRepositoryTestSuite struct {
	suite.Suite
	db          *sqlx.DB
	conf        *config.DBConfig
	dbMigration *migrate.Migrate
}

func (s *PgRepositoryTestSuite) SetupTest() {
	var err error
	s.conf, err = config.GetDBConfig()
	s.Require().NoError(err)

	s.db, err = pgStorage.NewPgConnection(s.conf)
	s.Require().NoError(err)

	s.dbMigration, err = migrate.New("file://../../../../migrations", pgStorage.GetConnectionString(s.conf))
	s.Require().NoError(err)

	if err = s.dbMigration.Up(); err != nil && err != migrate.ErrNoChange {
		s.Require().NoError(err)
	}
}

func (s *PgRepositoryTestSuite) TearDownTest() {
	s.NoError(s.dbMigration.Down())
	s.db.Close()
}

func TestPgRepositoriesTestSuite(t *testing.T) {
	suite.Run(t, new(PgRepositoryTestSuite))
}
