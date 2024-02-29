package app

import (
	"net/http"

	"github.com/ai-readone/go-url-shortner/configs"
	"github.com/ai-readone/go-url-shortner/internal/controller"
	"github.com/ai-readone/go-url-shortner/internal/database"
	"github.com/ai-readone/go-url-shortner/internal/models"
	"github.com/ai-readone/go-url-shortner/internal/service"
	"github.com/ai-readone/go-url-shortner/internal/store"
	"github.com/ai-readone/go-url-shortner/logger"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
)

func RegisterRoutes(conf *configs.Config) *gin.Engine {
	binding.Validator = new(models.DefaultValidator)

	router := gin.Default()
	router.MaxMultipartMemory = 2 << 20 // 2 KiB

	// response can contain only one Access-Control-Allow-Origin so we check if origin is allowed and set the header to the origin
	origins := make(map[string]bool, len(conf.AllowedOrigins))
	for _, i := range conf.AllowedOrigins {
		origins[i] = true
	}

	f := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			origin := c.GetHeader("Origin")
			if _, found := origins[origin]; found {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}

			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, accept, origin, Cache-Control")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")

			// if c.Request.Method != "POST" && c.Request.Method != "GET" {
			// 	c.AbortWithStatus(405)
			// 	return
			// }

			c.Next()
		}
	}
	router.Use(f())
	router.Use(gin.Recovery())
	router.Use(errorReporterMiddleware())

	urlStore := store.NewStore(database.GetSession())
	urlService := service.NewService(urlStore)
	urlController := controller.NewController(urlService)

	router.POST("/shorten", urlController.CreateShortenedUrl)
	router.GET("/:shortened", urlController.GetOriginalUrl)

	logger.Info("routes registered successfully!")

	return router
}

func errorReporterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		detectedErrors := c.Errors.ByType(gin.ErrorTypeAny)

		if len(detectedErrors) > 0 {
			logger.Error(detectedErrors)

			err := detectedErrors[0].Err

			var code int
			switch errors.Cause(err).(type) {
			case BadRequester:
				code = http.StatusBadRequest
			case NotFounder:
				code = http.StatusNotFound
			default:
				code = http.StatusInternalServerError
				err = errors.New("unexpected error encountered")
			}
			// Put the error into response
			c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
	}
}

type BadRequester interface {
	BadRequest()
}

type NotFounder interface {
	NotFound()
}
