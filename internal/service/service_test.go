package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ai-readone/go-url-shortner/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type storeMock struct {
	mock.Mock
}

func (m *storeMock) CreateShortenedUrl(ctx context.Context, url *models.Url) error {
	args := m.Called(ctx, url)
	return args.Error(0)
}

func (m *storeMock) GetOriginalUrl(ctx context.Context, shortenedUrl string) (*models.Url, error) {
	args := m.Called(ctx, shortenedUrl)
	return args.Get(0).(*models.Url), args.Error(1)
}

func (m *storeMock) GetExistingShortenedUrl(ctx context.Context, orignalUrl string) (*models.Url, error) {
	args := m.Called(ctx, orignalUrl)
	return args.Get(0).(*models.Url), args.Error(1)
}

func TestCreateShortenedUrl(t *testing.T) {
	testCases := []struct {
		name                     string
		storeError               error
		getExistingShortError    error
		createShortUrlStoreError error
		wantError                error
		url                      string
	}{
		{
			name:                  "failure:fails-when-get-existing-shortened-url-return-unexpected-error",
			getExistingShortError: errors.New("get existing error"),
			wantError:             errors.New("unexpected error while shortening url, Error: get existing error"),
			url:                   "https://nairametrics.com",
		},
		{
			name:      "sucess:suceeds-when-url-has-already-been-shortened",
			wantError: nil,
			url:       "perizer.com",
		},
		{
			name:      "success:service-ok",
			wantError: nil,
			url:       "https://nairametrics.com",
		},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := new(storeMock)
			store.On("GetExistingShortenedUrl", mock.Anything, "perizer.com").Return(&models.Url{ShortenedUrl: "dloejsa"}, nil).Once()
			store.On("GetExistingShortenedUrl", mock.Anything, mock.Anything).Return(&models.Url{}, tc.getExistingShortError).Once()
			store.On("CreateShortenedUrl", mock.Anything, mock.Anything).Return(tc.createShortUrlStoreError).Once()
			svc := NewService(store)

			ctx := context.TODO()

			_, err := svc.CreateShortenedUrl(ctx, tc.url)
			if tc.wantError == nil {
				assert.NoError(err, "Should have no error")
			} else {
				assert.EqualError(err, tc.wantError.Error(), "Should have equal error")
			}
		})
	}
}

func TestSaveShortenedUrl(t *testing.T) {
	testCases := []struct {
		name                     string
		storeError               error
		createShortUrlStoreError error
		wantError                error
		url                      string
	}{
		{
			name:                     "failure:fails-when-save-shortened-url-return-unexpected-error",
			createShortUrlStoreError: errors.New("save url error"),
			wantError:                errors.New("unexpected error while shortening url, Error: save url error"),
			url:                      "https://nairametrics.com",
		},
		{
			name:      "success:service-ok",
			wantError: nil,
			url:       "https://nairametrics.com",
		},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := new(storeMock)
			store.On("CreateShortenedUrl", mock.Anything, mock.Anything).Return(tc.createShortUrlStoreError).Once()
			svc := NewService(store)

			ctx := context.TODO()

			_, err := svc.saveShortenedURL(ctx, tc.url)
			if tc.wantError == nil {
				assert.NoError(err, "Should have no error")
			} else {
				assert.EqualError(err, tc.wantError.Error(), "Should have equal error")
			}
		})
	}
}

func TestGetOriginalUrl(t *testing.T) {
	testCases := []struct {
		name                string
		storeError          error
		getOriginalUrlError error
		wantError           error
		url                 string
	}{
		{
			name:                "failure:fails-when-get-original-url-return-error",
			getOriginalUrlError: errors.New("get original url error"),
			wantError:           errors.New("unexpected error while redirecting to original url, Error: get original url error"),
			url:                 "https://nairametrics.com",
		},
		{
			name:      "success:service-ok",
			wantError: nil,
			url:       "https://nairametrics.com",
		},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := new(storeMock)
			store.On("GetOriginalUrl", mock.Anything, mock.Anything).Return(&models.Url{}, tc.getOriginalUrlError).Once()
			svc := NewService(store)

			ctx := context.TODO()

			_, err := svc.GetOriginalUrl(ctx, tc.url)
			if tc.wantError == nil {
				assert.NoError(err, "Should have no error")
			} else {
				assert.EqualError(err, tc.wantError.Error(), "Should have equal error")
			}
		})
	}
}
