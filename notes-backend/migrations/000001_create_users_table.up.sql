CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    master_hash TEXT NULL,
    master_salt BYTEA NULL,
    wrapped_dek BYTEA NULL,
    default_role TEXT NULL,
    pref_default_tab TEXT NOT NULL DEFAULT 'notes',
    pref_notes_layout TEXT NOT NULL DEFAULT 'masonry',
    pref_mask_secrets BOOLEAN NOT NULL DEFAULT true,
    pref_autolock_min INTEGER NOT NULL DEFAULT 5,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- case-insensitive member/autocomplete search (§3.1)
CREATE INDEX idx_users_name_lower ON users (LOWER(name));
