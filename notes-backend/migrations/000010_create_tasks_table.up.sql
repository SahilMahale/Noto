CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scope_type TEXT NOT NULL CHECK (scope_type IN ('personal', 'team')),
    scope_id UUID NOT NULL,
    key_num INTEGER NOT NULL,
    title TEXT NOT NULL DEFAULT '',
    description_md TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'backlog' CHECK (status IN ('backlog', 'progress', 'review', 'done')),
    priority TEXT NOT NULL DEFAULT 'med' CHECK (priority IN ('high', 'med', 'low')),
    assignee_id UUID NULL REFERENCES users (id) ON DELETE SET NULL,
    epic_id UUID NULL REFERENCES epics (id) ON DELETE SET NULL,
    sprint_id UUID NULL REFERENCES sprints (id) ON DELETE SET NULL,
    created_by UUID NOT NULL REFERENCES users (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (scope_type, scope_id, key_num)
);

-- board view: filter by scope + sprint + status
CREATE INDEX idx_tasks_board ON tasks (scope_type, scope_id, sprint_id, status);
CREATE INDEX idx_tasks_epic_id ON tasks (epic_id);
CREATE INDEX idx_tasks_assignee_id ON tasks (assignee_id);
