package controller

import (
	"net/http"

	"github.com/ai-readone/go-url-shortner/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	service service.UrlService
}

func NewController(service service.UrlService) *Controller {
	return &Controller{service: service}
}

func (c *Controller) CreateShortenedUrl(ctx *gin.Context) {
	type originalUrl struct {
		Url string `json:"url"`
	}

	var original originalUrl
	err := ctx.ShouldBindJSON(&original)
	if err != nil {
		ctx.Error(badRequest{err})
		return
	}

	shortenedUrl, err := c.service.CreateShortenedUrl(ctx, original.Url)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, gin.H{"shortened_url": shortenedUrl})
}

func (c *Controller) GetOriginalUrl(ctx *gin.Context) {
	type shortenedUrl struct {
		Url string `uri:"shortened"`
	}

	var shortened shortenedUrl
	err := ctx.ShouldBindUri(&shortened)
	if err != nil {
		ctx.Error(badRequest{err})
		return
	}

	url, err := c.service.GetOriginalUrl(ctx, shortened.Url)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			ctx.Error(err)
		} else {
			ctx.Error(notFound{err})
		}
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, url)
}
