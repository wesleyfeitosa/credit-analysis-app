// Package handler holds the HTTP layer: it parses requests, delegates to the
// service layer and writes responses. It contains no business logic.
package handler

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"creditanalysis/internal/service"
)

// Handler wires the service layer to the HTTP routes.
type Handler struct {
	auth     *service.AuthService
	analyses *service.CreditAnalysisService
	prefs    *service.PreferencesService
	validate *validator.Validate
	log      *zap.Logger
}

// New builds a Handler.
func New(
	auth *service.AuthService,
	analyses *service.CreditAnalysisService,
	prefs *service.PreferencesService,
	log *zap.Logger,
) *Handler {
	return &Handler{
		auth:     auth,
		analyses: analyses,
		prefs:    prefs,
		validate: validator.New(validator.WithRequiredStructEnabled()),
		log:      log,
	}
}
