package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"creditanalysis/internal/sdui"
)

// LoginScreen returns the SDUI contract for the login screen.
func (h *Handler) LoginScreen(c *gin.Context) {
	c.JSON(http.StatusOK, sdui.LoginScreen())
}

// CreditAnalysesScreen returns the SDUI contract for the listing screen.
func (h *Handler) CreditAnalysesScreen(c *gin.Context) {
	c.JSON(http.StatusOK, sdui.CreditAnalysesScreen())
}
