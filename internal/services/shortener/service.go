package shortener

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/speps/go-hashids/v2"

	"github.com/phpCoder88/url-shortener-observable/internal/dto"
	"github.com/phpCoder88/url-shortener-observable/internal/entities"
	"github.com/phpCoder88/url-shortener-observable/internal/repositories/interfaces"
)

type Service struct {
	shortURLRepo interfaces.ShortURLRepository
	urlVisitRepo interfaces.URLVisitRepository
	tracer       opentracing.Tracer
}

func NewService(
	shortURLRepo interfaces.ShortURLRepository,
	urlVisitRepo interfaces.URLVisitRepository,
	tracer opentracing.Tracer,
) *Service {
	return &Service{
		shortURLRepo: shortURLRepo,
		urlVisitRepo: urlVisitRepo,
		tracer:       tracer,
	}
}

func (s *Service) FindAll(ctx context.Context, limit, offset int64) ([]dto.ShortURLReportDto, error) {
	span, newCtx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ShortenerService.FindAll")
	defer span.Finish()

	return s.shortURLRepo.FindAll(newCtx, limit, offset)
}

func (s *Service) CreateShortURL(ctx context.Context, urlStr string) (*entities.ShortURL, bool, error) {
	span, newCtx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ShortenerService.CreateShortURL")
	defer span.Finish()

	urlRecord, exists, err := s.IsURLExists(newCtx, urlStr)
	if err != nil {
		return nil, false, err
	}

	if exists {
		return urlRecord, true, nil
	}

	token, err := s.shortURL(newCtx, urlStr)
	if err != nil {
		return nil, false, err
	}

	urlRecord = &entities.ShortURL{
		LongURL:   urlStr,
		Token:     token,
		Enabled:   true,
		CreatedAt: time.Now(),
	}

	err = s.shortURLRepo.Add(newCtx, urlRecord)
	if err != nil {
		return nil, false, err
	}

	return urlRecord, false, nil
}

func (s *Service) IsURLExists(ctx context.Context, urlStr string) (*entities.ShortURL, bool, error) {
	span, newCtx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ShortenerService.IsURLExists")
	defer span.Finish()

	urlRecord, err := s.shortURLRepo.FindByURL(newCtx, urlStr)
	if err != nil {
		return nil, false, err
	}

	return urlRecord, urlRecord != nil, nil
}

func (s *Service) shortURL(ctx context.Context, urlStr string) (string, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ShortenerService.shortURL")
	defer span.Finish()

	hd := hashids.NewData()
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return "", err
	}

	token, err := h.Encode([]int{int(time.Now().UnixNano()) + len(urlStr)})
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) VisitFullURL(ctx context.Context, token, userIP string) (string, error) {
	span, newCtx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ShortenerService.VisitFullURL")
	defer span.Finish()

	urlRecord, err := s.shortURLRepo.FindByToken(newCtx, token)
	if err != nil {
		return "", err
	}

	err = s.urlVisitRepo.AddURLVisit(newCtx, urlRecord.ID, userIP)
	if err != nil {
		return "", err
	}

	return urlRecord.LongURL, nil
}

func (s *Service) ParseLimitOffsetQueryParams(ctx context.Context, query url.Values, param string, defaultVal int64) (int64, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ShortenerService.ParseLimitOffsetQueryParams_"+param)
	defer span.Finish()

	var paramInt int64
	var err error

	if paramSlice, ok := query[param]; ok {
		if len(paramSlice) > 1 {
			return 0, errors.New("too many values for param: " + param)
		}

		paramInt, err = strconv.ParseInt(paramSlice[0], 10, 64)
		if err != nil {
			return 0, errors.New(param + " param value isn't correct number")
		}

		if paramInt < 0 {
			return 0, errors.New(param + " param value is negative")
		}

		return paramInt, nil
	}

	return defaultVal, nil
}
