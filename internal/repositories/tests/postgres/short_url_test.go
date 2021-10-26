//go:build integration
// +build integration

package postgres

import (
	"github.com/phpCoder88/url-shortener/internal/entities"
	"github.com/phpCoder88/url-shortener/internal/repositories/postgres"
)

func (s *PgRepositoryTestSuite) TestShortURL_FindAll() {
	ShortURLRepo := postgres.NewPgShortURLRepository(s.db, s.conf.QueryTimeout)
	urls, err := ShortURLRepo.FindAll(10, 0)
	s.NoError(err)
	s.Empty(urls)

	url := &entities.ShortURL{
		LongURL: "https://golang.org/",
		Token:   "ztq6y0",
	}

	err = ShortURLRepo.Add(url)
	s.NoError(err)

	urls, err = ShortURLRepo.FindAll(10, 0)
	s.NoError(err)
	s.Equal(1, len(urls))
	s.Equal(url.LongURL, urls[0].LongURL)
	s.Equal(url.Token, urls[0].Token)
	s.NotEmpty(urls[0].ID)
	s.NotEmpty(urls[0].CreatedAt)
	s.True(urls[0].Enabled)
}

func (s *PgRepositoryTestSuite) TestShortURL_Add() {
	ShortURLRepo := postgres.NewPgShortURLRepository(s.db, s.conf.QueryTimeout)
	url := &entities.ShortURL{
		LongURL: "https://golang.org/",
		Token:   "ztq6y0",
	}
	err := ShortURLRepo.Add(url)
	s.NoError(err)

	dbURL, err := ShortURLRepo.FindByURL(url.LongURL)
	s.NoError(err)
	s.NotEmpty(dbURL)
	s.Equal(url.LongURL, dbURL.LongURL)
	s.Equal(url.Token, dbURL.Token)
	s.NotEmpty(dbURL.ID)
	s.NotEmpty(dbURL.CreatedAt)
	s.True(dbURL.Enabled)

	err = ShortURLRepo.Add(url)
	s.Error(err)
}

func (s *PgRepositoryTestSuite) TestShortURL_FindByURLOrToken() {
	ShortURLRepo := postgres.NewPgShortURLRepository(s.db, s.conf.QueryTimeout)
	url := &entities.ShortURL{
		LongURL: "https://golang.org/",
		Token:   "ztq6y0",
	}

	dbURL, err := ShortURLRepo.FindByURL(url.LongURL)
	s.NoError(err)
	s.Empty(dbURL)

	err = ShortURLRepo.Add(url)
	s.NoError(err)

	dbURL, err = ShortURLRepo.FindByURL(url.LongURL)
	s.NoError(err)
	s.NotEmpty(dbURL)
	s.Equal(url.LongURL, dbURL.LongURL)
	s.Equal(url.Token, dbURL.Token)
	s.NotEmpty(dbURL.ID)
	s.NotEmpty(dbURL.CreatedAt)
	s.True(dbURL.Enabled)

	dbURL, err = ShortURLRepo.FindByToken(url.Token)
	s.NoError(err)
	s.NotEmpty(dbURL)
	s.Equal(url.LongURL, dbURL.LongURL)
	s.Equal(url.Token, dbURL.Token)
	s.NotEmpty(dbURL.ID)
	s.NotEmpty(dbURL.CreatedAt)
	s.True(dbURL.Enabled)
}
