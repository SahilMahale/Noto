-- per-scope NT-{key_num} allocation; updated via
-- `UPDATE task_counters SET next = next + 1 WHERE ... RETURNING next`
-- inside the same transaction as the task INSERT (avoids MAX()+1 races, §3.10)
CREATE TABLE task_counters (
    scope_type TEXT NOT NULL CHECK (scope_type IN ('personal', 'team')),
    scope_id UUID NOT NULL,
    next INTEGER NOT NULL DEFAULT 1,
    PRIMARY KEY (scope_type, scope_id)
);
