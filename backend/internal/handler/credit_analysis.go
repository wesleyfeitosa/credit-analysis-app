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
// @Summary      Listar análises
// @Description  Lista análises de crédito com filtros, paginação e ordenação.
// @Tags         credit-analyses
// @Produce      json
// @Security     BearerAuth
// @Param        document    query     string  false  "Filtra por CPF/CNPJ"
// @Param        clientName  query     string  false  "Filtra por nome do cliente"
// @Param        status      query     string  false  "Filtra por status"  Enums(APROVADO, REPROVADO, EM_ANALISE, PENDENTE)
// @Param        scoreMin    query     int     false  "Score mínimo"
// @Param        scoreMax    query     int     false  "Score máximo"
// @Param        dateFrom    query     string  false  "Data inicial (RFC3339)"
// @Param        dateTo      query     string  false  "Data final (RFC3339)"
// @Param        page        query     int     false  "Página (default 1)"
// @Param        pageSize    query     int     false  "Itens por página (default 20)"
// @Param        sortBy      query     string  false  "Campo de ordenação (default createdAt)"
// @Param        sortDir     query     string  false  "Direção"  Enums(asc, desc)
// @Success      200         {object}  model.Page[model.CreditAnalysis]
// @Failure      401         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /credit-analyses [get]
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

type createAnalysisRequest struct {
	Document   string `json:"document" validate:"required,min=11"`
	ClientName string `json:"clientName" validate:"required,min=2"`
}

// CreateAnalysis registers a new analysis and emulates its lifecycle.
// @Summary      Criar análise
// @Description  Cria uma nova análise de crédito e emula o processo de análise (PENDENTE → EM_ANALISE → APROVADO/REPROVADO), retornando o resultado final com o histórico de eventos.
// @Tags         credit-analyses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        analysis  body      createAnalysisRequest  true  "Dados do cliente"
// @Success      201       {object}  model.CreditAnalysisDetail
// @Failure      400       {object}  map[string]string
// @Failure      401       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /credit-analyses [post]
func (h *Handler) CreateAnalysis(c *gin.Context) {
	var req createAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	detail, err := h.analyses.Create(c.Request.Context(), req.Document, req.ClientName)
	if err != nil {
		h.log.Error("create analysis", zapError(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create analysis"})
		return
	}
	c.JSON(http.StatusCreated, detail)
}

// GetAnalysis returns a single analysis with its event history.
// @Summary      Detalhe da análise
// @Description  Retorna uma análise de crédito com seu histórico de eventos.
// @Tags         credit-analyses
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID da análise"
// @Success      200  {object}  model.CreditAnalysisDetail
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /credit-analyses/{id} [get]
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
