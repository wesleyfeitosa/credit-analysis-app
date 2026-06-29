-- Seed data for local development and testing.
-- Test user credentials: admin@creditanalysis.com / senha123

INSERT INTO users (email, password_hash, name)
VALUES ('admin@creditanalysis.com', '$2a$10$iPT855tHUcpGaFNe4tytSOjgYagVHOx2qNJq.T2uGMN7LhXWhtd3m', 'Administrador')
ON CONFLICT (email) DO NOTHING;

INSERT INTO credit_analyses (document, client_name, status, score, created_at) VALUES
    ('123.456.789-00',     'Maria Silva',          'APROVADO',   820, now() - interval '1 day'),
    ('987.654.321-00',     'João Souza',           'REPROVADO',  340, now() - interval '2 days'),
    ('11.222.333/0001-44', 'Tech Solutions LTDA',  'EM_ANALISE', 610, now() - interval '3 days'),
    ('222.333.444-55',     'Ana Pereira',          'PENDENTE',   500, now() - interval '5 days'),
    ('44.555.666/0001-77', 'Comercial Boa Vista',  'APROVADO',   790, now() - interval '8 days')
ON CONFLICT DO NOTHING;

-- Event history for the first analysis.
INSERT INTO credit_analysis_events (analysis_id, status, note, created_at)
SELECT id, 'PENDENTE', 'Análise iniciada', created_at FROM credit_analyses WHERE client_name = 'Maria Silva'
UNION ALL
SELECT id, 'EM_ANALISE', 'Documentação em verificação', created_at + interval '2 hours' FROM credit_analyses WHERE client_name = 'Maria Silva'
UNION ALL
SELECT id, 'APROVADO', 'Crédito aprovado', created_at + interval '6 hours' FROM credit_analyses WHERE client_name = 'Maria Silva';
