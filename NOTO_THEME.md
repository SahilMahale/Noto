# Nōto — Theme & Build Guide

> Drop this file in your repo root (`notes-ui/NOTO_THEME.md`). Point Claude Code at it:
> *"Follow NOTO_THEME.md to apply the Nōto theme and build the notes screens."*
> If you keep a `CLAUDE.md`, add: `See NOTO_THEME.md for the design system.`

This guide turns the current scaffold into **Nōto** — a Google-Keep-style notes app with per-note **secret** and **private** modes, themed with the **carbonfox / oxocarbon** (IBM Carbon) dark palette. It maps that palette onto your existing shadcn CSS variables so every shadcn component inherits it automatically.

---

## 1. Stack (as-is — don't fight it)

- **Tailwind v4**, CSS-first config. Theme lives in `styles/globals.css` via `@theme` + `:root`/`.dark` CSS variables. **There is no `tailwind.config.js`** — do not create one.
- **shadcn/ui**, style `new-york`, base color `neutral`, `cssVariables: true`, icons = **lucide-react**. Add components with the shadcn CLI; they read the vars below.
- **TanStack Router** (file-based, `src/routes`, generated `routeTree.gen.ts` — don't hand-edit it).
- **Forms:** react-hook-form + zod + `@hookform/resolvers`. **Overlays:** `@headlessui/react` + `radix-ui` are available (use Radix/Headless Dialog for the editor modal).
- **Fonts:** **IBM Plex Sans** (`--font-sans`, variable wght 100–700) + **IBM Plex Mono** (`--font-mono`, 400/500/600) — self-hosted latin subsets in `assets/fonts`, the same files the mockup embeds. See §4.

### Scaffold cruft to remove
These are leftover from a starter template and are **not** part of Nōto:
- `src/routes/_auth.Home.tsx` — "Bun + React" demo card. Replace with the notes home.
- `NavBar.tsx` — the `Book Tables` / `Users` links and emoji. Rework into the app header (§6).
- Stock shadcn `:root`/`.dark` values in `globals.css` (pure black/white neutral) — overwrite with §3.

---

## 2. Palette — carbon / oxocarbon tokens

The raw ramp. Values are the canonical IBM Carbon / oxocarbon hexes (Tailwind v4 accepts hex in theme vars — no need to convert to oklch, though you may).

**Neutral ramp (backgrounds → text)**
| Token | Hex | Role |
|---|---|---|
| `gray-100` | `#161616` | app background (canvas) |
| `gray-90`  | `#262626` | card / surface, inputs |
| `#1f1f1f`  | `#1f1f1f` | raised surface (compose bar), hover row |
| `gray-80`  | `#393939` | borders, rules, dividers |
| `gray-70`  | `#525252` | strong border, faint text, placeholder |
| `gray-60`  | `#6f6f6f` | muted labels |
| `gray-50`  | `#8d8d8d` | secondary text / icons |
| `gray-40`  | `#a8a8a8` | tertiary text |
| `gray-30`  | `#c6c6c6` | note body text |
| `gray-20`  | `#dde1e6` | high-contrast body |
| `gray-10`  | `#f2f4f8` | primary text, "DONE" button fill |

**Accents (oxocarbon)**
| Token | Hex | Meaning in Nōto |
|---|---|---|
| `blue`   | `#33b1ff` | **primary** — logo, focus ring, active nav, links |
| `blue-40`| `#78a9ff` | link hover, "Work" label dot |
| `teal`   | `#3ddbd9` | **secret** mode — badge, generate, secret accents |
| `cyan`   | `#08bdba` | secret alt / label dot |
| `green`  | `#42be65` | success — restore, unarchive, toasts check |
| `purple` | `#be95ff` | **private** mode — shield badge, lock accents |
| `magenta`| `#ee5396` | **destructive** — delete/trash |
| `pink`   | `#ff7eb6` | label dot |

**Note card background tints** (subtle, desaturated — Keep-style, but dark):
`#262626` (default), `#21283b` (blue), `#16312f` (teal), `#2a2340` (violet), `#38202e` (rose), `#1e3226` (green).

---

## 3. Drop-in `styles/globals.css` theme block

**This theme is already applied** — `styles/globals.css` *is* the Nōto theme (the former `noto-globals.css`, merged 2026-07-12): the `@theme inline` bridge, the IBM Plex `@font-face` blocks, and everything below. `src/index.css` has also been fixed (`bg-zinc-950` → `bg-background`). The block below documents the values for reference.

The app runs **dark-only** — `:root` and `.dark` are set to the same values so it themes correctly even before a `.dark` class is applied, and `--radius: 0` gives the flat Carbon look (the mockup uses sharp corners, not rounded).

```css
:root, .dark {
  --radius: 0rem;                 /* flat / sharp — Carbon aesthetic. Use 0.25rem if you want a hair of softening */

  --background:            #161616;
  --foreground:            #f2f4f8;

  --card:                  #262626;
  --card-foreground:       #f2f4f8;
  --popover:               #1f1f1f;
  --popover-foreground:    #f2f4f8;

  --primary:               #33b1ff;   /* oxocarbon blue */
  --primary-foreground:    #161616;   /* dark text on blue fills */

  --secondary:             #262626;
  --secondary-foreground:  #f2f4f8;

  --muted:                 #262626;
  --muted-foreground:      #8d8d8d;

  --accent:                #1f1f1f;   /* hover surface */
  --accent-foreground:     #f2f4f8;

  --destructive:           #ee5396;   /* magenta */

  --border:                #393939;
  --input:                 #393939;
  --ring:                  #33b1ff;   /* blue focus ring */

  /* sidebar (used by shadcn sidebar primitives if you adopt them) */
  --sidebar:               #161616;
  --sidebar-foreground:    #a8a8a8;
  --sidebar-primary:       #33b1ff;
  --sidebar-primary-foreground: #161616;
  --sidebar-accent:        #262626;
  --sidebar-accent-foreground:  #f2f4f8;
  --sidebar-border:        #262626;
  --sidebar-ring:          #33b1ff;

  /* charts / label dots, reused for label swatches */
  --chart-1: #78a9ff;  /* Work    */
  --chart-2: #42be65;  /* Personal*/
  --chart-3: #be95ff;  /* Finance */
  --chart-4: #3ddbd9;
  --chart-5: #ff7eb6;
}
```

Feature semantic tokens follow the shadcn pattern — raw values in `:root` (`--secret`, `--private`, `--success`, `--surface`, `--note-*`), mapped in `@theme inline` (`--color-secret: var(--secret);` …) so `bg-secret`, `text-private`, `border-success`, `bg-note-blue` etc. work in markup. See `styles/globals.css` for the full set, including the six note-card tint tokens.

> Because shadcn components consume `--primary`, `--border`, `--ring`, etc., swapping these values re-themes **every** shadcn Button/Card/Input/Select without touching component code. Replace ad-hoc `zinc-*` / `sky-*` / `amber-*` utility classes in existing files (NavBar, Home) with the semantic ones (`bg-background`, `text-muted-foreground`, `border-border`, `text-primary`, `bg-destructive`).

---

## 4. Typography

**IBM Plex** is the Nōto typeface pair (matching the mockup): **IBM Plex Sans** for UI and body copy (`--font-sans` / default), **IBM Plex Mono** for technical values (`--font-mono` → `font-mono`). Both self-hosted in `assets/fonts` as latin subsets — Sans is a single variable font (wght 100–700; use 400/500/600/700), Mono ships as static 400/500/600. Note: latin subset does not cover the katakana `ノ` in the logo — it falls back to a system font, same as the mockup.

Conventions:
- **UI chrome, headings, buttons, note bodies:** IBM Plex Sans. Uppercase micro-labels use `text-[10px] tracking-[0.14em] text-muted-foreground` (e.g. `SITE`, `USERNAME`, `LABELS`).
- **Monospace (`font-mono`):** secret field values (site / username / password), password inputs, generator charset chips, and item counts — anywhere the mockup shows credential-like or numeric data.
- Note body text: `text-sm leading-relaxed text-[#c6c6c6]`.
- Do **not** round type scale below `12px`.

---

## 5. Look & feel rules (the "system")

- **Flat & ruled.** No rounded corners (`--radius: 0`), no drop shadows on cards. Structure comes from **1px `--border` rules** and flush-left alignment. The only shadow is on the editor modal overlay.
- **Active nav** = 2px left border in `--primary` + `--card` fill (see mockup). Inactive = transparent with `hover:bg-[#1f1f1f]`.
- **Focus:** always `outline: 2px solid var(--ring); outline-offset: 2px` — never remove it. shadcn's `focus-visible:ring-ring` already does this once `--ring` is set.
- **Icons:** lucide-react, `size={15}` in chrome, `strokeWidth={2}`. Exact icons used in the mockup:
  - search `Search`, lock/secret `Lock` / `LockOpen`, private `Shield`, pin `Pin`, archive `Archive`, unarchive `ArchiveRestore`, trash `Trash2`, restore `RotateCcw`, copy `Copy`, reveal `Eye` / `EyeOff`, generate `RefreshCw`, add `Plus`, menu `Menu`, close `X`, check (toast) `Check`.
- **Buttons:** primary = `bg-primary text-primary-foreground` (dark text on blue); confirm/DONE = `bg-foreground text-background`; secret action = `border-secret text-secret hover:bg-secret/10`; destructive = ghost that goes `text-destructive hover:bg-destructive/10`.
- **Toast:** bottom-left, `bg-foreground text-background`, green `Check` icon, ~1.8s.

---

## 6. Build tasks for Claude Code

Order matters — theme first, then screens. Each maps to files that already exist.

1. **Apply theme** — overwrite `:root`/`.dark` in `styles/globals.css` (§3), add semantic tokens, set `--radius: 0`. Verify existing pages still render, then strip `zinc/sky/amber` literals in favor of semantic classes.
2. **App shell / header** — rework `src/components/NavBar.tsx`: left = `Menu` toggle + logo (accent square with `ノ` + "Nōto" wordmark), center = search input (`bg-card border-border`), right = view switcher (MASONRY / GRID / LIST segmented) + lock button. Remove `Book Tables`/`Users`.
3. **Collapsible sidebar** — All notes, Secrets, dynamic Labels (with colored dots + a "New label" inline input), Archive, Trash. Collapses to a 54px icon rail (width transition); labels hide, dots remain. Active item = left blue border (§5).
4. **Notes grid** — three layouts sharing one card component:
   - *masonry* → CSS `column-width: 250px; column-gap: 14px` (cards `break-inside: avoid`).
   - *grid* → `grid-template-columns: repeat(auto-fill, minmax(250px, 1fr))`.
   - *list* → single `max-w-[680px]` column.
5. **Compose bar** — top of grid: text input + `+ NOTE` and `+ SECRET` (teal) quick-add buttons that open the editor.
6. **Note card** — title, body (or secret fields), badges (SECRET teal / PRIVATE purple), pin, footer actions (archive/unarchive, trash/restore, private lock). Locked private cards hide content behind a "Private — locked" panel with an UNLOCK button.
7. **Editor modal** (Radix/Headless `Dialog`) — NOTE ⇄ SECRET segmented toggle at top; NOTE = title + textarea; SECRET = structured **Site / Username / Password** fields, each with a copy button, password field with reveal + **Generate**. Footer: LABEL chips (+ New), color swatches, private-shield toggle, pin, archive, trash, DONE.
8. **Password generator** — inline panel in the editor: length slider (8–40), charset chips (A–Z / 0–9 / #!%), Regenerate. Use `crypto.getRandomValues` (excludes ambiguous chars `il1o0`). See snippet §7.
9. **Secret behavior** — passwords masked by default (`••••••••••`), click `Eye` to reveal per-card, `Copy` writes to `navigator.clipboard` and fires the toast. Structured fields, not free text.
10. **Private notes** — a `private` boolean per note; when the vault is locked (or on load) private notes render locked and their content is withheld until UNLOCK. Locking the vault re-locks all private notes.
11. **Lock screen** — `src/routes/login` already exists; restyle it as the vault unlock (accent `ノ` square, "Nōto — ノート · notes & secrets — locked", master-password input, blue UNLOCK). Wire to `authContext`.
12. **Labels** — dynamic list in context/state; creating one auto-assigns a dot color from `--chart-*` and it appears in both the sidebar filter and the editor chips.

### Suggested data shape
```ts
type Note = {
  id: string;
  secret: boolean;      // password/secret mode
  private: boolean;     // lockable
  title: string;
  body: string;         // notes
  site: string; username: string; password: string;  // secrets
  label: string | null;
  color: string;        // one of the card tints (§2)
  pinned: boolean; archived: boolean; trashed: boolean;
};
```

---

## 7. Password generator (reference)

```ts
export function generatePassword(len = 16, opts = { upper: true, nums: true, syms: true }) {
  let chars = "abcdefghijkmnopqrstuvwxyz";              /* no l */    // pragma: allowlist secret
  if (opts.upper) chars += "ABCDEFGHJKLMNPQRSTUVWXYZ";  /* no I, O */ // pragma: allowlist secret
  if (opts.nums)  chars += "23456789";                  // no 0, 1
  if (opts.syms)  chars += "!#$%&*+-=?@_";
  const buf = new Uint32Array(len);
  crypto.getRandomValues(buf);
  return Array.from(buf, (n) => chars[n % chars.length]).join("");
}
```

---

## 8. Reference mockup

An interactive HTML mockup of every screen/state lives alongside this repo (`Carbon Notes.dc.html` / the standalone bundle). Treat it as the visual source of truth for spacing, states, and interactions — match it, don't reinvent. When in doubt about a color's role, cross-check §2 against the mockup.
