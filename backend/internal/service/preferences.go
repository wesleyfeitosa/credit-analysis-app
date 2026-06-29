package service

import (
	"context"

	"creditanalysis/internal/repository"
)

// PreferencesService persists and restores per-user filter preferences.
type PreferencesService struct {
	repo repository.PreferencesRepository
}

// NewPreferencesService builds a PreferencesService.
func NewPreferencesService(repo repository.PreferencesRepository) *PreferencesService {
	return &PreferencesService{repo: repo}
}

// Save stores the raw filter JSON for the given user.
func (s *PreferencesService) Save(ctx context.Context, userID int64, filters []byte) error {
	return s.repo.SaveFilters(ctx, userID, filters)
}

// Get returns the stored filter JSON for the given user.
func (s *PreferencesService) Get(ctx context.Context, userID int64) ([]byte, error) {
	return s.repo.GetFilters(ctx, userID)
}
