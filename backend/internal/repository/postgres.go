package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"creditanalysis/internal/model"
)

// Postgres implements every repository interface backed by a pgx pool.
type Postgres struct {
	pool *pgxpool.Pool
}

// Compile-time guarantees that Postgres satisfies the repository contracts.
var (
	_ UserRepository           = (*Postgres)(nil)
	_ CreditAnalysisRepository = (*Postgres)(nil)
	_ PreferencesRepository    = (*Postgres)(nil)
)

// NewPostgres builds a Postgres repository from a connection pool.
func NewPostgres(pool *pgxpool.Pool) *Postgres {
	return &Postgres{pool: pool}
}

func (p *Postgres) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := p.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, name, created_at FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (p *Postgres) List(ctx context.Context, f model.ListFilter) (model.Page[model.CreditAnalysis], error) {
	var (
		conds []string
		args  []any
	)
	add := func(cond string, val any) {
		args = append(args, val)
		conds = append(conds, fmt.Sprintf(cond, len(args)))
	}

	if f.Document != "" {
		add("document ILIKE '%%' || $%d || '%%'", f.Document)
	}
	if f.ClientName != "" {
		add("client_name ILIKE '%%' || $%d || '%%'", f.ClientName)
	}
	if f.Status != "" {
		add("status = $%d", f.Status)
	}
	if f.ScoreMin != nil {
		add("score >= $%d", *f.ScoreMin)
	}
	if f.ScoreMax != nil {
		add("score <= $%d", *f.ScoreMax)
	}
	if f.DateFrom != nil {
		add("created_at >= $%d", *f.DateFrom)
	}
	if f.DateTo != nil {
		add("created_at <= $%d", *f.DateTo)
	}

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	var total int64
	if err := p.pool.QueryRow(ctx, "SELECT count(*) FROM credit_analyses "+where, args...).Scan(&total); err != nil {
		return model.Page[model.CreditAnalysis]{}, err
	}

	sortDir := "ASC"
	if strings.EqualFold(f.SortDir, "desc") {
		sortDir = "DESC"
	}

	args = append(args, f.PageSize, (f.Page-1)*f.PageSize)
	listSQL := fmt.Sprintf(
		`SELECT id, document, client_name, status, score, created_at
		 FROM credit_analyses %s
		 ORDER BY %s %s
		 LIMIT $%d OFFSET $%d`,
		where, sortColumn(f.SortBy), sortDir, len(args)-1, len(args),
	)

	rows, err := p.pool.Query(ctx, listSQL, args...)
	if err != nil {
		return model.Page[model.CreditAnalysis]{}, err
	}
	defer rows.Close()

	items := make([]model.CreditAnalysis, 0, f.PageSize)
	for rows.Next() {
		var a model.CreditAnalysis
		if err := rows.Scan(&a.ID, &a.Document, &a.ClientName, &a.Status, &a.Score, &a.CreatedAt); err != nil {
			return model.Page[model.CreditAnalysis]{}, err
		}
		items = append(items, a)
	}
	if err := rows.Err(); err != nil {
		return model.Page[model.CreditAnalysis]{}, err
	}

	return model.Page[model.CreditAnalysis]{
		Items:    items,
		Total:    total,
		Page:     f.Page,
		PageSize: f.PageSize,
	}, nil
}

func (p *Postgres) GetByID(ctx context.Context, id int64) (*model.CreditAnalysisDetail, error) {
	var d model.CreditAnalysisDetail
	err := p.pool.QueryRow(ctx,
		`SELECT id, document, client_name, status, score, created_at
		 FROM credit_analyses WHERE id = $1`, id,
	).Scan(&d.ID, &d.Document, &d.ClientName, &d.Status, &d.Score, &d.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	rows, err := p.pool.Query(ctx,
		`SELECT id, analysis_id, status, note, created_at
		 FROM credit_analysis_events WHERE analysis_id = $1 ORDER BY created_at ASC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	d.Events = []model.CreditAnalysisEvent{}
	for rows.Next() {
		var e model.CreditAnalysisEvent
		if err := rows.Scan(&e.ID, &e.AnalysisID, &e.Status, &e.Note, &e.CreatedAt); err != nil {
			return nil, err
		}
		d.Events = append(d.Events, e)
	}
	return &d, rows.Err()
}

func (p *Postgres) SaveFilters(ctx context.Context, userID int64, filters []byte) error {
	_, err := p.pool.Exec(ctx,
		`INSERT INTO user_filter_preferences (user_id, filters, updated_at)
		 VALUES ($1, $2, now())
		 ON CONFLICT (user_id) DO UPDATE SET filters = EXCLUDED.filters, updated_at = now()`,
		userID, filters)
	return err
}

func (p *Postgres) GetFilters(ctx context.Context, userID int64) ([]byte, error) {
	var filters []byte
	err := p.pool.QueryRow(ctx,
		`SELECT filters FROM user_filter_preferences WHERE user_id = $1`, userID,
	).Scan(&filters)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return filters, err
}

// sortColumn maps a client-facing sort key to a whitelisted column name,
// preventing SQL injection through the ORDER BY clause.
func sortColumn(key string) string {
	switch key {
	case "clientName":
		return "client_name"
	case "document":
		return "document"
	case "status":
		return "status"
	case "score":
		return "score"
	default:
		return "created_at"
	}
}
