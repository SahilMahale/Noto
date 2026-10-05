CREATE TABLE epics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scope_type TEXT NOT NULL CHECK (scope_type IN ('personal', 'team')),
    scope_id UUID NOT NULL,
    name TEXT NOT NULL,
    color TEXT NOT NULL,
    owner_id UUID NULL REFERENCES users (id) ON DELETE SET NULL,
    status TEXT NOT NULL DEFAULT 'todo' CHECK (status IN ('todo', 'progress', 'done')),
    start_sprint_id UUID NULL REFERENCES sprints (id) ON DELETE SET NULL,
    target_sprint_id UUID NULL REFERENCES sprints (id) ON DELETE SET NULL,
    description_md TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX uq_epics_scope_name_lower ON epics (scope_type, scope_id, LOWER(name));
