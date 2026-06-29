// Package repository defines data-access contracts and their PostgreSQL
// implementation. The service layer depends only on these interfaces.
package repository

import (
	"context"
	"errors"

	"creditanalysis/internal/model"
)

// ErrNotFound is returned when a requested record does not exist.
var ErrNotFound = errors.New("not found")

// UserRepository provides access to user credentials.
type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}

// CreditAnalysisRepository provides access to credit analyses.
type CreditAnalysisRepository interface {
	List(ctx context.Context, f model.ListFilter) (model.Page[model.CreditAnalysis], error)
	GetByID(ctx context.Context, id int64) (*model.CreditAnalysisDetail, error)
}

// PreferencesRepository persists per-user filter preferences.
type PreferencesRepository interface {
	SaveFilters(ctx context.Context, userID int64, filters []byte) error
	GetFilters(ctx context.Context, userID int64) ([]byte, error)
}
