package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health reports service liveness.
// @Summary      Health check
// @Description  Verifica se o serviço está no ar.
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
