package store

import (
	"context"

	"github.com/ai-readone/go-url-shortner/internal/models"
	"gorm.io/gorm"
)

type UrlStore interface {
	CreateShortenedUrl(ctx context.Context, url *models.Url) error
	GetOriginalUrl(ctx context.Context, shortenedUrl string) (*models.Url, error)
	GetExistingShortenedUrl(ctx context.Context, orignalUrl string) (*models.Url, error)
}

type Store struct {
	pgClient *gorm.DB
}

func NewStore(pgClient *gorm.DB) *Store {
	return &Store{pgClient: pgClient}
}

// saves the shortened url record to the database.
func (s *Store) CreateShortenedUrl(ctx context.Context, url *models.Url) error {
	result := s.pgClient.Model(url).Create(url).WithContext(ctx)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// retrieves a url record from the database,
//
//	using the shortened version provided by the user.
func (s *Store) GetOriginalUrl(ctx context.Context, shortenedUrl string) (*models.Url, error) {
	url := &models.Url{}
	result := s.pgClient.Model(url).First(url, "shortened_url = ?", shortenedUrl).WithContext(ctx)
	if result.Error != nil {
		return nil, result.Error
	}

	return url, nil
}

func (s *Store) GetExistingShortenedUrl(ctx context.Context, orignalUrl string) (*models.Url, error) {
	url := &models.Url{}
	result := s.pgClient.Model(url).First(url, "original_url = ?", orignalUrl).WithContext(ctx)
	if result.Error != nil {
		return nil, result.Error
	}

	return url, nil
}
