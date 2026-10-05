# Nōto
A simple notes taking app made as a reference, showing how to structure a go codebase, build and test it
while obeying the sonar cloud Quality gates.

## Status

- **Database: SQLite in production code today** (`notes-backend/internal/db/db.go`, file `Notes.DB`, via GORM).
  **PostgreSQL is the target** going forward — a local Postgres is fully wired up (Podman, see below) and the
  full schema exists as runnable migrations (`notes-backend/migrations/`), but the Go code has **not** been
  cut over yet: it still opens the SQLite file directly. See `roadmap/noto-db-schema.md` for the design and
  rationale (multi-user teams + concurrent writes + native partial unique indexes + full-text search are why
  Postgres was chosen over SQLite).
- The REST API surface is specified in `notes-backend/openAPI-spec/` (`noto-openapi.yaml` / `noto-api-contract.md`);
  implemented endpoints currently cover notes and users (see `notes-backend/server/`, `internal/notes`, `internal/user`).
  Teams, labels, sprints, epics, tasks and the secrets vault described in the schema/roadmap docs are **designed,
  not yet implemented**.
- `roadmap/` holds design docs that are ahead of the code: `noto-db-schema.md` (DB schema + Postgres migration
  rationale), `noto-ui-ux.md` / `google-keep-ui.md` (UI/UX references).

## Repo layout

```
notes-app/
├── compose.yaml              # podman/docker compose: postgres + pgadmin + backend + frontend (dev, hot-reload)
├── Makefile                  # make dev / backend / frontend / migrate-*
├── notes-backend/            # Go API (Fiber + GORM)
│   ├── cmd/main.go
│   ├── internal/             # db, notes, user packages
│   ├── server/                # HTTP routes, auth middleware
│   ├── migrations/           # golang-migrate SQL migrations (Postgres target schema)
│   ├── openAPI-spec/         # API contract (yaml + markdown)
│   ├── Dockerfile            # production image (SQLite build)
│   └── Dockerfile.dev         # dev image: air live-reload, used by compose.yaml
├── notes-ui/                  # React + Vite + TanStack Router frontend
│   ├── src/
│   └── Dockerfile.dev          # dev image: bun + vite dev server, used by compose.yaml
├── pgadmin/servers.json       # pre-registers the local Postgres connection in pgAdmin
├── roadmap/                   # design docs ahead of current implementation
└── mockup/                    # UI mockups
```

## Development

### Dependencies
* Backend
    * Go [1.23.2 or higher](https://go.dev/doc/)
    * fiber [http server](https://gofiber.io/)
    * sqlite3 [DB — current runtime store](https://www.sqlite.org/index.html)
    * PostgreSQL [DB — target store, schema ready, not yet wired into the Go code](https://www.postgresql.org/)
    * gorm [db ORM](https://gorm.io/)
    * [air](https://github.com/air-verse/air) (live reload)
    * [golang-migrate](https://github.com/golang-migrate/migrate) (DB migrations, Postgres schema only for now)
* Frontend
    * Typescript v5 or higher
    * Bun
    * React with [TanStack Router](https://tanstack.com/router)
    * [Vite](https://vite.dev/)
    * [Tailwind CSS](https://tailwindcss.com/) v4
* Local infra (optional, see "Containerized dev" below)
    * [Podman](https://podman.io/) (or Docker) + Compose

### Setup

#### Making RSA CERTs required for the JWTs
```bash
mkdir secrets
cd secrets
openssl genrsa -out private_key.pem 2048
openssl rsa -in private_key.pem -outform PEM -pubout -out public_key.pem.pub
```

#### Install dependencies
```bash
make install
```

### Running natively (bare-metal)

#### Start both backend and frontend (recommended)
```bash
make dev
```

#### Start backend only (air live reload on :8001)
```bash
make backend
```

#### Start frontend only (vite dev server on :5173)
```bash
make frontend
```

The frontend proxies `/user` and `/notes` to the backend. The proxy target is
configurable via `VITE_BACKEND_PROXY_TARGET` (see `notes-ui/vite.config.ts`); it
defaults to `http://localhost:8001`, which is correct for native, non-containerized
dev — you only need to set it when running the frontend inside a container (the
compose setup below does this for you).

#### Stop all services
```bash
make stop
```

#### Build backend binary
```bash
make build-backend
```

### Containerized dev (Podman Compose)

`compose.yaml` at the repo root runs the whole stack — Postgres, a Postgres admin UI,
the Go backend, and the React frontend — with live-reload, so editing source on the
host rebuilds/restarts the right thing automatically without a manual image rebuild:

```bash
podman compose up -d --build
```

| Service | URL | Notes |
|---|---|---|
| backend | http://localhost:8001 | Go binary rebuilt + restarted by `air` on every `.go` save (bind-mounted source) |
| frontend | http://localhost:5173 | Vite dev server, HMR on save (bind-mounted source) |
| postgres | localhost:5432 | db `noto`, user `noto`, password `noto_dev_password` |
| pgadmin | http://localhost:5050 | login `admin@noto.dev` / `admin`; the `noto` server connection is pre-registered (password: `noto_dev_password`) |

Editing `.go`/`.tsx`/etc. source is picked up live via bind mounts (`air` for Go,
Vite HMR for the frontend) — no rebuild needed. Only dependency-manifest changes
(`go.mod`/`go.sum`, `package.json`/`bun.lock`) need an image rebuild, which
`develop.watch` in `compose.yaml` triggers automatically if you run
`podman compose watch` alongside `up`; otherwise just re-run `up --build`.

Notes for Fedora/SELinux-enforcing hosts: the bind mounts in `compose.yaml` use the
`:Z` relabel flag already — if you fork this onto another SELinux host and see
"permission denied" from `air` or Vite, that flag is why it's there.

Stop everything with `podman compose down` (add `-v` to also drop the named
volumes — this deletes the local Postgres/pgAdmin data).

### Database & migrations

The authoritative schema design lives in `roadmap/noto-db-schema.md`. The runnable
version of that schema is in `notes-backend/migrations/` as
[golang-migrate](https://github.com/golang-migrate/migrate) SQL files (chosen over
Liquibase: single static binary, no JVM, plain versioned SQL instead of an XML/YAML
changelog DSL — a better fit for a Go backend).

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.1

make migrate-up       # apply all pending migrations (targets the compose postgres by default)
make migrate-down     # roll back one
make migrate-version  # print current version
```

`DB_URL` in the `Makefile` defaults to the local compose Postgres
(`postgres://noto:noto_dev_password@localhost:5432/noto?sslmode=disable`); override it
to target a different database. See `notes-backend/migrations/README.md` for the
table-by-table dependency order and how to add a new migration.

**This migration set is not yet consumed by the running app.** `internal/db/db.go`
still opens `Notes.DB` via `gorm.io/driver/sqlite` directly; cutting the backend over
to Postgres (reading a connection string, swapping the GORM driver, running migrations
on boot instead of `AutoMigrate`) is the next piece of work implied by the roadmap doc.

### CI / quality

Sonar Cloud quality gates run against this repo (`sonar-project.properties`); Go
tests report coverage via `notes-backend/coverage.out`. `lefthook.yml` configures
pre-commit hooks including `detect-secrets`.
