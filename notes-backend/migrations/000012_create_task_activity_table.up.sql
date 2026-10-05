CREATE TABLE task_activity (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES tasks (id) ON DELETE CASCADE,
    actor_id UUID NOT NULL REFERENCES users (id),
    kind TEXT NOT NULL CHECK (kind IN ('created', 'status', 'priority', 'assignee', 'epic', 'sprint', 'comment')),
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_task_activity_task_id_created_at ON task_activity (task_id, created_at);

-- append-only: written by the service layer on every logged PATCH and
-- on POST /comments; never edited (§3.11, §8).
