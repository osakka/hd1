-- Operational Transform tables for Phase 2: Collaboration

-- OT Documents
CREATE TABLE IF NOT EXISTS ot_documents (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    content TEXT NOT NULL DEFAULT '',
    version BIGINT NOT NULL DEFAULT 1,
    type VARCHAR(50) NOT NULL DEFAULT 'text',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS idx_ot_documents_session_id ON ot_documents(session_id);
CREATE INDEX IF NOT EXISTS idx_ot_documents_type ON ot_documents(type);
CREATE INDEX IF NOT EXISTS idx_ot_documents_created_at ON ot_documents(created_at);

-- OT Clients
CREATE TABLE IF NOT EXISTS ot_clients (
    id UUID PRIMARY KEY,
    document_id UUID NOT NULL REFERENCES ot_documents(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    version BIGINT NOT NULL DEFAULT 1,
    state VARCHAR(50) NOT NULL DEFAULT 'synchronized',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS idx_ot_clients_document_id ON ot_clients(document_id);
CREATE INDEX IF NOT EXISTS idx_ot_clients_user_id ON ot_clients(user_id);
CREATE INDEX IF NOT EXISTS idx_ot_clients_state ON ot_clients(state);

-- OT Operations
CREATE TABLE IF NOT EXISTS ot_operations (
    id UUID PRIMARY KEY,
    document_id UUID NOT NULL REFERENCES ot_documents(id) ON DELETE CASCADE,
    client_id UUID NOT NULL REFERENCES ot_clients(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    position INTEGER NOT NULL,
    content TEXT,
    length INTEGER,
    version BIGINT NOT NULL,
    data JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_ot_operations_document_id ON ot_operations(document_id);
CREATE INDEX IF NOT EXISTS idx_ot_operations_client_id ON ot_operations(client_id);
CREATE INDEX IF NOT EXISTS idx_ot_operations_user_id ON ot_operations(user_id);
CREATE INDEX IF NOT EXISTS idx_ot_operations_type ON ot_operations(type);
CREATE INDEX IF NOT EXISTS idx_ot_operations_version ON ot_operations(version);
CREATE INDEX IF NOT EXISTS idx_ot_operations_created_at ON ot_operations(created_at);

-- OT Snapshots (for performance optimization)
CREATE TABLE IF NOT EXISTS ot_snapshots (
    id UUID PRIMARY KEY,
    document_id UUID NOT NULL REFERENCES ot_documents(id) ON DELETE CASCADE,
    version BIGINT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS idx_ot_snapshots_document_id ON ot_snapshots(document_id);
CREATE INDEX IF NOT EXISTS idx_ot_snapshots_version ON ot_snapshots(version);
CREATE INDEX IF NOT EXISTS idx_ot_snapshots_created_at ON ot_snapshots(created_at);

-- OT Conflicts (for conflict resolution tracking)
CREATE TABLE IF NOT EXISTS ot_conflicts (
    id UUID PRIMARY KEY,
    document_id UUID NOT NULL REFERENCES ot_documents(id) ON DELETE CASCADE,
    operation_id UUID NOT NULL REFERENCES ot_operations(id) ON DELETE CASCADE,
    conflict_type VARCHAR(50) NOT NULL,
    resolution VARCHAR(50) NOT NULL,
    conflict_data JSONB NOT NULL DEFAULT '{}',
    resolution_data JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    resolved_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_ot_conflicts_document_id ON ot_conflicts(document_id);
CREATE INDEX IF NOT EXISTS idx_ot_conflicts_operation_id ON ot_conflicts(operation_id);
CREATE INDEX IF NOT EXISTS idx_ot_conflicts_conflict_type ON ot_conflicts(conflict_type);
CREATE INDEX IF NOT EXISTS idx_ot_conflicts_created_at ON ot_conflicts(created_at);

-- Triggers for updating timestamps
CREATE TRIGGER update_ot_documents_updated_at
    BEFORE UPDATE ON ot_documents
    FOR EACH ROW
    EXECUTE FUNCTION update_rtc_updated_at_column();

CREATE TRIGGER update_ot_clients_updated_at
    BEFORE UPDATE ON ot_clients
    FOR EACH ROW
    EXECUTE FUNCTION update_rtc_updated_at_column();