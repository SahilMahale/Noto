CREATE TABLE labels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    color TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- UNIQUE(user_id, name COLLATE NOCASE) in SQLite terms (§3.4)
CREATE UNIQUE INDEX uq_labels_user_id_name_lower ON labels (user_id, LOWER(name));
