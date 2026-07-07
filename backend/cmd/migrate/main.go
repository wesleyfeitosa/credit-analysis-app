// Command migrate applies the .sql files in a directory (default: migrations)
// against the database referenced by DATABASE_URL, in filename order.
//
// It runs SCHEMA migrations only by default — the deploy pipeline invokes it as
// `go run ./cmd/migrate`, and schema files are idempotent (IF NOT EXISTS) so
// re-running every push is safe. Seed data lives in ./seeds and is NOT applied
// automatically (re-running it would duplicate rows); apply it explicitly on a
// fresh database with `go run ./cmd/migrate -dir seeds`.
//
// The simple query protocol is used so that (a) multi-statement .sql files run
// in a single Exec and (b) it works through Supabase's transaction-mode pooler
// (port 6543), which does not support prepared statements.
package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"creditanalysis/internal/config"
)

func main() {
	dir := flag.String("dir", "migrations", "directory containing .sql migration files")
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()

	cfg := config.Load()

	connCfg, err := pgx.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("parse database url", zap.Error(err))
	}
	// Simple protocol: required for the transaction-mode pooler and lets each
	// file contain multiple semicolon-separated statements.
	connCfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	conn, err := pgx.ConnectConfig(ctx, connCfg)
	if err != nil {
		logger.Fatal("connect to database", zap.Error(err))
	}
	defer func() { _ = conn.Close(ctx) }()

	if err := conn.Ping(ctx); err != nil {
		logger.Fatal("ping database", zap.Error(err))
	}
	logger.Info("connected to database")

	files, err := filepath.Glob(filepath.Join(*dir, "*.sql"))
	if err != nil {
		logger.Fatal("glob migrations", zap.Error(err))
	}
	if len(files) == 0 {
		logger.Fatal("no migration files found", zap.String("dir", *dir))
	}
	sort.Strings(files)

	for _, f := range files {
		sql, err := os.ReadFile(f)
		if err != nil {
			logger.Fatal("read migration", zap.String("file", f), zap.Error(err))
		}
		if _, err := conn.Exec(ctx, string(sql)); err != nil {
			logger.Fatal("apply migration", zap.String("file", f), zap.Error(err))
		}
		logger.Info("applied migration", zap.String("file", filepath.Base(f)))
	}

	logger.Info("all migrations applied", zap.Int("count", len(files)))
}
