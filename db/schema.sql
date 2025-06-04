-- database schema for sqlc code generation for archival PostgreSQL

CREATE TABLE sessions_archive (
    id              UUID PRIMARY KEY,
    username        TEXT NOT NULL,
    ip              TEXT,
    user_agent      TEXT,
    status          TEXT CHECK (status IN ('active', 'expired')) NOT NULL,
    created_at      TIMESTAMP NOT NULL,
    last_active_at  TIMESTAMP NOT NULL,
    expiry          TIMESTAMP NOT NULL,
    completed_count INTEGER DEFAULT 0,
    failed_count    INTEGER DEFAULT 0,
    archived_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions_archive (
    id           UUID PRIMARY KEY,
    session_id   UUID REFERENCES sessions_archive(id) ON DELETE CASCADE,
    status       TEXT CHECK (status IN ('success', 'failure')) NOT NULL,
    output       TEXT,
    created_at   TIMESTAMP NOT NULL
);

-- indexes
CREATE INDEX idx_transactions_session_id ON transactions_archive(session_id);
CREATE INDEX idx_sessions_archive_time ON sessions_archive(archived_at DESC);
