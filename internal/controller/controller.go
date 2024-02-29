package controller

import (
	"github.com/ai-readone/go-url-shortner/internal/service"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service service.UrlService
}

func NewController(service service.UrlService) *Controller {
	return &Controller{service: service}
}

func (c *Controller) CreateShortenedUrl(ctx *gin.Context) {
	type url struct {
		Url string `json:"url"`
	}

	var requestUrl url
	err := ctx.ShouldBindJSON(&requestUrl)
	if err != nil {
		ctx.Error(badRequest{err})
		return
	}

	shortenedUrl, err := c.service.CreateShortenedUrl(ctx, requestUrl.Url)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, gin.H{"shortened_url": shortenedUrl})
}
