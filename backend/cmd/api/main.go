// Command api is the entrypoint of the credit analysis backend.
package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "creditanalysis/docs" // generated OpenAPI docs (swag init)
	"creditanalysis/internal/config"
	"creditanalysis/internal/handler"
	"creditanalysis/internal/repository"
	"creditanalysis/internal/service"
)

// @title           Credit Analysis API
// @version         1.0
// @description     API do Sistema de Análises de Crédito (autenticação, listagem e SDUI).
// @host            localhost:8080
// @BasePath        /
//
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Informe o token como: "Bearer <token>"

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
	r.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
	h.Register(r, cfg.JWTSecret)

	logger.Info("starting server", zap.String("port", cfg.Port))
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.Fatal("server stopped", zap.Error(err))
	}
}
