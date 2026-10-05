# Migrations

Schema migrations for the Nōto Postgres database (see `roadmap/noto-db-schema.md`
for the full design). Tool: [golang-migrate](https://github.com/golang-migrate/migrate) —
a free, single-binary migration tool. Chosen over Liquibase because it needs no JVM,
uses plain versioned SQL (`.up.sql` / `.down.sql`) instead of XML/YAML changelogs, and
is the de facto standard for Go backends.

## Install the CLI

```sh
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.1
```

## Run

Local dev Postgres (see `compose.yaml` at the repo root) listens on `localhost:5432`
with db `noto`, user `noto`, password `noto_dev_password`:

```sh
export DB_URL="postgres://noto:noto_dev_password@localhost:5432/noto?sslmode=disable"

migrate -path notes-backend/migrations -database "$DB_URL" up      # apply all pending
migrate -path notes-backend/migrations -database "$DB_URL" down 1  # roll back one step
migrate -path notes-backend/migrations -database "$DB_URL" version # current version
```

Or via the root `Makefile`: `make migrate-up` / `make migrate-down` / `make migrate-version`.

## Adding a new migration

```sh
migrate create -ext sql -dir notes-backend/migrations -seq <description>
```

This creates a new `NNNNNN_<description>.up.sql` / `.down.sql` pair. Keep each
migration focused on one table or one change, and always write the matching
`.down.sql` — `migrate down` replays them in reverse version order, so FK-dependent
tables (e.g. `task_activity` → `tasks`) must be numbered after what they reference.

## Table order (matches §3 of the design doc)

| # | table | depends on |
|---|---|---|
| 1 | `users` | — |
| 2 | `refresh_tokens` | users |
| 3 | `vault_tokens` | users |
| 4 | `labels` | users |
| 5 | `notes` | users, labels |
| 6 | `teams` | users |
| 7 | `team_members` | teams, users |
| 8 | `sprints` | — (scope_id is polymorphic: user or team id, no FK) |
| 9 | `epics` | users, sprints |
| 10 | `tasks` | users, epics, sprints |
| 11 | `task_counters` | — |
| 12 | `task_activity` | tasks, users |

`gen_random_uuid()` is built into Postgres ≥ 13 — no `CREATE EXTENSION` needed.
