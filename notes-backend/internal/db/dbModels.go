package db

import "time"

// Structs mirror the Postgres schema in notes-backend/migrations/ (see
// roadmap/noto-db-schema.md §3 for the design). Table names are pinned
// explicitly via TableName() to match the migrations exactly, since GORM's
// auto-pluralization would otherwise rename e.g. task_activity -> task_activities.

// ── users ──────────────────────────────────────────────────────────────

type User struct {
	ID              string    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Name            string    `gorm:"column:name;not null"`
	Email           string    `gorm:"column:email;not null;uniqueIndex"`
	PasswordHash    string    `gorm:"column:password_hash;not null"`
	MasterHash      *string   `gorm:"column:master_hash"`
	MasterSalt      []byte    `gorm:"column:master_salt"`
	WrappedDek      []byte    `gorm:"column:wrapped_dek"`
	DefaultRole     *string   `gorm:"column:default_role"`
	PrefDefaultTab  string    `gorm:"column:pref_default_tab;not null;default:notes"`
	PrefNotesLayout string    `gorm:"column:pref_notes_layout;not null;default:masonry"`
	PrefMaskSecrets bool      `gorm:"column:pref_mask_secrets;not null;default:true"`
	PrefAutolockMin int       `gorm:"column:pref_autolock_min;not null;default:5"`
	CreatedAt       time.Time `gorm:"column:created_at;not null"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null"`
}

func (User) TableName() string { return "users" }

// ── refresh_tokens ────────────────────────────────────────────────────

type RefreshToken struct {
	ID        string     `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    string     `gorm:"column:user_id;type:uuid;not null"`
	TokenHash string     `gorm:"column:token_hash;not null;uniqueIndex"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null"`
	RevokedAt *time.Time `gorm:"column:revoked_at"`
	CreatedAt time.Time  `gorm:"column:created_at;not null"`
}

func (RefreshToken) TableName() string { return "refresh_tokens" }

// ── vault_tokens ──────────────────────────────────────────────────────

type VaultToken struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    string    `gorm:"column:user_id;type:uuid;not null"`
	TokenHash string    `gorm:"column:token_hash;not null;uniqueIndex"`
	DekCache  []byte    `gorm:"column:dek_cache;not null"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
}

func (VaultToken) TableName() string { return "vault_tokens" }

// ── labels ────────────────────────────────────────────────────────────

type Label struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    string    `gorm:"column:user_id;type:uuid;not null"`
	Name      string    `gorm:"column:name;not null"`
	Color     string    `gorm:"column:color;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
}

func (Label) TableName() string { return "labels" }

// ── notes ─────────────────────────────────────────────────────────────

type Note struct {
	ID          string     `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      string     `gorm:"column:user_id;type:uuid;not null"`
	LabelID     *string    `gorm:"column:label_id;type:uuid"`
	IsSecret    bool       `gorm:"column:is_secret;not null;default:false"`
	IsPrivate   bool       `gorm:"column:is_private;not null;default:false"`
	Title       string     `gorm:"column:title;not null;default:''"`
	Body        string     `gorm:"column:body;not null;default:''"`
	Site        string     `gorm:"column:site;not null;default:''"`
	Username    string     `gorm:"column:username;not null;default:''"`
	PasswordEnc []byte     `gorm:"column:password_enc"`
	Color       string     `gorm:"column:color;not null;default:default"`
	Pinned      bool       `gorm:"column:pinned;not null;default:false"`
	ArchivedAt  *time.Time `gorm:"column:archived_at"`
	TrashedAt   *time.Time `gorm:"column:trashed_at"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;not null"`
	// search_tsv is a DB-generated (STORED) column — never written by the app.
	SearchTSV string `gorm:"column:search_tsv;->"`
}

func (Note) TableName() string { return "notes" }

// ── teams ─────────────────────────────────────────────────────────────

type Team struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `gorm:"column:name;not null"`
	Color     string    `gorm:"column:color;not null"`
	CreatedBy string    `gorm:"column:created_by;type:uuid;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

func (Team) TableName() string { return "teams" }

// ── team_members ──────────────────────────────────────────────────────

type TeamMember struct {
	TeamID   string    `gorm:"column:team_id;type:uuid;primaryKey"`
	UserID   string    `gorm:"column:user_id;type:uuid;primaryKey"`
	Role     string    `gorm:"column:role;not null;default:Member"`
	IsAdmin  bool      `gorm:"column:is_admin;not null;default:false"`
	JoinedAt time.Time `gorm:"column:joined_at;not null"`
}

func (TeamMember) TableName() string { return "team_members" }

// ── sprints ───────────────────────────────────────────────────────────

type Sprint struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	ScopeType string    `gorm:"column:scope_type;not null"` // "personal" | "team"
	ScopeID   string    `gorm:"column:scope_id;type:uuid;not null"`
	Name      string    `gorm:"column:name;not null"`
	StartDate time.Time `gorm:"column:start_date;not null;type:date"`
	EndDate   time.Time `gorm:"column:end_date;not null;type:date"`
	IsCurrent bool      `gorm:"column:is_current;not null;default:false"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
}

func (Sprint) TableName() string { return "sprints" }

// ── epics ─────────────────────────────────────────────────────────────

type Epic struct {
	ID             string    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	ScopeType      string    `gorm:"column:scope_type;not null"` // "personal" | "team"
	ScopeID        string    `gorm:"column:scope_id;type:uuid;not null"`
	Name           string    `gorm:"column:name;not null"`
	Color          string    `gorm:"column:color;not null"`
	OwnerID        *string   `gorm:"column:owner_id;type:uuid"`
	Status         string    `gorm:"column:status;not null;default:todo"` // todo | progress | done
	StartSprintID  *string   `gorm:"column:start_sprint_id;type:uuid"`
	TargetSprintID *string   `gorm:"column:target_sprint_id;type:uuid"`
	DescriptionMd  string    `gorm:"column:description_md;not null;default:''"`
	CreatedAt      time.Time `gorm:"column:created_at;not null"`
	UpdatedAt      time.Time `gorm:"column:updated_at;not null"`
}

func (Epic) TableName() string { return "epics" }

// ── tasks ─────────────────────────────────────────────────────────────

type Task struct {
	ID            string    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	ScopeType     string    `gorm:"column:scope_type;not null"` // "personal" | "team"
	ScopeID       string    `gorm:"column:scope_id;type:uuid;not null"`
	KeyNum        int       `gorm:"column:key_num;not null"` // per-scope counter -> NT-{key_num}
	Title         string    `gorm:"column:title;not null;default:''"`
	DescriptionMd string    `gorm:"column:description_md;not null;default:''"`
	Status        string    `gorm:"column:status;not null;default:backlog"` // backlog|progress|review|done
	Priority      string    `gorm:"column:priority;not null;default:med"`   // high|med|low
	AssigneeID    *string   `gorm:"column:assignee_id;type:uuid"`
	EpicID        *string   `gorm:"column:epic_id;type:uuid"`
	SprintID      *string   `gorm:"column:sprint_id;type:uuid"`
	CreatedBy     string    `gorm:"column:created_by;type:uuid;not null"`
	CreatedAt     time.Time `gorm:"column:created_at;not null"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null"`
}

func (Task) TableName() string { return "tasks" }

// ── task_counters ─────────────────────────────────────────────────────

// Per-scope NT-{key_num} allocation; `Next` is incremented via
// `UPDATE ... SET next = next + 1 ... RETURNING next` in the same
// transaction as the Task insert (avoids MAX()+1 races).
type TaskCounter struct {
	ScopeType string `gorm:"column:scope_type;primaryKey"` // "personal" | "team"
	ScopeID   string `gorm:"column:scope_id;type:uuid;primaryKey"`
	Next      int    `gorm:"column:next;not null;default:1"`
}

func (TaskCounter) TableName() string { return "task_counters" }

// ── task_activity ─────────────────────────────────────────────────────

// Append-only: written on every logged task PATCH and on POST /comments;
// never edited.
type TaskActivity struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	TaskID    string    `gorm:"column:task_id;type:uuid;not null"`
	ActorID   string    `gorm:"column:actor_id;type:uuid;not null"`
	Kind      string    `gorm:"column:kind;not null"` // created|status|priority|assignee|epic|sprint|comment
	Text      string    `gorm:"column:text;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
}

func (TaskActivity) TableName() string { return "task_activity" }
