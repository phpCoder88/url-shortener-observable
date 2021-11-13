package interfaces

import (
	"github.com/phpCoder88/url-shortener-observable/internal/dto"
	"github.com/phpCoder88/url-shortener-observable/internal/entities"
)

type ShortURLRepository interface {
	FindAll(int64, int64) ([]dto.ShortURLReportDto, error)
	FindByURL(string) (*entities.ShortURL, error)
	Add(*entities.ShortURL) error
	FindByToken(string) (*entities.ShortURL, error)
}
