package service

import (
	"context"

	"creditanalysis/internal/repository"
)

type PreferencesService struct {
	repo repository.PreferencesRepository
}

func NewPreferencesService(repo repository.PreferencesRepository) *PreferencesService {
	return &PreferencesService{repo: repo}
}

func (s *PreferencesService) Save(ctx context.Context, userID int64, filters []byte) error {
	return s.repo.SaveFilters(ctx, userID, filters)
}

func (s *PreferencesService) Get(ctx context.Context, userID int64) ([]byte, error) {
	return s.repo.GetFilters(ctx, userID)
}
