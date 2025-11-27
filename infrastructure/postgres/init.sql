-- TestPilot AI Database Initialization Script
-- PostgreSQL 15+

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable pgcrypto for password hashing (backup)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================
-- USERS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'user')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index on username for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- ============================================
-- API SPECIFICATIONS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS api_specifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    version VARCHAR(50) NOT NULL,
    source_type VARCHAR(50) NOT NULL CHECK (source_type IN ('file', 'postman', 'git', 'url')),
    source_path TEXT,
    content_hash VARCHAR(64) NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL
);

-- Unique constraint on name + version
CREATE UNIQUE INDEX IF NOT EXISTS idx_api_spec_name_version ON api_specifications(name, version);

-- Index on source_type for filtering
CREATE INDEX IF NOT EXISTS idx_api_spec_source_type ON api_specifications(source_type);

-- Index on created_by for user filtering
CREATE INDEX IF NOT EXISTS idx_api_spec_created_by ON api_specifications(created_by);

-- ============================================
-- ENVIRONMENTS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS environments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    base_url TEXT NOT NULL,
    auth_config JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index on name for faster lookups
CREATE INDEX IF NOT EXISTS idx_environments_name ON environments(name);

-- ============================================
-- TEST EXECUTIONS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS test_executions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    api_spec_id UUID REFERENCES api_specifications(id) ON DELETE SET NULL,
    environment_id UUID REFERENCES environments(id) ON DELETE SET NULL,
    natural_language_request TEXT NOT NULL,
    constructed_request JSONB NOT NULL,
    response JSONB,
    validation_result JSONB,
    status VARCHAR(50) NOT NULL CHECK (status IN ('success', 'failed', 'error')),
    execution_time_ms INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_test_exec_user_id ON test_executions(user_id);
CREATE INDEX IF NOT EXISTS idx_test_exec_api_spec_id ON test_executions(api_spec_id);
CREATE INDEX IF NOT EXISTS idx_test_exec_status ON test_executions(status);
CREATE INDEX IF NOT EXISTS idx_test_exec_created_at ON test_executions(created_at DESC);

-- Full-text search index on natural language requests
CREATE INDEX IF NOT EXISTS idx_test_exec_nl_request ON test_executions USING gin(to_tsvector('english', natural_language_request));

-- ============================================
-- VALIDATION RULES TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS validation_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    api_spec_id UUID NOT NULL REFERENCES api_specifications(id) ON DELETE CASCADE,
    rule_type VARCHAR(50) NOT NULL CHECK (rule_type IN ('schema', 'status', 'custom')),
    rule_definition JSONB NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index on api_spec_id for filtering
CREATE INDEX IF NOT EXISTS idx_validation_rules_api_spec_id ON validation_rules(api_spec_id);

-- Index on rule_type
CREATE INDEX IF NOT EXISTS idx_validation_rules_type ON validation_rules(rule_type);

-- ============================================
-- LEARNED PATTERNS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS learned_patterns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    api_spec_id UUID NOT NULL REFERENCES api_specifications(id) ON DELETE CASCADE,
    pattern_data JSONB NOT NULL,
    success_count INTEGER DEFAULT 0,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index on api_spec_id for filtering
CREATE INDEX IF NOT EXISTS idx_learned_patterns_api_spec_id ON learned_patterns(api_spec_id);

-- Index on success_count for threshold queries
CREATE INDEX IF NOT EXISTS idx_learned_patterns_success_count ON learned_patterns(success_count);

-- ============================================
-- SYSTEM CONFIG TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS system_config (
    key VARCHAR(255) PRIMARY KEY,
    value JSONB NOT NULL,
    description TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- INGESTION LOGS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS ingestion_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_type VARCHAR(50) NOT NULL,
    source_path TEXT,
    status VARCHAR(50) NOT NULL CHECK (status IN ('success', 'failed', 'partial')),
    apis_ingested INTEGER DEFAULT 0,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index on created_at for recent logs
CREATE INDEX IF NOT EXISTS idx_ingestion_logs_created_at ON ingestion_logs(created_at DESC);

-- ============================================
-- TRIGGERS FOR UPDATED_AT
-- ============================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers for tables with updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_api_specifications_updated_at BEFORE UPDATE ON api_specifications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_environments_updated_at BEFORE UPDATE ON environments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_validation_rules_updated_at BEFORE UPDATE ON validation_rules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_learned_patterns_updated_at BEFORE UPDATE ON learned_patterns
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_system_config_updated_at BEFORE UPDATE ON system_config
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- DEFAULT DATA
-- ============================================

-- Insert default system configuration
INSERT INTO system_config (key, value, description) VALUES
    ('learning_threshold', '5', 'Number of successful tests before learning patterns'),
    ('history_retention_days', '90', 'Number of days to retain test execution history'),
    ('max_request_size_mb', '10', 'Maximum request size in megabytes'),
    ('default_timeout_seconds', '30', 'Default timeout for API calls in seconds')
ON CONFLICT (key) DO NOTHING;

-- Insert default admin user (password: admin123)
-- Note: This should be changed immediately in production
INSERT INTO users (username, password_hash, role) VALUES
    ('admin', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5lAWCb0pQw8Bq', 'admin')
ON CONFLICT (username) DO NOTHING;

-- Insert default QA environment
INSERT INTO environments (name, base_url, auth_config) VALUES
    ('QA1', 'https://qa1.example.com', '{"type": "api_key", "header": "X-API-Key"}')
ON CONFLICT (name) DO NOTHING;

-- ============================================
-- VIEWS FOR ANALYTICS
-- ============================================

-- View for test execution statistics
CREATE OR REPLACE VIEW test_execution_stats AS
SELECT
    user_id,
    api_spec_id,
    DATE(created_at) as execution_date,
    COUNT(*) as total_tests,
    COUNT(*) FILTER (WHERE status = 'success') as successful_tests,
    COUNT(*) FILTER (WHERE status = 'failed') as failed_tests,
    COUNT(*) FILTER (WHERE status = 'error') as error_tests,
    AVG(execution_time_ms) as avg_execution_time_ms,
    MAX(execution_time_ms) as max_execution_time_ms,
    MIN(execution_time_ms) as min_execution_time_ms
FROM test_executions
GROUP BY user_id, api_spec_id, DATE(created_at);

-- View for API usage statistics
CREATE OR REPLACE VIEW api_usage_stats AS
SELECT
    a.id as api_spec_id,
    a.name as api_name,
    a.version as api_version,
    COUNT(t.id) as total_executions,
    COUNT(t.id) FILTER (WHERE t.status = 'success') as successful_executions,
    AVG(t.execution_time_ms) as avg_execution_time_ms,
    MAX(t.created_at) as last_execution_at
FROM api_specifications a
LEFT JOIN test_executions t ON a.id = t.api_spec_id
GROUP BY a.id, a.name, a.version;

-- ============================================
-- GRANT PERMISSIONS
-- ============================================

-- Grant necessary permissions to the testpilot user
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO testpilot;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO testpilot;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO testpilot;

-- ============================================
-- COMPLETION MESSAGE
-- ============================================

DO $$
BEGIN
    RAISE NOTICE 'TestPilot AI database initialized successfully!';
    RAISE NOTICE 'Default admin user created: username=admin, password=admin123';
    RAISE NOTICE 'Please change the default admin password immediately!';
END $$;

