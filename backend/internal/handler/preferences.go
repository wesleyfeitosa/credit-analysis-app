package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"creditanalysis/internal/middleware"
)

// SaveFilterPreferences stores the authenticated user's filter preferences.
// The body is an arbitrary JSON object describing the filter state.
// @Summary      Salvar preferências de filtros
// @Description  Salva as preferências de filtros do usuário autenticado (JSON arbitrário).
// @Tags         preferences
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        filters  body      object  true  "Estado dos filtros"
// @Success      204      "No Content"
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /users/preferences/filters [post]
func (h *Handler) SaveFilterPreferences(c *gin.Context) {
	userID := c.GetInt64(middleware.ContextUserID)

	body, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
		return
	}
	if !json.Valid(body) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "body must be valid JSON"})
		return
	}

	if err := h.prefs.Save(c.Request.Context(), userID, body); err != nil {
		h.log.Error("save preferences", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save preferences"})
		return
	}
	c.Status(http.StatusNoContent)
}

// zapError is a small helper so other handlers can log errors without importing
// zap directly.
func zapError(err error) zap.Field {
	return zap.Error(err)
}
