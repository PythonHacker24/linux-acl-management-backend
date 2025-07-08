-- database schema for sqlc code generation for archival PostgreSQL

CREATE TABLE IF NOT EXISTS sessions_archive (
    id              UUID PRIMARY KEY,
    username        TEXT NOT NULL,
    ip              TEXT,
    user_agent      TEXT,
    status          TEXT CHECK (status IN ('active', 'expired', 'pending')) NOT NULL,
    created_at      TIMESTAMP NOT NULL,
    last_active_at  TIMESTAMP NOT NULL,
    expiry          TIMESTAMP NOT NULL,
    completed_count INTEGER DEFAULT 0,
    failed_count    INTEGER DEFAULT 0,
    archived_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pending_transactions_archive (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    operation VARCHAR(20) NOT NULL CHECK (operation IN ('getfacl', 'setfacl')),
    target_path TEXT NOT NULL,
    entries JSONB NOT NULL DEFAULT '[]'::jsonb,
    status TEXT CHECK (status IN ('pending')) NOT NULL,
    error_msg TEXT,
    output TEXT,
    executed_by VARCHAR(255) NOT NULL,
    duration_ms BIGINT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS results_transactions_archive (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    operation VARCHAR(20) NOT NULL CHECK (operation IN ('getfacl', 'setfacl')),
    target_path TEXT NOT NULL,
    entries JSONB NOT NULL DEFAULT '[]'::jsonb,
    status TEXT CHECK (status IN ('success', 'failed')) NOT NULL,
    error_msg TEXT,
    output TEXT,
    executed_by VARCHAR(255) NOT NULL,
    duration_ms BIGINT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

/* add indexing for optimization */
