package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"creditanalysis/internal/model"
	"creditanalysis/internal/service"
)

// fakeAnalysisRepo is an in-memory repository.CreditAnalysisRepository that
// records what Create persisted so the emulation logic can be asserted in
// isolation from PostgreSQL.
type fakeAnalysisRepo struct {
	lastAnalysis model.CreditAnalysis
	lastEvents   []model.CreditAnalysisEvent
}

func (f *fakeAnalysisRepo) List(context.Context, model.ListFilter) (model.Page[model.CreditAnalysis], error) {
	return model.Page[model.CreditAnalysis]{}, nil
}

func (f *fakeAnalysisRepo) Create(_ context.Context, a model.CreditAnalysis, events []model.CreditAnalysisEvent) (int64, error) {
	f.lastAnalysis = a
	f.lastEvents = events
	return 42, nil
}

func (f *fakeAnalysisRepo) GetByID(_ context.Context, id int64) (*model.CreditAnalysisDetail, error) {
	return &model.CreditAnalysisDetail{
		CreditAnalysis: f.lastAnalysis,
		Events:         f.lastEvents,
	}, nil
}

func TestCreditAnalysisService_Create(t *testing.T) {
	repo := &fakeAnalysisRepo{}
	svc := service.NewCreditAnalysisService(repo)

	// Run many iterations: the score is random, so assert invariants that must
	// hold regardless of the drawn value.
	for i := 0; i < 200; i++ {
		detail, err := svc.Create(context.Background(), "123.456.789-00", "Maria Silva")
		require.NoError(t, err)
		require.NotNil(t, detail)

		require.Equal(t, "123.456.789-00", detail.Document)
		require.Equal(t, "Maria Silva", detail.ClientName)
		require.GreaterOrEqual(t, detail.Score, 0)
		require.LessOrEqual(t, detail.Score, 1000)

		// Final status is consistent with the score threshold (600).
		if detail.Score >= 600 {
			require.Equal(t, model.StatusApproved, detail.Status)
		} else {
			require.Equal(t, model.StatusRejected, detail.Status)
		}

		// The emulated lifecycle always produces three ordered events ending in
		// the same decision as the analysis.
		require.Len(t, detail.Events, 3)
		require.Equal(t, model.StatusPending, detail.Events[0].Status)
		require.Equal(t, model.StatusInReview, detail.Events[1].Status)
		require.Equal(t, detail.Status, detail.Events[2].Status)
		require.True(t, detail.Events[0].CreatedAt.Before(detail.Events[1].CreatedAt))
		require.True(t, detail.Events[1].CreatedAt.Before(detail.Events[2].CreatedAt))
	}
}
