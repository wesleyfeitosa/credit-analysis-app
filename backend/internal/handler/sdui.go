package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"creditanalysis/internal/middleware"
	"creditanalysis/internal/repository"
	"creditanalysis/internal/sdui"
)

// LoginScreen returns the SDUI contract for the login screen.
// @Summary      SDUI da tela de login
// @Description  Retorna o contrato SDUI que descreve a tela de login.
// @Tags         sdui
// @Produce      json
// @Success      200  {object}  sdui.Screen
// @Router       /sdui/screens/login [get]
func (h *Handler) LoginScreen(c *gin.Context) {
	c.JSON(http.StatusOK, sdui.LoginScreen())
}

// CreditAnalysesScreen returns the SDUI contract for the listing screen.
// @Summary      SDUI da tela de listagem
// @Description  Retorna o contrato SDUI que descreve a tela de listagem de análises.
// @Tags         sdui
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  sdui.Screen
// @Failure      401  {object}  map[string]string
// @Router       /sdui/screens/credit-analyses [get]
func (h *Handler) CreditAnalysesScreen(c *gin.Context) {
	userID := c.GetInt64(middleware.ContextUserID)

	// Pre-fill the screen with the user's saved filter preferences, if any. A
	// missing row is expected (first visit); any other read failure degrades
	// gracefully to an unfilled screen rather than failing the request.
	saved, err := h.prefs.Get(c.Request.Context(), userID)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFound) {
			h.log.Error("load filter preferences", zapError(err))
		}
		saved = nil
	}

	c.JSON(http.StatusOK, sdui.CreditAnalysesScreen(saved))
}
