package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

type Service interface {
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) ShortenUrl(ctx context.Context, url string) (string, error) {
	// encrypts the url using sha-256 hash function.
	sh := sha256.New()
	_, err := sh.Write([]byte(url))
	if err != nil {
		return "", err
	}

	// generate a random 8-byte string,
	// which will be appended to the hash output above to incrase its uniqueness,
	// thereby making it more unlikely for urls to produce same base64 encoding output
	randomBytes := make([]byte, 8)
	_, err = rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	fmt.Print(randomBytes)

	// encode the hash output to base64 format,
	// return the first 9 characters of the string,
	// which will be sued as the shortened version of the provided url
	shrotenedUrl := base64.URLEncoding.EncodeToString(sh.Sum(randomBytes))[:9]
	return shrotenedUrl, nil
}
