# Nōto — Database Schema & Solution Design

> Companion to `noto-openapi.yaml` / `noto-api-contract.md`. Target: **PostgreSQL via GORM**. Multi-user teams with concurrent writes (shared tasks/epics/sprints, team boards) and native support for partial unique indexes and full-text search make Postgres the primary target rather than SQLite. Drop in `Noto/notes-backend/docs/`.

## 1. Principles

- **One row per real thing the UI shows.** Every entity in the mockup (note, label, team, member, task, activity, epic, sprint, session, vault token) is a table; derived numbers (open counts, progress) are computed, not stored.
- **Two credentials, two secrets.** Login password and vault master password are hashed separately; the master password additionally derives the key that encrypts secret-note passwords. Compromising the login DB does not expose secrets.
- **Scope, not ownership.** Tasks, epics and sprints carry `scope_type` + `scope_id` (`personal`/user or `team`/team). One access rule covers all three.
- **Soft lifecycle for notes only.** `archived_at` / `trashed_at` timestamps (nullable) instead of booleans — free ordering and a future "auto-empty trash after 30 days" for nothing.
- **Native UUID + timestamptz.** `id` columns are `uuid` (app-generated or `gen_random_uuid()`, built into Postgres ≥ 13 — no extension needed). All timestamps are `timestamptz`, stored UTC.

## 2. Entity overview

```
users ─┬─< notes >─── labels
       ├─< refresh_tokens
       ├─< vault_tokens
       ├─< team_members >─── teams
       └─(scope personal)──┐
                            ├─< sprints
teams ─(scope team)────────┼─< epics ──┐
                            └─< tasks ──┼─< task_activity
                                        └── assignee → users
```

## 3. Tables

### 3.1 `users`
| column | type | notes |
|---|---|---|
| id | UUID PK DEFAULT gen_random_uuid() | |
| name | TEXT NOT NULL | avatar initials derived in API |
| email | TEXT NOT NULL UNIQUE | lowercased; used by `/directory` search |
| password_hash | TEXT NOT NULL | argon2id — **login** |
| master_hash | TEXT NULL | argon2id — **vault** verifier; null until first unlock setup |
| master_salt | BYTEA NULL | salt for KDF → vault key (§5) |
| wrapped_dek | BYTEA NULL | data-encryption key, wrapped by the master-derived key |
| default_role | TEXT NULL | pre-fills "Role" in member autocomplete |
| pref_default_tab | TEXT NOT NULL DEFAULT 'notes' | notes/tasks/teams |
| pref_notes_layout | TEXT NOT NULL DEFAULT 'masonry' | masonry/grid/list |
| pref_mask_secrets | BOOLEAN NOT NULL DEFAULT true | |
| pref_autolock_min | INTEGER NOT NULL DEFAULT 5 | 0 = never |
| created_at, updated_at | TIMESTAMPTZ NOT NULL DEFAULT now() | |

Index: `UNIQUE(email)`; for case-insensitive autocomplete `CREATE INDEX ON users (LOWER(name))`.

### 3.2 `refresh_tokens`
| column | type | notes |
|---|---|---|
| id | UUID PK DEFAULT gen_random_uuid() | |
| user_id | UUID FK users ON DELETE CASCADE | |
| token_hash | TEXT NOT NULL UNIQUE | sha256 of the opaque token |
| expires_at | TIMESTAMPTZ NOT NULL | |
| revoked_at | TIMESTAMPTZ NULL | `/auth/logout` sets this |
| created_at | TIMESTAMPTZ NOT NULL DEFAULT now() | |

### 3.3 `vault_tokens`
| column | type | notes |
|---|---|---|
| id | UUID PK DEFAULT gen_random_uuid() | |
| user_id | UUID FK users ON DELETE CASCADE | |
| token_hash | TEXT NOT NULL UNIQUE | |
| dek_cache | BYTEA NOT NULL | DEK re-wrapped with a server KEK for the token's lifetime (§5) |
| expires_at | TIMESTAMPTZ NOT NULL | = now + `pref_autolock_min` |
| created_at | TIMESTAMPTZ NOT NULL DEFAULT now() | |

`POST /vault/lock` and `PUT /me/password` delete all rows for the user.

### 3.4 `labels`
| column | type | notes |
|---|---|---|
| id | UUID PK DEFAULT gen_random_uuid() | |
| user_id | UUID FK users ON DELETE CASCADE | labels are per-user |
| name | TEXT NOT NULL | |
| color | TEXT NOT NULL | hex from chart palette |
| created_at | TIMESTAMPTZ NOT NULL DEFAULT now() | |

`UNIQUE(user_id, LOWER(name))` as a functional unique index → API `409`.

### 3.5 `notes`
| column | type | notes |
|---|---|---|
| id | UUID PK DEFAULT gen_random_uuid() | |
| user_id | UUID FK users ON DELETE CASCADE | |
| label_id | UUID NULL FK labels ON DELETE SET NULL | label delete clears, keeps note |
| is_secret | BOOLEAN NOT NULL DEFAULT false | NOTE ⇄ SECRET toggle |
| is_private | BOOLEAN NOT NULL DEFAULT false | lockable |
| title | TEXT NOT NULL DEFAULT '' | |
| body | TEXT NOT NULL DEFAULT '' | plain notes |
| site | TEXT NOT NULL DEFAULT '' | secrets |
| username | TEXT NOT NULL DEFAULT '' | secrets |
| password_enc | BYTEA NULL | AES-256-GCM(nonce‖ct‖tag) under DEK — **never plaintext** |
| color | TEXT NOT NULL DEFAULT 'default' | default/blue/teal/violet/rose/green |
| pinned | BOOLEAN NOT NULL DEFAULT false | |
| archived_at | TIMESTAMPTZ NULL | |
| trashed_at | TIMESTAMPTZ NULL | |
| created_at, updated_at | TIMESTAMPTZ NOT NULL DEFAULT now() | |

Indexes: `(user_id, trashed_at, archived_at, pinned DESC, updated_at DESC)` for the grid; `(user_id, label_id)`.
Search (`q`): generated `tsvector` column `search_tsv` over `title, body, site, username` with a `GIN` index (`CREATE INDEX ON notes USING GIN (search_tsv)`), kept current via `GENERATED ALWAYS AS (...) STORED` — no triggers needed. `password_enc` is never indexed or searched.

Invariants (enforce in service layer): archive/trash unpin; `is_private` body/site/username are returned empty unless a valid vault token is present (`locked: true`).

### 3.6 `teams`
| column | type | notes |
|---|---|---|
| id | UUID PK DEFAULT gen_random_uuid() | |
| name | TEXT NOT NULL | |
| color | TEXT NOT NULL | |
| created_by | UUID FK users | |
| created_at, updated_at | TIMESTAMPTZ NOT NULL DEFAULT now() | |

### 3.7 `team_members`
| column | type | notes |
|---|---|---|
| team_id | UUID FK teams ON DELETE CASCADE | |
| user_id | UUID FK users ON DELETE CASCADE | |
| role | TEXT NOT NULL DEFAULT 'Member' | free text: Design, SRE, Owner… |
| is_admin | BOOLEAN NOT NULL DEFAULT false | gates management modal, ADMIN badge |
| joined_at | TIMESTAMPTZ NOT NULL DEFAULT now() | |

`PRIMARY KEY(team_id, user_id)`; `INDEX(user_id)` for "my teams". Rule: a team must always have ≥1 admin (`409` on removing the last).

### 3.8 `sprints`
| column | type | notes |
|---|---|---|
| id | UUID PK DEFAULT gen_random_uuid() | |
| scope_type | TEXT NOT NULL | `personal` \| `team` |
| scope_id | UUID NOT NULL | user_id or team_id |
| name | TEXT NOT NULL | "Sprint 42" |
| start_date | DATE NOT NULL | |
| end_date | DATE NOT NULL | |
| is_current | BOOLEAN NOT NULL DEFAULT false | |
| created_at | TIMESTAMPTZ NOT NULL DEFAULT now() | |

`UNIQUE(scope_type, scope_id, name)`; native partial unique index `CREATE UNIQUE INDEX ON sprints (scope_type, scope_id) WHERE is_current` → exactly one current per scope, enforced by Postgres directly (no trigger).

### 3.9 `epics`
| column | type | notes |
|---|---|---|
| id | UUID PK DEFAULT gen_random_uuid() | |
| scope_type, scope_id | TEXT NOT NULL, UUID NOT NULL | |
| name | TEXT NOT NULL | |
| color | TEXT NOT NULL | |
| owner_id | UUID NULL FK users ON DELETE SET NULL | |
| status | TEXT NOT NULL DEFAULT 'todo' | todo/progress/done |
| start_sprint_id | UUID NULL FK sprints ON DELETE SET NULL | timeline bar start (explicit) |
| target_sprint_id | UUID NULL FK sprints ON DELETE SET NULL | timeline bar end (explicit) |
| description_md | TEXT NOT NULL DEFAULT '' | Markdown source |
| created_at, updated_at | TIMESTAMPTZ NOT NULL DEFAULT now() | |

`UNIQUE(scope_type, scope_id, LOWER(name))` as a functional unique index → `409`.

### 3.10 `tasks`
| column | type | notes |
|---|---|---|
| id | UUID PK DEFAULT gen_random_uuid() | |
| scope_type, scope_id | TEXT NOT NULL, UUID NOT NULL | immutable after create |
| key_num | INTEGER NOT NULL | per-scope counter → `NT-{key_num}` |
| title | TEXT NOT NULL DEFAULT '' | |
| description_md | TEXT NOT NULL DEFAULT '' | Markdown |
| status | TEXT NOT NULL DEFAULT 'backlog' | backlog/progress/review/done |
| priority | TEXT NOT NULL DEFAULT 'med' | high/med/low |
| assignee_id | UUID NULL FK users ON DELETE SET NULL | must be scope member |
| epic_id | UUID NULL FK epics ON DELETE SET NULL | epic delete clears chip |
| sprint_id | UUID NULL FK sprints ON DELETE SET NULL | defaults to current sprint |
| created_by | UUID FK users | |
| created_at, updated_at | TIMESTAMPTZ NOT NULL DEFAULT now() | |

`UNIQUE(scope_type, scope_id, key_num)`. Indexes: `(scope_type, scope_id, sprint_id, status)` for the board; `(epic_id)`; `(assignee_id)`.
Key allocation: `task_counters(scope_type, scope_id, next INTEGER)` row updated in the same transaction via `UPDATE … SET next = next + 1 RETURNING next` (avoids `MAX()+1` races; Postgres row-level locking makes this safe under concurrent inserts from multiple team members).

### 3.11 `task_activity`
| column | type | notes |
|---|---|---|
| id | UUID PK DEFAULT gen_random_uuid() | |
| task_id | UUID FK tasks ON DELETE CASCADE | |
| actor_id | UUID FK users | |
| kind | TEXT NOT NULL | created/status/priority/assignee/epic/sprint/comment |
| text | TEXT NOT NULL | "Status → In Review" or comment body |
| created_at | TIMESTAMPTZ NOT NULL DEFAULT now() | |

`INDEX(task_id, created_at)`. Written by the service layer on every logged `PATCH` and on `POST /comments`; never edited.

## 4. Access rules (single middleware)

| resource | read | write | admin-only write |
|---|---|---|---|
| notes, labels, prefs | owner | owner | — |
| team | any member | — | name, color, delete, members |
| tasks/epics/sprints (scope=personal) | owner | owner | — |
| tasks/epics (scope=team) | member | member | — |
| sprints (scope=team) | member | — | create/update/delete/mark current |

`scope_type='personal'` ⇒ `scope_id = current user`. Everything resolves to one query: `EXISTS(team_members WHERE team_id=? AND user_id=?)`.

## 5. Secrets: encryption design

```
login password  ──argon2id──▶ password_hash          (auth only)

master password ──argon2id(master_salt)──▶ KEK_user
                                             │ unwrap
users.wrapped_dek ────────────────────────▶ DEK  ──AES-256-GCM──▶ notes.password_enc
```

1. **First unlock** generates a random 32-byte DEK, wraps it with `KEK_user`, stores `wrapped_dek` + `master_hash`.
2. **`POST /vault/unlock`** verifies `master_hash`, unwraps DEK, re-wraps it with a **server KEK** (env/KMS) into `vault_tokens.dek_cache`, returns the opaque token. The master password is never stored and the DEK never sits in plaintext at rest.
3. **`GET /notes/{id}/secret`** loads `dek_cache` for the presented token, unwraps with the server KEK, decrypts `password_enc`, returns once. Missing/expired token → `403 VAULT_LOCKED`.
4. **`PUT /me/password`** (master) re-wraps the DEK under the new KEK_user — no re-encryption of notes — and deletes all `vault_tokens`.
5. Private notes reuse the same gate: with no valid token the API blanks `body/site/username` and sets `locked: true`.

Never log request bodies on `/vault/*`, `/notes/*/secret`, or `PUT /me/password`.

## 6. GORM notes (Postgres)

- Local dev runs Postgres in a container (Podman) — see §10.
- `id` columns: `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`; app can still pass a client-generated UUID on insert if needed.
- `time.Time` fields map directly to `timestamptz`; keep `archived_at`/`trashed_at`/`revoked_at` as `*time.Time`.
- The partial unique index (one current sprint per scope) and the `search_tsv` generated column + GIN index are raw SQL migrations run after `AutoMigrate` (GORM doesn't express either natively).
- `task_counters` update: `UPDATE … SET next = next + 1 RETURNING next` inside the create transaction.
- Functional unique indexes (`LOWER(name)`) are also raw SQL migrations, since GORM tags only cover plain column uniqueness.

## 7. Migration from today's schema

Current: `notes(noteID, title, body)`. Steps:
1. `AutoMigrate` adds new nullable/defaulted columns to `notes`; existing rows become plain, non-secret, non-private notes with `color='default'`.
2. Create `users` prefs columns + `refresh_tokens`, `vault_tokens`, `labels`.
3. Create `teams`, `team_members`; backfill nothing (users start with no teams).
4. Create `sprints`, `epics`, `tasks`, `task_counters`, `task_activity`.
5. Add `search_tsv` generated column + GIN index on `notes`; add partial unique index for current sprint.

## 8. Retention & housekeeping

- Purge `refresh_tokens`/`vault_tokens` where `expires_at < now() - interval '1 day'` (cron).
- Optional: hard-delete notes where `trashed_at < now() - interval '30 days'`.
- Activity is append-only; no purge.

## 9. Seed for local dev

Mirror the mockup: 1 user (Ada Keys), labels Work/Personal/Finance, 10 notes (4 secrets, 1 private), teams Orbit (admin) & Helix (not a member — create a second user), sprints 41–44 with 42 current, 5 epics, 9 tasks with activity. Keep in `cmd/seed`.

## 10. Local Postgres via Podman

Container setup for local development (no docker-compose required, though one can be added later):

```sh
podman pull docker.io/library/postgres:latest

podman run -d \
  --name noto-postgres \
  -e POSTGRES_DB=noto \
  -e POSTGRES_USER=noto \
  -e POSTGRES_PASSWORD=noto_dev_password \
  -p 5432:5432 \
  -v noto-postgres-data:/var/lib/postgresql \
  docker.io/library/postgres:latest
```

Note: the `postgres:latest` image is PG18+, which keeps data under `/var/lib/postgresql/<version>/...` (`pg_ctlcluster`-style layout) rather than directly in `/var/lib/postgresql/data`. Mount the volume at the parent `/var/lib/postgresql`, not `/var/lib/postgresql/data` — the latter fails to start with a "data directory... unused mount/volume" error. Pin to a specific major version (e.g. `postgres:18`) instead of `:latest` if you want this layout to stay stable across pulls.

- Data persists in the named volume `noto-postgres-data` across container restarts; `podman volume rm noto-postgres-data` to reset.
- Connection string for GORM: `host=localhost user=noto password=noto_dev_password dbname=noto port=5432 sslmode=disable`.
- `gen_random_uuid()` is built into Postgres ≥ 13 (the `postgres:latest` image); no `CREATE EXTENSION` needed for UUID generation.
- Stop/start: `podman stop noto-postgres` / `podman start noto-postgres`. Remove: `podman rm -f noto-postgres`.
