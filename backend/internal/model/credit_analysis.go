package model

import "time"

// AnalysisStatus is the lifecycle status of a credit analysis.
type AnalysisStatus string

const (
	StatusApproved AnalysisStatus = "APROVADO"
	StatusRejected AnalysisStatus = "REPROVADO"
	StatusInReview AnalysisStatus = "EM_ANALISE"
	StatusPending  AnalysisStatus = "PENDENTE"
)

// CreditAnalysis is a single credit analysis record.
type CreditAnalysis struct {
	ID         int64          `json:"id"`
	Document   string         `json:"document"`
	ClientName string         `json:"clientName"`
	Status     AnalysisStatus `json:"status"`
	Score      int            `json:"score"`
	CreatedAt  time.Time      `json:"createdAt"`
}

// CreditAnalysisEvent is an entry in the status/history timeline of an analysis.
type CreditAnalysisEvent struct {
	ID         int64          `json:"id"`
	AnalysisID int64          `json:"analysisId"`
	Status     AnalysisStatus `json:"status"`
	Note       string         `json:"note"`
	CreatedAt  time.Time      `json:"createdAt"`
}

// CreditAnalysisDetail is an analysis together with its event history.
type CreditAnalysisDetail struct {
	CreditAnalysis
	Events []CreditAnalysisEvent `json:"events"`
}

// ListFilter holds the supported query parameters for listing analyses.
type ListFilter struct {
	Document   string
	ClientName string
	Status     string
	ScoreMin   *int
	ScoreMax   *int
	DateFrom   *time.Time
	DateTo     *time.Time

	Page     int
	PageSize int
	SortBy   string
	SortDir  string // asc | desc
}

// Page is a generic paginated result envelope.
type Page[T any] struct {
	Items    []T   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}
