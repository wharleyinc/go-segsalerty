//go:generate mockgen -source=http_adapter.go -destination mocks/http_adapter_mock.go -package httprequestadapter_test . Service

package httprequestadapter

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-segsalerty/common/logger"
	"go-segsalerty/internal/model"
	"go.uber.org/zap"
)

const createShortLink = "/shorten"

const (
	ErrorInvalidRequest = "validation: invalid request"
)

type Service interface {
	ShortenUrl(ctx context.Context, newUrl model.Shortener) (*model.ShortenerDetails, error)
}

type apiShortUrl struct {
	LongUrl string `json:"longUrl" binding:"required"`
}

type apiError struct {
	Code    int    `json:"code" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type apiResult struct {
	Code int         `json:"code" binding:"required"`
	Data interface{} `json:"data,omitempty"`
}

type adapter struct {
	service Service
}

func NewHttpHandler(service Service) *adapter {
	return &adapter{service: service}
}

func (a adapter) ApplyRoutes(routes gin.IRouter) {
	routes.POST(createShortLink, create(a.service))
}

func create(service Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var newLink apiShortUrl
		if err := c.BindJSON(&newLink); err != nil {
			logger.Error(ctx, "error_decoding_request", zap.Error(err))
			writeApiError(c, 400, ErrorInvalidRequest)
			return
		}

		ctx = logger.With(ctx, zap.String("email", newLink.LongUrl))

		result, err := service.ShortenUrl(ctx, model.Shortener{
			OriginalUrl: newLink.LongUrl,
		})
		if err != nil {
			logger.Error(ctx, "error_shortening_url", zap.Error(err))
			writeApiError(c, 400, err.Error())
			return
		}

		writeApiResult(c, 200, result)
	}
}

func writeApiError(c *gin.Context, code int, message string) {
	c.Status(code)
	c.Header("Content-Type", "application/json")

	err := json.NewEncoder(c.Writer).Encode(apiError{
		Code:    code,
		Message: message,
	})

	if err != nil {
		logger.Error(c.Request.Context(), "error returning result", zap.Error(err))
	}
}

func writeApiResult(c *gin.Context, code int, body interface{}) {
	c.Status(code)
	c.Header("Content-Type", "application/json")

	err := json.NewEncoder(c.Writer).Encode(apiResult{
		Code: code,
		Data: body,
	})

	if err != nil {
		logger.Error(c.Request.Context(), "error returning result", zap.Error(err))
	}
}
