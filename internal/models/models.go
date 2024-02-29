package models

import (
	"time"
)

type Url struct {
	ShortenedUrl string    `json:"shortened_url" validate:"min=7,max=7"`
	OriginalUrl  string    `json:"original_url" binding:"url"`
	CreatedAt    time.Time `json:"created_at"`
}
