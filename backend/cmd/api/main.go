// Command api is the entrypoint of the credit analysis backend.
package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"creditanalysis/internal/config"
	"creditanalysis/internal/handler"
	"creditanalysis/internal/repository"
	"creditanalysis/internal/service"
)

func main() {
	cfg := config.Load()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("connect to database", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatal("ping database", zap.Error(err))
	}

	// Wire layers: repository -> service -> handler.
	repo := repository.NewPostgres(pool)
	authSvc := service.NewAuthService(repo, cfg.JWTSecret)
	analysisSvc := service.NewCreditAnalysisService(repo)
	prefsSvc := service.NewPreferencesService(repo)
	h := handler.New(authSvc, analysisSvc, prefsSvc, logger)

	r := gin.New()
	r.Use(gin.Recovery())
	h.Register(r, cfg.JWTSecret)

	logger.Info("starting server", zap.String("port", cfg.Port))
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.Fatal("server stopped", zap.Error(err))
	}
}
