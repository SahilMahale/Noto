# Nōto — Google Keep Style UI Roadmap

## Current Tech Stack

| Layer | Technology |
|---|---|
| Runtime / Bundler | Bun + Vite |
| UI Framework | React 19 + TypeScript |
| Styling | Tailwind CSS v4 |
| Component Library | shadcn/ui (new-york, neutral) |
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

## Phase 1 — Wire Up the Notes API

**Goal:** Replace stale bookings code and load real notes on the Home page.

### Tasks

1. **Rewrite `src/api/api.tsx`**
   - Remove all bookings/users stubs
   - Add typed functions: `getNotes`, `getNoteById`, `createNote`, `updateNote`, `deleteNote`
   - Auth token sourced from `authContext`

2. **Define shared types** in `src/types/notes.ts`
   ```ts
   export interface Note {
     noteID: string
     title: string
     body: string
   }
   export type CreateNotePayload = { title: string; body: string }
   export type UpdateNotePayload = { title?: string; body?: string }
   ```

3. **TanStack Query hooks** in `src/hooks/useNotes.ts`
   - `useNotes()` — GET all notes
   - `useCreateNote()` — mutation
   - `useUpdateNote()` — mutation
   - `useDeleteNote()` — mutation

---

## Phase 2 — Google Keep Cards Layout

**Goal:** Render notes as a masonry card grid, matching Keep's visual style.

### Layout

- **Top bar**: existing NavBar (keep as-is)
- **Create note input strip** at the top of the page (collapsed by default, expands on focus — like Keep's "Take a note…")
- **Card grid**: CSS columns masonry (`columns-1 sm:columns-2 lg:columns-3 xl:columns-4`) with `break-inside-avoid` per card

### Components to build

| Component | Location | Responsibility |
|---|---|---|
| `NoteCard` | `src/components/NoteCard.tsx` | Renders a single note card with hover controls |
| `NoteGrid` | `src/components/NoteGrid.tsx` | Masonry wrapper, maps notes → NoteCards |
| `CreateNoteInput` | `src/components/CreateNoteInput.tsx` | Collapsed/expanded inline create input |
| `NoteEditModal` | `src/components/NoteEditModal.tsx` | Full-screen modal for editing a note |
| `ColorPicker` | `src/components/ColorPicker.tsx` | Dot grid of background colour choices |

### NoteCard anatomy

```
┌─────────────────────────────┐  ← rounded-xl, coloured bg
│  [pin icon]          [⋮]   │  ← show on hover only
│                             │
│  Title (bold)               │
│  Body text (truncated,      │
│  max ~8 lines)              │
│                             │
│  ─────────────────────────  │
│  [🎨] [🗑] [📦 archive]    │  ← action bar, show on hover
└─────────────────────────────┘
```

### Card colour palette (stored client-side first, then backend)

```
default | red | orange | yellow | green | teal | blue | purple | pink | brown | grey
```

---

## Phase 3 — Full CRUD Interactions

**Goal:** Create, edit, delete, pin notes without page reloads.

### Create note

- `CreateNoteInput` starts as a single text row ("Take a note…")
- On focus → expands to show title + body fields + action bar
- On blur (click outside) or submit → calls `useCreateNote`, optimistically adds card to grid, collapses input

### Edit note

- Click anywhere on a `NoteCard` body → opens `NoteEditModal`
- Modal has title + full-height textarea + colour picker + close button
- Auto-saves on close (debounce 500 ms then `useUpdateNote`)

### Delete note

- Trash icon on NoteCard hover → `useDeleteNote` with optimistic removal
- No confirmation dialog (matches Keep behaviour; can undo via snackbar — see Phase 4)

### Pin note

- Pin icon top-right of NoteCard
- Pinned notes render in a separate "Pinned" section above the main grid
- Pin state stored in backend (`pinned bool` field — requires backend model change, see below)

### Backend model change needed for pin

Add `pinned` field to the notes table and expose it in `NoteRequest` / `NoteResp`. Migration is automatic via GORM `AutoMigrate`.

---

## Phase 4 — Polish & Extra Features

### Search / filter

- Search input in NavBar
- Client-side filter on `note.title + note.body` using a controlled query string
- No new API call needed (all notes already in memory via React Query cache)

### Archive

- Archive icon on NoteCard hover
- Archived notes hidden from main grid
- `/archive` route shows archived-only notes
- Backend: add `archived bool` field (same migration approach as pin)

### Undo snackbar

- After delete/archive → show a bottom snackbar "Note deleted · Undo" for 5 s
- Undo calls a restore endpoint or re-creates the note

### Empty states

- If no notes: centred illustration + "Your notes will appear here" (match Keep's empty state)
- If search returns 0 results: "No matching notes"

### Loading & skeleton

- Replace spinner with card-shaped skeletons (`animate-pulse`) while notes load

---

## Implementation Order

```
Phase 1  →  Phase 2 (layout only, hardcoded data)  →  Phase 2 (real data via Phase 1)
         →  Phase 3 (CRUD one feature at a time: create → edit → delete → pin)
         →  Phase 4 (search → archive → undo → empty states → skeletons)
```

---

## Files to Create / Modify

```
notes-ui/src/
  types/
    notes.ts                    ← NEW
  hooks/
    useNotes.ts                 ← NEW
  api/
    api.tsx                     ← REWRITE (remove bookings, add notes)
  components/
    NoteCard.tsx                ← NEW
    NoteGrid.tsx                ← NEW
    CreateNoteInput.tsx         ← NEW
    NoteEditModal.tsx           ← NEW
    ColorPicker.tsx             ← NEW
  routes/
    _auth.Home.tsx              ← REPLACE placeholder with NoteGrid + CreateNoteInput

notes-backend/
  models/models.go              ← ADD pinned, archived, color fields
  server/server.go              ← no route changes needed for Phase 1-3
```
