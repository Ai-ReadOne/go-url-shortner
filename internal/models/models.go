package models

import "time"

type Url struct {
	ShortenedUrl string    `json:"shortened_url"`
	OriginalUrl  string    `json:"original_url"`
	CreatedAt    time.Time `json:"created_at"`
}
