package handler

import (
	"github.com/gin-gonic/gin"

	"creditanalysis/internal/middleware"
)

// Register mounts every route on the given engine. Protected routes require a
// valid JWT via the Authorization: Bearer <token> header.
func (h *Handler) Register(r *gin.Engine, jwtSecret string) {
	r.GET("/health", h.Health)
	r.POST("/auth/login", h.Login)
	r.GET("/sdui/screens/login", h.LoginScreen)

	protected := r.Group("/")
	protected.Use(middleware.JWTAuth(jwtSecret))
	{
		protected.GET("/sdui/screens/credit-analyses", h.CreditAnalysesScreen)
		protected.GET("/credit-analyses", h.ListAnalyses)
		protected.GET("/credit-analyses/:id", h.GetAnalysis)
		protected.POST("/users/preferences/filters", h.SaveFilterPreferences)
	}
}
