-- database schema for sqlc code generation for archival PostgreSQL

CREATE TABLE IF NOT EXISTS sessions_archive (
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

CREATE TABLE IF NOT EXISTS transactions_archive (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    operation VARCHAR(20) NOT NULL CHECK (operation IN ('getfacl', 'setfacl')),
    target_path TEXT NOT NULL,
    entries JSONB NOT NULL DEFAULT '[]'::jsonb,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'success', 'failed')),
    error_msg TEXT,
    output TEXT,
    executed_by VARCHAR(255) NOT NULL,
    duration_ms BIGINT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- * Indexing for Sessions * --

-- Primary lookup indexes
CREATE INDEX IF NOT EXISTS idx_sessions_username ON sessions_archive(username);
CREATE INDEX IF NOT EXISTS idx_sessions_status ON sessions_archive(status);
CREATE INDEX IF NOT EXISTS idx_sessions_ip ON sessions_archive(ip);

-- Time-based indexes for chronological queries
CREATE INDEX IF NOT EXISTS idx_sessions_created_at ON sessions_archive(created_at);
CREATE INDEX IF NOT EXISTS idx_sessions_last_active_at ON sessions_archive(last_active_at);
CREATE INDEX IF NOT EXISTS idx_sessions_expiry ON sessions_archive(expiry);
CREATE INDEX IF NOT EXISTS idx_sessions_archived_at ON sessions_archive(archived_at);

-- Composite indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_sessions_username_status ON sessions_archive(username, status);
CREATE INDEX IF NOT EXISTS idx_sessions_username_created_at ON sessions_archive(username, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_sessions_username_last_active ON sessions_archive(username, last_active_at DESC);
CREATE INDEX IF NOT EXISTS idx_sessions_status_created_at ON sessions_archive(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_sessions_status_archived_at ON sessions_archive(status, archived_at DESC);

-- Performance indexes for analytics and monitoring
CREATE INDEX IF NOT EXISTS idx_sessions_completed_count ON sessions_archive(completed_count) WHERE completed_count > 0;
CREATE INDEX IF NOT EXISTS idx_sessions_failed_count ON sessions_archive(failed_count) WHERE failed_count > 0;
CREATE INDEX IF NOT EXISTS idx_sessions_user_agent ON sessions_archive(user_agent) WHERE user_agent IS NOT NULL;

-- Specialized composite indexes for complex queries
CREATE INDEX IF NOT EXISTS idx_sessions_username_ip ON sessions_archive(username, ip);
CREATE INDEX IF NOT EXISTS idx_sessions_ip_created_at ON sessions_archive(ip, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_sessions_expiry_status ON sessions_archive(expiry, status);

-- Partial indexes for active sessions monitoring
CREATE INDEX IF NOT EXISTS idx_sessions_active_last_active ON sessions_archive(last_active_at DESC) 
    WHERE status = 'active';
CREATE INDEX IF NOT EXISTS idx_sessions_expired_recent ON sessions_archive(expiry DESC) 
    WHERE status = 'expired';

-- Performance indexes for user activity analysis
CREATE INDEX IF NOT EXISTS idx_sessions_high_activity ON sessions_archive(username, completed_count DESC) 
    WHERE completed_count > 10;
CREATE INDEX IF NOT EXISTS idx_sessions_problematic ON sessions_archive(username, failed_count DESC) 
    WHERE failed_count > 5;

-- Indexes for cleanup and maintenance operations
CREATE INDEX IF NOT EXISTS idx_sessions_old_archived ON sessions_archive(archived_at) 
    WHERE archived_at < NOW() - INTERVAL '90 days';
CREATE INDEX IF NOT EXISTS idx_sessions_old_expired ON sessions_archive(expiry) 
    WHERE status = 'expired' AND expiry < NOW() - INTERVAL '30 days';

-- Security and audit indexes
CREATE INDEX IF NOT EXISTS idx_sessions_ip_count ON sessions_archive(ip, username, created_at) 
    WHERE ip IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_sessions_concurrent_users ON sessions_archive(username, created_at, expiry) 
    WHERE status = 'active';

-- * Indexing for Transactions * --

-- Primary lookup indexes
CREATE INDEX IF NOT EXISTS idx_transactions_session_id ON transactions_archive(session_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions_archive(status);
CREATE INDEX IF NOT EXISTS idx_transactions_operation ON transactions_archive(operation);

-- Time-based indexes for chronological queries
CREATE INDEX IF NOT EXISTS idx_transactions_timestamp ON transactions_archive(timestamp);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions_archive(created_at);

-- Composite indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_transactions_session_status ON transactions_archive(session_id, status);
CREATE INDEX IF NOT EXISTS idx_transactions_session_operation ON transactions_archive(session_id, operation);
CREATE INDEX IF NOT EXISTS idx_transactions_session_timestamp ON transactions_archive(session_id, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_session_created_at ON transactions_archive(session_id, created_at DESC);

-- Performance indexes for filtering and analytics
CREATE INDEX IF NOT EXISTS idx_transactions_target_path ON transactions_archive(target_path);
CREATE INDEX IF NOT EXISTS idx_transactions_executed_by ON transactions_archive(executed_by);
CREATE INDEX IF NOT EXISTS idx_transactions_duration ON transactions_archive(duration_ms) WHERE duration_ms IS NOT NULL;

-- Specialized composite indexes for complex queries
CREATE INDEX IF NOT EXISTS idx_transactions_session_path ON transactions_archive(session_id, target_path);
CREATE INDEX IF NOT EXISTS idx_transactions_status_timestamp ON transactions_archive(status, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_operation_timestamp ON transactions_archive(operation, timestamp DESC);

-- JSONB indexes for ACL entries queries (if you need to query within entries)
CREATE INDEX IF NOT EXISTS idx_transactions_entries_gin ON transactions_archive USING GIN (entries);

-- Partial indexes for active/recent data (performance optimization)
CREATE INDEX IF NOT EXISTS idx_transactions_recent_pending ON transactions_archive(session_id, timestamp DESC) 
    WHERE status = 'pending';
CREATE INDEX IF NOT EXISTS idx_transactions_recent_failed ON transactions_archive(session_id, timestamp DESC) 
    WHERE status = 'failed';

-- Index for cleanup operations (if you periodically clean old records)
CREATE INDEX IF NOT EXISTS idx_transactions_cleanup ON transactions_archive(created_at) 
    WHERE created_at < NOW() - INTERVAL '30 days';

CREATE INDEX idx_sessions_archive_time ON sessions_archive(archived_at DESC);
