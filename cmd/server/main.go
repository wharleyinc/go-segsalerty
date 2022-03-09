package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-segsalerty/common/logger"
	"go-segsalerty/internal/config"
	shorthy "go-segsalerty/internal/domain/shortener"
	httprequestadapter "go-segsalerty/internal/inbound/http-request-adapter"
	shortenermongodatabaseadapter "go-segsalerty/internal/outbound/shortener-mongo-database-adapter"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()
	router := gin.Default()

	// Used by shortener services
	appMongoDb, err := shortenermongodatabaseadapter.NewDatabaseAdapter(cfg)
	if err != nil {
		logger.Error(ctx, "error_starting_mongo_auth", zap.Error(err))
		return
	}

	shortenerService := shorthy.NewShortenerService(appMongoDb)

	handler := httprequestadapter.NewHttpHandler(shortenerService)
	handler.ApplyRoutes(router)

	router.Run(":" + cfg.Port)
}
