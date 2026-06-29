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
migrations/              # 0001_init.sql (schema) + 0002_seed.sql (test data)
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

`admin@creditanalysis.com` / `senha123` (seeded by `migrations/0002_seed.sql`).

## Run locally

```bash
go run ./cmd/api          # requires a reachable PostgreSQL (see DATABASE_URL)
go test ./...             # unit tests
go build ./cmd/api        # build binary
```

## Migrations

Two SQL files in `migrations/`. They run automatically when the postgres
container starts (mounted into `/docker-entrypoint-initdb.d`). To apply
manually against a running database:

```bash
psql "$DATABASE_URL" -f migrations/0001_init.sql
psql "$DATABASE_URL" -f migrations/0002_seed.sql
```
