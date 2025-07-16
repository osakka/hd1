-- Cross-platform Client System tables for Phase 4: Universal Platform

-- Clients
CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    platform VARCHAR(50) NOT NULL,
    version VARCHAR(50) NOT NULL,
    capabilities JSONB NOT NULL DEFAULT '[]',
    configuration JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    last_seen TIMESTAMP WITH TIME ZONE NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_clients_session_id ON clients(session_id);
CREATE INDEX IF NOT EXISTS idx_clients_user_id ON clients(user_id);
CREATE INDEX IF NOT EXISTS idx_clients_type ON clients(type);
CREATE INDEX IF NOT EXISTS idx_clients_platform ON clients(platform);
CREATE INDEX IF NOT EXISTS idx_clients_status ON clients(status);
CREATE INDEX IF NOT EXISTS idx_clients_last_seen ON clients(last_seen);
CREATE INDEX IF NOT EXISTS idx_clients_created_at ON clients(created_at);

-- Client Messages
CREATE TABLE IF NOT EXISTS client_messages (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    data JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_client_messages_session_id ON client_messages(session_id);
CREATE INDEX IF NOT EXISTS idx_client_messages_client_id ON client_messages(client_id);
CREATE INDEX IF NOT EXISTS idx_client_messages_type ON client_messages(type);
CREATE INDEX IF NOT EXISTS idx_client_messages_action ON client_messages(action);
CREATE INDEX IF NOT EXISTS idx_client_messages_created_at ON client_messages(created_at);

-- Client Capabilities
CREATE TABLE IF NOT EXISTS client_capabilities (
    id UUID PRIMARY KEY,
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    capability VARCHAR(100) NOT NULL,
    version VARCHAR(50),
    parameters JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    UNIQUE(client_id, capability)
);

CREATE INDEX IF NOT EXISTS idx_client_capabilities_client_id ON client_capabilities(client_id);
CREATE INDEX IF NOT EXISTS idx_client_capabilities_capability ON client_capabilities(capability);

-- Platform Adapters
CREATE TABLE IF NOT EXISTS platform_adapters (
    id UUID PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    platform VARCHAR(50) NOT NULL,
    version VARCHAR(50) NOT NULL,
    capabilities JSONB NOT NULL DEFAULT '[]',
    requirements JSONB NOT NULL DEFAULT '[]',
    configuration JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_platform_adapters_name ON platform_adapters(name);
CREATE INDEX IF NOT EXISTS idx_platform_adapters_platform ON platform_adapters(platform);
CREATE INDEX IF NOT EXISTS idx_platform_adapters_status ON platform_adapters(status);

-- Client Sessions (for tracking client activity)
CREATE TABLE IF NOT EXISTS client_sessions (
    id UUID PRIMARY KEY,
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ended_at TIMESTAMP WITH TIME ZONE,
    duration_seconds INTEGER,
    actions_count INTEGER NOT NULL DEFAULT 0,
    bytes_sent BIGINT NOT NULL DEFAULT 0,
    bytes_received BIGINT NOT NULL DEFAULT 0,
    metadata JSONB NOT NULL DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS idx_client_sessions_client_id ON client_sessions(client_id);
CREATE INDEX IF NOT EXISTS idx_client_sessions_session_id ON client_sessions(session_id);
CREATE INDEX IF NOT EXISTS idx_client_sessions_started_at ON client_sessions(started_at);
CREATE INDEX IF NOT EXISTS idx_client_sessions_ended_at ON client_sessions(ended_at);

-- Client Performance Metrics
CREATE TABLE IF NOT EXISTS client_performance_metrics (
    id UUID PRIMARY KEY,
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    metric_type VARCHAR(50) NOT NULL,
    value FLOAT NOT NULL,
    unit VARCHAR(20) NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_client_performance_metrics_client_id ON client_performance_metrics(client_id);
CREATE INDEX IF NOT EXISTS idx_client_performance_metrics_metric_type ON client_performance_metrics(metric_type);
CREATE INDEX IF NOT EXISTS idx_client_performance_metrics_created_at ON client_performance_metrics(created_at);

-- Client Synchronization States
CREATE TABLE IF NOT EXISTS client_sync_states (
    id UUID PRIMARY KEY,
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    last_sync_version BIGINT NOT NULL DEFAULT 0,
    last_sync_timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    sync_status VARCHAR(20) NOT NULL DEFAULT 'synced',
    pending_operations JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    UNIQUE(client_id, entity_type, entity_id)
);

CREATE INDEX IF NOT EXISTS idx_client_sync_states_client_id ON client_sync_states(client_id);
CREATE INDEX IF NOT EXISTS idx_client_sync_states_entity_type ON client_sync_states(entity_type);
CREATE INDEX IF NOT EXISTS idx_client_sync_states_entity_id ON client_sync_states(entity_id);
CREATE INDEX IF NOT EXISTS idx_client_sync_states_sync_status ON client_sync_states(sync_status);
CREATE INDEX IF NOT EXISTS idx_client_sync_states_last_sync_timestamp ON client_sync_states(last_sync_timestamp);

-- Triggers for updating timestamps
CREATE TRIGGER update_clients_updated_at
    BEFORE UPDATE ON clients
    FOR EACH ROW
    EXECUTE FUNCTION update_rtc_updated_at_column();

CREATE TRIGGER update_platform_adapters_updated_at
    BEFORE UPDATE ON platform_adapters
    FOR EACH ROW
    EXECUTE FUNCTION update_rtc_updated_at_column();

CREATE TRIGGER update_client_sync_states_updated_at
    BEFORE UPDATE ON client_sync_states
    FOR EACH ROW
    EXECUTE FUNCTION update_rtc_updated_at_column();