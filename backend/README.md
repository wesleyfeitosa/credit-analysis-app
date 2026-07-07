# Backend — Credit Analysis API

Go API for the credit analysis system. Layered as **handler → service → repository**.

## Layout

```
cmd/api/main.go          # entrypoint: loads config, wires layers, starts gin
internal/
  config/                # env-based configuration
  model/                 # domain types (User, CreditAnalysis, events, filters)
  repository/            # data-access interfaces + PostgreSQL (pgx) implementation
  service/               # business logic (auth/JWT, listing, preferences)
  handler/               # gin HTTP handlers + route registration
  middleware/            # JWT auth middleware
  sdui/                  # Server Driven UI screen contracts
migrations/              # 0001_init.sql — schema (applied on every deploy)
seeds/                   # 0001_seed.sql — test data (fresh DB only, not deployed)
```

Dependencies flow inward: `handler` depends on `service`, `service` depends on
the `repository` interfaces. `repository.Postgres` is the only place that knows
about the database.

## Endpoints

| Method | Route                             | Auth |
|--------|-----------------------------------|------|
| POST   | `/auth/login`                     | no   |
| GET    | `/sdui/screens/login`             | no   |
| GET    | `/health`                         | no   |
| GET    | `/sdui/screens/credit-analyses`   | yes  |
| GET    | `/credit-analyses`                | yes  |
| GET    | `/credit-analyses/:id`            | yes  |
| POST   | `/users/preferences/filters`      | yes  |

Protected routes require `Authorization: Bearer <token>`.

`GET /credit-analyses` accepts query params: `document`, `clientName`, `status`,
`scoreMin`, `scoreMax`, `dateFrom`/`dateTo` (RFC3339), `page`, `pageSize`,
`sortBy` (`clientName|document|status|score|createdAt`), `sortDir` (`asc|desc`).

## Configuration

Copy `.env.example` and adjust. Variables: `PORT`, `DATABASE_URL`, `JWT_SECRET`.

## Test credentials

`admin@creditanalysis.com` / `senha123` (seeded by `seeds/0001_seed.sql`).

## Run locally

```bash
go run ./cmd/api          # requires a reachable PostgreSQL (see DATABASE_URL)
go test ./...             # unit tests
go build ./cmd/api        # build binary
```

## Migrations & seeds

Schema lives in `migrations/` (idempotent — `IF NOT EXISTS`); seed/test data
lives in `seeds/`. On `docker compose up`, both run once on first boot of an
empty data volume. To apply against any database via the built-in runner:

```bash
go run ./cmd/migrate               # schema only (what the deploy runs on every push)
go run ./cmd/migrate -dir seeds    # seed a FRESH database (re-running duplicates rows)
```

The deploy pipeline applies **schema only** — seeds are never run automatically,
so pushes don't accumulate duplicate test data.
