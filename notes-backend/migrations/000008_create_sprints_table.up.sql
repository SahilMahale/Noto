CREATE TABLE sprints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scope_type TEXT NOT NULL CHECK (scope_type IN ('personal', 'team')),
    scope_id UUID NOT NULL,
    name TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_current BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (scope_type, scope_id, name)
);

-- exactly one current sprint per scope - native partial unique index (§3.8)
CREATE UNIQUE INDEX uq_sprints_one_current_per_scope ON sprints (scope_type, scope_id) WHERE is_current;
