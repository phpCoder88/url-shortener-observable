package interfaces

import (
	"context"

	"github.com/phpCoder88/url-shortener-observable/internal/dto"
	"github.com/phpCoder88/url-shortener-observable/internal/entities"
)

type ShortURLRepository interface {
	FindAll(context.Context, int64, int64) ([]dto.ShortURLReportDto, error)
	FindByURL(context.Context, string) (*entities.ShortURL, error)
	Add(context.Context, *entities.ShortURL) error
	FindByToken(context.Context, string) (*entities.ShortURL, error)
}
