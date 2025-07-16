-- HD1 Universal Platform Database Schema
-- Phase 1: Foundation
-- This schema supports multi-tenant sessions and service registry

-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    last_login TIMESTAMP,
    status VARCHAR(20) DEFAULT 'active',
    profile JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}'
);

-- Sessions table
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    owner_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'active',
    visibility VARCHAR(20) DEFAULT 'private',
    max_participants INTEGER DEFAULT 10,
    current_participants INTEGER DEFAULT 0,
    settings JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}'
);

-- Participants table
CREATE TABLE participants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    role VARCHAR(20) DEFAULT 'participant',
    joined_at TIMESTAMP DEFAULT NOW(),
    left_at TIMESTAMP,
    last_seen TIMESTAMP DEFAULT NOW(),
    permissions JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    UNIQUE(session_id, user_id)
);

-- Services table
CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    endpoint VARCHAR(500) NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    capabilities JSONB DEFAULT '[]',
    ui_mapping JSONB DEFAULT '{}',
    permissions JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    metadata JSONB DEFAULT '{}',
    UNIQUE(name)
);

-- Service health table
CREATE TABLE service_health (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL,
    response_time INTEGER,
    error_message TEXT,
    checked_at TIMESTAMP DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);

-- Permissions table
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    permission_type VARCHAR(50) NOT NULL,
    granted_at TIMESTAMP DEFAULT NOW(),
    granted_by UUID REFERENCES users(id),
    metadata JSONB DEFAULT '{}'
);

-- OAuth tokens table
CREATE TABLE oauth_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    access_token VARCHAR(1000) NOT NULL,
    refresh_token VARCHAR(1000),
    token_type VARCHAR(20) DEFAULT 'Bearer',
    expires_at TIMESTAMP,
    scope VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);

-- Audit log table
CREATE TABLE audit_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    session_id UUID REFERENCES sessions(id),
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);

-- Indexes for performance
CREATE INDEX idx_sessions_owner_id ON sessions(owner_id);
CREATE INDEX idx_sessions_status ON sessions(status);
CREATE INDEX idx_sessions_created_at ON sessions(created_at);

CREATE INDEX idx_participants_session_id ON participants(session_id);
CREATE INDEX idx_participants_user_id ON participants(user_id);
CREATE INDEX idx_participants_joined_at ON participants(joined_at);

CREATE INDEX idx_services_name ON services(name);
CREATE INDEX idx_services_type ON services(type);
CREATE INDEX idx_services_status ON services(status);

CREATE INDEX idx_service_health_service_id ON service_health(service_id);
CREATE INDEX idx_service_health_checked_at ON service_health(checked_at);

CREATE INDEX idx_permissions_resource_type ON permissions(resource_type);
CREATE INDEX idx_permissions_resource_id ON permissions(resource_id);
CREATE INDEX idx_permissions_user_id ON permissions(user_id);

CREATE INDEX idx_oauth_tokens_user_id ON oauth_tokens(user_id);
CREATE INDEX idx_oauth_tokens_access_token ON oauth_tokens(access_token);

CREATE INDEX idx_audit_log_user_id ON audit_log(user_id);
CREATE INDEX idx_audit_log_session_id ON audit_log(session_id);
CREATE INDEX idx_audit_log_created_at ON audit_log(created_at);

-- Triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sessions_updated_at BEFORE UPDATE ON sessions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_services_updated_at BEFORE UPDATE ON services
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to update participant count
CREATE OR REPLACE FUNCTION update_session_participant_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE sessions 
        SET current_participants = (
            SELECT COUNT(*) 
            FROM participants 
            WHERE session_id = NEW.session_id AND left_at IS NULL
        )
        WHERE id = NEW.session_id;
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        UPDATE sessions 
        SET current_participants = (
            SELECT COUNT(*) 
            FROM participants 
            WHERE session_id = NEW.session_id AND left_at IS NULL
        )
        WHERE id = NEW.session_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE sessions 
        SET current_participants = (
            SELECT COUNT(*) 
            FROM participants 
            WHERE session_id = OLD.session_id AND left_at IS NULL
        )
        WHERE id = OLD.session_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_session_participant_count_trigger
    AFTER INSERT OR UPDATE OR DELETE ON participants
    FOR EACH ROW EXECUTE FUNCTION update_session_participant_count();

-- Sample data for development
INSERT INTO users (email, username, password_hash) VALUES 
('admin@hd1.universe', 'admin', '$2a$10$example.hash.for.development');

INSERT INTO sessions (name, description, owner_id) VALUES 
('Default Session', 'Default session for testing', 
 (SELECT id FROM users WHERE username = 'admin'));

INSERT INTO services (name, description, type, endpoint, capabilities) VALUES 
('HD1 Core', 'Core HD1 3D rendering service', 'core', 'http://localhost:8080/api', 
 '["3d_rendering", "websocket_sync", "entity_management"]');