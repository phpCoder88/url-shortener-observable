package shortener

import (
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/phpCoder88/url-shortener-observable/internal/repositories/interfaces"

	"github.com/speps/go-hashids/v2"

	"github.com/phpCoder88/url-shortener-observable/internal/dto"
	"github.com/phpCoder88/url-shortener-observable/internal/entities"
)

type Service struct {
	shortURLRepo interfaces.ShortURLRepository
	urlVisitRepo interfaces.URLVisitRepository
}

func NewService(shortURLRepo interfaces.ShortURLRepository, urlVisitRepo interfaces.URLVisitRepository) *Service {
	return &Service{
		shortURLRepo: shortURLRepo,
		urlVisitRepo: urlVisitRepo,
	}
}

func (s *Service) FindAll(limit, offset int64) ([]dto.ShortURLReportDto, error) {
	return s.shortURLRepo.FindAll(limit, offset)
}

func (s *Service) CreateShortURL(urlStr string) (*entities.ShortURL, bool, error) {
	urlRecord, exists, err := s.IsURLExists(urlStr)
	if err != nil {
		return nil, false, err
	}

	if exists {
		return urlRecord, true, nil
	}

	token, err := s.shortURL(urlStr)
	if err != nil {
		return nil, false, err
	}

	urlRecord = &entities.ShortURL{
		LongURL:   urlStr,
		Token:     token,
		Enabled:   true,
		CreatedAt: time.Now(),
	}

	err = s.shortURLRepo.Add(urlRecord)
	if err != nil {
		return nil, false, err
	}

	return urlRecord, false, nil
}

func (s *Service) IsURLExists(urlStr string) (*entities.ShortURL, bool, error) {
	urlRecord, err := s.shortURLRepo.FindByURL(urlStr)
	if err != nil {
		return nil, false, err
	}

	return urlRecord, urlRecord != nil, nil
}

func (s *Service) shortURL(urlStr string) (string, error) {
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

func (s *Service) GetFullURL(token string) (string, error) {
	urlRecord, err := s.shortURLRepo.FindByToken(token)
	if err != nil {
		return "", err
	}

	return urlRecord.LongURL, nil
}

func (s *Service) VisitFullURL(token, userIP string) (string, error) {
	urlRecord, err := s.shortURLRepo.FindByToken(token)
	if err != nil {
		return "", err
	}

	err = s.urlVisitRepo.AddURLVisit(urlRecord.ID, userIP)
	if err != nil {
		return "", err
	}

	return urlRecord.LongURL, nil
}

func (s *Service) ParseLimitOffsetQueryParams(query url.Values, param string, defaultVal int64) (int64, error) {
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
