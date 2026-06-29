package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"creditanalysis/internal/model"
	"creditanalysis/internal/repository"
)

// ListAnalyses returns a paginated, filtered and sorted list of analyses.
func (h *Handler) ListAnalyses(c *gin.Context) {
	f := model.ListFilter{
		Document:   c.Query("document"),
		ClientName: c.Query("clientName"),
		Status:     c.Query("status"),
		SortBy:     c.DefaultQuery("sortBy", "createdAt"),
		SortDir:    c.DefaultQuery("sortDir", "desc"),
	}
	f.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	f.PageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if n, ok := parseInt(c.Query("scoreMin")); ok {
		f.ScoreMin = &n
	}
	if n, ok := parseInt(c.Query("scoreMax")); ok {
		f.ScoreMax = &n
	}
	if t, ok := parseTime(c.Query("dateFrom")); ok {
		f.DateFrom = &t
	}
	if t, ok := parseTime(c.Query("dateTo")); ok {
		f.DateTo = &t
	}

	page, err := h.analyses.List(c.Request.Context(), f)
	if err != nil {
		h.log.Error("list analyses", zapError(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list analyses"})
		return
	}
	c.JSON(http.StatusOK, page)
}

// GetAnalysis returns a single analysis with its event history.
func (h *Handler) GetAnalysis(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	detail, err := h.analyses.Get(c.Request.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "analysis not found"})
		return
	}
	if err != nil {
		h.log.Error("get analysis", zapError(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load analysis"})
		return
	}
	c.JSON(http.StatusOK, detail)
}

func parseInt(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	n, err := strconv.Atoi(s)
	return n, err == nil
}

func parseTime(s string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, false
	}
	t, err := time.Parse(time.RFC3339, s)
	return t, err == nil
}
