# Nōto — UI/UX Roadmap

> Supersedes `google-keep-ui.md`. Keep-style interaction model, but with the UX decisions
> locked in from the Nōto mockup (`mockup/Noto - standalone.html` — the visual source of truth)
> and the theme defined in `NOTO_THEME.md` + `notes-ui/styles/globals.css` (the applied Nōto theme, formerly `noto-globals.css`).

## Current Tech Stack

| Layer | Technology |
|---|---|
| Runtime / Bundler | Bun + Vite |
| UI Framework | React 19 + TypeScript |
| Styling | Tailwind CSS v4 (CSS-first, no tailwind.config) |
| Component Library | shadcn/ui (new-york, cssVariables) |
| Routing | TanStack Router (file-based) |
| Data Fetching | TanStack Query (React Query) |
| Forms | React Hook Form + Zod |
| HTTP | Axios |
| Icons | Lucide React |
| Backend | Go + Fiber v2 |
| ORM / DB | GORM + SQLite |
| Auth | JWT (RS256) via gofiber/contrib/jwt |

## Current State

- Auth (login, signup) works end-to-end
- Backend has full notes CRUD: `POST /notes/create`, `PATCH /notes/update/:noteID`, `DELETE /notes/delete/:noteID`, `GET /notes/get`, `GET /notes/get/:noteId`
- Home page (`/_auth/Home`) is a placeholder
- `src/api/api.tsx` still has leftover bookings code — no notes API calls exist yet
- Note data model: `{ noteID, title, body }`

---

## UX Decisions (locked)

These override anything Keep does differently. Rationale and pixel reference: the mockup.

1. **Theme** — carbonfox/oxocarbon dark-only palette (`NOTO_THEME.md`). Flat & ruled: `--radius: 0` (**no rounded corners anywhere** — replaces the old `rounded-xl` card spec), no card shadows; structure comes from 1px `--border` rules. IBM Plex Sans for UI/body, IBM Plex Mono (`font-mono`) for secret values, password inputs, and counts.
2. **Notes and secrets are one entity.** Every note has a `secret` toggle (NOTE ⇄ SECRET segmented control in the editor). Secret mode swaps the free-text body for **structured fields: Site / Username / Password** — not free text.
3. **Secret field behavior** — passwords **masked by default** (`••••••••••`), per-card `Eye` reveal toggle, one-click `Copy` buttons on username and password (→ `navigator.clipboard` + toast). Cards show a teal `SECRET` badge.
4. **Password generator** — inline panel in the editor: length slider 8–40, charset chips (A–Z / 0–9 / #!%), Regenerate. `crypto.getRandomValues`, ambiguous chars excluded (`il1o0IO`). New secrets get a generated password by default.
5. **Private notes** — a second, independent per-note flag (`private`, purple `Shield`). Locked private notes hide their entire content behind a "Private — locked" panel with an UNLOCK button; card can be re-locked from its footer. **Locking the vault re-locks all private notes.** Works on both plain notes and secrets.
6. **Vault lock screen** — app opens locked (master password → unlock); header lock button re-locks. Restyle the existing login route as the vault screen (`ノ` accent square + "Nōto" wordmark).
7. **View switcher** — MASONRY / GRID / LIST segmented control in the header. Masonry = CSS columns (`column-width: 250px`, `break-inside: avoid`); Grid = `repeat(auto-fill, minmax(250px, 1fr))`; List = single `max-w-[680px]` column. One shared `NoteCard`.
8. **Collapsible sidebar** — `Menu` hamburger in header toggles between full (216px: All notes, Secrets, Labels, Archive, Trash) and a 54px icon rail (labels collapse to their colored dots). Active item = 2px left border in `--primary` + card fill.
9. **Labels** — dynamic, user-creatable from **both** the sidebar ("New label" inline input) and the editor ("+ NEW" chip, which also assigns it to the open note). Auto-assigned dot color from the `--chart-*` cycle. Label chips toggle (click again to clear). Shown as an outline chip on the card footer.
10. **Compose bar** — always-expanded single input at top of grid with two explicit quick-add buttons: `+ NOTE` and `+ SECRET` (teal). Either opens the editor with the typed text as body/title. (Replaces Keep's collapse-on-blur pattern — discoverability of SECRET mode wins.)
11. **Card actions are always visible** in a bordered footer row (archive, trash, private lock; pin stays top-right) — **not** hover-only. Hover-reveal hides the secret/private affordances and fails on touch.
12. **Archive / Trash lifecycle** — archive and trash are flags, not deletion. Archive view offers green `UNARCHIVE`; Trash offers green `RESTORE` and `DELETE FOREVER` (magenta). Editor footer swaps Archive ⇄ Unarchive contextually. No confirm dialogs; the toast is the feedback.
13. **Card color tints** — 6 curated dark tints (`--note-*` tokens: default/blue/teal/violet/rose/green) as square swatches in the editor. Replaces the old 11-color Keep palette.
14. **Toast** — bottom-left, light-on-dark inversion (`bg-foreground text-background`), green check, ~1.8s. Used for copy/archive/trash/restore/label-created/lock feedback.
15. **Search** — header input, client-side filter over `title + body + site + username + label`.
16. **Pinning** — pin icon top-right of card; pinned sort first within the current view (mockup approach) rather than a separate "Pinned" section.

---

## Data Model (target)

```ts
export interface Note {
  noteID: string
  secret: boolean        // password/secret mode
  private: boolean       // lockable
  title: string
  body: string           // plain notes
  site: string           // secrets
  username: string       // secrets
  password: string       // secrets — see backend note below
  label: string | null
  color: NoteColor       // 'default'|'blue'|'teal'|'violet'|'rose'|'green'
  pinned: boolean
  archived: boolean
  trashed: boolean
}
```

**Backend:** add all fields to `models/models.go` (GORM `AutoMigrate`). ⚠️ `password` must be **encrypted at rest** (AES-GCM with a server-side key, or derived from the user's master secret) — never stored plaintext, never logged, and excluded from list-endpoint payloads if possible (fetch on reveal).

---

## Phase 1 — Theme + Wire Up the Notes API

1. ~~Apply `noto-globals.css` over `styles/globals.css`; fix `src/index.css` (`bg-zinc-950` → `bg-background`).~~ **Done 2026-07-12** (theme + IBM Plex fonts live in `styles/globals.css`). Remaining: strip `zinc/sky/amber` literals from NavBar/Home.
2. Rewrite `src/api/api.tsx` — remove bookings/users stubs; typed `getNotes`, `getNoteById`, `createNote`, `updateNote`, `deleteNote`; token from `authContext`.
3. `src/types/notes.ts` — model above. `src/hooks/useNotes.ts` — `useNotes()`, `useCreateNote()`, `useUpdateNote()`, `useDeleteNote()`.

## Phase 2 — Shell + Grid (hardcoded data first)

- Header (rework `NavBar.tsx`): menu toggle · `ノ` logo + wordmark · search · view switcher · lock button. Remove Book Tables/Users.
- `Sidebar.tsx` (collapsible, decision 8) · `NoteGrid.tsx` (3 layouts, decision 7) · `NoteCard.tsx` (badges, masked fields, footer actions — decisions 3, 5, 11) · `CreateNoteBar.tsx` (decision 10).

## Phase 3 — Editor + CRUD

- `NoteEditModal.tsx` (Radix Dialog): NOTE⇄SECRET toggle, structured secret fields with copy/reveal, `PasswordGenerator.tsx` (decision 4), label chips + create, color swatches, private/pin/archive/trash footer, DONE. Auto-save debounce 500ms.
- Optimistic create/edit/delete/pin via the Phase 1 hooks.

## Phase 4 — Lifecycle + Polish

- Archive & Trash routes/views with unarchive/restore/delete-forever (decision 12).
- Vault lock screen + private-note lock wiring (decisions 5–6).
- Labels end-to-end (decision 9) — needs `label` on backend or a client-side labels store first.
- Search filter (decision 15) · Toast system (decision 14) · empty states ("Nothing matches your search." / "Trash is empty.") · card-shaped `animate-pulse` skeletons.

## Implementation Order

```
Phase 1 → Phase 2 (hardcoded) → Phase 2 (real data)
        → Phase 3 (create → edit → secret fields → generator → delete → pin)
        → Phase 4 (archive/trash → vault+private → labels → search → toasts → polish)
```

## Files to Create / Modify

```
notes-ui/src/
  types/notes.ts                ← NEW
  hooks/useNotes.ts             ← NEW
  api/api.tsx                   ← REWRITE
  components/
    Sidebar.tsx                 ← NEW
    NoteCard.tsx                ← NEW
    NoteGrid.tsx                ← NEW
    CreateNoteBar.tsx           ← NEW
    NoteEditModal.tsx           ← NEW
    PasswordGenerator.tsx       ← NEW
    ColorPicker.tsx             ← NEW (6 square swatches)
    LabelPicker.tsx             ← NEW (chips + inline create)
    Toast.tsx                   ← NEW
  components/NavBar.tsx         ← REWORK (header, decision 8/15)
  routes/_auth.Home.tsx         ← REPLACE placeholder
  routes/login/                 ← RESTYLE as vault lock screen

notes-backend/
  models/models.go              ← ADD secret, private, site, username, password(encrypted),
                                   label, color, pinned, archived, trashed
```
