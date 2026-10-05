# Nōto API — Contract Guide

> Companion to `noto-openapi.yaml`. Drop both in `Noto/notes-backend/openAPI-spec/`.
> Every endpoint maps to a screen or interaction in the mockup (`Noto - standalone.html`);
> the UX decisions behind them are in `roadmap/noto-ui-ux.md`.

## Conventions

- Base path `/api/v1`, JSON everywhere, UUID ids, `camelCase` fields.
- **Auth:** `Authorization: Bearer <accessToken>` (JWT RS256, as the backend already issues). Refresh tokens rotate via `/auth/refresh`; `/auth/logout` revokes.
- **Vault:** a second, short-lived `X-Vault-Token` from `POST /vault/unlock` gates anything sensitive. This is what makes the lock screen, "Lock vault", auto-lock and private notes real rather than cosmetic.
- **Errors:** `{ code, message, fields? }`. Codes the UI branches on: `VAULT_LOCKED` (show lock screen), `FORBIDDEN` (non-admin clicked a team card), `CONFLICT` (duplicate label/team/epic name).
- Lists return pinned/sorted results ready to render; the client does not re-sort.

## Screen → endpoint map

**Sign up / Log in / Sign out**
`POST /auth/signup` · `POST /auth/login` · `POST /auth/refresh` · `POST /auth/logout`. Both auth calls return tokens + the `User` (with `preferences`) so the app can land on `defaultTab` immediately.

**Lock screen / Lock vault**
`POST /vault/unlock` (master password → `vaultToken`, expiry = `autoLockMinutes`) · `POST /vault/lock`. The header lock button and the profile-menu "Lock vault" call `/vault/lock` and drop the token client-side.

**Profile & settings**
`GET /me` · `PATCH /me` (name, email, `preferences.defaultTab / notesLayout / maskSecrets / autoLockMinutes`) · `PUT /me/password` (revokes all vault tokens). Team memberships section reads `GET /teams` (`myRole`, `iAmAdmin`).

**Notes grid, sidebar, search**
`GET /notes?view=all|secrets|archive|trash&label=&q=`. Secret `password` is never in list payloads (`hasPassword: true` instead). Private notes without a vault token come back with `locked: true` and empty content — the card renders the "Private — locked" panel from that flag alone.

**Compose bar / editor**
`POST /notes` (NOTE or SECRET via `secret: true`) · `PATCH /notes/{id}` for every editor field including the NOTE⇄SECRET toggle, label, color, pin, private flag · `DELETE /notes/{id}` (delete forever).
Lifecycle is explicit and idempotent: `POST/DELETE /notes/{id}/archive`, `POST/DELETE /notes/{id}/trash`.

**Secret reveal / copy / generator**
`GET /notes/{id}/secret` with `X-Vault-Token` → decrypted password, fetched on demand when the user clicks Eye or Copy. `POST /notes/generate-password` mirrors the generator panel (length 8–40, charset flags, ambiguous chars excluded); the client may also generate locally.

**Labels** (sidebar "New label", editor "+ NEW")
`GET/POST /labels` · `PATCH/DELETE /labels/{id}`. Server assigns the dot color from the chart palette when omitted.

**Teams tab**
`GET /teams` (members, `openTaskCount`, `iAmAdmin` → ADMIN badge) · `POST /teams` (creator = Owner/admin) · `PATCH/DELETE /teams/{id}` admin-only.
Management modal: `POST /teams/{id}/members` · `PATCH/DELETE /teams/{id}/members/{userId}`; removing the last admin returns `409`.
Member autocomplete: `GET /directory?q=&excludeTeam=&limit=5` — matches name **or** email, returns `defaultRole` to pre-fill the role field.

**Tasks tab — Board**
`GET /tasks?scope=personal|{teamId}&sprintId=` (sprint dropdown; omit for "All sprints") · `POST /tasks` (ADD → modal; key `NT-n` generated, `sprintId` defaults to current) · `PATCH /tasks/{id}` (←/→ column moves = `status`) · `DELETE /tasks/{id}`.
Card chips need only what `Task` carries: `priority`, embedded `epic {name,color}`, `assignee.initials`.

**Task detail modal**
`GET /tasks/{id}` → `TaskDetail` with Markdown `description` and `activity[]`. `PATCH` any of title / description / status / priority / assignee / epic / sprint — the server appends `status|priority|assignee` changes to activity so the timeline matches the mockup without client bookkeeping. Comments: `POST /tasks/{id}/comments`.

**Tasks tab — Epics view & Epic modal**
`GET /epics?scope=` (counts + `spanSprintIds` for the timeline bar) · `POST /epics` · `GET /epics/{id}` (`EpicDetail.tasks` for the progress bar and task list) · `PATCH /epics/{id}` (name, color, owner, status, start/target sprint, Markdown description) · `DELETE /epics/{id}` (tasks' `epicId` cleared).

**Tasks tab — Timeline**
Composed client-side from `GET /sprints?scope=` (columns, with `startDate`/`endDate` for the hover tooltip and "(current)" marker) + `GET /epics?scope=` (`spanSprintIds` → bar extent; explicit start→target wins over the derived span).

**Sprints**
`GET/POST /sprints` · `PATCH/DELETE /sprints/{id}`. Exactly one sprint per scope is `current`; setting it clears the others.

## Data model deltas vs. current backend

Current `models.go` has `Note{noteID,title,body}`. Additions required:

- **Note:** `secret`, `private`, `site`, `username`, `passwordEnc` (AES-GCM, key derived from master password or a server KEK — never plaintext, never logged), `labelId`, `color`, `pinned`, `archived`, `trashed`.
- **User:** `passwordHash` (login) is separate from the **master-password verifier** (vault). Add `preferences` JSON.
- **New tables:** `labels`, `teams`, `team_members(role, admin)`, `tasks(key_seq per user)`, `task_activity`, `epics`, `sprints`, `vault_tokens`, `refresh_tokens`.
- **Scope rule:** every task/epic/sprint row has `scope` = `personal` (owner = user) or a `teamId`; access = owner or team member; mutation of team/sprint config = admin.

## Recommended build order

1. Auth + `/me` + vault tokens — unblocks login/signup/lock screens and settings.
2. Notes + labels + lifecycle + `/secret` reveal — the Notes tab end-to-end.
3. Teams + members + directory — Teams tab and admin modal.
4. Sprints → Epics → Tasks (+ activity, comments) — Tasks tab in the order the dropdowns depend on each other.
