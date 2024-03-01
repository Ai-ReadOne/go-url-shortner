package controller

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ai-readone/go-url-shortner/internal/models"
	"github.com/ai-readone/go-url-shortner/logger"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
)

func TestMain(m *testing.M) {
	//Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Run the other tests
	os.Exit(m.Run())
}

// setup router engine to be used during testing
func getRouter() *gin.Engine {
	binding.Validator = new(models.DefaultValidator)
	router := gin.Default()
	router.Use(errorReporterMiddleware())
	return router
}

// recreates error reporter middleware to be used for test
// this is duplicated because import the errorReporterMiddleware
// from app package will lead to import cycle
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

func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {

	// Creates a response recorder
	w := httptest.NewRecorder()

	// Creates the service and process the above request.
	r.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}
