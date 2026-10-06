CREATE TABLE notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    label_id UUID NULL REFERENCES labels (id) ON DELETE SET NULL,
    is_secret BOOLEAN NOT NULL DEFAULT false,
    is_private BOOLEAN NOT NULL DEFAULT false,
    title TEXT NOT NULL DEFAULT '',
    body TEXT NOT NULL DEFAULT '',
    site TEXT NOT NULL DEFAULT '',
    username TEXT NOT NULL DEFAULT '',
    password_enc BYTEA NULL,
    color TEXT NOT NULL DEFAULT 'default',
    pinned BOOLEAN NOT NULL DEFAULT false,
    archived_at TIMESTAMPTZ NULL,
    trashed_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    -- password_enc is intentionally excluded: never indexed or searched (§3.5)
    search_tsv tsvector GENERATED ALWAYS AS (
        to_tsvector(
            'english',
            coalesce(title, '') || ' ' ||
            coalesce(body, '') || ' ' ||
            coalesce(site, '') || ' ' ||
            coalesce(username, '')
        )
    ) STORED
);

-- note grid listing (user's notes, sorted pinned-first then most recent)
CREATE INDEX idx_notes_grid ON notes (user_id, trashed_at, archived_at, pinned DESC, updated_at DESC);
CREATE INDEX idx_notes_user_label ON notes (user_id, label_id);
CREATE INDEX idx_notes_search_tsv ON notes USING GIN (search_tsv);
