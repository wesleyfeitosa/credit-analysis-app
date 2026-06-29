package service

import (
	"context"

	"creditanalysis/internal/model"
	"creditanalysis/internal/repository"
)

const defaultPageSize = 20
const maxPageSize = 100

// CreditAnalysisService implements listing and detail use cases.
type CreditAnalysisService struct {
	repo repository.CreditAnalysisRepository
}

// NewCreditAnalysisService builds a CreditAnalysisService.
func NewCreditAnalysisService(repo repository.CreditAnalysisRepository) *CreditAnalysisService {
	return &CreditAnalysisService{repo: repo}
}

// List normalizes pagination and delegates to the repository.
func (s *CreditAnalysisService) List(ctx context.Context, f model.ListFilter) (model.Page[model.CreditAnalysis], error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 || f.PageSize > maxPageSize {
		f.PageSize = defaultPageSize
	}
	return s.repo.List(ctx, f)
}

// Get returns a single analysis with its event history.
func (s *CreditAnalysisService) Get(ctx context.Context, id int64) (*model.CreditAnalysisDetail, error) {
	return s.repo.GetByID(ctx, id)
}
