package service

import (
	"context"
	"math/rand/v2"
	"time"

	"creditanalysis/internal/model"
	"creditanalysis/internal/repository"
)

const defaultPageSize = 20
const maxPageSize = 100

// approvalThreshold is the minimum score (out of maxScore) for an analysis to
// be approved during the emulated decision.
const (
	approvalThreshold = 600
	maxScore          = 1000
)

type CreditAnalysisService struct {
	repo repository.CreditAnalysisRepository
}

func NewCreditAnalysisService(repo repository.CreditAnalysisRepository) *CreditAnalysisService {
	return &CreditAnalysisService{repo: repo}
}

func (s *CreditAnalysisService) List(ctx context.Context, f model.ListFilter) (model.Page[model.CreditAnalysis], error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 || f.PageSize > maxPageSize {
		f.PageSize = defaultPageSize
	}
	return s.repo.List(ctx, f)
}

func (s *CreditAnalysisService) Get(ctx context.Context, id int64) (*model.CreditAnalysisDetail, error) {
	return s.repo.GetByID(ctx, id)
}

// Create registers a new credit analysis and emulates its full lifecycle
// synchronously: the record starts as PENDENTE, moves to EM_ANALISE, then a
// random score drives the final APROVADO/REPROVADO decision. Every transition
// is recorded as an event, and the completed analysis (with its history) is
// returned.
func (s *CreditAnalysisService) Create(ctx context.Context, document, clientName string) (*model.CreditAnalysisDetail, error) {
	score := rand.IntN(maxScore + 1) // 0..maxScore inclusive
	finalStatus := model.StatusRejected
	decisionNote := "Crédito reprovado: score abaixo do mínimo"
	if score >= approvalThreshold {
		finalStatus = model.StatusApproved
		decisionNote = "Crédito aprovado"
	}

	// Timestamps progress by a minute per step so the event history has a
	// stable, readable ordering.
	t0 := time.Now().UTC()
	events := []model.CreditAnalysisEvent{
		{Status: model.StatusPending, Note: "Análise iniciada", CreatedAt: t0},
		{Status: model.StatusInReview, Note: "Documentação em verificação", CreatedAt: t0.Add(1 * time.Minute)},
		{Status: finalStatus, Note: decisionNote, CreatedAt: t0.Add(2 * time.Minute)},
	}

	analysis := model.CreditAnalysis{
		Document:   document,
		ClientName: clientName,
		Status:     finalStatus,
		Score:      score,
		CreatedAt:  t0,
	}

	id, err := s.repo.Create(ctx, analysis, events)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}
