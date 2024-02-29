package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

type serviceMock struct {
	mock.Mock
}

func (m *serviceMock) CreateShortenedUrl(ctx context.Context, url string) (string, error) {
	args := m.Called(ctx, url)
	return args.Get(0).(string), args.Error(1)
}

func (m *serviceMock) GetOriginalUrl(ctx context.Context, shortenedUrl string) (string, error) {
	args := m.Called(ctx, shortenedUrl)
	return args.Get(0).(string), args.Error(1)
}

func TestCreateShortenedUrl(t *testing.T) {
	testCases := []struct {
		name         string
		wantStatus   int
		serviceError error
		shortenedUrl string
		payload      []byte
		originalUrl  string
	}{
		{
			name:         "failure:fails-when-url-does-not-have-host-name",
			wantStatus:   http.StatusBadRequest,
			serviceError: nil,
			payload:      []byte(`{"url":"//aransiolaibrahim"}`),
		},
		{
			name:         "failure:fails-when-create-shortened-url-service-function-returns-an-error",
			wantStatus:   http.StatusInternalServerError,
			payload:      []byte(`{"url":"https://linkedin.com/in/aransiolaibrahim"}`),
			serviceError: errors.New("failed"),
		},
		{
			name:         "success:redirects-to-original-url",
			wantStatus:   http.StatusCreated,
			serviceError: nil,
			shortenedUrl: "jdhfjka",
			payload:      []byte(`{"url":"https://linkedin.com/in/aransiolaibrahim"}`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := new(serviceMock)
			service.On("CreateShortenedUrl", mock.Anything, mock.Anything).Return(tc.shortenedUrl, tc.serviceError).Once()
			controller := NewController(service)

			r := getRouter()
			r.POST("/shorten", controller.CreateShortenedUrl)
			req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(tc.payload))

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				return w.Code == tc.wantStatus
			})
		})
	}
}

func TestGetOriginalUrls(t *testing.T) {
	testCases := []struct {
		name         string
		wantStatus   int
		serviceError error
		shortenedUrl string
		originalUrl  string
	}{
		{
			name:         "failure:fails-when-shortended-url-not-in-path",
			wantStatus:   http.StatusNotFound,
			serviceError: nil,
		},
		{
			name:         "failure:fails-when-shortended-url-lenght-less-than-7",
			wantStatus:   http.StatusBadRequest,
			serviceError: nil,
			shortenedUrl: "jdhf",
		},
		{
			name:         "failure:fails-when-shortended-url-lenght-more-than-7",
			wantStatus:   http.StatusBadRequest,
			serviceError: nil,
			shortenedUrl: "jdhfjkass",
		},
		{
			name:         "failure:fails-when-get-original-url-service-function-returns-unexpected-error",
			wantStatus:   http.StatusInternalServerError,
			serviceError: errors.New("failed"),
			shortenedUrl: "jdhfjka",
		},
		{
			name:         "failure:fails-when-get-original-url-service-function-returns-not-found-error",
			wantStatus:   http.StatusNotFound,
			serviceError: notFound{errors.New("not found")},
			shortenedUrl: "jdhfjka",
		},
		{
			name:         "success:redirects-to-original-url",
			wantStatus:   http.StatusMovedPermanently,
			serviceError: nil,
			shortenedUrl: "jdhfjka",
			originalUrl:  "https://linkedin.com/in/aransiolaibrahim",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := new(serviceMock)
			service.On("GetOriginalUrl", mock.Anything, mock.Anything).Return(tc.originalUrl, tc.serviceError).Once()
			controller := NewController(service)

			r := getRouter()

			r.GET("/:shortened", controller.GetOriginalUrl)
			req, _ := http.NewRequest("GET", fmt.Sprintf("/%s", tc.shortenedUrl), bytes.NewBuffer(nil))

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				return w.Code == tc.wantStatus
			})
		})
	}
}
