-- Schema for the credit analysis system.
-- Applied automatically on first start via the postgres docker entrypoint
-- (files in /docker-entrypoint-initdb.d run in alphabetical order).

CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    name          TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS credit_analyses (
    id          BIGSERIAL PRIMARY KEY,
    document    TEXT NOT NULL,
    client_name TEXT NOT NULL,
    status      TEXT NOT NULL CHECK (status IN ('APROVADO', 'REPROVADO', 'EM_ANALISE', 'PENDENTE')),
    score       INT  NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_credit_analyses_status     ON credit_analyses (status);
CREATE INDEX IF NOT EXISTS idx_credit_analyses_created_at ON credit_analyses (created_at);

CREATE TABLE IF NOT EXISTS credit_analysis_events (
    id          BIGSERIAL PRIMARY KEY,
    analysis_id BIGINT NOT NULL REFERENCES credit_analyses (id) ON DELETE CASCADE,
    status      TEXT NOT NULL,
    note        TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_events_analysis_id ON credit_analysis_events (analysis_id);

CREATE TABLE IF NOT EXISTS user_filter_preferences (
    user_id    BIGINT PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    filters    JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
