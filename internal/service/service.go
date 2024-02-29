package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/ai-readone/go-url-shortner/internal/models"
	"github.com/ai-readone/go-url-shortner/internal/store"
	"gorm.io/gorm"
)

type UrlService interface {
	CreateShortenedUrl(ctx context.Context, url string) (string, error)
	GetOriginalUrl(ctx context.Context, ShortenedUrl string) (string, error)
}

type Service struct {
	store store.UrlStore
}

func NewService(store store.UrlStore) *Service {
	return &Service{store: store}
}

func (s *Service) urlShortener(url string) (string, error) {
	// encrypts the url using sha-256 hash function.
	sh := sha256.New()
	_, err := sh.Write([]byte(url))
	if err != nil {
		return "", errors.New(fmt.Sprintf("error while shortening url, Error: %s", err))
	}

	// generate a random 8-byte string,
	// which will be appended to the hash output above to incrase its uniqueness,
	// thereby making it more unlikely for urls to produce same base64 encoding output
	randomBytes := make([]byte, 8)
	_, err = rand.Read(randomBytes)
	if err != nil {
		return "", errors.New(fmt.Sprintf("unexpected error while shortening url, Error: %s", err))
	}
	fmt.Print(randomBytes)

	// encode the hash output to base64 format,
	// return the first 7 characters of the string,
	// which will be used as the shortened version of the provided url
	shrotenedUrl := base64.URLEncoding.EncodeToString(sh.Sum(randomBytes))[:7]
	return shrotenedUrl, nil
}

func (s *Service) CreateShortenedUrl(ctx context.Context, url string) (string, error) {
	// checks the url for trailing slash,
	// removes the trailing slash if present.
	// this helps treat url w/o trailing slash as same.
	if url[len(url)-1] == '/' {
		url = url[:len(url)-2]
	}

	// checks if the url has been previoulsy shortened,
	// if yes retrieves and return the existing shotened url
	// if no, creates a new shortened url and saves it in the database.
	existing, err := s.store.GetExistingShortenedUrl(ctx, url)
	if err == nil {
		return existing.ShortenedUrl, nil
	}
	if err != gorm.ErrRecordNotFound {
		return "", errors.New(fmt.Sprintf("unexpected error while shortening url, Error: %s", err))
	}

	// generate a shortened_url for the provided url
	shortenedUrl, err := s.urlShortener(url)
	if err != nil {
		return "", errors.New(fmt.Sprintf("unexpected error while shortening url, Error: %s", err))
	}

	data := models.Url{
		ShortenedUrl: shortenedUrl,
		OriginalUrl:  url,
	}

	if err := s.store.CreateShortenedUrl(ctx, &data); err != nil {
		// checks if the error is a unique key constraints violation err,
		// meaning that the shortened_url column already exist in the database
		// since the original_url has been confirmed not to exist,
		// if true I recursively try to generate and store a unique shrotened url.
		if err == gorm.ErrDuplicatedKey {
			return s.CreateShortenedUrl(ctx, url)
		}

		return "", errors.New(fmt.Sprintf("unexpected error while shortening url, Error: %s", err))
	}

	return shortenedUrl, nil
}

func (s *Service) GetOriginalUrl(ctx context.Context, ShortenedUrl string) (string, error) {
	url, err := s.store.GetOriginalUrl(ctx, ShortenedUrl)
	if err != nil {
		return "", errors.New(fmt.Sprintf("unexpected error while redirecting to original url, Error: %s", err))
	}

	return url.OriginalUrl, nil
}
